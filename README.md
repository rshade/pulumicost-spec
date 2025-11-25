# PulumiCost Specification v0.1.0

## Production-ready specification for cloud cost source plugins

The PulumiCost specification provides a comprehensive, standardized protocol for cost source
plugins to integrate with cost management platforms. This specification enables consistent
cost data retrieval across AWS, Azure, GCP, Kubernetes, and custom providers.

[![CI](https://github.com/rshade/pulumicost-spec/actions/workflows/ci.yml/badge.svg)](https://github.com/rshade/pulumicost-spec/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/rshade/pulumicost-spec/branch/main/graph/badge.svg)](https://codecov.io/gh/rshade/pulumicost-spec)
[![Go Reference](https://pkg.go.dev/badge/github.com/rshade/pulumicost-spec/sdk/go.svg)](https://pkg.go.dev/github.com/rshade/pulumicost-spec/sdk/go)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## Overview

PulumiCost Specification v0.1.0 is a complete, enterprise-ready protocol for
standardizing cloud cost data retrieval. It provides:

### Core Features

- **Universal Protocol**: Standardized gRPC interface for cost plugins
- **Comprehensive Schema**: JSON Schema supporting 44+ billing models across all major cloud providers
- **Multi-Provider Support**: Native support for AWS, Azure, GCP, Kubernetes, and custom providers
- **Production SDK**: Complete Go SDK with automatic protobuf generation
- **Enterprise Testing**: Multi-level conformance testing framework (Basic, Standard, Advanced)
- **CI/CD Ready**: Complete GitHub Actions pipeline with validation and benchmarks

### Supported Billing Models

- **Time-based**: per_hour, per_minute, per_second, per_cpu_hour, per_memory_gb_hour
- **Storage**: per_gb_month, per_gb_hour, per_gb_day, per_iops, per_provisioned_iops
- **Usage-based**: per_request, per_invocation, per_operation, per_transaction
- **Commitment**: reserved, spot, preemptible, savings_plan, committed_use
- **Database**: per_rcu, per_wcu, per_dtu, per_ru
- **Network**: per_data_transfer_gb, per_bandwidth_gb, per_api_call
- **And more**: [See complete list](schemas/pricing_spec.schema.json)

### Key Capabilities

- **Tiered Pricing**: Volume-based pricing tiers with automatic discounts
- **Time Aggregation**: Flexible cost aggregation rules (calendar, billing, continuous)
- **Commitment Terms**: Reserved instance and spot pricing with discount tracking
- **Rich Metadata**: Resource tags, plugin metadata, and metric hints
- **Enterprise Tags**: Billing centers, cost centers, environments, projects

## Repository Structure

```text
pulumicost-spec/
├─ proto/pulumicost/v1/           # gRPC service definitions
│  └─ costsource.proto            # Complete v0.1.0 service specification
├─ schemas/                       # JSON schema validation
│  └─ pricing_spec.schema.json    # Comprehensive v0.1.0 pricing schema
├─ sdk/go/                        # Production Go SDK
│  ├─ proto/                      # Generated protobuf bindings (auto-generated)
│  ├─ types/                      # Helper types and validation
│  └─ testing/                    # Complete testing framework
├─ examples/                      # Cross-vendor examples
│  ├─ specs/                      # 8 comprehensive pricing examples
│  └─ requests/                   # Sample gRPC request payloads
├─ .github/workflows/             # Enterprise CI/CD pipeline
│  └─ ci.yml                      # Complete validation and testing
└─ docs/                          # Comprehensive documentation
```

### Core Components

- **[gRPC Service](proto/pulumicost/v1/costsource.proto)**: Complete v0.1.0 CostSourceService with 6 RPC methods
- **[JSON Schema](schemas/pricing_spec.schema.json)**: Comprehensive validation supporting all major cloud providers
- **[Go SDK](sdk/go/)**: Production-ready SDK with automatic protobuf generation
- **[Testing Framework](sdk/go/testing/)**: Multi-level conformance testing (Basic, Standard, Advanced)
- **[Examples](examples/)**: Cross-vendor examples demonstrating all major billing models
- **[CI/CD Pipeline](.github/workflows/ci.yml)**: Complete validation, testing, and performance benchmarks

## Quick Start

### Installation

```bash
# Add SDK to your Go project
go get github.com/rshade/pulumicost-spec/sdk/go/proto
go get github.com/rshade/pulumicost-spec/sdk/go/types
```

### Development Setup

```bash
# Clone repository
git clone https://github.com/rshade/pulumicost-spec.git
cd pulumicost-spec

# Generate Go SDK (installs buf automatically)
make generate

# Run complete test suite
make test

# Validate implementation
make validate
```

### Build System

- `make generate` - Generate Go SDK from protobuf (auto-installs buf CLI)
- `make test` - Run comprehensive test suite
- `make validate` - Run validation (tests + linting + schema validation)
- `make clean` - Clean generated files
- `make tidy` - Tidy Go dependencies

## Plugin Development

### Creating a Cost Source Plugin

Implement the `CostSourceServiceServer` interface:

#### Complete Plugin Implementation

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "time"
    
    "google.golang.org/grpc"
    "google.golang.org/protobuf/types/known/timestamppb"
    
    pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

type kubecostPlugin struct {
    pbc.UnimplementedCostSourceServiceServer
}

// Name returns the plugin name
func (k *kubecostPlugin) Name(ctx context.Context, req *pbc.NameRequest) (*pbc.NameResponse, error) {
    return &pbc.NameResponse{
        Name: "kubecost",
    }, nil
}

// Supports checks if a resource type is supported
func (k *kubecostPlugin) Supports(ctx context.Context, req *pbc.SupportsRequest) (*pbc.SupportsResponse, error) {
    resource := req.GetResource()
    
    // Example: Support Kubernetes namespaces
    if resource.GetProvider() == "kubernetes" && resource.GetResourceType() == "namespace" {
        return &pbc.SupportsResponse{
            Supported: true,
            Reason:    "",
        }, nil
    }
    
    return &pbc.SupportsResponse{
        Supported: false,
        Reason:    "Only Kubernetes namespaces are supported",
    }, nil
}

// GetActualCost retrieves historical cost data
func (k *kubecostPlugin) GetActualCost(ctx context.Context, req *pbc.GetActualCostRequest) (*pbc.GetActualCostResponse, error) {
    // Example: Return mock cost data for demonstration
    results := []*pbc.ActualCostResult{
        {
            Timestamp:   timestamppb.New(time.Now().Add(-24 * time.Hour)),
            Cost:        45.67,
            UsageAmount: 12.5,
            UsageUnit:   "cpu-hours",
            Source:      "kubecost",
        },
        {
            Timestamp:   timestamppb.New(time.Now().Add(-12 * time.Hour)),
            Cost:        52.34,
            UsageAmount: 14.2,
            UsageUnit:   "cpu-hours",
            Source:      "kubecost",
        },
    }
    
    return &pbc.GetActualCostResponse{
        Results: results,
    }, nil
}

// GetProjectedCost calculates projected costs
func (k *kubecostPlugin) GetProjectedCost(ctx context.Context, req *pbc.GetProjectedCostRequest) (*pbc.GetProjectedCostResponse, error) {
    return &pbc.GetProjectedCostResponse{
        UnitPrice:     0.03,
        Currency:      "USD",
        CostPerMonth:  216.00, // 30 days * 24 hours * 0.03
        BillingDetail: "kubecost-hourly-avg",
    }, nil
}

// GetPricingSpec returns detailed pricing specification
func (k *kubecostPlugin) GetPricingSpec(ctx context.Context, req *pbc.GetPricingSpecRequest) (*pbc.GetPricingSpecResponse, error) {
    spec := &pbc.PricingSpec{
        Provider:     "kubernetes",
        ResourceType: "namespace",
        BillingMode:  "per_cpu_hour",
        RatePerUnit:  0.03,
        Currency:     "USD",
        Description:  "Kubernetes namespace CPU pricing via Kubecost",
        MetricHints: []*pbc.UsageMetricHint{
            {
                Metric: "cpu_cores",
                Unit:   "hour",
            },
            {
                Metric: "memory_gb",
                Unit:   "hour",
            },
        },
        PluginMetadata: map[string]string{
            "cluster_name":     "production-cluster",
            "namespace_labels": "app=web,tier=frontend",
            "kubecost_version": "1.98.0",
        },
        Source: "kubecost",
    }
    
    return &pbc.GetPricingSpecResponse{
        Spec: spec,
    }, nil
}

// EstimateCost estimates the monthly cost for a resource before deployment
func (k *kubecostPlugin) EstimateCost(ctx context.Context, req *pbc.EstimateCostRequest) (*pbc.EstimateCostResponse, error) {
    // Example: Estimate cost for a Kubernetes workload
    // In production, this would analyze the attributes to calculate actual cost
    return &pbc.EstimateCostResponse{
        Currency:    "USD",
        CostMonthly: 216.00, // 30 days * 24 hours * 0.03 per hour
    }, nil
}

func main() {
    // Create gRPC server
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    
    server := grpc.NewServer()
    plugin := &kubecostPlugin{}
    
    pbc.RegisterCostSourceServiceServer(server, plugin)
    
    fmt.Println("Kubecost plugin server listening on :50051")
    if err := server.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
```

#### Consuming Cost Source Plugins

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "google.golang.org/protobuf/types/known/structpb"
    "google.golang.org/protobuf/types/known/timestamppb"

    pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

func main() {
    // Connect to cost source plugin
    conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()
    
    client := pbc.NewCostSourceServiceClient(conn)
    ctx := context.Background()
    
    // Get plugin name
    nameResp, err := client.Name(ctx, &pbc.NameRequest{})
    if err != nil {
        log.Fatalf("Failed to get name: %v", err)
    }
    fmt.Printf("Connected to plugin: %s\n", nameResp.GetName())
    
    // Check if plugin supports a resource
    resource := &pbc.ResourceDescriptor{
        Provider:     "kubernetes",
        ResourceType: "namespace",
        Region:       "us-east-1",
        Tags: map[string]string{
            "app":  "web",
            "tier": "frontend",
        },
    }
    
    supportsResp, err := client.Supports(ctx, &pbc.SupportsRequest{
        Resource: resource,
    })
    if err != nil {
        log.Fatalf("Failed to check support: %v", err)
    }
    
    if !supportsResp.GetSupported() {
        fmt.Printf("Resource not supported: %s\n", supportsResp.GetReason())
        return
    }
    
    // Get actual cost data
    actualCostResp, err := client.GetActualCost(ctx, &pbc.GetActualCostRequest{
        ResourceId: "namespace/production",
        Start:      timestamppb.New(time.Now().Add(-48 * time.Hour)),
        End:        timestamppb.New(time.Now()),
        Tags: map[string]string{
            "app": "web",
        },
    })
    if err != nil {
        log.Fatalf("Failed to get actual cost: %v", err)
    }
    
    fmt.Printf("\nActual cost data:\n")
    for _, result := range actualCostResp.GetResults() {
        fmt.Printf("  Time: %s, Cost: $%.2f, Usage: %.2f %s\n",
            result.GetTimestamp().AsTime().Format(time.RFC3339),
            result.GetCost(),
            result.GetUsageAmount(),
            result.GetUsageUnit())
    }
    
    // Get projected cost
    projectedResp, err := client.GetProjectedCost(ctx, &pbc.GetProjectedCostRequest{
        Resource: resource,
    })
    if err != nil {
        log.Fatalf("Failed to get projected cost: %v", err)
    }
    
    fmt.Printf("\nProjected cost:\n")
    fmt.Printf("  Unit price: $%.4f %s\n", projectedResp.GetUnitPrice(), projectedResp.GetCurrency())
    fmt.Printf("  Monthly cost: $%.2f\n", projectedResp.GetCostPerMonth())
    fmt.Printf("  Billing detail: %s\n", projectedResp.GetBillingDetail())
    
    // Get pricing specification
    specResp, err := client.GetPricingSpec(ctx, &pbc.GetPricingSpecRequest{
        Resource: resource,
    })
    if err != nil {
        log.Fatalf("Failed to get pricing spec: %v", err)
    }
    
    spec := specResp.GetSpec()
    fmt.Printf("\nPricing specification:\n")
    fmt.Printf("  Provider: %s\n", spec.GetProvider())
    fmt.Printf("  Resource Type: %s\n", spec.GetResourceType())
    fmt.Printf("  Billing Mode: %s\n", spec.GetBillingMode())
    fmt.Printf("  Rate per unit: $%.4f %s\n", spec.GetRatePerUnit(), spec.GetCurrency())
    fmt.Printf("  Description: %s\n", spec.GetDescription())
    
    if len(spec.GetMetricHints()) > 0 {
        fmt.Printf("  Metric hints:\n")
        for _, hint := range spec.GetMetricHints() {
            fmt.Printf("    - %s (%s)\n", hint.GetMetric(), hint.GetUnit())
        }
    }

    // Estimate cost for a resource before deployment ("what-if" analysis)
    estimateResp, err := client.EstimateCost(ctx, &pbc.EstimateCostRequest{
        ResourceType: "kubernetes:core/v1:Namespace",
        Attributes: &structpb.Struct{
            Fields: map[string]*structpb.Value{
                "cpu_limit": structpb.NewStringValue("2"),
                "memory_limit": structpb.NewStringValue("4Gi"),
            },
        },
    })
    if err != nil {
        log.Fatalf("Failed to estimate cost: %v", err)
    }

    fmt.Printf("\nEstimated cost (before deployment):\n")
    fmt.Printf("  Monthly cost: $%.2f %s\n", estimateResp.GetCostMonthly(), estimateResp.GetCurrency())
}
```

### Schema Validation

#### Validating PricingSpec Documents

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/rshade/pulumicost-spec/sdk/go/types"
)

func main() {
    // Example PricingSpec JSON document
    pricingSpecJSON := `{
        "provider": "aws",
        "resource_type": "ec2",
        "sku": "t3.micro",
        "region": "us-east-1",
        "billing_mode": "per_hour",
        "rate_per_unit": 0.0104,
        "currency": "USD",
        "description": "General Purpose t3.micro instance",
        "metric_hints": [
            {
                "metric": "vcpu_hours",
                "unit": "hour",
                "aggregation_method": "sum"
            }
        ],
        "source": "aws"
    }`
    
    // Validate the JSON document against the schema
    err := types.ValidatePricingSpec([]byte(pricingSpecJSON))
    if err != nil {
        log.Fatalf("Validation failed: %v", err)
    }
    
    fmt.Println("PricingSpec JSON is valid!")
}
```

#### Working with Billing Modes

The specification supports 44+ billing models across all major cloud providers:

```go
package main

import (
    "fmt"
    
    "github.com/rshade/pulumicost-spec/sdk/go/types"
)

func main() {
    // Check if a billing mode is valid
    validModes := []string{
        "per_hour",
        "per_gb_month", 
        "per_request",
        "spot",
        "reserved",
    }
    
    for _, mode := range validModes {
        if types.IsValidBillingMode(mode) {
            fmt.Printf("✓ %s is a valid billing mode\n", mode)
        } else {
            fmt.Printf("✗ %s is not a valid billing mode\n", mode)
        }
    }
    
    // Get all available billing modes
    fmt.Printf("\nAll available billing modes:\n")
    for _, mode := range types.GetAllBillingModes() {
        fmt.Printf("  - %s\n", mode)
    }
}
```

## Testing Framework

PulumiCost Specification includes a comprehensive testing framework with multi-level conformance validation.

### Plugin Testing

```go
package main

import (
    "testing"
    plugintesting "github.com/rshade/pulumicost-spec/sdk/go/testing"
)

func TestMyPlugin(t *testing.T) {
    plugin := &MyPluginImpl{}
    
    // Run basic integration tests
    harness := plugintesting.NewTestHarness(plugin)
    harness.Start(t)
    defer harness.Stop()
    
    // Test with mock client
    client := harness.Client()
    // ... test your plugin
}
```

### Conformance Testing

Three levels of conformance testing ensure production readiness:

#### Basic Conformance (Required)

- Core functionality validation
- Error handling requirements
- Basic response validation
- **Required for plugin submission**

#### Standard Conformance (Recommended)

- Data consistency validation
- Performance requirements (response times)
- Concurrency testing (10+ requests)
- **Recommended for enterprise deployments**

#### Advanced Conformance (High-Performance)

- Scalability testing (50+ concurrent requests)
- Large dataset handling (30+ days)
- Performance benchmarks
- **Required for high-throughput environments**

```go
func TestPluginConformance(t *testing.T) {
    plugin := &MyPluginImpl{}
    
    // Choose conformance level
    result := plugintesting.RunStandardConformanceTests(t, plugin)
    plugintesting.PrintConformanceReport(result)
    
    if result.FailedTests > 0 {
        t.Errorf("Plugin failed conformance: %s", result.Summary)
    }
}
```

### Performance Testing

```bash
# Run complete test suite
make test

# Run performance benchmarks
go test -bench=. -benchmem ./sdk/go/testing/

# Run conformance tests
go test -v -run TestConformance ./sdk/go/testing/

# Run integration tests
go test -v -tags=integration ./sdk/go/testing/
```

See [Testing Framework Documentation](sdk/go/testing/README.md) for complete details.

## Examples

The specification includes 8 comprehensive examples demonstrating all major cloud providers and billing models:

### Cloud Provider Examples

#### AWS (Amazon Web Services)

- **[EC2 t3.micro](examples/specs/aws-ec2-t3-micro.json)**: On-demand instance pricing
- **[S3 Tiered Storage](examples/specs/aws-s3-tiered-pricing.json)**: Volume-based tiered pricing
- **[Lambda Functions](examples/specs/aws-lambda-per-invocation.json)**: Serverless execution pricing

#### Azure (Microsoft Azure)

- **[Virtual Machine](examples/specs/azure-vm-per-second.json)**: Per-second billing with time aggregation
- **[SQL Database](examples/specs/azure-sql-dtu.json)**: DTU-based database pricing

#### GCP (Google Cloud Platform)

- **[Cloud Storage](examples/specs/gcp-storage-standard.json)**: Tiered storage pricing
- **[Preemptible Instances](examples/specs/gcp-preemptible-spot.json)**: Spot pricing with commitment terms

#### Kubernetes

- **[Namespace CPU](examples/specs/kubernetes-namespace-cpu.json)**: CPU-based pricing via Kubecost

### Billing Model Examples

- **Time-based**: Hourly, per-second, CPU hours, memory hours
- **Storage**: GB-month, GB-day, IOPS, provisioned IOPS
- **Usage-based**: Per request, per invocation, per transaction
- **Commitment**: Reserved, spot, preemptible with discount tracking
- **Tiered**: Volume-based pricing with automatic discounts
- **Database**: RCU, WCU, DTU, RU pricing models

See [Examples Documentation](examples/README.md) for detailed explanations.

## Architecture

### gRPC Service Interface

The CostSourceService provides 6 core RPC methods:

```protobuf
service CostSourceService {
  rpc Name(NameRequest) returns (NameResponse);                    // Plugin identification
  rpc Supports(SupportsRequest) returns (SupportsResponse);        // Resource support check
  rpc GetActualCost(GetActualCostRequest) returns (GetActualCostResponse);    // Historical costs
  rpc GetProjectedCost(GetProjectedCostRequest) returns (GetProjectedCostResponse); // Cost projections
  rpc GetPricingSpec(GetPricingSpecRequest) returns (GetPricingSpecResponse);       // Pricing specifications
  rpc EstimateCost(EstimateCostRequest) returns (EstimateCostResponse);       // Cost estimation before deployment
}
```

### JSON Schema Validation

The [pricing specification schema](schemas/pricing_spec.schema.json) validates:

- **Required fields**: provider, resource_type, billing_mode, rate_per_unit, currency
- **44+ billing models**: Supporting all major cloud provider patterns
- **Advanced features**: Tiered pricing, time aggregation, commitment terms
- **Rich metadata**: Resource tags, plugin metadata, metric hints

### SDK Architecture

- **[Generated Proto](sdk/go/proto/)**: Auto-generated from protobuf definitions
- **[Helper Types](sdk/go/types/)**: Domain types, validation, billing mode constants
- **[Testing Framework](sdk/go/testing/)**: Comprehensive plugin testing suite

## CI/CD Pipeline

Complete GitHub Actions pipeline with:

- **SDK Generation**: Automatic protobuf compilation with up-to-date verification
- **Testing**: Unit tests, integration tests, conformance tests
- **Performance**: Benchmark testing with artifact upload
- **Validation**: Go linting, buf linting, JSON schema validation
- **Breaking Changes**: Automatic detection with buf
- **Coverage**: Code coverage reporting with Codecov

## Production Readiness

### Enterprise Features

- ✅ **Complete specification**: All major cloud providers and billing models
- ✅ **Multi-level testing**: Basic, Standard, Advanced conformance levels
- ✅ **Performance validation**: Response time requirements and benchmarks
- ✅ **Schema validation**: Comprehensive JSON schema with all examples validated
- ✅ **Error handling**: Standardized gRPC error codes and patterns
- ✅ **Concurrency support**: Thread-safe plugin requirements
- ✅ **Documentation**: Complete API reference and developer guides
- ✅ **CI/CD pipeline**: Automated validation and testing

### Plugin Certification

Plugins can be certified at three levels:

1. **Basic**: Core functionality, required for submission
2. **Standard**: Production-ready, recommended for enterprise
3. **Advanced**: High-performance, required for scale

## Contributing

### Development Workflow

```bash
# 1. Make changes to proto or schema files
vim proto/pulumicost/v1/costsource.proto
vim schemas/pricing_spec.schema.json

# 2. Generate SDK and validate
make generate
make validate

# 3. Run complete test suite
make test

# 4. Submit PR
git add .
git commit -m "feat: add new billing model"
git push origin feature-branch
```

### Requirements

- Run `make validate` before submitting PRs
- All examples must pass schema validation
- Breaking changes require buf validation
- New billing models require examples and tests

### Versioning

Semantic versioning for proto changes:

- **MAJOR**: Breaking proto changes
- **MINOR**: Backward-compatible additions
- **PATCH**: Bug fixes, documentation updates

Current version: **v0.1.0** (production-ready)

## License

Apache License 2.0 - see [LICENSE](LICENSE) for details.

## Support

- **Documentation**: Complete API reference and examples
- **Testing**: Multi-level conformance testing framework
- **CI/CD**: Automated validation and performance testing
- **Issues**: GitHub issues for bug reports and feature requests

---

**PulumiCost Specification v0.1.0** - Production-ready protocol for cloud cost source plugins
