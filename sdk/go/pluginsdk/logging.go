package pluginsdk

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rshade/pulumicost-spec/sdk/go/pricing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// contextKey is the type for context keys to avoid collisions.
type contextKey string

const traceIDKey contextKey = "pulumicost-trace-id"

// TraceIDMetadataKey is the gRPC metadata header for trace ID propagation.
const TraceIDMetadataKey = "x-pulumicost-trace-id"

// Standard field names for structured logging consistency across plugins.
const (
	FieldTraceID       = "trace_id"
	FieldComponent     = "component"
	FieldOperation     = "operation"
	FieldDurationMs    = "duration_ms"
	FieldResourceURN   = "resource_urn"
	FieldResourceType  = "resource_type"
	FieldProvider      = "provider"
	FieldRegion        = "region"
	FieldPluginName    = "plugin_name"
	FieldPluginVersion = "plugin_version"
	FieldCostMonthly   = "cost_monthly"
	FieldAdapter       = "adapter"
	FieldErrorCode     = "error_code"
)

// newDefaultLogger creates a basic zerolog logger for internal SDK use.
// This is used when no custom logger is provided.
func newDefaultLogger() zerolog.Logger {
	return zerolog.New(os.Stderr).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()
}

// NewPluginLogger creates a configured zerolog logger for plugins.
//
// Parameters:
//   - pluginName: Identifier for the plugin (e.g., "aws-public")
//   - version: Plugin version (e.g., "v1.0.0")
//   - level: Minimum log level to output
//   - w: Output writer (nil defaults to os.Stderr)
//
// Returns a logger with plugin_name and plugin_version fields pre-configured.
func NewPluginLogger(pluginName, version string, level zerolog.Level, w io.Writer) zerolog.Logger {
	if w == nil {
		w = os.Stderr
	}

	return zerolog.New(w).
		Level(level).
		With().
		Timestamp().
		Str(FieldPluginName, pluginName).
		Str(FieldPluginVersion, version).
		Logger()
}

// TracingUnaryServerInterceptor returns a gRPC server interceptor that extracts
// trace_id from incoming request metadata, validates it, and adds it to the request context.
//
// The interceptor looks for the TraceIDMetadataKey header. If the trace_id is missing or
// invalid, a new valid trace_id is generated. The validated or generated trace_id is
// stored in the context for retrieval via TraceIDFromContext.
func TracingUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		var traceID string

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if values := md.Get(TraceIDMetadataKey); len(values) > 0 {
				traceID = values[0]
			}
		}

		// Validate the trace ID; generate a new one if invalid or missing
		if traceID == "" || pricing.ValidateTraceID(traceID) != nil {
			var err error
			traceID, err = GenerateTraceID()
			if err != nil {
				// If generation fails, proceed with empty trace ID
				// This maintains request flow even in extreme failure cases
				traceID = ""
			}
		}

		ctx = ContextWithTraceID(ctx, traceID)
		return handler(ctx, req)
	}
}

// TraceIDFromContext extracts the trace ID from the given context.
//
// Returns empty string if no trace ID is present in the context.
func TraceIDFromContext(ctx context.Context) string {
	if val := ctx.Value(traceIDKey); val != nil {
		if traceID, ok := val.(string); ok {
			return traceID
		}
	}
	return ""
}

// ContextWithTraceID returns a new context with the trace ID stored.
//
// This is typically called by the interceptor, but can be used directly
// for testing or manual trace ID injection.
func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// LogOperation returns a function that logs the operation duration when called.
//
// Usage:
//
//	done := LogOperation(logger, "GetProjectedCost")
//	defer done()
//	// ... perform operation ...
//
// The returned function logs at Info level with FieldOperation and FieldDurationMs.
func LogOperation(logger zerolog.Logger, operation string) func() {
	start := time.Now()
	return func() {
		logger.Info().
			Str(FieldOperation, operation).
			Int64(FieldDurationMs, time.Since(start).Milliseconds()).
			Msg("operation completed")
	}
}
