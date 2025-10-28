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
1. **Install Go deps**  
   ```bash
   cd backend
   go mod tidy
   ```
2. **Generate templ components**  
   ```bash
   templ generate ./internal/templ
   ```
3. **Generate Ent schema**  
   ```bash
   go run entgo.io/ent/cmd/ent generate ./ent/schema
   ```
4. **Build Tailwind assets**  
   ```bash
   cd web
   tailwindcss -i ./input.css -o ./static/app.css --minify
   cd ..
   ```
5. **Seed database**  
   ```bash
   go run ./cmd/duplynx seed --db-file var/duplynx.db
   ```

## Running the App
```bash
go run ./cmd/duplynx serve \
  --db-file var/duplynx.db \
  --addr :8080 \
  --embed-static
```

Visit `http://localhost:8080` to select the sample tenant, choose a machine, and open scan boards.

## Testing
- **Unit & integration**: `go test ./...`
- **Contract verification**: `go test ./tests/contract`
- **UI smoke** (Playwright): `npx playwright test`

## Development Tips
- Run `templ generate --watch` and `tailwindcss --watch` during UI work.
- Use `HTMX-Trigger` headers from the backend to refresh specific components after keeper assignments.
- Leverage `ENT_DEBUG=1` when tracing SQL during development.
