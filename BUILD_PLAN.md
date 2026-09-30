# App-Dir Matrix — Master Build Plan

## Baseline

```yaml
baseline:
  repository: /mnt/data/app-dir-matrix-build (isolated build workspace; no pre-existing Git repository)
  branch: null
  head: null
  dirty_state: not_applicable
  current_core_version: 0.2.0
  source_core_sha256: eb84b0403c240c519aa1ffb202dd916120433bc4e684ab1da0cf8c4f90d0bf7d
  builder_prompt_sha256: ef7c90be67bc49c48e75852bf23dc29c88625c27b73fb17f9a107deccd25a609
  roadmap_sha256: 873d02e7c4608fdda8e5b9803c104cceefa96321eafb5a87282be1384356e864
  design_reference_sha256: 944e1d2a4c0fd6711155f1ba63fcc13db1d0dbfb88ee1696799981a0636fdac4
  tests_found: existing HTML runtime only; no standalone test suite in supplied core artifact
  existing_phase_status:
    phase_0: partial
    phase_1: partial
    phase_2: candidate
    phase_3_to_9: not_started
  recovered_artifacts:
    - app-dir-matrix-core-v2.html
    - app-dir-matrix-core.html
    - kol-app-dir-matrix-v2.html
    - builder-plan prompt with embedded roadmap + Anti-AI reference
  unresolved_discrepancies:
    - supplied core has browser-side ingest envelope and pinned-tree hydration but no standalone reproducible packager/verifier
    - existing HTML can export hydrated snapshots, but independent reopen/compare verification was not separately implemented
    - live browser smoke is unavailable in the current container due Chromium headless/DBus startup failure
```

## Architecture decomposition

The product is split into five durable layers so UI, models, and deployment surfaces remain clients rather than authority sources:

1. **Project document** — versioned, portable repository/evidence model.
2. **Resolver/ingest** — ref-glob normalization, repository metadata resolution, pinned tree acquisition, named-path verification.
3. **Artifact engine** — deterministic HTML/manifest/context export and integrity hashing.
4. **Verification/drift** — reopen, validate, compare, and explain differences.
5. **Views/adapters** — browser UI, graph/semantic overlays, MCP/API, persistence, actors, teams, packs, enterprise controls.

Current implementation slice establishes layers 1–4 around the existing v0.2.0 HTML view.

## Phase dependency graph

```text
P0 Core Truth / Packaging
  ↓
P1 Project Intelligence
  ↓
P2 Agent Context
  ↓
P3 Persistent Memory
  ↓
P4 Hyper-Index / Scaffolds
  ↓
P5 Browser Companion
  ↓
P6 Actor / Task Runtime
  ↓
P7 Team Environment
  ↓
P8 Pack Ecosystem
  ↓
P9 Enterprise
```

Cross-cutting invariants: evidence classes remain separate; no context becomes execution authority; portable/local formats remain exportable under the repository rights notice; model/provider choice remains external to the project truth model.

---

## Phase 0 — Core Truth, Reproducibility, Packaging

```yaml
phase: 0
objective: Make ingest/pin/render/export/reopen/compare independently reproducible and verifiable.
prerequisites: Existing v0.2.0 HTML core and ref-glob process contract.
current_evidence:
  - v0.2.0 core embeds app_dir_matrix.ingest_ref_glob@1.0.0
  - browser core already distinguishes annotation/repo metadata/git tree/verified subpath/semantic graph/derived level
  - supplied core can hydrate pinned trees and export browser snapshots
existing_implementation:
  - keyword/regex structural index
  - structural graph + optional CodeGraph overlay
  - browser-side ref-glob normalization envelope
  - pinned-tree browser hydration/export
  - Phase 0 Go project schema/validator/normalizer/resolver/packager/verifier/comparator from this checkpoint
missing_work:
  - exercise pin command against real GitHub API and record fixtures from actual responses
  - version-negotiated project migrations
  - canonical manifest/context export from the Go layer
  - browser-runtime smoke in a healthy browser environment
  - malformed/escaped ref-glob corpus expansion
  - import of prior v0.2.0 browser snapshots into project/v1 documents
files_or_modules_likely_affected:
  - schema/
  - internal/matrix/
  - cmd/appdir-matrix/
  - web/
data_contracts:
  - app-dir-matrix.project/v1
  - app_dir_matrix.ingest_ref_glob@1.0.0
migrations:
  - future: v0.2 browser manifest/snapshot -> project/v1
  - future: project/v1 -> next schema version via explicit migration command
tests:
  - normalize escaped Markdown/URLs and dedupe repositories
  - default branch and explicit-ref selection
  - recursive tree truncation rejection
  - named subpath presence
  - authority inflation rejection
  - deterministic build byte equality
  - reopen + tree/pin comparison
  - drift detection
verification:
  - go test ./...
  - go vet ./...
  - scripts/smoke.sh
  - node --check over extracted runtime JS
  - browser smoke when environment supports it
risks:
  - GitHub API rate limits / branch ambiguity
  - schema drift between Go model and browser JS
  - self-contained HTML growing large for monorepos
security_or_authority_boundaries:
  - all evidence authority=none
  - truncated trees cannot be presented as complete
  - semantic edges are never inferred from filenames
  - failed resolver paths do not widen scope or fall back to mutation
 design_surface: Preserve existing matrix UI until Phase 0 contracts stabilize; no decoration expansion.
deployment_target: Local CLI + self-contained static HTML.
capital_stream_enabled_by_this_phase: Community/basic_4_life only.
explicit_non_goals:
  - billing
  - hosted state
  - agent execution
  - team permissions
exit_gate: A project can be ingested, pinned, rendered, exported, reopened, compared, and independently verified without hidden hosted state.
status: partial
```

