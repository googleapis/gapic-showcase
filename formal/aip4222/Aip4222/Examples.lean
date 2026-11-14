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

namespace Aip4222.Examples

/-- Canonical request payload modeled on `google/api/routing.proto` (plural `tables`). -/
def canonicalRequest : Request :=
  Request.ofFlat [
    ("table_name", "projects/proj_foo/instances/instance_bar/tables/table_baz"),
    ("app_profile_id", "profiles/prof_qux")
  ]

/-- Literal `table_name` from `google/api/routing.proto` line 58 (singular `table`). -/
def routingProtoLine58Request : Request :=
  Request.ofFlat [
    ("table_name", "projects/proj_foo/instances/instance_bar/table/table_baz"),
    ("app_profile_id", "profiles/prof_qux")
  ]

def ex1Rule : RoutingRule :=
  { parameters := [{ field := "app_profile_id" }] }

def ex2Rule : RoutingRule :=
  { parameters := [{ field := "app_profile_id", pathTemplate := "{routing_id=**}" }] }

def ex3aRule : RoutingRule :=
  { parameters := [{ field := "table_name", pathTemplate := "{table_name=projects/*/instances/*/**}" }] }

def ex3bRule : RoutingRule :=
  { parameters := [{ field := "table_name", pathTemplate := "{table_name=regions/*/zones/*/**}" }] }

def ex3cRule : RoutingRule :=
  { parameters := [
      { field := "table_name", pathTemplate := "{table_name=regions/*/zones/*/**}" },
      { field := "table_name", pathTemplate := "{table_name=projects/*/instances/*/**}" }
    ] }

def ex4Rule : RoutingRule :=
  { parameters := [{ field := "table_name", pathTemplate := "{routing_id=projects/*}/**" }] }

def ex5Rule : RoutingRule :=
  { parameters := [
      { field := "table_name", pathTemplate := "{routing_id=projects/*}/**" },
      { field := "table_name", pathTemplate := "{routing_id=projects/*/instances/*}/**" }
    ] }

def ex6aRule : RoutingRule :=
  { parameters := [
      { field := "table_name", pathTemplate := "{project_id=projects/*}/instances/*/**" },
      { field := "table_name", pathTemplate := "projects/*/{instance_id=instances/*}/**" }
    ] }

def ex6bRule : RoutingRule :=
  { parameters := [
      { field := "table_name", pathTemplate := "{project_id=projects/*}/**" },
      { field := "table_name", pathTemplate := "projects/*/{instance_id=instances/*}/**" }
    ] }

def ex7Rule : RoutingRule :=
  { parameters := [
      { field := "table_name", pathTemplate := "{project_id=projects/*}/**" },
      { field := "app_profile_id", pathTemplate := "{routing_id=**}" }
    ] }

def ex8Rule : RoutingRule :=
  { parameters := [
      { field := "table_name", pathTemplate := "{routing_id=projects/*}/**" },
      { field := "table_name", pathTemplate := "{routing_id=regions/*}/**" },
      { field := "app_profile_id", pathTemplate := "{routing_id=**}" }
    ] }

def ex9Rule : RoutingRule :=
  { parameters := [
      { field := "table_name", pathTemplate := "projects/*/{table_location=instances/*}/tables/*" },
      { field := "table_name", pathTemplate := "{table_location=regions/*/zones/*}/tables/*" },
      { field := "table_name", pathTemplate := "{routing_id=projects/*}/**" },
      { field := "app_profile_id", pathTemplate := "{routing_id=**}" },
      { field := "app_profile_id", pathTemplate := "profiles/{routing_id=*}" }
    ] }

private theorem evalRoutingRule_of_pairs (rule : RoutingRule) (req : Request)
    (pairs : List (String × String)) (encoded : String)
    (hPairs : evalRoutingRulePairs rule req = .ok pairs)
    (hFmt : formatHeaderPairs pairs = encoded) :
    evalRoutingRule rule req = .ok encoded := by
  unfold evalRoutingRule
  rw [hPairs]
  simp [Except.map, hFmt]

