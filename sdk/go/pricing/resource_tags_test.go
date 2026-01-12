package pricing_test

import (
	"testing"

	"github.com/rshade/finfocus-spec/sdk/go/pricing"
)

func TestValidatePricingSpec_ResourceTags(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		wantErr  bool
	}{
		{
			name: "Valid resource tags - billing center scenario",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0104,
				"currency": "USD",
				"resource_tags": {
					"billing_center": "engineering",
					"cost_center": "CC-1001",
					"environment": "production",
					"team": "platform",
					"project": "microservices"
				}
			}`,
			wantErr: false,
		},
		{
			name: "Valid resource tags - Flexera CCO scenario",
			jsonData: `{
				"provider": "azure",
				"resource_type": "vm",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.05,
				"currency": "USD",
				"resource_tags": {
					"flexera_billing_center": "Engineering-Prod",
					"department": "Research and Development",
					"budget_code": "R&D-2024-Q1",
					"approval_level": "manager",
					"charge_back_group": "product-development"
				}
			}`,
			wantErr: false,
		},
		{
			name: "Valid resource tags - empty tags object",
			jsonData: `{
				"provider": "gcp",
				"resource_type": "compute_engine",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.03,
				"currency": "USD",
				"resource_tags": {}
			}`,
			wantErr: false,
		},
		{
			name: "Invalid resource tags - non-string value",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0104,
				"currency": "USD",
				"resource_tags": {
					"billing_center": "engineering",
					"cost_amount": 1000.50
				}
			}`,
			wantErr: true,
		},
		{
			name: "Invalid resource tags - array instead of object",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0104,
				"currency": "USD",
				"resource_tags": ["tag1", "tag2"]
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pricing.ValidatePricingSpec([]byte(tt.jsonData))
			if (err != nil) != tt.wantErr {
				t.Errorf("pricing.ValidatePricingSpec() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePricingSpec_EnhancedPluginMetadata(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		wantErr  bool
	}{
		{
			name: "Valid enhanced metadata - AWS Reserved Instance",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "reserved",
				"rate_per_unit": 0.025,
				"currency": "USD",
				"plugin_metadata": {
					"aws_account_id": "123456789012",
					"reservation_id": "r-1234567890abcdef0",
					"availability_zone": "us-east-1a",
					"instance_family": "general_purpose",
					"commitment_discount": 0.4,
					"upfront_payment": 1000.00,
					"effective_hourly_rate": 0.025,
					"pricing_model_details": {
						"term_length": "1_year",
						"offering_class": "standard",
						"tenancy": "shared"
					}
				}
			}`,
			wantErr: false,
		},
		{
			name: "Valid enhanced metadata - Flexera complex structure",
			jsonData: `{
				"provider": "custom",
				"resource_type": "flexera_cco_resource",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.08,
				"currency": "USD",
				"plugin_metadata": {
					"flexera_org_id": "12345",
					"billing_center_hierarchy": {
						"root": "Corporate",
						"division": "Engineering",
						"department": "Platform",
						"team": "Infrastructure"
					},
					"cost_allocation_rules": [
						{
							"rule_id": "rule-001",
							"percentage": 60,
							"target_center": "Engineering-Prod"
						},
						{
							"rule_id": "rule-002",
							"percentage": 40,
							"target_center": "Engineering-Dev"
						}
					],
					"approval_workflow": {
						"required": true,
						"approver": "john.doe@company.com",
						"threshold": 500.00
					}
				}
			}`,
			wantErr: false,
		},
		{
			name: "Valid enhanced metadata - mixed data types",
			jsonData: `{
				"provider": "azure",
				"resource_type": "cosmos_db",
				"billing_mode": "per_ru",
				"rate_per_unit": 0.00008,
				"currency": "USD",
				"plugin_metadata": {
					"string_field": "test",
					"number_field": 42,
					"float_field": 3.14159,
					"boolean_field": true,
					"null_field": null,
					"array_field": [1, 2, 3],
					"object_field": {
						"nested_string": "value",
						"nested_number": 100
					}
				}
			}`,
			wantErr: false,
		},
		{
			name: "Valid enhanced metadata - empty object",
			jsonData: `{
				"provider": "gcp",
				"resource_type": "cloud_storage",
				"billing_mode": "per_gb_month",
				"rate_per_unit": 0.02,
				"currency": "USD",
				"plugin_metadata": {}
			}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pricing.ValidatePricingSpec([]byte(tt.jsonData))
			if (err != nil) != tt.wantErr {
				t.Errorf("pricing.ValidatePricingSpec() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePricingSpec_CombinedTagsAndMetadata(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		wantErr  bool
	}{
		{
			name: "Valid combined tags and enhanced metadata",
			jsonData: `{
				"provider": "aws",
				"resource_type": "s3",
				"billing_mode": "per_gb_month",
				"rate_per_unit": 0.023,
				"currency": "USD",
				"description": "S3 Standard storage with comprehensive tagging and metadata",
				"resource_tags": {
					"billing_center": "data-engineering",
					"cost_center": "CC-3001",
					"environment": "production",
					"data_classification": "confidential",
					"retention_policy": "7_years",
					"backup_required": "true"
				},
				"plugin_metadata": {
					"aws_account_id": "123456789012",
					"bucket_region": "us-east-1",
					"storage_class": "STANDARD",
					"versioning_enabled": true,
					"encryption": {
						"type": "AES256",
						"kms_key_id": "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"
					},
					"lifecycle_policies": [
						{
							"rule_name": "transition_to_ia",
							"days": 30,
							"target_storage_class": "STANDARD_IA"
						},
						{
							"rule_name": "transition_to_glacier",
							"days": 90,
							"target_storage_class": "GLACIER"
						}
					],
					"cost_optimization": {
						"intelligent_tiering_enabled": true,
						"estimated_monthly_savings": 125.50,
						"optimization_recommendations": [
							"Enable Intelligent Tiering",
							"Review lifecycle policies quarterly"
						]
					}
				},
				"pricing_tiers": [
					{
						"min_units": 0,
						"max_units": 50000,
						"rate_per_unit": 0.023
					},
					{
						"min_units": 50000,
						"rate_per_unit": 0.022
					}
				],
				"source": "aws"
			}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pricing.ValidatePricingSpec([]byte(tt.jsonData))
			if (err != nil) != tt.wantErr {
				t.Errorf("pricing.ValidatePricingSpec() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
