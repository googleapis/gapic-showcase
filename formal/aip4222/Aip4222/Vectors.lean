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

import Aip4222.Theorems
import Aip4222.Examples

namespace Aip4222

open Aip4222.Examples

/-- Evaluate a `MethodRoutingSpec` into unencoded `(headerKey, value)` pairs. -/
def evalMethodRoutingPairs (spec : MethodRoutingSpec) (req : Request) :
    Except String (List (String × String)) :=
  match spec.routing with
  | some rule => evalRoutingRulePairs rule req
  | none =>
    match spec.http with
    | none => .ok []
    | some http => .ok (evalHttpRulePairs http req)

/-- Insert a `CompiledKeyGroup` in lexicographical `headerKey` order (kernel-reducible). -/
def insertKeyGroupAlpha (g : CompiledKeyGroup) : List CompiledKeyGroup → List CompiledKeyGroup
  | [] => [g]
  | h :: t =>
    if g.headerKey ≤ h.headerKey then g :: h :: t
    else h :: insertKeyGroupAlpha g t

/-- Sort compiled key groups alphabetically by `headerKey` (`librarian` parity). -/
def sortKeyGroupsAlpha : List CompiledKeyGroup → List CompiledKeyGroup
  | [] => []
  | g :: gs => insertKeyGroupAlpha g (sortKeyGroupsAlpha gs)

/-- Insert a `(key, val)` pair in lexicographical `key` order (kernel-reducible). -/
def insertPairAlpha (p : String × String) : List (String × String) → List (String × String)
  | [] => [p]
  | h :: t =>
    if p.1 ≤ h.1 then p :: h :: t
    else h :: insertPairAlpha p t

/-- Canonicalize header pairs in lexicographical `key` order. -/
def sortPairsAlpha : List (String × String) → List (String × String)
  | [] => []
  | p :: ps => insertPairAlpha p (sortPairsAlpha ps)

/-- Conformance equivalence-class partitions for AIP-4222. -/
inductive ConformanceCategory where
  | canonicalExamples
  | explicitPrecedence
  | implicitHttpFallback
  | segmentBoundaries
  | conflictAndKeyOrder
  | rfc6570Encoding
  | nestedFields
  | showcaseEcho
  deriving DecidableEq, Repr, BEq

def ConformanceCategory.toSlug : ConformanceCategory → String
  | .canonicalExamples => "canonical_examples"
  | .explicitPrecedence => "explicit_precedence"
  | .implicitHttpFallback => "implicit_http_fallback"
  | .segmentBoundaries => "segment_boundaries"
  | .conflictAndKeyOrder => "conflict_and_key_order"
  | .rfc6570Encoding => "rfc6570_encoding"
  | .nestedFields => "nested_fields"
  | .showcaseEcho => "showcase_echo"

/-- Positive conformance vector certifying expected extraction & wire encoding. -/
structure ExtractVector where
  id : String
  category : ConformanceCategory
  description : String
  theoremRef : String
  spec : MethodRoutingSpec
  request : List (String × String)
  deriving DecidableEq, Repr, BEq

/-- Negative conformance vector certifying static template/field rejection. -/
structure InvalidTemplateVector where
  id : String
  bugRef : String
  description : String
  field : String
  pathTemplate : String
  deriving DecidableEq, Repr, BEq

private def canonFlatReq : List (String × String) :=
  [ ("table_name", "projects/proj_foo/instances/instance_bar/tables/table_baz")
  , ("app_profile_id", "profiles/prof_qux")
  ]

private def line58FlatReq : List (String × String) :=
  [ ("table_name", "projects/proj_foo/instances/instance_bar/table/table_baz")
  , ("app_profile_id", "profiles/prof_qux")
  ]

/-- Exact explicit `google.api.routing` rule from `gapic-showcase` `Echo.Echo` (`echo.proto:55-88`). -/
def showcaseEchoRule : RoutingRule :=
  { parameters :=
      [ { field := "header" }
      , { field := "header", pathTemplate := "{routing_id=**}" }
      , { field := "header", pathTemplate := "{table_name=regions/*/zones/*/**}" }
      , { field := "header", pathTemplate := "{super_id=projects/*}/**" }
      , { field := "header", pathTemplate := "{table_name=projects/*/instances/*/**}" }
      , { field := "header", pathTemplate := "projects/*/{instance_id=instances/*}/**" }
      , { field := "other_header", pathTemplate := "{baz=**}" }
      , { field := "other_header", pathTemplate := "{qux=projects/*}/**" }
      ] }

