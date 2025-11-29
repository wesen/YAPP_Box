---
Title: 'CLI Test Examples: DependencyGraphRule with Array Paths'
Ticket: YAPP-BUG-001
Status: active
Topics:
    - yapp
    - bug
    - resolver
    - testing
DocType: playbook
Intent: long-term
Owners: []
RelatedFiles:
    - Path: pkg/resolver/rules/dependency_graph.go
      Note: Rule being tested
    - Path: pkg/resolver/resolver.go
      Note: Fixed lookupPath function
ExternalSources: []
Summary: CLI test examples demonstrating DependencyGraphRule working correctly with array paths after bug fix
LastUpdated: 2025-11-29T12:45:00-05:00
---

# CLI Test Examples: DependencyGraphRule with Array Paths

## Purpose

This playbook provides ready-to-use YAML test files and CLI commands to verify that DependencyGraphRule works correctly with array paths after the YAPP-BUG-001 fix.

## Environment Assumptions

- Go development environment set up
- YAPP codebase built (`go build ./cmd/yappctl` or use `go run`)
- `jq` installed (optional, for JSON parsing)

## Verified Working Case

This case has been verified to work correctly after the bug fix:

**Test file:** `test-dependency-array-single.yaml`

```yaml
project: test
units: mm
yapp_version: 3.0

features:
  cutouts:
    - face: front
      from_face_left: 10
      from_face_top: 10
      width: max(vars.missing_height, 10)
      height: 20
      shape: rectangle
```

**Verified Command:**
```bash
go run ./cmd/yappctl resolve -i test-dependency-array-single.yaml --show-taxonomy --format json 2>&1 | jq '.context.missing_refs, .context.expression'
```

**Verified Output:**
```json
[
  "vars.missing_height"
]
"max(vars.missing_height, 10)"
```

**Verified Rule Output:**
```bash
go run ./cmd/yappctl resolve -i test-dependency-array-single.yaml 2>&1
```

Should show DependencyGraphRule with:
- Missing variables: `vars.missing_height`
- Error location: `features.cutouts.0.width`
- Expression: `max(vars.missing_height, 10)`
- Dependency chain: `vars.missing_height → features.cutouts.0.width`
- Suggested resolution order

**Unit Test Coverage:**
- `TestDependencyGraphRule_Render_ArrayPath_SingleMissingVar` in `pkg/resolver/rules/dependency_graph_test.go`
- `TestArrayExpressionLookupPath` in `pkg/resolver/resolver_test.go`

## Test Examples

### Example 1: Single Missing Variable in Array Element

**Test file:** `test-dependency-array-single.yaml`

```yaml
project: test
units: mm
yapp_version: 3.0

features:
  cutouts:
    - face: front
      from_face_left: 10
      from_face_top: 10
      width: max(vars.missing_height, 10)
      height: 20
      shape: rectangle
```

**Commands:**

```bash
# Create test file
cat > test-dependency-array-single.yaml << 'EOF'
project: test
units: mm
yapp_version: 3.0

features:
  cutouts:
    - face: front
      from_face_left: 10
      from_face_top: 10
      width: max(vars.missing_height, 10)
      height: 20
      shape: rectangle
EOF

# Test with taxonomy output
go run ./cmd/yappctl resolve -i test-dependency-array-single.yaml --show-taxonomy --format json 2>&1 | jq '.context.missing_refs'

# Expected output: ["vars.missing_height"]

# Test with rule output
go run ./cmd/yappctl resolve -i test-dependency-array-single.yaml 2>&1 | grep -A 20 "Missing variables"
```

**Expected Results:**
- Taxonomy shows `missing_refs: ["vars.missing_height"]`
- Taxonomy shows `expression: "max(vars.missing_height, 10)"`
- DependencyGraphRule displays dependency chain
- VarsScaffoldRule shows YAML scaffold

### Example 2: Multiple Missing Variables in Array Elements

**Test file:** `test-dependency-array-multiple.yaml`

```yaml
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
```

**Commands:**

```bash
# Create test file
cat > test-dependency-array-multiple.yaml << 'EOF'
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

# Test with taxonomy output
go run ./cmd/yappctl resolve -i test-dependency-array-multiple.yaml --show-taxonomy --format json 2>&1 | jq '.context'

# Test with rule output
go run ./cmd/yappctl resolve -i test-dependency-array-multiple.yaml 2>&1
```

