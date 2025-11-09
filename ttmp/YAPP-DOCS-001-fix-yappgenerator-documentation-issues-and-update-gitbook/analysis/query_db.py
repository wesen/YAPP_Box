#!/usr/bin/env python3
"""
YAPP Analysis Database Query Interface
Convenient interface to query the analysis database
"""

import sqlite3
import sys
from pathlib import Path
from tabulate import tabulate

DB_PATH = Path(__file__).parent / "yapp_analysis.db"

def get_db():
    return sqlite3.connect(DB_PATH)

def show_readme():
    """Show the README_FIRST table"""
    conn = get_db()
    cursor = conn.cursor()
    
    cursor.execute("SELECT section, content FROM readme_first ORDER BY id")
    
    print("\n" + "="*80)
    print("YAPP ANALYSIS DATABASE - README")
    print("="*80 + "\n")
    
    for section, content in cursor.fetchall():
        print(f"\n## {section}\n")
        print(content)
        print()
    
    conn.close()

def show_summary():
    """Show overall summary"""
    conn = get_db()
    cursor = conn.cursor()
    
    print("\n" + "="*80)
    print("YAPP ANALYSIS DATABASE - SUMMARY")
    print("="*80 + "\n")
    
    # Git commits
    cursor.execute("SELECT COUNT(*) FROM git_commits")
    print(f"Git Commits Analyzed: {cursor.fetchone()[0]}")
    
    cursor.execute("SELECT COUNT(*) FROM git_commits WHERE is_api_change = 1")
    print(f"Commits with API Changes: {cursor.fetchone()[0]}")
    
    # API changes
    cursor.execute("SELECT COUNT(*) FROM api_changes")
    print(f"\nAPI Changes Tracked: {cursor.fetchone()[0]}")
    
    cursor.execute("SELECT COUNT(*) FROM api_changes WHERE breaking_change = 1")
    print(f"Breaking Changes: {cursor.fetchone()[0]}")
    
    # Documentation issues
    cursor.execute("SELECT COUNT(*) FROM doc_issues WHERE status = 'open'")
    print(f"\nOpen Documentation Issues: {cursor.fetchone()[0]}")
    
    cursor.execute("""
        SELECT severity, COUNT(*) 
        FROM doc_issues 
        WHERE status = 'open'
        GROUP BY severity
    """)
    print("  By severity:")
    for severity, count in cursor.fetchall():
        print(f"    {severity}: {count}")
    
    # Example issues
    cursor.execute("SELECT COUNT(*) FROM example_issues WHERE status = 'open'")
    print(f"\nOpen Example Issues: {cursor.fetchone()[0]}")
    
    # Recommendations
    cursor.execute("SELECT COUNT(*) FROM fix_recommendations")
    print(f"\nFix Recommendations: {cursor.fetchone()[0]}")
    
    conn.close()

def show_api_changes(version=None):
    """Show API changes, optionally filtered by version"""
    conn = get_db()
    cursor = conn.cursor()
    
    if version:
        cursor.execute("""
            SELECT version, feature_name, change_type, description, severity, breaking_change
            FROM api_changes
            WHERE version = ?
            ORDER BY severity DESC, feature_name
        """, (version,))
        print(f"\nAPI Changes in {version}:")
    else:
        cursor.execute("""
            SELECT version, feature_name, change_type, description, severity, breaking_change
            FROM api_changes
            ORDER BY version DESC, severity DESC, feature_name
        """)
        print("\nAll API Changes:")
    
    rows = cursor.fetchall()
    if rows:
        headers = ['Version', 'Feature', 'Type', 'Description', 'Severity', 'Breaking']
        print(tabulate(rows, headers=headers, tablefmt='grid'))
    else:
        print("  No API changes found")
    
    conn.close()

def show_doc_issues(severity=None):
    """Show documentation issues, optionally filtered by severity"""
    conn = get_db()
    cursor = conn.cursor()
    
    if severity:
        cursor.execute("""
            SELECT page_name, issue_type, description, severity, status
            FROM doc_issues
            WHERE severity = ? AND status = 'open'
            ORDER BY page_name
        """, (severity,))
        print(f"\n{severity.upper()} Severity Documentation Issues:")
    else:
        cursor.execute("""
            SELECT page_name, issue_type, description, severity, status
            FROM doc_issues
            WHERE status = 'open'
            ORDER BY 
                CASE severity 
                    WHEN 'critical' THEN 1 
                    WHEN 'high' THEN 2 
                    WHEN 'medium' THEN 3 
                    WHEN 'low' THEN 4 
                END,
                page_name
        """)
        print("\nOpen Documentation Issues:")
    
    rows = cursor.fetchall()
    if rows:
        headers = ['Page', 'Type', 'Description', 'Severity', 'Status']
        print(tabulate(rows, headers=headers, tablefmt='grid', maxcolwidths=[20, 15, 50, 10, 10]))
    else:
        print("  No documentation issues found")
    
    conn.close()