theorem example_1_canonical :
    evalRoutingRulePairs ex1Rule canonicalRequest = .ok [("app_profile_id", "profiles/prof_qux")] ∧
    evalRoutingRule ex1Rule canonicalRequest = .ok "app_profile_id=profiles%2Fprof_qux" :=
  ⟨by rfl, evalRoutingRule_of_pairs _ _ _ _ (by rfl) (by rfl)⟩

theorem example_2_canonical :
    evalRoutingRulePairs ex2Rule canonicalRequest = .ok [("routing_id", "profiles/prof_qux")] ∧
    evalRoutingRule ex2Rule canonicalRequest = .ok "routing_id=profiles%2Fprof_qux" :=
  ⟨by rfl, evalRoutingRule_of_pairs _ _ _ _ (by rfl) (by rfl)⟩

theorem example_3a_canonical :
    evalRoutingRulePairs ex3aRule canonicalRequest = .ok [("table_name", "projects/proj_foo/instances/instance_bar/tables/table_baz")] ∧
    evalRoutingRule ex3aRule canonicalRequest = .ok "table_name=projects%2Fproj_foo%2Finstances%2Finstance_bar%2Ftables%2Ftable_baz" :=
  ⟨by rfl, evalRoutingRule_of_pairs _ _ _ _ (by rfl) (by rfl)⟩

theorem example_3b_canonical :
    evalRoutingRulePairs ex3bRule canonicalRequest = .ok [] ∧
    evalRoutingRule ex3bRule canonicalRequest = .ok "" :=
  ⟨by rfl, evalRoutingRule_of_pairs _ _ _ _ (by rfl) (by rfl)⟩

theorem example_3c_canonical :
    evalRoutingRulePairs ex3cRule canonicalRequest = .ok [("table_name", "projects/proj_foo/instances/instance_bar/tables/table_baz")] ∧
    evalRoutingRule ex3cRule canonicalRequest = .ok "table_name=projects%2Fproj_foo%2Finstances%2Finstance_bar%2Ftables%2Ftable_baz" :=
  ⟨by rfl, evalRoutingRule_of_pairs _ _ _ _ (by rfl) (by rfl)⟩

theorem example_4_canonical :
    evalRoutingRulePairs ex4Rule canonicalRequest = .ok [("routing_id", "projects/proj_foo")] ∧
    evalRoutingRule ex4Rule canonicalRequest = .ok "routing_id=projects%2Fproj_foo" :=
  ⟨by rfl, evalRoutingRule_of_pairs _ _ _ _ (by rfl) (by rfl)⟩

theorem example_5_canonical :
    evalRoutingRulePairs ex5Rule canonicalRequest = .ok [("routing_id", "projects/proj_foo/instances/instance_bar")] ∧
    evalRoutingRule ex5Rule canonicalRequest = .ok "routing_id=projects%2Fproj_foo%2Finstances%2Finstance_bar" :=
  ⟨by rfl, evalRoutingRule_of_pairs _ _ _ _ (by rfl) (by rfl)⟩

theorem example_6a_canonical :
    evalRoutingRulePairs ex6aRule canonicalRequest = .ok [("project_id", "projects/proj_foo"), ("instance_id", "instances/instance_bar")] ∧
    evalRoutingRule ex6aRule canonicalRequest = .ok "project_id=projects%2Fproj_foo&instance_id=instances%2Finstance_bar" :=
  ⟨by rfl, evalRoutingRule_of_pairs _ _ _ _ (by rfl) (by rfl)⟩

theorem example_6b_canonical :
    evalRoutingRulePairs ex6bRule canonicalRequest = .ok [("project_id", "projects/proj_foo"), ("instance_id", "instances/instance_bar")] ∧
    evalRoutingRule ex6bRule canonicalRequest = .ok "project_id=projects%2Fproj_foo&instance_id=instances%2Finstance_bar" :=
  ⟨by rfl, evalRoutingRule_of_pairs _ _ _ _ (by rfl) (by rfl)⟩

theorem example_7_canonical :
    evalRoutingRulePairs ex7Rule canonicalRequest = .ok [("project_id", "projects/proj_foo"), ("routing_id", "profiles/prof_qux")] ∧
    evalRoutingRule ex7Rule canonicalRequest = .ok "project_id=projects%2Fproj_foo&routing_id=profiles%2Fprof_qux" :=
  ⟨by rfl, evalRoutingRule_of_pairs _ _ _ _ (by rfl) (by rfl)⟩

