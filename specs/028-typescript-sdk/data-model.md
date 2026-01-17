# Data Model & SDK Architecture

**Feature**: TypeScript Client SDK
**Branch**: `028-typescript-sdk`

## Core Clients

### CostSourceClient

The primary entry point for cost operations.

**Methods**:

- `getActualCost(request: GetActualCostRequest): Promise<GetActualCostResponse>`
- `getProjectedCost(request: GetProjectedCostRequest): Promise<GetProjectedCostResponse>`
- `getRecommendations(request: GetRecommendationsRequest): Promise<GetRecommendationsResponse>`
- ... (and 8 others matching CostSourceService)

### ObservabilityClient

Handles plugin health and metrics.

**Methods**:

- `healthCheck()`: Checks plugin status.
- `getMetrics()`: Retrieves Prometheus-style metrics.
- `getServiceLevelIndicators()`: Retrieves SLIs.

### RegistryClient

Manages plugin lifecycle (install, update, remove).

**Methods**:

- `discoverPlugins()`
- `installPlugin()`
- ... (and 6 others)

## Builder Pattern Objects

### ResourceDescriptorBuilder

Fluent interface for constructing `ResourceDescriptor` messages.

```typescript
class ResourceDescriptorBuilder {
  withProvider(provider: string): this;
  withResourceType(type: string): this;
  withRegion(region: string): this;
  withTags(tags: Record<string, string>): this;
  build(): ResourceDescriptor;
}
```

### RecommendationFilterBuilder

Fluent interface for constructing `RecommendationFilter` messages.

```typescript
class RecommendationFilterBuilder {
  forProvider(provider: string): this;
  minSavings(amount: number, currency: string): this;
  withActionType(action: ActionType): this;
  build(): RecommendationFilter;
}
```

### FocusRecordBuilder

Constructs FOCUS 1.2/1.3 compliant cost records.

```typescript
class FocusRecordBuilder {
  withBilledCost(amount: number, currency: string): this;
  withBillingPeriod(start: Date, end: Date): this;
  withResourceId(id: string): this;
  build(): FocusCostRecord;
}
```

## Validation & Error Handling

### ValidationError

Thrown when input fails validation rules (ISO 4217, null checks, FOCUS rules).

```typescript
class ValidationError extends Error {
  constructor(message: string, public field?: string, public code?: string);
}
```

## REST Wrapper Model

The REST Wrapper exposes the RPCs as HTTP POST endpoints.

**Path Pattern**: `/finfocus.v1.<ServiceName>/<MethodName>`
**Body**: JSON representation of the Request Proto.
**Response**: JSON representation of the Response Proto.
