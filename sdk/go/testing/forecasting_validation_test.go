package testing_test

import (
	"errors"
	"testing"

	"github.com/rshade/pulumicost-spec/sdk/go/pricing"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// TestValidation_LinearWithoutRate validates that LINEAR without growth_rate returns error.
func TestValidation_LinearWithoutRate(t *testing.T) {
	err := pricing.ValidateGrowthParams(pbc.GrowthType_GROWTH_TYPE_LINEAR, nil)

	if err == nil {
		t.Error("Expected error for LINEAR without growth_rate, got nil")
	}

	if !errors.Is(err, pricing.ErrMissingGrowthRate) {
		t.Errorf("Expected ErrMissingGrowthRate, got %v", err)
	}
}

// TestValidation_ExponentialWithoutRate validates that EXPONENTIAL without growth_rate returns error.
func TestValidation_ExponentialWithoutRate(t *testing.T) {
	err := pricing.ValidateGrowthParams(pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, nil)

	if err == nil {
		t.Error("Expected error for EXPONENTIAL without growth_rate, got nil")
	}

	if !errors.Is(err, pricing.ErrMissingGrowthRate) {
		t.Errorf("Expected ErrMissingGrowthRate, got %v", err)
	}
}

// TestValidation_RateBelowNegativeOne validates that rate < -1.0 returns error.
func TestValidation_RateBelowNegativeOne(t *testing.T) {
	invalidRate := floatPtr(-1.5) // -150% decline is invalid

	err := pricing.ValidateGrowthParams(pbc.GrowthType_GROWTH_TYPE_LINEAR, invalidRate)

	if err == nil {
		t.Error("Expected error for growth_rate < -1.0, got nil")
	}

	if !errors.Is(err, pricing.ErrInvalidGrowthRate) {
		t.Errorf("Expected ErrInvalidGrowthRate, got %v", err)
	}
}

// TestValidation_NegativeRateAccepted validates that rate >= -1.0 is accepted.
func TestValidation_NegativeRateAccepted(t *testing.T) {
	tests := []struct {
		name string
		rate float64
	}{
		{"minus one (boundary)", -1.0},
		{"minus ninety percent", -0.9},
		{"minus ten percent", -0.10},
		{"zero", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pricing.ValidateGrowthParams(pbc.GrowthType_GROWTH_TYPE_LINEAR, floatPtr(tt.rate))
			if err != nil {
				t.Errorf("Expected no error for rate %v, got %v", tt.rate, err)
			}
		})
	}
}

// TestValidation_HighRateAccepted validates that very high rates are accepted.
func TestValidation_HighRateAccepted(t *testing.T) {
	tests := []struct {
		name string
		rate float64
	}{
		{"one hundred percent", 1.0},
		{"two hundred percent", 2.0},
		{"one thousand percent", 10.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pricing.ValidateGrowthParams(pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, floatPtr(tt.rate))
			if err != nil {
				t.Errorf("Expected no error for rate %v, got %v", tt.rate, err)
			}
		})
	}
}

// TestValidation_NoneWithoutRate validates that NONE without rate is valid.
func TestValidation_NoneWithoutRate(t *testing.T) {
	err := pricing.ValidateGrowthParams(pbc.GrowthType_GROWTH_TYPE_NONE, nil)
	if err != nil {
		t.Errorf("Expected no error for NONE without rate, got %v", err)
	}
}

// TestValidation_NoneWithRate validates that NONE with rate is valid (rate is ignored).
func TestValidation_NoneWithRate(t *testing.T) {
	err := pricing.ValidateGrowthParams(pbc.GrowthType_GROWTH_TYPE_NONE, floatPtr(0.10))
	if err != nil {
		t.Errorf("Expected no error for NONE with rate (ignored), got %v", err)
	}
}

// TestValidation_UnspecifiedWithoutRate validates that UNSPECIFIED without rate is valid.
func TestValidation_UnspecifiedWithoutRate(t *testing.T) {
	err := pricing.ValidateGrowthParams(pbc.GrowthType_GROWTH_TYPE_UNSPECIFIED, nil)
	if err != nil {
		t.Errorf("Expected no error for UNSPECIFIED without rate, got %v", err)
	}
}

