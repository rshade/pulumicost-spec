package pricing

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// ErrorCategory defines the category of plugin errors.
type ErrorCategory string

const (
	// TransientError indicates a temporary failure that may succeed on retry.
	TransientError ErrorCategory = "transient"
	// PermanentError indicates a failure that will not succeed on retry without changes.
	PermanentError ErrorCategory = "permanent"
	// ConfigurationError indicates an error due to invalid configuration or setup.
	ConfigurationError ErrorCategory = "configuration"
)

// ErrorCode defines standard error codes for plugin operations.
type ErrorCode string

const (
	// ErrorCodeNetworkTimeout indicates a network timeout occurred.
	ErrorCodeNetworkTimeout ErrorCode = "NETWORK_TIMEOUT"
	// ErrorCodeServiceUnavailable indicates the service is temporarily unavailable.
	ErrorCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	// ErrorCodeRateLimited indicates the request was rate limited.
	ErrorCodeRateLimited ErrorCode = "RATE_LIMITED"
	// ErrorCodeTemporaryFailure indicates a temporary failure occurred.
	ErrorCodeTemporaryFailure ErrorCode = "TEMPORARY_FAILURE"
	// ErrorCodeCircuitOpen indicates the circuit breaker is open.
	ErrorCodeCircuitOpen ErrorCode = "CIRCUIT_OPEN"

	// ErrorCodeInvalidResource indicates the resource specification is invalid.
	ErrorCodeInvalidResource ErrorCode = "INVALID_RESOURCE"
	// ErrorCodeResourceNotFound indicates the requested resource was not found.
	ErrorCodeResourceNotFound ErrorCode = "RESOURCE_NOT_FOUND"
	// ErrorCodeInvalidTimeRange indicates the time range is invalid.
	ErrorCodeInvalidTimeRange ErrorCode = "INVALID_TIME_RANGE"
	// ErrorCodeUnsupportedRegion indicates the region is not supported.
	ErrorCodeUnsupportedRegion ErrorCode = "UNSUPPORTED_REGION"
	// ErrorCodePermissionDenied indicates access is denied.
	ErrorCodePermissionDenied ErrorCode = "PERMISSION_DENIED"
	// ErrorCodeDataCorruption indicates data corruption was detected.
	ErrorCodeDataCorruption ErrorCode = "DATA_CORRUPTION"

	// ErrorCodeInvalidCredentials indicates authentication credentials are invalid.
	//nolint:gosec // This is an error code constant, not actual credentials
	ErrorCodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	// ErrorCodeMissingAPIKey indicates the API key is missing.
	ErrorCodeMissingAPIKey ErrorCode = "MISSING_API_KEY"
	// ErrorCodeInvalidEndpoint indicates the endpoint configuration is invalid.
	ErrorCodeInvalidEndpoint ErrorCode = "INVALID_ENDPOINT"
	// ErrorCodeInvalidProvider indicates the provider is invalid.
	ErrorCodeInvalidProvider ErrorCode = "INVALID_PROVIDER"
	// ErrorCodePluginNotConfigured indicates the plugin is not properly configured.
	ErrorCodePluginNotConfigured ErrorCode = "PLUGIN_NOT_CONFIGURED"
)

// Error message validation constants.
const (
	minErrorMessageLength = 10  // Minimum allowed error message length
	maxErrorMessageLength = 500 // Maximum allowed error message length
)

// Retry policy constants.
const (
	defaultMaxRetries   = 3                      // Default maximum number of retry attempts
	defaultBaseDelay    = 100 * time.Millisecond // Default base delay for exponential backoff
	defaultMaxDelay     = 30 * time.Second       // Default maximum delay between retries
	defaultMultiplier   = 2.0                    // Default exponential backoff multiplier
	defaultJitterFactor = 0.1                    // Default jitter factor (10% randomness)
	maxJitterFactor     = 0.5                    // Maximum allowed jitter factor (50%)

	// Conservative retry policy constants.
	conservativeMaxRetries   = 2                      // Conservative maximum retry attempts
	conservativeBaseDelay    = 200 * time.Millisecond // Conservative base delay
	conservativeMaxDelay     = 10 * time.Second       // Conservative maximum delay
	conservativeMultiplier   = 1.5                    // Conservative multiplier
	conservativeJitterFactor = 0.05                   // Conservative jitter factor

	// Aggressive retry policy constants.
	aggressiveMaxRetries   = 5                     // Aggressive maximum retry attempts
	aggressiveBaseDelay    = 50 * time.Millisecond // Aggressive base delay
	aggressiveMaxDelay     = 60 * time.Second      // Aggressive maximum delay
	aggressiveMultiplier   = 2.5                   // Aggressive multiplier
	aggressiveJitterFactor = 0.2                   // Aggressive jitter factor

	// Jitter calculation constants.
	jitterRangeMultiplier     = 2   // Multiplier for jitter range calculation
	secureRandomFallbackValue = 0.5 // Fallback value when secure random generation fails

	// Cryptographic random number generation constants.
	float64PrecisionBits = 53 // Number of bits for full float64 precision (2^53)

	// Circuit breaker constants.
	defaultFailureThreshold       = 5                // Default number of failures before opening circuit
	defaultRecoveryTimeout        = 60 * time.Second // Default timeout before attempting recovery
	defaultSuccessThreshold       = 3                // Default number of successes needed to close circuit
	defaultRequestVolumeThreshold = 10               // Default minimum requests before evaluating circuit state
	consecutiveFailureMultiplier  = 2                // Multiplier for consecutive failure limit calculation
	defaultFailureRateThreshold   = 0.5              // Default failure rate threshold (50%)

	// RPC method timeout constants.
	defaultNameTimeout             = 5 * time.Second  // Default timeout for Name RPC
	defaultSupportsTimeout         = 10 * time.Second // Default timeout for Supports RPC
	defaultGetActualCostTimeout    = 30 * time.Second // Default timeout for GetActualCost RPC
	defaultGetProjectedCostTimeout = 15 * time.Second // Default timeout for GetProjectedCost RPC
	defaultGetPricingSpecTimeout   = 20 * time.Second // Default timeout for GetPricingSpec RPC

	// Global timeout constants.
	defaultGlobalTimeout = 60 * time.Second  // Default global timeout for all operations
	minimumTimeout       = 1 * time.Second   // Minimum allowed timeout value
	maximumTimeout       = 300 * time.Second // Maximum allowed timeout value (5 minutes)

	// Fast timeout configuration constants.
	fastNameTimeout             = 2 * time.Second  // Fast timeout for Name RPC
	fastSupportsTimeout         = 5 * time.Second  // Fast timeout for Supports RPC
	fastGetActualCostTimeout    = 15 * time.Second // Fast timeout for GetActualCost RPC
	fastGetProjectedCostTimeout = 8 * time.Second  // Fast timeout for GetProjectedCost RPC
	fastGetPricingSpecTimeout   = 10 * time.Second // Fast timeout for GetPricingSpec RPC
	fastGlobalTimeout           = 30 * time.Second // Fast global timeout

	// Slow timeout configuration constants.
	slowNameTimeout             = 10 * time.Second  // Slow timeout for Name RPC
	slowSupportsTimeout         = 20 * time.Second  // Slow timeout for Supports RPC
	slowGetActualCostTimeout    = 120 * time.Second // Slow timeout for GetActualCost RPC
	slowGetProjectedCostTimeout = 30 * time.Second  // Slow timeout for GetProjectedCost RPC
	slowGetPricingSpecTimeout   = 45 * time.Second  // Slow timeout for GetPricingSpec RPC
	slowGlobalTimeout           = 180 * time.Second // Slow global timeout
)

