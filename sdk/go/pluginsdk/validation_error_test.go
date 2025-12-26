package pluginsdk_test

import (
	"testing"

	"github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *pluginsdk.ValidationError
		expected string
	}{
		{
			name: "cost hierarchy violation",
			err: &pluginsdk.ValidationError{
				FieldName:     "effective_cost",
				Constraint:    "must not exceed billed_cost",
				ActualValue:   "150.00",
				ExpectedValue: "<= 100.00",
			},
			expected: "effective_cost: must not exceed billed_cost (actual: 150.00, expected: <= 100.00)",
		},
		{
			name: "missing required field",
			err: &pluginsdk.ValidationError{
				FieldName:     "commitment_discount_status",
				Constraint:    "required when commitment_discount_id is set for usage charges",
				ActualValue:   "UNSPECIFIED",
				ExpectedValue: "USED or UNUSED",
			},
			expected: "commitment_discount_status: required when commitment_discount_id is set for usage charges (actual: UNSPECIFIED, expected: USED or UNUSED)",
		},
		{
			name: "empty values",
			err: &pluginsdk.ValidationError{
				FieldName:     "pricing_unit",
				Constraint:    "required when pricing_quantity > 0",
				ActualValue:   "",
				ExpectedValue: "non-empty string",
			},
			expected: "pricing_unit: required when pricing_quantity > 0 (actual: , expected: non-empty string)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestValidationError_ImplementsError(t *testing.T) {
	var err error = &pluginsdk.ValidationError{
		FieldName:     "test_field",
		Constraint:    "test constraint",
		ActualValue:   "actual",
		ExpectedValue: "expected",
	}

	// Verify it implements error interface
	require.Error(t, err)
	assert.Contains(t, err.Error(), "test_field")
}

func TestValidationError_ErrorsAs(t *testing.T) {
	// Create an error chain
	originalErr := &pluginsdk.ValidationError{
		FieldName:     "effective_cost",
		Constraint:    "must not exceed billed_cost",
		ActualValue:   "200.00",
		ExpectedValue: "<= 100.00",
	}

	// errors.As should work for extracting ValidationError
	var valErr *pluginsdk.ValidationError
	require.ErrorAs(t, originalErr, &valErr)
	assert.Equal(t, "effective_cost", valErr.FieldName)
	assert.Equal(t, "must not exceed billed_cost", valErr.Constraint)
	assert.Equal(t, "200.00", valErr.ActualValue)
	assert.Equal(t, "<= 100.00", valErr.ExpectedValue)
}

func TestValidationError_FieldAccess(t *testing.T) {
	err := &pluginsdk.ValidationError{
		FieldName:     "list_cost",
		Constraint:    "must be >= effective_cost",
		ActualValue:   "50.00",
		ExpectedValue: ">= 100.00",
	}

	// Verify all fields are accessible
	assert.Equal(t, "list_cost", err.FieldName)
	assert.Equal(t, "must be >= effective_cost", err.Constraint)
	assert.Equal(t, "50.00", err.ActualValue)
	assert.Equal(t, ">= 100.00", err.ExpectedValue)
}

func TestNewValidationError(t *testing.T) {
	err := pluginsdk.NewValidationError(
		"effective_cost",
		"must not exceed billed_cost",
		"150.00",
		"<= 100.00",
	)

	require.NotNil(t, err)
	assert.Equal(t, "effective_cost", err.FieldName)
	assert.Equal(t, "must not exceed billed_cost", err.Constraint)
	assert.Equal(t, "150.00", err.ActualValue)
	assert.Equal(t, "<= 100.00", err.ExpectedValue)

	// Verify error message format matches direct construction
	expected := "effective_cost: must not exceed billed_cost (actual: 150.00, expected: <= 100.00)"
	assert.Equal(t, expected, err.Error())
}