theorem example_8_canonical :
    evalRoutingRulePairs ex8Rule canonicalRequest = .ok [("routing_id", "profiles/prof_qux")] ∧
    evalRoutingRule ex8Rule canonicalRequest = .ok "routing_id=profiles%2Fprof_qux" :=
  ⟨by rfl, evalRoutingRule_of_pairs _ _ _ _ (by rfl) (by rfl)⟩

theorem example_9_canonical :
    evalRoutingRulePairs ex9Rule canonicalRequest = .ok [("table_location", "instances/instance_bar"), ("routing_id", "prof_qux")] ∧
    evalRoutingRule ex9Rule canonicalRequest = .ok "table_location=instances%2Finstance_bar&routing_id=prof_qux" :=
  ⟨by rfl, evalRoutingRule_of_pairs _ _ _ _ (by rfl) (by rfl)⟩

/-- All 9 canonical examples (12 sub-cases) from `google/api/routing.proto` match both raw pairs and encoded headers. -/
theorem examples_1_to_9_canonical :
    (evalRoutingRulePairs ex1Rule canonicalRequest = .ok [("app_profile_id", "profiles/prof_qux")] ∧
     evalRoutingRule ex1Rule canonicalRequest = .ok "app_profile_id=profiles%2Fprof_qux") ∧
    (evalRoutingRulePairs ex2Rule canonicalRequest = .ok [("routing_id", "profiles/prof_qux")] ∧
     evalRoutingRule ex2Rule canonicalRequest = .ok "routing_id=profiles%2Fprof_qux") ∧
    (evalRoutingRulePairs ex3aRule canonicalRequest = .ok [("table_name", "projects/proj_foo/instances/instance_bar/tables/table_baz")] ∧
     evalRoutingRule ex3aRule canonicalRequest = .ok "table_name=projects%2Fproj_foo%2Finstances%2Finstance_bar%2Ftables%2Ftable_baz") ∧
    (evalRoutingRulePairs ex3bRule canonicalRequest = .ok [] ∧
     evalRoutingRule ex3bRule canonicalRequest = .ok "") ∧
    (evalRoutingRulePairs ex3cRule canonicalRequest = .ok [("table_name", "projects/proj_foo/instances/instance_bar/tables/table_baz")] ∧
     evalRoutingRule ex3cRule canonicalRequest = .ok "table_name=projects%2Fproj_foo%2Finstances%2Finstance_bar%2Ftables%2Ftable_baz") ∧
    (evalRoutingRulePairs ex4Rule canonicalRequest = .ok [("routing_id", "projects/proj_foo")] ∧
     evalRoutingRule ex4Rule canonicalRequest = .ok "routing_id=projects%2Fproj_foo") ∧
    (evalRoutingRulePairs ex5Rule canonicalRequest = .ok [("routing_id", "projects/proj_foo/instances/instance_bar")] ∧
     evalRoutingRule ex5Rule canonicalRequest = .ok "routing_id=projects%2Fproj_foo%2Finstances%2Finstance_bar") ∧
    (evalRoutingRulePairs ex6aRule canonicalRequest = .ok [("project_id", "projects/proj_foo"), ("instance_id", "instances/instance_bar")] ∧
     evalRoutingRule ex6aRule canonicalRequest = .ok "project_id=projects%2Fproj_foo&instance_id=instances%2Finstance_bar") ∧
    (evalRoutingRulePairs ex6bRule canonicalRequest = .ok [("project_id", "projects/proj_foo"), ("instance_id", "instances/instance_bar")] ∧
     evalRoutingRule ex6bRule canonicalRequest = .ok "project_id=projects%2Fproj_foo&instance_id=instances%2Finstance_bar") ∧
    (evalRoutingRulePairs ex7Rule canonicalRequest = .ok [("project_id", "projects/proj_foo"), ("routing_id", "profiles/prof_qux")] ∧
     evalRoutingRule ex7Rule canonicalRequest = .ok "project_id=projects%2Fproj_foo&routing_id=profiles%2Fprof_qux") ∧
    (evalRoutingRulePairs ex8Rule canonicalRequest = .ok [("routing_id", "profiles/prof_qux")] ∧
     evalRoutingRule ex8Rule canonicalRequest = .ok "routing_id=profiles%2Fprof_qux") ∧
    (evalRoutingRulePairs ex9Rule canonicalRequest = .ok [("table_location", "instances/instance_bar"), ("routing_id", "prof_qux")] ∧
     evalRoutingRule ex9Rule canonicalRequest = .ok "table_location=instances%2Finstance_bar&routing_id=prof_qux") :=
  ⟨example_1_canonical, example_2_canonical, example_3a_canonical, example_3b_canonical,
   example_3c_canonical, example_4_canonical, example_5_canonical, example_6a_canonical,
   example_6b_canonical, example_7_canonical, example_8_canonical, example_9_canonical⟩

