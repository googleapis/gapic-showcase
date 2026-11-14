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

namespace Aip4222

/-- Split a list of characters on `sep`, preserving empty segments. -/
def splitCharsOn (sep : Char) : List Char → List (List Char)
  | [] => [[]]
  | c :: cs =>
    if c == sep then
      [] :: splitCharsOn sep cs
    else
      match splitCharsOn sep cs with
      | [] => [[c]]
      | cur :: rest => (c :: cur) :: rest

/-- Kernel-reducible string splitting on a single character delimiter. -/
def splitStringOn (sep : Char) (s : String) : List String :=
  (splitCharsOn sep s.toList).map String.ofList

/-- Kernel-reducible string join with `sep`. -/
def joinStringsWith (sep : String) : List String → String
  | [] => ""
  | [x] => x
  | x :: y :: rest => x ++ sep ++ joinStringsWith sep (y :: rest)

/-- Dot-separated field traversal path within a protobuf request message. -/
abbrev FieldPath := List String

def FieldPath.parse (s : String) : FieldPath :=
  if s == "" then [] else splitStringOn '.' s

def FieldPath.toDotted (fp : FieldPath) : String :=
  joinStringsWith "." fp

/--
Abstract request payload modeled as an association list from field path to string value.
Note: `Request.getField` uses first-match lookup if duplicate field paths are provided.
-/
abbrev Request := List (FieldPath × String)

def Request.getField : Request → FieldPath → Option String
  | [], _ => none
  | (k, v) :: rest, fp => if k == fp then some v else Request.getField rest fp

def Request.ofFlat (pairs : List (String × String)) : Request :=
  pairs.map fun (k, v) => (FieldPath.parse k, v)

/-- Atomic segment in an AIP-4222 path template. -/
inductive AtomSeg where
  | lit (s : String)
  | singleWildcard
  | multiWildcard
  deriving DecidableEq, Repr, BEq

/-- Parsed capture template decomposed into `(prefix, matching, suffix)`. -/
structure CaptureTemplate where
  headerKey : String
  prefixSegs : List AtomSeg
  matchSegs : List AtomSeg
  suffixSegs : List AtomSeg
  deriving DecidableEq, Repr, BEq

def CaptureTemplate.allSegs (t : CaptureTemplate) : List AtomSeg :=
  t.prefixSegs ++ t.matchSegs ++ t.suffixSegs

def isAsciiAlnum (c : Char) : Bool :=
  let n := c.toNat
  (n >= 65 && n <= 90) || (n >= 97 && n <= 122) || (n >= 48 && n <= 57)

/-- Allowed characters in AIP-4222 literals and header keys (ASCII alphanumeric plus `_`, `-`, `.`, `~`). -/
def isValidLitChar (c : Char) : Bool :=
  isAsciiAlnum c || c == '_' || c == '-' || c == '.' || c == '~'

def isValidHeaderKey (s : String) : Bool :=
  s != "" && s.toList.all isValidLitChar

def isValidAtom : AtomSeg → Bool
  | .lit s => s != "" && s.toList.all isValidLitChar
  | .singleWildcard => true
  | .multiWildcard => true

/--
Static validity predicate for AIP-4222 capture templates:
1. `headerKey` is a non-empty valid identifier.
2. `matchSegs` is non-empty.
3. Every atom segment is valid.
4. `**` appears at most once and only as the final segment of the entire template.
-/
def CaptureTemplate.isWellFormed (t : CaptureTemplate) : Bool :=
  let segs := t.allSegs
  isValidHeaderKey t.headerKey &&
  !t.matchSegs.isEmpty &&
  segs.all isValidAtom &&
  segs.dropLast.all (fun s => s != .multiWildcard)

def CaptureTemplate.WellFormed (t : CaptureTemplate) : Prop :=
  t.isWellFormed = true

instance (t : CaptureTemplate) : Decidable t.WellFormed :=
  inferInstanceAs (Decidable (t.isWellFormed = true))

def stripPrefixSlash? (s : String) : Option String :=
  match s.toList with
  | '/' :: rest => some (String.ofList rest)
  | _ => none

def stripSuffixSlash? (s : String) : Option String :=
  match s.toList.reverse with
  | '/' :: revRest => some (String.ofList revRest.reverse)
  | _ => none

