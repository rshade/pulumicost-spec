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

//nolint:testpackage // White-box testing to maintain consistent test package across test files
package mapping

import "testing"

// =============================================================================
// AWS Extraction Tests (User Story 1)
// =============================================================================

func TestExtractAWSSKU(t *testing.T) {
	tests := []struct {
		name       string
		properties map[string]string
		want       string
	}{
		{
			name:       "nil map",
			properties: nil,
			want:       "",
		},
		{
			name:       "empty map",
			properties: map[string]string{},
			want:       "",
		},
		{
			name:       "EC2 instanceType",
			properties: map[string]string{"instanceType": "t3.medium"},
			want:       "t3.medium",
		},
		{
			name:       "RDS instanceClass",
			properties: map[string]string{"instanceClass": "db.t3.micro"},
			want:       "db.t3.micro",
		},
		{
			name:       "generic type",
			properties: map[string]string{"type": "some-type"},
			want:       "some-type",
		},
		{
			name:       "EBS volumeType",
			properties: map[string]string{"volumeType": "gp3"},
			want:       "gp3",
		},
		{
			name: "instanceType takes priority over instanceClass",
			properties: map[string]string{
				"instanceType":  "t3.large",
				"instanceClass": "db.m5.large",
			},
			want: "t3.large",
		},
		{
			name: "instanceClass takes priority over type",
			properties: map[string]string{
				"instanceClass": "db.m5.large",
				"type":          "some-type",
			},
			want: "db.m5.large",
		},
		{
			name: "type takes priority over volumeType",
			properties: map[string]string{
				"type":       "some-type",
				"volumeType": "gp3",
			},
			want: "some-type",
		},
		{
			name:       "empty value treated as not found",
			properties: map[string]string{"instanceType": ""},
			want:       "",
		},
		{
			name:       "empty value falls through to next key",
			properties: map[string]string{"instanceType": "", "instanceClass": "db.t3.micro"},
			want:       "db.t3.micro",
		},
		{
			name:       "unrelated keys",
			properties: map[string]string{"someOtherKey": "value"},
			want:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractAWSSKU(tt.properties)
			if got != tt.want {
				t.Errorf("ExtractAWSSKU() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractAWSRegion(t *testing.T) {
	tests := []struct {
		name       string
		properties map[string]string
		want       string
	}{
		{
			name:       "nil map",
			properties: nil,
			want:       "",
		},
		{
			name:       "empty map",
			properties: map[string]string{},
			want:       "",
		},
		{
			name:       "explicit region",
			properties: map[string]string{"region": "us-east-1"},
			want:       "us-east-1",
		},
		{
			name:       "derived from availabilityZone",
			properties: map[string]string{"availabilityZone": "us-east-1a"},
			want:       "us-east-1",
		},
		{
			name:       "derived from availabilityZone us-west-2b",
			properties: map[string]string{"availabilityZone": "us-west-2b"},
			want:       "us-west-2",
		},
		{
			name:       "derived from availabilityZone eu-central-1c",
			properties: map[string]string{"availabilityZone": "eu-central-1c"},
			want:       "eu-central-1",
		},
		{
			name: "region takes priority over availabilityZone",
			properties: map[string]string{
				"region":           "us-west-1",
				"availabilityZone": "us-east-1a",
			},
			want: "us-west-1",
		},
		{
			name:       "empty region falls through to availabilityZone",
			properties: map[string]string{"region": "", "availabilityZone": "us-east-1a"},
			want:       "us-east-1",
		},
		{
			name:       "unrelated keys",
			properties: map[string]string{"someOtherKey": "value"},
			want:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractAWSRegion(tt.properties)
			if got != tt.want {
				t.Errorf("ExtractAWSRegion() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractAWSRegionFromAZ(t *testing.T) {
	tests := []struct {
		name             string
		availabilityZone string
		want             string
	}{
		{
			name:             "empty string",
			availabilityZone: "",
			want:             "",
		},
		{
			name:             "us-east-1a",
			availabilityZone: "us-east-1a",
			want:             "us-east-1",
		},
		{
			name:             "us-west-2b",
			availabilityZone: "us-west-2b",
			want:             "us-west-2",
		},
		{
			name:             "eu-central-1c",
			availabilityZone: "eu-central-1c",
			want:             "eu-central-1",
		},
		{
			name:             "ap-northeast-1d",
			availabilityZone: "ap-northeast-1d",
			want:             "ap-northeast-1",
		},
		{
			name:             "single letter zone suffix",
			availabilityZone: "us-east-1a",
			want:             "us-east-1",
		},
		{
			name:             "no suffix - returns as-is",
			availabilityZone: "us-east-1",
			want:             "us-east-1",
		},
		{
			name:             "just a letter returns empty",
			availabilityZone: "a",
			want:             "", // Single letter has no region prefix
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractAWSRegionFromAZ(tt.availabilityZone)
			if got != tt.want {
				t.Errorf("ExtractAWSRegionFromAZ(%q) = %q, want %q", tt.availabilityZone, got, tt.want)
			}
		})
	}
}

// =============================================================================
// Azure Extraction Tests (User Story 2)
// =============================================================================

func TestExtractAzureSKU(t *testing.T) {
	tests := []struct {
		name       string
		properties map[string]string
		want       string
	}{
		{
			name:       "nil map",
			properties: nil,
			want:       "",
		},
		{
			name:       "empty map",
			properties: map[string]string{},
			want:       "",
		},
		{
			name:       "vmSize",
			properties: map[string]string{"vmSize": "Standard_D2s_v3"},
			want:       "Standard_D2s_v3",
		},
		{
			name:       "sku",
			properties: map[string]string{"sku": "Standard"},
			want:       "Standard",
		},
		{
			name:       "tier",
			properties: map[string]string{"tier": "Premium"},
			want:       "Premium",
		},
		{
			name: "vmSize takes priority over sku",
			properties: map[string]string{
				"vmSize": "Standard_D2s_v3",
				"sku":    "Standard",
			},
			want: "Standard_D2s_v3",
		},
		{
			name: "sku takes priority over tier",
			properties: map[string]string{
				"sku":  "Standard",
				"tier": "Premium",
			},
			want: "Standard",
		},
		{
			name:       "empty value treated as not found",
			properties: map[string]string{"vmSize": ""},
			want:       "",
		},
		{
			name:       "unrelated keys",
			properties: map[string]string{"someOtherKey": "value"},
			want:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractAzureSKU(tt.properties)
			if got != tt.want {
				t.Errorf("ExtractAzureSKU() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractAzureRegion(t *testing.T) {
	tests := []struct {
		name       string
		properties map[string]string
		want       string
	}{
		{
			name:       "nil map",
			properties: nil,
			want:       "",
		},
		{
			name:       "empty map",
			properties: map[string]string{},
			want:       "",
		},
		{
			name:       "location",
			properties: map[string]string{"location": "eastus"},
			want:       "eastus",
		},
		{
			name:       "region",
			properties: map[string]string{"region": "westeurope"},
			want:       "westeurope",
		},
		{
			name: "location takes priority over region",
			properties: map[string]string{
				"location": "eastus",
				"region":   "westeurope",
			},
			want: "eastus",
		},
		{
			name:       "empty value treated as not found",
			properties: map[string]string{"location": ""},
			want:       "",
		},
		{
			name:       "empty location falls through to region",
			properties: map[string]string{"location": "", "region": "westeurope"},
			want:       "westeurope",
		},
		{
			name:       "unrelated keys",
			properties: map[string]string{"someOtherKey": "value"},
			want:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractAzureRegion(tt.properties)
			if got != tt.want {
				t.Errorf("ExtractAzureRegion() = %q, want %q", got, tt.want)
			}
		})
	}
}

// =============================================================================
// GCP Extraction Tests (User Story 3)
// =============================================================================

func TestExtractGCPSKU(t *testing.T) {
	tests := []struct {
		name       string
		properties map[string]string
		want       string
	}{
		{
			name:       "nil map",
			properties: nil,
			want:       "",
		},
		{
			name:       "empty map",
			properties: map[string]string{},
			want:       "",
		},
		{
			name:       "machineType",
			properties: map[string]string{"machineType": "n1-standard-4"},
			want:       "n1-standard-4",
		},
		{
			name:       "type",
			properties: map[string]string{"type": "pd-ssd"},
			want:       "pd-ssd",
		},
		{
			name:       "tier",
			properties: map[string]string{"tier": "db-custom"},
			want:       "db-custom",
		},
		{
			name: "machineType takes priority over type",
			properties: map[string]string{
				"machineType": "n1-standard-4",
				"type":        "pd-ssd",
			},
			want: "n1-standard-4",
		},
		{
			name: "type takes priority over tier",
			properties: map[string]string{
				"type": "pd-ssd",
				"tier": "db-custom",
			},
			want: "pd-ssd",
		},
		{
			name:       "empty value treated as not found",
			properties: map[string]string{"machineType": ""},
			want:       "",
		},
		{
			name:       "unrelated keys",
			properties: map[string]string{"someOtherKey": "value"},
			want:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractGCPSKU(tt.properties)
			if got != tt.want {
				t.Errorf("ExtractGCPSKU() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractGCPRegion(t *testing.T) {
	tests := []struct {
		name       string
		properties map[string]string
		want       string
	}{
		{
			name:       "nil map",
			properties: nil,
			want:       "",
		},
		{
			name:       "empty map",
			properties: map[string]string{},
			want:       "",
		},
		{
			name:       "explicit region",
			properties: map[string]string{"region": "us-central1"},
			want:       "us-central1",
		},
		{
			name:       "derived from zone",
			properties: map[string]string{"zone": "us-central1-a"},
			want:       "us-central1",
		},
		{
			name:       "derived from zone europe-west1-b",
			properties: map[string]string{"zone": "europe-west1-b"},
			want:       "europe-west1",
		},
		{
			name: "region takes priority over zone",
			properties: map[string]string{
				"region": "asia-east1",
				"zone":   "us-central1-a",
			},
			want: "asia-east1",
		},
		{
			name:       "empty region falls through to zone",
			properties: map[string]string{"region": "", "zone": "us-central1-a"},
			want:       "us-central1",
		},
		{
			name:       "invalid zone returns empty",
			properties: map[string]string{"zone": "invalid-zone-x"},
			want:       "",
		},
		{
			name:       "unrelated keys",
			properties: map[string]string{"someOtherKey": "value"},
			want:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractGCPRegion(tt.properties)
			if got != tt.want {
				t.Errorf("ExtractGCPRegion() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractGCPRegionFromZone(t *testing.T) {
	tests := []struct {
		name string
		zone string
		want string
	}{
		{
			name: "empty string",
			zone: "",
			want: "",
		},
		{
			name: "us-central1-a",
			zone: "us-central1-a",
			want: "us-central1",
		},
		{
			name: "europe-west1-b",
			zone: "europe-west1-b",
			want: "europe-west1",
		},
		{
			name: "asia-east1-c",
			zone: "asia-east1-c",
			want: "asia-east1",
		},
		{
			name: "invalid region after extraction",
			zone: "invalid-zone-x",
			want: "",
		},
		{
			name: "no hyphen",
			zone: "nohyphen",
			want: "",
		},
		{
			name: "single character",
			zone: "a",
			want: "",
		},
		{
			name: "hyphen at start",
			zone: "-a",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractGCPRegionFromZone(tt.zone)
			if got != tt.want {
				t.Errorf("ExtractGCPRegionFromZone(%q) = %q, want %q", tt.zone, got, tt.want)
			}
		})
	}
}

func TestIsValidGCPRegion(t *testing.T) {
	tests := []struct {
		name   string
		region string
		want   bool
	}{
		{
			name:   "empty string",
			region: "",
			want:   false,
		},
		{
			name:   "us-central1",
			region: "us-central1",
			want:   true,
		},
		{
			name:   "europe-west1",
			region: "europe-west1",
			want:   true,
		},
		{
			name:   "asia-east1",
			region: "asia-east1",
			want:   true,
		},
		{
			name:   "invalid region",
			region: "invalid-region",
			want:   false,
		},
		{
			name:   "partial match",
			region: "us-central",
			want:   false,
		},
		{
			name:   "zone instead of region",
			region: "us-central1-a",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidGCPRegion(tt.region)
			if got != tt.want {
				t.Errorf("IsValidGCPRegion(%q) = %v, want %v", tt.region, got, tt.want)
			}
		})
	}
}

// =============================================================================
// Generic Extraction Tests (User Story 4)
// =============================================================================

func TestExtractSKU(t *testing.T) {
	tests := []struct {
		name       string
		properties map[string]string
		keys       []string
		want       string
	}{
		{
			name:       "nil map",
			properties: nil,
			keys:       nil,
			want:       "",
		},
		{
			name:       "empty map with defaults",
			properties: map[string]string{},
			keys:       nil,
			want:       "",
		},
		{
			name:       "default key sku",
			properties: map[string]string{"sku": "custom-sku"},
			keys:       nil,
			want:       "custom-sku",
		},
		{
			name:       "default key type",
			properties: map[string]string{"type": "custom-type"},
			keys:       nil,
			want:       "custom-type",
		},
		{
			name:       "default key tier",
			properties: map[string]string{"tier": "custom-tier"},
			keys:       nil,
			want:       "custom-tier",
		},
		{
			name:       "custom key",
			properties: map[string]string{"customField": "custom-value"},
			keys:       []string{"customField"},
			want:       "custom-value",
		},
		{
			name:       "custom key priority",
			properties: map[string]string{"first": "value1", "second": "value2"},
			keys:       []string{"first", "second"},
			want:       "value1",
		},
		{
			name:       "custom key fallback",
			properties: map[string]string{"second": "value2"},
			keys:       []string{"first", "second"},
			want:       "value2",
		},
		{
			name:       "empty keys uses defaults",
			properties: map[string]string{"sku": "default-sku"},
			keys:       []string{},
			want:       "default-sku",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractSKU(tt.properties, tt.keys...)
			if got != tt.want {
				t.Errorf("ExtractSKU() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractRegion(t *testing.T) {
	tests := []struct {
		name       string
		properties map[string]string
		keys       []string
		want       string
	}{
		{
			name:       "nil map",
			properties: nil,
			keys:       nil,
			want:       "",
		},
		{
			name:       "empty map with defaults",
			properties: map[string]string{},
			keys:       nil,
			want:       "",
		},
		{
			name:       "default key region",
			properties: map[string]string{"region": "custom-region"},
			keys:       nil,
			want:       "custom-region",
		},
		{
			name:       "default key location",
			properties: map[string]string{"location": "custom-location"},
			keys:       nil,
			want:       "custom-location",
		},
		{
			name:       "default key zone",
			properties: map[string]string{"zone": "custom-zone"},
			keys:       nil,
			want:       "custom-zone",
		},
		{
			name:       "custom key",
			properties: map[string]string{"customRegion": "my-region"},
			keys:       []string{"customRegion"},
			want:       "my-region",
		},
		{
			name:       "custom key priority",
			properties: map[string]string{"first": "region1", "second": "region2"},
			keys:       []string{"first", "second"},
			want:       "region1",
		},
		{
			name:       "custom key fallback",
			properties: map[string]string{"second": "region2"},
			keys:       []string{"first", "second"},
			want:       "region2",
		},
		{
			name:       "empty keys uses defaults",
			properties: map[string]string{"region": "default-region"},
			keys:       []string{},
			want:       "default-region",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractRegion(tt.properties, tt.keys...)
			if got != tt.want {
				t.Errorf("ExtractRegion() = %q, want %q", got, tt.want)
			}
		})
	}
}
