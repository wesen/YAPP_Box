---
Title: 'Intern Guide: Understanding the Array Path Expression Bug'
Ticket: YAPP-BUG-001
Status: active
Topics:
    - yapp
    - bug
    - resolver
DocType: playbook
Intent: long-term
Owners: []
RelatedFiles:
    - Path: pkg/resolver/resolver.go
      Note: Main bug location - lookupPath and getPathString functions
    - Path: pkg/resolver/rules/dependency_graph.go
      Note: Rule affected by this bug
    - Path: design-doc/01-bug-analysis-expression-extraction-from-array-paths.md
      Note: Detailed bug analysis
ExternalSources: []
Summary: Step-by-step guide for understanding and reproducing the array path expression extraction bug
LastUpdated: 2025-11-29T12:31:45.804454324-05:00
---

# Intern Guide: Understanding the Array Path Expression Bug

## Purpose

This playbook helps new team members understand the bug where the resolver cannot extract expression strings from paths containing array indices. It provides step-by-step instructions to reproduce the bug, understand its impact, and verify when it's fixed.

## Environment Assumptions

- Go development environment set up
- YAPP codebase cloned and built
- Basic understanding of Go and YAML
- Familiarity with the resolver pipeline (see YAPP-ERROR-FEEDBACK-001 for context)

## Understanding the Bug

### What's Happening?

When a YAML file contains expressions in array elements (like `features.cutouts[0].width`), and those expressions reference missing variables, the resolver fails to extract the expression string for error reporting.

### Why It Matters

Without the expression string, we can't:
- Show which variables are missing
- Display dependency chains
- Provide helpful error messages

### The Root Cause

The `lookupPath()` function in `pkg/resolver/resolver.go` doesn't handle numeric path segments (array indices). When it encounters an array, it immediately returns `false`, assuming arrays won't appear in dotted paths. However, `collectUnresolved()` generates paths like `features.cutouts.0.height` which include numeric segments.

## Commands

### Step 1: Reproduce the Bug

Create a test file that triggers the bug:

```bash
cat > /tmp/test-bug-reproduce.yaml << 'EOF'
project: test
units: mm
yapp_version: 3.0

features:
  cutouts:
    - face: front
      from_face_left: 10
      from_face_top: 10
      width: max(vars.missing_height, 10)
      height: max(vars.missing_width, 20)
      shape: rectangle
EOF
```

### Step 2: Run Resolver and Check Taxonomy

```bash
# Run resolver - should fail with missing variables
go run ./cmd/yappctl resolve -i /tmp/test-bug-reproduce.yaml --show-taxonomy --format json 2>&1 | jq .
```

**Expected output (showing the bug):**
```json
{
  "context": {
    "column": 15,
    "expression": "",           // ❌ Empty - this is the bug!
    "iterations": 16,
    "line": 11,
    "missing_refs": null       // ❌ null - should contain ["vars.missing_height", "vars.missing_width"]
  },
  "path": "features.cutouts.0.height",
  "severity": "error",
  "stage": "expr.dependency.missing",
  "symptom": "dependency_missing"
}
```

**What should happen (after fix):**
```json
{
  "context": {
    "column": 15,
    "expression": "max(vars.missing_width, 20)",  // ✅ Expression extracted
    "iterations": 16,
    "line": 11,
    "missing_refs": ["vars.missing_width"]        // ✅ Missing refs populated
  },
  "path": "features.cutouts.0.height",
  "severity": "error",
  "stage": "expr.dependency.missing",
  "symptom": "dependency_missing"
}
```

### Step 3: Check Rule Output

```bash
# Run without --show-taxonomy to see rule output
go run ./cmd/yappctl resolve -i /tmp/test-bug-reproduce.yaml 2>&1
```

**Current behavior (bug):**
- `DependencyGraphRule` shows a warning because `missing_refs` is empty
- Cannot display dependency chains
- Error message is less helpful

