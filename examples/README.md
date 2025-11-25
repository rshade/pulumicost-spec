# PulumiCost Examples

This directory contains examples of PulumiCost specification documents including
PricingSpec JSON documents and plugin registry index files.

## Plugin Registry Index

The **[registry.json](registry.json)** file demonstrates the plugin registry index
format used by `pulumicost plugin install` to discover and install plugins.

### Registry Format

```json
{
  "schema_version": "1.0.0",
  "plugins": {
    "plugin-name": {
      "name": "plugin-name",
      "description": "Human-readable description (10-500 chars)",
      "repository": "owner/repo",
      "author": "Author Name",
      "supported_providers": ["aws", "azure", "gcp", "kubernetes", "custom"],
      "min_spec_version": "0.1.0",
      "capabilities": ["cost_retrieval", "cost_projection"],
      "security_level": "community",
      "keywords": ["keyword1", "keyword2"]
    }
  }
}
```

### Validation

```bash
# Validate registry against schema
npm run validate:registry
```

For detailed contribution guidelines, see
[specs/001-plugin-registry-schema/quickstart.md](../specs/001-plugin-registry-schema/quickstart.md).

---

## PricingSpec Examples

This section contains comprehensive examples of PricingSpec JSON documents
demonstrating various cloud provider pricing models and billing patterns.

## Example Categories

### Cloud Provider Examples

#### AWS (Amazon Web Services)

- **[aws-ec2-t3-micro.json](specs/aws-ec2-t3-micro.json)**: On-demand EC2 instance pricing
- **[aws-s3-tiered-pricing.json](specs/aws-s3-tiered-pricing.json)**: S3 storage with tiered pricing structure  
- **[aws-lambda-per-invocation.json](specs/aws-lambda-per-invocation.json)**: Lambda serverless function pricing

#### Azure (Microsoft Azure)

- **[azure-vm-per-second.json](specs/azure-vm-per-second.json)**: Virtual machine with per-second billing
- **[azure-sql-dtu.json](specs/azure-sql-dtu.json)**: SQL Database with DTU-based pricing

#### GCP (Google Cloud Platform)

- **[gcp-storage-standard.json](specs/gcp-storage-standard.json)**: Cloud Storage with tiered pricing
- **[gcp-preemptible-spot.json](specs/gcp-preemptible-spot.json)**: Preemptible VM instances with spot pricing

#### Kubernetes

- **[kubernetes-namespace-cpu.json](specs/kubernetes-namespace-cpu.json)**: Kubernetes namespace CPU pricing via Kubecost

## Billing Models Demonstrated

### Time-Based Billing

- **Per Hour**: `aws-ec2-t3-micro.json` - Standard hourly compute pricing
- **Per Second**: `azure-vm-per-second.json` - Granular second-level billing
- **Per CPU Hour**: `kubernetes-namespace-cpu.json` - CPU resource-based pricing

### Usage-Based Billing  

- **Per Invocation**: `aws-lambda-per-invocation.json` - Function execution pricing
- **Per GB Month**: `aws-s3-tiered-pricing.json`, `gcp-storage-standard.json` - Storage capacity pricing
- **Per DTU**: `azure-sql-dtu.json` - Database transaction unit pricing

### Pricing Models

- **On-Demand**: Standard pay-as-you-go pricing
- **Preemptible/Spot**: Discounted interruptible instances with spot market pricing
- **Tiered Pricing**: Volume-based pricing tiers with discounts at higher usage levels

## Key Features Demonstrated

### 1. Basic Pricing Structure

All examples include the required fields:

- `provider`: Cloud provider identifier
- `resource_type`: Type of resource being priced  
- `billing_mode`: How the resource is billed
- `rate_per_unit`: Price per billing unit
- `currency`: Pricing currency (ISO 4217 code)

### 2. Metadata and Tags

Examples show comprehensive metadata usage:

- **Resource Tags**: Billing centers, cost centers, environments, applications
- **Plugin Metadata**: Provider-specific details, account IDs, technical specifications
- **Metric Hints**: Usage metrics for cost calculation with aggregation methods

### 3. Advanced Pricing Features

- **Tiered Pricing**: Volume-based pricing structures (`aws-s3-tiered-pricing.json`, `gcp-storage-standard.json`)
- **Commitment Terms**: Reserved/spot pricing with discounts (`gcp-preemptible-spot.json`)
- **Time Aggregation**: Rules for aggregating costs over time periods

## Example Details

### AWS EC2 t3.micro (`aws-ec2-t3-micro.json`)

**Billing Model**: `per_hour`  
**Rate**: $0.0104 USD per hour  
**Use Case**: Standard on-demand compute instance pricing  

Features demonstrated:

- Basic hourly billing for compute resources
- Resource tagging for cost allocation
- Provider-specific metadata (account ID, availability zone, instance specs)
- Multiple metric hints (vCPU hours, memory hours)

