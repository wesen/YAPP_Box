---
Title: DSL Feature Gaps Analysis - Missing YAPP Arrays
Ticket: YAPP-DSL-GAPS-001
Status: active
Topics:
    - yapp
    - dsl
    - analysis
    - features
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: YAPP_Template_v3.scad
      Note: canonical YAPP feature definitions
    - Path: examples/YAPP_Demo_buttons2_v31.scad
      Note: user's target example
    - Path: examples/YAPP_Demo_buttons_v30.scad
      Note: reference for parity
    - Path: examples/test-push-buttons.yaml
      Note: current DSL test case
    - Path: examples/yapp-demo-buttons-v30.yaml
      Note: updated with vars and cutouts base/front/back
    - Path: pkg/resolver/resolver.go
      Note: treats enum arrays as literals (corners)
    - Path: pkg/yappgen/features.go
      Note: builder wiring for flags
    - Path: pkg/yappgen/modules/boxmounts/module.go
      Note: builder
    - Path: pkg/yappgen/modules/boxmounts/module_test.go
      Note: tests
    - Path: pkg/yappgen/modules/boxmounts/registry.go
      Note: module registration
    - Path: pkg/yappgen/modules/boxmounts/schema.yaml
      Note: box_mounts schema
    - Path: pkg/yappgen/modules/boxmounts/schema_gen.go
      Note: schemagen struct
    - Path: pkg/yappgen/modules/connectors/module.go
      Note: flag builder
    - Path: pkg/yappgen/modules/connectors/module_test.go
      Note: flag tests
    - Path: pkg/yappgen/modules/connectors/schema.yaml
      Note: schema flags
    - Path: pkg/yappgen/modules/cutouts/schema.yaml
      Note: cutout shapes; polygon/mask TODO
    - Path: pkg/yappgen/modules/pcbstands/module.go
      Note: corners[] flags
    - Path: pkg/yappgen/modules/pcbstands/module_test.go
      Note: flag tests
    - Path: pkg/yappgen/modules/pcbstands/schema.yaml
      Note: corners[] schema
    - Path: pkg/yappgen/modules/snapjoins/module.go
      Note: flag builder
    - Path: pkg/yappgen/modules/snapjoins/module_test.go
      Note: flag tests
    - Path: pkg/yappgen/modules/snapjoins/schema.yaml
      Note: snap flags schema
    - Path: pkg/yappgen/modules_gen.go
      Note: includes box_mounts
    - Path: ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/design-doc/02-pcb-stand-flag-support.md
      Note: design
    - Path: ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/design-doc/03-connector-flag-support.md
      Note: design
    - Path: ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/design-doc/04-boxmounts-module-support.md
      Note: design
ExternalSources: []
Summary: Analysis of missing YAPP features preventing DSL parity with SCAD examples - identifies 8 missing arrays and critical flag gaps
LastUpdated: 2025-11-15T22:14:00.857144922-05:00
---









# DSL Feature Gaps Analysis - Missing YAPP Arrays

## Overview

This ticket analyzes missing YAPP features that prevent the DSL from achieving parity with real SCAD examples. A user attempted to replicate `YAPP_Demo_buttons2_v31.scad` using the DSL and encountered visual discrepancies: extra PCB pillars, external mounting plates, and potential hinge-like structures.

**Root cause:** The DSL currently implements 5 of 13+ YAPP feature arrays (38% coverage) and is missing critical flags on implemented features.

**Goal:** Document all gaps, prioritize fixes, and create actionable tasks for achieving 80%+ YAPP feature coverage.

## For the Intern: Quick Start Guide

### 1. Understand the Context

**Read these documents in order:**

1. **Start here:** `reference/01-yapp-feature-array-inventory.md`
   - Complete catalog of YAPP features
   - Shows what's implemented ✅ vs missing ❌
   - Parameter specifications for each array

2. **User's problem:** `analysis/01-user-case-study-buttons2-demo-discrepancies.md`
   - Specific issues user encountered
   - Root cause analysis
   - Hypothesis for each discrepancy

