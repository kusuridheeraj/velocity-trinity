# Velocity Trinity: Development Roadmap

## Core Philosophy
**"Speed & Quality via Polyglot Architecture"**
- **Language:** Go (Golang) for all CLIs/Agents. fast, single-binary, cross-platform.
- **Delivery:** Docker containers for everything.
- **Target:** Kubernetes & Standard Docker Hosts.

---

## Phase 1: The Foundation (Core Architecture)

### Milestone 1.1: The Monorepo & Shared Core
**Goal:** A unified Go workspace where code can be shared across all 3 tools.
- **Task:** Initialize Go module `github.com/velocity-trinity/core`.
- **Task:** Implement structured logging (Zap/Logrus) that automatically scrubs secrets (PII/Tokens).
- **Task:** Build a `config` package that reads YAML/Env vars with validation.
- **Outcome:** A `hello-world` CLI that logs safely and reads config.

### Milestone 1.2: Product A - "Dependency-Aware CI" (The Brain)
**Goal:** Parse code dependencies without running it.
- **Task:** Integrate `tree-sitter` (Go bindings) for parsing TypeScript/Python.
- **Task:** Implement `ImportGraph` struct: Maps `file_A -> [file_B, file_C]`.
- **Task:** Build CLI command `dep-ci analyze <file>` that outputs affected files.
- **Outcome:** CLI that correctly identifies that changing `utils.ts` affects `main.ts`.

### Milestone 1.3: Product B - "LivePatch" (The Transport)
**Goal:** Secure, fast file syncing to a running container.
- **Task:** Build `livepatch-agent` (runs inside container) and `livepatch-cli` (runs on host).
- **Task:** Establish mTLS gRPC connection between CLI and Agent.
- **Task:** Implement binary-diff algorithm (like rsync) to send only changed bytes.
- **Outcome:** `livepatch-cli sync ./src` updates files inside a running Docker container in <500ms.

### Milestone 1.4: Product C - "Quantum Merge" (The Scheduler)
**Goal:** A queue that manages speculative jobs.
- **Task:** Setup Redis for job queue.
- **Task:** Implement `Scheduler` logic: "If PR #1 passes, promote PR #2's speculative build."
- **Task:** Build a mock "Executor" that simulates CI jobs (sleep 5s, pass/fail).
- **Outcome:** A simulation where 3 PRs are "merged" in parallel (virtually).

### Milestone 1.5: Integration & CI/CD
**Goal:** Automate *our* build process.
- **Task:** GitHub Actions to build binaries for Linux/Windows/Mac.
- **Task:** Docker multi-stage builds for the Agents.
- **Outcome:** Automatic release of `dep-ci.exe`, `livepatch.exe`, `quantum.exe`.

---

## Phase 2: The "Killer Features" (Advanced Functionality)

### Milestone 2.1: Product A - The "Smart Skip" Engine
- **Task:** Map tests to source files. (e.g., `user.test.ts` imports `user.ts`).
- **Task:** Implement `dep-ci run -- cmd="npm test"` which wraps the test runner and passes only relevant files.

### Milestone 2.2: Product B - Hot Reload Injection
- **Task:** Detect file changes (fsnotify).
- **Task:** Trigger process reload (SIGHUP or HMR signal) inside the container after sync.

### Milestone 2.3: Product C - GitHub Integration
- **Task:** Webhook handler for `pull_request` events.
- **Task:** Post status checks to GitHub ("Speculative Build: Pending").

### Milestone 2.4: Security Hardening
- **Task:** Audit all gRPC endpoints.
- **Task:** Ensure zero-trust auth between CLI and Agents.

---

## Phase 3: The "Enterprise" Layer (Management & Scale)

### Milestone 3.1: The Dashboard (Next.js + Go API)
- **Task:** Visual graph of dependencies (Product A).
- **Task:** Live logs from patch agents (Product B).
- **Task:** Queue visualization (Product C).

### Milestone 3.2: Advanced Reliability Features
- **Prod A:** Flaky test detection (retries).
- **Prod B:** Rollback on crash.
- **Prod C:** Conflict detection (Git merge simulation).

---

## Phase 4: Launch & Release

- **Documentation:** Write `getting-started`, `compliance-guide`, and API docs.
- **Final Polish:** Bug fixes and performance tuning.
- **Release:** Tag v1.0.0. Deploy demo environment.

---

## Monetization Strategy (Recap)

1.  **Dependency-Aware CI:**
    - **Free:** Open Source CLI.
    - **Paid:** "Enterprise Report" (Audit logs, Flake Analysis) - $49/mo.

2.  **LivePatch:**
    - **Free:** Local Docker syncing.
    - **Paid:** Kubernetes Cluster syncing (Team License) - $99/seat/mo.

3.  **Quantum Merge:**
    - **Free:** 1 concurrent speculation.
    - **Paid:** Unlimited parallelism - Usage based ($0.01/build minute).
