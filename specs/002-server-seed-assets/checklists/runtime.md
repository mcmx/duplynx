# Runtime Checklist: DupLynx Server Runtime

**Purpose**: Validate the quality of runtime server and CLI requirements before implementation
**Created**: 2025-11-01
**Feature**: [Link to spec.md](../spec.md)

## Requirement Completeness

- [ ] CHK001 Are CLI flag and environment variable mappings enumerated for both `serve` and `seed` commands so reviewers know every configurable input? [Completeness, Spec §FR-001]
- [ ] CHK002 Does the specification document all startup validation failure cases (config, schema migration, asset availability) along with expected handling outcomes? [Completeness, Spec §FR-002; Spec §FR-003; Spec §FR-006; Edge Cases]
- [ ] CHK003 Are the seeded dataset contents (tenants, machines, scans, duplicate groups, audit events) fully described with required attributes and relationships? [Completeness, Spec §FR-004; Data Model §Tenant/Machine/Scan/DuplicateGroup/AuditEvent]

## Requirement Clarity

- [ ] CHK004 Is the deterministic reseed behavior (wipe order, locking expectations, post-run state) articulated unambiguously for repeated executions? [Clarity, Spec §Clarifications; Spec §FR-004]
- [ ] CHK005 Are audit event fields (actor source, timestamp granularity, outcome vocabulary, metadata usage) precisely defined so logging expectations are clear? [Clarity, Spec §FR-007; Data Model §AuditEvent]

## Requirement Consistency

- [ ] CHK006 Do the Quickstart instructions and example commands match the functional requirements for CLI flags, defaults, and environment assumptions? [Consistency, Spec §FR-001; Quickstart §Serve the Dashboard]
- [ ] CHK007 Are timing expectations in success criteria consistent with the performance goals detailed in the implementation plan? [Consistency, Spec §SC-001–SC-003; Plan §Technical Context]

## Acceptance Criteria Quality

- [ ] CHK008 Are configuration validation outcomes expressed with observable pass/fail signals that reviewers can trace (e.g., explicit error messaging expectations)? [Acceptance Criteria, Spec §FR-002; Spec §User Story 1 Acceptance]
- [ ] CHK009 Does the smoke test requirement specify measurable UI markers or assertions to confirm the root page renders correctly? [Acceptance Criteria, Spec §FR-010]

## Scenario Coverage

- [ ] CHK010 Is the expected behavior defined when `serve` runs before any seeding (auto-migrate, instruct user, or fail)? [Coverage, Spec §FR-003; Spec §Assumptions]
- [ ] CHK011 Are multi-tenant route scenarios, including malformed or missing tenant identifiers, covered with explicit requirements? [Coverage, Spec §Edge Cases]

## Edge Case Coverage

- [ ] CHK012 Is handling for database file locks during seeding documented (retry strategy vs immediate abort with guidance)? [Edge Case, Spec §Edge Cases]
- [ ] CHK013 Are recovery steps outlined for partial seed failures (e.g., crash mid-population) to prevent inconsistent datasets? [Gap, Spec §FR-004]

## Non-Functional Requirements

- [ ] CHK014 Are logging/audit event requirements evaluated for storage or performance impact to ensure they meet observability goals without regressions? [Gap, Spec §FR-007; Spec §SC-003]
- [ ] CHK015 Are server responsiveness targets tied to specific monitoring or telemetry checkpoints for continued validation? [Non-Functional, Spec §SC-001; Spec §SC-003]

## Dependencies & Assumptions

- [ ] CHK016 Are Tailwind asset build prerequisites (tooling versions, build command outputs) captured so configuration reviewers can validate them? [Dependency, Quickstart §Build Tailwind Assets; Spec §Assumptions]
- [ ] CHK017 Are runtime environment dependencies (SQLite version, Playwright browsers, supported OS) enumerated and validated against assumptions? [Dependency, Spec §Assumptions; Quickstart §Prerequisites]

## Ambiguities & Conflicts

- [ ] CHK018 Is the phrase “standard developer machine” defined with representative hardware to make timing budgets testable? [Ambiguity, Spec §SC-002]
- [ ] CHK019 Is the requirement to fail when assets are missing explicitly incompatible with any embedded or fallback asset strategies to avoid conflicting expectations? [Ambiguity, Spec §Clarifications; Spec §FR-006]
