---
Title: "DSL-to-YAPP Translation Analysis"
Ticket: "YAPP-ENCL-DSL-001"
Status: "review"
Topics: ["yapp", "openscad", "translation", "design", "analysis"]
DocType: "analysis"
Intent: "mapping-and-design"
Owners: []
RelatedFiles:
    - Path: ttmp/YAPP-ENCL-DSL-001-create-enclosure-generation-playbook-and-yaml-dsl-for-electronics-projects/reference/01-enclosure-dsl-language-reference.md
      Note: DSL reference (source schema)
    - Path: YAPP_Template_v3.scad
      Note: Canonical YAPP template to generate from
    - Path: YAPPgenerator_v3.scad
      Note: Library used by all examples (main API)
    - Path: examples/YAPP_Demo_cutouts_all_coord_systems_v31.scad
      Note: Coordinate systems & cutouts catalog
    - Path: examples/YAPP_Demo_RealBox_v31.scad
      Note: Real-world example with many features
    - Path: examples/pcbStandTest.scad
      Note: Standoff variants and PCB array usage
    - Path: ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/analysis/yapp_analysis.db
      Note: Docs SQLite DB (pages, examples, API changes)
ExternalSources: []
Summary: "Design for translating the Enclosure DSL into a YAPP OpenSCAD file: mapping of fields to YAPP variables/arrays, file generation structure, coordinate/feature handling, and doc DB queries to quickly look up authoritative references."
LastUpdated: 2025-11-09T00:00:00Z
---

## 1) Purpose and Scope

Design a reliable, deterministic translation from the Enclosure DSL (reference document) to a YAPP-based OpenSCAD file that:
- Compiles in OpenSCAD without manual edits
- Uses idiomatic YAPP structure (variables + feature arrays + `YAPPgenerate();`)
- Maps DSL concepts (PCB, enclosure, features, coordinates, tolerances) to YAPP’s variables and arrays
- Links to canonical examples and documentation for quick lookup

Out of scope (for now):
- Non-core YAPP features not present in the DSL (e.g., pushButtons, labels, displayMounts) unless explicitly mapped later
- Advanced heuristics to infer unknown parameters; we’ll prefer explicit inputs or sensible defaults

## 2) Orientation: YAPP SCAD Anatomy

All example cases follow this anatomy:
- Include the generator at the top, then declare variables, then feature arrays, and finally call `YAPPgenerate();`.
- Key inputs:
  - PCB & box geometry: `pcbLength`, `pcbWidth`, `pcbThickness`, `padding{Front,Back,Left,Right}`, `wallThickness`, `basePlaneThickness`, `lidPlaneThickness`, `baseWallHeight`, `lidWallHeight`, `ridgeHeight`, `ridgeSlack`, `roundRadius`, etc.
  - Feature arrays (vectors of typed entries): `cutouts{Base,Lid,Front,Back,Left,Right}`, `pcbStands`, `connectors`, `snapJoins`, `boxMounts`, `lightTubes` (and more).
  - Coordinate flags per element: `{yappCoordPCB, yappCoordBox, yappCoordBoxInside}`, `{yappOrigin, yappCenter}`, `{yappGlobalOrigin, yappAltOrigin}`.

Useful exemplars in this repo:

```1:18:/home/manuel/code/others/YAPP_Box/examples/YAPP_Demo_cutouts_all_coord_systems_v31.scad
include <../YAPPgenerator_v3.scad>
```

```299:329:/home/manuel/code/others/YAPP_Box/examples/YAPP_Demo_RealBox_v31.scad
//---- This is where the magic happens ----
YAPPgenerate();
```

```55:99:/home/manuel/code/others/YAPP_Box/examples/pcbStandTest.scad
// The Following will be used as the first element in the pbc array
pcbLength           = 100; // Front to back
pcbWidth            = 20; // Side to side
pcbThickness        = 1.6;
...
pcb =
[
  ["Main",              pcbLength,pcbWidth,    0,0,    pcbThickness,  standoffHeight, standoffDiameter, standoffPinDiameter, standoffHoleSlack]
];
```

## 3) Mapping: DSL → YAPP

The DSL document defines structured fields for PCB, enclosure, features, coordinates, and tolerances. Below is the proposed canonical mapping.

- PCB
  - `pcb.length` → `pcbLength`
  - `pcb.width` → `pcbWidth`
  - `pcb.thickness` → `pcbThickness`
  - `pcb.standoffs.(diameter|height|screw_d)` → default standoff-related globals, and entries in `pcbStands` where positions are specified
  - Optional: also populate `pcb` array with a `"Main"` entry so other YAPP features can reference `pcbThickness()` etc.

- Enclosure
  - `enclosure.wall.thickness` → `wallThickness`
  - `enclosure.wall.fillet_radius` → `roundRadius`
  - `enclosure.base.thickness` → `basePlaneThickness`
  - `enclosure.lid.type`:
    - `screws` → generate `connectors` (screw posts) at inferred positions (by default, corners or provided custom positions)
    - `snap` → add `snapJoins` at default positions unless otherwise specified
  - `enclosure.lid.screws.(count|positions|head_clearance)` → emit entries in `connectors` with appropriate diameters/clearances; when `positions: corners`, place posts near each corner respecting padding and wall thickness
  - `enclosure.clearances` (if present) or `enclosure.wall.clearance` → symmetric `paddingFront`, `paddingBack`, `paddingLeft`, `paddingRight` (or allow per-side overrides if DSL provides them)

