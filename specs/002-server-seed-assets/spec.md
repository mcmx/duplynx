# Feature Specification: DupLynx Server Runtime

**Feature Branch**: `[002-server-seed-assets]`  
**Created**: 2025-11-01  
**Status**: Draft  
**Input**: User description: "Add runnable DupLynx server entrypoint, seed command, and static asset pipeline. Replace the stubbed `backend/cmd/duplynx/main.go` so `duplynx serve`: - loads config (flags/env), - opens/initializes SQLite via Ent migrations, - seeds demo data (tenants, machines, scans, duplicate groups), - mounts the existing Chi router and templ views, - serves Tailwind assets (embed or filesystem) on 0.0.0.0:8080. Introduce `duplynx seed` to create/update the demo SQLite file so users can run `go run ./cmd/duplynx seed --db-file var/duplynx.db`. Ensure the launch page, machine picker, scan list, and board render from the seeded database instead of in-memory slices. Keep tenant isolation middleware, log audit events, and update quickstart/documentation/tests accordingly (Playwright onboarding, contract/integration suites, new smoke test that hits `/`). Goal: after `make tidy && make ci`, users can run `go run ./cmd/duplynx serve --db-file var/duplynx.db` and view the dashboard at `http://localhost:8080` with seeded data."

## Clarifications

### Session 2025-11-01
- Q: What is the expected source for static assets during runtime? → A: Serve prebuilt assets from disk.
- Q: What audit trail should be generated when the seed command runs? → A: Event log entry per seed run.
- Q: How should the seed command handle existing records in the demo database? → A: Replace the entire demo dataset each run.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Launch demo environment (Priority: P1)

An evaluator wants to spin up DupLynx locally, run a single command, and immediately explore the dashboard with meaningful demo data.

**Why this priority**: Without a turnkey demo server, prospects and internal stakeholders cannot experience the product value, blocking adoption.

**Independent Test**: Can be fully tested by following the documented quickstart to start the server and confirming the dashboard renders seeded content end to end.

**Acceptance Scenarios**:

1. **Given** a clean workstation with prerequisites installed, **When** the evaluator runs the documented serve command with the provided database path, **Then** the dashboard loads in a browser with seeded tenants, machines, and duplicate insights.
2. **Given** an invalid configuration value, **When** the evaluator attempts to start the server, **Then** the process fails fast with a human-readable message and no partial state is left behind.

---

### User Story 2 - Refresh demo dataset (Priority: P2)

An enablement lead wants to regenerate the demo dataset so the sample tenants and scans stay current without manually editing database files.

**Why this priority**: Reusable data seeding keeps demos consistent and prevents drift between documentation, tests, and what users see.

**Independent Test**: Can be fully tested by running the seed command on an existing demo file and verifying all expected tenants and records exist without duplicates or missing links.

**Acceptance Scenarios**:

1. **Given** an empty or outdated demo database file, **When** the enablement lead runs the seed command, **Then** the file contains the full set of canonical tenants, machines, scans, duplicate groups, and audit history ready for serving.

---

### User Story 3 - CI confidence in runtime (Priority: P3)

A release engineer wants automated checks that ensure the seeded server responds and key journeys render before promoting changes.

**Why this priority**: Automated smoke coverage prevents regressions that would otherwise break demos or onboarding flows.

**Independent Test**: Can be fully tested by running the updated CI job that seeds the database, starts the server, and confirms the primary routes respond as expected.

**Acceptance Scenarios**:

1. **Given** a continuous integration run, **When** the pipeline executes the new smoke test suite, **Then** it seeds data, hits the root route, and fails the build if the page does not render expected elements.

### Edge Cases
- Startup fails fast when the asset directory is missing or unreadable, emitting a remediation hint to rerun `npm run build:tailwind`.
- Seeding overwrites the demo database, so manual edits are lost after reseed; operators are prompted to back up custom fixtures before running the command.
- Attempts to serve with a missing or unreadable configuration file terminate immediately with the list of required flags/environment variables and their defaults.
- Seeding into an existing database that already contains partial demo data first drops demo tables to prevent duplicate tenants, machines, or scans.
- Server start while the database file is locked by another process retries for up to 5 seconds before exiting with guidance to stop or wait for the competing process.
- Requests for static assets when the asset bundle is unavailable or outdated respond with HTTP 503 and log a warning instructing operators to rebuild assets.
- Multi-tenant view requests when the tenant identifier is absent or malformed return HTTP 400 with an error template rendered and the incident captured in audit logs.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST allow operators to launch the DupLynx server through a documented command that reads required configuration from flags and environment variables.
- **FR-002**: System MUST validate configuration at startup and surface actionable errors without mutating data when validation fails.
- **FR-003**: System MUST ensure the local demo database schema is up to date before serving user traffic.
- **FR-004**: System MUST recreate the demo database deterministically by clearing prior contents and seeding canonical tenants, machines, scans, duplicate groups, and associated audit records via a dedicated command.
- **FR-005**: System MUST make the seeded data visible across the launch page, tenant picker, scan list, and board views with tenant isolation enforced.
- **FR-006**: System MUST expose static assets necessary for the dashboard styling and interactions during runtime by serving a precompiled asset directory from disk and failing startup if assets are unavailable.
- **FR-007**: System MUST record audit events for each seeding run and server lifecycle action, including timestamp, actor, and outcome, to preserve traceability.
- **FR-008**: Documentation MUST describe prerequisites and step-by-step instructions to seed data, start the server, and verify the dashboard.
- **FR-009**: Automated test suites MUST cover the seeding workflow, core routes, and asset availability within continuous integration.
- **FR-010**: System MUST provide a smoke test that exercises the root route using the seeded database and fails if expected interface elements are missing.