def parseAtomSeg (s : String) : Except String AtomSeg :=
  if s == "" then
    .error "empty path segment"
  else if s == "*" then
    .ok .singleWildcard
  else if s == "**" then
    .ok .multiWildcard
  else if s.toList.all isValidLitChar then
    .ok (.lit s)
  else
    .error s!"invalid path segment: {s}"

def parseAtomSegs (s : String) : Except String (List AtomSeg) :=
  (splitStringOn '/' s).mapM parseAtomSeg

/-- Validate a protobuf `field` path string (`"a"` or `"a.b"`). -/
def isValidFieldPathStr (field : String) : Bool :=
  let fp := FieldPath.parse field
  field != "" && !fp.isEmpty && fp.all isValidHeaderKey

/--
Parse an AIP-4222 `(field, path_template)` pair into a validated `CaptureTemplate`.
Handles omitted `path_template` (`""`), omitted pattern shorthand (`"{parent}"`),
and ignores a single trailing `/` delimiter (`"projects/{parent}/"`).
-/
def parseRoutingParameter (field : String) (rawTemplate : String) :
    Except String CaptureTemplate := do
  if !isValidFieldPathStr field then
    throw s!"invalid routing field name: {field}"
  if rawTemplate == "" then
    return {
      headerKey := field
      prefixSegs := []
      matchSegs := [.multiWildcard]
      suffixSegs := []
    }
  -- AIP-4222: "The last symbol in a path_template may be a delimiter - it will be ignored."
  let normalizedTemplate :=
    match stripSuffixSlash? rawTemplate with
    | some trimmed => trimmed
    | none => rawTemplate
  match splitStringOn '{' normalizedTemplate with
  | [beforeOpen, afterOpen] =>
    if (splitStringOn '}' beforeOpen).length != 1 then
      throw s!"unexpected '}' before '\{' in {rawTemplate}"
    match splitStringOn '}' afterOpen with
    | [inside, afterClose] =>
      let prefixSegs ←
        if beforeOpen == "" then
          .ok []
        else
          match stripSuffixSlash? beforeOpen with
          | none => .error s!"expected '/' before '\{' in {rawTemplate}"
          | some p => parseAtomSegs p
      let (headerKey, matchSegs) ←
        match splitStringOn '=' inside with
        | [only] =>
          if isValidHeaderKey only then
            .ok (only, [.singleWildcard])
          else
            .error s!"invalid routing variable name: {only}"
        | [lhs, rhs] =>
          if isValidHeaderKey lhs then do
            let segs ← parseAtomSegs rhs
            .ok (lhs, segs)
          else
            .error s!"invalid routing variable name: {lhs}"
        | _ => .error s!"invalid variable expression: {inside}"
      let suffixSegs ←
        if afterClose == "" then
          .ok []
        else
          match stripPrefixSlash? afterClose with
          | none => .error s!"expected '/' after '}' in {rawTemplate}"
          | some s => parseAtomSegs s
      let candidate : CaptureTemplate := {
        headerKey := headerKey
        prefixSegs := prefixSegs
        matchSegs := matchSegs
        suffixSegs := suffixSegs
      }
      if candidate.isWellFormed then
        return candidate
      else
        throw s!"ill-formed capture template: {rawTemplate}"
    | _ => throw s!"expected single '}' in {rawTemplate}"
  | _ => throw s!"expected single '\{' in {rawTemplate}"

/-- Raw `google.api.RoutingParameter` entry. -/
structure RoutingParameter where
  field : String
  pathTemplate : String := ""
  deriving DecidableEq, Repr, BEq

/-- Raw `google.api.RoutingRule` annotation (`parameters := []` models `{}`). -/
structure RoutingRule where
  parameters : List RoutingParameter
  deriving DecidableEq, Repr, BEq

/--
Variable binding extracted from a `google.api.http` URI template.
Per AIP-4222, implicit HTTP routing transcribes the referenced field's value
directly when non-empty, ignoring HTTP URI sub-patterns.
-/
structure HttpVarBinding where
  fieldPath : FieldPath
  deriving DecidableEq, Repr, BEq

/-- Implicit HTTP bindings from `google.api.http` (primary rule + `additional_bindings`). -/
structure HttpRule where
  bindings : List (List HttpVarBinding)
  deriving DecidableEq, Repr, BEq

/-- Combined routing configuration on an RPC method. -/
structure MethodRoutingSpec where
  routing : Option RoutingRule
  http : Option HttpRule := none
  deriving DecidableEq, Repr, BEq

end Aip4222
