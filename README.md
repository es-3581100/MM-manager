# App-Dir Matrix Core

**Status:** active prototype / reusable core  
**Current core line:** v0.2.x  
**Product posture:** local-first, human/agent shared project intelligence  
**Primary artifact:** portable HTML project matrix with machine-readable contracts and export surfaces

> **Rights notice:** Copyright © 2026 es-3581100. **ALL RIGHTS RESERVED.** This repository is publicly viewable source, **not open-source licensed software**. No permission is granted to copy, modify, redistribute, host, train on, commercialize, create derivatives from, or incorporate the work into another system except under prior written permission. Patent, trademark/branding, commercial/application, and other IP rights are expressly reserved. See [`LICENSE`](LICENSE) and [`LEGAL.md`](LEGAL.md).

App-Dir Matrix Core turns repositories, references, notes, documentation, and machine-derived relationships into a shared project map that both a human and an AI agent can use.

The current core already treats structure, semantic evidence, human annotation, and derived display metadata as different evidence classes. The long-term product extends that foundation into persistent project intelligence, agent context delivery, stateful task execution, team coordination, and a pack ecosystem without making the local core dependent on one model provider or one hosted service.

## Product Thesis

The product is not intended to be only a repository visualizer, code graph, RAG wrapper, or AI coding frontend.

The target is:

```text
sources / repos / docs / notes / reference packs
                    ↓
        normalized + provenance-aware ingest
                    ↓
       pinned structural project evidence
                    ↓
 semantic / runtime / ownership / test overlays
                    ↓
         shared human + agent project map
                    ↓
   scoped context / memory / tasks / actors
                    ↓
      history / evidence / reusable packs
```

The durable value is the **project model around the model**: what exists, how it is connected, why the system believes a relationship, what changed, what is currently in scope, and what an agent is allowed to treat as evidence.

## Core Principles

- **Local-first by default.** A useful project should not require a hosted control plane.
- **Portable, rights-reserved core.** HTML, JSON, YAML, manifests, hashes, and interoperable exports remain first-class; portability and public visibility do not grant reuse rights.
- **Human and agent share the same map.** The UI is not a decorative layer over a separate hidden agent database.
- **Evidence classes stay separate.** Git structure is not silently promoted into semantic truth; annotations remain annotations; inferred and verified edges remain distinguishable.
- **Model/provider independence.** The project intelligence layer should outlive any individual LLM, coding agent, provider, or router.
- **No authority inflation.** Context, retrieval, semantic similarity, bookmarks, or a successful prior task do not automatically become execution permission.
- **Reproducible artifacts.** Important generated outputs should be versioned, exportable, and hashable.
- **Progressive complexity.** A new user should be able to add a project and explore it without understanding graphs, MCP, actors, vector stores, Q-tables, or runtime internals.

## Current Core

The current App-Dir Matrix Core establishes the reusable baseline:

```text
ref glob
  → normalize URLs / Markdown / annotations
  → resolve repository metadata
  → pin recursive Git trees
  → verify explicitly named subpaths
  → classify deployment / source / boundaries
  → render tree + matrix + structural graph
  → optionally overlay semantic edges
  → export portable/versioned artifacts
```

Current surfaces include repository deployment cards, keyword/regex search, tree and matrix navigation, structural parent/child graphing, optional CodeGraph semantic overlays, pinned-tree loading, live-branch refresh, manifest export, offline snapshot export, and an MCP-shaped ref-glob ingestion contract.

The following roadmap deliberately separates **what should exist eventually** from **what is already implemented**.

---

# Roadmap

## Phase 0 — Core Truth, Reproducibility, and Packaging

**Goal:** make the existing matrix a dependable reusable product core before adding intelligence layers.

### TODO

- Stabilize and version the App-Dir Matrix document/schema contract.
- Keep structural edges, semantic edges, annotations, derived display levels, and runtime evidence separately typed.
- Add deterministic build metadata and artifact SHA-256 output.
- Add explicit status vocabulary such as `implemented`, `verified`, `planned`, `candidate`, `stale`, `failed`, and `drifted` where useful.
- Add regression fixtures for ref-glob normalization, branch resolution, recursive-tree truncation, subpath verification, and export/import round-trips.
- Add snapshot comparison so a user can see branch/tree drift from the last pinned artifact.
- Preserve offline HTML as a complete human-readable project handoff.
- Formalize the core manifest so other applications can consume the matrix without scraping the UI.

