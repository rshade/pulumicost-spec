# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is the **examples directory** of the PulumiCost Specification repository, containing comprehensive examples of PricingSpec
JSON documents and gRPC request samples. The examples serve as validation references, documentation resources, and integration
test data for the PulumiCost ecosystem.

## Directory Structure & Purpose

### Two-Tier Example System

**`specs/` Directory - PricingSpec Examples**

- **10 comprehensive examples** across 4 cloud providers (AWS, Azure, GCP, Kubernetes)
- **Cross-vendor billing model coverage**: 8+ different billing modes demonstrated
- **Real pricing data**: Current, realistic pricing from actual cloud providers
- **Schema validation targets**: All examples validate against `../schemas/pricing_spec.schema.json`

**`requests/` Directory - gRPC Request Samples**

- **Sample request payloads** for `GetActualCost` and `GetProjectedCost` RPCs
- **Integration testing data** for plugin development
- **Realistic request patterns** with proper field structures

## Build Commands

### Schema Validation

```bash
# Validate all examples against JSON Schema (from repo root)
npm run validate:examples

# Validate single example
go run validate_examples.go examples/specs/aws-ec2-t3-micro.json

# Validate all examples with detailed output
for file in specs/*.json; do go run ../validate_examples.go "$file"; done
```

### Example Testing Integration

```bash
# Run example validation in CI pipeline (from repo root)
make validate-examples

# Test examples in Go SDK integration tests
cd ../sdk/go && go test -v -run TestValidatePricingSpec
```

### Development Workflow

```bash
# Add new example and validate immediately
cp template.json specs/new-example.json
# Edit new example...
go run ../validate_examples.js  # Validates all examples

# Verify example integration
cd .. && npm run validate
```

## Example Architecture Patterns

### Billing Mode Coverage Strategy

**Time-Based Examples**:

- `aws-ec2-t3-micro.json` - `per_hour` with basic hourly billing
- `azure-vm-per-second.json` - `per_second` with granular billing precision
- `kubernetes-namespace-cpu.json` - `per_cpu_hour` with resource-based pricing

**Usage-Based Examples**:

- `aws-lambda-per-invocation.json` - `per_invocation` with event-driven pricing
- `aws-s3-tiered-pricing.json` - `per_gb_month` with volume-based tiers
- `azure-sql-dtu.json` - `per_dtu` with database transaction unit pricing

**Commitment-Based Examples**:

- `gcp-preemptible-spot.json` - `preemptible` with spot pricing and discount terms

### Feature Demonstration Matrix

**Basic Features (All Examples)**:

- Required fields: `provider`, `resource_type`, `billing_mode`, `rate_per_unit`, `currency`
- Resource tags for cost allocation (`billing_center`, `cost_center`, `environment`)
- Plugin metadata with provider-specific details
- Metric hints for usage calculation

**Advanced Features (Selected Examples)**:

- **Tiered Pricing**: `aws-s3-tiered-pricing.json`, `gcp-storage-standard.json` - Volume discounts
- **Time Aggregation**: `azure-vm-per-second.json` - Billing window rules
- **Commitment Terms**: `gcp-preemptible-spot.json` - Discount percentage and duration

### Metadata Pattern Architecture

**Standardized Resource Tags**:

- `billing_center` - Cost allocation team (e.g., "engineering", "marketing")
- `cost_center` - Accounting code (e.g., "CC-1001", "CC-2002")
- `environment` - Deployment stage ("production", "development", "staging")
- Domain-specific tags vary by use case and provider

**Provider-Specific Plugin Metadata**:

- **AWS**: `aws_account_id`, `availability_zone`, `instance_family`, technical specs
- **Azure**: `subscription_id`, `resource_group`, service-specific configurations
- **GCP**: `project_id`, `zone`, commitment and pricing model details
- **Kubernetes**: `cluster_name`, `namespace_labels`, `kubecost_version`

### Cross-Provider Consistency

**Naming Conventions**:

- File format: `{provider}-{resource}-{distinguisher}.json`
- Resource types: Standardized across providers (`ec2`/`vm`/`compute_engine`)
- Currency: Consistent USD pricing for comparison
- Effective dates: Aligned timestamp formats

**Validation Integration**:

- All examples pass JSON Schema validation via `validate_examples.js`
- Integrated into CI/CD pipeline through `npm run validate:examples`
- Go SDK validation testing via `sdk/go/pricing/validate_test.go`
- Real-time validation during development workflow

## Example Selection Guide

### When Adding New Examples

**Choose Representative Scenarios**:

- Major cloud provider services with significant usage patterns
- Different billing models not yet demonstrated
- Advanced features requiring comprehensive documentation
- Common cost allocation and tagging scenarios

**Required Validation Steps**:

1. **Schema Compliance**: Must pass `validate_examples.js` validation
2. **Realistic Data**: Use current, accurate pricing from cloud providers
3. **Complete Metadata**: Include comprehensive resource tags and plugin metadata
4. **Metric Hints**: Provide relevant usage metrics with aggregation methods
5. **Documentation**: Update README.md with new example details

### Example Integration Points

**With Go SDK Testing**:

- Examples used in `sdk/go/pricing/validate_test.go` for comprehensive validation testing
- Cross-provider test coverage via realistic example data
- Schema validation integration testing

**With CI/CD Pipeline**:

- `make validate-examples` runs validation on all examples
- Breaking changes in schema detected via example validation failures
- Continuous validation ensures example-schema synchronization

**With Plugin Development**:

- Examples provide realistic `PricingSpec` response templates
- `requests/` directory samples provide gRPC request pattern guidance
- Cross-provider compatibility testing via diverse example set

## Validation Architecture

### Two-Tier Validation System

**JSON Schema Validation** (`validate_examples.js`):

- AJV-based validation with strict schema compliance
- Comprehensive error reporting with field-level details
- Integration with npm validation pipeline
- Automated CI/CD validation gates

**Go SDK Integration Validation**:

- Runtime validation via `pricing.ValidatePricingSpec()` function
- Cross-language validation consistency verification
- Integration test coverage across all example scenarios

### Example Quality Standards

**Data Accuracy Requirements**:

- Pricing rates must reflect current cloud provider pricing (within reasonable staleness)
- Resource specifications must match actual provider offerings
- Metadata fields must represent realistic deployment scenarios
- Currency and date formats must follow ISO standards

**Schema Compliance Requirements**:

- All required fields present and properly typed
- Advanced features properly structured (tiers, aggregation, terms)
- Provider enum values consistent with schema definitions
- Billing mode values from approved specification list

This examples directory serves as the validation foundation and integration reference for the entire PulumiCost ecosystem,
ensuring consistent pricing specification interpretation across all implementations.
