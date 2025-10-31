# Tasks: Create DupLynx

**Input**: Design documents from `/specs/001-create-duplynx/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Automated tests are required—plan linting, unit, integration, contract, UX, and performance coverage before implementation.

**Organization**: Tasks are grouped by user story so each slice is independently buildable, testable, and demoable.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Establish repo structure, tooling, and shared UI assets.

- [X] T001 Scaffold project directories per implementation plan in `backend/` and `tests/`
- [X] T002 Initialize Go module and workspace metadata in `backend/go.mod`
- [X] T003 Create primary server entrypoint with flag parsing in `backend/cmd/duplynx/main.go`
- [X] T004 [P] Configure linting/format tooling via `backend/.golangci.yml`
- [X] T005 [P] Add task runner targets for lint/test in `Makefile`
- [X] T006 [P] Configure Tailwind build pipeline in `backend/web/tailwind.config.js` and `backend/web/input.css`
- [X] T007 [P] Scaffold base templ layout and shared components in `backend/internal/templ/layout.templ`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure—Ent schemas, ingestion security, routing, CI, and performance baselines.

- [X] T008 Define tenant schema with mixins in `backend/ent/schema/tenant.go`
- [X] T009 [P] Define machine schema and relationships in `backend/ent/schema/machine.go`
- [X] T010 [P] Define scan, duplicate group, file instance, and action audit schemas in `backend/ent/schema/scan.go`, `backend/ent/schema/duplicategroup.go`, `backend/ent/schema/fileinstance.go`, `backend/ent/schema/actionaudit.go`
- [X] T011 Generate Ent client and migration files via `backend/ent/generate.go`
- [X] T012 Implement configuration bootstrap, including per-tenant HMAC signing secrets, in `backend/internal/app/config.go`
- [X] T013 Implement logging and metrics middleware scaffolding in `backend/internal/http/middleware/instrumentation.go`
- [X] T014 Build initial chi router wiring and templ renderer in `backend/internal/http/router.go`
- [X] T015 Implement ingestion handler with manifest validation and HMAC signature rejection in `backend/internal/ingestion/handler.go`
- [X] T016 [P] Add contract tests for signed/unsigned ingestion payloads in `tests/contract/ingestion_test.go`
- [X] T017 [P] Emit ingestion latency, signature failure, and error metrics in `backend/internal/ingestion/metrics.go`
- [X] T018 [P] Add ingestion performance benchmark for ≤500 ms acknowledgement in `tests/perf/ingestion_bench_test.go`
- [X] T019 [P] Implement seed CLI command for demo data in `backend/cmd/duplynx/seed/main.go`
- [X] T020 [P] Configure Playwright project and browsers in `tests/e2e/playwright.config.ts`
- [X] T021 [P] Add CI workflow for lint, go test, Playwright, and perf benchmarks in `.github/workflows/ci.yml`
- [X] T022 Enforce read-only database mode for GUI instances in `backend/internal/app/config.go`

---

## Phase 3: User Story 1 – Tenant & Machine Onboarding Flow (Priority: P1) 🎯 MVP

**Goal**: Visitors select a tenant, choose one of its machines, and reach the scan catalog without authentication.

**Independent Test**: Launch DupLynx, select “Sample Tenant A” then “Ares-Laptop,” and confirm the catalog appears within exactly three interactions.

### Tests

- [X] T023 [P] [US1] Write contract tests for `/tenants` and `/tenants/{tenantSlug}/machines` in `tests/contract/tenants_machines_test.go`
- [X] T024 [P] [US1] Create Playwright flow verifying three-click onboarding and header breadcrumb in `tests/e2e/onboarding.spec.ts`
- [X] T025 [P] [US1] Add integration test for tenant machine filtering logic in `tests/integration/tenancy_flow_test.go`

### Implementation

- [X] T026 [US1] Implement tenancy repository with scoped queries in `backend/internal/tenancy/repository.go`
- [X] T027 [US1] Implement `/tenants` list handler in `backend/internal/http/handlers/tenants.go`
- [X] T028 [US1] Implement tenant machine list handler in `backend/internal/http/handlers/machines.go`
- [X] T029 [US1] Build launch and machine picker templ views in `backend/internal/templ/launch.templ`
- [X] T030 [US1] Render header breadcrumb for tenant/machine in `backend/internal/templ/layout.templ`
- [X] T031 [US1] Wire tenant/machine routes and context injection in `backend/internal/http/router.go`
- [X] T032 [US1] Seed sample tenant and machine records for demo in `backend/internal/tenancy/seed.go`
- [X] T033 [P] [US1] Add unit test asserting header breadcrumb output in `tests/unit/layout_header_test.go`
- [X] T034 [US1] Log tenant selection events with structured metadata in `backend/internal/tenancy/audit.go`
- [X] T035 [P] [US1] Add integration test verifying tenant selection log emission in `tests/integration/tenancy_logging_test.go`
- [X] T036 [US1] Log machine selection events with machine metadata in `backend/internal/tenancy/audit.go`
- [X] T037 [P] [US1] Add integration test verifying machine selection log emission in `tests/integration/machine_logging_test.go`

---

## Phase 4: User Story 2 – Scan Management Board Overview (Priority: P1)

**Goal**: Users open a scan and view duplicate groups arranged by status lanes with performant rendering.

**Independent Test**: Open “Baseline Sweep 2025-10-01,” verify lane counts and card expansion with file metadata.

### Tests

- [X] T038 [P] [US2] Write contract tests for `/tenants/{tenantSlug}/scans` and `/scans/{scanId}` in `tests/contract/scans_test.go`
- [X] T039 [P] [US2] Add integration test assembling board data in `tests/integration/scan_board_test.go`
- [X] T040 [P] [US2] Add unit tests for board view model grouping in `tests/unit/board_view_test.go`

### Implementation

- [X] T041 [US2] Implement scan repository with status aggregations in `backend/internal/scans/repository.go`
- [X] T042 [US2] Implement board service composing lanes in `backend/internal/scans/service.go`
- [X] T043 [US2] Implement scan board HTTP handlers and JSON responses in `backend/internal/http/handlers/scan_board.go`
- [X] T044 [US2] Build templ components for board shell and lanes in `backend/internal/templ/board.templ`
- [X] T045 [US2] Add htmx partial handlers for lane refresh in `backend/internal/http/handlers/board_partials.go`
- [X] T046 [US2] Instrument board rendering timings and virtualized list hooks in `backend/internal/http/middleware/board_metrics.go`

---

## Phase 5: User Story 3 – Keeper Assignment & Duplicate Actions (Priority: P2)

**Goal**: Stewards assign keeper machines and trigger duplicate actions with audit logging and responsive UI feedback.

**Independent Test**: Assign “Helios-Server-02” as keeper, execute quarantine, verify audit trail and that status remains unchanged without manual reassignment.

### Tests

- [X] T047 [P] [US3] Write contract tests for `/duplicate-groups/{groupId}/keeper` and `/duplicate-groups/{groupId}/actions` in `tests/contract/actions_test.go`
- [X] T048 [P] [US3] Add integration test confirming keeper assignment and manual status retention in `tests/integration/keeper_assignment_test.go`
- [X] T049 [P] [US3] Add unit tests for action dispatcher outcomes in `tests/unit/actions_dispatcher_test.go`

### Implementation

- [X] T050 [US3] Implement action dispatcher coordinating keeper and action workflows in `backend/internal/actions/dispatcher.go`
- [X] T051 [US3] Implement action audit persistence layer in `backend/internal/actions/audit_store.go`
- [X] T052 [US3] Implement keeper and action HTTP handlers with validation in `backend/internal/http/handlers/actions.go`
- [X] T053 [US3] Enhance duplicate card templ with keeper selection UI in `backend/internal/templ/components/duplicate_card.templ`
- [X] T054 [US3] Implement htmx response fragments for action feedback in `backend/internal/http/handlers/actions_htmx.go`
- [X] T055 [US3] Record stubbed action audit entries with `stubbed=true` in `backend/internal/actions/dispatcher.go`
- [X] T056 [P] [US3] Unit test verifying stubbed audit payload in `tests/unit/actions_dispatcher_test.go`
- [X] T057 [US3] Log keeper assignment and duplicate action events in `backend/internal/actions/audit_logger.go`
- [X] T058 [P] [US3] Add integration test verifying keeper/action logging pipeline in `tests/integration/actions_logging_test.go`

---

## Phase 6: User Story 4 – Multi-Tenant Data Isolation (Priority: P3)

**Goal**: Guarantee tenant-scoped access so cross-tenant requests fail safely and are observable.

**Independent Test**: Attempt to fetch another tenant’s scan while scoped to “Sample Tenant A” and receive a tenant scope warning plus 404.

### Tests

- [X] T059 [P] [US4] Add unit tests for tenancy middleware context enforcement in `tests/unit/tenancy_middleware_test.go`
- [X] T060 [P] [US4] Add integration test for cross-tenant access rejection in `tests/integration/tenant_guard_test.go`
- [X] T061 [US4] Audit static asset routes to ensure tenant headers persist in `backend/internal/http/handlers/static.go`

### Implementation

- [X] T062 [US4] Implement tenancy scoping middleware attaching tenant context in `backend/internal/tenancy/middleware.go`
- [X] T063 [US4] Apply tenant filters across repositories in `backend/internal/tenancy/scoped_repository.go`
- [X] T064 [US4] Create tenant scope violation templ feedback in `backend/internal/templ/errors/tenant_scope.templ`

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, benchmarks, accessibility, and release readiness.

- [X] T065 Update quickstart instructions with final commands in `specs/001-create-duplynx/quickstart.md`
- [X] T066 Document SQLite writer constraints and deployment notes in `docs/duplynx-demo.md`
- [X] T067 Add board performance benchmark covering render latency in `tests/perf/board_bench_test.go`
- [X] T068 Run accessibility and contrast checks via Playwright axe audit in `tests/e2e/accessibility.spec.ts`
- [X] T069 Add final CI verification target for go test/playwright/perf combo with suite timing gates in `Makefile`
- [X] T070 Monitor CI e2e + integration runtime and fail when >8m in `.github/workflows/ci.yml`
- [X] T071 Document logging coverage for tenant/machine/action events in `docs/duplynx-demo.md`

---

## Dependencies & Execution Order

- **Phase 1 → Phase 2**: Foundational work depends on setup.
- **Phase 2 → Phases 3-6**: All user stories require Ent schemas, secure ingestion, routing, and CI tooling.
- **Phase 3 (US1)**: Enables tenant context and seeding—MVP checkpoint.
- **Phase 4 (US2)**: Depends on US1 data structures but can progress in parallel once foundational services exist.
- **Phase 5 (US3)**: Builds on duplicate group representations from US2 and ingestion logging from Phase 2.
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
# Tests
run tests/contract/tenants_machines_test.go
run tests/e2e/onboarding.spec.ts
run tests/integration/tenancy_flow_test.go

# Implementation
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

- **Total tasks**: 71
- **Per user story**:
  - US1: 15 tasks
  - US2: 9 tasks
  - US3: 12 tasks
  - US4: 6 tasks
- **Parallel opportunities**: Tests and component work within each story marked `[P]` can proceed concurrently once prerequisites land.
- **Independent test criteria**:
  - US1: Tenant and machine selection completes within three clicks and loads seeded scans.
  - US2: Scan board lanes render with correct counts and card expansion details.
  - US3: Keeper assignment succeeds, actions log audits, and statuses remain manual.
  - US4: Cross-tenant requests fail with scoped warnings and 404s.
- **Suggested MVP scope**: Deliver through Phase 3 (US1) to demonstrate tenant/machine onboarding end-to-end.

All tasks follow the required checklist format with sequential IDs, `[P]` markers where applicable, and explicit file paths.
