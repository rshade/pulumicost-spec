package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

// Embedded schema content (normally this would be generated or read from file)
const schemaJSON = `{
"$schema": "https://json-schema.org/draft/2020-12/schema",
"$id": "https://spec.pulumicost.dev/schemas/pricing_spec.schema.json",
"title": "PricingSpec",
"type": "object",
"required": ["provider", "resource_type", "billing_mode", "rate_per_unit", "currency"],
"properties": {
    "provider":      { "type": "string", "minLength": 1 },
    "resource_type": { "type": "string", "minLength": 1 },
    "sku":           { "type": "string" },
    "region":        { "type": "string" },
    "billing_mode":  { "type": "string", "enum": ["per_hour","per_gb_month","per_request","flat","per_day","per_cpu_hour"] },
    "rate_per_unit": { "type": "number", "minimum": 0 },
    "currency":      { "type": "string", "minLength": 1 },
    "description":   { "type": "string" },
    "metric_hints": {
    "type": "array",
    "items": {
        "type": "object",
        "required": ["metric","unit"],
        "properties": {
        "metric": { "type": "string", "minLength": 1 },
        "unit":   { "type": "string", "minLength": 1 }
        },
        "additionalProperties": false
    }
    },
    "plugin_metadata": {
    "type": "object",
    "additionalProperties": { "type": "string" }
    },
    "source": { "type": "string" }
},
"additionalProperties": false
}`

var compiled *jsonschema.Schema

func init() {
	c := jsonschema.NewCompiler()
	if err := c.AddResource("schema.json", strings.NewReader(schemaJSON)); err != nil {
		panic(err)
	}
	s, err := c.Compile("schema.json")
	if err != nil {
		panic(err)
	}
	compiled = s
}

func ValidatePricingSpec(doc []byte) error {
	if compiled == nil {
		return fmt.Errorf("schema not loaded")
	}
	var v interface{}
	if err := json.Unmarshal(doc, &v); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}
	return compiled.Validate(v)
}
