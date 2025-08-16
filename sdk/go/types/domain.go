package types

type BillingMode string

// Time-based billing modes
const (
	PerHour   BillingMode = "per_hour"
	PerMinute BillingMode = "per_minute"
	PerSecond BillingMode = "per_second"
	PerDay    BillingMode = "per_day"
	PerMonth  BillingMode = "per_month"
	PerYear   BillingMode = "per_year"
)

// Storage-based billing modes
const (
	PerGBMonth BillingMode = "per_gb_month"
	PerGBHour  BillingMode = "per_gb_hour"
	PerGBDay   BillingMode = "per_gb_day"
)

// Usage-based billing modes
const (
	PerRequest    BillingMode = "per_request"
	PerOperation  BillingMode = "per_operation"
	PerTransaction BillingMode = "per_transaction"
	PerExecution  BillingMode = "per_execution"
	PerInvocation BillingMode = "per_invocation"
	PerAPICall    BillingMode = "per_api_call"
	PerLookup     BillingMode = "per_lookup"
	PerQuery      BillingMode = "per_query"
)

// Compute-based billing modes
const (
	PerCPUHour      BillingMode = "per_cpu_hour"
	PerCPUMonth     BillingMode = "per_cpu_month"
	PerVCPUHour     BillingMode = "per_vcpu_hour"
	PerMemoryGBHour BillingMode = "per_memory_gb_hour"
	PerMemoryGBMonth BillingMode = "per_memory_gb_month"
)

// I/O-based billing modes
const (
	PerIOPS             BillingMode = "per_iops"
	PerProvisionedIOPS  BillingMode = "per_provisioned_iops"
	PerDataTransferGB   BillingMode = "per_data_transfer_gb"
	PerBandwidthGB      BillingMode = "per_bandwidth_gb"
)

// Database-specific billing modes
const (
	PerRCU BillingMode = "per_rcu" // DynamoDB Read Capacity Units
	PerWCU BillingMode = "per_wcu" // DynamoDB Write Capacity Units
	PerDTU BillingMode = "per_dtu" // Azure Database Transaction Units
	PerRU  BillingMode = "per_ru"  // Azure Cosmos DB Request Units
)

// Pricing model types
const (
	OnDemand     BillingMode = "on_demand"
	Reserved     BillingMode = "reserved"
	Spot         BillingMode = "spot"
	Preemptible  BillingMode = "preemptible"
	SavingsPlan  BillingMode = "savings_plan"
	CommittedUse BillingMode = "committed_use"
	HybridBenefit BillingMode = "hybrid_benefit"
	FlatRate     BillingMode = "flat"
)

// AllBillingModes contains all valid billing modes for validation
var AllBillingModes = []BillingMode{
	// Time-based
	PerHour, PerMinute, PerSecond, PerDay, PerMonth, PerYear,
	// Storage-based
	PerGBMonth, PerGBHour, PerGBDay,
	// Usage-based
	PerRequest, PerOperation, PerTransaction, PerExecution, PerInvocation,
	PerAPICall, PerLookup, PerQuery,
	// Compute-based
	PerCPUHour, PerCPUMonth, PerVCPUHour, PerMemoryGBHour, PerMemoryGBMonth,
	// I/O-based
	PerIOPS, PerProvisionedIOPS, PerDataTransferGB, PerBandwidthGB,
	// Database-specific
	PerRCU, PerWCU, PerDTU, PerRU,
	// Pricing models
	OnDemand, Reserved, Spot, Preemptible, SavingsPlan, CommittedUse, HybridBenefit, FlatRate,
}

func (b BillingMode) String() string { return string(b) }

func ValidBillingMode(s string) bool {
	mode := BillingMode(s)
	for _, validMode := range AllBillingModes {
		if mode == validMode {
			return true
		}
	}
	return false
}

// Provider enumeration for validation
type Provider string

const (
	AWS        Provider = "aws"
	Azure      Provider = "azure"
	GCP        Provider = "gcp"
	Kubernetes Provider = "kubernetes"
	Custom     Provider = "custom"
)

var AllProviders = []Provider{AWS, Azure, GCP, Kubernetes, Custom}

func (p Provider) String() string { return string(p) }

func ValidProvider(s string) bool {
	provider := Provider(s)
	for _, validProvider := range AllProviders {
		if provider == validProvider {
			return true
		}
	}
	return false
}
