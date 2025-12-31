# Quickstart: Add GetPluginInfo RPC

## Prerequisites

- Go 1.25+
- Protobuf compiler (`protoc`)
- `buf` CLI

## Steps

1. **Update Proto Definition**:
   - Modify `proto/pulumicost/v1/costsource.proto` to include `GetPluginInfo` RPC and messages.

2. **Regenerate SDK**:
   - Run `make generate` (or `buf generate`) to update `sdk/go/proto/`.

3. **Update SDK Implementation**:
   - In `sdk/go/pluginsdk/`:
     - Add `SpecVersion` constant (e.g., in `version.go`).
     - Implement `GetPluginInfo` in the base plugin handler (e.g., `base.go` or `server.go`).

4. **Verify**:
   - Run `go test ./sdk/go/...` to ensure no regressions.
   - Create a new conformance test case in `sdk/go/testing/` to call `GetPluginInfo` and validate the response.

## Usage Example (Client)

```go
client := pbc.NewCostSourceServiceClient(conn)
resp, err := client.GetPluginInfo(ctx, &pbc.GetPluginInfoRequest{})
if err != nil {
    if status.Code(err) == codes.Unimplemented {
        log.Warn("Plugin does not support GetPluginInfo (legacy version)")
        // Handle legacy plugin
    } else {
        log.Error("Failed to get plugin info", "error", err)
    }
} else {
    fmt.Printf("Plugin: %s %s (Spec: %s)\n", resp.Name, resp.Version, resp.SpecVersion)
}
```
