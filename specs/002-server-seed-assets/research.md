# Research Log: DupLynx Server Runtime

## Decision 1: Deterministic demo reseeding strategy
- **Decision**: Recreate the SQLite demo database from scratch each time the seed command runs.
- **Rationale**: Ensures the dataset never drifts from documentation/tests, avoids complex merge logic, and guarantees CI gets a predictable baseline.
- **Alternatives considered**:
  - Preserve existing rows and upsert only known records — rejected because stale or user-modified data could linger and break demos.
  - Abort when data exists — rejected because it complicates automation and forces manual cleanup.

## Decision 2: Static asset delivery approach
- **Decision**: Serve precompiled Tailwind assets from a configurable disk directory and fail fast if missing.
- **Rationale**: Keeps binaries lightweight, allows front-end tweaks without rebuilding Go binaries, and aligns with standard Go + Tailwind workflows using external build steps.
- **Alternatives considered**:
  - Embed assets via `embed.FS` — rejected for slower rebuilds and larger binaries; not needed for local demo.
  - Dynamically running Tailwind at startup — rejected due to toolchain overhead and longer startup times.

## Decision 3: Audit logging for seed/serve lifecycle
- **Decision**: Emit structured audit events for every seed run and server lifecycle transition (start/stop) including actor, timestamp, and outcome.
- **Rationale**: Provides traceability, supports constitution observability goals, and simplifies troubleshooting in CI/onboarding flows.
- **Alternatives considered**:
  - Simple stdout messages — rejected because they are harder to parse and do not integrate with existing audit viewers.
  - Logging only failures — rejected since success events help confirm readiness and detect skipped steps.

## Decision 4: CI smoke coverage pattern
- **Decision**: Add a CLI-driven smoke test that seeds the database, starts the server on a random port, and verifies the root page renders key elements before teardown.
- **Rationale**: Validates end-to-end readiness in under a minute, aligns with constitution testing standards, and matches success criteria.
- **Alternatives considered**:
  - Rely solely on Playwright onboarding tests — rejected because they are heavier and may not run in all CI jobs.
  - Skip smoke tests and trust unit coverage — rejected due to higher regression risk for wiring/config issues.
