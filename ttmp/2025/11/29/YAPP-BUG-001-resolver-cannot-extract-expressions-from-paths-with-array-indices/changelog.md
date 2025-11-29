# Changelog

## 2025-11-29

- Initial workspace created


## 2025-11-29

Created bug report and intern guide for array path expression extraction bug


## 2025-11-29

Fixed lookupPath() function to handle array indices in paths. Enhanced lookupPath() to parse numeric path segments as array indices when encountering []any types. Added integration tests TestArrayExpressionMissingDependencies and TestArrayExpressionLookupPath to verify expression extraction works correctly for array paths. Bug fix verified with manual test case - expressions and missing_refs now correctly extracted from paths like features.cutouts.0.height.

### Related Files

- pkg/resolver/resolver.go — Fixed lookupPath() to handle array indices
- pkg/resolver/resolver_test.go — Added tests for array path expression extraction


## 2025-11-29

Bug fix implemented and verified. All tests passing.


## 2025-11-29

Bug fixed: lookupPath() now handles array indices. DependencyGraphRule now works correctly with array paths.


## 2025-11-29

Added unit test TestDependencyGraphRule_Render_ArrayPath_SingleMissingVar for verified working case. Updated CLI examples playbook with verified output.


## 2025-11-29

All tasks complete: Bug fixed, tests added, CLI examples documented

