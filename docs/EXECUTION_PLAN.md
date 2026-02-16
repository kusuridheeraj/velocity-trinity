# 21-Day Execution Plan: Velocity Trinity

## Core Philosophy
**"Speed & Quality via Polyglot Architecture"**
- **Language:** Go (Golang) for all CLIs/Agents. fast, single-binary, cross-platform.
- **Delivery:** Docker containers for everything.
- **Target:** Kubernetes & Standard Docker Hosts.

---

## Phase 1: The Foundation (Days 1-5)

### Day 1: The Monorepo & Shared Core
**Goal:** A unified Go workspace where code can be shared across all 3 tools.
- **Task 1.1:** Initialize Go module `github.com/velocity-trinity/core`.
- **Task 1.2:** Implement structured logging (Zap/Logrus) that automatically scrubs secrets (PII/Tokens).
- **Task 1.3:** Build a `config` package that reads YAML/Env vars with validation.
- **Outcome:** A `hello-world` CLI that logs safely and reads config.

### Day 2: Product A - "Dependency-Aware CI" (The Brain)
**Goal:** Parse code dependencies without running it.
- **Task 2.1:** Integrate `tree-sitter` (Go bindings) for parsing TypeScript/Python.
- **Task 2.2:** Implement `ImportGraph` struct: Maps `file_A -> [file_B, file_C]`.
- **Task 2.3:** Build CLI command `dep-ci analyze <file>` that outputs affected files.
- **Outcome:** CLI that correctly identifies that changing `utils.ts` affects `main.ts`.

### Day 3: Product B - "LivePatch" (The Transport)
**Goal:** Secure, fast file syncing to a running container.
- **Task 3.1:** Build `livepatch-agent` (runs inside container) and `livepatch-cli` (runs on host).
- **Task 3.2:** Establish mTLS gRPC connection between CLI and Agent.
- **Task 3.3:** Implement binary-diff algorithm (like rsync) to send only changed bytes.
- **Outcome:** `livepatch-cli sync ./src` updates files inside a running Docker container in <500ms.

### Day 4: Product C - "Quantum Merge" (The Scheduler)
**Goal:** A queue that manages speculative jobs.
- **Task 4.1:** Setup Redis for job queue.
- **Task 4.2:** Implement `Scheduler` logic: "If PR #1 passes, promote PR #2's speculative build."
- **Task 4.3:** Build a mock "Executor" that simulates CI jobs (sleep 5s, pass/fail).
- **Outcome:** A simulation where 3 PRs are "merged" in parallel (virtually).

### Day 5: Integration & CI/CD
**Goal:** Automate *our* build process.
- **Task 5.1:** GitHub Actions to build binaries for Linux/Windows/Mac.
- **Task 5.2:** Docker multi-stage builds for the Agents.
- **Outcome:** Automatic release of `dep-ci.exe`, `livepatch.exe`, `quantum.exe`.

---

## Phase 2: The "Killer Features" (Days 6-12)

### Day 6-7: Product A - The "Smart Skip" Engine
- **Task:** Map tests to source files. (e.g., `user.test.ts` imports `user.ts`).
- **Task:** Implement `dep-ci run -- cmd="npm test"` which wraps the test runner and passes only relevant files.

### Day 8-9: Product B - Hot Reload Injection
- **Task:** Detect file changes (fsnotify).
- **Task:** Trigger process reload (SIGHUP or HMR signal) inside the container after sync.

### Day 10-11: Product C - GitHub Integration
- **Task:** Webhook handler for `pull_request` events.
- **Task:** Post status checks to GitHub ("Speculative Build: Pending").

### Day 12: Security Hardening
- **Task:** Audit all gRPC endpoints.
- **Task:** Ensure zero-trust auth between CLI and Agents.

---

## Phase 3: The "Enterprise" Layer (Days 13-18)

### Day 13-14: The Dashboard (Next.js + Go API)
- **Task:** Visual graph of dependencies (Product A).
- **Task:** Live logs from patch agents (Product B).
- **Task:** Queue visualization (Product C).

### Day 15-18: Advanced Features
- **Prod A:** Flaky test detection (retries).
- **Prod B:** Rollback on crash.
- **Prod C:** Conflict detection (Git merge simulation).

---

## Phase 4: Launch (Days 19-21)

- **Day 19:** Documentation (Docs site, Examples).
- **Day 20:** Final Polish & bug fixes.
- **Day 21:** "Launch" (Public Repo / Demo Video).

---

## Monetization Plan (Recap)

1.  **Dependency-Aware CI:**
    - **Free:** Open Source CLI.
    - **Paid:** "Enterprise Report" (Audit logs, Flake Analysis) - $49/mo.

2.  **LivePatch:**
    - **Free:** Local Docker syncing.
    - **Paid:** Kubernetes Cluster syncing (Team License) - $99/seat/mo.

3.  **Quantum Merge:**
    - **Free:** 1 concurrent speculation.
    - **Paid:** Unlimited parallelism - Usage based ($0.01/build minute).
