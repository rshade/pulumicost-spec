# SDK Conformance Testing

## Performance Conformance

### GetPluginInfo Latency

- **Standard Conformance**: ≤100ms average across 10 iterations
- **Advanced Conformance**: ≤50ms average across 10 iterations
- **Legacy Plugin Handling**: Unimplemented errors handled gracefully

### Running Performance Tests

```bash
# Standard conformance
go test -v ./sdk/go/testing -run Conformance

# Advanced conformance
CONFORMANCE_LEVEL=advanced go test -v ./sdk/go/testing -run Conformance
```
