package registry

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// DefaultValidator provides default implementation for plugin validation.
type DefaultValidator struct{}

// NewValidator creates a new default validator.
func NewValidator() PluginValidator {
	return &DefaultValidator{}
}

// ValidateManifest validates a plugin manifest against the schema and business rules.
func (v *DefaultValidator) ValidateManifest(manifest *PluginManifest) (*ValidationResult, error) {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []ValidationError{},
		Warnings: []ValidationWarning{},
	}

	// Validate required fields
	if manifest.Name == "" {
		result.addError("MISSING_NAME", "name", "Plugin name is required")
	}
	if manifest.Version == "" {
		result.addError("MISSING_VERSION", "version", "Plugin version is required")
	}
	if manifest.APIVersion == "" {
		result.addError("MISSING_API_VERSION", "api_version", "API version is required")
	}
	if manifest.PluginType == "" {
		result.addError("MISSING_PLUGIN_TYPE", "plugin_type", "Plugin type is required")
	}
	if len(manifest.Capabilities) == 0 {
		result.addError("MISSING_CAPABILITIES", "capabilities", "At least one capability is required")
	}

	// Validate name format
	if manifest.Name != "" {
		if err := validatePluginName(manifest.Name); err != nil {
			result.addError("INVALID_NAME_FORMAT", "name", err.Error())
		}
	}

	// Validate version format (semantic versioning)
	if manifest.Version != "" {
		if err := validateSemVer(manifest.Version); err != nil {
			result.addError("INVALID_VERSION_FORMAT", "version", err.Error())
		}
	}

	// Validate API version format
	if manifest.APIVersion != "" {
		if err := validateAPIVersion(manifest.APIVersion); err != nil {
			result.addError("INVALID_API_VERSION", "api_version", err.Error())
		}
	}

	// Validate plugin type
	if manifest.PluginType != "" {
		if err := validatePluginType(manifest.PluginType); err != nil {
			result.addError("INVALID_PLUGIN_TYPE", "plugin_type", err.Error())
		}
	}

	// Validate capabilities
	for i, capability := range manifest.Capabilities {
		if err := validateCapability(capability); err != nil {
			result.addError("INVALID_CAPABILITY", fmt.Sprintf("capabilities[%d]", i), err.Error())
		}
	}

	// Validate supported providers
	for i, provider := range manifest.SupportedProviders {
		if err := validateProvider(provider); err != nil {
			result.addError("INVALID_PROVIDER", fmt.Sprintf("supported_providers[%d]", i), err.Error())
		}
	}

	// Validate requirements
	if err := v.validateRequirements(&manifest.Requirements, result); err != nil {
		return nil, err
	}

	// Validate authentication config
	if err := v.validateAuthenticationConfig(&manifest.Authentication, result); err != nil {
		return nil, err
	}

	// Validate installation info
	if err := v.validateInstallationInfo(&manifest.Installation, result); err != nil {
		return nil, err
	}

	// Validate contact info
	if err := v.validateContactInfo(&manifest.Contacts, result); err != nil {
		return nil, err
	}

	// Add warnings for best practices
	v.addBestPracticeWarnings(manifest, result)

	result.Valid = len(result.Errors) == 0
	return result, nil
}

