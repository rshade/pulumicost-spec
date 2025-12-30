# Data Model: gRPC-Web Support

**Branch**: `029-grpc-web-support`

## Entities

No new data model entities are introduced. This feature extends the *transport layer* of the existing `CostSourceService`.

## API Contracts

The API contract remains defined by `proto/pulumicost/v1/costsource.proto`.

### Protocol Support

- **gRPC (Native)**: TCP/HTTP2 (Existing)
- **gRPC-Web**: HTTP/1.1 or HTTP/2 (New)
  - Content-Type: `application/grpc-web` or `application/grpc-web+proto`
  - Transports: `XMLHttpRequest` or `fetch` (Browser)

### Endpoints

- **Standard RPCs**: All existing RPCs (`GetActualCost`, `EstimateCost`, etc.) are exposed via gRPC-Web.
- **Health**:
  - `grpc.health.v1.Health/Check` (gRPC-Web)
  - `/healthz` (HTTP GET -> 200 OK)
