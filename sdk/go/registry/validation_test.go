package registry

import (
	"testing"
	"time"
)

func TestValidateManifest(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name          string
		manifest      *PluginManifest
		expectValid   bool
		expectErrors  int
		expectWarnings int
	}{
		{
			name: "valid manifest",
			manifest: &PluginManifest{
				Name:        "test-plugin",
				DisplayName: "Test Plugin",
				Description: "A test plugin for validation",
				Version:     "1.0.0",
				APIVersion:  "v1",
				PluginType:  PluginTypeCostSource,
				Capabilities: []Capability{
					CapabilityActualCost,
					CapabilityProjectedCost,
				},
				SupportedProviders: []string{"aws", "gcp"},
				Requirements: PluginRequirements{
					MinAPIVersion: "v1",
					MaxAPIVersion: "v1",
				},
				Authentication: AuthenticationConfig{
					Required: false,
					Methods:  []AuthenticationMethod{AuthMethodAPIKey},
				},
				Installation: InstallationInfo{
					BinaryURL: "https://example.com/plugin",
					Checksum:  "sha256:a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3",
					SizeBytes: 1024,
				},
				Contacts: ContactInfo{
					Maintainers: []Maintainer{
						{
							Name:  "Test Maintainer",
							Email: "test@example.com",
							URL:   "https://github.com/test",
						},
					},
					SupportURL:       "https://github.com/test/issues",
					DocumentationURL: "https://docs.example.com",
				},
				Metadata: PluginMetadata{
					Tags:     []string{"test", "example"},
					License:  "Apache-2.0",
					Created:  time.Now(),
					Updated:  time.Now(),
				},
			},
			expectValid:   true,
			expectErrors:  0,
			expectWarnings: 0,
		},
		{
			name: "missing required fields",
			manifest: &PluginManifest{
				// Missing all required fields
			},
			expectValid:   false,
			expectErrors:  5, // name, version, api_version, plugin_type, capabilities
			expectWarnings: 0,
		},
		{
			name: "invalid name format",
			manifest: &PluginManifest{
				Name:        "Invalid_Name",
				Version:     "1.0.0",
				APIVersion:  "v1",
				PluginType:  PluginTypeCostSource,
				Capabilities: []Capability{CapabilityActualCost},
			},
			expectValid:  false,
			expectErrors: 1,
		},
		{
			name: "invalid version format",
			manifest: &PluginManifest{
				Name:        "test-plugin",
				Version:     "invalid-version",
				APIVersion:  "v1",
				PluginType:  PluginTypeCostSource,
				Capabilities: []Capability{CapabilityActualCost},
			},
			expectValid:  false,
			expectErrors: 1,
		},
		{
			name: "invalid api version",
			manifest: &PluginManifest{
				Name:        "test-plugin",
				Version:     "1.0.0",
				APIVersion:  "invalid",
				PluginType:  PluginTypeCostSource,
				Capabilities: []Capability{CapabilityActualCost},
			},
			expectValid:  false,
			expectErrors: 1,
		},
		{
			name: "invalid plugin type",
			manifest: &PluginManifest{
				Name:        "test-plugin",
				Version:     "1.0.0",
				APIVersion:  "v1",
				PluginType:  PluginType("invalid"),
				Capabilities: []Capability{CapabilityActualCost},
			},
			expectValid:  false,
			expectErrors: 1,
		},
		{
			name: "invalid capability",
			manifest: &PluginManifest{
				Name:        "test-plugin",
				Version:     "1.0.0",
				APIVersion:  "v1",
				PluginType:  PluginTypeCostSource,
				Capabilities: []Capability{Capability("invalid")},
			},
			expectValid:  false,
			expectErrors: 1,
		},
		{
			name: "invalid provider",
			manifest: &PluginManifest{
				Name:               "test-plugin",
				Version:            "1.0.0",
				APIVersion:         "v1",
				PluginType:         PluginTypeCostSource,
				Capabilities:       []Capability{CapabilityActualCost},
				SupportedProviders: []string{"invalid"},
			},
			expectValid:  false,
			expectErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validator.ValidateManifest(tt.manifest)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Valid != tt.expectValid {
				t.Errorf("expected valid=%v, got valid=%v", tt.expectValid, result.Valid)
			}

			if len(result.Errors) != tt.expectErrors {
				t.Errorf("expected %d errors, got %d errors: %+v", tt.expectErrors, len(result.Errors), result.Errors)
			}

			if tt.expectWarnings > 0 && len(result.Warnings) < tt.expectWarnings {
				t.Errorf("expected at least %d warnings, got %d warnings", tt.expectWarnings, len(result.Warnings))
			}
		})
	}
}

