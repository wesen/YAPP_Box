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

