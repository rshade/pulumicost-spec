package pluginsdk_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
)

type testHealthChecker struct {
	err          error
	panicOnCheck bool
}

func (h *testHealthChecker) Check(_ context.Context) error {
	if h.panicOnCheck {
		panic("health check panic")
	}
	if h.err != nil {
		return h.err
	}
	return nil
}

func TestHealthChecker(t *testing.T) {
	t.Run("healthy checker returns nil", func(t *testing.T) {
		checker := &testHealthChecker{}
		ctx := context.Background()
		err := checker.Check(ctx)
		if err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("unhealthy checker returns error", func(t *testing.T) {
		checker := &testHealthChecker{err: errors.New("database unavailable")}
		ctx := context.Background()
		err := checker.Check(ctx)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if err.Error() != "database unavailable" {
			t.Errorf("error message mismatch: got %q, want %q", err.Error(), "database unavailable")
		}
	})

	t.Run("checker that panics is recovered", func(t *testing.T) {
		checker := &testHealthChecker{panicOnCheck: true}
		ctx := context.Background()
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic but did not panic")
			}
		}()
		err := checker.Check(ctx)
		t.Logf("err after panic: %v", err)
	})
}

func runHealthTest(
	t *testing.T,
	checker pluginsdk.HealthChecker,
	expectedCode int,
	expectedHealthy bool,
	expectedMsg string,
) {
	handler := pluginsdk.HealthHandler(checker)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != expectedCode {
		t.Errorf("expected %d, got %d", expectedCode, rec.Code)
	}

	var status pluginsdk.HealthStatus
	if err := json.NewDecoder(rec.Body).Decode(&status); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if status.Healthy != expectedHealthy {
		t.Errorf("expected healthy=%v, got %v", expectedHealthy, status.Healthy)
	}

	if expectedMsg != "" && status.Message != expectedMsg {
		t.Errorf("expected message %q, got %q", expectedMsg, status.Message)
	}
}

func TestHealthHandler_Integration(t *testing.T) {
	t.Run("healthy returns 200 with JSON", func(t *testing.T) {
		checker := &testHealthChecker{err: nil}
		runHealthTest(t, checker, http.StatusOK, true, "")
	})

	t.Run("unhealthy returns 503 with error message", func(t *testing.T) {
		checker := &testHealthChecker{err: errors.New("db down")}
		runHealthTest(t, checker, http.StatusServiceUnavailable, false, "db down")
	})

	t.Run("panic returns 503 with error message", func(t *testing.T) {
		checker := &testHealthChecker{panicOnCheck: true}
		runHealthTest(t, checker, http.StatusServiceUnavailable, false, "panic during health check: health check panic")
	})
}
