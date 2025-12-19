# Quickstart: GreenOps Integration

This guide explains how to implement the new GreenOps metrics in your cost source plugin.

## 1. Advertise Supported Metrics

In your `Supports` RPC implementation, populate the `SupportedMetrics` field:

```go
func (s *server) Supports(ctx context.Context, req *pbc.SupportsRequest) (*pbc.SupportsResponse, error) {
    return &pbc.SupportsResponse{
        Supported: true,
        SupportedMetrics: []pbc.MetricKind{
            pbc.MetricKind_METRIC_KIND_CARBON_FOOTPRINT,
            pbc.MetricKind_METRIC_KIND_ENERGY_CONSUMPTION,
        },
    }, nil
}
```

## 2. Handle Utilization in Cost Projections

In your `GetProjectedCost` RPC, retrieve the utilization percentage using the `pluginsdk.GetUtilization` helper:

```go
func (s *server) GetProjectedCost(ctx context.Context, req *pbc.GetProjectedCostRequest) (*pbc.GetProjectedCostResponse, error) {
    // Use the SDK helper to handle precedence (Resource > Global > Default 0.5)
    utilization := pluginsdk.GetUtilization(req)

    // Perform sustainability calculations using 'utilization'
    // ...
    
    return &pbc.GetProjectedCostResponse{
        UnitPrice:     unitPrice,
        Currency:      "USD",
        CostPerMonth:  costPerMonth,
        ImpactMetrics: []*pbc.ImpactMetric{
            {
                Kind:  pbc.MetricKind_METRIC_KIND_CARBON_FOOTPRINT,
                Value: 100.0 * utilization,
                Unit:  "gCO2e",
            },
        },
    }, nil
}
```

## 3. Standard Units

- **Carbon**: Grams of CO2 equivalent (gCO2e)
- **Energy**: Kilowatt-hours (kWh)
- **Water**: Liters (L)
