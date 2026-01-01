# Dry Run Examples

This directory contains example DryRun request/response payloads for different cloud providers
and resource types. These examples demonstrate how to use the dry-run capability to discover
plugin field mappings without performing actual cost data retrieval.

## Examples

### AWS EC2 Instance (`aws_ec2.json`)

Demonstrates dry-run for AWS EC2 instances.

**Key highlights:**

- All standard FOCUS fields are marked as SUPPORTED
- `availability_zone` is CONDITIONAL (only for multi-AZ VPCs)
- `resource_name` is CONDITIONAL (requires Name tag)
- `sku_price_id` is UNSUPPORTED (AWS doesn't expose this)

### Azure Virtual Machine (`azure_vm.json`)

Demonstrates dry-run for Azure VM resources.

**Key highlights:**

- All standard FOCUS fields are marked as SUPPORTED
- `availability_zone` is UNSUPPORTED (Azure uses availability sets, not zones)
- `resource_name` is CONDITIONAL (requires name tag)

### GCP Compute Engine (`gcp_compute.json`)

Demonstrates dry-run for GCP Compute Engine instances.

**Key highlights:**

- All standard FOCUS fields are marked as SUPPORTED
- `availability_zone` is CONDITIONAL (regional resources don't have zones)
- `resource_name` is CONDITIONAL (requires instance name)
- `sku_price_id` is UNSUPPORTED (GCP doesn't expose this)

### Kubernetes Pod (`k8s_pod.json`)

Demonstrates dry-run for Kubernetes Pod resources.

**Key highlights:**

- Many billing/account fields are UNSUPPORTED (Kubernetes cluster-scoped)
- `billed_cost` is CONDITIONAL (depends on cloud provider billing integration)
- `pricing_*` fields are UNSUPPORTED (Kubernetes doesn't define pricing directly)
- `service_provider_name` and `host_provider_name` are SUPPORTED/CONDITIONAL (for managed services like GKE/EKS)
- Cost allocation fields are CONDITIONAL (requires cluster configuration)

## Usage Patterns

### Cross-Provider Comparison

Query multiple plugins for the same resource type to compare capabilities:

```bash
# Query AWS plugin for EC2 field mappings
grpcurl -plaintext localhost:50051 pulumicost.v1.CostSourceService/DryRun \
  -d '{"resource":{"provider":"aws","resource_type":"ec2"}}'

# Query Azure plugin for VM field mappings
grpcurl -plaintext localhost:50052 pulumicost.v1.CostSourceService/DryRun \
  -d '{"resource":{"provider":"azure","resource_type":"vm"}}'

# Compare responses to understand provider differences
```

### Simulation Parameters

Use `simulation_parameters` to test different scenarios:

```json
{
  "resource": {
    "provider": "aws",
    "resource_type": "ec2"
  },
  "simulation_parameters": {
    "pricing_tier": "reserved",
    "commitment_term": "1-year"
  }
}
```

Plugins may use simulation parameters to adjust field statuses dynamically. For example,
a "reserved" pricing tier might change certain cost-related fields from DYNAMIC to SUPPORTED.

### Using dry_run Flag on Cost RPCs

Instead of calling DryRun RPC directly, you can use the `dry_run` flag on existing cost RPCs:

```bash
# GetActualCost with dry-run
grpcurl -plaintext localhost:50051 \
  pulumicost.v1.CostSourceService/GetActualCost \
  -d '{
    "resource_id": "i-1234567890",
    "start": "2024-01-01T00:00:00Z",
    "end": "2024-01-02T00:00:00Z",
    "dry_run": true
  }'

# GetProjectedCost with dry-run
grpcurl -plaintext localhost:50051 \
  pulumicost.v1.CostSourceService/GetProjectedCost \
  -d '{
    "resource": {
      "provider": "aws",
      "resource_type": "ec2",
      "region": "us-east-1"
    },
    "dry_run": true
  }'
```

### Check Plugin Capability Before Querying

Always check if plugin supports dry-run before making requests:

```bash
# Check capabilities
grpcurl -plaintext localhost:50051 \
  pulumicost.v1.CostSourceService/Supports \
  -d '{
    "resource": {
      "provider": "aws",
      "resource_type": "ec2"
    }
  }'

# Verify "dry_run": true in capabilities response
```

If `capabilities["dry_run"]` is false or missing, the plugin doesn't support dry-run
introspection. Fall back to using the Supports RPC to check basic resource type support.

## Field Status Reference

| Status      | Meaning                                 | When to Expect                                                      |
| ----------- | --------------------------------------- | ------------------------------------------------------------------- |
| SUPPORTED   | Always populated for this resource type | Core identity, billing, and charge fields                           |
| UNSUPPORTED | Never populated for this resource type  | Provider-specific fields, irrelevant fields                         |
| CONDITIONAL | Depends on resource configuration       | Tags, zones, names - check `condition_description`                  |
| DYNAMIC     | Requires runtime data                   | Actual cost values, computed fields - check `condition_description` |

## Troubleshooting

### Empty field_mappings

If `field_mappings` is empty or `resource_type_supported` is false:

- Verify `provider` and `resource_type` values match what the plugin expects
- Check if the resource type is supported by this plugin
- Query the Supports RPC to verify basic capability

### configuration_valid = false

If `configuration_valid` is false:

- Check `configuration_errors` array for specific issues
- Common issues: missing API keys, invalid endpoints, network connectivity
- Fix configuration issues and retry

### Performance Issues

If dry-run requests exceed 100ms:

- Plugin may be making external API calls (violates spec)
- Check plugin implementation for proper dry-run handling
- Dry-run should be stateless and fast (no external data retrieval)

## Integration Testing

Use these examples to test plugin implementations:

```bash
# Test AWS EC2 example
jq -r '.request' examples/requests/dry_run/aws_ec2.json | \
  grpcurl -plaintext localhost:50051 \
    pulumicost.v1.CostSourceService/DryRun \
    -d @- > /tmp/response.json

# Validate response against expected
jq -r '.response' examples/requests/dry_run/aws_ec2.json > /tmp/expected.json
diff /tmp/expected.json /tmp/response.json
```

This pattern can be automated in CI/CD pipelines for conformance testing across multiple providers.
