---
Title: "YAPP Documentation Study Playbook"
Ticket: "YAPP-ENCL-DSL-001"
Status: "review"
Topics: ["yapp", "openscad", "dsl", "docs", "analysis"]
DocType: "playbook"
Intent: "repeatable-docs-and-examples-study"
Owners: []
RelatedFiles:
    - Path: examples/YAPP_Demo_cutouts_all_coord_systems_v30.scad
      Note: SCAD demo - cutouts & coords
    - Path: examples/YAPP_Demo_lightTubes_v30.scad
      Note: SCAD demo - light tubes
    - Path: examples/YAPP_ViewShapes_v30.scad
      Note: SCAD demo - shapes
    - Path: ttmp/PICO-TEMP-001-raspberry-pi-pico-temperature-monitor-enclosure/pico_temp_monitor_box.scad
      Note: Case study model
    - Path: ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/analysis/init_yapp_analysis_db.sql
      Note: DB schema
    - Path: ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/analysis/query_db.py
      Note: DB query CLI
    - Path: ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/analysis/scrape_gitbook.py
      Note: Scraper script
    - Path: ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/analysis/yapp_analysis.db
      Note: Docs SQLite DB
ExternalSources: []
Summary: "Practical workflow to query YAPP docs DB, browse complex examples, and quickly locate implementation touchpoints in YAPPgenerator."
LastUpdated: 2025-11-08T00:00:00Z
---

## Purpose

Create a repeatable workflow to read and query YAPP documentation, browse examples, and locate implementation details when investigating behavior or planning changes.

## Locations

- Documentation DB and tools (ticket YAPP-DOCS-001):
  - analysis/query_db.py — Query CLI for the SQLite docs database
  - analysis/yapp_analysis.db — SQLite database with scraped docs, issues, API changes
  - analysis/scrape_gitbook.py — Scraper used to populate the DB
  - analysis/init_yapp_analysis_db.sql — Schema and initialization SQL
- Examples (local):
  - examples/ — YAPP demo .scad files (cutouts, light tubes, shapes, etc.)
