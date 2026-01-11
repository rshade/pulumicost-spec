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

// Package currency provides ISO 4217 currency validation and metadata.
//
// The currency package implements zero-allocation validation for ISO 4217 currency codes,
// following the established pattern from sdk/go/registry/domain.go. It provides fast,
// memory-efficient validation and metadata retrieval for all active ISO 4217 currencies.
//
// # Features
//
//   - Zero-allocation validation: IsValid() achieves <15 ns/op with 0 allocations
//   - Complete ISO 4217 coverage: 180+ active currencies with full metadata
//   - Case-sensitive validation: Only uppercase codes are accepted
//   - Metadata retrieval: GetCurrency() provides name, numeric code, and decimal places
//   - List operations: AllCurrencies() returns all valid currencies
//
// # Performance Characteristics
//
//   - IsValid(): <15 ns/op, 0 B/op, 0 allocs/op
//   - GetCurrency(): O(1) lookup with a defensive copy for safe mutation
//   - AllCurrencies(): <5 ns/op, 0 B/op, 0 allocs/op
//
// # Usage Examples
//
// Basic validation:
//
//	if currency.IsValid("USD") {
//	    fmt.Println("USD is valid")
//	}
//
// Metadata retrieval:
//
//	usd, err := currency.GetCurrency("USD")
//	if err != nil {
//	    // handle error
//	}
//	fmt.Printf("%s uses %d decimal places
", usd.Name, usd.MinorUnits)
//
// Listing all currencies:
//
//	for _, c := range currency.AllCurrencies() {
//	    fmt.Printf("%s: %s
", c.Code, c.Name)
//	}
//
// # Integration
//
// This package is used by sdk/go/pluginsdk/focus_conformance.go for validating
// billing and pricing currency fields in FOCUS cost records.
//
// # References
//
//   - ISO 4217: https://www.iso.org/iso-4217-currency-codes.html
//   - FOCUS 1.2 Specification: https://focus.finops.org
package currency