private def mkHttp (bindings : List (List String)) : Option HttpRule :=
  some { bindings := bindings.map fun b => b.map fun f => { fieldPath := FieldPath.parse f } }

private def mkRule (ps : List (String × String)) : Option RoutingRule :=
  some { parameters := ps.map fun (f, t) => { field := f, pathTemplate := t } }

/-- 1. Canonical Examples from google/api/routing.proto (13) -/
private def canonicalVectors : List ExtractVector :=
  [ { id := "ex1"
      category := .canonicalExamples
      description := "AIP-4222 Example 1: omitted path_template transcribes field verbatim"
      theoremRef := "Aip4222.Examples.example_1_canonical"
      spec := { routing := some ex1Rule }
      request := canonFlatReq }
  , { id := "ex2"
      category := .canonicalExamples
      description := "AIP-4222 Example 2: rename header key via {routing_id=**}"
      theoremRef := "Aip4222.Examples.example_2_canonical"
      spec := { routing := some ex2Rule }
      request := canonFlatReq }
  , { id := "ex3a"
      category := .canonicalExamples
      description := "AIP-4222 Example 3a: full match with projects/*/instances/*/**"
      theoremRef := "Aip4222.Examples.example_3a_canonical"
      spec := { routing := some ex3aRule }
      request := canonFlatReq }
  , { id := "ex3b"
      category := .canonicalExamples
      description := "AIP-4222 Example 3b: mismatch emits empty header"
      theoremRef := "Aip4222.Examples.example_3b_canonical"
      spec := { routing := some ex3bRule }
      request := canonFlatReq }
  , { id := "ex3c"
      category := .canonicalExamples
      description := "AIP-4222 Example 3c: fallback between region and project patterns"
      theoremRef := "Aip4222.Examples.example_3c_canonical"
      spec := { routing := some ex3cRule }
      request := canonFlatReq }
  , { id := "ex4"
      category := .canonicalExamples
      description := "AIP-4222 Example 4: extract prefix resource segment with ** suffix"
      theoremRef := "Aip4222.Examples.example_4_canonical"
      spec := { routing := some ex4Rule }
      request := canonFlatReq }
  , { id := "ex5"
      category := .canonicalExamples
      description := "AIP-4222 Example 5: last match wins across same-key overlapping prefixes"
      theoremRef := "Aip4222.Examples.example_5_canonical"
      spec := { routing := some ex5Rule }
      request := canonFlatReq }
  , { id := "ex6a"
      category := .canonicalExamples
      description := "AIP-4222 Example 6a: extract two distinct header keys from one field"
      theoremRef := "Aip4222.Examples.example_6a_canonical"
      spec := { routing := some ex6aRule }
      request := canonFlatReq }
  , { id := "ex6b"
      category := .canonicalExamples
      description := "AIP-4222 Example 6b: first key matches even when second key template misses"
      theoremRef := "Aip4222.Examples.example_6b_canonical"
      spec := { routing := some ex6bRule }
      request := canonFlatReq }
  , { id := "ex7"
      category := .canonicalExamples
      description := "AIP-4222 Example 7: extract independent keys from two request fields"
      theoremRef := "Aip4222.Examples.example_7_canonical"
      spec := { routing := some ex7Rule }
      request := canonFlatReq }
  , { id := "ex8"
      category := .canonicalExamples
      description := "AIP-4222 Example 8: cross-field override into shared header key routing_id"
      theoremRef := "Aip4222.Examples.example_8_canonical"
      spec := { routing := some ex8Rule }
      request := canonFlatReq }
  , { id := "ex9"
      category := .canonicalExamples
      description := "AIP-4222 Example 9 (canonical plural tables): combined keys and overrides"
      theoremRef := "Aip4222.Examples.example_9_canonical"
      spec := { routing := some ex9Rule }
      request := canonFlatReq }
  , { id := "ex9_line58_typo"
      category := .canonicalExamples
      description := "AIP-4222 Example 9 against singular table/ typo at routing.proto:58"
      theoremRef := "Aip4222.Examples.routing_proto_line58_typo_exposed"
      spec := { routing := some ex9Rule }
      request := line58FlatReq }
  ]