/-- Strict vs. loose templates from `google/api/routing.proto` Examples 6a & 6b when instance is absent. -/
theorem example_6a_vs_6b_without_instance :
    let req := Request.ofFlat [("table_name", "projects/proj_foo/tables/table_baz")]
    evalRoutingRule ex6aRule req = .ok "" ∧
    evalRoutingRule ex6bRule req = .ok "project_id=projects%2Fproj_foo" := by
  repeat (first | apply And.intro | rfl)

/-- Example 9 fallback to `table_name` when `app_profile_id` is empty. -/
theorem example_9_empty_app_profile_fallback :
    let req := Request.ofFlat [
      ("table_name", "projects/proj_foo/instances/instance_bar/tables/table_baz"),
      ("app_profile_id", "")
    ]
    evalRoutingRule ex9Rule req = .ok "table_location=instances%2Finstance_bar&routing_id=projects%2Fproj_foo" := by
  rfl

/--
Spec audit finding: `google/api/routing.proto` line 58 writes singular `table/table_baz`
while Example 9 (line 367) expects plural `tables/*`, causing `table_location` to be omitted.
-/
theorem routing_proto_line58_typo_exposed :
    evalRoutingRule ex9Rule routingProtoLine58Request = .ok "routing_id=prof_qux" := by
  rfl

/-- Hypothesis H1: Explicit empty `option (google.api.routing) = {}` always suppresses implicit HTTP routing. -/
theorem empty_routing_suppresses_http (http : Option HttpRule) (req : Request) :
    evalMethodRouting { routing := some { parameters := [] }, http := http } req = .ok "" := by
  rfl

/-- Hypothesis H2: Non-terminal `**` before a non-empty suffix (accepted by `librarian` test) is ill-formed. -/
theorem parser_rejects_mid_multi_wildcard :
    (parseRoutingParameter "default" "p/*/q/*/{routing=a/*/b/**}/s/*/u/*/v/*").isOk = false := by
  decide

private theorem parser_edge_1 :
    parseRoutingParameter "parent" "" =
      .ok { headerKey := "parent", prefixSegs := [], matchSegs := [.multiWildcard], suffixSegs := [] } := by rfl

private theorem parser_edge_2 :
    parseRoutingParameter "parent" "projects/{parent}" =
      .ok { headerKey := "parent", prefixSegs := [.lit "projects"], matchSegs := [.singleWildcard], suffixSegs := [] } := by rfl

private theorem parser_edge_3 :
    parseRoutingParameter "parent" "projects/{parent}/" =
      parseRoutingParameter "parent" "projects/{parent}" := by
  have h1 : parseRoutingParameter "parent" "projects/{parent}/" =
      .ok { headerKey := "parent", prefixSegs := [.lit "projects"], matchSegs := [.singleWildcard], suffixSegs := [] } := by rfl
  have h2 : parseRoutingParameter "parent" "projects/{parent}" =
      .ok { headerKey := "parent", prefixSegs := [.lit "projects"], matchSegs := [.singleWildcard], suffixSegs := [] } := by rfl
  exact h1.trans h2.symm

