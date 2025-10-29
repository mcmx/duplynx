# Implementation Plan: Create DupLynx

**Branch**: `001-create-duplynx` | **Date**: 2025-10-27 | **Spec**: `specs/001-create-duplynx/spec.md`
**Input**: Feature specification from `/specs/001-create-duplynx/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Build the initial DupLynx platform: a Go-based multi-tenant service backed by SQLite/Ent that seeds demo data, ingests HMAC-signed file scan results from predefined machines, and serves both an ingestion API and templ/htmx-powered dashboard to manage duplicate groups, keeper assignments, and lifecycle statuses.

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go 1.22  
**Primary Dependencies**: Ent ORM, templ, htmx, TailwindCSS, chi (HTTP router), golangci-lint  
**Storage**: SQLite (file-backed) managed via Ent migrations  
**Testing**: Go `testing` (unit + integration), httptest contract suites, Playwright/htmx UI smoke  
**Target Platform**: Linux/macOS servers for API + dashboard; stateless ingestion and GUI instances  
**Project Type**: Web application (Go backend with server-rendered frontend)  
**Performance Goals**: <1 s tenant/machine load, <200 ms board update, <500 ms ingestion acknowledgement with metrics coverage  
**Constraints**: Demo must run without external auth, SQLite persisted on a designated ingestion writer (single-writer, multi-reader) while additional GUI pods operate read-only  
**Scale/Scope**: Phase 0 seed: 1 tenant, 5 machines, 3 scans, up to 200 duplicate groups per scan; ready for multiple ingest GUI pods

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- `Code Quality Fidelity`: Enforce `golangci-lint` and `go fmt` in CI, require package-level docstrings for new modules (tenancy, scans, ingestion), and capture architecture notes in spec updates; code review approval mandatory.
- `Testing Standardization`: Run `go test ./...` (unit/integration), contract suites for ingestion API + dashboard JSON, and Playwright smoke tests; CI blocks on any failures or coverage <90% for touched packages.
- `Consistent User Experience`: Use templ components with Tailwind tokens, ensure ARIA landmarks and keyboard traversal for kanban board, and run axe-core accessibility checks on tenant pickers and scan board.
- `Performance Discipline`: Instrument middleware timing for ingestion + board endpoints, benchmark Ent queries for board load, log action acknowledgements, and capture ingestion latency metrics; budgets set to 1 s load, 200 ms updates, 500 ms ingestion response with alerts on breach.

_Post-design review (2025-10-27): Phase 1 artifacts uphold all gates; no exceptions required._

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
backend/
├── cmd/
│   └── duplynx/              # main server (API + dashboard) and seed commands
├── internal/
│   ├── app/                  # configuration, startup wiring
│   ├── tenancy/              # tenant & machine services, middleware
│   ├── ingestion/            # file scan ingestion handlers & validation
│   ├── scans/                # scan & duplicate group domain logic
│   ├── actions/              # keeper assignment & duplicate actions
│   ├── http/                 # routing, middleware, presenters
│   └── templ/                # templ components & view models
├── ent/                      # Ent schemas and generated code
└── web/
    ├── tailwind.config.js
    ├── input.css
    └── static/               # built CSS, htmx helpers

tests/
├── unit/
├── integration/
└── e2e/                      # Playwright/htmx scripts
```

**Structure Decision**: Consolidate backend and UI under `backend/` to keep Go modules cohesive while separating domain packages; dedicated ingestion package supports scaling multiple API instances, `tests/` hosts tiered suites aligned with constitution gates.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| _None_ | - | - |
