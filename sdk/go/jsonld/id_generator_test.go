package jsonld_test

import (
	"strings"
	"testing"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/rshade/finfocus-spec/sdk/go/jsonld"
	pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

func TestNewIDGenerator(t *testing.T) {
	gen := jsonld.NewIDGenerator()

	if gen == nil {
		t.Fatal("NewIDGenerator() returned nil")
	}

	// Verify it generates IDs with the expected prefix
	record := &pbc.FocusCostRecord{
		BillingAccountId: "123456789012",
		ChargePeriodStart: &timestamppb.Timestamp{
			Seconds: 1735689600,
		},
	}
	id := gen.Generate(record)
	if !strings.HasPrefix(id, "urn:focus:cost:") {
		t.Errorf("Expected ID to start with 'urn:focus:cost:', got '%s'", id)
	}
}

func TestIDGenerator_Generate_CompositeKey(t *testing.T) {
	gen := jsonld.NewIDGenerator()

	record := &pbc.FocusCostRecord{
		BillingAccountId: "123456789012",
		ChargePeriodStart: &timestamppb.Timestamp{
			Seconds: 1735689600, // 2025-01-01
		},
		ResourceId: "i-1234567890abcdef0",
	}

	id := gen.Generate(record)

	if id == "" {
		t.Fatal("Generate() returned empty string")
	}

	if !strings.HasPrefix(id, "urn:focus:cost:") {
		t.Errorf("Expected ID to start with 'urn:focus:cost:', got '%s'", id)
	}

	// ID should be deterministic for same record
	id2 := gen.Generate(record)
	if id != id2 {
		t.Error("Generate() returned different IDs for same record")
	}
}

func TestIDGenerator_Generate_DifferentRecords(t *testing.T) {
	gen := jsonld.NewIDGenerator()

	record1 := &pbc.FocusCostRecord{
		BillingAccountId: "123456789012",
		ChargePeriodStart: &timestamppb.Timestamp{
			Seconds: 1735689600,
		},
		ResourceId: "i-1234567890abcdef0",
	}

	record2 := &pbc.FocusCostRecord{
		BillingAccountId: "123456789012",
		ChargePeriodStart: &timestamppb.Timestamp{
			Seconds: 1735689600,
		},
		ResourceId: "i-abcdef1234567890", // Different resource
	}

	id1 := gen.Generate(record1)
	id2 := gen.Generate(record2)

	if id1 == id2 {
		t.Error("Generate() returned same ID for different records")
	}
}

func TestIDGenerator_GenerateCommitment(t *testing.T) {
	gen := jsonld.NewIDGenerator()

	record := &pbc.ContractCommitment{
		ContractCommitmentId: "cc-123456",
	}

	id := gen.GenerateCommitment(record)

	if id == "" {
		t.Fatal("GenerateCommitment() returned empty string")
	}

	if !strings.HasPrefix(id, "urn:focus:commitment:") {
		t.Errorf("Expected ID to start with 'urn:focus:commitment:', got '%s'", id)
	}
}

func TestIDGenerator_GenerateCommitment_Deterministic(t *testing.T) {
	gen := jsonld.NewIDGenerator()

	record := &pbc.ContractCommitment{
		ContractCommitmentId: "cc-123456",
	}

	id1 := gen.GenerateCommitment(record)
	id2 := gen.GenerateCommitment(record)

	if id1 != id2 {
		t.Error("GenerateCommitment() returned different IDs for same record")
	}
}

func TestIDGenerator_GenerateCommitment_EmptyID(t *testing.T) {
	gen := jsonld.NewIDGenerator()

	// Test empty commitment ID produces distinctive ID (not a hash collision)
	record := &pbc.ContractCommitment{
		ContractCommitmentId: "",
	}

	id := gen.GenerateCommitment(record)

	expected := "urn:focus:commitment:empty-commitment-id"
	if id != expected {
		t.Errorf("Expected %q for empty commitment ID, got %q", expected, id)
	}
}

func TestIDGenerator_GenerateCommitment_NilRecord(t *testing.T) {
	gen := jsonld.NewIDGenerator()

	id := gen.GenerateCommitment(nil)

	expected := "urn:focus:commitment:nil-commitment"
	if id != expected {
		t.Errorf("Expected %q for nil record, got %q", expected, id)
	}
}

func TestIDGenerator_WithUserIDField(t *testing.T) {
	// Test via Serializer option which applies the user ID field
	serializer := jsonld.NewSerializer(
		jsonld.WithUserIDField("invoice_id"),
	)

	record := &pbc.FocusCostRecord{
		BillingAccountId: "123456789012",
		InvoiceId:        "INV-2025-001",
		ChargePeriodStart: &timestamppb.Timestamp{
			Seconds: 1735689600,
		},
		ResourceId:      "i-1234567890abcdef0",
		BilledCost:      100.0,
		BillingCurrency: "USD",
	}

	output, err := serializer.Serialize(record)
	if err != nil {
		t.Fatalf("Serialize() failed: %v", err)
	}

	// Check that the output contains the invoice ID in the @id
	if !strings.Contains(string(output), "INV-2025-001") {
		t.Errorf("Expected @id to contain invoice_id, output: %s", string(output))
	}
}

func TestIDGenerator_CopyOnWrite(t *testing.T) {
	// Verify that With* methods don't mutate the original generator
	gen1 := jsonld.NewIDGenerator()

	record := &pbc.FocusCostRecord{
		BillingAccountId: "123456789012",
		InvoiceId:        "INV-2025-001",
		ChargePeriodStart: &timestamppb.Timestamp{
			Seconds: 1735689600,
		},
		ResourceId: "i-1234567890abcdef0",
	}

	// Get baseline ID from original generator
	id1 := gen1.Generate(record)

	// Cast to ConfigurableIDGenerator to access With* methods
	configGen, ok := gen1.(jsonld.ConfigurableIDGenerator)
	if !ok {
		t.Fatal("Expected NewIDGenerator to return ConfigurableIDGenerator")
	}

	// Create a new generator with different settings
	gen2 := configGen.WithUserIDField("invoice_id")

	// Original generator should still produce the same ID (not mutated)
	id1Again := gen1.Generate(record)
	if id1 != id1Again {
		t.Errorf("Original generator was mutated: expected %s, got %s", id1, id1Again)
	}

	// New generator should use invoice_id
	id2 := gen2.Generate(record)
	if !strings.Contains(id2, "INV-2025-001") {
		t.Errorf("New generator should use invoice_id, got: %s", id2)
	}

	// They should be different
	if id1 == id2 {
		t.Error("Original and modified generators should produce different IDs")
	}
}
