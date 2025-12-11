# Quickstart: Using Validation Helpers

## Introduction

The `pluginsdk` package provides validation helpers to ensure requests are well-formed before processing. These helpers are used by both the Core (before sending requests) and Plugins (defense-in-depth).

## Usage

### In a Plugin (Defense-in-Depth)

```go
import (
    "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
    pb "github.com/rshade/pulumicost-spec/proto/pulumicost/v1"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func (s *MyPluginServer) GetProjectedCost(ctx context.Context, req *pb.GetProjectedCostRequest) (*pb.GetProjectedCostResponse, error) {
    // Validate request first
    if err := pluginsdk.ValidateProjectedCostRequest(req); err != nil {
        return nil, status.Error(codes.InvalidArgument, err.Error())
    }

    // Process valid request...
}
```

### In Core (Pre-flight)

```go
func (c *Core) CallPlugin(req *pb.GetProjectedCostRequest) {
    if err := pluginsdk.ValidateProjectedCostRequest(req); err != nil {
        // Log error and return to user immediately
        log.Error().Err(err).Msg("Invalid request configuration")
        return
    }
    // Send to plugin...
}
```

## Error Messages

The validation errors are designed to be actionable.

Example Error:
> "resource.sku is required (use mapping.ExtractAWSSKU)"
