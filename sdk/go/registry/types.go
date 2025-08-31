package registry

import (
	"time"
)

// PluginType represents the type of plugin functionality.
type PluginType string

const (
	PluginTypeCostSource     PluginType = "cost_source"
	PluginTypeObservability  PluginType = "observability"
	PluginTypeRegistry       PluginType = "registry"
	PluginTypeAggregator     PluginType = "aggregator"
)

// Capability represents plugin capabilities.
type Capability string

const (
	CapabilityActualCost         Capability = "actual_cost"
	CapabilityProjectedCost      Capability = "projected_cost"
	CapabilityPricingSpec        Capability = "pricing_spec"
	CapabilityRealTimeCost       Capability = "real_time_cost"
	CapabilityHistoricalCost     Capability = "historical_cost"
	CapabilityCostForecasting    Capability = "cost_forecasting"
	CapabilityCostOptimization   Capability = "cost_optimization"
	CapabilityBudgetAlerts       Capability = "budget_alerts"
	CapabilityCustomMetrics      Capability = "custom_metrics"
	CapabilityMultiCloud         Capability = "multi_cloud"
	CapabilityKubernetesNative   Capability = "kubernetes_native"
	CapabilityTaggingSupport     Capability = "tagging_support"
	CapabilityDrillDown          Capability = "drill_down"
	CapabilityCostAllocation     Capability = "cost_allocation"
	CapabilityShowback           Capability = "showback"
	CapabilityChargeback         Capability = "chargeback"
)

// AuthenticationMethod represents supported authentication methods.
type AuthenticationMethod string

const (
	AuthMethodAPIKey         AuthenticationMethod = "api_key"
	AuthMethodOAuth2         AuthenticationMethod = "oauth2"
	AuthMethodServiceAccount AuthenticationMethod = "service_account"
	AuthMethodIAMRole        AuthenticationMethod = "iam_role"
	AuthMethodMutualTLS      AuthenticationMethod = "mutual_tls"
	AuthMethodBasicAuth      AuthenticationMethod = "basic_auth"
	AuthMethodBearerToken    AuthenticationMethod = "bearer_token"
)

// SecurityStatus represents the security status of a plugin.
type SecurityStatus string

const (
	SecurityStatusSecure     SecurityStatus = "secure"
	SecurityStatusAdvisory   SecurityStatus = "advisory"
	SecurityStatusVulnerable SecurityStatus = "vulnerable"
	SecurityStatusUnknown    SecurityStatus = "unknown"
)

// Severity represents the severity level of issues.
type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityModerate Severity = "moderate"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

// DependencyType represents the type of dependency.
type DependencyType string

const (
	DependencyTypeBinary  DependencyType = "binary"
	DependencyTypeLibrary DependencyType = "library"
	DependencyTypeService DependencyType = "service"
	DependencyTypePlugin  DependencyType = "plugin"
)

// PluginManifest contains complete plugin manifest information.
type PluginManifest struct {
	Name                string                `json:"name"`
	DisplayName         string                `json:"display_name"`
	Description         string                `json:"description"`
	Version             string                `json:"version"`
	APIVersion          string                `json:"api_version"`
	PluginType          PluginType            `json:"plugin_type"`
	Capabilities        []Capability          `json:"capabilities"`
	SupportedProviders  []string              `json:"supported_providers"`
	SupportedRegions    []string              `json:"supported_regions"`
	Requirements        PluginRequirements    `json:"requirements"`
	Authentication      AuthenticationConfig  `json:"authentication"`
	Installation        InstallationInfo      `json:"installation"`
	Configuration       ConfigurationSchema   `json:"configuration"`
	Contacts            ContactInfo           `json:"contacts"`
	Metadata            PluginMetadata        `json:"metadata"`
}

