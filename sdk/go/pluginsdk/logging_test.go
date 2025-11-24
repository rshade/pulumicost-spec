package pluginsdk_test

import (
	"bytes"
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// TestNewPluginLogger_DefaultStderr tests NewPluginLogger with default stderr.
func TestNewPluginLogger_DefaultStderr(t *testing.T) {
	var buf bytes.Buffer
	logger := pluginsdk.NewPluginLogger("test-plugin", "v1.0.0", zerolog.InfoLevel, &buf)

	logger.Info().Msg("test message")

	output := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte(`"plugin_name":"test-plugin"`)) {
		t.Errorf("Expected plugin_name field, got: %s", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"plugin_version":"v1.0.0"`)) {
		t.Errorf("Expected plugin_version field, got: %s", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"message":"test message"`)) {
		t.Errorf("Expected message field, got: %s", output)
	}
}

// TestNewPluginLogger_CustomWriter tests NewPluginLogger with custom io.Writer.
func TestNewPluginLogger_CustomWriter(t *testing.T) {
	var buf bytes.Buffer
	logger := pluginsdk.NewPluginLogger("custom-plugin", "v2.0.0", zerolog.DebugLevel, &buf)

	logger.Debug().Msg("debug message")

	output := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte(`"plugin_name":"custom-plugin"`)) {
		t.Errorf("Expected plugin_name field, got: %s", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"level":"debug"`)) {
		t.Errorf("Expected debug level, got: %s", output)
	}
}

// TestNewPluginLogger_LevelFiltering tests log level filtering (Debug vs Info).
func TestNewPluginLogger_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := pluginsdk.NewPluginLogger("test-plugin", "v1.0.0", zerolog.InfoLevel, &buf)

	// Debug should be filtered out at Info level
	logger.Debug().Msg("debug message")
	if buf.Len() > 0 {
		t.Errorf("Debug message should be filtered at Info level, got: %s", buf.String())
	}

	// Info should pass through
	logger.Info().Msg("info message")
	if !bytes.Contains(buf.Bytes(), []byte("info message")) {
		t.Errorf("Info message should pass through, got: %s", buf.String())
	}
}

// TestNewPluginLogger_EmptyNameVersion tests empty plugin name/version handling.
func TestNewPluginLogger_EmptyNameVersion(t *testing.T) {
	var buf bytes.Buffer
	logger := pluginsdk.NewPluginLogger("", "", zerolog.InfoLevel, &buf)

	logger.Info().Msg("test")

	output := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte(`"plugin_name":""`)) {
		t.Errorf("Expected empty plugin_name field, got: %s", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"plugin_version":""`)) {
		t.Errorf("Expected empty plugin_version field, got: %s", output)
	}
}

// TestNewPluginLogger_NilWriter tests file output configuration with nil writer.
func TestNewPluginLogger_NilWriter(t *testing.T) {
	// When nil is passed, should default to os.Stderr
	// We can't easily test stderr output, but we can verify no panic
	logger := pluginsdk.NewPluginLogger("test", "v1.0.0", zerolog.InfoLevel, nil)
	// Just verify it doesn't panic
	logger.Info().Msg("test")
	t.Log("Logger with nil writer did not panic")
}

// TestContextWithTraceID_TraceIDFromContext tests ContextWithTraceID/TraceIDFromContext.
func TestContextWithTraceID_TraceIDFromContext(t *testing.T) {
	ctx := context.Background()
	traceID := "abc123"

	ctx = pluginsdk.ContextWithTraceID(ctx, traceID)
	retrieved := pluginsdk.TraceIDFromContext(ctx)

	if retrieved != traceID {
		t.Errorf("Expected trace ID %q, got %q", traceID, retrieved)
	}
}

// TestTraceIDFromContext_EmptyContext tests TraceIDFromContext with empty context.
func TestTraceIDFromContext_EmptyContext(t *testing.T) {
	ctx := context.Background()
	retrieved := pluginsdk.TraceIDFromContext(ctx)

	if retrieved != "" {
		t.Errorf("Expected empty string for context without trace ID, got %q", retrieved)
	}
}

