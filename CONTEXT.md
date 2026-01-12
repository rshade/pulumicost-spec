# Project Context: finfocus-spec

## Core Architectural Identity

The `finfocus-spec` is a **Universal Specification and SDK Layer** for cloud cost observability.
It defines the standardized interfaces (Protobuf/gRPC) and provides the Go SDK (`pluginsdk`)
required to build and integrate cost-estimation plugins. It acts as the "narrow waist"
between cloud cost data sources and consumption layers.

## Technical Boundaries

- **No Financial Logic/Calculation:** This project does not calculate tiered pricing,
  amortization, or committed-use discounts (RI/SP). It defines the *schema* to hold these
  values, but the logic must reside in upstream providers.
- **No Stateful Data Storage:** The repository does not manage databases, perform
  migrations, or store historical cost data. It is a stateless interface and transport definition.
- **No Direct Cloud API Consumption:** The SDK and spec do not talk directly to AWS, Azure,
  or GCP billing APIs. That responsibility is strictly delegated to individual
  *plugin implementations* (e.g., `finfocus-plugins-aws`).
- **No Frontend/Visualization:** This is a backend/infrastructure-level specification.
  It does not contain UI components, dashboards, or graphing logic.
- **Observability is for Diagnostics, not End-Users:** Features like Prometheus metrics
  and structured logging within the SDK are intended for plugin maintainers to monitor
  performance and health, not as a primary cost-reporting feature for end-users.

## Data Source of Truth

The **Source of Truth** resides entirely within **Upstream Data Providers**
(e.g., Kubecost, Vantage, Cloud Provider Billing APIs).

- `finfocus-spec` is responsible for the **Standardized Model** (how data is structured).
- Plugin implementations are responsible for the **Data Retrieval** (fetching the data).
- The upstream service is responsible for the **Financial Accuracy** (the actual numbers).

## Interaction Model

- **Primary Protocol:** gRPC over Protobuf for high-performance, typed communication
  between the host and plugins.
- **Service Discovery:** JSON-based Plugin Manifests and Registries for plugin metadata
  and capability advertisement.
- **SDK Pattern:** A "Zero-Allocation" oriented Go SDK (`pluginsdk`) that provides
  boilerplate for gRPC server lifecycle, validation, and logging.
- **Configuration:** Environment-variable-driven configuration for SDK-level features
  (logging levels, ports).
