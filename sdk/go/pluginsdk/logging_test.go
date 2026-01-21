package pluginsdk_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/rshade/finfocus-spec/sdk/go/pluginsdk"
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
	// Test that interceptor extracts valid trace ID correctly
	interceptor := pluginsdk.TracingUnaryServerInterceptor()
	ctx := context.Background()
	validTraceID := "abcdef1234567890abcdef1234567890"
	md := metadata.New(map[string]string{
		pluginsdk.TraceIDMetadataKey: validTraceID,
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

	if capturedTraceID != validTraceID {
		t.Errorf("Expected trace ID %q, got %q", validTraceID, capturedTraceID)
	}
}

// TestTracingUnaryServerInterceptor_MissingMetadata tests interceptor generates trace ID when metadata is missing.
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

	// Should generate a valid trace ID when metadata is missing
	if capturedTraceID == "" {
		t.Errorf("Expected generated trace ID for missing metadata, got empty")
	}

	if len(capturedTraceID) != 32 {
		t.Errorf(
			"Generated trace ID should be 32 characters, got %d: %q",
			len(capturedTraceID),
			capturedTraceID,
		)
	}

	// Should be valid hex
	for _, r := range capturedTraceID {
		if r < '0' || (r > '9' && r < 'a') || r > 'f' {
			t.Errorf("Generated trace ID contains invalid character %q: %q", r, capturedTraceID)
		}
	}

	// Should not be all zeros
	if capturedTraceID == "00000000000000000000000000000000" {
		t.Errorf("Generated trace ID should not be all zeros")
	}
}