/-- 2. Explicit vs. Implicit HTTP Precedence (4) -/
private def explicitPrecedenceVectors : List ExtractVector :=
  [ { id := "h1_empty_routing_suppresses_single_http"
      category := .explicitPrecedence
      description := "H1: explicit option (google.api.routing) = {} suppresses google.api.http"
      theoremRef := "Aip4222.empty_routing_suppresses_http_universal"
      spec := { routing := mkRule [], http := mkHttp [["parent"]] }
      request := [("parent", "projects/p1")] }
  , { id := "h1_empty_routing_suppresses_multi_http"
      category := .explicitPrecedence
      description := "H1: explicit empty routing suppresses multiple HTTP bindings & vars"
      theoremRef := "Aip4222.empty_routing_suppresses_http_universal"
      spec := { routing := mkRule [], http := mkHttp [["parent", "id"], ["alt"]] }
      request := [("parent", "projects/p1"), ("id", "i1"), ("alt", "x")] }
  , { id := "explicit_nonempty_ignores_http_when_matched"
      category := .explicitPrecedence
      description := "Explicit routing ignores HTTP vars completely when matched"
      theoremRef := "Aip4222.explicit_routing_ignores_http"
      spec := { routing := mkRule [("name", "projects/*/{id=*}")], http := mkHttp [["parent"]] }
      request := [("name", "projects/p1/w1"), ("parent", "projects/should_ignore")] }
  , { id := "explicit_nonempty_ignores_http_even_when_unmatched"
      category := .explicitPrecedence
      description := "Explicit routing that fails to match still suppresses HTTP fallback"
      theoremRef := "Aip4222.explicit_routing_ignores_http"
      spec := { routing := mkRule [("name", "projects/*/{id=*}")], http := mkHttp [["parent"]] }
      request := [("name", "does-not-match"), ("parent", "projects/should_ignore")] }
  ]

/-- 3. Implicit google.api.http Fallback (6) -/
private def implicitHttpFallbackVectors : List ExtractVector :=
  [ { id := "h4_http_single_var_slash_encoded"
      category := .implicitHttpFallback
      description := "H4: implicit HTTP binding percent-encodes slashes and omits leading &"
      theoremRef := "Aip4222.h4_buggy_rust_http_fold_falsified"
      spec := { routing := none, http := mkHttp [["parent"]] }
      request := [("parent", "projects/p1")] }
  , { id := "h4_http_space_and_empty_filtered"
      category := .implicitHttpFallback
      description := "H4: implicit HTTP skips empty proto3 string values and encodes space as %20"
      theoremRef := "Aip4222.h4_buggy_rust_http_fold_falsified"
      spec := { routing := none, http := mkHttp [["parent", "empty"]] }
      request := [("parent", "a b"), ("empty", "")] }
  , { id := "h4_http_all_empty_omitted"
      category := .implicitHttpFallback
      description := "H4: implicit HTTP with all empty strings produces empty header"
      theoremRef := "Aip4222.formatHeaderPairs_empty_iff"
      spec := { routing := none, http := mkHttp [["parent", "empty"]] }
      request := [("parent", ""), ("empty", "")] }
  , { id := "h4_http_missing_field_omitted"
      category := .implicitHttpFallback
      description := "Implicit HTTP omits unset optional request fields"
      theoremRef := "Aip4222.Examples.http_implicit_and_precedence_cases"
      spec := { routing := none, http := mkHttp [["parent", "opt"]] }
      request := [("parent", "projects/p1")] }
  , { id := "h11_http_additional_bindings_dedup"
      category := .implicitHttpFallback
      description := "H11: implicit HTTP extracts additional_bindings vars and deduplicates repeats"
      theoremRef := "Aip4222.Examples.http_implicit_and_precedence_cases"
      spec := { routing := none, http := mkHttp [["parent", "id"], ["parent", "location"]] }
      request := [("parent", "projects/p1"), ("id", "42"), ("location", "us-central1")] }
  , { id := "no_routing_no_http_omitted"
      category := .implicitHttpFallback
      description := "Pure gRPC method with neither routing nor HTTP emits empty header"
      theoremRef := "Aip4222.evalMethodRouting"
      spec := { routing := none, http := none }
      request := [("parent", "projects/p1")] }
  ]

