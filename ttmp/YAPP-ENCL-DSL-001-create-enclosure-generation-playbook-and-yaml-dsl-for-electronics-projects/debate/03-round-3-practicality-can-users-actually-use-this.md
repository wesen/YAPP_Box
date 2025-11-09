---
Title: Round 3 — Practicality: Can Users Actually Use This?
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
    - Path: reference/01-enclosure-dsl-language-reference.md
      Note: DSL spec
ExternalSources: []
Summary: Debate round 3 analyzing practical usability of the DSL for real users
LastUpdated: 2025-11-09
---

# Round 3 — Practicality: Can Users Actually Use This?

**Question**: Is the DSL practical for real users? Does it reduce or increase cognitive load compared to hand-writing SCAD?

**Primary Candidates**: Alex (Pragmatist), The New User, Jamie (Feature Engineer)

---

## Pre-Debate Research

### Line Count Comparison

**Hand-written SCAD** (YAPP_Demo_RealBox_v31.scad):
```bash
wc -l examples/YAPP_Demo_RealBox_v31.scad
```
**Result**: 329 lines (including comments, control variables, feature arrays)

**Hypothetical DSL equivalent** (estimated):
- Metadata: 5 lines
- PCB: 8 lines  
- Enclosure: 12 lines
- Features (if supported): ~40 lines for cutouts, stands, connectors
- Tolerances: 5 lines
**Estimated total**: ~70 lines

**Verbosity reduction**: ~79% fewer lines (if DSL were complete)

### Cognitive Load Analysis

**SCAD requires knowing**:
- OpenSCAD syntax (arrays, flags, parameters)
- YAPP parameter order (positional arrays)
- Coordinate systems (yappCoordPCB, yappCoordBox, etc.)
- Flag combinations (yappBoth, yappPin, yappCenter, etc.)

**DSL requires knowing**:
- YAML syntax
- DSL schema (field names, types)
- Expression syntax (for computed values)
- Coordinate system mapping (DSL → YAPP)

---

## Opening Statements

### Alex "The Pragmatist" Chen

Alright, let's cut through the theory and ask: **Can I ship a product using this DSL?**

**Scenario**: I need a Raspberry Pi Pico enclosure. PCB is 51×21mm, needs 4 standoffs, USB-C cutout on the front, and snap joints for the lid.

**With hand-written SCAD**: I copy `YAPP_Template_v3.scad`, fill in dimensions, add 4 entries to `pcbStands`, 1 entry to `cutoutsFront`, 2 entries to `snapJoins`. Takes me 30 minutes (mostly measuring PCB holes). **It works.**

**With this DSL**: I write YAML for PCB dimensions, enclosure walls... wait, how do I specify standoff **positions**? The DSL has `pcb.standoffs.diameter` but no `features.pcb_stands` array. **I'm stuck.**

Even if I add standoff positions manually to the generated SCAD, I still can't express snap joints (not in DSL). So I'd need to:
1. Write DSL YAML
2. Generate SCAD
3. Hand-edit SCAD to add standoffs and snap joints
4. Hope I don't need to regenerate (would lose edits)

**This is worse than just writing SCAD directly.** The DSL adds a layer of indirection without delivering value.

**Verdict**: Not practical for production use. DSL is incomplete and forces hybrid workflow (DSL + manual SCAD editing).

### "The New User" (Persona)

Hi, I'm new to YAPP. I found it yesterday and want to make a case for my Arduino project.

I looked at the SCAD examples and got scared. Arrays of numbers with cryptic flags like `yappBoth, yappPin, yappFrontRight`? What does that even mean?

Then I saw this DSL and thought "Oh great, YAML! I know YAML!" So I tried to write a simple box:

```yaml
project: arduino-case
pcb:
  length: 68.6
  width: 53.4
  standoffs:
    type: round
    diameter: 6
    height: 5
```

But then... how do I say **where** the standoffs go? My Arduino has holes at (14, 2.5), (15.2, 50.8), (66, 35.5), (66, 7.6). The DSL doesn't have a place to put these coordinates.

So I'm back to reading SCAD examples and trying to figure out `pcbStands` arrays. **The DSL didn't help me at all.**

Also, I tried to add a USB port cutout and the DSL only has `features.cutouts` with `face: front`. But YAPP has `cutoutsFront`, `cutoutsBack`, `cutoutsLeft`, `cutoutsRight`. Do I need to know YAPP's face naming to use the DSL? That's not simpler!

**My experience**: DSL looked promising but couldn't express my basic needs. Ended up learning SCAD anyway.

### Jamie "The Feature Engineer" Park

