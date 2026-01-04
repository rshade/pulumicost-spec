# Research: JSON-LD / Schema.org Serialization

**Feature**: 032-jsonld-serialization
**Date**: 2025-12-31

## Research Topics

### 1. JSON-LD 1.1 Serialization Best Practices

**Decision**: Use Go's stdlib `encoding/json` with custom struct types for JSON-LD output

**Rationale**:

- JSON-LD is valid JSON - no special library needed for serialization-only use case
- Go's `encoding/json` is highly optimized and well-understood
- Custom struct types with `json` tags provide clean mapping from proto to JSON-LD
- Avoids external dependency complexity for a presentation-layer feature
- Full JSON-LD processing (expansion, compaction, framing) would need a library, but
  we only need serialization

**Alternatives Considered**:

- `github.com/piprate/json-gold`: Full JSON-LD processor, overkill for serialization-only
- `github.com/kazarena/json-gold`: Older fork, less maintained
- Custom RDF library: Wrong tool - we're outputting JSON-LD, not processing RDF

### 2. JSON-LD Context Structure

**Decision**: Embedded inline context with remote context URL option

**Rationale**:

- Inline context ensures self-contained documents (no network dependency)
- Remote context URL support enables enterprise customization
- JSON-LD 1.1 supports both patterns via `@context` array syntax
- Example: `"@context": ["https://schema.org", {"focus": "https://focus.finops.org/v1#"}]`

**Alternatives Considered**:

- Remote-only context: Requires network access, breaks offline use
- No context (expanded form): Verbose, harder to read, loses semantic web benefits

### 3. Schema.org Vocabulary Mappings

**Decision**: Map FOCUS fields to Schema.org where natural mappings exist

| FOCUS Field | Schema.org Type | Property |
|-------------|-----------------|----------|
| billed_cost, list_cost, effective_cost | schema:MonetaryAmount | schema:value + schema:currency |
| billing_currency | schema:MonetaryAmount | schema:currency (ISO 4217) |
| charge_period_start, charge_period_end | schema:DateTime | ISO 8601 format |
| service_name | schema:Service | schema:name |
| resource_name | schema:Thing | schema:name |
| region_name | schema:Place | schema:name |
| tags | schema:PropertyValue | schema:name + schema:value |

**Rationale**:

- Schema.org is widely recognized by search engines and knowledge graph tools
- MonetaryAmount provides proper currency semantics
- DateTime ensures ISO 8601 compliance
- Remaining FOCUS-specific fields use custom namespace

**Alternatives Considered**:

- GoodRelations ontology: More e-commerce focused, less universal
- DBpedia ontology: Academic focus, less tool support
- Pure custom ontology: Loses interoperability benefits

### 4. FOCUS Vocabulary Namespace

**Decision**: Use `https://focus.finops.org/v1#` as FOCUS vocabulary namespace

**Rationale**:

- Follows FinOps Foundation FOCUS specification domain
- Version-qualified (`v1#`) enables future evolution
- Fragment identifier (`#`) enables term definition without URL resolution
- Example terms: `focus:billingAccountId`, `focus:chargeCategory`, `focus:serviceCategory`

**Alternatives Considered**:

- `urn:focus:v1:`: URN scheme, less web-friendly
- `https://pulumicost.dev/focus/`: Project-specific, not spec-aligned
- No namespace: Breaks RDF semantics

### 5. ID Generation Strategy

**Decision**: User-provided ID with fallback to SHA256 composite key hash

**Rationale** (confirmed from clarification):

- User-provided IDs enable enterprise ID scheme integration
- Fallback ensures all records have deterministic IDs
- SHA256 of `billing_account_id + charge_period_start + resource_id` creates unique key
- Format: `urn:focus:cost:{hash}` for cost records
- Format: `urn:focus:commitment:{hash}` for contract commitments

**Alternatives Considered**:

- UUID v4: Random, prevents deduplication
- UUID v5: Namespace-based, but less transparent than direct hash
- Content-hash of entire record: Too sensitive to field ordering changes

### 6. Empty Value Handling

**Decision**: Omit fields with empty string, zero numeric, or nil values

**Rationale** (confirmed from clarification):

- JSON-LD open world assumption: missing = unknown, not empty
- Reduces output size significantly (typical records have many empty optional fields)
- Cleaner graph database ingestion (no noise from empty values)
- Aligns with proto3 default value semantics

**Alternatives Considered**:

- Include all with null: Creates noise, bloats output
- Include empty strings: Semantically different from missing
- Configurable: Adds API complexity for minimal benefit

### 7. Streaming Batch Serialization

**Decision**: Use `io.Writer` interface with JSON array streaming

**Rationale**:

- `io.Writer` enables any output target (file, network, buffer)
- JSON array streaming: `[` then record-by-record with `,` separator, then `]`
- Memory usage bounded to single record at a time
- Compatible with JSON-LD arrays (`@graph` pattern)

**Alternatives Considered**:

- Channel-based: More complex, harder to test
- Callback-based: Less idiomatic Go
- Full materialization: Memory explosion for large datasets

### 8. Deprecated Field Handling

**Decision**: Serialize deprecated fields with `@deprecated` annotation in context

**Rationale**:

- FOCUS proto marks `provider_name` and `publisher` as deprecated
- JSON-LD supports custom annotations via context
- Example: `"provider_name": {"@id": "focus:providerName", "@deprecated": true}`
- Enables consumers to detect and handle deprecations

**Alternatives Considered**:

- Omit deprecated fields: Breaks backward compatibility
- No annotation: Consumers unaware of deprecation
- Separate deprecation object: Over-engineering

### 9. Enum Serialization

**Decision**: Serialize enums as human-readable strings with IRI option

**Rationale**:

- Human-readable strings improve debuggability: `"CHARGE_CATEGORY_USAGE"` vs `1`
- IRI option for full semantic web compliance: `focus:ChargeCategoryUsage`
- Default to strings, configurable to IRIs
- Matches proto enum name conventions

**Alternatives Considered**:

- Numeric values: Loses semantics in JSON-LD
- Full IRI only: Verbose for typical use cases
- CamelCase only: Diverges from proto naming

### 10. Performance Optimization Patterns

**Decision**: Zero-allocation patterns where possible, sync.Pool for buffers

**Rationale**:

- Go stdlib JSON encoder is fast but allocates
- Pre-allocated struct fields avoid per-record allocation
- sync.Pool for byte buffers enables reuse across batch items
- Target: <1ms per record at p99, 0 allocs for hot path

**Alternatives Considered**:

- Third-party fast JSON (jsoniter): External dependency
- Code generation: Maintenance burden, diminishing returns
- No optimization: May miss <1ms target

## Summary

All technical decisions resolved. Key patterns:

1. **Stdlib-only implementation** with custom structs for JSON-LD mapping
2. **Embedded inline context** with Schema.org + FOCUS vocabulary
3. **Deterministic IDs** via user-provided or SHA256 composite key
4. **Omit empty values** following JSON-LD open world assumption
5. **Streaming batch** via io.Writer for bounded memory
6. **Human-readable enums** with IRI option for semantic web use cases