// ValidateBinary validates a plugin binary for security and compatibility.
func (v *DefaultValidator) ValidateBinary(binaryData []byte, manifest *PluginManifest) (*ValidationResult, error) {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []ValidationError{},
		Warnings: []ValidationWarning{},
	}

	// Validate binary size
	if len(binaryData) == 0 {
		result.addError("EMPTY_BINARY", "binary", "Plugin binary is empty")
		result.Valid = false
		return result, nil
	}

	// Validate checksum if provided
	if manifest.Installation.Checksum != "" {
		expectedChecksum := strings.TrimPrefix(manifest.Installation.Checksum, "sha256:")
		actualChecksum := calculateSHA256(binaryData)
		
		if strings.ToLower(expectedChecksum) != strings.ToLower(actualChecksum) {
			result.addError("CHECKSUM_MISMATCH", "checksum", 
				fmt.Sprintf("Binary checksum mismatch. Expected: %s, Got: %s", expectedChecksum, actualChecksum))
		}
	}

	// Validate binary size matches manifest
	if manifest.Installation.SizeBytes > 0 {
		actualSize := int64(len(binaryData))
		if actualSize != manifest.Installation.SizeBytes {
			result.addError("SIZE_MISMATCH", "size_bytes",
				fmt.Sprintf("Binary size mismatch. Expected: %d bytes, Got: %d bytes", 
					manifest.Installation.SizeBytes, actualSize))
		}
	}

	// Basic executable validation
	if err := validateExecutable(binaryData); err != nil {
		result.addError("INVALID_EXECUTABLE", "binary", err.Error())
	}

	result.Valid = len(result.Errors) == 0
	return result, nil
}

// ValidateCompatibility checks if a plugin is compatible with the current system.
func (v *DefaultValidator) ValidateCompatibility(manifest *PluginManifest) (*ValidationResult, error) {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []ValidationError{},
		Warnings: []ValidationWarning{},
	}

	// For this reference implementation, we assume compatibility
	// In a real implementation, this would check:
	// - Operating system compatibility
	// - Architecture compatibility  
	// - API version compatibility
	// - Dependency availability

	if manifest.Requirements.MinAPIVersion != "" {
		// Check minimum API version compatibility
		// This would compare against the current system's API version
		result.addWarning("API_VERSION_CHECK", "min_api_version", 
			"API version compatibility check not implemented in reference implementation")
	}

	result.Valid = len(result.Errors) == 0
	return result, nil
}

// ScanSecurity performs security scanning on a plugin binary.
func (v *DefaultValidator) ScanSecurity(binaryData []byte) (*SecurityScanResults, error) {
	// This is a placeholder implementation for security scanning
	// In a real implementation, this would:
	// - Scan for known vulnerabilities
	// - Check for malicious patterns
	// - Validate code signatures
	// - Perform static analysis

	results := &SecurityScanResults{
		Status:          SecurityStatusSecure,
		Vulnerabilities: []SecurityVulnerability{},
		LastScan:        time.Now(),
		ScanVersion:     "1.0.0-reference",
	}

	// Basic security checks
	if len(binaryData) == 0 {
		results.Status = SecurityStatusVulnerable
		results.Vulnerabilities = append(results.Vulnerabilities, SecurityVulnerability{
			ID:          "SEC-001",
			Severity:    SeverityCritical,
			Summary:     "Empty binary",
			Description: "Plugin binary is empty which poses a security risk",
		})
	}

	return results, nil
}

// Helper functions

func (r *ValidationResult) addError(code, field, message string) {
	r.Errors = append(r.Errors, ValidationError{
		Code:     code,
		Field:    field,
		Message:  message,
		Severity: SeverityHigh,
	})
}

func (r *ValidationResult) addWarning(code, field, message string) {
	r.Warnings = append(r.Warnings, ValidationWarning{
		Code:    code,
		Field:   field,
		Message: message,
	})
}

func validatePluginName(name string) error {
	// Plugin names must be lowercase, alphanumeric with hyphens, 3-64 characters
	pattern := `^[a-z0-9][a-z0-9-]*[a-z0-9]$`
	matched, err := regexp.MatchString(pattern, name)
	if err != nil {
		return fmt.Errorf("error validating plugin name: %w", err)
	}
	if !matched {
		return fmt.Errorf("plugin name must be lowercase alphanumeric with hyphens, 3-64 characters")
	}
	if len(name) < 3 || len(name) > 64 {
		return fmt.Errorf("plugin name must be between 3 and 64 characters")
	}
	return nil
}

