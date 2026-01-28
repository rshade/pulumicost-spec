package pluginsdk_test

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/rshade/finfocus-spec/sdk/go/pluginsdk"
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

func TestContextRemainingTime_UnexpiredDeadline(t *testing.T) {
	t.Parallel()

	// Create a context with a deadline 10 seconds in the future
	timeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Get the remaining time
	got := pluginsdk.ContextRemainingTime(ctx)

	// Verify remaining time is positive
	if got <= 0 {
		t.Errorf("expected positive duration for unexpired deadline, got %v", got)
	}

	// Verify remaining time is less than or equal to the original timeout
	// (allowing for slight timing variations)
	if got > timeout+time.Second {
		t.Errorf("duration %v should not exceed original timeout %v", got, timeout)
	}

	// Verify remaining time is reasonably close to the original timeout
	// (should be within 1 second of timeout since we just created the context)
	if got < timeout-time.Second {
		t.Errorf("expected duration close to %v, got %v (too much time passed)", timeout, got)
	}
}

func TestContextRemainingTime_ShortDeadline(t *testing.T) {
	t.Parallel()

	// Test with a very short deadline
	timeout := 100 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Get the remaining time immediately
	got := pluginsdk.ContextRemainingTime(ctx)

	// Allow small tolerance for timing variations (e.g., context creation overhead)
	tolerance := 10 * time.Millisecond

	// Should be positive and less than or equal to timeout (with tolerance)
	if got <= 0 {
		t.Errorf("expected positive duration, got %v", got)
	}
	if got > timeout+tolerance {
		t.Errorf("expected duration <= %v (with %v tolerance), got %v", timeout, tolerance, got)
	}

	// Wait deterministically for deadline to expire using context's Done channel
	<-ctx.Done()

	// After expiration, should return zero
	gotAfterExpiry := pluginsdk.ContextRemainingTime(ctx)
	if gotAfterExpiry != 0 {
		t.Errorf("expected zero duration after deadline expiry, got %v", gotAfterExpiry)
	}
}