// PluginError represents a standardized plugin error with category and retry information.
type PluginError struct {
	Code       ErrorCode
	Category   ErrorCategory
	Message    string
	Details    map[string]interface{}
	Timestamp  time.Time
	Retryable  bool
	RetryAfter *time.Duration
}

// NewTransientError creates a new transient error.
func NewTransientError(code ErrorCode, message string, retryAfter *time.Duration) *PluginError {
	return &PluginError{
		Code:       code,
		Category:   TransientError,
		Message:    message,
		Details:    make(map[string]interface{}),
		Timestamp:  time.Now(),
		Retryable:  true,
		RetryAfter: retryAfter,
	}
}

// NewPermanentError creates a new permanent error.
func NewPermanentError(code ErrorCode, message string) *PluginError {
	return &PluginError{
		Code:      code,
		Category:  PermanentError,
		Message:   message,
		Details:   make(map[string]interface{}),
		Timestamp: time.Now(),
		Retryable: false,
	}
}

// NewConfigurationError creates a new configuration error.
func NewConfigurationError(code ErrorCode, message string) *PluginError {
	return &PluginError{
		Code:      code,
		Category:  ConfigurationError,
		Message:   message,
		Details:   make(map[string]interface{}),
		Timestamp: time.Now(),
		Retryable: false,
	}
}

// NewFormattedTransientError creates a transient error with formatted message.
func NewFormattedTransientError(code ErrorCode, params map[string]string, retryAfter *time.Duration) *PluginError {
	message := FormatErrorMessage(code, params)
	return NewTransientError(code, message, retryAfter).WithDetails(convertToInterface(params))
}

// NewFormattedPermanentError creates a permanent error with formatted message.
func NewFormattedPermanentError(code ErrorCode, params map[string]string) *PluginError {
	message := FormatErrorMessage(code, params)
	return NewPermanentError(code, message).WithDetails(convertToInterface(params))
}

// NewFormattedConfigurationError creates a configuration error with formatted message.
func NewFormattedConfigurationError(code ErrorCode, params map[string]string) *PluginError {
	message := FormatErrorMessage(code, params)
	return NewConfigurationError(code, message).WithDetails(convertToInterface(params))
}

