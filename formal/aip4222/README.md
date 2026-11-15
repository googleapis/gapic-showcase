# AIP-4222 Formal Specification in Lean 4

Zero-dependency Lean 4 (`v4.35.0`) formal specification and executable oracle for
[AIP-4222: Routing Headers](https://google.aip.dev/client-libraries/4222).

## Prerequisites

Install Lean 4 (`v4.35.0`) via [`elan`](https://github.com/leanprover/elan) (matching `lean-toolchain`).
This package uses only Lean core (`Init`/`Std`) and builds completely offline with no external dependencies.

## Build & Verify Proofs

From `formal/aip4222`:

```bash
lake --wfail build
```

To print the axiom dependencies of all 33 exported theorems and verify zero `sorryAx` or custom axioms:

```bash
env lean Aip4222/AxiomAudit.lean
```

## Run the Executable Oracle (`aip4222-oracle`)

`Main.lean` builds the `aip4222-oracle` executable with three CLI modes:

- **Smoke suite** (default): verifies built-in examples and regression cases:

  ```bash
  lake exe aip4222-oracle
  # equivalent to:
  lake exe aip4222-oracle smoke
  ```

- **Batch JSON oracle (`--batch-json`)**: reads a JSON array of routing cases on `stdin` and writes evaluated pairs, headers, compiled templates, and validation errors as JSON on `stdout`:

  ```bash
  lake exe aip4222-oracle --batch-json < cases.json
  ```

- **Emit conformance suite (`--emit-conformance`)**: serializes the canonical machine-checked conformance vectors from `Aip4222/Vectors.lean` to `stdout` or a file path:

  ```bash
  lake exe aip4222-oracle --emit-conformance ../../schema/google/showcase/v1beta1/aip4222_routing_suite.json
  ```

## Verify Showcase Go Sync & Live gRPC Wire Conformance

From the `gapic-showcase` repository root, run the suite-sync check (which invokes the Lean oracle when `lake` is on `PATH`) and the live in-memory `bufconn` Echo client wire tests:

```bash
go test -count=1 -v ./server/services -run 'TestAip4222ConformanceSuiteSync|TestShowcaseEchoGrpcWireConformance'
```

- `TestAip4222ConformanceSuiteSync` verifies `schema/google/showcase/v1beta1/aip4222_routing_suite.json` matches `lake exe aip4222-oracle --emit-conformance`.
- `TestShowcaseEchoGrpcWireConformance` drives the generated Showcase `EchoClient` against an in-memory gRPC server, validating exact header emission and pinning three `gapic-generator-go` divergence subtests.

## Module Layout (`6,547` lines total)

- `Aip4222.lean` — Top-level package re-exporting all formalization modules.
- `Aip4222/Syntax.lean` — Path template AST, static `WellFormed` predicate, total parser, and `MethodRoutingSpec`.
- `Aip4222/Encoding.lean` — Kernel-reducible RFC 6570 §3.2.2 UTF-8 percent-encoding and `&`-joined header formatting.
- `Aip4222/Semantics.lean` — Declarative AIP-4222 reference evaluator ("last match wins" + explicit-over-implicit precedence).
- `Aip4222/Compiler.lean` — Executable compiled segment matcher IR, key-group compiler, and faithful bug models (`H1`, `H3a`/`H3b`/`H3c`, `H4`).
- `Aip4222/Theorems.lean` — Machine-checked universal refinement theorems (`compiled_matcher_refines_spec`, `compiled_rule_extensional_eq_spec`, `compiled_rule_single_variant_eq_spec`, `parseRoutingParameter_sound`, `h2_mid_multi_wildcard_always_dead`, encoding invariants) and concrete falsification/divergence witnesses.
- `Aip4222/Examples.lean` — Machine-checked proofs of `google/api/routing.proto` Examples 1–9, edge cases, and initial `librarian` hypotheses (`H1`, `H2`).
- `Aip4222/Vectors.lean` — Canonical positive extraction and negative parser-rejection conformance vectors verified against the theorems.
- `Aip4222/AxiomAudit.lean` — `#print axioms` verification gate covering 33 theorems across Phase 1 examples and Phase 2 universal theorems.
- `Main.lean` — `aip4222-oracle` CLI entrypoint (`smoke`, `--batch-json`, `--emit-conformance`).