func validateSemVer(version string) error {
	// Basic semantic version validation
	pattern := `^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`
	matched, err := regexp.MatchString(pattern, version)
	if err != nil {
		return fmt.Errorf("error validating version: %w", err)
	}
	if !matched {
		return fmt.Errorf("version must follow semantic versioning 2.0.0 format")
	}
	return nil
}

func validateAPIVersion(apiVersion string) error {
	pattern := `^v[1-9]\d*$`
	matched, err := regexp.MatchString(pattern, apiVersion)
	if err != nil {
		return fmt.Errorf("error validating API version: %w", err)
	}
	if !matched {
		return fmt.Errorf("API version must be in format v1, v2, etc.")
	}
	return nil
}

func validatePluginType(pluginType PluginType) error {
	validTypes := []PluginType{
		PluginTypeCostSource,
		PluginTypeObservability,
		PluginTypeRegistry,
		PluginTypeAggregator,
	}
	
	for _, validType := range validTypes {
		if pluginType == validType {
			return nil
		}
	}
	
	return fmt.Errorf("invalid plugin type: %s", pluginType)
}

func validateCapability(capability Capability) error {
	validCapabilities := []Capability{
		CapabilityActualCost,
		CapabilityProjectedCost,
		CapabilityPricingSpec,
		CapabilityRealTimeCost,
		CapabilityHistoricalCost,
		CapabilityCostForecasting,
		CapabilityCostOptimization,
		CapabilityBudgetAlerts,
		CapabilityCustomMetrics,
		CapabilityMultiCloud,
		CapabilityKubernetesNative,
		CapabilityTaggingSupport,
		CapabilityDrillDown,
		CapabilityCostAllocation,
		CapabilityShowback,
		CapabilityChargeback,
	}
	
	for _, validCap := range validCapabilities {
		if capability == validCap {
			return nil
		}
	}
	
	return fmt.Errorf("invalid capability: %s", capability)
}

func validateProvider(provider string) error {
	validProviders := []string{"aws", "azure", "gcp", "kubernetes", "alibaba", "oracle", "ibm", "custom"}
	
	for _, validProvider := range validProviders {
		if provider == validProvider {
			return nil
		}
	}
	
	return fmt.Errorf("invalid provider: %s", provider)
}

func (v *DefaultValidator) validateRequirements(req *PluginRequirements, result *ValidationResult) error {
	// Validate API version requirements
	if req.MinAPIVersion != "" {
		if err := validateAPIVersion(req.MinAPIVersion); err != nil {
			result.addError("INVALID_MIN_API_VERSION", "requirements.min_api_version", err.Error())
		}
	}
	
	if req.MaxAPIVersion != "" {
		if err := validateAPIVersion(req.MaxAPIVersion); err != nil {
			result.addError("INVALID_MAX_API_VERSION", "requirements.max_api_version", err.Error())
		}
	}
	
	return nil
}

func (v *DefaultValidator) validateAuthenticationConfig(auth *AuthenticationConfig, result *ValidationResult) error {
	// Validate authentication methods
	validMethods := []AuthenticationMethod{
		AuthMethodAPIKey,
		AuthMethodOAuth2,
		AuthMethodServiceAccount,
		AuthMethodIAMRole,
		AuthMethodMutualTLS,
		AuthMethodBasicAuth,
		AuthMethodBearerToken,
	}
	
	for i, method := range auth.Methods {
		valid := false
		for _, validMethod := range validMethods {
			if method == validMethod {
				valid = true
				break
			}
		}
		if !valid {
			result.addError("INVALID_AUTH_METHOD", fmt.Sprintf("authentication.methods[%d]", i), 
				fmt.Sprintf("invalid authentication method: %s", method))
		}
	}
	
	return nil
}