// Error implements the error interface.
func (e *PluginError) Error() string {
	if len(e.Details) > 0 {
		return fmt.Sprintf("[%s:%s] %s (details: %v)", e.Category, e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s:%s] %s", e.Category, e.Code, e.Message)
}

// IsRetryable returns whether the error should be retried.
func (e *PluginError) IsRetryable() bool {
	return e.Retryable
}

// GetRetryAfter returns the suggested retry delay.
func (e *PluginError) GetRetryAfter() *time.Duration {
	return e.RetryAfter
}

// GetGRPCStatus converts the plugin error to a gRPC status.
func (e *PluginError) GetGRPCStatus() *status.Status {
	var code codes.Code

	switch e.Code {
	case ErrorCodeNetworkTimeout, ErrorCodeServiceUnavailable:
		code = codes.Unavailable
	case ErrorCodeRateLimited:
		code = codes.ResourceExhausted
	case ErrorCodeInvalidResource, ErrorCodeInvalidTimeRange:
		code = codes.InvalidArgument
	case ErrorCodeResourceNotFound:
		code = codes.NotFound
	case ErrorCodePermissionDenied, ErrorCodeInvalidCredentials:
		code = codes.PermissionDenied
	case ErrorCodeUnsupportedRegion, ErrorCodeInvalidProvider:
		code = codes.Unimplemented
	case ErrorCodeTemporaryFailure, ErrorCodeCircuitOpen:
		code = codes.Unavailable
	case ErrorCodeDataCorruption:
		code = codes.DataLoss
	case ErrorCodeMissingAPIKey, ErrorCodeInvalidEndpoint, ErrorCodePluginNotConfigured:
		code = codes.FailedPrecondition
	default:
		code = codes.Internal
	}

	return status.New(code, e.Error())
}

// WithDetails adds details to the error.
func (e *PluginError) WithDetails(details map[string]interface{}) *PluginError {
	for k, v := range details {
		e.Details[k] = v
	}
	return e
}

// GetErrorMapping returns a map of error codes to their categories.
func GetErrorMapping() map[ErrorCode]ErrorCategory {
	return map[ErrorCode]ErrorCategory{
		// Transient errors
		ErrorCodeNetworkTimeout:     TransientError,
		ErrorCodeServiceUnavailable: TransientError,
		ErrorCodeRateLimited:        TransientError,
		ErrorCodeTemporaryFailure:   TransientError,
		ErrorCodeCircuitOpen:        TransientError,

		// Permanent errors
		ErrorCodeInvalidResource:   PermanentError,
		ErrorCodeResourceNotFound:  PermanentError,
		ErrorCodeInvalidTimeRange:  PermanentError,
		ErrorCodeUnsupportedRegion: PermanentError,
		ErrorCodePermissionDenied:  PermanentError,
		ErrorCodeDataCorruption:    PermanentError,

		// Configuration errors
		ErrorCodeInvalidCredentials:  ConfigurationError,
		ErrorCodeMissingAPIKey:       ConfigurationError,
		ErrorCodeInvalidEndpoint:     ConfigurationError,
		ErrorCodeInvalidProvider:     ConfigurationError,
		ErrorCodePluginNotConfigured: ConfigurationError,
	}
}

// IsTransientError checks if an error is transient.
func IsTransientError(err error) bool {
	var pluginErr *PluginError
	if errors.As(err, &pluginErr) {
		return pluginErr.Category == TransientError
	}
	return false
}

// IsPermanentError checks if an error is permanent.
func IsPermanentError(err error) bool {
	var pluginErr *PluginError
	if errors.As(err, &pluginErr) {
		return pluginErr.Category == PermanentError
	}
	return false
}

// IsConfigurationError checks if an error is configuration-related.
func IsConfigurationError(err error) bool {
	var pluginErr *PluginError
	if errors.As(err, &pluginErr) {
		return pluginErr.Category == ConfigurationError
	}
	return false
}

// ErrorMessageTemplate defines a template for error messages.
type ErrorMessageTemplate struct {
	Format      string
	Description string
	Examples    []string
}

// getTransientErrorTemplates returns templates for transient errors.
func getTransientErrorTemplates() map[ErrorCode]ErrorMessageTemplate {
	//nolint:exhaustive // Only handles transient error codes
	return map[ErrorCode]ErrorMessageTemplate{
		ErrorCodeNetworkTimeout: {
			Format:      "Network timeout occurred while {operation} for {resource_type} {resource_id}: {details}",
			Description: "Use when network operations exceed timeout limits",
			Examples: []string{
				"Network timeout occurred while retrieving costs for ec2 i-123456: connection timeout after 30s",
				"Network timeout occurred while querying pricing for s3 bucket-name: API timeout",
			},
		},
		ErrorCodeServiceUnavailable: {
			Format:      "Service temporarily unavailable for {operation}: {details}",
			Description: "Use when external services are temporarily down",
			Examples: []string{
				"Service temporarily unavailable for cost retrieval: AWS API returning 503",
				"Service temporarily unavailable for pricing query: maintenance window in progress",
			},
		},
		ErrorCodeRateLimited: {
			Format:      "Rate limit exceeded for {operation}: {details}. Retry after {retry_after}",
			Description: "Use when API rate limits are exceeded",
			Examples: []string{
				"Rate limit exceeded for cost retrieval: 1000 requests/minute limit reached. Retry after 60s",
				"Rate limit exceeded for pricing query: quota exceeded. Retry after 300s",
			},
		},
		ErrorCodeTemporaryFailure: {
			Format:      "Temporary failure in {operation} for {resource_type}: {details}",
			Description: "Use for transient failures without specific category",
			Examples: []string{
				"Temporary failure in cost calculation for ec2: backend service overloaded",
				"Temporary failure in pricing lookup for storage: database connection lost",
			},
		},
		ErrorCodeCircuitOpen: {
			Format:      "Circuit breaker open for {service}: {details}. Retry after {retry_after}",
			Description: "Use when circuit breaker prevents requests",
			Examples: []string{
				"Circuit breaker open for AWS Cost Explorer: 5 consecutive failures. Retry after 120s",
				"Circuit breaker open for pricing service: error rate threshold exceeded. Retry after 300s",
			},
		},
	}
}

// getPermanentErrorTemplates returns templates for permanent errors.
func getPermanentErrorTemplates() map[ErrorCode]ErrorMessageTemplate {
	//nolint:exhaustive // Only handles permanent error codes
	return map[ErrorCode]ErrorMessageTemplate{
		ErrorCodeInvalidResource: {
			Format:      "Invalid resource specification: {details}",
			Description: "Use when resource parameters are malformed or invalid",
			Examples: []string{
				"Invalid resource specification: missing required field 'provider'",
				"Invalid resource specification: unsupported resource_type 'invalid-type' for provider 'aws'",
			},
		},
		ErrorCodeResourceNotFound: {
			Format:      "Resource not found: {resource_type} {resource_id} in {region}",
			Description: "Use when specified resource cannot be located",
			Examples: []string{
				"Resource not found: ec2 i-nonexistent in us-east-1",
				"Resource not found: s3 missing-bucket in global",
			},
		},
		ErrorCodeInvalidTimeRange: {
			Format:      "Invalid time range: {details}",
			Description: "Use when time range parameters are invalid",
			Examples: []string{
				"Invalid time range: start time cannot be after end time",
				"Invalid time range: time range exceeds maximum allowed duration of 90 days",
			},
		},
		ErrorCodeUnsupportedRegion: {
			Format:      "Unsupported region: {region} for provider {provider}",
			Description: "Use when region is not supported by the provider",
			Examples: []string{
				"Unsupported region: mars-1 for provider aws",
				"Unsupported region: unknown-region for provider azure",
			},
		},
		ErrorCodePermissionDenied: {
			Format:      "Permission denied for {operation}: {details}",
			Description: "Use when access is denied due to insufficient permissions",
			Examples: []string{
				"Permission denied for cost retrieval: insufficient IAM permissions",
				"Permission denied for resource access: API key lacks required scope",
			},
		},
		ErrorCodeDataCorruption: {
			Format:      "Data corruption detected in {resource_type} {resource_id}: {details}",
			Description: "Use when data integrity issues are detected",
			Examples: []string{
				"Data corruption detected in cost data for i-123456: checksum mismatch",
				"Data corruption detected in pricing spec for compute: invalid rate values",
			},
		},
	}
}

// getConfigurationErrorTemplates returns templates for configuration errors.
func getConfigurationErrorTemplates() map[ErrorCode]ErrorMessageTemplate {
	//nolint:exhaustive // Only handles configuration error codes
	return map[ErrorCode]ErrorMessageTemplate{
		ErrorCodeInvalidCredentials: {
			Format:      "Invalid credentials for {provider}: {details}",
			Description: "Use when authentication credentials are invalid",
			Examples: []string{
				"Invalid credentials for aws: access key expired",
				"Invalid credentials for azure: service principal authentication failed",
			},
		},
		ErrorCodeMissingAPIKey: {
			Format:      "Missing API key for {service}: {details}",
			Description: "Use when required API key is not provided",
			Examples: []string{
				"Missing API key for pricing service: API_KEY environment variable not set",
				"Missing API key for cost retrieval: authentication token required",
			},
		},
		ErrorCodeInvalidEndpoint: {
			Format:      "Invalid endpoint configuration for {service}: {details}",
			Description: "Use when endpoint URLs or configurations are invalid",
			Examples: []string{
				"Invalid endpoint configuration for cost API: malformed URL 'invalid-url'",
				"Invalid endpoint configuration for pricing service: unreachable host",
			},
		},
		ErrorCodeInvalidProvider: {
			Format:      "Invalid provider '{provider}': {details}",
			Description: "Use when provider specification is invalid",
			Examples: []string{
				"Invalid provider 'unknown': supported providers are aws, azure, gcp, kubernetes",
				"Invalid provider '': provider cannot be empty",
			},
		},
		ErrorCodePluginNotConfigured: {
			Format:      "Plugin not configured: {details}",
			Description: "Use when plugin lacks required configuration",
			Examples: []string{
				"Plugin not configured: missing required configuration 'region'",
				"Plugin not configured: invalid configuration format in config.yaml",
			},
		},
	}
}

// GetErrorMessageTemplates returns standardized error message templates.
func GetErrorMessageTemplates() map[ErrorCode]ErrorMessageTemplate {
	templates := make(map[ErrorCode]ErrorMessageTemplate)

	// Merge all template maps
	for k, v := range getTransientErrorTemplates() {
		templates[k] = v
	}
	for k, v := range getPermanentErrorTemplates() {
		templates[k] = v
	}
	for k, v := range getConfigurationErrorTemplates() {
		templates[k] = v
	}

	return templates
}

// FormatErrorMessage formats an error message using the template for the given error code.
func FormatErrorMessage(code ErrorCode, params map[string]string) string {
	templates := GetErrorMessageTemplates()
	template, exists := templates[code]
	if !exists {
		return fmt.Sprintf("Unknown error code %s", code)
	}

	message := template.Format
	for key, value := range params {
		placeholder := "{" + key + "}"
		message = strings.ReplaceAll(message, placeholder, value)
	}

	return message
}

// StandardErrorDetails creates a standard details map for common error parameters.
func StandardErrorDetails(operation, resourceType, resourceID, region string) map[string]string {
	details := make(map[string]string)
	if operation != "" {
		details["operation"] = operation
	}
	if resourceType != "" {
		details["resource_type"] = resourceType
	}
	if resourceID != "" {
		details["resource_id"] = resourceID
	}
	if region != "" {
		details["region"] = region
	}
	return details
}

// convertToInterface converts string map to interface map.
func convertToInterface(params map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range params {
		result[k] = v
	}
	return result
}

// ValidateErrorMessage checks if an error message follows the standard format.
func ValidateErrorMessage(code ErrorCode, message string) error {
	templates := GetErrorMessageTemplates()
	_, exists := templates[code]
	if !exists {
		return fmt.Errorf("unknown error code: %s", code)
	}

	// Check if message follows the general structure
	if message == "" {
		return errors.New("error message cannot be empty")
	}

	// Check for minimum length
	if len(message) < minErrorMessageLength {
		return errors.New("error message too short: minimum 10 characters required")
	}

	// Check for maximum length
	if len(message) > maxErrorMessageLength {
		return errors.New("error message too long: maximum 500 characters allowed")
	}

	// Check that error code is referenced in the message
	if !strings.Contains(message, string(code)) {
		return fmt.Errorf("error message should contain error code %s", code)
	}

	return nil
}

// RetryPolicy defines the policy for retrying failed operations.
type RetryPolicy struct {
	MaxRetries      int           // Maximum number of retry attempts
	BaseDelay       time.Duration // Base delay for exponential backoff
	MaxDelay        time.Duration // Maximum delay between retries
	Multiplier      float64       // Exponential backoff multiplier
	JitterFactor    float64       // Jitter factor for randomizing delays (0.0-0.5)
	RetryableErrors []ErrorCode   // Specific error codes that should be retried
}

// NewDefaultRetryPolicy creates a retry policy with sensible defaults.
func NewDefaultRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxRetries:   defaultMaxRetries,
		BaseDelay:    defaultBaseDelay,
		MaxDelay:     defaultMaxDelay,
		Multiplier:   defaultMultiplier,
		JitterFactor: defaultJitterFactor,
		RetryableErrors: []ErrorCode{
			ErrorCodeNetworkTimeout,
			ErrorCodeServiceUnavailable,
			ErrorCodeRateLimited,
			ErrorCodeTemporaryFailure,
		},
	}
}