private theorem parser_edge_4 :
    parseRoutingParameter "table" "projects/{id=instances/*}/**/" =
      parseRoutingParameter "table" "projects/{id=instances/*}/**" := by
  have h1 : parseRoutingParameter "table" "projects/{id=instances/*}/**/" =
      .ok { headerKey := "id", prefixSegs := [.lit "projects"], matchSegs := [.lit "instances", .singleWildcard], suffixSegs := [.multiWildcard] } := by rfl
  have h2 : parseRoutingParameter "table" "projects/{id=instances/*}/**" =
      .ok { headerKey := "id", prefixSegs := [.lit "projects"], matchSegs := [.lit "instances", .singleWildcard], suffixSegs := [.multiWildcard] } := by rfl
  exact h1.trans h2.symm

private theorem parser_edge_5 :
    parseRoutingParameter "sub.field" "" =
      .ok { headerKey := "sub.field", prefixSegs := [], matchSegs := [.multiWildcard], suffixSegs := [] } := by rfl

private theorem parser_edge_6 : (parseRoutingParameter "default" "projects/{parent}//").isOk = false := by decide
private theorem parser_edge_7 : (parseRoutingParameter "default" "projects/*/{routing_id=**}foo").isOk = false := by decide
private theorem parser_edge_8 : (parseRoutingParameter "default" "projects/**/{a}").isOk = false := by decide
private theorem parser_edge_9 : (parseRoutingParameter "default" "{**}").isOk = false := by decide
private theorem parser_edge_10 : (parseRoutingParameter "default" "projects/{id=*}/{other=*}").isOk = false := by decide
private theorem parser_edge_11 : (parseRoutingParameter "" "").isOk = false := by decide

/-- Additional AIP-4222 parser conformance cases (shorthands, trailing `/`, nested fields, invalid templates). -/
theorem parser_edge_cases :
    parseRoutingParameter "parent" "" =
      .ok { headerKey := "parent", prefixSegs := [], matchSegs := [.multiWildcard], suffixSegs := [] } ∧
    parseRoutingParameter "parent" "projects/{parent}" =
      .ok { headerKey := "parent", prefixSegs := [.lit "projects"], matchSegs := [.singleWildcard], suffixSegs := [] } ∧
    parseRoutingParameter "parent" "projects/{parent}/" =
      parseRoutingParameter "parent" "projects/{parent}" ∧
    parseRoutingParameter "table" "projects/{id=instances/*}/**/" =
      parseRoutingParameter "table" "projects/{id=instances/*}/**" ∧
    parseRoutingParameter "sub.field" "" =
      .ok { headerKey := "sub.field", prefixSegs := [], matchSegs := [.multiWildcard], suffixSegs := [] } ∧
    (parseRoutingParameter "default" "projects/{parent}//").isOk = false ∧
    (parseRoutingParameter "default" "projects/*/{routing_id=**}foo").isOk = false ∧
    (parseRoutingParameter "default" "projects/**/{a}").isOk = false ∧
    (parseRoutingParameter "default" "{**}").isOk = false ∧
    (parseRoutingParameter "default" "projects/{id=*}/{other=*}").isOk = false ∧
    (parseRoutingParameter "" "").isOk = false :=
  ⟨parser_edge_1, parser_edge_2, parser_edge_3, parser_edge_4, parser_edge_5,
   parser_edge_6, parser_edge_7, parser_edge_8, parser_edge_9, parser_edge_10, parser_edge_11⟩

/-- RFC 6570 §3.2.2 encoding edge cases: `%20` (not `+`), unreserved `-._~`, `%25`, multi-byte UTF-8, dotted keys. -/
theorem encoding_rfc6570_edge_cases :
    percentEncode "Hello World! 50% ~_.-" = "Hello%20World%21%2050%25%20~_.-" ∧
    percentEncode "café €" = "caf%C3%A9%20%E2%82%AC" ∧
    formatHeaderPairs [("sub.field", "a/b c"), ("~k", "100%")] =
      "sub.field=a%2Fb%20c&~k=100%25" :=
  ⟨by rfl, by rfl, by rfl⟩

