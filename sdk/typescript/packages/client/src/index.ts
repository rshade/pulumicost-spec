// Generated proto types and services
// In Connect-ES v2, services are generated directly in *_pb.ts files
export * from "./generated/finfocus/v1/enums_pb.js";
export * from "./generated/finfocus/v1/budget_pb.js";
export * from "./generated/finfocus/v1/costsource_pb.js";
export * from "./generated/finfocus/v1/focus_pb.js";
export * from "./generated/finfocus/v1/registry_pb.js";

// Error handling - our custom ValidationError takes precedence
export { ValidationError } from "./errors/validation-error.js";

// Client implementations
export { CostSourceClient, CostSourceClientConfig } from "./clients/cost-source.js";
export { RegistryClient, ObservabilityClient, ClientConfig } from "./clients/auxiliary.js";

// Builder patterns
export { ResourceDescriptorBuilder } from "./builders/resource-descriptor.js";
export { RecommendationFilterBuilder } from "./builders/recommendation-filter.js";
export { FocusRecordBuilder } from "./builders/focus-record.js";

// Utilities
export { recommendationsIterator } from "./utils/pagination.js";