// NewConservativeRetryPolicy creates a retry policy with conservative settings.
func NewConservativeRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxRetries:   conservativeMaxRetries,
		BaseDelay:    conservativeBaseDelay,
		MaxDelay:     conservativeMaxDelay,
		Multiplier:   conservativeMultiplier,
		JitterFactor: conservativeJitterFactor,
		RetryableErrors: []ErrorCode{
			ErrorCodeNetworkTimeout,
			ErrorCodeServiceUnavailable,
		},
	}
}

// NewAggressiveRetryPolicy creates a retry policy with aggressive settings.
func NewAggressiveRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxRetries:   aggressiveMaxRetries,
		BaseDelay:    aggressiveBaseDelay,
		MaxDelay:     aggressiveMaxDelay,
		Multiplier:   aggressiveMultiplier,
		JitterFactor: aggressiveJitterFactor,
		RetryableErrors: []ErrorCode{
			ErrorCodeNetworkTimeout,
			ErrorCodeServiceUnavailable,
			ErrorCodeRateLimited,
			ErrorCodeTemporaryFailure,
			ErrorCodeCircuitOpen,
		},
	}
}

// Validate checks if the retry policy has valid parameters.
func (rp *RetryPolicy) Validate() error {
	if rp.MaxRetries < 0 {
		return errors.New("max retries cannot be negative")
	}
	if rp.BaseDelay <= 0 {
		return errors.New("base delay must be positive")
	}
	if rp.MaxDelay <= 0 {
		return errors.New("max delay must be positive")
	}
	if rp.BaseDelay > rp.MaxDelay {
		return errors.New("base delay cannot be greater than max delay")
	}
	if rp.Multiplier <= 1.0 {
		return errors.New("multiplier must be greater than 1.0")
	}
	if rp.JitterFactor < 0 || rp.JitterFactor > maxJitterFactor {
		return fmt.Errorf("jitter factor must be between 0.0 and %f", maxJitterFactor)
	}
	return nil
}

// ShouldRetry determines if an error should be retried based on the policy.
func (rp *RetryPolicy) ShouldRetry(err error, attempt int) bool {
	if attempt >= rp.MaxRetries {
		return false
	}

	var pluginErr *PluginError
	if !errors.As(err, &pluginErr) {
		return false // Only retry PluginError instances
	}

	// Check if the error category is retryable
	if pluginErr.Category != TransientError {
		return false
	}

	// Check if the specific error code is in the retryable list
	for _, retryableCode := range rp.RetryableErrors {
		if pluginErr.Code == retryableCode {
			return true
		}
	}

	return false
}

// CalculateDelay calculates the delay for the given retry attempt.
func (rp *RetryPolicy) CalculateDelay(attempt int) time.Duration {
	if attempt < 0 {
		return rp.BaseDelay
	}

	// Calculate exponential backoff delay
	delay := float64(rp.BaseDelay) * math.Pow(rp.Multiplier, float64(attempt))

	// Apply maximum delay limit
	if delay > float64(rp.MaxDelay) {
		delay = float64(rp.MaxDelay)
	}

	// Add jitter to prevent thundering herd problem
	if rp.JitterFactor > 0 {
		jitter := delay * rp.JitterFactor * (secureRandFloat()*jitterRangeMultiplier - 1) // Random value between -jitterFactor and +jitterFactor
		delay += jitter
	}

	// Ensure delay is not negative
	if delay < 0 {
		delay = float64(rp.BaseDelay)
	}

	return time.Duration(delay)
}

// RetryFunc represents a function that can be retried.
type RetryFunc func() error

