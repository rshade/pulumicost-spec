# Research: SDK Support for PULUMICOST_LOG_FILE

**Date**: 2025-12-08
**Feature**: 015-log-file

## Research Questions

### 1. Existing SDK Infrastructure

**Question**: What logging infrastructure already exists in the SDK?

**Findings**:

The `sdk/go/pluginsdk` package has established logging infrastructure:

1. **Environment Variable Support** (`env.go`):
   - `EnvLogFile = "PULUMICOST_LOG_FILE"` - constant already defined
   - `GetLogFile() string` - function already implemented, returns env var value
   - Empty string indicates stdout/stderr (per existing documentation)

2. **Logger Creation** (`logging.go`):
   - `NewPluginLogger(pluginName, version string, level zerolog.Level, w io.Writer)` - accepts
     writer parameter but requires manual file handling
   - `newDefaultLogger()` - internal function hardcoded to `os.Stderr`
   - Standard field names defined (`FieldTraceID`, `FieldPluginName`, etc.)

3. **Server Integration** (`sdk.go`):
   - `ServeConfig.Logger` - optional logger field
   - `NewServerWithOptions()` - uses `newDefaultLogger()` when logger is nil
   - `newDefaultLogger()` called internally when no custom logger provided

**Decision**: Enhance existing infrastructure rather than create new patterns.

**Rationale**: The building blocks exist - we need to wire them together with automatic file
handling.

---

### 2. Go File Handling Best Practices

**Question**: What are the best practices for log file handling in Go?

**Findings**:

1. **File Opening Flags**:
   - `os.O_APPEND|os.O_CREATE|os.O_WRONLY` - standard for log files
   - `O_APPEND` ensures atomic appends for concurrent writers
   - `O_CREATE` creates file if not exists
   - `O_WRONLY` prevents unnecessary read access

2. **File Permissions**:
   - `0644` (rw-r--r--) - standard for log files
   - Owner can read/write, others can read
   - Follows Unix convention for non-sensitive logs

3. **Error Handling**:
   - Check if path is a directory before opening
   - Fall back gracefully on permission errors
   - Log warning to stderr before switching output

4. **File Handle Lifecycle**:
   - Open once at startup
   - Keep handle open for duration of process
   - Close on process exit (handled by OS)
   - No explicit close needed for single-process logging

**Decision**: Use `os.OpenFile` with append/create/writeonly flags and 0644 permissions.

**Rationale**: This is the standard Go pattern for log file handling that handles concurrent
access correctly.

**Alternatives Considered**:

- `os.Create()` - rejected because it truncates existing files
- `bufio.Writer` wrapper - rejected as zerolog already buffers appropriately
- Explicit file rotation - rejected as out of scope (external tool responsibility)

---

### 3. Zerolog Writer Configuration

**Question**: How does zerolog handle writer configuration?

**Findings**:

1. **Writer Interface**:
   - Zerolog accepts any `io.Writer` for output
   - `zerolog.New(writer)` creates logger with specified writer
   - Writer can be changed at logger creation time only

2. **Multi-Writer Support**:
   - `zerolog.MultiLevelWriter()` for writing to multiple outputs
   - Not needed for this feature (single output destination)

3. **ConsoleWriter**:
   - `zerolog.ConsoleWriter{Out: os.Stderr}` for human-readable output
   - Currently used with `PULUMICOST_LOG_FORMAT=text`
   - Must wrap file writer in ConsoleWriter when format is text

**Decision**: Create logger with file writer at initialization; respect LOG_FORMAT setting.

**Rationale**: Zerolog's design expects writer configuration at creation time.

---

### 4. Concurrent Plugin Access

**Question**: What happens when multiple plugins write to the same log file?

**Findings**:

1. **POSIX Append Semantics**:
   - `O_APPEND` flag ensures atomic positioning and write
   - Each write is appended at current end of file
   - Works correctly for multiple processes

2. **Go's os.File**:
   - Uses underlying OS append semantics
   - Safe for concurrent writes from multiple processes
   - Each log line written atomically (up to PIPE_BUF size)

3. **Line Length Considerations**:
   - Zerolog produces single-line JSON output
   - Typical log lines well under 4KB (PIPE_BUF minimum)
   - Atomic writes guaranteed for reasonable log lines

**Decision**: Use append mode; document that concurrent access is safe.

**Rationale**: OS-level append semantics provide correct behavior without application-level
locking.

---

### 5. Invalid Path Handling

**Question**: How should the SDK handle invalid log file paths?

**Findings**:

1. **Error Categories**:
   - Path is a directory → cannot open as file
   - Parent directory doesn't exist → open fails
   - Permission denied → open fails
   - Disk full → write fails (after open succeeds)

2. **Detection Timing**:
   - Directory check: before opening
   - Open errors: at startup
   - Write errors: during logging (handled by zerolog)

3. **Fallback Strategy**:
   - Log warning to stderr
   - Continue with stderr as output
   - Don't crash the plugin

**Decision**: Validate path at startup, fall back to stderr with warning on any error.

**Rationale**: Plugin reliability is more important than strict log file enforcement.

---

## API Design

### New Public API

```go
// NewLogWriter returns an io.Writer configured based on PULUMICOST_LOG_FILE.
// If the environment variable is set to a valid path, returns a file writer.
// If not set or invalid, returns os.Stderr and logs a warning for invalid paths.
//
// The returned writer should be used with NewPluginLogger or passed to zerolog directly.
// The file is opened in append mode with 0644 permissions.
//
// Example:
//
//  writer := pluginsdk.NewLogWriter()
//  logger := pluginsdk.NewPluginLogger("my-plugin", "v1.0.0", zerolog.InfoLevel, writer)
//
func NewLogWriter() io.Writer
```

### Modified Internal API

```go
// newDefaultLogger creates a zerolog logger using NewLogWriter().
// This automatically respects PULUMICOST_LOG_FILE when set.
func newDefaultLogger() zerolog.Logger
```

## Conclusion

All technical questions resolved. Implementation path is clear:

1. Add `NewLogWriter()` public function in `logging.go`
2. Update `newDefaultLogger()` to use `NewLogWriter()`
3. Add comprehensive tests for all scenarios
4. Update README documentation
