package pluginsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
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
	Healthy     bool              `json:"healthy"`
	Message     string            `json:"message,omitempty"`
	Details     map[string]string `json:"details,omitempty"`
	LastChecked time.Time         `json:"last_checked"`
}

// executeCheck runs the health check safely, recovering from panics.
func executeCheck(ctx context.Context, checker HealthChecker) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("panic during health check: %v", rec)
		}
	}()
	return checker.Check(ctx)
}

// HealthHandler returns an http.Handler that responds to health check requests.
// If checker is provided, it runs the check and returns a JSON HealthStatus.
// If checker is nil, it returns 200 OK with "ok" body (legacy behavior).
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
				// Log encoding error to standard output since we don't have a logger context here
				_, _ = fmt.Fprintf(os.Stderr, "pluginsdk: failed to encode health status: %v\n", encodeErr)
			}
		}
	})
}