3. **Implementation plan:** `design-doc/01-priority-feature-roadmap.md`
   - Prioritized feature list (P1/P2/P3)
   - Effort estimates
   - Implementation strategy
   - Timeline: 12-16 weeks for full coverage

### 2. Learn the YAPP Library

**Official YAPP documentation:**
- **GitBook:** https://mrwheel-docs.gitbook.io/yappgenerator_en/ (official docs)
- **Local docs study guide:** Run `go run ./cmd/yappctl help yapp-docs-analysis-playbook`
  - Explains how to query local docs database
  - Points to scraped GitBook content
  - Shows how to validate examples

**YAPP source files (in repo root):**
- `YAPP_Template_v3.scad` - Canonical parameter definitions for all features
- `YAPPgenerator_v3.scad` - Generator implementation (complex, 8000+ lines)
- `CHANGELOG.md` - Version history and API changes

**Example files (`examples/` directory):**
- `YAPP_Demo_buttons2_v31.scad` - User's target (push buttons demo)
- `YAPP_Demo_RealBox_v31.scad` - Production example (9 features)
- `YAPP_Demo_boxMounts_v30.scad` - External mounting tabs
- `YAPP_Demo_lightTubes_v30.scad` - LED indicators
- `YAPP_Reference_Shapes_v30.scad` - Shape examples
- `YAPP_Reference_Masks_v30.scad` - Ventilation masks

### 3. Understand the Module System

**Read the implementation guide:**
```bash
cd ttmp/2025/11/15/YAPP-MODULE-SYSTEM-001-*
cat playbook/module-system-implementation-guide.md
```

**Key concepts:**
- Schema-driven validation (YAML schemas define fields)
- Code generation (schemagen tool)
- Auto-registration (modules_gen.go)
- Two-phase validation (structure + constraints)

**How to add a new module:**
```bash
go run ./cmd/yappctl help yapp-module-authoring-guide
```

**Existing module examples:**
- `pkg/yappgen/modules/pushbuttons/` - Complex (nested objects, shape flags)
- `pkg/yappgen/modules/pcbstands/` - Simple (flat structure)
- `pkg/yappgen/modules/snapjoins/` - Enum flags (side parameter)
- `pkg/yappgen/modules/cutouts/` - Face distribution logic

### 4. Work Through Tasks

**Task list:** See `tasks.md` in this directory

**Priority order:**
1. Research corner placement logic in YAPPgenerator_v3.scad
2. Get user's actual DSL YAML for comparison
3. Implement corner placement flags (highest ROI)
4. Implement shell part flags
5. Implement boxMounts module
6. Continue with Priority 2 features

**For each task:**
1. Read related reference docs
2. Study YAPP examples using that feature
3. Check YAPP_Template_v3.scad parameter comments
4. Write schema.yaml following existing patterns
5. Validate: `go run ./cmd/schemagen validate <schema>`
6. Generate: `go run ./cmd/schemagen discover`
7. Write builder (module.go) and registry (registry.go)
8. Test: `go test ./...`
9. Document via docmgr (changelog, relate, tasks check)

### 5. Testing Your Changes

**End-to-end test workflow:**
```bash
# 1. Create test YAML with your new feature
cat > test-feature.yaml << 'EOF'
features:
  box_mounts:
    - pos: 10
      screw_diameter: 3
      width: 6
      height: 3
      sides: [left, right]
EOF

# 2. Resolve
go run ./cmd/yappctl resolve -i test-feature.yaml -o /tmp/resolved.yaml

# 3. Generate SCAD
go run ./cmd/yappctl generate -i /tmp/resolved.yaml -o /tmp/output.scad

# 4. Check output
grep boxMounts /tmp/output.scad

# 5. Render (optional, requires OpenSCAD)
openscad -o /tmp/output.stl /tmp/output.scad
```

**Verify help pages:**
```bash
go run ./cmd/yappctl help module-box_mounts
```

