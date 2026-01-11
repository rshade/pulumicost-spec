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

package jsonld_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/rshade/pulumicost-spec/sdk/go/jsonld"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestSerializeStream_BasicBatch(t *testing.T) {
	serializer := jsonld.NewSerializer()

	records := make([]*pbc.FocusCostRecord, 100)
	for i := range 100 {
		records[i] = &pbc.FocusCostRecord{
			BillingAccountId: "123456789012",
			ChargePeriodStart: &timestamppb.Timestamp{
				Seconds: int64(1735689600 + i*86400),
			},
			ServiceName:     "Amazon EC2",
			BilledCost:      float64(100 + i),
			BillingCurrency: "USD",
		}
	}

	var buf bytes.Buffer
	result, err := serializer.SerializeSlice(records, &buf)
	if err != nil {
		t.Fatalf("SerializeSlice() failed: %v", err)
	}

	if result.RecordsWritten != 100 {
		t.Errorf("Expected 100 records written, got %d", result.RecordsWritten)
	}

	if result.HasErrors() {
		t.Errorf("Expected no errors, got %d errors", len(result.Errors))
	}

	// Verify output is valid JSON array
	var parsed []interface{}
	if unmarshalErr := json.Unmarshal(buf.Bytes(), &parsed); unmarshalErr != nil {
		t.Errorf("Output is not valid JSON array: %v", unmarshalErr)
	}

	if len(parsed) != 100 {
		t.Errorf("Expected 100 records in output, got %d", len(parsed))
	}
}

func TestSerializeStream_LargeBatch(t *testing.T) {
	serializer := jsonld.NewSerializer()

	// Test with 10,000 records as per spec
	numRecords := 10000
	records := make([]*pbc.FocusCostRecord, numRecords)
	for i := range numRecords {
		records[i] = &pbc.FocusCostRecord{
			BillingAccountId: "123456789012",
			ChargePeriodStart: &timestamppb.Timestamp{
				Seconds: int64(1735689600 + i*86400),
			},
			ServiceName:     "Amazon EC2",
			BilledCost:      float64(i),
			BillingCurrency: "USD",
		}
	}

	var buf bytes.Buffer
	result, err := serializer.SerializeSlice(records, &buf)
	if err != nil {
		t.Fatalf("SerializeSlice() failed: %v", err)
	}

	if result.RecordsWritten != numRecords {
		t.Errorf("Expected %d records written, got %d", numRecords, result.RecordsWritten)
	}

	if result.HasErrors() {
		t.Errorf("Expected no errors, got %d errors", len(result.Errors))
	}

	// Verify output is valid JSON
	var parsed []interface{}
	if unmarshalErr := json.Unmarshal(buf.Bytes(), &parsed); unmarshalErr != nil {
		t.Errorf("Output is not valid JSON array: %v", unmarshalErr)
	}

	if len(parsed) != numRecords {
		t.Errorf("Expected %d records in output, got %d", numRecords, len(parsed))
	}
}

func TestSerializeStream_ChannelInput(t *testing.T) {
	serializer := jsonld.NewSerializer()

	// Create a channel with records
	ch := make(chan *pbc.FocusCostRecord, 5)

	// Add valid records
	for i := range 5 {
		ch <- &pbc.FocusCostRecord{
			BillingAccountId: "123456789012",
			ChargePeriodStart: &timestamppb.Timestamp{
				Seconds: int64(1735689600 + i*86400),
			},
			ServiceName:     "Amazon EC2",
			BilledCost:      float64(100 + i),
			BillingCurrency: "USD",
		}
	}
	close(ch)

	var buf bytes.Buffer
	result, err := serializer.SerializeStream(ch, &buf)
	if err != nil {
		t.Fatalf("SerializeStream() failed: %v", err)
	}

	if result.RecordsWritten != 5 {
		t.Errorf("Expected 5 records written, got %d", result.RecordsWritten)
	}

	// Verify output is valid JSON
	var parsed []interface{}
	if unmarshalErr := json.Unmarshal(buf.Bytes(), &parsed); unmarshalErr != nil {
		t.Errorf("Output is not valid JSON array: %v", unmarshalErr)
	}
}

