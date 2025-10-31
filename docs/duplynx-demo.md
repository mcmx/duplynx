# DupLynx Demo Deployment Notes

This document captures the operational expectations for the Create DupLynx demo environment. It focuses on SQLite writer constraints, process layout, and the current logging surface so you can host a stable showcase without leaking tenant data.

## Process Topology

- **Ingestion writer**: run a single `duplynx` binary with `DUPLYNX_MODE=server` (default). This instance accepts signed ingestion payloads and performs all SQLite writes. Deploy it on a host with access to the shared database file, and expose the `/ingest` endpoints behind TLS plus any gateway auth you require.
- **Read-only dashboard replicas**: additional `duplynx` binaries can serve the dashboard with `DUPLYNX_MODE=gui`. The config forces the SQLite DSN into `mode=ro`, guaranteeing these pods never take database write locks. Point them at the same database file via a shared volume (NFS, SMB, or container volume) and front them with a load balancer.
- **Static assets**: the server serves Tailwind output from `backend/web/static/`. If you deploy from source, run `tailwindcss` ahead of time or set `--embed-static=false` and offload asset delivery to a CDN.

### SQLite Guidance

| Concern | Recommendation |
| --- | --- |
| File lock contention | Keep exactly one ingestion writer. Dashboard replicas should mount the database readonly using the `gui` mode or OS-level readonly volumes. |
| Durability | Place the database on resilient storage (e.g., SSD-backed persistent volume). Add periodic snapshotting for disaster recovery. |
| Busy timeouts | The Go configuration sets `_busy_timeout=5000` for basic write contention handling. Increase `DUPLYNX_SQLITE_BUSY_TIMEOUT` if ingestion payloads ever spike latency. |
| Migration hygiene | Ent migrations should run on the ingestion writer instance before scaling out additional GUI pods. |

## Required Configuration

| Variable | Purpose | Example |
| --- | --- | --- |
| `DUPLYNX_DB_FILE` | Absolute path to the shared SQLite database. | `/var/lib/duplynx/duplynx.db` |
| `DUPLYNX_ADDR` | HTTP bind address. | `:8080` |
| `DUPLYNX_EMBED_STATIC` | Toggle embedded static assets. | `true` |
| `DUPLYNX_MODE` | `server` (read/write) or `gui` (read-only dashboard). | `gui` |
| `DUPLYNX_TENANT_SECRETS` | Comma-delimited `tenantSlug:hexSecret` pairs for HMAC validation. | `sample-tenant-a:deadbeef` |

## Logging Coverage

DupLynx currently emits the following structured audit entries:

| Event | Package | Trigger |
| --- | --- | --- |
| `tenant_selection` | `internal/tenancy.AuditLogger` | When a user picks a tenant from the launch screen. |
| `machine_selection` | `internal/tenancy.AuditLogger` | When the UI records the active machine context. |
| `assign_keeper` | `internal/actions.Dispatcher` → `AuditLogger` | When a keeper machine is set on a duplicate group. |
| `delete_copies` / `create_hardlinks` / `quarantine` | `internal/actions.Dispatcher` → `AuditLogger` | When an action is triggered from the duplicate group card; entries include the payload and are marked `stubbed=true` in the current phase. |

Forward these logs to your observability stack (stdout collectors, Loki, etc.) to reconstruct user flows and prove tenant isolation. When running multiple GUI replicas, ensure each pod streams logs centrally so audit trails remain contiguous.

## Demo Checklist

1. Run the seed command to populate tenants, machines, scans, and duplicate groups:  
   ```bash
   go run ./cmd/duplynx seed --db-file var/duplynx.db
   ```
2. Start the ingestion writer:  
   ```bash
   DUPLYNX_MODE=server DUPLYNX_TENANT_SECRETS=sample-tenant-a:deadbeef \
   go run ./cmd/duplynx serve --db-file var/duplynx.db --addr :8080
   ```
3. (Optional) Launch a read-only dashboard replica:  
   ```bash
   DUPLYNX_MODE=gui go run ./cmd/duplynx serve --db-file var/duplynx.db --addr :8081
   ```
4. Configure your load balancer to direct ingestion traffic to the writer and dashboard traffic to GUI replicas.
5. Tail logs to verify tenant selection and keeper/action audit events while demoing.
