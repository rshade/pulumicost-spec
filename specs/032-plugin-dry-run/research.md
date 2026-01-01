# Research: Plugin Capability Dry Run Mode

**Feature**: 032-plugin-dry-run
**Date**: 2025-12-31

## Research Tasks Completed

### 1. FOCUS Field Enumeration

**Task**: Identify all FOCUS fields that need support status reporting

**Findings**:

The `FocusCostRecord` message in `proto/pulumicost/v1/focus.proto` contains ~50 fields
organized into these sections:

| Section | Field Range | Example Fields |
|---------|-------------|----------------|
| Identity & Hierarchy | 1-3, 24-25, 42-43 | provider_name, billing_account_id, sub_account_id |
| Billing Period | 18, 26-27 | billing_period_start, billing_period_end, billing_currency |
| Charge Period | 4-5 | charge_period_start, charge_period_end |
| Charge Details | 8, 28-30 | charge_category, charge_class, charge_description |
| Pricing Details | 9, 31-33, 50-54 | pricing_category, pricing_quantity, list_unit_price |
| Service & Product | 6-7, 55-56 | service_category, service_name, publisher |
| Resource Details | 12-13, 34 | resource_id, resource_name, resource_type |
| SKU Details | 14, 35, 57-58 | sku_id, sku_price_id, sku_meter |
| Location | 10-11, 36 | region_id, region_name, availability_zone |
| Financial Amounts | 15-17, 41 | billed_cost, list_cost, effective_cost |
| Consumption/Usage | 20-21 | consumed_quantity, consumed_unit |
| Commitment Discounts | 37-39, 46-49 | commitment_discount_id, commitment_discount_status |
| Capacity Reservation | 44-45 | capacity_reservation_id, capacity_reservation_status |
| Invoice Details | 19, 40 | invoice_id, invoice_issuer |
| Metadata | 22-23 | tags, extended_columns |
| FOCUS 1.3 Additions | 59-66 | service_provider_name, allocated_method_id, contract_applied |

**Decision**: Field enumeration should be based on `FocusCostRecord` message definition.
The DryRun response should report on fields 1-66 (all currently defined FOCUS fields).

### 2. Existing RPC Pattern Analysis

**Task**: Analyze existing CostSourceService RPCs for design patterns

**Findings**:

Current CostSourceService has 12 RPCs:

| RPC | Request Type | Response Type | Pattern |
|-----|--------------|---------------|---------|
| Name | NameRequest (empty) | NameResponse | Simple metadata |
| Supports | SupportsRequest (ResourceDescriptor) | SupportsResponse (bool + capabilities map) | Capability check |
| GetActualCost | GetActualCostRequest | GetActualCostResponse | Data retrieval |
| GetProjectedCost | GetProjectedCostRequest | GetProjectedCostResponse | Data retrieval |
| GetPricingSpec | GetPricingSpecRequest | GetPricingSpecResponse | Specification query |
| EstimateCost | EstimateCostRequest | EstimateCostResponse | Calculation |
| GetRecommendations | GetRecommendationsRequest | GetRecommendationsResponse | Data retrieval |
| DismissRecommendation | DismissRecommendationRequest | DismissRecommendationResponse | State mutation |
| GetBudgets | GetBudgetsRequest | GetBudgetsResponse | Data retrieval |
| GetPluginInfo | GetPluginInfoRequest (empty) | GetPluginInfoResponse | Metadata |

**Key Patterns Observed**:

1. **Capability discovery**: `Supports` RPC uses `capabilities` map (`map<string, bool>`)
2. **Metadata queries**: `GetPluginInfo` returns plugin-level info without data retrieval
3. **Request structure**: Most requests include `ResourceDescriptor` for resource targeting
4. **Error handling**: Uses gRPC status codes, documented in proto comments
5. **Optional interfaces**: Some RPCs are optional (return Unimplemented if not supported)

**Decision**: DryRun RPC follows the `GetPricingSpec` pattern (takes ResourceDescriptor,
returns specification data). Use `Supports.capabilities["dry_run"]` for capability check.

