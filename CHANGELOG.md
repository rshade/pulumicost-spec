# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Telemetry and Observability Specification**: Complete observability framework for PulumiCost plugins
  - New `ObservabilityService` gRPC service with health checks, metrics, and SLI endpoints
  - Comprehensive metrics specification covering latency, throughput, error rates, and resource usage
  - OpenTelemetry distributed tracing integration with context propagation
  - Structured logging standards with trace correlation and error classification
  - Service Level Indicators (SLIs) and Service Level Objectives (SLOs) framework
  - Three-tier conformance levels: Basic, Standard, and Advanced
- **Go SDK Observability Support**: Complete observability implementation in Go SDK
  - `observability.go`: Core observability types, metrics, SLIs, and validation functions
  - `observability_validate.go`: Comprehensive validation for telemetry data and log entries
  - Standard metrics constants and SLI definitions with target values
  - Helper functions for error rate calculation, availability measurement, and duration formatting
  - Conformance level requirements and validation framework
- **Testing Framework Extensions**: Observability testing capabilities
  - `ObservabilityTestSuite` for comprehensive conformance testing
  - Basic, Standard, and Advanced observability conformance test implementations
  - Mock interfaces for health checks, metrics collection, and SLI reporting
  - Performance testing under load and error condition validation
- **Comprehensive Examples and Documentation**
  - JSON examples for metrics, SLIs, health checks, and structured logs
  - Grafana dashboard template for Prometheus metrics visualization
  - CloudWatch dashboard CloudFormation template with alerts and custom metrics
  - Complete observability implementation guide with best practices
  - Dashboard templates for common monitoring tools (Grafana, CloudWatch)
- **Protocol Buffer Extensions**: Telemetry message types and service definitions
  - `TelemetryMetadata` for request correlation and processing metrics
  - `LogEntry` and `ErrorDetails` for structured logging
  - `Metric`, `ServiceLevelIndicator`, and health check message types
  - Time range and status enumerations for observability queries
- Comprehensive markdown linting configuration with proper exclusions for node_modules
- Changelog linting with keep-a-changelog tool
- Session learnings and solutions documentation in CLAUDE.md
- Enhanced CI/CD debugging commands and workflow optimizations
- JSON Schema validation improvements with ajv strict mode handling

### Changed

- Updated package.json with markdown and changelog linting commands
- Enhanced validate_examples.js to handle schema compilation issues
- Improved .markdownlintignore configuration for better exclusions

### Fixed

- Markdown linting errors reduced from 950+ to 0 by configuring proper exclusions
- JSON Schema validation issues with invalid `version` keyword
- AJV compilation errors with `$schema` reference resolution
- CI failures due to out-of-sync package-lock.json and go.mod files
- Line length violations and missing language specifiers in code blocks

## [0.1.0] - 2025-08-16

### Added

- Initial release of PulumiCost specification v0.1.0
- Complete gRPC service definitions in costsource.proto
- CostSourceService with 5 RPC methods: Name, Supports, GetActualCost, GetProjectedCost, GetPricingSpec
- Comprehensive message definitions for all request/response pairs
- ResourceDescriptor, ActualCostResult, PricingSpec, and UsageMetricHint messages
- Complete plugin developer guide with implementation examples
- JSON Schema for pricing specifications with 44+ billing models
- Cross-vendor examples for AWS, Azure, GCP, and Kubernetes
- Go SDK with generated protobuf code and helper types
- Comprehensive testing framework with integration and conformance tests

### Changed

- Established v0.1.0 as stable baseline for breaking change detection
- Finalized all service and message definitions

### Fixed

- buf lint compliance with zero errors
- Code formatting and build system compatibility
- Generated Go SDK compilation issues

[Unreleased]: https://github.com/rshade/pulumicost-spec/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/rshade/pulumicost-spec/releases/tag/v0.1.0