func (v *DefaultValidator) validateInstallationInfo(install *InstallationInfo, result *ValidationResult) error {
	// Validate URLs if provided
	if install.BinaryURL != "" && !isValidURL(install.BinaryURL) {
		result.addError("INVALID_BINARY_URL", "installation.binary_url", "invalid binary URL format")
	}
	
	if install.InstallScript != "" && !isValidURL(install.InstallScript) {
		result.addError("INVALID_INSTALL_SCRIPT_URL", "installation.install_script", "invalid install script URL format")
	}
	
	// Validate checksum format
	if install.Checksum != "" {
		if !strings.HasPrefix(install.Checksum, "sha256:") && !isValidSHA256(install.Checksum) {
			result.addError("INVALID_CHECKSUM_FORMAT", "installation.checksum", 
				"checksum must be SHA-256 in hex format or prefixed with 'sha256:'")
		}
	}
	
	return nil
}

func (v *DefaultValidator) validateContactInfo(contacts *ContactInfo, result *ValidationResult) error {
	// Validate maintainer emails
	for i, maintainer := range contacts.Maintainers {
		if maintainer.Email != "" && !isValidEmail(maintainer.Email) {
			result.addError("INVALID_MAINTAINER_EMAIL", fmt.Sprintf("contacts.maintainers[%d].email", i),
				"invalid email format")
		}
		if maintainer.URL != "" && !isValidURL(maintainer.URL) {
			result.addError("INVALID_MAINTAINER_URL", fmt.Sprintf("contacts.maintainers[%d].url", i),
				"invalid URL format")
		}
	}
	
	// Validate support and documentation URLs
	if contacts.SupportURL != "" && !isValidURL(contacts.SupportURL) {
		result.addError("INVALID_SUPPORT_URL", "contacts.support_url", "invalid support URL format")
	}
	
	if contacts.DocumentationURL != "" && !isValidURL(contacts.DocumentationURL) {
		result.addError("INVALID_DOCUMENTATION_URL", "contacts.documentation_url", "invalid documentation URL format")
	}
	
	return nil
}

func (v *DefaultValidator) addBestPracticeWarnings(manifest *PluginManifest, result *ValidationResult) {
	// Check if description is provided
	if manifest.Description == "" {
		result.addWarning("MISSING_DESCRIPTION", "description", "Plugin description is recommended for better discoverability")
	}
	
	// Check if display name is provided
	if manifest.DisplayName == "" {
		result.addWarning("MISSING_DISPLAY_NAME", "display_name", "Display name is recommended for better user experience")
	}
	
	// Check if tags are provided
	if len(manifest.Metadata.Tags) == 0 {
		result.addWarning("MISSING_TAGS", "metadata.tags", "Tags are recommended for better categorization")
	}
	
	// Check if license is provided
	if manifest.Metadata.License == "" {
		result.addWarning("MISSING_LICENSE", "metadata.license", "License information is recommended")
	}
}

func calculateSHA256(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func validateExecutable(data []byte) error {
	// Basic validation - check if it looks like an executable
	if len(data) < 4 {
		return fmt.Errorf("binary too small to be a valid executable")
	}
	
	// Check for common executable magic bytes
	// ELF: 7f 45 4c 46
	// PE: 4d 5a (MZ)
	// Mach-O: fe ed fa ce, fe ed fa cf, ce fa ed fe, cf fa ed fe
	
	if len(data) >= 4 {
		// ELF check
		if data[0] == 0x7f && data[1] == 'E' && data[2] == 'L' && data[3] == 'F' {
			return nil
		}
		// PE check
		if data[0] == 'M' && data[1] == 'Z' {
			return nil
		}
		// Basic Mach-O check
		if (data[0] == 0xfe && data[1] == 0xed && data[2] == 0xfa) ||
		   (data[0] == 0xce && data[1] == 0xfa && data[2] == 0xed) ||
		   (data[0] == 0xcf && data[1] == 0xfa && data[2] == 0xed) {
			return nil
		}
	}
	
	return fmt.Errorf("binary does not appear to be a valid executable")
}

func isValidURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func isValidSHA256(checksum string) bool {
	if len(checksum) != 64 {
		return false
	}
	pattern := `^[a-fA-F0-9]{64}$`
	matched, _ := regexp.MatchString(pattern, checksum)
	return matched
}