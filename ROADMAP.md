# Strategic Roadmap: finfocus-spec

## Vision

To provide the definitive, high-performance gRPC specification and Go SDK for cloud cost observability,
centered around the FinOps Foundation's FOCUS standard.

---

## [Past Milestones]

### Protocol & Modeling

- [x] **GreenOps Standardization** ([#176](https://github.com/rshade/finfocus-spec/issues/176)) -
  Sustainability metrics (carbon footprint, energy utilization) integrated into the core protocol.
- [x] **Recommendation Enhancements**
  ([#173](https://github.com/rshade/finfocus-spec/issues/173),
  [#171](https://github.com/rshade/finfocus-spec/issues/171),
  [#166](https://github.com/rshade/finfocus-spec/issues/166)) -
  Added resource-scoped targets, filtering, and extended action types.
- [x] **FOCUS 1.2 Integration**
  ([#100](https://github.com/rshade/finfocus-spec/issues/100),
  [#99](https://github.com/rshade/finfocus-spec/issues/99)) -
  Full schema coverage and builder API for FOCUS 1.2.
- [x] **RPC Expansion**
  ([#125](https://github.com/rshade/finfocus-spec/issues/125),
  [#123](https://github.com/rshade/finfocus-spec/issues/123),
  [#90](https://github.com/rshade/finfocus-spec/issues/90)) -
  Implemented `EstimateCost`, `GetBudgets`, and `GetRecommendations`.
- [x] **Zero-Allocation Enums**
  ([#63](https://github.com/rshade/finfocus-spec/issues/63),
  [#33](https://github.com/rshade/finfocus-spec/issues/33)) -
  High-performance validation for domain enums.

### SDK & Tooling

- [x] **Plugin Conformance Suite** ([#109](https://github.com/rshade/finfocus-spec/issues/109)) -
  Automated testing for plugin implementers to ensure spec compliance.
- [x] **SDK Foundation**
  ([#151](https://github.com/rshade/finfocus-spec/issues/151),
  [#148](https://github.com/rshade/finfocus-spec/issues/148),
  [#145](https://github.com/rshade/finfocus-spec/issues/145),
  [#139](https://github.com/rshade/finfocus-spec/issues/139)) -
  Centralized environment handling, Prometheus metrics, and unified file-based logging.
- [x] **Orchestration Support**
  ([#181](https://github.com/rshade/finfocus-spec/issues/181),
  [#143](https://github.com/rshade/finfocus-spec/issues/143),
  [#126](https://github.com/rshade/finfocus-spec/issues/126)) -
  Added reflection, `--port` flags, and `fallback-hints`.
- [x] **Multi-Protocol Support (gRPC-Web/Connect)**
  ([#189](https://github.com/rshade/finfocus-spec/issues/189),
  [#223](https://github.com/rshade/finfocus-spec/pull/223)) -
  Added connect-go integration enabling gRPC, gRPC-Web, and Connect protocols for browser compatibility.

---

## [Completed Milestones] (Q1 2026)

### Protocol & Modeling

- [x] **Forecasting Primitives**
  ([#241](https://github.com/rshade/finfocus-spec/issues/241),
  [#250](https://github.com/rshade/finfocus-spec/issues/250)) -
  Added `GrowthType` (Linear, Exponential) and `GrowthRate` to `CostResult` for cost projections.
- [x] **FOCUS 1.3 Migration** ([#199](https://github.com/rshade/finfocus-spec/issues/199)) -
  Audit new columns and entities in FOCUS 1.3 and update builder APIs.

### SDK & Tooling

- [x] **Plugin Capability Discovery** ([#242](https://github.com/rshade/finfocus-spec/issues/242)) -
  Implemented `GetPluginInfo` RPC for spec version compatibility and capability advertisement.
- [x] **Plugin Capability Dry Run Mode** ([#248](https://github.com/rshade/finfocus-spec/issues/248)) -
  Implemented `DryRun` for plugin field mapping discovery.
- [x] **Contextual FinOps Validation** ([#201](https://github.com/rshade/finfocus-spec/issues/201)) -
  Extended `pluginsdk/validation` to include contextual checks.
- [x] **Advanced SDK Patterns** ([#213](https://github.com/rshade/finfocus-spec/issues/213)) -
  Implementation examples for complex tiered pricing and multi-provider mapping.
- [x] **SDK Documentation Overhaul** ([#243](https://github.com/rshade/finfocus-spec/issues/243)) -
  Comprehensive Godoc, thread safety, rate limiting, and performance documentation.
- [x] **DismissRecommendation Consistency** ([#225](https://github.com/rshade/finfocus-spec/issues/225)) -
  Made DismissRecommendation follow the same pattern as other RPC methods.
- [x] **Standardized Benchmark Suite**
  ([#113](https://github.com/rshade/finfocus-spec/pull/113),
  [#142](https://github.com/rshade/finfocus-spec/pull/142)) -
  Formalized "Time to First Byte" and memory allocation benchmarks for plugins.
- [x] **JSON-LD / Schema.org Serialization** ([#187](https://github.com/rshade/finfocus-spec/issues/187),
  [#252](https://github.com/rshade/finfocus-spec/pull/252)) -
  Added `jsonld` package for FOCUS cost data serialization with Schema.org compatibility.

- [x] **v0.4.14 SDK Polish Release** ([#257](https://github.com/rshade/finfocus-spec/issues/257)) -
  Significant improvements to SDK developer experience, testing, and Connect protocol support.
  Includes:
  - **Connect Protocol**: Test coverage ([#227](https://github.com/rshade/finfocus-spec/issues/227)),
    CORS validation ([#234](https://github.com/rshade/finfocus-spec/issues/234)).
  - **Developer Experience**: Custom HealthChecker ([#230](https://github.com/rshade/finfocus-spec/issues/230)),
    Context helpers ([#232](https://github.com/rshade/finfocus-spec/issues/232)), ARN validation ([#203](https://github.com/rshade/finfocus-spec/issues/203)),
    Migration guide ([#246](https://github.com/rshade/finfocus-spec/issues/246)).
  - **Quality**: CI Benchmarks ([#224](https://github.com/rshade/finfocus-spec/issues/224)),
    Extreme value tests ([#212](https://github.com/rshade/finfocus-spec/issues/212)), Fuzzing ([#205](https://github.com/rshade/finfocus-spec/issues/205)).
- [x] **GetPluginInfo Performance Test** ([#244](https://github.com/rshade/finfocus-spec/issues/244)) -
  Added standalone conformance test for GetPluginInfo response time (<100ms assertion).
- [x] **User-Friendly Error Messages** ([#245](https://github.com/rshade/finfocus-spec/issues/245)) -
  Improved GetPluginInfo error messages for end-users.
- [x] **TypeScript SDK** ([#293](https://github.com/rshade/finfocus-spec/issues/293),
  [#302](https://github.com/rshade/finfocus-spec/pull/302)) -
  Initial TypeScript client SDK for browser and Node.js environments.
- [x] **Plugin Capability Enum** ([#287](https://github.com/rshade/finfocus-spec/issues/287)) -
  Added PluginCapability enum for granular feature discovery.
- [x] **Backward Compatibility for Environment Variables**
  ([#283](https://github.com/rshade/finfocus-spec/issues/283)) -
  Maintained legacy PULUMICOST_* environment variable support during migration.
- [x] **Migration Documentation** ([#282](https://github.com/rshade/finfocus-spec/issues/282)) -
  Added MIGRATION.md and llm-migration.json for PulumiCost to FinFocus migration.
- [x] **PulumiCost to FinFocus Rename** ([#272](https://github.com/rshade/finfocus-spec/issues/272)) -
  Complete project rename with backward compatibility shims.

---

## [Immediate Focus] (Q1 2026)

### TypeScript SDK

- [ ] **Migrate TypeScript SDK to Connect-ES v2**
  ([#304](https://github.com/rshade/finfocus-spec/issues/304)) -
  Update TypeScript SDK to use Connect-ES v2 for improved browser/Node.js support.

### Stability & Maintenance

- [ ] **Dependency Management** ([#13](https://github.com/rshade/finfocus-spec/issues/13)) -
  Automated tracking and updating of core proto and SDK dependencies.

### Bug Fixes (Capability Discovery)

- [ ] **Missing Nil Check in inferCapabilities**
  ([#297](https://github.com/rshade/finfocus-spec/issues/297)) -
  Fix potential nil pointer dereference in capability inference.
- [ ] **Missing Test Coverage for Capability Override Edge Cases**
  ([#296](https://github.com/rshade/finfocus-spec/issues/296)) -
  Add tests for edge cases in capability override handling.
- [ ] **Proto Field Type Inconsistency**
  ([#295](https://github.com/rshade/finfocus-spec/issues/295)) -
  Fix inconsistency between capabilities and metadata field types.
- [ ] **Missing DryRunHandler Interface Definition**
  ([#294](https://github.com/rshade/finfocus-spec/issues/294)) -
  Add proper interface definition for DryRunHandler.

### SDK Quality

- [ ] **Inconsistent Backward Compatibility Patterns**
  ([#299](https://github.com/rshade/finfocus-spec/issues/299)) -
  Standardize backward compatibility conversion patterns across SDK.
- [ ] **Optimize Slice Copying Performance**
  ([#301](https://github.com/rshade/finfocus-spec/issues/301)) -
  Use append pattern for improved slice copying performance.

### Documentation

- [ ] **Documentation Updates for Capability Discovery**
  ([#300](https://github.com/rshade/finfocus-spec/issues/300)) -
  Update docs to reflect new capability discovery system.
- [ ] **Unclear Exhaustive Nolint Directive**
  ([#298](https://github.com/rshade/finfocus-spec/issues/298)) -
  Clarify the nolint directive in legacyCapabilityMap.
- [ ] **Strengthen Log File Security Warning**
  ([#284](https://github.com/rshade/finfocus-spec/issues/284)) -
  Improve security guidance for log file handling (good first issue).

### Planned Features

- [ ] **Pricing Tier Intelligence** ([#217](https://github.com/rshade/finfocus-spec/issues/217))
  - Define `PricingTier` enum (`STANDARD`, `SPOT`, `RESERVED`)
  - Define `SpotRisk` enum (`LOW`, `MEDIUM`, `HIGH`) for interruption probability
- [ ] **Usage Profiles** ([#218](https://github.com/rshade/finfocus-spec/issues/218))
  - Add `UsageProfile` context to requests to allow plugins to adjust recommendations
    based on environment (e.g., Dev vs Prod)

---

## [Future Vision] (Long-Term)

### Active Research

- [ ] **Standardized Cost Allocation Lineage Metadata**
  ([#191](https://github.com/rshade/finfocus-spec/issues/191))
- [ ] **Standardized Recommendation Reasoning Metadata**
  ([#192](https://github.com/rshade/finfocus-spec/issues/192))
- [ ] **Distributed Tracing Propagation (Contextual Visibility)**
  ([#193](https://github.com/rshade/finfocus-spec/issues/193))
- [ ] **Authorization Middleware (OIDC/IAM)**
  ([#195](https://github.com/rshade/finfocus-spec/issues/195)) -
  Standardizing how plugins receive and validate identity without violating "Stateless" boundaries.
- [ ] **Cross-Language SDKs (Python/TS)**
  ([#196](https://github.com/rshade/finfocus-spec/issues/196)) -
  To be initiated after Go SDK reaches v1.0 stability and a critical mass of plugins exist.
- [ ] **Streaming Actual Cost (Streaming RPCs)**
  ([#197](https://github.com/rshade/finfocus-spec/issues/197)) -
  Evaluating the need for `StreamActualCost` for real-time anomaly detection.

### SDK Developer Experience

### Proposed for Discussion (Discovery)

- [ ] **Standardized Recommendation Reasoning** ([#188](https://github.com/rshade/finfocus-spec/issues/188))
- [ ] **Multi-Currency Segregation Pattern** ([#190](https://github.com/rshade/finfocus-spec/issues/190))
- [ ] **Validation Bypass Protocol** ([#216](https://github.com/rshade/finfocus-spec/issues/216)) -
  Add `BypassReason` and `OverrideMetadata` to `ValidationResult` for governance auditing.
- [ ] **Per-Request Credential Passing** ([#220](https://github.com/rshade/finfocus-spec/issues/220)) -
  Multi-tenant optimization allowing per-request cloud credentials.
- [ ] **Batch RPC for Multi-Resource Queries** ([#221](https://github.com/rshade/finfocus-spec/issues/221)) -
  Efficient batch API for querying costs across multiple resources.

---

## Boundary Safeguards (The "Hard No's")

- **No Orchestration Logic**: Multi-plugin aggregation and "Muxing" logic will stay in `finfocus-core`.
- **No Persistence**: The SDK will not provide built-in database support; it remains stateless.
- **No Native Math Engines**: The SDK will not perform financial amortization;
  it will only define the fields to store the results of such calculations.
