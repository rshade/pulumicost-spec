package types

import (
	"testing"
)

func TestValidatePricingSpec_ValidAWSExamples(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
	}{
		{
			name: "AWS EC2 On-Demand",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"sku": "t3.micro",
				"region": "us-east-1",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0104,
				"currency": "USD",
				"description": "General Purpose t3.micro instance",
				"metric_hints": [
					{
						"metric": "vcpu_hours",
						"unit": "hour",
						"aggregation_method": "sum"
					}
				],
				"source": "aws"
			}`,
		},
		{
			name: "AWS EC2 Reserved Instance",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"sku": "t3.micro",
				"region": "us-east-1",
				"billing_mode": "reserved",
				"rate_per_unit": 0.0062,
				"currency": "USD",
				"description": "Reserved t3.micro instance with 1-year commitment",
				"commitment_terms": {
					"duration": "1_year",
					"payment_option": "no_upfront",
					"discount_percentage": 40.4
				},
				"source": "aws"
			}`,
		},
		{
			name: "AWS S3 Storage",
			jsonData: `{
				"provider": "aws",
				"resource_type": "s3",
				"region": "us-east-1",
				"billing_mode": "per_gb_month",
				"rate_per_unit": 0.023,
				"currency": "USD",
				"description": "Standard storage pricing",
				"metric_hints": [
					{
						"metric": "storage_gb",
						"unit": "GB",
						"aggregation_method": "avg"
					}
				],
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
		},
		{
			name: "AWS Lambda",
			jsonData: `{
				"provider": "aws",
				"resource_type": "lambda",
				"billing_mode": "per_invocation",
				"rate_per_unit": 0.0000002,
				"currency": "USD",
				"description": "Lambda request charges",
				"metric_hints": [
					{
						"metric": "invocations",
						"unit": "count",
						"aggregation_method": "sum"
					}
				],
				"time_aggregation": {
					"window": "month",
					"method": "sum",
					"alignment": "billing"
				},
				"source": "aws"
			}`,
		},
		{
			name: "AWS DynamoDB",
			jsonData: `{
				"provider": "aws",
				"resource_type": "dynamodb",
				"region": "us-east-1",
				"billing_mode": "per_rcu",
				"rate_per_unit": 0.00013,
				"currency": "USD",
				"description": "DynamoDB read capacity units",
				"metric_hints": [
					{
						"metric": "read_capacity_units",
						"unit": "hour",
						"aggregation_method": "avg"
					}
				],
				"source": "aws"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePricingSpec([]byte(tt.jsonData))
			if err != nil {
				t.Errorf("ValidatePricingSpec() error = %v, expected nil", err)
			}
		})
	}
}

