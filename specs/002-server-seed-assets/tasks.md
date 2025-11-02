---

description: "Task list for DupLynx Server Runtime feature"
---

# Tasks: DupLynx Server Runtime

**Input**: Design documents from `/specs/002-server-seed-assets/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are explicitly requested in the specification for smoke coverage, Playwright onboarding, and seeded demo validation. Test tasks are included where requirements call them out.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Establish baseline tooling and directories required by all user stories.

- [X] T001 Update Tailwind toolchain and add `build:tailwind` script targeting `backend/web/dist` in `package.json`
- [X] T002 Add tracked asset output directory placeholder at `backend/web/dist/.gitkeep`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [X] T003 Replace the stub CLI entrypoint with Cobra root execution in `backend/cmd/duplynx/main.go`
- [X] T004 Create Cobra root command scaffolding with persistent config flags in `backend/cmd/duplynx/root.go`
- [X] T005 Implement runtime configuration schema with env/flag binding in `backend/internal/config/runtime.go`
- [X] T006 [P] Build SQLite connection and migration helper in `backend/internal/data/store.go`
- [X] T007 [P] Add static asset directory validation helper in `backend/internal/app/assets.go`
- [X] T008 [P] Introduce audit event writer utility for CLI commands in `backend/internal/observability/audit.go`
- [X] T009 Construct reusable HTTP server builder that wires Chi router and middleware in `backend/internal/app/server.go`

**Checkpoint**: Configuration, database, audit, and server scaffolding ready — user story implementation can now begin.

---

## Phase 3: User Story 1 - Launch demo environment (Priority: P1) 🎯 MVP

**Goal**: Enable evaluators to run `duplynx serve` and load the dashboard with tenant-aware routes backed by the database.

**Independent Test**: Run `go run ./cmd/duplynx serve --db-file var/duplynx.db --assets-dir backend/web/dist` against a seeded database stub and confirm `/` returns HTML with tenant, machine, and scan information while invalid configs fail fast.

### Tests for User Story 1 ⚠️

- [X] T010 [P] [US1] Create integration test covering serve command failure modes in `tests/integration/server/serve_config_validation_test.go`

### Implementation for User Story 1

- [X] T011 [US1] Implement `serve` subcommand to bootstrap config, DB, and HTTP server in `backend/cmd/duplynx/serve.go`
- [X] T012 [P] [US1] Refactor tenancy repository to query tenants and machines via Ent client in `backend/internal/tenancy/repository.go`
- [X] T013 [P] [US1] Refactor scans repository to query scan summaries via Ent client in `backend/internal/scans/repository.go`
- [X] T014 [P] [US1] Update HTTP router to mount disk-backed static assets and ent-powered handlers in `backend/internal/http/router.go`
- [X] T015 [US1] Emit audit events for server start/stop lifecycle in `backend/cmd/duplynx/serve.go`
- [X] T016 [US1] Document `duplynx serve` usage and flags in `docs/duplynx-demo.md`

**Checkpoint**: Launch command serves templ views using database data and respects asset gating.

### Parallel Execution Example (US1)

```bash
# Parallelizable tasks once T011 is underway
#   - T012 (tenancy repo refactor)
#   - T013 (scan repo refactor)
#   - T014 (router updates)
```

---

## Phase 4: User Story 2 - Refresh demo dataset (Priority: P2)

**Goal**: Provide `duplynx seed` command that rebuilds the SQLite demo database deterministically with canonical tenants, machines, scans, duplicate groups, and audit history.

**Independent Test**: Run `go run ./cmd/duplynx seed --db-file var/duplynx.db --assets-dir backend/web/dist` twice and verify the database contents match canonical fixtures with fresh audit entries and no duplicated records.

### Tests for User Story 2 ⚠️

- [ ] T017 [P] [US2] Add integration test validating seeded tenants and machines in `tests/integration/seeding/seed_demo_dataset_test.go`
- [ ] T018 [P] [US2] Add unit test ensuring seed audit events capture actor/outcome in `tests/unit/observability/seed_audit_test.go`

### Implementation for User Story 2

- [ ] T019 [US2] Define canonical demo dataset builders for tenants, machines, scans, and groups in `backend/internal/data/demo.go`
- [ ] T020 [P] [US2] Implement deterministic reseed workflow that drops and re-applies schema in `backend/internal/data/seed.go`
- [ ] T021 [P] [US2] Register `seed` subcommand with Cobra root in `backend/cmd/duplynx/root.go`
- [ ] T022 [P] [US2] Extend seed command to emit audit events and duration metrics in `backend/cmd/duplynx/seed.go`
- [ ] T023 [US2] Replace in-memory demo fixtures with database-backed helpers in `tests/integration/tenant_guard_test.go`
- [ ] T024 [P] [US2] Update contract tests to initialize seeded database fixtures in `tests/contract/actions_test.go`
- [ ] T025 [US2] Update documentation with `duplynx seed` workflow in `docs/duplynx-demo.md`

**Checkpoint**: Seed command recreates demo data deterministically and all tests exercise database-backed fixtures.

### Parallel Execution Example (US2)

```bash
# After T019 completes, run in parallel:
#   - T020 (reseed workflow)
#   - T021 (command registration)
#   - T022 (audit instrumentation)
```

---

## Phase 5: User Story 3 - CI confidence in runtime (Priority: P3)

**Goal**: Add automated smoke coverage so CI seeds the database, starts the server, and blocks promotion when core routes fail.

**Independent Test**: Execute `make smoke-demo` locally or in CI; ensure the command seeds the DB, starts the server, performs HTTP assertions on `/`, and tears down cleanly.

### Tests for User Story 3 ⚠️

- [ ] T026 [P] [US3] Create Go-based smoke test that seeds, serves, probes `/`, and asserts timing budgets in `tests/smoke/server_smoke_test.go`
- [ ] T027 [P] [US3] Update Playwright onboarding spec to consume seeded data in `tests/e2e/onboarding.spec.ts`

### Implementation for User Story 3

- [ ] T028 [US3] Add `smoke-demo` target chaining seed and serve verification with timing log output in `Makefile`
- [ ] T029 [P] [US3] Update CI workflow to call smoke target after `make ci` and surface timing metrics in `.github/workflows/ci.yml`
- [ ] T030 [US3] Document smoke test invocation and CI expectations in `docs/duplynx-demo.md`

**Checkpoint**: CI executes smoke tests automatically and documentation guides engineers through verification.

### Parallel Execution Example (US3)

```bash
# With seed/serve flow stable:
#   - T026 (smoke test) and T027 (Playwright updates) can proceed together.
```

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Remove deprecated artifacts and ensure quality across the updated runtime.

- [ ] T031 Remove legacy in-memory duplicate group store now superseded by Ent in `backend/internal/actions/store.go`
- [ ] T032 Update remaining unit tests to use seeded database helpers in `tests/unit/actions_dispatcher_test.go`
- [ ] T033 Run `go fmt` and `golangci-lint` over new packages in `backend/internal/...`
- [ ] T034 Verify quickstart instructions by executing documented flow, recording timing output, and capturing notes in `docs/duplynx-demo.md`
- [ ] T035 Document onboarding support baseline and tracking plan for SC-004 in `docs/duplynx-demo.md`
- [ ] T036 Create automated quickstart timing script in `scripts/measure_quickstart.sh` and document usage in `docs/duplynx-demo.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — execute immediately.
- **Foundational (Phase 2)**: Depends on Phase 1 — BLOCKS all user stories.
- **User Story Phases (3-5)**: Each depends on completion of Phase 2. Stories may proceed in priority order (US1 → US2 → US3) or in parallel once shared infrastructure is stable.
- **Polish (Phase 6)**: Depends on completion of desired user stories.

