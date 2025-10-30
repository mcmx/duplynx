# Quickstart — Create DupLynx

## Prerequisites
- Go 1.22+
- Node.js 20+ (for Tailwind CLI + templ tooling)
- templ CLI (`go install github.com/a-h/templ/cmd/templ@latest`)
- TailwindCSS CLI (`npm install -g tailwindcss` or `pnpm dlx tailwindcss`)
- SQLite3 client (for inspection) — optional

## Project Layout
```
backend/
  cmd/duplynx/           # main application entry
  internal/              # domain packages (tenancy, scans, actions, http, templ)
  ent/                   # Ent schema + generated code
  web/                   # Tailwind config, embedded assets
tests/                   # unit/integration/e2e harnesses
specs/001-create-duplynx # planning artifacts (this directory)
```

## Setup
1. **Install Node tooling** (from repo root)  
   ```bash
   npm install
   npx playwright install
   ```
2. **Install Go deps**  
   ```bash
   cd backend
   go mod tidy
   ```
3. **Generate templ components**  
   ```bash
   templ generate ./internal/templ
   ```
4. **Generate Ent schema**  
   ```bash
   go run entgo.io/ent/cmd/ent generate ./ent/schema
   ```
5. **Build Tailwind assets**  
   ```bash
   cd web
   tailwindcss -i ./input.css -o ./static/app.css --minify
   cd ..
   ```
6. **Seed demo data** (writes `var/duplynx.db`)  
   ```bash
   go run ./cmd/duplynx seed --db-file var/duplynx.db
   ```
7. **(Optional) Build a release binary**  
   ```bash
   go build -o ../bin/duplynx ./cmd/duplynx
   ```

## Running the App
```bash
go run ./cmd/duplynx serve \
  --db-file var/duplynx.db \
  --addr :8080 \
  --embed-static
```

Visit `http://localhost:8080` to select the sample tenant, choose a machine, and open scan boards. The server emits structured logs for tenant/machine selections and duplicate actions; tail them in another terminal while testing.

## Testing
- **Unit & integration**: `go test ./...`
- **Contract verification**: `go test ./tests/contract`
- **Benchmarks**: `go test -run=^$ -bench=. ./tests/perf`
- **UI smoke & accessibility** (Playwright): `npx playwright test`
- **Full CI sweep (lint + tests + perf)**: `make ci`

## Development Tips
- Run `templ generate --watch` and `tailwindcss --watch` during UI work (build artifacts land in `backend/web/static/` and are served via the Go binary).
- Use `HTMX-Trigger` headers from the backend to refresh specific components after keeper assignments.
- Leverage `ENT_DEBUG=1` when tracing SQL during development.
- To simulate multi-tenant access rules locally, send the `X-Duplynx-Tenant` header with requests; cross-tenant attempts return an HTML scope warning.
