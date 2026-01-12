package pricing_test

import (
	"testing"

	"github.com/rshade/finfocus-spec/sdk/go/pricing"
)

func TestValidBillingMode(t *testing.T) {
	tests := []struct {
		name        string
		billingMode string
		expectValid bool
	}{
		// Time-based
		{"per_hour valid", "per_hour", true},
		{"per_minute valid", "per_minute", true},
		{"per_second valid", "per_second", true},
		{"per_day valid", "per_day", true},
		{"per_month valid", "per_month", true},
		{"per_year valid", "per_year", true},

		// Storage-based
		{"per_gb_month valid", "per_gb_month", true},
		{"per_gb_hour valid", "per_gb_hour", true},
		{"per_gb_day valid", "per_gb_day", true},

		// Usage-based
		{"per_request valid", "per_request", true},
		{"per_operation valid", "per_operation", true},
		{"per_transaction valid", "per_transaction", true},
		{"per_execution valid", "per_execution", true},
		{"per_invocation valid", "per_invocation", true},
		{"per_api_call valid", "per_api_call", true},
		{"per_lookup valid", "per_lookup", true},
		{"per_query valid", "per_query", true},

		// Compute-based
		{"per_cpu_hour valid", "per_cpu_hour", true},
		{"per_cpu_month valid", "per_cpu_month", true},
		{"per_vcpu_hour valid", "per_vcpu_hour", true},
		{"per_memory_gb_hour valid", "per_memory_gb_hour", true},
		{"per_memory_gb_month valid", "per_memory_gb_month", true},

		// I/O-based
		{"per_iops valid", "per_iops", true},
		{"per_provisioned_iops valid", "per_provisioned_iops", true},
		{"per_data_transfer_gb valid", "per_data_transfer_gb", true},
		{"per_bandwidth_gb valid", "per_bandwidth_gb", true},

		// Database-specific
		{"per_rcu valid", "per_rcu", true},
		{"per_wcu valid", "per_wcu", true},
		{"per_dtu valid", "per_dtu", true},
		{"per_ru valid", "per_ru", true},

		// Pricing models
		{"on_demand valid", "on_demand", true},
		{"reserved valid", "reserved", true},
		{"spot valid", "spot", true},
		{"preemptible valid", "preemptible", true},
		{"savings_plan valid", "savings_plan", true},
		{"committed_use valid", "committed_use", true},
		{"hybrid_benefit valid", "hybrid_benefit", true},
		{"flat valid", "flat", true},

		// Invalid modes
		{"invalid mode", "invalid_mode", false},
		{"empty string", "", false},
		{"per hour with space", "per hour", false},
		{"case sensitive", "PER_HOUR", false},
		{"typo", "per_hor", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pricing.ValidBillingMode(tt.billingMode)
			if result != tt.expectValid {
				t.Errorf(
					"pricing.ValidBillingMode(%q) = %v, want %v",
					tt.billingMode,
					result,
					tt.expectValid,
				)
			}
		})
	}
}

func TestBillingModeString(t *testing.T) {
	tests := []struct {
		name         string
		billingMode  pricing.BillingMode
		expectString string
	}{
		{"PerHour string", pricing.PerHour, "per_hour"},
		{"PerGBMonth string", pricing.PerGBMonth, "per_gb_month"},
		{"PerRequest string", pricing.PerRequest, "per_request"},
		{"FlatRate string", pricing.FlatRate, "flat"},
		{"OnDemand string", pricing.OnDemand, "on_demand"},
		{"Reserved string", pricing.Reserved, "reserved"},
		{"Spot string", pricing.Spot, "spot"},
		{"PerRCU string", pricing.PerRCU, "per_rcu"},
		{"PerDTU string", pricing.PerDTU, "per_dtu"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.billingMode.String()
			if result != tt.expectString {
				t.Errorf("BillingMode.String() = %q, want %q", result, tt.expectString)
			}
		})
	}
}

func TestValidProvider(t *testing.T) {
	tests := []struct {
		name        string
		provider    string
		expectValid bool
	}{
		{"aws valid", "aws", true},
		{"azure valid", "azure", true},
		{"gcp valid", "gcp", true},
		{"kubernetes valid", "kubernetes", true},
		{"custom valid", "custom", true},

		{"invalid provider", "invalid_provider", false},
		{"empty string", "", false},
		{"case sensitive", "AWS", false},
		{"case sensitive azure", "Azure", false},
		{"case sensitive gcp", "GCP", false},
		{"typo", "azur", false},
		{"spaces", "g cp", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pricing.ValidProvider(tt.provider)
			if result != tt.expectValid {
				t.Errorf(
					"pricing.ValidProvider(%q) = %v, want %v",
					tt.provider,
					result,
					tt.expectValid,
				)
			}
		})
	}
}

func TestProviderString(t *testing.T) {
	tests := []struct {
		name         string
		provider     pricing.Provider
		expectString string
	}{
		{"AWS string", pricing.AWS, "aws"},
		{"Azure string", pricing.Azure, "azure"},
		{"GCP string", pricing.GCP, "gcp"},
		{"Kubernetes string", pricing.Kubernetes, "kubernetes"},
		{"Custom string", pricing.Custom, "custom"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.provider.String()
			if result != tt.expectString {
				t.Errorf("Provider.String() = %q, want %q", result, tt.expectString)
			}
		})
	}
}

func TestAllBillingModesCompleteness(t *testing.T) {
	// Test that GetAllBillingModes contains all expected billing modes
	allModesStr := pricing.GetAllBillingModes()
	allModes := make([]pricing.BillingMode, len(allModesStr))
	for i, mode := range allModesStr {
		allModes[i] = pricing.BillingMode(mode)
	}
	expectedCount := 40 // Based on the constants defined (added Tiered, NotImplemented)
	if len(allModes) != expectedCount {
		t.Errorf("GetAllBillingModes length = %d, want %d", len(allModes), expectedCount)
	}

	// Test that each mode in GetAllBillingModes validates as true
	for _, mode := range allModes {
		if !pricing.ValidBillingMode(mode.String()) {
			t.Errorf("GetAllBillingModes contains invalid mode: %q", mode)
		}
	}
}

func TestAllProvidersCompleteness(t *testing.T) {
	// Test that GetAllProviders contains all expected providers
	allProviders := pricing.GetAllProviders()
	expectedCount := 5 // aws, azure, gcp, kubernetes, custom
	if len(allProviders) != expectedCount {
		t.Errorf("GetAllProviders length = %d, want %d", len(allProviders), expectedCount)
	}

	// Test that each provider in GetAllProviders validates as true
	for _, provider := range allProviders {
		if !pricing.ValidProvider(provider.String()) {
			t.Errorf("GetAllProviders contains invalid provider: %q", provider)
		}
	}
}
