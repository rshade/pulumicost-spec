# Data Model: JSON-LD Serialization

**Feature**: 032-jsonld-serialization
**Date**: 2025-12-31

## Entities

### 1. Serializer

Primary interface for converting FOCUS records to JSON-LD.

| Field       | Type              | Description                        |
| ----------- | ----------------- | ---------------------------------- |
| context     | \*Context         | JSON-LD context configuration      |
| idGenerator | IDGenerator       | Strategy for generating @id values |
| options     | SerializerOptions | Serialization behavior options     |

**Methods**:

- `Serialize(record *pbc.FocusCostRecord) ([]byte, error)` - Single record
- `SerializeCommitment(record *pbc.ContractCommitment) ([]byte, error)` - Commitment record
- `SerializeStream(records <-chan *pbc.FocusCostRecord, w io.Writer) error` - Batch streaming

**Error Types**:

- `*ValidationError`: Context validation error with fields:
  - Field: string (field or configuration element causing error)
  - Message: string (human-readable error description)
  - Suggestion: string (optional fix recommendation)

### 2. Context

JSON-LD context configuration for vocabulary mapping.

| Field          | Type                   | Description                                   |
| -------------- | ---------------------- | --------------------------------------------- |
| baseURL        | string                 | Base IRI for relative references              |
| schemaOrg      | bool                   | Include Schema.org vocabulary (default: true) |
| focusNamespace | string                 | FOCUS vocabulary namespace URL                |
| customMappings | map[string]interface{} | Custom property mappings                      |
| remoteContexts | []string               | Additional remote context URLs                |

**Methods**:

- `Build() map[string]interface{}` - Generate @context object
- `WithCustomMapping(field, iri string) *Context` - Add custom mapping
- `WithRemoteContext(url string) *Context` - Add remote context

### 3. IDGenerator

Strategy interface for generating `@id` values.

| Field       | Type   | Description                              |
| ----------- | ------ | ---------------------------------------- |
| prefix      | string | IRI prefix (e.g., "urn:focus:cost:")     |
| userIDField | string | Optional field name for user-provided ID |

**Methods**:

- `Generate(record *pbc.FocusCostRecord) string` - Generate ID for cost record
- `GenerateCommitment(record *pbc.ContractCommitment) string` - Generate ID for commitment

**ID Generation Algorithm**:

```text
1. If userIDField is set and record has non-empty value → use that value
2. Otherwise, compute SHA256(billing_account_id + "|" + charge_period_start + "|" + resource_id)
3. Return prefix + hex(hash[:16]) (first 16 bytes = 32 hex chars)
```

### 4. SerializerOptions

Configuration options for serialization behavior.

| Field             | Type   | Default | Description                               |
| ----------------- | ------ | ------- | ----------------------------------------- |
| OmitEmptyFields   | bool   | true    | Skip empty/zero/nil fields                |
| UseIRIEnums       | bool   | false   | Serialize enums as full IRIs              |
| IncludeDeprecated | bool   | true    | Include deprecated fields with annotation |
| PrettyPrint       | bool   | false   | Indent JSON output                        |
| DateFormat        | string | RFC3339 | Timestamp format (ISO 8601)               |

### 5. JSONLDDocument

Output structure for serialized JSON-LD.

| Field      | JSON Key | Type        | Description                                       |
| ---------- | -------- | ----------- | ------------------------------------------------- |
| Context    | @context | interface{} | Vocabulary context                                |
| ID         | @id      | string      | Unique identifier                                 |
| Type       | @type    | string      | Record type (FocusCostRecord, ContractCommitment) |
| Properties | (inline) | various     | All FOCUS fields mapped to JSON-LD properties     |

## Field Mappings

### FocusCostRecord → JSON-LD

