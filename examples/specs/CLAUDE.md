# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is the **specs directory** containing the core PricingSpec JSON examples for the FinFocus specification. These 10
production-quality examples serve as the canonical reference implementations demonstrating all major cloud provider billing
models, advanced pricing features, and comprehensive metadata patterns.

## Example Architecture Matrix

### Cross-Provider Billing Model Coverage

**AWS Examples (3)**:

- `aws-ec2-t3-micro.json` - `per_hour` with basic compute pricing
- `aws-lambda-per-invocation.json` - `per_invocation` with serverless event-driven pricing
- `aws-s3-tiered-pricing.json` - `per_gb_month` with volume-based tiered pricing

**Azure Examples (2)**:

- `azure-vm-per-second.json` - `per_second` with granular time-based billing
- `azure-sql-dtu.json` - `per_dtu` with database transaction unit pricing

**GCP Examples (2)**:

- `gcp-storage-standard.json` - `per_gb_month` with tiered storage pricing
- `gcp-preemptible-spot.json` - `preemptible` with commitment terms and discount pricing

**Kubernetes Examples (1)**:

- `kubernetes-namespace-cpu.json` - `per_cpu_hour` with resource-based pricing via Kubecost

**Plus YAML Variants (2)**:

- `aws-ec2-t3-micro.yaml` - YAML format demonstration
- `kubecost-gke-node.yaml` - Node-level Kubernetes pricing

### Advanced Feature Demonstration

**Commitment Terms Pricing**:

- `gcp-preemptible-spot.json` - Spot pricing with 79% discount, duration constraints, payment options

**Tiered Pricing Structures**:

- `aws-s3-tiered-pricing.json` - 3-tier volume pricing (0-50K, 50K-450K, 450K+ GB)
- `gcp-storage-standard.json` - Multi-tier storage with progressive discounts

**Time Aggregation Rules**:

- `azure-vm-per-second.json` - Per-second billing with hourly aggregation windows
- `aws-lambda-per-invocation.json` - Monthly aggregation with billing alignment
- `kubernetes-namespace-cpu.json` - Continuous alignment with average aggregation

**Free Tier Integration**:

- `aws-lambda-per-invocation.json` - Free tier metadata (1M invocations, runtime specifications)

## Build Commands

### Validation Commands

```bash
# Validate single example (from repo root)
go run validate_examples.go examples/specs/aws-ec2-t3-micro.json

# Validate all examples with comprehensive output
for file in *.json; do echo "Validating $file..."; go run ../../validate_examples.js; done

# Schema validation via npm (from repo root)
npm run validate:examples
```

### Development Workflow

```bash
# Create new example from template pattern
cp aws-ec2-t3-micro.json new-provider-resource.json
# Edit with provider-specific details...

# Immediate validation
go run ../../validate_examples.js

# Integration with full validation pipeline
cd ../.. && make validate
```

### Example Testing Integration

```bash
# Test examples in Go SDK (from repo root)
cd sdk/go && go test -v -run TestValidatePricingSpec_ValidAWSExamples
cd sdk/go && go test -v -run TestValidatePricingSpec_ComplexValidExamples
```

## Metadata Pattern Architecture

### Standardized Resource Tags System

**Universal Tag Schema**:

- `billing_center` - Team/department for cost allocation (engineering, platform, data-platform)
- `cost_center` - Accounting code pattern (`CC-XXXX` format, sequential numbering)
- `environment` - Deployment stage (production, development, staging)
- Domain-specific tags vary by resource type and use case

**Cost Center Numbering Pattern**:

- CC-1001: Engineering general compute
- CC-2002: IT infrastructure services
- CC-4004: Platform serverless functions
- CC-5005: ML engineering batch processing
- CC-6006: Platform production workloads
- CC-7007: Data platform production databases

### Provider-Specific Plugin Metadata

**AWS Metadata Pattern**:

- `aws_account_id` - Always 12-digit account identifier
- `availability_zone` - Specific AZ within region (us-east-1a)
- Service-specific: `function_runtime`, `memory_size_mb`, `free_tier_eligible`

**Azure Metadata Pattern**:

- `azure_subscription_id` - UUID format subscription identifier
- `resource_group` - Logical grouping (dev-web-servers, prod-databases)
- Service-specific: `vm_size_family`, `dtu_capacity`, `service_tier`

**GCP Metadata Pattern**:

- `gcp_project_id` - Project identifier (my-project-123456)
- `machine_family` - Instance family (general-purpose)
- Pricing-specific: `preemptible`, `max_run_duration_hours`, `original_on_demand_price`

**Kubernetes/Kubecost Pattern**:

- `cluster_name` - Kubernetes cluster identifier
- `kubecost_version` - Cost analysis tool version
- `cost_allocation_method` - Allocation strategy (proportional)
- `namespace_labels` - Label selectors for resource grouping

### Metric Hints Architecture

**Time-Based Resources** (EC2, VM, Compute Engine):

- `vcpu_hours`/`vcpu_seconds` - CPU time consumption
- `memory_gb_hours`/`memory_gb_seconds` - Memory allocation time
- Aggregation: `sum` for cumulative resource time

**Usage-Based Resources** (Lambda, Functions):

- `invocations` - Event count with `sum` aggregation
- `duration_ms` - Execution time with `sum` aggregation
- `memory_mb` - Memory allocation with `avg` aggregation

**Storage Resources** (S3, Cloud Storage):

- `storage_gb` - Capacity with `avg` aggregation
- `get_requests`/`put_requests` - API operations with `sum` aggregation

**Database Resources** (SQL Database):

- `database_transaction_units` - DTU usage with `avg` aggregation
- `storage_gb` - Database size with `max` aggregation
- `cpu_percent` - Resource utilization with `avg` aggregation

## Example Quality Standards

### Data Accuracy Requirements

**Realistic Pricing Data**:

- AWS EC2 t3.micro: $0.0104/hour (current on-demand pricing)
- Azure B1s VM: $0.00000129/second ($0.00464/hour equivalent)
- GCP n1-standard-1 preemptible: $0.01/hour (79% discount from $0.0475 on-demand)
- Lambda: $0.0000002/invocation (current AWS pricing)

**Technical Specifications**:

- Memory, CPU, and storage specifications match actual provider offerings
- SKU identifiers correspond to real provider resource types
- Regional pricing variations reflected accurately

### Schema Compliance Architecture

**Required Field Validation**:

- All examples pass `../../schemas/pricing_spec.schema.json` validation
- Provider enum values: aws, azure, gcp, kubernetes
- Currency: USD for cross-provider comparison consistency
- ISO 8601 timestamp format for `effective_date`

**Advanced Feature Compliance**:

- Tiered pricing: min_units, max_units properly structured
- Commitment terms: duration, payment_option, discount_percentage valid
- Time aggregation: window/method/alignment combinations validated
- Metric hints: metric/unit/aggregation_method triples complete

This specs directory provides the foundational validation reference for the entire FinFocus ecosystem, ensuring consistent
pricing specification interpretation across all cloud providers and resource types.
