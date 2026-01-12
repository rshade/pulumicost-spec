# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is the **pricing package** of the FinFocus Go SDK, containing core domain types, validation logic, and billing mode
enumerations. The package provides the foundational types and validation for FinFocus pricing specifications across all
cloud providers.

## Core Components

### Domain Types (`domain.go`)

**BillingMode Enumeration System**:

- **38 comprehensive billing modes** organized by category:
  - **Time-based**: `per_hour`, `per_minute`, `per_second`, `per_day`, `per_month`, `per_year`
  - **Storage-based**: `per_gb_month`, `per_gb_hour`, `per_gb_day`
  - **Usage-based**: `per_request`, `per_operation`, `per_transaction`, `per_execution`, `per_invocation`,
    `per_api_call`, `per_lookup`, `per_query`
  - **Compute-based**: `per_cpu_hour`, `per_cpu_month`, `per_vcpu_hour`, `per_memory_gb_hour`,
    `per_memory_gb_month`
  - **I/O-based**: `per_iops`, `per_provisioned_iops`, `per_data_transfer_gb`, `per_bandwidth_gb`
  - **Database-specific**: `per_rcu` (DynamoDB), `per_wcu` (DynamoDB), `per_dtu` (Azure), `per_ru` (Cosmos DB)
  - **Pricing models**: `on_demand`, `reserved`, `spot`, `preemptible`, `savings_plan`, `committed_use`,
    `hybrid_benefit`, `flat`

**Provider Enumeration**:

- **5 supported providers**: `aws`, `azure`, `gcp`, `kubernetes`, `custom`
- Validation functions for both billing modes and providers
- String conversion methods for all enum types

**Key Functions**:

- `ValidBillingMode(string) bool` - Validates billing mode strings
- `GetAllBillingModes() []string` - Returns all valid billing modes
- `ValidProvider(string) bool` - Validates provider strings
- `GetAllProviders() []Provider` - Returns all valid providers

### JSON Schema Validation (`validate.go`)

**Embedded Schema System**:

- **Complete JSON Schema** embedded as `schemaJSON` constant (600+ lines)
- Synchronized with `../../../schemas/pricing_spec.schema.json`
- Comprehensive validation covering all PricingSpec fields and nested structures
- **Runtime validation** via `ValidatePricingSpec([]byte) error`

**Schema Features**:

- **Required fields**: provider, resource_type, billing_mode, rate_per_unit, currency
- **Advanced structures**: metric_hints, pricing_tiers, time_aggregation, commitment_terms
- **Flexible metadata**: resource_tags (string-only), plugin_metadata (any type)
- **ISO compliance**: 3-character currency codes, ISO 8601 date formats
- **Validation rules**: Non-negative rates, enum constraints, string length limits

## Build Commands

### Testing

```bash
# Run package tests from this directory
go test

# Run specific test categories
go test -run TestValidBillingMode    # Domain validation tests
go test -run TestValidatePricingSpec # JSON schema validation tests

# Run tests with coverage
go test -cover
go test -coverprofile=coverage.out && go tool cover -html=coverage.out

# Run from repository root (recommended)
cd ../../../ && make test
```

### Development

```bash
# Build package
go build

# Check formatting and imports
go fmt
go mod tidy

# Validate from repository root
cd ../../../ && make lint && make validate
```

## Architecture Patterns

### Billing Mode Validation Strategy

**Current Status**: Uses function-returned slices (not yet optimized)

The package uses a **centralized enumeration approach**:

1. **Constants Definitions**: All billing modes defined as typed constants
2. **Validation Registry**: `getAllBillingModes()` maintains canonical list
3. **String Conversion**: Type-safe conversion via `.String()` methods
4. **Validation Functions**: Multiple validation entry points for different use cases

```go
// Current usage patterns
if !pricing.ValidBillingMode("per_hour") {
    return errors.New("invalid billing mode")
}

allModes := pricing.GetAllBillingModes() // Get all 44 modes
```

**Performance Optimization Opportunity** 🔧:

The pricing package currently uses function-returned slices which allocate memory on each call. The registry package
has been optimized to use package-level variables for zero-allocation validation (see
`specs/001-domain-enum-optimization/validation-pattern.md`).

**Recommended Future Optimization**:

Apply the same pattern as registry package:

```go
// Recommended pattern (not yet implemented)
//nolint:gochecknoglobals // Intentional optimization for zero-allocation validation
var allBillingModes = []BillingMode{
    PerHour, PerMinute, PerSecond, /* ... all 44 values ... */
}

func getAllBillingModes() []BillingMode {
    return allBillingModes  // Returns reference, zero allocation
}

func ValidBillingMode(mode string) bool {
    billingMode := BillingMode(mode)
    for _, valid := range allBillingModes {  // Direct slice access
        if billingMode == valid {
            return true
        }
    }
    return false
}
```

