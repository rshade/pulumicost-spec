package jsonld_test

import (
	"encoding/json"
	"testing"

	"github.com/rshade/pulumicost-spec/sdk/go/jsonld"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestSchemaOrg_MonetaryAmountType(t *testing.T) {
	serializer := jsonld.NewSerializer()

	record := &pbc.FocusCostRecord{
		BillingAccountId: "123456789012",
		ChargePeriodStart: &timestamppb.Timestamp{
			Seconds: 1735689600,
		},
		BilledCost:      125.50,
		ListCost:        150.00,
		EffectiveCost:   125.50,
		ContractedCost:  120.00,
		BillingCurrency: "USD",
	}

	output, err := serializer.Serialize(record)
	if err != nil {
		t.Fatalf("Serialize() failed: %v", err)
	}

	var result map[string]interface{}
	if unmarshalErr := json.Unmarshal(output, &result); unmarshalErr != nil {
		t.Fatalf("Output is not valid JSON: %v", unmarshalErr)
	}

	// Check billedCost is serialized as MonetaryAmount
	costFields := []string{"billedCost", "listCost", "effectiveCost", "contractedCost"}
	for _, field := range costFields {
		costValue, ok := result[field]
		if !ok {
			t.Errorf("Output missing %s field", field)
			continue
		}

		costMap, isMap := costValue.(map[string]interface{})
		if !isMap {
			t.Errorf("%s should be a map (MonetaryAmount), got %T", field, costValue)
			continue
		}

		typeVal, hasType := costMap["@type"]
		if !hasType {
			t.Errorf("%s missing @type annotation", field)
			continue
		}
		if typeVal != "schema:MonetaryAmount" {
			t.Errorf("%s @type = %v, want schema:MonetaryAmount", field, typeVal)
		}

		if _, hasValue := costMap["value"]; !hasValue {
			t.Errorf("%s missing value field", field)
		}

		if _, hasCurrency := costMap["currency"]; !hasCurrency {
			t.Errorf("%s missing currency field", field)
		}
	}
}

func TestSchemaOrg_DateTimeFormatting(t *testing.T) {
	serializer := jsonld.NewSerializer()

	record := &pbc.FocusCostRecord{
		BillingAccountId: "123456789012",
		ChargePeriodStart: &timestamppb.Timestamp{
			Seconds: 1735689600, // 2025-01-01T00:00:00Z
		},
		ChargePeriodEnd: &timestamppb.Timestamp{
			Seconds: 1735776000, // 2025-01-02T00:00:00Z
		},
		BillingPeriodStart: &timestamppb.Timestamp{
			Seconds: 1735689600,
		},
		BillingPeriodEnd: &timestamppb.Timestamp{
			Seconds: 1738368000, // 2025-02-01T00:00:00Z
		},
		BilledCost:      100.0,
		BillingCurrency: "USD",
	}

	output, err := serializer.Serialize(record)
	if err != nil {
		t.Fatalf("Serialize() failed: %v", err)
	}

	var result map[string]interface{}
	if unmarshalErr := json.Unmarshal(output, &result); unmarshalErr != nil {
		t.Fatalf("Output is not valid JSON: %v", unmarshalErr)
	}

	// Check timestamp fields are ISO 8601 formatted
	timestampFields := []struct {
		name     string
		expected string
	}{
		{"chargePeriodStart", "2025-01-01T00:00:00Z"},
		{"chargePeriodEnd", "2025-01-02T00:00:00Z"},
		{"billingPeriodStart", "2025-01-01T00:00:00Z"},
		{"billingPeriodEnd", "2025-02-01T00:00:00Z"},
	}

	for _, tc := range timestampFields {
		value, ok := result[tc.name]
		if !ok {
			t.Errorf("Output missing %s field", tc.name)
			continue
		}

		strValue, isString := value.(string)
		if !isString {
			t.Errorf("%s should be a string, got %T", tc.name, value)
			continue
		}

		if strValue != tc.expected {
			t.Errorf("%s = %s, want %s", tc.name, strValue, tc.expected)
		}
	}
}

func TestSchemaOrg_FocusNamespaceFallback(t *testing.T) {
	serializer := jsonld.NewSerializer()

	record := &pbc.FocusCostRecord{
		BillingAccountId: "123456789012",
		ChargePeriodStart: &timestamppb.Timestamp{
			Seconds: 1735689600,
		},
		// These fields use FOCUS namespace (no Schema.org equivalent)
		SkuId:                 "sku-12345",
		CommitmentDiscountId:  "cdi-12345",
		CapacityReservationId: "cap-12345",
		AllocatedMethodId:     "alloc-method-001",
		AllocatedResourceId:   "alloc-resource-001",
		BilledCost:            100.0,
		BillingCurrency:       "USD",
	}

	output, err := serializer.Serialize(record)
	if err != nil {
		t.Fatalf("Serialize() failed: %v", err)
	}

	var result map[string]interface{}
	if unmarshalErr := json.Unmarshal(output, &result); unmarshalErr != nil {
		t.Fatalf("Output is not valid JSON: %v", unmarshalErr)
	}

	// Check that FOCUS-only fields are present in output
	focusOnlyFields := []string{
		"skuId",
		"commitmentDiscountId",
		"capacityReservationId",
		"allocatedMethodId",
		"allocatedResourceId",
	}

	for _, field := range focusOnlyFields {
		if _, ok := result[field]; !ok {
			t.Errorf("Output missing %s field (should use FOCUS namespace)", field)
		}
	}

	// Verify context includes FOCUS namespace
	context, ok := result["@context"].(map[string]interface{})
	if !ok {
		t.Error("@context should be a map")
		return
	}

	focusNS, hasFocus := context["focus"]
	if !hasFocus {
		t.Error("Context missing 'focus' namespace prefix")
		return
	}

	if focusNS != "https://focus.finops.org/v1#" {
		t.Errorf("FOCUS namespace = %v, want https://focus.finops.org/v1#", focusNS)
	}
}

func TestSchemaOrg_IsSchemaMapped(t *testing.T) {
	tests := []struct {
		field    string
		expected bool
	}{
		{"billed_cost", true},
		{"list_cost", true},
		{"effective_cost", true},
		{"contracted_cost", true},
		{"charge_period_start", true},
		{"charge_period_end", true},
		{"billing_period_start", true},
		{"billing_period_end", true},
		{"service_name", true},
		{"resource_name", true},
		{"region_name", true},
		// FOCUS-only fields
		{"sku_id", false},
		{"billing_account_id", false},
		{"commitment_discount_id", false},
		{"allocated_method_id", false},
	}

	for _, tc := range tests {
		t.Run(tc.field, func(t *testing.T) {
			got := jsonld.IsSchemaMapped(tc.field)
			if got != tc.expected {
				t.Errorf("IsSchemaMapped(%q) = %v, want %v", tc.field, got, tc.expected)
			}
		})
	}
}

func TestSchemaOrg_StandardPrefixes(t *testing.T) {
	prefixes := jsonld.StandardPrefixes()

	expectedPrefixes := map[string]string{
		"schema": "https://schema.org/",
		"focus":  "https://focus.finops.org/v1#",
		"xsd":    "http://www.w3.org/2001/XMLSchema#",
	}

	for prefix, expected := range expectedPrefixes {
		actual, ok := prefixes[prefix]
		if !ok {
			t.Errorf("StandardPrefixes() missing %q", prefix)
			continue
		}
		if actual != expected {
			t.Errorf("StandardPrefixes()[%q] = %q, want %q", prefix, actual, expected)
		}
	}
}

func TestConformance_JSONLD11Spec(t *testing.T) {
	serializer := jsonld.NewSerializer()

	record := &pbc.FocusCostRecord{
		BillingAccountId: "123456789012",
		ChargePeriodStart: &timestamppb.Timestamp{
			Seconds: 1735689600,
		},
		ServiceName:     "Amazon EC2",
		BilledCost:      125.50,
		BillingCurrency: "USD",
	}

	output, err := serializer.Serialize(record)
	if err != nil {
		t.Fatalf("Serialize() failed: %v", err)
	}

	var result map[string]interface{}
	if unmarshalErr := json.Unmarshal(output, &result); unmarshalErr != nil {
		t.Fatalf("Output is not valid JSON: %v", unmarshalErr)
	}

	// JSON-LD 1.1 requirements:
	// 1. @context MUST be present
	if _, ok := result["@context"]; !ok {
		t.Error("JSON-LD 1.1: @context MUST be present")
	}

	// 2. @type SHOULD be present for typed nodes
	if _, ok := result["@type"]; !ok {
		t.Error("JSON-LD 1.1: @type SHOULD be present for typed nodes")
	}

	// 3. @id SHOULD be present for identifiable resources
	if _, ok := result["@id"]; !ok {
		t.Error("JSON-LD 1.1: @id SHOULD be present for identifiable resources")
	}

	// 4. Context must include valid namespace prefixes
	context, contextOk := result["@context"].(map[string]interface{})
	if !contextOk {
		t.Fatal("JSON-LD 1.1: @context must be a map[string]interface{}")
	}
	if _, hasSchema := context["schema"]; !hasSchema {
		t.Error("JSON-LD 1.1: context should include 'schema' prefix for Schema.org")
	}
	if _, hasFocus := context["focus"]; !hasFocus {
		t.Error("JSON-LD 1.1: context should include 'focus' prefix for FOCUS vocabulary")
	}
}
