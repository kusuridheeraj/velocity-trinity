# ADR-009: GitHub Webhook Integration (Quantum Merge)

## Context
Product C (`quantum-merge`) needs to know when PRs are opened or updated.
Polling GitHub API is slow and rate-limited. Webhooks are real-time.

## Decision
We will expose an HTTP endpoint `/webhook` in the `quantum-merge` binary.

## Rationale
1.  **Simplicity**: Go's `net/http` is production-ready.
2.  **Payload**: We will parse `pull_request` JSON events.
3.  **Security**: We will validate the `X-Hub-Signature-256` header using a shared secret.

## Architecture
-   User exposes `quantum-merge` via `ngrok` (dev) or Load Balancer (prod).
-   GitHub sends JSON payload.
-   Handler extracts `number`, `head.sha`, `base.ref`.
-   Handler calls `Scheduler.Enqueue()`.

## Consequences
-   **Positive**: Real-time integration.
-   **Negative**: Binary now needs to be a long-running server (daemon), not just a CLI.
-   **Mitigation**: We already designed `serve` command for this.
