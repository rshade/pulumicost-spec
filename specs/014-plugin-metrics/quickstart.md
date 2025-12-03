# Quickstart: Plugin Metrics

This guide shows how to add Prometheus metrics to your PulumiCost plugin in under 5 minutes.

## Basic Usage (Recommended)

Add the metrics interceptor to your plugin's server configuration:

```go
package main

import (
    "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
    "google.golang.org/grpc"
)

func main() {
    plugin := &MyPlugin{}

    config := pluginsdk.ServeConfig{
        Plugin: plugin,
        UnaryInterceptors: []grpc.UnaryServerInterceptor{
            // Add metrics interceptor - that's it!
            pluginsdk.MetricsUnaryServerInterceptor("my-plugin"),
        },
    }

    if err := pluginsdk.Serve(context.Background(), config); err != nil {
        log.Fatal(err)
    }
}
```

## Exposing Metrics

### Option A: Use the Helper Server (Simple)

For quick setup, use the built-in metrics server:

```go
// Start metrics server on port 9090
server, err := pluginsdk.StartMetricsServer(pluginsdk.MetricsServerConfig{
    Port: 9090,
})
if err != nil {
    log.Fatal(err)
}
defer server.Shutdown(context.Background())
```

Then scrape `http://localhost:9090/metrics`.

### Option B: Custom HTTP Handler (Production)

For more control, access the registry directly:

```go
metrics := pluginsdk.NewPluginMetrics("my-plugin")
interceptor := pluginsdk.MetricsInterceptorWithRegistry(metrics)

// Add to your existing HTTP server
http.Handle("/metrics", promhttp.HandlerFor(
    metrics.Registry,
    promhttp.HandlerOpts{},
))
```

## Available Metrics

After enabling, you'll see these metrics at `/metrics`:

### Request Counter

```text
# HELP pulumicost_plugin_requests_total Total gRPC requests
# TYPE pulumicost_plugin_requests_total counter
pulumicost_plugin_requests_total{grpc_method="GetProjectedCost",grpc_code="OK",plugin_name="my-plugin"} 42
pulumicost_plugin_requests_total{grpc_method="GetProjectedCost",grpc_code="Internal",plugin_name="my-plugin"} 3
```

### Request Duration

```text
# HELP pulumicost_plugin_request_duration_seconds Request duration histogram
# TYPE pulumicost_plugin_request_duration_seconds histogram
pulumicost_plugin_request_duration_seconds_bucket{grpc_method="GetProjectedCost",plugin_name="my-plugin",le="0.005"} 10
pulumicost_plugin_request_duration_seconds_bucket{grpc_method="GetProjectedCost",plugin_name="my-plugin",le="0.01"} 25
...
pulumicost_plugin_request_duration_seconds_sum{grpc_method="GetProjectedCost",plugin_name="my-plugin"} 1.234
pulumicost_plugin_request_duration_seconds_count{grpc_method="GetProjectedCost",plugin_name="my-plugin"} 45
```

## Prometheus Configuration

Add to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'pulumicost-plugins'
    static_configs:
      - targets: ['localhost:9090']
```

## Useful PromQL Queries

### Request Rate by Method

```promql
sum(rate(pulumicost_plugin_requests_total[5m])) by (grpc_method)
```

### Error Rate

```promql
sum(rate(pulumicost_plugin_requests_total{grpc_code!="OK"}[5m]))
/ sum(rate(pulumicost_plugin_requests_total[5m]))
```

### P99 Latency

```promql
histogram_quantile(0.99, sum(rate(pulumicost_plugin_request_duration_seconds_bucket[5m])) by (le, grpc_method))
```

## Complete Example

```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
    "google.golang.org/grpc"
)

func main() {
    plugin := &MyPlugin{}

    // Create metrics with custom registry access
    metrics := pluginsdk.NewPluginMetrics("my-plugin")

    config := pluginsdk.ServeConfig{
        Plugin: plugin,
        UnaryInterceptors: []grpc.UnaryServerInterceptor{
            pluginsdk.MetricsInterceptorWithRegistry(metrics),
        },
    }

    // Start metrics HTTP server
    metricsServer, err := pluginsdk.StartMetricsServer(pluginsdk.MetricsServerConfig{
        Port:     9090,
        Registry: metrics.Registry,
    })
    if err != nil {
        log.Fatal(err)
    }

    // Graceful shutdown
    go func() {
        sigCh := make(chan os.Signal, 1)
        signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
        <-sigCh
        metricsServer.Shutdown(context.Background())
    }()

    // Start gRPC server
    if err := pluginsdk.Serve(context.Background(), config); err != nil {
        log.Fatal(err)
    }
}
```

## Next Steps

- Configure alerting rules based on error rates and latency
- Set up Grafana dashboards for visualization
- Review the [data model](data-model.md) for detailed metric specifications
