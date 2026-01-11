package mapping

// AWS property key constants for SKU extraction.
// These keys are checked in priority order by ExtractAWSSKU.
const (
	// AWSKeyInstanceType is the property key for EC2 instance types.
	AWSKeyInstanceType = "instanceType"
	// AWSKeyInstanceClass is the property key for RDS instance classes.
	AWSKeyInstanceClass = "instanceClass"
	// AWSKeyType is a generic type property key.
	AWSKeyType = "type"
	// AWSKeyVolumeType is the property key for EBS volume types.
	AWSKeyVolumeType = "volumeType"
)

// AWS property key constants for region extraction.
// These keys are checked in priority order by ExtractAWSRegion.
const (
	// AWSKeyRegion is the explicit region property key.
	AWSKeyRegion = "region"
	// AWSKeyAvailabilityZone is the availability zone property key.
	// Region is derived by removing the trailing letter.
	AWSKeyAvailabilityZone = "availabilityZone"
)

// Azure property key constants for SKU extraction.
// These keys are checked in priority order by ExtractAzureSKU.
const (
	// AzureKeyVMSize is the property key for Azure VM sizes.
	AzureKeyVMSize = "vmSize"
	// AzureKeySKU is the generic SKU property key.
	AzureKeySKU = "sku"
	// AzureKeyTier is the service tier property key.
	AzureKeyTier = "tier"
)

// Azure property key constants for region extraction.
// These keys are checked in priority order by ExtractAzureRegion.
const (
	// AzureKeyLocation is the primary Azure location property key.
	AzureKeyLocation = "location"
	// AzureKeyRegion is an alternative region property key.
	AzureKeyRegion = "region"
)

// GCP property key constants for SKU extraction.
// These keys are checked in priority order by ExtractGCPSKU.
const (
	// GCPKeyMachineType is the property key for GCP machine types.
	GCPKeyMachineType = "machineType"
	// GCPKeyType is a generic type property key.
	GCPKeyType = "type"
	// GCPKeyTier is the service tier property key.
	GCPKeyTier = "tier"
)

// GCP property key constants for region extraction.
// These keys are checked in priority order by ExtractGCPRegion.
const (
	// GCPKeyRegion is the explicit region property key.
	GCPKeyRegion = "region"
	// GCPKeyZone is the zone property key.
	// Region is derived by removing the last hyphen-delimited segment.
	GCPKeyZone = "zone"
)

// Default property keys for generic SKU extraction.
// Used by ExtractSKU when no custom keys are provided.
//
//nolint:gochecknoglobals // Intentional: read-only default configuration
var defaultSKUKeys = []string{"sku", "type", "tier"}

// Default property keys for generic region extraction.
// Used by ExtractRegion when no custom keys are provided.
//
//nolint:gochecknoglobals // Intentional: read-only default configuration
var defaultRegionKeys = []string{"region", "location", "zone"}
