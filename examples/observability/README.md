# FinFocus Plugin Observability Examples

This directory contains comprehensive examples demonstrating observability implementation for FinFocus plugins.

## Contents

### Metric Examples

- **[basic_metrics.json](basic_metrics.json)** - Example metrics that all plugins should expose
- **[service_level_indicators.json](service_level_indicators.json)** - SLI measurements for monitoring
- **[health_check.json](health_check.json)** - Health check response examples
- **[structured_logs.json](structured_logs.json)** - Structured logging format examples

### Dashboard Templates

- **[dashboards/grafana-dashboard.json](dashboards/grafana-dashboard.json)** - Grafana dashboard for Prometheus
  metrics
- **[dashboards/cloudwatch-dashboard.yaml](dashboards/cloudwatch-dashboard.yaml)** - CloudWatch dashboard
  CloudFormation template

## Usage

### Implementing Basic Observability

1. **Health Checks**: Implement the `HealthCheck` RPC method

   ```go
   func (p *Plugin) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
       // Check dependencies and return status
   }
   ```

2. **Metrics Collection**: Implement the `GetMetrics` RPC method

   ```go
   func (p *Plugin) GetMetrics(ctx context.Context, req *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {
       // Collect and return metrics
   }
   ```

3. **SLI Reporting**: Implement the `GetServiceLevelIndicators` RPC method

   ```go
   func (p *Plugin) GetServiceLevelIndicators(ctx context.Context, req *pb.GetSLIRequest) (*pb.GetSLIResponse, error) {
       // Calculate and return SLIs
   }
   ```

### Conformance Levels

#### Basic Conformance (Required)

- Health check endpoint
- Request count, error count, and latency metrics
- Availability and error rate SLIs

#### Standard Conformance (Production)

- All basic features
- Cache hit rate and connection metrics
- Latency percentiles (P95)
- Structured logging

#### Advanced Conformance (Enterprise)

- All standard features
- Custom metrics support
- Full distributed tracing
- Resource usage monitoring (memory, CPU)

### Testing Your Implementation

Use the provided test framework:

```go
import "github.com/rshade/finfocus-spec/sdk/go/testing"

func TestObservability(t *testing.T) {
    plugin := &MyPlugin{}
    suite := testing.NewObservabilityTestSuite(plugin, t)

    // Run conformance tests
    if !suite.RunStandardObservabilityTests() {
        t.Fatal("Observability tests failed")
    }
}
```

### Deploying Monitoring

#### Prometheus + Grafana

1. Configure Prometheus to scrape metrics:

   ```yaml
   scrape_configs:
     - job_name: "finfocus-plugin"
       static_configs:
         - targets: ["plugin:8080"]
       metrics_path: "/metrics"
   ```

2. Import the Grafana dashboard from `dashboards/grafana-dashboard.json`

#### AWS CloudWatch

Deploy the CloudWatch dashboard using CloudFormation:

```bash
aws cloudformation create-stack \
  --stack-name finfocus-observability \
  --template-body file://dashboards/cloudwatch-dashboard.yaml \
  --parameters ParameterKey=PluginName,ParameterValue=my-plugin
```

## Metrics Reference

### Standard Metrics

| Metric Name                           | Type      | Description                 | Labels                             |
| ------------------------------------- | --------- | --------------------------- | ---------------------------------- |
| `finfocus_requests_total`           | Counter   | Total requests processed    | `method`, `provider`, `status`     |
| `finfocus_errors_total`             | Counter   | Total errors encountered    | `method`, `provider`, `error_type` |
| `finfocus_request_duration_seconds` | Histogram | Request processing time     | `method`, `provider`               |
| `finfocus_cache_hit_rate_percent`   | Gauge     | Cache effectiveness         | `provider`, `cache_type`           |
| `finfocus_active_connections`       | Gauge     | Active external connections | `service`                          |

### Standard SLIs

| SLI Name         | Description               | Target | Unit                |
| ---------------- | ------------------------- | ------ | ------------------- |
| `availability`   | Service uptime percentage | 99.9%  | percentage          |
| `error_rate`     | Error rate percentage     | <0.1%  | percentage          |
| `latency_p95`    | 95th percentile latency   | <1s    | seconds             |
| `latency_p99`    | 99th percentile latency   | <2s    | seconds             |
| `throughput`     | Requests per second       | >100   | requests_per_second |
| `data_freshness` | Cost data age             | <24h   | hours               |

## Best Practices

### Metric Naming

- Use the `finfocus_` prefix for all metrics
- Follow Prometheus naming conventions
- Keep label cardinality low (<100 unique combinations)

### Health Checks

- Check critical dependencies (database, external APIs)
- Return detailed status information
- Cache results for short periods to reduce overhead

### Logging

- Use structured logging with consistent fields
- Include trace context for correlation
- Sanitize sensitive information

### Tracing

- Propagate trace context across service boundaries
- Add meaningful span attributes
- Use sampling for high-volume scenarios

## Troubleshooting

### High Cardinality Issues

If you see performance problems with metrics:

- Reduce the number of unique label values
- Use sampling for high-volume metrics
- Consider aggregating metrics at collection time

### Missing Trace Context

If traces aren't connecting across services:

- Verify OpenTelemetry propagation is configured
- Check that trace context is being extracted and injected
- Validate trace and span ID formats

### Health Check Timeouts

If health checks are timing out:

- Reduce dependency check timeouts
- Cache health check results
- Return partial status when some checks fail

## Additional Resources

- [Observability Implementation Guide](../../OBSERVABILITY_GUIDE.md)
- [Testing Framework Documentation](../sdk/go/testing/README.md)
- [OpenTelemetry Go Documentation](https://opentelemetry.io/docs/instrumentation/go/)
- [Prometheus Metric Types](https://prometheus.io/docs/concepts/metric_types/)
- [Grafana Dashboard Documentation](https://grafana.com/docs/grafana/latest/dashboards/)
