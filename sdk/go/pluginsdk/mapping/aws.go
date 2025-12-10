package mapping

// ExtractAWSRegionFromAZ derives the AWS region from a standard availability zone string.
//
// This function supports standard AWS availability zones of the form {region}{letter}
// (e.g., "us-east-1a"). It does NOT handle extended zone formats such as Local Zones
// (e.g., "us-west-2-lax-1a") or Wavelength Zones which have different naming patterns.
//
// Algorithm: Removes the trailing lowercase letter from the zone name.
//
// Examples:
//   - "us-east-1a" → "us-east-1"
//   - "eu-west-2b" → "eu-west-2"
//   - "ap-northeast-1c" → "ap-northeast-1"
//
// Returns empty string if input is empty or too short to be a valid AZ.
// Returns input as-is if no trailing letter suffix is found (assumes it's already a region).
// Never panics.
func ExtractAWSRegionFromAZ(availabilityZone string) string {
	if availabilityZone == "" {
		return ""
	}

	// AWS AZs end with a single lowercase letter (a-z)
	// Region is everything except the trailing letter
	length := len(availabilityZone)

	// Single character is not a valid AZ or region
	if length == 1 {
		return ""
	}

	// If the last character is a lowercase letter and length > 1, remove it
	lastChar := availabilityZone[length-1]
	if lastChar >= 'a' && lastChar <= 'z' {
		return availabilityZone[:length-1]
	}

	// Return as-is if no trailing letter (might already be a region)
	return availabilityZone
}

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
func ExtractAWSSKU(properties map[string]string) string {
	return extractFromKeys(properties, AWSKeyInstanceType, AWSKeyInstanceClass, AWSKeyType, AWSKeyVolumeType)
}

// ExtractAWSRegion extracts the region from AWS resource properties.
//
// Key priority order:
//  1. region - Explicit region setting
//  2. availabilityZone - Derived via ExtractAWSRegionFromAZ
//
// Returns empty string if no key found or input is nil/empty.
// Never panics.
func ExtractAWSRegion(properties map[string]string) string {
	if properties == nil {
		return ""
	}

	// Check explicit region first
	if region := properties[AWSKeyRegion]; region != "" {
		return region
	}

	// Try to derive from availability zone
	if az := properties[AWSKeyAvailabilityZone]; az != "" {
		return ExtractAWSRegionFromAZ(az)
	}

	return ""
}
