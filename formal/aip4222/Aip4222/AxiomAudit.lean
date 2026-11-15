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

import Aip4222.Examples
import Aip4222.Theorems
import Aip4222.Vectors

#print axioms Aip4222.Examples.examples_1_to_9_canonical
#print axioms Aip4222.Examples.example_6a_vs_6b_without_instance
#print axioms Aip4222.Examples.example_9_empty_app_profile_fallback
#print axioms Aip4222.Examples.routing_proto_line58_typo_exposed
#print axioms Aip4222.Examples.empty_routing_suppresses_http
#print axioms Aip4222.Examples.parser_rejects_mid_multi_wildcard
#print axioms Aip4222.Examples.parser_edge_cases
#print axioms Aip4222.Examples.encoding_rfc6570_edge_cases
#print axioms Aip4222.Examples.http_implicit_and_precedence_cases
#print axioms Aip4222.Examples.trailing_slash_and_colon_and_zero_seg_cases

-- Phase 2 universal refinement & bug witness theorems
#print axioms Aip4222.explicit_routing_ignores_http
#print axioms Aip4222.empty_routing_suppresses_http_universal
#print axioms Aip4222.h1_buggy_nil_check_falsified
#print axioms Aip4222.parseRoutingParameter_sound
#print axioms Aip4222.well_formed_forbids_mid_multi_wildcard
#print axioms Aip4222.h2_mid_multi_wildcard_always_dead
#print axioms Aip4222.h3a_buggy_rust_prefix_drops_slash
#print axioms Aip4222.h3b_buggy_rust_suffix_drops_slash
#print axioms Aip4222.h3c_buggy_rust_single_dstar_suffix_overmatches
#print axioms Aip4222.compiled_matcher_refines_spec
#print axioms Aip4222.reverse_first_match_eq_last_wins
#print axioms Aip4222.compileRoutingRuleDeclOrder_error_eq
#print axioms Aip4222.compiled_rule_extensional_eq_spec
#print axioms Aip4222.compiled_rule_single_variant_eq_spec
#print axioms Aip4222.static_vs_lazy_key_order_diverges
#print axioms Aip4222.percentEncode_chars_valid
#print axioms Aip4222.percentEncode_empty_iff
#print axioms Aip4222.formatHeaderPairs_empty_iff
#print axioms Aip4222.formatHeaderPairs_no_leading_amp
#print axioms Aip4222.h4_buggy_rust_http_fold_falsified
#print axioms Aip4222.extractVectors_all_valid
#print axioms Aip4222.extractVector_sound
#print axioms Aip4222.invalidTemplateVectors_all_rejected
