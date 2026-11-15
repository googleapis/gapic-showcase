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

import Aip4222.Compiler
import Aip4222.Examples

namespace Aip4222

/-!
## Group 2: Explicit-Over-Implicit Precedence & H1 Falsification
-/

/-- Explicit `google.api.routing` completely overrides `google.api.http` for any HTTP rule and request. -/
theorem explicit_routing_ignores_http (rule : RoutingRule) (h₁ h₂ : Option HttpRule) (req : Request) :
    evalMethodRouting { routing := some rule, http := h₁ } req =
      evalMethodRouting { routing := some rule, http := h₂ } req :=
  rfl

/-- An explicit empty `option (google.api.routing) = {}` universally suppresses implicit HTTP routing. -/
theorem empty_routing_suppresses_http_universal (h : Option HttpRule) (req : Request) :
    evalMethodRouting { routing := some { parameters := [] }, http := h } req = .ok "" :=
  rfl

/-- Concrete falsification of `librarian`'s H1 bug (`HasRouting()` nil/empty check). -/
theorem h1_buggy_nil_check_falsified :
    let spec : MethodRoutingSpec := {
      routing := some { parameters := [] }
      http := some { bindings := [[{ fieldPath := ["parent"] }]] }
    }
    let req := Request.ofFlat [("parent", "projects/p1")]
    evalMethodRoutingBuggyNilCheck spec req = .ok "parent=projects%2Fp1" ∧
    evalMethodRouting spec req = .ok "" ∧
    evalMethodRoutingBuggyNilCheck spec req ≠ evalMethodRouting spec req := by
  refine ⟨rfl, rfl, ?_⟩
  have h1 : evalMethodRoutingBuggyNilCheck
      { routing := some { parameters := [] }
        http := some { bindings := [[{ fieldPath := ["parent"] }]] } }
      (Request.ofFlat [("parent", "projects/p1")]) = .ok "parent=projects%2Fp1" := rfl
  have h2 : evalMethodRouting
      { routing := some { parameters := [] }
        http := some { bindings := [[{ fieldPath := ["parent"] }]] } }
      (Request.ofFlat [("parent", "projects/p1")]) = .ok "" := rfl
  rw [h1, h2]
  intro hEq
  injection hEq with hStr
  exact (by decide : ("parent=projects%2Fp1" : String) ≠ "") hStr

/-!
## Group 3: Parser Soundness, Well-Formedness, & H2 Dead-Matcher Proof
-/

private theorem splitCharsOn_ne_nil (sep : Char) (cs : List Char) : splitCharsOn sep cs ≠ [] := by
  induction cs with
  | nil => simp [splitCharsOn]
  | cons c rest ih =>
    unfold splitCharsOn
    by_cases h : (c == sep) = true
    · simp [h]
    · simp only [Bool.not_eq_true] at h
      simp only [h, Bool.false_eq_true, ↓reduceIte]
      cases hSplit : splitCharsOn sep rest with
      | nil => exact absurd hSplit ih
      | cons cur tails => simp

private theorem all_of_splitCharsOn_all (sep : Char) (cs : List Char) (p : Char → Bool)
    (hSep : p sep = true)
    (hSegs : ∀ seg ∈ splitCharsOn sep cs, seg.all p = true) :
    cs.all p = true := by
  induction cs with
  | nil => rfl
  | cons c rest ih =>
    unfold splitCharsOn at hSegs
    by_cases hEq : (c == sep) = true
    · have hc : c = sep := eq_of_beq hEq
      simp only [hEq, ↓reduceIte, List.mem_cons] at hSegs
      have hRest : rest.all p = true := ih (fun seg hs => hSegs seg (Or.inr hs))
      simp [List.all_cons, hc, hSep, hRest]
    · simp only [Bool.not_eq_true] at hEq
      simp only [hEq, Bool.false_eq_true, ↓reduceIte] at hSegs
      cases hSplit : splitCharsOn sep rest with
      | nil => exact absurd hSplit (splitCharsOn_ne_nil sep rest)
      | cons cur tails =>
        simp only [hSplit, List.mem_cons] at hSegs
        have hFirst : (c :: cur).all p = true := hSegs (c :: cur) (Or.inl rfl)
        simp only [List.all_cons, Bool.and_eq_true] at hFirst
        have hRest : rest.all p = true := ih (fun seg hs => by
          rw [hSplit, List.mem_cons] at hs
          rcases hs with rfl | hIn
          · exact hFirst.2
          · exact hSegs seg (Or.inr hIn))
        simp [List.all_cons, hFirst.1, hRest]

private theorem isValidLitChar_dot : isValidLitChar '.' = true := by decide

private theorem isValidHeaderKey_of_isValidFieldPathStr (s : String)
    (h : isValidFieldPathStr s = true) : isValidHeaderKey s = true := by
  unfold isValidFieldPathStr at h
  simp only [Bool.and_eq_true, bne_iff_ne, ne_eq] at h
  rcases h with ⟨⟨hNe, _⟩, hAll⟩
  have hEqFalse : (s == "") = false := beq_false_of_ne hNe
  have hFp : FieldPath.parse s = splitStringOn '.' s := by
    simp [FieldPath.parse, hEqFalse]
  rw [hFp, List.all_eq_true] at hAll
  simp only [isValidHeaderKey, Bool.and_eq_true, bne_iff_ne, ne_eq]
  refine ⟨hNe, ?_⟩
  apply all_of_splitCharsOn_all '.' s.toList isValidLitChar isValidLitChar_dot
  intro seg hMem
  have hSegValid := hAll (String.ofList seg) (by
    simp only [splitStringOn, List.mem_map]
    exact ⟨seg, hMem, rfl⟩)
  unfold isValidHeaderKey at hSegValid
  simp only [Bool.and_eq_true, String.toList_ofList] at hSegValid
  exact hSegValid.2

private theorem except_bind_eq_ok {ε α β : Type} (x : Except ε α) (f : α → Except ε β) (y : β)
    (h : (x >>= f) = .ok y) : ∃ a, x = .ok a ∧ f a = .ok y := by
  cases x with
  | error e => cases h
  | ok a => exact ⟨a, rfl, h⟩

private theorem tail_wf_check_sound (cand : CaptureTemplate) (msg : String) (t : CaptureTemplate)
    (h : (if cand.isWellFormed then pure cand else throw msg : Except String CaptureTemplate) = .ok t) :
    t.WellFormed := by
  by_cases hWF : cand.isWellFormed = true
  · simp only [hWF, ↓reduceIte, pure, Except.pure, Except.ok.injEq] at h
    subst h
    exact hWF
  · simp [hWF, throw, throwThe, MonadExceptOf.throw] at h

private theorem after_match_step_sound (raw afterClose hk : String)
    (p ms : List AtomSeg) (t : CaptureTemplate)
    (h : ((do
      let suffixSegs ←
        if afterClose == "" then Except.ok []
        else match stripPrefixSlash? afterClose with
          | none => Except.error s!"expected '/' after '}' in {raw}"
          | some s => parseAtomSegs s
      let candidate : CaptureTemplate := {
        headerKey := hk
        prefixSegs := p
        matchSegs := ms
        suffixSegs := suffixSegs
      }
      if candidate.isWellFormed then return candidate
      else throw s!"ill-formed capture template: {raw}") : Except String CaptureTemplate) = .ok t) :
    t.WellFormed := by
  cases hAfter : (afterClose == "") with
  | true =>
    simp only [hAfter, ↓reduceIte, bind, Except.bind] at h
    exact tail_wf_check_sound _ _ t h
  | false =>
    simp only [hAfter, Bool.false_eq_true, ↓reduceIte] at h
    cases hStrip : stripPrefixSlash? afterClose with
    | none => simp only [hStrip, bind, Except.bind] at h; cases h
    | some s =>
      simp only [hStrip] at h
      rcases except_bind_eq_ok _ _ _ h with ⟨ss, _, hEnd⟩
      exact tail_wf_check_sound _ _ t hEnd

private theorem after_prefix_step_sound (raw inside afterClose : String)
    (p : List AtomSeg) (t : CaptureTemplate)
    (h : ((do
      let (headerKey, matchSegs) ←
        match splitStringOn '=' inside with
        | [only] =>
          if isValidHeaderKey only then Except.ok (only, [AtomSeg.singleWildcard])
          else Except.error s!"invalid routing variable name: {only}"
        | [lhs, rhs] =>
          if isValidHeaderKey lhs then do
            let segs ← parseAtomSegs rhs
            Except.ok (lhs, segs)
          else Except.error s!"invalid routing variable name: {lhs}"
        | _ => Except.error s!"invalid variable expression: {inside}"
      let suffixSegs ←
        if afterClose == "" then Except.ok []
        else match stripPrefixSlash? afterClose with
          | none => Except.error s!"expected '/' after '}' in {raw}"
          | some s => parseAtomSegs s
      let candidate : CaptureTemplate := {
        headerKey := headerKey
        prefixSegs := p
        matchSegs := matchSegs
        suffixSegs := suffixSegs
      }
      if candidate.isWellFormed then return candidate
      else throw s!"ill-formed capture template: {raw}") : Except String CaptureTemplate) = .ok t) :
    t.WellFormed := by
  cases hEqSplit : splitStringOn '=' inside with
  | nil => simp [hEqSplit, bind, Except.bind] at h
  | cons lhs rest1 =>
    cases rest1 with
    | nil =>
      rw [hEqSplit] at h
      dsimp only at h
      cases hKey : isValidHeaderKey lhs with
      | false => simp [hKey, bind, Except.bind] at h
      | true =>
        simp only [hKey, ↓reduceIte, bind, Except.bind] at h
        exact after_match_step_sound raw afterClose lhs p [AtomSeg.singleWildcard] t h
    | cons rhs rest2 =>
      cases rest2 with
      | cons _ _ => simp [hEqSplit, bind, Except.bind] at h
      | nil =>
        rw [hEqSplit] at h
        dsimp only at h
        cases hKey : isValidHeaderKey lhs with
        | false => simp [hKey, bind, Except.bind] at h
        | true =>
          simp only [hKey, ↓reduceIte] at h
          rcases except_bind_eq_ok _ _ _ h with ⟨segs, _, hNext⟩
          exact after_match_step_sound raw afterClose lhs p segs t hNext

private theorem parseNonEmptyTemplate_sound (raw norm : String) (t : CaptureTemplate)
    (h : (match splitStringOn '{' norm with
      | [beforeOpen, afterOpen] =>
        if (splitStringOn '}' beforeOpen).length != 1 then
          throw s!"unexpected '}' before '\{' in {raw}"
        else
          match splitStringOn '}' afterOpen with
          | [inside, afterClose] => do
            let prefixSegs ←
              if beforeOpen == "" then Except.ok []
              else match stripSuffixSlash? beforeOpen with
                | none => Except.error s!"expected '/' before '\{' in {raw}"
                | some p => parseAtomSegs p
            let (headerKey, matchSegs) ←
              match splitStringOn '=' inside with
              | [only] =>
                if isValidHeaderKey only then Except.ok (only, [AtomSeg.singleWildcard])
                else Except.error s!"invalid routing variable name: {only}"
              | [lhs, rhs] =>
                if isValidHeaderKey lhs then do
                  let segs ← parseAtomSegs rhs
                  Except.ok (lhs, segs)
                else Except.error s!"invalid routing variable name: {lhs}"
              | _ => Except.error s!"invalid variable expression: {inside}"
            let suffixSegs ←
              if afterClose == "" then Except.ok []
              else match stripPrefixSlash? afterClose with
                | none => Except.error s!"expected '/' after '}' in {raw}"
                | some s => parseAtomSegs s
            let candidate : CaptureTemplate := {
              headerKey := headerKey
              prefixSegs := prefixSegs
              matchSegs := matchSegs
              suffixSegs := suffixSegs
            }
            if candidate.isWellFormed then return candidate
            else throw s!"ill-formed capture template: {raw}"
          | _ => throw s!"expected single '}' in {raw}"
      | _ => throw s!"expected single '\{' in {raw}") = .ok t) :
    t.WellFormed := by
  cases hOpen : splitStringOn '{' norm with
  | nil => simp [hOpen, throw, throwThe, MonadExceptOf.throw] at h
  | cons beforeOpen rest1 =>
    cases rest1 with
    | nil => simp [hOpen, throw, throwThe, MonadExceptOf.throw] at h
    | cons afterOpen rest2 =>
      cases rest2 with
      | cons _ _ => simp [hOpen, throw, throwThe, MonadExceptOf.throw] at h
      | nil =>
        rw [hOpen] at h
        dsimp only at h
        by_cases hBadClose : ((splitStringOn '}' beforeOpen).length != 1) = true
        · simp [hBadClose, throw, throwThe, MonadExceptOf.throw] at h
        · simp only [hBadClose, Bool.false_eq_true, ↓reduceIte] at h
          cases hClose : splitStringOn '}' afterOpen with
          | nil => simp [hClose, throw, throwThe, MonadExceptOf.throw] at h
          | cons inside rest3 =>
            cases rest3 with
            | nil => simp [hClose, throw, throwThe, MonadExceptOf.throw] at h
            | cons afterClose rest4 =>
              cases rest4 with
              | cons _ _ => simp [hClose, throw, throwThe, MonadExceptOf.throw] at h
              | nil =>
                rw [hClose] at h
                dsimp only at h
                cases hBefore : (beforeOpen == "") with
                | true =>
                  simp only [hBefore, ↓reduceIte, bind, Except.bind] at h
                  exact after_prefix_step_sound raw inside afterClose [] t h
                | false =>
                  simp only [hBefore, Bool.false_eq_true, ↓reduceIte] at h
                  cases hTrimPfx : stripSuffixSlash? beforeOpen with
                  | none =>
                    simp only [hTrimPfx, bind, Except.bind] at h
                    cases h
                  | some pStr =>
                    simp only [hTrimPfx] at h
                    rcases except_bind_eq_ok _ _ _ h with ⟨p, _, hp⟩
                    exact after_prefix_step_sound raw inside afterClose p t hp

/-- Every template accepted by `parseRoutingParameter` satisfies `CaptureTemplate.WellFormed`. -/
theorem parseRoutingParameter_sound (field raw : String) (t : CaptureTemplate)
    (h : parseRoutingParameter field raw = .ok t) : t.WellFormed := by
  unfold parseRoutingParameter at h
  cases hValidField : isValidFieldPathStr field with
  | false =>
    simp only [hValidField, Bool.not_false, ↓reduceIte, throw, throwThe, MonadExceptOf.throw,
      bind, Except.bind] at h
    cases h
  | true =>
    simp only [hValidField, Bool.not_true, Bool.false_eq_true, ↓reduceIte] at h
    cases hRaw : (raw == "") with
    | true =>
      simp only [hRaw, ↓reduceIte, pure, Except.pure, Except.ok.injEq] at h
      subst h
      simp only [CaptureTemplate.WellFormed, CaptureTemplate.isWellFormed, CaptureTemplate.allSegs,
        List.nil_append, List.append_nil, List.isEmpty_cons, Bool.not_false,
        List.all_cons, isValidAtom, List.all_nil, Bool.and_true, List.dropLast]
      exact isValidHeaderKey_of_isValidFieldPathStr field hValidField
    | false =>
      simp only [hRaw, Bool.false_eq_true, ↓reduceIte] at h
      cases hTrim : stripSuffixSlash? raw with
      | none =>
        simp only [hTrim] at h
        exact parseNonEmptyTemplate_sound raw raw t h
      | some trimmed =>
        simp only [hTrim] at h
        exact parseNonEmptyTemplate_sound raw trimmed t h