// TestValidation_UnspecifiedWithRate validates that UNSPECIFIED with rate is valid.
func TestValidation_UnspecifiedWithRate(t *testing.T) {
	err := pricing.ValidateGrowthParams(pbc.GrowthType_GROWTH_TYPE_UNSPECIFIED, floatPtr(0.10))
	if err != nil {
		t.Errorf("Expected no error for UNSPECIFIED with rate (ignored), got %v", err)
	}
}

// TestValidation_ErrorMessages validates error messages are descriptive.
func TestValidation_ErrorMessages(t *testing.T) {
	t.Run("MissingRateMessage", func(t *testing.T) {
		err := pricing.ValidateGrowthParams(pbc.GrowthType_GROWTH_TYPE_LINEAR, nil)
		if err == nil {
			t.Fatal("Expected error")
		}
		if err.Error() == "" {
			t.Error("Expected non-empty error message")
		}
		t.Logf("Error message: %s", err.Error())
	})

	t.Run("InvalidRateMessage", func(t *testing.T) {
		err := pricing.ValidateGrowthParams(pbc.GrowthType_GROWTH_TYPE_LINEAR, floatPtr(-2.0))
		if err == nil {
			t.Fatal("Expected error")
		}
		if err.Error() == "" {
			t.Error("Expected non-empty error message")
		}
		t.Logf("Error message: %s", err.Error())
	})
}

// TestValidation_ComprehensiveMatrix tests all growth type and rate combinations.
func TestValidation_ComprehensiveMatrix(t *testing.T) {
	tests := []struct {
		name        string
		growthType  pbc.GrowthType
		rate        *float64
		expectError bool
		errorType   error
	}{
		// UNSPECIFIED - always valid
		{"UNSPECIFIED nil", pbc.GrowthType_GROWTH_TYPE_UNSPECIFIED, nil, false, nil},
		{"UNSPECIFIED 0.0", pbc.GrowthType_GROWTH_TYPE_UNSPECIFIED, floatPtr(0.0), false, nil},
		{"UNSPECIFIED 0.10", pbc.GrowthType_GROWTH_TYPE_UNSPECIFIED, floatPtr(0.10), false, nil},

		// NONE - always valid
		{"NONE nil", pbc.GrowthType_GROWTH_TYPE_NONE, nil, false, nil},
		{"NONE 0.0", pbc.GrowthType_GROWTH_TYPE_NONE, floatPtr(0.0), false, nil},
		{"NONE 0.10", pbc.GrowthType_GROWTH_TYPE_NONE, floatPtr(0.10), false, nil},

		// LINEAR - requires rate
		{"LINEAR nil", pbc.GrowthType_GROWTH_TYPE_LINEAR, nil, true, pricing.ErrMissingGrowthRate},
		{"LINEAR 0.0", pbc.GrowthType_GROWTH_TYPE_LINEAR, floatPtr(0.0), false, nil},
		{"LINEAR 0.10", pbc.GrowthType_GROWTH_TYPE_LINEAR, floatPtr(0.10), false, nil},
		{"LINEAR -1.0", pbc.GrowthType_GROWTH_TYPE_LINEAR, floatPtr(-1.0), false, nil},
		{"LINEAR -1.5", pbc.GrowthType_GROWTH_TYPE_LINEAR, floatPtr(-1.5), true, pricing.ErrInvalidGrowthRate},

		// EXPONENTIAL - requires rate
		{"EXPONENTIAL nil", pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, nil, true, pricing.ErrMissingGrowthRate},
		{"EXPONENTIAL 0.0", pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, floatPtr(0.0), false, nil},
		{"EXPONENTIAL 0.10", pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, floatPtr(0.10), false, nil},
		{"EXPONENTIAL -1.0", pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, floatPtr(-1.0), false, nil},
		{
			"EXPONENTIAL -1.5",
			pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
			floatPtr(-1.5),
			true,
			pricing.ErrInvalidGrowthRate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pricing.ValidateGrowthParams(tt.growthType, tt.rate)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				if tt.errorType != nil && !errors.Is(err, tt.errorType) {
					t.Errorf("Expected error type %v, got %v", tt.errorType, err)
				}
			} else if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}
