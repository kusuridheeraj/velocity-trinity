# ADR-002: Go (Golang) as Primary Language

## Context
We are building CLI tools and system agents that need to:
1.  Run across multiple operating systems (Linux, Windows, macOS).
2.  Have minimal deployment dependencies (customers should just download a binary).
3.  Execute quickly (parse large codebases, handle concurrent network requests, process millions of lines of logs).
4.  Interact deeply with system resources (filesystems, network interfaces).

## Decision
We will use **Go (Golang)** as the primary language for all CLI tools and agents.

## Rationale
1.  **Performance**: Go compiles to efficient native machine code, critical for processing speed.
2.  **Concurrency**: Goroutines provide a lightweight model for handling thousands of concurrent operations (crucial for `quantum-merge`'s parallel execution and `live-patch`'s network handling).
3.  **Portability**: Go's cross-compilation is trivial (`GOOS=linux go build`).
4.  **Ecosystem**: Excellent support for CLIs (`cobra`, `viper`), Parsing (`tree-sitter` bindings), and Kubernetes (`client-go`).
5.  **Single Binary**: Delivers a zero-dependency executable, simplifying distribution and reducing customer friction.

## Consequences
-   **Positive**: High performance, easy distribution, strong typing for safety.
-   **Negative**: Lack of some high-level abstractions found in Python (e.g., dynamic typing for rapid scripting), but this is acceptable for system tools.
-   **Mitigation**: Use Python or Node.js only where absolutely necessary (e.g., specific language parsing libraries if Go lacks bindings), invoked as subprocesses if needed.