func TestValidatePricingSpec_ValidAzureExamples(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
	}{
		{
			name: "Azure VM",
			jsonData: `{
				"provider": "azure",
				"resource_type": "vm",
				"sku": "Standard_B1s",
				"region": "eastus",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.00464,
				"currency": "USD",
				"description": "Basic B1s virtual machine",
				"metric_hints": [
					{
						"metric": "vcpu_hours",
						"unit": "hour",
						"aggregation_method": "sum"
					}
				],
				"source": "azure"
			}`,
		},
		{
			name: "Azure Blob Storage",
			jsonData: `{
				"provider": "azure",
				"resource_type": "blob_storage",
				"region": "eastus",
				"billing_mode": "per_gb_month",
				"rate_per_unit": 0.0184,
				"currency": "USD",
				"description": "Hot access tier storage",
				"metric_hints": [
					{
						"metric": "storage_gb",
						"unit": "GB",
						"aggregation_method": "avg"
					}
				],
				"source": "azure"
			}`,
		},
		{
			name: "Azure SQL Database DTU",
			jsonData: `{
				"provider": "azure",
				"resource_type": "sql_database",
				"sku": "S1",
				"region": "eastus",
				"billing_mode": "per_dtu",
				"rate_per_unit": 0.02,
				"currency": "USD",
				"description": "Standard S1 SQL Database",
				"metric_hints": [
					{
						"metric": "database_transaction_units",
						"unit": "hour",
						"aggregation_method": "avg"
					}
				],
				"source": "azure"
			}`,
		},
		{
			name: "Azure Cosmos DB",
			jsonData: `{
				"provider": "azure",
				"resource_type": "cosmos_db",
				"region": "eastus",
				"billing_mode": "per_ru",
				"rate_per_unit": 0.00008,
				"currency": "USD",
				"description": "Cosmos DB request units",
				"metric_hints": [
					{
						"metric": "request_units",
						"unit": "hour",
						"aggregation_method": "avg"
					}
				],
				"source": "azure"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePricingSpec([]byte(tt.jsonData))
			if err != nil {
				t.Errorf("ValidatePricingSpec() error = %v, expected nil", err)
			}
		})
	}
}

func TestValidatePricingSpec_ValidGCPExamples(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
	}{
		{
			name: "GCP Compute Engine",
			jsonData: `{
				"provider": "gcp",
				"resource_type": "compute_engine",
				"sku": "n1-standard-1",
				"region": "us-central1",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0475,
				"currency": "USD",
				"description": "N1 standard machine type",
				"metric_hints": [
					{
						"metric": "vcpu_hours",
						"unit": "hour",
						"aggregation_method": "sum"
					}
				],
				"source": "gcp"
			}`,
		},
		{
			name: "GCP Preemptible Instance",
			jsonData: `{
				"provider": "gcp",
				"resource_type": "compute_engine",
				"sku": "n1-standard-1",
				"region": "us-central1",
				"billing_mode": "preemptible",
				"rate_per_unit": 0.01,
				"currency": "USD",
				"description": "Preemptible N1 standard machine",
				"commitment_terms": {
					"duration": "spot",
					"discount_percentage": 79.0
				},
				"source": "gcp"
			}`,
		},
		{
			name: "GCP Cloud Storage",
			jsonData: `{
				"provider": "gcp",
				"resource_type": "cloud_storage",
				"region": "us-central1",
				"billing_mode": "per_gb_month",
				"rate_per_unit": 0.02,
				"currency": "USD",
				"description": "Standard storage class",
				"metric_hints": [
					{
						"metric": "storage_gb",
						"unit": "GB",
						"aggregation_method": "avg"
					}
				],
				"source": "gcp"
			}`,
		},
		{
			name: "GCP Cloud Functions",
			jsonData: `{
				"provider": "gcp",
				"resource_type": "cloud_functions",
				"billing_mode": "per_invocation",
				"rate_per_unit": 0.0000004,
				"currency": "USD",
				"description": "Cloud Functions invocation charges",
				"metric_hints": [
					{
						"metric": "invocations",
						"unit": "count",
						"aggregation_method": "sum"
					}
				],
				"source": "gcp"
			}`,
		},
		{
			name: "GCP Committed Use Discount",
			jsonData: `{
				"provider": "gcp",
				"resource_type": "compute_engine",
				"sku": "n1-standard-1",
				"region": "us-central1",
				"billing_mode": "committed_use",
				"rate_per_unit": 0.025,
				"currency": "USD",
				"description": "1-year committed use discount",
				"commitment_terms": {
					"duration": "1_year",
					"payment_option": "monthly",
					"discount_percentage": 47.4
				},
				"source": "gcp"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePricingSpec([]byte(tt.jsonData))
			if err != nil {
				t.Errorf("ValidatePricingSpec() error = %v, expected nil", err)
			}
		})
	}
}

