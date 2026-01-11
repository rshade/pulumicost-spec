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

//nolint:testpackage // Internal package testing for extractFromKeys helper
package mapping

import "testing"

// Benchmark input data.
//
//nolint:gochecknoglobals // Test fixture data
var benchProps = map[string]string{
	"instanceType":     "t3.medium",
	"instanceClass":    "db.t3.micro",
	"type":             "some-type",
	"volumeType":       "gp3",
	"region":           "us-east-1",
	"availabilityZone": "us-east-1a",
	"vmSize":           "Standard_D2s_v3",
	"location":         "eastus",
	"machineType":      "n1-standard-4",
	"zone":             "us-central1-a",
	"sku":              "custom-sku",
	"tier":             "Premium",
}

// =============================================================================
// AWS Benchmarks
// =============================================================================

func BenchmarkExtractAWSSKU(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = ExtractAWSSKU(benchProps)
	}
}

func BenchmarkExtractAWSRegion(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = ExtractAWSRegion(benchProps)
	}
}

func BenchmarkExtractAWSRegionFromAZ(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = ExtractAWSRegionFromAZ("us-east-1a")
	}
}

func BenchmarkExtractAWSSKU_NilMap(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = ExtractAWSSKU(nil)
	}
}

func BenchmarkExtractAWSRegion_NilMap(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = ExtractAWSRegion(nil)
	}
}

// =============================================================================
// Azure Benchmarks
// =============================================================================

func BenchmarkExtractAzureSKU(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = ExtractAzureSKU(benchProps)
	}
}

func BenchmarkExtractAzureRegion(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = ExtractAzureRegion(benchProps)
	}
}

// =============================================================================
// GCP Benchmarks
// =============================================================================

func BenchmarkExtractGCPSKU(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = ExtractGCPSKU(benchProps)
	}
}

func BenchmarkExtractGCPRegion(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = ExtractGCPRegion(benchProps)
	}
}

func BenchmarkExtractGCPRegionFromZone(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = ExtractGCPRegionFromZone("us-central1-a")
	}
}

func BenchmarkIsValidGCPRegion(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = IsValidGCPRegion("us-central1")
	}
}

func BenchmarkIsValidGCPRegion_Invalid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = IsValidGCPRegion("invalid-region")
	}
}

func BenchmarkIsValidGCPRegion_LastRegion(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = IsValidGCPRegion("southamerica-west1")
	}
}

// =============================================================================
// Generic Benchmarks
// =============================================================================

func BenchmarkExtractSKU_DefaultKeys(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = ExtractSKU(benchProps)
	}
}

func BenchmarkExtractSKU_CustomKeys(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = ExtractSKU(benchProps, "customKey", "fallback1", "fallback2")
	}
}

func BenchmarkExtractRegion_DefaultKeys(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = ExtractRegion(benchProps)
	}
}

func BenchmarkExtractRegion_CustomKeys(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = ExtractRegion(benchProps, "customRegion", "fallback1", "fallback2")
	}
}

// =============================================================================
// Internal Helper Benchmarks
// =============================================================================

func BenchmarkExtractFromKeys(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = extractFromKeys(benchProps, "instanceType", "instanceClass", "type", "volumeType")
	}
}

func BenchmarkExtractFromKeys_NilMap(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = extractFromKeys(nil, "instanceType", "instanceClass")
	}
}