### Deployment target

- Single-file/offline HTML remains supported.
- Local file/project workflow remains the minimum viable deployment.
- Static hosting remains possible without requiring a server.

### Capital posture

**No paywall on the core.** This phase is the adoption and trust foundation.

The local matrix, basic repository indexing, basic search, portable exports, and base MCP/ref-glob contract are intended to remain available in the no-cost Community product tier. That product posture does not grant source-code reuse, derivative-work, hosting, training, or commercial rights.

### Exit gate

A project can be ingested, pinned, rendered, exported, reopened, compared, and independently verified without relying on hidden hosted state.

---

## Phase 1 — Project Intelligence Map

**Goal:** move from “browse the tree” to “understand the project.”

### TODO

- Add **two-node path queries**: select any two nodes and show every known structural/semantic path connecting them.
- Add **fan-in**, **fan-out**, and **blast-radius** analysis.
- Add semantic and structural **heatmaps**.
- Add edge filters for `structural`, `verified semantic`, `proposed`, `inferred`, `human annotation`, and future runtime evidence.
- Add Git-history overlays where available.
- Add test-coverage/test-location overlays.
- Add ownership/maintainer overlays.
- Add authority/trust-boundary overlays.
- Add deployment/runtime-boundary overlays.
- Add bookmarks and named node groups.
- Add **working sets** that let a user isolate a project slice for a task.
- Add saved views such as Security, Runtime, UI, Data, Agent, Deployment, Tests, and Ownership.
- Add dynamic graph modes including standard graph, radar/relationship view, and freeze-frame snapshots.
- Add comparison views for proposed-vs-verified relationships.
- Add a lightweight “why is this connected?” evidence panel for every non-structural edge.

### Deployment target

- Continue to run locally in the browser.
- Preserve static/offline readability even when richer analysis requires optional local helpers.

### Capital posture

Still primarily **Community / basic_4_life**.

This phase proves that the matrix is useful enough to become a daily project-navigation tool before subscriptions are asked to carry the business.

### Exit gate

A user can answer architecture questions from the project map without repeatedly re-reading the entire repository.

---

## Phase 2 — Agent Context Server and Scoped Exports

**Goal:** turn the project map into a reliable context provider for coding agents.

### TODO

- Add a stable local API/MCP query surface over the matrix.
- Add **agent-scope export**: export only the nodes, edges, notes, constraints, and evidence needed for one task.
- Add context capsules/working-set bundles that can be handed to another model or agent without rebuilding the project map.
- Add explicit scope metadata: included, excluded, read-only, reference-only, unresolved, and stale.
- Add query primitives for path, neighbors, dependencies, ownership, tests, provenance, related docs, and evidence.
- Add machine exports for JSON and YAML.
- Add graph/table exports designed for QTable, hash-table, LangGraph, and LangChain adapters without making those frameworks core dependencies.
- Add multi-application / multi-repository graphs with explicit project boundaries.
- Add generalized deployment-chain visualization such as:

```text
desktop → MCP → bridge → runtime → SQLite
```

- Add portable “agent brief” generation containing project purpose, relevant architecture, source pins, known boundaries, current working set, and unresolved questions.

### Deployment target

- Local MCP/API process paired with the static/browser UI.
- Optional headless CLI for CI, agent runtimes, and automation.
- Keep BYO-model/provider as the default posture.

### Planned capital stream

Introduce **Personal** as the first paid convenience tier, approximately **$10–15/month**, only after the local/no-cost Community product is useful on its own.

Potential Personal value:

- persistent private project workspaces;
- larger saved working sets;
- richer agent-context exports;
- history across sessions;
- optional managed indexing convenience.

The underlying local project format and basic MCP surface remain available in the no-cost Community product tier; their source remains rights-reserved.

### Exit gate

Two different agents can consume the same project model and receive consistent scoped project context with provenance.

---

## Phase 3 — Persistent Project Memory and Knowledge Layer

**Goal:** let the project remember what the repository alone cannot represent.

