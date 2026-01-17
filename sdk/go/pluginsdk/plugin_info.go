// Package pluginsdk provides plugin information types and builders.
package pluginsdk

import (
	"errors"
	"fmt"

	pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

// PluginInfo holds metadata about a plugin that can be reported via GetPluginInfo RPC.
// This information enables compatibility verification and diagnostic visibility.
type PluginInfo struct {
	// Name is the unique identifier for this plugin (required).
	Name string

	// Version is the plugin version following semantic versioning (required).
	// Example: "v1.0.0", "v2.1.3"
	Version string

	// SpecVersion is the finfocus-spec version this plugin implements (required).
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

	// Capabilities declares the functional capabilities this plugin provides.
	// When set, these are returned in GetPluginInfoResponse.capabilities.
	// When nil or empty, capabilities are auto-discovered from implemented interfaces.
	// This allows plugins to override auto-discovery if needed.
	Capabilities []pbc.PluginCapability
}

// PluginInfoOption is a functional option for configuring PluginInfo.
type PluginInfoOption func(*PluginInfo)

// NewPluginInfo creates a new PluginInfo with the given name and version.
// The SpecVersion is automatically set to the SDK's current SpecVersion constant.
// The PluginInfo should be treated as immutable after passing to ServeConfig.
// Modifying the original slices/maps after construction may cause data races.
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

// WithCapabilities sets the plugin capabilities.
// When set, these override auto-discovery from implemented interfaces.
// A defensive copy is made to prevent aliasing issues.
//
// Example:
//
//	WithCapabilities(pbc.PluginCapability_PLUGIN_CAPABILITY_PROJECTED_COSTS,
//	                 pbc.PluginCapability_PLUGIN_CAPABILITY_RECOMMENDATIONS)
func WithCapabilities(capabilities ...pbc.PluginCapability) PluginInfoOption {
	return func(info *PluginInfo) {
		// Defensive copy to prevent aliasing issues
		info.Capabilities = make([]pbc.PluginCapability, len(capabilities))
		copy(info.Capabilities, capabilities)
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

// inferCapabilities determines plugin capabilities by checking implemented interfaces.
// The type assertions themselves are zero-allocation, but inferCapabilities may
// allocate as needed when append grows the capabilities slice.
// Returns a slice of capabilities supported by the plugin.
//
// The base Plugin interface methods (GetProjectedCost, GetActualCost, etc.) are
// always assumed to be implemented since they are required by the interface.
// Only optional interfaces are checked via type assertion.
func inferCapabilities(plugin Plugin) []pbc.PluginCapability {
	var capabilities []pbc.PluginCapability

	// Base capabilities - always present since required by Plugin interface
	capabilities = append(capabilities,
		pbc.PluginCapability_PLUGIN_CAPABILITY_PROJECTED_COSTS,
		pbc.PluginCapability_PLUGIN_CAPABILITY_ACTUAL_COSTS,
		pbc.PluginCapability_PLUGIN_CAPABILITY_PRICING_SPEC,
		pbc.PluginCapability_PLUGIN_CAPABILITY_ESTIMATE_COST,
	)

	// Check optional interfaces using zero-allocation type assertions
	if _, ok := plugin.(RecommendationsProvider); ok {
		capabilities = append(capabilities, pbc.PluginCapability_PLUGIN_CAPABILITY_RECOMMENDATIONS)
	}
	if _, ok := plugin.(BudgetsProvider); ok {
		capabilities = append(capabilities, pbc.PluginCapability_PLUGIN_CAPABILITY_BUDGETS)
	}
	if _, ok := plugin.(DismissProvider); ok {
		capabilities = append(capabilities, pbc.PluginCapability_PLUGIN_CAPABILITY_DISMISS_RECOMMENDATIONS)
	}
	if _, ok := plugin.(DryRunHandler); ok {
		capabilities = append(capabilities, pbc.PluginCapability_PLUGIN_CAPABILITY_DRY_RUN)
	}

	return capabilities
}

// legacyCapabilityMap provides backward compatibility by mapping PluginCapability
// enums to the legacy string-based capability keys used in metadata.
//
// This maintains compatibility with clients that expect the old string-based
// capability reporting while the new enum-based approach is adopted.
//
// PLUGIN_CAPABILITY_UNSPECIFIED (value 0) is intentionally excluded from this map.
// It is the protobuf default/zero value and should never be used as an actual capability.
//
//nolint:exhaustive,gochecknoglobals // UNSPECIFIED intentionally excluded from legacy mapping
var legacyCapabilityMap = map[pbc.PluginCapability]string{
	pbc.PluginCapability_PLUGIN_CAPABILITY_PROJECTED_COSTS:         "projected_costs",
	pbc.PluginCapability_PLUGIN_CAPABILITY_ACTUAL_COSTS:            "actual_costs",
	pbc.PluginCapability_PLUGIN_CAPABILITY_PRICING_SPEC:            "pricing_spec",
	pbc.PluginCapability_PLUGIN_CAPABILITY_ESTIMATE_COST:           "estimate_cost",
	pbc.PluginCapability_PLUGIN_CAPABILITY_CARBON:                  "carbon",
	pbc.PluginCapability_PLUGIN_CAPABILITY_RECOMMENDATIONS:         "recommendations",
	pbc.PluginCapability_PLUGIN_CAPABILITY_DISMISS_RECOMMENDATIONS: "dismiss_recommendations",
	pbc.PluginCapability_PLUGIN_CAPABILITY_DRY_RUN:                 "dry_run",
	pbc.PluginCapability_PLUGIN_CAPABILITY_BUDGETS:                 "budgets",
	pbc.PluginCapability_PLUGIN_CAPABILITY_ENERGY:                  "energy",
	pbc.PluginCapability_PLUGIN_CAPABILITY_WATER:                   "water",
}

// capabilitiesToLegacyMetadata converts a slice of PluginCapability enums
// to the legacy metadata map format for backward compatibility.
//
// Returns a map[string]bool where the keys are legacy capability names
// and values are always true (presence indicates capability support).
func capabilitiesToLegacyMetadata(capabilities []pbc.PluginCapability) map[string]bool {
	if len(capabilities) == 0 {
		return nil
	}

	metadata := make(map[string]bool, len(capabilities))
	for _, cap := range capabilities {
		if key, exists := legacyCapabilityMap[cap]; exists {
			metadata[key] = true
		}
	}
	return metadata
}
