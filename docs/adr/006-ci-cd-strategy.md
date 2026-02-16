# ADR-006: CI/CD Pipeline Strategy (GitHub Actions)

## Context
We need to automate the build, test, and release process for all three products (`dependency-ci`, `live-patch`, `quantum-merge`).
Since we are using a **Monorepo**, a change in one tool shouldn't necessarily trigger a release for another unless they share code. However, for simplicity in the MVP phase, we will rebuild all artifacts on every push to `main`.

## Decision
We will use **GitHub Actions** with a matrix build strategy.

## Rationale
1.  **Matrix Builds**: Allows us to compile binaries for Linux (`amd64`, `arm64`), Windows (`amd64`), and macOS (`arm64`) in parallel.
2.  **Go Cache**: GitHub Actions has excellent support for caching `go mod` dependencies, speeding up builds.
3.  **Zero Infra**: No need to maintain a Jenkins server.

## Consequences
-   **Positive**: Every commit to `main` produces downloadable artifacts.
-   **Negative**: High usage of GitHub Actions minutes if commits are frequent.
-   **Mitigation**: We will only trigger full releases on tag pushes (e.g., `v*`).
