# Velocity Trinity: Comprehensive Codebase Guide

## 1. PROJECT OVERVIEW & ARCHITECTURE

### High-Level Architecture
*   **Type:** Distributed DevOps Suite (CLI + Agent + Server).
*   **Purpose:** Eliminates the "Integration Bottleneck" by optimizing testing (Dependency-CI), deployment (LivePatch), and merging (Quantum Merge).
*   **Tech Stack:** 
    *   **Go (Golang) 1.23+**: Chosen for performance, single-binary distribution, and concurrency.
    *   **Standard Library (`net/rpc`, `net/http`)**: Minimized external dependencies to keep builds fast and portable.
    *   **Gorilla Mux**: Routing for the Quantum Merge API.
    *   **Cobra/Viper**: CLI framework and configuration management.

### Architecture Diagram
```ascii
[Developer Machine]                  [Kubernetes Cluster / Remote Server]
       |                                          |
       |-- (1) dependency-ci ------------------> [Local Test Runner]
       |      (Analyzes Code -> Runs Tests)
       |
       |-- (2) live-patch CLI --(mTLS gRPC)--> [LivePatch Agent]
       |      (Syncs File)                        (Updates File & Restarts App)
       |
       |-- (3) quantum-merge CLI 
                  |
                  v
          [Quantum Merge Server] <--(Webhook)-- [GitHub]
                  |
                  |--(Dashboard API)
                  |--(Job Queue)
```

### Technology Trade-offs
1.  **Go vs. Node.js/Python**:
    *   **Why Go?** Static typing, compiled binaries (no `npm install` on customer machines), and goroutines for the merge queue.
    *   **Trade-off:** Slower iteration speed than Python for simple scripts. Harder to parse dynamic languages (AST complexity).
    *   **Alternative:** Node.js would allow sharing code with the dashboard frontend but would require the user to install a specific Node version.

2.  **`net/rpc` vs. gRPC (Protobuf)**:
    *   **Why `net/rpc`?** Zero dependencies. No `protoc` compiler needed. It's built into Go.
    *   **Trade-off:** Only works Go-to-Go. We cannot easily write a Python client for LivePatch later.
    *   **Alternative:** gRPC is industry standard but adds build complexity.

3.  **In-Memory Queue vs. Redis**:
    *   **Why In-Memory?** "Day 1" usability. Users can run the binary without setting up a Redis server.
    *   **Trade-off:** If the server crashes, all queued merge jobs are lost.
    *   **Alternative:** Redis is essential for production HA (High Availability).

---

## 2. PROJECT STRUCTURE DEEP DIVE

### Directory Structure
```
/velocity-trinity
├── /cmd                    # Entry points for binaries
│   ├── /dependency-ci      # Product A: Test optimization CLI
│   ├── /live-patch         # Product B: CLI & Agent
│   │   ├── /agent          # The server running inside containers
│   │   └── /cli            # The client tool
│   └── /quantum-merge      # Product C: Merge Queue Server
├── /pkg                    # Shared libraries
│   ├── /analyzer           # Dependency parsing logic (Product A)
│   ├── /config             # Viper configuration loader
│   ├── /dashboard          # Embedded Web Dashboard (Product C)
│   ├── /logger             # Zap structured logging
│   ├── /scheduler          # Queue logic (Product C)
│   ├── /transport          # RPC definitions (Product B)
│   ├── /utils              # Helper functions (TLS generation)
│   └── /webhook            # GitHub Webhook handler (Product C)
├── /docs                   # Documentation (Architecture, Usage)
└── go.mod                  # Dependencies
```

### Key Modules
*   **`pkg/analyzer`**:
    *   **Purpose:** Logic to read a file (TS/Python) and find its imports.
    *   **Entry Point:** `AnalyzeFile(path)`
    *   **Modification:** Change this if you need to support a new language (e.g., Rust/Java).
*   **`pkg/transport`**:
    *   **Purpose:** Defines the `FileSyncRequest` struct used by LivePatch.
    *   **Dependencies:** None (Pure Go structs).
    *   **Modification:** Change this to add new fields like `FileOwner` or `Permissions`.

---

## 3. CORE FUNCTIONALITY BREAKDOWN

### Feature: Smart Test Execution (Dependency-CI)
*   **What it does:** Scans changed source files (`foo.ts`) and finds their related tests (`foo.test.ts`).
*   **Why it exists:** Running a full 20-minute test suite for a 1-line change is wasteful.
*   **How it works:** 
    1.  User runs `dep-ci run --files="file1.ts file2.ts"`.
    2.  `pkg/analyzer` uses heuristics (naming convention) to find `file1.test.ts`.
    3.  It constructs a shell command: `npm test file1.test.ts`.
*   **Code Location:** `pkg/analyzer/test_finder.go`.

### Feature: Instant Container Sync (LivePatch)
*   **What it does:** Transfers a file from Host -> Container over TLS.
*   **How it works:**
    1.  CLI reads local file bytes.
    2.  CLI dials Agent via TLS.
    3.  Agent receives bytes, writes to disk, and executes optional `PostSyncCommand`.
*   **Security:** Uses mTLS (mutual authentication) to prevent unauthorized access.
*   **Code Location:** `pkg/transport/server.go` (Agent) and `client.go` (CLI).

