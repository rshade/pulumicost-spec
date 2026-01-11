// Copyright 2026 PulumiCost/FinFocus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
