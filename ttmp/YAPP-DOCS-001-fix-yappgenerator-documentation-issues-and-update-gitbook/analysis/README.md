# YAPP Documentation Analysis System

**Purpose**: Track API changes, documentation issues, and example problems in YAPPgenerator

**Database**: `yapp_analysis.db` (SQLite)  
**Created**: 2025-11-05  
**YAPPgenerator Version**: v3.3.8

---

## Quick Start

### View Database Contents

```bash
# Show README and usage guide
./query_db.py readme

# Show summary of all issues
./query_db.py summary

# Show all documentation issues
./query_db.py docs

# Show critical documentation issues only
./query_db.py docs critical

# Show API changes
./query_db.py api

# Show API changes for a specific version
./query_db.py api v3.3.8

# Show fix recommendations
./query_db.py recs

# Show high-priority recommendations (priority >= 8)
./query_db.py recs 8

# Show version information
./query_db.py versions

# Run custom SQL query
./query_db.py query "SELECT * FROM doc_issues WHERE severity = 'critical'"
```

### Direct SQLite Access

```bash
sqlite3 yapp_analysis.db

# Inside sqlite3:
.tables                    # List all tables
.schema doc_issues         # Show table schema
SELECT * FROM readme_first;  # Read the README
```

---

## Database Structure

### Core Tables

| Table | Purpose | Records |
|-------|---------|---------|
| `readme_first` | Usage guide and documentation | 5 sections |
| `git_commits` | All analyzed git commits | 100 commits |
| `api_changes` | Specific API changes by version | 9 changes |
| `doc_issues` | Problems in GitBook documentation | 8 issues |
| `example_issues` | Problems in example .scad files | 0 issues |
| `fix_recommendations` | Suggested fixes for issues | 2 recommendations |
| `version_info` | Version metadata and mismatches | 4 versions |
| `analysis_runs` | Metadata about analysis runs | 3 runs |

### Views (Pre-built Queries)

| View | Purpose |
|------|---------|
| `high_severity_issues` | All critical/high severity issues |
| `breaking_changes_by_version` | Breaking changes grouped by version |
| `doc_issues_with_fixes` | Doc issues with recommendations |
| `example_issues_summary` | Summary of example file problems |

---

## Current Status

### Git Analysis
- **100 commits** analyzed
- **43 commits** with potential API changes
- **9 API changes** tracked
- **0 breaking changes** found

### Documentation Issues
- **8 open issues** identified
  - **2 critical**: Version mismatch, missing include order docs
  - **3 high**: yappCenter usage, parameter notation, PCB stands syntax
  - **3 medium**: Coordinate diagrams, default values, example syntax

### Example Files
- **44 example files** analyzed
- **0 issues** found (examples are well-maintained!)

### Recommendations
- **2 fix recommendations** generated for critical issues

---

## Key Findings

### 1. Version Mismatch (CRITICAL)

**Problem**: GitBook documentation claims to be for v3.0 (February 2024) but the repository is at v3.3.8 (October 2024).

**Impact**: 
- Users may follow outdated syntax
- Parameter order may differ
- New features not documented

**Recommendation**: Update GitBook to v3.3.8 and add version tags to all examples.

### 2. Include Order Not Documented (CRITICAL)

**Problem**: Documentation doesn't explain that `include <YAPPgenerator_v3.scad>` must come BEFORE parameter definitions.

**Impact**:
- Variables defined before include get overwritten
- Causes "undefined variable" errors
- Major source of confusion for new users

**Recommendation**: Add prominent note in "Getting Started" section explaining include order requirements.

### 3. yappCenter Usage Unclear (HIGH)

**Problem**: Documentation doesn't clearly explain when to use `yappCenter` vs `yappOrigin` for circular cutouts.

**Impact**:
- Users manually calculate offsets unnecessarily
- Confusion about coordinate positioning
- More complex code

**Recommendation**: Add clear examples showing both approaches and when to use each.

---

## Analysis Scripts

### `analyze_git_history.py`
Analyzes git commits for API changes.

**Usage**:
```bash
python3 analyze_git_history.py
```

**What it does**:
- Scans last 100 commits
- Identifies commits with API changes
- Extracts function signature changes
- Tracks parameter modifications
- Populates `git_commits` and `api_changes` tables

### `analyze_examples.py`
Analyzes example .scad files for issues.

**Usage**:
```bash
python3 analyze_examples.py
```

**What it does**:
- Scans all files in `examples/` directory
- Checks for deprecated patterns
- Validates include order
- Detects missing `yappCenter` usage
- Tests syntax with OpenSCAD (if available)
- Populates `example_issues` table

### `populate_doc_issues.py`
Populates known documentation issues.

**Usage**:
```bash
python3 populate_doc_issues.py
```

**What it does**:
- Adds known issues from `yapp-llm-guidelines.md`
- Creates fix recommendations
- Populates `doc_issues` and `fix_recommendations` tables

### `query_db.py`
Interactive query interface.

