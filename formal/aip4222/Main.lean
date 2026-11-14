/-
Copyright 2026 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-/

import Aip4222

open Aip4222
open Aip4222.Examples

private def renderPairs (pairs : List (String × String)) : String :=
  if pairs.isEmpty then "<none>"
  else joinStringsWith ", " (pairs.map fun (k, v) => s!"{k}={v}")

def main : IO Unit := do
  let cases : List (String × RoutingRule) := [
    ("Example 1", ex1Rule),
    ("Example 2", ex2Rule),
    ("Example 3a", ex3aRule),
    ("Example 3b", ex3bRule),
    ("Example 3c", ex3cRule),
    ("Example 4", ex4Rule),
    ("Example 5", ex5Rule),
    ("Example 6a", ex6aRule),
    ("Example 6b", ex6bRule),
    ("Example 7", ex7Rule),
    ("Example 8", ex8Rule),
    ("Example 9", ex9Rule)
  ]
  IO.println "AIP-4222 Canonical Examples (google/api/routing.proto)"
  IO.println "Example    | Unencoded Pairs | Encoded x-goog-request-params"
  IO.println "-----------|-----------------|----------------------------"
  for (label, rule) in cases do
    match evalRoutingRulePairs rule canonicalRequest, evalRoutingRule rule canonicalRequest with
    | .ok pairs, .ok encoded =>
      let encDisplay := if encoded == "" then "<none>" else encoded
      IO.println s!"{label} | {renderPairs pairs} | {encDisplay}"
    | .error err, _ =>
      IO.println s!"{label} | ERROR: {err} |"
    | _, .error err =>
      IO.println s!"{label} | | ERROR: {err}"