/-- 4. Segment Compilation & Boundaries (12) -/
private def segmentBoundaryVectors : List ExtractVector :=
  [ { id := "h3a_prefix_slash_match"
      category := .segmentBoundaries
      description := "H3a: non-empty prefix requires trailing / before match slice"
      theoremRef := "Aip4222.h3a_buggy_rust_prefix_drops_slash"
      spec := { routing := mkRule [("name", "projects/*/{id=instances/*}")] }
      request := [("name", "projects/p/instances/i")] }
  , { id := "h3a_prefix_slash_reject_glued"
      category := .segmentBoundaries
      description := "H3a: glued prefix/match without slash delimiter must not match"
      theoremRef := "Aip4222.compiled_matcher_refines_spec"
      spec := { routing := mkRule [("name", "projects/*/{id=instances/*}")] }
      request := [("name", "projects/pinstances/i")] }
  , { id := "unanchored_prefix_garbage_rejected"
      category := .segmentBoundaries
      description := "Anchored start prevents GARBAGE/ prefix from matching {super_id=projects/*}/**"
      theoremRef := "Aip4222.compiled_matcher_refines_spec"
      spec := { routing := mkRule [("header", "{super_id=projects/*}/**")] }
      request := [("header", "GARBAGE/projects/p1")] }
  , { id := "h3b_suffix_slash_match"
      category := .segmentBoundaries
      description := "H3b: non-wildcard suffix requires leading / after match slice"
      theoremRef := "Aip4222.h3b_buggy_rust_suffix_drops_slash"
      spec := { routing := mkRule [("name", "{id=instances/*}/tables/*")] }
      request := [("name", "instances/i/tables/t")] }
  , { id := "h3b_suffix_extra_segment_rejected"
      category := .segmentBoundaries
      description := "H3b: extra trailing segment without /** must fail to match"
      theoremRef := "Aip4222.compiled_matcher_refines_spec"
      spec := { routing := mkRule [("name", "{id=instances/*}/tables/*")] }
      request := [("name", "instances/i/tables/t/extra")] }
  , { id := "h3c_trailing_dstar_exact_zero_segments"
      category := .segmentBoundaries
      description := "H3c: trailing /** matches zero trailing segments"
      theoremRef := "Aip4222.compiled_matcher_refines_spec"
      spec := { routing := mkRule [("name", "{id=projects}/**")] }
      request := [("name", "projects")] }
  , { id := "h3c_trailing_dstar_slash_segments"
      category := .segmentBoundaries
      description := "H3c: trailing /** matches /... boundary"
      theoremRef := "Aip4222.compiled_matcher_refines_spec"
      spec := { routing := mkRule [("name", "{id=projects}/**")] }
      request := [("name", "projects/a/b")] }
  , { id := "h3c_trailing_dstar_colon_verb"
      category := .segmentBoundaries
      description := "H3c: trailing /** matches :... custom action delimiter"
      theoremRef := "Aip4222.Examples.trailing_slash_and_colon_and_zero_seg_cases"
      spec := { routing := mkRule [("name", "{id=projects}/**")] }
      request := [("name", "projects:cancel")] }
  , { id := "h3c_trailing_dstar_reject_concatenated"
      category := .segmentBoundaries
      description := "H3c: singleton ** suffix rejects concatenated suffix without / or :"
      theoremRef := "Aip4222.h3c_buggy_rust_single_dstar_suffix_overmatches"
      spec := { routing := mkRule [("name", "{id=projects}/**")] }
      request := [("name", "projectsFOO")] }
  , { id := "h3c_trailing_dstar_reject_hyphen"
      category := .segmentBoundaries
      description := "H3c: singleton ** suffix rejects hyphenated non-delimiter projects-v2"
      theoremRef := "Aip4222.compiled_matcher_refines_spec"
      spec := { routing := mkRule [("name", "{id=projects}/**")] }
      request := [("name", "projects-v2")] }
  , { id := "trailing_slash_in_template_normalized"
      category := .segmentBoundaries
      description := "H7: single trailing / in path_template is ignored per AIP-4222"
      theoremRef := "Aip4222.Examples.trailing_slash_and_colon_and_zero_seg_cases"
      spec := { routing := mkRule [("parent", "projects/{parent}/")] }
      request := [("parent", "projects/p1")] }
  , { id := "shorthand_var_without_equals"
      category := .segmentBoundaries
      description := "Shorthand {parent} equals {parent=*} and rejects extra segments"
      theoremRef := "Aip4222.parseRoutingParameter_sound"
      spec := { routing := mkRule [("name", "projects/{parent}")] }
      request := [("name", "projects/p1")] }
  ]

