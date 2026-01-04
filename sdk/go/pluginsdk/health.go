// Package pluginsdk provides a development SDK for PulumiCost plugins.
package pluginsdk

import (
	"net/http"
)

// HealthHandler returns an http.Handler that responds to health check requests.
// It returns 200 OK with "ok" body for GET requests to any path.
// This is designed for use with /healthz endpoints.
//
// The handler is intentionally simple - it does not check plugin status or
// dependencies. For more sophisticated health checks, use the gRPC health
// protocol (grpc.health.v1.Health/Check).
func HealthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusOK)
		if r.Method == http.MethodGet {
			_, _ = w.Write([]byte("ok"))
		}
	})
}
