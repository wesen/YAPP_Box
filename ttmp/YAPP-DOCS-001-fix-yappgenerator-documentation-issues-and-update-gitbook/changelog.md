# Changelog

## 2025-11-05

- Initial workspace created


## 2025-11-05

Created comprehensive SQLite database system for tracking YAPP documentation issues. Database includes: git_commits (100 commits), api_changes (9 tracked), doc_issues (10 identified), version_info (4 versions), and analysis_runs tables.


## 2025-11-05

Built GitBook scraper (scrape_gitbook.py) that downloads all 22 documentation pages from https://mrwheel-docs.gitbook.io/yappgenerator_en/. Extracts 110 code blocks, analyzes for syntax issues, stores full text and HTML content. Found 8 code examples with potential issues (mainly yappCircle without yappCenter).


## 2025-11-05

Created git history analyzer (analyze_git_history.py) that scans last 100 commits, identifies 43 with potential API changes, extracts function signature changes and parameter modifications. Populates git_commits and api_changes tables.


## 2025-11-05

Populated known documentation issues from yapp-llm-guidelines.md: 2 critical (version mismatch, include order), 3 high (yappCenter usage, parameter notation, PCB stands), 3 medium. Added fix recommendations for critical issues.


## 2025-11-05

Created interactive query interface (query_db.py) with commands: readme, summary, docs, api, recs, versions, and custom SQL queries. Supports tabular output for easy viewing.


## 2025-11-05

Updated README_FIRST table in database with GitBook scraping information and usage examples. Database is now self-documenting.


## 2025-11-05

Created comprehensive fix documentation for critical and high-priority GitBook issues: version mismatch (v3.0→v3.3.8), include order requirements, yappCenter usage, parameter notation (p(n) vs n(a)), and PCB stands verification

### Related Files

- ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/fixes/01-version-mismatch-fix.md
- ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/fixes/02-include-order-fix.md
- ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/fixes/03-yappCenter-usage-fix.md
- ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/fixes/04-parameter-notation-fix.md
- ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/fixes/05-pcb-stands-parameter-order-fix.md
- ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/fixes/README.md


## 2025-11-05

Add reference docs: yappCircle+yappCenter corrections; coordinate systems; default values; PCB stands verification

### Related Files

- /home/manuel/code/others/YAPP_Box/ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/reference/code-examples-yappcircle-yappcenter-corrections.md
- /home/manuel/code/others/YAPP_Box/ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/reference/coordinate-systems-pcb-vs-box-vs-boxinside-comparison.md
- /home/manuel/code/others/YAPP_Box/ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/reference/default-values-for-optional-parameters.md
- /home/manuel/code/others/YAPP_Box/ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/reference/pcb-stands-parameter-order-verification.md


## 2025-11-05

Validated yappCenter corrections with OpenSCAD; fixed analyzer path; rendered two minimal tests

### Related Files

- /home/manuel/code/others/YAPP_Box/ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/analysis/analyze_examples.py
- /home/manuel/code/others/YAPP_Box/ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/various/test_light_tubes_center.scad
- /home/manuel/code/others/YAPP_Box/ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/various/test_push_buttons_center.scad


## 2025-11-05

Built all 44 example .scad files via OpenSCAD preview export; 0 failures


## 2025-11-05

Analyzer updated to treat OpenSCAD warnings as errors; re-ran examples: 23 issues across previews (warnings + nonzero exits).

### Related Files

- /home/manuel/code/others/YAPP_Box/ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/analysis/analyze_examples.py


## 2025-11-05

Examples: Fix ArduinoClone (yappTop→yappLid, correct STL path).

### Related Files

- /home/manuel/code/others/YAPP_Box/examples/YAPP_ArduinoClone_v30.scad


## 2025-11-05

Examples: Fix RealBox v30/v31 STL import paths.

### Related Files

- /home/manuel/code/others/YAPP_Box/examples/YAPP_Demo_RealBox_v30.scad
- /home/manuel/code/others/YAPP_Box/examples/YAPP_Demo_RealBox_v31.scad


## 2025-11-05

Examples: Fix includes (Labels v3, Connector Demo) and disable external import (PoolMonitor) for clean OpenSCAD validation.

### Related Files

- /home/manuel/code/others/YAPP_Box/examples/PoolMonitor_v30.scad
- /home/manuel/code/others/YAPP_Box/examples/YAPP_Connector_Demo.scad
- /home/manuel/code/others/YAPP_Box/examples/YAPP_Demo_Labels_v3.scad


## 2025-11-05

Draft GitBook fix docs: Light Tubes and Push Buttons (yappCenter corrections).

### Related Files

- /home/manuel/code/others/YAPP_Box/ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/fixes/06-gitbook-update-light-tubes.md
- /home/manuel/code/others/YAPP_Box/ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/fixes/07-gitbook-update-push-buttons.md


## 2025-11-05

Expand reference docs: coordinate diagrams + worked example; defaults for cutouts/boxMounts/connectors/labels; PCB stands param order (v3.3.8).

### Related Files

- /home/manuel/code/others/YAPP_Box/ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/reference/coordinate-systems-pcb-vs-box-vs-boxinside-comparison.md
- /home/manuel/code/others/YAPP_Box/ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/reference/default-values-for-optional-parameters.md
- /home/manuel/code/others/YAPP_Box/ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/reference/pcb-stands-parameter-order-verification.md


## 2025-11-05

Examples: RidgeExtDemo_v30 — replace yappLeftOrigin→yappAltOrigin to remove OpenSCAD warning.

### Related Files

- /home/manuel/code/others/YAPP_Box/examples/YAPP_RidgeExtDemo_v30.scad


## 2025-11-05

Add remaining cleanup tasks (GitBook fixes, diagrams, defaults, deprecations, analyzer, re-scrape).


## 2025-11-05

Tasks: Mark PCB Stands order verified; yappLeftOrigin sweep complete; remove early GitBook re-scrape task (superseded by post-maintainer task).


## 2025-11-05

Applied GitBook fix drafts (Light Tubes, Push Buttons), added coordinate diagrams, drafted PCB Stands v3.3.8 order fix, and confirmed yappLeftOrigin replaced in arrays.

### Related Files

- /home/manuel/code/others/YAPP_Box/ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/fixes/08-gitbook-update-pcb-stands-parameter-order.md — Draft for GitBook update to v3.3.8 param order
- /home/manuel/code/others/YAPP_Box/ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/reference/coordinate-systems-pcb-vs-box-vs-boxinside-comparison.md — Added and finalized ASCII diagrams


## 2025-11-05

Fixed RidgeExtDemo_v30: replaced ridgeExtTop expressions and defined demo-safe positions to eliminate OpenSCAD warning.

### Related Files

- /home/manuel/code/others/YAPP_Box/examples/RidgeExtDemo_v30.scad — Resolved ridgeExtTop warning (demo-only fix)


## 2025-11-05

Expanded defaults reference: added Ridge Extensions and PCB Stands defaults; clarified coordinate flags.

### Related Files

- /home/manuel/code/others/YAPP_Box/ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/reference/default-values-for-optional-parameters.md — Expanded defaults coverage


## 2025-11-05

Checked off duplicate coordinate diagrams task (side-by-side ASCII diagrams embedded).

