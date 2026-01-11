package pluginsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog/log"
)

const healthCheckTimeout = 5 * time.Second

// HealthChecker allows plugins to provide custom health check logic.
type HealthChecker interface {
	// Check performs a health check and returns nil if healthy,
	// or an error describing the health issue.
	Check(ctx context.Context) error
}

// HealthStatus provides detailed health information.
type HealthStatus struct {
	Healthy bool   `json:"healthy"`
	Message string `json:"message,omitempty"`
	// Details is an optional map meant to be populated by custom HealthChecker implementations
	// to provide component-specific key/value diagnostics. The default HealthHandler does not populate it.
	Details     map[string]string `json:"details,omitempty"`
	LastChecked time.Time         `json:"last_checked"`
}

// executeCheck runs the given HealthChecker's Check method and converts panics into an error.
// If the check panics, executeCheck recovers and returns an error describing the panic; otherwise
// it returns the error returned by Check.
func executeCheck(ctx context.Context, checker HealthChecker) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Debug().
				Bytes("stack", debug.Stack()).
				Interface("panic", rec).
				Msg("health check panic")
			err = fmt.Errorf("panic during health check: %v", rec)
		}
	}()
	return checker.Check(ctx)
}

// HealthHandler returns an http.Handler that serves plugin health checks using the given HealthChecker.
//
// HealthHandler accepts only GET and HEAD requests. If checker is nil, it preserves legacy behavior by
// responding with 200 OK and a plain "ok" body for GET requests (text/plain; charset=utf-8).
//
// When checker is non-nil, the handler runs checker.Check with the request context (applying a 5s timeout
// if the context has no deadline), and returns a JSON-encoded HealthStatus on GET requests. The handler
// sets X-Content-Type-Options: nosniff and Content-Type appropriately. It responds with HTTP 200 when the
// health check reports healthy and HTTP 503 when it reports unhealthy. Any error message from the check is
// included in the HealthStatus.Message field.
func HealthHandler(checker HealthChecker) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("X-Content-Type-Options", "nosniff")

		if checker == nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			if r.Method == http.MethodGet {
				_, _ = w.Write([]byte("ok"))
			}
			return
		}

		ctx := r.Context()
		// Ensure a timeout if not already present
		if _, ok := ctx.Deadline(); !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, healthCheckTimeout)
			defer cancel()
		}

		err := executeCheck(ctx, checker)

		status := HealthStatus{
			Healthy:     err == nil,
			LastChecked: time.Now(),
		}
		if err != nil {
			status.Message = err.Error()
		}

		w.Header().Set("Content-Type", "application/json")
		if status.Healthy {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		if r.Method == http.MethodGet {
			if encodeErr := json.NewEncoder(w).Encode(status); encodeErr != nil {
				// Log encoding error
				log.Error().Err(encodeErr).Msg("pluginsdk: failed to encode health status")
			}
		}
	})
}
