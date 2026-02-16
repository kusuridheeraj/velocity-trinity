# ADR-011: Reliability & Security Hardening (Advanced Features)

## Context
The MVP tools work, but they are fragile in enterprise environments.
-   **Security**: CA management for mTLS is manual.
-   **Reliability**: `dependency-ci` re-parses every file every time.

## Decision
We will implement:
1.  **Checksum Caching** for `dependency-ci` (avoid re-parsing unchanged files).
2.  **CA Helper Command** for `live-patch` (simplify mTLS setup).

## Rationale
1.  **Caching**: Parsing 100k files takes time. We will store a simple `.dep-ci-cache` JSON file mapping `file_path -> {hash, dependencies}`.
2.  **CA Helper**: Users struggle with `openssl`. A simple `live-patch init-ca` command reduces onboarding friction.

## Consequences
-   **Positive**: Faster second runs. Better UX.
-   **Negative**: Cache invalidation bugs are hard.
-   **Mitigation**: Add a `--clean` flag to force re-parsing.