### Azure VM per-second (`azure-vm-per-second.json`)  

**Billing Model**: `per_second`  
**Rate**: $0.00000129 USD per second  
**Use Case**: Granular billing for short-running workloads  

Features demonstrated:

- Per-second billing precision
- Time aggregation rules (hourly windows)
- Azure-specific metadata (subscription ID, resource group)
- Minimum billing increment handling

### GCP Storage Standard (`gcp-storage-standard.json`)

**Billing Model**: `per_gb_month`  
**Rate**: $0.02 USD per GB-month (with tiers)  
**Use Case**: Object storage with volume-based pricing  

Features demonstrated:

- Tiered pricing structure with volume discounts
- Storage-specific metrics (GB storage, request counts)
- GCP project and bucket metadata
- Monthly billing alignment

### AWS Lambda (`aws-lambda-per-invocation.json`)

**Billing Model**: `per_invocation`  
**Rate**: $0.0000002 USD per invocation  
**Use Case**: Serverless function execution pricing  

Features demonstrated:

- Event-driven pricing model
- Free tier metadata and limits
- Function runtime and configuration details
- Duration and memory metric tracking

### GCP Preemptible Instance (`gcp-preemptible-spot.json`)

**Billing Model**: `preemptible`  
**Rate**: $0.01 USD per hour (79% discount)  
**Use Case**: Spot pricing for interruptible workloads  

Features demonstrated:

- Commitment terms with spot pricing
- Discount percentage calculation
- Preemptible-specific constraints (max runtime)
- Cost comparison with on-demand pricing

## Validation

All examples have been validated against the PricingSpec JSON Schema:

```bash
# Validate individual example
go run validate_examples.go examples/specs/aws-ec2-t3-micro.json

# Validate all examples
for file in examples/specs/*.json; do 
    go run validate_examples.go "$file"
done
```

## RPC Request Examples

The `requests/` directory contains sample gRPC request payloads demonstrating how to call various
CostSource RPC methods. These examples are useful for testing, documentation, and integration development.

### EstimateCost Requests

EstimateCost enables "what-if" cost analysis before deploying resources. The following examples
demonstrate cross-provider cost estimation:

- **[estimate_cost_aws.json](requests/estimate_cost_aws.json)**: AWS EC2 instance estimation
  - Resource: `aws:ec2/instance:Instance`
  - Attributes: Instance type (t3.micro), region, EBS volumes, monitoring, tags

- **[estimate_cost_azure.json](requests/estimate_cost_azure.json)**: Azure VM estimation
  - Resource: `azure:compute/virtualMachine:VirtualMachine`
  - Attributes: VM size (Standard_B1s), location, OS disk, operating system, tags

- **[estimate_cost_gcp.json](requests/estimate_cost_gcp.json)**: GCP Compute instance estimation
  - Resource: `gcp:compute/instance:Instance`
  - Attributes: Machine type (e2-micro), zone, boot disk, scheduling, labels

#### EstimateCost Request Format

```json
{
  "resource_type": "provider:module/resource:Type",
  "attributes": {
    "configuration_key": "value",
    "nested_object": {
      "key": "value"
    }
  }
}
```

The `attributes` field uses protobuf `Struct` format, allowing flexible key-value pairs matching
Pulumi resource input properties.

### GetActualCost Requests

Historical cost retrieval examples:

- **[actual-cost-namespace.json](requests/actual-cost-namespace.json)**: Kubernetes namespace actual costs

### GetProjectedCost Requests

Cost projection examples:

- **[projected-cost-ec2.json](requests/projected-cost-ec2.json)**: AWS EC2 instance projected costs

## Usage in Testing

These examples are used in:

- JSON Schema validation tests
- CI/CD pipeline validation
- SDK integration testing
- Documentation and demo purposes

## Contributing New Examples

When adding new examples:

1. **Follow naming convention**: `{provider}-{resource}-{unique-identifier}.json`
2. **Include comprehensive metadata**: Tags, plugin metadata, metric hints
3. **Validate against schema**: Ensure the example passes JSON Schema validation
4. **Document unique features**: Update this README with any new billing models or features
5. **Add realistic pricing**: Use current, realistic pricing data from providers

Example template:

```json
{
  "provider": "provider_name",
  "resource_type": "resource_type",
  "billing_mode": "billing_mode",
  "rate_per_unit": 0.00,
  "currency": "USD",
  "description": "Clear description of the pricing model",
  "metric_hints": [
    {
      "metric": "relevant_metric",
      "unit": "unit",
      "aggregation_method": "sum|avg|max|min|p95|p99"
    }
  ],
  "resource_tags": {
    "billing_center": "team_name",
    "cost_center": "CC-XXXX", 
    "environment": "production|development|staging"
  },
  "plugin_metadata": {
    "provider_specific_field": "value"
  },
  "source": "provider_name",
  "effective_date": "YYYY-MM-DDTHH:mm:ssZ"
}
```
