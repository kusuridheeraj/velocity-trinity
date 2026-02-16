# ADR-010: Dashboard Implementation (Product C - Phase 3)

## Context
We need a unified dashboard to visualize the state of all three tools:
-   **CI**: Which tests ran?
-   **LivePatch**: Which containers were patched?
-   **Quantum**: What is the queue status?

## Decision
We will embed a **Single Page Application (SPA)** written in plain HTML/JS directly into the `quantum-merge` binary using Go's `embed` package.

## Rationale
1.  **Single Binary**: Users still only need to run `./quantum-merge serve`. The dashboard is served from `http://localhost:8090/`.
2.  **Simplicity**: No need for a separate `npm build` step for the frontend in the MVP phase.
3.  **Real-time**: The dashboard will poll the JSON API we exposed.

## API Design
-   `GET /api/jobs`: Returns JSON list of jobs in the queue.
-   `GET /api/stats`: Returns simple counters.

## Consequences
-   **Positive**: Zero deployment friction.
-   **Negative**: UI development is less ergonomic (no hot reloading of React components).
-   **Mitigation**: We will keep the UI very simple (Bootstrap + Vanilla JS).
