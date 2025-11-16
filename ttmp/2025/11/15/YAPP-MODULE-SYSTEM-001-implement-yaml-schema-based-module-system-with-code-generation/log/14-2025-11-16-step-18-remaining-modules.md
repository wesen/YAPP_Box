---
Title: 2025-11-16 Step 18 - remaining modules
Ticket: YAPP-MODULE-SYSTEM-001
Status: active
Topics:
    - yapp
    - dsl
    - codegen
    - architecture
DocType: log
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Converted connectors, snap_joins, cutouts to schema pipeline
LastUpdated: 2025-11-15T21:27:45.569837471-05:00
---


# 2025-11-16 Step 18 - remaining modules

<!-- Log entries in reverse chronological order (newest first) -->

## 2025-11-16 - Converted connectors, snap_joins, cutouts

Completed module system rollout by converting the remaining 3 modules:

**Connectors (order: 150):**
- 10 fields (x, y, stand_height, screw_d, screw_head_d, insert_d, outside_d, insert_depth, pcb_gap, fillet_radius)
- Typed builder following pcb_stands pattern
- 3 test cases (minimal, full, missing required)

**Snap Joins (order: 300):**
- 3 fields (pos, width, side)
- Side enum validated (left/right/front/back)
- Side flag appended after positional params
- 3 test cases including invalid side

**Cutouts (order: 400):**
- 9 fields (face, from_back, from_left, width, length, radius, shape, depth, angle)
- Face enum (front/back/left/right/top/lid/bottom/base)
- Shape enum (rectangle/circle/rounded_rect/circle_with_flats/circle_with_key)
- Special handling: returns map by face (cutoutsFront, cutoutsBack, etc.)
- Shape-specific dimension usage (circles ignore width/length)
- 3 test cases

**Generated:**
- `schemagen discover` now finds 5 modules
- All schema_gen.go/schema_gen_test.go files created
- `modules_gen.go` registers all 5 modules in order
- `go test ./...` passes for all modules

**Status:** All core YAPP features now use the schema-based module system!
