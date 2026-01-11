// Package contracts defines the public API for the Plugin Conformance Test Suite.
// This file documents the Go interfaces and types that plugin developers will use.
//
// NOTE: This is a design document, not compilable code. The actual implementation
// will be in sdk/go/testing/ package.
package contracts

import (
	"encoding/json"
	"time"

	pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

// ConformanceLevel defines the certification level for plugin validation.
type ConformanceLevel int

const (
	// ConformanceLevelBasic - Core functionality, required for all plugins.
	// Tests: Name validation, Supports handling, basic GetProjectedCost/GetPricingSpec.
	ConformanceLevelBasic ConformanceLevel = iota

	// ConformanceLevelStandard - Production readiness, recommended for deployment.
	// Includes Basic tests plus: error handling, consistency, 24h data, 10 concurrent.
	ConformanceLevelStandard

	// ConformanceLevelAdvanced - High performance, for demanding environments.
	// Includes Standard tests plus: latency thresholds, 50 concurrent, 30-day data.
	ConformanceLevelAdvanced
)

// TestCategory groups related conformance tests.
type TestCategory string

const (
	CategorySpecValidation TestCategory = "spec_validation"
	CategoryRPCCorrectness TestCategory = "rpc_correctness"
	CategoryPerformance    TestCategory = "performance"
	CategoryConcurrency    TestCategory = "concurrency"
)

// Configuration constants for default values.
const (
	defaultTimeoutSeconds           = 60
	defaultParallelRequests         = 10
	defaultBenchmarkDurationSeconds = 5

	// Performance baseline constants (milliseconds).
	nameStandardLatencyMs          = 100
	nameAdvancedLatencyMs          = 50
	supportsStandardLatencyMs      = 50
	supportsAdvancedLatencyMs      = 25
	projectedCostStandardLatencyMs = 200
	projectedCostAdvancedLatencyMs = 100
	pricingSpecStandardLatencyMs   = 200
	pricingSpecAdvancedLatencyMs   = 100
	actualCost24hStandardLatencyMs = 2000
	actualCost24hAdvancedLatencyMs = 1000
	actualCost30dAdvancedLatencyMs = 10000
)

// SuiteConfig configures conformance suite execution.
type SuiteConfig struct {
	// TargetLevel is the conformance level to validate against.
	// Default: ConformanceLevelStandard
	TargetLevel ConformanceLevel

	// Timeout is the maximum duration for each individual test.
	// Default: 60 * time.Second
	Timeout time.Duration

	// ParallelRequests is the number of concurrent requests for concurrency tests.
	// Default: 10 (Standard), 50 (Advanced)
	ParallelRequests int

	// EnableBenchmarks controls whether performance benchmarks are run.
	// Default: true
	EnableBenchmarks bool

	// BenchmarkDuration is how long to run each benchmark.
	// Default: 5 * time.Second
	BenchmarkDuration time.Duration
}

// DefaultConfig returns the default suite configuration.
func DefaultConfig() SuiteConfig {
	return SuiteConfig{
		TargetLevel:       ConformanceLevelStandard,
		Timeout:           defaultTimeoutSeconds * time.Second,
		ParallelRequests:  defaultParallelRequests,
		EnableBenchmarks:  true,
		BenchmarkDuration: defaultBenchmarkDurationSeconds * time.Second,
	}
}

// ConformanceSuite is the main entry point for running conformance tests.
type ConformanceSuite interface {
	// Run executes all conformance tests against the plugin implementation.
	// Returns a structured result suitable for JSON serialization.
	Run(impl pbc.CostSourceServiceServer) (*ConformanceResult, error)

	// RunCategory executes tests for a specific category only.
	RunCategory(impl pbc.CostSourceServiceServer, category TestCategory) (*CategoryResult, error)

	// SetConfig updates the suite configuration.
	SetConfig(config SuiteConfig)

	// GetConfig returns the current suite configuration.
	GetConfig() SuiteConfig
}

// NewConformanceSuite creates a new conformance suite with default configuration.
func NewConformanceSuite() ConformanceSuite {
	// Implementation in sdk/go/testing/conformance.go
	return nil
}

// NewConformanceSuiteWithConfig creates a new conformance suite with custom configuration.
func NewConformanceSuiteWithConfig(_ SuiteConfig) ConformanceSuite {
	// Implementation in sdk/go/testing/conformance.go
	return nil
}

// ConformanceResult is the complete result of suite execution.
// This struct is JSON-serializable for CI/CD integration.
type ConformanceResult struct {
	// Version is the report schema version (e.g., "1.0.0").
	Version string `json:"version"`

	// Timestamp is when the suite was executed.
	Timestamp time.Time `json:"timestamp"`

	// PluginName is the name returned by the plugin's Name() RPC.
	PluginName string `json:"plugin_name"`

	// LevelAchieved is the highest conformance level passed.
	LevelAchieved string `json:"level_achieved"`

	// Summary contains aggregate test counts.
	Summary ResultSummary `json:"summary"`

	// Categories contains results organized by test category.
	Categories map[TestCategory]*CategoryResult `json:"categories"`

	// Duration is the total execution time.
	Duration string `json:"duration"`
}

// ToJSON serializes the result to JSON format.
func (r *ConformanceResult) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

// Passed returns true if all tests passed at the target level.
func (r *ConformanceResult) Passed() bool {
	return r.Summary.Failed == 0
}

// ResultSummary contains aggregate test counts.
type ResultSummary struct {
	Total   int `json:"total"`
	Passed  int `json:"passed"`
	Failed  int `json:"failed"`
	Skipped int `json:"skipped"`
}

// CategoryResult contains results for a single test category.
type CategoryResult struct {
	Name    TestCategory `json:"name"`
	Passed  int          `json:"passed"`
	Failed  int          `json:"failed"`
	Skipped int          `json:"skipped"`
	Results []TestResult `json:"results"`
}

// TestResult contains the result of a single test execution.
type TestResult struct {
	Name     string       `json:"name"`
	Method   string       `json:"method"`
	Category TestCategory `json:"category"`
	Success  bool         `json:"success"`
	Error    string       `json:"error,omitempty"`
	Duration string       `json:"duration"`
	Details  string       `json:"details,omitempty"`
}

// ValidationError provides field-level error details for spec validation failures.
type ValidationError struct {
	Field    string      `json:"field"`
	Value    interface{} `json:"value"`
	Expected string      `json:"expected"`
	Message  string      `json:"message"`
}

// Error implements the error interface.
func (e ValidationError) Error() string {
	// Returns formatted error message
	return ""
}

// PerformanceBaseline defines latency thresholds for performance conformance.
type PerformanceBaseline struct {
	Method          string        `json:"method"`
	StandardLatency time.Duration `json:"standard_latency"`
	AdvancedLatency time.Duration `json:"advanced_latency"`
	MaxAllocBytes   int64         `json:"max_alloc_bytes,omitempty"`
}

// DefaultBaselines returns the canonical performance baselines.
func DefaultBaselines() []PerformanceBaseline {
	return []PerformanceBaseline{
		{
			Method:          "Name",
			StandardLatency: nameStandardLatencyMs * time.Millisecond,
			AdvancedLatency: nameAdvancedLatencyMs * time.Millisecond,
		},
		{
			Method:          "Supports",
			StandardLatency: supportsStandardLatencyMs * time.Millisecond,
			AdvancedLatency: supportsAdvancedLatencyMs * time.Millisecond,
		},
		{
			Method:          "GetProjectedCost",
			StandardLatency: projectedCostStandardLatencyMs * time.Millisecond,
			AdvancedLatency: projectedCostAdvancedLatencyMs * time.Millisecond,
		},
		{
			Method:          "GetPricingSpec",
			StandardLatency: pricingSpecStandardLatencyMs * time.Millisecond,
			AdvancedLatency: pricingSpecAdvancedLatencyMs * time.Millisecond,
		},
		{
			Method:          "GetActualCost_24h",
			StandardLatency: actualCost24hStandardLatencyMs * time.Millisecond,
			AdvancedLatency: actualCost24hAdvancedLatencyMs * time.Millisecond,
		},
		{
			Method:          "GetActualCost_30d",
			StandardLatency: 0,
			AdvancedLatency: actualCost30dAdvancedLatencyMs * time.Millisecond,
		},
	}
}

// Convenience functions for common usage patterns

// RunBasicConformance runs basic conformance tests and returns the result.
func RunBasicConformance(impl pbc.CostSourceServiceServer) (*ConformanceResult, error) {
	suite := NewConformanceSuiteWithConfig(SuiteConfig{
		TargetLevel: ConformanceLevelBasic,
	})
	return suite.Run(impl)
}

// RunStandardConformance runs standard conformance tests and returns the result.
func RunStandardConformance(impl pbc.CostSourceServiceServer) (*ConformanceResult, error) {
	suite := NewConformanceSuiteWithConfig(SuiteConfig{
		TargetLevel: ConformanceLevelStandard,
	})
	return suite.Run(impl)
}

// RunAdvancedConformance runs advanced conformance tests and returns the result.
func RunAdvancedConformance(impl pbc.CostSourceServiceServer) (*ConformanceResult, error) {
	suite := NewConformanceSuiteWithConfig(SuiteConfig{
		TargetLevel: ConformanceLevelAdvanced,
	})
	return suite.Run(impl)
}

// PrintReport prints a human-readable summary of the conformance result.
func PrintReport(_ *ConformanceResult) {
	// Implementation prints formatted report to stdout
}
