# Research: PluginSDK Mapping Package

**Feature**: 016-pluginsdk-mapping
**Date**: 2025-12-09

## Research Questions

### 1. AWS Availability Zone to Region Extraction

**Decision**: Strip trailing letter(s) from availability zone string

**Rationale**: AWS availability zones follow the pattern `{region}{zone-letter}` where region
is like `us-east-1` and zone-letter is `a`, `b`, `c`, etc. The extraction algorithm removes
the trailing alphabetic character(s) to derive the region.

**Alternatives considered**:

- Regex matching: More complex, potential edge cases with new region formats
- Lookup table: Requires maintenance as AWS adds regions
- String manipulation: Simple, handles all known and future formats

**Implementation**:

```go
func ExtractAWSRegionFromAZ(az string) string {
    if az == "" {
        return ""
    }
    // Remove trailing letters (a, b, c, etc.)
    for len(az) > 0 && az[len(az)-1] >= 'a' && az[len(az)-1] <= 'z' {
        az = az[:len(az)-1]
    }
    return az
}
```

### 2. GCP Zone to Region Extraction with Validation

**Decision**: Remove last hyphen-delimited segment and validate against known GCP regions list

**Rationale**: GCP zones follow the pattern `{region}-{zone-letter}` where region is like
`us-central1` and zone-letter is `a`, `b`, `c`, etc. Unlike AWS, the zone letter is separated
by a hyphen, making extraction straightforward via string manipulation. Validation against a
known regions list catches malformed inputs.

**Alternatives considered**:

- Regex matching: More complex than needed for hyphen-delimited format
- No validation: Would silently return invalid regions
- API-based validation: Adds external dependency, requires network

**Implementation**:

```go
var gcpRegions = []string{
    "asia-east1", "asia-east2", "asia-northeast1", "asia-northeast2", "asia-northeast3",
    "asia-south1", "asia-south2", "asia-southeast1", "asia-southeast2",
    "australia-southeast1", "australia-southeast2",
    "europe-central2", "europe-north1", "europe-southwest1", "europe-west1",
    "europe-west2", "europe-west3", "europe-west4", "europe-west6", "europe-west8",
    "europe-west9", "europe-west10", "europe-west12",
    "me-central1", "me-central2", "me-west1",
    "northamerica-northeast1", "northamerica-northeast2",
    "southamerica-east1", "southamerica-west1",
    "us-central1", "us-east1", "us-east4", "us-east5", "us-south1",
    "us-west1", "us-west2", "us-west3", "us-west4",
}

func ExtractGCPRegionFromZone(zone string) string {
    if zone == "" {
        return ""
    }
    lastHyphen := strings.LastIndex(zone, "-")
    if lastHyphen == -1 {
        return ""
    }
    region := zone[:lastHyphen]
    if !isValidGCPRegion(region) {
        return ""
    }
    return region
}
```

### 3. Property Key Priority Order

**Decision**: Provider-specific functions check keys in documented priority order

**Rationale**: Different Pulumi resource types use different property names for the same
concept. By checking multiple keys in a defined order, the extraction functions handle
various resource types without requiring caller knowledge of which key to check.

**AWS SKU Key Priority**:

1. `instanceType` - EC2 instances
2. `instanceClass` - RDS instances
3. `type` - Generic fallback
4. `volumeType` - EBS volumes

**AWS Region Key Priority**:

1. `region` - Explicit region setting
2. `availabilityZone` - Derived from AZ

**Azure SKU Key Priority**:

1. `vmSize` - Virtual machines
2. `sku` - Generic SKU field
3. `tier` - Service tier

**Azure Region Key Priority**:

1. `location` - Primary Azure location field
2. `region` - Alternative field name

**GCP SKU Key Priority**:

1. `machineType` - Compute instances
2. `type` - Generic type field
3. `tier` - Service tier

**GCP Region Key Priority**:

1. `region` - Explicit region
2. `zone` - Derived from zone

### 4. Zero-Allocation Pattern

**Decision**: Follow existing pluginsdk patterns with simple function design

**Rationale**: The registry package demonstrates zero-allocation validation using package-level
slices. For mapping functions, the primary allocation concern is string manipulation.
By using simple index-based operations instead of regex or string splitting with allocation,
we can achieve zero-allocation for most common cases.

**Performance Target**: <50 ns/op, 0 allocs/op (matching registry package performance)

### 5. Nil/Empty Input Handling

**Decision**: Return empty string for nil/empty inputs without panic

**Rationale**: Extraction functions are utility helpers that may receive data from various
sources. Defensive handling of nil/empty inputs prevents panic in edge cases and provides
consistent behavior. Empty string return indicates "no value found" which callers can check.

**Implementation Pattern**:

```go
func Extract*(props map[string]string, ...) string {
    if props == nil {
        return ""
    }
    // ... extraction logic
}
```

## Best Practices Research

### Go Package Design

- Follow existing `pluginsdk/` package naming convention
- Use `doc.go` for comprehensive package documentation
- Export only the public API; keep internal helpers unexported
- Table-driven tests for comprehensive coverage
- Separate benchmark tests in `*_benchmark_test.go`

### Cloud Provider Property Mapping

- AWS uses camelCase property names (`instanceType`, `availabilityZone`)
- Azure uses camelCase property names (`vmSize`, `location`)
- GCP uses camelCase property names (`machineType`, `zone`)
- All providers return strings for these properties in Pulumi

### SDK Consistency

- Match existing `env.go` function naming pattern (`Get*` prefix)
- Use `Extract*` prefix for mapping functions (distinguishes from getters)
- Document all public functions with godoc comments
- Include usage examples in package documentation

## Dependencies

No new external dependencies required. The mapping package uses only Go stdlib:

- `strings` - For string manipulation (LastIndex, etc.)

## Risks and Mitigations

| Risk | Mitigation |
|------|------------|
| GCP regions list becomes outdated | Include last-updated date; document refresh process |
| New property naming conventions | Generic extractors accept custom key lists |
| Performance regression | Benchmark tests with CI enforcement |
| Breaking pulumicost-core integration | Document migration path; maintain backward compat |