### 3. Proto Best Practices for New Enums

**Task**: Determine best practices for FieldSupportStatus enum

**Findings**:

From existing enums in `enums.proto`:

1. All enums have `_UNSPECIFIED = 0` as the default/unknown value
2. Enum names use SCREAMING_SNAKE_CASE with prefix matching enum name
3. Each enum value has clear documentation comment
4. Related enums are grouped together (FOCUS enums in enums.proto)

**Decision**: Add `FieldSupportStatus` enum to enums.proto:

```protobuf
enum FieldSupportStatus {
  FIELD_SUPPORT_STATUS_UNSPECIFIED = 0;
  FIELD_SUPPORT_STATUS_SUPPORTED = 1;
  FIELD_SUPPORT_STATUS_UNSUPPORTED = 2;
  FIELD_SUPPORT_STATUS_CONDITIONAL = 3;
  FIELD_SUPPORT_STATUS_DYNAMIC = 4;
}
```

### 4. Response Message Design

**Task**: Design DryRunRequest/Response messages

**Findings from similar messages**:

- `SupportsResponse`: Uses `map<string, bool> capabilities` for extensible flags
- `GetPluginInfoResponse`: Uses `map<string, string> metadata` for key-value pairs
- `RecommendationFilter`: Uses individual fields for structured queries

**Decision**: Structured message with repeated FieldMapping:

```protobuf
message DryRunRequest {
  ResourceDescriptor resource = 1;
  // Optional: simulate specific scenarios (e.g., different regions)
  map<string, string> simulation_parameters = 2;
}

message DryRunResponse {
  repeated FieldMapping field_mappings = 1;
  bool configuration_valid = 2;
  repeated string configuration_errors = 3;
  bool resource_type_supported = 4;
}

message FieldMapping {
  string field_name = 1;
  FieldSupportStatus support_status = 2;
  string condition_description = 3;
  string expected_type = 4;
}
```

### 5. Backward Compatibility Strategy

**Task**: Ensure backward compatibility with existing plugins

**Findings**:

From `GetPluginInfo` documentation:

```go
// Client-side error handling example (Go):
//   resp, err := client.GetPluginInfo(ctx, &GetPluginInfoRequest{})
//   if err != nil {
//       if status.Code(err) == codes.Unimplemented {
//           // Legacy plugin - use fallback values
```

**Decision**: DryRun RPC follows same pattern:

1. Plugins that don't implement DryRun return `codes.Unimplemented`
2. Hosts detect this and fall back to inferring capabilities from Supports RPC
3. `SupportsResponse.capabilities["dry_run"]` indicates explicit support
4. New optional `dry_run` field on GetActualCost/GetProjectedCost requests defaults to
   false, preserving existing behavior

## Alternatives Considered

### Alternative 1: Extend Supports RPC Instead of New DryRun RPC

**Rejected because**: Supports RPC answers "do you support this resource?" while DryRun
answers "what fields would you populate?". Conflating these creates semantic confusion
and makes the response message unwieldy.

### Alternative 2: Field Status via Map Instead of Repeated Message

**Rejected because**: `map<string, FieldSupportStatus>` loses the ability to include
`condition_description` and `expected_type` per field. The repeated message pattern is
more extensible.

### Alternative 3: Only Add dry_run Flag to Existing RPCs

**Rejected because**: This requires modifying response types (GetActualCostResponse would
need to conditionally include field mappings). A dedicated DryRun RPC has cleaner
semantics and doesn't pollute existing response types.

## Research Conclusions

All NEEDS CLARIFICATION items resolved:

1. **Field enumeration**: Use FocusCostRecord fields (1-66)
2. **RPC design**: New DryRun RPC + optional dry_run flag on cost RPCs
3. **Enum pattern**: FieldSupportStatus with UNSPECIFIED/SUPPORTED/UNSUPPORTED/CONDITIONAL/DYNAMIC
4. **Backward compatibility**: Unimplemented error handling + capabilities map

Ready for Phase 1: Design & Contracts.
