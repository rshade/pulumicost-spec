# Research: Validation Bypass Metadata

**Feature**: 034-validation-bypass-metadata
**Date**: 2026-01-24
**Phase**: 0 - Research

## Research Summary

This feature extends existing Go SDK types with bypass metadata. No external dependencies or
complex architectural decisions required. All technical approaches are well-established patterns
within the codebase.

## Decision 1: Struct Extension vs. New Type

**Decision**: Extend existing `ValidationResult` with a `Bypasses` slice field

**Rationale**:

- Maintains backward compatibility (existing consumers ignore new field)
- Zero breaking changes to existing API
- Follows Go convention for optional slice fields (nil = empty)
- JSON serialization handles empty slices gracefully (omitempty)

**Alternatives Considered**:

- New `ValidationResultWithBypass` type: Rejected - creates parallel APIs, migration burden
- Wrapper struct: Rejected - unnecessary indirection, complicates usage
- Embedding: Rejected - changes type identity, breaks type assertions

## Decision 2: Enum Implementation Pattern

**Decision**: Use typed string constants following established `registry` package pattern

**Rationale**:

- Consistent with existing patterns in `sdk/go/registry/domain.go`
- Zero-allocation validation using package-level slices
- JSON-serializable without custom marshaling
- Self-documenting in wire format

**Alternatives Considered**:

- Integer enums with String() method: Rejected - less readable in JSON output
- Proto enums: Rejected - not a gRPC message type, overkill for SDK-internal type
- Map-based validation: Rejected - 2x slower than slice-based per benchmarks

## Decision 3: Timestamp Representation

**Decision**: Use `time.Time` with RFC3339 JSON serialization

**Rationale**:

- Standard Go time type, familiar to all Go developers
- Built-in JSON marshaling to RFC3339 format
- Timezone-aware (UTC recommended)
- Supports zero-value detection for edge cases

**Alternatives Considered**:

- Unix timestamp int64: Rejected - less human-readable, timezone ambiguity
- Custom timestamp type: Rejected - unnecessary abstraction
- String timestamp: Rejected - parsing overhead, validation complexity

## Decision 4: Reason Field Truncation

**Decision**: Truncate at 500 characters with "..." suffix in constructor/setter

**Rationale**:

- Clarification session established 500-char limit
- Early truncation prevents downstream issues
- "..." suffix indicates truncation occurred
- Validation function can warn but not reject

**Alternatives Considered**:

- Hard rejection on overflow: Rejected - too strict for audit use case
- No limit (rely on downstream): Rejected - could cause storage/display issues
- Configurable limit: Rejected - overengineering for simple use case

## Decision 5: File Organization

**Decision**: Create new `bypass.go` file for bypass-specific types

**Rationale**:

- Keeps `observability.go` focused on metrics/SLI concerns
- Clear separation of concerns
- Easier to locate bypass-related code
- Follows single-responsibility principle

**Alternatives Considered**:

- Add everything to observability.go: Rejected - file would grow too large (500+ lines)
- Separate package: Rejected - unnecessary; types are closely related to validation

## Technical Findings

### Existing Patterns to Follow

1. **Enum Validation** (`sdk/go/registry/domain.go`):

   ```go
   //nolint:gochecknoglobals // Intentional optimization
   var allBypassSeverities = []BypassSeverity{...}

   func IsValidBypassSeverity(s string) bool {
       severity := BypassSeverity(s)
       for _, valid := range allBypassSeverities {
           if severity == valid {
               return true
           }
       }
       return false
   }
   ```

2. **Struct with JSON tags** (`observability.go`):

   ```go
   type BypassMetadata struct {
       Timestamp       time.Time       `json:"timestamp"`
       Reason          string          `json:"reason"`
       // ... etc
   }
   ```

3. **Optional slice field**:

   ```go
   type ValidationResult struct {
       Valid    bool           `json:"valid"`
       Errors   []string       `json:"errors,omitempty"`
       Warnings []string       `json:"warnings,omitempty"`
       Bypasses []BypassMetadata `json:"bypasses,omitempty"` // NEW
   }
   ```

### Performance Targets

Based on existing benchmarks in `sdk/go/registry/`:

- Enum validation: <10 ns/op, 0 allocs/op
- Struct creation: <50 ns/op for simple constructors
- JSON round-trip: Acceptable overhead for audit use case

### Serialization Verification

Verified that `encoding/json` handles:

- Empty slices with `omitempty` (field omitted)
- `time.Time` as RFC3339 string
- Typed string constants as plain strings

## Open Questions

None. All technical approaches resolved through research.

## Next Steps

Proceed to Phase 1: Data Model and Contracts design.