| Proto Field           | JSON-LD Property          | Schema.org Type        | Notes                |
| --------------------- | ------------------------- | ---------------------- | -------------------- |
| billing_account_id    | focus:billingAccountId    | -                      | Required             |
| billing_account_name  | focus:billingAccountName  | -                      | Optional             |
| charge_period_start   | focus:chargePeriodStart   | schema:DateTime        | ISO 8601             |
| charge_period_end     | focus:chargePeriodEnd     | schema:DateTime        | ISO 8601             |
| billed_cost           | focus:billedCost          | schema:MonetaryAmount  | With currency        |
| list_cost             | focus:listCost            | schema:MonetaryAmount  | With currency        |
| effective_cost        | focus:effectiveCost       | schema:MonetaryAmount  | With currency        |
| billing_currency      | (embedded)                | schema:currency        | ISO 4217             |
| service_name          | focus:serviceName         | schema:Service.name    | Optional             |
| resource_id           | focus:resourceId          | -                      | Optional             |
| resource_name         | focus:resourceName        | schema:Thing.name      | Optional             |
| region_id             | focus:regionId            | -                      | Optional             |
| region_name           | focus:regionName          | schema:Place.name      | Optional             |
| tags                  | focus:tags                | schema:PropertyValue[] | Key-value pairs      |
| extended_columns      | focus:extendedColumns     | -                      | Provider extensions  |
| provider_name         | focus:providerName        | -                      | @deprecated          |
| publisher             | focus:publisher           | -                      | @deprecated          |
| service_provider_name | focus:serviceProviderName | -                      | FOCUS 1.3            |
| host_provider_name    | focus:hostProviderName    | -                      | FOCUS 1.3            |
| allocated_method_id   | focus:allocatedMethodId   | -                      | FOCUS 1.3 allocation |
| allocated_resource_id | focus:allocatedResourceId | -                      | FOCUS 1.3 allocation |
| contract_applied      | focus:contractApplied     | -                      | Links to commitment  |

### ContractCommitment → JSON-LD

| Proto Field                  | JSON-LD Property                 | Schema.org Type       | Notes           |
| ---------------------------- | -------------------------------- | --------------------- | --------------- |
| contract_commitment_id       | focus:contractCommitmentId       | -                     | Primary key     |
| contract_id                  | focus:contractId                 | -                     | Parent contract |
| contract_commitment_category | focus:contractCommitmentCategory | -                     | SPEND/USAGE     |
| contract_commitment_cost     | focus:contractCommitmentCost     | schema:MonetaryAmount | Monetary        |
| contract_commitment_quantity | focus:contractCommitmentQuantity | -                     | Usage amount    |
| billing_currency             | (embedded)                       | schema:currency       | ISO 4217        |

### MonetaryAmount Structure

Cost fields are serialized as Schema.org MonetaryAmount:

```json
{
  "@type": "schema:MonetaryAmount",
  "schema:value": 123.45,
  "schema:currency": "USD"
}
```

### Tag/Map Serialization

Tags and maps are serialized as PropertyValue arrays:

```json
{
  "focus:tags": [
    {
      "@type": "schema:PropertyValue",
      "schema:name": "environment",
      "schema:value": "production"
    }
  ]
}
```

## Validation Rules

### Required Fields

- FocusCostRecord: `billing_account_id` must be non-empty for ID generation
- ContractCommitment: `contract_commitment_id` must be non-empty

### Field Constraints

- Timestamps: Must be valid google.protobuf.Timestamp (>= 1970-01-01)
- Currency: Must be 3-letter ISO 4217 code
- Enums: Must be valid proto enum values (not UNSPECIFIED unless explicitly allowed)

### State Transitions

N/A - Serialization is stateless transformation.

## Relationships

```text
FocusCostRecord ──────┬──── references ────> ContractCommitment
                      │    (via contract_applied)
                      │
                      └──── produces ────> JSONLDDocument
                                          ├── @context
                                          ├── @id
                                          ├── @type
                                          └── properties

Context ────────────────────> embeds ────> Schema.org vocabulary
                                          FOCUS vocabulary
                                          Custom mappings

IDGenerator ────────────────> generates ──> @id values
                                           (user-provided or hash)
```
