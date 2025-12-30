# Quickstart: Multi-Protocol Plugin Access

## Enabling Multi-Protocol Support in a Plugin

Modify your `main.go` to use the updated `pluginsdk.Serve` with web support:

```go
package main

import (
    "context"
    "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
)

func main() {
    ctx := context.Background()

    // Enable multi-protocol support (gRPC + gRPC-Web + Connect)
    err := pluginsdk.Serve(ctx, pluginsdk.ServeConfig{
        Plugin: &MyPlugin{},
        Web: pluginsdk.WebConfig{
            Enabled:              true,
            AllowedOrigins:       []string{"http://localhost:3000", "https://app.example.com"},
            EnableHealthEndpoint: true,
        },
    })
    if err != nil {
        panic(err)
    }
}
```

## Testing with curl (Connect Protocol)

The Connect protocol supports JSON, making debugging easy:

```bash
# Get plugin name
curl -X POST http://localhost:8080/pulumicost.v1.CostSourceService/Name \
  -H "Content-Type: application/json" \
  -d '{}'

# Estimate cost for a resource
curl -X POST http://localhost:8080/pulumicost.v1.CostSourceService/EstimateCost \
  -H "Content-Type: application/json" \
  -d '{
    "resource": {
      "id": "aws:ec2:i-1234567890abcdef0",
      "provider": "aws",
      "resource_type": "aws_instance",
      "region": "us-east-1"
    }
  }'

# Health check
curl http://localhost:8080/healthz
```

## Using the Go Client

```go
package main

import (
    "context"
    "fmt"
    "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk/client"
    pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

func main() {
    ctx := context.Background()

    // Connect to a plugin (uses Connect protocol by default)
    c, err := client.New(ctx, "http://localhost:8080")
    if err != nil {
        panic(err)
    }
    defer c.Close()

    // Get plugin name
    name, err := c.Name(ctx)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Connected to plugin: %s\n", name)

    // Batch estimate costs for multiple resources
    resources := []*pbc.ResourceDescriptor{
        {Id: "aws:ec2:i-123", Provider: "aws", ResourceType: "aws_instance"},
        {Id: "aws:s3:my-bucket", Provider: "aws", ResourceType: "aws_s3_bucket"},
    }

    results, err := c.BatchEstimateCost(ctx, resources, 10) // 10 concurrent requests
    if err != nil {
        panic(err)
    }

    for _, res := range results {
        if res.Error != nil {
            fmt.Printf("Resource %s failed: %v\n", res.Resource.Id, res.Error)
        } else {
            fmt.Printf("Resource %s cost: $%.2f/month\n", res.Resource.Id, res.Response.GetMonthlyEstimate())
        }
    }
}
```

## Browser Access (gRPC-Web/Connect)

For browser-based access, you can use:

1. **Connect Protocol with fetch** (simplest):

```javascript
const response = await fetch('http://localhost:8080/pulumicost.v1.CostSourceService/Name', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({})
});
const data = await response.json();
console.log('Plugin name:', data.name);
```

1. **@connectrpc/connect-web** (full typed client):

```typescript
import { createConnectTransport } from "@connectrpc/connect-web";
import { createClient } from "@connectrpc/connect";
import { CostSourceService } from "./gen/pulumicost/v1/costsource_connect";

const transport = createConnectTransport({
  baseUrl: "http://localhost:8080",
});

const client = createClient(CostSourceService, transport);
const response = await client.name({});
console.log('Plugin name:', response.name);
```

## Protocol Comparison

| Protocol   | Use Case                    | Content-Type              |
| ---------- | --------------------------- | ------------------------- |
| gRPC       | Go/Python/Java backends     | `application/grpc`        |
| gRPC-Web   | Browser with grpc-web lib   | `application/grpc-web`    |
| Connect    | curl, fetch, any HTTP client| `application/json`        |

All three protocols are served on the same port and endpoint!
