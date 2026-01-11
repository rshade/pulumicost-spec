package testing //nolint:testpackage // White-box testing required

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetPluginInfoPerformance_LegacyPlugin(t *testing.T) {
	// Create legacy plugin (does not implement PluginInfoProvider)
	plugin := NewMockLegacyPlugin("legacy-test-plugin")

	// Create test harness
	harness := NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	// Run GetPluginInfo performance test
	testFunc := createGetPluginInfoLatencyTest()
	result := testFunc(harness)

	// Should pass gracefully (Unimplemented error is acceptable for legacy plugins)
	assert.True(t, result.Success, "Legacy plugin performance test should pass")
	assert.Less(t, result.Duration, 100*time.Millisecond, "Should complete quickly even with Unimplemented error")
}
