package pluginsdk_test

import (
	"testing"

	"github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
	"github.com/stretchr/testify/assert"
)

func TestValidationMode_Constants(t *testing.T) {
	// Verify FailFast is the zero value (default)
	assert.Equal(t, pluginsdk.ValidationModeFailFast, pluginsdk.ValidationMode(0))

	// Verify Aggregate is distinct
	assert.Equal(t, pluginsdk.ValidationModeAggregate, pluginsdk.ValidationMode(1))

	// Verify they're not equal
	assert.NotEqual(t, pluginsdk.ValidationModeFailFast, pluginsdk.ValidationModeAggregate)
}

func TestValidationOptions_DefaultMode(t *testing.T) {
	// Default options should use FailFast mode
	opts := pluginsdk.ValidationOptions{}
	assert.Equal(t, pluginsdk.ValidationModeFailFast, opts.Mode)
}

func TestValidationOptions_AggregateMode(t *testing.T) {
	opts := pluginsdk.ValidationOptions{
		Mode: pluginsdk.ValidationModeAggregate,
	}
	assert.Equal(t, pluginsdk.ValidationModeAggregate, opts.Mode)
}

func TestValidationOptions_ExplicitFailFast(t *testing.T) {
	opts := pluginsdk.ValidationOptions{
		Mode: pluginsdk.ValidationModeFailFast,
	}
	assert.Equal(t, pluginsdk.ValidationModeFailFast, opts.Mode)
}

func TestValidationMode_String(t *testing.T) {
	tests := []struct {
		mode     pluginsdk.ValidationMode
		expected string
	}{
		{pluginsdk.ValidationModeFailFast, "fail_fast"},
		{pluginsdk.ValidationModeAggregate, "aggregate"},
		{pluginsdk.ValidationMode(99), "unknown"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.mode.String())
	}
}
