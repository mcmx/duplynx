# Quickstart: DupLynx Demo Runtime

## Prerequisites
- Go 1.22
- Node.js 20+ (for Tailwind build)
- `make`, `sqlite3`, and Playwright dependencies installed per repo README
- From repo root, run `make tidy && make ci` to ensure baseline readiness

## Build Tailwind Assets
```bash
npm install
npm run build:tailwind   # produces backend/web/dist assets
```
- Flags `--db-file` and `--assets-dir` default to `var/duplynx.db` and `backend/web/dist`; override via `DUPLYNX_DB_FILE` and `DUPLYNX_ASSETS_DIR` env vars. The `--log-level` flag defaults to `info` (`DUPLYNX_LOG_LEVEL`).

## Seed Demo Database
```bash
cd backend
go run ./cmd/duplynx seed --db-file ../var/duplynx.db --assets-dir ./web/dist
```
- Command recreates the database deterministically and logs audit events
- If the file is locked or unwritable, the CLI exits with a descriptive error

## Serve the Dashboard
```bash
cd backend
go run ./cmd/duplynx serve --db-file ../var/duplynx.db --assets-dir ./web/dist --addr 0.0.0.0:8080
```
- Wait for “HTTP server listening” message, then visit http://localhost:8080
- Use tenant switcher to explore seeded machines, scans, and duplicate groups
- You may also set `DUPLYNX_ADDR` to change the bind address without editing flags.

## Measure Quickstart Timing
```bash
./scripts/measure_quickstart.sh --db-file var/duplynx.db --assets-dir backend/web/dist
```
- Records seed/serve durations, writes results to `var/duplynx_bench.json`, and prints timing summaries.
- Script exits non-zero when budgets in the specification are exceeded.

## Run Smoke Verification
```bash
make smoke-demo   # wraps seeding + serve smoke test hitting /
```
- Fails fast if the root page does not render expected UI markers

## Troubleshooting
- **Missing assets**: rerun the Tailwind build; server refuses to start without `web/dist`
- **Locked database**: ensure no other DupLynx process is running; delete `var/duplynx.db` and reseed
- **Playwright failures**: update browsers via `npx playwright install` and rerun quickstart
- **Configuration overrides**: all CLI flags can be provided via `DUPLYNX_*` environment variables for CI scripts
