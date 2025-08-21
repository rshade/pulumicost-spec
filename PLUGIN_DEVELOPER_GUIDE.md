# PulumiCost Plugin Developer Guide

This guide provides comprehensive instructions for developing PulumiCost plugins using the canonical protocol specification.

## Table of Contents

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Implementing CostSource gRPC Service](#implementing-costsource-grpc-service)
   - [Service Interface](#service-interface)
   - [Request/Response Messages](#requestresponse-messages)
   - [Implementation Requirements](#implementation-requirements)
4. [Packaging and Manifest Format](#packaging-and-manifest-format)
   - [Plugin Structure](#plugin-structure)
   - [Manifest Configuration](#manifest-configuration)
   - [Distribution](#distribution)
5. [Example: Minimal Plugin Implementation](#example-minimal-plugin-implementation)
   - [Project Setup](#project-setup)
   - [Complete Code Example](#complete-code-example)
   - [Building and Running](#building-and-running)
6. [Testing and Validation](#testing-and-validation)
   - [Unit Testing](#unit-testing)
   - [Integration Testing](#integration-testing)
   - [Schema Validation](#schema-validation)
7. [Best Practices and Common Patterns](#best-practices-and-common-patterns)
   - [Error Handling](#error-handling)
   - [Performance Considerations](#performance-considerations)
   - [Security Guidelines](#security-guidelines)
8. [Troubleshooting](#troubleshooting)
   - [Common Issues](#common-issues)
   - [Debug Techniques](#debug-techniques)
   - [FAQ](#faq)

## Overview

PulumiCost plugins implement the `CostSourceService` gRPC interface to provide cost data from
various sources (cloud providers, cost management tools, custom pricing models).
This guide walks through the complete plugin development lifecycle from implementation to distribution.

### Architecture

```text
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   PulumiCost    │    │   Your Plugin    │    │   Cost Source   │
│     Core        │◄──►│  (gRPC Server)   │◄──►│  (AWS/GCP/etc)  │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │
                              ▼
                       Implements
                    CostSourceService
```

## Prerequisites

Before starting plugin development, ensure you have:

- Go 1.19+ installed
- Protocol Buffers compiler (protoc)
- buf CLI tool (<https://docs.buf.build/installation>)
- Basic understanding of gRPC and protocol buffers
- Access to cost data source (API credentials, database, etc.)

## Implementing CostSource gRPC Service

The `CostSourceService` gRPC interface is the core contract your plugin must implement.
This section covers each RPC method in detail.

### Service Interface

The service defines 5 RPC methods that provide different aspects of cost information:

```protobuf
service CostSourceService {
  rpc Name(NameRequest) returns (NameResponse);
  rpc Supports(SupportsRequest) returns (SupportsResponse);
  rpc GetActualCost(GetActualCostRequest) returns (GetActualCostResponse);
  rpc GetProjectedCost(GetProjectedCostRequest) returns (GetProjectedCostResponse);
  rpc GetPricingSpec(GetPricingSpecRequest) returns (GetPricingSpecResponse);
}
```

### Request/Response Messages

#### Name RPC

Returns the display name of your plugin.

**Request**: `NameRequest` (empty)
**Response**: `NameResponse`

```protobuf
message NameResponse {
  string name = 1; // e.g., "kubecost", "cloudability", "aws-pricing"
}
```

**Implementation Notes**:

- Choose a clear, descriptive name
- Use lowercase with hyphens for consistency
- Should be unique across all plugins

#### Supports RPC

Checks if your plugin can provide cost data for a specific resource type.

**Request**: `SupportsRequest`

```protobuf
message SupportsRequest {
  ResourceDescriptor resource = 1;
}
```

**Response**: `SupportsResponse`

```protobuf
message SupportsResponse {
  bool supported = 1;
  string reason = 2; // optional explanation if not supported
}
```

**Implementation Notes**:

- Check `resource.provider`, `resource.resource_type`, and `resource.region`
- Return `false` with a descriptive reason for unsupported resources
- Consider SKU-specific support (some plugins may only support certain instance types)

#### GetActualCost RPC

Retrieves historical cost data for a resource within a time range.

**Request**: `GetActualCostRequest`

```protobuf
message GetActualCostRequest {
  string resource_id = 1;         // flexible ID format per plugin
  google.protobuf.Timestamp start = 2;
  google.protobuf.Timestamp end = 3;
  map<string, string> tags = 4;   // optional filters
}
```

**Response**: `GetActualCostResponse`

```protobuf
message GetActualCostResponse {
  repeated ActualCostResult results = 1;
}

message ActualCostResult {
  google.protobuf.Timestamp timestamp = 1;
  double cost = 2;
  double usage_amount = 3;        // optional
  string usage_unit = 4;          // e.g., "hour", "GB"
  string source = 5;              // your plugin name
}
```

**Implementation Notes**:

- `resource_id` format is plugin-specific (e.g., AWS instance ID, K8s namespace)
- Return time-series data points within the requested range
- Include usage metrics when available for better cost analysis
- Handle time zone conversion appropriately

#### GetProjectedCost RPC

Calculates projected cost for a resource based on current pricing.

**Request**: `GetProjectedCostRequest`

```protobuf
message GetProjectedCostRequest {
  ResourceDescriptor resource = 1;
}
```

**Response**: `GetProjectedCostResponse`

```protobuf
message GetProjectedCostResponse {
  double unit_price = 1;          // price per billing unit
  string currency = 2;            // e.g., "USD"
  double cost_per_month = 3;      // convenience field
  string billing_detail = 4;      // e.g., "on-demand", "reserved"
}
```

**Implementation Notes**:

- Calculate based on current pricing tables
- `cost_per_month` should assume 30.44 days (365.25/12)
- Include billing context in `billing_detail`

#### GetPricingSpec RPC

Returns detailed pricing specification for a resource type.

**Request**: `GetPricingSpecRequest`

```protobuf
message GetPricingSpecRequest {
  ResourceDescriptor resource = 1;
}
```

**Response**: `GetPricingSpecResponse`

```protobuf
message GetPricingSpecResponse {
  PricingSpec spec = 1;
}
```

The `PricingSpec` message contains comprehensive pricing details:

```protobuf
message PricingSpec {
  string provider = 1;
  string resource_type = 2;
  string sku = 3;
  string region = 4;
  string billing_mode = 5;        // per_hour, per_gb_month, etc.
  double rate_per_unit = 6;
  string currency = 7;
  string description = 8;
  repeated UsageMetricHint metric_hints = 9;
  map<string, string> plugin_metadata = 10;
  string source = 11;
}
```

**Implementation Notes**:

- Use standard `billing_mode` values: `per_hour`, `per_gb_month`, `per_request`, `flat`, `per_day`, `per_cpu_hour`
- Provide helpful `metric_hints` for usage calculation
- Include plugin-specific metadata for debugging/auditing

### Implementation Requirements

#### Error Handling

- Use standard gRPC status codes
- Provide descriptive error messages
- Handle network timeouts gracefully
- Log errors for debugging

#### Resource Descriptor Validation

```go
func validateResourceDescriptor(rd *ResourceDescriptor) error {
    if rd.Provider == "" {
        return status.Error(codes.InvalidArgument, "provider is required")
    }
    if rd.ResourceType == "" {
        return status.Error(codes.InvalidArgument, "resource_type is required")
    }
    // Add other validations
    return nil
}
```

#### Authentication

- Support API keys, OAuth tokens, or service account credentials
- Load credentials from environment variables or config files
- Implement credential refresh logic for OAuth

#### Caching

- Cache pricing data to reduce API calls
- Implement TTL-based cache invalidation
- Consider regional pricing differences

#### Logging

- Use structured logging (JSON format recommended)
- Log request/response for debugging
- Include correlation IDs for tracing

## Packaging and Manifest Format

Plugins must follow a standardized packaging format for consistent deployment and discovery.

### Plugin Structure

A PulumiCost plugin follows this directory structure:

```text
my-plugin/
├── pulumicost-plugin.yaml    # Plugin manifest (required)
├── bin/                      # Plugin binaries
│   ├── plugin-linux-amd64   # Linux x64 binary
│   ├── plugin-darwin-amd64  # macOS x64 binary
│   └── plugin-windows-amd64 # Windows x64 binary
├── README.md                 # Plugin documentation
├── LICENSE                   # License file
└── examples/                 # Usage examples (optional)
    ├── config.yaml
    └── requests/
        ├── actual-cost-request.json
        └── projected-cost-request.json
```

### Manifest Configuration

The `pulumicost-plugin.yaml` file contains plugin metadata and configuration:

```yaml
# Plugin manifest version (required)
manifest_version: "1.0"

# Plugin metadata
name: "my-cost-plugin"
version: "1.0.0"
description: "Cost plugin for Custom Cost Source"
author: "Your Name <you@example.com>"
license: "Apache-2.0"
homepage: "https://github.com/yourorg/my-plugin"

# Plugin capabilities
provider: "custom"               # Primary provider this plugin supports
resource_types:                  # List of supported resource types
  - "vm"
  - "storage"
  - "database"
regions:                        # Supported regions (empty = all)
  - "us-east-1"
  - "us-west-2"
  - "eu-west-1"

# Runtime configuration
runtime:
  protocol: "grpc"              # Currently only "grpc" supported
  executable: "bin/plugin-{{.OS}}-{{.ARCH}}"
  port_range: "5000-6000"       # Preferred port range
  
# Authentication requirements
auth:
  required: true
  methods:                      # Supported auth methods
    - "api_key"
    - "oauth2"
  env_vars:                     # Required environment variables
    - "CUSTOM_API_KEY"
    - "CUSTOM_API_ENDPOINT"

# Plugin dependencies (optional)
dependencies:
  min_go_version: "1.19"
  external_tools: []

# Health check configuration
health_check:
  enabled: true
  timeout: "30s"
  endpoint: "/health"           # HTTP endpoint for health checks
```

### Manifest Field Reference

#### Required Fields

- **`manifest_version`**: Manifest format version (currently "1.0")
- **`name`**: Unique plugin name (lowercase, hyphens allowed)
- **`version`**: Semantic version (e.g., "1.2.3")
- **`description`**: Brief description of plugin functionality
- **`provider`**: Primary cloud provider or cost source
- **`runtime.protocol`**: Communication protocol ("grpc")
- **`runtime.executable`**: Path to plugin binary with template variables

#### Optional Fields

- **`author`**: Plugin author information
- **`license`**: SPDX license identifier
- **`homepage`**: Plugin homepage URL
- **`resource_types`**: Array of supported resource types
- **`regions`**: Array of supported regions (empty = all regions)
- **`runtime.port_range`**: Preferred port range for gRPC server
- **`auth`**: Authentication configuration
- **`dependencies`**: Runtime dependencies
- **`health_check`**: Health check configuration

#### Template Variables

The `runtime.executable` field supports template variables:

- `{{.OS}}`: Target operating system (linux, darwin, windows)
- `{{.ARCH}}`: Target architecture (amd64, arm64)
- `{{.VERSION}}`: Plugin version

Example: `"bin/plugin-{{.OS}}-{{.ARCH}}-{{.VERSION}}"`

### Binary Naming Convention

Plugin binaries should follow this naming pattern:

```text
plugin-<os>-<arch>[-<version>]
```

Examples:

- `plugin-linux-amd64`
- `plugin-darwin-arm64`
- `plugin-windows-amd64-1.0.0`

### Authentication Configuration

Plugins can specify authentication requirements:

```yaml
auth:
  required: true
  methods:
    - "api_key"      # API key authentication
    - "oauth2"       # OAuth 2.0 flow
    - "jwt"          # JWT token
    - "basic"        # Basic authentication
    - "cert"         # Client certificate
  env_vars:
    - "API_KEY"      # Required environment variables
    - "API_ENDPOINT"
  config_file:       # Optional config file path
    path: "~/.config/plugin/config.yaml"
    format: "yaml"
```

### Distribution

#### Package Formats

Plugins can be distributed in multiple formats:

1. **Tarball** (`.tar.gz`)

   ```bash
   tar -czf my-plugin-1.0.0.tar.gz my-plugin/
   ```

2. **ZIP Archive** (`.zip`)

   ```bash
   zip -r my-plugin-1.0.0.zip my-plugin/
   ```

3. **Container Image** (Docker)

   ```dockerfile
   FROM scratch
   COPY plugin-linux-amd64 /plugin
   COPY pulumicost-plugin.yaml /pulumicost-plugin.yaml
   ENTRYPOINT ["/plugin"]
   ```

#### Registry Publishing

Plugins can be published to registries for easy discovery:

```bash
# Publish to plugin registry
pulumicost plugin publish my-plugin-1.0.0.tar.gz

# Install from registry
pulumicost plugin install my-cost-plugin@1.0.0
```

### Validation

Validate your plugin package before distribution:

```bash
# Validate manifest syntax
pulumicost plugin validate pulumicost-plugin.yaml

# Test plugin package
pulumicost plugin test my-plugin-1.0.0.tar.gz
```

The validation checks:

- Manifest syntax and required fields
- Binary compatibility and architecture
- Authentication configuration
- Resource type format
- Version format compliance

## Example: Minimal Plugin Implementation

This section provides a complete, working plugin implementation that demonstrates all the concepts covered in this guide.

### Project Setup

Create a new Go project for your plugin:

```bash
mkdir my-cost-plugin
cd my-cost-plugin
go mod init github.com/yourorg/my-cost-plugin

# Add dependencies
go get github.com/rshade/pulumicost-spec/sdk/go/proto@latest
go get google.golang.org/grpc@latest
go get google.golang.org/protobuf@latest
```

### Complete Code Example

#### main.go

```go
package main

import (
 "context"
 "flag"
 "fmt"
 "log"
 "net"
 "os"
 "os/signal"
 "syscall"
 "time"

 pb "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
 "github.com/rshade/pulumicost-spec/sdk/go/types"
 "google.golang.org/grpc"
 "google.golang.org/grpc/codes"
 "google.golang.org/grpc/status"
 "google.golang.org/protobuf/types/known/timestamppb"
)

var (
 port = flag.Int("port", 50051, "The server port")
)

// server implements the CostSourceService
type server struct {
 pb.UnimplementedCostSourceServiceServer
 name string
}

// Name returns the plugin name
func (s *server) Name(ctx context.Context, req *pb.NameRequest) (*pb.NameResponse, error) {
 log.Println("Name RPC called")
 return &pb.NameResponse{
  Name: s.name,
 }, nil
}

// Supports checks if the plugin supports a resource type
func (s *server) Supports(ctx context.Context, req *pb.SupportsRequest) (*pb.SupportsResponse, error) {
 log.Printf("Supports RPC called for provider=%s, resource_type=%s", 
  req.Resource.Provider, req.Resource.ResourceType)

 // This example plugin supports "custom" provider with "vm" and "storage" resources
 if req.Resource.Provider != "custom" {
  return &pb.SupportsResponse{
   Supported: false,
   Reason:    fmt.Sprintf("Provider %s not supported", req.Resource.Provider),
  }, nil
 }

 switch req.Resource.ResourceType {
 case "vm", "storage":
  return &pb.SupportsResponse{
   Supported: true,
  }, nil
 default:
  return &pb.SupportsResponse{
   Supported: false,
   Reason:    fmt.Sprintf("Resource type %s not supported", req.Resource.ResourceType),
  }, nil
 }
}

// GetActualCost retrieves historical cost data
func (s *server) GetActualCost(ctx context.Context, req *pb.GetActualCostRequest) (*pb.GetActualCostResponse, error) {
 log.Printf("GetActualCost RPC called for resource_id=%s", req.ResourceId)

 // Validate request
 if req.ResourceId == "" {
  return nil, status.Error(codes.InvalidArgument, "resource_id is required")
 }
 if req.Start == nil || req.End == nil {
  return nil, status.Error(codes.InvalidArgument, "start and end timestamps are required")
 }

 // Mock cost data - replace with actual data source integration
 results := generateMockActualCostData(req.ResourceId, req.Start.AsTime(), req.End.AsTime())

 return &pb.GetActualCostResponse{
  Results: results,
 }, nil
}

// GetProjectedCost calculates projected costs
func (s *server) GetProjectedCost(ctx context.Context, req *pb.GetProjectedCostRequest) (*pb.GetProjectedCostResponse, error) {
 log.Printf("GetProjectedCost RPC called for provider=%s, resource_type=%s, sku=%s", 
  req.Resource.Provider, req.Resource.ResourceType, req.Resource.Sku)

 // Validate resource
 if err := validateResourceDescriptor(req.Resource); err != nil {
  return nil, err
 }

 // Calculate projected cost based on resource type
 var unitPrice float64
 var billingDetail string

 switch req.Resource.ResourceType {
 case "vm":
  unitPrice = 0.05 // $0.05 per hour
  billingDetail = "on-demand hourly"
 case "storage":
  unitPrice = 0.10 // $0.10 per GB per month
  billingDetail = "standard storage monthly"
 default:
  return nil, status.Error(codes.InvalidArgument, "unsupported resource type")
 }

 // Calculate cost per month (assume 30.44 days)
 hoursPerMonth := 24 * 30.44
 var costPerMonth float64
 if req.Resource.ResourceType == "vm" {
  costPerMonth = unitPrice * hoursPerMonth
 } else {
  costPerMonth = unitPrice // Already monthly rate
 }

 return &pb.GetProjectedCostResponse{
  UnitPrice:     unitPrice,
  Currency:      "USD",
  CostPerMonth:  costPerMonth,
  BillingDetail: billingDetail,
 }, nil
}

// GetPricingSpec returns detailed pricing specification
func (s *server) GetPricingSpec(ctx context.Context, req *pb.GetPricingSpecRequest) (*pb.GetPricingSpecResponse, error) {
 log.Printf("GetPricingSpec RPC called for provider=%s, resource_type=%s", 
  req.Resource.Provider, req.Resource.ResourceType)

 // Validate resource
 if err := validateResourceDescriptor(req.Resource); err != nil {
  return nil, err
 }

 // Build pricing spec based on resource type
 spec := &pb.PricingSpec{
  Provider:     req.Resource.Provider,
  ResourceType: req.Resource.ResourceType,
  Sku:          req.Resource.Sku,
  Region:       req.Resource.Region,
  Currency:     "USD",
  Description:  fmt.Sprintf("Pricing for %s %s", req.Resource.Provider, req.Resource.ResourceType),
  Source:       s.name,
  PluginMetadata: map[string]string{
   "plugin_version": "1.0.0",
   "last_updated":   time.Now().Format(time.RFC3339),
  },
 }

 switch req.Resource.ResourceType {
 case "vm":
  spec.BillingMode = string(types.PerHour)
  spec.RatePerUnit = 0.05
  spec.MetricHints = []*pb.UsageMetricHint{
   {
    Metric: "vcpu_hours",
    Unit:   "hour",
   },
  }
 case "storage":
  spec.BillingMode = string(types.PerGBMonth)
  spec.RatePerUnit = 0.10
  spec.MetricHints = []*pb.UsageMetricHint{
   {
    Metric: "storage_gb_month",
    Unit:   "GB",
   },
  }
 default:
  return nil, status.Error(codes.InvalidArgument, "unsupported resource type")
 }

 return &pb.GetPricingSpecResponse{
  Spec: spec,
 }, nil
}

// validateResourceDescriptor validates the resource descriptor
func validateResourceDescriptor(rd *pb.ResourceDescriptor) error {
 if rd == nil {
  return status.Error(codes.InvalidArgument, "resource descriptor is required")
 }
 if rd.Provider == "" {
  return status.Error(codes.InvalidArgument, "provider is required")
 }
 if rd.ResourceType == "" {
  return status.Error(codes.InvalidArgument, "resource_type is required")
 }
 return nil
}

// generateMockActualCostData creates sample historical cost data
func generateMockActualCostData(resourceID string, start, end time.Time) []*pb.ActualCostResult {
 var results []*pb.ActualCostResult
 
 // Generate daily cost data points
 current := start
 for current.Before(end) {
  // Mock cost calculation - replace with actual data source
  baseCost := 2.40 // $2.40 per day base cost
  variance := 0.20 * (0.5 - float64(current.Unix()%1000)/1000) // Add some variance
  dailyCost := baseCost + variance

  results = append(results, &pb.ActualCostResult{
   Timestamp:   timestamppb.New(current),
   Cost:        dailyCost,
   UsageAmount: 24.0, // 24 hours
   UsageUnit:   "hour",
   Source:      "my-cost-plugin",
  })

  current = current.AddDate(0, 0, 1) // Next day
 }

 return results
}

func main() {
 flag.Parse()

 // Create gRPC server
 lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
 if err != nil {
  log.Fatalf("Failed to listen: %v", err)
 }

 s := grpc.NewServer()
 pb.RegisterCostSourceServiceServer(s, &server{
  name: "my-cost-plugin",
 })

 log.Printf("Server listening at %v", lis.Addr())

 // Handle graceful shutdown
 c := make(chan os.Signal, 1)
 signal.Notify(c, os.Interrupt, syscall.SIGTERM)

 go func() {
  <-c
  log.Println("Shutting down gRPC server...")
  s.GracefulStop()
 }()

 // Start server
 if err := s.Serve(lis); err != nil {
  log.Fatalf("Failed to serve: %v", err)
 }
}
```

#### go.mod

```go
module github.com/yourorg/my-cost-plugin

go 1.19

require (
 github.com/rshade/pulumicost-spec/sdk/go/proto v0.1.0
 google.golang.org/grpc v1.58.3
 google.golang.org/protobuf v1.31.0
)

require (
 github.com/golang/protobuf v1.5.3 // indirect
 golang.org/x/net v0.12.0 // indirect
 golang.org/x/sys v0.10.0 // indirect
 golang.org/x/text v0.11.0 // indirect
 google.golang.org/genproto/googleapis/rpc v0.0.0-20230711160842-782d3b101e98 // indirect
)
```

#### pulumicost-plugin.yaml

```yaml
manifest_version: "1.0"

name: "my-cost-plugin"
version: "1.0.0"
description: "Example cost plugin demonstrating PulumiCost plugin development"
author: "Your Name <you@example.com>"
license: "Apache-2.0"
homepage: "https://github.com/yourorg/my-cost-plugin"

provider: "custom"
resource_types:
  - "vm"
  - "storage"

runtime:
  protocol: "grpc"
  executable: "bin/plugin-{{.OS}}-{{.ARCH}}"
  port_range: "5000-6000"

auth:
  required: false
  methods: []
  env_vars: []

dependencies:
  min_go_version: "1.19"

health_check:
  enabled: false
```

### Building and Running

#### Building the Plugin

```bash
# Build for current platform
go build -o bin/plugin .

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o bin/plugin-linux-amd64 .
GOOS=darwin GOARCH=amd64 go build -o bin/plugin-darwin-amd64 .
GOOS=windows GOARCH=amd64 go build -o bin/plugin-windows-amd64.exe .
```

#### Running the Plugin

```bash
# Run with default port (50051)
./bin/plugin

# Run with custom port
./bin/plugin -port 8080
```

#### Testing the Plugin

You can test the plugin using grpcurl or any gRPC client:

```bash
# Install grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Test Name RPC
grpcurl -plaintext -import-path . -proto proto/pulumicost/v1/costsource.proto \
  localhost:50051 pulumicost.v1.CostSourceService.Name

# Test Supports RPC
grpcurl -plaintext -import-path . -proto proto/pulumicost/v1/costsource.proto \
  -d '{"resource": {"provider": "custom", "resource_type": "vm"}}' \
  localhost:50051 pulumicost.v1.CostSourceService.Supports

# Test GetProjectedCost RPC
grpcurl -plaintext -import-path . -proto proto/pulumicost/v1/costsource.proto \
  -d '{"resource": {"provider": "custom", "resource_type": "vm", "sku": "small", "region": "us-east-1"}}' \
  localhost:50051 pulumicost.v1.CostSourceService.GetProjectedCost
```

### Plugin Package Structure

```text
my-cost-plugin-1.0.0/
├── pulumicost-plugin.yaml
├── bin/
│   ├── plugin-linux-amd64
│   ├── plugin-darwin-amd64
│   └── plugin-windows-amd64.exe
├── README.md
├── LICENSE
└── examples/
    └── requests/
        ├── supports-request.json
        ├── actual-cost-request.json
        └── projected-cost-request.json
```

### Key Implementation Details

1. **Error Handling**: Uses gRPC status codes for structured error reporting
2. **Validation**: Validates all input parameters before processing
3. **Logging**: Includes structured logging for debugging and monitoring
4. **Mock Data**: Demonstrates how to generate sample cost data (replace with real integration)
5. **Graceful Shutdown**: Handles SIGTERM/SIGINT for clean shutdowns
6. **Resource Support**: Shows how to implement resource type filtering
7. **Pricing Logic**: Demonstrates different billing models (hourly vs monthly)

## Testing and Validation

Comprehensive testing ensures your plugin works correctly and integrates smoothly with the PulumiCost ecosystem.

### Unit Testing

Create unit tests for each RPC method to verify functionality and edge cases.

#### test/plugin_test.go

```go
package test

import (
 "context"
 "testing"
 "time"

 pb "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
 "github.com/stretchr/testify/assert"
 "google.golang.org/grpc/codes"
 "google.golang.org/grpc/status"
 "google.golang.org/protobuf/types/known/timestamppb"
)

// Import your plugin server implementation
// import "github.com/yourorg/my-cost-plugin"

func TestName(t *testing.T) {
 server := &server{name: "test-plugin"}
 
 resp, err := server.Name(context.Background(), &pb.NameRequest{})
 
 assert.NoError(t, err)
 assert.NotNil(t, resp)
 assert.Equal(t, "test-plugin", resp.Name)
}

func TestSupports(t *testing.T) {
 server := &server{name: "test-plugin"}
 
 tests := []struct {
  name       string
  resource   *pb.ResourceDescriptor
  wantSupported bool
  wantReason    string
 }{
  {
   name: "supported vm resource",
   resource: &pb.ResourceDescriptor{
    Provider:     "custom",
    ResourceType: "vm",
   },
   wantSupported: true,
   wantReason:    "",
  },
  {
   name: "unsupported provider",
   resource: &pb.ResourceDescriptor{
    Provider:     "aws",
    ResourceType: "ec2",
   },
   wantSupported: false,
   wantReason:    "Provider aws not supported",
  },
  {
   name: "unsupported resource type",
   resource: &pb.ResourceDescriptor{
    Provider:     "custom",
    ResourceType: "database",
   },
   wantSupported: false,
   wantReason:    "Resource type database not supported",
  },
 }

 for _, tt := range tests {
  t.Run(tt.name, func(t *testing.T) {
   req := &pb.SupportsRequest{Resource: tt.resource}
   resp, err := server.Supports(context.Background(), req)
   
   assert.NoError(t, err)
   assert.NotNil(t, resp)
   assert.Equal(t, tt.wantSupported, resp.Supported)
   if tt.wantReason != "" {
    assert.Equal(t, tt.wantReason, resp.Reason)
   }
  })
 }
}

func TestGetActualCost(t *testing.T) {
 server := &server{name: "test-plugin"}
 
 tests := []struct {
  name      string
  request   *pb.GetActualCostRequest
  wantError bool
  errorCode codes.Code
 }{
  {
   name: "valid request",
   request: &pb.GetActualCostRequest{
    ResourceId: "test-resource-123",
    Start:      timestamppb.New(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
    End:        timestamppb.New(time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)),
   },
   wantError: false,
  },
  {
   name: "missing resource id",
   request: &pb.GetActualCostRequest{
    Start: timestamppb.New(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
    End:   timestamppb.New(time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)),
   },
   wantError: true,
   errorCode: codes.InvalidArgument,
  },
  {
   name: "missing timestamps",
   request: &pb.GetActualCostRequest{
    ResourceId: "test-resource-123",
   },
   wantError: true,
   errorCode: codes.InvalidArgument,
  },
 }

 for _, tt := range tests {
  t.Run(tt.name, func(t *testing.T) {
   resp, err := server.GetActualCost(context.Background(), tt.request)
   
   if tt.wantError {
    assert.Error(t, err)
    assert.Nil(t, resp)
    
    st, ok := status.FromError(err)
    assert.True(t, ok)
    assert.Equal(t, tt.errorCode, st.Code())
   } else {
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.NotEmpty(t, resp.Results)
    
    // Verify result structure
    for _, result := range resp.Results {
     assert.Greater(t, result.Cost, 0.0)
     assert.NotEmpty(t, result.Source)
     assert.NotNil(t, result.Timestamp)
    }
   }
  })
 }
}

func TestGetProjectedCost(t *testing.T) {
 server := &server{name: "test-plugin"}
 
 tests := []struct {
  name         string
  resource     *pb.ResourceDescriptor
  wantError    bool
  expectedUnit float64
 }{
  {
   name: "vm resource",
   resource: &pb.ResourceDescriptor{
    Provider:     "custom",
    ResourceType: "vm",
    Sku:          "small",
    Region:       "us-east-1",
   },
   wantError:    false,
   expectedUnit: 0.05,
  },
  {
   name: "storage resource",
   resource: &pb.ResourceDescriptor{
    Provider:     "custom",
    ResourceType: "storage",
    Sku:          "standard",
    Region:       "us-east-1",
   },
   wantError:    false,
   expectedUnit: 0.10,
  },
  {
   name: "unsupported resource",
   resource: &pb.ResourceDescriptor{
    Provider:     "custom",
    ResourceType: "database",
   },
   wantError: true,
  },
 }

 for _, tt := range tests {
  t.Run(tt.name, func(t *testing.T) {
   req := &pb.GetProjectedCostRequest{Resource: tt.resource}
   resp, err := server.GetProjectedCost(context.Background(), req)
   
   if tt.wantError {
    assert.Error(t, err)
    assert.Nil(t, resp)
   } else {
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, tt.expectedUnit, resp.UnitPrice)
    assert.Equal(t, "USD", resp.Currency)
    assert.Greater(t, resp.CostPerMonth, 0.0)
    assert.NotEmpty(t, resp.BillingDetail)
   }
  })
 }
}

func TestGetPricingSpec(t *testing.T) {
 server := &server{name: "test-plugin"}
 
 req := &pb.GetPricingSpecRequest{
  Resource: &pb.ResourceDescriptor{
   Provider:     "custom",
   ResourceType: "vm",
   Sku:          "small",
   Region:       "us-east-1",
  },
 }
 
 resp, err := server.GetPricingSpec(context.Background(), req)
 
 assert.NoError(t, err)
 assert.NotNil(t, resp)
 assert.NotNil(t, resp.Spec)
 
 spec := resp.Spec
 assert.Equal(t, "custom", spec.Provider)
 assert.Equal(t, "vm", spec.ResourceType)
 assert.Equal(t, "per_hour", spec.BillingMode)
 assert.Equal(t, 0.05, spec.RatePerUnit)
 assert.Equal(t, "USD", spec.Currency)
 assert.NotEmpty(t, spec.Description)
 assert.NotEmpty(t, spec.MetricHints)
 assert.Contains(t, spec.PluginMetadata, "plugin_version")
}

// Benchmark tests
func BenchmarkGetProjectedCost(b *testing.B) {
 server := &server{name: "test-plugin"}
 req := &pb.GetProjectedCostRequest{
  Resource: &pb.ResourceDescriptor{
   Provider:     "custom",
   ResourceType: "vm",
   Sku:          "small",
   Region:       "us-east-1",
  },
 }
 
 b.ResetTimer()
 for i := 0; i < b.N; i++ {
  _, err := server.GetProjectedCost(context.Background(), req)
  if err != nil {
   b.Fatal(err)
  }
 }
}
```

### Integration Testing

Test your plugin against the actual gRPC protocol using a client.

#### test/integration_test.go

```go
package test

import (
 "context"
 "log"
 "net"
 "testing"
 "time"

 pb "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
 "github.com/stretchr/testify/assert"
 "github.com/stretchr/testify/require"
 "google.golang.org/grpc"
 "google.golang.org/grpc/credentials/insecure"
 "google.golang.org/protobuf/types/known/timestamppb"
)

func TestIntegration(t *testing.T) {
 // Start server
 lis, err := net.Listen("tcp", ":0")
 require.NoError(t, err)
 
 s := grpc.NewServer()
 pb.RegisterCostSourceServiceServer(s, &server{name: "integration-test-plugin"})
 
 go func() {
  if err := s.Serve(lis); err != nil {
   log.Printf("Server error: %v", err)
  }
 }()
 defer s.Stop()
 
 // Connect client
 conn, err := grpc.NewClient(
  lis.Addr().String(),
  grpc.WithTransportCredentials(insecure.NewCredentials()),
 )
 require.NoError(t, err)
 defer conn.Close()
 
 client := pb.NewCostSourceServiceClient(conn)
 ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
 defer cancel()
 
 // Test Name RPC
 nameResp, err := client.Name(ctx, &pb.NameRequest{})
 require.NoError(t, err)
 assert.Equal(t, "integration-test-plugin", nameResp.Name)
 
 // Test Supports RPC
 supportsResp, err := client.Supports(ctx, &pb.SupportsRequest{
  Resource: &pb.ResourceDescriptor{
   Provider:     "custom",
   ResourceType: "vm",
  },
 })
 require.NoError(t, err)
 assert.True(t, supportsResp.Supported)
 
 // Test GetProjectedCost RPC
 projectedResp, err := client.GetProjectedCost(ctx, &pb.GetProjectedCostRequest{
  Resource: &pb.ResourceDescriptor{
   Provider:     "custom",
   ResourceType: "vm",
   Sku:          "small",
   Region:       "us-east-1",
  },
 })
 require.NoError(t, err)
 assert.Greater(t, projectedResp.UnitPrice, 0.0)
 assert.Equal(t, "USD", projectedResp.Currency)
 
 // Test GetActualCost RPC
 actualResp, err := client.GetActualCost(ctx, &pb.GetActualCostRequest{
  ResourceId: "test-vm-123",
  Start:      timestamppb.New(time.Now().AddDate(0, -1, 0)),
  End:        timestamppb.New(time.Now()),
 })
 require.NoError(t, err)
 assert.NotEmpty(t, actualResp.Results)
 
 // Test GetPricingSpec RPC
 specResp, err := client.GetPricingSpec(ctx, &pb.GetPricingSpecRequest{
  Resource: &pb.ResourceDescriptor{
   Provider:     "custom",
   ResourceType: "vm",
  },
 })
 require.NoError(t, err)
 assert.NotNil(t, specResp.Spec)
 assert.Equal(t, "custom", specResp.Spec.Provider)
}
```

### Schema Validation

Validate that your plugin generates valid PricingSpec documents.

#### test/schema_test.go

```go
package test

import (
 "encoding/json"
 "testing"

 "github.com/rshade/pulumicost-spec/sdk/go/types"
 "github.com/stretchr/testify/assert"
 "github.com/stretchr/testify/require"
)

func TestPricingSpecValidation(t *testing.T) {
 tests := []struct {
  name    string
  spec    map[string]interface{}
  wantErr bool
 }{
  {
   name: "valid vm spec",
   spec: map[string]interface{}{
    "provider":      "custom",
    "resource_type": "vm",
    "billing_mode":  "per_hour",
    "rate_per_unit": 0.05,
    "currency":      "USD",
    "description":   "Test VM pricing",
   },
   wantErr: false,
  },
  {
   name: "missing required fields",
   spec: map[string]interface{}{
    "provider":      "custom",
    "resource_type": "vm",
    // missing billing_mode, rate_per_unit, currency
   },
   wantErr: true,
  },
  {
   name: "invalid billing mode",
   spec: map[string]interface{}{
    "provider":      "custom",
    "resource_type": "vm",
    "billing_mode":  "invalid_mode",
    "rate_per_unit": 0.05,
    "currency":      "USD",
   },
   wantErr: true,
  },
  {
   name: "negative rate",
   spec: map[string]interface{}{
    "provider":      "custom",
    "resource_type": "vm",
    "billing_mode":  "per_hour",
    "rate_per_unit": -0.05,
    "currency":      "USD",
   },
   wantErr: true,
  },
 }

 for _, tt := range tests {
  t.Run(tt.name, func(t *testing.T) {
   // Convert to JSON bytes
   data, err := json.Marshal(tt.spec)
   require.NoError(t, err)
   
   // Validate using the schema
   err = types.ValidatePricingSpec(data)
   
   if tt.wantErr {
    assert.Error(t, err)
   } else {
    assert.NoError(t, err)
   }
  })
 }
}

func TestBillingModeValidation(t *testing.T) {
 validModes := []string{
  "per_hour",
  "per_gb_month",
  "per_request",
  "flat",
  "per_day",
  "per_cpu_hour",
 }
 
 for _, mode := range validModes {
  t.Run(mode, func(t *testing.T) {
   assert.True(t, types.ValidBillingMode(mode))
  })
 }
 
 invalidModes := []string{
  "per_second",
  "hourly",
  "",
  "invalid",
 }
 
 for _, mode := range invalidModes {
  t.Run(mode, func(t *testing.T) {
   assert.False(t, types.ValidBillingMode(mode))
  })
 }
}
```

### Load Testing

Test your plugin under realistic load conditions.

#### test/load_test.go

```go
package test

import (
 "context"
 "sync"
 "testing"
 "time"

 pb "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
 "github.com/stretchr/testify/assert"
)

func TestConcurrentRequests(t *testing.T) {
 server := &server{name: "load-test-plugin"}
 
 const numGoroutines = 50
 const requestsPerGoroutine = 20
 
 var wg sync.WaitGroup
 errorChan := make(chan error, numGoroutines*requestsPerGoroutine)
 
 for i := 0; i < numGoroutines; i++ {
  wg.Add(1)
  go func(workerID int) {
   defer wg.Done()
   
   for j := 0; j < requestsPerGoroutine; j++ {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    
    _, err := server.GetProjectedCost(ctx, &pb.GetProjectedCostRequest{
     Resource: &pb.ResourceDescriptor{
      Provider:     "custom",
      ResourceType: "vm",
      Sku:          "small",
      Region:       "us-east-1",
     },
    })
    
    cancel()
    
    if err != nil {
     errorChan <- err
    }
   }
  }(i)
 }
 
 wg.Wait()
 close(errorChan)
 
 // Check for errors
 var errors []error
 for err := range errorChan {
  errors = append(errors, err)
 }
 
 assert.Empty(t, errors, "Expected no errors during load test, got: %v", errors)
}
```

### Running Tests

#### Test Commands

```bash
# Run all tests
go test ./test/...

# Run tests with verbose output
go test -v ./test/...

# Run only unit tests
go test ./test/ -run TestName

# Run integration tests
go test ./test/ -run TestIntegration

# Run with coverage
go test -cover ./test/...

# Generate coverage report
go test -coverprofile=coverage.out ./test/...
go tool cover -html=coverage.out -o coverage.html

# Run benchmarks
go test -bench=. ./test/...

# Run load tests
go test ./test/ -run TestConcurrentRequests
```

#### Continuous Integration

Create a CI configuration to run tests automatically:

```yaml
# .github/workflows/test.yml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.19
    
    - name: Install dependencies
      run: go mod tidy
    
    - name: Run tests
      run: go test -race -cover ./test/...
    
    - name: Run benchmarks
      run: go test -bench=. ./test/...
```

### Test Data Management

#### Creating Test Fixtures

```go
// test/fixtures.go
package test

import (
 "time"
 pb "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
 "google.golang.org/protobuf/types/known/timestamppb"
)

func CreateTestResourceDescriptor() *pb.ResourceDescriptor {
 return &pb.ResourceDescriptor{
  Provider:     "custom",
  ResourceType: "vm",
  Sku:          "small",
  Region:       "us-east-1",
  Tags: map[string]string{
   "environment": "test",
   "team":        "platform",
  },
 }
}

func CreateTestActualCostRequest() *pb.GetActualCostRequest {
 return &pb.GetActualCostRequest{
  ResourceId: "test-vm-123",
  Start:      timestamppb.New(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
  End:        timestamppb.New(time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)),
  Tags: map[string]string{
   "app": "web",
   "env": "production",
  },
 }
}
```

### Testing Best Practices

1. **Test Coverage**: Aim for >90% code coverage
2. **Edge Cases**: Test boundary conditions and error scenarios
3. **Mock Dependencies**: Use mocks for external API calls
4. **Test Isolation**: Each test should be independent
5. **Performance**: Include benchmark and load tests
6. **CI/CD**: Run tests automatically on every commit
7. **Documentation**: Document test scenarios and expectations

## Best Practices and Common Patterns

Following established patterns ensures your plugin integrates well with the PulumiCost
ecosystem and provides a consistent experience.

### Error Response Handling

#### Use Standard gRPC Status Codes

```go
import (
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func (s *server) GetProjectedCost(ctx context.Context, req *pb.GetProjectedCostRequest) (*pb.GetProjectedCostResponse, error) {
    // Input validation
    if req.Resource == nil {
        return nil, status.Error(codes.InvalidArgument, "resource descriptor is required")
    }
    
    // Authentication/authorization errors
    if !s.isAuthorized(ctx) {
        return nil, status.Error(codes.Unauthenticated, "authentication required")
    }
    
    // Not found errors
    if !s.resourceExists(req.Resource) {
        return nil, status.Error(codes.NotFound, "resource not found")
    }
    
    // External service errors
    if !s.isExternalServiceAvailable() {
        return nil, status.Error(codes.Unavailable, "pricing service temporarily unavailable")
    }
    
    // Rate limiting
    if s.isRateLimited() {
        return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded")
    }
    
    // Implementation errors (should be rare in production)
    return nil, status.Error(codes.Internal, "internal server error")
}
```

#### Status Code Guidelines

| Code | When to Use | Example |
|------|-------------|---------|
| `InvalidArgument` | Bad request parameters | Missing required fields, invalid format |
| `Unauthenticated` | Authentication failed | Invalid API key, expired token |
| `PermissionDenied` | Authorization failed | Insufficient permissions |
| `NotFound` | Resource doesn't exist | Unknown resource ID |
| `AlreadyExists` | Resource already exists | Duplicate creation attempts |
| `FailedPrecondition` | System state issue | Service not initialized |
| `ResourceExhausted` | Rate limiting | Too many requests |
| `Unavailable` | Temporary failure | External service down |
| `Internal` | Implementation error | Unexpected exceptions |

### Performance Considerations

#### Caching Strategy

```go
import (
    "sync"
    "time"
)

type CacheEntry struct {
    Data      interface{}
    ExpiresAt time.Time
}

type PricingCache struct {
    mu    sync.RWMutex
    cache map[string]CacheEntry
    ttl   time.Duration
}

func NewPricingCache(ttl time.Duration) *PricingCache {
    return &PricingCache{
        cache: make(map[string]CacheEntry),
        ttl:   ttl,
    }
}

func (c *PricingCache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    entry, exists := c.cache[key]
    if !exists || time.Now().After(entry.ExpiresAt) {
        return nil, false
    }
    
    return entry.Data, true
}

func (c *PricingCache) Set(key string, data interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.cache[key] = CacheEntry{
        Data:      data,
        ExpiresAt: time.Now().Add(c.ttl),
    }
}

// Usage in plugin
func (s *server) GetProjectedCost(ctx context.Context, req *pb.GetProjectedCostRequest) (*pb.GetProjectedCostResponse, error) {
    // Create cache key
    cacheKey := fmt.Sprintf("projected:%s:%s:%s:%s",
        req.Resource.Provider,
        req.Resource.ResourceType,
        req.Resource.Sku,
        req.Resource.Region,
    )
    
    // Check cache first
    if cached, found := s.cache.Get(cacheKey); found {
        return cached.(*pb.GetProjectedCostResponse), nil
    }
    
    // Calculate cost
    response, err := s.calculateProjectedCost(req)
    if err != nil {
        return nil, err
    }
    
    // Cache result
    s.cache.Set(cacheKey, response)
    
    return response, nil
}
```

#### Connection Pooling

```go
import (
    "net/http"
    "time"
)

type PluginConfig struct {
    APIEndpoint    string
    APIKey         string
    MaxConnections int
    Timeout        time.Duration
}

func NewHTTPClient(config *PluginConfig) *http.Client {
    transport := &http.Transport{
        MaxIdleConns:        config.MaxConnections,
        MaxIdleConnsPerHost: config.MaxConnections,
        IdleConnTimeout:     90 * time.Second,
        DisableCompression:  false,
    }
    
    return &http.Client{
        Transport: transport,
        Timeout:   config.Timeout,
    }
}
```

#### Request Batching

```go
// Batch multiple resource requests when possible
func (s *server) batchGetPricing(resources []*pb.ResourceDescriptor) ([]*pb.PricingSpec, error) {
    // Group resources by provider/region for efficient API calls
    batches := s.groupResourcesByBatch(resources)
    
    var allSpecs []*pb.PricingSpec
    for _, batch := range batches {
        specs, err := s.fetchPricingBatch(batch)
        if err != nil {
            return nil, err
        }
        allSpecs = append(allSpecs, specs...)
    }
    
    return allSpecs, nil
}
```

### Security Guidelines

#### Credential Management

```go
import (
    "os"
    "fmt"
)

type Credentials struct {
    APIKey      string
    APIEndpoint string
}

func LoadCredentials() (*Credentials, error) {
    apiKey := os.Getenv("PLUGIN_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("PLUGIN_API_KEY environment variable is required")
    }
    
    endpoint := os.Getenv("PLUGIN_API_ENDPOINT")
    if endpoint == "" {
        endpoint = "https://api.example.com" // Default endpoint
    }
    
    return &Credentials{
        APIKey:      apiKey,
        APIEndpoint: endpoint,
    }, nil
}
```

#### Input Validation

```go
func validateResourceDescriptor(rd *pb.ResourceDescriptor) error {
    if rd == nil {
        return status.Error(codes.InvalidArgument, "resource descriptor is required")
    }
    
    // Provider validation
    if rd.Provider == "" {
        return status.Error(codes.InvalidArgument, "provider is required")
    }
    if len(rd.Provider) > 50 {
        return status.Error(codes.InvalidArgument, "provider name too long")
    }
    
    // Resource type validation
    if rd.ResourceType == "" {
        return status.Error(codes.InvalidArgument, "resource_type is required")
    }
    if len(rd.ResourceType) > 100 {
        return status.Error(codes.InvalidArgument, "resource_type name too long")
    }
    
    // SKU validation (optional field)
    if len(rd.Sku) > 100 {
        return status.Error(codes.InvalidArgument, "sku name too long")
    }
    
    // Region validation
    if len(rd.Region) > 50 {
        return status.Error(codes.InvalidArgument, "region name too long")
    }
    
    // Tags validation
    if len(rd.Tags) > 50 {
        return status.Error(codes.InvalidArgument, "too many tags (max 50)")
    }
    
    for key, value := range rd.Tags {
        if len(key) > 128 || len(value) > 256 {
            return status.Error(codes.InvalidArgument, "tag key/value too long")
        }
    }
    
    return nil
}
```

#### Rate Limiting

```go
import (
    "golang.org/x/time/rate"
    "context"
)

type RateLimiter struct {
    limiter *rate.Limiter
}

func NewRateLimiter(requestsPerSecond int) *RateLimiter {
    return &RateLimiter{
        limiter: rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond),
    }
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
    return rl.limiter.Wait(ctx)
}

// Usage in RPC methods
func (s *server) GetProjectedCost(ctx context.Context, req *pb.GetProjectedCostRequest) (*pb.GetProjectedCostResponse, error) {
    // Rate limiting
    if err := s.rateLimiter.Wait(ctx); err != nil {
        return nil, status.Error(codes.DeadlineExceeded, "request timeout due to rate limiting")
    }
    
    // Process request...
}
```

### Configuration Management

#### Configuration Structure

```go
type PluginConfig struct {
    // Server configuration
    Port         int           `yaml:"port" env:"PLUGIN_PORT" default:"50051"`
    Host         string        `yaml:"host" env:"PLUGIN_HOST" default:"localhost"`
    ReadTimeout  time.Duration `yaml:"read_timeout" env:"PLUGIN_READ_TIMEOUT" default:"30s"`
    WriteTimeout time.Duration `yaml:"write_timeout" env:"PLUGIN_WRITE_TIMEOUT" default:"30s"`
    
    // API configuration
    APIEndpoint string `yaml:"api_endpoint" env:"API_ENDPOINT" required:"true"`
    APIKey      string `yaml:"api_key" env:"API_KEY" required:"true"`
    APITimeout  time.Duration `yaml:"api_timeout" env:"API_TIMEOUT" default:"10s"`
    
    // Cache configuration
    CacheTTL     time.Duration `yaml:"cache_ttl" env:"CACHE_TTL" default:"5m"`
    CacheSize    int           `yaml:"cache_size" env:"CACHE_SIZE" default:"1000"`
    
    // Rate limiting
    RateLimit    int `yaml:"rate_limit" env:"RATE_LIMIT" default:"100"`
    
    // Logging
    LogLevel     string `yaml:"log_level" env:"LOG_LEVEL" default:"info"`
    LogFormat    string `yaml:"log_format" env:"LOG_FORMAT" default:"json"`
}

func LoadConfig() (*PluginConfig, error) {
    config := &PluginConfig{}
    
    // Load from file, then override with env vars
    if err := loadFromFile("config.yaml", config); err != nil {
        log.Printf("No config file found, using defaults: %v", err)
    }
    
    if err := loadFromEnv(config); err != nil {
        return nil, fmt.Errorf("failed to load configuration: %w", err)
    }
    
    return config, nil
}
```

### Logging Best Practices

#### Structured Logging

```go
import (
    "log/slog"
    "os"
)

func setupLogging(config *PluginConfig) *slog.Logger {
    var handler slog.Handler
    
    opts := &slog.HandlerOptions{
        Level: parseLogLevel(config.LogLevel),
    }
    
    if config.LogFormat == "json" {
        handler = slog.NewJSONHandler(os.Stdout, opts)
    } else {
        handler = slog.NewTextHandler(os.Stdout, opts)
    }
    
    return slog.New(handler)
}

func (s *server) GetProjectedCost(ctx context.Context, req *pb.GetProjectedCostRequest) (*pb.GetProjectedCostResponse, error) {
    logger := s.logger.With(
        "method", "GetProjectedCost",
        "provider", req.Resource.Provider,
        "resource_type", req.Resource.ResourceType,
        "region", req.Resource.Region,
    )
    
    logger.Info("processing projected cost request")
    
    start := time.Now()
    defer func() {
        logger.Info("completed projected cost request",
            "duration", time.Since(start),
        )
    }()
    
    // Process request...
}
```

#### Request Tracing

```go
import (
    "context"
    "github.com/google/uuid"
)

func withTraceID(ctx context.Context) context.Context {
    traceID := uuid.New().String()
    return context.WithValue(ctx, "trace_id", traceID)
}

func getTraceID(ctx context.Context) string {
    if traceID, ok := ctx.Value("trace_id").(string); ok {
        return traceID
    }
    return "unknown"
}

// Middleware for adding trace IDs
func (s *server) traceMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    ctx = withTraceID(ctx)
    
    s.logger.Info("request started",
        "trace_id", getTraceID(ctx),
        "method", info.FullMethod,
    )
    
    resp, err := handler(ctx, req)
    
    if err != nil {
        s.logger.Error("request failed",
            "trace_id", getTraceID(ctx),
            "method", info.FullMethod,
            "error", err,
        )
    }
    
    return resp, err
}
```

### Resource Mapping Patterns

#### Provider-Specific Mappings

```go
type ResourceMapper struct {
    providerMappings map[string]map[string]string
}

func NewResourceMapper() *ResourceMapper {
    return &ResourceMapper{
        providerMappings: map[string]map[string]string{
            "aws": {
                "ec2":             "virtual-machine",
                "s3":              "object-storage",
                "rds":             "database",
                "elasticache":     "cache",
                "elasticsearch":   "search-engine",
            },
            "azure": {
                "virtual-machines": "virtual-machine",
                "blob-storage":     "object-storage",
                "sql-database":     "database",
            },
            "gcp": {
                "compute-engine":   "virtual-machine",
                "cloud-storage":    "object-storage",
                "cloud-sql":        "database",
            },
        },
    }
}

func (rm *ResourceMapper) MapResourceType(provider, resourceType string) (string, error) {
    providerMap, exists := rm.providerMappings[provider]
    if !exists {
        return "", fmt.Errorf("unsupported provider: %s", provider)
    }
    
    mappedType, exists := providerMap[resourceType]
    if !exists {
        return "", fmt.Errorf("unsupported resource type: %s for provider: %s", resourceType, provider)
    }
    
    return mappedType, nil
}
```

### Cost Calculation Patterns

#### Time-Based Pricing

```go
func calculateHourlyToMonthly(hourlyRate float64) float64 {
    // Standard month = 365.25 days / 12 months = 30.4375 days
    const hoursPerMonth = 24 * 365.25 / 12
    return hourlyRate * hoursPerMonth
}

func calculateDailyToMonthly(dailyRate float64) float64 {
    const daysPerMonth = 365.25 / 12
    return dailyRate * daysPerMonth
}

func prorateCostForPeriod(unitPrice float64, billingMode string, hours int) float64 {
    switch billingMode {
    case "per_hour":
        return unitPrice * float64(hours)
    case "per_day":
        days := math.Ceil(float64(hours) / 24)
        return unitPrice * days
    case "per_gb_month":
        months := float64(hours) / (24 * 365.25 / 12)
        return unitPrice * months
    default:
        return unitPrice
    }
}
```

#### Usage-Based Calculations

```go
type UsageCalculator struct {
    baseCost      float64
    usageMetrics  map[string]float64
    billingTiers  []BillingTier
}

type BillingTier struct {
    MinUsage float64
    MaxUsage float64
    Rate     float64
}

func (uc *UsageCalculator) CalculateCost(usage map[string]float64) float64 {
    totalCost := uc.baseCost
    
    for metric, amount := range usage {
        if rate, exists := uc.usageMetrics[metric]; exists {
            totalCost += uc.calculateTieredCost(amount, rate)
        }
    }
    
    return totalCost
}

func (uc *UsageCalculator) calculateTieredCost(usage, baseRate float64) float64 {
    cost := 0.0
    remaining := usage
    
    for _, tier := range uc.billingTiers {
        tierUsage := math.Min(remaining, tier.MaxUsage-tier.MinUsage)
        cost += tierUsage * tier.Rate * baseRate
        remaining -= tierUsage
        
        if remaining <= 0 {
            break
        }
    }
    
    return cost
}
```

### Health Checks and Monitoring

#### Health Check Implementation

```go
import (
    "net/http"
    "encoding/json"
)

type HealthStatus struct {
    Status    string            `json:"status"`
    Version   string            `json:"version"`
    Timestamp time.Time         `json:"timestamp"`
    Checks    map[string]string `json:"checks"`
}

func (s *server) setupHealthCheck() {
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        status := s.checkHealth()
        
        w.Header().Set("Content-Type", "application/json")
        if status.Status == "healthy" {
            w.WriteHeader(http.StatusOK)
        } else {
            w.WriteHeader(http.StatusServiceUnavailable)
        }
        
        json.NewEncoder(w).Encode(status)
    })
}

func (s *server) checkHealth() *HealthStatus {
    checks := make(map[string]string)
    overallStatus := "healthy"
    
    // Check external API connectivity
    if err := s.pingExternalAPI(); err != nil {
        checks["external_api"] = "unhealthy: " + err.Error()
        overallStatus = "unhealthy"
    } else {
        checks["external_api"] = "healthy"
    }
    
    // Check cache connectivity
    if s.cache != nil {
        checks["cache"] = "healthy"
    } else {
        checks["cache"] = "unhealthy: cache not initialized"
        overallStatus = "unhealthy"
    }
    
    return &HealthStatus{
        Status:    overallStatus,
        Version:   "1.0.0",
        Timestamp: time.Now(),
        Checks:    checks,
    }
}
```

### Development Best Practices

1. **Code Organization**: Separate concerns into different packages (handlers, services, models)
2. **Interface Design**: Use interfaces for testability and modularity
3. **Error Wrapping**: Use `fmt.Errorf` with `%w` verb to wrap errors
4. **Context Propagation**: Pass context through all function calls
5. **Resource Cleanup**: Use defer statements for cleanup operations
6. **Graceful Shutdown**: Implement proper shutdown handling
7. **Configuration Validation**: Validate configuration on startup
8. **Documentation**: Use godoc comments for public APIs
9. **Version Compatibility**: Follow semantic versioning for plugin releases
10. **Testing**: Maintain high test coverage with unit and integration tests

## Troubleshooting

This section covers common issues encountered during plugin development and their solutions.

### Common Issues

#### Build and Compilation Issues

**Issue**: `package github.com/rshade/pulumicost-spec/sdk/go/proto is not in GOROOT`

**Cause**: Missing or incorrect Go module dependencies.

**Solution**:

```bash
# Initialize go module if not done
go mod init your-plugin-name

# Add the dependency
go get github.com/rshade/pulumicost-spec/sdk/go/proto@latest

# Clean up dependencies
go mod tidy
```

---

**Issue**: `cannot find package "google.golang.org/grpc"`

**Cause**: Missing gRPC dependencies.

**Solution**:

```bash
go get google.golang.org/grpc@latest
go get google.golang.org/protobuf@latest
```

---

**Issue**: `undefined: pb.UnimplementedCostSourceServiceServer`

**Cause**: Using an older version of the protobuf compiler or missing embedded interface.

**Solution**:

```go
// Ensure your server struct embeds the unimplemented server
type server struct {
    pb.UnimplementedCostSourceServiceServer
    // your fields
}

// Update protoc-gen-go if needed
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

#### Runtime Issues

**Issue**: `transport: Error while dialing dial tcp :50051: connect: connection refused`

**Cause**: gRPC server not starting or binding to wrong port.

**Solution**:

```go
// Check if port is already in use
func isPortAvailable(port int) bool {
    ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        return false
    }
    ln.Close()
    return true
}

// Use dynamic port allocation for testing
lis, err := net.Listen("tcp", ":0") // Let OS choose port
if err != nil {
    log.Fatalf("Failed to listen: %v", err)
}
log.Printf("Server listening on %s", lis.Addr().String())
```

---

**Issue**: `rpc error: code = Unimplemented desc = method Name not implemented`

**Cause**: RPC method not implemented in server struct.

**Solution**:

```go
// Ensure all required methods are implemented
func (s *server) Name(ctx context.Context, req *pb.NameRequest) (*pb.NameResponse, error) {
    return &pb.NameResponse{Name: "your-plugin-name"}, nil
}

func (s *server) Supports(ctx context.Context, req *pb.SupportsRequest) (*pb.SupportsResponse, error) {
    // Implementation required
}

func (s *server) GetActualCost(ctx context.Context, req *pb.GetActualCostRequest) (*pb.GetActualCostResponse, error) {
    // Implementation required
}

func (s *server) GetProjectedCost(ctx context.Context, req *pb.GetProjectedCostRequest) (*pb.GetProjectedCostResponse, error) {
    // Implementation required
}

func (s *server) GetPricingSpec(ctx context.Context, req *pb.GetPricingSpecRequest) (*pb.GetPricingSpecResponse, error) {
    // Implementation required
}
```

---

**Issue**: `context deadline exceeded`

**Cause**: Operations taking too long or external API timeouts.

**Solution**:

```go
// Set appropriate timeouts
func (s *server) GetActualCost(ctx context.Context, req *pb.GetActualCostRequest) (*pb.GetActualCostResponse, error) {
    // Check if context is already cancelled
    select {
    case <-ctx.Done():
        return nil, status.Error(codes.Canceled, "request cancelled")
    default:
    }

    // Create timeout for external API calls
    apiCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    // Use apiCtx for external calls
    data, err := s.fetchFromExternalAPI(apiCtx, req)
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            return nil, status.Error(codes.DeadlineExceeded, "external API timeout")
        }
        return nil, status.Error(codes.Internal, err.Error())
    }

    return data, nil
}
```

#### Authentication Issues

**Issue**: `rpc error: code = Unauthenticated desc = authentication required`

**Cause**: Missing or invalid API credentials.

**Solution**:

```bash
# Set required environment variables
export PLUGIN_API_KEY="your-api-key"
export PLUGIN_API_ENDPOINT="https://api.example.com"

# Verify credentials are loaded
echo $PLUGIN_API_KEY
```

```go
// Add credential validation on startup
func validateCredentials() error {
    if os.Getenv("PLUGIN_API_KEY") == "" {
        return fmt.Errorf("PLUGIN_API_KEY environment variable is required")
    }
    
    // Test API connectivity
    resp, err := http.Get(os.Getenv("PLUGIN_API_ENDPOINT") + "/health")
    if err != nil {
        return fmt.Errorf("cannot reach API endpoint: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode == 401 {
        return fmt.Errorf("invalid API credentials")
    }
    
    return nil
}
```

#### Data Format Issues

**Issue**: `json: cannot unmarshal string into Go value of type float64`

**Cause**: Type mismatch between API response and expected data types.

**Solution**:

```go
// Handle flexible JSON parsing
type APIResponse struct {
    Cost interface{} `json:"cost"` // Can be string or number
}

func parseFlexibleFloat(v interface{}) (float64, error) {
    switch val := v.(type) {
    case float64:
        return val, nil
    case string:
        return strconv.ParseFloat(val, 64)
    case int:
        return float64(val), nil
    default:
        return 0, fmt.Errorf("cannot convert %T to float64", v)
    }
}
```

---

**Issue**: Invalid timestamp formats

**Cause**: Different timestamp formats from external APIs.

**Solution**:

```go
func parseFlexibleTime(timeStr string) (*timestamppb.Timestamp, error) {
    formats := []string{
        time.RFC3339,
        "2006-01-02T15:04:05Z",
        "2006-01-02 15:04:05",
        "2006-01-02",
    }
    
    for _, format := range formats {
        if t, err := time.Parse(format, timeStr); err == nil {
            return timestamppb.New(t), nil
        }
    }
    
    return nil, fmt.Errorf("unable to parse time: %s", timeStr)
}
```

### Debug Techniques

#### Enabling Debug Logging

```go
import (
    "log/slog"
    "os"
)

func setupDebugLogging() *slog.Logger {
    level := slog.LevelInfo
    if os.Getenv("DEBUG") == "true" {
        level = slog.LevelDebug
    }
    
    opts := &slog.HandlerOptions{
        Level: level,
        AddSource: true, // Include source code locations
    }
    
    handler := slog.NewJSONHandler(os.Stdout, opts)
    return slog.New(handler)
}

// Usage in RPC methods
func (s *server) GetProjectedCost(ctx context.Context, req *pb.GetProjectedCostRequest) (*pb.GetProjectedCostResponse, error) {
    s.logger.Debug("GetProjectedCost called",
        "provider", req.Resource.Provider,
        "resource_type", req.Resource.ResourceType,
        "sku", req.Resource.Sku,
        "region", req.Resource.Region,
    )
    
    // Your implementation
    
    s.logger.Debug("GetProjectedCost response",
        "unit_price", response.UnitPrice,
        "currency", response.Currency,
    )
    
    return response, nil
}
```

#### Testing with grpcurl

```bash
# Test server is running
grpcurl -plaintext localhost:50051 list

# Test Name RPC
grpcurl -plaintext localhost:50051 pulumicost.v1.CostSourceService.Name

# Test with request data
grpcurl -plaintext \
    -d '{"resource": {"provider": "aws", "resource_type": "ec2", "sku": "t3.micro", "region": "us-east-1"}}' \
    localhost:50051 \
    pulumicost.v1.CostSourceService.GetProjectedCost

# Test with file input
echo '{"resource": {"provider": "aws", "resource_type": "ec2"}}' > request.json
grpcurl -plaintext -d @ localhost:50051 pulumicost.v1.CostSourceService.Supports < request.json
```

#### Using gRPC Health Checks

```go
import (
    "google.golang.org/grpc/health"
    "google.golang.org/grpc/health/grpc_health_v1"
)

func (s *server) setupHealthCheck() {
    healthServer := health.NewServer()
    healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
    grpc_health_v1.RegisterHealthServer(s.grpcServer, healthServer)
}

// Test health check
// grpcurl -plaintext localhost:50051 grpc.health.v1.Health.Check
```

#### Memory and Performance Profiling

```go
import (
    _ "net/http/pprof"
    "net/http"
)

func (s *server) setupProfiling() {
    go func() {
        log.Println("Starting profiling server on :6060")
        http.ListenAndServe("localhost:6060", nil)
    }()
}

// Access profiling data:
// go tool pprof http://localhost:6060/debug/pprof/profile
// go tool pprof http://localhost:6060/debug/pprof/heap
```

### FAQ

#### Q: How do I handle different currency conversions?

**A**: Store base currency rates and convert as needed:

```go
type CurrencyConverter struct {
    baseCurrency string
    rates       map[string]float64
}

func (cc *CurrencyConverter) Convert(amount float64, from, to string) float64 {
    if from == to {
        return amount
    }
    
    // Convert to base currency first, then to target
    baseAmount := amount / cc.rates[from]
    return baseAmount * cc.rates[to]
}
```

#### Q: Should I cache pricing data, and for how long?

**A**: Yes, cache pricing data to reduce API calls. Recommended TTLs:

- Static pricing (AWS public pricing): 24 hours
- Dynamic pricing (spot prices): 5-15 minutes
- Usage data: 1-5 minutes

```go
type CacheConfig struct {
    StaticPricingTTL  time.Duration // 24 hours
    DynamicPricingTTL time.Duration // 5 minutes
    UsageDataTTL      time.Duration // 1 minute
}
```

#### Q: How do I handle rate limiting from external APIs?

**A**: Implement exponential backoff and circuit breakers:

```go
func (s *server) makeAPICallWithRetry(ctx context.Context, req interface{}) (interface{}, error) {
    backoff := time.Second
    maxRetries := 3
    
    for i := 0; i < maxRetries; i++ {
        resp, err := s.makeAPICall(ctx, req)
        if err == nil {
            return resp, nil
        }
        
        // Check if it's a rate limit error
        if isRateLimitError(err) {
            time.Sleep(backoff)
            backoff *= 2 // Exponential backoff
            continue
        }
        
        return nil, err
    }
    
    return nil, status.Error(codes.ResourceExhausted, "max retries exceeded")
}
```

#### Q: What's the best way to handle missing data?

**A**: Return appropriate gRPC status codes:

```go
func (s *server) GetActualCost(ctx context.Context, req *pb.GetActualCostRequest) (*pb.GetActualCostResponse, error) {
    data, err := s.fetchCostData(req.ResourceId)
    if err != nil {
        if isNotFoundError(err) {
            return nil, status.Error(codes.NotFound, "resource not found")
        }
        return nil, status.Error(codes.Internal, err.Error())
    }
    
    if len(data) == 0 {
        // Return empty results, not an error
        return &pb.GetActualCostResponse{Results: []*pb.ActualCostResult{}}, nil
    }
    
    return &pb.GetActualCostResponse{Results: data}, nil
}
```

#### Q: How do I test my plugin with PulumiCost?

**A**: Use the integration testing approach:

1. Start your plugin server
2. Use PulumiCost CLI to connect to your plugin
3. Run test queries against your plugin

```bash
# Start your plugin
./my-plugin -port 50051

# Test with PulumiCost (example commands)
pulumicost plugin add my-plugin localhost:50051
pulumicost cost get --provider custom --resource-type vm --sku small
```

#### Q: What metrics should I monitor in production?

**A**: Key metrics to track:

```go
type Metrics struct {
    RequestCount    prometheus.Counter
    RequestDuration prometheus.Histogram
    ErrorRate       prometheus.Counter
    CacheHitRate    prometheus.Gauge
    APILatency      prometheus.Histogram
}

// Instrument your RPC methods
func (s *server) GetProjectedCost(ctx context.Context, req *pb.GetProjectedCostRequest) (*pb.GetProjectedCostResponse, error) {
    start := time.Now()
    s.metrics.RequestCount.Inc()
    
    defer func() {
        s.metrics.RequestDuration.Observe(time.Since(start).Seconds())
    }()
    
    // Your implementation
}
```

### Getting Help

If you encounter issues not covered in this guide:

1. **Check Logs**: Enable debug logging and examine error messages
2. **Review Examples**: Compare your implementation with the working examples
3. **Test Isolation**: Create minimal test cases to isolate the problem
4. **Community Support**: Reach out to the PulumiCost community
5. **GitHub Issues**: Report bugs or request features in the spec repository

### Useful Commands Reference

```bash
# Development workflow
go mod tidy                    # Clean dependencies
go build .                     # Build plugin
go test -v ./...              # Run tests
go test -race ./...           # Race condition detection

# gRPC testing
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext localhost:50051 ServiceName.MethodName

# Debugging
go run -race .                # Run with race detection
dlv debug                     # Debug with Delve
go tool pprof                 # Performance profiling

# Build for production
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o plugin .
```
