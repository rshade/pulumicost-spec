package pluginsdk

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"
)

// ValidateContext checks that a context is usable for RPC calls.
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
// or math.MaxInt64 duration if no deadline set.
func ContextRemainingTime(ctx context.Context) time.Duration {
	deadline, ok := ctx.Deadline()
	if !ok {
		return time.Duration(math.MaxInt64)
	}
	return time.Until(deadline)
}
