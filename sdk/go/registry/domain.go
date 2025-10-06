// Package registry provides domain types and validation for PulumiCost plugin registry management.
// It defines enums for discovery sources, plugin statuses, security levels, installation methods,
// capabilities, system permissions, authentication methods, and validation functions for plugin
// manifests and names.
package registry

import (
	"errors"
	"fmt"
)

// Provider represents supported cloud providers.
type Provider string

const (
	// ProviderAWS indicates Amazon Web Services provider.
	ProviderAWS Provider = "aws"
	// ProviderAzure indicates Microsoft Azure provider.
	ProviderAzure Provider = "azure"
	// ProviderGCP indicates Google Cloud Platform provider.
	ProviderGCP Provider = "gcp"
	// ProviderKubernetes indicates Kubernetes provider.
	ProviderKubernetes Provider = "kubernetes"
	// ProviderCustom indicates custom provider implementation.
	ProviderCustom Provider = "custom"
)

// validProviders is a map for O(1) provider validation lookups.
var validProviders = map[Provider]bool{
	ProviderAWS:        true,
	ProviderAzure:      true,
	ProviderGCP:        true,
	ProviderKubernetes: true,
	ProviderCustom:     true,
}

// AllProviders returns all valid providers.
func AllProviders() []Provider {
	return []Provider{
		ProviderAWS, ProviderAzure, ProviderGCP, ProviderKubernetes, ProviderCustom,
	}
}

// IsValidProvider checks if a provider string is valid.
func IsValidProvider(p string) bool {
	return validProviders[Provider(p)]
}

const (
	// MinPluginNameLength defines the minimum required length for plugin names.
	MinPluginNameLength = 2
	// MaxPluginNameLength defines the maximum allowed length for plugin names.
	MaxPluginNameLength = 50
)

// DiscoverySource represents the supported plugin discovery sources.
type DiscoverySource string

const (
	// DiscoverySourceFilesystem indicates filesystem-based plugin discovery.
	DiscoverySourceFilesystem DiscoverySource = "filesystem"
	// DiscoverySourceRegistry indicates registry-based plugin discovery.
	DiscoverySourceRegistry DiscoverySource = "registry"
	// DiscoverySourceURL indicates URL-based plugin discovery.
	DiscoverySourceURL DiscoverySource = "url"
	// DiscoverySourceGit indicates git repository-based plugin discovery.
	DiscoverySourceGit DiscoverySource = "git"
)

// validDiscoverySources is a map for O(1) discovery source validation lookups.
var validDiscoverySources = map[DiscoverySource]bool{
	DiscoverySourceFilesystem: true,
	DiscoverySourceRegistry:   true,
	DiscoverySourceURL:        true,
	DiscoverySourceGit:        true,
}

// AllDiscoverySources returns all valid discovery sources.
func AllDiscoverySources() []DiscoverySource {
	return []DiscoverySource{
		DiscoverySourceFilesystem,
		DiscoverySourceRegistry,
		DiscoverySourceURL,
		DiscoverySourceGit,
	}
}

// IsValidDiscoverySource checks if a discovery source is valid.
func IsValidDiscoverySource(source string) bool {
	return validDiscoverySources[DiscoverySource(source)]
}

// PluginStatus represents the current status of a plugin.
type PluginStatus string

const (
	// PluginStatusAvailable indicates plugin is available for installation.
	PluginStatusAvailable PluginStatus = "available"
	// PluginStatusInstalled indicates plugin is currently installed.
	PluginStatusInstalled PluginStatus = "installed"
	// PluginStatusActive indicates plugin is installed and running.
	PluginStatusActive PluginStatus = "active"
	// PluginStatusInactive indicates plugin is installed but not running.
	PluginStatusInactive PluginStatus = "inactive"
	// PluginStatusError indicates plugin is in error state.
	PluginStatusError PluginStatus = "error"
	// PluginStatusUpdating indicates plugin is currently being updated.
	PluginStatusUpdating PluginStatus = "updating"
)

// validPluginStatuses is a map for O(1) plugin status validation lookups.
var validPluginStatuses = map[PluginStatus]bool{
	PluginStatusAvailable: true,
	PluginStatusInstalled: true,
	PluginStatusActive:    true,
	PluginStatusInactive:  true,
	PluginStatusError:     true,
	PluginStatusUpdating:  true,
}

