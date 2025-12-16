# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is the **schemas directory** containing the canonical JSON Schema definition for PricingSpec validation in the PulumiCost
ecosystem. The single file `pricing_spec.schema.json` serves as the authoritative specification for all pricing data structures
across cloud providers and billing models.

## Schema Architecture

### Core Schema Design (`pricing_spec.schema.json`)

**Schema Specification**:

- **JSON Schema Draft 2020-12** compliance for modern validation features
- **Schema ID**: `https://spec.pulumicost.dev/schemas/pricing_spec.schema.json`
- **Version**: v0.4.6 production-ready specification
- **Validation Mode**: Strict with `additionalProperties: false` for most objects

**Required Fields Foundation**:

- `provider` - Cloud provider enum (aws, azure, gcp, kubernetes, custom)
- `resource_type` - Resource identifier (ec2, s3, vm, k8s-namespace)
- `billing_mode` - How resource is charged (44+ supported modes)
- `rate_per_unit` - Numeric pricing rate (non-negative)
- `currency` - ISO 4217 3-character currency code

### Billing Mode Enumeration System

**44 Comprehensive Billing Modes**:

**Time-Based Billing**:

- `per_hour`, `per_minute`, `per_second`
- `per_day`, `per_month`, `per_year`
- `per_cpu_hour`, `per_cpu_month`, `per_vcpu_hour`
- `per_memory_gb_hour`, `per_memory_gb_month`

**Usage-Based Billing**:

- `per_request`, `per_operation`, `per_transaction`
- `per_execution`, `per_invocation`
- `per_api_call`, `per_lookup`, `per_query`

**Storage-Based Billing**:

- `per_gb_month`, `per_gb_hour`, `per_gb_day`
- `per_iops`, `per_provisioned_iops`
- `per_data_transfer_gb`, `per_bandwidth_gb`

**Database-Specific Billing**:

- `per_rcu`, `per_wcu` (DynamoDB)
- `per_dtu` (Azure SQL)
- `per_ru` (Cosmos DB)

**Pricing Models**:

- `on_demand`, `reserved`, `spot`, `preemptible`
- `savings_plan`, `committed_use`, `hybrid_benefit`
- `flat` - Fixed pricing

### Advanced Features Architecture

**Tiered Pricing Structure**:

```json
"pricing_tiers": [
  {
    "min_units": 0,
    "max_units": 50000,
    "rate_per_unit": 0.023
  },
  {
    "min_units": 50000,
    "rate_per_unit": 0.021
  }
]
```

- Volume-based pricing with progressive discounts
- Optional `max_units` for highest tier (unlimited)
- Multiple tier support for complex pricing models

**Time Aggregation Rules**:

```json
"time_aggregation": {
  "window": "hour|day|month",
  "method": "sum|avg|prorated",
  "alignment": "calendar|billing|continuous"
}
```

- Flexible time window definitions
- Aggregation method specification
- Alignment strategy for billing periods

**Commitment Terms Modeling**:

```json
"commitment_terms": {
  "duration": "1_year|3_year|spot|on_demand",
  "payment_option": "all_upfront|partial_upfront|no_upfront|monthly",
  "discount_percentage": 79.0
}
```

- Reserved instance and savings plan support
- Payment schedule flexibility
- Discount calculation integration

### Metadata Architecture Systems

**Metric Hints System**:

```json
"metric_hints": [
  {
    "metric": "vcpu_hours",
    "unit": "hour",
    "aggregation_method": "sum|avg|max|min|p95|p99"
  }
]
```

- Usage calculation guidance for cost engines
- Standardized aggregation methods
- Flexible metric/unit combinations

**Resource Tags Standardization**:

```json
"resource_tags": {
  "billing_center": "engineering",
  "cost_center": "CC-1001",
  "environment": "production"
}
```

- String-only values for consistent querying
- Cost allocation and categorization support
- Standardized tag naming conventions

**Plugin Metadata Flexibility**:

```json
"plugin_metadata": {
  "additionalProperties": true
}
```

- Provider-specific data support
- Flexible schema for vendor variations
- Complex nested object support

## Build Commands

### Schema Validation

```bash
# Validate schema syntax (from repo root)
npm run validate:schema

# Validate all examples against schema
npm run validate:examples

# Direct AJV validation with strict mode disabled
npx ajv validate --strict=false -s schemas/pricing_spec.schema.json -d examples/specs/*.json
```

### Development Workflow

```bash
# Test schema changes with existing examples
for file in examples/specs/*.json; do
    echo "Validating $file against updated schema..."
    npx ajv validate --strict=false -s schemas/pricing_spec.schema.json -d "$file"
done

# Comprehensive validation pipeline
make validate-schema && make validate-examples
```

### Integration Testing

```bash
# Go SDK schema validation integration (from repo root)
cd sdk/go && go test -v -run TestValidatePricingSpec

# Test schema synchronization
cd sdk/go/pricing && go test -v -run TestEmbeddedSchemaSync
```

## Schema Evolution Patterns

### Backward Compatibility Strategy

**Additive Changes (Safe)**:

- New optional fields in existing objects
- New enum values in billing_mode
- Additional provider types
- Extended metadata properties

**Breaking Changes (Version Bump Required)**:

- New required fields
- Removal of enum values
- Field type changes
- additionalProperties changes

### Provider Extension Pattern

**Adding New Provider**:

1. Add provider to enum: `"provider": {"enum": ["aws", "azure", "gcp", "kubernetes", "new_provider", "custom"]}`
2. Update examples with new provider pricing
3. Validate all existing examples still pass
4. Update Go SDK provider validation

**Adding New Billing Mode**:

1. Add to billing_mode enum array
2. Create representative example in `../examples/specs/`
3. Update Go SDK billing mode validation
4. Add to conformance testing scenarios

### Schema Validation Integration

**Multi-Language Validation**:

- **JavaScript/Node.js**: AJV with format validation
- **Go**: Embedded schema in `sdk/go/pricing/validate.go`
- **CI/CD Pipeline**: npm script integration with make targets
- **Plugin Development**: Runtime validation in testing framework

**Validation Error Handling**:

- Field-level error reporting with JSON paths
- Comprehensive error messages for debugging
- Format validation (ISO dates, currency codes, etc.)
- Constraint validation (minimums, patterns, enums)

### Schema Synchronization Architecture

**Single Source of Truth**:

- Schema file serves as canonical definition
- Go SDK embeds schema as string constant
- JavaScript validation uses schema directly
- Examples validated against current schema

**Change Propagation**:

1. **Schema Update**: Modify `pricing_spec.schema.json`
2. **Go SDK Sync**: Update embedded schema in `sdk/go/pricing/validate.go`
3. **Example Validation**: Ensure all examples pass new schema
4. **Integration Testing**: Verify cross-language validation consistency

This schema serves as the foundation for consistent pricing specification validation across the entire PulumiCost ecosystem,
ensuring data integrity and interoperability across all implementations and cloud providers.
