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

package pricing_test

import (
	"testing"
	"time"

	"github.com/rshade/pulumicost-spec/sdk/go/pricing"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// TestErrorCodeMapping tests error code to category mapping.
func TestErrorCodeMapping(t *testing.T) {
	mapping := pricing.GetErrorMapping()

	// Test some key mappings
	testCases := []struct {
		code     pricing.ErrorCode
		category pricing.ErrorCategory
	}{
		{pricing.ErrorCodeNetworkTimeout, pricing.TransientError},
		{pricing.ErrorCodeInvalidResource, pricing.PermanentError},
		{pricing.ErrorCodeInvalidCredentials, pricing.ConfigurationError},
	}

	for _, tc := range testCases {
		if mapping[tc.code] != tc.category {
			t.Errorf("Expected %v to map to %v, got %v", tc.code, tc.category, mapping[tc.code])
		}
	}
}

// TestErrorCreation tests creating different types of errors.
func TestErrorCreation(t *testing.T) {
	t.Run("TransientError", func(t *testing.T) {
		err := pricing.NewTransientError(
			pricing.ErrorCodeServiceUnavailable,
			"Service unavailable",
			nil,
		)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if !pricing.IsTransientError(err) {
			t.Error("Expected error to be transient")
		}

		if err.IsRetryable() != true {
			t.Error("Expected transient error to be retryable")
		}
	})

	t.Run("PermanentError", func(t *testing.T) {
		err := pricing.NewPermanentError(
			pricing.ErrorCodeInvalidResource,
			"Invalid resource",
		)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if !pricing.IsPermanentError(err) {
			t.Error("Expected error to be permanent")
		}

		if err.IsRetryable() != false {
			t.Error("Expected permanent error to not be retryable")
		}
	})

	t.Run("ConfigurationError", func(t *testing.T) {
		err := pricing.NewConfigurationError(
			pricing.ErrorCodeInvalidCredentials,
			"Invalid credentials",
		)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if !pricing.IsConfigurationError(err) {
			t.Error("Expected error to be configuration error")
		}

		if err.IsRetryable() != false {
			t.Error("Expected configuration error to not be retryable")
		}
	})
}

// TestErrorProtoConversion tests conversion between error types and protobuf.
func TestErrorProtoConversion(t *testing.T) {
	// Test error code conversion
	testCases := []struct {
		errorCode pricing.ErrorCode
		protoCode pbc.ErrorCode
	}{
		{pricing.ErrorCodeNetworkTimeout, pbc.ErrorCode_ERROR_CODE_NETWORK_TIMEOUT},
		{pricing.ErrorCodeInvalidResource, pbc.ErrorCode_ERROR_CODE_INVALID_RESOURCE},
		{pricing.ErrorCodeInvalidCredentials, pbc.ErrorCode_ERROR_CODE_INVALID_CREDENTIALS},
	}

	for _, tc := range testCases {
		// Test ErrorCode -> Proto conversion
		protoCode := tc.errorCode.ToProto()
		if protoCode != tc.protoCode {
			t.Errorf("Expected proto code %v, got %v", tc.protoCode, protoCode)
		}

		// Test Proto -> ErrorCode conversion
		errorCode := pricing.FromProtoErrorCode(tc.protoCode)
		if errorCode != tc.errorCode {
			t.Errorf("Expected error code %v, got %v", tc.errorCode, errorCode)
		}
	}
}

// TestErrorCategoryProtoConversion tests category conversion.
func TestErrorCategoryProtoConversion(t *testing.T) {
	testCases := []struct {
		category      pricing.ErrorCategory
		protoCategory pbc.ErrorCategory
	}{
		{pricing.TransientError, pbc.ErrorCategory_ERROR_CATEGORY_TRANSIENT},
		{pricing.PermanentError, pbc.ErrorCategory_ERROR_CATEGORY_PERMANENT},
		{pricing.ConfigurationError, pbc.ErrorCategory_ERROR_CATEGORY_CONFIGURATION},
	}

	for _, tc := range testCases {
		// Test category -> proto conversion
		protoCategory := tc.category.ToProto()
		if protoCategory != tc.protoCategory {
			t.Errorf("Expected proto category %v, got %v", tc.protoCategory, protoCategory)
		}

		// Test proto -> category conversion
		category := pricing.FromProtoErrorCategory(tc.protoCategory)
		if category != tc.category {
			t.Errorf("Expected category %v, got %v", tc.category, category)
		}
	}
}

// TestPluginErrorDetails tests PluginError structure and methods.
func TestPluginErrorDetails(t *testing.T) {
	retryAfter := 30 * time.Second
	pluginErr := pricing.NewTransientError(
		pricing.ErrorCodeRateLimited,
		"Rate limit exceeded",
		&retryAfter,
	)

	// Add some details
	pluginErr = pluginErr.WithDetails(map[string]interface{}{
		"limit":  "100",
		"window": "60s",
	})

	// Test basic properties (the error message includes category and code formatting)
	errorMsg := pluginErr.Error()
	if errorMsg == "" {
		t.Error("Expected non-empty error message")
	}
	// Should contain the original message
	if !contains(errorMsg, "Rate limit exceeded") {
		t.Errorf("Expected error message to contain 'Rate limit exceeded', got %s", errorMsg)
	}

	if !pluginErr.IsRetryable() {
		t.Error("Expected error to be retryable")
	}

	if pluginErr.GetRetryAfter() == nil || *pluginErr.GetRetryAfter() != retryAfter {
		t.Errorf("Expected retry after %v, got %v", retryAfter, pluginErr.GetRetryAfter())
	}
}

// TestPluginErrorProtoRoundTrip tests full round-trip conversion.
func TestPluginErrorProtoRoundTrip(t *testing.T) {
	retryAfter := 45 * time.Second
	originalErr := pricing.NewTransientError(
		pricing.ErrorCodeServiceUnavailable,
		"Service temporarily unavailable",
		&retryAfter,
	).WithDetails(map[string]interface{}{
		"service": "pricing-api",
		"region":  "us-east-1",
	})

	// Convert to proto
	protoDetail := originalErr.ToProtoErrorDetail()

	// Validate key proto fields
	if protoDetail.GetCode() != pbc.ErrorCode_ERROR_CODE_SERVICE_UNAVAILABLE {
		t.Errorf("Expected SERVICE_UNAVAILABLE code, got %v", protoDetail.GetCode())
	}

	if protoDetail.GetMessage() != "Service temporarily unavailable" {
		t.Errorf("Expected specific message, got %s", protoDetail.GetMessage())
	}

	if protoDetail.GetRetryAfterSeconds() != 45 {
		t.Errorf("Expected retry after 45s, got %d", protoDetail.GetRetryAfterSeconds())
	}

	// Convert back from proto
	convertedErr := pricing.FromProtoErrorDetail(protoDetail)
	if convertedErr == nil {
		t.Fatal("Expected converted error, got nil")
	}

	if convertedErr.Error() != originalErr.Error() {
		t.Errorf("Expected message %s, got %s", originalErr.Error(), convertedErr.Error())
	}

	if !convertedErr.IsRetryable() {
		t.Error("Expected converted error to be retryable")
	}
}

// TestRetryPolicyBasics tests basic retry policy functionality.
func TestRetryPolicyBasics(t *testing.T) {
	policy := pricing.NewDefaultRetryPolicy()

	if policy == nil {
		t.Fatal("Expected policy to be created, got nil")
	}

	// Test that it can validate itself
	if err := policy.Validate(); err != nil {
		t.Errorf("Expected valid policy, got error: %v", err)
	}

	// Test different configurations
	aggressivePolicy := pricing.NewAggressiveRetryPolicy()
	if aggressivePolicy == nil {
		t.Fatal("Expected aggressive policy to be created, got nil")
	}

	conservativePolicy := pricing.NewConservativeRetryPolicy()
	if conservativePolicy == nil {
		t.Fatal("Expected conservative policy to be created, got nil")
	}
}

// TestCircuitBreakerBasics tests basic circuit breaker functionality.
func TestCircuitBreakerBasics(t *testing.T) {
	breaker := pricing.NewDefaultCircuitBreaker("test-breaker")

	if breaker == nil {
		t.Fatal("Expected circuit breaker to be created, got nil")
	}

	// Test name
	if breaker.Name() != "test-breaker" {
		t.Errorf("Expected name 'test-breaker', got %s", breaker.Name())
	}

	// Test initial state
	if breaker.State() != pricing.CircuitClosed {
		t.Errorf("Expected initial state Closed, got %v", breaker.State())
	}

	// Test metrics
	metrics := breaker.Metrics()
	if metrics.TotalRequests != 0 {
		t.Errorf("Expected 0 total requests initially, got %d", metrics.TotalRequests)
	}
}

// TestTimeoutConfigBasics tests basic timeout configuration functionality.
func TestTimeoutConfigBasics(t *testing.T) {
	config := pricing.NewDefaultTimeoutConfig()

	if config == nil {
		t.Fatal("Expected timeout config to be created, got nil")
	}

	// Test that it returns reasonable timeouts
	nameTimeout := config.GetTimeoutForMethod("Name")
	if nameTimeout <= 0 {
		t.Errorf("Expected positive timeout for Name method, got %v", nameTimeout)
	}

	actualCostTimeout := config.GetTimeoutForMethod("GetActualCost")
	if actualCostTimeout <= nameTimeout {
		t.Error("Expected GetActualCost timeout to be longer than Name timeout")
	}

	// Test different configurations
	fastConfig := pricing.NewFastTimeoutConfig()
	fastTimeout := fastConfig.GetTimeoutForMethod("Name")

	slowConfig := pricing.NewSlowTimeoutConfig()
	slowTimeout := slowConfig.GetTimeoutForMethod("Name")

	if fastTimeout >= slowTimeout {
		t.Error("Expected fast config to have shorter timeout than slow config")
	}
}

// TestCreateProtoErrorDetail tests direct proto error creation.
func TestCreateProtoErrorDetail(t *testing.T) {
	details := map[string]string{
		"resource": "test-resource",
		"region":   "us-west-2",
	}

	protoDetail := pricing.CreateProtoErrorDetail(
		pricing.ErrorCodeResourceNotFound,
		pricing.PermanentError,
		"Resource not found",
		details,
	)

	if protoDetail.GetCode() != pbc.ErrorCode_ERROR_CODE_RESOURCE_NOT_FOUND {
		t.Errorf("Expected RESOURCE_NOT_FOUND code, got %v", protoDetail.GetCode())
	}

	if protoDetail.GetMessage() != "Resource not found" {
		t.Errorf("Expected 'Resource not found' message, got %s", protoDetail.GetMessage())
	}

	if len(protoDetail.GetDetails()) != 2 {
		t.Errorf("Expected 2 details, got %d", len(protoDetail.GetDetails()))
	}

	// Test with nil details
	protoDetailNil := pricing.CreateProtoErrorDetail(
		pricing.ErrorCodeInvalidCredentials,
		pricing.ConfigurationError,
		"Invalid credentials",
		nil,
	)

	if protoDetailNil.GetDetails() == nil {
		t.Error("Expected empty details map, got nil")
	}

	if len(protoDetailNil.GetDetails()) != 0 {
		t.Errorf("Expected empty details map, got %d items", len(protoDetailNil.GetDetails()))
	}
}

// TestFormatErrorMessage tests error message formatting.
func TestFormatErrorMessage(t *testing.T) {
	params := map[string]string{
		"resource": "test-resource",
		"region":   "us-east-1",
	}

	message := pricing.FormatErrorMessage(pricing.ErrorCodeResourceNotFound, params)

	if message == "" {
		t.Error("Expected non-empty formatted message")
	}

	// Message should contain the resource from params
	if len(message) < 10 { // Basic sanity check
		t.Errorf("Expected reasonable message length, got %d chars: %s", len(message), message)
	}
}

// TestStandardErrorDetails tests standard error detail creation.
func TestStandardErrorDetails(t *testing.T) {
	details := pricing.StandardErrorDetails(
		"GetPricing",
		"ec2",
		"i-1234567890abcdef0",
		"us-east-1",
	)

	if details == nil {
		t.Fatal("Expected details map, got nil")
	}

	if len(details) == 0 {
		t.Error("Expected non-empty details map")
	}

	// Should contain the provided values
	expectedKeys := []string{"operation", "resource_type", "resource_id", "region"}
	for _, key := range expectedKeys {
		if _, exists := details[key]; !exists {
			t.Errorf("Expected details to contain key %s", key)
		}
	}
}

// contains checks if a string contains a substring.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && findInString(s, substr)
}

// findInString searches for substring in string.
func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
