# ADR-003: Dependency Parsing Strategy (MVP)

## Context
Product A (`dependency-ci`) needs to build a dependency graph of the codebase to determine which tests to run. Accurate parsing requires understanding the Abstract Syntax Tree (AST) of various languages (TypeScript, Python, Go).

## Options
1.  **Tree-sitter (C-bindings)**: High accuracy, but requires CGO and C compiler on the host machine, complicating cross-compilation and Windows support.
2.  **Language Servers (LSP)**: Extremely accurate but heavy and slow to start up for a quick CLI check.
3.  **Regex / Tokenization**: Fast, zero-dependency, but potential for false positives (e.g., matching imports inside comments).

## Decision
We will use **Regex/Tokenization** for the initial MVP phase.

## Rationale
1.  **Speed**: Regex scanning is orders of magnitude faster than spinning up LSPs.
2.  **Portability**: Pure Go implementation ensures the binary works on Windows/Linux/Mac without `gcc`.
3.  **MVP Scope**: We can achieve 90% accuracy with smart regexes (ignoring comments). The goal is to *reduce* the test set, not be perfect. If we accidentally run an extra test (false positive), it's acceptable. Missing a test (false negative) is the risk, which we mitigate by falling back to "run all" on ambiguous parsing.

## Consequences
-   **Positive**: Rapid development, single binary.
-   **Negative**: Edge cases in complex import syntax might be missed.
-   **Mitigation**: We will add a `--strict` mode later that uses a Dockerized parser if higher accuracy is needed.
