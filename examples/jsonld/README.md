# JSON-LD Examples

This directory contains example JSON-LD outputs for FOCUS cost data serialization.

## Examples

### focus_cost_record.jsonld

A single FocusCostRecord serialized to JSON-LD format. Demonstrates:

- Required `@context` with Schema.org and FOCUS vocabularies
- `@type` and `@id` for linked data semantics
- Cost fields serialized as Schema.org `MonetaryAmount`
- Tags as key-value pairs
- FOCUS 1.3 provider fields

### batch_output.jsonld

Multiple FocusCostRecords serialized as a JSON-LD array. Demonstrates:

- Batch serialization for enterprise pipelines
- Streaming output format
- Shared context across records

### contract_commitment.jsonld

A ContractCommitment record with linked data references. Demonstrates:

- Commitment dataset serialization
- Cross-record linking via `@id` references
- FOCUS 1.3 contract commitment fields

## Validation

Validate examples against JSON-LD specification:

```bash
# Using jsonld.js playground
# Visit: https://json-ld.org/playground/

# Or validate programmatically
go test ./sdk/go/jsonld/... -v
```

## Usage

See the [JSON-LD Serialization Package README](../../sdk/go/jsonld/README.md) for usage examples.
