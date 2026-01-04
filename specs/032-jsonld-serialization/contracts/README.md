# JSON-LD Contracts

This directory contains JSON-LD context definitions for FOCUS data serialization.

## Files

### focus-context.jsonld

The canonical JSON-LD 1.1 context for FOCUS cost data. This context:

- Defines the `focus:` namespace for FOCUS-specific terms
- Integrates `schema:` (Schema.org) for common types
- Maps all 66 FocusCostRecord fields to JSON-LD properties
- Maps all 12 ContractCommitment fields
- Declares type coercions (xsd:dateTime, xsd:decimal)
- Marks deprecated fields with schema:supersededBy

## Usage

### Embedded Context (Default)

The serializer embeds this context inline for self-contained documents:

```json
{
  "@context": { ... },
  "@type": "FocusCostRecord",
  "@id": "urn:focus:cost:abc123...",
  "billingAccountId": "123456789012"
}
```

### Remote Context

For bandwidth optimization, reference the context by URL:

```json
{
  "@context": "https://focus.finops.org/context/v1.jsonld",
  "@type": "FocusCostRecord",
  "@id": "urn:focus:cost:abc123..."
}
```

## Versioning

Context versions follow FOCUS specification versions:

- `v1` - FOCUS 1.2 and 1.3 compatibility
- Future: `v2` when FOCUS 2.0 is released

## Validation

Validate context syntax:

```bash
# Using jsonld.js playground
npx jsonld normalize contracts/focus-context.jsonld

# Or online at https://json-ld.org/playground/
```
