# Memory Matrix Runtime Contracts — v1

Canonical Phase Zero.2 flow:

```text
evidence -> SemanticDescriptor/v1 -> PredictionReceipt/v1 -> AdmissionDecision/v1 -> TierEvent/v1 -> Outcome/v1 -> RuntimeTrace/v1
```

Rules: evidence never carries execution authority; structural tree membership does not imply semantic dependency; unknown mutability/recomputability stays UNKNOWN; configured policy evidence is distinct from observed actor events; Oracle evidence is lab/test-only; RuntimeTrace/v1 is an observability record, not a command surface.
