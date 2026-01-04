# Feature Specification: JSON-LD / Schema.org Serialization

**Feature Branch**: `032-jsonld-serialization`
**Created**: 2025-12-31
**Status**: Draft
**Input**: User description: "discovery: JSON-LD / Schema.org Serialization - Support
high-performance serialization of FOCUS data into JSON-LD formats for enterprise knowledge
graph indexing."

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Basic FOCUS Record Serialization (Priority: P1)

A FinOps practitioner wants to export individual FOCUS cost records as JSON-LD documents
so they can be ingested into their enterprise knowledge graph system for cross-referencing
with other organizational data (projects, teams, applications).

**Why this priority**: This is the foundational capability - serializing a single FOCUS
cost record to JSON-LD. Without this, no other JSON-LD functionality is possible. It
delivers immediate value for organizations starting knowledge graph adoption.

**Independent Test**: Can be fully tested by serializing a single FocusCostRecord and
validating the output is well-formed JSON-LD that passes JSON-LD validation tools.
Delivers value by enabling knowledge graph ingestion.

**Acceptance Scenarios**:

1. **Given** a populated FocusCostRecord with required fields (billing_account_id,
   charge_period_start, service_name, billed_cost), **When** serialized to JSON-LD,
   **Then** the output includes a valid `@context` declaration and correctly mapped
   properties
2. **Given** a FocusCostRecord with FOCUS 1.3 allocation fields populated, **When**
   serialized to JSON-LD, **Then** the allocation relationships are expressed as linked
   data references
3. **Given** a FocusCostRecord with tags and extended_columns, **When** serialized to
   JSON-LD, **Then** custom properties are included using an appropriate extension
   namespace

---

### User Story 2 - Batch Serialization for Enterprise Scale (Priority: P2)

A data engineer needs to export large volumes of FOCUS cost records (thousands to
millions) efficiently so they can be bulk-loaded into their knowledge graph database
without performance bottlenecks.

**Why this priority**: High-performance batch serialization is critical for enterprise
adoption. Single-record serialization is useful but insufficient for real-world data
volumes. This enables practical use cases.

**Independent Test**: Can be tested by serializing 10,000 FocusCostRecords in a batch
and measuring throughput. Delivers value by enabling enterprise-scale data pipelines.

**Acceptance Scenarios**:

1. **Given** a collection of 10,000 FocusCostRecords, **When** batch serialized to
   JSON-LD, **Then** processing completes within acceptable time bounds for the data
   volume
2. **Given** a batch serialization operation, **When** processing large datasets,
   **Then** memory usage remains bounded (streaming output rather than full
   materialization)
3. **Given** a batch serialization with mixed record completeness, **When** some records
   have validation errors, **Then** valid records are still output and errors are
   reported without blocking the entire batch

---

### User Story 3 - Schema.org Vocabulary Mapping (Priority: P2)

A data scientist wants FOCUS cost data to use Schema.org vocabulary where applicable so
it can be discovered and understood by standard semantic web tools and search engines
that already understand Schema.org.

**Why this priority**: Schema.org compatibility enables interoperability with existing
semantic web infrastructure. This expands the ecosystem of tools that can consume the
data.

**Independent Test**: Can be tested by serializing a FocusCostRecord and validating that
Schema.org types and properties are correctly applied where mappings exist. Delivers
value by enabling semantic web tool integration.

**Acceptance Scenarios**:

1. **Given** a FocusCostRecord with monetary values, **When** serialized to JSON-LD,
   **Then** cost fields use Schema.org MonetaryAmount or PriceSpecification types where
   appropriate
2. **Given** a FocusCostRecord with timestamp fields, **When** serialized to JSON-LD,
   **Then** dates use Schema.org DateTime formatting (ISO 8601)
3. **Given** FOCUS-specific fields with no Schema.org equivalent, **When** serialized to
   JSON-LD, **Then** a custom FOCUS vocabulary namespace is used for those properties

---

### User Story 4 - Custom Context Configuration (Priority: P3)

A platform architect wants to customize the JSON-LD context to align with their
organization's existing ontology so the FOCUS data integrates seamlessly with their
domain-specific knowledge graph schema.

**Why this priority**: Organizations have existing ontologies and need to map FOCUS data
to their schema. This enables adoption in enterprises with mature knowledge graph
infrastructure.

**Independent Test**: Can be tested by providing a custom context configuration and
validating the output uses the specified property mappings. Delivers value by enabling
integration with existing enterprise ontologies.

**Acceptance Scenarios**:

1. **Given** a custom context configuration that remaps property names, **When**
   serializing a FocusCostRecord, **Then** the output uses the custom property names
2. **Given** a context configuration that adds organization-specific types, **When**
   serializing, **Then** records include the additional type declarations
3. **Given** an invalid context configuration, **When** attempting to serialize, **Then**
   a \*jsonld.ValidationError is returned containing:
   - Field name causing validation failure
   - Specific issue (malformed URL, non-IRI format, circular reference)
   - Suggested fix or correction

---

### User Story 5 - Contract Commitment Dataset Serialization (Priority: P3)

A finance analyst wants to export ContractCommitment records as JSON-LD with proper links
to associated FocusCostRecords so they can analyze commitment utilization through graph
queries.

**Why this priority**: ContractCommitment is a separate FOCUS 1.3 dataset that must be
serializable independently. Linking to cost records enables powerful graph-based
analysis.

