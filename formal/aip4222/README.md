# AIP-4222 Formal Specification in Lean 4

Zero-dependency Lean 4 (`v4.35.0`) formal specification and executable oracle for
[AIP-4222: Routing Headers](https://google.aip.dev/client-libraries/4222).

## Prerequisites

Install Lean 4 (`v4.35.0`) via [`elan`](https://github.com/leanprover/elan) (matching `lean-toolchain`).
This package uses only Lean core (`Init`/`Std`) and builds completely offline with no external dependencies.

## Build & Verify Proofs

```bash
lake --wfail build
```

All exported theorems in `Aip4222/AxiomAudit.lean` are checked via `#print axioms` with zero `sorry` or compiler-trust axioms.

## Run Executable Summary

```bash
lake exe aip4222-oracle
```

## Module Layout

- `Aip4222/Syntax.lean` — Path template AST, static `WellFormed` predicate, total parser, and `MethodRoutingSpec`.
- `Aip4222/Encoding.lean` — Kernel-reducible RFC 6570 §3.2.2 UTF-8 percent-encoding and `&`-joined header formatting.
- `Aip4222/Semantics.lean` — Declarative AIP-4222 reference evaluator ("last match wins" + explicit-over-implicit precedence).
- `Aip4222/Compiler.lean` — Executable compiled segment matcher IR, key-group compiler, and faithful bug models (`H1`, `H3a`/`H3b`/`H3c`, `H4`).
- `Aip4222/Theorems.lean` — Machine-checked universal refinement theorems (`compiled_matcher_refines_spec`, `compiled_rule_extensional_eq_spec`, `compiled_rule_single_variant_eq_spec`, `parseRoutingParameter_sound`, `h2_mid_multi_wildcard_always_dead`, encoding invariants) and concrete falsification/divergence witnesses.
- `Aip4222/Examples.lean` — Machine-checked proofs of `google/api/routing.proto` Examples 1–9, edge cases, and initial `librarian` hypotheses (`H1`, `H2`).
- `Aip4222/AxiomAudit.lean` — `#print axioms` verification gate covering Phase 1 examples and Phase 2 universal theorems.
- `Main.lean` — `aip4222-oracle` CLI entrypoint.
