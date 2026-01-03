# Quickstart: JSON-LD Serialization

Get started with JSON-LD serialization for FOCUS cost data in under 5 minutes.

## Installation

The jsonld package is part of the pulumicost-spec Go SDK:

```go
import "github.com/rshade/pulumicost-spec/sdk/go/jsonld"
```

## Basic Usage

### Serialize a Single Record

```go
package main

import (
    "fmt"
    "time"

    "github.com/rshade/pulumicost-spec/sdk/go/jsonld"
    "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
    pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

func main() {
    // Build a FOCUS cost record using the builder
    record := pluginsdk.NewFocusRecordBuilder().
        WithBillingAccount("123456789012", "My AWS Account").
        WithChargePeriod(time.Now().AddDate(0, 0, -1), time.Now()).
        WithService(pbc.FocusServiceCategory_FOCUS_SERVICE_CATEGORY_COMPUTE, "Amazon EC2").
        WithBilledCost(125.50, "USD").
        Build()

    // Create serializer with default options
    serializer := jsonld.NewSerializer()

    // Serialize to JSON-LD
    output, err := serializer.Serialize(record)
    if err != nil {
        panic(err)
    }

    fmt.Println(string(output))
}
```

**Output:**

```json
{
  "@context": {
    "schema": "https://schema.org/",
    "focus": "https://focus.finops.org/v1#"
  },
  "@type": "FocusCostRecord",
  "@id": "urn:focus:cost:a1b2c3d4e5f6...",
  "billingAccountId": "123456789012",
  "billingAccountName": "My AWS Account",
  "chargePeriodStart": "2025-12-30T00:00:00Z",
  "chargePeriodEnd": "2025-12-31T00:00:00Z",
  "serviceCategory": "COMPUTE",
  "serviceName": "Amazon EC2",
  "billedCost": {
    "@type": "MonetaryAmount",
    "value": 125.50,
    "currency": "USD"
  }
}
```

### Batch Serialization (Streaming)

For large datasets, use streaming to bound memory usage:

```go
func serializeBatch(records []*pbc.FocusCostRecord, w io.Writer) error {
    serializer := jsonld.NewSerializer()

    // Stream records to writer (file, network, etc.)
    return serializer.SerializeStream(
        sliceToChannel(records),
        w,
    )
}

func sliceToChannel(records []*pbc.FocusCostRecord) <-chan *pbc.FocusCostRecord {
    ch := make(chan *pbc.FocusCostRecord)
    go func() {
        defer close(ch)
        for _, r := range records {
            ch <- r
        }
    }()
    return ch
}
```

### Custom Context Configuration

Override default vocabulary mappings for enterprise integration:

```go
// Create custom context with additional ontology
ctx := jsonld.NewContext().
    WithRemoteContext("https://your-org.com/ontology/v1").
    WithCustomMapping("billingAccountId", "yourOrg:accountIdentifier").
    WithCustomMapping("serviceName", "yourOrg:cloudService")

serializer := jsonld.NewSerializer(
    jsonld.WithContext(ctx),
)
```

### User-Provided IDs

Supply your own record identifiers:

```go
serializer := jsonld.NewSerializer(
    jsonld.WithUserIDField("invoice_id"), // Use invoice_id as @id
)
```

### Pretty-Print Output

For debugging or human-readable output:

```go
serializer := jsonld.NewSerializer(
    jsonld.WithPrettyPrint(true),
)
```

## Configuration Options

| Option | Default | Description |
|--------|---------|-------------|
| `WithContext(ctx)` | Default FOCUS context | Custom JSON-LD context |
| `WithOmitEmpty(bool)` | `true` | Omit empty/zero/nil fields |
| `WithIRIEnums(bool)` | `false` | Serialize enums as full IRIs |
| `WithDeprecated(bool)` | `true` | Include deprecated fields |
| `WithPrettyPrint(bool)` | `false` | Indent JSON output |
| `WithUserIDField(field)` | none | Field to use as @id |
| `WithIDPrefix(prefix)` | `urn:focus:cost:` | Prefix for generated IDs |

## Performance Tips

1. **Reuse serializer instances** - Create once, use many times
2. **Use streaming for batches** - Avoids loading all records into memory
3. **Disable pretty-print in production** - Saves ~20% output size
4. **Pre-allocate channels** - Buffer size = expected batch size / 10

## Next Steps

- See [data-model.md](data-model.md) for complete field mappings
- See [contracts/](contracts/) for JSON-LD context definitions
- Run benchmarks: `go test -bench=. ./sdk/go/jsonld/`