// TestTracingUnaryServerInterceptor_Integration tests TracingUnaryServerInterceptor with bufconn.
func TestTracingUnaryServerInterceptor_Integration(t *testing.T) {
	// Test that interceptor extracts trace ID correctly
	interceptor := pluginsdk.TracingUnaryServerInterceptor()
	ctx := context.Background()
	md := metadata.New(map[string]string{
		pluginsdk.TraceIDMetadataKey: "test-trace-123",
	})
	ctx = metadata.NewIncomingContext(ctx, md)

	var capturedTraceID string
	handler := func(ctx context.Context, _ interface{}) (interface{}, error) {
		capturedTraceID = pluginsdk.TraceIDFromContext(ctx)
		return struct{}{}, nil
	}

	_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{}, handler)
	if err != nil {
		t.Fatalf("Interceptor failed: %v", err)
	}

	if capturedTraceID != "test-trace-123" {
		t.Errorf("Expected trace ID 'test-trace-123', got %q", capturedTraceID)
	}
}

// TestTracingUnaryServerInterceptor_MissingMetadata tests interceptor with missing metadata.
func TestTracingUnaryServerInterceptor_MissingMetadata(t *testing.T) {
	interceptor := pluginsdk.TracingUnaryServerInterceptor()
	ctx := context.Background()
	// No metadata set

	var capturedTraceID string
	handler := func(ctx context.Context, _ interface{}) (interface{}, error) {
		capturedTraceID = pluginsdk.TraceIDFromContext(ctx)
		return struct{}{}, nil
	}

	_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{}, handler)
	if err != nil {
		t.Fatalf("Interceptor failed: %v", err)
	}

	if capturedTraceID != "" {
		t.Errorf("Expected empty trace ID for missing metadata, got %q", capturedTraceID)
	}
}

// TestTracingUnaryServerInterceptor_ConcurrentRequests tests concurrent requests with different trace_ids.
func TestTracingUnaryServerInterceptor_ConcurrentRequests(t *testing.T) {
	interceptor := pluginsdk.TracingUnaryServerInterceptor()

	var wg sync.WaitGroup
	results := make(map[string]string)
	var mu sync.Mutex

	traceIDs := []string{"trace-1", "trace-2", "trace-3", "trace-4", "trace-5"}

	for _, traceID := range traceIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()

			ctx := context.Background()
			md := metadata.New(map[string]string{
				pluginsdk.TraceIDMetadataKey: id,
			})
			ctx = metadata.NewIncomingContext(ctx, md)

			handler := func(ctx context.Context, _ interface{}) (interface{}, error) {
				captured := pluginsdk.TraceIDFromContext(ctx)
				mu.Lock()
				results[id] = captured
				mu.Unlock()
				return struct{}{}, nil
			}

			_, _ = interceptor(ctx, nil, &grpc.UnaryServerInfo{}, handler)
		}(traceID)
	}

	wg.Wait()

	for _, expected := range traceIDs {
		if results[expected] != expected {
			t.Errorf("Trace ID mismatch: expected %q, got %q", expected, results[expected])
		}
	}
}

// TestTracingUnaryServerInterceptor_MultipleTraceIDs tests multiple trace_id values in metadata (use first).
func TestTracingUnaryServerInterceptor_MultipleTraceIDs(t *testing.T) {
	interceptor := pluginsdk.TracingUnaryServerInterceptor()
	ctx := context.Background()

	// Create metadata with multiple values for the same key
	md := metadata.Pairs(
		pluginsdk.TraceIDMetadataKey, "first-trace",
		pluginsdk.TraceIDMetadataKey, "second-trace",
	)
	ctx = metadata.NewIncomingContext(ctx, md)

	var capturedTraceID string
	handler := func(ctx context.Context, _ interface{}) (interface{}, error) {
		capturedTraceID = pluginsdk.TraceIDFromContext(ctx)
		return struct{}{}, nil
	}

	_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{}, handler)
	if err != nil {
		t.Fatalf("Interceptor failed: %v", err)
	}

	if capturedTraceID != "first-trace" {
		t.Errorf("Expected first trace ID 'first-trace', got %q", capturedTraceID)
	}
}

