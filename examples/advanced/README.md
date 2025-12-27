# Advanced Implementation Examples

This directory contains runnable Go examples demonstrating advanced plugin
implementation patterns for the PulumiCost ecosystem.

## Examples

### tieredpricing/

Demonstrates tiered pricing calculation patterns for cloud resources with
volume-based pricing (like AWS S3 storage).

**Key concepts:**

- Parsing `pricing_tiers` from PricingSpec JSON
- Calculating costs across multiple tiers
- Generating detailed cost breakdowns
- Cross-provider pricing comparison

**Run:**

```bash
go run ./tieredpricing/
```

### multiprovider/

Demonstrates building plugins that support multiple cloud providers (AWS, Azure,
GCP) using the mapping package.

**Key concepts:**

- Provider-specific property extraction
- Strategy pattern for provider handling
- ResourceMatcher for multi-provider support
- Cross-provider region normalization

**Run:**

```bash
go run ./multiprovider/
```

## Related Documentation

- [Advanced Patterns Guide](../../docs/ADVANCED_PATTERNS.md) - Full documentation
- [Mapping Package](../../sdk/go/pluginsdk/mapping/doc.go) - Property extraction API
- [Plugin Developer Guide](../../PLUGIN_DEVELOPER_GUIDE.md) - Core plugin development

## Logging

These examples use Go's standard library `log/slog` package for simplicity and
portability. Production plugins should use `zerolog` as shown in the main SDK
documentation, but `slog` is perfectly acceptable and provides similar structured
logging capabilities with no external dependencies.

## Building and Testing

```bash
# Run all examples with full output
go run ./tieredpricing/
go run ./multiprovider/

# Quick validation (--quiet mode) - useful for CI/CD
go run ./tieredpricing/ --quiet
go run ./multiprovider/ --quiet

# Build to verify compilation
go build ./...

# Run tests
go test ./...
```

The `--quiet` flag is useful for CI/CD pipelines where you only need to verify
the examples execute successfully without verbose logging output.

## Pattern Summary

| Pattern | Example | Use Case |
|---------|---------|----------|
| Tiered Pricing | `tieredpricing/` | Storage, bandwidth with volume discounts |
| Multi-Provider | `multiprovider/` | Plugins supporting AWS, Azure, GCP |
| Property Extraction | `multiprovider/` | SKU/region from resource properties |
| Region Normalization | `multiprovider/` | Cross-provider cost comparison |
