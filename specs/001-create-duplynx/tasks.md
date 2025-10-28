# Tasks: Create DupLynx

**Input**: Design documents from `/specs/001-create-duplynx/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Automated tests are required under the DupLynx constitution—plan linting, unit, integration, contract, UX, and performance coverage before implementation.

**Organization**: Tasks are grouped by user story so each slice is independently buildable, testable, and demoable.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Establish repo structure, tooling, and shared UI assets.

- [ ] T001 Scaffold project directories per implementation plan in `backend/` and `tests/`
- [ ] T002 Initialize Go module and workspace metadata in `backend/go.mod`
- [ ] T003 Create primary server entrypoint with flag parsing in `backend/cmd/duplynx/main.go`
- [ ] T004 [P] Configure linting/format tooling via `backend/.golangci.yml`
- [ ] T005 [P] Add task runner targets for lint/test in `Makefile`
- [ ] T006 [P] Configure Tailwind build pipeline in `backend/web/tailwind.config.js` and `backend/web/input.css`
- [ ] T007 [P] Scaffold base templ layout and shared components in `backend/internal/templ/layout.templ`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that all stories depend on—Ent schemas, routing, ingestion skeleton, CI, and tooling.

- [ ] T008 Define tenant schema with mixins in `backend/ent/schema/tenant.go`
- [ ] T009 [P] Define machine schema and relationships in `backend/ent/schema/machine.go`
- [ ] T010 [P] Define scan, duplicate group, file instance, and action audit schemas in `backend/ent/schema/scan.go`, `backend/ent/schema/duplicategroup.go`, `backend/ent/schema/fileinstance.go`, `backend/ent/schema/actionaudit.go`
- [ ] T011 Generate Ent client and migration files via `backend/ent/generate.go`
- [ ] T012 Implement configuration and SQLite connection bootstrap in `backend/internal/app/config.go`
- [ ] T013 Implement logging and metrics middleware scaffolding in `backend/internal/http/middleware/instrumentation.go`
- [ ] T014 Build initial chi router wiring templ renderer in `backend/internal/http/router.go`
- [ ] T015 [P] Implement ingestion handler stub accepting scan payloads in `backend/internal/ingestion/handler.go`
- [ ] T016 [P] Implement seed CLI command for demo data in `backend/cmd/duplynx/seed/main.go`
- [ ] T017 [P] Configure Playwright project and browsers in `tests/e2e/playwright.config.ts`
- [ ] T018 [P] Add CI workflow for lint, go test, Playwright smoke in `.github/workflows/ci.yml`

---

## Phase 3: User Story 1 – Tenant & Machine Onboarding Flow (Priority: P1) 🎯 MVP

**Goal**: Visitors select a tenant, choose one of its predefined machines, and reach the scan catalog without authentication.

**Independent Test**: Launch DupLynx, select “Sample Tenant A” then “Ares-Laptop,” verify machine metadata and seeded scans render.

### Tests

- [ ] T019 [P] [US1] Write contract tests for `/tenants` and `/tenants/{tenantSlug}/machines` in `tests/contract/tenants_machines_test.go`
- [ ] T020 [P] [US1] Create Playwright flow covering tenant and machine selection in `tests/e2e/onboarding.spec.ts`
- [ ] T021 [P] [US1] Add integration test for tenant machine filtering logic in `tests/integration/tenancy_flow_test.go`

### Implementation

- [ ] T022 [US1] Implement tenancy repository with scoped queries in `backend/internal/tenancy/repository.go`
- [ ] T023 [US1] Implement `/tenants` list handler in `backend/internal/http/handlers/tenants.go`
- [ ] T024 [US1] Implement tenant machine list handler in `backend/internal/http/handlers/machines.go`
- [ ] T025 [US1] Build launch and machine picker templ views in `backend/internal/templ/launch.templ`
- [ ] T026 [US1] Wire tenant/machine routes and context injection in `backend/internal/http/router.go`
- [ ] T027 [US1] Seed sample tenant and machine records for demo in `backend/internal/tenancy/seed.go`

---

## Phase 4: User Story 2 – Scan Management Board Overview (Priority: P1)

**Goal**: Users open a scan and view duplicate groups arranged by status lanes with performant rendering.

**Independent Test**: Open “Baseline Sweep 2025-10-01” board, verify each lane shows counts and metadata, expand a card to inspect file instances.

### Tests

- [ ] T028 [P] [US2] Write contract tests for `/scans` and `/duplicate-groups` endpoints in `tests/contract/scans_test.go`
- [ ] T029 [P] [US2] Add integration test assembling board data in `tests/integration/scan_board_test.go`
- [ ] T030 [P] [US2] Add unit tests for board view model grouping in `tests/unit/board_view_test.go`

### Implementation

- [ ] T031 [US2] Implement scan repository with status aggregations in `backend/internal/scans/repository.go`
- [ ] T032 [US2] Implement board service composing lanes in `backend/internal/scans/service.go`
- [ ] T033 [US2] Implement scan board HTTP handlers and JSON responses in `backend/internal/http/handlers/scan_board.go`
- [ ] T034 [US2] Build templ components for board shell and lanes in `backend/internal/templ/board.templ`
- [ ] T035 [US2] Add htmx partial handlers for lane refresh in `backend/internal/http/handlers/board_partials.go`
- [ ] T036 [US2] Instrument board rendering timings and virtualized list hooks in `backend/internal/http/middleware/board_metrics.go`

---

## Phase 5: User Story 3 – Keeper Assignment & Duplicate Actions (Priority: P2)

**Goal**: Stewards assign keeper machines and trigger duplicate actions with audit logging and responsive UI feedback.

**Independent Test**: Assign “Helios-Server-02” as keeper for a group, execute a quarantine action, verify audit trail and lane status remain consistent.

### Tests

- [ ] T037 [P] [US3] Write contract tests for `/duplicate-groups/{groupId}/keeper` and `/duplicate-groups/{groupId}/actions` in `tests/contract/actions_test.go`
- [ ] T038 [P] [US3] Add integration test for keeper assignment state transitions in `tests/integration/keeper_assignment_test.go`
- [ ] T039 [P] [US3] Add unit tests for action dispatcher outcomes in `tests/unit/actions_dispatcher_test.go`

### Implementation

- [ ] T040 [US3] Implement action dispatcher coordinating keeper and action workflows in `backend/internal/actions/dispatcher.go`
- [ ] T041 [US3] Implement action audit persistence layer in `backend/internal/actions/audit_store.go`
- [ ] T042 [US3] Implement keeper and action HTTP handlers with validation in `backend/internal/http/handlers/actions.go`
- [ ] T043 [US3] Enhance duplicate card templ with keeper selection UI in `backend/internal/templ/components/duplicate_card.templ`
- [ ] T044 [US3] Implement htmx response fragments for action feedback in `backend/internal/http/handlers/actions_htmx.go`

---

## Phase 6: User Story 4 – Multi-Tenant Data Isolation (Priority: P3)

**Goal**: Guarantee tenant-scoped access so cross-tenant requests fail safely and are observable.

**Independent Test**: When authenticated for “Sample Tenant A,” attempting to fetch a “Sample Tenant B” scan returns a tenant scope warning and 404.

### Tests

- [ ] T045 [P] [US4] Add unit tests for tenancy middleware context enforcement in `tests/unit/tenancy_middleware_test.go`
- [ ] T046 [P] [US4] Add integration test for cross-tenant access rejection in `tests/integration/tenant_guard_test.go`

### Implementation

- [ ] T047 [US4] Implement tenancy scoping middleware attaching tenant context in `backend/internal/tenancy/middleware.go`
- [ ] T048 [US4] Apply tenant filters across repositories in `backend/internal/tenancy/scoped_repository.go`
- [ ] T049 [US4] Create tenant scope violation templ feedback in `backend/internal/templ/errors/tenant_scope.templ`

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, benchmarks, accessibility, and release readiness.

- [ ] T050 Update quickstart instructions with final commands in `specs/001-create-duplynx/quickstart.md`
- [ ] T051 Document deployment and scaling notes in `docs/duplynx-demo.md`
- [ ] T052 Add performance benchmark covering board render latency in `tests/perf/board_bench_test.go`
- [ ] T053 Run accessibility and contrast checks via Playwright axe audit in `tests/e2e/accessibility.spec.ts`
- [ ] T054 Add final CI verification target for go test/playwright combo in `Makefile`

---

## Dependencies & Execution Order

- **Phase 1 → Phase 2**: Foundational work depends on setup.
- **Phase 2 → Phases 3-6**: All user stories require Ent schemas, routing, ingestion, and CI tooling.
- **Phase 3 (US1)**: Enables tenant context and seeding—MVP checkpoint.
- **Phase 4 (US2)**: Depends on US1 data structures but can progress in parallel once foundational data services exist.
- **Phase 5 (US3)**: Depends on US2 duplicate group representations to attach actions.
- **Phase 6 (US4)**: Builds on tenancy infrastructure from US1 and repository layers from US2/US3.
- **Phase 7**: Runs after desired user stories complete; captures documentation and performance compliance.

### Story Dependency Graph

`US1 → {US2, US4}`  
`US2 → US3`  
`US3 → Phase 7 readiness`  
`US4 → Phase 7 readiness`

---

## Parallel Execution Examples

### User Story 1
```bash
# Parallelizable tasks:
# Tests
run tests/contract/tenants_machines_test.go
run tests/e2e/onboarding.spec.ts
run tests/integration/tenancy_flow_test.go

