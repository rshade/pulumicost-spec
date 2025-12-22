# Strategic Roadmap: pulumicost-spec

## Vision

To provide the definitive, high-performance gRPC specification and Go SDK for cloud cost observability,
centered around the FinOps Foundation's FOCUS standard.

---

## [Past Milestones]

- **[Done]** **FOCUS 1.2 Integration**: Full schema coverage and builder API for FOCUS 1.2.
- **[Done]** **RPC Expansion**: Implemented `EstimateCost`, `GetBudgets`, and `GetRecommendations`.
- **[Done]** **Zero-Allocation Enums**: High-performance validation for domain enums.
- **[Done]** **Plugin Conformance Suite**: Automated testing for plugin implementers to ensure spec compliance.
- **[Done]** **Orchestration Support**: Added `--port` flags and `fallback-hints` to enable `pulumicost-core`
  to manage multiple plugin lifecycles.

---

## [Immediate Focus] (Q1 2026)

- **[Planned]** **FOCUS 1.3 Migration**:
  - Audit new columns and entities in FOCUS 1.3 (e.g., improved commitment modeling).
  - Update Protobuf definitions and Go builder APIs to support 1.3 specs.
- **[Planned]** **Advanced Validation Helpers**:
  - Extend `pluginsdk/validation` to include contextual checks (e.g., ensuring `ListPrice` is not less than
    `BilledCost` in a FOCUS record).
- **[In Progress]** **GreenOps Standardization**:
  - Refining the metrics schema for carbon footprint and energy utilization based on initial GreenOps integration specs.
- **[Planned]** **SDK Documentation Overhaul**:
  - Comprehensive Godoc and implementation examples for the new `mapping` and `validation` packages.

---

## [Future Vision] (Long-Term)

- **[Researching]** **Authorization Middleware**:
  - Standardizing how plugins receive and validate identity (OIDC/IAM) via gRPC metadata without violating "Stateless" boundaries.
- **[Planned]** **Performance Benchmarking Suite**:
  - Standardized benchmarks for plugins to measure "Time to First Byte" for large-scale cost estimations.
- **[Researching]** **Cross-Language SDKs (Python/TS)**:
  - *Constraint*: Only to be initiated after the Go SDK reaches v1.0 stability and a critical mass of plugins (10+) exist.
  - Focus on maintaining 1:1 parity with the gRPC spec via automated generation.
- **[Researching]** **Streaming RPCs**:
  - Evaluating the need for `StreamActualCost` for real-time anomaly detection where request/response latency is too high.

---

## Boundary Safeguards (The "Hard No's")

- **No Orchestration Logic**: Multi-plugin aggregation and "Muxing" logic will stay in `pulumicost-core`.
- **No Persistence**: The SDK will not provide built-in database support.
- **No Native Math Engines**: The SDK will not perform financial amortization; it will only define the fields
  to store the results of such calculations.