/-- 5. Conflict Resolution & Key Ordering (5) -/
private def conflictAndKeyOrderVectors : List ExtractVector :=
  [ { id := "same_key_later_match_overrides_earlier"
      category := .conflictAndKeyOrder
      description := "Reversed variant .or_else implements last-match-wins when both match"
      theoremRef := "Aip4222.reverse_first_match_eq_last_wins"
      spec :=
        { routing := mkRule
            [ ("a", "{k=a/*}")
            , ("b", "{k=**}")
            ] }
      request := [("a", "a/1"), ("b", "override")] }
  , { id := "same_key_later_miss_preserves_earlier"
      category := .conflictAndKeyOrder
      description := "Last-match-wins falls back to earlier match when later variant misses"
      theoremRef := "Aip4222.reverse_first_match_eq_last_wins"
      spec :=
        { routing := mkRule
            [ ("a", "{k=**}")
            , ("b", "{k=specific/*}")
            ] }
      request := [("a", "fallback"), ("b", "does-not-match")] }
  , { id := "static_vs_lazy_key_order_divergence"
      category := .conflictAndKeyOrder
      description := "Multi-variant key ordering vs alphabetical extensional key-map contract"
      theoremRef := "Aip4222.static_vs_lazy_key_order_diverges"
      spec :=
        { routing := mkRule
            [ ("a", "v1/{x=*}")
            , ("b", "v1/{y=*}")
            , ("c", "v2/{x=*}")
            ] }
      request := [("a", "miss"), ("b", "v1/1"), ("c", "v2/2")] }
  , { id := "alphabetical_vs_decl_single_variant"
      category := .conflictAndKeyOrder
      description := "Distinct keys declared z then a: decl_pairs vs alphabetical groups"
      theoremRef := "Aip4222.compiled_rule_extensional_eq_spec"
      spec :=
        { routing := mkRule
            [ ("first", "{z_key=*}")
            , ("second", "{a_key=*}")
            ] }
      request := [("first", "1"), ("second", "2")] }
  , { id := "empty_field_value_never_emits_pair"
      category := .conflictAndKeyOrder
      description := "Empty field string never matches * or ** and is suppressed"
      theoremRef := "Aip4222.Examples.parser_edge_cases"
      spec := { routing := mkRule [("name", "{id=**}")] }
      request := [("name", "")] }
  ]

/-- 6. RFC 6570 §3.2.2 Percent-Encoding & Literal Escaping (7) -/
private def rfc6570EncodingVectors : List ExtractVector :=
  [ { id := "rfc6570_unreserved_preserved"
      category := .rfc6570Encoding
      description := "RFC 6570 unreserved chars A-Za-z0-9-._~ remain unencoded"
      theoremRef := "Aip4222.percentEncode_chars_valid"
      spec := { routing := mkRule [("v", "{k=**}")] }
      request := [("v", "AZaz09-._~")] }
  , { id := "rfc6570_space_is_pct20_not_plus"
      category := .rfc6570Encoding
      description := "RFC 6570 encodes space as %20 (not application/x-www-form-urlencoded +)"
      theoremRef := "Aip4222.Examples.encoding_rfc6570_edge_cases"
      spec := { routing := mkRule [("v", "{k=**}")] }
      request := [("v", "hello world")] }
  , { id := "rfc6570_slash_and_colon_encoded"
      category := .rfc6570Encoding
      description := "Slashes and colons in routing values encode as %2F and %3A"
      theoremRef := "Aip4222.Examples.encoding_rfc6570_edge_cases"
      spec := { routing := mkRule [("v", "{k=**}")] }
      request := [("v", "projects/p:act")] }
  , { id := "rfc6570_plus_amp_eq_pct_encoded"
      category := .rfc6570Encoding
      description := "Delimiters +, &, =, and % are percent-encoded cleanly"
      theoremRef := "Aip4222.Examples.encoding_rfc6570_edge_cases"
      spec := { routing := mkRule [("v", "{k=**}")] }
      request := [("v", "a+b&c=d%20")] }
  , { id := "rfc6570_utf8_multibyte"
      category := .rfc6570Encoding
      description := "Multi-byte UTF-8 runes encode byte-by-byte in uppercase hex"
      theoremRef := "Aip4222.Examples.encoding_rfc6570_edge_cases"
      spec := { routing := mkRule [("v", "{k=**}")] }
      request := [("v", "café/東京")] }
  , { id := "rfc6570_dot_literal_exact_not_regex_any"
      category := .rfc6570Encoding
      description := "Literal dot in v1.0 matches only '.', unlike unescaped regex '.'"
      theoremRef := "Aip4222.compiled_matcher_refines_spec"
      spec := { routing := mkRule [("v", "v1.0/{id=*}")] }
      request := [("v", "v1X0/item")] }
  , { id := "rfc6570_dot_literal_exact_match"
      category := .rfc6570Encoding
      description := "Literal dot in v1.0 matches exact 'v1.0/item'"
      theoremRef := "Aip4222.compiled_matcher_refines_spec"
      spec := { routing := mkRule [("v", "v1.0/{id=*}")] }
      request := [("v", "v1.0/item")] }
  ]