/-- Well-formed templates prohibit `.multiWildcard` anywhere before the final position. -/
theorem well_formed_forbids_mid_multi_wildcard (t : CaptureTemplate) (hWF : t.WellFormed) :
    AtomSeg.multiWildcard ∉ t.prefixSegs ∧
    (t.suffixSegs ≠ [] → AtomSeg.multiWildcard ∉ t.matchSegs) ∧
    AtomSeg.multiWildcard ∉ t.suffixSegs.dropLast := by
  unfold CaptureTemplate.WellFormed CaptureTemplate.isWellFormed CaptureTemplate.allSegs at hWF
  simp only [Bool.and_eq_true, Bool.not_eq_true'] at hWF
  rcases hWF with ⟨⟨⟨_, hMatchNonEmpty⟩, _⟩, hNoMidMulti⟩
  rw [List.all_eq_true] at hNoMidMulti
  have hMatchNe : t.matchSegs ≠ [] := by
    intro hm
    simp [hm] at hMatchNonEmpty
  have notMulti : ∀ x ∈ (t.prefixSegs ++ t.matchSegs ++ t.suffixSegs).dropLast, x ≠ AtomSeg.multiWildcard := by
    intro x hx hEq
    have hb := hNoMidMulti x hx
    subst hEq
    revert hb
    decide
  cases hSuf : t.suffixSegs with
  | nil =>
    have hDrop : (t.prefixSegs ++ t.matchSegs ++ []).dropLast = t.prefixSegs ++ t.matchSegs.dropLast := by
      rw [List.append_nil]
      exact List.dropLast_append_of_ne_nil hMatchNe
    refine ⟨?_, ?_, ?_⟩
    · intro hMem
      apply notMulti .multiWildcard _ rfl
      rw [hSuf, hDrop]
      exact List.mem_append_left _ hMem
    · intro hNe
      exact absurd rfl hNe
    · simp
  | cons s ss =>
    have hDrop : (t.prefixSegs ++ t.matchSegs ++ (s :: ss)).dropLast =
        (t.prefixSegs ++ t.matchSegs) ++ (s :: ss).dropLast :=
      List.dropLast_append_of_ne_nil (by simp)
    refine ⟨?_, ?_, ?_⟩
    · intro hMem
      apply notMulti .multiWildcard _ rfl
      rw [hSuf, hDrop]
      exact List.mem_append_left _ (List.mem_append_left _ hMem)
    · intro _ hMem
      apply notMulti .multiWildcard _ rfl
      rw [hSuf, hDrop]
      exact List.mem_append_left _ (List.mem_append_right _ hMem)
    · intro hMem
      apply notMulti .multiWildcard _ rfl
      rw [hSuf, hDrop]
      exact List.mem_append_right _ hMem

/-- Concrete ill-formed H2 template accepted by `librarian/parser/routing_info_test.go:403`. -/
def h2BadTemplate : CaptureTemplate :=
  {
    headerKey := "routing"
    prefixSegs := [.lit "p", .singleWildcard, .lit "q", .singleWildcard]
    matchSegs := [.lit "a", .singleWildcard, .lit "b", .multiWildcard]
    suffixSegs := [.lit "s", .singleWildcard, .lit "u", .singleWildcard, .lit "v", .singleWildcard]
  }

private theorem matchSegList_append (xs ys : List RuntimeSeg) (cs : List Char) :
    matchSegList (xs ++ ys) cs =
      match matchSegList xs cs with
      | none => none
      | some (nx, remx) =>
        match matchSegList ys remx with
        | none => none
        | some (ny, remy) => some (nx + ny, remy) := by
  induction xs generalizing cs with
  | nil =>
    simp only [List.nil_append, matchSegList]
    cases matchSegList ys cs <;> simp
  | cons x rest ih =>
    simp only [List.cons_append, matchSegList, bind, Option.bind]
    cases hx : x.matchSize cs with
    | none => simp
    | some n =>
      dsimp only
      rw [ih (cs.drop n)]
      cases hRest : matchSegList rest (cs.drop n) with
      | none => simp
      | some p =>
        simp only []
        cases hTail : matchSegList ys p.2 with
        | none => simp [hTail, pure]
        | some q => simp [hTail, pure, Nat.add_assoc]

private theorem trail_matchSize_drop_nil (trail : RuntimeSeg)
    (hTrail : trail = .multiWildcard ∨ trail = .trailingMultiWildcard)
    (cs : List Char) (n : Nat) (h : trail.matchSize cs = some n) : cs.drop n = [] := by
  rcases hTrail with rfl | rfl
  · simp [RuntimeSeg.matchSize] at h
    subst h
    simp
  · simp only [RuntimeSeg.matchSize] at h
    split at h <;> cases h <;> simp

private theorem matchSegList_trailingMulti_rem_nil (segs : List RuntimeSeg) (trail : RuntimeSeg)
    (hTrail : trail = .multiWildcard ∨ trail = .trailingMultiWildcard)
    (cs rem : List Char) (k : Nat)
    (hMatch : matchSegList (segs ++ [trail]) cs = some (k, rem)) : rem = [] := by
  rw [matchSegList_append] at hMatch
  cases hPre : matchSegList segs cs with
  | none => simp [hPre] at hMatch
  | some p =>
    simp only [hPre, matchSegList, bind, Option.bind] at hMatch
    cases hLast : trail.matchSize p.2 with
    | none => simp [hLast] at hMatch
    | some n =>
      rw [hLast] at hMatch
      dsimp only at hMatch
      cases hMatch
      exact trail_matchSize_drop_nil trail hTrail p.2 n hLast

/--
Because greedy `**` (`matchSegList_trailingMulti_rem_nil`) consumes the entire remaining string
(`rem = []`), `librarian`'s H2 test template `h2BadTemplate` (`p/*/q/*/{routing=a/*/b/**}/s/*/u/*/v/*`)
is rejected by AIP-4222 (`parseRoutingParameter ...` fails, `¬ h2BadTemplate.WellFormed`) AND compiles
under both `compileTemplate` and `compileTemplateBuggyRust` to a universally dead matcher yielding
`none` on every `val : String`.
-/
theorem h2_mid_multi_wildcard_always_dead :
    (parseRoutingParameter "default" "p/*/q/*/{routing=a/*/b/**}/s/*/u/*/v/*").isOk = false ∧
    ¬ h2BadTemplate.WellFormed ∧
    ∀ (val : String),
      evalCompiled (compileTemplate h2BadTemplate) val = none ∧
      evalCompiled (compileTemplateBuggyRust h2BadTemplate) val = none := by
  have hm1 : (compileTemplate h2BadTemplate).matchSegs =
      [.literal "a/", .singleWildcard, .literal "/b"] ++ [.trailingMultiWildcard] := by decide
  have hs1 : (compileTemplate h2BadTemplate).suffixSegs =
      [.literal "/s/", .singleWildcard, .literal "/u/", .singleWildcard, .literal "/v/", .singleWildcard] := by decide
  have hm2 : (compileTemplateBuggyRust h2BadTemplate).matchSegs =
      [.literal "a/", .singleWildcard, .literal "/b"] ++ [.trailingMultiWildcard] := by decide
  have hs2 : (compileTemplateBuggyRust h2BadTemplate).suffixSegs =
      [.literal "s/", .singleWildcard, .literal "/u/", .singleWildcard, .literal "/v/", .singleWildcard] := by decide
  have hSuf1Nil : matchSegList
      [.literal "/s/", .singleWildcard, .literal "/u/", .singleWildcard, .literal "/v/", .singleWildcard] [] = none := rfl
  have hSuf2Nil : matchSegList
      [.literal "s/", .singleWildcard, .literal "/u/", .singleWildcard, .literal "/v/", .singleWildcard] [] = none := rfl
  refine ⟨by decide, by decide, ?_⟩
  intro val
  constructor
  · unfold evalCompiled
    simp only [bind, Option.bind]
    cases hp : matchSegList (compileTemplate h2BadTemplate).prefixSegs val.toList with
    | none => rfl
    | some p1 =>
      dsimp only
      rw [hm1]
      cases hm : matchSegList ([.literal "a/", .singleWildcard, .literal "/b"] ++ [.trailingMultiWildcard]) p1.2 with
      | none => rfl
      | some p2 =>
        dsimp only
        have hp2Nil := matchSegList_trailingMulti_rem_nil _ .trailingMultiWildcard (Or.inr rfl) p1.2 p2.2 p2.1 hm
        rw [hp2Nil, hs1, hSuf1Nil]
  · unfold evalCompiled
    simp only [bind, Option.bind]
    cases hp : matchSegList (compileTemplateBuggyRust h2BadTemplate).prefixSegs val.toList with
    | none => rfl
    | some p1 =>
      dsimp only
      rw [hm2]
      cases hm : matchSegList ([.literal "a/", .singleWildcard, .literal "/b"] ++ [.trailingMultiWildcard]) p1.2 with
      | none => rfl
      | some p2 =>
        dsimp only
        have hp2Nil := matchSegList_trailingMulti_rem_nil _ .trailingMultiWildcard (Or.inr rfl) p1.2 p2.2 p2.1 hm
        rw [hp2Nil, hs2, hSuf2Nil]

/-!
## Group 5: Compiler Refinement (`compiled_matcher_refines_spec`) & H3 Falsifications
*(Placed before Group 4 because Group 4 depends on `compiled_matcher_refines_spec`.)*
-/

/-- Concrete falsification of H3a: `librarian` Rust drops trailing `/` on `PrefixSegments`. -/
theorem h3a_buggy_rust_prefix_drops_slash :
    let t : CaptureTemplate := {
      headerKey := "id"
      prefixSegs := [.lit "projects", .singleWildcard]
      matchSegs := [.lit "instances", .singleWildcard]
      suffixSegs := []
    }
    let val := "projects/p/instances/i"
    t.WellFormed ∧
    evalCompiled (compileTemplateBuggyRust t) val = none ∧
    evalCompiled (compileTemplate t) val = some "instances/i" ∧
    extractCapture t val = some "instances/i" := by
  exact ⟨by decide, by rfl, by rfl, by rfl⟩

/-- Concrete falsification of H3b: `librarian` Rust drops leading `/` on `SuffixSegments`. -/
theorem h3b_buggy_rust_suffix_drops_slash :
    let t : CaptureTemplate := {
      headerKey := "id"
      prefixSegs := []
      matchSegs := [.lit "instances", .singleWildcard]
      suffixSegs := [.lit "tables", .singleWildcard]
    }
    let val := "instances/i/tables/t"
    t.WellFormed ∧
    evalCompiled (compileTemplateBuggyRust t) val = none ∧
    evalCompiled (compileTemplate t) val = some "instances/i" ∧
    extractCapture t val = some "instances/i" := by
  exact ⟨by decide, by rfl, by rfl, by rfl⟩

/-- Concrete falsification of H3c: `librarian` Rust emits `MultiWildcard` for singleton `**` suffix. -/
theorem h3c_buggy_rust_single_dstar_suffix_overmatches :
    let t : CaptureTemplate := {
      headerKey := "id"
      prefixSegs := []
      matchSegs := [.lit "projects"]
      suffixSegs := [.multiWildcard]
    }
    let val := "projectsFOO"
    t.WellFormed ∧
    evalCompiled (compileTemplateBuggyRust t) val = some "projects" ∧
    evalCompiled (compileTemplate t) val = none ∧
    extractCapture t val = none := by
  exact ⟨by decide, by rfl, by rfl, by rfl⟩

private theorem matchSegList_drop_eq {segs : List RuntimeSeg} {cs rem : List Char} {n : Nat}
    (h : matchSegList segs cs = some (n, rem)) : rem = cs.drop n ∧ n ≤ cs.length := by
  induction segs generalizing cs n rem with
  | nil =>
    simp only [matchSegList, Option.some.injEq, Prod.mk.injEq] at h
    rcases h with ⟨rfl, rfl⟩
    simp
  | cons seg rest ih =>
    simp only [matchSegList, bind, Option.bind] at h
    cases hk : seg.matchSize cs with
    | none => simp [hk] at h
    | some k =>
      rw [hk] at h
      dsimp only at h
      cases hr : matchSegList rest (cs.drop k) with
      | none => simp [hr] at h
      | some p =>
        rw [hr] at h
        simp only [pure, Option.some.injEq, Prod.mk.injEq] at h
        rcases h with ⟨rfl, rfl⟩
        rcases ih hr with ⟨hDrop, hLe⟩
        have hkLe : k ≤ cs.length := by
          cases seg with
          | literal s =>
            simp only [RuntimeSeg.matchSize] at hk
            split at hk <;> cases hk
            next hPref =>
              rw [List.isPrefixOf_iff_prefix] at hPref
              exact hPref.length_le
          | singleWildcard =>
            simp only [RuntimeSeg.matchSize] at hk
            split at hk <;> cases hk
            have hLen := congrArg List.length (List.takeWhile_append_dropWhile (p := (· != '/')) (l := cs))
            rw [List.length_append] at hLen
            omega
          | multiWildcard =>
            simp only [RuntimeSeg.matchSize, Option.some.injEq] at hk
            omega
          | trailingMultiWildcard =>
            simp only [RuntimeSeg.matchSize] at hk
            split at hk <;> cases hk <;> simp
        refine ⟨?_, ?_⟩
        · rw [hDrop, List.drop_drop, Nat.add_comm]
        · have hLen : (cs.drop k).length = cs.length - k := List.length_drop
          omega

private theorem isPrefixOf_append_chars (p q r : List Char) :
    List.isPrefixOf (p ++ q) r = (List.isPrefixOf p r && List.isPrefixOf q (r.drop p.length)) := by
  induction p generalizing r with
  | nil => simp [List.isPrefixOf]
  | cons a as ih =>
    cases r with
    | nil => simp [List.isPrefixOf]
    | cons b bs =>
      simp only [List.cons_append, List.isPrefixOf, List.length_cons, List.drop_succ_cons, ih]
      simp [Bool.and_assoc]

private theorem matchSegList_append_left_cong (xs t1 t2 : List RuntimeSeg)
    (hEquiv : ∀ cs, matchSegList t1 cs = matchSegList t2 cs) (cs : List Char) :
    matchSegList (xs ++ t1) cs = matchSegList (xs ++ t2) cs := by
  rw [matchSegList_append, matchSegList_append]
  cases matchSegList xs cs with
  | none => rfl
  | some p => simp [hEquiv p.2]

private theorem matchSegList_lit_empty (tail : List RuntimeSeg) (cs : List Char) :
    matchSegList (.literal "" :: tail) cs = matchSegList tail cs := by
  simp only [matchSegList, RuntimeSeg.matchSize, String.toList_empty, List.isPrefixOf,
    ↓reduceIte, List.length_nil, List.drop_zero, bind, Option.bind]
  cases matchSegList tail cs <;> simp [pure]

private theorem matchSegList_lit_append (a b : String) (tail : List RuntimeSeg) (cs : List Char) :
    matchSegList (.literal a :: .literal b :: tail) cs =
      matchSegList (.literal (a ++ b) :: tail) cs := by
  simp only [matchSegList, RuntimeSeg.matchSize, String.toList_append, List.length_append,
    isPrefixOf_append_chars, bind, Option.bind]
  by_cases ha : a.toList.isPrefixOf cs = true
  · by_cases hb : b.toList.isPrefixOf (cs.drop a.toList.length) = true
    · simp only [ha, hb, ↓reduceIte, Bool.and_self, List.drop_drop]
      cases matchSegList tail (cs.drop (a.toList.length + b.toList.length)) <;> simp [pure, Nat.add_assoc]
    · have hb' : ¬ (b.toList <+: cs.drop a.toList.length) := by simpa [List.isPrefixOf_iff_prefix] using hb
      simp [ha, hb, hb']
  · simp [ha]

private theorem matchSegList_flushLitBuffer (buf : String) (acc tail : List RuntimeSeg) (cs : List Char) :
    matchSegList (flushLitBuffer buf acc ++ tail) cs =
      matchSegList (acc ++ [.literal buf] ++ tail) cs := by
  unfold flushLitBuffer
  by_cases hb : (buf == "") = true
  · have hEq : buf = "" := eq_of_beq hb
    subst hEq
    simp only [hb, ↓reduceIte, List.append_assoc]
    exact (matchSegList_append_left_cong acc _ tail (fun r => matchSegList_lit_empty tail r) cs).symm
  · simp [hb, List.append_assoc]

private theorem matchSegList_annotateSegmentsLoop (segs : List AtomSeg) (totalLen index : Nat)
    (isPrefix isSuffix : Bool) (litBuf : String) (acc : List RuntimeSeg) (cs : List Char) :
    matchSegList (annotateSegmentsLoop segs totalLen index isPrefix isSuffix litBuf acc) cs =
      matchSegList (acc ++ [.literal litBuf] ++ emitUnmergedLoop segs totalLen index isPrefix isSuffix) cs := by
  induction segs generalizing index litBuf acc cs with
  | nil =>
    unfold annotateSegmentsLoop emitUnmergedLoop
    simp only []
    by_cases hp : (isPrefix && decide (totalLen > 0)) = true
    · simp only [hp, ↓reduceIte]
      have hFlush := matchSegList_flushLitBuffer (litBuf ++ "/") acc [] cs
      simp only [List.append_nil] at hFlush
      rw [hFlush, List.append_assoc]
      apply matchSegList_append_left_cong acc
      intro r
      exact (matchSegList_lit_append litBuf "/" [] r).symm
    · simp only [hp, Bool.false_eq_true, ↓reduceIte, List.append_nil]
      have hFlush := matchSegList_flushLitBuffer litBuf acc [] cs
      simpa using hFlush
  | cons seg rest ih =>
    unfold annotateSegmentsLoop emitUnmergedLoop
    dsimp only
    cases seg with
    | multiWildcard =>
      dsimp only
      rw [ih]
      let kind := if totalLen == 1 && !isSuffix then RuntimeSeg.multiWildcard else if totalLen != index + 1 then RuntimeSeg.multiWildcard else RuntimeSeg.trailingMultiWildcard
      have h1 : matchSegList (flushLitBuffer litBuf acc ++ [kind, .literal ""] ++
          emitUnmergedLoop rest totalLen (index + 1) isPrefix isSuffix) cs =
          matchSegList (acc ++ [.literal litBuf, kind] ++
            emitUnmergedLoop rest totalLen (index + 1) isPrefix isSuffix) cs := by
        calc matchSegList (flushLitBuffer litBuf acc ++ [kind, .literal ""] ++
                emitUnmergedLoop rest totalLen (index + 1) isPrefix isSuffix) cs
            = matchSegList (flushLitBuffer litBuf acc ++
                ([kind] ++ (.literal "" :: emitUnmergedLoop rest totalLen (index + 1) isPrefix isSuffix))) cs := by
              simp [List.append_assoc]
          _ = matchSegList (flushLitBuffer litBuf acc ++
                ([kind] ++ emitUnmergedLoop rest totalLen (index + 1) isPrefix isSuffix)) cs := by
              apply matchSegList_append_left_cong
              intro r
              apply matchSegList_append_left_cong [kind]
              intro r2
              exact matchSegList_lit_empty _ r2
          _ = matchSegList (acc ++ [.literal litBuf] ++
                (kind :: emitUnmergedLoop rest totalLen (index + 1) isPrefix isSuffix)) cs := by
              have h := matchSegList_flushLitBuffer litBuf acc
                (kind :: emitUnmergedLoop rest totalLen (index + 1) isPrefix isSuffix) cs
              simpa [List.append_assoc] using h
          _ = matchSegList (acc ++ [.literal litBuf, kind] ++
                emitUnmergedLoop rest totalLen (index + 1) isPrefix isSuffix) cs := by
              simp [List.append_assoc]
      simpa [kind, List.append_assoc] using h1
    | singleWildcard =>
      dsimp only
      rw [ih]
      let buf1 := if index != 0 || isSuffix then litBuf ++ "/" else litBuf
      have hSlash : ∀ r, matchSegList (.literal buf1 :: .singleWildcard ::
            emitUnmergedLoop rest totalLen (index + 1) isPrefix isSuffix) r =
          matchSegList (.literal litBuf :: (if index != 0 || isSuffix then [.literal "/"] else []) ++
            [.singleWildcard] ++ emitUnmergedLoop rest totalLen (index + 1) isPrefix isSuffix) r := by
        intro r
        by_cases hc : (index != 0 || isSuffix) = true
        · simp only [buf1, hc, ↓reduceIte, List.cons_append, List.nil_append]
          exact (matchSegList_lit_append litBuf "/" _ r).symm
        · simp only [buf1, hc, Bool.false_eq_true, ↓reduceIte, List.nil_append, List.cons_append]
      calc matchSegList ((flushLitBuffer buf1 acc ++ [.singleWildcard]) ++ [.literal ""] ++
              emitUnmergedLoop rest totalLen (index + 1) isPrefix isSuffix) cs
          = matchSegList (flushLitBuffer buf1 acc ++ ([.singleWildcard] ++
              (.literal "" :: emitUnmergedLoop rest totalLen (index + 1) isPrefix isSuffix))) cs := by
            simp [List.append_assoc]
        _ = matchSegList (flushLitBuffer buf1 acc ++ ([.singleWildcard] ++
              emitUnmergedLoop rest totalLen (index + 1) isPrefix isSuffix)) cs := by
            apply matchSegList_append_left_cong
            intro r
            apply matchSegList_append_left_cong [.singleWildcard]
            intro r2
            exact matchSegList_lit_empty _ r2
        _ = matchSegList (acc ++ [.literal buf1] ++ ([.singleWildcard] ++
              emitUnmergedLoop rest totalLen (index + 1) isPrefix isSuffix)) cs := by
            exact matchSegList_flushLitBuffer buf1 acc _ cs
        _ = matchSegList (acc ++ ([.literal buf1] ++ (.singleWildcard ::
              emitUnmergedLoop rest totalLen (index + 1) isPrefix isSuffix))) cs := by
            simp [List.append_assoc]
        _ = matchSegList (acc ++ [.literal litBuf] ++
              ((if index != 0 || isSuffix then [.literal "/"] else []) ++ [.singleWildcard] ++
                emitUnmergedLoop rest totalLen (index + 1) isPrefix isSuffix)) cs := by
            simp only [List.append_assoc]
            apply matchSegList_append_left_cong acc
            intro r
            simpa [List.append_assoc] using hSlash r
    | lit s =>
      dsimp only
      rw [ih]
      let buf1 := if index != 0 || isSuffix then litBuf ++ "/" else litBuf
      have hLit : ∀ r, matchSegList (.literal (buf1 ++ s) ::
            emitUnmergedLoop rest totalLen (index + 1) isPrefix isSuffix) r =
          matchSegList (.literal litBuf ::
            ((if index != 0 || isSuffix then [.literal "/"] else []) ++ [.literal s] ++
              emitUnmergedLoop rest totalLen (index + 1) isPrefix isSuffix)) r := by
        intro r
        by_cases hc : (index != 0 || isSuffix) = true
        · simp only [buf1, hc, ↓reduceIte, List.cons_append, List.nil_append]
          rw [← matchSegList_lit_append (litBuf ++ "/") s _ r]
          rw [← matchSegList_lit_append litBuf "/" _ r]
        · simp only [buf1, hc, Bool.false_eq_true, ↓reduceIte, List.nil_append, List.cons_append]
          exact (matchSegList_lit_append litBuf s _ r).symm
      simpa [buf1, List.append_assoc] using matchSegList_append_left_cong acc _ _ hLit cs

private theorem matchSegList_annotateSegments (segs : List AtomSeg) (isPrefix isSuffix : Bool)
    (cs : List Char) :
    matchSegList (annotateSegments segs isPrefix isSuffix) cs =
      matchSegList (emitUnmerged segs isPrefix isSuffix) cs := by
  unfold annotateSegments emitUnmerged
  split
  · rfl
  · rw [matchSegList_annotateSegmentsLoop]
    simp only [List.nil_append]
    exact matchSegList_lit_empty _ cs

private theorem evalCompiled_eq_unmerged (t : CaptureTemplate) (val : String) :
    evalCompiled (compileTemplate t) val =
      (do
        let (_, rem1) ← matchSegList (emitUnmerged t.prefixSegs true false) val.toList
        let (mCount, rem2) ← matchSegList (emitUnmerged t.matchSegs false false) rem1
        let (_, rem3) ← matchSegList (emitUnmerged t.suffixSegs false true) rem2
        if !rem3.isEmpty || mCount == 0 then none
        else some (String.ofList (rem1.take mCount))) := by
  unfold evalCompiled compileTemplate
  simp only [bind, Option.bind]
  rw [matchSegList_annotateSegments]
  cases matchSegList (emitUnmerged t.prefixSegs true false) val.toList with
  | none => rfl
  | some p1 =>
    dsimp only
    rw [matchSegList_annotateSegments]
    cases matchSegList (emitUnmerged t.matchSegs false false) p1.2 with
    | none => rfl
    | some p2 =>
      dsimp only
      rw [matchSegList_annotateSegments]

/-- Runtime atom corresponding to a single non-`**` `AtomSeg`. -/
private def atomToRuntime : AtomSeg → RuntimeSeg
  | .lit s => .literal s
  | .singleWildcard => .singleWildcard
  | .multiWildcard => .multiWildcard

private theorem atomToRuntime_literal (s : String) : atomToRuntime (.lit s) = .literal s := rfl
private theorem atomToRuntime_singleWildcard : atomToRuntime .singleWildcard = .singleWildcard := rfl

private theorem emitUnmergedLoop_prefix (segs : List AtomSeg) (totalLen index : Nat)
    (hNoMulti : AtomSeg.multiWildcard ∉ segs) :
    emitUnmergedLoop segs totalLen index true false =
      emitUnmergedLoop segs totalLen index false false ++
        (if totalLen > 0 then [.literal "/"] else []) := by
  induction segs generalizing index with
  | nil =>
    unfold emitUnmergedLoop
    by_cases ht : totalLen > 0 <;> simp [ht]
  | cons seg rest ih =>
    have hRest : AtomSeg.multiWildcard ∉ rest := fun h => hNoMulti (List.mem_cons_of_mem seg h)
    have hSeg : seg ≠ AtomSeg.multiWildcard := fun h => hNoMulti (by simp [h])
    cases seg with
    | multiWildcard => exact absurd rfl hSeg
    | singleWildcard =>
      simp only [emitUnmergedLoop, Bool.or_false]
      rw [ih (index + 1) hRest]
      simp [List.append_assoc]
    | lit s =>
      simp only [emitUnmergedLoop, Bool.or_false]
      rw [ih (index + 1) hRest]
      simp [List.append_assoc]

private theorem emitUnmerged_prefix (segs : List AtomSeg) (hNe : segs ≠ [])
    (hNoMulti : AtomSeg.multiWildcard ∉ segs) :
    emitUnmerged segs true false = emitUnmerged segs false false ++ [.literal "/"] := by
  cases segs with
  | nil => exact absurd rfl hNe
  | cons s ss =>
    unfold emitUnmerged
    simp only [List.isEmpty_cons, Bool.false_eq_true, ↓reduceIte]
    rw [emitUnmergedLoop_prefix (s :: ss) (s :: ss).length 0 hNoMulti]
    simp

private theorem emitUnmergedLoop_suffix_gt_zero (segs : List AtomSeg) (totalLen index : Nat)
    (hIdx : 1 ≤ index) (hTot : 2 ≤ totalLen) :
    emitUnmergedLoop segs totalLen index false true =
      emitUnmergedLoop segs totalLen index false false := by
  induction segs generalizing index with
  | nil => simp [emitUnmergedLoop]
  | cons seg rest ih =>
    have hIdxNe : (index != 0) = true := by simp [bne]; omega
    have hTotNe : (totalLen == 1) = false := by simp; omega
    cases seg <;> simp only [emitUnmergedLoop, hIdxNe, Bool.true_or, hTotNe, Bool.false_and,
      Bool.false_eq_true, ↓reduceIte] <;> rw [ih (index + 1) (by omega)]

private theorem emitUnmerged_suffix_cons (sHead : AtomSeg) (sTail : List AtomSeg)
    (hHead : sHead ≠ .multiWildcard) :
    emitUnmerged (sHead :: sTail) false true =
      [.literal "/"] ++ emitUnmerged (sHead :: sTail) false false := by
  unfold emitUnmerged
  simp only [List.isEmpty_cons, Bool.false_eq_true, ↓reduceIte, List.length_cons]
  cases sHead with
  | multiWildcard => exact absurd rfl hHead
  | singleWildcard =>
    simp only [emitUnmergedLoop, bne_self_eq_false, Bool.false_or, Bool.or_false, ↓reduceIte,
      List.nil_append, List.cons_append, Nat.zero_add]
    cases sTail with
    | nil => simp [emitUnmergedLoop]
    | cons x xs =>
      simp only [Bool.false_eq_true, ↓reduceIte, List.nil_append, List.cons_append]
      rw [emitUnmergedLoop_suffix_gt_zero (x :: xs) ((x :: xs).length + 1) 1 (by omega) (by simp)]
  | lit s =>
    simp only [emitUnmergedLoop, bne_self_eq_false, Bool.false_or, Bool.or_false, ↓reduceIte,
      List.nil_append, List.cons_append, Nat.zero_add]
    cases sTail with
    | nil => simp [emitUnmergedLoop]
    | cons x xs =>
      simp only [Bool.false_eq_true, ↓reduceIte, List.nil_append, List.cons_append]
      rw [emitUnmergedLoop_suffix_gt_zero (x :: xs) ((x :: xs).length + 1) 1 (by omega) (by simp)]

private theorem emitUnmergedLoop_append_multi (xs : List AtomSeg) (index totalLen : Nat)
    (hTot : totalLen = index + xs.length + 1) (hTotGe2 : 2 ≤ totalLen)
    (hNoMulti : AtomSeg.multiWildcard ∉ xs) (isSuffix : Bool) :
    emitUnmergedLoop (xs ++ [.multiWildcard]) totalLen index false isSuffix =
      emitUnmergedLoop xs (index + xs.length) index false isSuffix ++ [.trailingMultiWildcard] := by
  induction xs generalizing index with
  | nil =>
    simp only [List.nil_append, List.length_nil, Nat.add_zero] at hTot ⊢
    have hTotNe : (totalLen == 1) = false := by simp; omega
    have hEnd : (totalLen != index + 1) = false := by simp; omega
    simp [emitUnmergedLoop, hTotNe, hEnd]
  | cons x rest ih =>
    have hNoRest : AtomSeg.multiWildcard ∉ rest := fun h => hNoMulti (List.mem_cons_of_mem x h)
    have hx : x ≠ .multiWildcard := fun h => hNoMulti (by simp [h])
    have hLen' : totalLen = index + 1 + rest.length + 1 := by simp [hTot]; omega
    cases x with
    | multiWildcard => exact absurd rfl hx
    | singleWildcard =>
      simp only [List.cons_append, emitUnmergedLoop]
      rw [ih (index + 1) hLen' hNoRest]
      have hArg : index + 1 + rest.length = index + (rest.length + 1) := by omega
      simp [List.append_assoc, hArg]
    | lit s =>
      simp only [List.cons_append, emitUnmergedLoop]
      rw [ih (index + 1) hLen' hNoRest]
      have hArg : index + 1 + rest.length = index + (rest.length + 1) := by omega
      simp [List.append_assoc, hArg]

private theorem emitUnmerged_append_multi (xs : List AtomSeg) (hNe : xs ≠ [])
    (hNoMulti : AtomSeg.multiWildcard ∉ xs) :
    emitUnmerged (xs ++ [AtomSeg.multiWildcard]) false false =
      emitUnmerged xs false false ++ [.trailingMultiWildcard] := by
  unfold emitUnmerged
  have hNe2 : (xs ++ [AtomSeg.multiWildcard]).isEmpty = false := by
    cases xs <;> simp
  have hNe1 : xs.isEmpty = false := List.isEmpty_eq_false_iff.mpr hNe
  simp only [hNe2, hNe1, Bool.false_eq_true, ↓reduceIte, List.length_append, List.length_cons,
    List.length_nil]
  have hGe : 2 ≤ xs.length + 1 := by
    cases xs with
    | nil => exact absurd rfl hNe
    | cons _ _ => simp
  have hRes := emitUnmergedLoop_append_multi xs 0 (xs.length + 1) (by simp) hGe hNoMulti false
  simpa using hRes

private theorem emitUnmerged_singleton (p : AtomSeg) (hp : p ≠ .multiWildcard) :
    emitUnmerged [p] false false = [atomToRuntime p] := by
  cases p with
  | multiWildcard => exact absurd rfl hp
  | singleWildcard => rfl
  | lit s => rfl

private theorem emitUnmergedLoop_shift_idx (zs : List AtomSeg) (idx : Nat) (hIdx : 1 ≤ idx) :
    emitUnmergedLoop zs (idx + zs.length) idx false false =
      emitUnmergedLoop zs (1 + zs.length) 1 false false := by
  induction zs generalizing idx with
  | nil => simp [emitUnmergedLoop]
  | cons z rest ih =>
    have hIdxNe : (idx != 0) = true := by simp [bne]; omega
    have hOneNe : ((1 : Nat) != 0) = true := by decide
    have hGe2 : (idx + (z :: rest).length == 1) = false := by simp; try omega
    have hGe2' : (1 + (z :: rest).length == 1) = false := by simp; try omega
    have hEnd1 : (idx + (z :: rest).length != idx + 1) = (rest.length != 0) := by
      cases rest <;> simp [bne] <;> try omega
    have hEnd2 : (1 + (z :: rest).length != 1 + 1) = (rest.length != 0) := by
      cases rest <;> simp [bne] <;> try omega
    have hLen1 : idx + (z :: rest).length = (idx + 1) + rest.length := by simp; try omega
    have hLen2 : 1 + (z :: rest).length = (1 + 1) + rest.length := by simp; try omega
    cases z with
    | multiWildcard =>
      simp only [emitUnmergedLoop, hGe2, hGe2', Bool.false_and, Bool.false_eq_true, ↓reduceIte,
        hEnd1, hEnd2]
      rw [hLen1, ih (idx + 1) (by omega), hLen2, ih 2 (by omega)]
    | singleWildcard =>
      simp only [emitUnmergedLoop, hIdxNe, hOneNe, Bool.true_or, ↓reduceIte]
      rw [hLen1, ih (idx + 1) (by omega), hLen2, ih 2 (by omega)]
    | lit s =>
      simp only [emitUnmergedLoop, hIdxNe, hOneNe, Bool.true_or, ↓reduceIte]
      rw [hLen1, ih (idx + 1) (by omega), hLen2, ih 2 (by omega)]

private theorem emitUnmerged_cons_slash (p q : AtomSeg) (qs : List AtomSeg)
    (hp : p ≠ .multiWildcard) (hq : q ≠ .multiWildcard) :
    emitUnmerged (p :: q :: qs) false false =
      [atomToRuntime p, .literal "/"] ++ emitUnmerged (q :: qs) false false := by
  unfold emitUnmerged
  simp only [List.isEmpty_cons, Bool.false_eq_true, ↓reduceIte, List.length_cons]
  have hTail : emitUnmergedLoop (q :: qs) (qs.length + 1 + 1) 1 false false =
      [.literal "/"] ++ emitUnmergedLoop (q :: qs) (qs.length + 1) 0 false false := by
    have hL1 : qs.length + 1 + 1 = 2 + qs.length := by omega
    have hL2 : qs.length + 1 = 1 + qs.length := by omega
    cases q with
    | multiWildcard => exact absurd rfl hq
    | singleWildcard =>
      have hLHS : emitUnmergedLoop (AtomSeg.singleWildcard :: qs) (qs.length + 1 + 1) 1 false false =
          [.literal "/"] ++ ([.singleWildcard] ++ emitUnmergedLoop qs (qs.length + 1 + 1) 2 false false) := rfl
      have hRHS : emitUnmergedLoop (AtomSeg.singleWildcard :: qs) (qs.length + 1) 0 false false =
          [.singleWildcard] ++ emitUnmergedLoop qs (qs.length + 1) 1 false false := rfl
      rw [hLHS, hRHS, hL1, hL2,
        emitUnmergedLoop_shift_idx qs 2 (by omega), emitUnmergedLoop_shift_idx qs 1 (by omega)]
    | lit s =>
      have hLHS : emitUnmergedLoop (AtomSeg.lit s :: qs) (qs.length + 1 + 1) 1 false false =
          [.literal "/"] ++ ([.literal s] ++ emitUnmergedLoop qs (qs.length + 1 + 1) 2 false false) := rfl
      have hRHS : emitUnmergedLoop (AtomSeg.lit s :: qs) (qs.length + 1) 0 false false =
          [.literal s] ++ emitUnmergedLoop qs (qs.length + 1) 1 false false := rfl
      rw [hLHS, hRHS, hL1, hL2,
        emitUnmergedLoop_shift_idx qs 2 (by omega), emitUnmergedLoop_shift_idx qs 1 (by omega)]
  cases p with
  | multiWildcard => exact absurd rfl hp
  | singleWildcard =>
    have hOuter : emitUnmergedLoop (AtomSeg.singleWildcard :: q :: qs) (qs.length + 1 + 1) 0 false false =
        [.singleWildcard] ++ emitUnmergedLoop (q :: qs) (qs.length + 1 + 1) 1 false false := rfl
    change emitUnmergedLoop (AtomSeg.singleWildcard :: q :: qs) (qs.length + 1 + 1) 0 false false =
        [atomToRuntime .singleWildcard, .literal "/"] ++
          emitUnmergedLoop (q :: qs) (qs.length + 1) 0 false false
    rw [hOuter, hTail]
    rfl
  | lit s =>
    have hOuter : emitUnmergedLoop (AtomSeg.lit s :: q :: qs) (qs.length + 1 + 1) 0 false false =
        [.literal s] ++ emitUnmergedLoop (q :: qs) (qs.length + 1 + 1) 1 false false := rfl
    change emitUnmergedLoop (AtomSeg.lit s :: q :: qs) (qs.length + 1 + 1) 0 false false =
        [atomToRuntime (.lit s), .literal "/"] ++
          emitUnmergedLoop (q :: qs) (qs.length + 1) 0 false false
    rw [hOuter, hTail]
    rfl

/-! ### Character and Segment Split Bridge Lemmas -/

private theorem isValidLitChar_ne_slash_colon (c : Char) (h : isValidLitChar c = true) :
    c ≠ '/' ∧ c ≠ ':' := by
  constructor <;> {
    intro heq
    subst heq
    revert h
    decide
  }

private theorem ValidLitNoSlash {s : String} (h : isValidAtom (.lit s) = true) :
    s ≠ "" ∧ ∀ c ∈ s.toList, c ≠ '/' ∧ c ≠ ':' := by
  unfold isValidAtom at h
  simp only [Bool.and_eq_true, bne_iff_ne, ne_eq, List.all_eq_true] at h
  refine ⟨h.1, fun c hc => isValidLitChar_ne_slash_colon c (h.2 c hc)⟩

private theorem splitCharsOn_takeWhile_dropWhile (sep : Char) (cs : List Char) :
    splitCharsOn sep cs =
      match cs.dropWhile (· != sep) with
      | [] => [cs.takeWhile (· != sep)]
      | _ :: rem => cs.takeWhile (· != sep) :: splitCharsOn sep rem := by
  induction cs with
  | nil => simp [splitCharsOn]
  | cons c rest ih =>
    conv => lhs; unfold splitCharsOn
    by_cases hEq : (c == sep) = true
    · have hc : c = sep := eq_of_beq hEq
      subst hc
      simp [bne_self_eq_false]
    · simp only [Bool.not_eq_true] at hEq
      have hBne : (c != sep) = true := by simp [bne, hEq]
      simp only [hEq, Bool.false_eq_true, ↓reduceIte, List.takeWhile_cons, List.dropWhile_cons, hBne]
      rw [ih]
      cases rest.dropWhile (· != sep) with
      | nil => simp
      | cons _ _ => simp

private theorem dropWhile_ne_slash_head {cs : List Char} {c : Char} {rem : List Char}
    (h : cs.dropWhile (· != '/') = c :: rem) : c = '/' := by
  have hNot := @List.head_dropWhile_not _ (· != '/') cs (by simp [h])
  have hHead : (cs.dropWhile (· != '/')).head (by simp [h]) = c := by simp [h]
  rw [hHead] at hNot
  simpa [bne] using hNot

private theorem takeWhile_dropWhile_append_of_no_slash (xs ys : List Char)
    (hxs : ∀ c ∈ xs, c ≠ '/') :
    (xs ++ ys).takeWhile (· != '/') = xs ++ ys.takeWhile (· != '/') ∧
    (xs ++ ys).dropWhile (· != '/') = ys.dropWhile (· != '/') := by
  induction xs with
  | nil => simp
  | cons x xs ih =>
    have hx : (x != '/') = true := by
      simp [bne]
      exact hxs x (by simp)
    have hRest : ∀ c ∈ xs, c ≠ '/' := fun c hc => hxs c (by simp [hc])
    rcases ih hRest with ⟨ih1, ih2⟩
    simp [hx, ih1, ih2]

private theorem takeWhile_mem {α : Type} {p : α → Bool} :
    ∀ {xs : List α} {x : α}, x ∈ xs.takeWhile p → p x = true
  | [], _, h => by simp at h
  | y :: ys, x, h => by
    simp only [List.takeWhile_cons] at h
    split at h
    · next hy =>
      simp only [List.mem_cons] at h
      rcases h with rfl | hIn
      · exact hy
      · exact takeWhile_mem hIn
    · simp at h

private def splitSegsOf (cs : List Char) : List String :=
  (splitCharsOn '/' cs).map String.ofList

private theorem splitSegsOf_eq_splitStringOn (cs : List Char) :
    splitSegsOf cs = splitStringOn '/' (String.ofList cs) := by
  simp [splitSegsOf, splitStringOn, String.toList_ofList]

private theorem splitSegsOf_ne_nil (cs : List Char) : splitSegsOf cs ≠ [] := by
  simp [splitSegsOf, splitCharsOn_ne_nil]

private theorem splitSegsOf_eq_take_drop (cs : List Char) :
    splitSegsOf cs =
      match cs.dropWhile (· != '/') with
      | [] => [String.ofList (cs.takeWhile (· != '/'))]
      | _ :: rem => String.ofList (cs.takeWhile (· != '/')) :: splitSegsOf rem := by
  unfold splitSegsOf
  rw [splitCharsOn_takeWhile_dropWhile]
  cases cs.dropWhile (· != '/') with
  | nil => simp
  | cons _ rem => simp

private theorem joinStringsWith_slash_of_drop_match : ∀ (n : Nat) (cs : List Char), cs.length ≤ n →
    (joinStringsWith "/" (splitSegsOf cs)).toList = cs
  | 0, [], _ => rfl
  | 0, _ :: _, h => by simp at h
  | n + 1, cs, hLen => by
    rw [splitSegsOf_eq_take_drop]
    cases hDrop : cs.dropWhile (· != '/') with
    | nil =>
      simp only [joinStringsWith, String.toList_ofList]
      have hSum := List.takeWhile_append_dropWhile (p := (· != '/')) (l := cs)
      rw [hDrop, List.append_nil] at hSum
      exact hSum
    | cons d rem =>
      have hd : d = '/' := dropWhile_ne_slash_head hDrop
      subst hd
      have hSum := List.takeWhile_append_dropWhile (p := (· != '/')) (l := cs)
      rw [hDrop] at hSum
      have hRemLt : rem.length ≤ n := by
        have hLenEq := congrArg List.length hSum
        simp only [List.length_append, List.length_cons] at hLenEq
        omega
      have ih := joinStringsWith_slash_of_drop_match n rem hRemLt
      dsimp only
      cases hSegs : splitSegsOf rem with
      | nil => exact absurd hSegs (splitSegsOf_ne_nil rem)
      | cons first rest =>
        rw [hSegs] at ih
        simp only [joinStringsWith, String.toList_append, String.toList_ofList] at ih ⊢
        rw [ih]
        simpa [List.append_assoc] using hSum

private theorem string_ofList_join_splitSegs (cs : List Char) :
    joinStringsWith "/" (splitSegsOf cs) = String.ofList cs := by
  rw [String.ext_iff, String.toList_ofList]
  exact joinStringsWith_slash_of_drop_match cs.length cs (Nat.le_refl _)

private theorem joinStringsWith_cons_ne_nil (x : String) (xs : List String) (h : xs ≠ []) :
    joinStringsWith "/" (x :: xs) = x ++ "/" ++ joinStringsWith "/" xs := by
  cases xs with
  | nil => exact absurd rfl h
  | cons y ys => simp [joinStringsWith, String.append_assoc]

private theorem joinStringsWith_slash_nonempty (xs : List String) (hNe : xs ≠ [])
    (hNonEmpty : ∃ s ∈ xs, s ≠ "") :
    joinStringsWith "/" xs ≠ "" := by
  cases xs with
  | nil => exact absurd rfl hNe
  | cons x xs =>
    cases xs with
    | nil =>
      simp only [joinStringsWith]
      rcases hNonEmpty with ⟨s, hs, hne⟩
      simp only [List.mem_singleton] at hs
      subst hs
      exact hne
    | cons y ys =>
      intro hEmpty
      have hList : (x ++ "/" ++ joinStringsWith "/" (y :: ys)).toList = [] := by
        rw [← joinStringsWith_cons_ne_nil x (y :: ys) (by simp), hEmpty]
        rfl
      simp [String.toList_append] at hList

private theorem matchAtomSegs_cons_no_multi (a : AtomSeg) (pats : List AtomSeg)
    (s : String) (segs : List String) (hNoMulti : AtomSeg.multiWildcard ∉ a :: pats)
    (ha : a ≠ .multiWildcard) (hPatsNe : pats ≠ []) :
    matchAtomSegs (a :: pats) (s :: segs) =
      (matchAtomSegs [a] [s] && matchAtomSegs pats segs) := by
  cases a with
  | multiWildcard => exact absurd rfl ha
  | singleWildcard =>
    cases pats with
    | nil => exact absurd rfl hPatsNe
    | cons p ps => cases segs <;> simp [matchAtomSegs]
  | lit expected =>
    cases pats with
    | nil => exact absurd rfl hPatsNe
    | cons p ps =>
      have hp : p ≠ .multiWildcard := fun heq => hNoMulti (by simp [heq])
      cases p with
      | multiWildcard => exact absurd rfl hp
      | singleWildcard => cases segs <;> simp [matchAtomSegs]
      | lit e => cases segs <;> simp [matchAtomSegs]

private theorem match_one_slash_eq (p : AtomSeg) (hpValid : isValidAtom p = true)
    (hpNo : p ≠ .multiWildcard) (cs : List Char) :
    matchSegList [atomToRuntime p, .literal "/"] cs =
      match cs.dropWhile (· != '/') with
      | [] => none
      | _ :: rem =>
        let h := cs.takeWhile (· != '/')
        if matchAtomSegs [p] [String.ofList h] then some (h.length + 1, rem) else none := by
  have hDecomp := List.takeWhile_append_dropWhile (p := (· != '/')) (l := cs)
  cases p with
  | multiWildcard => exact absurd rfl hpNo
  | singleWildcard =>
    simp only [atomToRuntime, matchSegList, RuntimeSeg.matchSize, bind, Option.bind]
    have hDropH : cs.drop (cs.takeWhile (· != '/')).length = cs.dropWhile (· != '/') := by
      have : cs.drop (cs.takeWhile (· != '/')).length =
          (cs.takeWhile (· != '/') ++ cs.dropWhile (· != '/')).drop (cs.takeWhile (· != '/')).length := by
        rw [hDecomp]
      rw [this, List.drop_left]
    cases hEmpty : cs.takeWhile (· != '/') with
    | nil =>
      simp only [List.length_nil, beq_self_eq_true, ↓reduceIte]
      cases cs.dropWhile (· != '/') with
      | nil => rfl
      | cons c rem => simp [matchAtomSegs]
    | cons c0 cs0 =>
      rw [hEmpty, List.length_cons] at hDropH
      have hStrNe : (String.ofList (c0 :: cs0) != "") = true := by simp [bne, String.ext_iff]
      simp only [List.length_cons, matchAtomSegs, hStrNe, Bool.and_true]
      cases hd : cs.dropWhile (· != '/') with
      | nil =>
        have hEmp : cs.drop (cs0.length + 1) = [] := by rw [hDropH, hd]
        simp [hEmp]
      | cons c rem =>
        have hc : c = '/' := dropWhile_ne_slash_head hd
        subst hc
        have hNext : cs.drop (cs0.length + 1) = '/' :: rem := by rw [hDropH, hd]
        have hNz : (cs0.length + 1 == 0) = false := by simp
        have hSlashToList : ("/" : String).toList = ['/'] := rfl
        simp only [hNz, Bool.false_eq_true, ↓reduceIte, hNext, hSlashToList,
          List.isPrefixOf, beq_self_eq_true, List.length_singleton, pure]
        simp [Nat.add_comm]
  | lit expected =>
    rcases ValidLitNoSlash hpValid with ⟨hExpNe, hNoSlash⟩
    have hLitNoSlash : ∀ c ∈ expected.toList, c ≠ '/' := fun c hc => (hNoSlash c hc).1
    have hSlashEq : (expected ++ "/").toList = expected.toList ++ ['/'] := by
      simp [String.toList_append]
    simp only [atomToRuntime]
    rw [matchSegList_lit_append]
    simp only [matchSegList, RuntimeSeg.matchSize, bind, Option.bind, List.isPrefixOf_iff_prefix,
      hSlashEq]
    cases hd : cs.dropWhile (· != '/') with
    | nil =>
      simp only []
      by_cases hPref : (expected.toList ++ ['/']) <+: cs
      · exfalso
        rcases hPref with ⟨tail, hEq⟩
        have hIn : '/' ∈ cs := by rw [← hEq]; simp
        rw [← hDecomp, hd, List.append_nil] at hIn
        have hBad := takeWhile_mem hIn
        simp at hBad
      · simp [hPref]
    | cons c rem =>
      have hc : c = '/' := dropWhile_ne_slash_head hd
      subst hc
      let h := cs.takeWhile (· != '/')
      have hSum : cs = h ++ '/' :: rem := by
        change cs = cs.takeWhile (· != '/') ++ '/' :: rem
        simpa [hd] using hDecomp.symm
      by_cases hMatch : String.ofList h = expected
      · have hEqH : h = expected.toList := by
          rw [← String.toList_ofList (l := h), hMatch]
        have hPref : (expected.toList ++ ['/']) <+: cs := by
          rw [hSum, hEqH]
          exact ⟨rem, by simp [List.append_assoc]⟩
        have hTakeRem : (expected.toList ++ '/' :: rem).takeWhile (· != '/') = expected.toList := by
          rw [(takeWhile_dropWhile_append_of_no_slash expected.toList ('/' :: rem) hLitNoSlash).1]
          simp
        simp only [hPref, ↓reduceIte, pure]
        rw [hSum, hEqH]
        simp only [List.length_append, List.length_singleton, matchAtomSegs,
          hTakeRem, String.ofList_toList, beq_self_eq_true, Bool.and_self, ↓reduceIte]
        have hDrop1 : (expected.toList ++ '/' :: rem).drop (expected.toList.length + 1) = rem := by
          simp
        rw [hDrop1]
      · have hNePref : ¬ (expected.toList ++ ['/']) <+: cs := by
          intro hPref
          rcases hPref with ⟨tail, hEq⟩
          have hSplit := takeWhile_dropWhile_append_of_no_slash expected.toList ('/' :: tail) hLitNoSlash
          have hTake : h = expected.toList := by
            have hh : cs.takeWhile (· != '/') = h := rfl
            rw [← hEq] at hh
            simpa [hSplit.1, List.takeWhile_cons] using hh.symm
          apply hMatch
          rw [hTake, String.ofList_toList]
        have hBeqFalse : (String.ofList (cs.takeWhile (· != '/')) == expected) = false :=
          beq_false_of_ne hMatch
        simp [hNePref, matchAtomSegs, hBeqFalse]

private theorem match_one_exact_eq (p : AtomSeg) (hpValid : isValidAtom p = true)
    (hpNo : p ≠ .multiWildcard) (cs : List Char) :
    (match matchSegList [atomToRuntime p] cs with
      | some (n, []) => some n
      | _ => none) =
      (match cs.dropWhile (· != '/') with
        | [] =>
          if cs != [] && matchAtomSegs [p] [String.ofList cs] then some cs.length else none
        | _ :: _ => none) := by
  have hDecomp := List.takeWhile_append_dropWhile (p := (· != '/')) (l := cs)
  cases p with
  | multiWildcard => exact absurd rfl hpNo
  | singleWildcard =>
    simp only [atomToRuntime, matchSegList, RuntimeSeg.matchSize, bind, Option.bind, pure]
    have hDropH : cs.drop (cs.takeWhile (· != '/')).length = cs.dropWhile (· != '/') := by
      have : cs.drop (cs.takeWhile (· != '/')).length =
          (cs.takeWhile (· != '/') ++ cs.dropWhile (· != '/')).drop (cs.takeWhile (· != '/')).length := by
        rw [hDecomp]
      rw [this, List.drop_left]
    cases hEmpty : cs.takeWhile (· != '/') with
    | nil =>
      simp only [List.length_nil, beq_self_eq_true, ↓reduceIte]
      cases hd : cs.dropWhile (· != '/') with
      | nil =>
        have hCsNil : cs = [] := by simpa [hEmpty, hd] using hDecomp.symm
        simp [hCsNil]
      | cons _ _ => rfl
    | cons c0 cs0 =>
      cases hd : cs.dropWhile (· != '/') with
      | nil =>
        have hCsEq : cs = c0 :: cs0 := by simpa [hEmpty, hd] using hDecomp.symm
        rw [hEmpty, List.length_cons] at hDropH
        simp [hCsEq, matchAtomSegs, bne]
      | cons _ _ =>
        rw [hEmpty, List.length_cons] at hDropH
        simp [hDropH, hd]
  | lit expected =>
    rcases ValidLitNoSlash hpValid with ⟨hExpNe, hNoSlash⟩
    have hLitNoSlash : ∀ c ∈ expected.toList, c ≠ '/' := fun c hc => (hNoSlash c hc).1
    simp only [atomToRuntime, matchSegList, RuntimeSeg.matchSize, bind, Option.bind, pure,
      List.isPrefixOf_iff_prefix]
    cases hd : cs.dropWhile (· != '/') with
    | nil =>
      by_cases hEqExp : cs = expected.toList
      · have hPref : expected.toList <+: cs := by rw [hEqExp]; exact List.prefix_rfl
        simp [hEqExp, List.drop_length, matchAtomSegs, String.ofList_toList, hExpNe]
      · by_cases hPref : expected.toList <+: cs
        · have hPrefCopy := hPref
          rcases hPref with ⟨tail, hTailEq⟩
          have hTailNe : tail ≠ [] := by
            intro ht; subst ht; simp only [List.append_nil] at hTailEq; exact hEqExp hTailEq.symm
          have hDropNe : cs.drop expected.toList.length ≠ [] := by
            rw [← hTailEq, List.drop_left]; exact hTailNe
          cases hDr : cs.drop expected.toList.length with
          | nil => exact absurd hDr hDropNe
          | cons _ _ =>
            simp only [hPrefCopy, ↓reduceIte, hDr]
            have hMatchFalse : (String.ofList cs == expected) = false := by
              apply beq_false_of_ne
              intro hs
              apply hEqExp
              rw [← hs, String.toList_ofList]
            simp [matchAtomSegs, hMatchFalse]
        · have hMatchFalse : (String.ofList cs == expected) = false := by
            apply beq_false_of_ne
            intro hs
            apply hEqExp
            rw [← hs, String.toList_ofList]
          simp [hPref, matchAtomSegs, hMatchFalse]
    | cons c rem =>
      have hc : c = '/' := dropWhile_ne_slash_head hd
      subst hc
      by_cases hPref : expected.toList <+: cs
      · have hPrefCopy := hPref
        rcases hPref with ⟨tail, hTailEq⟩
        have hIn : '/' ∈ cs := by rw [← hDecomp, hd]; simp
        have hNotIn : '/' ∉ expected.toList := fun hm => hLitNoSlash '/' hm rfl
        have hTailIn : '/' ∈ tail := by
          rw [← hTailEq, List.mem_append] at hIn
          rcases hIn with hBad | hOk
          · exact absurd hBad hNotIn
          · exact hOk
        have hDropNe : cs.drop expected.toList.length ≠ [] := by
          rw [← hTailEq, List.drop_left]
          intro hNil; rw [hNil] at hTailIn; cases hTailIn
        cases hDr : cs.drop expected.toList.length with
        | nil => exact absurd hDr hDropNe
        | cons _ _ => simp [hPrefCopy, hDr]
      · simp [hPref]

/--
Shared characterization for literal+delimiter+`**`:
`expected` matches a prefix of `cs` followed by `[]`, `'/' :: _`, or `':' :: _`
iff the first slash-delimited segment of `cs` equals `expected` or has `expected ++ ":"` prefix.
-/
private theorem lit_trailingMulti_iff (expected : String)
    (hValid : isValidAtom (.lit expected) = true) (cs : List Char) :
    (∃ rem2, expected.toList ++ rem2 = cs ∧
        (rem2 = [] ∨ (∃ r, rem2 = '/' :: r) ∨ (∃ r, rem2 = ':' :: r))) ↔
      (let s0 := String.ofList (cs.takeWhile (· != '/'))
       s0 == expected || hasColonSuffix expected s0) = true := by
  rcases ValidLitNoSlash hValid with ⟨hExpNe, hNoChar⟩
  have hNoSlash : ∀ c ∈ expected.toList, c ≠ '/' := fun c hc => (hNoChar c hc).1
  have hDecomp := List.takeWhile_append_dropWhile (p := (· != '/')) (l := cs)
  constructor
  · rintro ⟨rem2, hEq, hRem⟩
    rcases hRem with rfl | ⟨r, rfl⟩ | ⟨r, rfl⟩
    · simp only [List.append_nil] at hEq
      have hTake : cs.takeWhile (· != '/') = expected.toList := by
        rw [← hEq]
        simpa using (takeWhile_dropWhile_append_of_no_slash expected.toList [] hNoSlash).1
      simp [hTake, String.ofList_toList]
    · have hTake : cs.takeWhile (· != '/') = expected.toList := by
        rw [← hEq, (takeWhile_dropWhile_append_of_no_slash expected.toList ('/' :: r) hNoSlash).1]
        simp
      simp [hTake, String.ofList_toList]
    · have hSlashColon : (':' != '/') = true := by decide
      have hTake : cs.takeWhile (· != '/') = expected.toList ++ ':' :: r.takeWhile (· != '/') := by
        rw [← hEq, (takeWhile_dropWhile_append_of_no_slash expected.toList (':' :: r) hNoSlash).1]
        simp [hSlashColon]
      rw [Bool.or_eq_true]
      right
      unfold hasColonSuffix
      simp only [hTake, String.toList_ofList, List.length_append, List.length_cons,
        Bool.and_eq_true, decide_eq_true_eq, beq_iff_eq]
      refine ⟨by omega, ?_⟩
      have hSplitList : expected.toList ++ ':' :: r.takeWhile (· != '/') =
          (expected.toList ++ [':']) ++ r.takeWhile (· != '/') := by
        simp [List.append_assoc]
      have hLenEq : expected.toList.length + 1 = (expected.toList ++ [':']).length := by
        simp
      rw [hSplitList, hLenEq]
      exact List.take_left
  · intro hMatch
    simp only [Bool.or_eq_true, beq_iff_eq] at hMatch
    rcases hMatch with hEq | hColon
    · have hTake : cs.takeWhile (· != '/') = expected.toList := by
        rw [← String.toList_ofList (l := cs.takeWhile (· != '/')), hEq]
      cases hd : cs.dropWhile (· != '/') with
      | nil =>
        refine ⟨[], ?_, Or.inl rfl⟩
        simpa [hTake, hd] using hDecomp
      | cons c r =>
        have hc : c = '/' := dropWhile_ne_slash_head hd
        subst hc
        refine ⟨'/' :: r, ?_, Or.inr (Or.inl ⟨r, rfl⟩)⟩
        simpa [hTake, hd] using hDecomp
    · unfold hasColonSuffix at hColon
      simp only [String.toList_ofList, Bool.and_eq_true, decide_eq_true_eq, beq_iff_eq] at hColon
      rcases hColon with ⟨_, hPref⟩
      let h := cs.takeWhile (· != '/')
      have hEqPref : h.take (expected.toList.length + 1) ++ h.drop (expected.toList.length + 1) = h :=
        List.take_append_drop _ _
      rw [hPref] at hEqPref
      have hRest : expected.toList ++ ':' :: (h.drop (expected.toList.length + 1) ++
          cs.dropWhile (· != '/')) = cs := by
        calc expected.toList ++ ':' :: (h.drop (expected.toList.length + 1) ++ cs.dropWhile (· != '/'))
            = (expected.toList ++ [':']) ++ h.drop (expected.toList.length + 1) ++ cs.dropWhile (· != '/') := by
              simp [List.append_assoc]
          _ = h ++ cs.dropWhile (· != '/') := by rw [hEqPref]
          _ = cs := hDecomp
      exact ⟨':' :: (h.drop (expected.toList.length + 1) ++ cs.dropWhile (· != '/')),
        hRest, Or.inr (Or.inr ⟨_, rfl⟩)⟩

/-! ### Theorem Group A: Multi-segment Slash-Terminated Matching (`prefixSegs` & pre-suffix `matchSegs`) -/

private theorem matchUnmergedSlashEnd_spec (pats : List AtomSeg) (hNe : pats ≠ [])
    (hValid : pats.all isValidAtom = true) (hNoMulti : AtomSeg.multiWildcard ∉ pats)
    (cs : List Char) :
    (match matchSegList (emitUnmerged pats false false ++ [.literal "/"]) cs with
      | none => none
      | some (n, rem) => some (cs.take (n - 1), splitSegsOf rem)) =
    (let segs := splitSegsOf cs
     if pats.length < segs.length && matchAtomSegs pats (segs.take pats.length) then
       some ((joinStringsWith "/" (segs.take pats.length)).toList, segs.drop pats.length)
     else none) := by
  induction pats generalizing cs with
  | nil => exact absurd rfl hNe
  | cons p ps ih =>
    have hpValid : isValidAtom p = true := List.all_eq_true.mp hValid p (by simp)
    have hpNo : p ≠ .multiWildcard := fun h => hNoMulti (by simp [h])
    cases ps with
    | nil =>
      rw [emitUnmerged_singleton p hpNo]
      simp only [List.singleton_append]
      rw [match_one_slash_eq p hpValid hpNo cs]
      rw [splitSegsOf_eq_take_drop cs]
      cases hd : cs.dropWhile (· != '/') with
      | nil => simp
      | cons c rem =>
        have hc : c = '/' := dropWhile_ne_slash_head hd
        subst hc
        have hSum := List.takeWhile_append_dropWhile (p := (· != '/')) (l := cs)
        rw [hd] at hSum
        have hSegsNe : splitSegsOf rem ≠ [] := splitSegsOf_ne_nil rem
        have hLenGt : 1 < (String.ofList (cs.takeWhile (· != '/')) :: splitSegsOf rem).length := by
          cases hS : splitSegsOf rem with
          | nil => exact absurd hS hSegsNe
          | cons _ _ => simp
        simp only [List.length_singleton, hLenGt, List.take_succ_cons,
          List.take_zero, List.drop_succ_cons, List.drop_zero]
        by_cases hm : matchAtomSegs [p] [String.ofList (cs.takeWhile (· != '/'))] = true
        · simp only [hm, ↓reduceIte, joinStringsWith, String.toList_ofList]
          have hTake : cs.take (cs.takeWhile (· != '/')).length = cs.takeWhile (· != '/') := by
            have hTmp := @List.take_left Char (cs.takeWhile (· != '/')) ('/' :: rem)
            rwa [hSum] at hTmp
          simp [hTake]
        · simp [hm]
    | cons q qs =>
      have hPsValid : (q :: qs).all isValidAtom = true := by
        rw [List.all_eq_true] at hValid ⊢
        intro a ha; exact hValid a (List.mem_cons_of_mem p ha)
      have hPsNo : AtomSeg.multiWildcard ∉ (q :: qs) :=
        fun h => hNoMulti (List.mem_cons_of_mem p h)
      have hqNo : q ≠ .multiWildcard := fun h => hPsNo (by simp [h])
      rw [emitUnmerged_cons_slash p q qs hpNo hqNo, List.append_assoc, matchSegList_append,
        match_one_slash_eq p hpValid hpNo cs]
      rw [splitSegsOf_eq_take_drop cs]
      cases hd : cs.dropWhile (· != '/') with
      | nil => simp
      | cons c rem =>
        have hc : c = '/' := dropWhile_ne_slash_head hd
        subst hc
        have hSum := List.takeWhile_append_dropWhile (p := (· != '/')) (l := cs)
        rw [hd] at hSum
        let s0 := String.ofList (cs.takeWhile (· != '/'))
        by_cases hmp : matchAtomSegs [p] [s0] = true
        · have hmp' : matchAtomSegs [p] [String.ofList (cs.takeWhile (· != '/'))] = true := hmp
          simp only [hmp', ↓reduceIte]
          have ihRem := ih (by simp) hPsValid hPsNo rem
          have hMatchCons : matchAtomSegs (p :: q :: qs)
              (String.ofList (cs.takeWhile (· != '/')) :: (splitSegsOf rem).take (qs.length + 1)) =
              matchAtomSegs (q :: qs) ((splitSegsOf rem).take (qs.length + 1)) := by
            rw [matchAtomSegs_cons_no_multi p (q :: qs)
              (String.ofList (cs.takeWhile (· != '/'))) _ hNoMulti hpNo (by simp)]
            have hmp0 : matchAtomSegs [p] [String.ofList (cs.takeWhile (· != '/'))] = true := hmp
            simp [hmp0]
          cases hSub : matchSegList (emitUnmerged (q :: qs) false false ++ [.literal "/"]) rem with
          | none =>
            simp only [hSub] at ihRem ⊢
            simp only [List.length_cons, Nat.add_lt_add_iff_right, List.take_succ_cons,
              List.drop_succ_cons]
            split at ihRem
            · simp at ihRem
            · next hCond =>
              simp only [Bool.and_eq_true] at hCond
              rw [hMatchCons]
              split
              · next hC2 =>
                simp only [Bool.and_eq_true] at hC2
                exact absurd ⟨hC2.1, hC2.2⟩ hCond
              · rfl
          | some pSub =>
            rcases pSub with ⟨ny, remy⟩
            simp only [hSub] at ihRem
            have hNyPos : 1 ≤ ny := by
              rw [matchSegList_append] at hSub
              cases hA : matchSegList (emitUnmerged (q :: qs) false false) rem with
              | none => simp [hA] at hSub
              | some pA =>
                rcases pA with ⟨nx, remx⟩
                simp only [hA] at hSub
                cases remx with
                | nil => simp [matchSegList, RuntimeSeg.matchSize] at hSub
                | cons c cs' =>
                  by_cases hc : c = '/'
                  · subst hc
                    have hLitOne : matchSegList [RuntimeSeg.literal "/"] ('/' :: cs') = some (1, cs') := by
                      rfl
                    simp only [hLitOne] at hSub
                    cases hSub
                    omega
                  · have hLitNone : matchSegList [RuntimeSeg.literal "/"] (c :: cs') = none := by
                      simp [matchSegList, RuntimeSeg.matchSize, List.isPrefixOf, Ne.symm hc]
                    simp [hLitNone] at hSub
            split at ihRem
            · simp only [Option.some.injEq, Prod.mk.injEq] at ihRem
              rcases ihRem with ⟨hTakeRem, hDropRem⟩
              next hCond =>
              simp only [Bool.and_eq_true] at hCond
              simp only [List.length_cons, Nat.add_lt_add_iff_right, List.take_succ_cons,
                List.drop_succ_cons]
              rw [hMatchCons]
              have hCond1 : qs.length + 1 < (splitSegsOf rem).length :=
                decide_eq_true_eq.mp hCond.1
              have hCond2 : matchAtomSegs (q :: qs) ((splitSegsOf rem).take (qs.length + 1)) = true := hCond.2
              simp only [hCond1, decide_true, hCond2, Bool.and_self, ↓reduceIte,
                Option.some.injEq, Prod.mk.injEq]
              refine ⟨?_, hDropRem⟩
              have hArith : (cs.takeWhile (· != '/')).length + 1 + ny - 1 =
                  (cs.takeWhile (· != '/')).length + (ny - 1 + 1) := by omega
              have hDropCs : cs.drop (cs.takeWhile (· != '/')).length = '/' :: rem := by
                have hTmp := @List.drop_left Char (cs.takeWhile (· != '/')) ('/' :: rem)
                rwa [hSum] at hTmp
              have hTakeSum : cs.take ((cs.takeWhile (· != '/')).length + (ny - 1 + 1)) =
                  cs.take (cs.takeWhile (· != '/')).length ++
                    '/' :: rem.take (ny - 1) := by
                rw [List.take_add, hDropCs, List.take_succ_cons]
              have hTakeFst : cs.take (cs.takeWhile (· != '/')).length = cs.takeWhile (· != '/') := by
                have hTmp := @List.take_left Char (cs.takeWhile (· != '/')) ('/' :: rem)
                rwa [hSum] at hTmp
              have hTakeNe : (splitSegsOf rem).take (qs.length + 1) ≠ [] := by
                cases hs : splitSegsOf rem with
                | nil => exact absurd hs (splitSegsOf_ne_nil rem)
                | cons _ _ => simp
              rw [hArith, hTakeSum, hTakeFst, hTakeRem,
                joinStringsWith_cons_ne_nil (String.ofList (cs.takeWhile (· != '/'))) _ hTakeNe,
                String.toList_append, String.toList_append]
              simp [String.toList_ofList, List.append_assoc]
            · cases ihRem
        · have hmp' : matchAtomSegs [p] [String.ofList (cs.takeWhile (· != '/'))] = false :=
            Bool.eq_false_iff.mpr hmp
          simp only [hmp', Bool.false_eq_true, ↓reduceIte, List.length_cons,
            Nat.add_lt_add_iff_right, List.take_succ_cons]
          rw [matchAtomSegs_cons_no_multi p (q :: qs)
            (String.ofList (cs.takeWhile (· != '/'))) _ hNoMulti hpNo (by simp)]
          simp [hmp']

/-! ### Theorem Group B: Exact Full-String Matching (`suffixSegs = []` without `**`, or exact `suffixSegs`) -/

private theorem matchAtomSegs_singleton_valid_ne_nil (p : AtomSeg)
    (hpValid : isValidAtom p = true) (hpNo : p ≠ .multiWildcard) :
    matchAtomSegs [p] [""] = false := by
  cases p with
  | multiWildcard => exact absurd rfl hpNo
  | singleWildcard => simp [matchAtomSegs]
  | lit s =>
    rcases ValidLitNoSlash hpValid with ⟨hNe, _⟩
    have hBeq : ("" == s) = false := beq_false_of_ne hNe.symm
    simp [matchAtomSegs, hBeq]

private theorem matchAtomSegs_valid_join_ne_empty :
    ∀ (pats : List AtomSeg) (segs : List String),
      pats ≠ [] → pats.all isValidAtom = true → AtomSeg.multiWildcard ∉ pats →
      segs.length = pats.length → matchAtomSegs pats segs = true →
      joinStringsWith "/" segs ≠ ""
  | [], _, hNe, _, _, _, _ => absurd rfl hNe
  | [p], segs, _, hValid, hNoMulti, hLen, hMatch => by
    cases segs with
    | nil => simp at hLen
    | cons s ss =>
      cases ss with
      | cons _ _ => simp at hLen
      | nil =>
        simp only [joinStringsWith]
        intro hEmp
        subst hEmp
        have hpValid : isValidAtom p = true := List.all_eq_true.mp hValid p (by simp)
        have hpNo : p ≠ .multiWildcard := fun h => hNoMulti (by simp [h])
        rw [matchAtomSegs_singleton_valid_ne_nil p hpValid hpNo] at hMatch
        cases hMatch
  | _ :: _ :: _, segs, _, _, _, hLen, _ => by
    cases segs with
    | nil => simp at hLen
    | cons s ss =>
      cases ss with
      | nil => simp at hLen
      | cons s2 ss2 =>
        intro hEmp
        have hList : (s ++ "/" ++ joinStringsWith "/" (s2 :: ss2)).toList = [] := by
          rw [← joinStringsWith_cons_ne_nil s (s2 :: ss2) (by simp), hEmp]
          rfl
        simp [String.toList_append] at hList

private theorem emitUnmerged_cons_pos (pats : List AtomSeg) (hNe : pats ≠ [])
    (hValid : pats.all isValidAtom = true) (hNoMulti : AtomSeg.multiWildcard ∉ pats)
    (cs : List Char) {n : Nat} {rem : List Char}
    (hMatch : matchSegList (emitUnmerged pats false false) cs = some (n, rem)) :
    1 ≤ n ∧ cs ≠ [] := by
  cases pats with
  | nil => exact absurd rfl hNe
  | cons p ps =>
    have hpValid : isValidAtom p = true := List.all_eq_true.mp hValid p (by simp)
    have hpNo : p ≠ .multiWildcard := fun h => hNoMulti (by simp [h])
    have hHeadPos : ∀ (rest : List RuntimeSeg),
        matchSegList (atomToRuntime p :: rest) cs = some (n, rem) → 1 ≤ n ∧ cs ≠ [] := by
      intro rest hM
      simp only [matchSegList, bind, Option.bind] at hM
      cases p with
      | multiWildcard => exact absurd rfl hpNo
      | singleWildcard =>
        simp only [atomToRuntime, RuntimeSeg.matchSize] at hM
        cases hZero : ((cs.takeWhile (· != '/')).length == 0) with
        | true => simp [hZero] at hM
        | false =>
          simp only [hZero, Bool.false_eq_true, ↓reduceIte] at hM
          have hTakePos : 0 < (cs.takeWhile (· != '/')).length := by
            cases hL : (cs.takeWhile (· != '/')).length with
            | zero => simp [hL] at hZero
            | succ _ => omega
          have hCsNe : cs ≠ [] := by
            intro hNil
            subst hNil
            simp at hTakePos
          cases hRest : matchSegList rest (cs.drop (cs.takeWhile (· != '/')).length) with
          | none => simp [hRest] at hM
          | some pR =>
            rcases pR with ⟨m, rem'⟩
            simp only [hRest, pure, Option.some.injEq, Prod.mk.injEq] at hM
            exact ⟨by omega, hCsNe⟩
      | lit s =>
        rcases ValidLitNoSlash hpValid with ⟨hsNe, _⟩
        simp only [atomToRuntime, RuntimeSeg.matchSize] at hM
        cases hPrefB : List.isPrefixOf s.toList cs with
        | false => simp [hPrefB] at hM
        | true =>
          simp only [hPrefB, ↓reduceIte] at hM
          rw [List.isPrefixOf_iff_prefix] at hPrefB
          rcases hPrefB with ⟨tail, hEq⟩
          have hsPos : 1 ≤ s.toList.length := by
            cases hSL : s.toList with
            | nil => exact absurd (String.ext (by simp [hSL])) hsNe
            | cons _ _ => simp
          have hCsNe : cs ≠ [] := by
            intro hNil
            subst hNil
            simp at hEq
            exact absurd (String.ext (by simp [hEq.1])) hsNe
          cases hRest : matchSegList rest (cs.drop s.toList.length) with
          | none => simp [hRest] at hM
          | some pR =>
            rcases pR with ⟨m, rem'⟩
            simp only [hRest, pure, Option.some.injEq, Prod.mk.injEq] at hM
            exact ⟨by omega, hCsNe⟩
    cases ps with
    | nil =>
      rw [emitUnmerged_singleton p hpNo] at hMatch
      exact hHeadPos [] hMatch
    | cons q qs =>
      have hqNo : q ≠ .multiWildcard := fun h => hNoMulti (by simp [h])
      rw [emitUnmerged_cons_slash p q qs hpNo hqNo] at hMatch
      exact hHeadPos _ hMatch

private theorem matchUnmergedExact_spec (pats : List AtomSeg) (hNe : pats ≠ [])
    (hValid : pats.all isValidAtom = true) (hNoMulti : AtomSeg.multiWildcard ∉ pats)
    (cs : List Char) :
    (match matchSegList (emitUnmerged pats false false) cs with
      | some (n, []) => if n == 0 then none else some (String.ofList (cs.take n))
      | _ => none) =
    (let segs := splitSegsOf cs
     if segs.length == pats.length && matchAtomSegs pats segs then
       let cap := joinStringsWith "/" segs
       if cap == "" then none else some cap
     else none) := by
  induction pats generalizing cs with
  | nil => exact absurd rfl hNe
  | cons p ps ih =>
    have hpValid : isValidAtom p = true := List.all_eq_true.mp hValid p (by simp)
    have hpNo : p ≠ .multiWildcard := fun h => hNoMulti (by simp [h])
    cases ps with
    | nil =>
      rw [emitUnmerged_singleton p hpNo]
      have hOne := match_one_exact_eq p hpValid hpNo cs
      rw [splitSegsOf_eq_take_drop cs]
      cases hd : cs.dropWhile (· != '/') with
      | cons c rem =>
        rw [hd] at hOne
        have hSegsNe : splitSegsOf rem ≠ [] := splitSegsOf_ne_nil rem
        have hLenNe : ((String.ofList (cs.takeWhile (· != '/')) :: splitSegsOf rem).length == 1) = false := by
          cases hs : splitSegsOf rem with
          | nil => exact absurd hs hSegsNe
          | cons _ _ => simp
        simp only [List.length_singleton, hLenNe, Bool.false_and, Bool.false_eq_true, ↓reduceIte]
        cases hMatch : matchSegList [atomToRuntime p] cs with
        | none => rfl
        | some pRes =>
          rcases pRes with ⟨n, rem'⟩
          cases rem' with
          | cons _ _ => rfl
          | nil => rw [hMatch] at hOne; cases hOne
      | nil =>
        rw [hd] at hOne
        have hSum := List.takeWhile_append_dropWhile (p := (· != '/')) (l := cs)
        rw [hd, List.append_nil] at hSum
        simp only [hSum, List.length_singleton, beq_self_eq_true, Bool.true_and, joinStringsWith]
        cases hMatch : matchSegList [atomToRuntime p] cs with
        | none =>
          rw [hMatch] at hOne
          dsimp only at hOne
          split at hOne
          · cases hOne
          · next hCond =>
            by_cases hm : matchAtomSegs [p] [String.ofList cs] = true
            · by_cases hCsNil : cs = []
              · subst hCsNil
                rw [matchAtomSegs_singleton_valid_ne_nil p hpValid hpNo] at hm
                cases hm
              · exfalso
                apply hCond
                simp [bne, hCsNil, hm]
            · simp [hm]
        | some pRes =>
          rcases pRes with ⟨n, rem'⟩
          rw [hMatch] at hOne
          cases rem' with
          | cons r rs =>
            dsimp only at hOne
            split at hOne
            · cases hOne
            · next hCond =>
              by_cases hm : matchAtomSegs [p] [String.ofList cs] = true
              · by_cases hCsNil : cs = []
                · subst hCsNil
                  rw [matchAtomSegs_singleton_valid_ne_nil p hpValid hpNo] at hm
                  cases hm
                · exfalso
                  apply hCond
                  simp [bne, hCsNil, hm]
              · simp [hm]
          | nil =>
            dsimp only at hOne
            split at hOne
            · next hCond =>
              simp only [Option.some.injEq] at hOne
              subst hOne
              simp only [Bool.and_eq_true, bne_iff_ne, ne_eq] at hCond
              have hm : matchAtomSegs [p] [String.ofList cs] = true := hCond.2
              have hLenNz : (cs.length == 0) = false := by
                cases hCs : cs with
                | nil => exact absurd hCs hCond.1
                | cons _ _ => simp
              have hStrNe : (String.ofList cs == "") = false := by
                apply beq_false_of_ne
                intro hE
                apply hCond.1
                rw [← String.toList_ofList (l := cs), hE]
                rfl
              simp [hLenNz, List.take_length, hm, hStrNe]
            · cases hOne
    | cons q qs =>
      have hPsValid : (q :: qs).all isValidAtom = true := by
        rw [List.all_eq_true] at hValid ⊢
        intro a ha; exact hValid a (List.mem_cons_of_mem p ha)
      have hPsNo : AtomSeg.multiWildcard ∉ (q :: qs) :=
        fun h => hNoMulti (List.mem_cons_of_mem p h)
      have hqNo : q ≠ .multiWildcard := fun h => hPsNo (by simp [h])
      rw [emitUnmerged_cons_slash p q qs hpNo hqNo, matchSegList_append,
        match_one_slash_eq p hpValid hpNo cs, splitSegsOf_eq_take_drop cs]
      cases hd : cs.dropWhile (· != '/') with
      | nil => simp
      | cons c rem =>
        have hc : c = '/' := dropWhile_ne_slash_head hd
        subst hc
        have hSum := List.takeWhile_append_dropWhile (p := (· != '/')) (l := cs)
        rw [hd] at hSum
        let s0 := String.ofList (cs.takeWhile (· != '/'))
        have ihRem := ih (by simp) hPsValid hPsNo rem
        by_cases hmp : matchAtomSegs [p] [s0] = true
        · have hmp' : matchAtomSegs [p] [String.ofList (cs.takeWhile (· != '/'))] = true := hmp
          simp only [hmp', ↓reduceIte]
          have hCapNe : (joinStringsWith "/" (s0 :: splitSegsOf rem) == "") = false := by
            apply beq_false_of_ne
            rw [joinStringsWith_cons_ne_nil s0 (splitSegsOf rem) (splitSegsOf_ne_nil rem)]
            intro hE
            have hL : (s0 ++ "/" ++ joinStringsWith "/" (splitSegsOf rem)).toList = [] := by
              rw [hE]; rfl
            simp [String.toList_append] at hL
          have hMatchCons : matchAtomSegs (p :: q :: qs) (s0 :: splitSegsOf rem) =
              matchAtomSegs (q :: qs) (splitSegsOf rem) := by
            rw [matchAtomSegs_cons_no_multi p (q :: qs) s0 _ hNoMulti hpNo (by simp)]
            have : matchAtomSegs [p] [s0] = true := hmp
            simp [this]
          cases hSub : matchSegList (emitUnmerged (q :: qs) false false) rem with
          | none =>
            simp only [hSub] at ihRem ⊢
            have hCondFalse :
                ((splitSegsOf rem).length == (q :: qs).length &&
                  matchAtomSegs (q :: qs) (splitSegsOf rem)) = false := by
              cases hC : (splitSegsOf rem).length == (q :: qs).length &&
                  matchAtomSegs (q :: qs) (splitSegsOf rem)
              · rfl
              · exfalso
                simp only [hC, ↓reduceIte] at ihRem
                simp only [Bool.and_eq_true, beq_iff_eq] at hC
                have hNeRem := matchAtomSegs_valid_join_ne_empty (q :: qs) (splitSegsOf rem)
                  (by simp) hPsValid hPsNo hC.1 hC.2
                have hBeqF : (joinStringsWith "/" (splitSegsOf rem) == "") = false :=
                  beq_false_of_ne hNeRem
                simp [hBeqF] at ihRem
            simp only [List.length_cons]
            have hMatchC' : matchAtomSegs (p :: q :: qs)
                (String.ofList (cs.takeWhile (· != '/')) :: splitSegsOf rem) =
                matchAtomSegs (q :: qs) (splitSegsOf rem) := hMatchCons
            rw [hMatchC']
            split
            · next hT =>
              simp only [Nat.add_right_cancel_iff, Bool.and_eq_true, beq_iff_eq] at hT
              have hTrue : ((splitSegsOf rem).length == (q :: qs).length &&
                  matchAtomSegs (q :: qs) (splitSegsOf rem)) = true := by
                simp [hT.1, hT.2]
              rw [hTrue] at hCondFalse
              cases hCondFalse
            · rfl
          | some pSub =>
            rcases pSub with ⟨ny, remy⟩
            simp only [hSub] at ihRem ⊢
            cases remy with
            | cons r rs =>
              have hCondFalse :
                  ((splitSegsOf rem).length == (q :: qs).length &&
                    matchAtomSegs (q :: qs) (splitSegsOf rem)) = false := by
                cases hC : (splitSegsOf rem).length == (q :: qs).length &&
                    matchAtomSegs (q :: qs) (splitSegsOf rem)
                · rfl
                · exfalso
                  simp only [hC, ↓reduceIte] at ihRem
                  simp only [Bool.and_eq_true, beq_iff_eq] at hC
                  have hNeRem := matchAtomSegs_valid_join_ne_empty (q :: qs) (splitSegsOf rem)
                    (by simp) hPsValid hPsNo hC.1 hC.2
                  have hBeqF : (joinStringsWith "/" (splitSegsOf rem) == "") = false :=
                    beq_false_of_ne hNeRem
                  simp [hBeqF] at ihRem
              simp only [List.length_cons]
              have hMatchC' : matchAtomSegs (p :: q :: qs)
                  (String.ofList (cs.takeWhile (· != '/')) :: splitSegsOf rem) =
                  matchAtomSegs (q :: qs) (splitSegsOf rem) := hMatchCons
              rw [hMatchC']
              split
              · next hT =>
                simp only [Nat.add_right_cancel_iff, Bool.and_eq_true, beq_iff_eq] at hT
                have hTrue : ((splitSegsOf rem).length == (q :: qs).length &&
                    matchAtomSegs (q :: qs) (splitSegsOf rem)) = true := by
                  simp [hT.1, hT.2]
                rw [hTrue] at hCondFalse
                cases hCondFalse
              · rfl
            | nil =>
              have hNz1 : ((cs.takeWhile (· != '/')).length + 1 + ny == 0) = false := by simp
              have hNyPos : 1 ≤ ny := (emitUnmerged_cons_pos (q :: qs) (by simp) hPsValid hPsNo rem hSub).1
              have hNzNy : (ny == 0) = false := by
                cases ny with
                | zero => omega
                | succ _ => simp
              simp only [hNz1, hNzNy, Bool.false_eq_true, ↓reduceIte] at ihRem ⊢
              rcases matchSegList_drop_eq hSub with ⟨hDrop, hLe⟩
              have hEqNy : ny = rem.length :=
                Nat.le_antisymm hLe (List.drop_eq_nil_iff.mp hDrop.symm)
              subst hEqNy
              split at ihRem
              · next hCond =>
                simp only [List.length_cons, beq_iff_eq, Bool.and_eq_true, Nat.add_right_cancel_iff] at hCond ⊢
                have hMatchC' : matchAtomSegs (p :: q :: qs)
                    (String.ofList (cs.takeWhile (· != '/')) :: splitSegsOf rem) =
                    matchAtomSegs (q :: qs) (splitSegsOf rem) := hMatchCons
                simp only [hMatchC', hCond.1, hCond.2, and_self, ↓reduceIte]
                have hAllLen : cs.take ((cs.takeWhile (· != '/')).length + 1 + rem.length) = cs := by
                  have hTotLen : (cs.takeWhile (· != '/')).length + 1 + rem.length = cs.length := by
                    have hL := congrArg List.length hSum
                    simp only [List.length_append, List.length_cons] at hL
                    omega
                  rw [hTotLen, List.take_length]
                have hLhsCap : cs.take ((cs.takeWhile (· != '/')).length + 1 + rem.length) =
                    cs.takeWhile (· != '/') ++ '/' :: rem := by
                  rw [hAllLen, hSum]
                rw [joinStringsWith_cons_ne_nil (String.ofList (cs.takeWhile (· != '/')))
                    (splitSegsOf rem) (splitSegsOf_ne_nil rem),
                  string_ofList_join_splitSegs]
                simp [hLhsCap, String.ext_iff, String.toList_append, String.toList_ofList]
              · cases ihRem
        · have hmp' : matchAtomSegs [p] [String.ofList (cs.takeWhile (· != '/'))] = false :=
            Bool.eq_false_iff.mpr hmp
          simp only [hmp', Bool.false_eq_true, ↓reduceIte]
          rw [matchAtomSegs_cons_no_multi p (q :: qs)
            (String.ofList (cs.takeWhile (· != '/'))) _ hNoMulti hpNo (by simp)]
          simp [hmp']

/-! ### Theorem Group C1: Terminal `**` Matching (`matchSegs` or `suffixSegs` ending in `**`) -/

private theorem matchAtomSegs_length_of_no_multi :
    ∀ (pats : List AtomSeg) (segs : List String),
      AtomSeg.multiWildcard ∉ pats →
      matchAtomSegs pats segs = true →
      segs.length = pats.length
  | [], segs, _, h => by
    cases segs with
    | nil => rfl
    | cons _ _ => simp [matchAtomSegs] at h
  | p :: ps, segs, hNo, h => by
    have hp : p ≠ .multiWildcard := fun heq => hNo (by simp [heq])
    have hPsNo : AtomSeg.multiWildcard ∉ ps := fun hm => hNo (by simp [hm])
    cases ps with
    | nil =>
      cases p with
      | multiWildcard => exact absurd rfl hp
      | lit _ =>
        cases segs with
        | nil => simp [matchAtomSegs] at h
        | cons _ segRest =>
          cases segRest with
          | nil => rfl
          | cons _ _ => simp [matchAtomSegs] at h
      | singleWildcard =>
        cases segs with
        | nil => simp [matchAtomSegs] at h
        | cons _ segRest =>
          cases segRest with
          | nil => rfl
          | cons _ _ => simp [matchAtomSegs] at h
    | cons q qs =>
      cases segs with
      | nil => cases p <;> simp [matchAtomSegs] at h
      | cons s segRest =>
        rw [matchAtomSegs_cons_no_multi p (q :: qs) s segRest hNo hp (by simp)] at h
        simp only [Bool.and_eq_true] at h
        have ih := matchAtomSegs_length_of_no_multi (q :: qs) segRest hPsNo h.2
        simp [ih]

private theorem matchAtomSegs_cons_two_nonmulti (p q : AtomSeg) (rest : List AtomSeg)
    (s0 : String) (segRest : List String)
    (hp : p ≠ .multiWildcard) (hq : q ≠ .multiWildcard) :
    matchAtomSegs (p :: q :: rest) (s0 :: segRest) =
      (matchAtomSegs [p] [s0] && matchAtomSegs (q :: rest) segRest) := by
  cases p with
  | multiWildcard => exact absurd rfl hp
  | singleWildcard =>
    cases q with
    | multiWildcard => exact absurd rfl hq
    | singleWildcard => cases segRest <;> simp [matchAtomSegs]
    | lit _ => cases segRest <;> simp [matchAtomSegs]
  | lit _ =>
    cases q with
    | multiWildcard => exact absurd rfl hq
    | singleWildcard => cases segRest <;> simp [matchAtomSegs]
    | lit _ => cases segRest <;> simp [matchAtomSegs]

private theorem trailingMulti_matchSize_iff (cs : List Char) {n : Nat} :
    RuntimeSeg.trailingMultiWildcard.matchSize cs = some n ↔
      n = cs.length ∧ (cs = [] ∨ (∃ r, cs = '/' :: r) ∨ (∃ r, cs = ':' :: r)) := by
  cases cs with
  | nil => simp [RuntimeSeg.matchSize, eq_comm]
  | cons c r =>
    simp only [RuntimeSeg.matchSize]
    by_cases h1 : c = '/'
    · subst h1; simp [Option.some.injEq, eq_comm]
    · by_cases h2 : c = ':'
      · subst h2; simp [Option.some.injEq, eq_comm]
      · simp [h1, h2]

private theorem matchSegList_trailingMulti_iff (cs rem : List Char) (n : Nat) :
    matchSegList [.trailingMultiWildcard] cs = some (n, rem) ↔
      n = cs.length ∧ rem = [] ∧ (cs = [] ∨ (∃ r, cs = '/' :: r) ∨ (∃ r, cs = ':' :: r)) := by
  simp only [matchSegList, bind, Option.bind, pure]
  constructor
  · intro h
    cases hS : RuntimeSeg.trailingMultiWildcard.matchSize cs with
    | none => simp [hS] at h
    | some k =>
      simp only [hS, Option.some.injEq, Prod.mk.injEq] at h
      rcases (trailingMulti_matchSize_iff cs).mp hS with ⟨rfl, hCond⟩
      exact ⟨by omega, by simp [← h.2], hCond⟩
  · rintro ⟨rfl, rfl, hCond⟩
    have hS : RuntimeSeg.trailingMultiWildcard.matchSize cs = some cs.length :=
      (trailingMulti_matchSize_iff cs).mpr ⟨rfl, hCond⟩
    simp [hS]

private theorem match_one_trailing_eq (p : AtomSeg) (hpValid : isValidAtom p = true)
    (hpNo : p ≠ .multiWildcard) (cs : List Char) :
    matchSegList [atomToRuntime p, .trailingMultiWildcard] cs =
      if matchAtomSegs [p, .multiWildcard] (splitSegsOf cs) then some (cs.length, []) else none := by
  have hDecomp := List.takeWhile_append_dropWhile (p := (· != '/')) (l := cs)
  cases p with
  | multiWildcard => exact absurd rfl hpNo
  | singleWildcard =>
    rw [splitSegsOf_eq_take_drop cs]
    simp only [atomToRuntime, matchSegList, RuntimeSeg.matchSize, bind, Option.bind]
    cases hZero : ((cs.takeWhile (· != '/')).length == 0) with
    | true =>
      have hTakeNil : cs.takeWhile (· != '/') = [] := by
        cases ht : cs.takeWhile (· != '/') with
        | nil => rfl
        | cons _ _ => simp [ht] at hZero
      simp only [hTakeNil, ↓reduceIte, String.ofList_nil]
      cases cs.dropWhile (· != '/') with
      | nil => simp [matchAtomSegs]
      | cons _ _ => simp [matchAtomSegs]
    | false =>
      simp only [Bool.false_eq_true, ↓reduceIte]
      have hStrNe : (String.ofList (cs.takeWhile (· != '/')) != "") = true := by
        simp only [bne_iff_ne, ne_eq, String.ext_iff, String.toList_ofList]
        cases ht : cs.takeWhile (· != '/') with
        | nil => simp [ht] at hZero
        | cons _ _ => simp
      have hDropEq : cs.drop (cs.takeWhile (· != '/')).length = cs.dropWhile (· != '/') := by
        have hTmp := @List.drop_left Char (cs.takeWhile (· != '/')) (cs.dropWhile (· != '/'))
        rwa [hDecomp] at hTmp
      cases hd : cs.dropWhile (· != '/') with
      | nil =>
        have hSum : cs.takeWhile (· != '/') = cs := by simpa [hd] using hDecomp
        have hTrailing : RuntimeSeg.trailingMultiWildcard.matchSize [] = some 0 := rfl
        simp only [hDropEq, hd, matchAtomSegs, hStrNe]
        rw [hSum]
        rfl
      | cons c rem =>
        have hc : c = '/' := dropWhile_ne_slash_head hd
        subst hc
        have hSum := hDecomp
        rw [hd] at hSum
        have hMatch : matchAtomSegs [.singleWildcard, .multiWildcard]
            (String.ofList (cs.takeWhile (· != '/')) :: splitSegsOf rem) = true := by
          simp [matchAtomSegs, hStrNe]
        simp only [hDropEq, hd, hMatch, ↓reduceIte, pure, Option.some.injEq, Prod.mk.injEq]
        have hLen := congrArg List.length hSum
        simp only [List.length_append, List.length_cons] at hLen
        simp [Nat.add_comm, hLen]
  | lit expected =>
    have hLit := lit_trailingMulti_iff expected hpValid cs
    cases hSuf : splitSegsOf cs with
    | nil => exact absurd hSuf (splitSegsOf_ne_nil cs)
    | cons s0 segRest =>
      have hs0Eq : s0 = String.ofList (cs.takeWhile (· != '/')) := by
        rw [splitSegsOf_eq_take_drop cs] at hSuf
        cases hd : cs.dropWhile (· != '/') with
        | nil => simp [hd] at hSuf; exact hSuf.1.symm
        | cons _ _ => simp [hd] at hSuf; exact hSuf.1.symm
      rw [← hs0Eq] at hLit
      rw [show [atomToRuntime (.lit expected), .trailingMultiWildcard] =
          [.literal expected] ++ [.trailingMultiWildcard] from rfl, matchSegList_append]
      cases hPrefB : List.isPrefixOf expected.toList cs with
      | false =>
        have hLitNone : matchSegList [.literal expected] cs = none := by
          simp [matchSegList, RuntimeSeg.matchSize, hPrefB]
        rw [hLitNone]
        dsimp only
        by_cases hM : matchAtomSegs [.lit expected, .multiWildcard] (s0 :: segRest) = true
        · exfalso
          simp only [matchAtomSegs] at hM
          rcases hLit.mpr hM with ⟨rem2, hEq, _⟩
          have hPrefTrue : List.isPrefixOf expected.toList cs = true :=
            List.isPrefixOf_iff_prefix.mpr ⟨rem2, hEq⟩
          rw [hPrefTrue] at hPrefB
          cases hPrefB
        · simp [hM]
      | true =>
        have hPref : expected.toList <+: cs := List.isPrefixOf_iff_prefix.mp hPrefB
        rcases hPref with ⟨tail, hTailEq⟩
        have hLitSome : matchSegList [.literal expected] cs =
            some (expected.toList.length, tail) := by
          have hDrop : cs.drop expected.toList.length = tail := by
            rw [← hTailEq, List.drop_left]
          simp [matchSegList, RuntimeSeg.matchSize, hPrefB, hDrop]
        rw [hLitSome]
        dsimp only
        cases hT : matchSegList [.trailingMultiWildcard] tail with
        | none =>
          by_cases hM : matchAtomSegs [.lit expected, .multiWildcard] (s0 :: segRest) = true
          · exfalso
            simp only [matchAtomSegs] at hM
            rcases hLit.mpr hM with ⟨rem2, hEq, hCond⟩
            have hSame : tail = rem2 := by
              have hEq2 : expected.toList ++ tail = expected.toList ++ rem2 := by
                rw [hTailEq, hEq]
              exact List.append_cancel_left hEq2
            subst hSame
            have hSome := (matchSegList_trailingMulti_iff tail [] tail.length).mpr
              ⟨rfl, rfl, hCond⟩
            rw [hSome] at hT
            cases hT
          · simp [hM]
        | some pT =>
          rcases pT with ⟨m, rem'⟩
          rcases (matchSegList_trailingMulti_iff tail rem' m).mp hT with ⟨rfl, rfl, hCond⟩
          have hM : matchAtomSegs [.lit expected, .multiWildcard] (s0 :: segRest) = true := by
            simp only [matchAtomSegs]
            exact hLit.mp ⟨tail, hTailEq, hCond⟩
          simp only [hM, ↓reduceIte, Option.some.injEq, Prod.mk.injEq]
          have hLen := congrArg List.length hTailEq
          simp only [List.length_append] at hLen
          exact ⟨hLen, trivial⟩

private theorem matchAtomSegs_append_multi_join_ne_empty :
    ∀ (pats : List AtomSeg) (segs : List String),
      pats ≠ [] → pats.all isValidAtom = true → AtomSeg.multiWildcard ∉ pats →
      matchAtomSegs (pats ++ [.multiWildcard]) segs = true →
      joinStringsWith "/" segs ≠ ""
  | [], _, hNe, _, _, _ => absurd rfl hNe
  | [p], segs, _, hValid, hNo, hMatch => by
    have hpValid : isValidAtom p = true := List.all_eq_true.mp hValid p (by simp)
    have hpNo : p ≠ .multiWildcard := fun h => hNo (by simp [h])
    cases segs with
    | nil => cases p <;> simp [matchAtomSegs] at hMatch
    | cons s ss =>
      have hsNe : s ≠ "" := by
        cases p with
        | multiWildcard => exact absurd rfl hpNo
        | singleWildcard =>
          simp only [List.singleton_append, matchAtomSegs, Bool.and_true] at hMatch
          simp [bne] at hMatch
          exact hMatch
        | lit e =>
          rcases ValidLitNoSlash hpValid with ⟨heNe, _⟩
          simp only [List.singleton_append, matchAtomSegs, Bool.or_eq_true, beq_iff_eq] at hMatch
          rcases hMatch with rfl | hCol
          · exact heNe
          · unfold hasColonSuffix at hCol
            simp only [Bool.and_eq_true, decide_eq_true_eq] at hCol
            intro hEmpty
            subst hEmpty
            simp at hCol
      exact joinStringsWith_slash_nonempty (s :: ss) (by simp) ⟨s, by simp, hsNe⟩
  | p :: q :: qs, segs, _, _, hNo, hMatch => by
    have hpNo : p ≠ .multiWildcard := fun h => hNo (by simp [h])
    have hqNo : q ≠ .multiWildcard := fun h => hNo (by simp [h])
    cases segs with
    | nil => cases p <;> simp [matchAtomSegs] at hMatch
    | cons s ss =>
      rw [show (p :: q :: qs) ++ [.multiWildcard] = p :: q :: (qs ++ [.multiWildcard]) from rfl] at hMatch
      have hStep := matchAtomSegs_cons_two_nonmulti p q (qs ++ [AtomSeg.multiWildcard]) s ss hpNo hqNo
      rw [hStep] at hMatch
      simp only [Bool.and_eq_true] at hMatch
      cases ss with
      | nil => cases q <;> simp [matchAtomSegs] at hMatch
      | cons s2 ss2 =>
        intro hEmpty
        have hList : (s ++ "/" ++ joinStringsWith "/" (s2 :: ss2)).toList = [] := by
          rw [← joinStringsWith_cons_ne_nil s (s2 :: ss2) (by simp), hEmpty]
          rfl
        simp [String.toList_append] at hList

private theorem matchAtomSegs_nonmulti_nil {p : AtomSeg} (hp : p ≠ .multiWildcard)
    (ps : List AtomSeg) : matchAtomSegs (p :: ps) [] = false := by
  cases p with
  | multiWildcard => exact absurd rfl hp
  | singleWildcard => simp [matchAtomSegs]
  | lit _ => cases ps <;> simp [matchAtomSegs]

private theorem emitUnmerged_append_multi_pos (pats : List AtomSeg) (hNe : pats ≠ [])
    (hValid : pats.all isValidAtom = true) (hNoMulti : AtomSeg.multiWildcard ∉ pats)
    (cs : List Char) {n : Nat} {rem : List Char}
    (hMatch : matchSegList (emitUnmerged (pats ++ [.multiWildcard]) false false) cs = some (n, rem)) :
    1 ≤ n ∧ cs ≠ [] := by
  rw [emitUnmerged_append_multi pats hNe hNoMulti, matchSegList_append] at hMatch
  cases hA : matchSegList (emitUnmerged pats false false) cs with
  | none => simp [hA] at hMatch
  | some pA =>
    rcases pA with ⟨nx, remx⟩
    rcases emitUnmerged_cons_pos pats hNe hValid hNoMulti cs hA with ⟨hNxPos, hCsNe⟩
    simp only [hA] at hMatch
    cases hB : matchSegList [.trailingMultiWildcard] remx with
    | none => simp [hB] at hMatch
    | some pB =>
      simp only [hB, Option.some.injEq, Prod.mk.injEq] at hMatch
      exact ⟨by omega, hCsNe⟩

private theorem matchUnmergedAppendMulti_spec (pats : List AtomSeg) (hNe : pats ≠ [])
    (hValid : pats.all isValidAtom = true) (hNoMulti : AtomSeg.multiWildcard ∉ pats)
    (cs : List Char) :
    (match matchSegList (emitUnmerged (pats ++ [.multiWildcard]) false false) cs with
      | some (n, []) => if n == 0 then none else some (String.ofList (cs.take n))
      | _ => none) =
    (let segs := splitSegsOf cs
     if matchAtomSegs (pats ++ [.multiWildcard]) segs then
       let cap := joinStringsWith "/" segs
       if cap == "" then none else some cap
     else none) := by
  induction pats generalizing cs with
  | nil => exact absurd rfl hNe
  | cons p ps ih =>
    have hpValid : isValidAtom p = true := List.all_eq_true.mp hValid p (by simp)
    have hpNo : p ≠ .multiWildcard := fun h => hNoMulti (by simp [h])
    cases ps with
    | nil =>
      rw [emitUnmerged_append_multi [p] (by simp) hNoMulti, emitUnmerged_singleton p hpNo]
      simp only [List.singleton_append]
      rw [match_one_trailing_eq p hpValid hpNo cs]
      by_cases hM : matchAtomSegs [p, .multiWildcard] (splitSegsOf cs) = true
      · have hCapNe := matchAtomSegs_append_multi_join_ne_empty [p] (splitSegsOf cs)
          (by simp) hValid hNoMulti (by simpa using hM)
        have hCsNe : cs ≠ [] := by
          intro hNil
          subst hNil
          rw [string_ofList_join_splitSegs] at hCapNe
          exact hCapNe rfl
        have hNz : (cs.length == 0) = false := by
          cases cs with
          | nil => exact absurd rfl hCsNe
          | cons _ _ => simp
        simp [hM, hNz, hCsNe, List.take_length, string_ofList_join_splitSegs]
      · simp [hM]
    | cons q qs =>
      have hPsValid : (q :: qs).all isValidAtom = true := by
        rw [List.all_eq_true] at hValid ⊢
        intro a ha; exact hValid a (List.mem_cons_of_mem p ha)
      have hPsNo : AtomSeg.multiWildcard ∉ (q :: qs) :=
        fun h => hNoMulti (List.mem_cons_of_mem p h)
      have hqNo : q ≠ .multiWildcard := fun h => hPsNo (by simp [h])
      rw [emitUnmerged_append_multi (p :: q :: qs) (by simp) hNoMulti,
        emitUnmerged_cons_slash p q qs hpNo hqNo, List.append_assoc, matchSegList_append,
        match_one_slash_eq p hpValid hpNo cs, splitSegsOf_eq_take_drop cs]
      rw [← emitUnmerged_append_multi (q :: qs) (by simp) hPsNo]
      cases hd : cs.dropWhile (· != '/') with
      | nil =>
        dsimp only
        rw [show (p :: q :: qs) ++ [.multiWildcard] = p :: q :: (qs ++ [AtomSeg.multiWildcard]) from rfl]
        have hStep := matchAtomSegs_cons_two_nonmulti p q (qs ++ [AtomSeg.multiWildcard])
          (String.ofList (cs.takeWhile (· != '/'))) [] hpNo hqNo
        rw [hStep]
        rw [matchAtomSegs_nonmulti_nil hqNo _]
        simp
      | cons c rem =>
        have hc : c = '/' := dropWhile_ne_slash_head hd
        subst hc
        dsimp only
        have hSum := List.takeWhile_append_dropWhile (p := (· != '/')) (l := cs)
        rw [hd] at hSum
        have hMatchCons : matchAtomSegs ((p :: q :: qs) ++ [AtomSeg.multiWildcard])
            (String.ofList (cs.takeWhile (· != '/')) :: splitSegsOf rem) =
            (matchAtomSegs [p] [String.ofList (cs.takeWhile (· != '/'))] &&
              matchAtomSegs ((q :: qs) ++ [AtomSeg.multiWildcard]) (splitSegsOf rem)) :=
          matchAtomSegs_cons_two_nonmulti p q (qs ++ [AtomSeg.multiWildcard])
            (String.ofList (cs.takeWhile (· != '/'))) (splitSegsOf rem) hpNo hqNo
        rw [hMatchCons]
        by_cases hmp : matchAtomSegs [p] [String.ofList (cs.takeWhile (· != '/'))] = true
        · simp only [hmp, ↓reduceIte, Bool.true_and]
          have ihRem := ih (by simp) hPsValid hPsNo rem
          cases hSub : matchSegList (emitUnmerged ((q :: qs) ++ [AtomSeg.multiWildcard]) false false) rem with
          | none =>
            simp only [hSub] at ihRem ⊢
            have hMFalse : matchAtomSegs ((q :: qs) ++ [AtomSeg.multiWildcard]) (splitSegsOf rem) = false := by
              cases hC : matchAtomSegs ((q :: qs) ++ [AtomSeg.multiWildcard]) (splitSegsOf rem)
              · rfl
              · exfalso
                simp only [hC, ↓reduceIte] at ihRem
                have hNeRem := matchAtomSegs_append_multi_join_ne_empty (q :: qs) (splitSegsOf rem)
                  (by simp) hPsValid hPsNo hC
                simp [beq_false_of_ne hNeRem] at ihRem
            simp only [hMFalse, Bool.false_eq_true, ↓reduceIte]
          | some pSub =>
            rcases pSub with ⟨ny, remy⟩
            simp only [hSub] at ihRem ⊢
            cases remy with
            | cons r rs =>
              have hMFalse : matchAtomSegs ((q :: qs) ++ [AtomSeg.multiWildcard]) (splitSegsOf rem) = false := by
                cases hC : matchAtomSegs ((q :: qs) ++ [AtomSeg.multiWildcard]) (splitSegsOf rem)
                · rfl
                · exfalso
                  simp only [hC, ↓reduceIte] at ihRem
                  have hNeRem := matchAtomSegs_append_multi_join_ne_empty (q :: qs) (splitSegsOf rem)
                    (by simp) hPsValid hPsNo hC
                  simp [beq_false_of_ne hNeRem] at ihRem
              simp only [hMFalse, Bool.false_eq_true, ↓reduceIte]
            | nil =>
              have hNz1 : ((cs.takeWhile (· != '/')).length + 1 + ny == 0) = false := by simp
              have hNyPos : 1 ≤ ny :=
                (emitUnmerged_append_multi_pos (q :: qs) (by simp) hPsValid hPsNo rem hSub).1
              have hNzNy : (ny == 0) = false := by
                cases ny with
                | zero => omega
                | succ _ => simp
              simp only [hNz1, hNzNy, Bool.false_eq_true, ↓reduceIte] at ihRem ⊢
              rcases matchSegList_drop_eq hSub with ⟨hDrop, hLe⟩
              have hEqNy : ny = rem.length :=
                Nat.le_antisymm hLe (List.drop_eq_nil_iff.mp hDrop.symm)
              subst hEqNy
              by_cases hC : matchAtomSegs ((q :: qs) ++ [AtomSeg.multiWildcard]) (splitSegsOf rem) = true
              · simp only [hC, ↓reduceIte]
                have hAllLen : cs.take ((cs.takeWhile (· != '/')).length + 1 + rem.length) = cs := by
                  have hTotLen : (cs.takeWhile (· != '/')).length + 1 + rem.length = cs.length := by
                    have hL := congrArg List.length hSum
                    simp only [List.length_append, List.length_cons] at hL
                    omega
                  rw [hTotLen, List.take_length]
                have hLhsCap : cs.take ((cs.takeWhile (· != '/')).length + 1 + rem.length) =
                    cs.takeWhile (· != '/') ++ '/' :: rem := by
                  rw [hAllLen, hSum]
                rw [joinStringsWith_cons_ne_nil (String.ofList (cs.takeWhile (· != '/')))
                    (splitSegsOf rem) (splitSegsOf_ne_nil rem),
                  string_ofList_join_splitSegs]
                simp [hLhsCap, String.ext_iff, String.toList_append, String.toList_ofList]
              · have hCF : matchAtomSegs ((q :: qs) ++ [AtomSeg.multiWildcard]) (splitSegsOf rem) = false :=
                  Bool.eq_false_iff.mpr hC
                simp only [hCF, Bool.false_eq_true, ↓reduceIte] at ihRem
                cases ihRem
        · have hmp' : matchAtomSegs [p] [String.ofList (cs.takeWhile (· != '/'))] = false :=
            Bool.eq_false_iff.mpr hmp
          simp [hmp']

/-! ### Theorem Group C2: Multi-Segment Match + Trailing `**` Suffix -/

private theorem splitSegsOf_head (cs : List Char) :
    ∃ rest, splitSegsOf cs = String.ofList (cs.takeWhile (· != '/')) :: rest := by
  rw [splitSegsOf_eq_take_drop cs]
  cases cs.dropWhile (· != '/') with
  | nil => exact ⟨[], rfl⟩
  | cons _ rem => exact ⟨splitSegsOf rem, rfl⟩

private theorem matchPrefixAndCapture_ignore_hdr_pfx (t : CaptureTemplate) (segs : List String) :
    matchPrefixAndCapture t segs =
      matchPrefixAndCapture
        { headerKey := "", prefixSegs := [], matchSegs := t.matchSegs, suffixSegs := t.suffixSegs }
        segs := rfl

private theorem matchSegs_beq_single_multi {xs : List AtomSeg}
    (h : (xs == [AtomSeg.multiWildcard]) = true) : xs = [AtomSeg.multiWildcard] := by
  cases xs with
  | nil => cases h
  | cons x rest =>
    cases rest with
    | cons _ _ => cases x <;> cases h
    | nil =>
      cases x with
      | multiWildcard => rfl
      | singleWildcard => cases h
      | lit _ => cases h

private theorem matchPrefixAndCapture_multi_cons
    (p q : AtomSeg) (qs : List AtomSeg)
    (hValid : (q :: qs).all isValidAtom = true)
    (hpNo : p ≠ .multiWildcard) (hNoMulti : AtomSeg.multiWildcard ∉ (p :: q :: qs))
    (s0 : String) (remSegs : List String) :
    matchPrefixAndCapture
      { headerKey := "", prefixSegs := [], matchSegs := p :: q :: qs, suffixSegs := [.multiWildcard] }
      (s0 :: remSegs) =
    if matchAtomSegs [p] [s0] then
      (matchPrefixAndCapture
        { headerKey := "", prefixSegs := [], matchSegs := q :: qs, suffixSegs := [.multiWildcard] }
        remSegs).map (fun cap => s0 ++ "/" ++ cap)
    else none := by
  have hPsNo : AtomSeg.multiWildcard ∉ (q :: qs) :=
    fun h => hNoMulti (List.mem_cons_of_mem p h)
  have hSingleMulti : ∀ (xs : List String), matchAtomSegs [.multiWildcard] xs = true := by
    intro xs
    cases xs <;> simp [matchAtomSegs]
  have hMultiBeq : ([AtomSeg.multiWildcard] == [AtomSeg.multiWildcard]) = true := rfl
  by_cases hLen : remSegs.length < (q :: qs).length
  · have hLen' : remSegs.length < qs.length + 1 := hLen
    unfold matchPrefixAndCapture
    by_cases hmp : matchAtomSegs [p] [s0] = true <;> simp [hLen', hmp]
  · have hLen' : ¬ remSegs.length < qs.length + 1 := hLen
    unfold matchPrefixAndCapture
    dsimp only
    simp only [List.isEmpty_cons, Bool.false_eq_true, ↓reduceIte, List.length_cons,
      Nat.add_lt_add_iff_right, hLen', List.take_succ_cons, List.drop_succ_cons]
    have hGe : qs.length + 1 ≤ remSegs.length := Nat.le_of_not_lt hLen'
    have hMatchedLen : (remSegs.take (qs.length + 1)).length = (q :: qs).length := by
      simp only [List.length_take, List.length_cons]
      omega
    have hMatchedNe : remSegs.take (qs.length + 1) ≠ [] := by
      intro hNil
      have hL : (remSegs.take (qs.length + 1)).length = 0 := by simp [hNil]
      simp [hMatchedLen] at hL
    have hMatchSplit : matchAtomSegs (p :: q :: qs) (s0 :: remSegs.take (qs.length + 1)) =
        (matchAtomSegs [p] [s0] && matchAtomSegs (q :: qs) (remSegs.take (qs.length + 1))) :=
      matchAtomSegs_cons_no_multi p (q :: qs) s0 (remSegs.take (qs.length + 1)) hNoMulti hpNo (by simp)
    rw [hMatchSplit]
    simp only [hSingleMulti, Bool.and_true, hMultiBeq, ↓reduceIte]
    cases hrP : (q :: qs).reverse with
    | nil =>
      have : ((q :: qs).reverse).length = 0 := by simp [hrP]
      simp at this
    | cons rPat revSubP =>
      have hPatEq : q :: qs = revSubP.reverse ++ [rPat] := by
        calc q :: qs = (q :: qs).reverse.reverse := by simp
          _ = (rPat :: revSubP).reverse := by rw [hrP]
          _ = revSubP.reverse ++ [rPat] := by simp
      have hBigPatRev : (p :: q :: qs).reverse = rPat :: (revSubP ++ [p]) := by
        simp [List.reverse_cons, hrP]
      cases hrS : (remSegs.take (qs.length + 1)).reverse with
      | nil =>
        have : remSegs.take (qs.length + 1) = [] := by
          simpa using congrArg List.reverse hrS
        exact absurd this hMatchedNe
      | cons lastSeg revSubS =>
        have hSegEq : remSegs.take (qs.length + 1) = revSubS.reverse ++ [lastSeg] := by
          calc remSegs.take (qs.length + 1) = (remSegs.take (qs.length + 1)).reverse.reverse := by simp
            _ = (lastSeg :: revSubS).reverse := by rw [hrS]
            _ = revSubS.reverse ++ [lastSeg] := by simp
        have hBigSegRev : (s0 :: remSegs.take (qs.length + 1)).reverse = lastSeg :: (revSubS ++ [s0]) := by
          simp [List.reverse_cons, hrS]
        have hSubLenEq : revSubP.reverse.length = revSubS.reverse.length := by
          have h1 := congrArg List.length hPatEq
          have h2 := congrArg List.length hSegEq
          simp only [List.length_append, List.length_singleton] at h1 h2
          omega
        have hNoSubP : AtomSeg.multiWildcard ∉ revSubP.reverse := by
          intro hMem
          apply hPsNo
          rw [hPatEq]
          exact List.mem_append_left _ hMem
        have hValidRPat : isValidAtom rPat = true := by
          rw [List.all_eq_true] at hValid
          apply hValid
          rw [hPatEq]
          simp
        by_cases hmp : matchAtomSegs [p] [s0] = true
        · simp only [hmp, ↓reduceIte, Bool.true_and]
          by_cases hB1 : matchAtomSegs (q :: qs) (remSegs.take (qs.length + 1)) = true
          · simp only [hB1, ↓reduceIte]
            have hCapRemNe := matchAtomSegs_valid_join_ne_empty (q :: qs) (remSegs.take (qs.length + 1))
              (by simp) hValid hPsNo hMatchedLen hB1
            have hBigCapNe : (joinStringsWith "/" (s0 :: remSegs.take (qs.length + 1)) == "") = false := by
              apply beq_false_of_ne
              rw [joinStringsWith_cons_ne_nil s0 (remSegs.take (qs.length + 1)) hMatchedNe]
              intro hE
              have hL : (s0 ++ "/" ++ joinStringsWith "/" (remSegs.take (qs.length + 1))).toList = [] := by
                rw [hE]; rfl
              simp [String.toList_append] at hL
            have hCapRemF : (joinStringsWith "/" (remSegs.take (qs.length + 1)) == "") = false :=
              beq_false_of_ne hCapRemNe
            simp only [hCapRemF, hBigCapNe, Bool.false_eq_true, ↓reduceIte, Option.map_some]
            congr 1
            exact joinStringsWith_cons_ne_nil s0 (remSegs.take (qs.length + 1)) hMatchedNe
          · have hB1F : matchAtomSegs (q :: qs) (remSegs.take (qs.length + 1)) = false :=
              Bool.eq_false_iff.mpr hB1
            simp only [hB1F, Bool.false_eq_true, ↓reduceIte, hBigPatRev, hBigSegRev]
            cases rPat with
            | multiWildcard => simp
            | singleWildcard => simp
            | lit expected =>
              have hexpNe : expected ≠ "" := (ValidLitNoSlash hValidRPat).1
              have hTailCons : revSubP.reverse = [] → revSubS.reverse = [] := by
                intro hNil
                cases hRS : revSubS.reverse with
                | nil => rfl
                | cons _ _ => simp [hNil, hRS] at hSubLenEq
              have hMatchPrefix : matchAtomSegs (revSubP ++ [p]).reverse (revSubS ++ [s0]).reverse =
                  matchAtomSegs revSubP.reverse revSubS.reverse := by
                simp only [List.reverse_append, List.reverse_singleton, List.singleton_append]
                cases hRP : revSubP.reverse with
                | nil =>
                  have hRS := hTailCons hRP
                  rw [hRS]
                  simp [hmp, matchAtomSegs]
                | cons a as =>
                  have hNoMid' : AtomSeg.multiWildcard ∉ (p :: a :: as) := by
                    intro hM
                    rcases List.mem_cons.mp hM with rfl | hIn
                    · exact hpNo rfl
                    · exact hNoSubP (by rw [hRP]; exact hIn)
                  rw [matchAtomSegs_cons_no_multi p (a :: as) s0 revSubS.reverse
                    hNoMid' hpNo (by simp)]
                  simp [hmp]
              simp only [List.reverse_append, List.reverse_singleton, List.singleton_append] at hMatchPrefix ⊢
              have hMatchP' : matchAtomSegs (p :: revSubP.reverse) (s0 :: revSubS.reverse) =
                  matchAtomSegs revSubP.reverse revSubS.reverse := hMatchPrefix
              rw [hMatchP']
              by_cases hCol : (hasColonSuffix expected lastSeg &&
                  matchAtomSegs revSubP.reverse revSubS.reverse) = true
              · simp only [hCol, ↓reduceIte, List.cons_append]
                have hSmallCapNe : (joinStringsWith "/" (revSubS.reverse ++ [expected]) == "") = false :=
                  beq_false_of_ne (joinStringsWith_slash_nonempty (revSubS.reverse ++ [expected])
                    (by simp) ⟨expected, by simp, hexpNe⟩)
                have hBigCapNe : (joinStringsWith "/" (s0 :: (revSubS.reverse ++ [expected])) == "") = false :=
                  beq_false_of_ne (joinStringsWith_slash_nonempty (s0 :: (revSubS.reverse ++ [expected]))
                    (by simp) ⟨expected, by simp, hexpNe⟩)
                simp only [hSmallCapNe, hBigCapNe, Bool.false_eq_true, ↓reduceIte, Option.map_some]
                congr 1
                exact joinStringsWith_cons_ne_nil s0 (revSubS.reverse ++ [expected]) (by simp)
              · have hColF : (hasColonSuffix expected lastSeg &&
                    matchAtomSegs revSubP.reverse revSubS.reverse) = false :=
                  Bool.eq_false_iff.mpr hCol
                simp only [hColF, Bool.false_eq_true, ↓reduceIte, Option.map_none]
        · have hmp' : matchAtomSegs [p] [s0] = false := Bool.eq_false_iff.mpr hmp
          simp only [hmp', Bool.false_eq_true, ↓reduceIte, Bool.false_and, hBigPatRev, hBigSegRev]
          cases rPat with
          | multiWildcard => simp
          | singleWildcard => simp
          | lit expected =>
            have hMatchFalse : matchAtomSegs (revSubP ++ [p]).reverse (revSubS ++ [s0]).reverse = false := by
              simp only [List.reverse_append, List.reverse_singleton, List.singleton_append]
              cases hRP : revSubP.reverse with
              | nil =>
                have hRS : revSubS.reverse = [] := by
                  cases hS : revSubS.reverse with
                  | nil => rfl
                  | cons _ _ => simp [hRP, hS] at hSubLenEq
                rw [hRS]
                exact hmp'
              | cons a as =>
                have hNoMid' : AtomSeg.multiWildcard ∉ (p :: a :: as) := by
                  intro hM
                  rcases List.mem_cons.mp hM with rfl | hIn
                  · exact hpNo rfl
                  · exact hNoSubP (by rw [hRP]; exact hIn)
                rw [matchAtomSegs_cons_no_multi p (a :: as) s0 revSubS.reverse
                  hNoMid' hpNo (by simp)]
                simp [hmp']
            simp only [List.reverse_append, List.reverse_singleton, List.singleton_append] at hMatchFalse ⊢
            simp only [hMatchFalse, Bool.and_false, Bool.false_eq_true, ↓reduceIte]

private theorem matchUnmergedWithMultiSuffix_spec (pats : List AtomSeg) (hNe : pats ≠ [])
    (hValid : pats.all isValidAtom = true) (hNoMulti : AtomSeg.multiWildcard ∉ pats)
    (cs : List Char) :
    (do
      let (mCount, rem2) ← matchSegList (emitUnmerged pats false false) cs
      let (_, rem3) ← matchSegList [.trailingMultiWildcard] rem2
      if !rem3.isEmpty || mCount == 0 then none
      else some (String.ofList (cs.take mCount))) =
    matchPrefixAndCapture
      { headerKey := "", prefixSegs := [], matchSegs := pats, suffixSegs := [.multiWildcard] }
      (splitSegsOf cs) := by
  induction pats generalizing cs with
  | nil => exact absurd rfl hNe
  | cons p ps ih =>
    have hpValid : isValidAtom p = true := List.all_eq_true.mp hValid p (by simp)
    have hpNo : p ≠ .multiWildcard := fun h => hNoMulti (by simp [h])
    cases ps with
    | nil =>
      rw [emitUnmerged_singleton p hpNo]
      rcases splitSegsOf_head cs with ⟨segRest, hSegsEq⟩
      rw [hSegsEq]
      let s0 := String.ofList (cs.takeWhile (· != '/'))
      have hDecomp := List.takeWhile_append_dropWhile (p := (· != '/')) (l := cs)
      cases p with
      | multiWildcard => exact absurd rfl hpNo
      | singleWildcard =>
        cases hZero : ((cs.takeWhile (· != '/')).length == 0) with
        | true =>
          have hTakeNil : cs.takeWhile (· != '/') = [] := by
            cases ht : cs.takeWhile (· != '/') with
            | nil => rfl
            | cons _ _ => simp [ht] at hZero
          have hNone : matchSegList [atomToRuntime .singleWildcard] cs = none := by
            simp [atomToRuntime, matchSegList, RuntimeSeg.matchSize, hZero]
          rw [hNone]
          simp [hTakeNil, bind, Option.bind, matchPrefixAndCapture, matchAtomSegs]
        | false =>
          have hStrNe : s0 ≠ "" := by
            simp only [s0, ne_eq, String.ext_iff, String.toList_ofList]
            cases ht : cs.takeWhile (· != '/') with
            | nil => simp [ht] at hZero
            | cons _ _ => simp
          have hDropEq : cs.drop (cs.takeWhile (· != '/')).length = cs.dropWhile (· != '/') := by
            have hTmp := @List.drop_left Char (cs.takeWhile (· != '/')) (cs.dropWhile (· != '/'))
            rwa [hDecomp] at hTmp
          have hCondDrop : cs.dropWhile (· != '/') = [] ∨
              (∃ r, cs.dropWhile (· != '/') = '/' :: r) ∨
              (∃ r, cs.dropWhile (· != '/') = ':' :: r) := by
            cases hd : cs.dropWhile (· != '/') with
            | nil => exact Or.inl rfl
            | cons c r => exact Or.inr (Or.inl ⟨r, by rw [dropWhile_ne_slash_head hd]⟩)
          have hM1 : matchSegList [atomToRuntime .singleWildcard] cs =
              some ((cs.takeWhile (· != '/')).length, cs.drop (cs.takeWhile (· != '/')).length) := by
            simp [atomToRuntime, matchSegList, RuntimeSeg.matchSize, hZero]
          have hM2 : matchSegList [.trailingMultiWildcard] (cs.drop (cs.takeWhile (· != '/')).length) =
              some ((cs.drop (cs.takeWhile (· != '/')).length).length, []) :=
            (matchSegList_trailingMulti_iff _ [] _).mpr ⟨rfl, rfl, by rwa [hDropEq]⟩
          have hSingleMulti : matchAtomSegs [.multiWildcard] segRest = true := by
            cases segRest <;> simp [matchAtomSegs]
          have hTakeEq : cs.take (cs.takeWhile (· != '/')).length = cs.takeWhile (· != '/') := by
            have hTmp := @List.take_left Char (cs.takeWhile (· != '/')) (cs.dropWhile (· != '/'))
            rwa [hDecomp] at hTmp
          rw [hM1]
          simp only [bind, Option.bind, hM2, List.isEmpty_nil, Bool.not_true, hZero,
            Bool.false_or, Bool.false_eq_true, ↓reduceIte, hTakeEq]
          simp [matchPrefixAndCapture, matchAtomSegs, bne_iff_ne, hStrNe,
            joinStringsWith, beq_false_of_ne hStrNe, s0]
      | lit expected =>
        rcases ValidLitNoSlash hpValid with ⟨hExpNe, _⟩
        have hLit := lit_trailingMulti_iff expected hpValid cs
        have hSingleMulti : matchAtomSegs [.multiWildcard] segRest = true := by
          cases segRest <;> simp [matchAtomSegs]
        have hStrNe : (expected == "") = false := beq_false_of_ne hExpNe
        cases hPrefB : List.isPrefixOf expected.toList cs with
        | false =>
          have hM1None : matchSegList [atomToRuntime (.lit expected)] cs = none := by
            simp [atomToRuntime, matchSegList, RuntimeSeg.matchSize, hPrefB]
          rw [hM1None]
          simp only [bind, Option.bind]
          unfold matchPrefixAndCapture
          simp only [List.isEmpty_cons, Bool.false_eq_true, ↓reduceIte, List.length_cons,
            List.length_nil, Nat.zero_add, Nat.lt_one_iff, reduceCtorEq, List.take_succ_cons,
            List.take_zero, matchAtomSegs, Bool.and_true, List.reverse_cons,
            List.reverse_nil, List.nil_append]
          by_cases hEq : String.ofList (cs.takeWhile (· != '/')) = expected
          · exfalso
            have hOr : (String.ofList (cs.takeWhile (· != '/')) == expected ||
                hasColonSuffix expected (String.ofList (cs.takeWhile (· != '/')))) = true := by
              simp [hEq]
            rcases hLit.mpr hOr with ⟨rem2, hApp, _⟩
            have hP := List.isPrefixOf_iff_prefix.mpr ⟨rem2, hApp⟩
            rw [hP] at hPrefB
            cases hPrefB
          · have hBeqF : (String.ofList (cs.takeWhile (· != '/')) == expected) = false :=
              beq_false_of_ne hEq
            simp only [hBeqF, Bool.false_eq_true, ↓reduceIte]
            by_cases hCol : hasColonSuffix expected (String.ofList (cs.takeWhile (· != '/'))) = true
            · exfalso
              have hOr : (String.ofList (cs.takeWhile (· != '/')) == expected ||
                  hasColonSuffix expected (String.ofList (cs.takeWhile (· != '/')))) = true := by
                simp [hCol]
              rcases hLit.mpr hOr with ⟨rem2, hApp, _⟩
              have hP := List.isPrefixOf_iff_prefix.mpr ⟨rem2, hApp⟩
              rw [hP] at hPrefB
              cases hPrefB
            · have hColF : hasColonSuffix expected (String.ofList (cs.takeWhile (· != '/'))) = false :=
                Bool.eq_false_iff.mpr hCol
              simp [hColF]
        | true =>
          have hPref : expected.toList <+: cs := List.isPrefixOf_iff_prefix.mp hPrefB
          rcases hPref with ⟨tail, hTailEq⟩
          have hDrop : cs.drop expected.toList.length = tail := by
            rw [← hTailEq, List.drop_left]
          have hTake : cs.take expected.toList.length = expected.toList := by
            rw [← hTailEq, List.take_left]
          have hLenNz : (expected.toList.length == 0) = false := by
            cases hEL : expected.toList with
            | nil => exact absurd (String.ext (by simp [hEL])) hExpNe
            | cons _ _ => simp
          have hM1Some : matchSegList [atomToRuntime (.lit expected)] cs =
              some (expected.toList.length, tail) := by
            simp [atomToRuntime, matchSegList, RuntimeSeg.matchSize, hPrefB, hDrop]
          rw [hM1Some]
          simp only [bind, Option.bind]
          cases hT : matchSegList [.trailingMultiWildcard] tail with
          | none =>
            unfold matchPrefixAndCapture
            simp only [List.isEmpty_cons, Bool.false_eq_true, ↓reduceIte, List.length_cons,
              List.length_nil, Nat.zero_add, Nat.lt_one_iff, reduceCtorEq, List.take_succ_cons,
              List.take_zero, matchAtomSegs, Bool.and_true, List.reverse_cons,
              List.reverse_nil, List.nil_append]
            have hNotLit : (String.ofList (cs.takeWhile (· != '/')) == expected ||
                hasColonSuffix expected (String.ofList (cs.takeWhile (· != '/')))) = false := by
              cases hOr : (String.ofList (cs.takeWhile (· != '/')) == expected ||
                  hasColonSuffix expected (String.ofList (cs.takeWhile (· != '/'))))
              · rfl
              · exfalso
                rcases hLit.mpr hOr with ⟨rem2, hApp, hCond⟩
                have hSame : tail = rem2 := by
                  have hEq2 : expected.toList ++ tail = expected.toList ++ rem2 := by
                    rw [hTailEq, hApp]
                  exact List.append_cancel_left hEq2
                subst hSame
                have hSome := (matchSegList_trailingMulti_iff tail [] tail.length).mpr
                  ⟨rfl, rfl, hCond⟩
                rw [hSome] at hT
                cases hT
            simp only [Bool.or_eq_false_iff] at hNotLit
            simp [hNotLit.1, hNotLit.2]
          | some pT =>
            rcases pT with ⟨m, rem'⟩
            rcases (matchSegList_trailingMulti_iff tail rem' m).mp hT with ⟨rfl, rfl, hCond⟩
            simp only [List.isEmpty_nil, Bool.not_true, hLenNz, Bool.false_or,
              Bool.false_eq_true, ↓reduceIte, hTake, String.ofList_toList]
            have hOr := hLit.mp ⟨tail, hTailEq, hCond⟩
            unfold matchPrefixAndCapture
            simp only [List.isEmpty_cons, Bool.false_eq_true, ↓reduceIte, List.length_cons,
              List.length_nil, Nat.zero_add, Nat.lt_one_iff, reduceCtorEq, List.take_succ_cons,
              List.take_zero, matchAtomSegs, Bool.and_true, List.reverse_cons,
              List.reverse_nil, List.nil_append, joinStringsWith, hStrNe,
              show ([AtomSeg.multiWildcard] == [AtomSeg.multiWildcard]) = true from rfl]
            simp only [Bool.or_eq_true] at hOr
            rcases hOr with hEq | hCol
            · have hEqS : String.ofList (cs.takeWhile (· != '/')) = expected := eq_of_beq hEq
              simp [hEqS, hStrNe]
            · by_cases hEq : String.ofList (cs.takeWhile (· != '/')) = expected
              · simp [hEq, hStrNe]
              · simp [beq_false_of_ne hEq, hCol]
    | cons q qs =>
      have hPsValid : (q :: qs).all isValidAtom = true := by
        rw [List.all_eq_true] at hValid ⊢
        intro a ha; exact hValid a (List.mem_cons_of_mem p ha)
      have hPsNo : AtomSeg.multiWildcard ∉ (q :: qs) :=
        fun h => hNoMulti (List.mem_cons_of_mem p h)
      have hqNo : q ≠ .multiWildcard := fun h => hPsNo (by simp [h])
      have hMatchJoin :
          (do
            let (mCount, rem2) ← matchSegList (emitUnmerged (p :: q :: qs) false false) cs
            let (_, rem3) ← matchSegList [.trailingMultiWildcard] rem2
            if !rem3.isEmpty || mCount == 0 then none
            else some (String.ofList (cs.take mCount))) =
          match matchSegList [atomToRuntime p, .literal "/"] cs with
          | none => none
          | some (nx, remx) =>
            match matchSegList (emitUnmerged (q :: qs) false false) remx with
            | none => none
            | some (mRem, rem2) =>
              match matchSegList [.trailingMultiWildcard] rem2 with
              | none => none
              | some (_, rem3) =>
                if !rem3.isEmpty || (nx + mRem == 0) then none
                else some (String.ofList (cs.take (nx + mRem))) := by
        rw [emitUnmerged_cons_slash p q qs hpNo hqNo, matchSegList_append]
        simp only [bind, Option.bind]
        cases matchSegList [atomToRuntime p, .literal "/"] cs with
        | none => rfl
        | some pX =>
          rcases pX with ⟨nx, remx⟩
          dsimp only
          cases matchSegList (emitUnmerged (q :: qs) false false) remx with
          | none => rfl
          | some pM =>
            rcases pM with ⟨mRem, rem2⟩
            dsimp only
            cases matchSegList [.trailingMultiWildcard] rem2 <;> rfl
      rw [hMatchJoin, match_one_slash_eq p hpValid hpNo cs, splitSegsOf_eq_take_drop cs]
      cases hd : cs.dropWhile (· != '/') with
      | nil =>
        dsimp only
        simp [matchPrefixAndCapture]
      | cons c rem =>
        have hc : c = '/' := dropWhile_ne_slash_head hd
        subst hc
        dsimp only
        have hSum := List.takeWhile_append_dropWhile (p := (· != '/')) (l := cs)
        rw [hd] at hSum
        rw [matchPrefixAndCapture_multi_cons p q qs hPsValid hpNo hNoMulti
          (String.ofList (cs.takeWhile (· != '/'))) (splitSegsOf rem)]
        by_cases hmp : matchAtomSegs [p] [String.ofList (cs.takeWhile (· != '/'))] = true
        · simp only [hmp, ↓reduceIte]
          have ihRem := ih (by simp) hPsValid hPsNo rem
          simp only [bind, Option.bind] at ihRem
          rw [← ihRem]
          cases hMRem : matchSegList (emitUnmerged (q :: qs) false false) rem with
          | none => rfl
          | some pM =>
            rcases pM with ⟨mRem, rem2⟩
            dsimp only
            cases hTm : matchSegList [.trailingMultiWildcard] rem2 with
            | none => rfl
            | some pT =>
              rcases pT with ⟨_, rem3⟩
              dsimp only
              have hmRemPos : 1 ≤ mRem :=
                (emitUnmerged_cons_pos (q :: qs) (by simp) hPsValid hPsNo rem hMRem).1
              have hNz1 : ((cs.takeWhile (· != '/')).length + 1 + mRem == 0) = false := by simp
              have hNz2 : (mRem == 0) = false := by
                cases mRem with
                | zero => omega
                | succ _ => simp
              cases rem3 with
              | cons _ _ =>
                simp only [List.isEmpty_cons, Bool.not_false, Bool.true_or, ↓reduceIte,
                  Option.map_none]
              | nil =>
                simp only [List.isEmpty_nil, Bool.not_true, hNz1, hNz2, Bool.false_or,
                  Bool.false_eq_true, ↓reduceIte, Option.map_some]
                congr 1
                have hDropCs : cs.drop (cs.takeWhile (· != '/')).length = '/' :: rem := by
                  have hTmp := @List.drop_left Char (cs.takeWhile (· != '/')) ('/' :: rem)
                  rwa [hSum] at hTmp
                have hTakeSum : cs.take ((cs.takeWhile (· != '/')).length + 1 + mRem) =
                    cs.take (cs.takeWhile (· != '/')).length ++ '/' :: rem.take mRem := by
                  rw [Nat.add_assoc, List.take_add, hDropCs, Nat.add_comm 1 mRem, List.take_succ_cons]
                have hTakeLeft : cs.take (cs.takeWhile (· != '/')).length = cs.takeWhile (· != '/') := by
                  have hTmp := @List.take_left Char (cs.takeWhile (· != '/')) ('/' :: rem)
                  rwa [hSum] at hTmp
                rw [hTakeSum, hTakeLeft]
                simp [String.ext_iff, String.toList_append, String.toList_ofList]
        · have hmp' : matchAtomSegs [p] [String.ofList (cs.takeWhile (· != '/'))] = false :=
            Bool.eq_false_iff.mpr hmp
          simp [hmp']

/-! ### Full Template Refinement (`evalCompiled (compileTemplate t) val = extractCapture t val`) -/

private theorem split_dropLast_no_multi (xs : List AtomSeg)
    (hNoDrop : AtomSeg.multiWildcard ∉ xs.dropLast) :
    AtomSeg.multiWildcard ∉ xs ∨ xs = [.multiWildcard] ∨
      ∃ init, init ≠ [] ∧ AtomSeg.multiWildcard ∉ init ∧ xs = init ++ [.multiWildcard] := by
  cases hr : xs.reverse with
  | nil =>
    have hNil : xs = [] := by simpa using congrArg List.reverse hr
    left; simp [hNil]
  | cons r revInit =>
    have hEq : xs = revInit.reverse ++ [r] := by
      calc xs = xs.reverse.reverse := by simp
        _ = (r :: revInit).reverse := by rw [hr]
        _ = revInit.reverse ++ [r] := by simp
    have hDropEq : xs.dropLast = revInit.reverse := by
      rw [hEq, List.dropLast_concat]
    rw [hDropEq] at hNoDrop
    cases r with
    | multiWildcard =>
      cases hRev : revInit.reverse with
      | nil =>
        right; left; simp [hEq, hRev]
      | cons x xs' =>
        right; right
        refine ⟨revInit.reverse, by simp [hRev], hNoDrop, hEq⟩
    | singleWildcard =>
      left
      rw [hEq]
      intro hM
      simp only [List.mem_append, List.mem_singleton] at hM
      rcases hM with h1 | h2
      · exact hNoDrop h1
      · cases h2
    | lit _ =>
      left
      rw [hEq]
      intro hM
      simp only [List.mem_append, List.mem_singleton] at hM
      rcases hM with h1 | h2
      · exact hNoDrop h1
      · cases h2

private theorem well_formed_decomp (t : CaptureTemplate) (hWF : t.WellFormed) :
    isValidHeaderKey t.headerKey = true ∧
    t.matchSegs ≠ [] ∧
    t.prefixSegs.all isValidAtom = true ∧
    t.matchSegs.all isValidAtom = true ∧
    t.suffixSegs.all isValidAtom = true ∧
    AtomSeg.multiWildcard ∉ t.prefixSegs ∧
    AtomSeg.multiWildcard ∉ t.suffixSegs.dropLast ∧
    (t.suffixSegs = [] → AtomSeg.multiWildcard ∉ t.matchSegs.dropLast) ∧
    (t.suffixSegs ≠ [] → AtomSeg.multiWildcard ∉ t.matchSegs) := by
  rcases well_formed_forbids_mid_multi_wildcard t hWF with ⟨hPfxNo, hMatchIfSfx, hSfxDrop⟩
  unfold CaptureTemplate.WellFormed CaptureTemplate.isWellFormed CaptureTemplate.allSegs at hWF
  simp only [Bool.and_eq_true, Bool.not_eq_true'] at hWF
  rcases hWF with ⟨⟨⟨hKey, hMatchNeBool⟩, hAllValid⟩, hAllDrop⟩
  have hMatchNe : t.matchSegs ≠ [] := by
    intro h
    simp [h] at hMatchNeBool
  rw [List.all_eq_true] at hAllValid hAllDrop
  have hV1 : t.prefixSegs.all isValidAtom = true := by
    rw [List.all_eq_true]
    intro x hx; exact hAllValid x (List.mem_append_left _ (List.mem_append_left _ hx))
  have hV2 : t.matchSegs.all isValidAtom = true := by
    rw [List.all_eq_true]
    intro x hx; exact hAllValid x (List.mem_append_left _ (List.mem_append_right _ hx))
  have hV3 : t.suffixSegs.all isValidAtom = true := by
    rw [List.all_eq_true]
    intro x hx; exact hAllValid x (List.mem_append_right _ hx)
  refine ⟨hKey, hMatchNe, hV1, hV2, hV3, hPfxNo, hSfxDrop, ?_, hMatchIfSfx⟩
  intro hSfxNil hMem
  have hDropNil : (t.prefixSegs ++ t.matchSegs ++ t.suffixSegs).dropLast =
      t.prefixSegs ++ t.matchSegs.dropLast := by
    rw [hSfxNil, List.append_nil]
    exact List.dropLast_append_of_ne_nil hMatchNe
  have hIn : AtomSeg.multiWildcard ∈ (t.prefixSegs ++ t.matchSegs ++ t.suffixSegs).dropLast := by
    rw [hDropNil]
    exact List.mem_append_right _ hMem
  have hb := hAllDrop .multiWildcard hIn
  revert hb
  decide

private theorem matchPrefixAndCapture_nil (t : CaptureTemplate) (hMatchNe : t.matchSegs ≠ []) :
    matchPrefixAndCapture t [] = none := by
  unfold matchPrefixAndCapture
  cases hSfx : t.suffixSegs with
  | nil =>
    simp only [List.isEmpty_nil, ↓reduceIte, joinStringsWith]
    cases hM : t.matchSegs with
    | nil => exact absurd hM hMatchNe
    | cons m ms =>
      cases m with
      | multiWildcard => cases ms <;> simp [matchAtomSegs]
      | singleWildcard => simp [matchAtomSegs]
      | lit _ => cases ms <;> simp [matchAtomSegs]
  | cons _ _ =>
    have hPos : 0 < t.matchSegs.length := by
      cases hM : t.matchSegs with
      | nil => exact absurd hM hMatchNe
      | cons _ _ => simp
    simp [List.isEmpty_cons, List.length_nil, hPos]

private theorem matchPrefixAndCapture_eq_compiled (t : CaptureTemplate) (hWF : t.WellFormed)
    (rem1 : List Char) :
    (do
      let (mCount, rem2) ← matchSegList (emitUnmerged t.matchSegs false false) rem1
      let (_, rem3) ← matchSegList (emitUnmerged t.suffixSegs false true) rem2
      if !rem3.isEmpty || mCount == 0 then none
      else some (String.ofList (rem1.take mCount))) =
    matchPrefixAndCapture t (splitSegsOf rem1) := by
  rcases well_formed_decomp t hWF with
    ⟨_, hMatchNe, _, hMatchValid, hSfxValid, _, hSfxDrop, hMatchDropNo, hMatchIfSfx⟩
  cases hSfx : t.suffixSegs with
  | nil =>
    have hDropNo : AtomSeg.multiWildcard ∉ t.matchSegs.dropLast := hMatchDropNo hSfx
    have hEmitNil : emitUnmerged [] false true = [] := rfl
    have hBindSimplify :
        (do
          let (mCount, rem2) ← matchSegList (emitUnmerged t.matchSegs false false) rem1
          let (_, rem3) ← matchSegList (emitUnmerged [] false true) rem2
          if !rem3.isEmpty || mCount == 0 then none
          else some (String.ofList (rem1.take mCount))) =
        match matchSegList (emitUnmerged t.matchSegs false false) rem1 with
        | some (n, []) => if n == 0 then none else some (String.ofList (rem1.take n))
        | _ => none := by
      rw [hEmitNil]
      simp only [bind, Option.bind, matchSegList]
      cases matchSegList (emitUnmerged t.matchSegs false false) rem1 with
      | none => rfl
      | some pM =>
        rcases pM with ⟨n, rem2⟩
        cases rem2 <;> simp
    rw [hBindSimplify]
    rcases split_dropLast_no_multi t.matchSegs hDropNo with
      hNoMulti | hBareMulti | ⟨init, hInitNe, hInitNo, hInitEq⟩
    · rw [matchUnmergedExact_spec t.matchSegs hMatchNe hMatchValid hNoMulti rem1]
      rw [matchPrefixAndCapture_ignore_hdr_pfx t (splitSegsOf rem1)]
      unfold matchPrefixAndCapture
      simp only [hSfx, List.isEmpty_nil, ↓reduceIte]
      by_cases hm : matchAtomSegs t.matchSegs (splitSegsOf rem1) = true
      · have hLenEq : (splitSegsOf rem1).length = t.matchSegs.length :=
          matchAtomSegs_length_of_no_multi t.matchSegs (splitSegsOf rem1) hNoMulti hm
        simp [hm, hLenEq]
      · have hmf : matchAtomSegs t.matchSegs (splitSegsOf rem1) = false :=
          Bool.eq_false_iff.mpr hm
        simp [hmf]
    · have hEmitMulti : emitUnmerged t.matchSegs false false = [.multiWildcard] := by
        rw [hBareMulti]; rfl
      rw [hEmitMulti]
      simp only [matchSegList, RuntimeSeg.matchSize, bind, Option.bind, pure, List.drop_length]
      rw [matchPrefixAndCapture_ignore_hdr_pfx t (splitSegsOf rem1)]
      unfold matchPrefixAndCapture
      have hMultiTrue : matchAtomSegs [.multiWildcard] (splitSegsOf rem1) = true := by
        cases splitSegsOf rem1 <;> simp [matchAtomSegs]
      simp only [hSfx, hBareMulti, List.isEmpty_nil, ↓reduceIte, hMultiTrue,
        string_ofList_join_splitSegs]
      cases rem1 with
      | nil => simp
      | cons c cs => simp
    · have hInitValid : init.all isValidAtom = true := by
        rw [List.all_eq_true] at hMatchValid ⊢
        intro a ha; apply hMatchValid; rw [hInitEq]; exact List.mem_append_left _ ha
      rw [hInitEq, matchUnmergedAppendMulti_spec init hInitNe hInitValid hInitNo rem1]
      rw [matchPrefixAndCapture_ignore_hdr_pfx t (splitSegsOf rem1)]
      unfold matchPrefixAndCapture
      simp [hSfx, hInitEq]
  | cons sHead sTail =>
    have hMatchNo : AtomSeg.multiWildcard ∉ t.matchSegs :=
      hMatchIfSfx (by rw [hSfx]; simp)
    cases hTail : sTail with
    | nil =>
      by_cases hsHead : sHead = .multiWildcard
      · subst hsHead
        have hSfxSingle : emitUnmerged [.multiWildcard] false true = [.trailingMultiWildcard] := rfl
        have hEq1 : t.suffixSegs = [.multiWildcard] := by rw [hSfx, hTail]
        rw [hSfxSingle, matchPrefixAndCapture_ignore_hdr_pfx t]
        simp only [hEq1]
        exact matchUnmergedWithMultiSuffix_spec t.matchSegs hMatchNe hMatchValid hMatchNo rem1
      · have hsTailVal : [sHead].all isValidAtom = true := by
          simpa [hSfx, hTail] using hSfxValid
        have hsNoMulti : AtomSeg.multiWildcard ∉ [sHead] := by simp [Ne.symm hsHead]
        have hSfxEmit : emitUnmerged (sHead :: sTail) false true =
            [.literal "/"] ++ emitUnmerged (sHead :: sTail) false false :=
          emitUnmerged_suffix_cons sHead sTail hsHead
        rw [hTail] at hSfxEmit
        rw [matchPrefixAndCapture_ignore_hdr_pfx t (splitSegsOf rem1), hSfx, hTail]
        have hJoin :
            (do
              let (mCount, rem2) ← matchSegList (emitUnmerged t.matchSegs false false) rem1
              let (_, rem3) ← matchSegList (emitUnmerged [sHead] false true) rem2
              if !rem3.isEmpty || mCount == 0 then none
              else some (String.ofList (rem1.take mCount))) =
            (match matchSegList (emitUnmerged t.matchSegs false false ++ [.literal "/"]) rem1 with
              | none => none
              | some (n, remAfter) =>
                let mCount := n - 1
                match matchSegList (emitUnmerged [sHead] false false) remAfter with
                | some (sCount, []) =>
                  if sCount == 0 || mCount == 0 then none
                  else some (String.ofList (rem1.take mCount))
                | _ => none) := by
          rw [hSfxEmit, matchSegList_append]
          cases hM : matchSegList (emitUnmerged t.matchSegs false false) rem1 with
          | none => simp [bind, Option.bind]
          | some pM =>
            rcases pM with ⟨mCount, rem2⟩
            simp only [bind, Option.bind]
            rw [matchSegList_append]
            cases hSL : matchSegList [.literal "/"] rem2 with
            | none => simp
            | some pSL =>
              rcases pSL with ⟨kSlash, remAfter⟩
              have hk1 : kSlash = 1 := by
                cases rem2 with
                | nil => simp [matchSegList, RuntimeSeg.matchSize] at hSL
                | cons c cs' =>
                  by_cases hc : c = '/'
                  · subst hc
                    have hLit1 : matchSegList [RuntimeSeg.literal "/"] ('/' :: cs') = some (1, cs') := rfl
                    rw [hLit1] at hSL
                    cases hSL; rfl
                  · have hLitN : matchSegList [RuntimeSeg.literal "/"] (c :: cs') = none := by
                      simp [matchSegList, RuntimeSeg.matchSize, List.isPrefixOf, Ne.symm hc]
                    rw [hLitN] at hSL
                    cases hSL
              subst hk1
              dsimp only
              cases hSE : matchSegList (emitUnmerged [sHead] false false) remAfter with
              | none => simp
              | some pSE =>
                rcases pSE with ⟨sCount, rem3⟩
                have hsPos : 1 ≤ sCount :=
                  (emitUnmerged_cons_pos [sHead] (by simp) hsTailVal hsNoMulti remAfter hSE).1
                have hsNz : (sCount == 0) = false := by
                  cases sCount with
                  | zero => omega
                  | succ _ => simp
                cases rem3 with
                | nil => simp [hsNz]
                | cons _ _ => simp
        rw [hJoin]
        have hSlashSpec := matchUnmergedSlashEnd_spec t.matchSegs hMatchNe hMatchValid hMatchNo rem1
        have hExactSfx := fun r => matchUnmergedExact_spec [sHead] (by simp) hsTailVal hsNoMulti r
        have hNotMultiSfx : ([sHead] == [.multiWildcard]) = false := by
          cases sHead with
          | multiWildcard => exact absurd rfl hsHead
          | singleWildcard => rfl
          | lit _ => rfl
        unfold matchPrefixAndCapture
        simp only [hNotMultiSfx]
        by_cases hLenLt : t.matchSegs.length < (splitSegsOf rem1).length
        · have hNotSmall : ¬ (splitSegsOf rem1).length < t.matchSegs.length := by omega
          simp only [List.isEmpty_cons, Bool.false_eq_true, ↓reduceIte, hNotSmall]
          by_cases hMat : matchAtomSegs t.matchSegs ((splitSegsOf rem1).take t.matchSegs.length) = true
          · cases hSLRes : matchSegList (emitUnmerged t.matchSegs false false ++ [.literal "/"]) rem1 with
            | none =>
              simp [hSLRes, hLenLt, hMat] at hSlashSpec
            | some pSL =>
              rcases pSL with ⟨n, remAfter⟩
              simp only [hSLRes, hLenLt, decide_true, hMat, Bool.and_self, ↓reduceIte,
                Option.some.injEq, Prod.mk.injEq] at hSlashSpec
              rcases hSlashSpec with ⟨hTakeN, hSegsDrop⟩
              have hmPos : 1 ≤ n - 1 := by
                rw [matchSegList_append] at hSLRes
                cases hA : matchSegList (emitUnmerged t.matchSegs false false) rem1 with
                | none => simp [hA] at hSLRes
                | some pA =>
                  rcases pA with ⟨nx, remx⟩
                  have hNxPos := (emitUnmerged_cons_pos t.matchSegs hMatchNe hMatchValid hMatchNo rem1 hA).1
                  simp only [hA] at hSLRes
                  cases remx with
                  | nil => simp [matchSegList, RuntimeSeg.matchSize] at hSLRes
                  | cons c cs' =>
                    by_cases hc : c = '/'
                    · subst hc
                      have hLit1 : matchSegList [RuntimeSeg.literal "/"] ('/' :: cs') = some (1, cs') := rfl
                      simp only [hLit1, Option.some.injEq, Prod.mk.injEq] at hSLRes
                      omega
                    · have hLitN : matchSegList [RuntimeSeg.literal "/"] (c :: cs') = none := by
                        simp [matchSegList, RuntimeSeg.matchSize, List.isPrefixOf, Ne.symm hc]
                      simp [hLitN] at hSLRes
              have hmNz : (n - 1 == 0) = false := by
                cases hm : n - 1 with
                | zero => omega
                | succ _ => simp
              dsimp only
              rw [hTakeN]
              have hDropSfx := hExactSfx remAfter
              rw [hSegsDrop] at hDropSfx
              cases hSfxMatch : matchSegList (emitUnmerged [sHead] false false) remAfter with
              | none =>
                simp only [hSfxMatch] at hDropSfx ⊢
                have hCondF :
                    (((splitSegsOf rem1).drop t.matchSegs.length).length == [sHead].length &&
                      matchAtomSegs [sHead] ((splitSegsOf rem1).drop t.matchSegs.length)) = false := by
                  cases hC : ((splitSegsOf rem1).drop t.matchSegs.length).length == [sHead].length &&
                      matchAtomSegs [sHead] ((splitSegsOf rem1).drop t.matchSegs.length)
                  · rfl
                  · exfalso
                    simp only [hC, ↓reduceIte] at hDropSfx
                    simp only [Bool.and_eq_true, beq_iff_eq] at hC
                    have hNeS := matchAtomSegs_valid_join_ne_empty [sHead]
                      ((splitSegsOf rem1).drop t.matchSegs.length)
                      (by simp) hsTailVal hsNoMulti hC.1 hC.2
                    simp [beq_false_of_ne hNeS] at hDropSfx
                by_cases hSfxBool : matchAtomSegs [sHead] ((splitSegsOf rem1).drop t.matchSegs.length) = true
                · have hSfxLen := matchAtomSegs_length_of_no_multi [sHead]
                    ((splitSegsOf rem1).drop t.matchSegs.length) hsNoMulti hSfxBool
                  simp [hSfxLen, hSfxBool] at hCondF
                · simp [Bool.eq_false_iff.mpr hSfxBool]
              | some pSfx =>
                rcases pSfx with ⟨sCount, rem3⟩
                cases rem3 with
                | cons _ _ =>
                  simp only [hSfxMatch] at hDropSfx ⊢
                  have hCondF :
                      (((splitSegsOf rem1).drop t.matchSegs.length).length == [sHead].length &&
                        matchAtomSegs [sHead] ((splitSegsOf rem1).drop t.matchSegs.length)) = false := by
                    cases hC : ((splitSegsOf rem1).drop t.matchSegs.length).length == [sHead].length &&
                        matchAtomSegs [sHead] ((splitSegsOf rem1).drop t.matchSegs.length)
                    · rfl
                    · exfalso
                      simp only [hC, ↓reduceIte] at hDropSfx
                      simp only [Bool.and_eq_true, beq_iff_eq] at hC
                      have hNeS := matchAtomSegs_valid_join_ne_empty [sHead]
                        ((splitSegsOf rem1).drop t.matchSegs.length)
                        (by simp) hsTailVal hsNoMulti hC.1 hC.2
                      simp [beq_false_of_ne hNeS] at hDropSfx
                  by_cases hSfxBool : matchAtomSegs [sHead] ((splitSegsOf rem1).drop t.matchSegs.length) = true
                  · have hSfxLen := matchAtomSegs_length_of_no_multi [sHead]
                      ((splitSegsOf rem1).drop t.matchSegs.length) hsNoMulti hSfxBool
                    simp [hSfxLen, hSfxBool] at hCondF
                  · simp [Bool.eq_false_iff.mpr hSfxBool]
                | nil =>
                  simp only [hSfxMatch] at hDropSfx
                  have hsPos : 1 ≤ sCount :=
                    (emitUnmerged_cons_pos [sHead] (by simp) hsTailVal hsNoMulti remAfter hSfxMatch).1
                  have hsNz : (sCount == 0) = false := by
                    cases sCount with
                    | zero => omega
                    | succ _ => simp
                  simp only [hsNz, Bool.false_eq_true, ↓reduceIte] at hDropSfx
                  split at hDropSfx
                  · next hC =>
                    simp only [Bool.and_eq_true, beq_iff_eq] at hC
                    have hTakeLenEq : ((splitSegsOf rem1).take t.matchSegs.length).length = t.matchSegs.length := by
                      simp; omega
                    have hCapNe := matchAtomSegs_valid_join_ne_empty t.matchSegs
                      ((splitSegsOf rem1).take t.matchSegs.length)
                      hMatchNe hMatchValid hMatchNo hTakeLenEq hMat
                    simp [hsNz, hmNz, String.ofList_toList, hMat, hC.2, beq_false_of_ne hCapNe]
                  · simp at hDropSfx
          · have hMatF : matchAtomSegs t.matchSegs ((splitSegsOf rem1).take t.matchSegs.length) = false :=
              Bool.eq_false_iff.mpr hMat
            cases hSLRes : matchSegList (emitUnmerged t.matchSegs false false ++ [.literal "/"]) rem1 with
            | none => simp [hMatF]
            | some _ => simp [hSLRes, hMatF] at hSlashSpec
        · have hSLNone : matchSegList (emitUnmerged t.matchSegs false false ++ [.literal "/"]) rem1 = none := by
            cases hSLRes : matchSegList (emitUnmerged t.matchSegs false false ++ [.literal "/"]) rem1 with
            | none => rfl
            | some _ => simp [hSLRes, hLenLt] at hSlashSpec
          rw [hSLNone]
          dsimp only
          by_cases hLt : (splitSegsOf rem1).length < t.matchSegs.length
          · simp [hLt]
          · have hDropNil : (splitSegsOf rem1).drop t.matchSegs.length = [] := by
              rw [List.drop_eq_nil_iff]
              omega
            have hNilF : matchAtomSegs [sHead] [] = false := matchAtomSegs_nonmulti_nil hsHead []
            simp [hLt, hDropNil, hNilF]
    | cons sSecond sRest =>
      have hsHead : sHead ≠ .multiWildcard := by
        intro hEq
        apply hSfxDrop
        rw [hSfx, hTail, hEq]
        simp [List.dropLast]
      have hSfxTailValid : (sHead :: sSecond :: sRest).all isValidAtom = true := by
        simpa [hSfx, hTail] using hSfxValid
      have hSfxDropNo : AtomSeg.multiWildcard ∉ (sHead :: sSecond :: sRest).dropLast := by
        simpa [hSfx, hTail] using hSfxDrop
      have hSfxEmit : emitUnmerged (sHead :: sSecond :: sRest) false true =
          [.literal "/"] ++ emitUnmerged (sHead :: sSecond :: sRest) false false :=
        emitUnmerged_suffix_cons sHead (sSecond :: sRest) hsHead
      rw [matchPrefixAndCapture_ignore_hdr_pfx t (splitSegsOf rem1), hSfx, hTail]
      have hSfxSpec : ∀ (r : List Char),
          (match matchSegList (emitUnmerged (sHead :: sSecond :: sRest) false false) r with
            | some (n, []) => if n == 0 then none else some (String.ofList (r.take n))
            | _ => none) =
          (if matchAtomSegs (sHead :: sSecond :: sRest) (splitSegsOf r) then
            some (joinStringsWith "/" (splitSegsOf r)) else none) := by
        intro r
        rcases split_dropLast_no_multi (sHead :: sSecond :: sRest) hSfxDropNo with
          hSNo | hSBare | ⟨sInit, hsInitNe, hsInitNo, hsInitEq⟩
        · rw [matchUnmergedExact_spec (sHead :: sSecond :: sRest) (by simp) hSfxTailValid hSNo r]
          by_cases hm : matchAtomSegs (sHead :: sSecond :: sRest) (splitSegsOf r) = true
          · have hLen := matchAtomSegs_length_of_no_multi _ _ hSNo hm
            have hCapNe := matchAtomSegs_valid_join_ne_empty (sHead :: sSecond :: sRest)
              (splitSegsOf r) (by simp) hSfxTailValid hSNo hLen hm
            simp [hm, hLen, beq_false_of_ne hCapNe]
          · simp [Bool.eq_false_iff.mpr hm]
        · cases hSBare
        · have hsInitValid : sInit.all isValidAtom = true := by
            rw [List.all_eq_true] at hSfxTailValid ⊢
            intro a ha; apply hSfxTailValid; rw [hsInitEq]; exact List.mem_append_left _ ha
          rw [hsInitEq, matchUnmergedAppendMulti_spec sInit hsInitNe hsInitValid hsInitNo r]
          by_cases hm : matchAtomSegs (sInit ++ [.multiWildcard]) (splitSegsOf r) = true
          · have hCapNe := matchAtomSegs_append_multi_join_ne_empty sInit (splitSegsOf r)
              hsInitNe hsInitValid hsInitNo hm
            simp [hm, beq_false_of_ne hCapNe]
          · simp [Bool.eq_false_iff.mpr hm]
      have hSfxPos : ∀ (r : List Char) {sCount : Nat} {rem3 : List Char},
          matchSegList (emitUnmerged (sHead :: sSecond :: sRest) false false) r = some (sCount, rem3) →
          1 ≤ sCount := by
        intro r sCount rem3 hSE
        rcases split_dropLast_no_multi (sHead :: sSecond :: sRest) hSfxDropNo with
          hSNo | hSBare | ⟨sInit, hsInitNe, hsInitNo, hsInitEq⟩
        · exact (emitUnmerged_cons_pos _ (by simp) hSfxTailValid hSNo r hSE).1
        · cases hSBare
        · have hsInitValid : sInit.all isValidAtom = true := by
            rw [List.all_eq_true] at hSfxTailValid ⊢
            intro a ha; apply hSfxTailValid; rw [hsInitEq]; exact List.mem_append_left _ ha
          rw [hsInitEq] at hSE
          exact (emitUnmerged_append_multi_pos sInit hsInitNe hsInitValid hsInitNo r hSE).1
      have hJoin :
          (do
            let (mCount, rem2) ← matchSegList (emitUnmerged t.matchSegs false false) rem1
            let (_, rem3) ← matchSegList (emitUnmerged (sHead :: sSecond :: sRest) false true) rem2
            if !rem3.isEmpty || mCount == 0 then none
            else some (String.ofList (rem1.take mCount))) =
          (match matchSegList (emitUnmerged t.matchSegs false false ++ [.literal "/"]) rem1 with
            | none => none
            | some (n, remAfter) =>
              let mCount := n - 1
              match matchSegList (emitUnmerged (sHead :: sSecond :: sRest) false false) remAfter with
              | some (sCount, []) =>
                if sCount == 0 || mCount == 0 then none
                else some (String.ofList (rem1.take mCount))
              | _ => none) := by
        rw [hSfxEmit, matchSegList_append]
        cases hM : matchSegList (emitUnmerged t.matchSegs false false) rem1 with
        | none => simp [bind, Option.bind]
        | some pM =>
          rcases pM with ⟨mCount, rem2⟩
          simp only [bind, Option.bind]
          rw [matchSegList_append]
          cases hSL : matchSegList [.literal "/"] rem2 with
          | none => simp
          | some pSL =>
            rcases pSL with ⟨kSlash, remAfter⟩
            have hk1 : kSlash = 1 := by
              cases rem2 with
              | nil => simp [matchSegList, RuntimeSeg.matchSize] at hSL
              | cons c cs' =>
                by_cases hc : c = '/'
                · subst hc
                  have hLit1 : matchSegList [RuntimeSeg.literal "/"] ('/' :: cs') = some (1, cs') := rfl
                  rw [hLit1] at hSL
                  cases hSL; rfl
                · have hLitN : matchSegList [RuntimeSeg.literal "/"] (c :: cs') = none := by
                    simp [matchSegList, RuntimeSeg.matchSize, List.isPrefixOf, Ne.symm hc]
                  rw [hLitN] at hSL
                  cases hSL
            subst hk1
            dsimp only
            cases hSE : matchSegList (emitUnmerged (sHead :: sSecond :: sRest) false false) remAfter with
            | none => simp
            | some pSE =>
              rcases pSE with ⟨sCount, rem3⟩
              have hsPos : 1 ≤ sCount := hSfxPos remAfter hSE
              have hsNz : (sCount == 0) = false := by
                cases sCount with
                | zero => omega
                | succ _ => simp
              cases rem3 with
              | nil => simp [hsNz]
              | cons _ _ => simp
      rw [hJoin]
      have hSlashSpec := matchUnmergedSlashEnd_spec t.matchSegs hMatchNe hMatchValid hMatchNo rem1
      have hNotMultiSfx : (sHead :: sSecond :: sRest == [AtomSeg.multiWildcard]) = false := by
        cases sHead <;> rfl
      unfold matchPrefixAndCapture
      simp only [hNotMultiSfx, List.isEmpty_cons, Bool.false_eq_true, ↓reduceIte]
      by_cases hLenLt : t.matchSegs.length < (splitSegsOf rem1).length
      · have hNotSmall : ¬ (splitSegsOf rem1).length < t.matchSegs.length := by omega
        simp only [hNotSmall, ↓reduceIte]
        by_cases hMat : matchAtomSegs t.matchSegs ((splitSegsOf rem1).take t.matchSegs.length) = true
        · cases hSLRes : matchSegList (emitUnmerged t.matchSegs false false ++ [.literal "/"]) rem1 with
          | none => simp [hSLRes, hLenLt, hMat] at hSlashSpec
          | some pSL =>
            rcases pSL with ⟨n, remAfter⟩
            simp only [hSLRes, hLenLt, decide_true, hMat, Bool.and_self, ↓reduceIte,
              Option.some.injEq, Prod.mk.injEq] at hSlashSpec
            rcases hSlashSpec with ⟨hTakeN, hSegsDrop⟩
            have hmPos : 1 ≤ n - 1 := by
              rw [matchSegList_append] at hSLRes
              cases hA : matchSegList (emitUnmerged t.matchSegs false false) rem1 with
              | none => simp [hA] at hSLRes
              | some pA =>
                rcases pA with ⟨nx, remx⟩
                have hNxPos := (emitUnmerged_cons_pos t.matchSegs hMatchNe hMatchValid hMatchNo rem1 hA).1
                simp only [hA] at hSLRes
                cases remx with
                | nil => simp [matchSegList, RuntimeSeg.matchSize] at hSLRes
                | cons c cs' =>
                  by_cases hc : c = '/'
                  · subst hc
                    have hLit1 : matchSegList [RuntimeSeg.literal "/"] ('/' :: cs') = some (1, cs') := rfl
                    simp only [hLit1, Option.some.injEq, Prod.mk.injEq] at hSLRes
                    omega
                  · have hLitN : matchSegList [RuntimeSeg.literal "/"] (c :: cs') = none := by
                      simp [matchSegList, RuntimeSeg.matchSize, List.isPrefixOf, Ne.symm hc]
                    simp [hLitN] at hSLRes
            have hmNz : (n - 1 == 0) = false := by
              cases hm : n - 1 with
              | zero => omega
              | succ _ => simp
            dsimp only
            rw [hTakeN]
            have hDropSfx := hSfxSpec remAfter
            rw [hSegsDrop] at hDropSfx
            cases hSfxMatch : matchSegList (emitUnmerged (sHead :: sSecond :: sRest) false false) remAfter with
            | none =>
              simp only [hSfxMatch] at hDropSfx ⊢
              by_cases hSfxBool : matchAtomSegs (sHead :: sSecond :: sRest)
                  ((splitSegsOf rem1).drop t.matchSegs.length) = true
              · simp [hSfxBool] at hDropSfx
              · simp [Bool.eq_false_iff.mpr hSfxBool]
            | some pSfx =>
              rcases pSfx with ⟨sCount, rem3⟩
              cases rem3 with
              | cons _ _ =>
                simp only [hSfxMatch] at hDropSfx ⊢
                by_cases hSfxBool : matchAtomSegs (sHead :: sSecond :: sRest)
                    ((splitSegsOf rem1).drop t.matchSegs.length) = true
                · simp [hSfxBool] at hDropSfx
                · simp [Bool.eq_false_iff.mpr hSfxBool]
              | nil =>
                simp only [hSfxMatch] at hDropSfx
                have hsPos : 1 ≤ sCount := hSfxPos remAfter hSfxMatch
                have hsNz : (sCount == 0) = false := by
                  cases sCount with
                  | zero => omega
                  | succ _ => simp
                simp only [hsNz, Bool.false_eq_true, ↓reduceIte] at hDropSfx
                by_cases hC : matchAtomSegs (sHead :: sSecond :: sRest)
                    ((splitSegsOf rem1).drop t.matchSegs.length) = true
                · have hTakeLenEq : ((splitSegsOf rem1).take t.matchSegs.length).length = t.matchSegs.length := by
                    simp; omega
                  have hCapNe := matchAtomSegs_valid_join_ne_empty t.matchSegs
                    ((splitSegsOf rem1).take t.matchSegs.length)
                    hMatchNe hMatchValid hMatchNo hTakeLenEq hMat
                  simp [hsNz, hmNz, String.ofList_toList, hMat, hC, beq_false_of_ne hCapNe]
                · simp [Bool.eq_false_iff.mpr hC] at hDropSfx
        · have hMatF : matchAtomSegs t.matchSegs ((splitSegsOf rem1).take t.matchSegs.length) = false :=
            Bool.eq_false_iff.mpr hMat
          cases hSLRes : matchSegList (emitUnmerged t.matchSegs false false ++ [.literal "/"]) rem1 with
          | none => simp [hMatF]
          | some _ => simp [hSLRes, hMatF] at hSlashSpec
      · have hSLNone : matchSegList (emitUnmerged t.matchSegs false false ++ [.literal "/"]) rem1 = none := by
          cases hSLRes : matchSegList (emitUnmerged t.matchSegs false false ++ [.literal "/"]) rem1 with
          | none => rfl
          | some _ => simp [hSLRes, hLenLt] at hSlashSpec
        rw [hSLNone]
        dsimp only
        by_cases hLt : (splitSegsOf rem1).length < t.matchSegs.length
        · simp [hLt]
        · have hDropNil : (splitSegsOf rem1).drop t.matchSegs.length = [] := by
            rw [List.drop_eq_nil_iff]
            omega
          rw [hDropNil]
          have hNilF : matchAtomSegs (sHead :: sSecond :: sRest) [] = false :=
            matchAtomSegs_nonmulti_nil hsHead (sSecond :: sRest)
          simp [hLt, hNilF]

/--
**Universal Compiler Matcher Refinement**: for every valid template (`t.WellFormed`) and every
field value `val : String`, evaluating the compiled segment IR matches `extractCapture t val`.
-/
theorem compiled_matcher_refines_spec (t : CaptureTemplate) (val : String)
    (hWF : t.WellFormed) :
    evalCompiled (compileTemplate t) val = extractCapture t val := by
  rcases well_formed_decomp t hWF with
    ⟨_, hMatchNe, hPfxValid, _, _, hPfxNo, _, _, _⟩
  rw [evalCompiled_eq_unmerged t val]
  by_cases hValEmpty : val = ""
  · subst hValEmpty
    have hEmpList : ("" : String).toList = [] := rfl
    simp only [extractCapture, beq_self_eq_true, ↓reduceIte, hEmpList]
    cases hPfx : t.prefixSegs with
    | nil =>
      have hEmitNil : emitUnmerged [] true false = [] := rfl
      rw [hEmitNil]
      simp only [matchSegList, bind, Option.bind]
      cases hM : matchSegList (emitUnmerged t.matchSegs false false) [] with
      | none => rfl
      | some pM =>
        rcases pM with ⟨mCount, rem2⟩
        rcases matchSegList_drop_eq hM with ⟨_, hLe⟩
        have hm0 : mCount = 0 := by simp at hLe; exact hLe
        subst hm0
        dsimp only
        cases matchSegList (emitUnmerged t.suffixSegs false true) rem2 <;> simp
    | cons p ps =>
      have hEmitPfx : emitUnmerged (p :: ps) true false =
          emitUnmerged (p :: ps) false false ++ [.literal "/"] :=
        emitUnmerged_prefix (p :: ps) (by simp) (by simpa [hPfx] using hPfxNo)
      rw [hEmitPfx, matchSegList_append]
      cases hA : matchSegList (emitUnmerged (p :: ps) false false) [] with
      | none => simp [bind, Option.bind]
      | some _ =>
        have hBad := (emitUnmerged_cons_pos (p :: ps) (by simp)
          (by simpa [hPfx] using hPfxValid) (by simpa [hPfx] using hPfxNo) [] hA).2
        exact absurd rfl hBad
  · have hValEmptyF : (val == "") = false := beq_false_of_ne hValEmpty
    by_cases hBare : t.prefixSegs.isEmpty && (t.matchSegs == [.multiWildcard]) && t.suffixSegs.isEmpty = true
    · simp only [Bool.and_eq_true] at hBare
      rcases hBare with ⟨⟨hPfxNilB, hMatchBareB⟩, hSfxNilB⟩
      have hMatchBare : t.matchSegs = [AtomSeg.multiWildcard] := matchSegs_beq_single_multi hMatchBareB
      have hPfxNil : t.prefixSegs = [] := by cases hP : t.prefixSegs <;> simp_all
      have hSfxNil : t.suffixSegs = [] := by cases hS : t.suffixSegs <;> simp_all
      have hValListNe : val.toList ≠ [] := by
        intro hNil
        apply hValEmpty
        exact String.ext (by simp [hNil])
      rw [hPfxNil, hMatchBare, hSfxNil]
      simp [extractCapture, hValEmpty, hPfxNil, hMatchBare, hSfxNil, emitUnmerged,
        emitUnmergedLoop, matchSegList, RuntimeSeg.matchSize, bind, Option.bind,
        List.take_length, String.ofList_toList,
        show (AtomSeg.multiWildcard == AtomSeg.multiWildcard) = true from rfl]
    · have hBareF : (t.prefixSegs.isEmpty && (t.matchSegs == [.multiWildcard]) && t.suffixSegs.isEmpty) = false :=
        Bool.eq_false_iff.mpr hBare
      simp only [extractCapture, hValEmptyF, hBareF, Bool.false_eq_true, ↓reduceIte]
      cases hPfx : t.prefixSegs with
      | nil =>
        have hEmitNil : emitUnmerged [] true false = [] := rfl
        rw [hEmitNil]
        simp only [List.length_nil, Nat.not_lt_zero, ↓reduceIte, List.take_zero,
          matchAtomSegs, Bool.not_true, Bool.false_eq_true, List.drop_zero,
          bind, Option.bind, matchSegList]
        have hMEq := matchPrefixAndCapture_eq_compiled t hWF val.toList
        simp only [bind, Option.bind] at hMEq
        rw [hMEq, splitSegsOf_eq_splitStringOn, String.ofList_toList]
      | cons p ps =>
        have hPfxEmit : emitUnmerged (p :: ps) true false =
            emitUnmerged (p :: ps) false false ++ [.literal "/"] :=
          emitUnmerged_prefix (p :: ps) (by simp) (by simpa [hPfx] using hPfxNo)
        rw [hPfxEmit]
        have hPrefixSpec := matchUnmergedSlashEnd_spec (p :: ps) (by simp)
          (by simpa [hPfx] using hPfxValid) (by simpa [hPfx] using hPfxNo) val.toList
        rw [splitSegsOf_eq_splitStringOn, String.ofList_toList] at hPrefixSpec
        cases hPfxMatch : matchSegList (emitUnmerged (p :: ps) false false ++ [.literal "/"]) val.toList with
        | none =>
          simp only [hPfxMatch] at hPrefixSpec ⊢
          simp only [bind, Option.bind]
          by_cases hLt : (splitStringOn '/' val).length < (p :: ps).length
          · simp only [hLt, ↓reduceIte]
          · by_cases hMat : matchAtomSegs (p :: ps) ((splitStringOn '/' val).take (p :: ps).length) = true
            · have hNotGt : ¬ (p :: ps).length < (splitStringOn '/' val).length := by
                intro hGt
                simp only [hGt, decide_true, hMat, Bool.and_self, ↓reduceIte] at hPrefixSpec
                cases hPrefixSpec
              have hEqLen : (splitStringOn '/' val).length = (p :: ps).length := by omega
              have hDropNil : (splitStringOn '/' val).drop (p :: ps).length = [] := by
                rw [List.drop_eq_nil_iff]
                omega
              simp only [hLt, ↓reduceIte, hMat, Bool.not_true, Bool.false_eq_true, hDropNil]
              exact (matchPrefixAndCapture_nil t hMatchNe).symm
            · have hMatF : matchAtomSegs (p :: ps) ((splitStringOn '/' val).take (p :: ps).length) = false :=
                Bool.eq_false_iff.mpr hMat
              simp only [hLt, ↓reduceIte, hMatF, Bool.not_false]
        | some p1 =>
          rcases p1 with ⟨n1, rem1⟩
          simp only [hPfxMatch] at hPrefixSpec
          split at hPrefixSpec
          · next hCond =>
            simp only [Option.some.injEq, Prod.mk.injEq] at hPrefixSpec
            rcases hPrefixSpec with ⟨_, hSegs1⟩
            simp only [Bool.and_eq_true] at hCond
            have hNotLt : ¬ (splitStringOn '/' val).length < (p :: ps).length := by
              have hGt : (p :: ps).length < (splitStringOn '/' val).length := decide_eq_true_eq.mp hCond.1
              omega
            simp only [hNotLt, ↓reduceIte, hCond.2, Bool.not_true, Bool.false_eq_true,
              bind, Option.bind]
            have hMEq := matchPrefixAndCapture_eq_compiled t hWF rem1
            simp only [bind, Option.bind] at hMEq
            rw [hMEq, hSegs1]
          · cases hPrefixSpec

/-!
## Group 4: Conflict Resolution & Key-Order Subtlety
-/

private theorem firstMatchWins_append {α : Type} (xs ys : List (Option α)) :
    firstMatchWins (xs ++ ys) =
      match firstMatchWins xs with
      | some v => some v
      | none => firstMatchWins ys := by
  induction xs with
  | nil => simp [firstMatchWins]
  | cons x rest ih =>
    cases x with
    | none => simp only [List.cons_append, firstMatchWins]; exact ih
    | some v => simp [firstMatchWins]

private theorem lastMatchWins_foldl_eq {α : Type} (init : Option α) (vals : List (Option α)) :
    vals.foldl (fun acc x => match x with | some v => some v | none => acc) init =
      match firstMatchWins vals.reverse with
      | some v => some v
      | none => init := by
  induction vals generalizing init with
  | nil => rfl
  | cons x rest ih =>
    simp only [List.foldl_cons, List.reverse_cons]
    rw [ih, firstMatchWins_append]
    cases firstMatchWins rest.reverse with
    | some v => rfl
    | none => cases x <;> rfl

/-- Reversing per-key variants converts declaration-order "last match wins" into `.or_else` "first match wins". -/
theorem reverse_first_match_eq_last_wins {α : Type} (vals : List (Option α)) :
    firstMatchWins vals.reverse = lastMatchWins vals := by
  have h := lastMatchWins_foldl_eq none vals
  exact Eq.trans (by cases firstMatchWins vals.reverse <;> rfl) h.symm

private theorem compileParams_error_iff_evalAux_error (req : Request) :
    ∀ (params : List RoutingParameter) (groups : List CompiledKeyGroup)
      (acc : List (String × String)) (e : String),
      compileRoutingParamsDeclOrder params groups = .error e ↔
        evalRoutingParamsAux req params acc = .error e
  | [], _, _, _ => by simp [compileRoutingParamsDeclOrder, evalRoutingParamsAux]
  | param :: rest, groups, acc, e => by
    unfold compileRoutingParamsDeclOrder evalRoutingParamsAux
    cases hParse : parseRoutingParameter param.field param.pathTemplate with
    | error err => simp [bind, Except.bind]
    | ok tmpl =>
      simp only [bind, Except.bind]
      exact compileParams_error_iff_evalAux_error req rest _ _ e

/-- Any parse error during rule compilation agrees identically with `evalRoutingRule`. -/
theorem compileRoutingRuleDeclOrder_error_eq (rule : RoutingRule) (req : Request) (e : String)
    (hErr : compileRoutingRuleDeclOrder rule = .error e) :
    evalRoutingRule rule req = .error e := by
  unfold evalRoutingRule evalRoutingRulePairs
  unfold compileRoutingRuleDeclOrder at hErr
  have hAux := (compileParams_error_iff_evalAux_error req rule.parameters [] [] e).mp hErr
  rw [hAux]
  rfl

private def groupLookup (k : String) (req : Request) : List CompiledKeyGroup → Option String
  | [] => none
  | g :: rest =>
    if g.headerKey == k then
      firstMatchWins (g.variants.map (fun v => evalCompiledVariant v req))
    else
      groupLookup k req rest

private theorem mem_map_headerKey_insertVariant (groups : List CompiledKeyGroup)
    (fp : FieldPath) (ct : CompiledTemplate) (k : String) :
    k ∈ (insertVariantDeclOrder groups fp ct).map CompiledKeyGroup.headerKey ↔
      k = ct.headerKey ∨ k ∈ groups.map CompiledKeyGroup.headerKey := by
  induction groups with
  | nil => simp [insertVariantDeclOrder]
  | cons g rest ih =>
    unfold insertVariantDeclOrder
    by_cases hEq : (g.headerKey == ct.headerKey) = true
    · have hSame : g.headerKey = ct.headerKey := eq_of_beq hEq
      simp only [hEq, ↓reduceIte, List.map_cons, List.mem_cons]
      constructor
      · intro h; rcases h with rfl | hR
        · left; exact hSame
        · right; right; exact hR
      · intro h; rcases h with rfl | hR
        · left; exact hSame.symm
        · exact hR
    · have hEqF : (g.headerKey == ct.headerKey) = false := Bool.eq_false_iff.mpr hEq
      simp only [hEqF, Bool.false_eq_true, ↓reduceIte, List.map_cons, List.mem_cons, ih]
      constructor
      · rintro (rfl | rfl | hR)
        · right; left; rfl
        · left; rfl
        · right; right; exact hR
      · rintro (rfl | rfl | hR)
        · right; left; rfl
        · left; rfl
        · right; right; exact hR

private theorem insertVariantDeclOrder_nodup (groups : List CompiledKeyGroup)
    (fp : FieldPath) (ct : CompiledTemplate)
    (hNodup : (groups.map CompiledKeyGroup.headerKey).Nodup) :
    ((insertVariantDeclOrder groups fp ct).map CompiledKeyGroup.headerKey).Nodup := by
  induction groups with
  | nil => simp [insertVariantDeclOrder]
  | cons g rest ih =>
    rw [List.map_cons, List.nodup_cons] at hNodup
    rcases hNodup with ⟨hNotMem, hRestNodup⟩
    unfold insertVariantDeclOrder
    by_cases hEq : (g.headerKey == ct.headerKey) = true
    · simp [hEq, List.nodup_cons, hNotMem, hRestNodup]
    · have hEqF : (g.headerKey == ct.headerKey) = false := Bool.eq_false_iff.mpr hEq
      have hNe : g.headerKey ≠ ct.headerKey := fun heq => hEq (by simp [heq])
      simp only [hEqF, Bool.false_eq_true, ↓reduceIte, List.map_cons, List.nodup_cons]
      refine ⟨?_, ih hRestNodup⟩
      rw [mem_map_headerKey_insertVariant]
      rintro (hEq2 | hIn)
      · exact hNe hEq2
      · exact hNotMem hIn

private theorem evalCompiledKeyGroupsPairs_fst_subset (groups : List CompiledKeyGroup)
    (req : Request) (k : String) :
    k ∈ (evalCompiledKeyGroupsPairs groups req).map Prod.fst →
      k ∈ groups.map CompiledKeyGroup.headerKey := by
  induction groups with
  | nil => simp [evalCompiledKeyGroupsPairs]
  | cons g rest ih =>
    unfold evalCompiledKeyGroupsPairs
    simp only [List.filterMap_cons, List.map_cons, List.mem_cons]
    cases hG : evalCompiledKeyGroup g req with
    | none => intro h; exact Or.inr (ih h)
    | some kv =>
      unfold evalCompiledKeyGroup at hG
      cases hm : firstMatchWins (g.variants.map (fun v => evalCompiledVariant v req)) with
      | none => simp [hm] at hG
      | some v =>
        simp only [hm, Option.map_some, Option.some.injEq] at hG
        subst hG
        simp only [List.map_cons, List.mem_cons]
        rintro (rfl | hR)
        · left; rfl
        · right; exact ih hR

private theorem evalCompiledKeyGroupsPairs_nodup (groups : List CompiledKeyGroup) (req : Request)
    (hNodup : (groups.map CompiledKeyGroup.headerKey).Nodup) :
    NoDupKeys (evalCompiledKeyGroupsPairs groups req) := by
  unfold NoDupKeys
  induction groups with
  | nil => simp [evalCompiledKeyGroupsPairs]
  | cons g rest ih =>
    rw [List.map_cons, List.nodup_cons] at hNodup
    rcases hNodup with ⟨hNotMem, hRestNodup⟩
    unfold evalCompiledKeyGroupsPairs
    simp only [List.filterMap_cons]
    cases hG : evalCompiledKeyGroup g req with
    | none => exact ih hRestNodup
    | some kv =>
      unfold evalCompiledKeyGroup at hG
      cases hm : firstMatchWins (g.variants.map (fun v => evalCompiledVariant v req)) with
      | none => simp [hm] at hG
      | some v =>
        simp only [hm, Option.map_some, Option.some.injEq] at hG
        subst hG
        simp only [List.map_cons, List.nodup_cons]
        refine ⟨?_, ih hRestNodup⟩
        intro hIn
        exact hNotMem (evalCompiledKeyGroupsPairs_fst_subset rest req g.headerKey hIn)

private theorem lookupHeaderKey_eq_none_of_not_mem (k : String) (pairs : List (String × String))
    (hNotMem : k ∉ pairs.map Prod.fst) :
    lookupHeaderKey k pairs = none := by
  induction pairs with
  | nil => rfl
  | cons kv rest ih =>
    simp only [List.map_cons, List.mem_cons, not_or] at hNotMem
    unfold lookupHeaderKey
    have hNe : (kv.1 == k) = false := beq_false_of_ne (Ne.symm hNotMem.1)
    simp only [hNe, Bool.false_eq_true, ↓reduceIte]
    exact ih hNotMem.2

private theorem lookupHeaderKey_evalCompiledKeyGroupsPairs (groups : List CompiledKeyGroup)
    (req : Request) (hNodup : (groups.map CompiledKeyGroup.headerKey).Nodup) (k : String) :
    lookupHeaderKey k (evalCompiledKeyGroupsPairs groups req) = groupLookup k req groups := by
  induction groups with
  | nil => rfl
  | cons g rest ih =>
    rw [List.map_cons, List.nodup_cons] at hNodup
    rcases hNodup with ⟨hNotMem, hRestNodup⟩
    unfold evalCompiledKeyGroupsPairs groupLookup
    simp only [List.filterMap_cons]
    by_cases hEq : (g.headerKey == k) = true
    · have hSame : g.headerKey = k := eq_of_beq hEq
      simp only [hEq, ↓reduceIte]
      unfold evalCompiledKeyGroup
      cases hm : firstMatchWins (g.variants.map (fun v => evalCompiledVariant v req)) with
      | none =>
        simp only [Option.map_none]
        apply lookupHeaderKey_eq_none_of_not_mem
        intro hIn
        apply hNotMem
        rw [hSame]
        exact evalCompiledKeyGroupsPairs_fst_subset rest req k hIn
      | some v =>
        simp [hEq, lookupHeaderKey]
    · have hEqF : (g.headerKey == k) = false := Bool.eq_false_iff.mpr hEq
      simp only [hEqF, Bool.false_eq_true, ↓reduceIte]
      unfold evalCompiledKeyGroup
      cases hm : firstMatchWins (g.variants.map (fun v => evalCompiledVariant v req)) with
      | none => exact ih hRestNodup
      | some v =>
        simp only [Option.map_some, lookupHeaderKey, hEqF, Bool.false_eq_true, ↓reduceIte]
        exact ih hRestNodup

private theorem groupLookup_insertVariantDeclOrder (groups : List CompiledKeyGroup)
    (fp : FieldPath) (ct : CompiledTemplate) (req : Request)
    (hNodup : (groups.map CompiledKeyGroup.headerKey).Nodup) (k : String) :
    groupLookup k req (insertVariantDeclOrder groups fp ct) =
      if ct.headerKey == k then
        match evalCompiledVariant ⟨fp, ct⟩ req with
        | some v => some v
        | none => groupLookup k req groups
      else groupLookup k req groups := by
  induction groups with
  | nil =>
    unfold insertVariantDeclOrder groupLookup
    by_cases hEq : (ct.headerKey == k) = true
    · simp only [hEq, ↓reduceIte, List.map_cons, List.map_nil]
      cases evalCompiledVariant ⟨fp, ct⟩ req <;> simp [firstMatchWins]
    · have hEqF : (ct.headerKey == k) = false := Bool.eq_false_iff.mpr hEq
      simp [hEqF, groupLookup]
  | cons g rest ih =>
    rw [List.map_cons, List.nodup_cons] at hNodup
    rcases hNodup with ⟨_, hRestNodup⟩
    unfold insertVariantDeclOrder
    by_cases hEqG : (g.headerKey == ct.headerKey) = true
    · have hSame : g.headerKey = ct.headerKey := eq_of_beq hEqG
      simp only [hEqG, ↓reduceIte]
      unfold groupLookup
      by_cases hEqK : (ct.headerKey == k) = true
      · have hGK : (g.headerKey == k) = true := by rw [hSame]; exact hEqK
        simp only [hEqK, hGK, ↓reduceIte, List.map_cons]
        cases evalCompiledVariant ⟨fp, ct⟩ req with
        | none => rfl
        | some v => simp [firstMatchWins]
      · have hEqKF : (ct.headerKey == k) = false := Bool.eq_false_iff.mpr hEqK
        have hGKF : (g.headerKey == k) = false := by rw [hSame]; exact hEqKF
        simp [hEqKF, hGKF]
    · have hEqGF : (g.headerKey == ct.headerKey) = false := Bool.eq_false_iff.mpr hEqG
      simp only [hEqGF, Bool.false_eq_true, ↓reduceIte]
      unfold groupLookup
      by_cases hGK : (g.headerKey == k) = true
      · have hSameGK : g.headerKey = k := eq_of_beq hGK
        have hCKF : (ct.headerKey == k) = false := by
          apply beq_false_of_ne
          intro hCEq
          apply hEqG
          rw [hSameGK, hCEq]
          simp
        simp [hGK, hCKF]
      · have hGKF : (g.headerKey == k) = false := Bool.eq_false_iff.mpr hGK
        simp only [hGKF, Bool.false_eq_true, ↓reduceIte]
        exact ih hRestNodup

private theorem upsertHeaderPair_nodup (pairs : List (String × String)) (k v : String)
    (hNodup : NoDupKeys pairs) :
    NoDupKeys (upsertHeaderPair pairs k v) := by
  unfold NoDupKeys at hNodup ⊢
  induction pairs with
  | nil => simp [upsertHeaderPair]
  | cons kv rest ih =>
    rw [List.map_cons, List.nodup_cons] at hNodup
    rcases hNodup with ⟨hNotMem, hRestNodup⟩
    unfold upsertHeaderPair
    by_cases hEq : (kv.1 == k) = true
    · have hSame : kv.1 = k := eq_of_beq hEq
      simp only [hEq, ↓reduceIte, List.map_cons, List.nodup_cons]
      refine ⟨?_, hRestNodup⟩
      simpa [hSame] using hNotMem
    · have hEqF : (kv.1 == k) = false := Bool.eq_false_iff.mpr hEq
      have hNe : kv.1 ≠ k := fun heq => hEq (by simp [heq])
      simp only [hEqF, Bool.false_eq_true, ↓reduceIte, List.map_cons, List.nodup_cons]
      refine ⟨?_, ih hRestNodup⟩
      intro hIn
      apply hNotMem
      clear ih hRestNodup hNotMem
      induction rest with
      | nil =>
        simp only [upsertHeaderPair, List.map_cons, List.map_nil, List.mem_cons,
          List.not_mem_nil, or_false] at hIn
        exact absurd hIn hNe
      | cons kv2 r2 ih2 =>
        unfold upsertHeaderPair at hIn
        by_cases h2 : (kv2.1 == k) = true
        · have hS2 : kv2.1 = k := eq_of_beq h2
          simp only [h2, ↓reduceIte, List.map_cons, List.mem_cons] at hIn
          simp only [List.map_cons, List.mem_cons]
          rcases hIn with hEqKv | hInR2
          · exact absurd hEqKv hNe
          · exact Or.inr hInR2
        · have h2F : (kv2.1 == k) = false := Bool.eq_false_iff.mpr h2
          simp only [h2F, Bool.false_eq_true, ↓reduceIte, List.map_cons, List.mem_cons] at hIn
          simp only [List.map_cons, List.mem_cons]
          rcases hIn with hEqKv | hInR2
          · exact Or.inl hEqKv
          · exact Or.inr (ih2 hInR2)

private theorem lookupHeaderKey_upsertHeaderPair (pairs : List (String × String)) (k v q : String) :
    lookupHeaderKey q (upsertHeaderPair pairs k v) =
      if k == q then some v else lookupHeaderKey q pairs := by
  induction pairs with
  | nil =>
    unfold upsertHeaderPair lookupHeaderKey
    by_cases hEq : (k == q) = true
    · simp [hEq]
    · have hEqF : (k == q) = false := Bool.eq_false_iff.mpr hEq
      simp [hEqF, lookupHeaderKey]
  | cons kv rest ih =>
    unfold upsertHeaderPair
    by_cases hEq : (kv.1 == k) = true
    · have hSame : kv.1 = k := eq_of_beq hEq
      simp only [hEq, ↓reduceIte, lookupHeaderKey]
      rw [hSame]
      by_cases hKQ : (k == q) = true <;> simp [hKQ]
    · have hEqF : (kv.1 == k) = false := Bool.eq_false_iff.mpr hEq
      simp only [hEqF, Bool.false_eq_true, ↓reduceIte, lookupHeaderKey]
      by_cases hKvQ : (kv.1 == q) = true
      · have hSameKvQ : kv.1 = q := eq_of_beq hKvQ
        have hKQF : (k == q) = false := by
          apply beq_false_of_ne
          intro hKQ
          apply hEq
          rw [hSameKvQ, hKQ]
          simp
        simp [hKvQ, hKQF]
      · have hKvQF : (kv.1 == q) = false := Bool.eq_false_iff.mpr hKvQ
        simp [hKvQF, ih]

private theorem compiled_params_extensional_aux (req : Request) :
    ∀ (params : List RoutingParameter) (groups : List CompiledKeyGroup)
      (specAcc : List (String × String)),
      (groups.map CompiledKeyGroup.headerKey).Nodup →
      NoDupKeys specAcc →
      (∀ k, groupLookup k req groups = lookupHeaderKey k specAcc) →
      match compileRoutingParamsDeclOrder params groups, evalRoutingParamsAux req params specAcc with
      | .ok finalGroups, .ok specPairs =>
        let compiledPairs := evalCompiledKeyGroupsPairs finalGroups req
        NoDupKeys compiledPairs ∧
        NoDupKeys specPairs ∧
        ∀ k, lookupHeaderKey k compiledPairs = lookupHeaderKey k specPairs
      | .error e1, .error e2 => e1 = e2
      | _, _ => False
  | [], groups, specAcc, hGrpNodup, hAccNodup, hAgr => by
    simp only [compileRoutingParamsDeclOrder, evalRoutingParamsAux]
    refine ⟨evalCompiledKeyGroupsPairs_nodup groups req hGrpNodup, hAccNodup, fun k => ?_⟩
    rw [lookupHeaderKey_evalCompiledKeyGroupsPairs groups req hGrpNodup k, hAgr k]
  | param :: rest, groups, specAcc, hGrpNodup, hAccNodup, hAgr => by
    unfold compileRoutingParamsDeclOrder evalRoutingParamsAux
    cases hParse : parseRoutingParameter param.field param.pathTemplate with
    | error err => simp [bind, Except.bind]
    | ok tmpl =>
      simp only [bind, Except.bind]
      have hWF : tmpl.WellFormed := parseRoutingParameter_sound param.field param.pathTemplate tmpl hParse
      let ct := compileTemplate tmpl
      let fp := FieldPath.parse param.field
      have hKeyEq : ct.headerKey = tmpl.headerKey := rfl
      have hVarEq : evalCompiledVariant ⟨fp, ct⟩ req =
          match req.getField fp with
          | none => none
          | some val => extractCapture tmpl val := by
        unfold evalCompiledVariant
        cases req.getField fp with
        | none => rfl
        | some val => exact compiled_matcher_refines_spec tmpl val hWF
      have hGrpNodup' := insertVariantDeclOrder_nodup groups fp ct hGrpNodup
      cases hField : req.getField fp with
      | none =>
        dsimp only
        have hVarNone : evalCompiledVariant ⟨fp, ct⟩ req = none := by
          simpa [hField] using hVarEq
        apply compiled_params_extensional_aux req rest (insertVariantDeclOrder groups fp ct)
          specAcc hGrpNodup' hAccNodup
        intro k
        rw [groupLookup_insertVariantDeclOrder groups fp ct req hGrpNodup k, hVarNone]
        dsimp only
        cases (ct.headerKey == k) <;> simp [hAgr k]
      | some val =>
        dsimp only
        cases hCap : extractCapture tmpl val with
        | none =>
          dsimp only
          have hVarNone : evalCompiledVariant ⟨fp, ct⟩ req = none := by
            simpa [hField, hCap] using hVarEq
          apply compiled_params_extensional_aux req rest (insertVariantDeclOrder groups fp ct)
            specAcc hGrpNodup' hAccNodup
          intro k
          rw [groupLookup_insertVariantDeclOrder groups fp ct req hGrpNodup k, hVarNone]
          dsimp only
          cases (ct.headerKey == k) <;> simp [hAgr k]
        | some captured =>
          dsimp only
          have hVarSome : evalCompiledVariant ⟨fp, ct⟩ req = some captured := by
            simpa [hField, hCap] using hVarEq
          apply compiled_params_extensional_aux req rest (insertVariantDeclOrder groups fp ct)
            (upsertHeaderPair specAcc tmpl.headerKey captured) hGrpNodup'
            (upsertHeaderPair_nodup specAcc tmpl.headerKey captured hAccNodup)
          intro k
          rw [groupLookup_insertVariantDeclOrder groups fp ct req hGrpNodup k, hVarSome,
            lookupHeaderKey_upsertHeaderPair, hKeyEq]
          dsimp only
          by_cases hEq : (tmpl.headerKey == k) = true
          · simp [hEq]
          · have hEqF : (tmpl.headerKey == k) = false := Bool.eq_false_iff.mpr hEq
            simp [hEqF, hAgr k]

/--
Extensional equivalence between the compiled key-group pipeline (`evalCompiledKeyGroupsPairs`) and
the declarative specification (`evalRoutingRulePairs`) over all header keys and uniqueness invariants.
-/
theorem compiled_rule_extensional_eq_spec (rule : RoutingRule) (req : Request) :
    match compileRoutingRuleDeclOrder rule, evalRoutingRulePairs rule req with
    | .ok groups, .ok specPairs =>
      let compiledPairs := evalCompiledKeyGroupsPairs groups req
      NoDupKeys compiledPairs ∧
      NoDupKeys specPairs ∧
      ∀ k, lookupHeaderKey k compiledPairs = lookupHeaderKey k specPairs
    | .error e1, .error e2 => e1 = e2
    | _, _ => False := by
  unfold compileRoutingRuleDeclOrder evalRoutingRulePairs
  exact compiled_params_extensional_aux req rule.parameters [] []
    (by simp) (by simp [NoDupKeys]) (fun _ => rfl)

private theorem insertVariant_preserves_ge_two (groups : List CompiledKeyGroup)
    (fp : FieldPath) (ct : CompiledTemplate)
    (hEx : ∃ g ∈ groups, 2 ≤ g.variants.length) :
    ∃ g ∈ insertVariantDeclOrder groups fp ct, 2 ≤ g.variants.length := by
  induction groups with
  | nil => rcases hEx with ⟨_, hMem, _⟩; cases hMem
  | cons g rest ih =>
    unfold insertVariantDeclOrder
    rcases hEx with ⟨g0, hMem, hGe2⟩
    simp only [List.mem_cons] at hMem
    by_cases hEq : (g.headerKey == ct.headerKey) = true
    · simp only [hEq, ↓reduceIte]
      rcases hMem with rfl | hRest
      · refine ⟨{ g0 with variants := { fieldPath := fp, compiled := ct } :: g0.variants },
          by simp, ?_⟩
        simp; omega
      · exact ⟨g0, by simp [hRest], hGe2⟩
    · have hEqF : (g.headerKey == ct.headerKey) = false := Bool.eq_false_iff.mpr hEq
      simp only [hEqF, Bool.false_eq_true, ↓reduceIte]
      rcases hMem with rfl | hRest
      · exact ⟨g0, by simp, hGe2⟩
      · rcases ih ⟨g0, hRest, hGe2⟩ with ⟨g', hG', hGe'⟩
        exact ⟨g', by simp [hG'], hGe'⟩

private theorem compileParams_preserves_ge_two (params : List RoutingParameter)
    (groups finalGroups : List CompiledKeyGroup)
    (hComp : compileRoutingParamsDeclOrder params groups = .ok finalGroups)
    (hEx : ∃ g ∈ groups, 2 ≤ g.variants.length) :
    ∃ g ∈ finalGroups, 2 ≤ g.variants.length := by
  induction params generalizing groups with
  | nil =>
    simp only [compileRoutingParamsDeclOrder, Except.ok.injEq] at hComp
    subst hComp
    exact hEx
  | cons param rest ih =>
    unfold compileRoutingParamsDeclOrder at hComp
    cases hParse : parseRoutingParameter param.field param.pathTemplate with
    | error _ => simp [hParse, bind, Except.bind] at hComp
    | ok tmpl =>
      simp only [hParse, bind, Except.bind] at hComp
      exact ih _ hComp (insertVariant_preserves_ge_two groups _ _ hEx)

private theorem insertVariant_mem_ge_two_of_mem_key (groups : List CompiledKeyGroup)
    (fp : FieldPath) (ct : CompiledTemplate)
    (hNonEmpty : ∀ g ∈ groups, 1 ≤ g.variants.length)
    (hIn : ct.headerKey ∈ groups.map CompiledKeyGroup.headerKey) :
    ∃ g ∈ insertVariantDeclOrder groups fp ct, 2 ≤ g.variants.length := by
  induction groups with
  | nil => cases hIn
  | cons g rest ih =>
    simp only [List.map_cons, List.mem_cons] at hIn
    unfold insertVariantDeclOrder
    by_cases hEq : (g.headerKey == ct.headerKey) = true
    · simp only [hEq, ↓reduceIte]
      have hG1 := hNonEmpty g (by simp)
      refine ⟨{ g with variants := { fieldPath := fp, compiled := ct } :: g.variants },
        by simp, ?_⟩
      simp; omega
    · have hEqF : (g.headerKey == ct.headerKey) = false := Bool.eq_false_iff.mpr hEq
      have hNe : g.headerKey ≠ ct.headerKey := fun heq => hEq (by simp [heq])
      simp only [hEqF, Bool.false_eq_true, ↓reduceIte]
      rcases hIn with hHead | hTail
      · exact absurd hHead.symm hNe
      · have hRestNE : ∀ g0 ∈ rest, 1 ≤ g0.variants.length := fun g0 h0 => hNonEmpty g0 (by simp [h0])
        rcases ih hRestNE hTail with ⟨g', hG', hGe'⟩
        exact ⟨g', by simp [hG'], hGe'⟩

private theorem insertVariant_nonempty (groups : List CompiledKeyGroup)
    (fp : FieldPath) (ct : CompiledTemplate)
    (hNonEmpty : ∀ g ∈ groups, 1 ≤ g.variants.length) :
    ∀ g ∈ insertVariantDeclOrder groups fp ct, 1 ≤ g.variants.length := by
  induction groups with
  | nil => simp [insertVariantDeclOrder]
  | cons g rest ih =>
    unfold insertVariantDeclOrder
    by_cases hEq : (g.headerKey == ct.headerKey) = true
    · simp only [hEq, ↓reduceIte, List.mem_cons]
      rintro g0 (rfl | hR)
      · simp
      · exact hNonEmpty g0 (by simp [hR])
    · have hEqF : (g.headerKey == ct.headerKey) = false := Bool.eq_false_iff.mpr hEq
      simp only [hEqF, Bool.false_eq_true, ↓reduceIte, List.mem_cons]
      rintro g0 (rfl | hR)
      · exact hNonEmpty g0 (by simp)
      · exact ih (fun x hx => hNonEmpty x (by simp [hx])) g0 hR

private theorem insertVariant_append_of_not_mem (groups : List CompiledKeyGroup)
    (fp : FieldPath) (ct : CompiledTemplate)
    (hNotMem : ct.headerKey ∉ groups.map CompiledKeyGroup.headerKey) :
    insertVariantDeclOrder groups fp ct =
      groups ++ [{ headerKey := ct.headerKey, variants := [{ fieldPath := fp, compiled := ct }] }] := by
  induction groups with
  | nil => rfl
  | cons g rest ih =>
    simp only [List.map_cons, List.mem_cons, not_or] at hNotMem
    unfold insertVariantDeclOrder
    have hNe : (g.headerKey == ct.headerKey) = false := beq_false_of_ne (Ne.symm hNotMem.1)
    simp only [hNe, Bool.false_eq_true, ↓reduceIte, List.cons_append]
    rw [ih hNotMem.2]

private theorem upsertHeaderPair_append_of_not_mem (pairs : List (String × String))
    (k v : String) (hNotMem : k ∉ pairs.map Prod.fst) :
    upsertHeaderPair pairs k v = pairs ++ [(k, v)] := by
  induction pairs with
  | nil => rfl
  | cons kv rest ih =>
    simp only [List.map_cons, List.mem_cons, not_or] at hNotMem
    unfold upsertHeaderPair
    have hNe : (kv.1 == k) = false := beq_false_of_ne (Ne.symm hNotMem.1)
    simp only [hNe, Bool.false_eq_true, ↓reduceIte, List.cons_append]
    rw [ih hNotMem.2]

private theorem compiled_params_single_variant_aux (req : Request) :
    ∀ (params : List RoutingParameter) (groups finalGroups : List CompiledKeyGroup)
      (specAcc : List (String × String)),
      compileRoutingParamsDeclOrder params groups = .ok finalGroups →
      (∀ g ∈ groups, 1 ≤ g.variants.length) →
      (∀ g ∈ finalGroups, g.variants.length ≤ 1) →
      evalCompiledKeyGroupsPairs groups req = specAcc →
      evalRoutingParamsAux req params specAcc = .ok (evalCompiledKeyGroupsPairs finalGroups req)
  | [], groups, finalGroups, specAcc, hComp, _, _, hEqPairs => by
    simp only [compileRoutingParamsDeclOrder, Except.ok.injEq] at hComp
    subst hComp
    simp [evalRoutingParamsAux, hEqPairs]
  | param :: rest, groups, finalGroups, specAcc, hComp, hNonEmpty, hSingle, hEqPairs => by
    unfold compileRoutingParamsDeclOrder at hComp
    unfold evalRoutingParamsAux
    cases hParse : parseRoutingParameter param.field param.pathTemplate with
    | error _ => simp [hParse, bind, Except.bind] at hComp
    | ok tmpl =>
      simp only [hParse, bind, Except.bind] at hComp ⊢
      have hWF : tmpl.WellFormed := parseRoutingParameter_sound param.field param.pathTemplate tmpl hParse
      let ct := compileTemplate tmpl
      let fp := FieldPath.parse param.field
      have hFresh : ct.headerKey ∉ groups.map CompiledKeyGroup.headerKey := by
        intro hIn
        have hGe2 := insertVariant_mem_ge_two_of_mem_key groups fp ct hNonEmpty hIn
        have hFinal := compileParams_preserves_ge_two rest (insertVariantDeclOrder groups fp ct)
          finalGroups hComp hGe2
        rcases hFinal with ⟨gBad, hMemBad, hLen2⟩
        have hLen1 := hSingle gBad hMemBad
        omega
      have hIns := insertVariant_append_of_not_mem groups fp ct hFresh
      have hNotMemAcc : tmpl.headerKey ∉ specAcc.map Prod.fst := by
        intro hIn
        apply hFresh
        rw [← hEqPairs] at hIn
        exact evalCompiledKeyGroupsPairs_fst_subset groups req tmpl.headerKey hIn
      have hVarEq : evalCompiledVariant ⟨fp, ct⟩ req =
          match req.getField fp with
          | none => none
          | some val => extractCapture tmpl val := by
        unfold evalCompiledVariant
        cases req.getField fp with
        | none => rfl
        | some val => exact compiled_matcher_refines_spec tmpl val hWF
      have hNextNE := insertVariant_nonempty groups fp ct hNonEmpty
      have hPairsStep : evalCompiledKeyGroupsPairs (insertVariantDeclOrder groups fp ct) req =
          match req.getField fp with
          | none => specAcc
          | some val =>
            match extractCapture tmpl val with
            | none => specAcc
            | some captured => upsertHeaderPair specAcc tmpl.headerKey captured := by
        rw [hIns]
        unfold evalCompiledKeyGroupsPairs
        rw [List.filterMap_append]
        have hSingleEval : [ { headerKey := ct.headerKey, variants := [{ fieldPath := fp, compiled := ct }] } ].filterMap
            (fun g => evalCompiledKeyGroup g req) =
            match evalCompiledVariant ⟨fp, ct⟩ req with
            | none => []
            | some v => [(ct.headerKey, v)] := by
          unfold evalCompiledKeyGroup
          cases hVar : evalCompiledVariant ⟨fp, ct⟩ req <;>
            (dsimp only [List.filterMap, List.map]; rw [hVar]; rfl)
        rw [hSingleEval, hVarEq]
        have hBase : List.filterMap (fun g => evalCompiledKeyGroup g req) groups = specAcc := hEqPairs
        rw [hBase]
        cases req.getField fp with
        | none => simp
        | some val =>
          dsimp only
          cases hCap : extractCapture tmpl val with
          | none => simp
          | some captured =>
            dsimp only
            rw [upsertHeaderPair_append_of_not_mem specAcc tmpl.headerKey captured hNotMemAcc]
            rfl
      exact compiled_params_single_variant_aux req rest (insertVariantDeclOrder groups fp ct)
        finalGroups _ hComp hNextNE hSingle hPairsStep

/--
Exact equality (`=`) between compiled and declarative evaluation when every key group has at most
one variant (distinct header keys): in the single-variant regime, static group order coincides with
lazy first-match key order.
-/
theorem compiled_rule_single_variant_eq_spec (rule : RoutingRule) (req : Request)
    (groups : List CompiledKeyGroup)
    (hComp : compileRoutingRuleDeclOrder rule = .ok groups)
    (hSingle : ∀ g ∈ groups, g.variants.length ≤ 1) :
    evalCompiledRuleDeclOrder rule req = evalRoutingRule rule req := by
  unfold evalCompiledRuleDeclOrder evalRoutingRule evalRoutingRulePairs
  unfold compileRoutingRuleDeclOrder at hComp
  have hPairs := compiled_params_single_variant_aux req rule.parameters [] groups []
    hComp (fun _ h => nomatch h) hSingle rfl
  unfold compileRoutingRuleDeclOrder
  rw [hComp, hPairs]
  rfl

/--
Concrete witness of the **static vs lazy key-order divergence**:
when `p1` (key `"x"`) misses and `p2` (key `"y"`) hits before `p3` (key `"x"`) hits,
static key grouping emits `"x=2&y=1"` while declarative `upsertHeaderPair` emits `"y=1&x=2"`,
even though both lookup identically on every key (`compiled_rule_extensional_eq_spec`).
-/
theorem static_vs_lazy_key_order_diverges :
    let rule : RoutingRule := {
      parameters := [
        { field := "a", pathTemplate := "v1/{x=*}" },
        { field := "b", pathTemplate := "v1/{y=*}" },
        { field := "c", pathTemplate := "v2/{x=*}" }
      ]
    }
    let req := Request.ofFlat [("a", "v2/miss"), ("b", "v1/1"), ("c", "v2/2")]
    evalCompiledRuleDeclOrder rule req = .ok "x=2&y=1" ∧
    evalRoutingRule rule req = .ok "y=1&x=2" ∧
    ("x=2&y=1" : String) ≠ "y=1&x=2" := by
  refine ⟨rfl, rfl, ?_⟩
  decide

/-!
## Group 6: RFC 6570 Encoding & Formatting Invariants (+ H4 Falsification)
-/

private theorem char_ofNat_toNat_lt_128 {n : Nat} (h : n < 128) : (Char.ofNat n).toNat = n := by
  have hValid : Nat.isValidChar n := by
    left
    omega
  simp [Char.ofNat, hValid, Char.toNat, Char.ofNatAux, UInt32.toNat_ofNatLT]

private theorem hexDigitUpper_char_valid (n : UInt8) :
    (hexDigitUpper n).toNat < 128 ∧ isUnreservedByte (UInt8.ofNat (hexDigitUpper n).toNat) = true := by
  unfold hexDigitUpper
  dsimp only
  have hv : n.toNat % 16 < 16 := Nat.mod_lt _ (by omega)
  split
  · next hlt =>
    have hLt : 48 + n.toNat % 16 < 128 := by omega
    rw [char_ofNat_toNat_lt_128 hLt]
    refine ⟨hLt, ?_⟩
    unfold isUnreservedByte
    simp only [Bool.or_eq_true, Bool.and_eq_true, decide_eq_true_eq,
      UInt8.le_iff_toNat_le, UInt8.toNat_ofNat', UInt8.toNat_ofNat, Nat.reducePow]
    have hMod : (48 + n.toNat % 16) % 256 = 48 + n.toNat % 16 := Nat.mod_eq_of_lt (by omega)
    rw [hMod]
    left; left; left; left; right
    omega
  · next hge =>
    have hLt : 55 + n.toNat % 16 < 128 := by omega
    rw [char_ofNat_toNat_lt_128 hLt]
    refine ⟨hLt, ?_⟩
    unfold isUnreservedByte
    simp only [Bool.or_eq_true, Bool.and_eq_true, decide_eq_true_eq,
      UInt8.le_iff_toNat_le, UInt8.toNat_ofNat', UInt8.toNat_ofNat, Nat.reducePow]
    have hMod : (55 + n.toNat % 16) % 256 = 55 + n.toNat % 16 := Nat.mod_eq_of_lt (by omega)
    rw [hMod]
    left; left; left; left; left; left
    omega

private theorem percentEncodeByte_char_valid (b : UInt8) (c : Char)
    (hc : c ∈ percentEncodeByte b) :
    (c.toNat < 128 ∧ isUnreservedByte (UInt8.ofNat c.toNat) = true) ∨ c = '%' := by
  unfold percentEncodeByte at hc
  split at hc
  · next hUnres =>
    simp only [List.mem_singleton] at hc
    subst hc
    have hLt : b.toNat < 128 := by
      unfold isUnreservedByte at hUnres
      simp only [Bool.or_eq_true, Bool.and_eq_true, decide_eq_true_eq, beq_iff_eq,
        UInt8.le_iff_toNat_le, ← UInt8.toNat_inj] at hUnres
      have h57 : (57 : UInt8).toNat = 57 := rfl
      have h90 : (90 : UInt8).toNat = 90 := rfl
      have h122 : (122 : UInt8).toNat = 122 := rfl
      have h45 : (45 : UInt8).toNat = 45 := rfl
      have h46 : (46 : UInt8).toNat = 46 := rfl
      have h95 : (95 : UInt8).toNat = 95 := rfl
      have h126 : (126 : UInt8).toNat = 126 := rfl
      omega
    left
    have hNat : (Char.ofNat b.toNat).toNat = b.toNat := char_ofNat_toNat_lt_128 hLt
    rw [hNat]
    refine ⟨hLt, ?_⟩
    simpa [UInt8.ofNat_toNat] using hUnres
  · simp only [List.mem_cons, List.not_mem_nil, or_false] at hc
    rcases hc with rfl | rfl | rfl
    · right; rfl
    · left; exact hexDigitUpper_char_valid (b >>> 4)
    · left; exact hexDigitUpper_char_valid (b &&& 0x0F)

/-- Every character produced by `percentEncode` is either RFC 6570 unreserved or `'%'`. -/
theorem percentEncode_chars_valid (s : String) (c : Char)
    (hc : c ∈ (percentEncode s).toList) :
    (c.toNat < 128 ∧ isUnreservedByte (UInt8.ofNat c.toNat) = true) ∨ c = '%' := by
  unfold percentEncode utf8Bytes at hc
  rw [String.toList_ofList, List.mem_flatMap] at hc
  rcases hc with ⟨b, _, hb⟩
  exact percentEncodeByte_char_valid b c hb

private theorem utf8EncodeChar_ne_nil (c : Char) : String.utf8EncodeChar c ≠ [] := by
  dsimp [String.utf8EncodeChar]
  split <;> try simp
  split <;> try simp
  split <;> simp

private theorem percentEncodeByte_ne_nil (b : UInt8) : percentEncodeByte b ≠ [] := by
  unfold percentEncodeByte
  split <;> simp

/-- `percentEncode s` is empty if and only if `s` is empty. -/
theorem percentEncode_empty_iff (s : String) : percentEncode s = "" ↔ s = "" := by
  constructor
  · intro h
    have hList : (percentEncode s).toList = [] := by simp [h]
    unfold percentEncode utf8Bytes at hList
    rw [String.toList_ofList] at hList
    rw [String.ext_iff]
    cases hs : s.toList with
    | nil => rfl
    | cons c rest =>
      rw [hs, List.flatMap_cons, List.flatMap_append] at hList
      have hBytesNe : String.utf8EncodeChar c ≠ [] := utf8EncodeChar_ne_nil c
      cases hUtf : String.utf8EncodeChar c with
      | nil => exact absurd hUtf hBytesNe
      | cons b bs =>
        rw [hUtf, List.flatMap_cons] at hList
        have hEncNe := percentEncodeByte_ne_nil b
        simp only [List.append_eq_nil_iff] at hList
        exact absurd hList.1.1 hEncNe
  · intro h
    subst h
    rfl

private theorem formatKeyValue_ne_nil (k v : String) : (formatKeyValue k v).toList ≠ [] := by
  unfold formatKeyValue
  simp [String.toList_append]

private theorem joinStringsWith_amp_nonempty (s : String) (ss : List String)
    (hs : s.toList ≠ []) :
    (joinStringsWith "&" (s :: ss)).toList ≠ [] ∧
    ((joinStringsWith "&" (s :: ss)).toList.head? = s.toList.head?) := by
  cases ss with
  | nil => simp [joinStringsWith, hs]
  | cons t ts =>
    simp only [joinStringsWith, String.toList_append]
    constructor
    · intro hNil
      simp only [List.append_eq_nil_iff] at hNil
      exact hs hNil.1.1
    · cases hsList : s.toList with
      | nil => exact absurd hsList hs
      | cons c cs => simp

/-- Formatted `x-goog-request-params` header is empty iff the extracted pair list is empty. -/
theorem formatHeaderPairs_empty_iff (pairs : List (String × String)) :
    formatHeaderPairs pairs = "" ↔ pairs = [] := by
  constructor
  · intro h
    cases pairs with
    | nil => rfl
    | cons kv rest =>
      unfold formatHeaderPairs at h
      simp only [List.map_cons] at h
      have hNe := (joinStringsWith_amp_nonempty (formatKeyValue kv.1 kv.2)
        (rest.map fun (k, v) => formatKeyValue k v)
        (formatKeyValue_ne_nil kv.1 kv.2)).1
      exfalso
      apply hNe
      simp [h]
  · intro h
    subst h
    rfl

/-- Formatted `x-goog-request-params` headers never begin with `'&'`. -/
theorem formatHeaderPairs_no_leading_amp (pairs : List (String × String)) :
    (formatHeaderPairs pairs).toList.head? ≠ some '&' := by
  cases pairs with
  | nil => simp [formatHeaderPairs, joinStringsWith]
  | cons kv rest =>
    unfold formatHeaderPairs
    simp only [List.map_cons]
    have hHead := (joinStringsWith_amp_nonempty (formatKeyValue kv.1 kv.2)
      (rest.map fun (k, v) => formatKeyValue k v)
      (formatKeyValue_ne_nil kv.1 kv.2)).2
    rw [hHead]
    unfold formatKeyValue
    simp only [String.toList_append]
    cases hKey : (percentEncode kv.1).toList with
    | nil =>
      intro hContra
      simp only [List.nil_append, String.toList] at hContra
      injection hContra with hEq
      have hBad : ('=' : Char) ≠ '&' := by decide
      exact hBad hEq
    | cons c cs =>
      have hMem : c ∈ (percentEncode kv.1).toList := by simp [hKey]
      have hValid := percentEncode_chars_valid kv.1 c hMem
      simp only [List.append_assoc, List.cons_append, List.head?_cons]
      intro hAmp
      injection hAmp with hAmpEq
      subst hAmpEq
      rcases hValid with ⟨_, hUnres⟩ | hPct
      · revert hUnres; decide
      · cases hPct

/-- Concrete falsification of `librarian`'s H4 implicit HTTP fold (`transport.rs.gotmpl`). -/
theorem h4_buggy_rust_http_fold_falsified :
    let httpRule : HttpRule := {
      bindings := [[{ fieldPath := ["parent"] }, { fieldPath := ["empty_field"] }]]
    }
    let req := Request.ofFlat [("parent", "projects/my project"), ("empty_field", "")]
    formatHttpBuggyRust httpRule req = "&parent=projects/my project&empty_field=" ∧
    formatHeaderPairs (evalHttpRulePairs httpRule req) = "parent=projects%2Fmy%20project" ∧
    formatHttpBuggyRust httpRule req ≠ formatHeaderPairs (evalHttpRulePairs httpRule req) := by
  refine ⟨rfl, rfl, ?_⟩
  have h1 : formatHttpBuggyRust
      { bindings := [[{ fieldPath := ["parent"] }, { fieldPath := ["empty_field"] }]] }
      (Request.ofFlat [("parent", "projects/my project"), ("empty_field", "")]) =
      "&parent=projects/my project&empty_field=" := rfl
  have h2 : formatHeaderPairs (evalHttpRulePairs
      { bindings := [[{ fieldPath := ["parent"] }, { fieldPath := ["empty_field"] }]] }
      (Request.ofFlat [("parent", "projects/my project"), ("empty_field", "")])) =
      "parent=projects%2Fmy%20project" := rfl
  rw [h1, h2]
  decide

end Aip4222
