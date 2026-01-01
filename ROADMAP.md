# Strategic Roadmap: pulumicost-spec

## Vision

To provide the definitive, high-performance gRPC specification and Go SDK for cloud cost observability,
centered around the FinOps Foundation's FOCUS standard.

---

## [Past Milestones]

### Protocol & Modeling

- [x] **GreenOps Standardization** ([#176](https://github.com/rshade/pulumicost-spec/issues/176)) -
  Sustainability metrics (carbon footprint, energy utilization) integrated into the core protocol.
- [x] **Recommendation Enhancements**
  ([#173](https://github.com/rshade/pulumicost-spec/issues/173),
  [#171](https://github.com/rshade/pulumicost-spec/issues/171),
  [#166](https://github.com/rshade/pulumicost-spec/issues/166)) -
  Added resource-scoped targets, filtering, and extended action types.
- [x] **FOCUS 1.2 Integration**
  ([#100](https://github.com/rshade/pulumicost-spec/issues/100),
  [#99](https://github.com/rshade/pulumicost-spec/issues/99)) -
  Full schema coverage and builder API for FOCUS 1.2.
- [x] **RPC Expansion**
  ([#125](https://github.com/rshade/pulumicost-spec/issues/125),
  [#123](https://github.com/rshade/pulumicost-spec/issues/123),
  [#90](https://github.com/rshade/pulumicost-spec/issues/90)) -
  Implemented `EstimateCost`, `GetBudgets`, and `GetRecommendations`.
- [x] **Zero-Allocation Enums**
  ([#63](https://github.com/rshade/pulumicost-spec/issues/63),
  [#33](https://github.com/rshade/pulumicost-spec/issues/33)) -
  High-performance validation for domain enums.

### SDK & Tooling

- [x] **Plugin Conformance Suite** ([#109](https://github.com/rshade/pulumicost-spec/issues/109)) -
  Automated testing for plugin implementers to ensure spec compliance.
- [x] **SDK Foundation**
  ([#151](https://github.com/rshade/pulumicost-spec/issues/151),
  [#148](https://github.com/rshade/pulumicost-spec/issues/148),
  [#145](https://github.com/rshade/pulumicost-spec/issues/145),
  [#139](https://github.com/rshade/pulumicost-spec/issues/139)) -
  Centralized environment handling, Prometheus metrics, and unified file-based logging.
- [x] **Orchestration Support**
  ([#181](https://github.com/rshade/pulumicost-spec/issues/181),
  [#143](https://github.com/rshade/pulumicost-spec/issues/143),
  [#126](https://github.com/rshade/pulumicost-spec/issues/126)) -
  Added reflection, `--port` flags, and `fallback-hints`.
- [x] **Multi-Protocol Support (gRPC-Web/Connect)**
  ([#189](https://github.com/rshade/pulumicost-spec/issues/189),
  [#223](https://github.com/rshade/pulumicost-spec/pull/223)) -
  Added connect-go integration enabling gRPC, gRPC-Web, and Connect protocols for browser compatibility.

---

## [Completed Milestones] (Q1 2026)

### Protocol & Modeling

- [x] **Forecasting Primitives** ([#241](https://github.com/rshade/pulumicost-spec/issues/241), [#250](https://github.com/rshade/pulumicost-spec/issues/250)) -
  Added `GrowthType` (Linear, Exponential) and `GrowthRate` to `CostResult` for cost projections.
- [x] **FOCUS 1.3 Migration** ([#199](https://github.com/rshade/pulumicost-spec/issues/199)) -
  Audit new columns and entities in FOCUS 1.3 and update builder APIs.

### SDK & Tooling

- [x] **Plugin Capability Discovery** ([#242](https://github.com/rshade/pulumicost-spec/issues/242)) -
  Implemented `GetPluginInfo` RPC for spec version compatibility and capability advertisement.
- [x] **Plugin Capability Dry Run Mode** ([#248](https://github.com/rshade/pulumicost-spec/issues/248)) -
  Implemented `DryRun` for plugin field mapping discovery.
- [x] **Contextual FinOps Validation** ([#201](https://github.com/rshade/pulumicost-spec/issues/201)) -
  Extended `pluginsdk/validation` to include contextual checks.
- [x] **Advanced SDK Patterns** ([#213](https://github.com/rshade/pulumicost-spec/issues/213)) -
  Implementation examples for complex tiered pricing and multi-provider mapping.
- [x] **SDK Documentation Overhaul** ([#243](https://github.com/rshade/pulumicost-spec/issues/243)) -
  Comprehensive Godoc, thread safety, rate limiting, and performance documentation.

---

## [Immediate Focus] (Q1 2026)

### Stability & Maintenance

- [ ] **Dependency Management** ([#13](https://github.com/rshade/pulumicost-spec/issues/13)) -
  Automated tracking and updating of core proto and SDK dependencies.

### Planned Features

- [ ] **Pricing Tier Intelligence** ([#217](https://github.com/rshade/pulumicost-spec/issues/217))
  - Define `PricingTier` enum (`STANDARD`, `SPOT`, `RESERVED`)
  - Define `SpotRisk` enum (`LOW`, `MEDIUM`, `HIGH`) for interruption probability
- [ ] **Usage Profiles** ([#218](https://github.com/rshade/pulumicost-spec/issues/218))
  - Add `UsageProfile` context to requests to allow plugins to adjust recommendations
    based on environment (e.g., Dev vs Prod)

### In Progress

- [/] **Standardized Benchmark Suite**
  ([#113](https://github.com/rshade/pulumicost-spec/issues/113),
  [#142](https://github.com/rshade/pulumicost-spec/issues/142)) -
  Formalizing "Time to First Byte" and memory allocation benchmarks for plugins.

### Connect Protocol Enhancements

- [ ] **DismissRecommendation Consistency** ([#225](https://github.com/rshade/pulumicost-spec/issues/225)) -
  Make DismissRecommendation follow the same pattern as other RPC methods.
- [ ] **Configurable Client Timeouts** ([#226](https://github.com/rshade/pulumicost-spec/issues/226)) -
  Add per-request timeout configuration for client operations.
- [ ] **Connect Protocol Test Coverage** ([#227](https://github.com/rshade/pulumicost-spec/issues/227)) -
  Expand test coverage for edge cases (concurrent requests, large payloads, shutdown).
- [ ] **CORS Configuration** ([#228](https://github.com/rshade/pulumicost-spec/issues/228),
  [#229](https://github.com/rshade/pulumicost-spec/issues/229)) -
  Make CORS headers and max-age configurable.

---

## [Future Vision] (Long-Term)

### Active Research

- [ ] **Standardized Cost Allocation Lineage Metadata**
  ([#191](https://github.com/rshade/pulumicost-spec/issues/191))
- [ ] **Standardized Recommendation Reasoning Metadata**
  ([#192](https://github.com/rshade/pulumicost-spec/issues/192))
- [ ] **Distributed Tracing Propagation (Contextual Visibility)**
  ([#193](https://github.com/rshade/pulumicost-spec/issues/193))
- [ ] **Authorization Middleware (OIDC/IAM)**
  ([#195](https://github.com/rshade/pulumicost-spec/issues/195)) -
  Standardizing how plugins receive and validate identity without violating "Stateless" boundaries.
- [ ] **Cross-Language SDKs (Python/TS)**
  ([#196](https://github.com/rshade/pulumicost-spec/issues/196)) -
  To be initiated after Go SDK reaches v1.0 stability and a critical mass of plugins exist.
- [ ] **Streaming Actual Cost (Streaming RPCs)**
  ([#197](https://github.com/rshade/pulumicost-spec/issues/197)) -
  Evaluating the need for `StreamActualCost` for real-time anomaly detection.

### SDK Developer Experience

- [ ] **Custom Health Checker Interface** ([#230](https://github.com/rshade/pulumicost-spec/issues/230)) -
  Allow plugins to implement custom health check logic beyond simple "ok" responses.
- [ ] **Context Validation Helpers** ([#232](https://github.com/rshade/pulumicost-spec/issues/232)) -
  Add utilities for validating context deadlines and detecting already-expired contexts.

### Proposed for Discussion (Discovery)

- [ ] **JSON-LD / Schema.org Serialization** ([#187](https://github.com/rshade/pulumicost-spec/issues/187))
- [ ] **Standardized Recommendation Reasoning** ([#188](https://github.com/rshade/pulumicost-spec/issues/188))
- [ ] **Multi-Currency Segregation Pattern** ([#190](https://github.com/rshade/pulumicost-spec/issues/190))
- [ ] **Validation Bypass Protocol** ([#216](https://github.com/rshade/pulumicost-spec/issues/216)) -
  Add `BypassReason` and `OverrideMetadata` to `ValidationResult` for governance auditing.

---

## Boundary Safeguards (The "Hard No's")

- **No Orchestration Logic**: Multi-plugin aggregation and "Muxing" logic will stay in `pulumicost-core`.
- **No Persistence**: The SDK will not provide built-in database support; it remains stateless.
- **No Native Math Engines**: The SDK will not perform financial amortization;
  it will only define the fields to store the results of such calculations.