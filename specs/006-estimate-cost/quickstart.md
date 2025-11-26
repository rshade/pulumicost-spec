# Quickstart: EstimateCost RPC

This guide shows how to use the EstimateCost RPC method for proactive cost estimation.

## Overview

The EstimateCost RPC enables you to get cost estimates for Pulumi resources **before** deploying
them. This allows you to:

- Compare costs between different resource configurations
- Make budget-informed infrastructure decisions during development
- Validate cost assumptions before committing to resource specifications

## Prerequisites

- A PulumiCost plugin that implements the EstimateCost RPC
- The resource type you want to estimate must be supported by the plugin
- gRPC client library for your language (examples use Go)

## Basic Usage

### Step 1: Check Resource Support

First, verify the plugin supports your resource type using the existing `Supports` RPC:

```go
import (
    pb "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
    "google.golang.org/grpc"
)

conn, _ := grpc.Dial("localhost:50051", grpc.WithInsecure())
client := pb.NewCostSourceClient(conn)

supportsResp, _ := client.Supports(ctx, &pb.SupportsRequest{
    ResourceType: "aws:ec2/instance:Instance",
})

if !supportsResp.GetSupported() {
    // Resource type not supported by this plugin
}
```

### Step 2: Estimate Cost

Call `EstimateCost` with the resource type and configuration attributes:

```go
req := &pb.EstimateCostRequest{
    ResourceType: "aws:ec2/instance:Instance",
    Attributes: &structpb.Struct{
        Fields: map[string]*structpb.Value{
            "instanceType": structpb.NewStringValue("t3.micro"),
            "region":       structpb.NewStringValue("us-east-1"),
        },
    },
}

resp, err := client.EstimateCost(ctx, req)
if err != nil {
    // Handle error (invalid format, unsupported resource, missing attributes, etc.)
    log.Fatalf("EstimateCost failed: %v", err)
}

fmt.Printf("Estimated monthly cost: %s %s\n", resp.GetCostMonthly(), resp.GetCurrency())
// Output: Estimated monthly cost: 7.30 USD
```

### Step 3: Compare Configurations

Compare costs across different configurations:

```go
instances := []string{"t3.micro", "t3.small", "t3.medium", "t3.large"}

for _, instanceType := range instances {
    req := &pb.EstimateCostRequest{
        ResourceType: "aws:ec2/instance:Instance",
        Attributes: &structpb.Struct{
            Fields: map[string]*structpb.Value{
                "instanceType": structpb.NewStringValue(instanceType),
                "region":       structpb.NewStringValue("us-east-1"),
            },
        },
    }

    resp, err := client.EstimateCost(ctx, req)
    if err != nil {
        continue
    }

    fmt.Printf("%s: %s %s/month\n",
        instanceType, resp.GetCostMonthly(), resp.GetCurrency())
}

// Output:
// t3.micro: 7.30 USD/month
// t3.small: 14.60 USD/month
// t3.medium: 29.20 USD/month
// t3.large: 58.40 USD/month
```

## Cross-Provider Examples

### AWS EC2 Instance

```go
req := &pb.EstimateCostRequest{
    ResourceType: "aws:ec2/instance:Instance",
    Attributes: &structpb.Struct{
        Fields: map[string]*structpb.Value{
            "instanceType": structpb.NewStringValue("t3.micro"),
            "region":       structpb.NewStringValue("us-east-1"),
        },
    },
}
```

### Azure Virtual Machine

```go
req := &pb.EstimateCostRequest{
    ResourceType: "azure:compute/virtualMachine:VirtualMachine",
    Attributes: &structpb.Struct{
        Fields: map[string]*structpb.Value{
            "vmSize":   structpb.NewStringValue("Standard_B1s"),
            "location": structpb.NewStringValue("eastus"),
            "osType":   structpb.NewStringValue("Linux"),
        },
    },
}
```

### GCP Compute Instance

```go
req := &pb.EstimateCostRequest{
    ResourceType: "gcp:compute/instance:Instance",
    Attributes: &structpb.Struct{
        Fields: map[string]*structpb.Value{
            "machineType": structpb.NewStringValue("e2-micro"),
            "zone":        structpb.NewStringValue("us-central1-a"),
        },
    },
}
```

## Error Handling

### Invalid Resource Type Format

```go
resp, err := client.EstimateCost(ctx, &pb.EstimateCostRequest{
    ResourceType: "invalid-format",
    Attributes:   nil,
})

if err != nil {
    st := status.Convert(err)
    if st.Code() == codes.InvalidArgument {
        // Error: "resource_type must follow provider:module/resource:Type format"
    }
}
```

