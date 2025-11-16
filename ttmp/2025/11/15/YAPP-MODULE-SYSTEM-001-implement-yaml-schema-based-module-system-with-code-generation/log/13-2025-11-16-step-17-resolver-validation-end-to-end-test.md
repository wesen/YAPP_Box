---
Title: 2025-11-16 Step 17 - resolver validation + end-to-end test
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
Summary: Wired resolver validation and tested YAML to SCAD pipeline
LastUpdated: 2025-11-15T21:10:26.688741094-05:00
---


# 2025-11-16 Step 17 - resolver validation + end-to-end test

<!-- Log entries in reverse chronological order (newest first) -->

## 2025-11-16 - End-to-end YAML → SCAD pipeline working

Successfully wired resolver validation and tested the full pipeline:

**Resolver Integration:**
- Created `pkg/resolver/validation.go` with `validateStructure()` and `validateConstraints()`
- Modified `Resolve()` to call validation before and after expression resolution
- Moved strict-mode helpers to `pkg/resolver/strict.go` to avoid duplication
- Registry schemas now validate features during resolution

**End-to-End Test:**
1. Created `examples/test-push-buttons.yaml` with:
   - 4 pcb_stands using expressions (pcb_length - 3, etc.)
   - 2 push_buttons (rectangle + circle shapes)
   - Variables for dimensions
2. Ran `yappctl resolve -i examples/test-push-buttons.yaml -o /tmp/test-resolved.yaml`
   - ✅ Expressions resolved (pcb_length - 3 → 62)
   - ✅ Schema validation passed (no errors)
3. Ran `yappctl generate -i /tmp/test-resolved.yaml -o /tmp/test-output.scad`
   - ✅ Generated valid OpenSCAD code
   - ✅ `pcbStands` array with 4 entries
   - ✅ `pushButtons` array with 2 entries (correct positional params)
   - ✅ `printSwitchExtenders = true` set automatically

**Generated SCAD snippet:**
```
pushButtons =
[
  [15, 10, 8, 6, 0.5, 2.5, 5.5, 1, 3.5, undef, yappRectangle, undef, undef, 2, 2.5, 0.25, undef, yappCoordPCB],
  [50, 10, 10, 10, 5, 3, 5.5, 1, 3.5, undef, yappCircle, undef, undef, undef, undef, undef, undef, yappCoordPCB]
];
```

**Status:** The module system MVP is now functional end-to-end for push_buttons and pcb_stands!
