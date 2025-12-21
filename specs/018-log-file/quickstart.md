# Quickstart: SDK Log File Support

**Feature**: 015-log-file

## Overview

The PulumiCost SDK supports redirecting plugin logs to a file via the `PULUMICOST_LOG_FILE`
environment variable. This enables the Core CLI to orchestrate plugins without log output
polluting the user-facing interface.

## Basic Usage

### For Plugin Developers

No code changes required. If your plugin uses the SDK's default logging, it automatically
respects `PULUMICOST_LOG_FILE`:

```go
// Your existing plugin code works unchanged
func main() {
    flag.Parse()
    ctx := context.Background()
    if err := pluginsdk.Serve(ctx, pluginsdk.ServeConfig{
        Plugin: &MyPlugin{},
    }); err != nil {
        log.Fatal(err)
    }
}
```

When `PULUMICOST_LOG_FILE` is set, all SDK logs go to the specified file.

### For Core CLI Developers

Set the environment variable before spawning plugins:

```bash
# Direct all plugin logs to a single file
export PULUMICOST_LOG_FILE=/var/log/pulumicost/plugins.log
./my-plugin
```

Or set per-plugin for separate log files:

```bash
PULUMICOST_LOG_FILE=/var/log/pulumicost/aws.log ./aws-plugin &
PULUMICOST_LOG_FILE=/var/log/pulumicost/azure.log ./azure-plugin &
```

## Custom Logger Configuration

If you create a custom logger, use `NewLogWriter()` to respect the environment variable:

```go
import (
    "github.com/rs/zerolog"
    "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
)

func main() {
    // Get writer that respects PULUMICOST_LOG_FILE
    writer := pluginsdk.NewLogWriter()

    // Create logger with custom configuration
    logger := pluginsdk.NewPluginLogger(
        "my-plugin",
        "v1.0.0",
        zerolog.DebugLevel, // or parse from PULUMICOST_LOG_LEVEL
        writer,
    )

    // Use the logger
    logger.Info().Msg("Plugin started")
}
```

## Environment Variables

| Variable | Purpose | Default |
|----------|---------|---------|
| `PULUMICOST_LOG_FILE` | Log file path | stderr |
| `PULUMICOST_LOG_LEVEL` | Log verbosity | info |
| `PULUMICOST_LOG_FORMAT` | Output format (json/text) | json |

## Behavior

### Valid Path

```bash
export PULUMICOST_LOG_FILE=/tmp/plugin.log
./my-plugin
# Logs written to /tmp/plugin.log
# stderr is clean
```

### Invalid Path

```bash
export PULUMICOST_LOG_FILE=/nonexistent/dir/plugin.log
./my-plugin
# Warning logged to stderr: "failed to open log file, falling back to stderr"
# Subsequent logs go to stderr
```

### Not Set

```bash
unset PULUMICOST_LOG_FILE
./my-plugin
# All logs go to stderr (default behavior)
```

## Testing Log File Configuration

```go
func TestLogFileConfiguration(t *testing.T) {
    // Set up test file
    tmpFile := filepath.Join(t.TempDir(), "test.log")
    t.Setenv("PULUMICOST_LOG_FILE", tmpFile)

    // Get writer (should be file)
    writer := pluginsdk.NewLogWriter()

    // Create logger and log a message
    logger := zerolog.New(writer).With().Timestamp().Logger()
    logger.Info().Msg("test message")

    // Verify log was written to file
    content, err := os.ReadFile(tmpFile)
    require.NoError(t, err)
    assert.Contains(t, string(content), "test message")
}
```

## Best Practices

1. **Let SDK handle logging**: Use default logger when possible
2. **Use absolute paths**: Relative paths depend on working directory
3. **Create parent directories**: SDK creates the file but not parent directories
4. **External rotation**: Use logrotate or similar for log file management
5. **Aggregate logs**: Consider sending all plugins to one file for easier debugging
