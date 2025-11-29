---
Title: Film Developer Enclosure Dimension Reference
Ticket: FILM-DEV-ENCLOSURE-001
Status: active
Topics:
    - film-developer
    - enclosure
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: pkg/yappgen/assets/YAPPgenerator_v3.scad
      Note: Main YAPP generator SCAD file
    - Path: ttmp/2025/11/28/FILM-DEV-ENCLOSURE-001-film-developer-enclosure-dimension-reference/reference/02-cutout-offset-computation.md
      Note: Detailed explanation of cutout offset computation
    - Path: ttmp/2025/11/28/FILM-DEV-ENCLOSURE-001-film-developer-enclosure-dimension-reference/reference/03-yapp-dsl-cutout-position-transformation.md
      Note: DSL-side coordinate transformation documentation
ExternalSources: []
Summary: Documentation for YAPP enclosure dimension derivation and cutout offset computation, especially for circle cutouts
LastUpdated: 2025-11-28T19:20:42.639102131-05:00
---





# Film Developer Enclosure Dimension Reference

## Overview

This ticket documents how YAPP generator computes enclosure dimensions and cutout positions, with special focus on understanding circle cutout positioning. This is critical for the film developer project which uses 40mm button cutouts that must be precisely positioned.

The ticket also includes design work for extending the YAML DSL to support per-side clearance configuration and direct final dimension specification.

## Reference Documents

- **[Enclosure Dimension Derivation](./reference/01-enclosure-dimension-derivation.md)** - How PCB geometry and enclosure settings convert to inside/outer dimensions
- **[Cutout Offset Computation](./reference/02-cutout-offset-computation.md)** - Detailed explanation of how cutout positions and offsets are computed in SCAD generator, especially for circle cutouts. **Key finding:** Circle cutouts specify the bottom-left corner of the bounding box, not the center, unless `yappCenter` flag is used.
- **[YAPP DSL Cutout Position Transformation](./reference/03-yapp-dsl-cutout-position-transformation.md)** - Explains how DSL face-relative coordinates (`from_face_left`, `from_face_bottom`, `from_face_back`) transform into SCAD array positions, and how the `origin: center` flag affects positioning.

## Analysis Documents

- **[Enclosure Dimension Computation Analysis](./analysis/01-enclosure-dimension-computation-analysis.md)** - Comprehensive analysis of how YAPPgenerator_v3.scad computes final enclosure dimensions, including YAML DSL configuration capabilities and limitations

## Design Documents

- **[Per-Side Clearance and Final Dimensions](./design/01-per-side-clearance-and-final-dimensions.md)** - Design for adding per-side clearance and final dimension configuration to YAML DSL with mutual exclusivity validation

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- film-developer
- enclosure

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