### TODO

- Add local **SQLite** persistence for project metadata, history, bookmarks, saved views, working sets, evidence records, and user annotations.
- Add wiki-like project indexing: entities, concepts, decisions, architecture notes, reference packs, and backlinks.
- Add time-aware snapshots and change history.
- Add source freshness and drift tracking.
- Add claim/evidence records that can answer “why does the system believe this?”
- Add optional lightweight vector retrieval / ChromaDB adapter for similarity search while keeping vector similarity non-authoritative.
- Add optional QTable storage for learned routing/workflow hints without confusing learned utility with factual truth or permission.
- Add hashed lookup/index surfaces for fast deterministic retrieval.
- Add reference-pack versioning, provenance, compatibility, and update status.
- Add import/export bridges for external knowledge stores without requiring one canonical vendor.

### Deployment target

- Local database by default.
- Portable backup/export bundle containing the human-readable matrix plus machine-readable state.
- Optional encrypted/private hosted sync later, but never required for single-user operation.

### Planned capital stream

Introduce **Pro** at approximately **$20–25/month** when persistent intelligence provides clear value beyond the static matrix.

Potential Pro value:

- larger persistent histories;
- advanced semantic indexing;
- project-memory tooling;
- automated drift checks;
- richer multi-project intelligence;
- managed compute for expensive indexing operations.

Optional hosted semantic/indexing compute can be usage-metered, but ordinary local use should not be token-metered by App-Dir Matrix.

### Exit gate

The system can preserve useful project knowledge across agent/model/session changes without replacing source code or verified evidence as authority.

---

## Phase 4 — Hyper-Index and Scaffold Generation

**Goal:** make the matrix useful both for understanding an existing project and creating the next one.

### TODO

- Expand wiki-like indexes into navigable human/agent hyper-indexes.
- Generate project/reference indexes from repos, docs, specs, and curated reference globs.
- Add reusable reference-pack templates.
- Add “project skeleton from map” / scaffold generation.
- Generate starter README, AGENT.md, manifest, source-map/index, test map, and architecture pages from a selected project profile.
- Allow a working set or reference pack to become a reusable scaffold for a new project.
- Add template inheritance without silently copying runtime authority or secrets.
- Add schema-aware “create a new app shaped like this architecture” output.
- Add a diff explaining what was copied as structure, what was only referenced, and what still requires implementation.

### Deployment target

- Local generator/CLI plus browser UI.
- Export complete project starter bundles.
- Make generated artifacts consumable by ordinary Git workflows rather than requiring App-Dir Matrix at runtime.

### Planned capital stream

Scaffold generation remains available in the no-cost Community product tier at a useful baseline; repository source rights remain reserved.

Premium value may include curated/maintained templates, advanced private project profiles, managed compatibility checks, and richer regeneration/diff tooling.

### Exit gate

A user can turn a proven project/reference map into a transparent, editable starting scaffold without treating the source project as a magical one-shot prompt.

---

## Phase 5 — Browser Companion and Lightweight Full Runtime

**Goal:** make capture and navigation available where developers already browse, while preserving the local-first architecture.

### TODO

### 15. Companion Chrome/Chromium extension

- Capture GitHub repositories, subpaths, docs, links, and selected text directly into an App-Dir Matrix ref glob.
- Add current page/repo to a working set.
- Bookmark nodes/reference pages into the active project.
- Open the local matrix at the matching project/node.
- Show basic project membership and stale/pinned status for a repository page.
- Keep extension permissions minimal and explicit.
- Do not make browser content an execution authority merely because it was captured by the extension.

### 16. Light full-extension runtime

- Allow a lightweight App-Dir Matrix project surface to run directly as an extension UI where practical.
- Cache/open selected portable matrix snapshots.
- Support local project switching, bookmarks, saved views, and scope export without requiring the full desktop/runtime stack.
- Use a local companion process only for capabilities the browser sandbox cannot safely provide.
- Keep the extension runtime intentionally smaller than the full project-intelligence runtime.

### Additional deployment work

- Add **Go-App** / WebAssembly/PWA deployment research for a compiled portable browser surface.
- Preserve ordinary static HTML as a lowest-common-denominator artifact.
- Add Linux-friendly local installation and user-level service options where a helper process is required.

