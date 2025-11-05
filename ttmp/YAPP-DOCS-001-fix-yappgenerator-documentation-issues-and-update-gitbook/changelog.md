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

