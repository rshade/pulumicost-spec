package registry_test

import (
	"testing"

	"github.com/rshade/pulumicost-spec/sdk/go/registry"
)

func TestDiscoverySource(t *testing.T) {
	tests := []struct {
		source   string
		expected bool
	}{
		{"filesystem", true},
		{"registry", true},
		{"url", true},
		{"git", true},
		{"invalid", false},
		{"", false},
	}

	for _, test := range tests {
		result := registry.IsValidDiscoverySource(test.source)
		if result != test.expected {
			t.Errorf("registry.IsValidDiscoverySource(%q) = %v, expected %v", test.source, result, test.expected)
		}
	}
}

func TestPluginStatus(t *testing.T) {
	tests := []struct {
		status   string
		expected bool
	}{
		{"available", true},
		{"installed", true},
		{"active", true},
		{"inactive", true},
		{"error", true},
		{"updating", true},
		{"invalid", false},
		{"", false},
	}

	for _, test := range tests {
		result := registry.IsValidPluginStatus(test.status)
		if result != test.expected {
			t.Errorf("registry.IsValidPluginStatus(%q) = %v, expected %v", test.status, result, test.expected)
		}
	}
}

func TestSecurityLevel(t *testing.T) {
	tests := []struct {
		level    string
		expected bool
	}{
		{"untrusted", true},
		{"community", true},
		{"verified", true},
		{"official", true},
		{"invalid", false},
		{"", false},
	}

	for _, test := range tests {
		result := registry.IsValidSecurityLevel(test.level)
		if result != test.expected {
			t.Errorf("registry.IsValidSecurityLevel(%q) = %v, expected %v", test.level, result, test.expected)
		}
	}
}

func TestInstallationMethod(t *testing.T) {
	tests := []struct {
		method   string
		expected bool
	}{
		{"binary", true},
		{"container", true},
		{"script", true},
		{"package", true},
		{"invalid", false},
		{"", false},
	}

	for _, test := range tests {
		result := registry.IsValidInstallationMethod(test.method)
		if result != test.expected {
			t.Errorf("registry.IsValidInstallationMethod(%q) = %v, expected %v", test.method, result, test.expected)
		}
	}
}

func TestPluginCapability(t *testing.T) {
	tests := []struct {
		capability string
		expected   bool
	}{
		{"cost_retrieval", true},
		{"cost_projection", true},
		{"pricing_specs", true},
		{"historical_data", true},
		{"real_time_data", true},
		{"batch_processing", true},
		{"rate_limiting", true},
		{"caching", true},
		{"encryption", true},
		{"compression", true},
		{"filtering", true},
		{"aggregation", true},
		{"multi_tenancy", true},
		{"audit_logging", true},
		{"invalid", false},
		{"", false},
	}

	for _, test := range tests {
		result := registry.IsValidPluginCapability(test.capability)
		if result != test.expected {
			t.Errorf("registry.IsValidPluginCapability(%q) = %v, expected %v", test.capability, result, test.expected)
		}
	}
}

func TestSystemPermission(t *testing.T) {
	tests := []struct {
		permission string
		expected   bool
	}{
		{"network_access", true},
		{"filesystem_read", true},
		{"filesystem_write", true},
		{"environment_read", true},
		{"process_spawn", true},
		{"system_info", true},
		{"temp_files", true},
		{"config_read", true},
		{"metrics_collect", true},
		{"invalid", false},
		{"", false},
	}

	for _, test := range tests {
		result := registry.IsValidSystemPermission(test.permission)
		if result != test.expected {
			t.Errorf("registry.IsValidSystemPermission(%q) = %v, expected %v", test.permission, result, test.expected)
		}
	}
}

func TestAuthMethod(t *testing.T) {
	tests := []struct {
		method   string
		expected bool
	}{
		{"none", true},
		{"api_key", true},
		{"jwt", true},
		{"oauth2", true},
		{"mtls", true},
		{"basic_auth", true},
		{"invalid", false},
		{"", false},
	}

	for _, test := range tests {
		result := registry.IsValidAuthMethod(test.method)
		if result != test.expected {
			t.Errorf("registry.IsValidAuthMethod(%q) = %v, expected %v", test.method, result, test.expected)
		}
	}
}

func TestValidatePluginName(t *testing.T) {
	tests := []struct {
		name        string
		expectError bool
	}{
		{"valid-plugin", false},
		{"plugin123", false},
		{"ab", false},
		{"plugin-name-123", false},
		{"", true},
		{"a", true},              // single character not allowed (min length 2)
		{"x", true},              // single character not allowed (min length 2)
		{"invalid_plugin", true}, // underscore not allowed
		{"Invalid-Plugin", true}, // uppercase not allowed
		{"-invalid", true},       // cannot start with dash
		{"invalid-", true},       // cannot end with dash
		{"plugin--name", false},  // double dash allowed
		{"plugin name", true},    // space not allowed
		{"plugin.name", true},    // dot not allowed
		{"this-is-a-very-long-plugin-name-that-exceeds-fifty-characters", true}, // too long
	}

	for _, test := range tests {
		err := registry.ValidatePluginName(test.name)
		hasError := err != nil
		if hasError != test.expectError {
			if test.expectError {
				t.Errorf("registry.ValidatePluginName(%q) expected error, but got none", test.name)
			} else {
				t.Errorf("registry.ValidatePluginName(%q) unexpected error: %v", test.name, err)
			}
		}
	}
}

func TestAllFunctions(t *testing.T) {
	// Test that all enumeration functions return expected counts
	tests := []struct {
		name     string
		count    int
		function func() int
	}{
		{"registry.AllDiscoverySources", 4, func() int { return len(registry.AllDiscoverySources()) }},
		{"registry.AllPluginStatuses", 6, func() int { return len(registry.AllPluginStatuses()) }},
		{"registry.AllSecurityLevels", 4, func() int { return len(registry.AllSecurityLevels()) }},
		{"registry.AllInstallationMethods", 4, func() int { return len(registry.AllInstallationMethods()) }},
		{"registry.AllPluginCapabilities", 14, func() int { return len(registry.AllPluginCapabilities()) }},
		{"registry.AllSystemPermissions", 9, func() int { return len(registry.AllSystemPermissions()) }},
		{"registry.AllAuthMethods", 6, func() int { return len(registry.AllAuthMethods()) }},
	}

	for _, test := range tests {
		count := test.function()
		if count != test.count {
			t.Errorf("%s returned %d items, expected %d", test.name, count, test.count)
		}
	}
}
