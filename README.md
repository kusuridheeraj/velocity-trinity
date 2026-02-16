# Welcome to Velocity Trinity

Velocity Trinity is a suite of high-performance DevOps tools designed to fix the "Integration Bottleneck."

## The Suite

### 1. Dependency-CI (`dependency-ci`)
**Smart Test Orchestration.**
Stop running 100% of your tests for a 1-line change.
- [Documentation](./dependency-ci/README.md)
- [Architecture](./docs/adr/007-smart-test-execution.md)

### 2. LivePatch (`live-patch`)
**Instant Container Sync.**
Stop waiting for Docker builds. Patch running containers in milliseconds.
- [Documentation](./live-patch/README.md)
- [Architecture](./docs/adr/004-communication-protocol.md)

### 3. Quantum Merge (`quantum-merge`)
**Speculative Merge Queue.**
Run CI for PR #2 assuming PR #1 passes.
- [Documentation](./quantum-merge/README.md)
- [Architecture](./docs/adr/005-queue-strategy.md)

## Getting Started

### Prerequisites
- Go 1.23+
- Docker (optional, for testing agents)

### Installation
```bash
git clone https://github.com/velocity-trinity/velocity-trinity
cd velocity-trinity
go build ./...
```

### Running the Dashboard
```bash
./quantum-merge serve
# Open http://localhost:8090
```
