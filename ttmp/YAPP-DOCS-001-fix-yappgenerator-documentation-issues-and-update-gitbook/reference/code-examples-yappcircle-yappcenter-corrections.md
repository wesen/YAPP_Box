# Code Examples: yappCircle + yappCenter Corrections

Status: draft
Owner: manuel
LastUpdated: 2025-11-05

## Overview

Several GitBook pages use `yappCircle` without enabling center-based positioning via `yappCenter`. This can lead to off-by-radius placement and confusing coordinates.

This reference consolidates all affected examples and shows the corrected usage pattern. Update the GitBook pages accordingly and then re-run the analysis.

## Affected Pages

- Light Tubes — `https://mrwheel-docs.gitbook.io/yappgenerator_en/light-tubes`
- Push Buttons — `https://mrwheel-docs.gitbook.io/yappgenerator_en/push-buttons`

## Correct Usage Pattern

- When using circular features with `yappCircle`, add the center-based flag `yappCenter` in the named-parameter flags list for that feature.
- Keep the coordinate system explicit via `n(a)` (e.g., `yappCoordPCB`, `yappCoordBox`, or `yappCoordBoxInside`).
- Preserve existing flags like `yappNoFillet`.

Conceptually, turn:

- `... yappCircle ...`

into:

- `... yappCircle, yappCenter ...`

The exact slot for `yappCenter` depends on the feature’s signature (named flag list). In examples that already pass `n(a)` and `n(b)`, add `yappCenter` to the next available named flag position for that feature.

## Light Tubes — Example Correction

Original snippets (from analysis DB):

- `lightTubes = [ [15, 20, 1.5, 5, 1, 2, yappCircle, 0.5], ... ];`
- `lightTubes = [ [15, 20, 2, 6, 1, 5, yappCircle, 0.5], ... ];`

Recommended (conceptual) correction:

- `lightTubes = [ [15, 20, 1.5, 5, 1, 2, yappCircle, 0.5, /* n(a)=coord? */, yappCenter], ... ];`

Notes:

- Ensure the coordinate system (`n(a)`) is set as intended (default is `yappCoordPCB`).
- Insert `yappCenter` in the named flags position for center-based placement.

## Push Buttons — Example Correction

Original snippet (from analysis DB):

- `pushButtons = [ [84.2, 30.7, 8, 8, 0, 2, 1, 3.5, yappCircle], ... ];`

Recommended (conceptual) correction:

- `pushButtons = [ [84.2, 30.7, 8, 8, 0, 2, 1, 3.5, yappCircle, yappCenter], ... ];`

Notes:

- Keep any existing coordinate system flags and other named flags intact.
- If the shape is chosen via a positional parameter elsewhere, move `yappCenter` to the named flags list.

## Validation Steps

1. Update the GitBook pages with the corrected examples.
2. Re-scrape documentation and verify no `yappCircle` examples lack `yappCenter`:
   - `python3 analysis/scrape_gitbook.py`
   - `analysis/query_db.py query "SELECT COUNT(1) FROM gitbook_code_examples WHERE code_block LIKE '%yappCircle%' AND (issue_description LIKE '%yappCenter%' OR has_syntax_issue = 1)"`
3. If issues remain, refine the example placement of `yappCenter` based on the feature’s signature.

## Appendix: Discovery Queries

Use these from `ttmp/YAPP-DOCS-001-.../analysis/`:

```
./query_db.py query "SELECT DISTINCT p.title, p.url FROM gitbook_pages p JOIN gitbook_code_examples c ON c.page_id = p.id WHERE c.issue_description LIKE '%yappCenter%';"
./query_db.py query "SELECT p.title, c.code_block FROM gitbook_pages p JOIN gitbook_code_examples c ON c.page_id = p.id WHERE c.issue_description LIKE '%yappCenter%';"
```

---
Title: 'Code Examples: yappCircle + yappCenter Corrections'
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
LastUpdated: 2025-11-05T17:17:36.834835386-05:00
---


# Code Examples: yappCircle + yappCenter Corrections

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
