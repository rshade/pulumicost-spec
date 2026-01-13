package pluginsdk

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rshade/finfocus-spec/sdk/go/pricing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// contextKey is the type for context keys to avoid collisions.
type contextKey string

const traceIDKey contextKey = "finfocus-trace-id"

// TraceIDMetadataKey is the gRPC metadata header for trace ID propagation.
const TraceIDMetadataKey = "x-finfocus-trace-id"

// Log file configuration constants.
const (
	// LogFilePermissions is the file permission mode for created log files (rw-r--r--).
	// Note: If log files may contain sensitive information (API keys, resource details),
	// consider using 0640 and running plugins under a dedicated user account.
	LogFilePermissions = 0644

	// LogFileFlags are the flags used when opening log files for writing.
	// O_APPEND ensures atomic appends when multiple processes write to the same file (POSIX).
	// O_CREATE creates the file if it doesn't exist.
	// O_WRONLY opens for write-only access.
	LogFileFlags = os.O_APPEND | os.O_CREATE | os.O_WRONLY
)

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

	// GetRecommendations-specific fields.
	FieldRecommendationCount = "recommendation_count"
	FieldFilterCategory      = "filter_category"
	FieldFilterActionType    = "filter_action_type"
	FieldPageSize            = "page_size"
	FieldTotalSavings        = "total_savings"

	// GetBudgets-specific fields.
	FieldIncludeStatus   = "include_status"
	FieldTotalBudgets    = "total_budgets"
	FieldBudgetsOk       = "budgets_ok"
	FieldBudgetsWarning  = "budgets_warning"
	FieldBudgetsCritical = "budgets_critical"
	FieldBudgetsExceeded = "budgets_exceeded"
)

//nolint:gochecknoglobals // Intentional singleton for log file handle reuse (process lifetime)
var (
	logWriterOnce sync.Once
	logWriter     io.Writer
)

// NewLogWriter returns an io.Writer configured based on the FINFOCUS_LOG_FILE
// environment variable. The file is opened once on first call and reused for
// subsequent calls (singleton pattern). The file remains open for the lifetime
// of the process, which is intentional as plugins log continuously during their
// execution. The OS cleans up the file descriptor when the process exits.
//
// If the environment variable is not set, empty, or invalid, returns os.Stderr.
//
// When the path is invalid or inaccessible, a warning is logged to stderr
// before falling back to stderr as the output destination.
//
// The file is opened in append mode with 0644 permissions, allowing multiple
// plugins to safely write to the same log file concurrently.
//
// Example usage:
//
//	writer := pluginsdk.NewLogWriter()
//	logger := pluginsdk.NewPluginLogger("my-plugin", "v1.0.0", zerolog.InfoLevel, writer)
//
//	// Note: The SDK's internal default logger (used when no custom logger is
//	// provided to ServeConfig) respects FINFOCUS_LOG_FILE automatically.
func NewLogWriter() io.Writer {
	logWriterOnce.Do(func() {
		logFile := GetLogFile()

		// Return stderr if env var is not set or empty
		if logFile == "" {
			logWriter = os.Stderr
			return
		}

		warnLogger := zerolog.New(os.Stderr).With().Timestamp().Logger()

		// Check for absolute path
		if !filepath.IsAbs(logFile) {
			warnLogger.Warn().Str("path", logFile).
				Msg("FINFOCUS_LOG_FILE should be absolute path, falling back to stderr")
			logWriter = os.Stderr
			return
		}

		// Attempt to open/create the log file
		file, err := os.OpenFile(logFile, LogFileFlags, LogFilePermissions)
		if err != nil {
			// Check if it's a directory error for better error message
			if info, statErr := os.Stat(logFile); statErr == nil && info.IsDir() {
				warnLogger.Warn().
					Str("path", logFile).
					Msg("FINFOCUS_LOG_FILE points to a directory, falling back to stderr")
				logWriter = os.Stderr
				return
			}

			warnLogger.Warn().
				Str("path", logFile).
				Err(err).
				Msg("failed to open log file, falling back to stderr")
			logWriter = os.Stderr
			return
		}

		logWriter = file
	})

	return logWriter
}

// ResetLogWriter resets the singleton log writer.
// This is primarily for testing purposes to allow re-initialization
// when the environment configuration changes.
func ResetLogWriter() {
	logWriterOnce = sync.Once{}
	logWriter = nil
}

// newDefaultLogger creates a basic zerolog logger for internal SDK use.
// This is used when no custom logger is provided.
// It respects the FINFOCUS_LOG_FILE environment variable via NewLogWriter().
func newDefaultLogger() zerolog.Logger {
	return zerolog.New(NewLogWriter()).
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