/-- Implicit `google.api.http` fallback, `additional_bindings` deduplication, and explicit override. -/
theorem http_implicit_and_precedence_cases :
    let httpRule : HttpRule := {
      bindings := [
        [{ fieldPath := ["parent"] }, { fieldPath := ["missing"] }],
        [{ fieldPath := ["parent"] }, { fieldPath := ["topic", "name"] }]
      ]
    }
    let req := Request.ofFlat [("parent", "projects/p1"), ("topic.name", "t/1"), ("unused", "x")]
    evalMethodRouting { routing := none, http := some httpRule } req =
      .ok "parent=projects%2Fp1&topic.name=t%2F1" ∧
    evalMethodRouting { routing := some ex1Rule, http := some httpRule } canonicalRequest =
      .ok "app_profile_id=profiles%2Fprof_qux" :=
  ⟨by rfl, by rfl⟩

private def colonTestRule : RoutingRule :=
  { parameters := [{ field := "name", pathTemplate := "{id=projects/foo}/**" }] }

private def tableTestRule : RoutingRule :=
  { parameters := [{ field := "name", pathTemplate := "projects/{id=*}/tables/**" }] }

private def nestedTestRule : RoutingRule :=
  { parameters := [{ field := "book.author.name", pathTemplate := "authors/{author=*}" }] }

private theorem zero_seg_1 :
    evalRoutingRule ex4Rule (Request.ofFlat [("table_name", "projects/proj_foo")]) =
      .ok "routing_id=projects%2Fproj_foo" := by rfl

private theorem zero_seg_2 :
    evalRoutingRule ex4Rule (Request.ofFlat [("table_name", "projects/proj_foo/")]) =
      .ok "routing_id=projects%2Fproj_foo" := by rfl

private theorem colon_1 :
    evalRoutingRule colonTestRule (Request.ofFlat [("name", "projects/foo:verb")]) =
      .ok "id=projects%2Ffoo" := by rfl

private theorem colon_2 :
    evalRoutingRule colonTestRule (Request.ofFlat [("name", "projects/foo-bar")]) =
      .ok "" := by rfl

private theorem colon_3 :
    evalRoutingRule tableTestRule (Request.ofFlat [("name", "projects/p1/tables:batch")]) =
      .ok "id=p1" := by rfl

private theorem colon_4 :
    evalRoutingRule tableTestRule (Request.ofFlat [("name", "projects/p1/tables-v2")]) =
      .ok "" := by rfl

private theorem nested_1 :
    evalRoutingRule nestedTestRule (Request.ofFlat [("book.author.name", "authors/ada")]) =
      .ok "author=ada" := by rfl

private theorem unset_1 :
    evalRoutingRule ex1Rule (Request.ofFlat []) = .ok "" := by rfl

/-- Zero-segment `**`, `[:/]` colon delimiter vs `-`, dotted field runtime lookup, and unset omission. -/
theorem trailing_slash_and_colon_and_zero_seg_cases :
    evalRoutingRule ex4Rule (Request.ofFlat [("table_name", "projects/proj_foo")]) =
      .ok "routing_id=projects%2Fproj_foo" ∧
    evalRoutingRule ex4Rule (Request.ofFlat [("table_name", "projects/proj_foo/")]) =
      .ok "routing_id=projects%2Fproj_foo" ∧
    evalRoutingRule colonTestRule (Request.ofFlat [("name", "projects/foo:verb")]) =
      .ok "id=projects%2Ffoo" ∧
    evalRoutingRule colonTestRule (Request.ofFlat [("name", "projects/foo-bar")]) =
      .ok "" ∧
    evalRoutingRule tableTestRule (Request.ofFlat [("name", "projects/p1/tables:batch")]) =
      .ok "id=p1" ∧
    evalRoutingRule tableTestRule (Request.ofFlat [("name", "projects/p1/tables-v2")]) =
      .ok "" ∧
    evalRoutingRule nestedTestRule (Request.ofFlat [("book.author.name", "authors/ada")]) =
      .ok "author=ada" ∧
    evalRoutingRule ex1Rule (Request.ofFlat []) = .ok "" :=
  ⟨zero_seg_1, zero_seg_2, colon_1, colon_2, colon_3, colon_4, nested_1, unset_1⟩

end Aip4222.Examples
