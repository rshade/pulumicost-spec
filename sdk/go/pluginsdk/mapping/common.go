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