// RetryWithPolicy executes a function with retry logic based on the provided policy.
func RetryWithPolicy(ctx context.Context, policy *RetryPolicy, fn RetryFunc) error {
	if policy == nil {
		policy = NewDefaultRetryPolicy()
	}

	if err := policy.Validate(); err != nil {
		return fmt.Errorf("invalid retry policy: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= policy.MaxRetries; attempt++ {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Execute the function
		err := fn()
		if err == nil {
			return nil // Success
		}

		lastErr = err

		// Check if we should retry
		if !policy.ShouldRetry(err, attempt) {
			break
		}

		// If this is the last allowed attempt, don't wait
		if attempt >= policy.MaxRetries {
			break
		}

		// Calculate and wait for the delay
		delay := policy.CalculateDelay(attempt)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			// Continue to next attempt
		}
	}

	return lastErr
}

// RetryWithDefaultPolicy executes a function with the default retry policy.
func RetryWithDefaultPolicy(ctx context.Context, fn RetryFunc) error {
	return RetryWithPolicy(ctx, NewDefaultRetryPolicy(), fn)
}

// CircuitBreakerState represents the state of a circuit breaker.
type CircuitBreakerState int

const (
	// CircuitClosed indicates the circuit is closed and requests are allowed.
	CircuitClosed CircuitBreakerState = iota
	// CircuitOpen indicates the circuit is open and requests are blocked.
	CircuitOpen
	// CircuitHalfOpen indicates the circuit is testing if the service has recovered.
	CircuitHalfOpen
)

// secureRandFloat generates a cryptographically secure random float64 between 0.0 and 1.0.
func secureRandFloat() float64 {
	// Generate a random 64-bit integer
	maxValue := big.NewInt(1 << float64PrecisionBits) // Use 2^53 for full float64 precision
	n, err := rand.Int(rand.Reader, maxValue)
	if err != nil {
		// Fallback to a fixed value if random generation fails
		return secureRandomFallbackValue
	}

	// Convert to float64 between 0.0 and 1.0
	return float64(n.Int64()) / float64(maxValue.Int64())
}

// CircuitBreakerConfig defines the configuration for a circuit breaker.
type CircuitBreakerConfig struct {
	FailureThreshold        int           // Number of failures before opening circuit
	RecoveryTimeout         time.Duration // Timeout before attempting recovery
	SuccessThreshold        int           // Number of successes needed to close circuit
	RequestVolumeThreshold  int           // Minimum requests before evaluating circuit state
	FailureRateThreshold    float64       // Failure rate threshold (0.0-1.0) for opening circuit
	ConsecutiveFailureLimit int           // Maximum consecutive failures before forcing open
}

// NewDefaultCircuitBreakerConfig creates a circuit breaker config with sensible defaults.
func NewDefaultCircuitBreakerConfig() *CircuitBreakerConfig {
	return &CircuitBreakerConfig{
		FailureThreshold:        defaultFailureThreshold,
		RecoveryTimeout:         defaultRecoveryTimeout,
		SuccessThreshold:        defaultSuccessThreshold,
		RequestVolumeThreshold:  defaultRequestVolumeThreshold,
		FailureRateThreshold:    defaultFailureRateThreshold,                            // 50% failure rate
		ConsecutiveFailureLimit: defaultFailureThreshold * consecutiveFailureMultiplier, // Double the failure threshold
	}
}

// Validate checks if the circuit breaker config has valid parameters.
func (cbc *CircuitBreakerConfig) Validate() error {
	if cbc.FailureThreshold <= 0 {
		return errors.New("failure threshold must be positive")
	}
	if cbc.RecoveryTimeout <= 0 {
		return errors.New("recovery timeout must be positive")
	}
	if cbc.SuccessThreshold <= 0 {
		return errors.New("success threshold must be positive")
	}
	if cbc.RequestVolumeThreshold <= 0 {
		return errors.New("request volume threshold must be positive")
	}
	if cbc.FailureRateThreshold < 0 || cbc.FailureRateThreshold > 1 {
		return errors.New("failure rate threshold must be between 0.0 and 1.0")
	}
	if cbc.ConsecutiveFailureLimit <= 0 {
		return errors.New("consecutive failure limit must be positive")
	}
	return nil
}

// CircuitBreakerMetrics holds metrics for circuit breaker monitoring.
type CircuitBreakerMetrics struct {
	TotalRequests       int64     // Total number of requests
	SuccessfulRequests  int64     // Number of successful requests
	FailedRequests      int64     // Number of failed requests
	ConsecutiveFailures int       // Current consecutive failures
	LastFailureTime     time.Time // Time of last failure
	LastSuccessTime     time.Time // Time of last success
	StateTransitions    int64     // Number of state transitions
}

// FailureRate calculates the current failure rate.
func (cbm *CircuitBreakerMetrics) FailureRate() float64 {
	if cbm.TotalRequests == 0 {
		return 0.0
	}
	return float64(cbm.FailedRequests) / float64(cbm.TotalRequests)
}

// CircuitBreaker implements the circuit breaker pattern for plugin reliability.
type CircuitBreaker struct {
	name      string
	state     CircuitBreakerState
	config    *CircuitBreakerConfig
	metrics   *CircuitBreakerMetrics
	stateTime time.Time // Time of last state change
}

// NewCircuitBreaker creates a new circuit breaker with the given configuration.
func NewCircuitBreaker(name string, config *CircuitBreakerConfig) (*CircuitBreaker, error) {
	if config == nil {
		config = NewDefaultCircuitBreakerConfig()
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid circuit breaker config: %w", err)
	}

	return &CircuitBreaker{
		name:      name,
		state:     CircuitClosed,
		config:    config,
		metrics:   &CircuitBreakerMetrics{},
		stateTime: time.Now(),
	}, nil
}

// NewDefaultCircuitBreaker creates a circuit breaker with default configuration.
func NewDefaultCircuitBreaker(name string) *CircuitBreaker {
	cb, _ := NewCircuitBreaker(name, NewDefaultCircuitBreakerConfig()) // Default config is always valid
	return cb
}

// Name returns the circuit breaker name.
func (cb *CircuitBreaker) Name() string {
	return cb.name
}

// State returns the current circuit breaker state.
func (cb *CircuitBreaker) State() CircuitBreakerState {
	return cb.state
}

// Metrics returns a copy of the current metrics.
func (cb *CircuitBreaker) Metrics() CircuitBreakerMetrics {
	return *cb.metrics // Return copy to prevent external modification
}

// IsRequestAllowed determines if a request should be allowed based on circuit state.
func (cb *CircuitBreaker) IsRequestAllowed() bool {
	switch cb.state {
	case CircuitClosed:
		return true
	case CircuitOpen:
		// Check if recovery timeout has passed
		if time.Since(cb.stateTime) >= cb.config.RecoveryTimeout {
			cb.setState(CircuitHalfOpen)
			return true
		}
		return false
	case CircuitHalfOpen:
		// Allow limited requests to test if service has recovered
		return true
	default:
		return false
	}
}

