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
    - Path: /home/manuel/code/others/YAPP_Box/cmd/yappctl/generate_command.go
      Note: flags + timeout
    - Path: /home/manuel/code/others/YAPP_Box/examples/compare/cutouts_all_faces.yaml
      Note: add pos_z to raise side-face cutouts
    - Path: /home/manuel/code/others/YAPP_Box/examples/yapp-demo-buttons-v30.yaml
      Note: polygon+mask base update
    - Path: /home/manuel/code/others/YAPP_Box/examples/yapp-demo-lighttubes.yaml
      Note: updated with pos_z for side face cutouts
    - Path: /home/manuel/code/others/YAPP_Box/examples/yapp-demo-masks.yaml
      Note: example demonstrating masks in DSL
    - Path: /home/manuel/code/others/YAPP_Box/pkg/cli/generatorcli/generator.go
      Note: STL modes + copy
    - Path: /home/manuel/code/others/YAPP_Box/pkg/docs/tutorials/yapp-dsl-reference.md
      Note: document light_tubes
    - Path: /home/manuel/code/others/YAPP_Box/pkg/resolver/resolver.go
      Note: treat cutouts.mask.preset as string literal for resolution
    - Path: /home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/cutouts/module.go
      Note: mask encoding
    - Path: /home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/cutouts/registry.go
      Note: ValidateConstraints implemented
    - Path: /home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/cutouts/schema.yaml
      Note: add mask object and presets
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/log/03-2025-11-17-implementation-diary-pushbuttons-flags.md
      Note: push_buttons flags
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/scripts/validate-task16-pushbuttons-flags.sh
      Note: 'validation script for Task #16'
    - Path: YAPP_Template_v3.scad
      Note: canonical YAPP feature definitions
    - Path: examples/YAPP_Demo_buttons2_v31.scad
      Note: user's target example
    - Path: examples/YAPP_Demo_buttons_v30.scad
      Note: reference for parity
    - Path: examples/YAPP_Demo_buttons_v30_computed.scad
      Note: debug version with computed values
    - Path: examples/test-push-buttons.yaml
      Note: current DSL test case
    - Path: examples/yapp-demo-buttons-v30.yaml
      Note: updated with vars and cutouts base/front/back
    - Path: examples/yapp-demo-lighttubes-with-errors.yaml
      Note: test case with resolver errors (snap_joins uses 'sides' array instead of 'side' string) - useful for testing resolver error message improvements
    - Path: examples/yapp-demo-lighttubes.yaml
      Note: Fixed cutouts structure (flat list with face field) and changed polygon to rounded_rect
    - Path: examples/yapp-demo-polygon-test.yaml
      Note: Test YAML with polygon examples
    - Path: log/01-2025-11-16-implementation-diary-lighttubes-module.md
      Note: Complete implementation diary documenting lightTubes module implementation process
    - Path: pkg/resolver/resolver.go
      Note: treats enum arrays as literals (corners)
    - Path: pkg/resolver/validation.go
      Note: Calls ValidateStructure/ValidateConstraints but they return nil without checking
    - Path: pkg/yappgen/emit.go
      Note: emit wall heights
    - Path: pkg/yappgen/features.go
      Note: Migrated cutouts to module Build
    - Path: pkg/yappgen/map.go
      Note: Removed distributeCutouts/buildCutoutParams and friends
    - Path: pkg/yappgen/model.go
      Note: Cutouts type change; removed Cutout struct
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
    - Path: pkg/yappgen/modules/cutouts/module.go
      Note: Added polygonPresetFlag function and polygon preset emission
    - Path: pkg/yappgen/modules/cutouts/registry.go
      Note: ValidateStructure and ValidateConstraints are stubbed out (TODO) - enum validation not implemented
    - Path: pkg/yappgen/modules/cutouts/schema.yaml
      Note: Added polygon shape enum and polygon preset field
    - Path: pkg/yappgen/modules/lighttubes/module.go
      Note: Build function converting DSL to YAPP array format with positional params and flags
    - Path: pkg/yappgen/modules/lighttubes/registry.go
      Note: FeatureModule registration
    - Path: pkg/yappgen/modules/lighttubes/schema.yaml
      Note: lightTubes schema with required fields (x
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
    - Path: pkg/yappgen/schema.go
      Note: Added polygon support to ShapeFlag for legacy code path
    - Path: pkg/yappgen/yappgen_test.go
      Note: Adjusted cutouts tests to module path
    - Path: ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/design-doc/02-pcb-stand-flag-support.md
      Note: design
    - Path: ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/design-doc/03-connector-flag-support.md
      Note: design
    - Path: ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/design-doc/04-boxmounts-module-support.md
      Note: design
ExternalSources: []
Summary: Implemented Priority 1 flags (corners, shell parts, snap joins) and boxMounts module. DSL now achieves ~70% YAPP feature coverage with working v30 demo.
LastUpdated: 2025-11-16T00:30:35.203590956-05:00
---


























# DSL Feature Gaps Analysis - Missing YAPP Arrays

## Overview

This ticket analyzes missing YAPP features that prevent the DSL from achieving parity with real SCAD examples. A user attempted to replicate `YAPP_Demo_buttons2_v31.scad` using the DSL and encountered visual discrepancies: extra PCB pillars, external mounting plates, and potential hinge-like structures.

**Root cause (initial):** The DSL implemented 5 of 13+ YAPP feature arrays (38% coverage) and was missing critical flags on implemented features.

**Current status:** DSL now implements 6 modules with full Priority 1 flag support (~70% coverage). Working v30 demo validates corner placement, shell parts, snap joins, and vertical dimensions. Remaining work focuses on LED indicators, text labels, and advanced cutout shapes.

**Goal:** Achieve 80%+ YAPP feature coverage for production-ready enclosures.

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

### 3. Understand the Module System ⚠️ CRITICAL

**⚠️ READ THIS FIRST:** The module system is the foundation for all feature implementation. Every new feature follows this pattern.

**Primary documentation:**
- **Module System Implementation Guide:** `../YAPP-MODULE-SYSTEM-001-*/playbook/module-system-implementation-guide.md`
  - Complete guide for newcomers
  - Explains architecture, code generation, validation
  - Step-by-step module authoring instructions
  - Troubleshooting and best practices

**Quick reference:**
```bash
# View the full guide
cat ../YAPP-MODULE-SYSTEM-001-*/playbook/module-system-implementation-guide.md

# Or use the help command
go run ./cmd/yappctl help yapp-module-authoring-guide
```

**Key concepts (from the guide):**
- Schema-driven validation (YAML schemas define fields)
- Code generation (schemagen tool)
- Auto-registration (modules_gen.go)
- Two-phase validation (structure + constraints)
- Module = schema.yaml + module.go + registry.go

**How to add a new module (summary):**
1. Create `pkg/yappgen/modules/yourmodule/schema.yaml`
2. Run `go run ./cmd/schemagen discover` (generates schema_gen.go)
3. Write `module.go` with Build() function
4. Write `registry.go` with NewModule() function
5. Add to `features.go` and `model.go`
6. Test end-to-end

**See the guide for complete details!**

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
| **⚠️ Module System Implementation Guide** | **CRITICAL - Complete guide for implementing modules** | `../YAPP-MODULE-SYSTEM-001-*/playbook/module-system-implementation-guide.md` |
| Module authoring guide | Quick reference for adding features | `../../pkg/docs/tutorials/yapp-module-authoring-guide.md` |
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

### Immediate Actions (Phase 2)

1. **Implement lightTubes** - LED light pipes (2-3 days, task #8)
   - Schema: pos, diameter, height, wall, shape, lens_thickness
   - Flags: coordinate, origin, no_fillet, pcb_name, through_lid
   - Example: `YAPP_Demo_lightTubes_v30.scad`

2. **Implement cutout polygon shapes** - Custom shapes (2-3 days, task #13)
   - Add `polygon` shape enum value
   - Support shape presets (shapeHexagon, shapeArrow, etc.)
   - Emit `yappPolygon` flag + preset reference

3. **Implement cutout masks** - Ventilation patterns (3-4 days, task #14)
   - Add mask field with presets (maskHoneycomb, maskCircles, etc.)
   - Emit `yappMaskDef` or `[yappMaskDef, hOffset, vOffset, rotation]`
   - Test with base cutout from v30 demo

### Short-term (Phase 2 continued)

4. **Implement labelsPlane** - Text labels (5-7 days, task #9)
   - Schema: pos, text, font, size, depth, face, rotation
   - Flags: alignment, direction
   - Complex: 12 parameters + text rendering

5. **Implement ridgeExt** - Split openings (3-4 days, task #10)
   - Four separate arrays (Front/Back/Left/Right)
   - Schema: pos, width, height
   - Flags: origin, coordinate, pcb_name

### Medium-term (Phase 3)

6. **Implement displayMounts** - Display mounting (7-10 days)
   - 17 parameters (most complex module)
   - Window cutouts, pin mounts, bevels
   - Flags: origin, coordinate, self_threading

7. **Multi-PCB support** - pcb array (low priority)
   - Requires schema changes to support multiple PCB definitions
   - All modules already support pcb_name flag

## Resources

### Documentation

- **YAPP GitBook:** https://mrwheel-docs.gitbook.io/yappgenerator_en/
- **Local docs:** `go run ./cmd/yappctl help yapp-docs-analysis-playbook`
- **Module authoring:** `go run ./cmd/yappctl help yapp-module-authoring-guide`
- **DSL reference:** `go run ./cmd/yappctl help yapp-dsl-reference`

### Related Tickets

- **⚠️ YAPP-MODULE-SYSTEM-001** - Module system implementation (completed) - **CRITICAL REFERENCE**
  - Location: `../YAPP-MODULE-SYSTEM-001-*/`
  - Status: 28/32 tasks complete, end-to-end pipeline working
  - **Key docs:** `playbook/module-system-implementation-guide.md` - **READ THIS FIRST**
  - This is the authoritative guide for implementing any new YAPP DSL module
  - Explains the entire architecture, code generation, validation, and module authoring workflow
  - Every feature implementation in this ticket follows patterns from this guide

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
- ✅ Priority 1 flags implemented (corners, shell parts, treatment, coordinate)
- ✅ boxMounts module implemented
- ✅ snap_joins flags implemented (alignment, symmetric, diamond)
- ✅ Vertical dimensions support (baseWallHeight, lidWallHeight, ridgeHeight)
- ✅ Working v30 demo example with STL validation
- ⏳ Priority 2 modules pending (lightTubes, labelsPlane, ridgeExt)
- ⏳ Advanced cutout features pending (polygon, masks)

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

### What Works Now (6 modules + full flag support)

✅ **pcb_stands** - Full flag support (corners[], shell_part, treatment, coordinate, no_fillet, pcb_name, self_threading)
✅ **connectors** - Full flag support (corners[], coordinate, countersink, through_lid, no_fillet, no_internal_fillet, pcb_name, self_threading)
✅ **snap_joins** - Full flag support (alignment, symmetric, diamond)
✅ **cutouts** - Basic shapes (rectangle, circle, rounded_rect, circle_with_flats, circle_with_key) on all 6 faces
✅ **push_buttons** - Full nested object support with shape presets
✅ **box_mounts** - External mounting tabs with face selection

### Implemented (Priority 1) ✅

✅ **Corner placement flags** - `corners: [front_left, back_right]` auto-generates mirrored standoffs
✅ **Shell part flags** - `shell_part: base_only/lid_only/both`
✅ **Standoff treatment** - `treatment: pin/hole/top_pin`
✅ **Coordinate systems** - `coordinate: pcb/box/box_inside`
✅ **boxMounts module** - External mounting tabs
✅ **Snap join flags** - `alignment: center`, `symmetric: true`, `diamond: true`
✅ **Vertical dimensions** - `baseWallHeight`, `lidWallHeight`, `ridgeHeight`

**Status:** Phase 1 complete! Can replicate most common YAPP patterns.

### Implemented (Priority 2) ✅

✅ **lightTubes** - LED light pipes with circle/rectangle shapes, lens thickness, coordinate flags
✅ **Cutout polygons** - yappPolygon + shape presets (hexagon, arrow, 6pt_star, iso_triangle, triangle)
✅ **Cutout coordinate/origin flags** - yappCoordBox, yappCoordPCB, yappCoordBoxInside, yappCenter, yappAltOrigin
✅ **Face-relative cutout naming** - from_face_left, from_face_bottom, from_face_back with strict validation

**Status:** Phase 2 partially complete! Can add LED indicators, custom polygon cutouts, and precise coordinate control.

### Remaining (Priority 2)

❌ **labelsPlane** - Text labels (5-7 days)
❌ **ridgeExt*** - Ridge extensions for split openings (3-4 days)
❌ **Cutout masks** - yappMaskDef + ventilation patterns (3-4 days)

**Impact:** Can't add text labels or ventilation patterns

**Effort:** 2-3 weeks

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
- **⚠️ READ FIRST:** `../YAPP-MODULE-SYSTEM-001-*/playbook/module-system-implementation-guide.md`
  - Complete guide covering architecture, code generation, validation, module authoring
  - Written for newcomers with step-by-step instructions
  - Includes troubleshooting and best practices
- Quick reference: `go run ./cmd/yappctl help yapp-module-authoring-guide`
- Study existing modules in `pkg/yappgen/modules/` (especially `boxmounts/` and `lighttubes/`)

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