## Phase 1 — Project Intelligence Map

```yaml
phase: 1
objective: Answer architecture questions from typed structural and semantic evidence.
prerequisites: Phase 0 stable project document and drift semantics.
current_evidence: Existing UI has structural tree, graph, regex filters, node focus, semantic CodeGraph overlay.
existing_implementation: Basic structural browsing/search/graph; permissive CodeGraph JSON import.
missing_work:
  - two-node path queries
  - fan-in/fan-out/blast radius
  - typed edge filters and why-connected evidence inspector
  - Git/test/ownership/authority/deployment overlays
  - bookmarks, named groups, working sets, saved views
  - proposed-vs-verified comparison
files_or_modules_likely_affected: internal/graph, web/intelligence, schema overlay contracts
data_contracts: typed edge records with provenance and evidence class
migrations: additive overlay collections only after Phase 0 migration contract exists
tests: path correctness, cycle handling, overlay isolation, unresolved semantic references
verification: fixture graphs + artifact inspection + narrow/wide UI checks
risks: conflating structure with semantics; graph readability on large repos
security_or_authority_boundaries: overlays are descriptive; trust-boundary labels do not grant permissions
design_surface: dense instrumentation map using existing dark technical material language
deployment_target: local HTML/runtime
capital_stream_enabled_by_this_phase: Community
explicit_non_goals: persistent knowledge, agents acting on map
exit_gate: A user can answer architecture questions from the map without repeatedly re-reading the repository.
status: partial
```

## Phase 2 — Agent Context Server and Scoped Exports

```yaml
phase: 2
objective: Serve the same provenance-aware project model to multiple agents.
prerequisites: P0 stable document; P1 query semantics.
current_evidence: MCP-shaped ingest contract exists; no live project query server is supplied.
existing_implementation: Candidate contract/export ideas only.
missing_work: local API/MCP, task scopes, context capsules, path/neighbor/dependency/ownership/test/provenance queries, JSON/YAML exports, headless CLI, multi-repo views
files_or_modules_likely_affected: cmd/server, internal/query, schema/context-capsule
data_contracts: scoped context capsule with source pins and query reason
migrations: none until context capsule contract is versioned
tests: two-client consistency, scope reduction, stale pin reporting
verification: independent clients consume identical fixture state
risks: context widening; adapter-specific semantics leaking into core
security_or_authority_boundaries: context is never permission; execution endpoints are out of scope
 design_surface: agent/user query inspector, optional
deployment_target: localhost API/MCP + CLI
capital_stream_enabled_by_this_phase: earliest Personal seam; open basic query surface stays free
explicit_non_goals: hosted-only context; provider lock-in
exit_gate: Two different agents receive consistent scoped context with provenance.
status: candidate
```

## Phase 3 — Persistent Project Memory and Knowledge Layer

```yaml
phase: 3
objective: Persist knowledge source trees cannot carry while retaining explainability.
prerequisites: P0 schema/migrations; P2 query contract.
current_evidence: Roadmap only.
existing_implementation: none in this workspace.
missing_work: SQLite store, history, snapshots, claims/evidence, wiki entities, optional vector/Chroma and QTable adapters, hashed indexes, pack versioning
files_or_modules_likely_affected: internal/store, migrations, adapters/vector, adapters/qtable
data_contracts: immutable evidence IDs + temporal validity + source freshness
migrations: first durable DB migration set
tests: create/migrate/reopen/crash safety/export-import/provenance survival
verification: deterministic DB export + replay
risks: similarity/utility treated as truth; migration data loss
security_or_authority_boundaries: vector/QTable records explicitly non-authoritative
design_surface: history/claim/evidence inspector
deployment_target: local SQLite first
capital_stream_enabled_by_this_phase: Pro + optional compute
explicit_non_goals: mandatory hosted database
exit_gate: Project explains what/why/when and what changed.
status: not_started
```

