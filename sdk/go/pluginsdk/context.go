package pluginsdk

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"
)

// ValidateContext checks that the provided context is usable for RPC calls.
// It returns an error if ctx is nil or if the context has already been cancelled or its deadline exceeded; otherwise it returns nil.
func ValidateContext(ctx context.Context) error {
	if ctx == nil {
		return errors.New("context cannot be nil")
	}
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context already cancelled or expired: %w", err)
	}
	return nil
}

// ContextRemainingTime returns time until deadline, negative if expired,
// ContextRemainingTime reports the duration until ctx's deadline.
// If ctx has no deadline it returns a duration equal to math.MaxInt64; if the deadline has passed the result is negative.
func ContextRemainingTime(ctx context.Context) time.Duration {
	deadline, ok := ctx.Deadline()
	if !ok {
		return time.Duration(math.MaxInt64)
	}
	return time.Until(deadline)
}