### Runtime Configuration & Operations

- **Config Surface**: `duplynx serve` and `duplynx seed` MUST expose flags `--db-file` (default `var/duplynx.db`), `--assets-dir` (default `backend/web/dist`), `--addr` (serve only, default `0.0.0.0:8080`), and `--log-level` (default `info`). Each flag MUST map to environment overrides (`DUPLYNX_DB_FILE`, `DUPLYNX_ASSETS_DIR`, `DUPLYNX_ADDR`, `DUPLYNX_LOG_LEVEL`).
- **Startup Validation**: On launch, the CLI MUST validate configuration (paths exist, values parse) and exit with non-zero status and actionable messaging when validation fails. Schema migration errors MUST abort before the HTTP server binds. Missing or unreadable asset directories MUST halt startup with guidance to rebuild Tailwind assets.
- **Serve Before Seed**: When `duplynx serve` runs against an empty database, it MUST apply migrations, log a warning that no demo data exists, and instruct the operator to run `duplynx seed`. The command MUST refuse to fabricate partial demo data.
- **Deterministic Seed Flow**: `duplynx seed` MUST obtain an exclusive lock on the database file, drop existing demo tables, apply migrations, insert canonical tenants → machines → scans → duplicate groups in deterministic order, and emit audit events summarizing operations. On partial failure it MUST roll back outstanding transactions, clear half-populated tables, and exit with guidance to rerun the command.
- **Audit Logging**: Seed and serve events MUST emit structured stdout logs describing actor (CLI user or CI identifier), ISO8601 timestamps, duration, outcome (`success`, `validation_error`, `io_error`), and metadata (db path, git revision) so centralized pipelines can persist them. Log volume MUST remain under 1 MB per run on the standard workstation.
- **Telemetry & Timings**: The smoke test and quickstart timing script MUST record seed and serve durations, emit results to CI job logs, and write measurements to `var/duplynx_bench.json` for local inspection. Alerts MUST trigger when budgets exceed thresholds defined in Success Criteria.

### Key Entities *(include if feature involves data)*

- **Tenant**: Represents an organization using DupLynx; includes name, demo status, and links to machines and scans.
- **Machine**: Represents a host or endpoint within a tenant; tracks identifiers, display labels, and associated scan runs.
- **Scan**: Represents a deduplication analysis performed on a machine; includes timestamp, status, and linked duplicate groups.
- **Duplicate Group**: Represents a set of related findings surfaced by a scan; includes severity, suggestion summary, and owning tenant context.
- **Audit Event**: Represents logged actions such as seeding, server start, and configuration changes; includes actor, timestamp, outcome, and description or metadata captured via structured stdout events.

## Assumptions

- Evaluators run the commands on a workstation where required runtime dependencies are already installed.
- The demo database file is stored on local disk and can be overwritten during seeding without impacting production data.
- Existing routing, templating, and tenant isolation components remain available for integration with the new entrypoint.
- The DupLynx CLI is delivered as a single Cobra-based binary with `serve` and `seed` subcommands accessed via `go run ./cmd/duplynx <command>`.
- Tailwind assets are built externally and supplied via `backend/web/dist`; embedded or fallback asset strategies are intentionally out of scope for this feature.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: On the standard DupLynx developer workstation (Apple M2, 16 GB RAM, SSD), evaluators complete the documented seed and serve workflow and view the dashboard with demo data in under 5 minutes, as recorded by the automated quickstart timing script located at `scripts/measure_quickstart.sh`.
- **SC-002**: The seeding command populates the full demo dataset in under 60 seconds on the standard workstation, with each run logging its measured duration into the audit event stream.
- **SC-003**: CI pipelines complete the new automated smoke test suite with zero manual intervention, capturing seed and serve durations in the job logs, and block releases on failure.
- **SC-004**: Onboarding support requests related to starting DupLynx locally decrease by at least 50% within one release cycle compared to the baseline volume captured in the pre-release support log, with progress noted in the release summary.

## Demo Dataset Composition

- **Tenants**: Three canonical tenants (`sample-tenant-a`, `sample-tenant-b`, `sample-tenant-c`) each flagged as demo and paired with consistent branding colors for UI cards.
- **Machines**: Four machines per tenant with deterministic IDs, hostnames, and display names covering laptop, server, and appliance personas. Machine display names MUST sort alphabetically for picker stability.
- **Scans**: Each tenant receives two canonical scans with IDs formatted `baseline-sweep-YYYY-MM-DD` and `followup-sweep-YYYY-MM-DD`, statuses set to `succeeded`, and timestamps spaced 24 hours apart.
- **Duplicate Groups**: Each scan contains at least three duplicate groups covering `low`, `medium`, and `high` severity with stable example paths leveraged by the board render tests.
- **Audit Events**: Seed runs append `seed_start`, `seed_complete`, and `seed_error` entries; serve runs append `serve_start` and `serve_stop`. All events include duration metrics and metadata for db path and git revision.

## Smoke Test Assertions

- Root route MUST render tenant summary cards displaying tenant name, machine count, and scan count.
- Tenant picker endpoint MUST respond with JSON containing expected tenant slugs and names.
- Scan list view MUST include duplicate group counts per scan and at least one templ-based CTA button for launching the board view.
- Smoke test MUST assert presence of the Tailwind-generated stylesheet link and fail the run if the asset link is missing.
