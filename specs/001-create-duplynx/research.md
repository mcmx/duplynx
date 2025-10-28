# Phase 0 Research — Create DupLynx

## Decision 1: Layered Go service with Ent over SQLite
- **Decision**: Implement the backend as a layered Go 1.22 service using Ent-generated repositories backed by a file-scoped SQLite database.
- **Rationale**: Ent provides schema-first modeling, graph relationships, and migration helpers that dovetail with Go modules. SQLite keeps the demo footprint minimal while supporting transactions and foreign keys needed for tenant isolation. Layering the service (`tenancy`, `scans`, `actions`, `http`) matches Go best practices for separation of concerns and testability.
- **Alternatives considered**:
  - *Raw database/sql + SQLC*: More boilerplate for relations and migrations; slower to iterate on evolving schema.
  - *PostgreSQL backend*: Heavier operational footprint for a demo; adds deployment friction without delivering new capabilities for the seeded scenario.

## Decision 2: templ-driven server rendering with htmx progressive enhancement
- **Decision**: Use templ components to render all primary views (tenant selection, machine list, scan board) and layer htmx for inline interactions like keeper assignment and board updates.
- **Rationale**: templ integrates with Go templates at compile time, enabling type-safe components and precompiled assets. htmx keeps interactivity lightweight, leveraging server-rendered fragments and reducing the need for a separate SPA build. This combination aligns with the requirement for a server-first dashboard while enabling partial updates.
- **Alternatives considered**:
  - *Pure Go `html/template`*: Lacks templ’s component abstraction and type checking; harder to reuse UI blocks.
  - *Client-side SPA (React/Vue)*: Increases complexity, adds bundle tooling, and conflicts with the requirement to use templ and htmx where possible.

## Decision 3: TailwindCSS pipeline embedded via Go’s `embed` package
- **Decision**: Compile TailwindCSS using the Tailwind CLI into a static asset directory embedded into the Go binary and served from `/web/static`.
- **Rationale**: Embedding avoids external asset servers and keeps deployment self-contained. Tailwind tokens ensure consistency with DupLynx design guidance, and embedding compiled CSS with digested filenames supports caching while keeping the demo simple.
- **Alternatives considered**:
  - *Runtime JIT Tailwind*: Adds overhead and complicates binary distribution.
  - *Manual CSS*: Slows iteration and risks diverging from design tokens, reducing UX consistency.

## Decision 4: Tenant isolation enforced through request context middleware
- **Decision**: Introduce HTTP middleware that resolves tenant and machine IDs from route params, validates ownership via tenancy services, and injects scoped context objects used by repositories and view models.
- **Rationale**: Middleware centralizes authorization logic, ensuring every handler interacts with tenant-scoped repositories. Ent queries can automatically filter by `tenant_id`, preventing leakage. Context scoping simplifies auditing and logging.
- **Alternatives considered**:
  - *Separate database per tenant*: Overkill for demo, increases migration complexity.
  - *Handler-level checks only*: Prone to omissions and difficult to test comprehensively.

## Decision 5: Action workflow queued via in-memory dispatcher with audit logging
- **Decision**: Execute keeper assignments and duplicate actions through an in-memory dispatcher that records `ActionAudit` entries and updates state synchronously for the demo.
- **Rationale**: For Phase 0, immediate feedback with consistent state is more valuable than asynchronous orchestration. The dispatcher abstraction preserves an upgrade path to background workers while keeping the UI responsive (<500 ms acknowledgement).
- **Alternatives considered**:
  - *Full job queue (e.g., Redis-based)*: Adds infrastructure without real remote agents in this phase.
  - *Direct handler mutations without dispatcher*: Harder to add instrumentation and future async execution modes.
