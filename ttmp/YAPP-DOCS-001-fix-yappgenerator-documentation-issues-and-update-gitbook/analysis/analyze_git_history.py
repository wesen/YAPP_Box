#!/usr/bin/env python3
"""
YAPP Git History Analyzer
Analyzes git commits to identify API changes and populate the database
"""

import sqlite3
import subprocess
import re
import json
from datetime import datetime
from pathlib import Path

DB_PATH = Path(__file__).parent / "yapp_analysis.db"
REPO_PATH = Path(__file__).parent.parent.parent

def get_db():
    """Get database connection"""
    return sqlite3.connect(DB_PATH)

def run_git_command(cmd):
    """Run a git command and return output"""
    result = subprocess.run(
        cmd,
        cwd=REPO_PATH,
        capture_output=True,
        text=True,
        shell=isinstance(cmd, str)
    )
    return result.stdout.strip()

def analyze_commits(limit=50):
    """Analyze recent commits for API changes"""
    conn = get_db()
    cursor = conn.cursor()
    
    # Start analysis run
    cursor.execute("""
        INSERT INTO analysis_runs (run_type, status, items_processed)
        VALUES ('git_analysis', 'running', 0)
    """)
    run_id = cursor.lastrowid
    conn.commit()
    
    try:
        # Get commit log
        log_format = "%H|%aI|%an|%s"
        commits_raw = run_git_command(
            f'git log --format="{log_format}" -n {limit} --'
        )
        
        commits = []
        for line in commits_raw.split('\n'):
            if not line:
                continue
            parts = line.split('|', 3)
            if len(parts) == 4:
                commits.append({
                    'hash': parts[0],
                    'date': parts[1],
                    'author': parts[2],
                    'message': parts[3]
                })
        
        issues_found = 0
        
        for commit in commits:
            # Get commit stats
            stats = run_git_command(
                f'git show --stat --format="" {commit["hash"]}'
            )
            
            # Parse stats
            files_changed = 0
            insertions = 0
            deletions = 0
            
            stats_match = re.search(r'(\d+) files? changed', stats)
            if stats_match:
                files_changed = int(stats_match.group(1))
            
            ins_match = re.search(r'(\d+) insertions?', stats)
            if ins_match:
                insertions = int(ins_match.group(1))
            
            del_match = re.search(r'(\d+) deletions?', stats)
            if del_match:
                deletions = int(del_match.group(1))
            
            # Check if it's an API change
            is_api_change = any(keyword in commit['message'].lower() for keyword in [
                'api', 'parameter', 'function', 'breaking', 'change', 'update', 'fix'
            ])
            
            # Extract version tag if present
            version_tag = None
            version_match = re.search(r'v?\d+\.\d+\.?\d*', commit['message'])
            if version_match:
                version_tag = version_match.group(0)
                if not version_tag.startswith('v'):
                    version_tag = 'v' + version_tag
            
            # Insert commit
            cursor.execute("""
                INSERT OR IGNORE INTO git_commits 
                (commit_hash, commit_date, author, message, version_tag, 
                 files_changed, insertions, deletions, is_api_change)
                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
            """, (
                commit['hash'], commit['date'], commit['author'], commit['message'],
                version_tag, files_changed, insertions, deletions, is_api_change
            ))
            
            # Analyze the diff for API changes
            if is_api_change:
                analyze_commit_diff(cursor, commit['hash'], version_tag)
                issues_found += 1
        
        # Update analysis run
        cursor.execute("""
            UPDATE analysis_runs 
            SET status = 'completed', 
                end_time = CURRENT_TIMESTAMP,
                items_processed = ?,
                issues_found = ?
            WHERE id = ?
        """, (len(commits), issues_found, run_id))
        
        conn.commit()
        print(f"✓ Analyzed {len(commits)} commits, found {issues_found} potential API changes")
        
    except Exception as e:
        cursor.execute("""
            UPDATE analysis_runs 
            SET status = 'failed', 
                end_time = CURRENT_TIMESTAMP,
                errors = ?
            WHERE id = ?
        """, (str(e), run_id))
        conn.commit()
        print(f"✗ Error: {e}")
        raise
    finally:
        conn.close()

def analyze_commit_diff(cursor, commit_hash, version):
    """Analyze a specific commit's diff for API changes"""
    # Get the diff
    diff = run_git_command(f'git show {commit_hash} -- YAPPgenerator_v3.scad')
    
    if not diff:
        return
    
    # Look for function signature changes
    function_changes = re.findall(r'[-+]\s*module\s+(\w+)\s*\(', diff)
    for func_name in set(function_changes):
        # Check if it's a real change (both - and +)
        if f'-  module {func_name}' in diff or f'- module {func_name}' in diff:
            cursor.execute("""
                INSERT OR IGNORE INTO api_changes 
                (commit_hash, version, change_type, feature_name, description, severity)
                VALUES (?, ?, 'behavior_changed', ?, ?, 'medium')
            """, (
                commit_hash,
                version or 'unknown',
                func_name,
                f'Function {func_name} was modified in this commit'
            ))
    
    # Look for parameter changes in comments
    param_changes = re.findall(r'[-+]//\s*p\(\d+\)\s*=\s*(.+)', diff)
    if param_changes:
        cursor.execute("""
            INSERT OR IGNORE INTO api_changes 
            (commit_hash, version, change_type, feature_name, description, severity)
            VALUES (?, ?, 'parameter_changed', ?, ?, 'medium')
        """, (
            commit_hash,
            version or 'unknown',
            'parameters',
            f'Parameter documentation changed: {", ".join(set(param_changes[:3]))}'
        ))

def print_summary():
    """Print a summary of the analysis"""
    conn = get_db()
    cursor = conn.cursor()
    
    print("\n" + "="*70)
    print("YAPP Git Analysis Summary")
    print("="*70)
    
    # Total commits
    cursor.execute("SELECT COUNT(*) FROM git_commits")
    total_commits = cursor.fetchone()[0]
    print(f"\nTotal commits analyzed: {total_commits}")
    
    # API changes
    cursor.execute("SELECT COUNT(*) FROM git_commits WHERE is_api_change = 1")
    api_commits = cursor.fetchone()[0]
    print(f"Commits with potential API changes: {api_commits}")
    
    # By version
    cursor.execute("""
        SELECT version_tag, COUNT(*) 
        FROM git_commits 
        WHERE version_tag IS NOT NULL 
        GROUP BY version_tag 
        ORDER BY version_tag DESC
    """)
    print("\nCommits by version:")
    for row in cursor.fetchall():
        print(f"  {row[0]}: {row[1]} commits")
    
    # Recent API changes
    cursor.execute("""
        SELECT feature_name, description, version 
        FROM api_changes 
        ORDER BY id DESC 
        LIMIT 5
    """)
    print("\nRecent API changes:")
    for row in cursor.fetchall():
        print(f"  [{row[2]}] {row[0]}: {row[1]}")
    
    conn.close()

if __name__ == "__main__":
    print("Analyzing YAPP git history...")
    analyze_commits(limit=100)
    print_summary()

