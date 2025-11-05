# YAPP Analysis Database - Usage Examples

This document shows practical examples of how to use the database to fix YAPP documentation.

## Quick Start

```bash
# Show the built-in README
./query_db.py readme

# Show overall summary
./query_db.py summary

# Show all documentation issues
./query_db.py docs
```

## Finding Issues

### 1. Find All Critical Issues

```bash
./query_db.py docs critical
```

Or with SQL:
```sql
SELECT page_name, description, suggested_fix 
FROM doc_issues 
WHERE severity = 'critical' AND status = 'open';
```

### 2. Find Code Examples with Problems

```bash
./query_db.py query "
SELECT p.title, p.url, c.issue_description, c.code_block
FROM gitbook_pages p 
JOIN gitbook_code_examples c ON c.page_id = p.id 
WHERE c.has_syntax_issue = 1
"
```

### 3. Find Pages That Need Updates

```sql
SELECT DISTINCT p.title, p.url, COUNT(c.id) as issue_count
FROM gitbook_pages p 
JOIN gitbook_code_examples c ON c.page_id = p.id 
WHERE c.has_syntax_issue = 1
GROUP BY p.id
ORDER BY issue_count DESC;
```

## Searching Documentation

### 1. Find All Pages About a Topic

```sql
-- Find pages mentioning cutouts
SELECT title, url, word_count 
FROM gitbook_pages 
WHERE content LIKE '%cutout%' 
ORDER BY word_count DESC;
```

### 2. Search for Specific API Usage

```sql
-- Find all code examples using yappCircle
SELECT p.title, c.code_block 
FROM gitbook_pages p 
JOIN gitbook_code_examples c ON c.page_id = p.id 
WHERE c.code_block LIKE '%yappCircle%';
```

### 3. Find Pages with Most Code Examples

```sql
SELECT title, url, code_blocks 
FROM gitbook_pages 
WHERE code_blocks > 0 
ORDER BY code_blocks DESC 
LIMIT 10;
```

## Analyzing API Changes

### 1. See What Changed in a Version

```bash
./query_db.py api v3.3.8
```

Or:
```sql
SELECT feature_name, change_type, description, breaking_change
FROM api_changes 
WHERE version = 'v3.3.8';
```

### 2. Find Breaking Changes

```sql
SELECT version, feature_name, description, old_syntax, new_syntax
FROM api_changes 
WHERE breaking_change = 1
ORDER BY version DESC;
```

### 3. Compare Two Versions

```sql
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

## Practical Workflows

### Workflow 1: Fix a Specific Page

**Goal**: Fix the "Cutouts" page

```bash
# 1. Get the page content
sqlite3 yapp_analysis.db "SELECT content FROM gitbook_pages WHERE title = 'Cutouts'"

# 2. Find issues on this page
sqlite3 yapp_analysis.db "
SELECT c.code_block, c.issue_description 
FROM gitbook_pages p 
JOIN gitbook_code_examples c ON c.page_id = p.id 
WHERE p.title = 'Cutouts' AND c.has_syntax_issue = 1
"

# 3. Get the URL to update
sqlite3 yapp_analysis.db "SELECT url FROM gitbook_pages WHERE title = 'Cutouts'"

# 4. Mark as resolved after fixing
sqlite3 yapp_analysis.db "
UPDATE doc_issues 
SET status = 'resolved', updated_at = CURRENT_TIMESTAMP 
WHERE page_name = 'Cutouts'
"
```

### Workflow 2: Update All yappCircle Examples

**Goal**: Find and fix all examples missing `yappCenter`

```bash
# 1. Find all affected pages
./query_db.py query "
SELECT DISTINCT p.title, p.url 
FROM gitbook_pages p 
JOIN gitbook_code_examples c ON c.page_id = p.id 
WHERE c.issue_description LIKE '%yappCenter%'
"

# 2. Get the specific code blocks
./query_db.py query "
SELECT p.title, c.code_block 
FROM gitbook_pages p 
JOIN gitbook_code_examples c ON c.page_id = p.id 
WHERE c.issue_description LIKE '%yappCenter%'
"

