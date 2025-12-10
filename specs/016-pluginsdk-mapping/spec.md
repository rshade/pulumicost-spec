# Feature Specification: PluginSDK Mapping Package

**Feature Branch**: `016-pluginsdk-mapping`
**Created**: 2025-12-09
**Status**: Draft
**Input**: User description: "Create a pluginsdk/mapping/ package that provides shared helper
functions for extracting SKU, region, and other pricing-relevant fields from Pulumi resource
properties."

## Clarifications

### Session 2025-12-09

- Q: GCP zone-to-region extraction algorithm? → A: Remove last hyphen segment and validate
  against known GCP regions list

## User Scenarios & Testing _(mandatory)_

### User Story 1 - AWS Plugin Developer Extracts Resource Properties (Priority: P1)

A plugin developer building an AWS cost estimation plugin needs to extract SKU and region
information from EC2 instance properties. The developer imports the mapping package and uses
AWS-specific helper functions to consistently extract `instanceType` as the SKU and derive
the region from `availabilityZone`.

**Why this priority**: AWS is the most common cloud provider, and EC2/EBS are fundamental
resource types. Enabling AWS property extraction provides immediate value for the majority
of plugin developers.

**Independent Test**: Can be fully tested by passing a map of AWS EC2 properties and verifying
correct SKU and region extraction. Delivers value as a standalone AWS extraction utility.

**Acceptance Scenarios**:

1. **Given** a map containing `{"instanceType": "t3.medium", "availabilityZone": "us-east-1a"}`,
   **When** `ExtractAWSSKU` is called, **Then** return `"t3.medium"`
2. **Given** a map containing `{"availabilityZone": "us-west-2b"}`,
   **When** `ExtractAWSRegion` is called, **Then** return `"us-west-2"`
3. **Given** a map containing `{"volumeType": "gp3"}`,
   **When** `ExtractAWSSKU` is called for EBS volumes, **Then** return `"gp3"`
4. **Given** a map containing `{"instanceClass": "db.t3.micro"}`,
   **When** `ExtractAWSSKU` is called for RDS, **Then** return `"db.t3.micro"`

---

### User Story 2 - Azure Plugin Developer Extracts VM Properties (Priority: P2)

A plugin developer building an Azure cost estimation plugin needs to extract SKU and region
information from Azure VM properties. The developer uses Azure-specific helper functions to
extract `vmSize` as the SKU and the location/region directly.

**Why this priority**: Azure is the second most common cloud provider. Supporting Azure
extraction enables multi-cloud scenarios and expands the plugin developer audience.

**Independent Test**: Can be fully tested by passing a map of Azure VM properties and verifying
correct SKU and region extraction.

**Acceptance Scenarios**:

1. **Given** a map containing `{"vmSize": "Standard_D2s_v3", "location": "eastus"}`,
   **When** `ExtractAzureSKU` is called, **Then** return `"Standard_D2s_v3"`
2. **Given** a map containing `{"location": "westeurope"}`,
   **When** `ExtractAzureRegion` is called, **Then** return `"westeurope"`

---

### User Story 3 - GCP Plugin Developer Extracts Compute Properties (Priority: P2)

A plugin developer building a GCP cost estimation plugin needs to extract SKU and region
information from GCP Compute Engine properties. The developer uses GCP-specific helper
functions to extract `machineType` as the SKU and derive the region from zone or direct
region property.

**Why this priority**: GCP is the third major cloud provider. Supporting GCP completes the
multi-cloud coverage for plugin developers.

**Independent Test**: Can be fully tested by passing a map of GCP Compute properties and
verifying correct SKU and region extraction.

**Acceptance Scenarios**:

1. **Given** a map containing `{"machineType": "n1-standard-4", "zone": "us-central1-a"}`,
   **When** `ExtractGCPSKU` is called, **Then** return `"n1-standard-4"`
2. **Given** a map containing `{"zone": "europe-west1-b"}`,
   **When** `ExtractGCPRegion` is called, **Then** return `"europe-west1"`
3. **Given** a map containing `{"region": "asia-east1"}`,
   **When** `ExtractGCPRegion` is called, **Then** return `"asia-east1"`

---

### User Story 4 - FinOps Plugin Developer Uses Generic Extractors (Priority: P3)

A plugin developer building a FinOps tool integration (e.g., Kubecost, Vantage) needs to
extract properties from multiple cloud providers or custom resource types. The developer
uses generic fallback extractors that check multiple possible property keys.

**Why this priority**: Generic extractors provide flexibility for edge cases and multi-cloud
FinOps tools, but most plugin developers will use provider-specific functions first.

**Independent Test**: Can be fully tested by passing maps with various property key names and
verifying the generic extractor finds values using fallback key lists.

**Acceptance Scenarios**:

1. **Given** a map containing `{"sku": "custom-sku-value"}`,
   **When** `ExtractSKU` is called with default keys, **Then** return `"custom-sku-value"`
2. **Given** a map containing `{"region": "custom-region"}`,
   **When** `ExtractRegion` is called with default keys, **Then** return `"custom-region"`
3. **Given** a map containing `{"customField": "value"}`,
   **When** `ExtractSKU` is called with custom keys `["customField"]`, **Then** return `"value"`

---

### User Story 5 - Core System Decoupled from Cloud-Specific Logic (Priority: P1)

