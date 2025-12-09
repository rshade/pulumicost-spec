# Data Model: FallbackHint

## Entities

### FallbackHint (Enum)

Represents the recommendation from a plugin regarding whether the core system should attempt to
query other plugins (fallback) for a given resource.

| Value Name | Number | Description | Semantics |
| :--- | :--- | :--- | :--- |
| `FALLBACK_HINT_UNSPECIFIED` | 0 | Default value. Treat as "No Fallback". | Plugin has data or is the definitive source. Do not fallback. |
| `FALLBACK_HINT_NONE` | 1 | Explicitly no fallback needed. | Plugin has data (even if $0.00). Do not fallback. |
| `FALLBACK_HINT_RECOMMENDED` | 2 | Plugin has no data, suggests checking others. | Core SHOULD try other plugins if configured. |
| `FALLBACK_HINT_REQUIRED` | 3 | Plugin cannot handle this request type. | Core MUST try other plugins (if any) or fail. |

### GetActualCostResponse (Message)

Updated to include the fallback hint.

| Field Name | Type | Number | Description |
| :--- | :--- | :--- | :--- |
| `results` | `repeated ActualCostResult` | 1 | Existing field. The cost data. |
| `fallback_hint` | `FallbackHint` | 2 | **NEW**. The fallback recommendation. |

## Relationships

- `GetActualCostResponse` **has-one** `FallbackHint`.
- `GetActualCostResponse` contains `ActualCostResult` (0..N).

## Validation Rules

1. **Default Safety**: If `fallback_hint` is missing (wire format), it defaults to 0
   (`UNSPECIFIED`), which means "No Fallback". This prevents accidental infinite loops or
   fallback chains.
2. **Data Precedence**: If `results` is non-empty AND `fallback_hint` is `RECOMMENDED` or
   `REQUIRED`, the Core system SHOULD prefer the data and log a warning. Data presence implies
   the plugin did its job.
3. **Empty Data**: If `results` is empty `[]`, `fallback_hint` SHOULD be `RECOMMENDED` or
   `REQUIRED` unless the cost is genuinely unknown/zero and no fallback is desired (rare).
