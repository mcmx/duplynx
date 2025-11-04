# Implementation Plan: DupLynx Server Runtime

**Branch**: `[002-server-seed-assets]` | **Date**: 2025-11-01 | **Spec**: [link](../spec.md)  
**Input**: Feature specification from `/specs/002-server-seed-assets/spec.md`

## Summary

Deliver a production-like DupLynx CLI that offers `serve` and `seed` commands, recreates the SQLite demo database deterministically, serves the existing Chi/templ UI with Tailwind assets from disk, and adds CI smoke coverage with timing instrumentation so evaluators can reach a fully populated dashboard in minutes while meeting performance budgets.

## Technical Context

**Language/Version**: Go 1.22  
**Primary Dependencies**: Chi router, Ent ORM, templ view engine, TailwindCSS build artifacts, Cobra-style CLI scaffolding  
**Storage**: SQLite demo database file on local disk (`var/duplynx.db`)  
**Testing**: `go test` suites (unit/integration), Playwright onboarding tests, contract/integration suites, new HTTP smoke test for `/`, automated timing assertions for seed/serve workflows, quickstart timing script (`scripts/measure_quickstart.sh`)  
**Target Platform**: Localhost server on Linux/macOS (bind `0.0.0.0:8080`)  
**Project Type**: Backend CLI-managed web server with server-rendered UI  
**Performance Goals**: Seed workflow <60s; dashboard available <5 minutes from kickoff; root route loads <2s on standard dev hardware  
**Constraints**: Fail fast on invalid config, require Tailwind assets present on disk before serving, reseed database on demand without partial leftovers, log measurable timing evidence for success criteria  
**Scale/Scope**: Single-node demo usage with <20 canonical tenants/machines/scans groups

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Code Quality Fidelity**: Changes extend existing Go module with documented CLI commands, eliminate stubbed data, and add docs/tests to maintain clarity. ✅
- **Testing Standardization**: Plan includes unit/integration/Playwright updates plus a new deterministic smoke test to protect coverage. ✅
- **Consistent User Experience**: Uses shared templ layouts and Tailwind tokens; quickstart documents the expected flows to preserve consistency. ✅
- **Performance Discipline**: Success criteria define timing budgets; plan adds audit logging and smoke checks to monitor runtime responsiveness. ✅

Gate verdict: PASS – proceed to Phase 0 research.

Post-Phase 1 verification: Completed design artifacts (data model, contracts, quickstart) uphold code quality, testing, UX, and performance principles with explicit tasks for logging, asset validation, smoke coverage, and performance instrumentation. ✅

## Project Structure

### Documentation (this feature)

```text
specs/002-server-seed-assets/
├── plan.md              # This file (/speckit.plan output)
├── research.md          # Phase 0 research synthesis
├── data-model.md        # Phase 1 entity/state design
├── quickstart.md        # Phase 1 operator instructions
├── contracts/           # Phase 1 API/route contracts
└── tasks.md             # Produced later by /speckit.tasks
```

### Source Code (repository root)

```text
backend/
├── cmd/
│   └── duplynx/          # CLI entrypoint, will add serve/seed commands and flag handling
├── internal/
│   ├── app/              # Existing router, middleware, templ renderers to be wired
│   ├── config/           # Config loading utilities to extend for CLI flags/env
│   ├── data/             # (New) demo dataset builders/seed routines
│   └── observability/    # (New or extended) audit logging helpers
├── ent/                  # Generated ORM models & migrations
├── web/
│   ├── tailwind.config   # Asset pipeline sources
│   └── dist/             # Built CSS/JS served from disk
└── go.mod

tests/
├── integration/          # Extend with seeded DB scenarios
├── contract/             # Update for router expectations
└── smoke/                # Add CLI-driven HTTP smoke test hitting `/`
```

**Structure Decision**: Extend the existing `backend` Go module with CLI commands and supporting packages while centralizing demo seeding under `internal/data`. Serve Tailwind assets from `backend/web/dist`, and place new smoke checks alongside existing test suites under `tests/`.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| _None_ | — | — |
