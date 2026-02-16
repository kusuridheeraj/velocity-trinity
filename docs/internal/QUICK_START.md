# Quick Start Guide: Velocity Trinity (5 Minutes)

## Goal
Run all three tools (`dependency-ci`, `live-patch`, `quantum-merge`) locally in 5 minutes.

## Prerequisites
*   Go 1.23+
*   Git

## Step 1: Clone & Build (1 Minute)
```bash
git clone https://github.com/kusuridheeraj/velocity-trinity
cd velocity-trinity

# Build binaries
go build ./...
```

## Step 2: Run Dependency-CI (1 Minute)
```bash
# Analyze a file
./dependency-ci analyze ./pkg/analyzer/analyzer.go

# Find tests for a file
./dependency-ci run --files="./pkg/analyzer/analyzer.go" --cmd="echo [TEST]"
```

## Step 3: Run LivePatch (2 Minutes)
```bash
# Start Agent (Server)
export PORT=8080
./live-patch-agent &

# Create a test file
echo "Hello World" > test.txt

# Sync file to Agent
./live-patch sync test.txt --target localhost:8080
```

## Step 4: Run Quantum Merge (1 Minute)
```bash
# Start Server
./quantum-merge serve &

# Open Dashboard
# Visit http://localhost:8090/
```

## Troubleshooting
*   **"command not found"**: Ensure `.` is in your PATH or use `./binary-name`.
*   **"bind: address already in use"**: Kill processes on port 8080/8090 (`lsof -i :8080`).
*   **"connection refused"**: Agent crashed? Check logs.
