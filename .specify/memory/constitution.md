<!--
Sync Impact Report - Constitution v1.0.0
========================================
Version Change: [INITIAL] → 1.0.0
Modified Principles: Initial constitution creation
Added Sections: All sections (initial creation)
Removed Sections: None
Templates Requiring Updates:
  ✅ .specify/templates/plan-template.md - Constitution Check section aligned
  ✅ .specify/templates/spec-template.md - Requirements structure aligned
  ✅ .specify/templates/tasks-template.md - Test-first workflow aligned
Follow-up TODOs: None
-->

# PulumiCost Specification Constitution

## Core Principles

### I. gRPC Proto Specification-First Development

Every change to the protocol MUST begin with protobuf specification updates before implementation:

- **Proto definitions are the source of truth**: All gRPC service definitions in
  `proto/pulumicost/v1/costsource.proto` define the contract
- Proto changes require corresponding JSON schema updates for PricingSpec messages
- All protobuf message types MUST have comprehensive validation
- Breaking changes MUST be detected via buf and documented
- SDK code is generated from proto definitions, never hand-written
- gRPC service methods MUST use proper status codes and error handling

**Rationale**: As a gRPC protobuf specification repository, the proto files define the wire protocol and
service contract. All implementation (Go SDK, other language bindings) is generated from these canonical proto
definitions.

### II. Multi-Provider gRPC Consistency

The gRPC service specification MUST maintain feature parity across all major cloud providers:

- New billing modes in PricingSpec messages MUST include cross-provider examples (AWS, Azure, GCP, Kubernetes)
- Proto message field additions MUST support existing provider patterns
- Examples MUST demonstrate real-world provider use cases as gRPC request/response payloads
- Provider-specific extensions allowed only via protobuf map fields and documented as such
- ResourceDescriptor message fields MUST be provider-agnostic

**Rationale**: Plugin developers implement the gRPC CostSourceService interface. Provider-specific proto fields
fragment the ecosystem and break interoperability.

### III. Test-First Protocol (NON-NEGOTIABLE)

TDD is mandatory for all gRPC specification changes:

1. Write conformance tests defining expected gRPC behavior (request → response)
2. Tests MUST fail against current proto/implementation
3. Update proto definitions to make tests pass
4. Regenerate SDK via buf and validate examples
5. Test gRPC error conditions and status codes
6. Red-Green-Refactor cycle strictly enforced

**Rationale**: gRPC protocol specifications have high downstream impact on plugin implementations. Tests define
the RPC contract before proto changes prevent breaking existing clients and servers.

### IV. Protobuf Backward Compatibility

Breaking changes to protobuf definitions are strictly controlled:

- MAJOR version bump required for breaking proto changes (field removals, type changes, renaming)
- buf breaking change detection MUST pass in CI
- Deprecated protobuf fields MUST remain for one MAJOR version
- Use `reserved` keyword for removed fields to prevent field number reuse
- Migration guides required for all breaking proto changes
- `UnimplementedCostSourceServiceServer` embedding required for forward compatibility
- Follow protobuf field numbering best practices (reserve 1-15 for frequent fields)

**Rationale**: Protobuf wire format compatibility is critical. Breaking changes cascade through all plugin
implementations and client applications.

### V. Comprehensive Documentation

Every gRPC specification element MUST be documented:

- Proto messages and fields require inline comments (used for generated docs)
- RPC methods require description of request/response contract
- JSON schema fields require descriptions matching proto comments
- Examples require README explanations with sample gRPC request/response payloads
- Billing modes require cross-provider coverage matrix
- API reference auto-generated from proto comments via protoc-gen-doc or buf

**Rationale**: gRPC specifications are only useful if understood. Plugin developers need clear proto field
semantics and RPC method contracts.

### VI. Performance as a gRPC Requirement

gRPC protocol design MUST consider performance implications:

- Conformance tests include RPC response time requirements
- Benchmarks track SDK generation, serialization, and gRPC call performance
- Large dataset handling (30+ days) tested at Advanced conformance level
- Concurrent RPC request requirements specified (10+ for Standard, 50+ for Advanced)
- Consider streaming RPCs for large dataset queries in future versions
- Protobuf message size considerations for network efficiency

**Rationale**: gRPC cost data queries can involve large datasets and high request volumes. Performance
requirements prevent inefficient proto message designs and RPC patterns.

### VII. Validation at Multiple Levels

Multi-layer validation ensures gRPC specification quality:

- **Protobuf layer**: Buf validates proto syntax, style, and breaking changes
- **Data layer**: JSON Schema validates PricingSpec message JSON representation
- **Service layer**: Conformance tests validate gRPC service behavior (Basic/Standard/Advanced)
- **SDK layer**: Integration tests validate generated gRPC client/server code
- **CI layer**: All validation layers run together in GitHub Actions