// PluginInfo contains basic plugin information for listings.
type PluginInfo struct {
	Name               string         `json:"name"`
	DisplayName        string         `json:"display_name"`
	Description        string         `json:"description"`
	LatestVersion      string         `json:"latest_version"`
	PluginType         PluginType     `json:"plugin_type"`
	Capabilities       []Capability   `json:"capabilities"`
	SupportedProviders []string       `json:"supported_providers"`
	Downloads          DownloadStats  `json:"downloads"`
	Rating             Rating         `json:"rating"`
	SecurityStatus     SecurityStatus `json:"security_status"`
	Verified           bool           `json:"verified"`
	Tags               []string       `json:"tags"`
	LastUpdated        time.Time      `json:"last_updated"`
}

// PluginRequirements contains plugin runtime requirements.
type PluginRequirements struct {
	MinAPIVersion        string                `json:"min_api_version"`
	MaxAPIVersion        string                `json:"max_api_version"`
	OperatingSystems     []string              `json:"os"`
	Architectures        []string              `json:"arch"`
	Dependencies         []Dependency          `json:"dependencies"`
	EnvironmentVariables []EnvironmentVariable `json:"environment_variables"`
}

// Dependency represents an external dependency.
type Dependency struct {
	Name     string         `json:"name"`
	Version  string         `json:"version"`
	Type     DependencyType `json:"type"`
	Optional bool           `json:"optional"`
}

// EnvironmentVariable represents a required environment variable.
type EnvironmentVariable struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Required     bool   `json:"required"`
	DefaultValue string `json:"default_value"`
}

// AuthenticationConfig contains authentication and security information.
type AuthenticationConfig struct {
	Required    bool                     `json:"required"`
	Methods     []AuthenticationMethod   `json:"methods"`
	Scopes      []string                 `json:"scopes"`
	Permissions []string                 `json:"permissions"`
}

// InstallationInfo contains plugin installation information.
type InstallationInfo struct {
	BinaryURL     string `json:"binary_url"`
	Checksum      string `json:"checksum"`
	SizeBytes     int64  `json:"size_bytes"`
	InstallScript string `json:"install_script"`
	DockerImage   string `json:"docker_image"`
}

// ConfigurationSchema contains plugin configuration information.
type ConfigurationSchema struct {
	SchemaURL string                 `json:"schema_url"`
	Defaults  map[string]interface{} `json:"defaults"`
	Examples  []map[string]interface{} `json:"examples"`
}

// ContactInfo contains plugin contact information.
type ContactInfo struct {
	Maintainers      []Maintainer `json:"maintainers"`
	SupportURL       string       `json:"support_url"`
	DocumentationURL string       `json:"documentation_url"`
}

// Maintainer represents a plugin maintainer.
type Maintainer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	URL   string `json:"url"`
}

// PluginMetadata contains additional plugin metadata.
type PluginMetadata struct {
	Tags            []string               `json:"tags"`
	License         string                 `json:"license"`
	Homepage        string                 `json:"homepage"`
	Repository      string                 `json:"repository"`
	ReleaseNotesURL string                 `json:"release_notes_url"`
	Created         time.Time              `json:"created"`
	Updated         time.Time              `json:"updated"`
	Custom          map[string]interface{} `json:"custom"`
}

// DownloadStats contains plugin download statistics.
type DownloadStats struct {
	Total       int64 `json:"total"`
	Last30Days  int64 `json:"last_30_days"`
}

// Rating contains plugin rating information.
type Rating struct {
	Average float64 `json:"average"`
	Count   int32   `json:"count"`
}

// PluginRegistry contains registry metadata and plugin listings.
type PluginRegistry struct {
	RegistryVersion  string           `json:"registry_version"`
	RegistryMetadata RegistryMetadata `json:"registry_metadata"`
	Plugins          []PluginEntry    `json:"plugins"`
	Categories       []Category       `json:"categories"`
}

// RegistryMetadata contains metadata about a registry instance.
type RegistryMetadata struct {
	Name                   string      `json:"name"`
	Description            string      `json:"description"`
	URL                    string      `json:"url"`
	Maintainer             Maintainer  `json:"maintainer"`
	LastUpdated            time.Time   `json:"last_updated"`
	SupportedAPIVersions   []string    `json:"supported_api_versions"`
}