# Views and handlers
edit backend/internal/http/handlers/tenants.go
edit backend/internal/templ/launch.templ
```

### User Story 2
```bash
# Tests
run tests/contract/scans_test.go
run tests/integration/scan_board_test.go
run tests/unit/board_view_test.go

# Implementation
edit backend/internal/scans/repository.go
edit backend/internal/templ/board.templ
```

### User Story 3
```bash
# Tests
run tests/contract/actions_test.go
run tests/integration/keeper_assignment_test.go
run tests/unit/actions_dispatcher_test.go

# Implementation
edit backend/internal/actions/dispatcher.go
edit backend/internal/http/handlers/actions.go
```

### User Story 4
```bash
# Tests
run tests/unit/tenancy_middleware_test.go
run tests/integration/tenant_guard_test.go

# Implementation
edit backend/internal/tenancy/middleware.go
edit backend/internal/templ/errors/tenant_scope.templ
```

---

## Implementation Strategy

### MVP First (deliver US1)
1. Complete Phases 1-2 to establish infrastructure.
2. Execute Phase 3 (US1) to unlock tenant selection MVP.
3. Validate onboarding via contract + e2e tests before proceeding.

### Incremental Delivery
1. After MVP, deliver Phase 4 (US2) to expose the scan board.
2. Layer Phase 5 (US3) for keeper actions and auditing.
3. Finish with Phase 6 (US4) to harden tenant boundaries.
4. Close with Phase 7 polish for documentation, accessibility, and performance signage.

### Parallel Team Strategy
1. Run Phases 1-2 with shared effort.
2. Once Phase 2 completes:
   - Developer A: Phase 3 (US1)
   - Developer B: Phase 4 (US2)
   - Developer C: Phase 5 (US3)
   - Developer D: Phase 6 (US4)
3. Sync on Phase 7 for cross-cutting polish and release prep.

---

## Task Counts & Coverage

- **Total tasks**: 54
- **Per user story**:
  - US1: 9 tasks
  - US2: 9 tasks
  - US3: 8 tasks
  - US4: 5 tasks
- **Parallel opportunities**: Tests and component work within each story marked `[P]` can proceed concurrently once prerequisites land.
- **Independent test criteria**:
  - US1: Tenant and machine selection flow validates seed data and navigation.
  - US2: Scan board lanes render with correct counts and card expansion details.
  - US3: Keeper assignment and actions persist state with audit trail feedback.
  - US4: Cross-tenant requests fail with scoped warnings and 404s.
- **Suggested MVP scope**: Deliver through Phase 3 (US1) to demonstrate tenant/machine onboarding end-to-end.

All tasks follow the required checklist format with sequential IDs, `[P]` markers where applicable, and explicit file paths.