// AllPluginStatuses returns all valid plugin statuses.
func AllPluginStatuses() []PluginStatus {
	return []PluginStatus{
		PluginStatusAvailable,
		PluginStatusInstalled,
		PluginStatusActive,
		PluginStatusInactive,
		PluginStatusError,
		PluginStatusUpdating,
	}
}

// IsValidPluginStatus checks if a plugin status is valid.
func IsValidPluginStatus(status string) bool {
	return validPluginStatuses[PluginStatus(status)]
}

// SecurityLevel represents plugin security trust levels.
type SecurityLevel string

const (
	// SecurityLevelUntrusted indicates untrusted plugin requiring explicit approval.
	SecurityLevelUntrusted SecurityLevel = "untrusted"
	// SecurityLevelCommunity indicates community-verified plugin.
	SecurityLevelCommunity SecurityLevel = "community"
	// SecurityLevelVerified indicates officially verified plugin.
	SecurityLevelVerified SecurityLevel = "verified"
	// SecurityLevelOfficial indicates official PulumiCost plugin.
	SecurityLevelOfficial SecurityLevel = "official"
)

// validSecurityLevels is a map for O(1) security level validation lookups.
var validSecurityLevels = map[SecurityLevel]bool{
	SecurityLevelUntrusted: true,
	SecurityLevelCommunity: true,
	SecurityLevelVerified:  true,
	SecurityLevelOfficial:  true,
}

// AllSecurityLevels returns all valid security levels.
func AllSecurityLevels() []SecurityLevel {
	return []SecurityLevel{
		SecurityLevelUntrusted,
		SecurityLevelCommunity,
		SecurityLevelVerified,
		SecurityLevelOfficial,
	}
}

// IsValidSecurityLevel checks if a security level is valid.
func IsValidSecurityLevel(level string) bool {
	return validSecurityLevels[SecurityLevel(level)]
}

// InstallationMethod represents different plugin installation methods.
type InstallationMethod string

const (
	// InstallationMethodBinary indicates direct binary installation.
	InstallationMethodBinary InstallationMethod = "binary"
	// InstallationMethodContainer indicates container image deployment.
	InstallationMethodContainer InstallationMethod = "container"
	// InstallationMethodScript indicates script-based installation.
	InstallationMethodScript InstallationMethod = "script"
	// InstallationMethodPackage indicates system package manager installation.
	InstallationMethodPackage InstallationMethod = "package"
)

// validInstallationMethods is a map for O(1) installation method validation lookups.
var validInstallationMethods = map[InstallationMethod]bool{
	InstallationMethodBinary:    true,
	InstallationMethodContainer: true,
	InstallationMethodScript:    true,
	InstallationMethodPackage:   true,
}

// AllInstallationMethods returns all valid installation methods.
func AllInstallationMethods() []InstallationMethod {
	return []InstallationMethod{
		InstallationMethodBinary,
		InstallationMethodContainer,
		InstallationMethodScript,
		InstallationMethodPackage,
	}
}

// IsValidInstallationMethod checks if an installation method is valid.
func IsValidInstallationMethod(method string) bool {
	return validInstallationMethods[InstallationMethod(method)]
}

// PluginCapability represents plugin capabilities.
type PluginCapability string

