# PulumiCost Plugin Observability Implementation Guide

This guide provides comprehensive instructions for implementing observability features in PulumiCost plugins,
including telemetry, health checks, metrics, and distributed tracing.

## Table of Contents

- [Overview](#overview)
- [Conformance Levels](#conformance-levels)
- [Implementation Checklist](#implementation-checklist)
- [Health Checks](#health-checks)
- [Metrics Collection](#metrics-collection)
- [Service Level Indicators (SLIs)](#service-level-indicators-slis)
- [Distributed Tracing](#distributed-tracing)
- [Structured Logging](#structured-logging)
- [Testing Observability](#testing-observability)
- [Dashboard Templates](#dashboard-templates)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

## Overview

Observability enables monitoring, debugging, and performance optimization of PulumiCost plugins. The specification
defines three conformance levels with increasing requirements:

- **Basic**: Essential observability for all plugins
- **Standard**: Production-ready observability features
- **Advanced**: High-performance and enterprise-grade observability

## Conformance Levels

### Basic Conformance (Required for All Plugins)

**Required Features:**

- Health check endpoint
- Basic metrics collection (requests, errors, latency)
- Standard error codes and messages

**Required Metrics:**

- `pulumicost_requests_total` - Total requests processed
- `pulumicost_errors_total` - Total errors encountered  
- `pulumicost_request_duration_seconds` - Request processing time

**Required SLIs:**

- `availability` - Service availability percentage
- `error_rate` - Error rate percentage

### Standard Conformance (Recommended for Production)

**Additional Features:**

- Metrics endpoint with multiple formats
- SLI reporting endpoint
- Structured logging
- Basic tracing context

**Additional Metrics:**

- `pulumicost_latency_p95_seconds` - 95th percentile latency
- `pulumicost_cache_hit_rate_percent` - Cache effectiveness
- `pulumicost_active_connections` - Active external connections

**Additional SLIs:**

- `latency_p95` - 95th percentile response time
- `throughput` - Requests per second

### Advanced Conformance (High-Performance Requirements)

**Additional Features:**

- Custom metrics registration
- Full distributed tracing
- Performance monitoring
- Resource usage tracking

**Additional Metrics:**

- `pulumicost_latency_p99_seconds` - 99th percentile latency
- `pulumicost_memory_usage_bytes` - Memory consumption
- `pulumicost_cpu_usage_percent` - CPU utilization
- `pulumicost_data_source_latency_seconds` - External API latency

**Additional SLIs:**

- `latency_p99` - 99th percentile response time
- `data_freshness` - Age of cost data

## Implementation Checklist

### 1. Implement ObservabilityService Interface

```protobuf
service ObservabilityService {
  rpc HealthCheck(HealthCheckRequest) returns (HealthCheckResponse);
  rpc GetMetrics(GetMetricsRequest) returns (GetMetricsResponse);
  rpc GetServiceLevelIndicators(GetSLIRequest) returns (GetSLIResponse);
}
```

### 2. Basic Health Check Implementation

```go
func (p *Plugin) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
    status := pb.HealthCheckResponse_SERVING
    message := "All systems operational"
    
    // Check critical dependencies
    if err := p.checkDatabase(ctx); err != nil {
        status = pb.HealthCheckResponse_NOT_SERVING
        message = fmt.Sprintf("Database unavailable: %v", err)
    }
    
    if err := p.checkExternalAPIs(ctx); err != nil {
        if status == pb.HealthCheckResponse_SERVING {
            status = pb.HealthCheckResponse_SERVING // Degraded but still serving
            message = fmt.Sprintf("External APIs degraded: %v", err)
        }
    }
    
    return &pb.HealthCheckResponse{
        Status:        status,
        Message:       message,
        LastCheckTime: timestamppb.Now(),
    }, nil
}
```

### 3. Metrics Collection Implementation

```go
func (p *Plugin) GetMetrics(ctx context.Context, req *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {
    var metrics []*pb.Metric
    
    // Collect standard metrics
    if shouldIncludeMetric(req.MetricNames, pricing.StandardMetrics.RequestsTotal) {
        metrics = append(metrics, p.collectRequestMetrics())
    }
    
    if shouldIncludeMetric(req.MetricNames, pricing.StandardMetrics.ErrorsTotal) {
        metrics = append(metrics, p.collectErrorMetrics())
    }
    
    if shouldIncludeMetric(req.MetricNames, pricing.StandardMetrics.RequestDurationSeconds) {
        metrics = append(metrics, p.collectLatencyMetrics())
    }
    
    return &pb.GetMetricsResponse{
        Metrics:   metrics,
        Timestamp: timestamppb.Now(),
        Format:    req.Format,
    }, nil
}

func (p *Plugin) collectRequestMetrics() *pb.Metric {
    return &pb.Metric{
        Name: pricing.StandardMetrics.RequestsTotal,
        Help: "Total number of requests processed by the plugin",
        Type: pb.MetricType_COUNTER,
        Samples: []*pb.MetricSample{
            {
                Labels: map[string]string{
                    "method":   "GetActualCost",
                    "provider": "aws",
                    "status":   "success",
                },
                Value:     float64(p.requestCounts["GetActualCost"]["aws"]["success"]),
                Timestamp: timestamppb.Now(),
            },
            // ... more samples
        },
    }
}
```

### 4. SLI Reporting Implementation

```go
func (p *Plugin) GetServiceLevelIndicators(ctx context.Context, req *pb.GetSLIRequest) (*pb.GetSLIResponse, error) {
    var slis []*pb.ServiceLevelIndicator
    
    // Calculate availability SLI
    availability := p.calculateAvailability(req.TimeRange)
    slis = append(slis, &pb.ServiceLevelIndicator{
        Name:        "availability",
        Description: "Percentage of successful requests over total requests",
        Value:       availability,
        Unit:        "percentage",
        TargetValue: 99.9,
        Status:      p.getSLIStatus(availability, 99.9),
    })
    
    // Calculate error rate SLI
    errorRate := p.calculateErrorRate(req.TimeRange)
    slis = append(slis, &pb.ServiceLevelIndicator{
        Name:        "error_rate",
        Description: "Percentage of requests that result in errors",
        Value:       errorRate,
        Unit:        "percentage",
        TargetValue: 0.1,
        Status:      p.getSLIStatus(0.1, errorRate), // Lower is better
    })
    
    return &pb.GetSLIResponse{
        Slis:            slis,
        MeasurementTime: timestamppb.Now(),
    }, nil
}
```

## Health Checks

### Implementation Strategy

1. **Check Critical Dependencies**: Database, external APIs, cache
2. **Return Detailed Status**: Overall status plus individual component status
3. **Include Response Times**: Help identify performance issues
4. **Use Appropriate Status Codes**: `SERVING`, `NOT_SERVING`, `SERVICE_UNKNOWN`

### Example Health Check Response

```json
{
  "status": "SERVING",
  "message": "All systems operational",
  "last_check_time": "2024-01-15T10:30:00Z",
  "checks": {
    "database": {
      "status": "SERVING",
      "response_time_ms": 12.5
    },
    "aws_billing_api": {
      "status": "SERVING",
      "response_time_ms": 245.3
    }
  }
}
```

## Metrics Collection

### Standard Metrics

All plugins must implement these core metrics:

| Metric Name | Type | Description | Labels |
|-------------|------|-------------|---------|
| `pulumicost_requests_total` | Counter | Total requests processed | `method`, `provider`, `status` |
| `pulumicost_errors_total` | Counter | Total errors encountered | `method`, `provider`, `error_type` |
| `pulumicost_request_duration_seconds` | Histogram | Request processing time | `method`, `provider` |

### Advanced Metrics

For Standard and Advanced conformance:

| Metric Name | Type | Description | Labels |
|-------------|------|-------------|---------|
| `pulumicost_cache_hit_rate_percent` | Gauge | Cache effectiveness | `provider`, `cache_type` |
| `pulumicost_active_connections` | Gauge | Active external connections | `service` |
| `pulumicost_memory_usage_bytes` | Gauge | Memory consumption | `component` |

### Metric Implementation Best Practices

```go
// Use consistent labeling
labels := map[string]string{
    pricing.StandardLabels.Method:       "GetActualCost",
    pricing.StandardLabels.Provider:     "aws", 
    pricing.StandardLabels.ResourceType: "ec2",
    pricing.StandardLabels.Region:       "us-east-1",
}

// Validate metric names and values
if err := pricing.ValidateMetricNameStrict(metricName); err != nil {
    return fmt.Errorf("invalid metric name: %w", err)
}

if err := pricing.ValidateMetricValue(value, pricing.MetricTypeCounter); err != nil {
    return fmt.Errorf("invalid metric value: %w", err)
}
```

## Service Level Indicators (SLIs)

### Core SLIs

| SLI Name | Description | Target | Unit |
|----------|-------------|--------|------|
| `availability` | Service uptime percentage | 99.9% | percentage |
| `error_rate` | Error rate percentage | <0.1% | percentage |
| `latency_p95` | 95th percentile latency | <1s | seconds |
| `latency_p99` | 99th percentile latency | <2s | seconds |
| `throughput` | Requests per second | >100 | requests_per_second |
| `data_freshness` | Cost data age | <24h | hours |

### SLI Calculation Examples

```go
func (p *Plugin) calculateAvailability(timeRange *pb.TimeRange) float64 {
    total := p.getTotalRequests(timeRange)
    successful := p.getSuccessfulRequests(timeRange)
    return pricing.CalculateAvailability(successful, total)
}

func (p *Plugin) calculateErrorRate(timeRange *pb.TimeRange) float64 {
    total := p.getTotalRequests(timeRange)
    successful := p.getSuccessfulRequests(timeRange)
    return pricing.CalculateErrorRate(total, successful)
}
```

## Distributed Tracing

### OpenTelemetry Integration

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

func (p *Plugin) GetActualCost(ctx context.Context, req *pb.GetActualCostRequest) (*pb.GetActualCostResponse, error) {
    tracer := otel.Tracer("pulumicost-plugin")
    ctx, span := tracer.Start(ctx, "GetActualCost")
    defer span.End()
    
    // Add span attributes
    span.SetAttributes(
        attribute.String("provider", req.Resource.Provider),
        attribute.String("resource_type", req.Resource.ResourceType),
        attribute.String("region", req.Resource.Region),
    )
    
    // Process request
    result, err := p.processActualCostRequest(ctx, req)
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return nil, err
    }
    
    // Add telemetry metadata to response
    if result.Telemetry == nil {
        result.Telemetry = &pb.TelemetryMetadata{}
    }
    
    result.Telemetry.TraceId = span.SpanContext().TraceID().String()
    result.Telemetry.SpanId = span.SpanContext().SpanID().String()
    result.Telemetry.ProcessingTimeMs = time.Since(startTime).Milliseconds()
    
    return result, nil
}
```

### Trace Context Propagation

```go
// Extract trace context from incoming gRPC metadata
func extractTraceContext(ctx context.Context) context.Context {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return ctx
    }
    
    return otel.GetTextMapPropagator().Extract(ctx, &metadataCarrier{md: md})
}

// Inject trace context into outgoing requests
func injectTraceContext(ctx context.Context) context.Context {
    md := metadata.New(nil)
    otel.GetTextMapPropagator().Inject(ctx, &metadataCarrier{md: md})
    return metadata.NewOutgoingContext(ctx, md)
}
```

## Structured Logging

### Log Entry Format

```go
type LogEntry struct {
    Timestamp    time.Time         `json:"timestamp"`
    Level        string            `json:"level"`
    Message      string            `json:"message"`
    Component    string            `json:"component"`
    TraceID      string            `json:"trace_id,omitempty"`
    SpanID       string            `json:"span_id,omitempty"`
    Fields       map[string]string `json:"fields,omitempty"`
    ErrorDetails *ErrorDetails     `json:"error_details,omitempty"`
}
```

### Logging Implementation

```go
func (p *Plugin) logWithContext(ctx context.Context, level pricing.LogLevel, component, message string, fields map[string]string) {
    entry := &pb.LogEntry{
        Timestamp: timestamppb.Now(),
        Level:     string(level),
        Message:   message,
        Component: component,
        Fields:    fields,
    }
    
    // Extract trace context
    span := trace.SpanFromContext(ctx)
    if span.SpanContext().IsValid() {
        entry.TraceId = span.SpanContext().TraceID().String()
        entry.SpanId = span.SpanContext().SpanID().String()
    }
    
    // Validate and emit log
    if suite := pricing.ValidateLogEntry(entry.Level, entry.Message, entry.Component, entry.TraceId, entry.SpanId, entry.Fields); !suite.IsValid() {
        // Handle validation errors
        return
    }
    
    p.logger.Log(entry)
}

// Usage example
p.logWithContext(ctx, pricing.LogLevelInfo, "cost-service", "Processing cost query", map[string]string{
    "method":        "GetActualCost",
    "provider":      req.Resource.Provider,
    "resource_type": req.Resource.ResourceType,
    "request_id":    requestID,
})
```

## Testing Observability

### Unit Testing

```go
func TestObservabilityMetrics(t *testing.T) {
    plugin := NewTestPlugin()
    
    // Test metrics collection
    suite := plugintesting.NewObservabilityTestSuite(plugin, t)
    if !suite.RunBasicObservabilityTests() {
        t.Fatal("Basic observability tests failed")
    }
}

func TestHealthCheck(t *testing.T) {
    plugin := NewTestPlugin()
    
    resp, err := plugin.HealthCheck(context.Background(), &pb.HealthCheckRequest{})
    if err != nil {
        t.Fatalf("Health check failed: %v", err)
    }
    
    if resp.Status != pb.HealthCheckResponse_SERVING {
        t.Errorf("Expected SERVING status, got %v", resp.Status)
    }
}
```

### Integration Testing

```go
func TestObservabilityConformance(t *testing.T) {
    plugin := NewTestPlugin()
    
    // Test conformance level
    level := pricing.ConformanceStandard
    requirements := pricing.GetObservabilityRequirements(level)
    
    // Verify required metrics
    metrics, err := plugin.GetMetrics(context.Background(), &pb.GetMetricsRequest{})
    if err != nil {
        t.Fatalf("Failed to get metrics: %v", err)
    }
    
    foundMetrics := make(map[string]bool)
    for _, metric := range metrics.Metrics {
        foundMetrics[metric.Name] = true
    }
    
    for _, required := range requirements.RequiredMetrics {
        if !foundMetrics[required] {
            t.Errorf("Required metric '%s' not found", required)
        }
    }
}
```

## Dashboard Templates

### Prometheus/Grafana Dashboard

```json
{
  "dashboard": {
    "title": "PulumiCost Plugin Observability",
    "panels": [
      {
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(pulumicost_requests_total[5m])",
            "legendFormat": "{{method}} - {{provider}}"
          }
        ]
      },
      {
        "title": "Error Rate",
        "type": "singlestat",
        "targets": [
          {
            "expr": "rate(pulumicost_errors_total[5m]) / rate(pulumicost_requests_total[5m]) * 100"
          }
        ]
      },
      {
        "title": "Latency Distribution",
        "type": "heatmap",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(pulumicost_request_duration_seconds_bucket[5m]))"
          }
        ]
      }
    ]
  }
}
```

### CloudWatch Dashboard

```yaml
Resources:
  PulumiCostDashboard:
    Type: AWS::CloudWatch::Dashboard
    Properties:
      DashboardName: PulumiCost-Plugin-Observability
      DashboardBody: !Sub |
        {
          "widgets": [
            {
              "type": "metric",
              "properties": {
                "metrics": [
                  ["AWS/ApplicationELB", "RequestCount"],
                  ["AWS/ApplicationELB", "HTTPCode_Target_2XX_Count"],
                  ["AWS/ApplicationELB", "HTTPCode_Target_4XX_Count"],
                  ["AWS/ApplicationELB", "HTTPCode_Target_5XX_Count"]
                ],
                "period": 300,
                "stat": "Sum",
                "region": "us-east-1",
                "title": "Request Metrics"
              }
            }
          ]
        }
```

## Best Practices

### 1. Metric Naming and Labeling

- Use consistent metric naming: `pulumicost_<component>_<metric>_<unit>`
- Apply standard labels: `method`, `provider`, `resource_type`, `region`
- Avoid high-cardinality labels (e.g., user IDs, request IDs)
- Validate metric names and label values

### 2. Error Handling and Logging

- Use structured logging with consistent fields
- Include trace context in all log entries
- Classify errors by category: `network`, `auth`, `data`, `validation`
- Provide actionable error messages

### 3. Performance Considerations

- Minimize observability overhead (<5% of request time)
- Use sampling for high-volume traces
- Batch metric updates where possible
- Cache health check results for short periods

### 4. Security

- Sanitize sensitive data from logs and metrics
- Don't include credentials or PII in telemetry
- Use correlation IDs instead of user identifiers
- Implement proper access controls for observability endpoints

## Troubleshooting

### Common Issues

1. **High Cardinality Metrics**
   - **Problem**: Too many unique label combinations
   - **Solution**: Reduce label diversity, use sampling

2. **Missing Trace Context**
   - **Problem**: Traces not connected across services
   - **Solution**: Ensure proper context propagation

3. **Slow Health Checks**
   - **Problem**: Health checks timing out
   - **Solution**: Use shorter timeouts, cache results

4. **Log Volume Issues**
   - **Problem**: Excessive logging impacting performance
   - **Solution**: Adjust log levels, use sampling

### Debugging Steps

1. **Verify Observability Endpoints**

   ```bash
   # Test health check
   grpcurl -plaintext localhost:8080 pulumicost.v1.ObservabilityService/HealthCheck
   
   # Test metrics collection
   grpcurl -plaintext localhost:8080 pulumicost.v1.ObservabilityService/GetMetrics
   ```

2. **Check Metric Values**

   ```go
   // Validate metric implementation
   if err := pricing.ValidateMetricNameStrict(metricName); err != nil {
       log.Printf("Invalid metric name: %v", err)
   }
   ```

3. **Trace Validation**

   ```go
   // Verify trace context
   suite := pricing.ValidateObservabilityMetadata(traceID, spanID, requestID, processingTime, qualityScore)
   if !suite.IsValid() {
       log.Printf("Invalid telemetry: %v", suite.Errors)
   }
   ```

### Performance Monitoring

Monitor observability overhead:

```go
// Measure observability impact
start := time.Now()
// ... business logic ...
businessDuration := time.Since(start)

observabilityStart := time.Now()
// ... observability collection ...
observabilityDuration := time.Since(observabilityStart)

overhead := float64(observabilityDuration) / float64(businessDuration) * 100
if overhead > 5.0 {
    log.Printf("High observability overhead: %.2f%%", overhead)
}
```

This guide provides a comprehensive foundation for implementing observability in PulumiCost plugins. Follow the
conformance levels appropriate for your deployment environment and continuously monitor the effectiveness of your
observability implementation.
