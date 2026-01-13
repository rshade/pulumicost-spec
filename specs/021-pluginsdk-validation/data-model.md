# Data Model: Validation Logic

## Logical Entities

### ProjectedCostValidation

Validates `finfocus.v1.GetProjectedCostRequest`.

**Fields Checked**:

1. **Request**: Must not be nil.
2. **Resource**: Must not be nil.
3. **Resource.Provider**: Must not be empty.
4. **Resource.ResourceType**: Must not be empty.
5. **Resource.Sku**: Must not be empty.
   - *Error Guidance*: "use mapping.ExtractAWSSKU()" (if provider is AWS, or generic guidance).
6. **Resource.Region**: Must not be empty.
   - *Error Guidance*: "use mapping.ExtractAWSRegion()".

### ActualCostValidation

Validates `finfocus.v1.GetActualCostRequest`.

**Fields Checked**:

1. **Request**: Must not be nil.
2. **ResourceId**: Must not be empty.
3. **StartTime**: Must not be nil (and valid Timestamp).
4. **EndTime**: Must not be nil (and valid Timestamp).
5. **TimeRange**: EndTime > StartTime.

## Error Schema

Standard Go `error` interface.
Format: `"[field_name] is required[: guidance]"` or `"[rule description]"`
