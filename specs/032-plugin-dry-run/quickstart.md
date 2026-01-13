# Quickstart: Plugin Capability Dry Run Mode

This guide shows how to use the dry-run capability to discover plugin field mappings.

## Overview

The dry-run feature allows hosts to query a plugin for its FOCUS field mapping logic
without performing actual cost data retrieval. This is useful for:

- **Debugging**: Understand which fields a plugin populates for a resource type
- **Validation**: Verify plugin configuration before production deployment
- **Comparison**: Compare capabilities across different plugins

## Using the DryRun RPC

### Basic Usage (Go)

```go
package main

import (
    "context"
    "fmt"
    "log"

    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"

    pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

func main() {
    // Connect to plugin
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    client := pbc.NewCostSourceServiceClient(conn)
    ctx := context.Background()

    // Check if plugin supports dry-run
    supportsResp, err := client.Supports(ctx, &pbc.SupportsRequest{
        Resource: &pbc.ResourceDescriptor{
            Provider:     "aws",
            ResourceType: "ec2",
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    if !supportsResp.Capabilities["dry_run"] {
        log.Println("Plugin does not support dry-run introspection")
        return
    }

    // Query field mappings
    resp, err := client.DryRun(ctx, &pbc.DryRunRequest{
        Resource: &pbc.ResourceDescriptor{
            Provider:     "aws",
            ResourceType: "ec2",
            Region:       "us-east-1",
        },
    })
    if err != nil {
        if status.Code(err) == codes.Unimplemented {
            log.Println("Plugin does not implement DryRun RPC")
            return
        }
        log.Fatal(err)
    }

    // Check configuration validity
    if !resp.ConfigurationValid {
        fmt.Println("Configuration errors:")
        for _, e := range resp.ConfigurationErrors {
            fmt.Printf("  - %s\n", e)
        }
        return
    }

    // List supported fields
    fmt.Printf("Field mappings for aws/ec2 (%d fields):\n", len(resp.FieldMappings))
    for _, fm := range resp.FieldMappings {
        status := statusToString(fm.SupportStatus)
        fmt.Printf("  %-30s %s\n", fm.FieldName, status)
        if fm.ConditionDescription != "" {
            fmt.Printf("    └─ %s\n", fm.ConditionDescription)
        }
    }
}

func statusToString(s pbc.FieldSupportStatus) string {
    switch s {
    case pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_SUPPORTED:
        return "✓ SUPPORTED"
    case pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_UNSUPPORTED:
        return "✗ UNSUPPORTED"
    case pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_CONDITIONAL:
        return "? CONDITIONAL"
    case pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_DYNAMIC:
        return "~ DYNAMIC"
    default:
        return "? UNKNOWN"
    }
}
```

### Using dry_run Flag on Cost RPCs

You can also use the `dry_run` flag on existing GetActualCost/GetProjectedCost calls:

```go
// Inline validation during cost query
resp, err := client.GetActualCost(ctx, &pbc.GetActualCostRequest{
    ResourceId: "i-1234567890abcdef0",
    Start:      startTime,
    End:        endTime,
    DryRun:     true,  // Return field mappings instead of cost data
})

if resp.DryRunResult != nil {
    // Introspection mode - check field mappings
    for _, fm := range resp.DryRunResult.FieldMappings {
        fmt.Printf("%s: %v\n", fm.FieldName, fm.SupportStatus)
    }
} else {
    // Normal mode - process cost data
    for _, result := range resp.Results {
        fmt.Printf("Cost: %f %s\n", result.Cost, result.Source)
    }
}
```

## Implementing DryRun in a Plugin

### Plugin Implementation (Go)

```go
package main

import (
    "context"

    pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

type MyPlugin struct {
    pbc.UnimplementedCostSourceServiceServer
}

func (p *MyPlugin) DryRun(
    ctx context.Context,
    req *pbc.DryRunRequest,
) (*pbc.DryRunResponse, error) {
    // Check if resource type is supported
    if !p.supportsResourceType(req.Resource) {
        return &pbc.DryRunResponse{
            ResourceTypeSupported: false,
            ConfigurationValid:    true,
        }, nil
    }

    // Build field mappings based on resource type
    mappings := p.getFieldMappings(req.Resource)

    return &pbc.DryRunResponse{
        FieldMappings:         mappings,
        ConfigurationValid:    true,
        ResourceTypeSupported: true,
    }, nil
}

func (p *MyPlugin) getFieldMappings(res *pbc.ResourceDescriptor) []*pbc.FieldMapping {
    // Return mappings based on resource type
    switch res.ResourceType {
    case "ec2":
        return []*pbc.FieldMapping{
            {
                FieldName:     "service_category",
                SupportStatus: pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_SUPPORTED,
                ExpectedType:  "enum",
            },
            {
                FieldName:     "billed_cost",
                SupportStatus: pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_SUPPORTED,
                ExpectedType:  "double",
            },
            {
                FieldName:            "availability_zone",
                SupportStatus:        pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_CONDITIONAL,
                ConditionDescription: "Populated for instances in multi-AZ VPCs",
                ExpectedType:         "string",
            },
            // ... more fields
        }
    default:
        return nil
    }
}

// Advertise dry-run support in capabilities
func (p *MyPlugin) Supports(
    ctx context.Context,
    req *pbc.SupportsRequest,
) (*pbc.SupportsResponse, error) {
    return &pbc.SupportsResponse{
        Supported: p.supportsResourceType(req.Resource),
        Capabilities: map[string]bool{
            "dry_run":         true,
            "recommendations": false,
        },
    }, nil
}
```

## Field Status Reference

| Status | Meaning | Example |
|--------|---------|---------|
| SUPPORTED | Always populated for this resource | service_category, billed_cost |
| UNSUPPORTED | Never populated for this resource | kubernetes-specific fields on AWS |
| CONDITIONAL | Depends on resource configuration | availability_zone (multi-AZ only) |
| DYNAMIC | Requires runtime data | Actual cost values |

## Troubleshooting

### "Plugin does not implement DryRun RPC"

The plugin is a legacy implementation that doesn't support dry-run introspection.
Options:

1. Upgrade the plugin to a version that supports dry-run
2. Use `Supports` RPC to check basic resource type support instead

### "Configuration errors" in response

The plugin detected configuration issues. Check:

1. Required environment variables (API keys, endpoints)
2. Plugin configuration file syntax
3. Network connectivity to backend services

### Empty field_mappings

When `resource_type_supported` is false, `field_mappings` will be empty.
Verify the resource type string matches what the plugin expects.
