# Implementation Plan: Create DupLynx

**Branch**: `001-create-duplynx` | **Date**: 2025-10-27 | **Spec**: `specs/001-create-duplynx/spec.md`
**Input**: Feature specification from `/specs/001-create-duplynx/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Deliver a demo-ready DupLynx platform that seeds a sample tenant with five predefined machines, three scans, and duplicate group data accessible through a web dashboard. The backend Go services persist tenant-scoped metadata in SQLite via Ent, expose HTTP endpoints for dashboards and actions, and render templ-powered views enhanced by htmx and TailwindCSS to manage duplicate groups, keeper assignments, and lifecycle states.

## Technical Context

**Language/Version**: Go 1.22  
**Primary Dependencies**: templ, htmx, TailwindCSS, Ent (ORM)  
**Storage**: SQLite (file-backed) via Ent generated schema  
**Testing**: Go `testing` (unit/integration), Playwright or htmx harness for UI smoke, go test CI  
**Target Platform**: Linux containers and macOS dev environments via web dashboard  
**Project Type**: Web application with Go backend + server-rendered templ frontend  
**Performance Goals**: <1 s initial tenant/machine load, <200 ms board repaint with 200 duplicate groups, <500 ms action acknowledgement  
**Constraints**: Single-process demo deploy, no authentication, tenant isolation enforced at service/data layer, offline agents stubbed  
**Scale/Scope**: Phase 0 demo: 1 tenant, 5 machines, 3 scans, 200 duplicate groups per scan; architecture ready for multi-tenant growth

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- `Code Quality Fidelity`: Enforce `golangci-lint` and `go fmt` pre-commit, document new HTTP handlers, Ent schema, and templ components in package docstrings. Require code review sign-off with architecture notes captured in spec artifacts.
- `Testing Standardization`: Mandate `go test ./...` (unit + repository integration with SQLite in-memory), generated contract tests for API handlers, and Playwright smoke covering tenant→machine→scan selection; CI blocks on failures.
- `Consistent User Experience`: Apply shared Tailwind token presets, templ layout components, ARIA landmarks, and keyboard navigation for kanban lanes; run axe-core accessibility checks on critical screens.
- `Performance Discipline`: Track request timings via middleware logging, benchmark Ent query path for board load, and log action acknowledgement latency; budgets set to 1 s load, 200 ms repaint, 500 ms action response with follow-up alerts.

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
│   └── duplynx/          # main entrypoint
├── internal/
│   ├── app/              # configuration, bootstrap
│   ├── tenancy/          # tenant/machine context services
│   ├── scans/            # scan, duplicate group domain logic
│   ├── actions/          # keeper assignment and file action workflows
│   ├── http/             # handlers, middleware, routers
│   └── templ/            # templ components + view models
├── ent/                  # generated Ent schema + migrations
└── web/
    ├── tailwind.config.js
    └── static/           # compiled CSS, htmx helpers

tests/
├── unit/                 # package-scoped go tests
├── integration/          # sqlite-backed repository + HTTP tests
└── e2e/                  # Playwright/htmx harness scripts
```

**Structure Decision**: Adopt combined backend/web layout under `backend/` with shared templ/Tailwind assets and dedicated testing dirs to align Go monorepo conventions while isolating feature domains.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| _None_ | - | - |