// TestLogOperation_TimingAccuracy tests LogOperation timing accuracy.
func TestLogOperation_TimingAccuracy(t *testing.T) {
	var buf bytes.Buffer
	logger := pluginsdk.NewPluginLogger("test", "v1.0.0", zerolog.InfoLevel, &buf)

	done := pluginsdk.LogOperation(logger, "test-operation")
	time.Sleep(50 * time.Millisecond)
	done()

	// Parse the JSON output
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	durationMs, ok := logEntry[pluginsdk.FieldDurationMs].(float64)
	if !ok {
		t.Fatalf("Expected duration_ms field, got: %v", logEntry)
	}

	// Allow some tolerance for timing
	if durationMs < 40 || durationMs > 150 {
		t.Errorf("Expected duration around 50ms, got %.0fms", durationMs)
	}
}

// TestLogOperation_OutputFormat tests LogOperation log output format.
func TestLogOperation_OutputFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := pluginsdk.NewPluginLogger("test", "v1.0.0", zerolog.InfoLevel, &buf)

	done := pluginsdk.LogOperation(logger, "GetProjectedCost")
	done()

	output := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte(`"operation":"GetProjectedCost"`)) {
		t.Errorf("Expected operation field, got: %s", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"duration_ms":`)) {
		t.Errorf("Expected duration_ms field, got: %s", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"message":"operation completed"`)) {
		t.Errorf("Expected message field, got: %s", output)
	}
}

// TestFieldConstants_Values tests all 11 field constants have correct string values.
func TestFieldConstants_Values(t *testing.T) {
	tests := []struct {
		constant string
		expected string
	}{
		{pluginsdk.FieldTraceID, "trace_id"},
		{pluginsdk.FieldComponent, "component"},
		{pluginsdk.FieldOperation, "operation"},
		{pluginsdk.FieldDurationMs, "duration_ms"},
		{pluginsdk.FieldResourceURN, "resource_urn"},
		{pluginsdk.FieldResourceType, "resource_type"},
		{pluginsdk.FieldPluginName, "plugin_name"},
		{pluginsdk.FieldPluginVersion, "plugin_version"},
		{pluginsdk.FieldCostMonthly, "cost_monthly"},
		{pluginsdk.FieldAdapter, "adapter"},
		{pluginsdk.FieldErrorCode, "error_code"},
	}

	for _, tt := range tests {
		if tt.constant != tt.expected {
			t.Errorf("Field constant mismatch: expected %q, got %q", tt.expected, tt.constant)
		}
	}
}

// TestTraceIDMetadataKey_Value tests TraceIDMetadataKey constant value.
func TestTraceIDMetadataKey_Value(t *testing.T) {
	expected := "x-pulumicost-trace-id"
	if pluginsdk.TraceIDMetadataKey != expected {
		t.Errorf("TraceIDMetadataKey mismatch: expected %q, got %q", expected, pluginsdk.TraceIDMetadataKey)
	}
}

// BenchmarkNewPluginLogger measures logger construction performance.
func BenchmarkNewPluginLogger(b *testing.B) {
	var buf bytes.Buffer
	b.ResetTimer()
	for range b.N {
		_ = pluginsdk.NewPluginLogger("bench-plugin", "v1.0.0", zerolog.InfoLevel, &buf)
		buf.Reset()
	}
}

// BenchmarkLogCall measures log call overhead.
func BenchmarkLogCall(b *testing.B) {
	var buf bytes.Buffer
	logger := pluginsdk.NewPluginLogger("bench-plugin", "v1.0.0", zerolog.InfoLevel, &buf)
	b.ResetTimer()
	for range b.N {
		logger.Info().Str("key", "value").Msg("benchmark message")
		buf.Reset()
	}
}

// BenchmarkInterceptor measures interceptor overhead.
func BenchmarkInterceptor(b *testing.B) {
	interceptor := pluginsdk.TracingUnaryServerInterceptor()
	ctx := context.Background()
	md := metadata.New(map[string]string{
		pluginsdk.TraceIDMetadataKey: "bench-trace-id",
	})
	ctx = metadata.NewIncomingContext(ctx, md)

	handler := func(ctx context.Context, _ interface{}) (interface{}, error) {
		_ = pluginsdk.TraceIDFromContext(ctx)
		return struct{}{}, nil
	}

	b.ResetTimer()
	for range b.N {
		_, _ = interceptor(ctx, nil, &grpc.UnaryServerInfo{}, handler)
	}
}