**Expected behavior (after fix):**
- `DependencyGraphRule` shows full dependency chain
- Lists missing variables: `vars.missing_height`, `vars.missing_width`
- Shows suggested resolution order

### Step 4: Inspect the Code

```bash
# Look at the problematic function
grep -A 20 "func lookupPath" pkg/resolver/resolver.go

# Look at how paths are generated
grep -A 30 "func collectUnresolved" pkg/resolver/resolver.go

# See where getPathString is called
grep -B 5 -A 5 "getPathString" pkg/resolver/resolver.go
```

### Step 5: Understand the Flow

1. **Resolution fails** → `collectUnresolved()` finds `features.cutouts.0.height` still contains a string
2. **Path extraction** → `getPathString(state, "features.cutouts.0.height")` is called
3. **Lookup fails** → `lookupPath()` encounters array at `features.cutouts` and bails out
4. **Empty expression** → `getPathString()` returns `("", false)`
5. **Missing deps fail** → `findMissingDependencies("", state)` returns `nil`
6. **Poor error** → Taxonomy created with empty `MissingRefs`

## Exit Criteria

### To Verify the Bug Exists

- [ ] Test file with array expression fails resolution
- [ ] Taxonomy shows `expression: ""` (empty string)
- [ ] Taxonomy shows `missing_refs: null`
- [ ] `DependencyGraphRule` shows warning instead of dependency chain

### To Verify the Bug is Fixed

- [ ] `lookupPath()` handles numeric path segments correctly
- [ ] Test file shows `expression` populated with actual expression string
- [ ] Test file shows `missing_refs` array with variable names
- [ ] `DependencyGraphRule` displays full dependency chain
- [ ] Unit tests pass for array path lookups

## Testing the Fix

Once the fix is implemented, run these tests:

```bash
# Unit tests for lookupPath with array indices
go test ./pkg/resolver/... -v -run TestLookupPath

# Integration test with array expressions
go test ./pkg/resolver/... -v -run TestResolveResult_ArrayExpressions

# Test DependencyGraphRule with array paths
go test ./pkg/resolver/rules/... -v -run TestDependencyGraphRule

# Manual test with real YAML
go run ./cmd/yappctl resolve -i /tmp/test-bug-reproduce.yaml --show-taxonomy --format json | jq '.context.missing_refs'
# Should output: ["vars.missing_height", "vars.missing_width"]
```

## Notes

### Related Work

- **YAPP-ERROR-FEEDBACK-001:** Implemented `DependencyGraphRule` which is affected by this bug
- **Design doc:** See `design-doc/01-bug-analysis-expression-extraction-from-array-paths.md` for detailed analysis

### Workaround

Currently, `DependencyGraphRule` handles the empty `missing_refs` case gracefully by showing a warning. However, users don't get the helpful dependency chain visualization.

### Impact

- **Severity:** Medium - Resolver still correctly identifies failures, but error messages are less helpful
- **Frequency:** Occurs whenever missing variables are referenced in array element expressions
- **User Impact:** Users see generic "missing dependencies" errors instead of specific variable names

### Debugging Tips

If you're debugging this issue:

1. **Add logging** to `lookupPath()` to see where it fails:
   ```go
   log.Printf("lookupPath: path=%s, current type=%T, segment=%s", path, cur, p)
   ```

2. **Check path format** - verify `collectUnresolved()` generates paths with numeric segments

3. **Test lookupPath directly**:
   ```go
   state := map[string]any{
       "features": map[string]any{
           "cutouts": []any{
               map[string]any{"height": "max(vars.x, 10)"},
           },
       },
   }
   v, ok := lookupPath(state, "features.cutouts.0.height")
   // Should return ("max(vars.x, 10)", true) after fix
   ```

### Common Mistakes

- **Don't** assume paths always use bracket notation (`features.cutouts[0].height`) - they use dot notation (`features.cutouts.0.height`)
- **Don't** forget to handle edge cases: negative indices, out-of-bounds, non-numeric segments in arrays
- **Do** test with deeply nested arrays to ensure recursive handling works