// RecordSuccess records a successful request.
func (cb *CircuitBreaker) RecordSuccess() {
	cb.metrics.TotalRequests++
	cb.metrics.SuccessfulRequests++
	cb.metrics.ConsecutiveFailures = 0
	cb.metrics.LastSuccessTime = time.Now()

	// Check if we should close the circuit
	if cb.state == CircuitHalfOpen {
		// Count consecutive successes since half-open
		if cb.metrics.SuccessfulRequests >= int64(cb.config.SuccessThreshold) {
			cb.setState(CircuitClosed)
		}
	}
}

// RecordFailure records a failed request and updates circuit state if necessary.
func (cb *CircuitBreaker) RecordFailure(_ error) {
	cb.metrics.TotalRequests++
	cb.metrics.FailedRequests++
	cb.metrics.ConsecutiveFailures++
	cb.metrics.LastFailureTime = time.Now()

	// Check if circuit should be opened
	cb.evaluateCircuitState()
}

// evaluateCircuitState checks if the circuit should be opened based on failure metrics.
func (cb *CircuitBreaker) evaluateCircuitState() {
	// Don't evaluate if we don't have enough requests
	if cb.metrics.TotalRequests < int64(cb.config.RequestVolumeThreshold) {
		return
	}

	// Check consecutive failures
	shouldOpen := cb.metrics.ConsecutiveFailures >= cb.config.ConsecutiveFailureLimit ||
		// Check failure threshold
		cb.metrics.FailedRequests >= int64(cb.config.FailureThreshold) ||
		// Check failure rate
		cb.metrics.FailureRate() >= cb.config.FailureRateThreshold

	if shouldOpen && cb.state != CircuitOpen {
		cb.setState(CircuitOpen)
	}
}

// setState changes the circuit breaker state and updates metrics.
func (cb *CircuitBreaker) setState(newState CircuitBreakerState) {
	if cb.state != newState {
		cb.state = newState
		cb.stateTime = time.Now()
		cb.metrics.StateTransitions++

		// Reset success counter when entering half-open state
		if newState == CircuitHalfOpen {
			cb.metrics.SuccessfulRequests = 0
		}

		// Reset metrics when closing circuit
		if newState == CircuitClosed {
			cb.resetMetrics()
		}
	}
}

// resetMetrics resets the circuit breaker metrics.
func (cb *CircuitBreaker) resetMetrics() {
	cb.metrics.TotalRequests = 0
	cb.metrics.SuccessfulRequests = 0
	cb.metrics.FailedRequests = 0
	cb.metrics.ConsecutiveFailures = 0
	// Keep LastFailureTime and LastSuccessTime for monitoring
	// Keep StateTransitions for monitoring
}

// Execute wraps a function call with circuit breaker logic.
func (cb *CircuitBreaker) Execute(fn func() error) error {
	if !cb.IsRequestAllowed() {
		return NewTransientError(
			ErrorCodeCircuitOpen,
			fmt.Sprintf("Circuit breaker '%s' is open", cb.name),
			&cb.config.RecoveryTimeout,
		)
	}

	err := fn()
	if err != nil {
		cb.RecordFailure(err)
		return err
	}

	cb.RecordSuccess()
	return nil
}

// ForceOpen forces the circuit breaker to open state.
func (cb *CircuitBreaker) ForceOpen() {
	cb.setState(CircuitOpen)
}

// ForceClose forces the circuit breaker to closed state and resets metrics.
func (cb *CircuitBreaker) ForceClose() {
	cb.setState(CircuitClosed)
}

// String returns a string representation of the circuit breaker state.
func (cb *CircuitBreaker) String() string {
	var stateStr string
	switch cb.state {
	case CircuitClosed:
		stateStr = "CLOSED"
	case CircuitOpen:
		stateStr = "OPEN"
	case CircuitHalfOpen:
		stateStr = "HALF_OPEN"
	default:
		stateStr = "UNKNOWN"
	}

	return fmt.Sprintf("CircuitBreaker{name=%s, state=%s, failures=%d, rate=%.2f}",
		cb.name, stateStr, cb.metrics.ConsecutiveFailures, cb.metrics.FailureRate())
}

// TimeoutConfig defines timeout settings for RPC methods.
type TimeoutConfig struct {
	NameTimeout             time.Duration // Timeout for Name RPC
	SupportsTimeout         time.Duration // Timeout for Supports RPC
	GetActualCostTimeout    time.Duration // Timeout for GetActualCost RPC
	GetProjectedCostTimeout time.Duration // Timeout for GetProjectedCost RPC
	GetPricingSpecTimeout   time.Duration // Timeout for GetPricingSpec RPC
	GlobalTimeout           time.Duration // Global timeout for all operations
}

// NewDefaultTimeoutConfig creates a timeout config with sensible defaults.
func NewDefaultTimeoutConfig() *TimeoutConfig {
	return &TimeoutConfig{
		NameTimeout:             defaultNameTimeout,
		SupportsTimeout:         defaultSupportsTimeout,
		GetActualCostTimeout:    defaultGetActualCostTimeout,
		GetProjectedCostTimeout: defaultGetProjectedCostTimeout,
		GetPricingSpecTimeout:   defaultGetPricingSpecTimeout,
		GlobalTimeout:           defaultGlobalTimeout,
	}
}

// NewFastTimeoutConfig creates a timeout config optimized for fast responses.
func NewFastTimeoutConfig() *TimeoutConfig {
	return &TimeoutConfig{
		NameTimeout:             fastNameTimeout,
		SupportsTimeout:         fastSupportsTimeout,
		GetActualCostTimeout:    fastGetActualCostTimeout,
		GetProjectedCostTimeout: fastGetProjectedCostTimeout,
		GetPricingSpecTimeout:   fastGetPricingSpecTimeout,
		GlobalTimeout:           fastGlobalTimeout,
	}
}

// NewSlowTimeoutConfig creates a timeout config for slower or batch operations.
func NewSlowTimeoutConfig() *TimeoutConfig {
	return &TimeoutConfig{
		NameTimeout:             slowNameTimeout,
		SupportsTimeout:         slowSupportsTimeout,
		GetActualCostTimeout:    slowGetActualCostTimeout,
		GetProjectedCostTimeout: slowGetProjectedCostTimeout,
		GetPricingSpecTimeout:   slowGetPricingSpecTimeout,
		GlobalTimeout:           slowGlobalTimeout,
	}
}

