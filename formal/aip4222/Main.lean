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
import Lean.Data.Json

open Aip4222
open Aip4222.Examples

private def renderPairs (pairs : List (String × String)) : String :=
  if pairs.isEmpty then "<none>"
  else joinStringsWith ", " (pairs.map fun (k, v) => s!"{k}={v}")

private def runSmoke : IO Unit := do
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

private def runtimeSegToJson : RuntimeSeg → Lean.Json
  | .literal s => Lean.Json.mkObj [("kind", Lean.Json.str "lit"), ("val", Lean.Json.str s)]
  | .singleWildcard => Lean.Json.mkObj [("kind", Lean.Json.str "star")]
  | .multiWildcard => Lean.Json.mkObj [("kind", Lean.Json.str "dstar")]
  | .trailingMultiWildcard => Lean.Json.mkObj [("kind", Lean.Json.str "trailing_dstar")]

private def runtimeSegsToJson (segs : List RuntimeSeg) : Lean.Json :=
  Lean.Json.arr (segs.map runtimeSegToJson).toArray

private def compiledTemplateToJson (fp : FieldPath) (ct : CompiledTemplate) : Lean.Json :=
  Lean.Json.mkObj [
    ("header_key", Lean.Json.str ct.headerKey),
    ("field_path", Lean.Json.arr (fp.map Lean.Json.str).toArray),
    ("prefix", runtimeSegsToJson ct.prefixSegs),
    ("match", runtimeSegsToJson ct.matchSegs),
    ("suffix", runtimeSegsToJson ct.suffixSegs)
  ]

private def compileTemplatesForJson (params : List RoutingParameter) : Except String (List Lean.Json) :=
  params.mapM fun p => do
    let tmpl ← parseRoutingParameter p.field p.pathTemplate
    let ct := compileTemplate tmpl
    return compiledTemplateToJson (FieldPath.parse p.field) ct

private def parseRoutingParamJson (j : Lean.Json) : Except String RoutingParameter := do
  let field ← j.getObjValAs? String "field"
  let pathTemplate ← j.getObjValAs? String "path_template"
  return { field := field, pathTemplate := pathTemplate }

private def parseRoutingJson (j : Lean.Json) : Except String (Option RoutingRule) := do
  match j with
  | .null => return none
  | _ =>
    let paramsJson ← j.getObjVal? "parameters"
    let paramsArr ← paramsJson.getArr?
    let params ← paramsArr.toList.mapM parseRoutingParamJson
    return some { parameters := params }

private def parseHttpJson (j : Lean.Json) : Except String (Option HttpRule) := do
  match j with
  | .null => return none
  | _ =>
    let bindingsJson ← j.getObjVal? "bindings"
    let bindingsArr ← bindingsJson.getArr?
    let bindings ← bindingsArr.toList.mapM fun bJson => do
      let varsArr ← bJson.getArr?
      varsArr.toList.mapM fun vJson => do
        let s ← vJson.getStr?
        return { fieldPath := FieldPath.parse s }
    return some { bindings := bindings }

private def parseRequestJson (j : Lean.Json) : Except String Request := do
  match j with
  | .obj kvs =>
    let pairs ← kvs.toList.mapM fun (k, v) => do
      let s ← v.getStr?
      return (k, s)
    return Request.ofFlat pairs
  | _ => .error "expected request to be a JSON object"

