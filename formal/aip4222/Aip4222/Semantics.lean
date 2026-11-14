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

import Aip4222.Encoding

namespace Aip4222

/-- True iff `seg` starts with `expected ++ ":"` (supporting `([:/].*)?` for terminal `**`). -/
def hasColonSuffix (expected seg : String) : Bool :=
  let expChars := expected.toList
  let segChars := seg.toList
  segChars.length > expChars.length &&
  segChars.take (expChars.length + 1) == expChars ++ [':']

/--
Match a well-formed list of atomic segments against `/`-split input segments.
A terminal `.multiWildcard` matches zero or more remaining segments, and when
preceded by `.lit expected` also accepts a `:` delimiter (`expected:...`) per AIP-4222's `([:/].*)?`.
-/
def matchAtomSegs : List AtomSeg → List String → Bool
  | [], [] => true
  | [], _ :: _ => false
  | [.multiWildcard], _ => true
  | [.lit expected, .multiWildcard], seg :: _ =>
    seg == expected || hasColonSuffix expected seg
  | .multiWildcard :: _, _ => false
  | _, [] => false
  | .lit expected :: patRest, seg :: segRest =>
    (seg == expected) && matchAtomSegs patRest segRest
  | .singleWildcard :: patRest, seg :: segRest =>
    (seg != "") && matchAtomSegs patRest segRest

private def matchPrefixAndCapture (t : CaptureTemplate) (rem : List String) : Option String :=
  if t.suffixSegs.isEmpty then
    if matchAtomSegs t.matchSegs rem then
      let captured := joinStringsWith "/" rem
      if captured == "" then none else some captured
    else
      none
  else
    let mLen := t.matchSegs.length
    if rem.length < mLen then
      none
    else
      let matchedSegs := rem.take mLen
      let restSegs := rem.drop mLen
      if matchAtomSegs t.matchSegs matchedSegs && matchAtomSegs t.suffixSegs restSegs then
        let captured := joinStringsWith "/" matchedSegs
        if captured == "" then none else some captured
      else if t.suffixSegs == [.multiWildcard] then
        -- Handle `.../{k=.../lit}/**` against `.../lit:verb...` per `([:/].*)?`
        match t.matchSegs.reverse, matchedSegs.reverse with
        | .lit expected :: revPatInit, lastSeg :: revSegInit =>
          if hasColonSuffix expected lastSeg &&
              matchAtomSegs revPatInit.reverse revSegInit.reverse then
            let captured := joinStringsWith "/" (revSegInit.reverse ++ [expected])
            if captured == "" then none else some captured
          else
            none
        | _, _ => none
      else
        none

/--
Declarative extraction semantics for a single `CaptureTemplate` against a field value `val`.
Empty values or unmatched templates yield `none`.
-/
def extractCapture (t : CaptureTemplate) (val : String) : Option String :=
  if val == "" then
    none
  else if t.prefixSegs.isEmpty && t.matchSegs == [.multiWildcard] && t.suffixSegs.isEmpty then
    some val
  else
    let segs := splitStringOn '/' val
    let pLen := t.prefixSegs.length
    if segs.length < pLen then
      none
    else if !matchAtomSegs t.prefixSegs (segs.take pLen) then
      none
    else
      matchPrefixAndCapture t (segs.drop pLen)

/--
Insert or overwrite `(k, v)` in an ordered key-value list.
Updating in-place preserves first-seen key declaration order while implementing
AIP-4222's "last match wins" rule for duplicate keys.
-/
def upsertHeaderPair : List (String × String) → String → String → List (String × String)
  | [], k, v => [(k, v)]
  | (k', v') :: rest, k, v =>
    if k' == k then
      (k, v) :: rest
    else
      (k', v') :: upsertHeaderPair rest k v

def evalRoutingParamsAux (req : Request) :
    List RoutingParameter → List (String × String) → Except String (List (String × String))
  | [], acc => .ok acc
  | param :: rest, acc => do
    let tmpl ← parseRoutingParameter param.field param.pathTemplate
    let nextAcc :=
      match req.getField (FieldPath.parse param.field) with
      | none => acc
      | some val =>
        match extractCapture tmpl val with
        | none => acc
        | some captured => upsertHeaderPair acc tmpl.headerKey captured
    evalRoutingParamsAux req rest nextAcc

/-- Evaluate all explicit `RoutingParameter`s in declaration order ("last match wins"). -/
def evalRoutingRulePairs (rule : RoutingRule) (req : Request) :
    Except String (List (String × String)) :=
  evalRoutingParamsAux req rule.parameters []

/-- Evaluate an explicit `RoutingRule` into the encoded `x-goog-request-params` value. -/
def evalRoutingRule (rule : RoutingRule) (req : Request) : Except String String :=
  (evalRoutingRulePairs rule req).map formatHeaderPairs

/--
Evaluate implicit routing parameters derived from a `google.api.http` rule.
Each referenced field whose value is present and non-empty is transcribed directly.
-/
def evalHttpRulePairs (http : HttpRule) (req : Request) : List (String × String) :=
  http.bindings.foldl (fun acc binding =>
    binding.foldl (fun acc2 vb =>
      if vb.fieldPath.isEmpty then
        acc2
      else
        match req.getField vb.fieldPath with
        | none => acc2
        | some val =>
          if val == "" then acc2 else upsertHeaderPair acc2 vb.fieldPath.toDotted val
    ) acc
  ) []

/--
Top-level AIP-4222 method routing evaluation.
When `spec.routing` is `some rule`, `spec.http` is completely ignored—even when
`rule.parameters` is empty (`option (google.api.routing) = {}`).
-/
def evalMethodRouting (spec : MethodRoutingSpec) (req : Request) : Except String String :=
  match spec.routing with
  | some rule => evalRoutingRule rule req
  | none =>
    match spec.http with
    | none => .ok ""
    | some http => .ok (formatHeaderPairs (evalHttpRulePairs http req))

end Aip4222