The pulumicost-core maintainer needs to remove cloud-specific extraction logic from the core
adapter to keep the core system cloud-agnostic. By importing the mapping package from
pulumicost-spec, the core delegates all property extraction to the shared SDK.

**Why this priority**: Decoupling is a primary architectural goal. Without this, the core
remains tightly coupled to cloud-specific knowledge, making maintenance difficult.

**Independent Test**: Can be verified by removing cloud-specific code from core adapter and
replacing with mapping package imports while maintaining existing functionality.

**Acceptance Scenarios**:

1. **Given** the mapping package is published,
   **When** pulumicost-core imports and uses mapping functions,
   **Then** all existing extraction tests continue to pass
2. **Given** a new cloud property naming convention is discovered,
   **When** the mapping package is updated,
   **Then** pulumicost-core receives the fix without code changes

---

### Edge Cases

- What happens when a property map is empty? Return empty string, no panic.
- What happens when none of the expected keys exist in the map? Return empty string.
- What happens when availability zone format is unexpected (e.g., missing letter suffix)?
  Return the input as-is for region extraction.
- How does system handle nil map input? Return empty string, no panic.
- What happens when a property value is an empty string? Return empty string (treat as
  not found).
- What happens when GCP zone produces invalid region after extraction? Return empty string
  if derived region is not in known GCP regions list.

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: System MUST provide `ExtractAWSSKU(properties map[string]string) string` function
  that extracts SKU from AWS resource properties by checking keys: `instanceType`,
  `instanceClass`, `type`, `volumeType` in priority order
- **FR-002**: System MUST provide `ExtractAWSRegion(properties map[string]string) string`
  function that extracts region from AWS resource properties by checking keys: `region`,
  `availabilityZone` and deriving region from availability zone
- **FR-003**: System MUST provide `ExtractAWSRegionFromAZ(availabilityZone string) string`
  function that derives region by removing the trailing letter from availability zone
  (e.g., `us-east-1a` -> `us-east-1`)
- **FR-004**: System MUST provide `ExtractAzureSKU(properties map[string]string) string`
  function that extracts SKU from Azure resource properties by checking keys: `vmSize`,
  `sku`, `tier`
- **FR-005**: System MUST provide `ExtractAzureRegion(properties map[string]string) string`
  function that extracts region from Azure resource properties by checking keys: `location`,
  `region`
- **FR-006**: System MUST provide `ExtractGCPSKU(properties map[string]string) string`
  function that extracts SKU from GCP resource properties by checking keys: `machineType`,
  `type`, `tier`
- **FR-007**: System MUST provide `ExtractGCPRegion(properties map[string]string) string`
  function that extracts region from GCP resource properties by checking keys: `region`,
  `zone` and deriving region from zone by removing the last hyphen-delimited segment
  (e.g., `us-central1-a` → `us-central1`) and validating against a known GCP regions list
- **FR-008**: System MUST provide `ExtractSKU(properties map[string]string, keys ...string)`
  generic function that checks provided keys in order and returns first non-empty value;
  if no keys provided, uses default keys: `sku`, `type`, `tier`
- **FR-009**: System MUST provide `ExtractRegion(properties map[string]string, keys ...string)`
  generic function that checks provided keys in order and returns first non-empty value;
  if no keys provided, uses default keys: `region`, `location`, `zone`
- **FR-010**: All extraction functions MUST return empty string when no matching key is found
  or map is nil/empty
- **FR-011**: All extraction functions MUST NOT panic on nil or empty input
- **FR-012**: System MUST be importable as a Go package at
  `github.com/rshade/pulumicost-spec/sdk/go/pluginsdk/mapping`
- **FR-013**: System MUST provide a known GCP regions list for validation during zone-to-region
  extraction; if derived region is not in list, return empty string
- **FR-014**: System MUST provide `ExtractGCPRegionFromZone(zone string) string` function
  that derives region by removing the last hyphen-delimited segment from a zone string
  (e.g., `us-central1-a` → `us-central1`) and validates against the known GCP regions list;
  returns empty string if derived region is invalid
- **FR-015**: System MUST provide `IsValidGCPRegion(region string) bool` function that
  returns true if the provided region exists in the known GCP regions list

### Key Entities

- **PropertyMap**: A `map[string]string` representing Pulumi resource properties where keys
  are property names and values are property values
- **SKU**: A string representing the pricing-relevant resource type identifier (e.g., instance
  type, VM size, machine type)
- **Region**: A string representing the cloud region where the resource is deployed

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Plugin developers can extract SKU and region from AWS, Azure, and GCP resources
  using 2 function calls or fewer
- **SC-002**: All extraction functions achieve 90%+ code coverage through unit tests
- **SC-003**: Package provides complete documentation with usage examples for each cloud
  provider
- **SC-004**: Existing pulumicost-core extraction logic can be fully replaced by mapping
  package functions with no behavior change
- **SC-005**: Zero panics occur when handling edge cases (nil input, empty maps, missing keys)
- **SC-006**: Plugin developers can discover correct extraction function within 30 seconds
  through package documentation

## Assumptions

- Property maps use string keys and string values (consistent with Pulumi resource property
  representation)
- AWS availability zones follow the standard format of `{region}{letter}` (e.g., `us-east-1a`)
- GCP zones follow the standard format of `{region}-{letter}` (e.g., `us-central1-a`)
- Azure locations are direct region identifiers without transformation needed
- The mapping package will be versioned and released alongside other pulumicost-spec SDK
  components
