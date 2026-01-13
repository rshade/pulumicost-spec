# Phase 0: Research & Design Decisions

**Feature**: Add ARN to GetActualCostRequest
**Date**: 2025-12-14

## Decision 1: Field Naming

* **Decision**: Use `arn` as the field name.
* **Rationale**: The feature specification explicitly proposed `arn`. While "ARN" is AWS-specific
  terminology (Amazon Resource Name), it is often used colloquially in multi-cloud contexts to mean
  "the long, unique ID". However, to be strictly provider-agnostic, `canonical_id` might be better.
  BUT, `resource_id` is already generic. The prompt/spec specifically asked for `arn` to enable
  `finfocus-plugin-aws-ce`.
* **Alternatives Considered**:
  * `canonical_id`: More generic, but `arn` was requested.
  * `provider_id`: Could be confused with the plugin ID.
  * `cloud_identifier`: Too long.

## Decision 2: Field Number

* **Decision**: Use field number `5`.
* **Rationale**: `resource_id` (1), `start` (2), `end` (3), `tags` (4) are currently defined. `5` is
  the next available number and falls within the optimized 1-15 range for protobuf.

## Decision 3: Backward Compatibility

* **Decision**: The field will be `optional` (implicit in Proto3).
* **Rationale**: Existing clients will simply not send this field. Existing servers will ignore it
  if they don't know it, or receive it as unknown field. Updated servers can check for its presence
  (empty string check).

## Decision 4: Test Strategy

* **Decision**: Use `buf` for breaking change detection and create a Go test case in `sdk/go/testing`
  that constructs a request with the new field.
* **Rationale**: Ensures the generated code works as expected.
