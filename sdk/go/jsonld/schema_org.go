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

package jsonld

// Schema.org vocabulary type definitions for JSON-LD context.
//
// Schema.org provides widely-understood vocabulary types for linked data.
// This package defines constants for Schema.org types and properties used in
// FOCUS cost data serialization where natural mappings exist.
const (
	// SchemaNamespace is the Schema.org namespace IRI.
	SchemaNamespace = "https://schema.org/"

	// Schema.org types for JSON-LD @type declarations.
	MonetaryAmountType = "schema:MonetaryAmount"
	DateTimeType       = "schema:DateTime"
	PropertyValueType  = "schema:PropertyValue"
	ServiceType        = "schema:Service"
	ThingType          = "schema:Thing"
	PlaceType          = "schema:Place"

	// Schema.org properties for compact IRIs.
	SchemaValue        = "schema:value"
	SchemaCurrency     = "schema:currency"
	SchemaName         = "schema:name"
	SchemaSupersededBy = "schema:supersededBy"
)

// MonetaryAmountTypeCoercion returns a JSON-LD type coercion for Schema.org MonetaryAmount.
func MonetaryAmountTypeCoercion() map[string]interface{} {
	return map[string]interface{}{
		"@type": MonetaryAmountType,
		"schema:value": map[string]interface{}{
			"@id":   SchemaValue,
			"@type": "http://www.w3.org/2001/XMLSchema#decimal",
		},
		"schema:currency": SchemaCurrency,
	}
}

// DateTimeTypeCoercion returns a JSON-LD type coercion for ISO 8601 dates.
func DateTimeTypeCoercion() map[string]interface{} {
	return map[string]interface{}{
		"@id":   "schema:DateTime",
		"@type": "http://www.w3.org/2001/XMLSchema#dateTime",
	}
}

//nolint:gochecknoglobals // Intentional optimization for zero-allocation validation
var schemaMappedFields = map[string]bool{
	"billed_cost":          true,
	"list_cost":            true,
	"effective_cost":       true,
	"contracted_cost":      true,
	"charge_period_start":  true,
	"charge_period_end":    true,
	"billing_period_start": true,
	"billing_period_end":   true,
	"service_name":         true,
	"resource_name":        true,
	"region_name":          true,
}

// IsSchemaMapped checks if a field has a natural Schema.org mapping.
//
// Returns true for fields that should use Schema.org types (MonetaryAmount, DateTime, etc.)
// instead of FOCUS namespace.
// Uses zero-allocation lookup via package-level map.
func IsSchemaMapped(fieldName string) bool {
	return schemaMappedFields[fieldName]
}
