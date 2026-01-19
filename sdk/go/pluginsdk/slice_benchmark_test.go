package pluginsdk_test

import (
	"testing"

	pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

// BenchmarkSliceCopyPatterns compares the performance of different slice copying patterns.
//
// Pattern 1: make + copy
// Pattern 2: append(nil, src...)
//
// We test with typical sizes for providers (1-5) and capabilities (4-12).
func BenchmarkSliceCopyPatterns(b *testing.B) {
	// Setup typical data
	providers := []string{"aws", "azure", "gcp", "kubernetes"}
	capabilities := []pbc.PluginCapability{
		pbc.PluginCapability_PLUGIN_CAPABILITY_DRY_RUN,
		pbc.PluginCapability_PLUGIN_CAPABILITY_RECOMMENDATIONS,
		pbc.PluginCapability_PLUGIN_CAPABILITY_BUDGETS,
		pbc.PluginCapability_PLUGIN_CAPABILITY_PROJECTED_COSTS,
		pbc.PluginCapability_PLUGIN_CAPABILITY_ACTUAL_COSTS,
		pbc.PluginCapability_PLUGIN_CAPABILITY_PRICING_SPEC,
		pbc.PluginCapability_PLUGIN_CAPABILITY_ESTIMATE_COST,
		pbc.PluginCapability_PLUGIN_CAPABILITY_DISMISS_RECOMMENDATIONS,
	}

	b.Run("Providers_MakeCopy", func(b *testing.B) {
		for range b.N {
			dst := make([]string, len(providers))
			copy(dst, providers)
		}
	})

	b.Run("Providers_Append", func(b *testing.B) {
		for range b.N {
			_ = append([]string(nil), providers...)
		}
	})

	b.Run("Capabilities_MakeCopy", func(b *testing.B) {
		for range b.N {
			dst := make([]pbc.PluginCapability, len(capabilities))
			copy(dst, capabilities)
		}
	})

	b.Run("Capabilities_Append", func(b *testing.B) {
		for range b.N {
			_ = append([]pbc.PluginCapability(nil), capabilities...)
		}
	})
}
