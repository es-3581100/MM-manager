# Phase Zero.2 — Semantic Runtime Contract Integration

## Purpose

This chunk changes the architectural center of gravity from “repository matrix with optional semantics” toward a provenance-first semantic object/tier runtime whose Matrix is the human/agent observatory.

It does **not** add a new predictor, claim physical-storage performance, or allow structural evidence to trigger movement directly.

## Canonical flow

```text
project/repository evidence
  -> SemanticDescriptor/v1
  -> PredictionReceipt/v1
  -> AdmissionDecision/v1
  -> TierEvent/v1
  -> Outcome/v1
  -> RuntimeTrace/v1
  -> Matrix observatory (next integration)
```

## Implemented

- Six versioned runtime contracts plus descriptor-set envelope.
- `project-to-semantic.go` adapter from `app-dir-matrix.project/v1`.
- Adapter hard-rejects unknown project schema versions.
- Structural tree entries become semantic objects without inventing dependency edges.
- Imported mutability/recomputability remain `UNKNOWN` unless explicitly evidenced.
- All imported provenance remains `authority: none`.
- Kotlin runtime uses typed mutability/recomputability and a formal `PredictionReceipt`.
- `CONFIGURED_POLICY` is separated from observed `ACTOR_EVENT` evidence.
- The tested runtime emits per-policy `RuntimeTrace/v1` JSON containing receipts, admissions, tier events, outcomes, and compact metrics.
- Non-oracle traces are checked for oracle-evidence leakage.

## Regression result

The generated `fairness-summary.csv` is byte-for-byte identical to the owner-provided repaired Phase Zero.1 bundle. Contract integration did not change policy performance/accounting.

## Explicit inference boundary

1. tree membership and blob/tree identity are structural evidence only;
2. path adjacency does not create semantic dependencies;
3. mutability and recomputability remain `UNKNOWN` without explicit evidence;
4. imported evidence has `authority: none`.

## Verification result

```text
PHASE_ZERO_2_INTEGRATION_PASS
PHASE_ZERO_2_CONTRACT_VALIDATION_PASS
project_descriptors=6
runtime_traces=6
FAIRNESS_PARITY_PASS
```

## Next gate

Make the Matrix ingest `RuntimeTrace/v1` as a read-only observatory overlay: object state, prediction/admission/outcome trail, queue pressure, evidence source, and per-policy accounting. Runtime traces remain evidence; they do not become execution authority.
