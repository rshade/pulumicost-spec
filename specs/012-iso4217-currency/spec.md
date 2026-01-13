# Feature Specification: ISO 4217 Currency Validation Package

**Feature Branch**: `013-iso4217-currency`
**Created**: 2025-11-30
**Status**: Draft
**Input**: Extract ISO 4217 currency validation logic from `sdk/go/pluginsdk/focus_conformance.go`
into a separate, reusable package.

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Validate Currency Codes in FOCUS Records (Priority: P1)

As a plugin developer validating FOCUS cost records, I need to validate that currency codes conform
to ISO 4217 so that billing data maintains international financial standards compliance.

**Why this priority**: Currency validation is the core functionality that enables FOCUS compliance.
The existing implementation in `focus_conformance.go` already depends on this capability, making it
the essential baseline feature.

**Independent Test**: Can be fully tested by providing various currency codes (valid and invalid)
and verifying the validation response returns the expected boolean result.

**Acceptance Scenarios**:

1. **Given** a valid ISO 4217 currency code like "USD", **When** validating the code,
   **Then** the validation returns success (true).
2. **Given** an invalid currency code like "XYZ" or "US", **When** validating the code,
   **Then** the validation returns failure (false).
3. **Given** a lowercase currency code like "usd", **When** validating the code,
   **Then** the validation returns failure (ISO 4217 codes are uppercase).
4. **Given** an empty string, **When** validating the code,
   **Then** the validation returns failure.

---

### User Story 2 - Retrieve Currency Metadata (Priority: P2)

As a developer building financial reporting features, I need to retrieve currency metadata (name,
numeric code, minor units) so that I can format monetary values correctly for different currencies.

**Why this priority**: While not strictly required for basic validation, currency metadata enables
proper formatting (e.g., knowing JPY uses 0 decimal places while USD uses 2) and enhances user
experience in cost reporting.

**Independent Test**: Can be fully tested by retrieving metadata for known currencies and verifying
the returned information matches ISO 4217 specifications.

**Acceptance Scenarios**:

1. **Given** a valid currency code "USD", **When** retrieving currency metadata,
   **Then** the system returns name "US Dollar", numeric code "840", and minor units "2".
2. **Given** a valid currency code "JPY", **When** retrieving currency metadata,
   **Then** the system returns name "Yen", numeric code "392", and minor units "0".
3. **Given** an invalid currency code, **When** retrieving currency metadata,
   **Then** the system returns an appropriate error.

---

### User Story 3 - List All Valid Currencies (Priority: P3)

As a developer building currency selection interfaces or generating documentation, I need to
retrieve a complete list of all valid ISO 4217 currencies so that I can present options to users
or validate against the complete set.

**Why this priority**: Listing all currencies supports UI development and documentation generation
but is not required for the core validation use case.

**Independent Test**: Can be fully tested by retrieving the currency list and verifying it contains
all expected ISO 4217 currencies (180+ codes).

**Acceptance Scenarios**:

1. **Given** a request for all currencies, **When** retrieving the list,
   **Then** the system returns all active ISO 4217 currencies (180+ entries).
2. **Given** the returned currency list, **When** inspecting any entry,
   **Then** each entry includes code, name, numeric code, and minor units.

---

### User Story 4 - Migrate Existing Validation (Priority: P4)

As a maintainer of the FinFocus SDK, I need the existing `focus_conformance.go` to use the new
currency package so that validation logic is centralized and maintainable.

**Why this priority**: This is a refactoring task that depends on the new package being complete.
It ensures the ecosystem benefits from the reusable package.

**Independent Test**: Can be tested by running existing FOCUS conformance tests and verifying they
pass with the migrated implementation.

**Acceptance Scenarios**:

1. **Given** the existing FOCUS conformance tests, **When** running tests after migration,
   **Then** all tests pass without modification.
2. **Given** the `focus_conformance.go` file, **When** inspecting after migration,
   **Then** it imports and uses the new currency package instead of inline validation.

---

### Edge Cases

- What happens when currency code contains leading/trailing whitespace?
  (Should return invalid - codes must be exactly 3 uppercase letters)
- How does system handle historic/obsolete currency codes (e.g., "DEM" for Deutsche Mark)?
  (Should return invalid - only active ISO 4217 codes are included)
- How does system handle supranational currencies (e.g., "XDR" for Special Drawing Rights)?
  (Should return valid - these are in ISO 4217)
- How does system handle test currencies (e.g., "XTS")?
  (Should return valid - test codes are in ISO 4217)
- How does system handle no-currency code "XXX"?
  (Should return valid - this is a valid ISO 4217 code for "no currency")

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: Package MUST provide a function to validate if a string is a valid ISO 4217
  currency code.
- **FR-002**: Package MUST include all active ISO 4217 currency codes (180+ currencies as of
  current ISO 4217 maintenance).
- **FR-003**: Package MUST provide a function to retrieve currency metadata (code, name,
  numeric code, minor units) for a valid currency code.
- **FR-004**: Package MUST provide a function to list all valid currency codes and their metadata.
- **FR-005**: Validation MUST match exact uppercase 3-letter codes only (case-sensitive).
- **FR-006**: Package MUST use zero-allocation validation pattern consistent with
  `sdk/go/registry/domain.go` for performance.
- **FR-007**: Package MUST provide comprehensive unit tests achieving >90% coverage.
- **FR-008**: Package MUST provide benchmarks consistent with registry package standards
  (measuring ns/op and allocs/op).
- **FR-009**: The existing `focus_conformance.go` MUST be updated to import and use the new
  currency package.

### Key Entities

- **Currency**: Represents an ISO 4217 currency with code (3-letter string), name (human-readable
  string), numeric code (3-digit string), and minor units (integer representing decimal places
  for the currency).

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Currency validation completes in under 15 nanoseconds per operation (matching
  registry package performance of 5-12 ns/op).
- **SC-002**: Currency validation performs zero memory allocations per operation
  (0 B/op, 0 allocs/op).
- **SC-003**: Package includes complete ISO 4217 currency list (180+ active currencies verified
  against official specification).
- **SC-004**: All existing FOCUS conformance tests pass after migration to new package.
- **SC-005**: Unit test coverage exceeds 90% for the new currency package.

## Assumptions

- **Currency List Scope**: The package includes only active ISO 4217 currencies. Historic/withdrawn
  currencies (e.g., DEM, FRF) are excluded unless specifically required for backward compatibility
  with existing billing data.
- **Update Frequency**: Currency codes change infrequently (ISO 4217 is updated periodically).
  The initial implementation uses a static list; future maintenance will update as needed.
- **Numeric Codes**: Numeric codes are stored as strings to preserve leading zeros
  (e.g., "008" for Albanian Lek).
- **Case Sensitivity**: Validation is case-sensitive following ISO 4217 standard
  (codes are always uppercase).
