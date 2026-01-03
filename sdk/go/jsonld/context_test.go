package jsonld_test

import (
	"testing"

	"github.com/rshade/pulumicost-spec/sdk/go/jsonld"
)

func TestNewContext(t *testing.T) {
	ctx := jsonld.NewContext()

	if ctx == nil {
		t.Fatal("NewContext() returned nil")
	}

	// Verify default context includes expected keys (no remote = map)
	result, ok := ctx.Build().(map[string]interface{})
	if !ok {
		t.Fatal("Expected Build() to return map when no remote contexts")
	}
	_, hasSchema := result["schema"]
	if !hasSchema {
		t.Error("Expected schema to be present in default context")
	}
	if result["focus"] != "https://focus.finops.org/v1#" {
		t.Errorf("Expected focus namespace 'https://focus.finops.org/v1#', got '%v'", result["focus"])
	}
}

func TestContextWithCustomMapping(t *testing.T) {
	ctx := jsonld.NewContext().
		WithCustomMapping("billingAccountId", "yourOrg:accountIdentifier")

	result, ok := ctx.Build().(map[string]interface{})
	if !ok {
		t.Fatal("Expected Build() to return map when no remote contexts")
	}
	if result["billingAccountId"] != "yourOrg:accountIdentifier" {
		t.Errorf("Custom mapping not set correctly, got: %v", result["billingAccountId"])
	}
}

func TestContextWithRemoteContext(t *testing.T) {
	url := "https://your-org.com/ontology/v1"
	ctx := jsonld.NewContext().
		WithRemoteContext(url)

	// Per JSON-LD 1.1: remote context returns array [url, inline-object]
	result, ok := ctx.Build().([]interface{})
	if !ok {
		t.Fatal("Expected Build() to return array when remote contexts exist")
	}

	if len(result) != 2 {
		t.Fatalf("Expected array length 2, got %d", len(result))
	}

	// First element should be the remote URL string
	if result[0] != url {
		t.Errorf("Expected first element to be remote URL %q, got %v", url, result[0])
	}

	// Second element should be inline context object
	inline, inlineOK := result[1].(map[string]interface{})
	if !inlineOK {
		t.Fatal("Expected second element to be inline context object")
	}

	// Inline should have schema, focus, xsd
	_, hasSchema := inline["schema"]
	if !hasSchema {
		t.Error("Inline context missing 'schema'")
	}
}

func TestContextBuild_DefaultContext(t *testing.T) {
	ctx := jsonld.NewContext()
	result, ok := ctx.Build().(map[string]interface{})
	if !ok {
		t.Fatal("Expected Build() to return map when no remote contexts")
	}

	wantKeys := []string{"schema", "focus", "xsd"}
	for _, key := range wantKeys {
		if _, exists := result[key]; !exists {
			t.Errorf("Context missing expected key: %s", key)
		}
	}
}

func TestContextBuild_WithRemoteContext(t *testing.T) {
	ctx := jsonld.NewContext().WithRemoteContext("https://example.com/context.json")
	result, ok := ctx.Build().([]interface{})
	if !ok {
		t.Fatal("Expected Build() to return array when remote contexts exist")
	}

	if len(result) != 2 {
		t.Fatalf("Expected array length 2, got %d", len(result))
	}

	if result[0] != "https://example.com/context.json" {
		t.Errorf("Expected first element to be remote URL, got %v", result[0])
	}

	inline, ok := result[1].(map[string]interface{})
	if !ok {
		t.Fatal("Expected last element to be inline context map")
	}

	wantKeys := []string{"schema", "focus", "xsd"}
	for _, key := range wantKeys {
		if _, exists := inline[key]; !exists {
			t.Errorf("Inline context missing expected key: %s", key)
		}
	}
}

