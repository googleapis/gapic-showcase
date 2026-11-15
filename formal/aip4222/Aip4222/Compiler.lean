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

import Aip4222.Semantics

namespace Aip4222

/--
Compiled character-stream segment IR matching `gaxi::routing_parameter::Segment`
in `google-cloud-rust`.
-/
inductive RuntimeSeg where
  | literal (s : String)
  | singleWildcard
  | multiWildcard
  | trailingMultiWildcard
  deriving DecidableEq, Repr, BEq

/--
Consume a prefix of `remains` according to a single `RuntimeSeg`, returning the number of
characters consumed. Mirrors `Segment::match_size` in `google-cloud-rust`.
-/
def RuntimeSeg.matchSize (seg : RuntimeSeg) (remains : List Char) : Option Nat :=
  match seg with
  | .literal lit =>
    if List.isPrefixOf lit.toList remains then some lit.toList.length else none
  | .singleWildcard =>
    let n := (remains.takeWhile (· != '/')).length
    if n == 0 then none else some n
  | .multiWildcard =>
    some remains.length
  | .trailingMultiWildcard =>
    match remains with
    | [] => some 0
    | '/' :: _ => some remains.length
    | ':' :: _ => some remains.length
    | _ => none

/-- Match a sequence of `RuntimeSeg`s sequentially, returning `(consumedLen, remainder)`. -/
def matchSegList : List RuntimeSeg → List Char → Option (Nat × List Char)
  | [], remains => some (0, remains)
  | seg :: rest, remains => do
    let n ← seg.matchSize remains
    let (m, rem') ← matchSegList rest (remains.drop n)
    return (n + m, rem')

/-- Compiled `(prefix, matching, suffix)` triple of runtime segments. -/
structure CompiledTemplate where
  headerKey : String
  prefixSegs : List RuntimeSeg
  matchSegs : List RuntimeSeg
  suffixSegs : List RuntimeSeg
  deriving DecidableEq, Repr, BEq

/--
Evaluate a `CompiledTemplate` over a character stream, mirroring `gaxi::routing_parameter::value`
in `google-cloud-rust`. Requires all input characters to be consumed (`rem3.isEmpty`) and a
non-empty match (`mCount != 0`).
-/
def evalCompiled (ct : CompiledTemplate) (val : String) : Option String := do
  let chars := val.toList
  let (_, rem1) ← matchSegList ct.prefixSegs chars
  let (mCount, rem2) ← matchSegList ct.matchSegs rem1
  let (_, rem3) ← matchSegList ct.suffixSegs rem2
  if !rem3.isEmpty || mCount == 0 then
    none
  else
    some (String.ofList (rem1.take mCount))

/-- Flush non-empty accumulated literal buffer into `acc`. -/
def flushLitBuffer (buf : String) (acc : List RuntimeSeg) : List RuntimeSeg :=
  if buf == "" then acc else acc ++ [.literal buf]

/--
Go-style imperative-buffer segment compiler mirroring `annotateSegments` in
`librarian/internal/sidekick/swift/format_path.go` (when called with `isPrefix`/`isSuffix`)
and `librarian/internal/sidekick/rust/annotate_method.go` (when called with `(false, false)`).
-/
def annotateSegmentsLoop (segs : List AtomSeg) (totalLen index : Nat)
    (isPrefix isSuffix : Bool) (litBuf : String) (acc : List RuntimeSeg) :
    List RuntimeSeg :=
  match segs with
  | [] =>
    let finalBuf := if isPrefix && totalLen > 0 then litBuf ++ "/" else litBuf
    flushLitBuffer finalBuf acc
  | seg :: rest =>
    match seg with
    | .multiWildcard =>
      let acc1 := flushLitBuffer litBuf acc
      let nextSeg :=
        if totalLen == 1 && !isSuffix then
          RuntimeSeg.multiWildcard
        else if totalLen != index + 1 then
          RuntimeSeg.multiWildcard
        else
          RuntimeSeg.trailingMultiWildcard
      annotateSegmentsLoop rest totalLen (index + 1) isPrefix isSuffix "" (acc1 ++ [nextSeg])
    | .singleWildcard =>
      let buf1 := if index != 0 || isSuffix then litBuf ++ "/" else litBuf
      let acc1 := flushLitBuffer buf1 acc
      annotateSegmentsLoop rest totalLen (index + 1) isPrefix isSuffix "" (acc1 ++ [.singleWildcard])
    | .lit s =>
      let buf1 := if index != 0 || isSuffix then litBuf ++ "/" else litBuf
      annotateSegmentsLoop rest totalLen (index + 1) isPrefix isSuffix (buf1 ++ s) acc

/-- Compile an `AtomSeg` slice (`prefix`, `matching`, or `suffix`) into `RuntimeSeg`s. -/
def annotateSegments (segs : List AtomSeg) (isPrefix isSuffix : Bool) : List RuntimeSeg :=
  if segs.isEmpty then []
  else annotateSegmentsLoop segs segs.length 0 isPrefix isSuffix "" []

/-- Structural unmerged emission loop used as a clean intermediate specification for `annotateSegmentsLoop`. -/
def emitUnmergedLoop (segs : List AtomSeg) (totalLen index : Nat)
    (isPrefix isSuffix : Bool) : List RuntimeSeg :=
  match segs with
  | [] => if isPrefix && totalLen > 0 then [.literal "/"] else []
  | seg :: rest =>
    let tail := emitUnmergedLoop rest totalLen (index + 1) isPrefix isSuffix
    match seg with
    | .multiWildcard =>
      let nextSeg :=
        if totalLen == 1 && !isSuffix then
          RuntimeSeg.multiWildcard
        else if totalLen != index + 1 then
          RuntimeSeg.multiWildcard
        else
          RuntimeSeg.trailingMultiWildcard
      nextSeg :: tail
    | .singleWildcard =>
      let slash := if index != 0 || isSuffix then [.literal "/"] else []
      slash ++ [.singleWildcard] ++ tail
    | .lit s =>
      let slash := if index != 0 || isSuffix then [.literal "/"] else []
      slash ++ [.literal s] ++ tail

/-- Unmerged counterpart of `annotateSegments`. -/
def emitUnmerged (segs : List AtomSeg) (isPrefix isSuffix : Bool) : List RuntimeSeg :=
  if segs.isEmpty then []
  else emitUnmergedLoop segs segs.length 0 isPrefix isSuffix

/-- Verified template compiler setting `isPrefix = true` and `isSuffix = true` on the outer slices. -/
def compileTemplate (t : CaptureTemplate) : CompiledTemplate :=
  {
    headerKey := t.headerKey
    prefixSegs := annotateSegments t.prefixSegs true false
    matchSegs := annotateSegments t.matchSegs false false
    suffixSegs := annotateSegments t.suffixSegs false true
  }

/--
Unpatched `librarian` Rust compiler (`rust/annotate_method.go`), which passes `(false, false)` on all
three slices—dropping prefix trailing `/`, suffix leading `/`, and single-segment `**` suffix
boundary checks (**H3a**–**H3c**).
-/
def compileTemplateBuggyRust (t : CaptureTemplate) : CompiledTemplate :=
  {
    headerKey := t.headerKey
    prefixSegs := annotateSegments t.prefixSegs false false
    matchSegs := annotateSegments t.matchSegs false false
    suffixSegs := annotateSegments t.suffixSegs false false
  }

/-- Compiled variant bound to a request field path. -/
structure CompiledVariant where
  fieldPath : FieldPath
  compiled : CompiledTemplate
  deriving DecidableEq, Repr, BEq

/-- Variants grouped per header key in **reverse** declaration order (so first match wins). -/
structure CompiledKeyGroup where
  headerKey : String
  variants : List CompiledVariant
  deriving DecidableEq, Repr, BEq

/-- First non-`none` match from a list (modeling generated `.or_else(...)` chains). -/
def firstMatchWins {α : Type} : List (Option α) → Option α
  | [] => none
  | some x :: _ => some x
  | none :: rest => firstMatchWins rest

/-- Last non-`none` match from a list (modeling declaration-order `last match wins`). -/
def lastMatchWins {α : Type} (vals : List (Option α)) : Option α :=
  vals.foldl (fun acc x => match x with | some v => some v | none => acc) none

/--
Insert a compiled variant into the group for `ct.headerKey`, preserving first-seen header key
declaration order while prepending variants so conflict resolution reverses order per key.
-/
def insertVariantDeclOrder (groups : List CompiledKeyGroup) (fp : FieldPath)
    (ct : CompiledTemplate) : List CompiledKeyGroup :=
  match groups with
  | [] => [{ headerKey := ct.headerKey, variants := [{ fieldPath := fp, compiled := ct }] }]
  | g :: rest =>
    if g.headerKey == ct.headerKey then
      { g with variants := { fieldPath := fp, compiled := ct } :: g.variants } :: rest
    else
      g :: insertVariantDeclOrder rest fp ct

def compileRoutingParamsDeclOrder :
    List RoutingParameter → List CompiledKeyGroup → Except String (List CompiledKeyGroup)
  | [], groups => .ok groups
  | param :: rest, groups => do
    let tmpl ← parseRoutingParameter param.field param.pathTemplate
    let ct := compileTemplate tmpl
    let fp := FieldPath.parse param.field
    compileRoutingParamsDeclOrder rest (insertVariantDeclOrder groups fp ct)

/-- Compile a `RoutingRule` into key groups preserving first-seen key declaration order. -/
def compileRoutingRuleDeclOrder (rule : RoutingRule) : Except String (List CompiledKeyGroup) :=
  compileRoutingParamsDeclOrder rule.parameters []

def evalCompiledVariant (v : CompiledVariant) (req : Request) : Option String :=
  match req.getField v.fieldPath with
  | none => none
  | some val => evalCompiled v.compiled val

def evalCompiledKeyGroup (g : CompiledKeyGroup) (req : Request) : Option (String × String) :=
  let matched := firstMatchWins (g.variants.map (fun v => evalCompiledVariant v req))
  matched.map (fun val => (g.headerKey, val))

/-- Evaluate compiled key groups into unencoded `(headerKey, capturedValue)` pairs. -/
def evalCompiledKeyGroupsPairs (groups : List CompiledKeyGroup) (req : Request) :
    List (String × String) :=
  groups.filterMap (fun g => evalCompiledKeyGroup g req)

/-- Evaluate a `RoutingRule` through the compiled declaration-order key-group pipeline. -/
def evalCompiledRuleDeclOrder (rule : RoutingRule) (req : Request) : Except String String :=
  (compileRoutingRuleDeclOrder rule).map fun groups =>
    formatHeaderPairs (evalCompiledKeyGroupsPairs groups req)

/-- Lookup a header key's value in an ordered `(key, value)` pair list. -/
def lookupHeaderKey (k : String) : List (String × String) → Option String
  | [] => none
  | (k', v') :: rest => if k' == k then some v' else lookupHeaderKey k rest

/-- Pairwise distinctness of header keys in a `(key, value)` pair list. -/
def NoDupKeys (pairs : List (String × String)) : Prop :=
  (pairs.map Prod.fst).Nodup

/--
Model of `librarian`'s **H1** nil/empty routing bug (`parser/routing_info.go` & `api/method.go`):
`option (google.api.routing) = {}` produces an empty slice, so `HasRouting()` returns `false`
and execution mistakenly falls back to `google.api.http`.
-/
def evalMethodRoutingBuggyNilCheck (spec : MethodRoutingSpec) (req : Request) :
    Except String String :=
  match spec.routing with
  | some rule =>
    if rule.parameters.isEmpty then
      match spec.http with
      | none => .ok ""
      | some http => .ok (formatHeaderPairs (evalHttpRulePairs http req))
    else
      evalRoutingRule rule req
  | none =>
    match spec.http with
    | none => .ok ""
    | some http => .ok (formatHeaderPairs (evalHttpRulePairs http req))

/--
Model of `librarian`'s **H4** implicit HTTP routing fold (`rust/templates/crate/src/transport.rs.gotmpl`):
deduplicates HTTP variables by field path, then unconditionally prepends `"&"` and skips both empty
string filtering and RFC 6570 percent-encoding (`TODO(#2548)`).
-/
def formatHttpBuggyRust (http : HttpRule) (req : Request) : String :=
  let uniqueFields := http.bindings.foldl (fun acc binding =>
    binding.foldl (fun acc2 vb =>
      if vb.fieldPath.isEmpty || acc2.contains vb.fieldPath then acc2
      else acc2 ++ [vb.fieldPath]
    ) acc
  ) []
  uniqueFields.foldl (fun acc fp =>
    match req.getField fp with
    | none => acc
    | some val => acc ++ "&" ++ fp.toDotted ++ "=" ++ val
  ) ""

end Aip4222
