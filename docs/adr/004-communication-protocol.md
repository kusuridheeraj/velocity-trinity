# ADR-004: LivePatch Communication Protocol

## Context
Product B (`live-patch`) requires a high-performance, bidirectional channel between the developer's CLI and the agent running inside a container/VM.
Key Requirements:
1.  **Low Latency**: File sync must happen in < 500ms.
2.  **Security**: Data flows over potentially public networks (e.g., from dev laptop to cloud Kubernetes cluster).
3.  **Simplicity**: Minimize external dependencies (like `protoc`) to keep the build process simple and portable.

## Options
1.  **gRPC (Protobuf)**: Industry standard, strictly typed, but requires `protoc` compiler installation and generated code management.
2.  **REST (HTTP/JSON)**: Simple, universal, but high overhead for binary file transfer (Base64 encoding bloat).
3.  **Go `net/rpc` (Gob)**: Built-in to Go standard library, efficient binary encoding, zero external dependencies.

## Decision
We will use **Go `net/rpc` over TLS** for the MVP (Days 1-21).

## Rationale
1.  **Zero Dependency**: No need for users or CI to install the Protocol Buffers compiler.
2.  **Performance**: `encoding/gob` is highly optimized for Go-to-Go communication (perfect since both CLI and Agent are written in Go).
3.  **Speed of Implementation**: We can implement the full transport layer in < 100 lines of code using standard libraries.

## Consequences
-   **Positive**: Extremely fast development cycle, single binary distribution remains intact.
-   **Negative**: Communication is limited to Go clients (can't easily write a Python client later).
-   **Mitigation**: If we need non-Go clients in the future, we can wrap the RPC server with a gRPC gateway or HTTP/JSON adapter.
