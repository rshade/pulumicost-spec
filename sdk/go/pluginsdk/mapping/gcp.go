package mapping

import "strings"

// gcpRegions contains all known GCP regions as of 2025-12.
// Used for validation during zone-to-region extraction.
// This is unexported to prevent external mutation. Use AllGCPRegions() to get a copy.
//
//nolint:gochecknoglobals // Intentional: read-only reference data
var gcpRegions = []string{
	// Asia Pacific
	"asia-east1",      // Taiwan
	"asia-east2",      // Hong Kong
	"asia-northeast1", // Tokyo
	"asia-northeast2", // Osaka
	"asia-northeast3", // Seoul
	"asia-south1",     // Mumbai
	"asia-south2",     // Delhi
	"asia-southeast1", // Singapore
	"asia-southeast2", // Jakarta

	// Australia
	"australia-southeast1", // Sydney
	"australia-southeast2", // Melbourne

	// Europe
	"europe-central2",   // Warsaw
	"europe-north1",     // Finland
	"europe-southwest1", // Madrid
	"europe-west1",      // Belgium
	"europe-west2",      // London
	"europe-west3",      // Frankfurt
	"europe-west4",      // Netherlands
	"europe-west6",      // Zurich
	"europe-west8",      // Milan
	"europe-west9",      // Paris
	"europe-west10",     // Berlin
	"europe-west12",     // Turin

	// Middle East
	"me-central1", // Doha
	"me-central2", // Dammam
	"me-west1",    // Tel Aviv

	// North America
	"northamerica-northeast1", // Montreal
	"northamerica-northeast2", // Toronto
	"us-central1",             // Iowa
	"us-east1",                // South Carolina
	"us-east4",                // Virginia
	"us-east5",                // Columbus
	"us-south1",               // Dallas
	"us-west1",                // Oregon
	"us-west2",                // Los Angeles
	"us-west3",                // Salt Lake City
	"us-west4",                // Las Vegas

	// South America
	"southamerica-east1", // São Paulo
	"southamerica-west1", // Santiago
}

// AllGCPRegions returns a copy of the known GCP regions list.
// This returns a fresh copy to prevent external mutation of the internal list.
func AllGCPRegions() []string {
	result := make([]string, len(gcpRegions))
	copy(result, gcpRegions)
	return result
}

// IsValidGCPRegion checks if a string is a known valid GCP region.
// Returns false for empty string or unknown regions.
func IsValidGCPRegion(region string) bool {
	if region == "" {
		return false
	}
	for _, r := range gcpRegions {
		if r == region {
			return true
		}
	}
	return false
}

// ExtractGCPRegionFromZone derives the GCP region from a zone string.
//
// Algorithm:
//  1. Remove the last hyphen-delimited segment (e.g., "us-central1-a" → "us-central1")
//  2. Validate against known GCP regions list
//  3. Return empty string if not a valid GCP region
//
// Returns empty string if input is empty, has no hyphen, or produces invalid region.
// Never panics.
func ExtractGCPRegionFromZone(zone string) string {
	if zone == "" {
		return ""
	}

	// Find the last hyphen
	lastIdx := strings.LastIndex(zone, "-")
	if lastIdx <= 0 {
		return ""
	}

	// Extract everything before the last hyphen
	region := zone[:lastIdx]

	// Validate against known regions
	if !IsValidGCPRegion(region) {
		return ""
	}

	return region
}

// ExtractGCPSKU extracts the SKU (machine type, tier, etc.) from GCP resource properties.
//
// Key priority order:
//  1. machineType - Compute instances
//  2. type - Generic type field
//  3. tier - Service tier
//
// Returns empty string if no key found or input is nil/empty.
// Never panics.
func ExtractGCPSKU(properties map[string]string) string {
	return extractFromKeys(properties, GCPKeyMachineType, GCPKeyType, GCPKeyTier)
}

// ExtractGCPRegion extracts the region from GCP resource properties.
//
// Key priority order:
//  1. region - Explicit region setting
//  2. zone - Derived via ExtractGCPRegionFromZone
//
// Returns empty string if no key found or input is nil/empty.
// Never panics.
func ExtractGCPRegion(properties map[string]string) string {
	if properties == nil {
		return ""
	}

	// Check explicit region first
	if region := properties[GCPKeyRegion]; region != "" {
		return region
	}

	// Try to derive from zone
	if zone := properties[GCPKeyZone]; zone != "" {
		return ExtractGCPRegionFromZone(zone)
	}

	return ""
}