// PluginEntry represents a plugin entry in a registry.
type PluginEntry struct {
	Name               string                `json:"name"`
	DisplayName        string                `json:"display_name"`
	Description        string                `json:"description"`
	LatestVersion      string                `json:"latest_version"`
	PluginType         PluginType            `json:"plugin_type"`
	SupportedProviders []string              `json:"supported_providers"`
	Capabilities       []Capability          `json:"capabilities"`
	Versions           []PluginVersion       `json:"versions"`
	Downloads          DownloadStats         `json:"downloads"`
	Rating             Rating                `json:"rating"`
	SecurityStatus     PluginSecurityStatus  `json:"security_status"`
	Verified           bool                  `json:"verified"`
	Tags               []string              `json:"tags"`
}

// PluginVersion represents a specific version of a plugin.
type PluginVersion struct {
	Version           string                `json:"version"`
	ManifestURL       string                `json:"manifest_url"`
	Published         time.Time             `json:"published"`
	Deprecated        bool                  `json:"deprecated"`
	SecurityAdvisory  *SecurityAdvisory     `json:"security_advisory,omitempty"`
	Yanked            bool                  `json:"yanked"`
	YankReason        string                `json:"yank_reason,omitempty"`
	APICompatibility  APICompatibility      `json:"api_compatibility"`
}

// SecurityAdvisory contains security advisory information.
type SecurityAdvisory struct {
	Severity Severity `json:"severity"`
	Summary  string   `json:"summary"`
	URL      string   `json:"url"`
}

// APICompatibility contains API compatibility information.
type APICompatibility struct {
	MinAPIVersion string `json:"min_api_version"`
	MaxAPIVersion string `json:"max_api_version"`
}

// PluginSecurityStatus contains overall security status of a plugin.
type PluginSecurityStatus struct {
	Status           SecurityStatus `json:"status"`
	LastAudit        time.Time      `json:"last_audit"`
	AdvisoriesCount  int            `json:"advisories_count"`
}

// Category represents a plugin category for organization.
type Category struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Plugins     []string `json:"plugins"`
}

// PluginFilter contains criteria for filtering plugin searches.
type PluginFilter struct {
	PluginType     PluginType     `json:"plugin_type,omitempty"`
	Providers      []string       `json:"providers,omitempty"`
	Capabilities   []Capability   `json:"capabilities,omitempty"`
	Tags           []string       `json:"tags,omitempty"`
	Verified       *bool          `json:"verified,omitempty"`
	MinRating      *float64       `json:"min_rating,omitempty"`
	APIVersion     string         `json:"api_version,omitempty"`
	SecurityStatus SecurityStatus `json:"security_status,omitempty"`
}

// ValidationError represents a plugin validation error.
type ValidationError struct {
	Code     string   `json:"code"`
	Field    string   `json:"field"`
	Message  string   `json:"message"`
	Severity Severity `json:"severity"`
}

// ValidationWarning represents a plugin validation warning.
type ValidationWarning struct {
	Code    string `json:"code"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationResult contains validation results.
type ValidationResult struct {
	Valid    bool                 `json:"valid"`
	Errors   []ValidationError    `json:"errors"`
	Warnings []ValidationWarning  `json:"warnings"`
}

// SecurityScanResults contains security scan results.
type SecurityScanResults struct {
	Status          SecurityStatus          `json:"status"`
	Vulnerabilities []SecurityVulnerability `json:"vulnerabilities"`
	LastScan        time.Time               `json:"last_scan"`
	ScanVersion     string                  `json:"scan_version"`
}

// SecurityVulnerability represents a security vulnerability.
type SecurityVulnerability struct {
	ID              string   `json:"id"`
	Severity        Severity `json:"severity"`
	Summary         string   `json:"summary"`
	Description     string   `json:"description"`
	URL             string   `json:"url"`
	FixedInVersion  string   `json:"fixed_in_version"`
}

// InstallOptions contains options for plugin installation.
type InstallOptions struct {
	RegistryURL     string            `json:"registry_url"`
	VerifySignature bool              `json:"verify_signature"`
	Force           bool              `json:"force"`
	Config          map[string]string `json:"config"`
}

// UpdateOptions contains options for plugin updates.
type UpdateOptions struct {
	Version string `json:"version"`
	Force   bool   `json:"force"`
}