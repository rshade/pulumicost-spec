# Quickstart: Using the ARN Field

## Overview

The `GetActualCostRequest` message now includes an optional `arn` field (field number 5). This field
is intended to carry the canonical cloud identifier (e.g., AWS ARN) when available, enabling plugins
to make more precise queries to cloud provider APIs.

## Usage

### Go SDK

When constructing a `GetActualCostRequest`, populate the `Arn` field:

```go
req := &pbc.GetActualCostRequest{
    ResourceId: "urn:pulumi:stack:resource:type::name",
    Start:      startTime,
    End:        endTime,
    Arn:        "arn:aws:ec2:us-east-1:123456789012:instance/i-1234567890abcdef0",
}
```

### Plugin Implementation

Plugins should check if `Arn` is populated and use it as the primary identifier for cloud API calls.
If empty, fall back to `ResourceId`.

```go
func (s *MyServer) GetActualCost(ctx context.Context, req *pbc.GetActualCostRequest) (*pbc.GetActualCostResponse, error) {
    identifier := req.Arn
    if identifier == "" {
        identifier = req.ResourceId
    }
    // ... query cost source using identifier ...
}
```
