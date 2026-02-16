# Development Guide: Velocity Trinity

## Current State: Core MVP (Phase 1 Complete)
The Velocity Trinity suite is currently in the **MVP (Minimum Viable Product)** phase.
The following core binaries are functional:

- `dependency-ci.exe`: Parses TypeScript/Python dependencies and detects relevant tests.
- `live-patch.exe` + `live-patch-agent.exe`: Secure file sync with TLS and remote command execution.
- `quantum-merge.exe`: Runs a speculative job queue in-memory, listens to GitHub webhooks, and serves a dashboard API.

---

## Phase 2: Core Enhancements (Active Development)

### Milestone 2.1: Product A - The "Smart Skip" Engine
**Task:** Connect `dependency-ci` to your test runner.
1.  Modify `pkg/analyzer/analyzer.go` to return a list of *test files* affected by a change, not just source files.
2.  Implement `dependency-ci run -- cmd="npm test"`:
    - Parse changed files from Git (`git diff --name-only`).
    - Calculate affected tests.
    - Run `npm test <list of tests>`.

### Milestone 2.2: Product B - Hot Reload Injection
**Task:** Make the changes *live*.
1.  Update `pkg/transport/server.go` to trigger a command after sync.
2.  Add a `PostSyncCommand` field to `FileSyncRequest`.
3.  Example: `live-patch sync server.js --restart="systemctl restart node-app"`.

### Milestone 2.3: Product C - GitHub Integration
**Task:** Connect `quantum-merge` to real PRs.
1.  Create a GitHub App.
2.  Implement a webhook handler in `cmd/quantum-merge/main.go` for `pull_request` events.
3.  When a PR is opened, `Enqueue()` a job.

### Milestone 2.4: Security Hardening
1.  Replace self-signed certs with a proper CA management command (`live-patch init-ca`).
2.  Add API Keys for `quantum-merge` so only GitHub can trigger builds.

---

## Phase 3: The "Enterprise" Layer (Future)

### Milestone 3.1: The Dashboard
**Task:** Visualize the pipeline.
1.  Create a new `web/dashboard` folder (Next.js).
2.  Expose a REST API in `cmd/quantum-merge` to list jobs: `GET /jobs`.

### Milestone 3.2: Reliability & Scale
1.  **Product A:** Add caching (don't re-analyze if files haven't changed).
2.  **Product B:** Add checksum verification to ensure file integrity.
3.  **Product C:** Switch `MemoryQueue` to `RedisQueue` (implement the interface).

---

## Phase 4: Launch & Release
1.  Write `CONTRIBUTING.md`.
2.  Create release scripts in `.github/workflows`.
3.  Prepare demo assets.
