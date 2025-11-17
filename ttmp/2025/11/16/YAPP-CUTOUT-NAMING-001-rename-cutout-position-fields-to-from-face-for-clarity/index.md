---
Title: Rename cutout position fields to from_face_* for clarity
Ticket: YAPP-CUTOUT-NAMING-001
Status: active
Topics:
    - yapp
    - dsl
    - api
    - naming
DocType: index
Intent: short-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Replace from_back/from_left cutout fields with from_face_left (all faces), from_face_bottom (sides), from_face_back (base/lid) for self-documenting positioning with strict validation"
LastUpdated: 2025-11-16T16:52:35.757727267-05:00
---


# Rename cutout position fields to from_face_* for clarity

## Overview

This ticket implements a breaking change to cutout position field naming based on unanimous approval from a 4-participant debate. The current `from_back`/`from_left` fields are ambiguous and change meaning by face type. We're replacing them with self-documenting face-relative names.

**Current problems:**
- `from_back` doesn't mean "from back face"—it means "from back edge of box globally"
- `from_left` means vertical (posz) for side faces, horizontal (posy) for base/lid
- `pos_z` is a workaround override for the `from_left` confusion

**New naming scheme:**
- `from_face_left` (all faces) — horizontal position along the face
- `from_face_bottom` (front/back/left/right only) — vertical position from bottom edge
- `from_face_back` (base/lid only) — depth position from back edge

**Validation:**
- Using wrong field for face type produces clear error with examples
- No backwards compatibility—clean break for clarity

**Scope:**
- Update schema with conditional field requirements
- Implement CustomValidate() with face-type checking
- Migrate 60+ usages across 9 example files
- Regenerate comparison STLs
- Update DSL reference documentation

**Rationale:** Debate (YAPP-DSL-GAPS-001) showed unanimous support. Names become self-documenting, eliminate mental model complexity, and teach users through validation errors.

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- yapp
- dsl
- api
- naming

## Tasks

See [tasks.md](./tasks.md) for the current task list.

## Changelog

See [changelog.md](./changelog.md) for recent changes and decisions.

## Structure

- design/ - Architecture and design documents
- reference/ - Prompt packs, API contracts, context summaries
- playbooks/ - Command sequences and test procedures
- scripts/ - Temporary code and tooling
- various/ - Working notes and research
- archive/ - Deprecated or reference-only artifacts