**Expected Results:**
- Taxonomy shows `missing_refs` containing both variables
- DependencyGraphRule shows both variables in dependency chain
- Suggested resolution order lists both variables

### Example 3: Nested Arrays with Missing Variables

**Test file:** `test-dependency-array-nested.yaml`

```yaml
project: test
units: mm
yapp_version: 3.0

features:
  connectors:
    - positions:
        - x: max(vars.connector_x, 10)
          y: max(vars.connector_y, 20)
      stand_height: 5
      screw_d: 3
```

**Commands:**

```bash
# Create test file
cat > test-dependency-array-nested.yaml << 'EOF'
project: test
units: mm
yapp_version: 3.0

features:
  connectors:
    - positions:
        - x: max(vars.connector_x, 10)
          y: max(vars.connector_y, 20)
      stand_height: 5
      screw_d: 3
EOF

# Test with taxonomy output
go run ./cmd/yappctl resolve -i test-dependency-array-nested.yaml --show-taxonomy --format json 2>&1 | jq '.context.missing_refs'

# Test with rule output
go run ./cmd/yappctl resolve -i test-dependency-array-nested.yaml 2>&1
```

**Expected Results:**
- Taxonomy correctly extracts expressions from nested array paths
- DependencyGraphRule shows dependency chain for nested paths
- Path format: `features.connectors.0.positions.0.x` or similar

### Example 4: Multiple Array Elements with Different Missing Variables

**Test file:** `test-dependency-array-multiple-elements.yaml`

```yaml
project: test
units: mm
yapp_version: 3.0

features:
  cutouts:
    - face: front
      width: max(vars.cutout1_width, 10)
      height: 20
      shape: rectangle
    - face: back
      width: max(vars.cutout2_width, 15)
      height: 25
      shape: circle
```

**Commands:**

```bash
# Create test file
cat > test-dependency-array-multiple-elements.yaml << 'EOF'
project: test
units: mm
yapp_version: 3.0

features:
  cutouts:
    - face: front
      width: max(vars.cutout1_width, 10)
      height: 20
      shape: rectangle
    - face: back
      width: max(vars.cutout2_width, 15)
      height: 25
      shape: circle
EOF

# Test with taxonomy output
go run ./cmd/yappctl resolve -i test-dependency-array-multiple-elements.yaml --show-taxonomy --format json 2>&1 | jq '.context.missing_refs'

# Test with rule output
go run ./cmd/yappctl resolve -i test-dependency-array-multiple-elements.yaml 2>&1
```

**Expected Results:**
- First unresolved path shows its missing variables
- DependencyGraphRule shows dependency chain for first error
- Note: Currently only first unresolved path is reported (see resolver.go TODO)

## Verification Checklist

For each test example, verify:

- [ ] Taxonomy shows `expression` field populated (not empty)
- [ ] Taxonomy shows `missing_refs` array populated (not null)
- [ ] `missing_refs` contains correct variable names
- [ ] DependencyGraphRule displays dependency chain
- [ ] DependencyGraphRule shows suggested resolution order
- [ ] VarsScaffoldRule shows YAML scaffold (if applicable)
- [ ] Error message is helpful and actionable

## Exit Criteria

All test examples should:
1. Successfully extract expression strings from array paths
2. Populate `missing_refs` correctly
3. Display helpful dependency chains via DependencyGraphRule
4. Show suggested resolution order

## Notes

### Path Format

The resolver uses dot notation with numeric segments for array indices:
- `features.cutouts.0.width` (not `features.cutouts[0].width`)
- `features.connectors.0.positions.0.x` (nested arrays)

### Current Limitations

- Only the first unresolved path is reported (see `resolver.go` line 83 TODO)
- Multiple unresolved paths in different array elements will only show the first one

### Related Work

- **YAPP-ERROR-FEEDBACK-001:** DependencyGraphRule implementation
- **YAPP-BUG-001:** Bug fix for array path expression extraction

### Unit Test Coverage

The following unit tests verify this functionality:

- `TestDependencyGraphRule_Render_ArrayPath_SingleMissingVar` - Tests DependencyGraphRule with exact verified case
- `TestArrayExpressionLookupPath` - Tests lookupPath() handles array indices correctly
- `TestArrayExpressionMissingDependencies` - Tests full resolution pipeline with array expressions

Run tests with:
```bash
go test ./pkg/resolver/rules/... -v -run TestDependencyGraphRule
go test ./pkg/resolver/... -v -run TestArrayExpression
```

