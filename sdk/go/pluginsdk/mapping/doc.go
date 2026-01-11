// Package mapping provides helper functions for extracting SKU, region, and other
// pricing-relevant fields from Pulumi resource properties.
//
// This package supports AWS, Azure, and GCP cloud providers with provider-specific
// extraction functions, plus generic fallback extractors for custom resource types.
//
// # Provider-Specific Functions
//
// AWS:
//   - ExtractAWSSKU: Extracts SKU from instanceType, instanceClass, type, or volumeType
//   - ExtractAWSRegion: Extracts region from region or availabilityZone
//   - ExtractAWSRegionFromAZ: Derives region from availability zone string
//
// Azure:
//   - ExtractAzureSKU: Extracts SKU from vmSize, sku, or tier
//   - ExtractAzureRegion: Extracts region from location or region
//
// GCP:
//   - ExtractGCPSKU: Extracts SKU from machineType, type, or tier
//   - ExtractGCPRegion: Extracts region from region or zone
//   - ExtractGCPRegionFromZone: Derives region from zone string with validation
//   - IsValidGCPRegion: Validates against known GCP regions list
//
// # Generic Functions
//
//   - ExtractSKU: Generic SKU extraction with custom or default keys
//   - ExtractRegion: Generic region extraction with custom or default keys
//
// # Usage
//
// All functions accept a map[string]string representing Pulumi resource properties.
// All functions return empty string for missing or invalid input and never panic.
//
// AWS Example:
//
//	props := map[string]string{
//	    "instanceType":     "t3.medium",
//	    "availabilityZone": "us-east-1a",
//	}
//	sku := mapping.ExtractAWSSKU(props)       // "t3.medium"
//	region := mapping.ExtractAWSRegion(props) // "us-east-1"
//
// Azure Example:
//
//	props := map[string]string{
//	    "vmSize":   "Standard_D2s_v3",
//	    "location": "eastus",
//	}
//	sku := mapping.ExtractAzureSKU(props)       // "Standard_D2s_v3"
//	region := mapping.ExtractAzureRegion(props) // "eastus"
//
// GCP Example:
//
//	props := map[string]string{
//	    "machineType": "n1-standard-4",
//	    "zone":        "us-central1-a",
//	}
//	sku := mapping.ExtractGCPSKU(props)       // "n1-standard-4"
//	region := mapping.ExtractGCPRegion(props) // "us-central1"
//
// Generic Example:
//
//	props := map[string]string{"customSKU": "my-sku"}
//	sku := mapping.ExtractSKU(props, "customSKU", "fallback") // "my-sku"
//
// # Error Handling
//
// Functions in this package are designed to never panic and always return
// safe default values (empty string or false) for invalid input:
//
//	mapping.ExtractAWSSKU(nil)                 // ""
//	mapping.ExtractAWSSKU(map[string]string{}) // ""
//	mapping.ExtractAWSRegionFromAZ("")         // ""
//
// # Performance
//
// All functions are optimized for minimal allocation (<50 ns/op, 0 allocs/op target):
//   - No regex compilation at runtime
//   - Simple string operations only
//   - Package-level constants for key names
package mapping