- Case study (local):
  - ttmp/PICO-TEMP-001-.../*.scad — Real-world box used to validate parameters and workflows

## Prerequisites

- Python 3 available in PATH
- Python package: tabulate (install with: `python3 -m pip install --user tabulate`)
- sqlite3 CLI optional (for ad-hoc queries)
- ripgrep optional (for code search). Install on Debian/Ubuntu: `sudo apt-get install ripgrep`
- OpenSCAD optional for rendering previews. Install from your OS package manager if needed

### Verify prerequisites quickly

```bash
# Python + tabulate (quiet install if missing), sqlite3, ripgrep, OpenSCAD presence (non-fatal)
python3 --version >/dev/null 2>&1 && python3 -m pip install --user -q tabulate && \
sqlite3 --version >/dev/null 2>&1 || echo "sqlite3 not installed" && \
rg --version >/dev/null 2>&1 || echo "rg not installed" && \
command -v openscad >/dev/null 2>&1 || echo "openscad not installed"
```

## Quickstart: Query the Docs DB

From the repo root:

```bash
cd ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/analysis
python3 ./query_db.py readme
python3 ./query_db.py summary
python3 ./query_db.py api            # all API changes
python3 ./query_db.py api v3.3.8     # API changes for specific version
python3 ./query_db.py docs           # open documentation issues
python3 ./query_db.py recs 8         # fix recommendations with priority >= 8
python3 ./query_db.py versions       # YAPP versions tracked
```

For ad-hoc SQL:

```bash
sqlite3 ./yapp_analysis.db '.schema'
sqlite3 ./yapp_analysis.db 'SELECT * FROM api_changes ORDER BY version DESC LIMIT 10'
```

## Read Specific GitBook Pages

The DB contains 22 scraped GitBook pages with full text and code examples. Query by title or search content.

List all available pages:

```bash
sqlite3 ./yapp_analysis.db "SELECT title FROM gitbook_pages WHERE title IS NOT NULL ORDER BY title"
```

Read a specific page (e.g., Cutouts):

```bash
sqlite3 ./yapp_analysis.db "SELECT content FROM gitbook_pages WHERE title = 'Cutouts'"
```

Search for a topic across all pages:

```bash
sqlite3 ./yapp_analysis.db "SELECT title, url FROM gitbook_pages WHERE content LIKE '%yappCenter%'"
```

Get code examples from a page:

```bash
sqlite3 ./yapp_analysis.db "SELECT c.code_block FROM gitbook_pages p JOIN gitbook_code_examples c ON c.page_id = p.id WHERE p.title = 'Standoffs'"
```

### Most Useful Pages (by code examples)

- **Standoffs** (22 code blocks): PCB mounting, pin/hole configuration
- **YAPP Box Settings** (14 blocks): Core parameters (pcbLength, wallThickness, etc.)
- **Getting Started** (11 blocks): Initial setup and first box
- **Cutouts** (varies): Shape types, coordinate systems, masks
- **Coordinate Systems** (varies): yappCoordBox vs yappCoordPCB, origin handling

### Database Schema (Key Tables)

**gitbook_pages** — Scraped documentation pages
- `id`, `title`, `url`, `content`, `html_content`
- `word_count`, `code_blocks`, `has_examples`

**gitbook_code_examples** — Code blocks extracted from pages
- `id`, `page_id` (FK to gitbook_pages)
- `code_block`, `language`, `context`
- `has_syntax_issue`, `issue_description`

**api_changes** — Tracked API changes by version
- `id`, `version`, `commit_hash`
- `change_type` (parameter_added, parameter_removed, behavior_changed, bug_fix)
- `feature_name` (cutoutsLid, pcbStands, connectorNew, etc.)
- `description`, `old_syntax`, `new_syntax`
- `breaking_change`, `severity` (low, medium, high, critical)

**doc_issues** — Documentation problems
- `id`, `page_name`, `page_url`, `section`
- `issue_type` (outdated_syntax, wrong_version, incorrect_example, broken_link)
- `description`, `suggested_fix`
- `severity`, `status` (open, in_progress, resolved)

**version_info** — YAPP version metadata
- `version`, `release_date`, `commit_hash`
- `docs_version`, `actual_version`, `version_mismatch`

Common joins:
```sql
-- Get code examples with page context
SELECT p.title, c.code_block, c.language
FROM gitbook_pages p
JOIN gitbook_code_examples c ON c.page_id = p.id
WHERE p.title = 'Standoffs';

-- Find issues with their related API changes
SELECT d.page_name, d.description, a.version, a.feature_name
FROM doc_issues d
LEFT JOIN api_changes a ON d.related_api_change_id = a.id
WHERE d.status = 'open';
```

## Browse Examples (SCAD)

- SCAD demos (glance by feature):
  - examples/YAPP_Demo_cutouts_all_coord_systems_v30.scad
  - examples/YAPP_Demo_cutouts_all_coord_systems_v31.scad
  - examples/YAPP_Demo_lightTubes_v30.scad
  - examples/YAPP_ViewShapes_v30.scad
Render previews with OpenSCAD:

```bash
openscad -o examples/preview.png examples/YAPP_Demo_cutouts_all_coord_systems_v30.scad
# If openscad is not installed, open the file in the GUI or install it via your package manager.
```

## Find Implementation Touchpoints

Because YAPP is structured as an OpenSCAD library, example files include the library at the repo root:

```scad
include <../YAPPgenerator_v3.scad>
```

To locate relevant parts with ripgrep:

```bash
rg -n '(use|include) <.*YAPP.*' examples/         # find includes of the YAPP library
rg -n 'YAPPgenerate\\(|lightTubes|pcbStands|cutouts|snapJoins|boxMounts' examples/
                                                   # find principal feature arrays and main call
```

Then open the referenced include paths in your local environment or upstream YAPP repository.

### Troubleshooting
- If `openscad` is not installed, skip CLI rendering and open `.scad` in the GUI, or install via your OS package manager.
- If `rg` is missing, use `grep -R` as a fallback.

## When Investigating a Specific Behavior

1) Check DB for documented guidance or known issues:
```bash
python3 ./query_db.py docs high
python3 ./query_db.py api v3.3.8
```
2) Reproduce with the closest SCAD demo in examples/.
3) Reproduce with the real case (PICO-TEMP-001) by adjusting a single parameter.
4) If still unclear, search includes/usages to find the YAPP definitions and read the corresponding functions/modules.

### Reading complex examples quickly
- YAPP generator is invoked via `YAPPgenerate();` at the bottom of examples.
- Feature arrays near the top control behavior:
  - `cutoutsBase`, `cutoutsLid`, `cutoutsFront`, `cutoutsBack`, `cutoutsLeft`, `cutoutsRight`
  - `lightTubes`, `pcbStands`, `snapJoins`, `boxMounts`
- Coordinates and flags:
  - Coordinate system: `yappCoordBox`, `yappCoordPCB`, and alternates like `yappAltOrigin`/`yappGlobalOrigin`
  - Origin handling: `yappOrigin` vs `yappCenter`
- Good study files:
  - Cutouts and coordinate systems: `examples/YAPP_Demo_cutouts_all_coord_systems_v31.scad`
  - Light tubes and lid treatment: `examples/YAPP_Demo_lightTubes_v30.scad`

## Reference

- Coordinate systems: ../../YAPP-DOCS-001-.../reference/coordinate-systems-pcb-vs-box-vs-boxinside-comparison.md
- Stands ordering caveats: ../../YAPP-DOCS-001-.../reference/pcb-stands-parameter-order-verification.md
- Defaults for optional parameters: ../../YAPP-DOCS-001-.../reference/default-values-for-optional-parameters.md

## Exit Criteria

- You can query the DB for API changes, doc issues, and versions.
- You can locate and render a relevant YAPP demo example locally.
- You can map example SCAD parameters/features to their YAPP modules/functions to study behavior further.