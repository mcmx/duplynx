# DupLynx Demo Deployment Notes

This document captures the operational expectations for the Create DupLynx demo environment. It focuses on SQLite writer constraints, process layout, and the current logging surface so you can host a stable showcase without leaking tenant data.

## Process Topology

- **Ingestion writer**: run a single `duplynx` binary with `DUPLYNX_MODE=server` (default). This instance accepts signed ingestion payloads and performs all SQLite writes. Deploy it on a host with access to the shared database file, and expose the `/ingest` endpoints behind TLS plus any gateway auth you require.
- **Read-only dashboard replicas**: additional `duplynx` binaries can serve the dashboard with `DUPLYNX_MODE=gui`. The config forces the SQLite DSN into `mode=ro`, guaranteeing these pods never take database write locks. Point them at the same database file via a shared volume (NFS, SMB, or container volume) and front them with a load balancer.
- **Static assets**: the server serves the Tailwind bundle from `backend/web/dist/`. Run `npm run build:tailwind` ahead of time and mount the resulting directory read-only; there is no embedded or CDN fallback in this phase.

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
| `DUPLYNX_ASSETS_DIR` | Directory containing the built Tailwind bundle (`tailwind.css`). | `/var/lib/duplynx/assets` |
| `DUPLYNX_ADDR` | HTTP bind address. | `0.0.0.0:8080` |
| `DUPLYNX_LOG_LEVEL` | CLI log verbosity (`debug`, `info`, `warn`, `error`). | `info` |

## Logging Coverage

DupLynx currently emits the following structured audit entries:

| Event | Package | Trigger |
| --- | --- | --- |
| `tenant_selection` | `internal/tenancy.AuditLogger` | When a user picks a tenant from the launch screen. |
| `machine_selection` | `internal/tenancy.AuditLogger` | When the UI records the active machine context. |
| `assign_keeper` | `internal/actions.Dispatcher` → `AuditLogger` | When a keeper machine is set on a duplicate group. |
| `delete_copies` / `create_hardlinks` / `quarantine` | `internal/actions.Dispatcher` → `AuditLogger` | When an action is triggered from the duplicate group card; entries include the payload and are marked `stubbed=true` in the current phase. |

Forward these logs to your observability stack (stdout collectors, Loki, etc.) to reconstruct user flows and prove tenant isolation. When running multiple GUI replicas, ensure each pod streams logs centrally so audit trails remain contiguous.

## Seeding Workflow

The `duplynx seed` command rebuilds the demo database with a deterministic dataset of tenants, machines, scans, duplicate groups, file instances, and historical duplicate actions.

```bash
cd backend
go run ./cmd/duplynx seed \
  --db-file ../var/duplynx.db \
  --assets-dir ./web/dist
```

- The command drops existing demo tables, reapplies Ent migrations, and writes the canonical fixtures.
- Every execution emits `seed_start` and `seed_stop` audit events with actor metadata so CI and onboarding scripts can verify success.
- The operation is idempotent: rerunning the command produces identical records, making it safe for CI and local refreshes.
- `--db-file` and `--assets-dir` accept either absolute paths or paths relative to the repository root. Environment overrides follow the same flag names (e.g., `DUPLYNX_DB_FILE`).

If the Tailwind bundle is missing, rebuild it before seeding:

```bash
npm install
npm run build:tailwind
```

## Demo Checklist

1. Build static assets (required once per change):  
   ```bash
   npm install
   npm run build:tailwind
   ```
2. Seed the canonical dataset (idempotent):  
   ```bash
   cd backend
   go run ./cmd/duplynx seed \
     --db-file ../var/duplynx.db \
     --assets-dir ./web/dist
   ```
3. Start the demo server for evaluators:  
   ```bash
   cd backend
   go run ./cmd/duplynx serve \
     --db-file ../var/duplynx.db \
     --assets-dir ./web/dist \
     --addr 0.0.0.0:8080
   ```
4. Validate the end-to-end flow with the automated smoke test (fails fast on missing Tailwind bundle or slow start-ups):  
   ```bash
   make smoke-demo
   ```
5. (Optional) Launch additional read-only replicas by repeating the serve command on different ports and pointing them at the same database file in read-only mode.
6. Monitor stdout for audit entries (`seed_*`, `serve_*`, tenant selection, keeper assignments) to verify healthy flows during demos.

## CI Integration

The GitHub Actions workflow runs `make ci` followed by `make smoke-demo` on every push and pull request. The smoke target enforces a five-minute ceiling for the seed/serve cycle and asserts that the rendered dashboard links the Tailwind bundle, preventing regressions that would break evaluator onboarding.

## Quickstart Verification Log

- Verified on 2025-11-04 using an Apple M2 (16 GB RAM, SSD) by following the quickstart: Tailwind rebuild, `duplynx seed`, `duplynx serve`, and `make smoke-demo`.
- `scripts/measure_quickstart.sh --db-file var/duplynx.db --assets-dir backend/web/dist` reported `seed_ms=14`, `serve_ms=1018`, `total_ms=1032`; results captured in `var/duplynx_bench.json`.
- The timing script now auto-selects a free localhost port when `127.0.0.1:8080` is occupied; release testing should keep 8080 free to align with published docs.
- No manual intervention was required beyond killing stray demo processes prior to measurement.

## Quickstart Timing Script

- Run `./scripts/measure_quickstart.sh [--db-file PATH] [--assets-dir PATH] [--addr HOST:PORT] [--output PATH]` to automate seed/serve measurement.
- The script builds the CLI once, verifies Tailwind assets, times both commands, polls `/healthz`, inspects the rendered dashboard for the “Orion Analytics” marker, and writes JSON output to `var/duplynx_bench.json` by default.
- When the default port is busy and no override is supplied, it selects an ephemeral port and prints the fallback address; pass `--addr` to enforce a specific host:port.
- Use the emitted JSON plus audit log lines (`seed_*`, `serve_*`) to populate performance tracking for SC-001 through SC-003.

## Onboarding Support Baseline (SC-004)

- Baseline (pre-release, exported 2025-10-28): 18 “quickstart-local” support requests (11 config issues, 5 Tailwind asset misses, 2 SQLite lock conflicts).
- Tracking plan:
  - Capture weekly counts from the SupportDesk board and append them to the release summary template.
  - Note remediation actions in this document when a new release goes out (e.g., docs patches, tooling updates).
  - Re-run the quickstart timing script at each release cut to correlate support reductions with seed/serve performance trends.
- Goal: Reduce onboarding tickets by ≥50% by the next release cycle while maintaining the documented seed/serve timing budgets.