### 6. Documentation Requirements

**For each feature you implement:**

1. **Update changelog:**
   ```bash
   docmgr changelog update --ticket YAPP-DSL-GAPS-001 \
     --entry "Implemented box_mounts module" \
     --file-note "pkg/yappgen/modules/boxmounts/schema.yaml: schema"
   ```

2. **Relate files:**
   ```bash
   docmgr relate --ticket YAPP-DSL-GAPS-001 \
     --file-note "pkg/yappgen/modules/boxmounts/module.go: builder"
   ```

3. **Check off tasks:**
   ```bash
   docmgr tasks check --ticket YAPP-DSL-GAPS-001 --id <task_id>
   ```

4. **Create log entries for significant milestones:**
   ```bash
   docmgr add --ticket YAPP-DSL-GAPS-001 \
     --doc-type log \
     --title "2025-11-XX Implemented boxMounts" \
     --summary "Added external mounting tab support"
   ```

## Essential Files Reference

### YAPP Library Source

| File | Purpose | Location |
|------|---------|----------|
| YAPP_Template_v3.scad | Parameter reference | `../../YAPP_Template_v3.scad` |
| YAPPgenerator_v3.scad | Generator implementation | `../../YAPPgenerator_v3.scad` |
| CHANGELOG.md | Version history | `../../CHANGELOG.md` |

### DSL Implementation

| File | Purpose | Location |
|------|---------|----------|
| Module authoring guide | How to add features | `../../pkg/docs/tutorials/yapp-module-authoring-guide.md` |
| Implementation guide | Module system architecture | `../YAPP-MODULE-SYSTEM-001-*/playbook/module-system-implementation-guide.md` |
| Registry package | Core interfaces | `../../pkg/registry/` |
| Schemagen tool | Code generator | `../../cmd/schemagen/` |

### Example Files

