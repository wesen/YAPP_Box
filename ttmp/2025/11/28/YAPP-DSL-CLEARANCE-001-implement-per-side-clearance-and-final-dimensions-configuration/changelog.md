# Changelog

## 2025-11-28

- Initial workspace created


## 2025-11-28

Created ticket and design document for implementing per-side clearance and final dimensions configuration. Added 15 tasks for implementation.


## 2025-11-28

Implemented per-side clearance and final dimensions configuration. Added validation for mutual exclusivity, helper functions to detect configuration modes, and computation logic for reverse padding calculation. All tests passing.

### Related Files

- pkg/yappgen/model.go — Updated Model struct and BuildModel to handle per-side clearance and final dimensions
- pkg/yappgen/model_test.go — Added comprehensive unit tests for all new functionality
- pkg/yappgen/validation.go — Added validateEnclosureDimensions function for mutual exclusivity checks


## 2025-11-28

Updated YAPP DSL reference documentation with per-side clearance and final dimensions sections. Created comprehensive examples directory with valid examples (uniform, per-side, final dimensions, partial final dimensions) and error examples demonstrating validation failures. All SCAD generation verified with correct computations.

### Related Files

- examples/yapp/enclosure-dimensions/ — Created example YAML files and generated SCAD files demonstrating all three clearance modes
- pkg/docs/tutorials/yapp-dsl-reference.md — Added documentation for per-side clearance and final dimensions configuration modes


## 2025-11-28

Created SCAD verification script that parses generated SCAD files and verifies padding values and shell dimension computations. All examples verified correct: uniform clearance (97.8×77.8mm), per-side clearance (98.3×76.8mm), final dimensions (100.0×80.0mm), and partial final dimensions (100.0×76.8mm).

### Related Files

- examples/yapp/enclosure-dimensions/verify_scad.py — SCAD file parser and verification script

