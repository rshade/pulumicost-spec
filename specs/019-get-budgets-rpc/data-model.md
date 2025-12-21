# Data Model: GetBudgets RPC

**Date**: 2025-12-09
**Feature**: GetBudgets RPC for Plugin-Provided Budget Information

## Entity Definitions

### Budget

Represents a spending limit with alert thresholds from cloud cost management services.

**Fields**:

- `id` (string): Unique identifier for the budget (required, non-empty)
- `name` (string): Human-readable budget name (required, non-empty)
- `source` (string): Provider identifier (e.g., "aws-budgets", "gcp-billing") (required, non-empty)
- `amount` (BudgetAmount): Monetary limit and currency (required)
- `period` (BudgetPeriod): Time interval for budget calculations (required)
- `filter` (BudgetFilter): Scope restrictions (optional)
- `thresholds` ([]BudgetThreshold): Alert points with percentages (optional, repeatable)
- `status` (BudgetStatus): Current spending state (optional, populated when include_status=true)
- `created_at` (google.protobuf.Timestamp): Budget creation time (optional)
- `updated_at` (google.protobuf.Timestamp): Last budget modification time (optional)
- `metadata` (map<string,string>): Provider-specific additional data (optional)

**Relationships**:

- Contains one BudgetAmount
- Contains one BudgetPeriod
- Contains one BudgetFilter
- Contains zero or more BudgetThresholds
- Contains zero or one BudgetStatus

**Validation Rules**:

- id must be non-empty string
- name must be non-empty string
- source must be non-empty string
- amount must be present with positive limit
- period must be valid enum value (not UNSPECIFIED)

### BudgetAmount

Specifies the monetary limit and currency for a budget.

**Fields**:

- `limit` (double): Maximum spending amount (required, > 0)
- `currency` (string): ISO 4217 currency code (required, 3 characters)

**Relationships**:

- Belongs to Budget

**Validation Rules**:

- limit must be greater than 0
- currency must be exactly 3 characters (ISO 4217 format)

### BudgetPeriod

Defines the time interval for budget calculations.

**Fields**:

- `period` (enum): Time interval type

**Enum Values**:

- `BUDGET_PERIOD_UNSPECIFIED` (0): Invalid/unset
- `BUDGET_PERIOD_DAILY` (1): Daily budget cycle
- `BUDGET_PERIOD_WEEKLY` (2): Weekly budget cycle
- `BUDGET_PERIOD_MONTHLY` (3): Monthly budget cycle
- `BUDGET_PERIOD_QUARTERLY` (4): Quarterly budget cycle
- `BUDGET_PERIOD_ANNUALLY` (5): Annual budget cycle

**Relationships**:

- Belongs to Budget

**Validation Rules**:

- Must not be UNSPECIFIED

### BudgetFilter

Allows narrowing down budgets by provider, region, resource type, or tags.

**Fields**:

- `providers` ([]string): Cloud provider restrictions (optional)
- `regions` ([]string): Geographic region restrictions (optional)
- `resource_types` ([]string): Resource type restrictions (optional)
- `tags` (map<string,string>): Tag-based filtering (optional)

**Relationships**:

- Belongs to Budget

**Validation Rules**:

- No specific validation (all fields optional)

### BudgetThreshold

Defines alert points with percentages and trigger types.

**Fields**:

- `percentage` (double): Alert threshold percentage (0-100) (required)
- `type` (ThresholdType): Actual vs forecasted spending (required)
- `triggered` (bool): Whether threshold has been crossed (optional)
- `triggered_at` (google.protobuf.Timestamp): When threshold was crossed (optional)

**Enum Values for ThresholdType**:

- `THRESHOLD_TYPE_UNSPECIFIED` (0): Invalid/unset
- `THRESHOLD_TYPE_ACTUAL` (1): Based on actual spending
- `THRESHOLD_TYPE_FORECASTED` (2): Based on forecasted spending

**Relationships**:

- Belongs to Budget (many-to-one)

**Validation Rules**:

- percentage must be between 0 and 100
- type must not be UNSPECIFIED

### BudgetStatus

Shows current spending state and utilization metrics.

**Fields**:

- `current_spend` (double): Actual spending to date (required)
- `forecasted_spend` (double): Predicted end-of-period spending (optional)
- `percentage_used` (double): Current utilization percentage (required)
- `percentage_forecasted` (double): Forecasted utilization percentage (optional)
- `currency` (string): Currency for spend amounts (required, 3 characters)
- `health` (BudgetHealthStatus): Overall budget health assessment (required)

**Enum Values for BudgetHealthStatus**:

- `BUDGET_HEALTH_UNSPECIFIED` (0): Invalid/unset
- `BUDGET_HEALTH_OK` (1): Within normal thresholds
- `BUDGET_HEALTH_WARNING` (2): Approaching limits
- `BUDGET_HEALTH_CRITICAL` (3): Near or at limits
- `BUDGET_HEALTH_EXCEEDED` (4): Over budget

**Relationships**:

- Belongs to Budget

**Validation Rules**:

- current_spend must be >= 0
- forecasted_spend must be >= current_spend if present
- percentage_used must be between 0 and 100
- percentage_forecasted must be between 0 and 100 if present
- currency must be exactly 3 characters
- health must not be UNSPECIFIED

### BudgetSummary

Provides aggregated statistics across multiple budgets.

**Fields**:

- `total_budgets` (int32): Total number of budgets (required)
- `budgets_ok` (int32): Number of healthy budgets (required)
- `budgets_warning` (int32): Number of warning budgets (required)
- `budgets_exceeded` (int32): Number of exceeded budgets (required)

**Relationships**:

- Returned in GetBudgetsResponse

**Validation Rules**:

- All counts must be >= 0
- budgets_ok + budgets_warning + budgets_exceeded <= total_budgets

## Request/Response Structures

### GetBudgetsRequest

**Fields**:

- `filter` (BudgetFilter): Optional filtering criteria
- `include_status` (bool): Whether to fetch current spending status

### GetBudgetsResponse

**Fields**:

- `budgets` ([]Budget): List of budget information
- `summary` (BudgetSummary): Aggregated statistics

## State Transitions

Budgets have implicit state based on spending patterns:

- OK → WARNING: When spending crosses warning threshold
- WARNING → CRITICAL: When spending approaches limit
- CRITICAL → EXCEEDED: When spending exceeds budget
- Any state → OK: When budget resets on new period

## Data Volume Assumptions

- Typical deployment: 100-1000 budgets per department/user
- Peak concurrent requests: 10-50 (Standard conformance level)
- Message size: Budget definitions ~1-5KB each, status data ~0.5KB each
- Response time target: <5 seconds for full budget queries
