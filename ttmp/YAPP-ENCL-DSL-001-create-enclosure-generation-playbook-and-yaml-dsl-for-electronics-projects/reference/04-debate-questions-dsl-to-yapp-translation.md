---
Title: Debate Questions — DSL to YAPP Translation
Ticket: YAPP-ENCL-DSL-001
Status: active
Topics:
    - yapp
    - openscad
    - dsl
DocType: reference
Intent: long-term
Owners:
    - manuel
RelatedFiles:
    - Path: reference/03-debate-format-and-candidates-dsl-to-yapp-translation.md
      Note: Candidate profiles
ExternalSources: []
Summary: Four debate questions for analyzing DSL-to-YAPP translation correctness
LastUpdated: 2025-11-09
---

# Debate Questions — DSL to YAPP Translation

## Question Progression

### Round 1: Completeness — Can the DSL Express Real YAPP Boxes?

**Primary candidates**: Jamie (Feature Engineer), YAPP_Demo_RealBox_v31.scad, The Template

**Question**: Does the DSL language reference cover all essential YAPP features needed to generate production-quality enclosures? What's missing?

**Key sub-questions**:
- Can we express pcbStands, connectors, cutouts (all 6 faces), snapJoins, lightTubes?
- What about advanced features: boxMounts, pushButtons, labels, ridge extensions?
- Are coordinate systems (yappCoordPCB, yappCoordBox, yappCoordBoxInside) mappable?

**Research needed**:
- Compare DSL language reference against YAPP_Template_v3.scad feature arrays
- Analyze YAPP_Demo_RealBox_v31.scad for features used
- Check DSL for gaps

---

### Round 2: Correctness — Are the Mappings Semantically Accurate?

**Primary candidates**: Dr. Sarah (Architect), YAPPgenerator_v3.scad, The Historian

**Question**: Does the proposed DSL-to-YAPP translation preserve semantic meaning? Are there impedance mismatches or coordinate system confusion?

**Key sub-questions**:
- Do DSL coordinate origins map correctly to YAPP's three systems?
- Are parameter orders and optional values handled correctly?
- Does the DSL expose or hide YAPP's complexity appropriately?

**Research needed**:
- Review analysis/02-dsl-to-yapp-translation-analysis.md mappings
- Check YAPP docs DB for coordinate system explanations
- Trace DSL `coordinates.origin` to YAPP flags (yappCoordPCB, etc.)

---

### Round 3: Practicality — Can Users Actually Use This?

**Primary candidates**: Alex (Pragmatist), The New User, Jamie (Feature Engineer)

**Question**: Is the DSL practical for real users? Does it reduce or increase cognitive load compared to hand-writing SCAD?

**Key sub-questions**:
- Is YAML with expressions easier than OpenSCAD arrays?
- What's the error experience when users make mistakes?
- Can users migrate existing SCAD to DSL (or vice versa)?

**Research needed**:
- Compare DSL YAML example to equivalent SCAD
- Count lines, measure verbosity
- Identify common user errors (wrong coordinate system, missing required fields)

---

### Round 4: Risks and Evolution — What Could Go Wrong?

**Primary candidates**: The Historian, Dr. Sarah (Architect), YAPPgenerator_v3.scad

**Question**: What are the risks of this approach, and how do we mitigate them? How does this handle YAPP version evolution?

**Key sub-questions**:
- What happens when YAPP v3.4 adds new features or changes APIs?
- Are we creating a maintenance burden (DSL + generator + YAPP)?
- What if users need features the DSL doesn't support?

**Research needed**:
- Query yapp_analysis.db for API changes and breaking changes
- Check version_info table for version drift patterns
- Identify escape hatches (can users drop to raw SCAD?)

---

## Candidate Mapping

| Round | Primary Candidates | Wildcard Interrupters |
|-------|-------------------|----------------------|
| 1 | Jamie, RealBox, Template | New User, Generator |
| 2 | Sarah, Generator, Historian | Template, Alex |
| 3 | Alex, New User, Jamie | RealBox, Sarah |
| 4 | Historian, Sarah, Generator | Alex, Template |

## Expected Outcomes

After 4 rounds:
- **Round 1**: List of missing DSL features or confirmed completeness
- **Round 2**: Validated or corrected semantic mappings
- **Round 3**: User experience assessment and friction points
- **Round 4**: Risk register and mitigation strategies
