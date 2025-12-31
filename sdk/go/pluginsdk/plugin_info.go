// Package pluginsdk provides plugin information types and builders.
package pluginsdk

import (
	"errors"
	"fmt"
)

// PluginInfo holds metadata about a plugin that can be reported via GetPluginInfo RPC.
// This information enables compatibility verification and diagnostic visibility.
type PluginInfo struct {
	// Name is the unique identifier for this plugin (required).
	Name string

	// Version is the plugin version following semantic versioning (required).
	// Example: "v1.0.0", "v2.1.3"
	Version string

	// SpecVersion is the pulumicost-spec version this plugin implements (required).
	// This value should match the SDK's SpecVersion constant unless the plugin
	// is designed to work with a specific older/newer spec version.
	// Example: "v0.4.11"
	SpecVersion string

	// Providers is a list of cloud providers this plugin supports.
	// Examples: ["aws"], ["azure", "gcp"], ["kubernetes"]
	// An empty list indicates the plugin is provider-agnostic.
	Providers []string

	// Metadata is optional key-value pairs for diagnostic information.
	// Examples: {"build_date": "2024-01-15", "git_commit": "abc123"}
	Metadata map[string]string
}

// PluginInfoOption is a functional option for configuring PluginInfo.
type PluginInfoOption func(*PluginInfo)

// NewPluginInfo creates a new PluginInfo with the given name and version.
// The SpecVersion is automatically set to the SDK's current SpecVersion constant.
// Use functional options to customize additional fields.
//
// Example:
//
//	info := NewPluginInfo("aws-cost-plugin", "v1.0.0",
//	    WithProviders("aws"),
//	    WithMetadata("build_date", "2024-01-15"),
//	)
func NewPluginInfo(name, version string, opts ...PluginInfoOption) *PluginInfo {
	info := &PluginInfo{
		Name:        name,
		Version:     version,
		SpecVersion: SpecVersion, // Use SDK's current spec version
	}

	for _, opt := range opts {
		opt(info)
	}

	return info
}

// WithSpecVersion overrides the spec version for the plugin.
// Use this only if your plugin needs to report a different spec version than the SDK's default.
func WithSpecVersion(specVersion string) PluginInfoOption {
	return func(info *PluginInfo) {
		info.SpecVersion = specVersion
	}
}

// WithProviders sets the list of cloud providers this plugin supports.
// Multiple providers can be specified in a single call.
// A defensive copy is made to prevent aliasing issues if the caller modifies
// the slice after construction.
//
// Example:
//
//	WithProviders("aws", "azure", "gcp")
func WithProviders(providers ...string) PluginInfoOption {
	return func(info *PluginInfo) {
		// Defensive copy to prevent aliasing issues
		info.Providers = make([]string, len(providers))
		copy(info.Providers, providers)
	}
}

// WithMetadata adds a key-value pair to the plugin's metadata.
// Multiple calls can be chained to add multiple metadata entries.
//
// Example:
//
//	NewPluginInfo("plugin", "v1.0.0",
//	    WithMetadata("build_date", "2024-01-15"),
//	    WithMetadata("git_commit", "abc123"),
//	)
func WithMetadata(key, value string) PluginInfoOption {
	return func(info *PluginInfo) {
		if info.Metadata == nil {
			info.Metadata = make(map[string]string)
		}
		info.Metadata[key] = value
	}
}

// WithMetadataMap sets the entire metadata map at once.
// This replaces any existing metadata.
// A defensive copy is made to prevent aliasing issues if the caller modifies
// the map after construction.
//
// Example:
//
//	WithMetadataMap(map[string]string{
//	    "build_date": "2024-01-15",
//	    "git_commit": "abc123",
//	})
func WithMetadataMap(metadata map[string]string) PluginInfoOption {
	return func(info *PluginInfo) {
		// Defensive copy to prevent aliasing issues
		info.Metadata = make(map[string]string, len(metadata))
		for k, v := range metadata {
			info.Metadata[k] = v
		}
	}
}

// Validate checks that the PluginInfo has all required fields and they are valid.
// Returns an error if validation fails. Returns an error if info is nil.
//
// Error messages include the actual invalid values to aid debugging. This is safe
// because plugin names, versions, and spec versions are not sensitive data.
func (info *PluginInfo) Validate() error {
	if info == nil {
		return errors.New("plugin info validation failed: PluginInfo is nil")
	}
	if info.Name == "" {
		return errors.New("plugin info validation failed: name is required (got empty string)")
	}
	if info.Version == "" {
		return fmt.Errorf("plugin info validation failed: version is required for plugin %q (got empty string)",
			info.Name)
	}
	if info.SpecVersion == "" {
		return fmt.Errorf("plugin info validation failed: spec_version is required for plugin %q (got empty string)",
			info.Name)
	}

	// Validate spec_version is a valid SemVer
	if err := ValidateSpecVersion(info.SpecVersion); err != nil {
		return fmt.Errorf("plugin info validation failed for plugin %q: invalid spec_version %q: %w",
			info.Name, info.SpecVersion, err)
	}

	return nil
}
