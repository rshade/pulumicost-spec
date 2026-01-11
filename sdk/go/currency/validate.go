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

package currency

// IsValid checks if code is a valid ISO 4217 currency code.
// Returns true for active ISO 4217 currencies, false otherwise.
//
// Validation rules:
//   - Case-sensitive: "USD" is valid, "usd" is not
//   - No whitespace: " USD" is not valid
//   - Active only: Historic currencies (e.g., "DEM") return false
//
// Performance: <15 ns/op, 0 allocs/op
//
// Example:
//
//	currency.IsValid("USD") // true
//	currency.IsValid("XYZ") // false
//	currency.IsValid("usd") // false (case-sensitive)
func IsValid(code string) bool {
	// Map lookup for O(1) validation with zero allocations
	// The currencyByCode map is built once at package init
	_, ok := currencyByCode[code]
	return ok
}
