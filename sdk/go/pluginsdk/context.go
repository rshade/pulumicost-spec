package pluginsdk

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"
)

// Sentinel errors for context validation.
var (
	// ErrNilContext is returned when a nil context is provided to ValidateContext.
	ErrNilContext = errors.New("context cannot be nil")

	// ErrContextCanceled is returned when the context has already been canceled or expired.
	ErrContextCanceled = errors.New("context already canceled or expired")
)

// ValidateContext checks that the provided context is usable for RPC calls.
// It returns ErrNilContext if ctx is nil, or wraps ErrContextCanceled with the
// underlying context error if the context has already been canceled or its deadline exceeded;
// otherwise it returns nil.
func ValidateContext(ctx context.Context) error {
	if ctx == nil {
		return ErrNilContext
	}
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("%w: %w", ErrContextCanceled, err)
	}
	return nil
}

// ContextRemainingTime returns the duration until the context's deadline.
// It returns 0 if the context is nil, already canceled, or expired.
// It returns math.MaxInt64 if the context has no deadline.
// Otherwise, it returns the time remaining until the deadline (which may be negative
// if the deadline has just passed but the context error hasn't propagated yet).
func ContextRemainingTime(ctx context.Context) time.Duration {
	if ctx == nil {
		return 0
	}
	// Return 0 for already canceled/expired contexts
	if ctx.Err() != nil {
		return 0
	}
	deadline, ok := ctx.Deadline()
	if !ok {
		return time.Duration(math.MaxInt64)
	}
	return time.Until(deadline)
}