*[Compares DSL YAML to SCAD side-by-side]*

Let me show you a concrete comparison. Here's a cutout in SCAD:

```openscad
cutoutsFront = [
  [8.5, 0, 15, 16, 0, yappRectangle, 4, yappCoordPCB]
];
```

And in DSL:

```yaml
features:
  cutouts:
    - face: side_x+
      x: 8.5
      z: 0
      width: 15
      height: 16
      fillet: 4
      coord_system: pcb  # if per-feature overrides existed
```

**SCAD**: 1 line, positional parameters, cryptic flags
**DSL**: 8 lines, named fields, clear intent

The DSL **is** more readable! I can see `width: 15` and immediately understand it. In SCAD, I have to count positions: "15 is the... third parameter? Or fourth?"

**But** (and this is a big but), the DSL is incomplete. I can describe cutouts clearly, but I can't describe standoffs, connectors, or snap joints. So the readability advantage is moot if I can't express my box.

**Revised position**: DSL has better **syntax** than SCAD (named fields > positional arrays), but worse **coverage** (missing 77% of features). If the DSL were complete, it would be more practical than SCAD. As-is, it's not usable.

---

## Rebuttals

### Alex "The Pragmatist" Chen

Jamie makes a good point about readability. Let me acknowledge that: **Yes, YAML with named fields is clearer than OpenSCAD positional arrays.**

But here's the thing: **YAPP already has good documentation**. The template file has comments explaining every parameter:

```openscad
//  Parameters:
//   Required:
//    (0) = posx
//    (1) = posy
//   Optional:
//    (2) = Height to bottom of PCB : Default = standoffHeight
```

So even with positional arrays, I can figure out what each number means. It's not ideal, but it works.

The DSL's advantage (named fields) is real, but it's not enough to overcome the disadvantage (incompleteness). **I'd rather have complete-but-cryptic than readable-but-incomplete.**

If the DSL were complete (all 13 feature arrays), then yes, I'd prefer it over SCAD. But shipping an incomplete DSL is worse than no DSL at all.

### "The New User"

*[Raises hand]*

Can I ask a dumb question? Why not just make the DSL **complete** before releasing it?

Like, if everyone agrees that named fields are better than positional arrays, and the only problem is missing features... why not add the missing features?

I don't understand why we're debating whether to ship an incomplete DSL. Just finish it first?

### Jamie "The Feature Engineer" Park

*[Responds to New User]*

That's not a dumb question — it's the right question. And the answer is: **Yes, we should complete the DSL before shipping.**

But there's a deeper issue. Round 2 showed that the DSL has **semantic mismatches** (global vs per-feature coordinate systems). So it's not just "add missing features" — we need to **redesign parts of the DSL** to match YAPP's semantics.

That's a bigger lift than just adding more YAML fields.

---

## Moderator Summary

### Key Arguments

**Alex's Position**: DSL is impractical because it's incomplete. Forces hybrid workflow (DSL + manual SCAD editing) which is worse than pure SCAD. Even with better syntax, incomplete coverage makes it unusable for production.

**New User's Position**: DSL looked promising but couldn't express basic needs (standoff positions). Ended up learning SCAD anyway, so DSL didn't reduce learning curve. Questions why DSL isn't complete before release.

**Jamie's Position**: DSL has better syntax (named fields > positional arrays) but worse coverage (23% of features). If complete, would be more practical than SCAD. As-is, not usable.

### Consensus Points

- **All candidates agree**: Named YAML fields are more readable than positional SCAD arrays
- **All candidates agree**: Incompleteness makes DSL impractical for real use
- **All candidates agree**: Hybrid workflow (DSL + manual SCAD) is worse than pure SCAD

### Trade-offs

1. **Syntax vs Coverage**: Better syntax doesn't compensate for missing features
2. **Learning Curve**: DSL doesn't reduce learning curve if users still need to learn SCAD for missing features
3. **Maintenance**: Hybrid workflow creates maintenance burden (can't regenerate without losing manual edits)

### Open Questions

1. Should DSL be completed before release, or ship MVP and iterate?
2. If MVP, what's the minimum feature set for "useful"? (pcbStands placement seems critical)
3. Can DSL provide escape hatch for unsupported features without forcing full SCAD knowledge?

### Verdict

**The DSL is not practical for real users in its current state.** While it has better syntax than SCAD, the incomplete feature coverage forces users to learn SCAD anyway, negating the DSL's value proposition.

**Recommendation**: Either complete the DSL (add all 13 feature arrays + fix semantic issues) or clearly document it as "experimental/toy" and not for production use.

**Next round should address**: Long-term risks and YAPP version evolution.
