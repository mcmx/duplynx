# Data Model: DupLynx Server Runtime

## Tenant
- **Purpose**: Represents an organization within the demo dataset.
- **Key Fields**:
  - `id` (UUID) — primary identifier.
  - `slug` (string) — URL-safe unique key used in routing.
  - `name` (string) — display label.
  - `demo` (bool) — marks demo tenants for filtering.
  - `created_at` / `updated_at` (timestamps).
- **Relationships**:
  - Has many `Machine` records.
  - Has many `Scan` records through machines.
  - Has many `DuplicateGroup` records scoped by tenant.
- **Validation Rules**:
  - `slug` unique across all tenants (case-insensitive).
  - `name` must be non-empty and <= 100 characters.
- **Lifecycle**:
  - Created via seed script only; deletions handled by full reseed.

## Machine
- **Purpose**: Represents a host/end-point analyzed by DupLynx.
- **Key Fields**:
  - `id` (UUID) — primary identifier.
  - `tenant_id` (UUID) — foreign key to `Tenant`.
  - `hostname` (string) — canonical machine identifier.
  - `display_name` (string) — friendly label for UI.
  - `last_scan_at` (timestamp) — most recent scan completion.
- **Relationships**:
  - Belongs to `Tenant`.
  - Has many `Scan` records.
- **Validation Rules**:
  - `hostname` unique within a tenant.
  - `display_name` required for UI drop-downs.
- **Lifecycle**:
  - Created during seeding; reseed replaces entire set.

## Scan
- **Purpose**: Captures a deduplication analysis run for a machine.
- **Key Fields**:
  - `id` (UUID).
  - `tenant_id` (UUID) — duplicates tenant context for quick filtering.
  - `machine_id` (UUID) — owning machine.
  - `started_at` / `completed_at` (timestamps).
  - `status` (enum: `pending`, `running`, `succeeded`, `failed`).
  - `finding_count` (int) — number of duplicate groups uncovered.
- **Relationships**:
  - Belongs to `Machine` and `Tenant`.
  - Has many `DuplicateGroup` records.
- **Validation Rules**:
  - `completed_at` must be >= `started_at` when status finalized.
  - `status` must be terminal (`succeeded`/`failed`) before seeding duplicate groups.
- **Lifecycle**:
  - Seed script inserts canonical succeeded scans per machine.

## DuplicateGroup
- **Purpose**: Represents a cluster of related duplicate findings.
- **Key Fields**:
  - `id` (UUID).
  - `tenant_id` (UUID).
  - `scan_id` (UUID).
  - `title` (string) — summary displayed on board.
  - `severity` (enum: `low`, `medium`, `high`).
  - `example_path` (string) — sample duplicate reference.
  - `suggested_action` (string) — remediation guidance.
- **Relationships**:
  - Belongs to `Scan` and `Tenant`.
- **Validation Rules**:
  - `title` required, <= 140 characters.
  - `severity` must be within enumerated set.
- **Lifecycle**:
  - Created by seed script in deterministic ordering for stable board layouts.

## AuditEvent
- **Purpose**: Logs seeding and server lifecycle events.
- **Key Fields**:
  - `id` (UUID).
  - `timestamp` (timestamp).
  - `actor` (string) — CLI user or CI identifier.
  - `action` (enum: `seed_start`, `seed_complete`, `serve_start`, `serve_stop`, `seed_error`, `serve_error`).
  - `outcome` (string) — success/failure detail or error summary.
  - `metadata` (JSON) — optional context (db path, version info).
- **Relationships**:
  - Not linked to tenants; global audit table.
- **Validation Rules**:
  - `action` required and enumerated.
  - `outcome` required when `action` ends with `_error`.
- **Lifecycle**:
  - Entries appended by CLI commands; cleared only on full reseed when database recreated.

## Derived Views
- **Machine Picker View**: joins Tenant → Machine for dropdown; requires machines sorted by `display_name`.
- **Board View**: joins Scan → DuplicateGroup filtered by tenant; boards sorted by severity desc then title.
- **Launch Dashboard**: aggregates counts of scans and groups per tenant for summary cards.