func TestSerializeBatch(t *testing.T) {
	serializer := jsonld.NewSerializer()

	records := []*pbc.FocusCostRecord{
		{
			BillingAccountId:  "123456789012",
			ChargePeriodStart: &timestamppb.Timestamp{Seconds: 1735689600},
			ServiceName:       "Amazon EC2",
			BilledCost:        100.0,
			BillingCurrency:   "USD",
		},
		{
			BillingAccountId:  "123456789012",
			ChargePeriodStart: &timestamppb.Timestamp{Seconds: 1735776000},
			ServiceName:       "Amazon S3",
			BilledCost:        50.0,
			BillingCurrency:   "USD",
		},
	}

	data, result, err := serializer.SerializeBatch(records)
	if err != nil {
		t.Fatalf("SerializeBatch() failed: %v", err)
	}

	if result.RecordsWritten != 2 {
		t.Errorf("Expected 2 records written, got %d", result.RecordsWritten)
	}

	// Verify output is valid JSON array
	var parsed []map[string]interface{}
	if unmarshalErr := json.Unmarshal(data, &parsed); unmarshalErr != nil {
		t.Errorf("Output is not valid JSON array: %v", unmarshalErr)
	}

	if len(parsed) != 2 {
		t.Errorf("Expected 2 records in output, got %d", len(parsed))
	}

	// Check first record has expected fields
	if parsed[0]["serviceName"] != "Amazon EC2" {
		t.Errorf("First record serviceName = %v, want Amazon EC2", parsed[0]["serviceName"])
	}

	// Check second record has expected fields
	if parsed[1]["serviceName"] != "Amazon S3" {
		t.Errorf("Second record serviceName = %v, want Amazon S3", parsed[1]["serviceName"])
	}
}

func TestSerializeStream_EmptyInput(t *testing.T) {
	serializer := jsonld.NewSerializer()

	ch := make(chan *pbc.FocusCostRecord)
	close(ch) // Empty channel

	var buf bytes.Buffer
	result, err := serializer.SerializeStream(ch, &buf)
	if err != nil {
		t.Fatalf("SerializeStream() failed: %v", err)
	}

	if result.RecordsWritten != 0 {
		t.Errorf("Expected 0 records written, got %d", result.RecordsWritten)
	}

	// Verify output is empty JSON array
	var parsed []interface{}
	if unmarshalErr := json.Unmarshal(buf.Bytes(), &parsed); unmarshalErr != nil {
		t.Errorf("Output is not valid JSON array: %v", unmarshalErr)
	}

	if len(parsed) != 0 {
		t.Errorf("Expected empty array, got %d elements", len(parsed))
	}
}

func TestStreamResult_HasErrors(t *testing.T) {
	tests := []struct {
		name     string
		errors   []*jsonld.StreamError
		expected bool
	}{
		{
			name:     "no errors",
			errors:   nil,
			expected: false,
		},
		{
			name:     "empty errors",
			errors:   []*jsonld.StreamError{},
			expected: false,
		},
		{
			name: "with errors",
			errors: []*jsonld.StreamError{
				{Index: 0, Message: "test error"},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &jsonld.StreamResult{Errors: tt.errors}
			if got := result.HasErrors(); got != tt.expected {
				t.Errorf("HasErrors() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStreamError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *jsonld.StreamError
		expected string
	}{
		{
			name: "error without wrapped error",
			err: &jsonld.StreamError{
				Index:   0,
				Message: "test error",
				Err:     nil,
			},
			expected: "stream error at index 0: test error",
		},
		{
			name: "error with wrapped error",
			err: &jsonld.StreamError{
				Index:   5,
				Message: "serialization failed",
				Err:     errors.New("underlying error"),
			},
			expected: "stream error at index 5: serialization failed: underlying error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestStreamError_Unwrap(t *testing.T) {
	underlyingErr := errors.New("underlying error")

	tests := []struct {
		name     string
		err      *jsonld.StreamError
		expected error
	}{
		{
			name: "unwrap nil error",
			err: &jsonld.StreamError{
				Index:   0,
				Message: "test",
				Err:     nil,
			},
			expected: nil,
		},
		{
			name: "unwrap wrapped error",
			err: &jsonld.StreamError{
				Index:   1,
				Message: "test",
				Err:     underlyingErr,
			},
			expected: underlyingErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Unwrap()
			// Use errors.Is for proper error comparison (handles wrapped errors correctly)
			if tt.expected == nil {
				if got != nil {
					t.Errorf("Unwrap() = %v, want nil", got)
				}
			} else if !errors.Is(got, tt.expected) {
				t.Errorf("Unwrap() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSerializeStream_PrettyPrint(t *testing.T) {
	serializer := jsonld.NewSerializer(jsonld.WithPrettyPrint(true))

	records := []*pbc.FocusCostRecord{
		{
			BillingAccountId:  "123456789012",
			ChargePeriodStart: &timestamppb.Timestamp{Seconds: 1735689600},
			ServiceName:       "Amazon EC2",
			BilledCost:        100.0,
			BillingCurrency:   "USD",
		},
	}

	data, _, err := serializer.SerializeBatch(records)
	if err != nil {
		t.Fatalf("SerializeBatch() failed: %v", err)
	}

	// Pretty print should include newlines
	output := string(data)
	if !bytes.Contains(data, []byte("
")) {
		t.Error("Expected pretty-printed output to contain newlines")
	}

	// Should still be valid JSON
	var parsed []interface{}
	if unmarshalErr := json.Unmarshal([]byte(output), &parsed); unmarshalErr != nil {
		t.Errorf("Output is not valid JSON: %v", unmarshalErr)
	}
}