### Feature: Speculative Merge Queue (Quantum Merge)
*   **What it does:** Queues Pull Requests to run CI in parallel.
*   **How it works:**
    1.  GitHub Webhook triggers `pkg/webhook/server.go`.
    2.  Job is added to `pkg/scheduler/queue.go` (In-Memory).
    3.  Workers (`pkg/scheduler/worker.go`) pick up jobs and simulate CI execution.
*   **Code Location:** `cmd/quantum-merge/main.go` wires the Webhook + Dashboard together.

---

## 4. DATA FLOW DOCUMENTATION

### LivePatch Sync Flow
1.  **Entry:** Developer runs `live-patch sync server.js`.
2.  **Validation:** CLI checks if file exists.
3.  **Transport:** CLI opens TCP connection to Agent (port 8080). TLS Handshake occurs.
4.  **RPC Call:** `LivePatchServer.SyncFile` is invoked.
5.  **Agent Action:** 
    *   Validates path (prevents `../` traversal).
    *   Writes file to disk.
    *   Runs `PostSyncCommand` (e.g., `npm restart`).
6.  **Response:** Returns `Success: true` or error message.

---

## 5. SETUP & RUNNING GUIDE

### Prerequisites
*   Go 1.23+
*   Docker (optional, for testing agents)

### Installation
```bash
git clone https://github.com/kusuridheeraj/velocity-trinity
cd velocity-trinity
go build ./...
```

### Running Locally (Development)
**1. Dependency-CI**
```bash
./dependency-ci analyze ./pkg/analyzer/analyzer.go
```

**2. LivePatch**
```bash
# Terminal 1: Start Agent
export PORT=8080
./live-patch-agent

# Terminal 2: Sync File
./live-patch sync README.md --target localhost:8080
```

**3. Quantum Merge**
```bash
# Start Server
./quantum-merge serve
# Visit http://localhost:8090/
```

---

## 6. TESTING GUIDE

### Test Structure
*   **Unit Tests:** Located alongside code (e.g., `pkg/analyzer/analyzer_test.go`).
*   **Integration Tests:** Currently manual (via the `test_project` scenarios).

### Running Tests
```bash
# Run all unit tests
go test ./...

# Run with coverage
go test -cover ./...
```

### Writing New Tests
*   Use the standard `testing` package.
*   For `LivePatch`, mock the `net.Conn` to test RPC without real networking.

---

## 7. COMMON DEVELOPMENT SCENARIOS

### Adding a New Language to Dependency-CI
1.  Create `pkg/analyzer/languages/rust.go`.
2.  Implement `Parse()` interface using Regex for `use crate::...`.
3.  Register it in `pkg/analyzer/analyzer.go` switch statement.
4.  Add a test case in `analyzer_test.go`.

### Debugging a "Stuck" Merge Queue
1.  Check the Dashboard (`http://localhost:8090`).
2.  Look at the logs (stdout of `quantum-merge` process).
3.  If a job is "RUNNING" forever, the worker goroutine might have panicked or deadlocked.

---

## 8. DEPENDENCIES
*   **`spf13/cobra`**: CLI scaffolding. Standard for Go CLIs.
*   **`spf13/viper`**: Configuration. Handles env vars/flags automatically.
*   **`gorilla/mux`**: HTTP Router. Flexible and standard.
*   **`go.uber.org/zap`**: Logging. Extremely fast structured logging.

---

## 9. SECURITY & BEST PRACTICES

*   **LivePatch Authentication**: Currently uses Self-Signed Certs for MVP. In production, this **MUST** be replaced with a proper CA (Certificate Authority) to prevent Man-in-the-Middle attacks.
*   **Path Traversal**: `pkg/transport/server.go` includes a check `filepath.HasPrefix` to ensure clients can't write to `/etc/passwd`.
*   **Secrets**: Secrets are scrubbed from logs by the Zap logger configuration in `pkg/logger`.

---

## 10. PERFORMANCE CONSIDERATIONS

*   **LivePatch Latency**: The handshake (TLS) takes ~50ms. The transfer is fast. Keep connections alive (Keep-Alive) in future versions to skip handshake.
*   **Dependency Parsing**: Regex parsing large files (1MB+) is slow. Use `bufio.Scanner` to stream line-by-line instead of loading full file into RAM.

---

## 11. DEPLOYMENT & DEVOPS

*   **CI Pipeline**: GitHub Actions (`.github/workflows/build.yml`) builds binaries for Linux/Windows/Mac on every push.
*   **Docker**: The Agent is designed to be `COPY --from=...` into existing containers.

---

## 12. REAL-WORLD PROBLEM SOLVING

**How do I add authentication to the Dashboard?**
1.  Add a Middleware in `pkg/dashboard/server.go`.
2.  Wrap the `router` in `RegisterRoutes`.
3.  Check for a Basic Auth header or Session Cookie.

**How do I fix a "Connection Refused" in LivePatch?**
1.  Check if Agent is running (`ps aux | grep agent`).
2.  Check if Port 8080 is exposed in Docker (`-p 8080:8080`).
3.  Check if TLS certs match (Client and Server must trust the same CA).

---

## 13. CRITICAL "IF REMOVED" ANALYSIS

*   **`pkg/transport/server.go` (Path Validation)**: If removed, a malicious developer could overwrite the host's `/etc/shadow` file via the Agent. **Critical Vulnerability.**
*   **`pkg/scheduler/worker.go` (Goroutine)**: If removed, jobs will sit in the queue forever. The API will accept them, but nothing will process them.
*   **`pkg/logger`**: If replaced with `fmt.Println`, we lose timestamps and log levels, making debugging impossible in production.
