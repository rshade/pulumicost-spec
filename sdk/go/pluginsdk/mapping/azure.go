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
