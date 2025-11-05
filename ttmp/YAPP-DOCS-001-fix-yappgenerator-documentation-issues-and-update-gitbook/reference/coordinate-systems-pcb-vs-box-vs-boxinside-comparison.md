# Coordinate Systems: PCB vs Box vs BoxInside Comparison

Status: draft
Owner: manuel
LastUpdated: 2025-11-05

## Overview

YAPPgenerator uses multiple coordinate systems to position features. This document provides a concise comparison and guidance for choosing the right one.

Default origin in most examples is `yappCoordPCB` (PCB[0,0,0]). Be explicit in documentation and examples to avoid confusion.

## Quick Comparison

| System             | Origin reference       | Typical use                                           | Notes |
|--------------------|------------------------|-------------------------------------------------------|-------|
| `yappCoordPCB`     | PCB [0,0,0]            | Place features relative to PCB coordinates            | Default in many examples |
| `yappCoordBox`     | Box outer origin       | Place features relative to outer box coordinates      | Useful for exterior features |
| `yappCoordBoxInside` | Box inner origin     | Place features relative to inner cavity               | Good for internal clearances |

## Selection Guidance

- Prefer `yappCoordPCB` for cutouts and components that reference PCB positions.
- Use `yappCoordBox` for external features (vents, labels) that align to the enclosure shell.
- Use `yappCoordBoxInside` for interior standoffs, supports, and cable channels.

## Center vs Origin Flags

- For circular features, prefer center-based positioning via `yappCenter` where available to reduce manual offset calculations.
- When not using `yappCenter`, ensure the radius/diameter is accounted for in offsets from the chosen origin.

## Diagram Placeholders

Add three simple diagrams showing the same point across systems:

- Diagram A: PCB coordinate frame
- Diagram B: Box (outer) coordinate frame
- Diagram C: Box (inside) coordinate frame

These diagrams should share an example feature (e.g., a hole at PCB (x,y)) and show how it maps to each system.

## Validation

Verify examples explicitly set the coordinate system via the named parameter slot (`n(a)` in many features). Avoid relying on implicit defaults in documentation.

---
Title: 'Coordinate Systems: PCB vs Box vs BoxInside Comparison'
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
LastUpdated: 2025-11-05T17:17:36.981955816-05:00
---


# Coordinate Systems: PCB vs Box vs BoxInside Comparison

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
