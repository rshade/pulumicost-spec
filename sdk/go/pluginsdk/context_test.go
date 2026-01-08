package pluginsdk_test

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
)

func TestValidateContext(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		expectError bool
	}{
		{
			name:        "nil context returns error",
			ctx:         nil,
			expectError: true,
		},
		{
			name: "cancelled context returns error",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			expectError: true,
		},
		{
			name: "expired context returns error",
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), -1*time.Hour)
				cancel()
				return ctx
			}(),
			expectError: true,
		},
		{
			name:        "valid context returns nil error",
			ctx:         context.Background(),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pluginsdk.ValidateContext(tt.ctx)
			if tt.expectError && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}

func TestContextRemainingTime(t *testing.T) {
	cancelledCtx, cancel := context.WithTimeout(context.Background(), time.Hour)
	cancel()

	tests := []struct {
		name             string
		ctx              context.Context
		expectedDuration time.Duration
	}{
		{
			name:             "nil context returns zero",
			ctx:              nil,
			expectedDuration: 0,
		},
		{
			name:             "cancelled context returns zero",
			ctx:              cancelledCtx,
			expectedDuration: 0,
		},
		{
			name:             "context without deadline returns MaxInt64 duration",
			ctx:              context.Background(),
			expectedDuration: time.Duration(math.MaxInt64),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pluginsdk.ContextRemainingTime(tt.ctx)
			gotSeconds := got.Seconds()
			expectedSeconds := tt.expectedDuration.Seconds()
			if math.Abs(gotSeconds-expectedSeconds) > 1.0 {
				t.Errorf("duration mismatch: got %v, expected %v", got, tt.expectedDuration)
			}
		})
	}
}
