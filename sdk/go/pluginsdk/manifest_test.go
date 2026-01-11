// Copyright 2026 PulumiCost/FinFocus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pluginsdk_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//nolint:gocognit // Test complexity is acceptable for comprehensive table-driven test
func TestManifestSaveLoad(t *testing.T) {
	now := timestamppb.New(time.Date(2025, time.January, 1, 12, 0, 0, 0, time.UTC))

	expectedManifest := &pbc.PluginManifest{
		Metadata: &pbc.PluginMetadata{
			Name:        "test-plugin",
			Version:     "1.0.0",
			Description: "A test plugin for PulumiCost",
			Author:      "Test Author",
			Homepage:    "https://example.com",
			Repository:  "https://github.com/example/test-plugin",
			License:     "Apache-2.0",
			Keywords:    []string{"cost", "cloud", "test"},
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		Specification: &pbc.PluginSpecification{
			SpecVersion:        "0.1.0",
			SupportedProviders: []string{"aws", "azure"},
			SupportedResources: map[string]*pbc.ProviderResources{
				"aws": {
					ResourceTypes: []string{"ec2", "s3"},
					BillingModes:  []string{"hourly", "monthly"},
					Regions:       []string{"us-east-1", "us-west-2"},
				},
			},
			Capabilities: []string{"projected_cost", "actual_cost"},
			ServiceDefinition: &pbc.ServiceDefinition{
				ServiceName:     "CostSourceService",
				PackageName:     "pulumicost.v1",
				Methods:         []string{"GetProjectedCost", "GetActualCost"},
				Port:            50051,
				HealthCheckPath: "/healthz",
			},
			ObservabilitySupport: &pbc.ObservabilitySupport{
				MetricsEnabled:      true,
				TracingEnabled:      true,
				LoggingEnabled:      true,
				HealthChecksEnabled: true,
				SliSupport:          false,
			},
		},
		Security: &pbc.PluginSecurity{
			Signature:        "some-signature",
			PublicKey:        "some-public-key",
			CertificateChain: []string{"cert1", "cert2"},
			SecurityLevel:    pbc.SecurityLevel_SECURITY_LEVEL_VERIFIED,
			Permissions:      []string{"read_cloud_creds"},
			SandboxRequired:  false,
		},
		Installation: &pbc.InstallationSpec{
			InstallationMethod: pbc.InstallationMethod_INSTALLATION_METHOD_BINARY,
			DownloadUrl:        "https://example.com/download",
			Checksum:           "abc123def456",
			ChecksumAlgorithm:  "sha256",
			InstallScript:      "install.sh",
			PreInstallChecks:   []string{"check_docker"},
			PostInstallSteps:   []string{"restart_service"},
		},
		Configuration: &pbc.ConfigurationSpec{
			Schema:         "{}",
			DefaultConfig:  "{}",
			RequiredFields: []string{"api_key"},
			Examples: []*pbc.ConfigurationExample{
				{
					Name:        "example-config",
					Description: "A basic configuration",
					Config:      `{\"api_key\": \"123\"}`,
				},
			},
		},
	}

	testCases := []struct {
		name         string
		ext          string
		marshalErr   string // Substring to expect in marshal error
		unmarshalErr string // Substring to expect in unmarshal error
	}{
		{
			name: "YAML format",
			ext:  ".yaml",
		},
		{
			name: "JSON format",
			ext:  ".json",
		},
		{
			name:         "Unsupported format",
			ext:          ".txt",
			marshalErr:   "unsupported manifest file extension",
			unmarshalErr: "unsupported manifest file extension",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			manifestPath := filepath.Join(tmpDir, "test-manifest"+tc.ext)

			// Test SaveManifest
			err := pluginsdk.SaveManifest(manifestPath, expectedManifest)
			if tc.marshalErr != "" {
				if err == nil || !ErrorContains(err, tc.marshalErr) {
					t.Errorf("SaveManifest() expected error containing %q, got %v", tc.marshalErr, err)
				}
				// If we expected a marshal error, we can't proceed to load
				return
			}
			if err != nil {
				t.Fatalf("SaveManifest() failed: %v", err)
			}

			// Test LoadManifest
			loadedManifest, err := pluginsdk.LoadManifest(manifestPath)
			if tc.unmarshalErr != "" {
				if err == nil || !ErrorContains(err, tc.unmarshalErr) {
					t.Errorf("LoadManifest() expected error containing %q, got %v", tc.unmarshalErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("LoadManifest() failed: %v", err)
			}

			// Compare loaded manifest with expected manifest
			if diff := cmp.Diff(expectedManifest, loadedManifest, protocmp.Transform()); diff != "" {
				t.Errorf("Loaded manifest mismatch (-want +got):
%s", diff)
			}
		})
	}
}

func TestLoadManifestErrors(t *testing.T) {
	testCases := []struct {
		name        string
		filePath    string
		fileContent string
		expectError string // Substring to expect in error message
	}{
		{
			name:        "non-existent file",
			filePath:    "non-existent.yaml",
			expectError: "reading manifest file",
		},
		{
			name:        "invalid YAML content",
			filePath:    "invalid.yaml",
			fileContent: "invalid: yaml: content: [",
			expectError: "parsing YAML manifest",
		},
		{
			name:        "invalid JSON content",
			filePath:    "invalid.json",
			fileContent: "{ \"metadata\": { \"name\": 123 } }", // Name should be string, not int
			expectError: "parsing JSON manifest",
		},
		{
			name:        "unsupported file extension",
			filePath:    "test.txt",
			fileContent: "some content",
			expectError: "unsupported manifest file extension",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			fullPath := filepath.Join(tmpDir, tc.filePath)

			if tc.fileContent != "" {
				err := os.WriteFile(fullPath, []byte(tc.fileContent), 0o600)
				if err != nil {
					t.Fatalf("Failed to write file for test %q: %v", tc.name, err)
				}
			}

			_, err := pluginsdk.LoadManifest(fullPath)
			if err == nil || !ErrorContains(err, tc.expectError) {
				t.Errorf("LoadManifest() expected error containing %q, got %v", tc.expectError, err)
			}
		})
	}
}

// ErrorContains checks if an error's message contains a specific substring.
// This is a helper function to avoid direct string comparison on error messages.
func ErrorContains(err error, s string) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), s)
}

func TestValidationErrors(t *testing.T) {
	// Test empty errors
	var emptyErrs pluginsdk.ValidationErrors
	if emptyErrs.Error() != "no validation errors" {
		t.Errorf("Expected 'no validation errors', got %q", emptyErrs.Error())
	}

	// Test with single error
	singleErr := pluginsdk.ValidationErrors{
		&pbc.ValidationError{Message: "first error"},
	}
	errMsg := singleErr.Error()
	if !strings.Contains(errMsg, "validation failed with 1 error(s)") {
		t.Errorf("Expected '1 error(s)' in message, got %q", errMsg)
	}
	if !strings.Contains(errMsg, "first error") {
		t.Errorf("Expected 'first error' in message, got %q", errMsg)
	}

	// Test with multiple errors
	multiErr := pluginsdk.ValidationErrors{
		&pbc.ValidationError{Message: "error one"},
		&pbc.ValidationError{Message: "error two"},
	}
	errMsg = multiErr.Error()
	if !strings.Contains(errMsg, "validation failed with 2 error(s)") {
		t.Errorf("Expected '2 error(s)' in message, got %q", errMsg)
	}
	if !strings.Contains(errMsg, "error one") || !strings.Contains(errMsg, "error two") {
		t.Errorf("Expected both errors in message, got %q", errMsg)
	}
}
