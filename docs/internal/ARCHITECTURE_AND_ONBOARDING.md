# Architecture & Onboarding: The Senior Engineer's Handbook

## 1. The High-Level Cartography (Architecture & Trade-offs)

### System Diagram
We are building a **Distributed DevOps Toolchain** composed of three independent pillars that can be deployed together or separately.

*   **Pillar A (Dependency-CI):** A CLI tool that statically analyzes code to determine the minimal test set.
*   **Pillar B (LivePatch):** A Client-Server architecture for file synchronization over secure TLS tunnels.
*   **Pillar C (Quantum Merge):** An Event-Driven microservice (Webhook Listener -> Queue -> Worker) for predictive CI.

### Directory Structure & Responsibilities
*   `cmd/`: **The Entry Points.** Each subfolder (`dependency-ci`, `live-patch`, `quantum-merge`) builds into a separate binary.
*   `pkg/analyzer`: **The Brain.** Parses source code (AST/Regex) to understand dependencies.
*   `pkg/transport`: **The Nervous System.** Defines the RPC protocol and data structures for file transfer.
*   `pkg/scheduler`: **The Heart.** Manages the job queue and worker pool for Quantum Merge.
*   `pkg/dashboard`: **The Face.** Serves the embedded UI and API.

### Key Trade-offs
1.  **Monorepo vs. Polyrepo**:
    *   **Choice:** Monorepo.
    *   **Why:** Rapid iteration across shared libraries (`pkg/logger`, `pkg/config`).
    *   **Downside:** Build times increase as the codebase grows. Tooling complexity (checking which service changed).

2.  **Go `net/rpc` vs. gRPC**:
    *   **Choice:** Native Go RPC.
    *   **Why:** Zero dependencies. No `protoc` installation required for developers.
    *   **Downside:** Not interoperable with Python/Node.js clients easily.

3.  **In-Memory Queue vs. Redis**:
    *   **Choice:** In-Memory (Go Channels).
    *   **Why:** "Day 1" usability. Users can run the binary instantly.
    *   **Downside:** Data loss on crash. Not HA.

---

## 2. The "Chesterton's Fence" Analysis (Component Deep Dive)

### Component: `pkg/transport/server.go` (LivePatch Agent)
1.  **Purpose:** Receives file bytes and writes them to disk inside a container.
2.  **The "Bus Factor":** If deleted, **LivePatch stops working entirely.** The CLI cannot connect.
3.  **Dependencies:** Relies on `pkg/logger` for audit trails and `crypto/tls` for security.

### Component: `pkg/analyzer/test_finder.go` (Dependency-CI)
1.  **Purpose:** Maps `foo.ts` to `foo.test.ts`.
2.  **The "Bus Factor":** If buggy, **CI pipelines might skip relevant tests**, leading to bugs in production.
3.  **Dependencies:** None (Pure logic).

### Component: `pkg/scheduler/worker.go` (Quantum Merge)
1.  **Purpose:** Background workers that pick up CI jobs and execute them.
2.  **The "Bus Factor":** If removed, the API accepts jobs but **nothing ever runs.** The queue fills up infinitely.
3.  **Dependencies:** `pkg/scheduler/queue.go`.

---

## 3. The Operator's Manual (Run & Test)

### Prerequisites
*   **Go 1.23+**: Required for `go.mod`.
*   **Docker**: Optional (for testing Agent inside containers).
*   **Make** (Optional): For running complex build scripts.

### Environment Variables
*   `PORT`: Port for LivePatch Agent (Default: `8080`) or Quantum Merge (Default: `8090`).
*   `LOG_LEVEL`: `debug`, `info`, `warn`, `error` (Default: `info`).
*   `ENV`: `production` or `development` (Default: `development`).

### Step-by-Step Start
```bash
# 1. Install Dependencies
go mod tidy

# 2. Build All Tools
go build ./...

# 3. Run Unit Tests
go test ./...
```

### Critical Tests to Write (If Missing)
1.  **LivePatch Integration Test:** Spin up a real Docker container, run the Agent inside, and try to `sync` a file from the host.
2.  **Dependency-CI Parser Test:** Create a complex file with weird imports (e.g., `import { a } from "./b"; // comment`) and ensure regex parses it correctly.
3.  **Quantum Merge Concurrency Test:** Enqueue 100 jobs simultaneously and ensure the worker pool processes them without race conditions.

---

## 4. The "Black Box" Simulation (Real-World Debugging)

### Scenario 1: The Silent Failure
*   **Context:** A user says "I ran `live-patch sync`, it said 'Success', but the file didn't change in the container!"
*   **Clue:** Check `pkg/transport/server.go`. Is it writing to the correct path? Is it silently ignoring errors?
*   **Task:** Add debug logs to print the *absolute path* where the file is being written.

### Scenario 2: The Infinite Loop
*   **Context:** The Quantum Merge Dashboard shows a job as "RUNNING" for 4 hours.
*   **Clue:** Check `pkg/scheduler/worker.go`. Did the worker panic? Is there a timeout?
*   **Task:** Implement a context with timeout (e.g., 10 minutes) for every job execution.

### Scenario 3: The False Positive
*   **Context:** Dependency-CI is skipping `user.test.ts` even though `user.ts` changed.
*   **Clue:** Check `pkg/analyzer/test_finder.go`. Does it support the specific file extension or directory structure?
*   **Task:** Add a test case for `src/models/user.ts` -> `tests/unit/user.test.ts` mapping.

---

## 5. Critique & Roadmap

### The Messiest Part
*   **Dependency Parsing (Regex):** Using Regex to parse code is fragile. It will break on multi-line imports or complex syntax.
*   **Refactor:** Replace Regex with a proper AST parser (e.g., `tree-sitter` or `go/ast`) for robust language support.

### Immediate Improvement
*   **Configuration:** Currently, flags are somewhat hardcoded.
*   **Action:** Centralize all config into `pkg/config` using Viper to support `config.yaml`, env vars, and flags uniformly across all 3 tools.
