# Feature Specification: Create DupLynx

**Feature Branch**: `[001-create-duplynx]`  
**Created**: 2025-10-27  
**Status**: Draft  
**Input**: User description: "/prompts:speckit.specify Develop DupLynx, a distributed duplicate file detector platform. It should allow users to scan filesystems on multiple machines, upload checksums, and other file metadata to a central server, identify duplicate files across machines based on cryptographic hashes, and manage them by deleting copies, creating hardlinks, or quarantining. The system must be multi-tenant, supporting multiple isolated tenants (e.g., organizations or teams) where each tenant has segregated data, machines, scans, and access controls to ensure privacy and scalability across shared infrastructure. In this initial phase for this feature, let's call it "Create DupLynx," let's have multiple machines but the machines will be declared ahead of time, predefined per tenant. I want five machines in two different categories, one personal laptop and four server instances, assigned to a sample tenant. Let's create three different sample scans. Let's have the standard duplicate management views for the status of each group, such as "Review," "Action Needed," "Resolved," and "Archived." There will be no login for this application as this is just the very first testing thing to ensure that our basic features are set up. The system must have a web dashboard as the primary interface for all interactions. You should be able to, from that group card, assign one of the valid machines as the "keeper" for the master copy. When you first launch DupLynx, it's going to give you a list of the sample tenants to pick from, followed by a list of the five machines within the selected tenant to pick from. There will be no password required. When you click on a machine, you go into the main view, which displays the list of scans for that tenant. When you click on a scan, you open the management board for that scan in the web dashboard."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Tenant & Machine Onboarding Flow (Priority: P1)

First-time visitors see the launch screen, pick a sample tenant, choose from that tenant’s predefined machines, and land on the scan catalog for validation or demo use without entering credentials.

**Why this priority**: Without tenant and machine selection, no other functionality is reachable; this is the minimal path to interact with DupLynx.

**Independent Test**: Automated UI smoke test navigates from launch to scan list by selecting tenant “Sample Tenant A” and machine “Ares-Laptop,” asserting machine metadata and scan list render successfully.

**Constitution Alignment**:
- UX: Reuse DupLynx global layout, typography tokens, and card list component for tenant/machine pickers.
- Testing: Cypress (or equivalent) e2e path for launch → tenant pick → machine pick, plus unit tests for tenant store segregation.
- Performance: Initial tenant/machine load must stay under 1s for ≤10 tenants with ≤25 machines cached locally.

**Acceptance Scenarios**:

1. **Given** DupLynx launches fresh, **When** the user selects “Sample Tenant A,” **Then** only that tenant’s predefined machines (1 laptop, 4 servers) appear with category labels.
2. **Given** the machine list is visible, **When** the user selects “Ares-Laptop,” **Then** the scan list view loads showing three sample scans tied to that tenant.

---

### User Story 2 - Scan Management Board Overview (Priority: P1)

Operators need to open a scan and review duplicate groups organized by lifecycle status (“Review,” “Action Needed,” “Resolved,” “Archived”) within the web dashboard.

**Why this priority**: Core value of DupLynx is surfacing duplicate groups per scan; without the board, no insight or triage can occur.

**Independent Test**: Component integration test loads scan “Baseline Sweep 2025-10-01” and verifies each status column renders with seeded duplicate cards and counts.

**Constitution Alignment**:
- UX: Apply DupLynx kanban board pattern with accessible column headers and status badges from the design system.
- Testing: Jest/component tests validating column rendering, plus mocked API contract test ensuring status-filtered data separation.
- Performance: Board rendering should virtualize duplicate cards to keep column paint time ≤200 ms for up to 200 groups.

**Acceptance Scenarios**:

1. **Given** the scan “Baseline Sweep 2025-10-01” is open, **When** the dashboard renders, **Then** each status lane displays its seeded duplicate groups with counts and summary metadata.
2. **Given** a duplicate group card is visible, **When** a user expands it, **Then** file instances, sizes, and hosting machines are listed with accessible labels.

---

### User Story 3 - Keeper Assignment & Duplicate Actions (Priority: P2)

Stewards adjust each duplicate group by designating a keeper machine for the master copy and triggering duplicate management actions (delete extra copies, create hardlinks, quarantine) from the group card.

**Why this priority**: Assigning a keeper establishes ownership for deduplication workflows, and action controls validate the end-to-end management concept.

**Independent Test**: Automated test mocks action services, selects a duplicate group, assigns keeper “Helios-Server-02,” and asserts state changes and audit log entries.

**Constitution Alignment**:
- UX: Use standard action menu pattern with confirmation modals and consistent iconography from the token library.
- Testing: Integration tests covering keeper assignment state transitions and action dispatch; contract tests for action API payloads.
- Performance: Action execution queue must respond within 500 ms for acknowledgement even when remote operations are stubbed.

**Acceptance Scenarios**:

1. **Given** a duplicate group without a keeper, **When** the user assigns “Helios-Server-02” as keeper, **Then** the card indicates the keeper, persists the assignment, and audits the change.
2. **Given** a duplicate group with non-keeper copies, **When** the user chooses “Quarantine,” **Then** the system marks the relevant file instances as quarantined and moves the card to “Action Needed.”

---

