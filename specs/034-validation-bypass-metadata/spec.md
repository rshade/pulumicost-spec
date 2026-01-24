# Feature Specification: Validation Bypass Metadata

**Feature Branch**: `034-validation-bypass-metadata`
**Created**: 2026-01-24
**Status**: Draft
**Input**: User description: "Add fields to ValidationResult to carry metadata about why a policy
was bypassed (if --yolo was used), ensuring audit trails survive the stateless boundary."

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Audit Trail for Bypassed Validations (Priority: P1)

As a compliance officer, I need to see a clear audit trail when cost validation policies were
bypassed, so that I can understand who bypassed validations, when, and why, for compliance
reporting and incident investigation purposes.

**Why this priority**: This is the core value proposition of the feature. Without audit trail
metadata surviving the stateless boundary, organizations cannot maintain compliance records or
investigate why certain validations were skipped during cost processing.

**Independent Test**: Can be fully tested by triggering a validation bypass (e.g., `--yolo` flag)
and verifying the ValidationResult contains complete bypass metadata including timestamp, reason,
and operator identity.

**Acceptance Scenarios**:

1. **Given** a validation check fails due to policy violation, **When** an operator bypasses the
   validation using a force flag, **Then** the ValidationResult includes bypass metadata with
   reason, timestamp, and operator identifier.
2. **Given** a ValidationResult with bypass metadata, **When** the result crosses a stateless
   service boundary (e.g., gRPC response), **Then** the bypass metadata is preserved and
   accessible on the receiving side.
3. **Given** multiple validations are bypassed in a single operation, **When** the ValidationResult
   is generated, **Then** each bypassed validation has its own metadata entry with specific details.

---

### User Story 2 - Display Bypass Information in CLI Output (Priority: P2)

As an operator using the CLI, I want to see clear warnings when I bypass validations, so that I
understand the implications of my actions and have a record of what was bypassed.

**Why this priority**: User feedback is essential for safe operation. Operators need immediate
visibility into which validations they bypassed and why this might be risky.

**Independent Test**: Can be tested by running a command with validation bypass enabled and
verifying the CLI output displays all bypassed validations with their severity and original
failure reasons.

**Acceptance Scenarios**:

1. **Given** an operator runs a command with `--yolo` flag, **When** validations are bypassed,
   **Then** the CLI displays a summary of bypassed validations with their severity levels.
2. **Given** bypass metadata exists in a ValidationResult, **When** formatting the result for
   display, **Then** the output clearly distinguishes between passed, failed, and bypassed
   validations.

---

### User Story 3 - Query Historical Bypass Events (Priority: P3)

As a security analyst, I want to query which validations were bypassed over a time period, so
that I can identify patterns of policy violations and assess risk exposure.

**Why this priority**: Historical analysis enables proactive security management and helps
identify operators who frequently bypass validations or specific policies that are commonly
overridden.

**Independent Test**: Can be tested by generating multiple ValidationResults with bypass metadata
over simulated time periods and verifying the data supports filtering by time range, operator,
and bypass reason.

**Acceptance Scenarios**:

1. **Given** ValidationResults with bypass metadata are stored, **When** querying by time range,
   **Then** all bypass events within that range are returned with full metadata.
2. **Given** multiple operators have bypassed validations, **When** filtering by operator
   identifier, **Then** only that operator's bypass events are returned.

---

### Edge Cases

- What happens when a bypass is requested but no validations actually fail?
  (No bypass metadata should be recorded)
- How does the system handle bypass metadata when the operator identifier is not available?
  (Use "unknown" with a warning)
- What happens when bypass reason exceeds 500 characters?
  (Truncate at 500 characters with "..." suffix indicating truncation)
- How does the system behave when timestamp cannot be determined?
  (Use zero-value timestamp with warning flag)

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: ValidationResult MUST include optional fields for recording bypass events
- **FR-002**: Bypass metadata MUST include a timestamp indicating when the bypass occurred
- **FR-003**: Bypass metadata MUST include a human-readable reason explaining why the bypass
  was performed
- **FR-004**: Bypass metadata SHOULD include an operator identifier when available (user,
  service account, or system)
- **FR-005**: Bypass metadata MUST include the name/identifier of the validation that was
  bypassed
- **FR-006**: Bypass metadata MUST include the original validation failure message that would
  have been shown
- **FR-007**: ValidationResult MUST preserve bypass metadata across gRPC serialization and
  deserialization
- **FR-008**: System MUST record bypass severity level (warning, error, critical) for each
  bypassed validation
- **FR-009**: Bypass metadata MUST NOT be recorded when no validations were actually bypassed
- **FR-010**: System MUST support multiple bypass entries in a single ValidationResult (one
  per bypassed validation)
- **FR-011**: System MUST include a mechanism type field indicating how the bypass was triggered
  (flag, environment variable, configuration)
- **FR-012**: Bypass metadata MUST be retained for a minimum of 90 days to support quarterly
  compliance reviews and historical analysis (caller responsibility, not SDK)
- **FR-013**: System MUST include a truncation indicator when the bypass reason exceeds the
  maximum length and is truncated

### Key Entities

- **ValidationResult**: Extended structure containing validation outcomes plus optional bypass
  metadata. Represents the complete result of a validation operation including any policy
  overrides.
- **BypassMetadata**: Individual bypass event record containing timestamp, reason (max 500
  characters), operator, validation name, original failure message, severity, and mechanism type.
- **BypassMechanism**: Enumeration of how the bypass was triggered (command-line flag,
  environment variable, configuration file, programmatic override).
- **BypassSeverity**: Enumeration indicating the risk level of the bypassed validation
  (warning, error, critical).

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: 100% of validation bypass events are captured with complete metadata when bypass
  mechanisms are used
- **SC-002**: Bypass metadata survives round-trip serialization with zero data loss across
  service boundaries
- **SC-003**: Operators can identify the specific validations bypassed within 5 seconds of
  viewing output
- **SC-004**: Compliance reports can enumerate all bypass events for a given time period with
  accurate timestamps
- **SC-005**: System correctly distinguishes between "no validations run", "validations passed",
  and "validations bypassed" states

## Clarifications

### Session 2026-01-24

- Q: What is the maximum length for the bypass reason field? → A: 500 characters (room for
  context while bounded)
- Q: How long should bypass metadata be retained? → A: 90 days (quarterly compliance window)

## Assumptions

- The existing ValidationResult structure in `sdk/go/pricing/observability.go` will be extended
  rather than replaced
- Bypass metadata follows the same serialization patterns as existing ValidationResult fields
- Operator identity will be provided by the calling context (CLI extracts from environment, API
  from auth context)
- Timestamp precision of seconds is sufficient for audit purposes
- Bypass events are expected to be infrequent in normal operations (less than 5% of validation
  runs)
- The `--yolo` flag mentioned in the issue is one example of a bypass mechanism; the solution
  should support multiple bypass mechanisms