**Rationale**: gRPC specifications require validation at protocol definition, data serialization, service
behavior, and code generation levels. Each layer catches different proto error classes.

## Quality Standards

### gRPC Code Generation

- Generated code (sdk/go/proto/) MUST NOT be manually edited
- Generated gRPC service stubs and message types MUST be up-to-date in CI (verified via buf generate check)
- buf CLI version pinned locally (bin/buf v1.32.1) to ensure consistent proto compilation
- Proto changes automatically trigger gRPC SDK regeneration via `make generate`
- Generated code includes both protobuf message types and gRPC service interfaces

### Testing Requirements

- **Unit tests**: SDK helper code (types, validation, billing mode enums)
- **Integration tests**: In-memory gRPC server/client via bufconn harness
- **Conformance tests**: gRPC service behavior at three levels (Basic/Standard/Advanced)
- **Performance benchmarks**: gRPC call latency and message serialization with memory profiling
- **RPC error testing**: Proper gRPC status code handling for all error conditions
- All tests MUST pass before merge

### Schema Validation

- All JSON PricingSpec examples MUST validate against JSON schema
- JSON schema MUST match protobuf PricingSpec message definition
- Schema changes MUST maintain backward compatibility with existing proto messages
- AJV validation in CI with strict mode disabled for protobuf compatibility
- Cross-vendor example coverage required for all billing modes

## Development Workflow

### gRPC Specification Change Process

1. **Identify Need**: Issue describes required RPC capability or protobuf message change
2. **Research**: Analyze provider APIs and existing proto patterns
3. **Propose**: Draft proto service/message changes with example payloads
4. **Test**: Write conformance tests defining expected gRPC RPC behavior
5. **Implement**: Update proto, regenerate gRPC SDK via buf, validate examples
6. **Review**: PR includes proto changes, generated code diff, tests, examples, and documentation
7. **Validate**: CI runs all validation layers (buf lint/breaking, schema, conformance, benchmarks)

### Breaking Change Protocol (gRPC-specific)

1. Buf detects breaking protobuf change in CI (field removal, type change, etc.)
2. PR description documents proto breakage and migration path for plugin implementations
3. CHANGELOG.md updated with MAJOR version entry
4. Migration guide created in docs/ with before/after proto examples
5. Deprecation notices added to proto comments with protobuf `deprecated` option
6. Deprecated protobuf fields retained for one MAJOR version
7. Use `reserved` keyword for permanently removed field numbers

### Example Contribution Requirements

- JSON files in examples/specs/ representing protobuf PricingSpec messages
- Sample gRPC request payloads in examples/requests/
- README.md entry documenting billing model and RPC usage
- Proto validation passing (buf lint)
- JSON schema validation passing for all examples
- Cross-reference to related billing modes and RPC methods

## Governance

### Amendment Process

Constitution changes require:

1. Proposal documenting rationale and impact on gRPC development workflow
2. Review of affected templates and proto practices
3. Version bump per semantic versioning rules:
   - MAJOR: Backward incompatible governance changes affecting proto development
   - MINOR: New principles or sections added
   - PATCH: Clarifications, wording, typo fixes
4. Update to all dependent template files
5. Sync Impact Report prepended to constitution

### Compliance Review

All PRs and reviews MUST verify:

- **gRPC proto-first approach**: Proto definitions updated before SDK code
- **Cross-provider consistency**: Protobuf messages support all major providers
- **Test-first protocol**: gRPC conformance tests written and failed before proto changes
- **Backward compatibility**: buf breaking check passes (or MAJOR version justified)
- **Documentation complete**: Proto comments, examples, gRPC request/response docs, README
- **Performance requirements**: Conformance level appropriate for RPC patterns
- **Validation passing**: buf lint/breaking, JSON schema, conformance tests, benchmarks

### Complexity Justification

Any deviation from these principles MUST be justified:

- Document in PR description why simpler proto approach insufficient
- Review by maintainer required for principle violations
- Consider if new RPC method needed vs extending existing messages
- Complexity tracking in implementation plans when constitutional gates fail
- Prefer simplicity: YAGNI principles apply to gRPC service design (avoid premature streaming, complex
  message hierarchies)

### Runtime Development Guidance

For day-to-day gRPC development guidance not covered by this constitution, refer to `CLAUDE.md` in the repository
root. The constitution defines non-negotiable protobuf and gRPC principles; CLAUDE.md provides practical buf
commands, proto generation patterns, and workflow tips.

**Version**: 1.0.0 | **Ratified**: 2025-08-11 | **Last Amended**: 2025-11-17
