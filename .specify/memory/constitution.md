<!--
Sync Impact Report
Version: N/A → 1.0.0
Modified Principles:
- None (initial adoption)
Added Sections:
- I. Code Quality Fidelity
- II. Testing Standardization
- III. Consistent User Experience
- IV. Performance Discipline
- Delivery Standards
- Workflow Integration
Removed Sections:
- None
Templates Requiring Updates:
- ✅ .specify/templates/plan-template.md
- ✅ .specify/templates/spec-template.md
- ✅ .specify/templates/tasks-template.md
Follow-up TODOs:
- None
-->
# Duplynx Constitution

## Core Principles

### I. Code Quality Fidelity
- Every change MUST pass automated linting, static analysis, and formatting checks configured for the repository. Failing gates block the merge.
- Teams MUST document public APIs, error contracts, and architectural decisions in-line or in companion docs before merging.
- Dead code, unused feature flags, and deprecated pathways MUST be removed or sunset as part of the change set, with exceptions logged in the release plan.
**Rationale:** Enforced quality gates keep the codebase understandable, reduce unexpected regressions, and ensure future contributors inherit clear intent.

### II. Testing Standardization
- Feature work MUST begin with executable tests (unit, integration, contract) that capture expected behavior before implementation, following a strict red-green-refactor loop.
- CI pipelines MUST execute the full regression suite on every merge request and block merges when coverage for touched modules drops below 90% or any test fails.
- All automated tests MUST run deterministically and complete within 15 minutes for standard pull requests; longer runs require mitigation tasks in the plan.
**Rationale:** Predictable testing validates behavior early, prevents regressions, and produces fast feedback loops that protect delivery velocity.

### III. Consistent User Experience
- All user-facing surfaces MUST use the shared Duplynx design tokens, typography, spacing, and component library; new patterns require prior UX sign-off.
- Accessibility acceptance tests MUST verify WCAG 2.1 AA compliance (keyboard navigation, contrast, screen reader semantics) before release.
- Product copy, error messages, and interaction flows MUST remain consistent across platforms; deviations demand documented rationale and approved experiment IDs.
**Rationale:** Consistency builds user trust, reduces usability defects, and ensures inclusive experiences across devices and locales.

### IV. Performance Discipline
- Every feature MUST define baseline and target performance budgets (latency, throughput, memory) in the specification before development begins.
- Automated performance regression tests or benchmarks MUST run in CI for critical paths and fail the pipeline when thresholds exceed agreed budgets.
- Production deployments MUST include observability instrumentation (metrics, tracing, logging) that exposes the defined performance indicators within 24 hours of launch.
**Rationale:** Intentional performance management keeps the product responsive at scale and enables quick diagnosis when degradations occur.

## Delivery Standards

- Specifications and implementation plans MUST map requirements, tests, accessibility checks, and performance budgets to the four core principles.
- Implementation plans MUST identify quality gates (linting, testing, UX review, performance baselines) before Phase 0 research begins.
- Release checklists MUST attach evidence—CI runs, accessibility audits, benchmark reports, and monitoring dashboards—demonstrating compliance for each principle.

## Workflow Integration

- Phase 0 research MUST log anticipated UX, quality, testing, and performance risks, with mitigation tasks added to the plan template.
- Code reviews MUST include explicit approvals for quality, testing, UX, and performance guardianship; missing approvals block the merge.
- Post-release reviews MUST compare live telemetry and UX feedback to committed budgets, creating follow-up tasks within two business days when drift is detected.

## Governance

- This constitution supersedes conflicting guidelines for Duplynx delivery; product and engineering leads are accountable for enforcement.
- Amendments require an RFC citing impacted sections, review by architecture, QA, and UX leads, and recorded migration steps for in-flight work.
- Version updates follow semantic versioning: MAJOR for principle removals or breaking governance changes, MINOR for new principles or material guidance, PATCH for clarifications.
- Compliance audits occur at least quarterly, sampling recent releases. Findings require remediation plans within five business days and tracking until closure.
- Temporary deviations MUST be documented with owner, expiration date, and mitigation plan; expired deviations automatically trigger review.

**Version**: 1.0.0 | **Ratified**: 2025-10-27 | **Last Amended**: 2025-10-27