# 3. Create a fix list
./query_db.py query "
SELECT 
    p.title as page,
    p.url,
    'Add yappCenter flag to circle positioning' as fix,
    COUNT(c.id) as occurrences
FROM gitbook_pages p 
JOIN gitbook_code_examples c ON c.page_id = p.id 
WHERE c.issue_description LIKE '%yappCenter%'
GROUP BY p.id
"
```

### Workflow 3: Generate Documentation Update Report

```bash
# Create a comprehensive report
./query_db.py query "
SELECT 
    'Critical' as priority,
    page_name as page,
    issue_type as type,
    description,
    suggested_fix as fix
FROM doc_issues 
WHERE severity = 'critical' AND status = 'open'
UNION ALL
SELECT 
    'High' as priority,
    page_name,
    issue_type,
    description,
    suggested_fix
FROM doc_issues 
WHERE severity = 'high' AND status = 'open'
ORDER BY priority, page
"
```

## Advanced Queries

### 1. Find Pages Needing Most Work

```sql
-- Combine doc issues and code issues
SELECT 
    p.title,
    p.url,
    COUNT(DISTINCT d.id) as doc_issues,
    COUNT(DISTINCT c.id) as code_issues,
    COUNT(DISTINCT d.id) + COUNT(DISTINCT c.id) as total_issues
FROM gitbook_pages p
LEFT JOIN doc_issues d ON d.page_url = p.url
LEFT JOIN gitbook_code_examples c ON c.page_id = p.id AND c.has_syntax_issue = 1
GROUP BY p.id
HAVING total_issues > 0
ORDER BY total_issues DESC;
```

### 2. Generate Fix Checklist

```sql
-- Create a prioritized fix list
SELECT 
    ROW_NUMBER() OVER (ORDER BY 
        CASE severity 
            WHEN 'critical' THEN 1 
            WHEN 'high' THEN 2 
            WHEN 'medium' THEN 3 
            ELSE 4 
        END
    ) as priority,
    '[ ] ' || page_name || ': ' || description as task,
    severity,
    suggested_fix
FROM doc_issues 
WHERE status = 'open'
ORDER BY priority;
```

### 3. Track Progress Over Time

```sql
-- See analysis runs
SELECT 
    run_type,
    start_time,
    status,
    items_processed,
    issues_found,
    ROUND((JULIANDAY(end_time) - JULIANDAY(start_time)) * 24 * 60, 2) as duration_minutes
FROM analysis_runs
ORDER BY start_time DESC;
```

### 4. Find Related Issues

```sql
-- Find all issues related to a specific feature
SELECT 
    'API Change' as source,
    version as context,
    feature_name as item,
    description
FROM api_changes 
WHERE feature_name LIKE '%cutout%'
UNION ALL
SELECT 
    'Doc Issue' as source,
    page_name as context,
    issue_type as item,
    description
FROM doc_issues 
WHERE description LIKE '%cutout%'
UNION ALL
SELECT 
    'Code Example' as source,
    (SELECT title FROM gitbook_pages WHERE id = page_id) as context,
    'Code Issue' as item,
    issue_description as description
FROM gitbook_code_examples 
WHERE code_block LIKE '%cutout%' AND has_syntax_issue = 1;
```

## Export and Reporting

### 1. Export to CSV

```bash
# Export all issues
sqlite3 -header -csv yapp_analysis.db "
SELECT page_name, issue_type, severity, description, suggested_fix 
FROM doc_issues 
WHERE status = 'open'
" > doc_issues.csv

# Export code issues
sqlite3 -header -csv yapp_analysis.db "
SELECT p.title, c.issue_description, c.code_block 
FROM gitbook_pages p 
JOIN gitbook_code_examples c ON c.page_id = p.id 
WHERE c.has_syntax_issue = 1
" > code_issues.csv
```

### 2. Generate Markdown Report

```bash
# Create a markdown checklist
sqlite3 yapp_analysis.db "
SELECT '- [ ] **' || page_name || '** (' || severity || '): ' || description
FROM doc_issues 
WHERE status = 'open'
ORDER BY 
    CASE severity 
        WHEN 'critical' THEN 1 
        WHEN 'high' THEN 2 
        WHEN 'medium' THEN 3 
        ELSE 4 
    END