// Validate checks if the timeout config has valid values.
func (tc *TimeoutConfig) Validate() error {
	timeouts := map[string]time.Duration{
		"NameTimeout":             tc.NameTimeout,
		"SupportsTimeout":         tc.SupportsTimeout,
		"GetActualCostTimeout":    tc.GetActualCostTimeout,
		"GetProjectedCostTimeout": tc.GetProjectedCostTimeout,
		"GetPricingSpecTimeout":   tc.GetPricingSpecTimeout,
		"GlobalTimeout":           tc.GlobalTimeout,
	}

	for name, timeout := range timeouts {
		if timeout <= 0 {
			return fmt.Errorf("%s must be positive", name)
		}
		if timeout < minimumTimeout {
			return fmt.Errorf("%s must be at least %v", name, minimumTimeout)
		}
		if timeout > maximumTimeout {
			return fmt.Errorf("%s must be at most %v", name, maximumTimeout)
		}
	}

	// Validate that global timeout is larger than individual method timeouts
	methodTimeouts := []time.Duration{
		tc.NameTimeout, tc.SupportsTimeout, tc.GetActualCostTimeout,
		tc.GetProjectedCostTimeout, tc.GetPricingSpecTimeout,
	}

	for _, methodTimeout := range methodTimeouts {
		if tc.GlobalTimeout < methodTimeout {
			return fmt.Errorf("GlobalTimeout (%v) must be at least as large as the largest method timeout (%v)",
				tc.GlobalTimeout, methodTimeout)
		}
	}

	return nil
}

// GetTimeoutForMethod returns the appropriate timeout for a given RPC method.
func (tc *TimeoutConfig) GetTimeoutForMethod(method string) time.Duration {
	switch method {
	case "Name":
		return tc.NameTimeout
	case "Supports":
		return tc.SupportsTimeout
	case "GetActualCost":
		return tc.GetActualCostTimeout
	case "GetProjectedCost":
		return tc.GetProjectedCostTimeout
	case "GetPricingSpec":
		return tc.GetPricingSpecTimeout
	default:
		return tc.GlobalTimeout
	}
}

// TimeoutWrapper provides timeout-aware execution for functions.
type TimeoutWrapper struct {
	config *TimeoutConfig
}

// NewTimeoutWrapper creates a new timeout wrapper with the given configuration.
func NewTimeoutWrapper(config *TimeoutConfig) (*TimeoutWrapper, error) {
	if config == nil {
		config = NewDefaultTimeoutConfig()
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid timeout config: %w", err)
	}

	return &TimeoutWrapper{config: config}, nil
}

// NewDefaultTimeoutWrapper creates a timeout wrapper with default configuration.
func NewDefaultTimeoutWrapper() *TimeoutWrapper {
	wrapper, _ := NewTimeoutWrapper(NewDefaultTimeoutConfig()) // Default config is always valid
	return wrapper
}

// ExecuteWithTimeout executes a function with a method-specific timeout.
func (tw *TimeoutWrapper) ExecuteWithTimeout(ctx context.Context, method string, fn func(context.Context) error) error {
	timeout := tw.config.GetTimeoutForMethod(method)

	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Execute function in a goroutine to handle timeout
	errChan := make(chan error, 1)
	go func() {
		errChan <- fn(timeoutCtx)
	}()

	// Wait for completion or timeout
	select {
	case err := <-errChan:
		return err
	case <-timeoutCtx.Done():
		if timeoutCtx.Err() == context.DeadlineExceeded {
			return NewTransientError(
				ErrorCodeNetworkTimeout,
				fmt.Sprintf("Operation '%s' timed out after %v", method, timeout),
				nil,
			)
		}
		return timeoutCtx.Err()
	}
}

// ExecuteWithGlobalTimeout executes a function with the global timeout.
func (tw *TimeoutWrapper) ExecuteWithGlobalTimeout(ctx context.Context, fn func(context.Context) error) error {
	return tw.ExecuteWithTimeout(ctx, "", fn) // Empty method name uses global timeout
}

// CreateTimeoutContext creates a context with timeout for the specified method.
func (tw *TimeoutWrapper) CreateTimeoutContext(
	ctx context.Context,
	method string,
) (context.Context, context.CancelFunc) {
	timeout := tw.config.GetTimeoutForMethod(method)
	return context.WithTimeout(ctx, timeout)
}

// GetConfig returns a copy of the current timeout configuration.
func (tw *TimeoutWrapper) GetConfig() TimeoutConfig {
	return *tw.config // Return copy to prevent external modification
}

// UpdateConfig updates the timeout configuration after validation.
func (tw *TimeoutWrapper) UpdateConfig(config *TimeoutConfig) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid timeout config: %w", err)
	}
	tw.config = config
	return nil
}

// TimeoutAwareRetry combines retry logic with timeout handling.
func TimeoutAwareRetry(
	ctx context.Context,
	timeoutWrapper *TimeoutWrapper,
	retryPolicy *RetryPolicy,
	method string,
	fn func(context.Context) error,
) error {
	if timeoutWrapper == nil {
		timeoutWrapper = NewDefaultTimeoutWrapper()
	}
	if retryPolicy == nil {
		retryPolicy = NewDefaultRetryPolicy()
	}

	return RetryWithPolicy(ctx, retryPolicy, func() error {
		return timeoutWrapper.ExecuteWithTimeout(ctx, method, fn)
	})
}

// Proto Error Integration Functions

// ToProto converts ErrorCategory to protobuf ErrorCategory.
func (ec ErrorCategory) ToProto() pbc.ErrorCategory {
	switch ec {
	case TransientError:
		return pbc.ErrorCategory_ERROR_CATEGORY_TRANSIENT
	case PermanentError:
		return pbc.ErrorCategory_ERROR_CATEGORY_PERMANENT
	case ConfigurationError:
		return pbc.ErrorCategory_ERROR_CATEGORY_CONFIGURATION
	default:
		return pbc.ErrorCategory_ERROR_CATEGORY_UNSPECIFIED
	}
}

// FromProtoErrorCategory converts protobuf ErrorCategory to ErrorCategory.
func FromProtoErrorCategory(protoCategory pbc.ErrorCategory) ErrorCategory {
	switch protoCategory {
	case pbc.ErrorCategory_ERROR_CATEGORY_TRANSIENT:
		return TransientError
	case pbc.ErrorCategory_ERROR_CATEGORY_PERMANENT:
		return PermanentError
	case pbc.ErrorCategory_ERROR_CATEGORY_CONFIGURATION:
		return ConfigurationError
	case pbc.ErrorCategory_ERROR_CATEGORY_UNSPECIFIED:
		return ""
	default:
		return ""
	}
}

