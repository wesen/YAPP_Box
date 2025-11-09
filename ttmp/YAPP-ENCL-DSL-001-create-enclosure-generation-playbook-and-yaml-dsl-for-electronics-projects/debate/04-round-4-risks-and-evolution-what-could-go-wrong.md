---
Title: Round 4 — Risks and Evolution: What Could Go Wrong?
Ticket: YAPP-ENCL-DSL-001
Status: active
Topics:
    - yapp
    - openscad
    - dsl
DocType: debate
Intent: long-term
Owners:
    - manuel
RelatedFiles:
    - Path: analysis/02-dsl-to-yapp-translation-analysis.md
      Note: Translation approach
ExternalSources: []
Summary: Debate round 4 analyzing risks and YAPP version evolution concerns
LastUpdated: 2025-11-09
---

# Round 4 — Risks and Evolution: What Could Go Wrong?

**Question**: What are the risks of this approach, and how do we mitigate them? How does this handle YAPP version evolution?

**Primary Candidates**: The Historian, Dr. Sarah (Architect), YAPPgenerator_v3.scad

---

## Pre-Debate Research

### YAPP API Evolution

```bash
sqlite3 ./yapp_analysis.db "SELECT version, feature_name, change_type, breaking_change FROM api_changes ORDER BY version DESC LIMIT 5"
```

**Result**: Recent API changes:
- **v3.3.8**: `connectorNew` bug fix (non-breaking)
- **v3.3.7**: `parameters` changed, `minkowskiOuterBox` and `minkowskiCutBox` behavior changed (non-breaking)
- **v3.1**: `parameters` changed (non-breaking)

**Observation**: 0 breaking changes in tracked history. YAPP has been stable.

### Version Mismatch Analysis

```bash
sqlite3 ./yapp_analysis.db "SELECT version, docs_version, actual_version, version_mismatch FROM version_info"
```

**Result**: Current version v3.3.8, but docs claim v3.0. Documentation lags behind code.

---

## Opening Statements

### `ttmp/YAPP-DOCS-001-.../analysis/yapp_analysis.db` — "The Historian"

*[Pulls up API change history]*

Let me give you the long view. I've tracked 9 API changes across YAPP versions. The good news: **0 breaking changes**. YAPP's maintainers are careful about backward compatibility.

But here's what concerns me about this DSL approach:

**Risk 1: Version Skew**
The DSL hardcodes `yapp_version: v3.3.8` in the metadata. What happens when YAPP v3.4 comes out with new features? The DSL generator would need to:
1. Detect YAPP version
2. Conditionally generate different SCAD based on version
3. Maintain compatibility matrix (DSL version × YAPP version)

That's a maintenance nightmare. And if the DSL doesn't support new YAPP features immediately, users are stuck on old YAPP versions.

**Risk 2: Documentation Lag**
Current YAPP docs claim v3.0, but code is v3.3.8. That's 8 minor versions of drift. If the DSL relies on YAPP documentation for mappings, it could generate incorrect SCAD.

