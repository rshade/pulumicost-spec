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

package pluginsdk_test

import (
	"testing"

	"github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
)

// TestTestPluginIntegration tests the TestPlugin and TestServer utilities.
func TestTestPluginIntegration(t *testing.T) {
	// Create a plugin using BasePlugin
	plugin := pluginsdk.NewBasePlugin("integration-test-plugin")
	plugin.Matcher().AddProvider("aws")
	plugin.Matcher().AddResourceType("aws:ec2:Instance")

	// Create test plugin which starts a real gRPC server
	tp := pluginsdk.NewTestPlugin(t, plugin)

	// Test the Name method
	tp.TestName("integration-test-plugin")
}

// TestTestServerClose tests the TestServer Close method.
func TestTestServerClose(t *testing.T) {
	plugin := pluginsdk.NewBasePlugin("close-test-plugin")

	// Create test server directly
	ts := pluginsdk.NewTestServer(t, plugin)

	// Get client and verify it works
	client := ts.Client()
	if client == nil {
		t.Fatal("Expected non-nil client")
	}

	// Close should not panic or error
	ts.Close()
}

// TestTestServerCloseNil tests that Close handles nil gracefully.
func TestTestServerCloseNil(_ *testing.T) {
	var ts *pluginsdk.TestServer
	// This should not panic
	ts.Close()
}

// TestTestPluginProjectedCost tests the TestProjectedCost helper.
func TestTestPluginProjectedCost(t *testing.T) {
	plugin := pluginsdk.NewBasePlugin("cost-test-plugin")

	tp := pluginsdk.NewTestPlugin(t, plugin)

	// Test projected cost - should return error since BasePlugin returns "not supported"
	resource := pluginsdk.CreateTestResource("aws", "aws:ec2:Instance", nil)
	tp.TestProjectedCost(resource, true) // expect error
}

// TestTestPluginActualCost tests the TestActualCost helper.
func TestTestPluginActualCost(t *testing.T) {
	plugin := pluginsdk.NewBasePlugin("actual-cost-test-plugin")

	tp := pluginsdk.NewTestPlugin(t, plugin)

	// Test actual cost - should return error since BasePlugin returns "no data"
	tp.TestActualCost("test-resource-id", 0, 3600, true) // expect error
}