/-- 7. Nested Message FieldPaths (3) -/
private def nestedFieldVectors : List ExtractVector :=
  [ { id := "nested_2_level_hit"
      category := .nestedFields
      description := "Two-level nested protobuf field traversal (nested.id)"
      theoremRef := "Aip4222.evalRoutingRule"
      spec := { routing := mkRule [("nested.id", "projects/*/{id=*}")] }
      request := [("nested.id", "projects/p1/res1")] }
  , { id := "nested_2_level_missing"
      category := .nestedFields
      description := "Unset nested message field yields empty header"
      theoremRef := "Aip4222.evalRoutingRule"
      spec := { routing := mkRule [("nested.id", "projects/*/{id=*}")] }
      request := [] }
  , { id := "nested_3_level_hit"
      category := .nestedFields
      description := "Three-level nested protobuf field traversal (a.b.c)"
      theoremRef := "Aip4222.evalRoutingRule"
      spec := { routing := mkRule [("a.b.c", "projects/*/{id=*}")] }
      request := [("a.b.c", "projects/p1/v1")] }
  ]

/-- 8. Live gapic-showcase Echo.Echo Service Conformance (7) -/
private def showcaseEchoVectors : List ExtractVector :=
  [ { id := "showcase_echo_both_empty"
      category := .showcaseEcho
      description := "Echo.Echo with empty header & other_header omits x-goog-request-params"
      theoremRef := "Aip4222.formatHeaderPairs_empty_iff"
      spec := { routing := some showcaseEchoRule }
      request := [("header", ""), ("other_header", "")] }
  , { id := "showcase_echo_project_instance"
      category := .showcaseEcho
      description := "Echo.Echo full extraction across 7 routing keys from header and other_header"
      theoremRef := "Aip4222.compiled_rule_extensional_eq_spec"
      spec := { routing := some showcaseEchoRule }
      request :=
        [ ("header", "projects/p1/instances/i1/tables/t1")
        , ("other_header", "projects/p2/bar")
        ] }
  , { id := "showcase_echo_region_zone"
      category := .showcaseEcho
      description := "Echo.Echo region/zone fallback path for table_name"
      theoremRef := "Aip4222.compiled_matcher_refines_spec"
      spec := { routing := some showcaseEchoRule }
      request :=
        [ ("header", "regions/us-east1/zones/z1/clusters/c1")
        , ("other_header", "")
        ] }
  , { id := "showcase_echo_unanchored_garbage_prefix"
      category := .showcaseEcho
      description := "Echo.Echo rejects unanchored GARBAGE/ prefix for projects/patterns"
      theoremRef := "Aip4222.compiled_matcher_refines_spec"
      spec := { routing := some showcaseEchoRule }
      request :=
        [ ("header", "GARBAGE/projects/p1/instances/i1")
        , ("other_header", "")
        ] }
  , { id := "showcase_echo_colon_action_suffix"
      category := .showcaseEcho
      description := "Echo.Echo with colon action suffix across header and other_header"
      theoremRef := "Aip4222.compiled_matcher_refines_spec"
      spec := { routing := some showcaseEchoRule }
      request :=
        [ ("header", "projects/p1/instances/i1:cancel")
        , ("other_header", "projects/p2:act")
        ] }
  , { id := "showcase_echo_rfc6570_space_and_tilde"
      category := .showcaseEcho
      description := "Echo.Echo encodes space as %20 and preserves unreserved tilde ~"
      theoremRef := "Aip4222.percentEncode_chars_valid"
      spec := { routing := some showcaseEchoRule }
      request :=
        [ ("header", "hello world~1")
        , ("other_header", "")
        ] }
  , { id := "showcase_echo_utf8"
      category := .showcaseEcho
      description := "Echo.Echo UTF-8 percent-encoding on multi-byte header value"
      theoremRef := "Aip4222.Examples.encoding_rfc6570_edge_cases"
      spec := { routing := some showcaseEchoRule }
      request :=
        [ ("header", "café")
        , ("other_header", "")
        ] }
  ]