// ToProto converts ErrorCode to protobuf ErrorCode.
func (ec ErrorCode) ToProto() pbc.ErrorCode {
	switch ec {
	case ErrorCodeNetworkTimeout:
		return pbc.ErrorCode_ERROR_CODE_NETWORK_TIMEOUT
	case ErrorCodeServiceUnavailable:
		return pbc.ErrorCode_ERROR_CODE_SERVICE_UNAVAILABLE
	case ErrorCodeRateLimited:
		return pbc.ErrorCode_ERROR_CODE_RATE_LIMITED
	case ErrorCodeTemporaryFailure:
		return pbc.ErrorCode_ERROR_CODE_TEMPORARY_FAILURE
	case ErrorCodeCircuitOpen:
		return pbc.ErrorCode_ERROR_CODE_CIRCUIT_OPEN
	case ErrorCodeInvalidResource:
		return pbc.ErrorCode_ERROR_CODE_INVALID_RESOURCE
	case ErrorCodeResourceNotFound:
		return pbc.ErrorCode_ERROR_CODE_RESOURCE_NOT_FOUND
	case ErrorCodeInvalidTimeRange:
		return pbc.ErrorCode_ERROR_CODE_INVALID_TIME_RANGE
	case ErrorCodeUnsupportedRegion:
		return pbc.ErrorCode_ERROR_CODE_UNSUPPORTED_REGION
	case ErrorCodePermissionDenied:
		return pbc.ErrorCode_ERROR_CODE_PERMISSION_DENIED
	case ErrorCodeDataCorruption:
		return pbc.ErrorCode_ERROR_CODE_DATA_CORRUPTION
	case ErrorCodeInvalidCredentials:
		return pbc.ErrorCode_ERROR_CODE_INVALID_CREDENTIALS
	case ErrorCodeMissingAPIKey:
		return pbc.ErrorCode_ERROR_CODE_MISSING_API_KEY
	case ErrorCodeInvalidEndpoint:
		return pbc.ErrorCode_ERROR_CODE_INVALID_ENDPOINT
	case ErrorCodeInvalidProvider:
		return pbc.ErrorCode_ERROR_CODE_INVALID_PROVIDER
	case ErrorCodePluginNotConfigured:
		return pbc.ErrorCode_ERROR_CODE_PLUGIN_NOT_CONFIGURED
	default:
		return pbc.ErrorCode_ERROR_CODE_UNSPECIFIED
	}
}

// FromProtoErrorCode converts protobuf ErrorCode to ErrorCode.
func FromProtoErrorCode(protoCode pbc.ErrorCode) ErrorCode {
	switch protoCode {
	case pbc.ErrorCode_ERROR_CODE_UNSPECIFIED:
		return ""
	case pbc.ErrorCode_ERROR_CODE_NETWORK_TIMEOUT:
		return ErrorCodeNetworkTimeout
	case pbc.ErrorCode_ERROR_CODE_SERVICE_UNAVAILABLE:
		return ErrorCodeServiceUnavailable
	case pbc.ErrorCode_ERROR_CODE_RATE_LIMITED:
		return ErrorCodeRateLimited
	case pbc.ErrorCode_ERROR_CODE_TEMPORARY_FAILURE:
		return ErrorCodeTemporaryFailure
	case pbc.ErrorCode_ERROR_CODE_CIRCUIT_OPEN:
		return ErrorCodeCircuitOpen
	case pbc.ErrorCode_ERROR_CODE_INVALID_RESOURCE:
		return ErrorCodeInvalidResource
	case pbc.ErrorCode_ERROR_CODE_RESOURCE_NOT_FOUND:
		return ErrorCodeResourceNotFound
	case pbc.ErrorCode_ERROR_CODE_INVALID_TIME_RANGE:
		return ErrorCodeInvalidTimeRange
	case pbc.ErrorCode_ERROR_CODE_UNSUPPORTED_REGION:
		return ErrorCodeUnsupportedRegion
	case pbc.ErrorCode_ERROR_CODE_PERMISSION_DENIED:
		return ErrorCodePermissionDenied
	case pbc.ErrorCode_ERROR_CODE_DATA_CORRUPTION:
		return ErrorCodeDataCorruption
	case pbc.ErrorCode_ERROR_CODE_INVALID_CREDENTIALS:
		return ErrorCodeInvalidCredentials
	case pbc.ErrorCode_ERROR_CODE_MISSING_API_KEY:
		return ErrorCodeMissingAPIKey
	case pbc.ErrorCode_ERROR_CODE_INVALID_ENDPOINT:
		return ErrorCodeInvalidEndpoint
	case pbc.ErrorCode_ERROR_CODE_INVALID_PROVIDER:
		return ErrorCodeInvalidProvider
	case pbc.ErrorCode_ERROR_CODE_PLUGIN_NOT_CONFIGURED:
		return ErrorCodePluginNotConfigured
	default:
		return ""
	}
}

// ToProtoErrorDetail converts a PluginError to protobuf ErrorDetail.
func (e *PluginError) ToProtoErrorDetail() *pbc.ErrorDetail {
	detail := &pbc.ErrorDetail{
		Code:      e.Code.ToProto(),
		Category:  e.Category.ToProto(),
		Message:   e.Message,
		Details:   make(map[string]string),
		Timestamp: timestamppb.New(e.Timestamp),
	}

	// Convert details from interface{} to string
	for k, v := range e.Details {
		if str, ok := v.(string); ok {
			detail.Details[k] = str
		} else {
			detail.Details[k] = fmt.Sprintf("%v", v)
		}
	}

	// Add retry_after_seconds for transient errors
	if e.RetryAfter != nil && e.Category == TransientError {
		seconds := int32(e.RetryAfter.Seconds())
		detail.RetryAfterSeconds = &seconds
	}

	return detail
}

// FromProtoErrorDetail converts protobuf ErrorDetail to PluginError.
func FromProtoErrorDetail(detail *pbc.ErrorDetail) *PluginError {
	if detail == nil {
		return nil
	}

	pluginErr := &PluginError{
		Code:      FromProtoErrorCode(detail.GetCode()),
		Category:  FromProtoErrorCategory(detail.GetCategory()),
		Message:   detail.GetMessage(),
		Details:   make(map[string]interface{}),
		Timestamp: time.Now(),
		Retryable: FromProtoErrorCategory(detail.GetCategory()) == TransientError,
	}

	// Convert details from string to interface{}
	for k, v := range detail.GetDetails() {
		pluginErr.Details[k] = v
	}

	// Set timestamp if provided
	if detail.GetTimestamp() != nil {
		pluginErr.Timestamp = detail.GetTimestamp().AsTime()
	}

	// Set retry delay if provided
	if detail.GetRetryAfterSeconds() > 0 {
		retryAfter := time.Duration(detail.GetRetryAfterSeconds()) * time.Second
		pluginErr.RetryAfter = &retryAfter
	}

	return pluginErr
}

// CreateProtoErrorDetail creates a protobuf ErrorDetail from components.
func CreateProtoErrorDetail(
	code ErrorCode,
	category ErrorCategory,
	message string,
	details map[string]string,
) *pbc.ErrorDetail {
	detail := &pbc.ErrorDetail{
		Code:      code.ToProto(),
		Category:  category.ToProto(),
		Message:   message,
		Details:   details,
		Timestamp: timestamppb.New(time.Now()),
	}

	if details == nil {
		detail.Details = make(map[string]string)
	}

	return detail
}
