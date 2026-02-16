# ADR-001: Monorepo Architecture

## Context
We are building a suite of three interconnected DevOps tools: `dependency-ci`, `live-patch`, and `quantum-merge`. These tools share significant common functionality (logging, configuration, authentication, Kubernetes clients). The project timeline is aggressive (21 days) and requires rapid iteration.

## Decision
We will adopt a **Monorepo** structure for all three projects.

## Rationale
1.  **Code Reuse**: Shared libraries in `pkg/` can be imported directly by all tools without complex versioning or publishing to external package registries.
2.  **Atomic Commits**: Features that span multiple tools (e.g., a shared logging update) can be committed in a single transaction.
3.  **Unified CI**: A single CI pipeline can build, test, and release all tools simultaneously, reducing maintenance overhead.
4.  **Simplicity**: Developer onboarding is faster with a single repository to clone and explore.

## Consequences
-   **Positive**: Faster development velocity, easier refactoring across tools.
-   **Negative**: Potential for build times to increase as the codebase grows (mitigated by Go's fast compilation).
-   **Mitigation**: We will use Go modules and strictly separate `cmd/` directories to keep binaries distinct.