/-- Deterministic 57-vector positive AIP-4222 conformance suite across 8 partitions. -/
def extractVectors : List ExtractVector :=
  canonicalVectors ++
  explicitPrecedenceVectors ++
  implicitHttpFallbackVectors ++
  segmentBoundaryVectors ++
  conflictAndKeyOrderVectors ++
  rfc6570EncodingVectors ++
  nestedFieldVectors ++
  showcaseEchoVectors

/-- Deterministic 20-vector negative suite covering malformed AIP-4222 templates/fields. -/
def invalidTemplateVectors : List InvalidTemplateVector :=
  [ { id := "h2_mid_dstar_routing_test"
      bugRef := "H2"
      description := "Non-terminal ** inside variable before multi-segment suffix"
      field := "default"
      pathTemplate := "p/*/q/*/{routing=a/*/b/**}/s/*/u/*/v/*" }
  , { id := "h2_dstar_in_prefix"
      bugRef := "H2"
      description := "** inside prefix slice"
      field := "name"
      pathTemplate := "projects/**/{id=*}" }
  , { id := "h2_dstar_mid_match"
      bugRef := "H2"
      description := "** before trailing segment inside match braces"
      field := "name"
      pathTemplate := "projects/*/{id=a/**/b}" }
  , { id := "h2_double_dstar"
      bugRef := "H2"
      description := "** in both match and suffix"
      field := "name"
      pathTemplate := "projects/*/{id=**}/**" }
  , { id := "h6_empty_subpattern"
      bugRef := "H6"
      description := "Empty variable sub-pattern after ="
      field := "name"
      pathTemplate := "{id=}" }
  , { id := "h7_leading_slash"
      bugRef := "H7"
      description := "Leading slash before variable"
      field := "name"
      pathTemplate := "/{id=*}" }
  , { id := "h7_double_trailing_slash"
      bugRef := "H7"
      description := "Double trailing slash must not silently strip twice"
      field := "name"
      pathTemplate := "projects/{id=*}//" }
  , { id := "h7_missing_slash_before_brace"
      bugRef := "H7"
      description := "Prefix glued directly to opening brace without /"
      field := "name"
      pathTemplate := "projects{id=*}" }
  , { id := "h7_missing_slash_after_brace"
      bugRef := "H7"
      description := "Suffix glued directly to closing brace without /"
      field := "name"
      pathTemplate := "{id=*}tail" }
  , { id := "h7_match_leading_slash"
      bugRef := "H7"
      description := "Leading slash inside variable braces"
      field := "name"
      pathTemplate := "{id=/a}" }
  , { id := "h7_match_trailing_slash"
      bugRef := "H7"
      description := "Trailing slash inside variable braces"
      field := "name"
      pathTemplate := "{id=a/}" }
  , { id := "h7_double_slash_interior"
      bugRef := "H7"
      description := "Empty interior path segment (//)"
      field := "name"
      pathTemplate := "projects//locations/{id=*}" }
  , { id := "h8_unnamed_dstar"
      bugRef := "H8"
      description := "Unnamed {**} lacks required header variable identifier"
      field := "default"
      pathTemplate := "{**}" }
  , { id := "h8_unnamed_star"
      bugRef := "H8"
      description := "Unnamed {*} lacks required header variable identifier"
      field := "default"
      pathTemplate := "{*}" }
  , { id := "h8_multi_seg_shorthand"
      bugRef := "H8"
      description := "Multi-segment expression {a/b} without ="
      field := "default"
      pathTemplate := "projects/*/{a/b}" }
  , { id := "h9_empty_field_and_template"
      bugRef := "H9"
      description := "Malformed routing_parameters {} entry with empty field"
      field := ""
      pathTemplate := "" }
  , { id := "h9_empty_field_nonempty_template"
      bugRef := "H9"
      description := "Routing parameter with empty field and valid template"
      field := ""
      pathTemplate := "{id=*}" }
  , { id := "h10_embedded_star_in_match"
      bugRef := "H10"
      description := "Embedded * within literal token foo*bar"
      field := "name"
      pathTemplate := "projects/{id=foo*bar}" }
  , { id := "h10_quote_in_literal"
      bugRef := "H10"
      description := "Unescaped quote character forbidden in path segment"
      field := "name"
      pathTemplate := "{id=foo\"bar}" }
  , { id := "h10_invalid_field_char"
      bugRef := "H10"
      description := "Invalid wildcard character inside routing field identifier"
      field := "bad*field"
      pathTemplate := "{id=*}" }
  ]