### Planned capital stream

The companion extension itself should be a no-cost acquisition surface or included with Community/Personal; no source-code license is implied.

Pro may include cross-project persistence, managed sync, larger indexes, and advanced browser-to-project workflows.

### Exit gate

A user can encounter a useful repository/document in the browser and add it to the correct project intelligence space with minimal friction and clear provenance.

---

## Phase 6 — Actor / Task Runtime

**Goal:** move from passive project knowledge to explicit, stateful work without turning the matrix into an uncontrolled autonomous agent.

### TODO

- Add task/actor state machines with explicit states and allowed transitions.
- Record task start, state transitions, inputs, outputs, evidence, model/provider identity when available, and completion/abort reason.
- Bind actors/workers to explicit project scopes and working sets.
- Add task receipts and reproducible handoffs.
- Add queued work, pause/resume/review states, and failure recovery.
- Add “propose → review → approve → execute/emit” patterns for actions that cross project boundaries.
- Keep execution capabilities separate from descriptive project knowledge.
- Add multi-agent handoff without allowing agents to silently widen each other’s scopes.
- Add task replay/inspection where deterministic or safe to do so.

### Deployment target

- Local actor/task process integrated with the matrix API.
- Optional connectors for OpenCode, Codex, Claude, GPT, Gemini, and other agents through shared project contracts rather than provider-specific project formats.

### Planned capital stream

This phase supports the stronger **Pro** offering and begins to justify a separate automation/compute add-on for users who want hosted workers.

The base product should continue to support BYO agent/model so users are paying for durable project intelligence and orchestration—not for artificial token lock-in.

### Exit gate

Long-running or multi-step work has inspectable state, bounded scope, durable history, and explicit handoff evidence.

---

## Phase 7 — Team Project Operating Environment

**Goal:** let multiple humans and agents share one durable project understanding.

### TODO

- Shared project matrices and reference packs.
- Team bookmarks, saved views, scopes, and working sets.
- Human/agent handoff records.
- Permissions for read, annotate, propose, approve, administer, and execute where execution surfaces exist.
- Shared project history and change review.
- CI integration for matrix rebuilds, drift checks, semantic refresh, and artifact publication.
- Review workflows for proposed-vs-verified edges and annotations.
- Signed/versioned project snapshots.
- Conflict handling for concurrent project knowledge edits.
- Private/self-hosted team deployment.

### Deployment target

- Self-hostable team server first.
- Optional managed cloud collaboration second.
- Git-compatible project artifacts remain exportable even when collaboration is hosted.

### Planned capital stream

Introduce **Team** at approximately **$30–40/user/month** when genuine collaboration exists.

Team value should come from shared state, governance, history, orchestration, and reduced duplicated agent work—not by removing local features from individual users.

### Exit gate

A team can share a project map and agent context without each person/model independently reconstructing project truth.

---

## Phase 8 — Reference-Pack / Project-Pack Ecosystem

**Goal:** turn curated project intelligence into a reusable ecosystem rather than a private collection of one-off maps.

### TODO

- Formalize installable **project packs / knowledge packs / reference packs**.
- Define pack metadata, compatibility, provenance, source pins, signatures/hashes, update channels, and trust levels.
- Support community, official, private-team, and vendor-maintained packs.
- Add pack composition/inheritance.
- Add compatibility tests and stale-source warnings.
- Allow packs to carry references, verified paths, semantic maps, docs, search presets, tests, agent instructions, workflow hints, and scaffolds while keeping execution authority separate.
- Build a searchable pack registry/catalog.

Example future packs could cover a framework, language ecosystem, infrastructure stack, security domain, game/tooling ecosystem, or internal corporate platform.

### Deployment target

- Portable, documented pack format; applicable rights remain governed by each pack's license or rights notice.
- Local install/import first.
- Optional hosted catalog/registry later.

### Planned capital streams

- No-cost community product packs, subject to the rights/license terms attached to each pack.
- Paid official/maintained packs.
- Private organization packs as part of Team/Enterprise.
- Optional marketplace revenue share for third-party pack authors.