## Phase 4 — Hyper-Index and Scaffold Generation

```yaml
phase: 4
objective: Generate reusable human/agent references and project scaffolds from explicit verified knowledge.
prerequisites: P3 durable knowledge + P0 deterministic artifact generation.
current_evidence: Existing Matrix itself is a human/agent hyper-index pattern.
existing_implementation: no generalized scaffold engine.
missing_work: linked README/AGENT/reference generation, templates/inheritance, safe pack import, user-authored region preservation, deterministic manifests
files_or_modules_likely_affected: internal/scaffold, templates/, generators/
data_contracts: generated-region ownership metadata + template provenance
migrations: template version transitions
tests: regeneration idempotence, user-content preservation, secret exclusion
verification: golden fixtures and diff stability
risks: overwriting user content; propagating secrets or authority claims
security_or_authority_boundaries: generated instructions do not create permissions
design_surface: scaffold preview/diff
 deployment_target: CLI + local UI
capital_stream_enabled_by_this_phase: Pro
explicit_non_goals: opaque one-shot project generation
exit_gate: Reusable knowledge creates safe project structures with provenance.
status: not_started
```

## Phase 5 — Browser Companion and Lightweight Full Runtime

```yaml
phase: 5
objective: Bring browser-discovered references into project space with low friction and explicit provenance.
prerequisites: P2 scoped ingest API; P4 ref/scaffold contracts.
current_evidence: Browser static HTML runtime exists; extension does not.
existing_implementation: none.
missing_work: MV3 companion, minimal permissions, project selection, capture inbox, local runtime, local project switching, bookmarks/views, scoped export, Go-App/WASM research
files_or_modules_likely_affected: extension/, web/runtime/, cmd/helper
data_contracts: capture envelope + provenance
migrations: extension storage versioning only for convenience state
tests: permissions, capture fidelity, no credential collection, uninstall safety
verification: extension smoke + local helper integration
risks: permission creep; browser storage leakage
security_or_authority_boundaries: no cookies/tokens/broad injection; local-first privacy
design_surface: compact companion, not full IDE
 deployment_target: Chromium MV3 + optional PWA/WASM
capital_stream_enabled_by_this_phase: acquisition/Pro convenience
explicit_non_goals: browser as authority store
exit_gate: Browser reference enters correct project with provenance.
status: not_started
```

## Phase 6 — Actor / Task Runtime

```yaml
phase: 6
objective: Make long/multi-step work explicit, inspectable, bounded, and recoverable.
prerequisites: P2 scoped context; P3 persistence; P0 verification receipts.
current_evidence: separate micro-state prototype supplied, not integrated.
existing_implementation: none in App-Dir Matrix.
missing_work: typed task states/transitions, scopes, receipts, pause/resume/review/abort/failure, provider identity, handoff, bounded replay
files_or_modules_likely_affected: internal/tasks, internal/actors, schema/task
data_contracts: task state machine + immutable receipt chain
migrations: persisted task schema
tests: invalid transition rejection, interruption recovery, stale context, handoff scope invariance
verification: replay fixture and forced failure recovery
risks: authority widening; hidden autonomy
security_or_authority_boundaries: propose/review/approve/execute surfaces separate; execution authority external and explicit
design_surface: state/receipt timeline
 deployment_target: local worker runtime first
capital_stream_enabled_by_this_phase: Pro/compute add-on
explicit_non_goals: unrestricted autonomous agent
exit_gate: Long work has inspectable state, bounded scope, durable history, explicit handoff evidence.
status: not_started
```

## Phase 7 — Team Project Operating Environment

```yaml
phase: 7
objective: Share one durable project understanding across humans and agents.
prerequisites: P3 persistence; P6 receipts; stable conflict semantics.
current_evidence: roadmap only.
existing_implementation: none.
missing_work: shared matrices/scopes/views, permission roles, handoffs, history, CI refresh, signed snapshots, conflict handling, self-hosted server
files_or_modules_likely_affected: server/team, authz, sync, ci
data_contracts: actor identity, role grants, signed snapshot metadata
migrations: team DB + permission migrations
tests: concurrent edit, least privilege, stale snapshot, role transitions
verification: multi-user integration fixture
risks: authorization complexity; sync conflicts
security_or_authority_boundaries: read/annotate/propose/approve/administer/execute explicitly distinct
design_surface: review/conflict/permission views
 deployment_target: self-hosted team server first
capital_stream_enabled_by_this_phase: Team
explicit_non_goals: cloud-only collaboration
exit_gate: Team shares one project map/context without reconstructing truth independently.
status: not_started
```