### User Story 4 - Multi-Tenant Data Isolation (Priority: P3)

Platform admins must ensure each tenant’s data, machines, scans, and duplicate actions remain isolated even when running on shared infrastructure during demos.

**Why this priority**: Isolation protects privacy and is foundational for scaling beyond the sample tenant.

**Independent Test**: API contract test attempts to request scans for “Sample Tenant A” while authenticated to “Sample Tenant B” context and expects a 404.

**Constitution Alignment**:
- UX: Tenant name persistently displayed in shell header informing users of current context.
- Testing: Unit tests for tenancy middleware plus multi-tenant integration tests covering machine/scan segregation.
- Performance: Tenant scoping queries must keep database response ≤150 ms for basic list endpoints with indexes on tenant foreign keys.

**Acceptance Scenarios**:

1. **Given** the user has selected “Sample Tenant A,” **When** they attempt to navigate directly to a “Sample Tenant B” scan URL, **Then** the app rejects the navigation with a tenant scope warning.
2. **Given** background sync jobs run, **When** duplicate metadata is processed, **Then** no tenant’s data appears in another tenant’s board or API payloads.

### Edge Cases

- What happens when a scan contains zero duplicate groups? → Display an empty state in each status lane with guidance to rerun or adjust scan parameters.
- How does system handle a keeper assignment for a machine not in the tenant? → Prevent selection, surface validation error, and log the access attempt for audit.
- How does the UI respond if an action endpoint is temporarily unavailable? → Show non-blocking toast with retry option and keep group in previous status.
- What occurs when sample seed data fails to load at launch? → Present fallback message with “Retry Seeding” action and diagnostic logging.

## Requirements *(mandatory)*

### Quality Guardrails

- **QG-001**: Code MUST meet project lint/static analysis rules and include public API documentation updates.
- **QG-002**: Automated tests (unit, integration, contract as applicable) MUST cover the new behavior and run in CI as merge blockers.
- **QG-003**: UX changes MUST reference shared components/tokens and include accessibility acceptance criteria.
- **QG-004**: Performance budgets or SLAs impacted by this feature MUST be defined along with measurement/monitoring hooks.

### Functional Requirements

- **FR-001**: System MUST seed at least one sample tenant (“Sample Tenant A”) with five predefined machines (1 personal laptop, 4 server instances) labeled by category at launch.
- **FR-002**: System MUST present a tenant selection view on first load showing all sample tenants with descriptive metadata and no authentication.
- **FR-003**: Selecting a tenant MUST show only that tenant’s machines, and choosing a machine MUST navigate to the tenant’s scan catalog with three seeded scans.
- **FR-004**: Each scan entry MUST display metadata (name, executed machine, timestamp, duplicate count snapshot) and link to the management board.
- **FR-005**: Management board MUST organize duplicate groups into the statuses “Review,” “Action Needed,” “Resolved,” and “Archived” with counts per lane.
- **FR-006**: Duplicate group cards MUST allow assigning any machine from the tenant’s roster as the keeper for the master copy with audit trail.
- **FR-007**: Duplicate group cards MUST expose actions to delete redundant files, create hardlinks back to the keeper copy, and quarantine suspicious files; initial implementation MAY stub side effects but MUST persist intended state.
- **FR-008**: All data access APIs MUST enforce tenant scoping, preventing machines, scans, or duplicate groups from leaking across tenants.
- **FR-009**: Web dashboard MUST surface current tenant context and machine selection in the global header for clarity.
- **FR-010**: System MUST log key user actions (tenant selection, machine selection, keeper assignment, duplicate actions) for future monitoring hooks.

### Key Entities *(include if feature involves data)*

- **Tenant**: Represents an organization; attributes include `id`, `name`, descriptive metadata; relates to many `Machine`, `Scan`, and `DuplicateGroup` records.
- **Machine**: Represents a predefined host; attributes include `id`, `tenant_id`, `name`, `category`, `role`, `last_scan_at`.
- **Scan**: Represents a deduplication run; attributes include `id`, `tenant_id`, `initiated_machine_id`, `name`, `started_at`, `completed_at`, `summary`.
- **DuplicateGroup**: Represents a set of identical file hashes; attributes include `id`, `tenant_id`, `scan_id`, `status`, `keeper_machine_id`, `hash`, `total_size`.
- **FileInstance**: Represents a specific file copy; attributes include `id`, `duplicate_group_id`, `machine_id`, `path`, `size`, `last_seen_at`, `quarantined`.
- **ActionAudit**: Records user-driven actions; attributes include `id`, `tenant_id`, `duplicate_group_id`, `action_type`, `payload`, `performed_by`, `performed_at`.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Demo users can navigate from launch to a scan board and identify a duplicate group within 3 clicks using the standard dashboard layout.
- **SC-002**: Linting/static analysis passes with zero new warnings, and minimum 90% coverage is maintained for modules handling tenant scoping and duplicate board logic.
- **SC-003**: Automated UI/integration suite exercising tenant selection, scan navigation, and keeper assignment completes in ≤8 minutes on CI hardware.
- **SC-004**: Seeded board renders in ≤1 s and action acknowledgements respond in ≤500 ms for up to 200 duplicate groups per scan during demos.