func TestContextBuild_WithCustomMapping(t *testing.T) {
	ctx := jsonld.NewContext().WithCustomMapping("customField", "custom:iri")
	result, ok := ctx.Build().(map[string]interface{})
	if !ok {
		t.Fatal("Expected Build() to return map when no remote contexts")
	}

	wantKeys := []string{"schema", "focus", "xsd", "customField"}
	for _, key := range wantKeys {
		if _, exists := result[key]; !exists {
			t.Errorf("Context missing expected key: %s", key)
		}
	}
}

func TestContextBuild_WithMultipleRemoteContexts(t *testing.T) {
	ctx := jsonld.NewContext().
		WithRemoteContext("https://example.com/context1.json").
		WithRemoteContext("https://example.com/context2.json")
	result, ok := ctx.Build().([]interface{})
	if !ok {
		t.Fatal("Expected Build() to return array when remote contexts exist")
	}

	if len(result) != 3 {
		t.Fatalf("Expected array length 3, got %d", len(result))
	}

	// Verify remote URLs in order
	if result[0] != "https://example.com/context1.json" {
		t.Errorf("First remote URL mismatch, got %v", result[0])
	}
	if result[1] != "https://example.com/context2.json" {
		t.Errorf("Second remote URL mismatch, got %v", result[1])
	}

	// Last element is inline context
	inline, ok := result[2].(map[string]interface{})
	if !ok {
		t.Fatal("Expected last element to be inline context map")
	}

	wantKeys := []string{"schema", "focus", "xsd"}
	for _, key := range wantKeys {
		if _, exists := inline[key]; !exists {
			t.Errorf("Inline context missing expected key: %s", key)
		}
	}
}

func TestContextBuild_SchemaOrgValue(t *testing.T) {
	ctx := jsonld.NewContext()
	result, ok := ctx.Build().(map[string]interface{})
	if !ok {
		t.Fatal("Expected Build() to return map when no remote contexts")
	}

	schemaValue, ok := result["schema"]
	if !ok {
		t.Fatal("Context missing 'schema' key")
	}

	if schemaValue != "https://schema.org/" {
		t.Errorf("Expected schema namespace 'https://schema.org/', got '%v'", schemaValue)
	}
}

func TestContextBuild_FocusNamespaceValue(t *testing.T) {
	ctx := jsonld.NewContext()
	result, ok := ctx.Build().(map[string]interface{})
	if !ok {
		t.Fatal("Expected Build() to return map when no remote contexts")
	}

	focusValue, ok := result["focus"]
	if !ok {
		t.Fatal("Context missing 'focus' key")
	}

	if focusValue != "https://focus.finops.org/v1#" {
		t.Errorf("Expected focus namespace 'https://focus.finops.org/v1#', got '%v'", focusValue)
	}
}

func TestContextBuild_Chaining(t *testing.T) {
	ctx := jsonld.NewContext().
		WithCustomMapping("field1", "iri1").
		WithCustomMapping("field2", "iri2").
		WithRemoteContext("https://example.com/context.json").
		WithRemoteContext("https://example.com/context2.json")

	// With remote contexts, Build() returns array
	result, ok := ctx.Build().([]interface{})
	if !ok {
		t.Fatal("Expected Build() to return array when remote contexts exist")
	}

	// Should be [url1, url2, inline-object]
	if len(result) != 3 {
		t.Fatalf("Expected array length 3, got %d", len(result))
	}

	// Verify remote contexts in order
	if result[0] != "https://example.com/context.json" {
		t.Error("Remote context 1 not in correct position")
	}
	if result[1] != "https://example.com/context2.json" {
		t.Error("Remote context 2 not in correct position")
	}

	// Verify custom mappings in inline context
	inline, ok := result[2].(map[string]interface{})
	if !ok {
		t.Fatal("Expected last element to be inline context map")
	}
	if inline["field1"] != "iri1" {
		t.Error("Custom mapping field1 not preserved")
	}
	if inline["field2"] != "iri2" {
		t.Error("Custom mapping field2 not preserved")
	}
}

