// Package contracts defines the public API for the pluginsdk/mapping package.
// This file documents the function signatures and expected behavior.
// It is NOT executable code - it serves as a contract specification.
package contracts

// =============================================================================
// AWS Extraction Functions
// =============================================================================

// ExtractAWSSKU extracts the SKU (instance type, volume type, etc.) from AWS
// resource properties.
//
// Key priority order:
//  1. instanceType - EC2 instances
//  2. instanceClass - RDS instances
//  3. type - Generic fallback
//  4. volumeType - EBS volumes
//
// Returns empty string if no key found or input is nil/empty.
// Never panics.
//
// Example:
//
//	props := map[string]string{"instanceType": "t3.medium"}
//	sku := mapping.ExtractAWSSKU(props) // Returns "t3.medium"
func ExtractAWSSKU(_ map[string]string) string { return "" }

// ExtractAWSRegion extracts the region from AWS resource properties.
//
// Key priority order:
//  1. region - Explicit region setting
//  2. availabilityZone - Derived via ExtractAWSRegionFromAZ
//
// Returns empty string if no key found or input is nil/empty.
// Never panics.
//
// Example:
//
//	props := map[string]string{"availabilityZone": "us-east-1a"}
//	region := mapping.ExtractAWSRegion(props) // Returns "us-east-1"
func ExtractAWSRegion(_ map[string]string) string { return "" }

// ExtractAWSRegionFromAZ derives the AWS region from an availability zone string.
//
// Algorithm: Removes trailing lowercase letter(s) from the zone name.
//
// Examples:
//   - "us-east-1a" → "us-east-1"
//   - "eu-west-2b" → "eu-west-2"
//   - "ap-northeast-1c" → "ap-northeast-1"
//
// Returns empty string if input is empty.
// Never panics.
func ExtractAWSRegionFromAZ(_ string) string { return "" }

// =============================================================================
// Azure Extraction Functions
// =============================================================================

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
//
// Example:
//
//	props := map[string]string{"vmSize": "Standard_D2s_v3"}
//	sku := mapping.ExtractAzureSKU(props) // Returns "Standard_D2s_v3"
func ExtractAzureSKU(_ map[string]string) string { return "" }

// ExtractAzureRegion extracts the region (location) from Azure resource properties.
//
// Key priority order:
//  1. location - Primary Azure location field
//  2. region - Alternative field name
//
// Returns empty string if no key found or input is nil/empty.
// Never panics.
//
// Example:
//
//	props := map[string]string{"location": "eastus"}
//	region := mapping.ExtractAzureRegion(props) // Returns "eastus"
func ExtractAzureRegion(_ map[string]string) string { return "" }

// =============================================================================
// GCP Extraction Functions
// =============================================================================

// ExtractGCPSKU extracts the SKU (machine type, tier, etc.) from GCP resource
// properties.
//
// Key priority order:
//  1. machineType - Compute instances
//  2. type - Generic type field
//  3. tier - Service tier
//
// Returns empty string if no key found or input is nil/empty.
// Never panics.
//
// Example:
//
//	props := map[string]string{"machineType": "n1-standard-4"}
//	sku := mapping.ExtractGCPSKU(props) // Returns "n1-standard-4"
func ExtractGCPSKU(_ map[string]string) string { return "" }

// ExtractGCPRegion extracts the region from GCP resource properties.
//
// Key priority order:
//  1. region - Explicit region setting
//  2. zone - Derived via ExtractGCPRegionFromZone
//
// Returns empty string if no key found or input is nil/empty.
// Never panics.
//
// Example:
//
//	props := map[string]string{"zone": "us-central1-a"}
//	region := mapping.ExtractGCPRegion(props) // Returns "us-central1"
func ExtractGCPRegion(_ map[string]string) string { return "" }

// ExtractGCPRegionFromZone derives the GCP region from a zone string.
//
// Algorithm:
//  1. Remove the last hyphen-delimited segment (e.g., "us-central1-a" → "us-central1")
//  2. Validate against known GCP regions list
//  3. Return empty string if not a valid GCP region
//
// Returns empty string if input is empty, has no hyphen, or produces invalid region.
// Never panics.
//
// Example:
//
//	region := mapping.ExtractGCPRegionFromZone("europe-west1-b") // Returns "europe-west1"
//	region := mapping.ExtractGCPRegionFromZone("invalid-zone")   // Returns ""
func ExtractGCPRegionFromZone(_ string) string { return "" }

// IsValidGCPRegion checks if a string is a known valid GCP region.
//
// Uses the package-level GCPRegions slice for validation.
// Returns false for empty string.
//
// Example:
//
//	valid := mapping.IsValidGCPRegion("us-central1")  // Returns true
//	valid := mapping.IsValidGCPRegion("invalid")     // Returns false
func IsValidGCPRegion(_ string) bool { return false }

// =============================================================================
// Generic Extraction Functions
// =============================================================================

// ExtractSKU extracts a value from properties using the provided key list.
//
// Checks keys in order and returns the first non-empty value found.
// If no keys provided, uses default SKU keys: "sku", "type", "tier".
//
// Returns empty string if no matching key found or input is nil/empty.
// Never panics.
//
// Example:
//
//	props := map[string]string{"customSKU": "my-sku"}
//	sku := mapping.ExtractSKU(props, "customSKU", "fallbackSKU") // Returns "my-sku"
func ExtractSKU(_ map[string]string, _ ...string) string { return "" }

// ExtractRegion extracts a value from properties using the provided key list.
//
// Checks keys in order and returns the first non-empty value found.
// If no keys provided, uses default region keys: "region", "location", "zone".
//
// Returns empty string if no matching key found or input is nil/empty.
// Never panics.
//
// Example:
//
//	props := map[string]string{"customRegion": "us-east"}
//	region := mapping.ExtractRegion(props, "customRegion", "fallback") // Returns "us-east"
func ExtractRegion(_ map[string]string, _ ...string) string { return "" }

// =============================================================================
// Package Variables
// =============================================================================

// GCPRegions is the list of known valid GCP regions.
// Used by IsValidGCPRegion and ExtractGCPRegionFromZone for validation.
// Updated: 2025-12
//
//nolint:gochecknoglobals // Intentional: read-only reference data
var GCPRegions []string