func TestValidatePricingSpec_ValidKubernetesExamples(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
	}{
		{
			name: "Kubernetes Namespace",
			jsonData: `{
				"provider": "kubernetes",
				"resource_type": "namespace",
				"billing_mode": "per_cpu_hour",
				"rate_per_unit": 0.03,
				"currency": "USD",
				"description": "Kubernetes namespace CPU pricing",
				"metric_hints": [
					{
						"metric": "cpu_cores",
						"unit": "hour",
						"aggregation_method": "avg"
					}
				],
				"source": "kubecost"
			}`,
		},
		{
			name: "Kubernetes Memory",
			jsonData: `{
				"provider": "kubernetes",
				"resource_type": "namespace",
				"billing_mode": "per_memory_gb_hour",
				"rate_per_unit": 0.004,
				"currency": "USD",
				"description": "Kubernetes namespace memory pricing",
				"metric_hints": [
					{
						"metric": "memory_gb",
						"unit": "hour",
						"aggregation_method": "avg"
					}
				],
				"source": "kubecost"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePricingSpec([]byte(tt.jsonData))
			if err != nil {
				t.Errorf("ValidatePricingSpec() error = %v, expected nil", err)
			}
		})
	}
}

func TestValidatePricingSpec_ComplexValidExamples(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
	}{
		{
			name: "Complete AWS S3 with tiers and time aggregation",
			jsonData: `{
				"provider": "aws",
				"resource_type": "s3",
				"region": "us-east-1",
				"billing_mode": "per_gb_month",
				"rate_per_unit": 0.023,
				"currency": "USD",
				"description": "Standard storage with tiered pricing",
				"metric_hints": [
					{
						"metric": "storage_gb",
						"unit": "GB",
						"aggregation_method": "avg"
					},
					{
						"metric": "requests",
						"unit": "count",
						"aggregation_method": "sum"
					}
				],
				"pricing_tiers": [
					{
						"min_units": 0,
						"max_units": 50000,
						"rate_per_unit": 0.023
					},
					{
						"min_units": 50000,
						"max_units": 450000,
						"rate_per_unit": 0.022
					},
					{
						"min_units": 450000,
						"rate_per_unit": 0.021
					}
				],
				"time_aggregation": {
					"window": "month",
					"method": "avg",
					"alignment": "billing"
				},
				"resource_tags": {
					"billing_center": "engineering",
					"cost_center": "CC-1001",
					"environment": "production",
					"team": "storage"
				},
				"plugin_metadata": {
					"storage_class": "STANDARD",
					"availability": "99.999999999",
					"aws_account_id": "123456789012",
					"pricing_tier_details": {
						"tier_1_limit": 50000,
						"tier_2_limit": 450000,
						"tier_3_unlimited": true
					}
				},
				"source": "aws",
				"effective_date": "2024-01-01T00:00:00Z",
				"expiration_date": "2024-12-31T23:59:59Z"
			}`,
		},
		{
			name: "Azure VM with hybrid benefit",
			jsonData: `{
				"provider": "azure",
				"resource_type": "vm",
				"sku": "Standard_D2s_v3",
				"region": "eastus",
				"billing_mode": "hybrid_benefit",
				"rate_per_unit": 0.096,
				"currency": "USD",
				"description": "Windows VM with Azure Hybrid Benefit",
				"metric_hints": [
					{
						"metric": "vcpu_hours",
						"unit": "hour",
						"aggregation_method": "sum"
					},
					{
						"metric": "memory_gb_hours",
						"unit": "hour",
						"aggregation_method": "sum"
					}
				],
				"commitment_terms": {
					"duration": "on_demand",
					"discount_percentage": 49.0
				},
				"time_aggregation": {
					"window": "hour",
					"method": "sum",
					"alignment": "calendar"
				},
				"resource_tags": {
					"billing_center": "it-infrastructure", 
					"cost_center": "CC-2002",
					"environment": "production",
					"application": "web-server",
					"owner": "platform-team"
				},
				"plugin_metadata": {
					"os_type": "windows",
					"license_included": "false",
					"azure_subscription_id": "12345678-1234-1234-1234-123456789012",
					"resource_group": "prod-web-servers",
					"vm_size_family": "Dsv3",
					"hybrid_benefit_savings": {
						"original_rate": 0.192,
						"discounted_rate": 0.096,
						"savings_percentage": 50.0
					}
				},
				"source": "azure"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePricingSpec([]byte(tt.jsonData))
			if err != nil {
				t.Errorf("ValidatePricingSpec() error = %v, expected nil", err)
			}
		})
	}
}