**Usage**: See "Quick Start" section above.

---

## Useful Queries

### Find All Issues for a Feature

```sql
-- Find all issues related to cutouts
SELECT 'API Change' as source, feature_name, description 
FROM api_changes 
WHERE feature_name LIKE '%cutout%'
UNION ALL
SELECT 'Doc Issue' as source, page_name, description 
FROM doc_issues 
WHERE description LIKE '%cutout%'
UNION ALL
SELECT 'Example Issue' as source, file_name, description 
FROM example_issues 
WHERE description LIKE '%cutout%';
```

### Compare Versions

```sql
-- See what changed between versions
SELECT 
    a1.feature_name,
    a1.version as old_version,
    a2.version as new_version,
    a1.description as old_desc,
    a2.description as new_desc
FROM api_changes a1
JOIN api_changes a2 ON a1.feature_name = a2.feature_name
WHERE a1.version = 'v3.3.7' AND a2.version = 'v3.3.8';
```

### Get Actionable Items

```sql
-- Get all high-priority items that need fixing
SELECT 
    d.page_name,
    d.issue_type,
    d.description,
    r.recommendation,
    r.estimated_effort
FROM doc_issues d
JOIN fix_recommendations r ON r.issue_type = 'doc_issue' AND r.issue_id = d.id
WHERE d.status = 'open' AND r.priority >= 8
ORDER BY r.priority DESC;
```

### Track Progress

```sql
-- See analysis run history
SELECT 
    run_type,
    start_time,
    end_time,
    status,
    items_processed,
    issues_found
FROM analysis_runs
ORDER BY start_time DESC;
```

---

## Extending the Database

### Add a New Issue Type

```sql
-- Example: Add a new issue type for performance problems
INSERT INTO doc_issues 
(page_name, issue_type, description, severity, status)
VALUES 
('Rendering', 'performance_issue', 'Documentation does not mention render time optimization', 'low', 'open');
```

### Mark Issue as Resolved

```sql
-- Mark issue as resolved
UPDATE doc_issues 
SET status = 'resolved', updated_at = CURRENT_TIMESTAMP 
WHERE id = 1;
```

### Add a Recommendation

```sql
-- Add a fix recommendation
INSERT INTO fix_recommendations 
(issue_type, issue_id, priority, recommendation, estimated_effort)
VALUES 
('doc_issue', 1, 8, 'Add performance optimization section to docs', 'small');
```

---

## Maintenance

### Re-run Analysis

To update the database with latest information:

```bash
# Re-analyze git history (last 100 commits)
python3 analyze_git_history.py

# Re-analyze examples
python3 analyze_examples.py

# Add any new known issues
python3 populate_doc_issues.py
```

### Backup Database

```bash
# Create backup
cp yapp_analysis.db yapp_analysis.db.backup

# Or export to SQL
sqlite3 yapp_analysis.db .dump > yapp_analysis_backup.sql
```

### Reset Database

```bash
# Delete and recreate
rm yapp_analysis.db
sqlite3 yapp_analysis.db < init_yapp_analysis_db.sql

# Re-run all analysis scripts
python3 analyze_git_history.py
python3 analyze_examples.py
python3 populate_doc_issues.py
```

---

## Integration with docmgr

This analysis system is part of the YAPP_Box project documentation managed by docmgr.

**Ticket**: PICO-TEMP-001  
**Location**: `/home/manuel/.cursor/worktrees/YAPP_Box/RST6w/docs/analysis/`

### Related Documentation
- `../../yapp-llm-guidelines.md` - Comprehensive LLM guide for YAPP
- `../../ttmp/PICO-TEMP-001-*/` - Example enclosure project
- `../../README.md` - Main YAPP_Box README

---

## Future Enhancements

### Planned Features
- [ ] Web scraper for GitBook documentation
- [ ] Automated comparison with online docs
- [ ] CI/CD integration for continuous analysis
- [ ] Export reports to markdown/HTML
- [ ] Track issue resolution over time
- [ ] Link to GitHub issues

### Potential Improvements
- Add more sophisticated API change detection
- Parse OpenSCAD AST for deeper analysis
- Generate automated fix patches
- Create visual diff reports
- Add test coverage tracking

---

## Contributing

To add new analysis capabilities:

1. Create a new Python script in this directory
2. Follow the existing patterns for database access
3. Add appropriate tables/columns to `init_yapp_analysis_db.sql`
4. Update this README with usage instructions
5. Add the script to the maintenance workflow

---

## References

- [YAPPgenerator GitHub](https://github.com/mrWheel/YAPP_Box)
- [YAPPgenerator GitBook](https://mrwheel-docs.gitbook.io/yappgenerator_en/)
- [YAPP LLM Guidelines](../../yapp-llm-guidelines.md)
- [SQLite Documentation](https://www.sqlite.org/docs.html)

---

**Last Updated**: 2025-11-05  
**Database Version**: 1.0  
**Status**: Active