" > fix_checklist.md
```

### 3. Generate HTML Report

```bash
# Export to HTML table
sqlite3 -html yapp_analysis.db "
SELECT page_name, issue_type, severity, description 
FROM doc_issues 
WHERE status = 'open'
" > issues_report.html
```

## Updating the Database

### After Fixing Issues

```sql
-- Mark a specific issue as resolved
UPDATE doc_issues 
SET status = 'resolved', updated_at = CURRENT_TIMESTAMP 
WHERE id = 1;

-- Mark all issues for a page as resolved
UPDATE doc_issues 
SET status = 'resolved', updated_at = CURRENT_TIMESTAMP 
WHERE page_name = 'Cutouts';

-- Mark code examples as fixed
UPDATE gitbook_code_examples 
SET has_syntax_issue = 0, issue_description = NULL 
WHERE page_id = (SELECT id FROM gitbook_pages WHERE title = 'Light Tubes');
```

### Re-scraping After Updates

```bash
# Re-scrape GitBook to verify fixes
python3 scrape_gitbook.py

# Check if issues are resolved
./query_db.py query "
SELECT p.title, COUNT(c.id) as remaining_issues
FROM gitbook_pages p 
JOIN gitbook_code_examples c ON c.page_id = p.id 
WHERE c.has_syntax_issue = 1
GROUP BY p.id
"
```

## Tips and Tricks

### 1. Use Views for Common Queries

The database includes pre-built views:

```sql
-- High severity issues across all sources
SELECT * FROM high_severity_issues;

-- Breaking changes by version
SELECT * FROM breaking_changes_by_version;

-- Doc issues with recommendations
SELECT * FROM doc_issues_with_fixes;

-- Example issues summary
SELECT * FROM example_issues_summary;
```

### 2. Full-Text Search

```sql
-- Search across all content
SELECT title, url 
FROM gitbook_pages 
WHERE content LIKE '%coordinate system%' 
   OR content LIKE '%yappCoordPCB%';
```

### 3. Combine Multiple Filters

```sql
-- Find high-priority pages with multiple issues
SELECT 
    p.title,
    p.url,
    p.code_blocks,
    COUNT(c.id) as issues
FROM gitbook_pages p 
JOIN gitbook_code_examples c ON c.page_id = p.id 
WHERE c.has_syntax_issue = 1 
  AND p.code_blocks > 5
GROUP BY p.id
ORDER BY issues DESC;
```

## Integration with Git

### Track Fixes in Git

```bash
# After fixing issues, commit with reference
git add docs/
git commit -m "Fix yappCenter documentation issues

Resolved issues:
- Light Tubes: Added yappCenter examples
- Push Buttons: Clarified center-based positioning

Database query:
SELECT * FROM doc_issues WHERE page_name IN ('Light Tubes', 'Push Buttons')
"

# Update database to mark as resolved
sqlite3 yapp_analysis.db "
UPDATE doc_issues 
SET status = 'resolved', 
    notes = 'Fixed in commit $(git rev-parse HEAD)',
    updated_at = CURRENT_TIMESTAMP 
WHERE page_name IN ('Light Tubes', 'Push Buttons')
"
```

---

## Quick Reference

**Most Useful Commands:**

```bash
# Show README
./query_db.py readme

# Show summary
./query_db.py summary

# Show critical issues
./query_db.py docs critical

# Show API changes
./query_db.py api v3.3.8

# Custom query
./query_db.py query "YOUR SQL HERE"

# Direct SQLite
sqlite3 yapp_analysis.db
```

**Most Useful Queries:**

```sql
-- All open issues
SELECT * FROM doc_issues WHERE status = 'open';

-- Code examples with problems
SELECT * FROM gitbook_code_examples WHERE has_syntax_issue = 1;

-- Pages about a topic
SELECT * FROM gitbook_pages WHERE content LIKE '%topic%';

-- API changes in version
SELECT * FROM api_changes WHERE version = 'v3.3.8';
```

---

**For more information, see:**
- `README.md` - Full documentation
- `./query_db.py readme` - Built-in guide
- SQLite documentation: https://www.sqlite.org/docs.html