### User Story Dependencies

- **US1 (P1)**: Requires foundational configuration/DB/audit scaffolding.
- **US2 (P2)**: Requires US1 database read paths to exist; produces deterministic data for downstream flows.
- **US3 (P3)**: Requires US1 and US2 to ensure smoke tests have functioning seed and serve commands.

### Story Completion Graph

```
Phase 1 → Phase 2 → US1 → US2 → US3 → Phase 6
```

### Parallel Opportunities Summary

- Setup tasks can run sequentially while dependencies resolve.
- Foundational tasks T006–T008 can run in parallel once T003–T005 begin.
- Within US1, repository refactors (T012–T014) are parallelizable after T011 stubs the command.
- Within US2, seeding workflow (T020), command registration (T021), and audit wiring (T022) can run concurrently once dataset builders (T019) exist.
- Within US3, smoke and Playwright updates (T026, T027) can progress in parallel before CI wiring (T029).

---

## Implementation Strategy

### MVP First (Deliver US1)

1. Complete Phases 1 and 2 to establish configuration, database access, and audit scaffolding.
2. Implement US1 (Phase 3) to deliver a runnable server backed by the database.
3. Validate using T010 integration test and manual quickstart steps.

### Incremental Delivery

- **Increment 1**: Deliver US1 — evaluators can launch the server with seeded data stubs.
- **Increment 2**: Deliver US2 — deterministic seeding command available and tests migrated to database fixtures.
- **Increment 3**: Deliver US3 — CI smoke coverage and Playwright updates guarding regressions.

### Parallel Team Strategy

- Developer A: Focus on Foundational tasks (T003–T009) then US1 command implementation.
- Developer B: After foundational work, tackle US1 repository refactors (T012–T014) and US2 dataset builders (T019–T020).
- Developer C: Own US3 automation (T026–T029) once US1/US2 are integrated.

---
