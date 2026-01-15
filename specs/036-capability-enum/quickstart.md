# Quickstart: Using Plugin Capabilities

## For Plugin Developers

If you are using the FinFocus Go SDK, capabilities are automatically detected. Simply implement the relevant interfaces:

```go
type MyPlugin struct {
    pluginsdk.UnimplementedPlugin
}

// Implementing this interface automatically adds PLUGIN_CAPABILITY_RECOMMENDATIONS
func (p *MyPlugin) GetRecommendations(ctx context.Context, req *pbc.GetRecommendationsRequest) (*pbc.GetRecommendationsResponse, error) {
    // ...
}
```

## For Consumer Application Developers

Query the plugin info to see what it supports:

```go
info, err := client.GetPluginInfo(ctx, &pbc.GetPluginInfoRequest{})
if err != nil {
    // Handle error
}

for _, cap := range info.Capabilities {
    if cap == pbc.PluginCapability_PLUGIN_CAPABILITY_RECOMMENDATIONS {
        // Plugin supports recommendations!
    }
}
```
