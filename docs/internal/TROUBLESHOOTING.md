# Troubleshooting Guide: Velocity Trinity

## Common Issues & Solutions

### 1. Build Failures
**Error:** `go: module not found` or `missing go.sum`.
*   **Cause:** Go modules are out of sync.
*   **Solution:** Run `go mod tidy` in the root directory.

**Error:** `cgo: C compiler "gcc" not found`.
*   **Cause:** You are trying to compile a dependency that uses CGO on Windows without GCC installed.
*   **Solution:** Install TDM-GCC or MinGW. Or, ensure `CGO_ENABLED=0` environment variable is set.

### 2. LivePatch Issues
**Error:** `rpc: can't find service LivePatchServer.SyncFile`.
*   **Cause:** Mismatch between Client and Agent versions (different struct definitions).
*   **Solution:** Rebuild both binaries (`go build ./...`) to ensure they use the same `pkg/transport`.

**Error:** `x509: certificate signed by unknown authority`.
*   **Cause:** Client does not trust the Agent's self-signed certificate.
*   **Solution:** Run with `--insecure` (if implemented) or add the CA to your system trust store. For dev, ensure `InsecureSkipVerify: true` is set in client code.

### 3. Dependency-CI Issues
**Error:** `panic: runtime error: index out of range`.
*   **Cause:** Regex parser failed on an empty line or malformed import.
*   **Solution:** Check the input file. The regex parser is fragile. Add a guard clause in `pkg/analyzer/languages/*.go`.

**Error:** `exec: "npm": executable file not found in %PATH%`.
*   **Cause:** `npm` is not in your system PATH or you are on Windows and need `npm.cmd`.
*   **Solution:** Use `npm.cmd` on Windows. Or ensure Node.js is installed.

### 4. Quantum Merge Issues
**Error:** `bind: address already in use`.
*   **Cause:** Another instance of `quantum-merge` is running on port 8090.
*   **Solution:** `killall quantum-merge` or `taskkill /F /IM quantum-merge.exe`.

**Error:** Dashboard shows blank page / 404 API errors.
*   **Cause:** The embedded UI files are missing or the router configuration is wrong.
*   **Solution:** Ensure `//go:embed ui/*` is working correctly and the `ui/` folder exists at build time.
