# Quickstart: GetBudgets RPC Implementation

**Date**: 2025-12-09
**Feature**: GetBudgets RPC for Plugin-Provided Budget Information

## Overview

This guide provides a quick start for implementing the GetBudgets RPC in plugins and integrating
budget functionality into pulumicost-core.

## For Plugin Developers

### 1. Implement BudgetsProvider Interface

```go
type BudgetsProvider interface {
    GetBudgets(ctx context.Context, req *pbc.GetBudgetsRequest) (*pbc.GetBudgetsResponse, error)
}
```

### 2. Handle Optional RPC

```go
func (s *Server) GetBudgets(ctx context.Context, req *pbc.GetBudgetsRequest) (*pbc.GetBudgetsResponse, error) {
    if bp, ok := s.impl.(BudgetsProvider); ok {
        return bp.GetBudgets(ctx, req)
    }
    return nil, status.Error(codes.Unimplemented, "plugin does not support GetBudgets")
}
```

### 3. Map Provider Data to Proto Messages

```go
// Example AWS Budgets mapping
budget := &pbc.Budget{
    Id:   awsBudget.BudgetName,
    Name: awsBudget.BudgetName,
    Source: "aws-budgets",
    Amount: &pbc.BudgetAmount{
        Limit:    awsBudget.BudgetLimit.Amount,
        Currency: awsBudget.BudgetLimit.Unit,
    },
    Period: mapAWSTimeUnit(awsBudget.TimeUnit),
    // ... map other fields
}
```

## For Core Developers

### 1. Update Proto Definitions

Add budget messages to `proto/pulumicost/v1/budget.proto`:

```protobuf
message Budget {
  string id = 1;
  string name = 2;
  string source = 3;
  // ... complete message definition
}
```

Add RPC to `proto/pulumicost/v1/costsource.proto`:

```protobuf
service CostSource {
  // ... existing RPCs
  rpc GetBudgets(GetBudgetsRequest) returns (GetBudgetsResponse);
}
```

### 2. Regenerate SDK

```bash
make generate
```

This updates `sdk/go/proto/` with new message types and gRPC service interfaces.

### 3. Add Conformance Tests

Create tests in `sdk/go/testing/` following existing patterns:

```go
func TestGetBudgets(t *testing.T) {
    // Test basic RPC structure
    // Test cross-provider mapping
    // Test performance with large budget sets
}
```

### 4. Update Documentation

Add examples to `examples/` and update README files.

## Development Workflow

1. **Proto First**: Update protobuf definitions
2. **Generate**: Run `make generate` to create SDK
3. **Test**: Write conformance tests before implementation
4. **Implement**: Add plugin support following interface
5. **Validate**: Run full test suite with `make validate`

## Key Considerations

- **Optional RPC**: Plugins may return `Unimplemented` if budgets not supported
- **Performance**: Use `include_status` flag to optimize live data fetching
- **Cross-Provider**: Maintain consistent field mapping across AWS, GCP, Azure, Kubecost
- **Backward Compatibility**: New RPC doesn't break existing plugins

## Testing Checklist

- [ ] Proto definitions compile without errors
- [ ] SDK generates successfully
- [ ] Basic conformance tests pass
- [ ] Cross-provider examples validate
- [ ] Performance benchmarks meet requirements
- [ ] buf linting and breaking change checks pass
