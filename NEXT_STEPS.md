# Next Steps: Completing the 21-Day Plan

Congratulations! You have successfully built the **Core Foundation** (Days 1-4) of the Velocity Trinity suite.
All three CLI tools compile and run:

- `dependency-ci.exe`: Parses TypeScript/Python dependencies.
- `live-patch.exe` + `live-patch-agent.exe`: Syncs files securely over TLS.
- `quantum-merge.exe`: Runs a speculative job queue in-memory.

---

## Phase 2: Core Mechanics (Days 6-12)

### Day 6-7: Product A - The "Smart Skip" Engine
**Task:** Connect `dependency-ci` to your test runner.
1.  Modify `pkg/analyzer/analyzer.go` to return a list of *test files* affected by a change, not just source files.
2.  Implement `dependency-ci run -- cmd="npm test"`:
    - Parse changed files from Git (`git diff --name-only`).
    - Calculate affected tests.
    - Run `npm test <list of tests>`.

### Day 8-9: Product B - Hot Reload Injection
**Task:** Make the changes *live*.
1.  Update `pkg/transport/server.go` to trigger a command after sync.
2.  Add a `PostSyncCommand` field to `FileSyncRequest`.
3.  Example: `live-patch sync server.js --restart="systemctl restart node-app"`.

### Day 10-11: Product C - GitHub Integration
**Task:** Connect `quantum-merge` to real PRs.
1.  Create a GitHub App.
2.  Implement a webhook handler in `cmd/quantum-merge/main.go` for `pull_request` events.
3.  When a PR is opened, `Enqueue()` a job.

### Day 12: Security Hardening
1.  Replace self-signed certs with a proper CA management command (`live-patch init-ca`).
2.  Add API Keys for `quantum-merge` so only GitHub can trigger builds.

---

## Phase 3: The "Enterprise" Layer (Days 13-18)

### Day 13-14: The Dashboard
**Task:** Visualize the pipeline.
1.  Create a new `web/dashboard` folder (Next.js).
2.  Expose a REST API in `cmd/quantum-merge` to list jobs: `GET /jobs`.

### Day 15-18: Reliability
1.  **Product A:** Add caching (don't re-analyze if files haven't changed).
2.  **Product B:** Add checksum verification to ensure file integrity.
3.  **Product C:** Switch `MemoryQueue` to `RedisQueue` (implement the interface).

---

## Phase 4: Launch (Days 19-21)
1.  Write `CONTRIBUTING.md`.
2.  Create release scripts in `.github/workflows`.
3.  Record a demo video of all 3 tools working together!
