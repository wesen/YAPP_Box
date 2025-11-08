# PCB Stands: Parameter Order Verification

Status: draft
Owner: manuel
LastUpdated: 2025-11-05

## Purpose

Confirm that the parameter order documented for PCB stands matches the current implementation (v3.3.8). Differences between v3.0-era documentation and v3.3.8 may cause confusion.

## Verification Checklist

1. Identify the source module implementing PCB stands in the repository.
2. Extract the parameter signature (positional and named flags) from v3.3.8.
3. Compare against the GitBook page for PCB Stands.
4. Update the GitBook page to match the implementation.
5. Add migration notes if order changed since v3.0.

## Findings (initial)

- Known risk: Documentation may list parameters in a different order than the actual module. Confirm exact order before updating examples.

### v3.3.8 Parameter Order (from template/generator comments)

Required:
- `p(0)` posx
- `p(1)` posy

Optional:
- `p(2)` heightToBottomOfPCB (Default = `standoffHeight`)
- `p(3)` pcbGap (Default = `-1`; resolves to `pcbThickness` for `yappCoordPCB`, else `0`)
- `p(4)` standoffDiameter (Default = `standoffDiameter`)
- `p(5)` standoffPinDiameter (Default = `standoffPinDiameter`)
- `p(6)` standoffHoleSlack (Default = `standoffHoleSlack`)
- `p(7)` filletRadius (`0` = Auto)
- Named flags (order-insensitive):
  - `n(a)` `{ <yappBoth> | yappLidOnly | yappBaseOnly }`
  - `n(b)` `{ <yappPin>, yappHole, yappTopPin }`
  - `n(c)` corners `{ yappAllCorners | yappFrontLeft | yappFrontRight | yappBackLeft | yappBackRight }`
  - `n(d)` coords `{ <yappCoordPCB> | yappCoordBox | yappCoordBoxInside }`
  - `n(e)` `{ yappNoFillet }`
  - (optional) `[yappPCBName, "Main"]`

Recommendation for docs: present the required pair first, then group optionals by role (geometry, placement, flags) and mark defaults inline.

## Next Actions

- Cross-check implementation in YAPP_Box (v3.3.8 tag or latest main).
- Draft corrected parameter table and example usage.
- Submit PR or change request to the GitBook maintainers.

---
Title: 'PCB Stands: Parameter Order Verification'
Ticket: YAPP-DOCS-001
Status: active
Topics:
    - yappgenerator
    - documentation
DocType: reference
Intent: long-term
Owners:
    - manuel
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2025-11-05T17:17:37.292584554-05:00
---


# PCB Stands: Parameter Order Verification

## Goal

<!-- What is the purpose of this reference document? -->

## Context

<!-- Provide background context needed to use this reference -->

## Quick Reference

<!-- Provide copy/paste-ready content, API contracts, or quick-look tables -->

## Usage Examples

<!-- Show how to use this reference in practice -->

## Related

<!-- Link to related documents or resources -->
