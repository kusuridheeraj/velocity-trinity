# ADR-005: Quantum Merge Queue Strategy

## Context
Product C (`quantum-merge`) requires a job queue to manage speculative CI runs. The original plan called for Redis.
However, we want the initial MVP to be easily runnable by users without setting up external infrastructure like Redis or PostgreSQL.

## Decision
We will define a `Queue` interface and provide two implementations:
1.  **In-Memory (Channel-based)**: The default for the MVP CLI.
2.  **Redis (Future)**: The scalable backend for production deployments.

## Rationale
1.  **Zero Friction**: Users can run `./quantum-merge serve` and see it work immediately.
2.  **Testability**: In-memory queues are trivial to test unit tests against.
3.  **Go Channels**: Go's native concurrency primitives are robust enough for single-node prototypes.

## Consequences
-   **Positive**: Instant "Time to Hello World".
-   **Negative**: If the process restarts, the queue is lost (acceptable for MVP).
-   **Mitigation**: We will add a `--redis-url` flag later to switch backends.
