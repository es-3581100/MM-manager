# Phase 2 — Read-only Agent Context CLI

This slice introduces the first executable Phase-2 context contract for MM-manager.

## Boundary

The agent CLI is an evidence/query surface only.

```text
project/v1
  -> validated AgentIndex
  -> deterministic read-only query
  -> JSON context / receipt
```

It does not mutate the project, execute project actions, infer permissions, or promote evidence into authority.

Every returned context object and node carries:

```json
{"authority":"none"}
```

## Commands

```text
appdir-matrix agent bootstrap
appdir-matrix agent resolve
appdir-matrix agent inspect
appdir-matrix agent expand
appdir-matrix agent scope
appdir-matrix agent verify-receipt
```

All commands require `--project <project.json>`.

Resolution is deliberately deterministic and dependency-free in this first slice: query tokens are matched against normalized project node metadata. No model, embedding store, vector score, or Q-table participates in truth or permission.

## Pointers

Nodes use stable local pointers:

```text
mem://repo/<owner/repo>
mem://path/<owner/repo>/<path>
mem://evidence/<evidence-id>
```

Unknown pointers fail closed.

## Context capsules

`agent scope` emits `app-dir-matrix.context-capsule/v1` containing only the bounded matched nodes plus an embedded `app-dir-matrix.retrieval-receipt/v1`.

The receipt binds:

- project ID;
- SHA-256 of the exact project artifact;
- query;
- selected node IDs;
- SHA-256 of each selected node.

`agent verify-receipt` rejects project drift, unknown nodes, node-hash drift, duplicate receipt nodes, schema mismatch, receipt-ID mismatch, or authority inflation.

## OpenCode consumer

The companion native plugin is developed at:

```text
https://github.com/es-3581100/mm-manager-opencode
```

The plugin is intentionally a transport adapter. Query semantics remain in this Go core.
