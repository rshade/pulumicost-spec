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

// Package pricing provides domain types and validation for PulumiCost pricing specifications.
// It includes billing mode constants, unit types, and validation helpers for ensuring
// pricing data conforms to the PulumiCost schema.
package pricing

// BillingMode represents the billing model for a cloud resource.
// It defines how the resource is charged (e.g., per hour, per GB-month, etc.).
type BillingMode string

// Time-based billing modes.
const (
	PerHour   BillingMode = "per_hour"
	PerMinute BillingMode = "per_minute"
	PerSecond BillingMode = "per_second"
	PerDay    BillingMode = "per_day"
	PerMonth  BillingMode = "per_month"
	PerYear   BillingMode = "per_year"
)

// Storage-based billing modes.
const (
	PerGBMonth BillingMode = "per_gb_month"
	PerGBHour  BillingMode = "per_gb_hour"
	PerGBDay   BillingMode = "per_gb_day"
)

// Usage-based billing modes.
const (
	PerRequest     BillingMode = "per_request"
	PerOperation   BillingMode = "per_operation"
	PerTransaction BillingMode = "per_transaction"
	PerExecution   BillingMode = "per_execution"
	PerInvocation  BillingMode = "per_invocation"
	PerAPICall     BillingMode = "per_api_call"
	PerLookup      BillingMode = "per_lookup"
	PerQuery       BillingMode = "per_query"
)

// Compute-based billing modes.
const (
	PerCPUHour       BillingMode = "per_cpu_hour"
	PerCPUMonth      BillingMode = "per_cpu_month"
	PerVCPUHour      BillingMode = "per_vcpu_hour"
	PerMemoryGBHour  BillingMode = "per_memory_gb_hour"
	PerMemoryGBMonth BillingMode = "per_memory_gb_month"
)

// I/O-based billing modes.
const (
	PerIOPS            BillingMode = "per_iops"
	PerProvisionedIOPS BillingMode = "per_provisioned_iops"
	PerDataTransferGB  BillingMode = "per_data_transfer_gb"
	PerBandwidthGB     BillingMode = "per_bandwidth_gb"
)

// Database-specific billing modes.
const (
	PerRCU BillingMode = "per_rcu" // DynamoDB Read Capacity Units
	PerWCU BillingMode = "per_wcu" // DynamoDB Write Capacity Units
	PerDTU BillingMode = "per_dtu" // Azure Database Transaction Units
	PerRU  BillingMode = "per_ru"  // Azure Cosmos DB Request Units
)

// Pricing model types.
const (
	OnDemand       BillingMode = "on_demand"
	Reserved       BillingMode = "reserved"
	Spot           BillingMode = "spot"
	Preemptible    BillingMode = "preemptible"
	SavingsPlan    BillingMode = "savings_plan"
	CommittedUse   BillingMode = "committed_use"
	HybridBenefit  BillingMode = "hybrid_benefit"
	FlatRate       BillingMode = "flat"
	Tiered         BillingMode = "tiered"
	NotImplemented BillingMode = "not_implemented"
)

// Unit represents the unit of measurement for pricing.
type Unit string

// Unit constants for pricing specifications.
const (
	UnitHour    Unit = "hour"
	UnitGBMonth Unit = "GB-month"
	UnitRequest Unit = "request"
	UnitUnknown Unit = "unknown"
	UnitDTU     Unit = "DTU"
	UnitRCU     Unit = "RCU"
	UnitWCU     Unit = "WCU"
	UnitRU      Unit = "RU"
)

// String returns the unit as its string value.
func (u Unit) String() string { return string(u) }

// getAllBillingModes returns all valid billing modes for validation.
func getAllBillingModes() []BillingMode {
	return []BillingMode{
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
		Tiered, NotImplemented,
	}
}

// String returns the billing mode as its string value.
func (b BillingMode) String() string { return string(b) }

// ValidBillingMode checks if the given string represents a valid billing mode.
func ValidBillingMode(s string) bool {
	mode := BillingMode(s)
	for _, validMode := range getAllBillingModes() {
		if mode == validMode {
			return true
		}
	}
	return false
}

// IsValidBillingMode checks if a billing mode string is valid.
func IsValidBillingMode(s string) bool {
	return ValidBillingMode(s)
}

// GetAllBillingModes returns all available billing modes as strings.
func GetAllBillingModes() []string {
	allModes := getAllBillingModes()
	modes := make([]string, len(allModes))
	for i, mode := range allModes {
		modes[i] = mode.String()
	}
	return modes
}

// Provider enumeration for validation.
type Provider string

const (
	AWS        Provider = "aws"
	Azure      Provider = "azure"
	GCP        Provider = "gcp"
	Kubernetes Provider = "kubernetes"
	Custom     Provider = "custom"
)

// GetAllProviders returns all valid providers.
func GetAllProviders() []Provider {
	return []Provider{AWS, Azure, GCP, Kubernetes, Custom}
}

// String returns the provider name as its string value.
func (p Provider) String() string { return string(p) }

// ValidProvider checks if the given string represents a valid cloud provider.
func ValidProvider(s string) bool {
	provider := Provider(s)
	for _, validProvider := range GetAllProviders() {
		if provider == validProvider {
			return true
		}
	}
	return false
}
