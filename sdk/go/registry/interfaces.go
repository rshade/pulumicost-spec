package registry

import (
	"context"
)

// PluginRegistry defines the interface for plugin registries.
type PluginRegistry interface {
	// ListPlugins returns available plugins based on search criteria.
	ListPlugins(ctx context.Context, filter *PluginFilter, limit, offset int) ([]*PluginInfo, int, error)
	
	// GetPlugin returns detailed information about a specific plugin.
	GetPlugin(ctx context.Context, name, version string) (*PluginManifest, error)
	
	// GetPluginManifest returns the manifest for a specific plugin version.
	GetPluginManifest(ctx context.Context, name, version string) (*PluginManifest, error)
	
	// SearchPlugins performs full-text search across plugin metadata.
	SearchPlugins(ctx context.Context, query string, filter *PluginFilter, limit, offset int) ([]*PluginInfo, int, error)
	
	// RegisterPlugin registers a new plugin or updates an existing one.
	RegisterPlugin(ctx context.Context, manifest *PluginManifest, binaryData []byte, signature []byte) error
	
	// UnregisterPlugin removes a plugin from the registry.
	UnregisterPlugin(ctx context.Context, name, version, reason string) error
	
	// ValidatePlugin validates a plugin manifest and binary.
	ValidatePlugin(ctx context.Context, manifest *PluginManifest, binaryData []byte) (*ValidationResult, error)
}

// PluginDiscovery defines the interface for plugin discovery from various sources.
type PluginDiscovery interface {
	// Discover finds and returns all available plugins from the discovery source.
	Discover(ctx context.Context) ([]*PluginManifest, error)
	
	// Watch monitors for changes in plugin availability and calls the callback for updates.
	Watch(ctx context.Context, callback func(*PluginManifest)) error
	
	// Stop stops any ongoing discovery or watch operations.
	Stop() error
}

// PluginInstaller defines the interface for plugin installation and lifecycle management.
type PluginInstaller interface {
	// Install installs a plugin from a source (registry, URL, etc.).
	Install(ctx context.Context, source string, opts *InstallOptions) error
	
	// Update updates an installed plugin to a newer version.
	Update(ctx context.Context, name string, opts *UpdateOptions) error
	
	// Remove removes an installed plugin.
	Remove(ctx context.Context, name string) error
	
	// List returns all installed plugins.
	List(ctx context.Context) ([]*PluginInfo, error)
	
	// Validate validates an installed plugin.
	Validate(ctx context.Context, name string) (*ValidationResult, error)
	
	// Enable enables a disabled plugin.
	Enable(ctx context.Context, name string) error
	
	// Disable disables an active plugin.
	Disable(ctx context.Context, name string) error
}

// PluginValidator defines the interface for plugin validation.
type PluginValidator interface {
	// ValidateManifest validates a plugin manifest against the schema.
	ValidateManifest(manifest *PluginManifest) (*ValidationResult, error)
	
	// ValidateBinary validates a plugin binary for security and compatibility.
	ValidateBinary(binaryData []byte, manifest *PluginManifest) (*ValidationResult, error)
	
	// ValidateCompatibility checks if a plugin is compatible with the current system.
	ValidateCompatibility(manifest *PluginManifest) (*ValidationResult, error)
	
	// ScanSecurity performs security scanning on a plugin binary.
	ScanSecurity(binaryData []byte) (*SecurityScanResults, error)
}

// PluginLoader defines the interface for loading and executing plugins.
type PluginLoader interface {
	// Load loads a plugin into memory for execution.
	Load(ctx context.Context, name, version string) (Plugin, error)
	
	// Unload unloads a plugin from memory.
	Unload(ctx context.Context, name string) error
	
	// IsLoaded checks if a plugin is currently loaded.
	IsLoaded(name string) bool
	
	// GetLoaded returns all currently loaded plugins.
	GetLoaded() []Plugin
}

// Plugin defines the interface that all PulumiCost plugins must implement.
type Plugin interface {
	// Name returns the plugin name.
	Name() string
	
	// Version returns the plugin version.
	Version() string
	
	// Manifest returns the plugin manifest.
	Manifest() *PluginManifest
	
	// Start starts the plugin and makes it ready to serve requests.
	Start(ctx context.Context) error
	
	// Stop stops the plugin and cleans up resources.
	Stop(ctx context.Context) error
	
	// IsHealthy returns the current health status of the plugin.
	IsHealthy(ctx context.Context) bool
}

// RegistryClient defines the interface for interacting with remote plugin registries.
type RegistryClient interface {
	// GetRegistry retrieves the complete registry information.
	GetRegistry(ctx context.Context) (*PluginRegistry, error)
	
	// ListPlugins lists available plugins with optional filtering.
	ListPlugins(ctx context.Context, filter *PluginFilter, limit, offset int) ([]*PluginInfo, int, error)
	
	// GetPlugin retrieves detailed information about a specific plugin.
	GetPlugin(ctx context.Context, name, version string) (*PluginInfo, error)
	
	// GetManifest retrieves the manifest for a specific plugin version.
	GetManifest(ctx context.Context, name, version string) (*PluginManifest, error)
	
	// SearchPlugins performs a search across the registry.
	SearchPlugins(ctx context.Context, query string, limit, offset int) ([]*PluginInfo, int, error)
	
	// DownloadPlugin downloads a plugin binary from the registry.
	DownloadPlugin(ctx context.Context, name, version string) ([]byte, error)
}

// RegistryServer defines the interface for implementing plugin registry servers.
type RegistryServer interface {
	// Serve starts the registry server.
	Serve(ctx context.Context, address string) error
	
	// Stop stops the registry server.
	Stop(ctx context.Context) error
	
	// AddPlugin adds a plugin to the registry.
	AddPlugin(ctx context.Context, manifest *PluginManifest, binaryData []byte) error
	
	// RemovePlugin removes a plugin from the registry.
	RemovePlugin(ctx context.Context, name, version string) error
	
	// UpdatePlugin updates a plugin in the registry.
	UpdatePlugin(ctx context.Context, manifest *PluginManifest, binaryData []byte) error
}