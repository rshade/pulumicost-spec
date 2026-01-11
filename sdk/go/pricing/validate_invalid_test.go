package pricing_test

import (
	"testing"

	"github.com/rshade/finfocus-spec/sdk/go/pricing"
)

func TestValidatePricingSpec_InvalidSchemas(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		wantErr  bool
	}{
		{
			name: "Missing required field: provider",
			jsonData: `{
				"resource_type": "ec2",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0104,
				"currency": "USD"
			}`,
			wantErr: true,
		},
		{
			name: "Missing required field: resource_type",
			jsonData: `{
				"provider": "aws",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0104,
				"currency": "USD"
			}`,
			wantErr: true,
		},
		{
			name: "Missing required field: billing_mode",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"rate_per_unit": 0.0104,
				"currency": "USD"
			}`,
			wantErr: true,
		},
		{
			name: "Missing required field: rate_per_unit",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "per_hour",
				"currency": "USD"
			}`,
			wantErr: true,
		},
		{
			name: "Missing required field: currency",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0104
			}`,
			wantErr: true,
		},
		{
			name: "Invalid provider",
			jsonData: `{
				"provider": "invalid_provider",
				"resource_type": "ec2",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0104,
				"currency": "USD"
			}`,
			wantErr: true,
		},
		{
			name: "Invalid billing_mode",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "invalid_billing_mode",
				"rate_per_unit": 0.0104,
				"currency": "USD"
			}`,
			wantErr: true,
		},
		{
			name: "Negative rate_per_unit",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "per_hour",
				"rate_per_unit": -0.0104,
				"currency": "USD"
			}`,
			wantErr: true,
		},
		{
			name: "Invalid currency format (too short)",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0104,
				"currency": "US"
			}`,
			wantErr: true,
		},
		{
			name: "Invalid currency format (too long)",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0104,
				"currency": "USDD"
			}`,
			wantErr: true,
		},
		{
			name: "Invalid currency format (lowercase)",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0104,
				"currency": "usd"
			}`,
			wantErr: true,
		},
		{
			name: "Invalid provider enum value (empty string)",
			jsonData: `{
				"provider": "",
				"resource_type": "ec2",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0104,
				"currency": "USD"
			}`,
			wantErr: true,
		},
		{
			name: "Invalid resource_type (empty string)",
			jsonData: `{
				"provider": "aws",
				"resource_type": "",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0104,
				"currency": "USD"
			}`,
			wantErr: true,
		},
		{
			name: "Invalid JSON syntax",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2"
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0104,
				"currency": "USD"
			}`,
			wantErr: true,
		},
		{
			name: "Additional properties not allowed",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0104,
				"currency": "USD",
				"unknown_field": "value"
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

func TestValidatePricingSpec_InvalidMetricHints(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		wantErr  bool
	}{
		{
			name: "metric_hints missing required metric field",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0104,
				"currency": "USD",
				"metric_hints": [
					{
						"unit": "hour"
					}
				]
			}`,
			wantErr: true,
		},
		{
			name: "metric_hints missing required unit field",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0104,
				"currency": "USD",
				"metric_hints": [
					{
						"metric": "vcpu_hours"
					}
				]
			}`,
			wantErr: true,
		},
		{
			name: "metric_hints invalid aggregation_method",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0104,
				"currency": "USD",
				"metric_hints": [
					{
						"metric": "vcpu_hours",
						"unit": "hour",
						"aggregation_method": "invalid_method"
					}
				]
			}`,
			wantErr: true,
		},
		{
			name: "metric_hints empty metric string",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0104,
				"currency": "USD",
				"metric_hints": [
					{
						"metric": "",
						"unit": "hour"
					}
				]
			}`,
			wantErr: true,
		},
		{
			name: "metric_hints empty unit string",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0104,
				"currency": "USD",
				"metric_hints": [
					{
						"metric": "vcpu_hours",
						"unit": ""
					}
				]
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

func TestValidatePricingSpec_InvalidPricingTiers(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		wantErr  bool
	}{
		{
			name: "pricing_tiers missing required min_units",
			jsonData: `{
				"provider": "aws",
				"resource_type": "s3",
				"billing_mode": "per_gb_month",
				"rate_per_unit": 0.023,
				"currency": "USD",
				"pricing_tiers": [
					{
						"rate_per_unit": 0.023
					}
				]
			}`,
			wantErr: true,
		},
		{
			name: "pricing_tiers missing required rate_per_unit",
			jsonData: `{
				"provider": "aws",
				"resource_type": "s3",
				"billing_mode": "per_gb_month",
				"rate_per_unit": 0.023,
				"currency": "USD",
				"pricing_tiers": [
					{
						"min_units": 0
					}
				]
			}`,
			wantErr: true,
		},
		{
			name: "pricing_tiers negative min_units",
			jsonData: `{
				"provider": "aws",
				"resource_type": "s3",
				"billing_mode": "per_gb_month",
				"rate_per_unit": 0.023,
				"currency": "USD",
				"pricing_tiers": [
					{
						"min_units": -1,
						"rate_per_unit": 0.023
					}
				]
			}`,
			wantErr: true,
		},
		{
			name: "pricing_tiers negative rate_per_unit",
			jsonData: `{
				"provider": "aws",
				"resource_type": "s3",
				"billing_mode": "per_gb_month",
				"rate_per_unit": 0.023,
				"currency": "USD",
				"pricing_tiers": [
					{
						"min_units": 0,
						"rate_per_unit": -0.023
					}
				]
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

func TestValidatePricingSpec_InvalidTimeAggregation(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		wantErr  bool
	}{
		{
			name: "time_aggregation invalid window",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0104,
				"currency": "USD",
				"time_aggregation": {
					"window": "invalid_window",
					"method": "sum",
					"alignment": "calendar"
				}
			}`,
			wantErr: true,
		},
		{
			name: "time_aggregation invalid method",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0104,
				"currency": "USD",
				"time_aggregation": {
					"window": "hour",
					"method": "invalid_method",
					"alignment": "calendar"
				}
			}`,
			wantErr: true,
		},
		{
			name: "time_aggregation invalid alignment",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0104,
				"currency": "USD",
				"time_aggregation": {
					"window": "hour",
					"method": "sum",
					"alignment": "invalid_alignment"
				}
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

func TestValidatePricingSpec_InvalidCommitmentTerms(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		wantErr  bool
	}{
		{
			name: "commitment_terms invalid duration",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "reserved",
				"rate_per_unit": 0.0062,
				"currency": "USD",
				"commitment_terms": {
					"duration": "invalid_duration",
					"payment_option": "no_upfront",
					"discount_percentage": 40.4
				}
			}`,
			wantErr: true,
		},
		{
			name: "commitment_terms invalid payment_option",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "reserved",
				"rate_per_unit": 0.0062,
				"currency": "USD",
				"commitment_terms": {
					"duration": "1_year",
					"payment_option": "invalid_payment",
					"discount_percentage": 40.4
				}
			}`,
			wantErr: true,
		},
		{
			name: "commitment_terms negative discount_percentage",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "reserved",
				"rate_per_unit": 0.0062,
				"currency": "USD",
				"commitment_terms": {
					"duration": "1_year",
					"payment_option": "no_upfront",
					"discount_percentage": -10
				}
			}`,
			wantErr: true,
		},
		{
			name: "commitment_terms discount_percentage over 100",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "reserved",
				"rate_per_unit": 0.0062,
				"currency": "USD",
				"commitment_terms": {
					"duration": "1_year",
					"payment_option": "no_upfront",
					"discount_percentage": 110
				}
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

func TestValidatePricingSpec_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		wantErr  bool
	}{
		{
			name: "valid effective_date format",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0104,
				"currency": "USD",
				"effective_date": "2024-01-01T00:00:00Z"
			}`,
			wantErr: false,
		},
		{
			name: "valid expiration_date format",
			jsonData: `{
				"provider": "aws",
				"resource_type": "ec2",
				"billing_mode": "per_hour",
				"rate_per_unit": 0.0104,
				"currency": "USD",
				"expiration_date": "2024-12-31T23:59:59Z"
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
