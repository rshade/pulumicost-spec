# Data Model: PluginSDK Mapping Package

**Feature**: 016-pluginsdk-mapping
**Date**: 2025-12-09

## Overview

The mapping package operates on simple data types with no complex domain entities. The primary
data structures are Go's built-in `map[string]string` for property maps and `string` for
extracted values.

## Core Types

### PropertyMap (Input Type)

```go
// PropertyMap represents Pulumi resource properties as a string-to-string map.
// This matches the common representation of resource properties in Pulumi outputs.
type PropertyMap = map[string]string
```

**Usage**: All extraction functions accept `PropertyMap` as their primary input.

**Validation Rules**:

- May be nil (functions return empty string)
- May be empty (functions return empty string)
- Keys and values are strings (no type conversion needed)

### Extracted Values (Output Type)

All extraction functions return `string`:

- Empty string indicates "not found" or "invalid input"
- Non-empty string is the extracted value (SKU, region, etc.)

## GCP Regions List

A package-level slice containing known valid GCP regions for validation:

```go
// GCPRegions contains all known GCP regions as of 2025-12.
// Used for validation during zone-to-region extraction.
//
//nolint:gochecknoglobals // Intentional: read-only reference data
var GCPRegions = []string{
    // Asia Pacific
    "asia-east1",           // Taiwan
    "asia-east2",           // Hong Kong
    "asia-northeast1",      // Tokyo
    "asia-northeast2",      // Osaka
    "asia-northeast3",      // Seoul
    "asia-south1",          // Mumbai
    "asia-south2",          // Delhi
    "asia-southeast1",      // Singapore
    "asia-southeast2",      // Jakarta

    // Australia
    "australia-southeast1", // Sydney
    "australia-southeast2", // Melbourne

    // Europe
    "europe-central2",      // Warsaw
    "europe-north1",        // Finland
    "europe-southwest1",    // Madrid
    "europe-west1",         // Belgium
    "europe-west2",         // London
    "europe-west3",         // Frankfurt
    "europe-west4",         // Netherlands
    "europe-west6",         // Zurich
    "europe-west8",         // Milan
    "europe-west9",         // Paris
    "europe-west10",        // Berlin
    "europe-west12",        // Turin

    // Middle East
    "me-central1",          // Doha
    "me-central2",          // Dammam
    "me-west1",             // Tel Aviv

    // North America
    "northamerica-northeast1", // Montreal
    "northamerica-northeast2", // Toronto
    "us-central1",          // Iowa
    "us-east1",             // South Carolina
    "us-east4",             // Virginia
    "us-east5",             // Columbus
    "us-south1",            // Dallas
    "us-west1",             // Oregon
    "us-west2",             // Los Angeles
    "us-west3",             // Salt Lake City
    "us-west4",             // Las Vegas

    // South America
    "southamerica-east1",   // São Paulo
    "southamerica-west1",   // Santiago
}
```

**Update Process**: When GCP adds new regions, update this list and bump the package version.

## Property Key Constants

Provider-specific property key constants for documentation and testing:

### AWS Property Keys

```go
const (
    // AWS SKU property keys (checked in priority order)
    AWSKeyInstanceType  = "instanceType"   // EC2 instances
    AWSKeyInstanceClass = "instanceClass"  // RDS instances
    AWSKeyType          = "type"           // Generic type
    AWSKeyVolumeType    = "volumeType"     // EBS volumes

    // AWS Region property keys (checked in priority order)
    AWSKeyRegion           = "region"           // Explicit region
    AWSKeyAvailabilityZone = "availabilityZone" // Derived from AZ
)
```

### Azure Property Keys

```go
const (
    // Azure SKU property keys (checked in priority order)
    AzureKeyVMSize = "vmSize"   // Virtual machines
    AzureKeySKU    = "sku"      // Generic SKU
    AzureKeyTier   = "tier"     // Service tier

    // Azure Region property keys (checked in priority order)
    AzureKeyLocation = "location" // Primary location field
    AzureKeyRegion   = "region"   // Alternative field
)
```

### GCP Property Keys

```go
const (
    // GCP SKU property keys (checked in priority order)
    GCPKeyMachineType = "machineType" // Compute instances
    GCPKeyType        = "type"        // Generic type
    GCPKeyTier        = "tier"        // Service tier

    // GCP Region property keys (checked in priority order)
    GCPKeyRegion = "region" // Explicit region
    GCPKeyZone   = "zone"   // Derived from zone
)
```

## State Transitions

N/A - The mapping package is stateless. All functions are pure with no side effects.

## Relationships

```text
┌─────────────────────┐
│    PropertyMap      │
│  map[string]string  │
└─────────┬───────────┘
          │
          │ input to
          ▼
┌─────────────────────┐     ┌─────────────────┐
│  Provider-Specific  │────▶│  string (SKU)   │
│  Extract*SKU()      │     └─────────────────┘
└─────────────────────┘
          │
          │ input to
          ▼
┌─────────────────────┐     ┌─────────────────┐
│  Provider-Specific  │────▶│ string (Region) │
│  Extract*Region()   │     └─────────────────┘
└─────────────────────┘
          │
          │ may call
          ▼
┌─────────────────────┐
│  ExtractAWSRegion   │
│     FromAZ()        │
│  ExtractGCPRegion   │
│     FromZone()      │
└─────────────────────┘
          │
          │ validates against
          ▼
┌─────────────────────┐
│    GCPRegions[]     │
│  (validation list)  │
└─────────────────────┘
```

## Validation Rules

### Input Validation

| Input | Behavior |
|-------|----------|
| `nil` map | Return empty string |
| Empty map | Return empty string |
| Missing keys | Return empty string |
| Empty value for key | Return empty string (treat as not found) |

### GCP Region Validation

| Input Zone | Derived Region | Valid? | Output |
|------------|----------------|--------|--------|
| `us-central1-a` | `us-central1` | Yes | `us-central1` |
| `europe-west1-b` | `europe-west1` | Yes | `europe-west1` |
| `invalid-zone-x` | `invalid-zone` | No | `` (empty) |
| `no-hyphen` | N/A | N/A | `` (empty) |
| `` (empty) | N/A | N/A | `` (empty) |

### AWS Availability Zone Extraction

| Input AZ | Output Region |
|----------|---------------|
| `us-east-1a` | `us-east-1` |
| `us-west-2b` | `us-west-2` |
| `eu-central-1c` | `eu-central-1` |
| `ap-northeast-1d` | `ap-northeast-1` |
| `` (empty) | `` (empty) |