func TestValidateBinary(t *testing.T) {
	validator := NewValidator()

	manifest := &PluginManifest{
		Installation: InstallationInfo{
			Checksum:  "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			SizeBytes: 0,
		},
	}

	tests := []struct {
		name         string
		binaryData   []byte
		manifest     *PluginManifest
		expectValid  bool
		expectErrors int
	}{
		{
			name:         "empty binary",
			binaryData:   []byte{},
			manifest:     manifest,
			expectValid:  false,
			expectErrors: 1,
		},
		{
			name:         "valid checksum",
			binaryData:   []byte{}, // empty data has known SHA-256
			manifest:     manifest,
			expectValid:  false, // Still invalid due to empty binary
			expectErrors: 1,
		},
		{
			name:       "checksum mismatch",
			binaryData: []byte("test data"),
			manifest: &PluginManifest{
				Installation: InstallationInfo{
					Checksum:  "sha256:wrong_checksum",
					SizeBytes: 9,
				},
			},
			expectValid:  false,
			expectErrors: 2, // Invalid checksum format + checksum mismatch would be caught by format validation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validator.ValidateBinary(tt.binaryData, tt.manifest)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Valid != tt.expectValid {
				t.Errorf("expected valid=%v, got valid=%v", tt.expectValid, result.Valid)
			}

			if len(result.Errors) < 1 && !tt.expectValid {
				t.Errorf("expected at least 1 error for invalid binary, got %d errors", len(result.Errors))
			}
		})
	}
}

func TestValidateCompatibility(t *testing.T) {
	validator := NewValidator()

	manifest := &PluginManifest{
		Requirements: PluginRequirements{
			MinAPIVersion: "v1",
		},
	}

	result, err := validator.ValidateCompatibility(manifest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// In the reference implementation, this should always return valid with warnings
	if !result.Valid {
		t.Errorf("expected compatibility validation to be valid in reference implementation")
	}

	if len(result.Warnings) == 0 {
		t.Errorf("expected at least one warning about unimplemented compatibility checks")
	}
}

func TestScanSecurity(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name           string
		binaryData     []byte
		expectedStatus SecurityStatus
		expectedVulns  int
	}{
		{
			name:           "empty binary",
			binaryData:     []byte{},
			expectedStatus: SecurityStatusVulnerable,
			expectedVulns:  1,
		},
		{
			name:           "non-empty binary",
			binaryData:     []byte("test binary data"),
			expectedStatus: SecurityStatusSecure,
			expectedVulns:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := validator.ScanSecurity(tt.binaryData)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if results.Status != tt.expectedStatus {
				t.Errorf("expected status %s, got %s", tt.expectedStatus, results.Status)
			}

			if len(results.Vulnerabilities) != tt.expectedVulns {
				t.Errorf("expected %d vulnerabilities, got %d", tt.expectedVulns, len(results.Vulnerabilities))
			}

			if results.ScanVersion == "" {
				t.Error("expected scan version to be set")
			}
		})
	}
}

func TestValidationHelpers(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func() error
		valid    bool
	}{
		{
			name:     "valid plugin name",
			testFunc: func() error { return validatePluginName("test-plugin") },
			valid:    true,
		},
		{
			name:     "invalid plugin name - uppercase",
			testFunc: func() error { return validatePluginName("Test-Plugin") },
			valid:    false,
		},
		{
			name:     "invalid plugin name - too short",
			testFunc: func() error { return validatePluginName("ab") },
			valid:    false,
		},
		{
			name:     "valid semver",
			testFunc: func() error { return validateSemVer("1.0.0") },
			valid:    true,
		},
		{
			name:     "valid semver with prerelease",
			testFunc: func() error { return validateSemVer("1.0.0-beta.1") },
			valid:    true,
		},
		{
			name:     "invalid semver",
			testFunc: func() error { return validateSemVer("1.0") },
			valid:    false,
		},
		{
			name:     "valid api version",
			testFunc: func() error { return validateAPIVersion("v1") },
			valid:    true,
		},
		{
			name:     "invalid api version",
			testFunc: func() error { return validateAPIVersion("v0") },
			valid:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.testFunc()
			if tt.valid && err != nil {
				t.Errorf("expected validation to pass, but got error: %v", err)
			}
			if !tt.valid && err == nil {
				t.Error("expected validation to fail, but got no error")
			}
		})
	}
}