const (
	// PluginCapabilityCostRetrieval indicates cost data retrieval capability.
	PluginCapabilityCostRetrieval PluginCapability = "cost_retrieval"
	// PluginCapabilityCostProjection indicates cost projection capability.
	PluginCapabilityCostProjection PluginCapability = "cost_projection"
	// PluginCapabilityPricingSpecs indicates pricing specification capability.
	PluginCapabilityPricingSpecs PluginCapability = "pricing_specs"
	// PluginCapabilityHistoricalData indicates historical data support.
	PluginCapabilityHistoricalData PluginCapability = "historical_data"
	// PluginCapabilityRealTimeData indicates real-time data support.
	PluginCapabilityRealTimeData PluginCapability = "real_time_data"
	// PluginCapabilityBatchProcessing indicates batch processing support.
	PluginCapabilityBatchProcessing PluginCapability = "batch_processing"
	// PluginCapabilityRateLimiting indicates rate limiting support.
	PluginCapabilityRateLimiting PluginCapability = "rate_limiting"
	// PluginCapabilityCaching indicates caching support.
	PluginCapabilityCaching PluginCapability = "caching"
	// PluginCapabilityEncryption indicates encryption support.
	PluginCapabilityEncryption PluginCapability = "encryption"
	// PluginCapabilityCompression indicates compression support.
	PluginCapabilityCompression PluginCapability = "compression"
	// PluginCapabilityFiltering indicates filtering support.
	PluginCapabilityFiltering PluginCapability = "filtering"
	// PluginCapabilityAggregation indicates aggregation support.
	PluginCapabilityAggregation PluginCapability = "aggregation"
	// PluginCapabilityMultiTenancy indicates multi-tenancy support.
	PluginCapabilityMultiTenancy PluginCapability = "multi_tenancy"
	// PluginCapabilityAuditLogging indicates audit logging support.
	PluginCapabilityAuditLogging PluginCapability = "audit_logging"
)

// validPluginCapabilities is a map for O(1) plugin capability validation lookups.
var validPluginCapabilities = map[PluginCapability]bool{
	PluginCapabilityCostRetrieval:    true,
	PluginCapabilityCostProjection:   true,
	PluginCapabilityPricingSpecs:     true,
	PluginCapabilityHistoricalData:   true,
	PluginCapabilityRealTimeData:     true,
	PluginCapabilityBatchProcessing:  true,
	PluginCapabilityRateLimiting:     true,
	PluginCapabilityCaching:          true,
	PluginCapabilityEncryption:       true,
	PluginCapabilityCompression:      true,
	PluginCapabilityFiltering:        true,
	PluginCapabilityAggregation:      true,
	PluginCapabilityMultiTenancy:     true,
	PluginCapabilityAuditLogging:     true,
}

// AllPluginCapabilities returns all valid plugin capabilities.
func AllPluginCapabilities() []PluginCapability {
	return []PluginCapability{
		PluginCapabilityCostRetrieval,
		PluginCapabilityCostProjection,
		PluginCapabilityPricingSpecs,
		PluginCapabilityHistoricalData,
		PluginCapabilityRealTimeData,
		PluginCapabilityBatchProcessing,
		PluginCapabilityRateLimiting,
		PluginCapabilityCaching,
		PluginCapabilityEncryption,
		PluginCapabilityCompression,
		PluginCapabilityFiltering,
		PluginCapabilityAggregation,
		PluginCapabilityMultiTenancy,
		PluginCapabilityAuditLogging,
	}
}

// IsValidPluginCapability checks if a plugin capability is valid.
func IsValidPluginCapability(capability string) bool {
	return validPluginCapabilities[PluginCapability(capability)]
}

// SystemPermission represents required system permissions.
type SystemPermission string

const (
	// SystemPermissionNetworkAccess indicates outbound network connection permission.
	SystemPermissionNetworkAccess SystemPermission = "network_access"
	// SystemPermissionFilesystemRead indicates filesystem read permission.
	SystemPermissionFilesystemRead SystemPermission = "filesystem_read"
	// SystemPermissionFilesystemWrite indicates filesystem write permission.
	SystemPermissionFilesystemWrite SystemPermission = "filesystem_write"
	// SystemPermissionEnvironmentRead indicates environment variable read permission.
	SystemPermissionEnvironmentRead SystemPermission = "environment_read"
	// SystemPermissionProcessSpawn indicates process spawn permission.
	SystemPermissionProcessSpawn SystemPermission = "process_spawn"
	// SystemPermissionSystemInfo indicates system information access permission.
	SystemPermissionSystemInfo SystemPermission = "system_info"
	// SystemPermissionTempFiles indicates temporary file creation permission.
	SystemPermissionTempFiles SystemPermission = "temp_files"
	// SystemPermissionConfigRead indicates configuration file read permission.
	SystemPermissionConfigRead SystemPermission = "config_read"
	// SystemPermissionMetricsCollect indicates metrics collection permission.
	SystemPermissionMetricsCollect SystemPermission = "metrics_collect"
)