**Expected Performance** (based on registry package scaling):

- Current: ~40-60 ns/op with allocation overhead
- Optimized: ~20-25 ns/op (44 values × 0.5 ns/value + 5 ns base)
- Memory: 0 allocs/op (vs current 1 alloc/op)

**Recommendation**: Apply optimization in future PR when performance becomes critical or for consistency with registry
package pattern.

### Schema Validation Architecture

The validation system uses **embedded schema compilation**:

1. **Schema Embedding**: JSON Schema embedded as string constant
2. **Runtime Compilation**: Schema compiled to validator on first use
3. **Validation Pipeline**: JSON unmarshaling → schema validation
4. **Error Propagation**: Detailed validation errors with field context

```go
// Validation usage
err := pricing.ValidatePricingSpec(jsonBytes)
if err != nil {
    // Handle validation error with specific field information
}
```

### Test Organization Strategy

Tests are organized by **functional categories**:

- **Domain tests** (`domain_test.go`): Enum validation, string conversion, completeness
- **Schema validation** (`validate_test.go`): Valid examples across all providers
- **Invalid input tests** (`validate_invalid_test.go`): Error condition coverage
- **Resource tags tests** (`resource_tags_test.go`): Tag validation and metadata handling

**Test Coverage**:

- **Cross-provider examples**: AWS, Azure, GCP, Kubernetes scenarios
- **Complex structures**: Multi-tier pricing, commitment terms, metadata
- **Error conditions**: Missing fields, invalid formats, constraint violations
- **Edge cases**: Empty objects, boundary values, format variations

## Common Development Patterns

### Adding New Billing Modes

1. **Add constant** in appropriate category block in `domain.go`
2. **Update `getAllBillingModes()`** function to include new mode
3. **Update embedded schema** in `validate.go` (billing_mode enum array)
4. **Add test case** in `domain_test.go` for the new mode
5. **Update expected count** in `TestAllBillingModesCompleteness`

### Schema Updates

1. **Modify embedded `schemaJSON`** constant in `validate.go`
2. **Ensure synchronization** with `../../../schemas/pricing_spec.schema.json`
3. **Add validation tests** for new fields in `validate_test.go`
4. **Add invalid cases** in `validate_invalid_test.go`
5. **Verify completeness** by running full test suite

### Provider Extensions

1. **Add provider constant** to Provider enum in `domain.go`
2. **Update `GetAllProviders()`** to include new provider
3. **Update schema enum** in embedded JSON schema
4. **Add test coverage** in `domain_test.go` and validation tests
5. **Update expected count** in provider completeness test

## Test Execution Patterns

### Schema Validation Testing

The package includes **comprehensive validation test suites**:

```bash
# Valid examples across providers
go test -run TestValidatePricingSpec_ValidAWSExamples
go test -run TestValidatePricingSpec_ValidAzureExamples
go test -run TestValidatePricingSpec_ValidGCPExamples
go test -run TestValidatePricingSpec_ValidKubernetesExamples

# Invalid input validation
go test -run TestValidatePricingSpec_InvalidSchemas
go test -run TestValidatePricingSpec_InvalidMetricHints
go test -run TestValidatePricingSpec_InvalidPricingTiers

# Complex scenarios
go test -run TestValidatePricingSpec_ComplexValidExamples
go test -run TestValidatePricingSpec_CombinedTagsAndMetadata
```

### Domain Validation Testing

```bash
# Billing mode validation
go test -run TestValidBillingMode
go test -run TestBillingModeString
go test -run TestAllBillingModesCompleteness

# Provider validation
go test -run TestValidProvider
go test -run TestProviderString
go test -run TestAllProvidersCompleteness
```

## Key Design Decisions

### Type Safety Approach

- **Typed enums** (`BillingMode`, `Provider`) prevent string errors
- **Validation functions** provide runtime checking for external input
- **String conversion** maintains compatibility with JSON/protobuf

### Schema Synchronization

- **Embedded schema** ensures validation consistency
- **Single source of truth** in `schemas/pricing_spec.schema.json`
- **Runtime compilation** balances performance and flexibility

### Test Coverage Strategy

- **Real-world examples** from all major cloud providers
- **Comprehensive error cases** covering all validation paths
- **Structural completeness** ensuring enum coverage and consistency

### Extensibility Patterns

- **Category-based organization** for billing modes
- **Flexible metadata support** via `additionalProperties: true`
- **Provider-agnostic design** supporting custom implementations
