package testing_test

import (
	"testing"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
	"google.golang.org/protobuf/proto"
)

func TestGetActualCostRequest_ArnField(t *testing.T) {
	req := &pbc.GetActualCostRequest{}

	// Verify that the Arn field exists and is of type string
	if _, ok := interface{}(req).(interface{ GetArn() string }); !ok {
		t.Errorf("GetActualCostRequest does not have a GetArn() method returning string")
	}

	// Verify setting and getting the Arn field
	expectedArn := "arn:aws:ec2:us-east-1:123456789012:instance/i-1234567890abcdef0"
	req.Arn = expectedArn
	if req.GetArn() != expectedArn {
		t.Errorf("Arn field not set or retrieved correctly. Expected %s, got %s", expectedArn, req.GetArn())
	}
}

func TestGetActualCostRequest_BackwardCompatibility(t *testing.T) {
	// Simulate an old client sending a request without the Arn field
	oldReq := &pbc.GetActualCostRequest{
		ResourceId: "i-123",
		Tags:       map[string]string{"env": "dev"},
	}

	// Marshal and Unmarshal to simulate network transmission
	marshaledData, err := proto.Marshal(oldReq)
	if err != nil {
		t.Fatalf("Failed to marshal old request: %v", err)
	}

	newReq := &pbc.GetActualCostRequest{}
	err = proto.Unmarshal(marshaledData, newReq)
	if err != nil {
		t.Fatalf("Failed to unmarshal data into new request: %v", err)
	}

	// Verify that the new request correctly parses the old data and Arn is empty
	if newReq.GetResourceId() != oldReq.GetResourceId() {
		t.Errorf(
			"ResourceId mismatch. Expected %s, got %s",
			oldReq.GetResourceId(),
			newReq.GetResourceId(),
		)
	}
	if len(newReq.GetTags()) != len(oldReq.GetTags()) {
		t.Errorf(
			"Tags length mismatch. Expected %d, got %d",
			len(oldReq.GetTags()),
			len(newReq.GetTags()),
		)
	}
	if newReq.GetArn() != "" {
		t.Errorf("Arn field should be empty for old request. Got %s", newReq.GetArn())
	}

	// Simulate a new client sending a request with the Arn field
	newReqWithArn := &pbc.GetActualCostRequest{
		ResourceId: "i-456",
		Arn:        "arn:aws:ec2:us-west-2:987654321098:instance/i-newarn",
	}

	marshaledDataWithArn, err := proto.Marshal(newReqWithArn)
	if err != nil {
		t.Fatalf("Failed to marshal new request with Arn: %v", err)
	}

	unmarshaledNewReq := &pbc.GetActualCostRequest{}
	err = proto.Unmarshal(marshaledDataWithArn, unmarshaledNewReq)
	if err != nil {
		t.Fatalf("Failed to unmarshal new request with Arn: %v", err)
	}

	if unmarshaledNewReq.GetArn() != newReqWithArn.GetArn() {
		t.Errorf(
			"Arn field mismatch after round-trip. Expected %s, got %s",
			newReqWithArn.GetArn(),
			unmarshaledNewReq.GetArn(),
		)
	}
}