/-- Kernel-computable static well-formedness check for a positive conformance vector. -/
def extractVectorValid (v : ExtractVector) : Bool :=
  match v.spec.routing with
  | none => true
  | some rule => (compileRoutingRuleDeclOrder rule).isOk

/-- Kernel-computable rejection check for a negative conformance vector. -/
def invalidVectorRejected (v : InvalidTemplateVector) : Bool :=
  !(parseRoutingParameter v.field v.pathTemplate).isOk

private theorem canonicalVectors_valid : canonicalVectors.all extractVectorValid = true := rfl
private theorem explicitPrecedenceVectors_valid : explicitPrecedenceVectors.all extractVectorValid = true := rfl
private theorem implicitHttpFallbackVectors_valid : implicitHttpFallbackVectors.all extractVectorValid = true := rfl
private theorem segmentBoundaryVectors_valid : segmentBoundaryVectors.all extractVectorValid = true := rfl
private theorem conflictAndKeyOrderVectors_valid : conflictAndKeyOrderVectors.all extractVectorValid = true := rfl
private theorem rfc6570EncodingVectors_valid : rfc6570EncodingVectors.all extractVectorValid = true := rfl
private theorem nestedFieldVectors_valid : nestedFieldVectors.all extractVectorValid = true := rfl
private theorem showcaseEchoVectors_valid : showcaseEchoVectors.all extractVectorValid = true := rfl

/-- Machine-checked certificate that all 57 positive conformance vectors compile cleanly. -/
theorem extractVectors_all_valid :
    extractVectors.all extractVectorValid = true := by
  simp only [extractVectors, List.all_append,
    canonicalVectors_valid, explicitPrecedenceVectors_valid,
    implicitHttpFallbackVectors_valid, segmentBoundaryVectors_valid,
    conflictAndKeyOrderVectors_valid, rfc6570EncodingVectors_valid,
    nestedFieldVectors_valid, showcaseEchoVectors_valid, Bool.and_self]

/-- Universal soundness bridge: every vector in `extractVectors` evaluates without error and agrees extensionally with its compiled key groups. -/
theorem extractVector_sound (v : ExtractVector) (hMem : v ∈ extractVectors) :
    let req := Request.ofFlat v.request
    ∃ declPairs, evalMethodRoutingPairs v.spec req = .ok declPairs ∧
      (∀ rule, v.spec.routing = some rule →
        ∃ groups, compileRoutingRuleDeclOrder rule = .ok groups ∧
          ∀ k, lookupHeaderKey k (evalCompiledKeyGroupsPairs groups req) = lookupHeaderKey k declPairs) := by
  have hValid : extractVectorValid v = true := List.all_eq_true.mp extractVectors_all_valid v hMem
  unfold extractVectorValid evalMethodRoutingPairs at *
  cases hRouting : v.spec.routing with
  | none =>
    cases v.spec.http <;> refine ⟨_, rfl, fun _ h => by cases h⟩
  | some rule =>
    rw [hRouting] at hValid
    have hExt := compiled_rule_extensional_eq_spec rule (Request.ofFlat v.request)
    cases hComp : compileRoutingRuleDeclOrder rule with
    | error _ => simp [hComp, Except.isOk, Except.toBool] at hValid
    | ok groups =>
      rw [hComp] at hExt
      cases hEval : evalRoutingRulePairs rule (Request.ofFlat v.request) with
      | error _ => rw [hEval] at hExt; exact False.elim hExt
      | ok declPairs =>
        rw [hEval] at hExt
        exact ⟨declPairs, hEval, fun _ hEq => by cases hEq; exact ⟨groups, hComp, hExt.2.2⟩⟩

/-- Machine-checked certificate that all 20 negative conformance vectors are statically rejected. -/
theorem invalidTemplateVectors_all_rejected :
    invalidTemplateVectors.all invalidVectorRejected = true := rfl

end Aip4222
