# GetBudgets RPC Contracts

This directory contains example request and response payloads for the GetBudgets RPC method.

## Files

### Request Examples

- `get-budgets-request-aws.json` - AWS provider budget query
- `get-budgets-request-gcp.json` - GCP provider budget query
- `get-budgets-request-kubecost.json` - Kubecost provider budget query

### Response Examples

- `get-budgets-response-aws.json` - AWS budget data with status
- `get-budgets-response-gcp.json` - GCP budget data with status
- `get-budgets-response-kubecost.json` - Kubecost budget data with status

## Usage

These examples demonstrate:

- Cross-provider budget data unification
- Optional status inclusion via `include_status` flag
- Provider-specific metadata preservation
- Threshold triggering and health status calculation
- Currency handling and percentage calculations

## Validation

All examples validate against the JSON schema derived from the protobuf definitions and pass buf linting.
