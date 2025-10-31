# Data Model — Create DupLynx

## Overview

DupLynx persists multi-tenant deduplication metadata in SQLite via Ent. All tables include `id`, `created_at`, and `updated_at` fields managed by Ent mixins. Foreign key constraints enforce tenant isolation.

## Entities

### Tenant
- **Table**: `tenants`
- **Fields**:
  - `id` (UUID, PK)
  - `slug` (string, unique, lower-kebab)
  - `name` (string, required)
  - `description` (text, optional)
  - `primary_contact` (string, optional for future demo extensions)
- **Relationships**:
  - `machines` (1:N Machine)
  - `scans` (1:N Scan)
  - `action_audits` (1:N ActionAudit)
- **Validation/Rules**:
  - Slug must be unique and URL-safe.
  - Soft delete not required for Phase 0.

### Machine
- **Table**: `machines`
- **Fields**:
  - `id` (UUID, PK)
  - `tenant_id` (UUID, FK → tenants.id, cascade delete)
  - `name` (string, required)
  - `category` (enum: `personal_laptop`, `server`)
  - `hostname` (string, optional; seeded demo metadata)
  - `role` (string, e.g., “ingest”, optional)
  - `last_scan_at` (datetime, nullable)
- **Relationships**:
  - `tenant` (N:1 Tenant)
  - `initiated_scans` (1:N Scan via `initiated_machine_id`)
  - `keeper_groups` (1:N DuplicateGroup via `keeper_machine_id`)
  - `file_instances` (1:N FileInstance)
- **Validation/Rules**:
  - Category + name combination unique per tenant.
  - `last_scan_at` auto-updates when a scan completes on the machine.

### Scan
- **Table**: `scans`
- **Fields**:
  - `id` (UUID, PK)
  - `tenant_id` (UUID, FK → tenants.id)
  - `initiated_machine_id` (UUID, FK → machines.id)
  - `name` (string, required, unique per tenant)
  - `description` (text, optional)
  - `started_at` (datetime, required)
  - `completed_at` (datetime, nullable)
  - `duplicate_group_count` (int, denormalized summary)
- **Relationships**:
  - `tenant` (N:1 Tenant)
  - `initiated_machine` (N:1 Machine)
  - `duplicate_groups` (1:N DuplicateGroup)
- **Validation/Rules**:
  - `completed_at` must be ≥ `started_at` when present.
  - Demo seed ensures exactly three scans per tenant for Phase 0.

### DuplicateGroup
- **Table**: `duplicate_groups`
- **Fields**:
  - `id` (UUID, PK)
  - `tenant_id` (UUID, FK → tenants.id)
  - `scan_id` (UUID, FK → scans.id, cascade delete)
  - `hash` (string, required, SHA-256 hex)
  - `status` (enum: `review`, `action_needed`, `resolved`, `archived`)
  - `keeper_machine_id` (UUID, FK → machines.id, nullable)
  - `total_size_bytes` (bigint, required)
  - `file_count` (int, required)
- **Relationships**:
  - `scan` (N:1 Scan)
  - `tenant` (N:1 Tenant)
  - `keeper_machine` (N:1 Machine, optional)
  - `file_instances` (1:N FileInstance)
  - `action_audits` (1:N ActionAudit)
- **Validation/Rules**:
  - Status defaults to `review`.
  - Keeper must belong to same tenant; enforced via middleware and DB constraint.
  - `file_count` ≥ 2 for duplicates.

#### Status Transitions
```
review → action_needed (when detection requires manual triage)
review → resolved (keeper assignment with automated cleanup)
action_needed → resolved (tasks completed)
resolved → archived (historical record)
action_needed → archived (manual override)
```
All other transitions invalid and rejected at service layer.

### FileInstance
- **Table**: `file_instances`
- **Fields**:
  - `id` (UUID, PK)
  - `duplicate_group_id` (UUID, FK → duplicate_groups.id, cascade delete)
  - `machine_id` (UUID, FK → machines.id)
  - `path` (string, required)
  - `size_bytes` (bigint, required)
  - `checksum` (string, required, SHA-256 hex)
  - `last_seen_at` (datetime, required)
  - `quarantined` (bool, default false)
- **Relationships**:
  - `duplicate_group` (N:1 DuplicateGroup)
  - `machine` (N:1 Machine)
- **Validation/Rules**:
  - `checksum` must match group `hash`.
  - `quarantined` update triggers audit entry.

### ActionAudit
- **Table**: `action_audits`
- **Fields**:
  - `id` (UUID, PK)
  - `tenant_id` (UUID, FK → tenants.id)
  - `duplicate_group_id` (UUID, FK → duplicate_groups.id)
  - `actor` (string, default `system` for demo)
  - `action_type` (enum: `assign_keeper`, `delete_copies`, `create_hardlinks`, `quarantine`, `retry`, `note`)
  - `payload` (JSON, optional)
  - `performed_at` (datetime, default now)
- **Relationships**:
  - `duplicate_group` (N:1 DuplicateGroup)
  - `tenant` (N:1 Tenant)
- **Validation/Rules**:
  - Payload schema validated per action type (e.g., list of machine IDs).
  - All writes logged to console for monitoring budget.

## Indexing Strategy
- `machines`: composite index `(tenant_id, category, name)`
- `scans`: unique `(tenant_id, name)`, index `(initiated_machine_id)`
- `duplicate_groups`: composite index `(tenant_id, status)`, unique `(scan_id, hash)`
- `file_instances`: composite index `(duplicate_group_id, machine_id)`
- `action_audits`: index `(tenant_id, performed_at DESC)`

## Seed Data Snapshot
- Tenant: “Sample Tenant A” (slug `sample-tenant-a`)
- Machines: 5 total (1 `personal_laptop`, 4 `server`) with canonical names.
- Scans: “Baseline Sweep 2025-10-01”, “Media Audit 2025-10-10”, “Archive Sync 2025-10-20”.
- Duplicate groups: ~50 per scan distributed across statuses; keeper preassigned for a subset to demonstrate state variety.