### Unsupported Resource Type

```go
resp, err := client.EstimateCost(ctx, &pb.EstimateCostRequest{
    ResourceType: "aws:lambda/function:Function",
    Attributes:   nil,
})

if err != nil {
    st := status.Convert(err)
    if st.Code() == codes.NotFound {
        // Error: "resource type ... is not supported by this plugin"
    }
}
```

### Missing Required Attributes

```go
resp, err := client.EstimateCost(ctx, &pb.EstimateCostRequest{
    ResourceType: "aws:ec2/instance:Instance",
    Attributes:   nil, // Empty attributes
})

if err != nil {
    st := status.Convert(err)
    if st.Code() == codes.InvalidArgument {
        // Error: "missing required attributes: [instanceType, region]"
    }
}
```

### Pricing Source Unavailable

```go
resp, err := client.EstimateCost(ctx, &pb.EstimateCostRequest{
    ResourceType: "aws:ec2/instance:Instance",
    Attributes: &structpb.Struct{
        Fields: map[string]*structpb.Value{
            "instanceType": structpb.NewStringValue("t3.micro"),
            "region":       structpb.NewStringValue("us-east-1"),
        },
    },
})

if err != nil {
    st := status.Convert(err)
    if st.Code() == codes.Unavailable {
        // Transient error - pricing source unavailable
        // Retry with exponential backoff in your application layer
    }
}
```

## Best Practices

### 1. Validate Resource Type Format

Before calling EstimateCost, validate the resource type format to avoid unnecessary RPC calls:

```go
import "regexp"

var resourceTypePattern = regexp.MustCompile(`^[a-z0-9]+:[a-z0-9]+/[a-z0-9]+:[A-Z][a-zA-Z0-9]*$`)

if !resourceTypePattern.MatchString(resourceType) {
    return fmt.Errorf("invalid resource type format: %s", resourceType)
}
```

### 2. Check Support Before Estimation

Always use the `Supports` RPC first to verify the plugin can handle your resource type:

```go
supported, err := checkSupports(client, resourceType)
if err != nil || !supported {
    return fmt.Errorf("resource type %s not supported", resourceType)
}

// Now safe to call EstimateCost
```

### 3. Handle Null Attributes

The `attributes` field may be null or missing. Plugins interpret this as empty attributes and will
return errors if required attributes are missing:

```go
// Null attributes are valid - plugin will validate
req := &pb.EstimateCostRequest{
    ResourceType: "aws:s3/bucket:Bucket",
    Attributes:   nil, // Valid for resources with no required attributes
}
```

### 4. Implement Retry Logic

For `Unavailable` errors (pricing source failures), implement retry logic in your application:

```go
import "time"

func estimateWithRetry(client pb.CostSourceClient, req *pb.EstimateCostRequest, maxRetries int) (*pb.EstimateCostResponse, error) {
    var resp *pb.EstimateCostResponse
    var err error

    for i := 0; i < maxRetries; i++ {
        resp, err = client.EstimateCost(context.Background(), req)
        if err == nil {
            return resp, nil
        }

        st := status.Convert(err)
        if st.Code() != codes.Unavailable {
            return nil, err // Non-retryable error
        }

        // Exponential backoff
        backoff := time.Duration(1<<uint(i)) * time.Second
        time.Sleep(backoff)
    }

    return nil, fmt.Errorf("max retries exceeded: %w", err)
}
```

### 5. Cache Estimates

EstimateCost is idempotent - identical inputs always produce identical outputs. Cache results to
reduce RPC calls:

```go
type estimateCache struct {
    mu    sync.RWMutex
    cache map[string]*pb.EstimateCostResponse
}

func (c *estimateCache) get(resourceType string, attributes *structpb.Struct) *pb.EstimateCostResponse {
    c.mu.RLock()
    defer c.mu.RUnlock()

    key := cacheKey(resourceType, attributes)
    return c.cache[key]
}
```

## Performance Considerations

- **Target Response Time**: <500ms for standard resource types (SC-002)
- **Idempotency**: Results are deterministic - safe to cache
- **Concurrency**: Supports multiple concurrent requests
- **No Retry in SDK**: Retry logic is handled by plugins/core, not the SDK

## Next Steps

- See [data-model.md](./data-model.md) for complete protobuf message definitions
- See [contracts/examples.md](./contracts/examples.md) for more request/response examples
- Implement EstimateCost in your plugin (see plugin implementation guide)
- Add conformance tests for your plugin implementation
