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

package registry_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/rshade/pulumicost-spec/sdk/go/registry"
)

func TestValidatePluginManifest_ValidManifests(t *testing.T) {
	validManifests := []string{
		// Minimal valid manifest
		`{
			"metadata": {
				"name": "test-plugin",
				"version": "1.0.0",
				"description": "Test plugin for validation",
				"author": "Test Author"
			},
			"specification": {
				"spec_version": "0.1.0",
				"supported_providers": ["custom"],
				"service_definition": {
					"service_name": "CostSourceService",
					"package_name": "pulumicost.v1",
					"methods": ["Name"]
				}
			},
			"installation": {
				"installation_method": "binary"
			}
		}`,
		// Full AWS plugin example
		`{
			"metadata": {
				"name": "aws-cost-plugin",
				"version": "2.1.0",
				"description": "Comprehensive AWS cost source plugin supporting EC2, S3, Lambda, RDS, and DynamoDB with real-time and historical cost data retrieval",
				"author": "PulumiCost Team"
			},
			"specification": {
				"spec_version": "0.1.0",
				"supported_providers": ["aws"],
				"service_definition": {
					"service_name": "CostSourceService",
					"package_name": "pulumicost.v1",
					"methods": ["Name", "Supports", "GetActualCost", "GetProjectedCost", "GetPricingSpec"]
				}
			},
			"installation": {
				"installation_method": "container"
			}
		}`,
	}

	for i, manifestJSON := range validManifests {
		err := registry.ValidatePluginManifest([]byte(manifestJSON))
		if err != nil {
			t.Errorf("Valid manifest %d failed validation: %v", i+1, err)
		}
	}
}

func TestValidatePluginManifest_InvalidManifests(t *testing.T) {
	invalidManifests := []struct {
		name                  string
		manifest              string
		expectedErrorContains string
	}{
		{
			name:     "empty object",
			manifest: `{}`,
		},
		{
			name: "missing metadata",
			manifest: `{
				"specification": {
					"spec_version": "0.1.0",
					"supported_providers": ["custom"],
					"service_definition": {
						"service_name": "CostSourceService",
						"package_name": "pulumicost.v1",
						"methods": ["Name"]
					}
				},
				"installation": {
					"installation_method": "binary"
				}
			}`,
		},
		{
			name: "invalid plugin name",
			manifest: `{
				"metadata": {
					"name": "Invalid_Plugin",
					"version": "1.0.0",
					"description": "Test plugin for validation",
					"author": "Test Author"
				},
				"specification": {
					"spec_version": "0.1.0",
					"supported_providers": ["custom"],
					"service_definition": {
						"service_name": "CostSourceService",
						"package_name": "pulumicost.v1",
						"methods": ["Name"]
					}
				},
				"installation": {
					"installation_method": "binary"
				}
			}`,
		},
		{
			name: "invalid version",
			manifest: `{
				"metadata": {
					"name": "test-plugin",
					"version": "invalid-version",
					"description": "Test plugin for validation",
					"author": "Test Author"
				},
				"specification": {
					"spec_version": "0.1.0",
					"supported_providers": ["custom"],
					"service_definition": {
						"service_name": "CostSourceService",
						"package_name": "pulumicost.v1",
						"methods": ["Name"]
					}
				},
				"installation": {
					"installation_method": "binary"
				}
			}`,
		},
		{
			name: "description too short",
			manifest: `{
				"metadata": {
					"name": "test-plugin",
					"version": "1.0.0",
					"description": "Short",
					"author": "Test Author"
				},
				"specification": {
					"spec_version": "0.1.0",
					"supported_providers": ["custom"],
					"service_definition": {
						"service_name": "CostSourceService",
						"package_name": "pulumicost.v1",
						"methods": ["Name"]
					}
				},
				"installation": {
					"installation_method": "binary"
				}
			}`,
		},
		{
			name: "invalid provider",
			manifest: `{
				"metadata": {
					"name": "test-plugin",
					"version": "1.0.0",
					"description": "Test plugin for validation",
					"author": "Test Author"
				},
				"specification": {
					"spec_version": "0.1.0",
					"supported_providers": ["invalid-provider"],
					"service_definition": {
						"service_name": "CostSourceService",
						"package_name": "pulumicost.v1",
						"methods": ["Name"]
					}
				},
				"installation": {
					"installation_method": "binary"
				}
			}`,
			expectedErrorContains: "must be one of: aws, azure, gcp, kubernetes, custom",
		},
		{
			name: "invalid service name",
			manifest: `{
				"metadata": {
					"name": "test-plugin",
					"version": "1.0.0",
					"description": "Test plugin for validation",
					"author": "Test Author"
				},
				"specification": {
					"spec_version": "0.1.0",
					"supported_providers": ["custom"],
					"service_definition": {
						"service_name": "invalidServiceName",
						"package_name": "pulumicost.v1",
						"methods": ["Name"]
					}
				},
				"installation": {
					"installation_method": "binary"
				}
			}`,
		},
		{
			name: "invalid package name",
			manifest: `{
				"metadata": {
					"name": "test-plugin",
					"version": "1.0.0",
					"description": "Test plugin for validation",
					"author": "Test Author"
				},
				"specification": {
					"spec_version": "0.1.0",
					"supported_providers": ["custom"],
					"service_definition": {
						"service_name": "CostSourceService",
						"package_name": "Invalid.Package.Name",
						"methods": ["Name"]
					}
				},
				"installation": {
					"installation_method": "binary"
				}
			}`,
		},
		{
			name: "invalid method",
			manifest: `{
				"metadata": {
					"name": "test-plugin",
					"version": "1.0.0",
					"description": "Test plugin for validation",
					"author": "Test Author"
				},
				"specification": {
					"spec_version": "0.1.0",
					"supported_providers": ["custom"],
					"service_definition": {
						"service_name": "CostSourceService",
						"package_name": "pulumicost.v1",
						"methods": ["InvalidMethod"]
					}
				},
				"installation": {
					"installation_method": "binary"
				}
			}`,
		},
		{
			name: "invalid installation method",
			manifest: `{
				"metadata": {
					"name": "test-plugin",
					"version": "1.0.0",
					"description": "Test plugin for validation",
					"author": "Test Author"
				},
				"specification": {
					"spec_version": "0.1.0",
					"supported_providers": ["custom"],
					"service_definition": {
						"service_name": "CostSourceService",
						"package_name": "pulumicost.v1",
						"methods": ["Name"]
					}
				},
				"installation": {
					"installation_method": "invalid-method"
				}
			}`,
		},
	}

	for _, test := range invalidManifests {
		err := registry.ValidatePluginManifest([]byte(test.manifest))
		if err == nil {
			t.Errorf("Invalid manifest '%s' should have failed validation", test.name)
		} else if test.expectedErrorContains != "" {
			if !strings.Contains(err.Error(), test.expectedErrorContains) {
				t.Errorf("Invalid manifest '%s' error message should contain '%s', got: %v",
					test.name, test.expectedErrorContains, err)
			}
		}
	}
}