The format itself should remain open enough that paid packs compete on curation, maintenance, verification, and convenience rather than lock-in.

### Exit gate

A new project can gain high-quality domain context from an installable pack while preserving traceable sources and project-local control.

---

## Phase 9 — Enterprise Governance and Large-Scale Deployment

**Goal:** support organizations that need centralized policy and evidence without changing the core project model into a closed SaaS-only format.

### TODO

- SSO and RBAC.
- Organization-level policy and trust configuration.
- Signed provenance and artifact-attestation workflows.
- Audit export and retention controls.
- Air-gapped/offline deployment.
- Organization-hosted model/provider routing.
- Large repository/monorepo scaling.
- Central pack registry and controlled update channels.
- Compliance-oriented evidence retention without pretending the product itself guarantees compliance.
- Support/SLA and controlled upgrade channels.

### Deployment target

- Self-hosted and air-gapped supported as first-class options.
- Managed enterprise service optional.

### Planned capital stream

**Enterprise: custom pricing.**

Enterprise revenue should come from deployment, governance, support, scale, policy controls, integration, and managed infrastructure—not from making the core project format proprietary.

### Exit gate

An organization can adopt App-Dir Matrix without surrendering source/project knowledge to an external SaaS and without forking the product to meet ordinary governance requirements.

---

# Planned Capital Streams

The capital model should activate in the same order as product capability.

| Stream | Earliest sensible phase | Planned model | What pays for it |
|---|---:|---|---|
| Community / `basic_4_life` | 0 | No-cost product tier / rights-reserved source | Adoption, trust, ecosystem growth |
| Personal | 2 | ~**$10–15/month** | Persistent private project convenience, richer agent context, managed indexing options |
| Pro | 3 | ~**$20–25/month** | Advanced project memory, semantic computation, drift/history, multi-project tooling, automation |
| Team | 7 | ~**$30–40/user/month** | Shared project state, handoffs, governance, collaboration, CI, agent coordination |
| Hosted compute add-ons | 3+ | Usage-based / bundled allowance | Expensive semantic indexing, optional hosted workers, large graph computation |
| Official project/reference packs | 8 | One-time and/or subscription | Maintained, curated, versioned domain intelligence |
| Pack marketplace | 8 | Revenue share | Distribution/payment/discovery for third-party pack authors |
| Enterprise | 9 | Custom | Self-hosting, SSO/RBAC, policy, audit, air-gap, support, integrations, scale |

## Capital guardrails

The project should avoid monetizing the things that create adoption and portability:

```text
basic local HTML export
basic repository indexing
basic local search
basic portable project format
basic MCP/query contract
BYO model/provider support
```

The durable paid value should instead accumulate around:

```text
persistence
scale
semantic computation
collaboration
history
automation
managed infrastructure
verified/maintained packs
governance
enterprise integration
```

This keeps the no-cost Community product useful, makes paid tiers additive rather than coercive, and avoids tying the business to reselling model tokens. Product-tier pricing is separate from source-code licensing.

---

# Product Evolution Summary

```text
Phase 0  reproducible repo hyper-index
    ↓
Phase 1  project intelligence map
    ↓
Phase 2  agent context server
    ↓
Phase 3  persistent project memory
    ↓
Phase 4  hyper-index + scaffold generator
    ↓
Phase 5  browser companion + portable runtime
    ↓
Phase 6  actor / task runtime
    ↓
Phase 7  team project operating environment
    ↓
Phase 8  reference-pack ecosystem
    ↓
Phase 9  enterprise governance / deployment
```

The important sequencing rule is that **App-Dir Matrix should become useful before it becomes large**.

The earliest commercially credible product is not the full multi-agent operating environment. It is the combination of **Phase 1 + Phase 2**: a project intelligence map that can reliably feed scoped, provenance-aware context to humans and coding agents.

Everything after that should compound the same project model rather than replace it.

---

# Long-Term Product Position

The long-term product can be summarized as:

> **A local-first project intelligence and control layer that gives humans and AI agents a shared, persistent, verifiable understanding of a software project.**

The graph is a view. The LLM is a client. The repository is a source. The durable product is the project model, its evidence, its state, and the ability to move that understanding between humans, agents, tools, sessions, and deployments without losing provenance.