private def pairsToJson (pairs : List (String × String)) : Lean.Json :=
  Lean.Json.arr (pairs.map fun (k, v) => Lean.Json.arr #[Lean.Json.str k, Lean.Json.str v]).toArray

private def routingRuleToJson : Option RoutingRule → Lean.Json
  | none => Lean.Json.null
  | some rule =>
    let paramsJson := rule.parameters.map fun p =>
      Lean.Json.mkObj
        [ ("field", Lean.Json.str p.field)
        , ("path_template", Lean.Json.str p.pathTemplate)
        ]
    Lean.Json.mkObj [("parameters", Lean.Json.arr paramsJson.toArray)]

private def httpRuleToJson : Option HttpRule → Lean.Json
  | none => Lean.Json.null
  | some http =>
    let bindingsJson := http.bindings.map fun b =>
      Lean.Json.arr (b.map fun vb => Lean.Json.str vb.fieldPath.toDotted).toArray
    Lean.Json.mkObj [("bindings", Lean.Json.arr bindingsJson.toArray)]

private def requestToJson (reqPairs : List (String × String)) : Lean.Json :=
  Lean.Json.mkObj (reqPairs.map fun (k, v) => (k, Lean.Json.str v))

private def evalCaseJson (item : Lean.Json) : Except String Lean.Json := do
  let id ← item.getObjValAs? String "id"
  let routingVal ← item.getObjVal? "routing"
  let routing ← parseRoutingJson routingVal
  let httpVal ← item.getObjVal? "http"
  let http ← parseHttpJson httpVal
  let reqVal ← item.getObjVal? "request"
  let req ← parseRequestJson reqVal
  match routing with
  | none =>
    let declPairs :=
      match http with
      | none => []
      | some h => evalHttpRulePairs h req
    let declHeader := formatHeaderPairs declPairs
    return Lean.Json.mkObj [
      ("id", Lean.Json.str id),
      ("valid", Lean.Json.bool true),
      ("error", Lean.Json.null),
      ("decl_header", Lean.Json.str declHeader),
      ("decl_pairs", pairsToJson declPairs),
      ("alpha_header", Lean.Json.str declHeader),
      ("compiled_templates", Lean.Json.arr #[])
    ]
  | some rule =>
    match evalRoutingRulePairs rule req, compileRoutingRuleDeclOrder rule,
        compileTemplatesForJson rule.parameters with
    | .ok declPairs, .ok groups, .ok tmplsJson =>
      let declHeader := formatHeaderPairs declPairs
      let sortedGroups := sortKeyGroupsAlpha groups
      let alphaHeader := formatHeaderPairs (evalCompiledKeyGroupsPairs sortedGroups req)
      return Lean.Json.mkObj [
        ("id", Lean.Json.str id),
        ("valid", Lean.Json.bool true),
        ("error", Lean.Json.null),
        ("decl_header", Lean.Json.str declHeader),
        ("decl_pairs", pairsToJson declPairs),
        ("alpha_header", Lean.Json.str alphaHeader),
        ("compiled_templates", Lean.Json.arr tmplsJson.toArray)
      ]
    | .error err, _, _ =>
      return Lean.Json.mkObj [
        ("id", Lean.Json.str id),
        ("valid", Lean.Json.bool false),
        ("error", Lean.Json.str err),
        ("decl_header", Lean.Json.str ""),
        ("decl_pairs", Lean.Json.arr #[]),
        ("alpha_header", Lean.Json.str ""),
        ("compiled_templates", Lean.Json.arr #[])
      ]
    | _, .error err, _ | _, _, .error err =>
      return Lean.Json.mkObj [
        ("id", Lean.Json.str id),
        ("valid", Lean.Json.bool false),
        ("error", Lean.Json.str err),
        ("decl_header", Lean.Json.str ""),
        ("decl_pairs", Lean.Json.arr #[]),
        ("alpha_header", Lean.Json.str ""),
        ("compiled_templates", Lean.Json.arr #[])
      ]

private def extractVectorToJson (v : ExtractVector) : Except String Lean.Json := do
  let req := Request.ofFlat v.request
  let declPairs ← evalMethodRoutingPairs v.spec req
  let declHeader := formatHeaderPairs declPairs
  let (staticDeclHeader, alphaHeader, tmplsJson) ←
    match v.spec.routing with
    | none => .ok (declHeader, declHeader, [])
    | some rule => do
      let groups ← compileRoutingRuleDeclOrder rule
      let tmpls ← compileTemplatesForJson rule.parameters
      let staticDecl := formatHeaderPairs (evalCompiledKeyGroupsPairs groups req)
      let alpha := formatHeaderPairs (evalCompiledKeyGroupsPairs (sortKeyGroupsAlpha groups) req)
      .ok (staticDecl, alpha, tmpls)
  return Lean.Json.mkObj [
    ("id", Lean.Json.str v.id),
    ("category", Lean.Json.str v.category.toSlug),
    ("description", Lean.Json.str v.description),
    ("theorem_ref", Lean.Json.str v.theoremRef),
    ("routing", routingRuleToJson v.spec.routing),
    ("http", httpRuleToJson v.spec.http),
    ("request", requestToJson v.request),
    ("expected_pairs", pairsToJson (sortPairsAlpha declPairs)),
    ("expected_decl_header", Lean.Json.str declHeader),
    ("expected_static_decl_header", Lean.Json.str staticDeclHeader),
    ("expected_alpha_header", Lean.Json.str alphaHeader),
    ("expect_header_sent", Lean.Json.bool (alphaHeader != "")),
    ("compiled_templates", Lean.Json.arr tmplsJson.toArray)
  ]

private def invalidVectorToJson (v : InvalidTemplateVector) : Except String Lean.Json :=
  match parseRoutingParameter v.field v.pathTemplate with
  | .ok _ => .error s!"negative vector {v.id} unexpectedly parsed successfully"
  | .error err =>
    .ok <| Lean.Json.mkObj [
      ("id", Lean.Json.str v.id),
      ("bug_ref", Lean.Json.str v.bugRef),
      ("description", Lean.Json.str v.description),
      ("field", Lean.Json.str v.field),
      ("path_template", Lean.Json.str v.pathTemplate),
      ("expected_error", Lean.Json.str err)
    ]

private def buildConformanceSuiteJson : Except String Lean.Json := do
  let extArr ← extractVectors.mapM extractVectorToJson
  let invArr ← invalidTemplateVectors.mapM invalidVectorToJson
  return Lean.Json.mkObj [
    ("schema_version", Lean.Json.str "1.0"),
    ("spec", Lean.Json.str "AIP-4222"),
    ("generator", Lean.Json.str "formal/aip4222/Aip4222/Vectors.lean"),
    ("extract_vectors", Lean.Json.arr extArr.toArray),
    ("invalid_template_vectors", Lean.Json.arr invArr.toArray)
  ]

private def runEmitConformance (outPath : Option String) : IO Unit := do
  let suite ← IO.ofExcept buildConformanceSuiteJson
  let rendered := Lean.Json.pretty suite 2 ++ "\n"
  match outPath with
  | none => IO.print rendered
  | some path => IO.FS.writeFile path rendered

private def runBatchJson : IO Unit := do
  let input ← (← IO.getStdin).readToEnd
  let parsed ← IO.ofExcept (Lean.Json.parse input)
  let arr ← IO.ofExcept parsed.getArr?
  let results ← IO.ofExcept (arr.toList.mapM evalCaseJson)
  IO.println (Lean.Json.compress (Lean.Json.arr results.toArray))

def main (args : List String) : IO Unit := do
  match args with
  | [] | ["smoke"] => runSmoke
  | ["--batch-json"] => runBatchJson
  | ["--emit-conformance"] => runEmitConformance none
  | ["--emit-conformance", outPath] => runEmitConformance (some outPath)
  | _ =>
    IO.eprintln s!"unknown arguments: {args}; expected [smoke|--batch-json|--emit-conformance [path]]"
    IO.Process.exit 2