func TestValidatePluginManifest_InvalidJSON(t *testing.T) {
	invalidJSON := `{ "invalid": json }`
	err := registry.ValidatePluginManifest([]byte(invalidJSON))
	if err == nil {
		t.Errorf("Invalid JSON should have failed validation")
	}
}

func TestValidatePluginManifest_Examples(t *testing.T) {
	// Test that our example plugin manifests are valid
	exampleManifests := []string{
		// This should match the format from our examples directory
		`{
			"metadata": {
				"name": "aws-cost-plugin",
				"version": "2.1.0",
				"description": "Comprehensive AWS cost source plugin supporting EC2, S3, Lambda, RDS, and DynamoDB with real-time and historical cost data retrieval",
				"author": "PulumiCost Team",
				"homepage": "https://pulumicost.dev/plugins/aws",
				"repository": "https://github.com/pulumicost/aws-plugin",
				"license": "Apache-2.0",
				"keywords": ["aws", "cost", "ec2", "s3", "lambda", "rds", "dynamodb"],
				"created_at": "2024-01-15T10:00:00Z",
				"updated_at": "2024-03-20T14:30:00Z"
			},
			"specification": {
				"spec_version": "0.1.0",
				"supported_providers": ["aws"],
				"supported_resources": {
					"aws": {
						"resource_types": ["ec2", "s3", "lambda", "rds", "dynamodb"],
						"billing_modes": [
							"per_hour", "per_gb_month", "per_invocation",
							"per_rcu", "per_wcu", "reserved", "spot"
						],
						"regions": [
							"us-east-1", "us-west-2", "eu-west-1",
							"ap-southeast-1", "ap-northeast-1"
						]
					}
				},
				"capabilities": [
					"cost_retrieval", "cost_projection", "pricing_specs",
					"historical_data", "real_time_data", "caching", "filtering"
				],
				"service_definition": {
					"service_name": "CostSourceService",
					"package_name": "pulumicost.v1",
					"methods": ["Name", "Supports", "GetActualCost", "GetProjectedCost", "GetPricingSpec"],
					"port": 50051,
					"health_check_path": "/health"
				},
				"observability_support": {
					"metrics_enabled": true,
					"tracing_enabled": true,
					"logging_enabled": true,
					"health_checks_enabled": true,
					"sli_support": true
				}
			},
			"security": {
				"signature": "MEUCIQCx7HjRFkL3+Y8XrGQm4nW2V9iE2fP8jS6bK1qN7dR9yAIgD8sJ5tK2oP1mN3vB4cE5fG6hI7jK8lM9nO0pQ1rS2tU=",
				"public_key": "-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1234567890abcdef...
-----END PUBLIC KEY-----",
				"security_level": "verified",
				"permissions": [
					"network_access", "filesystem_read", "config_read", "temp_files"
				],
				"sandbox_required": false
			},
			"installation": {
				"installation_method": "binary",
				"download_url": "https://releases.pulumicost.dev/aws-plugin/v2.1.0/aws-plugin-linux-amd64.tar.gz",
				"checksum": "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456",
				"checksum_algorithm": "sha256",
				"pre_install_checks": [
					"verify_aws_credentials",
					"check_network_connectivity",
					"validate_permissions"
				],
				"post_install_steps": [
					"create_config_directory",
					"setup_log_rotation",
					"register_health_check"
				]
			}
		}`,
	}

	for i, manifestJSON := range exampleManifests {
		// First check if it's valid JSON
		var manifest map[string]interface{}
		if err := json.Unmarshal([]byte(manifestJSON), &manifest); err != nil {
			t.Errorf("Example manifest %d is not valid JSON: %v", i+1, err)
			continue
		}

		// Then validate against our validation function
		// Note: This validates against our simplified schema embedded in validate.go
		// The full schema validation would be done by the JSON schema validator
		err := registry.ValidatePluginManifest([]byte(manifestJSON))
		if err != nil {
			t.Errorf("Example manifest %d failed validation: %v", i+1, err)
		}
	}
}