## Phase 8 — Reference-Pack / Project-Pack Ecosystem

```yaml
phase: 8
objective: Install reusable domain intelligence without losing source traceability.
prerequisites: P4 scaffolds; P7 trust/signature model.
current_evidence: current KoL instance demonstrates a ref-pack concept but not a generic registry format.
existing_implementation: informal ref pack data only.
missing_work: open pack format, metadata, compatibility, pins, hashes/signatures, update channels, trust levels, composition, tests, registry/catalog
files_or_modules_likely_affected: schema/pack, internal/packs, registry/
data_contracts: signed pack manifest and compatibility range
migrations: pack schema migrations
tests: composition conflicts, stale source warnings, signature failures
verification: install/uninstall/upgrade fixture
risks: trusted pack mistaken for execution permission
security_or_authority_boundaries: packs cannot silently acquire execution authority
design_surface: pack provenance/trust inspector
 deployment_target: local registry/import first
capital_stream_enabled_by_this_phase: paid packs/marketplace
explicit_non_goals: proprietary core format
exit_gate: Domain context installs with traceable sources and project-local control.
status: not_started
```

## Phase 9 — Enterprise Governance and Large-Scale Deployment

```yaml
phase: 9
objective: Add organization governance without turning the core into SaaS-only software.
prerequisites: P7 permissions/history; P8 signed packs; mature P0 migration/attestation.
current_evidence: roadmap only.
existing_implementation: none.
missing_work: SSO/RBAC, organization policy, provenance attestation, audit retention/export, air-gap, provider routing, monorepo scale, controlled registries, upgrade channels, support tooling
files_or_modules_likely_affected: enterprise/, auth/, attestation/, deployment/
data_contracts: organization policy + attestations + retention config
migrations: governed upgrade framework
tests: offline install, RBAC matrix, audit export, policy conflict, scale fixtures
verification: controlled deployment acceptance suite
risks: compliance overclaims; proprietary lock-in
security_or_authority_boundaries: governance config explicit and auditable; no hidden vendor authority
 design_surface: admin/audit surfaces
 deployment_target: self-hosted/air-gapped first; managed optional
capital_stream_enabled_by_this_phase: Enterprise custom
explicit_non_goals: claiming product alone guarantees compliance
exit_gate: Organization adopts without surrendering project knowledge or forking for ordinary governance.
status: not_started
```

## Test strategy

Use the narrowest useful test first: format → unit → resolver contract → deterministic artifact generation → reopen/verify → drift compare → JS syntax → browser smoke → packaging. Every phase adds fixtures at its own boundary; generated output is inspected, not merely trusted because a command returned 0.

## Design validation plan

Visual direction remains the existing dense technical instrument-panel language: dark graphite planes, semantic evidence colors, monospace operational typography, and restrained local interaction. Before any new visual scope, validate identity, hierarchy, consistency, interaction states, and responsive resilience against the embedded Anti-AI reference. Phase 0 intentionally changes no visual styling.

## Rollout/deployment strategy

Local CLI and self-contained HTML first. Hosted services remain optional additive layers. Stable files and schemas precede daemons. Every hosted/team surface must preserve export back to the documented portable project document.

## Capital activation gates

Community remains Phase 0+. Personal cannot activate before a useful scoped context service exists. Pro cannot activate before persistent project memory/semantic computation exists. Team waits for real shared state/permissions. Marketplace waits for the pack format. Enterprise waits for governance/deployment capability.

## Next executable slice

Phase Zero.2 semantic runtime contracts are now established and regression-verified. Next, make the Matrix ingest `memory-matrix.runtime-trace/v1` as a read-only observatory overlay: object/tier state, prediction → admission → outcome trail, queue pressure, evidence source, and policy accounting. Do not add predictor sophistication or let runtime evidence become execution authority before this observability path is proven. Project schema migration/version negotiation and browser runtime smoke remain Phase 0 blockers.

## Direction pivot — Phase Zero.2

The project architecture now separates three durable planes:

- **Go evidence/package plane:** ingest, normalize, pin, verify, replay, drift, portable project truth.
- **Semantic runtime/policy plane:** descriptor → prediction receipt → admission → tier event → outcome.
- **Matrix observatory plane:** human/agent visualization and query over structural evidence plus read-only runtime traces.

The Q-table or any learned policy belongs behind the prediction/admission boundary. It may propose utility; it does not become factual truth or execution authority.
