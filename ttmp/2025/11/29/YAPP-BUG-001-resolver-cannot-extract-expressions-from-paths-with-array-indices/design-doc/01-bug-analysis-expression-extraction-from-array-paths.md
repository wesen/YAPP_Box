---
Title: 'Bug Analysis: Expression Extraction from Array Paths'
Ticket: YAPP-BUG-001
Status: active
Topics:
    - yapp
    - bug
    - resolver
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: pkg/resolver/resolver.go
      Note: Contains getPathString and lookupPath functions that fail on array indices
    - Path: pkg/resolver/resolver.go
      Note: Contains collectUnresolved which generates paths with array indices
    - Path: pkg/resolver/rules/dependency_graph.go
      Note: DependencyGraphRule affected by this bug - can't show dependency chains properly
ExternalSources: []
Summary: Resolver cannot extract expression strings from paths containing array indices (e.g., features.cutouts.0.height), causing missing dependency information in error reporting
LastUpdated: 2025-11-29T12:31:44.183821622-05:00
---

# Bug Analysis: Expression Extraction from Array Paths

## Executive Summary

When the resolver encounters unresolved expressions in array elements (e.g., `features.cutouts[0].height`), it fails to extract the expression string from the document state. This prevents proper error reporting, specifically affecting the `DependencyGraphRule` which cannot show dependency chains when missing variables occur in array contexts.

**Impact:** Medium - Error messages are less helpful when variables are missing in array elements, but the resolver still correctly identifies that resolution failed.

**Root Cause:** The `getPathString()` function uses `lookupPath()` which doesn't handle array indices in dotted paths (e.g., `features.cutouts.0.height`).

## Problem Statement

### Current Behavior

When resolution fails for an expression in an array element:

1. `collectUnresolved()` correctly identifies the path as `features.cutouts.0.height`
2. `getPathString()` is called with this path to extract the expression string
3. `lookupPath()` fails because it doesn't handle numeric path segments (array indices)
4. `findMissingDependencies()` receives an empty expression string
5. Missing dependency extraction fails, returning `nil`
6. Error taxonomy is created with empty `MissingRefs` array
7. `DependencyGraphRule` cannot display helpful dependency chains

### Example Failure Case

**Input YAML:**
```yaml
features:
  cutouts:
    - face: front
      width: max(vars.missing_height, 10)
      height: max(vars.missing_width, 20)
```

**Expected:** Error should show `vars.missing_height` and `vars.missing_width` as missing dependencies.

**Actual:** Error shows empty `missing_refs` because expression extraction fails.

### Code Location

The bug is in `pkg/resolver/resolver.go`:

```go
func getPathString(state map[string]any, path string) (string, bool) {
	v, ok := lookupPath(state, path)  // ❌ Fails for paths with array indices
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

func lookupPath(state any, path string) (any, bool) {
	parts := strings.Split(path, ".")
	cur := state
	for _, p := range parts {
		switch t := cur.(type) {
		case map[string]any:
			v, ok := t[p]
			if !ok {
				return nil, false
			}
			cur = v
		case []any:
			// ❌ indexes are not expected in dotted paths here; bail
			return nil, false
		default:
			return nil, false
		}
	}
	return cur, true
}
```

The `lookupPath()` function explicitly bails out when encountering array types, assuming array indices won't appear in dotted paths. However, `collectUnresolved()` generates paths like `features.cutouts.0.height` which include numeric segments representing array indices.

## Proposed Solution

### Option 1: Enhance `lookupPath()` to Handle Array Indices (Recommended)

Modify `lookupPath()` to parse numeric path segments as array indices:

```go
func lookupPath(state any, path string) (any, bool) {
	parts := strings.Split(path, ".")
	cur := state
	for _, p := range parts {
		switch t := cur.(type) {
		case map[string]any:
			v, ok := t[p]
			if !ok {
				return nil, false
			}
			cur = v
		case []any:
			// Try to parse segment as array index
			idx, err := strconv.Atoi(p)
			if err != nil || idx < 0 || idx >= len(t) {
				return nil, false
			}
			cur = t[idx]
		default:
			return nil, false
		}
	}
	return cur, true
}
```

**Pros:**
- Minimal code changes
- Fixes the root cause
- Reusable for other path-based lookups

**Cons:**
- Requires careful handling of edge cases (negative indices, out-of-bounds)

### Option 2: Store Expression Strings During Resolution

Track expression strings during the resolution pass and store them in a map for later lookup:

```go
type ResolveState struct {
	state map[string]any
	expressions map[string]string  // path -> expression string
}

func resolvePass(state map[string]any, ..., expressions map[string]string) {
	// When encountering an expression, store it:
	expressions[path] = exprStr
	// ...
}
```

**Pros:**
- More reliable (doesn't depend on path lookup)
- Can handle complex nested structures

**Cons:**
- More invasive changes
- Requires threading `expressions` map through resolution pipeline

### Option 3: Parse Expression from Unresolved String Value

When `getPathString()` fails, fall back to checking if the unresolved value is already a string (the expression):

```go
func getPathString(state map[string]any, path string) (string, bool) {
	v, ok := lookupPath(state, path)
	if !ok {
		// Fallback: try to find unresolved string value directly
		// This works because unresolved expressions remain as strings
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}
```

**Pros:**
- Simple fallback

**Cons:**
- Doesn't solve the root cause
- May not work if value was transformed during resolution

## Design Decisions

**Recommended Approach:** Option 1 - Enhance `lookupPath()` to handle array indices.

**Rationale:**
- Fixes the root cause with minimal changes
- Makes the path lookup system more robust
- Benefits other code paths that use `lookupPath()`
- Aligns with how `collectUnresolved()` generates paths

## Alternatives Considered

See Proposed Solution section above for alternatives. Option 1 is recommended as the most straightforward fix.

## Implementation Plan

1. **Update `lookupPath()` function** (`pkg/resolver/resolver.go`)
   - Add array index parsing logic
   - Handle edge cases (negative indices, out-of-bounds)
   - Add unit tests

2. **Add unit tests** (`pkg/resolver/resolver_test.go`)
   - Test `lookupPath()` with array indices: `features.cutouts.0.height`
   - Test edge cases: negative indices, out-of-bounds, non-numeric segments
   - Test `getPathString()` with array paths

3. **Integration test** (`pkg/resolver/resolver_test.go`)
   - Test full resolution failure with array element expressions
   - Verify `MissingRefs` are correctly populated in taxonomy

4. **Update `DependencyGraphRule` tests** (`pkg/resolver/rules/dependency_graph_test.go`)
   - Add test case with array path to verify rule works after fix

5. **Manual testing**
   - Test with real YAML files containing array expressions
   - Verify error messages show correct dependency chains

## Open Questions

1. **Path format consistency:** Should we standardize on a path format (e.g., `features.cutouts[0].height` vs `features.cutouts.0.height`)? Currently `collectUnresolved()` uses dot notation with numeric segments.

2. **Performance:** Does parsing numeric segments on every lookup have performance implications? Likely negligible, but worth measuring.

3. **Nested arrays:** How should we handle deeply nested arrays? Current approach should handle this recursively.

## References

- Related ticket: YAPP-ERROR-FEEDBACK-001 (DependencyGraphRule implementation)
- Code location: `pkg/resolver/resolver.go:583-610`
- Affected rule: `pkg/resolver/rules/dependency_graph.go`
