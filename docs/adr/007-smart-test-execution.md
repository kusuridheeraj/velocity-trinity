# ADR-007: Smart Test Execution (Product A - Phase 2)

## Context
Product A (`dependency-ci`) currently analyzes code dependencies (`src/foo.ts` depends on `src/bar.ts`). Now, it needs to understand *tests* (`src/foo.test.ts` imports `src/foo.ts`).

## Decision
We will implement a simple **naming convention heuristic** alongside the import graph.

## Rationale
1.  **Complexity**: Parsing test frameworks (Jest/Pytest) ASTs to find exactly which test calls which function is extremely complex.
2.  **Accuracy**: A naming convention (`foo.ts` -> `foo.test.ts` or `test_foo.py`) covers 95% of use cases.
3.  **Speed**: Checking filenames is O(1).

## Strategy
1.  When `foo.ts` changes:
    - Find all files that import `foo.ts` (using our existing `analyzer`).
    - For each affected file, check if a corresponding test file exists (e.g., `foo.test.ts`).
    - Add test file to the execution list.

## Consequences
-   **Positive**: Extremely fast execution.
-   **Negative**: Doesn't catch "integration tests" that test multiple modules without importing them directly (e.g., API tests).
-   **Mitigation**: Users can tag integration tests as "always run" via configuration.
