# Strategic Roadmap: pulumicost-spec

## Vision

To provide the definitive, high-performance gRPC specification and Go SDK for cloud cost observability, centered around the FinOps Foundation's FOCUS standard.

---

## [Past Milestones]

### Protocol & Modeling
- [x] **GreenOps Standardization** ([#176](https://github.com/rshade/pulumicost-spec/issues/176)) - Sustainability metrics (carbon footprint, energy utilization) integrated into the core protocol.
- [x] **Recommendation Enhancements** ([#173](https://github.com/rshade/pulumicost-spec/issues/173), [#171](https://github.com/rshade/pulumicost-spec/issues/171), [#166](https://github.com/rshade/pulumicost-spec/issues/166)) - Added resource-scoped targets, filtering, and extended action types.
- [x] **FOCUS 1.2 Integration** ([#100](https://github.com/rshade/pulumicost-spec/issues/100), [#99](https://github.com/rshade/pulumicost-spec/issues/99)) - Full schema coverage and builder API for FOCUS 1.2.
- [x] **RPC Expansion** ([#125](https://github.com/rshade/pulumicost-spec/issues/125), [#123](https://github.com/rshade/pulumicost-spec/issues/123), [#90](https://github.com/rshade/pulumicost-spec/issues/90)) - Implemented `EstimateCost`, `GetBudgets`, and `GetRecommendations`.
- [x] **Zero-Allocation Enums** ([#63](https://github.com/rshade/pulumicost-spec/issues/63), [#33](https://github.com/rshade/pulumicost-spec/issues/33)) - High-performance validation for domain enums.

### SDK & Tooling
- [x] **Plugin Conformance Suite** ([#109](https://github.com/rshade/pulumicost-spec/issues/109)) - Automated testing for plugin implementers to ensure spec compliance.
- [x] **SDK Foundation** ([#151](https://github.com/rshade/pulumicost-spec/issues/151), [#148](https://github.com/rshade/pulumicost-spec/issues/148), [#145](https://github.com/rshade/pulumicost-spec/issues/145), [#139](https://github.com/rshade/pulumicost-spec/issues/139)) - Centralized environment handling, Prometheus metrics, and unified file-based logging.
- [x] **Orchestration Support** ([#181](https://github.com/rshade/pulumicost-spec/issues/181), [#143](https://github.com/rshade/pulumicost-spec/issues/143), [#126](https://github.com/rshade/pulumicost-spec/issues/126)) - Added reflection, `--port` flags, and `fallback-hints`.

---

## [Immediate Focus] (Q1 2026)

### High Priority Research
- [ ] **Plugin Capability Discovery (Feature Flagging)** ([#194](https://github.com/rshade/pulumicost-spec/issues/194)) - Implementing a discovery protocol for advertisement of supported RPCs.

### Stability & Maintenance
- [ ] **Dependency Management** ([#13](https://github.com/rshade/pulumicost-spec/issues/13)) - Automated tracking and updating of core proto and SDK dependencies.

### Planned Features
- [ ] **FOCUS 1.3 Migration** ([#183](https://github.com/rshade/pulumicost-spec/issues/183)) - Audit new columns and entities in FOCUS 1.3 and update builder APIs.
- [ ] **Contextual FinOps Validation** ([#184](https://github.com/rshade/pulumicost-spec/issues/184)) - Extend `pluginsdk/validation` to include contextual checks (e.g., ensuring `ListPrice` is not less than `BilledCost`).
- [ ] **Advanced SDK Patterns** ([#185](https://github.com/rshade/pulumicost-spec/issues/185)) - Implementation examples for complex tiered pricing and multi-provider mapping.
- [ ] **SDK Documentation Overhaul** - Comprehensive Godoc and implementation examples for the new `mapping` and `validation` packages.

### In Progress
- [/] **Standardized Benchmark Suite** ([#113](https://github.com/rshade/pulumicost-spec/issues/113), [#142](https://github.com/rshade/pulumicost-spec/issues/142)) - Formalizing "Time to First Byte" and memory allocation benchmarks for plugins.

---

## [Future Vision] (Long-Term)

### Active Research
- [ ] **Standardized Cost Allocation Lineage Metadata** ([#191](https://github.com/rshade/pulumicost-spec/issues/191))
- [ ] **Standardized Recommendation Reasoning Metadata** ([#192](https://github.com/rshade/pulumicost-spec/issues/192))
- [ ] **Distributed Tracing Propagation (Contextual Visibility)** ([#193](https://github.com/rshade/pulumicost-spec/issues/193))
- [ ] **Authorization Middleware (OIDC/IAM)** ([#195](https://github.com/rshade/pulumicost-spec/issues/195)) - Standardizing how plugins receive and validate identity without violating "Stateless" boundaries.
- [ ] **Cross-Language SDKs (Python/TS)** ([#196](https://github.com/rshade/pulumicost-spec/issues/196)) - To be initiated after Go SDK reaches v1.0 stability and a critical mass of plugins exist.
- [ ] **Streaming Actual Cost (Streaming RPCs)** ([#197](https://github.com/rshade/pulumicost-spec/issues/197)) - Evaluating the need for `StreamActualCost` for real-time anomaly detection.

### Proposed for Discussion (Discovery)
- [ ] **Plugin Capability Dry Run Mode** ([#186](https://github.com/rshade/pulumicost-spec/issues/186))
- [ ] **JSON-LD / Schema.org Serialization** ([#187](https://github.com/rshade/pulumicost-spec/issues/187))
- [ ] **Standardized Recommendation Reasoning** ([#188](https://github.com/rshade/pulumicost-spec/issues/188))
- [ ] **gRPC-Web support (Pulumi Insights)** ([#189](https://github.com/rshade/pulumicost-spec/issues/189))
- [ ] **Multi-Currency Segregation Pattern** ([#190](https://github.com/rshade/pulumicost-spec/issues/190))

---

## Boundary Safeguards (The "Hard No's")

- **No Orchestration Logic**: Multi-plugin aggregation and "Muxing" logic will stay in `pulumicost-core`.
- **No Persistence**: The SDK will not provide built-in database support; it remains stateless.
- **No Native Math Engines**: The SDK will not perform financial amortization; it will only define the fields to store the results of such calculations.