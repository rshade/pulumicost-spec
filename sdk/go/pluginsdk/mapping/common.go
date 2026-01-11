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

package mapping

// extractFromKeys checks the provided keys in order and returns the first non-empty value.
// Returns empty string if no matching key found or map is nil/empty.
// This is an internal helper function used by all provider-specific extractors.
func extractFromKeys(properties map[string]string, keys ...string) string {
	if properties == nil {
		return ""
	}
	for _, key := range keys {
		if value := properties[key]; value != "" {
			return value
		}
	}
	return ""
}

// ExtractSKU extracts a value from properties using the provided key list.
//
// Checks keys in order and returns the first non-empty value found.
// If no keys provided, uses default SKU keys: "sku", "type", "tier".
//
// Returns empty string if no matching key found or input is nil/empty.
// Never panics.
func ExtractSKU(properties map[string]string, keys ...string) string {
	if len(keys) == 0 {
		keys = defaultSKUKeys
	}
	return extractFromKeys(properties, keys...)
}

// ExtractRegion extracts a value from properties using the provided key list.
//
// Checks keys in order and returns the first non-empty value found.
// If no keys provided, uses default region keys: "region", "location", "zone".
//
// Returns empty string if no matching key found or input is nil/empty.
// Never panics.
func ExtractRegion(properties map[string]string, keys ...string) string {
	if len(keys) == 0 {
		keys = defaultRegionKeys
	}
	return extractFromKeys(properties, keys...)
}