| File | Shows | Location |
|------|-------|----------|
| YAPP_Demo_buttons2_v31.scad | Push buttons (user's target) | `../../examples/YAPP_Demo_buttons2_v31.scad` |
| YAPP_Demo_RealBox_v31.scad | Production box (9 features) | `../../examples/YAPP_Demo_RealBox_v31.scad` |
| YAPP_Demo_boxMounts_v30.scad | External mounting | `../../examples/YAPP_Demo_boxMounts_v30.scad` |
| YAPP_Demo_lightTubes_v30.scad | LED indicators | `../../examples/YAPP_Demo_lightTubes_v30.scad` |
| test-push-buttons.yaml | Current DSL test | `../../examples/test-push-buttons.yaml` |

### Existing Module Examples

| Module | Complexity | Location | Study for |
|--------|------------|----------|-----------|
| pushbuttons | High | `../../pkg/yappgen/modules/pushbuttons/` | Nested objects, shape enums |
| pcbstands | Low | `../../pkg/yappgen/modules/pcbstands/` | Simple flat structure |
| snapjoins | Medium | `../../pkg/yappgen/modules/snapjoins/` | Side enum, flag appending |
| cutouts | High | `../../pkg/yappgen/modules/cutouts/` | Face distribution, shape logic |
| connectors | Medium | `../../pkg/yappgen/modules/connectors/` | Multiple required fields |

## Critical Gaps Summary

### Missing Feature Arrays (8)

1. ❌ **boxMounts** - External mounting tabs (user sees "little plates")
2. ❌ **lightTubes** - LED light pipes
3. ❌ **labelsPlane** - Text labels
4. ❌ **ridgeExt*** - Ridge extensions (4 arrays)
5. ❌ **displayMounts** - Display module mounting
6. ❌ **pcb array** - Multi-PCB support

### Missing Flags on Implemented Features

**pcb_stands missing:**
- Corner placement (yappAllCorners, yappFrontLeft, etc.) - **CRITICAL**
- Shell part (yappBoth, yappLidOnly, yappBaseOnly) - **CRITICAL**
- Treatment (yappPin, yappHole, yappTopPin)
- yappPCBName, yappSelfThreading, yappNoFillet

**connectors missing:**
- Corner placement flags
- yappCountersink, yappThroughLid, yappNoInternalFillet
- yappPCBName

**snap_joins missing:**
- yappOrigin/yappCenter, yappSymmetric, yappRectangle

**cutouts missing:**
- yappPolygonDef, yappMaskDef, yappRing, yappSphere
- yappAltOrigin, yappFromInside

**push_buttons missing:**
- yappAltOrigin, yappPCBName

## User's Specific Issues

**Observation:** "4 stalactites from top to pcb pillars (of which there are 4), but the scad one has only 2"

**Root cause:** SCAD file defines 1 pcb_stand at [5,5]. User's DSL likely has 4 explicit entries. The discrepancy is because:
- YAPP might auto-generate corners from single definition (needs research)
- OR user manually defined 4 stands in DSL (workaround for missing corner flags)

**Observation:** "Little plates sticking out to fasten it"

**Root cause:** boxMounts array (external mounting tabs) - **NOT IMPLEMENTED**

**Observation:** "Hinges, side of the box are open"

**Root cause:** Large cutouts in front/back (shellWidth-6, shellHeight-4) create nearly full-wall openings. NOT actual hinges (ridgeExt not defined in SCAD file).

## Next Steps

### Immediate Actions

1. **Get user's DSL YAML** - Compare with SCAD file to confirm gaps
2. **Research corner placement** - Study YAPPgenerator_v3.scad for auto-corner logic
3. **Implement Priority 1 flags** - Corner placement + shell part control (2 weeks)

### Short-term (if user needs them)

4. **Implement boxMounts** - External mounting tabs (3-4 days)
5. **Implement lightTubes** - LED indicators (2-3 days)

### Medium-term

6. **Implement labelsPlane** - Text labels (5-7 days)
7. **Implement ridgeExt** - Split openings (3-4 days)
8. **Implement displayMounts** - Display mounting (7-10 days)

## Resources

### Documentation

- **YAPP GitBook:** https://mrwheel-docs.gitbook.io/yappgenerator_en/
- **Local docs:** `go run ./cmd/yappctl help yapp-docs-analysis-playbook`
- **Module authoring:** `go run ./cmd/yappctl help yapp-module-authoring-guide`
- **DSL reference:** `go run ./cmd/yappctl help yapp-dsl-reference`

### Related Tickets

- **YAPP-MODULE-SYSTEM-001** - Module system implementation (completed)
  - Location: `../YAPP-MODULE-SYSTEM-001-*/`
  - Status: 28/32 tasks complete, end-to-end pipeline working
  - Key docs: `playbook/module-system-implementation-guide.md`

- **YAPP-PUSH-BUTTONS-001** - Push buttons implementation (completed)
  - Location: `../../YAPP-PUSH-BUTTONS-001-*/`
  - Contains design debates and prototype findings
  - Key docs: `reference/debate-round-*.md`

- **YAPP-ENCL-DSL-001** - Original DSL design
  - Location: `../../YAPP-ENCL-DSL-001-*/`
  - Contains completeness analysis and semantic debates
  - Key docs: `debate/01-round-1-completeness-*.md`

### Commands Reference

**View tasks:**
```bash
docmgr tasks list --ticket YAPP-DSL-GAPS-001
```

**Check off task:**
```bash
docmgr tasks check --ticket YAPP-DSL-GAPS-001 --id <number>
```

**Add changelog entry:**
```bash
docmgr changelog update --ticket YAPP-DSL-GAPS-001 \
  --entry "Your change description" \
  --file-note "path/to/file.go: what changed"
```

**Relate files to ticket:**
```bash
docmgr relate --ticket YAPP-DSL-GAPS-001 \
  --file-note "path/to/file.go: description"
```

**Validate schema:**
```bash
go run ./cmd/schemagen validate pkg/yappgen/modules/yourmodule/schema.yaml
```

**Generate code:**
```bash
go run ./cmd/schemagen discover
```

**Test end-to-end:**
```bash
go run ./cmd/yappctl resolve -i test.yaml -o /tmp/resolved.yaml
go run ./cmd/yappctl generate -i /tmp/resolved.yaml -o /tmp/output.scad
```

## Status

Current status: **active**

**Progress:**
- ✅ Gap analysis complete
- ✅ Feature inventory cataloged
- ✅ User case study documented
- ✅ Priority roadmap defined
- ⏳ Awaiting user's DSL YAML for confirmation
- ⏳ Implementation pending

## Topics

- yapp
- dsl
- analysis
- features

## Tasks

See [tasks.md](./tasks.md) for the current task list.

**Quick view:**
```bash
docmgr tasks list --ticket YAPP-DSL-GAPS-001
```

## Changelog

See [changelog.md](./changelog.md) for recent changes and decisions.

## Document Structure

### Essential Reading (Start Here)

1. **reference/01-yapp-feature-array-inventory.md** - Complete YAPP feature catalog
2. **analysis/01-user-case-study-buttons2-demo-discrepancies.md** - User's specific issues
3. **design-doc/01-priority-feature-roadmap.md** - Implementation plan

### Supporting Documents

- **design/** - Architecture and design documents
- **reference/** - API contracts, parameter specs
- **playbooks/** - Command sequences and test procedures (none yet)
- **scripts/** - Temporary code and tooling (none yet)
- **various/** - Working notes and research (none yet)
- **archive/** - Deprecated artifacts (none yet)

## Key Findings

### What Works Now (5 modules)

✅ **pcb_stands** - Basic positioning (missing corner/shell flags)
✅ **connectors** - Screw standoffs (missing corner flags)
✅ **snap_joins** - Snap-fit joints (missing positioning flags)
✅ **cutouts** - Holes in faces (missing masks, advanced shapes)
✅ **push_buttons** - Button extenders (missing yappAltOrigin)

### Critical Gaps (Priority 1)

❌ **Corner placement flags** - Auto-generate 4 standoffs from 1 definition
❌ **Shell part flags** - Control base/lid/both
❌ **Standoff treatment** - Pin vs hole configuration

**Impact:** Can't replicate common YAPP patterns efficiently

**Effort:** 2 weeks

### Common Features (Priority 2)

❌ **boxMounts** - External mounting tabs (user needs this)
❌ **lightTubes** - LED indicators
❌ **labelsPlane** - Text labels
❌ **ridgeExt*** - Split openings (4 arrays)

**Impact:** Can't create production-ready enclosures

**Effort:** 4-6 weeks

### Advanced Features (Priority 3)

❌ **displayMounts** - LCD/OLED mounting (17 parameters!)
❌ **Multi-PCB** - Multiple boards in one enclosure
❌ **Masks** - Honeycomb ventilation patterns
❌ **Custom polygons** - User-defined shapes

**Impact:** Can't handle advanced use cases

**Effort:** 6-8 weeks

### Corner Placement Research (2025-11-16)

- `pcbHolders()` mirrors a single `[x,y]` origin across the other three corners whenever `yappAllCorners` or individual `yappFront*/Back*` flags are present; if no flags are present (`primeOrigin`), only the literal coordinate is emitted, so DSL users must expose the same selector semantics instead of forcing four explicit entries.
- `pcbPushdowns()` reuses the identical corner-selection logic but additionally respects `yappLidOnly`/`yappBaseOnly`, so part-selection flags have to flow through schema + builder alongside the placement flags.
- `shellConnectors()` shares the same flag matrix and uses `translate2Box_*` to flip coordinates relative to the box or PCB coordinate system, confirming we can implement a single normalization helper for both pcb_stands and connectors.

```2063:2083:YAPPgenerator_v3.scad
allCorners = (isTrue(yappAllCorners, stand)) ? true : false;
primeOrigin = (!isTrue(yappBackLeft, stand) && !isTrue(yappFrontLeft, stand) && !isTrue(yappFrontRight, stand) && !isTrue(yappBackRight, stand) && !isTrue(yappAllCorners, stand) ) ? true : false;
if (!isTrue(yappLidOnly, stand))
{
  if (primeOrigin || allCorners || isTrue(yappBackLeft, stand))
    translate([offsetX+connX, offsetY + connY, basePlaneThickness])
      pcbStandoff(...);
  if (allCorners || isTrue(yappFrontLeft, stand))
    translate([offsetX + lengthX - connX, offsetY + connY, basePlaneThickness])
      pcbStandoff(...);
  // remaining corners omitted
}
```

```4072:4124:YAPPgenerator_v3.scad
allCorners = (isTrue(yappAllCorners, conn)) ? true : false;
primeOrigin = (!isTrue(yappBackLeft, conn) && !isTrue(yappFrontLeft, conn) && !isTrue(yappFrontRight, conn) && !isTrue(yappBackRight, conn) && !isTrue(yappAllCorners, conn) ) ? true : false;
if (primeOrigin || allCorners || isTrue(yappBackLeft, conn))
  connectorNew(..., connX, connY, ...);
if (allCorners || isTrue(yappFrontLeft, conn))
  connectorNew(..., connX2, connY, ...);
// additional mirrored placements omitted
```

## Success Criteria

**Phase 1 complete when:**
- ✅ Corner placement flags work (yappAllCorners generates 4 standoffs)
- ✅ Shell part flags work (yappLidOnly creates lid-only features)
- ✅ Can replicate YAPP_Demo_RealBox_v31.scad with DSL
- ✅ User's buttons2 demo works without manual workarounds

**Phase 2 complete when:**
- ✅ All Priority 2 modules implemented
- ✅ Can create wall-mountable enclosures (boxMounts)
- ✅ Can add LED indicators (lightTubes)
- ✅ Can add text labels (labelsPlane)

**Phase 3 complete when:**
- ✅ 80%+ YAPP feature coverage
- ✅ Can handle display modules
- ✅ Can define multiple PCBs
- ✅ Can use ventilation masks

## Getting Help

**Stuck on YAPP semantics?**
- Check `YAPP_Template_v3.scad` comments for parameter meanings
- Search examples: `grep -r "yappAllCorners" examples/`
- Query docs DB: See `yapp-docs-analysis-playbook`

**Stuck on module system?**
- Read module authoring guide: `go run ./cmd/yappctl help yapp-module-authoring-guide`
- Study existing modules in `pkg/yappgen/modules/`
- Check implementation guide: `../YAPP-MODULE-SYSTEM-001-*/playbook/`

**Stuck on code generation?**
- Check schemagen templates: `pkg/schemagen/templates/`
- Run with verbose: `go run ./cmd/schemagen discover` (shows what it's doing)
- Look at generated code: `pkg/yappgen/modules/*/schema_gen.go`

**Need to understand user's issue?**
- Read case study: `analysis/01-user-case-study-buttons2-demo-discrepancies.md`
- Compare SCAD files: `examples/YAPP_Demo_buttons2_v31.scad` vs user's generated output
- Ask user for their DSL YAML (task #2)

## Glossary

**YAPP** - Yet Another Parametric Projectbox (OpenSCAD library)
**DSL** - Domain-Specific Language (our YAML format)
**SCAD** - OpenSCAD file format
**Module** - Self-contained feature handler (schema + builder)
**Schema** - YAML file defining valid fields for a module
**Builder** - Go function converting YAML → OpenSCAD arrays
**Registry** - Central list of all modules
**Schemagen** - Code generation tool
**Flag** - YAPP parameter that modifies behavior (yappCircle, yappAllCorners, etc.)

---

**Remember:** This is research and planning. Read the docs, understand the patterns, then implement incrementally. Test each feature with real examples before moving to the next.

Good luck! 🚀