- Coordinates
  - `coordinates.origin: pcb|box|boxinside` → per-element coordinate flags:
    - `pcb` → `yappCoordPCB`
    - `box` → `yappCoordBox`
    - `boxinside` → `yappCoordBoxInside` (available in newer YAPP)
  - `coordinates.reference_plane: base|lid` → used to derive Z positions where applicable; for side faces we still pass X/Y/Z using the appropriate coordinate flag
  - Centering flags via DSL hints map to `{yappOrigin|yappCenter}` accordingly

- Features
  - `features.holes` (circular) on a face → the corresponding `cutouts{Face}` entries with shape `yappCircle` (map diameter to radius/diameter as needed by YAPP shape signature)
  - `features.cutouts` (rectangular/rounded) → `yappRectangle` or `yappRoundedRect`
  - `features.light_tubes` → `lightTubes` with `[posx, posy, tubeLength, tubeWidth, tubeWall, gapAbovePcb, {yappCircle|yappRectangle}, …]`
  - Each feature entry carries coordinate flags: `{yappCoordPCB|yappCoordBox|yappCoordBoxInside}`, `{yappOrigin|yappCenter}`, and optional `{yappGlobalOrigin|yappAltOrigin}`

- Tolerances
  - `tolerances.perimeter` → base padding defaults if `enclosure.wall.clearance` not specified
  - `tolerances.holes` → expand cutout/holes slightly (apply to radii/width where sensible)
  - `tolerances.mating` → adjust `ridgeSlack` (and, for newer YAPP: `ridgeGap`)

Version caveats:
- Newer YAPP adds parameters like `ridgeGap` (seen in v3.3.7). If the DSL specifies a mating gap separate from ridge slack, map to `ridgeGap`; otherwise fall back to `ridgeSlack`.

## 4) File Generation Structure (proposed)

Generated SCAD should follow YAPP template patterns:
1. Header:
   - `include <../YAPPgenerator_v3.scad>` (adjust relative path if the generated file lives under project subfolders)
2. Box/PCB variables:
   - `pcbLength`, `pcbWidth`, `pcbThickness`, `padding*`, `wallThickness`, `basePlaneThickness`, `lidPlaneThickness`, `baseWallHeight`, `lidWallHeight`, `ridgeHeight`, `ridgeSlack`, `roundRadius`, `printerLayerHeight` (optional)
3. Optional PCB array:
   - `pcb = [["Main", pcbLength, pcbWidth, 0, 0, pcbThickness, standoffHeight, standoffDiameter, standoffPinDiameter, standoffHoleSlack]];`
4. Feature arrays:
   - Initialize empty arrays for all canonical arrays and populate based on DSL features:
     - `cutoutsBase`, `cutoutsLid`, `cutoutsFront`, `cutoutsBack`, `cutoutsLeft`, `cutoutsRight`
     - `pcbStands`, `connectors`, `snapJoins`, `boxMounts`, `lightTubes`
5. Final call:
   - `YAPPgenerate();`

Notes:
- Keep arrays present even if empty; this is consistent with template style and avoids missing-symbol issues when toggling features later.
- When generating screw-based lids, ensure connector positions don’t collide with the ridge or snap joins. Favor conservative offsets.

## 5) Detailed Elements and Signatures

- Cutouts and coordinate systems catalog

```120:154:/home/manuel/code/others/YAPP_Box/examples/YAPP_Demo_cutouts_all_coord_systems_v31.scad
//  yappRectangle | yappCircle | yappRoundedRect | yappCircleWithFlats | yappCircleWithKey | yappPolygon
// Params: [fromBack, fromLeft, width, length, radius, shape, ...flags...]
```

- Real-world example (many arrays and flags in use)

```202:237:/home/manuel/code/others/YAPP_Box/examples/YAPP_Demo_RealBox_v31.scad
// Cutouts and flags include: yappCenter, yappCoordPCB, yappGlobalOrigin/yappAltOrigin
cutoutsBase = [
  [pcbLength/2, pcbWidth/2 ,25, 25, 0, yappPolygon, shapeHexagon, maskHoneycomb, yappCenter, yappCoordPCB]
];
```

- PCB array & standoff examples

```95:103:/home/manuel/code/others/YAPP_Box/examples/pcbStandTest.scad
pcb = [
  ["Main", pcbLength, pcbWidth, 0, 0, pcbThickness, standoffHeight, standoffDiameter, standoffPinDiameter, standoffHoleSlack]
];
```

## 6) Generation Algorithm (pseudocode)

