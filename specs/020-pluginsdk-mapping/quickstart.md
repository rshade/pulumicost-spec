# Quickstart: PluginSDK Mapping Package

**Feature**: 016-pluginsdk-mapping
**Date**: 2025-12-09

## Installation

The mapping package is part of the finfocus-spec SDK:

```go
import "github.com/rshade/finfocus-spec/sdk/go/pluginsdk/mapping"
```

## Basic Usage

### AWS Resource Property Extraction

```go
package main

import (
    "fmt"
    "github.com/rshade/finfocus-spec/sdk/go/pluginsdk/mapping"
)

func main() {
    // EC2 Instance properties from Pulumi
    ec2Props := map[string]string{
        "instanceType":     "t3.medium",
        "availabilityZone": "us-east-1a",
    }

    sku := mapping.ExtractAWSSKU(ec2Props)       // "t3.medium"
    region := mapping.ExtractAWSRegion(ec2Props) // "us-east-1"

    fmt.Printf("SKU: %s, Region: %s\n", sku, region)
}
```

### Azure Resource Property Extraction

```go
// Azure VM properties from Pulumi
vmProps := map[string]string{
    "vmSize":   "Standard_D2s_v3",
    "location": "eastus",
}

sku := mapping.ExtractAzureSKU(vmProps)       // "Standard_D2s_v3"
region := mapping.ExtractAzureRegion(vmProps) // "eastus"
```

### GCP Resource Property Extraction

```go
// GCP Compute Instance properties from Pulumi
computeProps := map[string]string{
    "machineType": "n1-standard-4",
    "zone":        "us-central1-a",
}

sku := mapping.ExtractGCPSKU(computeProps)       // "n1-standard-4"
region := mapping.ExtractGCPRegion(computeProps) // "us-central1"
```

### Generic Extraction with Custom Keys

```go
// Custom resource with non-standard property names
customProps := map[string]string{
    "customSKUField":    "my-sku-value",
    "customRegionField": "custom-region",
}

// Use generic extractors with custom key lists
sku := mapping.ExtractSKU(customProps, "customSKUField", "fallbackSKU")
region := mapping.ExtractRegion(customProps, "customRegionField", "fallbackRegion")
```

## Common Patterns

### Plugin Implementation

```go
package myplugin

import (
    "github.com/rshade/finfocus-spec/sdk/go/pluginsdk/mapping"
    pb "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

func (p *MyPlugin) GetProjectedCost(ctx context.Context, req *pb.GetProjectedCostRequest) (
    *pb.GetProjectedCostResponse, error,
) {
    for _, resource := range req.Resources {
        // Convert properties to string map
        props := make(map[string]string)
        for k, v := range resource.Properties {
            props[k] = v
        }

        // Extract SKU and region based on provider
        var sku, region string
        switch resource.Provider {
        case "aws":
            sku = mapping.ExtractAWSSKU(props)
            region = mapping.ExtractAWSRegion(props)
        case "azure":
            sku = mapping.ExtractAzureSKU(props)
            region = mapping.ExtractAzureRegion(props)
        case "gcp":
            sku = mapping.ExtractGCPSKU(props)
            region = mapping.ExtractGCPRegion(props)
        default:
            // Use generic extractors for unknown providers
            sku = mapping.ExtractSKU(props)
            region = mapping.ExtractRegion(props)
        }

        // Use sku and region for pricing lookup...
    }
    // ...
}
```

### Handling Missing Data

```go
props := map[string]string{
    "someOtherField": "value",
}

sku := mapping.ExtractAWSSKU(props)
if sku == "" {
    // No SKU found - handle gracefully
    log.Warn().Msg("no SKU property found in resource")
    // Use default or skip resource
}

// Empty/nil inputs are safe
var nilProps map[string]string
sku = mapping.ExtractAWSSKU(nilProps) // Returns ""
```

### Direct Zone/AZ Extraction

```go
// When you have the zone/AZ value directly (not from a property map)
awsRegion := mapping.ExtractAWSRegionFromAZ("us-west-2a")    // "us-west-2"
gcpRegion := mapping.ExtractGCPRegionFromZone("europe-west1-b") // "europe-west1"
```

### GCP Region Validation

```go
// Check if a region string is a valid GCP region
if mapping.IsValidGCPRegion("us-central1") {
    fmt.Println("Valid GCP region")
}

// List all known GCP regions
for _, region := range mapping.GCPRegions {
    fmt.Println(region)
}
```

## Error Handling

The mapping package is designed to never panic. All functions return empty string
for invalid or missing inputs:

```go
// All of these return "" safely
mapping.ExtractAWSSKU(nil)                    // ""
mapping.ExtractAWSSKU(map[string]string{})    // ""
mapping.ExtractAWSRegionFromAZ("")            // ""
mapping.ExtractGCPRegionFromZone("invalid")   // "" (not a valid region)
```

## Performance

The mapping functions are optimized for minimal allocation:

- Target: <50 ns/op, 0 allocs/op
- No regex compilation at runtime
- Simple string operations only
- Package-level constants for key names

## Next Steps

1. Import the mapping package in your plugin
2. Use provider-specific functions for known cloud resources
3. Use generic functions for custom or unknown resource types
4. Handle empty returns gracefully (indicates missing/invalid data)