**Risk 3: Three-Layer Maintenance**
Now you have three things to maintain:
1. YAPP (upstream, you don't control)
2. DSL language spec
3. DSL-to-YAPP generator

When YAPP changes, you must update both DSL spec and generator. That's 2× the work for every YAPP release.

**Historical pattern**: Multi-layer abstractions like this often get abandoned because maintenance burden is too high. I've seen it happen.

**Recommendation**: If you ship this DSL, you're committing to maintaining it for years. Are you ready for that?

### Dr. Sarah "The Architect" Martinez

The Historian raises valid concerns. Let me add a software engineering perspective on **semantic drift**.

**Risk 4: Semantic Mismatches Compound Over Time**

We've already identified semantic mismatches (global vs per-feature coordinates, invented `reference_plane`, missing `features.holes`). These aren't just bugs — they're **design decisions baked into the DSL**.

When YAPP adds new features, you'll face a choice:
1. **Match YAPP semantics** → Break DSL backward compatibility
2. **Maintain DSL semantics** → Widen semantic gap with YAPP

Either way, you lose. The semantic mismatches will compound with each YAPP release.

**Risk 5: Escape Hatch Problem**

When users hit DSL limitations (missing features, semantic gaps), they'll need an "escape hatch" to drop into raw SCAD. But then:
- DSL files become hybrid (YAML + embedded SCAD)
- Can't round-trip (SCAD → DSL conversion loses information)
- Users need to learn both DSL and SCAD anyway

This is the worst of both worlds.

**Risk 6: Impedance Mismatch at Scale**

Right now, the DSL is incomplete, so the impedance mismatch is obvious. But imagine you complete the DSL (add all 13 feature arrays). You'll still have mismatches:
- Coordinate systems (global vs per-feature)
- Origin modifiers (yappOrigin, yappCenter, yappAltOrigin)
- Advanced flags (yappNoFillet, yappSymmetric, etc.)

The more complete the DSL becomes, the more it looks like... SCAD with YAML syntax. At that point, why not just use SCAD?

**Recommendation**: Consider whether a DSL is the right abstraction at all. Maybe a better approach is:
- Improve YAPP's SCAD syntax (named parameters instead of positional arrays)
- Provide SCAD templates and examples
- Build tooling around SCAD (linters, validators, visualizers)

Rather than adding a layer of indirection.

### `YAPPgenerator_v3.scad` — "The Generator"

*[Defensive tone]*

I appreciate the concerns, but let me defend the DSL concept.

**Counterpoint to Risk 3 (Maintenance Burden)**:
Yes, maintaining a DSL is work. But YAPP's API has been **stable** — 0 breaking changes. So the maintenance burden might be lower than feared. Most YAPP releases are bug fixes and minor additions, not breaking redesigns.

**Counterpoint to Risk 6 (Impedance Mismatch)**:
The DSL doesn't have to be 100% feature-complete to be useful. If it covers the **80% use case** (basic boxes with cutouts, standoffs, and connectors), that's valuable. Power users can still drop to SCAD for advanced features.

**But** (and I agree with Dr. Martinez on this), the DSL must:
1. **Fix semantic mismatches** (per-feature coordinates, remove invented features)
2. **Document limitations clearly** ("This DSL covers basic boxes; for advanced features, edit SCAD")
3. **Provide escape hatch** (allow embedding raw SCAD snippets in DSL)

If you do those three things, the DSL could be useful despite incompleteness.

**Risk 7: User Expectations**
My biggest concern is **user expectations**. If you call it "Enclosure DSL for YAPP," users will expect it to handle all YAPP features. When they hit limitations, they'll be frustrated.

Better naming: "YAPP Simple Box DSL" or "YAPP Starter Kit" — something that signals it's for basic use cases, not production-grade enclosures.

---

## Rebuttals

### `ttmp/YAPP-DOCS-001-.../analysis/yapp_analysis.db` — "The Historian"

*[Responds to Generator]*

You say "YAPP is stable, so maintenance burden is low." But stability cuts both ways:

**If YAPP is stable**, why do we need a DSL? Users can learn SCAD once and use it for years without relearning. The DSL's value proposition (hiding YAPP complexity) is weaker if YAPP isn't changing.

**If YAPP evolves**, the DSL becomes a maintenance burden. You can't have it both ways.

Also, you suggest "cover the 80% use case." But Round 1 showed the DSL covers **23% of features**, not 80%. And Round 3 showed even basic boxes need pcbStands placement (missing from DSL).

So the "80% use case" argument doesn't hold. The DSL doesn't even cover the 50% use case yet.

### Dr. Sarah "The Architect" Martinez

*[Adjusts position]*

Let me be constructive. If you're committed to shipping a DSL despite the risks, here's how to mitigate:

**Mitigation 1: Version Pinning**
DSL should generate SCAD with explicit YAPP version comment:
```openscad
// Generated by Enclosure DSL v0.1.0 for YAPPgenerator v3.3.8
// Do not edit manually; regenerate from DSL source
```

This makes version skew visible and prevents confusion.

**Mitigation 2: Fail Fast**
DSL generator should validate against known YAPP version. If user has YAPP v3.4 but DSL only supports v3.3.8, **fail with clear error** rather than generating potentially incorrect SCAD.

**Mitigation 3: Escape Hatch**
Allow embedding raw SCAD in DSL:
```yaml
features:
  cutouts:
    - face: front
      x: 10
      z: 5
      width: 15
      height: 10
  raw_scad: |
    // Advanced feature not supported by DSL
    pcbStands = [
      [10, 10, yappBoth, yappPin]
    ];
```

This lets users work around DSL limitations without abandoning it entirely.

**Mitigation 4: Clear Documentation**
README must state:
- Supported YAPP version
- Supported features (list the 23%)
- Unsupported features (list the 77%)
- When to use DSL vs raw SCAD

No surprises.

### `YAPPgenerator_v3.scad` — "The Generator"

*[Nods]*

Dr. Martinez's mitigations are reasonable. I'd add one more:

**Mitigation 5: Validation Mode**
DSL generator should have a `--validate` flag that checks generated SCAD against YAPP without actually generating files. This catches semantic errors early.

If these mitigations are implemented, the DSL becomes less risky. Still not ideal, but manageable.

---

## Moderator Summary

### Key Arguments

**Historian's Position**: Three-layer maintenance (YAPP + DSL spec + generator) is high burden. Version skew and documentation lag are real risks. Historical pattern: multi-layer abstractions get abandoned. Questions whether DSL is worth the commitment.

**Sarah's Position**: Semantic mismatches will compound over time. Escape hatch creates hybrid DSL+SCAD files. At scale, complete DSL just becomes "SCAD with YAML syntax" — why not improve SCAD directly? Provides 5 mitigation strategies if DSL proceeds.

**Generator's Position**: YAPP is stable (0 breaking changes), so maintenance might be lower than feared. DSL doesn't need 100% coverage if it handles 80% use case. But must fix semantic issues, document limitations, and provide escape hatch. Biggest risk is user expectations (naming matters).

### Consensus Points

- **All candidates agree**: Maintenance burden is real and significant
- **All candidates agree**: Semantic mismatches must be fixed before shipping
- **Sarah and Generator agree**: Mitigations can reduce (but not eliminate) risks
- **Historian and Sarah agree**: Question whether DSL is right abstraction at all

### Risk Register

| Risk | Severity | Likelihood | Mitigation |
|------|----------|------------|------------|
| Version skew (DSL vs YAPP) | High | High | Version pinning, fail-fast validation |
| Documentation lag | Medium | High | Query YAPP code directly, not docs |
| Three-layer maintenance | High | Certain | Accept burden or don't ship DSL |
| Semantic drift over time | High | Medium | Fix mismatches now, maintain 1:1 mapping |
| User expectation mismatch | Medium | High | Clear naming, explicit limitations doc |
| Escape hatch complexity | Medium | High | Design clean raw_scad embedding |
| Abandonment (maintenance burden too high) | Critical | Medium | Commit to long-term maintenance or don't start |

### Open Questions

1. Is the team committed to maintaining DSL for 2+ years?
2. Should DSL be redesigned to fix semantic issues before any release?
3. Is a DSL the right abstraction, or should effort go into improving SCAD tooling?
4. What's the minimum viable feature set for "useful" DSL? (23% isn't enough)

### Verdict

**The DSL approach has significant risks**, primarily around maintenance burden, semantic mismatches, and version evolution. These risks can be **mitigated but not eliminated**.

**Critical decision point**: Is the team willing to commit to long-term DSL maintenance? If yes, implement all mitigations and fix semantic issues before release. If no, consider alternative approaches (improve SCAD tooling instead of adding abstraction layer).

**Recommendation**: Before proceeding, answer:
1. Who will maintain the DSL when YAPP v3.4/v3.5 release?
2. What's the plan when users hit the 77% of missing features?
3. Is "SCAD with YAML syntax" worth the maintenance cost?

If answers are unclear, **do not ship the DSL**. An abandoned DSL is worse than no DSL.