// validSystemPermissions is a map for O(1) system permission validation lookups.
var validSystemPermissions = map[SystemPermission]bool{
	SystemPermissionNetworkAccess:   true,
	SystemPermissionFilesystemRead:  true,
	SystemPermissionFilesystemWrite: true,
	SystemPermissionEnvironmentRead: true,
	SystemPermissionProcessSpawn:    true,
	SystemPermissionSystemInfo:      true,
	SystemPermissionTempFiles:       true,
	SystemPermissionConfigRead:      true,
	SystemPermissionMetricsCollect:  true,
}

// AllSystemPermissions returns all valid system permissions.
func AllSystemPermissions() []SystemPermission {
	return []SystemPermission{
		SystemPermissionNetworkAccess,
		SystemPermissionFilesystemRead,
		SystemPermissionFilesystemWrite,
		SystemPermissionEnvironmentRead,
		SystemPermissionProcessSpawn,
		SystemPermissionSystemInfo,
		SystemPermissionTempFiles,
		SystemPermissionConfigRead,
		SystemPermissionMetricsCollect,
	}
}

// IsValidSystemPermission checks if a system permission is valid.
func IsValidSystemPermission(permission string) bool {
	return validSystemPermissions[SystemPermission(permission)]
}

// AuthMethod represents supported authentication methods.
type AuthMethod string

const (
	// AuthMethodNone indicates no authentication required.
	AuthMethodNone AuthMethod = "none"
	// AuthMethodAPIKey indicates API key authentication.
	AuthMethodAPIKey AuthMethod = "api_key"
	// AuthMethodJWT indicates JWT token authentication.
	AuthMethodJWT AuthMethod = "jwt"
	// AuthMethodOAuth2 indicates OAuth2 authentication.
	AuthMethodOAuth2 AuthMethod = "oauth2"
	// AuthMethodMTLS indicates mutual TLS authentication.
	AuthMethodMTLS AuthMethod = "mtls"
	// AuthMethodBasicAuth indicates basic HTTP authentication.
	AuthMethodBasicAuth AuthMethod = "basic_auth"
)

// validAuthMethods is a map for O(1) authentication method validation lookups.
var validAuthMethods = map[AuthMethod]bool{
	AuthMethodNone:      true,
	AuthMethodAPIKey:    true,
	AuthMethodJWT:       true,
	AuthMethodOAuth2:    true,
	AuthMethodMTLS:      true,
	AuthMethodBasicAuth: true,
}

// AllAuthMethods returns all valid authentication methods.
func AllAuthMethods() []AuthMethod {
	return []AuthMethod{
		AuthMethodNone,
		AuthMethodAPIKey,
		AuthMethodJWT,
		AuthMethodOAuth2,
		AuthMethodMTLS,
		AuthMethodBasicAuth,
	}
}

// IsValidAuthMethod checks if an authentication method is valid.
func IsValidAuthMethod(method string) bool {
	return validAuthMethods[AuthMethod(method)]
}

// ValidatePluginName validates a plugin name according to registry rules.
func ValidatePluginName(name string) error {
	if name == "" {
		return errors.New("plugin name cannot be empty")
	}
	if len(name) < MinPluginNameLength {
		return fmt.Errorf("plugin name must be at least %d characters long", MinPluginNameLength)
	}
	if len(name) > MaxPluginNameLength {
		return fmt.Errorf("plugin name must be no more than %d characters long", MaxPluginNameLength)
	}

	// Check pattern: ^[a-z0-9]([a-z0-9-]*[a-z0-9])?$
	if !isValidPluginNameChar(rune(name[0])) || name[0] == '-' {
		return errors.New("plugin name must start with a lowercase letter or digit")
	}
	if len(name) > 1 {
		if !isValidPluginNameChar(rune(name[len(name)-1])) || name[len(name)-1] == '-' {
			return errors.New("plugin name must end with a lowercase letter or digit")
		}
	}

	for i, r := range name {
		if !isValidPluginNameChar(r) && r != '-' {
			return fmt.Errorf("plugin name contains invalid character '%c' at position %d", r, i)
		}
	}

	return nil
}

// isValidPluginNameChar checks if a rune is valid for plugin names (a-z, 0-9).
func isValidPluginNameChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
}