def show_recommendations(priority_min=None):
    """Show fix recommendations"""
    conn = get_db()
    cursor = conn.cursor()
    
    if priority_min:
        cursor.execute("""
            SELECT r.issue_type, r.priority, r.recommendation, r.estimated_effort
            FROM fix_recommendations r
            WHERE r.priority >= ?
            ORDER BY r.priority DESC
        """, (priority_min,))
        print(f"\nFix Recommendations (priority >= {priority_min}):")
    else:
        cursor.execute("""
            SELECT r.issue_type, r.priority, r.recommendation, r.estimated_effort
            FROM fix_recommendations r
            ORDER BY r.priority DESC
        """)
        print("\nAll Fix Recommendations:")
    
    rows = cursor.fetchall()
    if rows:
        headers = ['Issue Type', 'Priority', 'Recommendation', 'Effort']
        print(tabulate(rows, headers=headers, tablefmt='grid', maxcolwidths=[15, 8, 60, 10]))
    else:
        print("  No recommendations found")
    
    conn.close()

def show_versions():
    """Show version information"""
    conn = get_db()
    cursor = conn.cursor()
    
    cursor.execute("""
        SELECT version, release_date, docs_version, actual_version, version_mismatch, notes
        FROM version_info
        ORDER BY version DESC
    """)
    
    print("\nVersion Information:")
    rows = cursor.fetchall()
    if rows:
        headers = ['Version', 'Release Date', 'Docs Ver', 'Actual Ver', 'Mismatch', 'Notes']
        # Normalize None values to empty strings for tabulate rendering
        normalized = []
        for version, release_date, docs_version, actual_version, version_mismatch, notes in rows:
            normalized.append([
                version or '',
                release_date or '',
                docs_version or '',
                actual_version or '',
                version_mismatch if version_mismatch is not None else '',
                notes or ''
            ])
        print(tabulate(normalized, headers=headers, tablefmt='grid', maxcolwidths=[10, 12, 10, 10, 8, 40]))
    
    conn.close()

def custom_query(query):
    """Execute a custom SQL query"""
    conn = get_db()
    cursor = conn.cursor()
    
    try:
        cursor.execute(query)
        rows = cursor.fetchall()
        
        if rows:
            # Get column names
            headers = [description[0] for description in cursor.description]
            print(tabulate(rows, headers=headers, tablefmt='grid'))
        else:
            print("Query returned no results")
    except Exception as e:
        print(f"Error executing query: {e}")
    finally:
        conn.close()

def main():
    """Main CLI interface"""
    if len(sys.argv) < 2:
        print("Usage: query_db.py <command> [options]")
        print("\nCommands:")
        print("  readme              - Show README_FIRST")
        print("  summary             - Show overall summary")
        print("  api [version]       - Show API changes (optionally for a version)")
        print("  docs [severity]     - Show doc issues (optionally filter by severity)")
        print("  recs [min_priority] - Show recommendations (optionally with min priority)")
        print("  versions            - Show version information")
        print("  query '<SQL>'       - Execute custom SQL query")
        print("\nExamples:")
        print("  query_db.py readme")
        print("  query_db.py api v3.3.8")
        print("  query_db.py docs critical")
        print("  query_db.py recs 8")
        print("  query_db.py query 'SELECT * FROM git_commits LIMIT 5'")
        sys.exit(1)
    
    command = sys.argv[1]
    
    if command == 'readme':
        show_readme()
    elif command == 'summary':
        show_summary()
    elif command == 'api':
        version = sys.argv[2] if len(sys.argv) > 2 else None
        show_api_changes(version)
    elif command == 'docs':
        severity = sys.argv[2] if len(sys.argv) > 2 else None
        show_doc_issues(severity)
    elif command == 'recs':
        priority = int(sys.argv[2]) if len(sys.argv) > 2 else None
        show_recommendations(priority)
    elif command == 'versions':
        show_versions()
    elif command == 'query':
        if len(sys.argv) < 3:
            print("Error: query command requires SQL query as argument")
            sys.exit(1)
        custom_query(sys.argv[2])
    else:
        print(f"Unknown command: {command}")
        sys.exit(1)

if __name__ == "__main__":
    main()