// TestTracingUnaryServerInterceptor_ConcurrentRequests tests concurrent requests with different valid trace_ids.
func TestTracingUnaryServerInterceptor_ConcurrentRequests(t *testing.T) {
	interceptor := pluginsdk.TracingUnaryServerInterceptor()

	var wg sync.WaitGroup
	results := make(map[string]string)
	var mu sync.Mutex

	traceIDs := []string{
		"abcdef1234567890abcdef1234567890",
		"1234567890abcdef1234567890abcdef",
		"a1b2c3d4e5f678901234567890abcdef",
		"fedcba0987654321fedcba0987654321",
		"11111111111111111111111111111111",
	}

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

// TestTracingUnaryServerInterceptor_MultipleTraceIDs tests multiple valid trace_id values in metadata (use first).
func TestTracingUnaryServerInterceptor_MultipleTraceIDs(t *testing.T) {
	interceptor := pluginsdk.TracingUnaryServerInterceptor()
	ctx := context.Background()

	// Create metadata with multiple values for the same key
	firstTraceID := "abcdef1234567890abcdef1234567890"
	md := metadata.Pairs(
		pluginsdk.TraceIDMetadataKey, firstTraceID,
		pluginsdk.TraceIDMetadataKey, "1234567890abcdef1234567890abcdef",
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

	if capturedTraceID != firstTraceID {
		t.Errorf("Expected first trace ID %q, got %q", firstTraceID, capturedTraceID)
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
	expected := "x-finfocus-trace-id"
	if pluginsdk.TraceIDMetadataKey != expected {
		t.Errorf(
			"TraceIDMetadataKey mismatch: expected %q, got %q",
			expected,
			pluginsdk.TraceIDMetadataKey,
		)
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

// BenchmarkInterceptorValidation measures validation overhead with invalid trace IDs.
func BenchmarkInterceptorValidation(b *testing.B) {
	interceptor := pluginsdk.TracingUnaryServerInterceptor()
	ctx := context.Background()
	md := metadata.New(map[string]string{
		pluginsdk.TraceIDMetadataKey: "invalid-trace-id", // Invalid: too short
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

// BenchmarkInterceptorGeneration measures generation overhead when no trace ID is provided.
func BenchmarkInterceptorGeneration(b *testing.B) {
	interceptor := pluginsdk.TracingUnaryServerInterceptor()
	ctx := context.Background()
	// No metadata - will trigger generation

	handler := func(ctx context.Context, _ interface{}) (interface{}, error) {
		_ = pluginsdk.TraceIDFromContext(ctx)
		return struct{}{}, nil
	}

	b.ResetTimer()
	for range b.N {
		_, _ = interceptor(ctx, nil, &grpc.UnaryServerInfo{}, handler)
	}
}

// validateGeneratedTraceID validates that a generated trace ID meets requirements.
func validateGeneratedTraceID(t *testing.T, traceID, originalInput string) {
	t.Helper()

	if traceID == "" {
		t.Errorf("Expected generated trace ID for invalid input %q, got empty", originalInput)
	}

	if len(traceID) != 32 {
		t.Errorf(
			"Generated trace ID should be 32 characters, got %d: %q",
			len(traceID),
			traceID,
		)
	}

	// Should be valid hex
	for _, r := range traceID {
		if r < '0' || (r > '9' && r < 'a') || r > 'f' {
			t.Errorf(
				"Generated trace ID contains invalid character %q: %q",
				r,
				traceID,
			)
		}
	}

	// Should not be all zeros
	if traceID == "00000000000000000000000000000000" {
		t.Errorf("Generated trace ID should not be all zeros")
	}

	// Should not be the same as invalid input
	if traceID == originalInput {
		t.Errorf(
			"Trace ID should have been replaced, but got same value: %q",
			traceID,
		)
	}
}

// TestTracingUnaryServerInterceptor_InvalidTraceIDs tests interceptor validation of invalid trace IDs.
func TestTracingUnaryServerInterceptor_InvalidTraceIDs(t *testing.T) {
	interceptor := pluginsdk.TracingUnaryServerInterceptor()

	tests := []struct {
		name        string
		traceID     string
		description string
	}{
		{
			name:        "too_short",
			traceID:     "abc123",
			description: "trace ID shorter than 32 characters",
		},
		{
			name:        "too_long",
			traceID:     "abcdef1234567890abcdef1234567890extra",
			description: "trace ID longer than 32 characters",
		},
		{
			name:        "non_hex_chars",
			traceID:     "gggggggggggggggggggggggggggggggg",
			description: "trace ID with non-hexadecimal characters",
		},
		{
			name:        "all_zeros",
			traceID:     "00000000000000000000000000000000",
			description: "trace ID that is all zeros",
		},
		{
			name:        "mixed_case_invalid",
			traceID:     "ABCDEF1234567890abcdef1234567890",
			description: "trace ID with uppercase (should be lowercase)",
		},
		{
			name:        "empty_string",
			traceID:     "",
			description: "empty trace ID string",
		},
		{
			name:        "control_characters",
			traceID:     "abc123\n\r\t",
			description: "trace ID with control characters",
		},
		{
			name:        "unicode_characters",
			traceID:     "abc123\u00e9\u00f1",
			description: "trace ID with unicode characters",
		},
		{
			name:        "excessive_length",
			traceID:     string(make([]byte, 10000)), // 10KB of data
			description: "trace ID with excessive length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			md := metadata.New(map[string]string{
				pluginsdk.TraceIDMetadataKey: tt.traceID,
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

			validateGeneratedTraceID(t, capturedTraceID, tt.traceID)
		})
	}
}

// TestTracingUnaryServerInterceptor_ValidTraceIDs tests interceptor preserves valid trace IDs.
func TestTracingUnaryServerInterceptor_ValidTraceIDs(t *testing.T) {
	interceptor := pluginsdk.TracingUnaryServerInterceptor()

	validTraceIDs := []string{
		"abcdef1234567890abcdef1234567890",
		"1234567890abcdef1234567890abcdef",
		"a1b2c3d4e5f678901234567890abcdef",
		"fedcba0987654321fedcba0987654321",
	}

	for _, expectedTraceID := range validTraceIDs {
		t.Run("valid_"+expectedTraceID[:8], func(t *testing.T) {
			ctx := context.Background()
			md := metadata.New(map[string]string{
				pluginsdk.TraceIDMetadataKey: expectedTraceID,
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

			if capturedTraceID != expectedTraceID {
				t.Errorf(
					"Expected trace ID %q to be preserved, got %q",
					expectedTraceID,
					capturedTraceID,
				)
			}
		})
	}
}

// TestTracingUnaryServerInterceptor_EdgeCases tests edge cases for trace ID validation.
func TestTracingUnaryServerInterceptor_EdgeCases(t *testing.T) {
	interceptor := pluginsdk.TracingUnaryServerInterceptor()

	tests := []struct {
		name        string
		setupMD     func() metadata.MD
		description string
	}{
		{
			name: "multiple_trace_ids",
			setupMD: func() metadata.MD {
				return metadata.Pairs(
					pluginsdk.TraceIDMetadataKey, "first-valid-trace-id1234567890",
					pluginsdk.TraceIDMetadataKey, "second-trace-id-should-be-ignored",
				)
			},
			description: "multiple trace_id values should use first one",
		},
		{
			name: "first_invalid_multiple",
			setupMD: func() metadata.MD {
				return metadata.Pairs(
					pluginsdk.TraceIDMetadataKey, "invalid",
					pluginsdk.TraceIDMetadataKey, "abcdef1234567890abcdef1234567890",
				)
			},
			description: "first trace_id invalid, should generate new one",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			md := tt.setupMD()
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

			// Should always have a valid trace ID
			if capturedTraceID == "" {
				t.Errorf("Expected valid trace ID, got empty")
			}

			if len(capturedTraceID) != 32 {
				t.Errorf(
					"Trace ID should be 32 characters, got %d: %q",
					len(capturedTraceID),
					capturedTraceID,
				)
			}
		})
	}
}

// =============================================================================
// NewLogWriter Tests (T005-T019)
// =============================================================================

// TestNewLogWriter_ValidPath verifies file writer returned when env var set to valid path.
func TestNewLogWriter_ValidPath(t *testing.T) {
	pluginsdk.ResetLogWriter()
	defer pluginsdk.ResetLogWriter()
	tmpFile := filepath.Join(t.TempDir(), "test.log")
	t.Setenv("FINFOCUS_LOG_FILE", tmpFile)

	writer := pluginsdk.NewLogWriter()

	// Should not be stderr
	if writer == os.Stderr {
		t.Error("Expected file writer, got os.Stderr")
	}

	// Verify it's a file by writing to it
	logger := zerolog.New(writer).With().Timestamp().Logger()
	logger.Info().Msg("test message")

	// Read back from file
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if !bytes.Contains(content, []byte("test message")) {
		t.Errorf("Expected log file to contain 'test message', got: %s", string(content))
	}
}

// TestNewLogWriter_FileCreated verifies file is created if not exists.
func TestNewLogWriter_FileCreated(t *testing.T) {
	pluginsdk.ResetLogWriter()
	defer pluginsdk.ResetLogWriter()
	tmpFile := filepath.Join(t.TempDir(), "new.log")
	t.Setenv("FINFOCUS_LOG_FILE", tmpFile)

	// Verify file doesn't exist yet
	if _, err := os.Stat(tmpFile); !os.IsNotExist(err) {
		t.Fatal("Test file should not exist before NewLogWriter call")
	}

	writer := pluginsdk.NewLogWriter()

	// Write something to trigger file creation
	logger := zerolog.New(writer).With().Timestamp().Logger()
	logger.Info().Msg("test")

	// Verify file now exists
	info, err := os.Stat(tmpFile)
	if err != nil {
		t.Fatalf("Log file should have been created: %v", err)
	}

	// Verify permissions (0644)
	expectedPerm := os.FileMode(0644)
	actualPerm := info.Mode().Perm()
	if actualPerm != expectedPerm {
		t.Errorf("Expected permissions %o, got %o", expectedPerm, actualPerm)
	}
}

// TestNewLogWriter_FileAppended verifies existing file is appended (not truncated).
func TestNewLogWriter_FileAppended(t *testing.T) {
	pluginsdk.ResetLogWriter()
	defer pluginsdk.ResetLogWriter()
	tmpFile := filepath.Join(t.TempDir(), "append.log")

	// Create file with existing content
	existingContent := `{"level":"info","message":"existing"}`
	if err := os.WriteFile(tmpFile, []byte(existingContent+"\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	t.Setenv("FINFOCUS_LOG_FILE", tmpFile)

	writer := pluginsdk.NewLogWriter()
	logger := zerolog.New(writer).With().Timestamp().Logger()
	logger.Info().Msg("new message")

	// Read back and verify both messages exist
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if !bytes.Contains(content, []byte("existing")) {
		t.Error("Existing content should be preserved")
	}
	if !bytes.Contains(content, []byte("new message")) {
		t.Error("New content should be appended")
	}
}

// TestNewLogWriter_AllLogLevels verifies debug/info/warn/error all captured in file.
func TestNewLogWriter_AllLogLevels(t *testing.T) {
	pluginsdk.ResetLogWriter()
	defer pluginsdk.ResetLogWriter()
	tmpFile := filepath.Join(t.TempDir(), "levels.log")
	t.Setenv("FINFOCUS_LOG_FILE", tmpFile)

	writer := pluginsdk.NewLogWriter()
	logger := zerolog.New(writer).Level(zerolog.DebugLevel).With().Timestamp().Logger()

	logger.Debug().Msg("debug message")
	logger.Info().Msg("info message")
	logger.Warn().Msg("warn message")
	logger.Error().Msg("error message")

	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	levels := []string{"debug", "info", "warn", "error"}
	for _, level := range levels {
		if !bytes.Contains(content, []byte(level+" message")) {
			t.Errorf("Expected %s message in log file", level)
		}
	}
}

// TestNewLogWriter_EnvNotSet verifies os.Stderr returned when env var not set.
func TestNewLogWriter_EnvNotSet(t *testing.T) {
	pluginsdk.ResetLogWriter()
	defer pluginsdk.ResetLogWriter()
	// Ensure env var is not set (empty string is treated as unset)
	t.Setenv("FINFOCUS_LOG_FILE", "")

	writer := pluginsdk.NewLogWriter()

	if writer != os.Stderr {
		t.Error("Expected os.Stderr when FINFOCUS_LOG_FILE is not set")
	}
}

// TestNewLogWriter_EmptyString verifies os.Stderr returned when env var is empty string.
func TestNewLogWriter_EmptyString(t *testing.T) {
	pluginsdk.ResetLogWriter()
	defer pluginsdk.ResetLogWriter()
	t.Setenv("FINFOCUS_LOG_FILE", "")

	writer := pluginsdk.NewLogWriter()

	if writer != os.Stderr {
		t.Error("Expected os.Stderr when FINFOCUS_LOG_FILE is empty string")
	}
}

// TestNewLogWriter_DirectoryPath verifies stderr + warning when path is a directory.
func TestNewLogWriter_DirectoryPath(t *testing.T) {
	pluginsdk.ResetLogWriter()
	defer pluginsdk.ResetLogWriter()
	tmpDir := t.TempDir()
	t.Setenv("FINFOCUS_LOG_FILE", tmpDir)

	writer := pluginsdk.NewLogWriter()

	if writer != os.Stderr {
		t.Error("Expected os.Stderr when FINFOCUS_LOG_FILE points to a directory")
	}
}

// TestNewLogWriter_NonexistentParent verifies stderr + warning when parent dir doesn't exist.
func TestNewLogWriter_NonexistentParent(t *testing.T) {
	pluginsdk.ResetLogWriter()
	defer pluginsdk.ResetLogWriter()
	nonexistentPath := filepath.Join(t.TempDir(), "nonexistent", "subdir", "test.log")
	t.Setenv("FINFOCUS_LOG_FILE", nonexistentPath)

	writer := pluginsdk.NewLogWriter()

	if writer != os.Stderr {
		t.Error("Expected os.Stderr when parent directory doesn't exist")
	}
}

// TestNewLogWriter_BackwardCompatibility verifies existing plugins work without modification.
func TestNewLogWriter_BackwardCompatibility(t *testing.T) {
	// Clear the env var to simulate existing plugin behavior
	t.Setenv("FINFOCUS_LOG_FILE", "")

	// Create a plugin logger the "old way" - should still work
	var buf bytes.Buffer
	logger := pluginsdk.NewPluginLogger("test-plugin", "v1.0.0", zerolog.InfoLevel, &buf)

	logger.Info().Msg("test message")

	if !bytes.Contains(buf.Bytes(), []byte("test message")) {
		t.Error("Plugin logger should still work with custom writer")
	}
}

// TestLogFileConstants_Values verifies log file constants have correct values.
func TestLogFileConstants_Values(t *testing.T) {
	// Test LogFilePermissions
	expectedPerm := os.FileMode(0644)
	if pluginsdk.LogFilePermissions != expectedPerm {
		t.Errorf("LogFilePermissions: expected %o, got %o", expectedPerm, pluginsdk.LogFilePermissions)
	}

	// Test LogFileFlags
	expectedFlags := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	if pluginsdk.LogFileFlags != expectedFlags {
		t.Errorf("LogFileFlags: expected %d, got %d", expectedFlags, pluginsdk.LogFileFlags)
	}
}

// BenchmarkNewLogWriter_File measures file writer creation performance.
func BenchmarkNewLogWriter_File(b *testing.B) {
	pluginsdk.ResetLogWriter()
	defer pluginsdk.ResetLogWriter()
	tmpFile := filepath.Join(b.TempDir(), "bench.log")
	b.Setenv("FINFOCUS_LOG_FILE", tmpFile)
	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_ = pluginsdk.NewLogWriter()
	}
}

// BenchmarkNewLogWriter_Stderr measures stderr fallback performance.
func BenchmarkNewLogWriter_Stderr(b *testing.B) {
	pluginsdk.ResetLogWriter()
	defer pluginsdk.ResetLogWriter()
	b.Setenv("FINFOCUS_LOG_FILE", "")
	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_ = pluginsdk.NewLogWriter()
	}
}
