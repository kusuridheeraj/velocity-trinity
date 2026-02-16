# ADR-008: LivePatch Command Execution

## Context
Syncing files (`live-patch sync`) is only half the battle. If I update `server.js`, I need the node process to restart.
We need a way to execute commands inside the container *after* the file sync is complete.

## Decision
We will add a `--restart` flag to the `sync` command.

## Protocol Update
We will modify the `FileSyncRequest` struct to include a `PostSyncCommand` string.
- Client: Sends `sync server.js --restart="npm restart"`
- Agent: Writes file -> Executes `npm restart` -> Returns "Sync & Restart Successful".

## Security Implications
-   **Risk**: Remote Code Execution (RCE).
-   **Mitigation**: The Agent is *designed* to be a development tool. It should **NEVER** be installed in production containers accessible from the public internet without mutual TLS (mTLS) and strict firewall rules.
-   **Warning**: We will add a startup warning log: "WARNING: LivePatch Agent allows remote command execution. Do not expose to public internet."
