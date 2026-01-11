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

// ExtractAzureSKU extracts the SKU (VM size, tier, etc.) from Azure resource
// properties.
//
// Key priority order:
//  1. vmSize - Virtual machines
//  2. sku - Generic SKU field
//  3. tier - Service tier
//
// Returns empty string if no key found or input is nil/empty.
// Never panics.
func ExtractAzureSKU(properties map[string]string) string {
	return extractFromKeys(properties, AzureKeyVMSize, AzureKeySKU, AzureKeyTier)
}

// ExtractAzureRegion extracts the region (location) from Azure resource properties.
//
// Key priority order:
//  1. location - Primary Azure location field
//  2. region - Alternative field name
//
// Returns empty string if no key found or input is nil/empty.
// Never panics.
func ExtractAzureRegion(properties map[string]string) string {
	return extractFromKeys(properties, AzureKeyLocation, AzureKeyRegion)
}
