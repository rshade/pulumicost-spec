# Quickstart: Validation Bypass Metadata

**Feature**: 034-validation-bypass-metadata
**Date**: 2026-01-24
**Phase**: 1 - Design

## Overview

This quickstart shows how to use the validation bypass metadata feature to create audit trails
when validation policies are intentionally bypassed.

## Basic Usage

### Recording a Bypass Event (Recommended: Builder Pattern)

The builder pattern is the **recommended approach** for creating bypass metadata. It provides
validation, automatic defaults, and reason truncation handling.

```go
package main

import (
    "github.com/rshade/finfocus-spec/sdk/go/pricing"
)

func validateWithBypass(forceBypass bool) pricing.ValidationResult {
    // Run validation
    err := checkBudgetLimit()
    if err != nil {
        if forceBypass {
            // Record the bypass for audit trail (recommended builder pattern)
            bypass := pricing.NewBypassMetadata(
                "budget_limit",
                err.Error(),
                pricing.WithReason("Emergency deployment approved by manager"),
                pricing.WithOperator(getCurrentUser()),
                pricing.WithSeverity(pricing.BypassSeverityError),
                pricing.WithMechanism(pricing.BypassMechanismFlag),
            )

            return pricing.ValidationResult{
                Valid:    true, // Passes because bypass was authorized
                Bypasses: []pricing.BypassMetadata{bypass},
            }
        }

        // Normal failure (no bypass)
        return pricing.ValidationResult{
            Valid:  false,
            Errors: []string{err.Error()},
        }
    }

    // Normal success
    return pricing.ValidationResult{Valid: true}
}
```

### Alternative: Direct Struct Initialization

For advanced use cases where you need full control, you can initialize the struct directly.
Note that this bypasses automatic timestamp setting and reason truncation.

```go
import "time"

// Direct struct initialization (advanced usage)
bypass := pricing.BypassMetadata{
    Timestamp:      time.Now().UTC(),
    ValidationName: "budget_limit",
    OriginalError:  "Cost exceeds budget by $500",
    Reason:         "Emergency deployment approved",
    Operator:       "user@example.com",
    Severity:       pricing.BypassSeverityError,
    Mechanism:      pricing.BypassMechanismFlag,
}

result := pricing.ValidationResult{
    Valid:    true,
    Bypasses: []pricing.BypassMetadata{bypass},
}
```

### Checking for Bypasses

```go
func handleValidationResult(result pricing.ValidationResult) {
    if len(result.Bypasses) > 0 {
        fmt.Println("⚠️  Validations were bypassed:")
        for _, b := range result.Bypasses {
            fmt.Printf("  - %s (%s): %s\n",
                b.ValidationName,
                b.Severity,
                b.Reason,
            )
        }
    }

    if !result.Valid {
        fmt.Println("❌ Validation failed:")
        for _, e := range result.Errors {
            fmt.Printf("  - %s\n", e)
        }
    }
}
```

## CLI Integration Example

```go
package main

import (
    "flag"
    "os"

    "github.com/rshade/finfocus-spec/sdk/go/pricing"
)

func main() {
    yolo := flag.Bool("yolo", false, "Bypass all validations")
    flag.Parse()

    // Determine bypass mechanism
    mechanism := pricing.BypassMechanismFlag
    if os.Getenv("FINFOCUS_BYPASS_VALIDATIONS") == "true" {
        mechanism = pricing.BypassMechanismEnvVar
    }

    result := runValidations(*yolo, mechanism)

    // Display results
    if len(result.Bypasses) > 0 {
        fmt.Println("\n⚠️  BYPASSED VALIDATIONS:")
        for _, b := range result.Bypasses {
            fmt.Printf("  [%s] %s\n", b.Severity, b.ValidationName)
            fmt.Printf("    Reason: %s\n", b.Reason)
            fmt.Printf("    Original error: %s\n", b.OriginalError)
            fmt.Printf("    Operator: %s\n", b.Operator)
            fmt.Printf("    Time: %s\n", b.Timestamp.Format(time.RFC3339))
        }
    }
}
```

## JSON Serialization

Bypass metadata serializes cleanly to JSON for logging and transmission:

```go
result := pricing.ValidationResult{
    Valid: true,
    Bypasses: []pricing.BypassMetadata{{
        Timestamp:      time.Date(2026, 1, 24, 10, 30, 0, 0, time.UTC),
        ValidationName: "budget_limit",
        OriginalError:  "Cost exceeds budget by $500",
        Reason:         "Emergency deployment",
        Operator:       "user@example.com",
        Severity:       pricing.BypassSeverityError,
        Mechanism:      pricing.BypassMechanismFlag,
    }},
}

jsonBytes, _ := json.MarshalIndent(result, "", "  ")
fmt.Println(string(jsonBytes))
```

Output:

```json
{
  "valid": true,
  "bypasses": [
    {
      "timestamp": "2026-01-24T10:30:00Z",
      "validation_name": "budget_limit",
      "original_error": "Cost exceeds budget by $500",
      "reason": "Emergency deployment",
      "operator": "user@example.com",
      "severity": "error",
      "mechanism": "flag"
    }
  ]
}
```

## Validation Helpers

```go
// Validate bypass metadata before recording
func recordBypass(b pricing.BypassMetadata) error {
    if err := pricing.ValidateBypassMetadata(b); err != nil {
        return fmt.Errorf("invalid bypass metadata: %w", err)
    }
    // Record to audit log...
    return nil
}

// Check if a severity is valid
if !pricing.IsValidBypassSeverity("error") {
    return errors.New("invalid severity")
}

// Check if a mechanism is valid
if !pricing.IsValidBypassMechanism("flag") {
    return errors.New("invalid mechanism")
}
```

## Best Practices

1. **Always set a reason**: Even brief reasons help with audit reviews
2. **Use UTC timestamps**: Avoid timezone confusion in distributed systems
3. **Capture operator identity**: Use service accounts for automated bypasses
4. **Choose appropriate severity**: Match to actual risk level
5. **Log bypasses**: Send to observability system for compliance tracking
6. **Don't bypass without need**: Only create bypass records for actual failures
7. **Implement retention policy**: The SDK provides data structures; callers must implement
   retention policies (90-day minimum recommended for quarterly compliance reviews)

## Common Patterns

### Multiple Bypasses

```go
result := pricing.ValidationResult{
    Valid: true,
    Bypasses: []pricing.BypassMetadata{
        {ValidationName: "budget_limit", Severity: pricing.BypassSeverityError, ...},
        {ValidationName: "region_policy", Severity: pricing.BypassSeverityWarning, ...},
    },
}
```

### Programmatic Bypass (API)

```go
bypass := pricing.BypassMetadata{
    Mechanism: pricing.BypassMechanismProgrammatic,
    Operator:  "api-service-account",
    Reason:    "Automated scaling event",
    // ...
}
```

### Config-Based Bypass

```go
bypass := pricing.BypassMetadata{
    Mechanism: pricing.BypassMechanismConfig,
    Operator:  "config:production.yaml",
    Reason:    "Pre-approved in configuration",
    // ...
}
```