func TestContext_Validation(t *testing.T) {
	tests := []struct {
		name        string
		ctx         *jsonld.Context
		expectError bool
	}{
		{
			name:        "valid default context",
			ctx:         jsonld.NewContext(),
			expectError: false,
		},
		{
			name: "valid with remote context",
			ctx: jsonld.NewContext().
				WithRemoteContext("https://example.com/context.json"),
			expectError: false,
		},
		{
			name: "valid with custom mapping",
			ctx: jsonld.NewContext().
				WithCustomMapping("field", "https://example.com/vocab#field"),
			expectError: false,
		},
		{
			name: "invalid remote context URL - not a URL",
			ctx: jsonld.NewContext().
				WithRemoteContext("not-a-valid-url"),
			expectError: true,
		},
		{
			name: "invalid remote context URL - ftp scheme",
			ctx: jsonld.NewContext().
				WithRemoteContext("ftp://example.com/context.json"),
			expectError: true,
		},
		{
			name: "invalid remote context URL - relative path",
			ctx: jsonld.NewContext().
				WithRemoteContext("/relative/path"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ctx.Validate()
			if tt.expectError && err == nil {
				t.Error("Expected validation error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected validation error: %v", err)
			}
		})
	}
}

func TestContext_OverrideExistingMapping(t *testing.T) {
	// Custom mapping should override default namespace mappings
	ctx := jsonld.NewContext().
		WithCustomMapping("billingAccountId", "myOrg:customAccountID")

	result, ok := ctx.Build().(map[string]interface{})
	if !ok {
		t.Fatal("Expected Build() to return map when no remote contexts")
	}

	if result["billingAccountId"] != "myOrg:customAccountID" {
		t.Errorf("Custom mapping did not override, got: %v", result["billingAccountId"])
	}
}

func TestContext_CopyOnWrite_CustomMapping(t *testing.T) {
	// Verify that WithCustomMapping doesn't mutate the original context
	ctx1 := jsonld.NewContext()

	// Get baseline result from original context
	result1, ok := ctx1.Build().(map[string]interface{})
	if !ok {
		t.Fatal("Expected Build() to return map")
	}
	_, hasCustom := result1["customField"]
	if hasCustom {
		t.Fatal("Original context should not have customField")
	}

	// Create a new context with custom mapping
	ctx2 := ctx1.WithCustomMapping("customField", "custom:iri")

	// Original context should still NOT have the custom field (not mutated)
	result1Again, ok := ctx1.Build().(map[string]interface{})
	if !ok {
		t.Fatal("Expected Build() to return map")
	}
	_, hasCustomAfter := result1Again["customField"]
	if hasCustomAfter {
		t.Error("Original context was mutated by WithCustomMapping")
	}

	// New context SHOULD have the custom field
	result2, ok := ctx2.Build().(map[string]interface{})
	if !ok {
		t.Fatal("Expected Build() to return map")
	}
	if result2["customField"] != "custom:iri" {
		t.Errorf("New context should have customField, got: %v", result2["customField"])
	}
}

func TestContext_CopyOnWrite_RemoteContext(t *testing.T) {
	// Verify that WithRemoteContext doesn't mutate the original context
	ctx1 := jsonld.NewContext()

	// Original context should return a map (no remote contexts)
	result1, ok := ctx1.Build().(map[string]interface{})
	if !ok {
		t.Fatal("Expected Build() to return map when no remote contexts")
	}
	_ = result1

	// Create a new context with a remote context
	ctx2 := ctx1.WithRemoteContext("https://example.com/context.json")

	// Original context should still return a map (not mutated)
	result1Again, ok := ctx1.Build().(map[string]interface{})
	if !ok {
		t.Error("Original context was mutated by WithRemoteContext - now returns array")
	}
	_ = result1Again

	// New context should return an array
	result2, ok := ctx2.Build().([]interface{})
	if !ok {
		t.Fatal("Expected new context Build() to return array")
	}
	if len(result2) != 2 {
		t.Errorf("Expected array length 2, got %d", len(result2))
	}
	if result2[0] != "https://example.com/context.json" {
		t.Errorf("Expected first element to be remote URL, got %v", result2[0])
	}
}
