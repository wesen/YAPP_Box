# Implementation Diary: Cutout Field Renaming

**Date:** 2025-11-16  
**Ticket:** YAPP-CUTOUT-NAMING-001

## Session 1: Setup and Schema Update

### 16:52 - Ticket Creation

Created ticket using docmgr:
```bash
docmgr create-ticket --ticket "YAPP-CUTOUT-NAMING-001" \
  --title "Rename cutout position fields to from_face_* for clarity" \
  --topics "yapp,dsl,api,naming"
```

Updated index.md with:
- Summary of the problem and solution
- New naming scheme specification
- Scope and rationale from debate

Created tasks.md with implementation checklist and 18+ test cases.

### 17:00 - Schema Update Complete

**What I did:**
1. Updated `pkg/yappgen/modules/cutouts/schema.yaml`:
   - Removed: `from_back`, `from_left`, `pos_z`
   - Added: `from_face_left` (required), `from_face_bottom` (optional), `from_face_back` (optional)
   - Updated field descriptions with clear face-type-specific guidance
   - Updated all 6 test cases to use new field names

2. Ran `go run ./cmd/schemagen discover` to regenerate `schema_gen.go`

3. Implemented `CustomValidate()` with strict face-type checking:
   - Side faces (front/back/left/right) require `from_face_bottom`, reject `from_face_back`
   - Horizontal faces (base/lid) require `from_face_back`, reject `from_face_bottom`
   - Error messages include helpful examples showing correct usage
   - Added imports for `fmt` and `strings`

4. Updated `Build()` method in `module.go`:
   - Removed old `FromBack`/`FromLeft`/`PosZ` field access
   - Added logic to determine pos0/pos1 based on face type and new fields:
     - Side faces: `from_face_left` → pos0, `from_face_bottom` → pos1
     - Base/lid: `from_face_back` → pos0, `from_face_left` → pos1
   - Verified module compiles successfully

**Findings:**
- Schema generation worked perfectly—new fields have correct types (FromFaceLeft float64, others *float64)
- CustomValidate() can check nil pointers to detect missing required fields per face type
- Position mapping logic is straightforward once face type is determined

**What worked:**
- Face-type-specific validation at CustomValidate() level (not schema level) is clean
- Error messages with embedded YAML examples will help users learn
- The field order swap for base/lid (from_face_back→pos0, from_face_left→pos1) matches YAPP's [X,Y] expectation

### Next Steps

Now need to migrate example YAML files. Found 60+ usages across 9 files:
- examples/compare/cutouts_all_faces.yaml
- examples/yapp-demo-lighttubes.yaml
- examples/compare/cutouts_polygons.yaml
- examples/yapp-demo-polygon-test.yaml
- examples/yapp-demo-lighttubes-with-errors.yaml
- examples/yapp-demo-buttons-v30.yaml
- examples/yapp-demo-buttons.yaml
- examples/yapp-mvp-combined.yaml
- examples/yapp-mvp-pcbstands-cutouts.yaml

Will update them systematically, then regenerate STLs to verify no visual regressions.

---