```text
input: dsl (resolved values), target_scad_path

write "include <../YAPPgenerator_v3.scad>"

// --- Core variables
emit pcbLength = dsl.pcb.length
emit pcbWidth  = dsl.pcb.width
emit pcbThickness = dsl.pcb.thickness

emit wallThickness = dsl.enclosure.wall.thickness
emit basePlaneThickness = dsl.enclosure.base.thickness
emit lidPlaneThickness  = dsl.enclosure.lid.thickness? else default
emit baseWallHeight, lidWallHeight = derive from desired clearance above pcb
emit ridgeHeight, ridgeSlack (and ridgeGap if needed)
emit roundRadius = dsl.enclosure.wall.fillet_radius

emit paddingFront/Back/Left/Right = from enclosure.wall.clearance or tolerances.perimeter

// --- Optional PCB array (recommended)
emit pcb = [["Main", pcbLength, pcbWidth, 0, 0, pcbThickness, standoffHeight, standoffDiameter, standoffPinDiameter, standoffHoleSlack]]

// --- Features
init arrays: cutouts*, pcbStands, connectors, snapJoins, boxMounts, lightTubes
for each dsl.features.holes on face F:
  push into cutoutsF: [x, y/z, width/diameter->radius mapping, ..., yappCircle, coord/origin flags]

for each dsl.features.cutouts on face F:
  shape = rectangle | roundedRect -> yappRectangle | yappRoundedRect
  push into cutoutsF with flags

for each dsl.features.light_tubes:
  push into lightTubes with geometry + flags

if lid.type == screws:
  synthesize connectors at corners (or provided positions)
if lid.type == snap:
  add symmetric snapJoins along edges

emit YAPPgenerate()
```

## 7) Queries (YAPP-DOCS-001 SQLite DB)

From `ttmp/YAPP-DOCS-001-.../analysis/`:

```bash
# List key pages
sqlite3 ./yapp_analysis.db "SELECT title FROM gitbook_pages WHERE title IN ('Getting Started','YAPP Box Settings','Cutouts','Coordinate Systems','Standoffs') ORDER BY title"

# Read specific pages
sqlite3 ./yapp_analysis.db "SELECT content FROM gitbook_pages WHERE title = 'YAPP Box Settings'"
sqlite3 ./yapp_analysis.db "SELECT content FROM gitbook_pages WHERE title = 'Cutouts'"
sqlite3 ./yapp_analysis.db "SELECT content FROM gitbook_pages WHERE title = 'Coordinate Systems'"
sqlite3 ./yapp_analysis.db "SELECT content FROM gitbook_pages WHERE title = 'Standoffs'"

# Pull code examples for 'Standoffs'
sqlite3 ./yapp_analysis.db \"SELECT c.code_block FROM gitbook_pages p JOIN gitbook_code_examples c ON c.page_id = p.id WHERE p.title = 'Standoffs'\"

# Track API/version details (e.g., ridgeGap in v3.3.7+)
sqlite3 ./yapp_analysis.db \"SELECT * FROM api_changes WHERE version >= 'v3.3.7' ORDER BY version DESC\"
sqlite3 ./yapp_analysis.db \"SELECT * FROM version_info ORDER BY release_date DESC\"
```

## 8) Quick Links (repo paths)
- Template: `YAPP_Template_v3.scad`
- Library: `YAPPgenerator_v3.scad`
- Coordinates & cutouts demo: `examples/YAPP_Demo_cutouts_all_coord_systems_v31.scad`
- Real box demo: `examples/YAPP_Demo_RealBox_v31.scad`
- PCB stands variants: `examples/pcbStandTest.scad`
- Docs DB & tooling:
  - `ttmp/YAPP-DOCS-001-.../analysis/yapp_analysis.db`
  - `ttmp/YAPP-DOCS-001-.../analysis/query_db.py`

## 9) Open Questions / Assumptions
- Lid “screws” mapping: default to four corner posts sized by `enclosure.lid.screws.*`; allow per-position overrides in the DSL.
- Tolerances distribution: apply `tolerances.holes` uniformly to circular features; review if rectangular cutouts should receive an equivalent offset.
- Coordinate flags per feature: default from `coordinates.origin`, with per-feature override allowed.

## 10) Next Steps
- [ ] Implement a small generator that emits a `.scad` from a resolved DSL instance using the mapping above
- [ ] Validate against:
  - `examples/YAPP_Demo_RealBox_v31.scad` geometry patterns
  - `examples/YAPP_Demo_cutouts_all_coord_systems_v31.scad` coordinate flags
  - Standoff methods shown in `examples/pcbStandTest.scad`
- [ ] Render sanity previews in OpenSCAD; verify lid/base fit and clearances
- [ ] Iterate mapping for edge cases (masks, polygons, alternate origins)

## 11) Optional docmgr workflow
For contributors to add updates consistently:

```bash
docmgr add --ticket YAPP-ENCL-DSL-001 --doc-type analysis --title "DSL-to-YAPP Translation Analysis"
DOC="ttmp/YAPP-ENCL-DSL-001-*/analysis/02-dsl-to-yapp-translation-analysis.md"
docmgr meta update --doc "$DOC" --field Summary --value "Design for translating the Enclosure DSL into a YAPP OpenSCAD file"
```


