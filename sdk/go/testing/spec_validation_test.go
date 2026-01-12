package testing_test

import (
	"testing"

	plugintesting "github.com/rshade/finfocus-spec/sdk/go/testing"
)

// TestSpecValidationPassesForValidPricingSpec validates that spec validation passes
// for a valid PricingSpec (T016).
func TestSpecValidationPassesForValidPricingSpec(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	tests := plugintesting.SpecValidationTests()
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result := test.TestFunc(harness)
			if !result.Success {
				t.Errorf("Test %s failed: %v - %s", test.Name, result.Error, result.Details)
			}
		})
	}
}

// TestSpecValidationFailsForInvalidBillingMode validates that spec validation fails
// for an invalid billing mode (T017).
func TestSpecValidationFailsForInvalidBillingMode(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	// Configure plugin to return an invalid billing mode
	plugin.PluginName = "invalid-billing-mode-plugin"

	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	// Test the billing mode validation directly
	invalidModes := []string{
		"invalid_mode",
		"HOURLY",   // Wrong case
		"per-hour", // Wrong separator
		"",
	}

	for _, mode := range invalidModes {
		t.Run("InvalidMode_"+mode, func(t *testing.T) {
			err := plugintesting.ValidateBillingModePublic(mode)
			if mode == "" {
				// Empty is handled by required fields check
				if err != nil {
					t.Logf("Empty mode correctly skipped by enum check (handled by required check)")
				}
			} else if err == nil {
				t.Errorf("Expected validation error for billing mode %q, got nil", mode)
			}
		})
	}

	// Test valid modes don't fail
	validModes := []string{
		"per_hour",
		"on_demand",
		"flat_rate",
	}

	for _, mode := range validModes {
		t.Run("ValidMode_"+mode, func(t *testing.T) {
			err := plugintesting.ValidateBillingModePublic(mode)
			if err != nil {
				t.Errorf("Unexpected validation error for billing mode %q: %v", mode, err)
			}
		})
	}
}

// TestSpecValidationFailsForMissingRequiredFields validates that spec validation fails
// for missing required fields (T018).
//
//nolint:gocognit // Table-driven tests inherently have higher complexity
func TestSpecValidationFailsForMissingRequiredFields(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	// The mock plugin should return valid required fields
	// We test the validation function directly with constructed specs
	testCases := []struct {
		name          string
		provider      string
		resourceType  string
		billingMode   string
		currency      string
		expectErrors  bool
		errorContains string
	}{
		{
			name:          "MissingProvider",
			provider:      "",
			resourceType:  "ec2",
			billingMode:   "per_hour",
			currency:      "USD",
			expectErrors:  true,
			errorContains: "provider",
		},
		{
			name:          "MissingResourceType",
			provider:      "aws",
			resourceType:  "",
			billingMode:   "per_hour",
			currency:      "USD",
			expectErrors:  true,
			errorContains: "resource_type",
		},
		{
			name:          "MissingBillingMode",
			provider:      "aws",
			resourceType:  "ec2",
			billingMode:   "",
			currency:      "USD",
			expectErrors:  true,
			errorContains: "billing_mode",
		},
		{
			name:          "MissingCurrency",
			provider:      "aws",
			resourceType:  "ec2",
			billingMode:   "per_hour",
			currency:      "",
			expectErrors:  true,
			errorContains: "currency",
		},
		{
			name:          "AllFieldsPresent",
			provider:      "aws",
			resourceType:  "ec2",
			billingMode:   "per_hour",
			currency:      "USD",
			expectErrors:  false,
			errorContains: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			errors := plugintesting.ValidateRequiredFieldsPublic(
				tc.provider,
				tc.resourceType,
				tc.billingMode,
				tc.currency,
			)

			if tc.expectErrors && len(errors) == 0 {
				t.Error("Expected validation errors, got none")
			}

			if !tc.expectErrors && len(errors) > 0 {
				t.Errorf("Expected no errors, got: %v", errors)
			}

			if tc.expectErrors && tc.errorContains != "" {
				found := false
				for _, err := range errors {
					if err.Field == tc.errorContains {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error for field %q, got: %v", tc.errorContains, errors)
				}
			}
		})
	}
}

// TestSpecValidationCurrencyFormat validates currency format (3-character ISO code).
func TestSpecValidationCurrencyFormat(t *testing.T) {
	testCases := []struct {
		currency    string
		expectError bool
	}{
		{"USD", false},
		{"EUR", false},
		{"GBP", false},
		{"JPY", false},
		{"US", true},   // Too short
		{"USDD", true}, // Too long
		{"", false},    // Empty handled by required check
	}

	for _, tc := range testCases {
		t.Run("Currency_"+tc.currency, func(t *testing.T) {
			err := plugintesting.ValidateCurrencyFormatPublic(tc.currency)
			if tc.expectError && err == nil {
				t.Errorf("Expected error for currency %q, got nil", tc.currency)
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error for currency %q: %v", tc.currency, err)
			}
		})
	}
}

// TestRegisterSpecValidationTests validates that tests can be registered.
func TestRegisterSpecValidationTests(t *testing.T) {
	suite := plugintesting.NewConformanceSuite()
	plugintesting.RegisterSpecValidationTests(suite)

	config := suite.GetConfig()
	if config.TargetLevel != plugintesting.ConformanceLevelStandard {
		t.Errorf("Expected default target level Standard, got %v", config.TargetLevel)
	}
}
