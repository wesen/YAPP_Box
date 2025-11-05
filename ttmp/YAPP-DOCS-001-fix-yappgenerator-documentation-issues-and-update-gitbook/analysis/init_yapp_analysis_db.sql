-- YAPP Documentation Analysis Database
-- Created: 2025-11-05
-- Purpose: Track API changes, documentation issues, and example problems

-- ============================================================================
-- README_FIRST: How to use this database
-- ============================================================================
CREATE TABLE IF NOT EXISTS readme_first (
    id INTEGER PRIMARY KEY,
    section TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO readme_first (section, content) VALUES
('Overview', 
'This SQLite database tracks issues and changes in the YAPPgenerator project.
It contains:
1. Git commit analysis for API changes
2. Documentation issues from the GitBook
3. Example file problems
4. API change tracking between versions
5. Recommendations for fixes

Use this database to understand what changed, what''s broken, and what needs updating.'),

('Quick Queries',
'-- Get all API changes in a version:
SELECT * FROM api_changes WHERE version = ''v3.3.8'';

-- Find all broken examples:
SELECT * FROM example_issues WHERE status = ''broken'';

-- Get documentation issues by severity:
SELECT * FROM doc_issues WHERE severity = ''high'' ORDER BY page_name;

-- See what changed in a specific commit:
SELECT * FROM git_commits WHERE commit_hash LIKE ''f9400c4%'';

-- Get recommendations for a specific issue:
SELECT * FROM fix_recommendations WHERE issue_type = ''api_change'';'),

('Table Descriptions',
'git_commits: All commits with metadata
api_changes: Specific API changes extracted from commits
doc_issues: Problems found in GitBook documentation
example_issues: Problems in example .scad files
fix_recommendations: Suggested fixes for issues
version_info: Version metadata and release info'),

('Workflow',
'1. Run analysis scripts to populate the database
2. Query for specific issues or changes
3. Use fix_recommendations to guide updates
4. Mark issues as resolved when fixed
5. Re-run analysis to verify fixes'),

('Maintenance',
'To update the database:
1. Run git log analysis: populate git_commits and api_changes
2. Run example scanner: populate example_issues
3. Run doc analyzer: populate doc_issues
4. Generate recommendations: populate fix_recommendations

Database location: docs/analysis/yapp_analysis.db
Schema: docs/analysis/init_yapp_analysis_db.sql');

-- ============================================================================
-- Git Commit Analysis
-- ============================================================================
CREATE TABLE IF NOT EXISTS git_commits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    commit_hash TEXT NOT NULL UNIQUE,
    commit_date TIMESTAMP,
    author TEXT,
    message TEXT,
    version_tag TEXT,
    files_changed INTEGER,
    insertions INTEGER,
    deletions INTEGER,
    is_api_change BOOLEAN DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_commits_hash ON git_commits(commit_hash);
CREATE INDEX idx_commits_version ON git_commits(version_tag);
CREATE INDEX idx_commits_api ON git_commits(is_api_change);

-- ============================================================================
-- API Changes Tracking
-- ============================================================================
CREATE TABLE IF NOT EXISTS api_changes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    commit_hash TEXT,
    version TEXT NOT NULL,
    change_type TEXT NOT NULL, -- 'parameter_added', 'parameter_removed', 'parameter_renamed', 'behavior_changed', 'bug_fix'
    feature_name TEXT NOT NULL, -- e.g., 'cutoutsLid', 'pcbStands', 'connectorNew'
    description TEXT NOT NULL,
    old_syntax TEXT,
    new_syntax TEXT,
    breaking_change BOOLEAN DEFAULT 0,
    affects_examples BOOLEAN DEFAULT 0,
    affects_docs BOOLEAN DEFAULT 0,
    severity TEXT DEFAULT 'medium', -- 'low', 'medium', 'high', 'critical'
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (commit_hash) REFERENCES git_commits(commit_hash)
);

CREATE INDEX idx_api_version ON api_changes(version);
CREATE INDEX idx_api_feature ON api_changes(feature_name);
CREATE INDEX idx_api_breaking ON api_changes(breaking_change);

-- ============================================================================
-- Documentation Issues
-- ============================================================================
CREATE TABLE IF NOT EXISTS doc_issues (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    page_name TEXT NOT NULL,
    page_url TEXT,
    section TEXT,
    issue_type TEXT NOT NULL, -- 'outdated_syntax', 'wrong_version', 'missing_parameter', 'incorrect_example', 'broken_link'
    description TEXT NOT NULL,
    current_content TEXT,
    suggested_fix TEXT,
    severity TEXT DEFAULT 'medium', -- 'low', 'medium', 'high', 'critical'
    status TEXT DEFAULT 'open', -- 'open', 'in_progress', 'resolved', 'wont_fix'
    related_api_change_id INTEGER,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (related_api_change_id) REFERENCES api_changes(id)
);

CREATE INDEX idx_doc_page ON doc_issues(page_name);
CREATE INDEX idx_doc_type ON doc_issues(issue_type);
CREATE INDEX idx_doc_severity ON doc_issues(severity);
CREATE INDEX idx_doc_status ON doc_issues(status);

-- ============================================================================
-- Example File Issues
-- ============================================================================
CREATE TABLE IF NOT EXISTS example_issues (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    file_name TEXT NOT NULL,
    file_path TEXT NOT NULL,
    issue_type TEXT NOT NULL, -- 'syntax_error', 'deprecated_api', 'render_error', 'outdated_pattern', 'performance_issue'
    line_number INTEGER,
    description TEXT NOT NULL,
    code_snippet TEXT,
    suggested_fix TEXT,
    severity TEXT DEFAULT 'medium',
    status TEXT DEFAULT 'open',
    tested_version TEXT,
    related_api_change_id INTEGER,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (related_api_change_id) REFERENCES api_changes(id)
);

CREATE INDEX idx_example_file ON example_issues(file_name);
CREATE INDEX idx_example_type ON example_issues(issue_type);
CREATE INDEX idx_example_status ON example_issues(status);

-- ============================================================================
-- Fix Recommendations
-- ============================================================================
CREATE TABLE IF NOT EXISTS fix_recommendations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    issue_type TEXT NOT NULL, -- 'api_change', 'doc_issue', 'example_issue'
    issue_id INTEGER NOT NULL, -- References id from api_changes, doc_issues, or example_issues
    priority INTEGER DEFAULT 5, -- 1-10, 10 being highest
    recommendation TEXT NOT NULL,
    implementation_steps TEXT,
    estimated_effort TEXT, -- 'trivial', 'small', 'medium', 'large'
    dependencies TEXT, -- Other issues that must be fixed first
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_fix_type ON fix_recommendations(issue_type);
CREATE INDEX idx_fix_priority ON fix_recommendations(priority);

-- ============================================================================
-- Version Information
-- ============================================================================
CREATE TABLE IF NOT EXISTS version_info (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    version TEXT NOT NULL UNIQUE,
    release_date DATE,
    commit_hash TEXT,
    major_changes TEXT,
    breaking_changes TEXT,
    deprecated_features TEXT,
    new_features TEXT,
    bug_fixes TEXT,
    docs_version TEXT, -- What version the docs claim to be
    actual_version TEXT, -- What version the code actually is
    version_mismatch BOOLEAN DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (commit_hash) REFERENCES git_commits(commit_hash)
);

CREATE INDEX idx_version ON version_info(version);
CREATE INDEX idx_version_mismatch ON version_info(version_mismatch);

-- ============================================================================
-- Analysis Metadata
-- ============================================================================
CREATE TABLE IF NOT EXISTS analysis_runs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    run_type TEXT NOT NULL, -- 'git_analysis', 'doc_analysis', 'example_analysis', 'full_analysis'
    start_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    end_time TIMESTAMP,
    status TEXT DEFAULT 'running', -- 'running', 'completed', 'failed'
    items_processed INTEGER DEFAULT 0,
    issues_found INTEGER DEFAULT 0,
    errors TEXT,
    notes TEXT
);

CREATE INDEX idx_analysis_type ON analysis_runs(run_type);
CREATE INDEX idx_analysis_status ON analysis_runs(status);

-- ============================================================================
-- Initial Data: Known Version Information
-- ============================================================================
INSERT INTO version_info (version, release_date, commit_hash, docs_version, actual_version, version_mismatch, notes)
VALUES 
('v3.3.8', '2025-10-24', 'f9400c4', 'v3.0', 'v3.3.8', 1, 'Current version. Docs claim v3.0 but code is v3.3.8'),
('v3.3.7', '2025-04-17', '047abd3', 'v3.0', 'v3.3.7', 1, 'Previous version'),
('v3.3', NULL, '45c9d22', 'v3.0', 'v3.3', 1, 'Tagged version v3.3'),
('v3.0', '2023-12-01', NULL, 'v3.0', 'v3.0', 0, 'Major rewrite with new API');

-- ============================================================================
-- Initial Data: Known API Change from v3.3.7 to v3.3.8
-- ============================================================================
INSERT INTO api_changes (commit_hash, version, change_type, feature_name, description, old_syntax, new_syntax, breaking_change, severity, notes)
VALUES 
('f9400c4', 'v3.3.8', 'bug_fix', 'connectorNew', 
 'Fixed pcbGap calculation in connectorNew function. Was always returning 0, now correctly uses pcbGapTmp or calculates based on coordinate system.',
 'pcbGap = (pcbGapTmp == undef ) ? ((theCoordSystem[0]==yappCoordPCB) ? pcb_Thickness : 0) : 0;',
 'pcbGap = (pcbGapTmp == undef ) ? ((theCoordSystem[0]==yappCoordPCB) ? pcb_Thickness : 0) : pcbGapTmp;',
 0, 'low', 'Bug fix, not a breaking change. Improves connector positioning accuracy.');

-- ============================================================================
-- Views for Common Queries
-- ============================================================================

-- View: All high-severity issues
CREATE VIEW IF NOT EXISTS high_severity_issues AS
SELECT 
    'doc_issue' as issue_source,
    id,
    page_name as item_name,
    issue_type,
    description,
    severity,
    status
FROM doc_issues
WHERE severity IN ('high', 'critical')
UNION ALL
SELECT 
    'example_issue' as issue_source,
    id,
    file_name as item_name,
    issue_type,
    description,
    severity,
    status
FROM example_issues
WHERE severity IN ('high', 'critical');

-- View: Breaking changes by version
CREATE VIEW IF NOT EXISTS breaking_changes_by_version AS
SELECT 
    version,
    COUNT(*) as breaking_change_count,
    GROUP_CONCAT(feature_name, ', ') as affected_features
FROM api_changes
WHERE breaking_change = 1
GROUP BY version
ORDER BY version DESC;

-- View: Documentation issues with recommendations
CREATE VIEW IF NOT EXISTS doc_issues_with_fixes AS
SELECT 
    d.id,
    d.page_name,
    d.issue_type,
    d.description,
    d.severity,
    d.status,
    r.recommendation,
    r.priority,
    r.estimated_effort
FROM doc_issues d
LEFT JOIN fix_recommendations r ON r.issue_type = 'doc_issue' AND r.issue_id = d.id
WHERE d.status != 'resolved'
ORDER BY d.severity DESC, r.priority DESC;

-- View: Example issues summary
CREATE VIEW IF NOT EXISTS example_issues_summary AS
SELECT 
    file_name,
    COUNT(*) as issue_count,
    SUM(CASE WHEN severity IN ('high', 'critical') THEN 1 ELSE 0 END) as critical_issues,
    GROUP_CONCAT(DISTINCT issue_type) as issue_types,
    MAX(CASE WHEN status = 'open' THEN 1 ELSE 0 END) as has_open_issues
FROM example_issues
GROUP BY file_name
ORDER BY critical_issues DESC, issue_count DESC;

-- ============================================================================
-- Useful Functions (via SQL queries)
-- ============================================================================

-- To mark an issue as resolved:
-- UPDATE doc_issues SET status = 'resolved', updated_at = CURRENT_TIMESTAMP WHERE id = ?;
-- UPDATE example_issues SET status = 'resolved', updated_at = CURRENT_TIMESTAMP WHERE id = ?;

-- To find all issues related to a specific feature:
-- SELECT * FROM api_changes WHERE feature_name LIKE '%cutouts%';
-- SELECT * FROM doc_issues WHERE description LIKE '%cutouts%';
-- SELECT * FROM example_issues WHERE description LIKE '%cutouts%';

-- To get a version comparison:
-- SELECT a1.feature_name, a1.version as old_version, a2.version as new_version, 
--        a1.description as old_desc, a2.description as new_desc
-- FROM api_changes a1
-- JOIN api_changes a2 ON a1.feature_name = a2.feature_name
-- WHERE a1.version = 'v3.3.7' AND a2.version = 'v3.3.8';

