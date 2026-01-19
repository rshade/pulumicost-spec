# Quickstart: Capability Discovery

## How to Enable Capabilities

### Method 1: Automatic Discovery (Recommended)

Simply implement the standard interfaces in your plugin struct. The SDK will automatically detect them.

```go
type MyPlugin struct {
    pbc.UnimplementedCostSourceServiceServer
}

// Implementing DryRunHandler automatically enables PLUGIN_CAPABILITY_DRY_RUN
func (p *MyPlugin) HandleDryRun(req *pbc.DryRunRequest) (*pbc.DryRunResponse, error) {
    return pluginsdk.NewDryRunResponse(...), nil
}

func main() {
    // No extra config needed!
    info := pluginsdk.NewPluginInfo("my-plugin", "v1.0.0")
    
    // SDK starts serving...
    pluginsdk.Serve(info, &MyPlugin{})
}
```

### Method 2: Manual Override (Advanced)

If you need to force capabilities (e.g., to disable a feature dynamically or if you are proxying), use `WithCapabilities`.

**Warning**: This completely replaces auto-discovery. You must list ALL capabilities you want to support.

```go
func main() {
    // Only advertising DryRun, ignoring any other implemented interfaces
    info := pluginsdk.NewPluginInfo("my-plugin", "v1.0.0",
        pluginsdk.WithCapabilities([]pbc.PluginCapability{
            pbc.PluginCapability_PLUGIN_CAPABILITY_DRY_RUN,
        }),
    )
    
    pluginsdk.Serve(info, &MyPlugin{})
}
```