**Independent Test**: Can be tested by serializing a ContractCommitment and validating
it has proper linked data references. Delivers value by enabling commitment analysis in
knowledge graphs.

**Acceptance Scenarios**:

1. **Given** a ContractCommitment record, **When** serialized to JSON-LD, **Then** the
   output is a valid JSON-LD document with appropriate type declarations
2. **Given** a ContractCommitment and associated FocusCostRecords with contract_applied
   references, **When** both are serialized, **Then** the JSON-LD documents have
   resolvable references between them

---

### Edge Cases

- What happens when a FocusCostRecord has only required fields and all optional fields
  are empty/zero?
- How does the system handle deprecated fields (provider_name, publisher) - are they
  serialized with deprecation annotations?
- What happens when tags or extended_columns contain characters that need JSON escaping?
- Empty/null values are omitted from output (not serialized as empty strings or nulls)
- What happens when serializing records with invalid UTF-8 sequences in string fields?

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: System MUST serialize FocusCostRecord protocol buffer messages to valid
  JSON-LD 1.1 format
- **FR-002**: System MUST include a `@context` declaration in all JSON-LD output that
  defines FOCUS vocabulary terms
- **FR-003**: System MUST map FOCUS monetary fields (billed_cost, list_cost,
  effective_cost, contracted_cost) to appropriate linked data representations with
  currency information
- **FR-004**: System MUST serialize timestamp fields (google.protobuf.Timestamp) to
  ISO 8601 formatted strings
- **FR-005**: System MUST handle enum fields (FocusChargeCategory, FocusServiceCategory,
  etc.) by serializing them as IRIs or controlled vocabulary terms
- **FR-006**: System MUST serialize map fields (tags, extended_columns, allocated_tags)
  as JSON objects preserving all key-value pairs
- **FR-007**: System MUST omit fields with empty string, zero numeric, or nil values
  from JSON-LD output (following JSON-LD open world assumption)
- **FR-008**: System MUST provide a batch serialization interface for processing
  multiple records efficiently
- **FR-009**: System MUST serialize ContractCommitment protocol buffer messages to valid
  JSON-LD 1.1 format
- **FR-010**: System MUST generate unique `@id` values for each serialized record using
  user-provided ID if available, otherwise falling back to composite key hash (SHA256 of
  billing_account_id + charge_period_start + resource_id)
- **FR-011**: System MUST support Schema.org vocabulary mappings for common properties
  (dates, monetary amounts, organization references)
- **FR-012**: System MUST define a FOCUS-specific vocabulary namespace for FOCUS-only
  concepts that lack Schema.org equivalents
- **FR-013**: System MUST allow configuration of custom JSON-LD contexts to support
  enterprise ontology integration
- **FR-014**: System MUST include deprecated field annotations in JSON-LD output for
  fields marked as deprecated in the proto schema
- **FR-015**: System MUST support streaming output for large batch operations to bound
  memory usage

### Key Entities

- **FocusCostRecord**: The primary cost record entity with 66 fields across identity,
  billing, pricing, service, resource, financial, and metadata categories
- **ContractCommitment**: The FOCUS 1.3 contract commitment entity with commitment
  terms, periods, and financial obligations
- **JSON-LD Context**: The vocabulary mapping definition that controls how proto fields
  map to linked data properties
- **FOCUS Vocabulary**: The custom RDF vocabulary namespace defining FOCUS-specific
  terms not in Schema.org

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Single record serialization completes in under 1 millisecond for a
  fully-populated FocusCostRecord
- **SC-002**: Batch serialization of 10,000 records completes in under 5 seconds
- **SC-003**: Memory usage during batch serialization stays bounded regardless of total
  record count (streaming design)
- **SC-004**: 100% of generated JSON-LD output passes validation against JSON-LD 1.1
  specification
- **SC-005**: Schema.org type annotations are correctly applied to 100% of mappable
  fields
- **SC-006**: Custom context configurations work correctly for 100% of supported
  override scenarios
- **SC-007**: Zero data loss - all populated fields from source records appear in
  JSON-LD output

## Assumptions

- JSON-LD 1.1 specification is the target (not JSON-LD 1.0) for better support of
  nested contexts and @included
- Schema.org vocabulary will be used where natural mappings exist; custom FOCUS
  vocabulary fills gaps
- Record IDs (`@id` values) support user-provided IDs with fallback to deterministic
  composite key hash (not random UUIDs)
- The Go SDK is the primary implementation target; other language SDKs may follow the
  same patterns
- Consumers are expected to have JSON-LD processing capabilities (JSON-LD processors,
  RDF libraries, or graph databases)
- Performance targets assume typical server-class hardware (multi-core CPU, adequate
  RAM)

## Clarifications

### Session 2025-12-31

- Q: What strategy should be used for generating `@id` values? → A: User-provided ID
  with fallback to composite key hash (SHA256 of billing_account_id +
  charge_period_start + resource_id)
- Q: How should empty/null values be handled in JSON-LD output? → A: Omit fields with
  empty string, zero numeric, or nil values

## Non-Goals / Out of Scope

- RDF/Turtle or other RDF serialization formats (JSON-LD only for this feature)
- Full OWL ontology definition (vocabulary/context only)
- SPARQL query interface implementation
- Knowledge graph database integration (serialization only, not loading)
- UI for viewing/editing JSON-LD output
- JSON-LD framing or compaction algorithms (output is expanded form)
