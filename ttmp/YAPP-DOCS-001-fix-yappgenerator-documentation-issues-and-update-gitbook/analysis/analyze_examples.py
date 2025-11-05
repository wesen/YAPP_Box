#!/usr/bin/env python3
"""
YAPP Example Files Analyzer
Analyzes example .scad files for issues and outdated patterns
"""

import sqlite3
import re
from pathlib import Path
import subprocess

DB_PATH = Path(__file__).parent / "yapp_analysis.db"
REPO_PATH = Path(__file__).parent.parent.parent
EXAMPLES_PATH = REPO_PATH / "examples"

def get_db():
    """Get database connection"""
    return sqlite3.connect(DB_PATH)

def analyze_example_file(file_path):
    """Analyze a single example file for issues"""
    issues = []
    
    try:
        content = file_path.read_text()
        lines = content.split('\n')
        
        # Check for version in comments
        version_match = re.search(r'Version\s+(\d+\.\d+)', content)
        file_version = version_match.group(1) if version_match else 'unknown'
        
        # Check include statement
        include_line = None
        for i, line in enumerate(lines, 1):
            if 'include' in line and 'YAPPgenerator' in line:
                include_line = i
                # Check if include comes first (before parameter definitions)
                # Look for parameter definitions before include
                for j in range(i):
                    if re.match(r'^\s*(pcbLength|pcbWidth|wallThickness|paddingFront)\s*=', lines[j]):
                        issues.append({
                            'type': 'outdated_pattern',
                            'line': j + 1,
                            'description': 'Parameter defined before include statement - may cause issues',
                            'snippet': lines[j].strip(),
                            'severity': 'medium'
                        })
                        break
                break
        
        # Check for deprecated patterns
        deprecated_patterns = [
            (r'printSideBySide', 'showSideBySide', 'printSideBySide renamed to showSideBySide'),
            (r'pcbX\s*=', 'pcbLength', 'pcbX deprecated, use pcbLength'),
            (r'pcbY\s*=', 'pcbWidth', 'pcbY deprecated, use pcbWidth'),
        ]
        
        for pattern, replacement, message in deprecated_patterns:
            for i, line in enumerate(lines, 1):
                if re.search(pattern, line):
                    issues.append({
                        'type': 'deprecated_api',
                        'line': i,
                        'description': message,
                        'snippet': line.strip(),
                        'severity': 'high'
                    })
        
        # Check for parameter order in cutouts
        cutout_arrays = ['cutoutsLid', 'cutoutsBase', 'cutoutsFront', 'cutoutsBack', 'cutoutsLeft', 'cutoutsRight']
        for array_name in cutout_arrays:
            pattern = rf'{array_name}\s*=\s*\['
            if re.search(pattern, content):
                # Find the array definition
                in_array = False
                array_start = 0
                for i, line in enumerate(lines, 1):
                    if re.search(pattern, line):
                        in_array = True
                        array_start = i
                    elif in_array and '];' in line:
                        in_array = False
                    elif in_array:
                        # Check for common mistakes
                        # Missing yappCenter for circles
                        if 'yappCircle' in line and 'yappCenter' not in line:
                            # Check if position looks like it might need yappCenter
                            issues.append({
                                'type': 'outdated_pattern',
                                'line': i,
                                'description': f'{array_name}: yappCircle without yappCenter - may need center-based positioning',
                                'snippet': line.strip(),
                                'severity': 'low'
                            })
        
        # Check for calculated values defined before include
        calc_pattern = r'(\w+)\s*=\s*\([^)]+\)\s*/\s*2'  # e.g., x = (pcbWidth - 10) / 2
        for i, line in enumerate(lines, 1):
            if re.search(calc_pattern, line):
                if include_line and i < include_line:
                    issues.append({
                        'type': 'outdated_pattern',
                        'line': i,
                        'description': 'Calculated value defined before include - may cause undefined variable issues',
                        'snippet': line.strip(),
                        'severity': 'high'
                    })
        
        # Try to render the file (just check syntax)
        try:
            result = subprocess.run(
                ['openscad', '-o', '/dev/null', '--check-only', str(file_path)],
                capture_output=True,
                text=True,
                timeout=10
            )
            if result.returncode != 0:
                # Parse errors
                errors = result.stderr
                if 'WARNING' in errors or 'ERROR' in errors:
                    issues.append({
                        'type': 'syntax_error',
                        'line': None,
                        'description': 'OpenSCAD reports warnings or errors',
                        'snippet': errors[:200],
                        'severity': 'high'
                    })
        except (subprocess.TimeoutExpired, FileNotFoundError):
            pass  # OpenSCAD not available or timeout
        
        return file_version, issues
        
    except Exception as e:
        return 'unknown', [{
            'type': 'syntax_error',
            'line': None,
            'description': f'Failed to analyze file: {str(e)}',
            'snippet': '',
            'severity': 'critical'
        }]

def analyze_all_examples():
    """Analyze all example files"""
    conn = get_db()
    cursor = conn.cursor()
    
    # Start analysis run
    cursor.execute("""
        INSERT INTO analysis_runs (run_type, status, items_processed)
        VALUES ('example_analysis', 'running', 0)
    """)
    run_id = cursor.lastrowid
    conn.commit()
    
    try:
        example_files = list(EXAMPLES_PATH.glob("*.scad"))
        total_issues = 0
        
        for file_path in example_files:
            print(f"Analyzing {file_path.name}...")
            file_version, issues = analyze_example_file(file_path)
            
            for issue in issues:
                cursor.execute("""
                    INSERT INTO example_issues 
                    (file_name, file_path, issue_type, line_number, description, 
                     code_snippet, severity, tested_version)
                    VALUES (?, ?, ?, ?, ?, ?, ?, ?)
                """, (
                    file_path.name,
                    str(file_path.relative_to(REPO_PATH)),
                    issue['type'],
                    issue['line'],
                    issue['description'],
                    issue['snippet'],
                    issue['severity'],
                    file_version
                ))
                total_issues += 1
        
        # Update analysis run
        cursor.execute("""
            UPDATE analysis_runs 
            SET status = 'completed', 
                end_time = CURRENT_TIMESTAMP,
                items_processed = ?,
                issues_found = ?
            WHERE id = ?
        """, (len(example_files), total_issues, run_id))
        
        conn.commit()
        print(f"\n✓ Analyzed {len(example_files)} example files, found {total_issues} issues")
        
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

def print_summary():
    """Print summary of example analysis"""
    conn = get_db()
    cursor = conn.cursor()
    
    print("\n" + "="*70)
    print("YAPP Example Analysis Summary")
    print("="*70)
    
    # Total issues by severity
    cursor.execute("""
        SELECT severity, COUNT(*) 
        FROM example_issues 
        GROUP BY severity 
        ORDER BY 
            CASE severity 
                WHEN 'critical' THEN 1 
                WHEN 'high' THEN 2 
                WHEN 'medium' THEN 3 
                WHEN 'low' THEN 4 
            END
    """)
    print("\nIssues by severity:")
    for row in cursor.fetchall():
        print(f"  {row[0]}: {row[1]}")
    
    # Issues by type
    cursor.execute("""
        SELECT issue_type, COUNT(*) 
        FROM example_issues 
        GROUP BY issue_type 
        ORDER BY COUNT(*) DESC
    """)
    print("\nIssues by type:")
    for row in cursor.fetchall():
        print(f"  {row[0]}: {row[1]}")
    
    # Files with most issues
    cursor.execute("""
        SELECT file_name, COUNT(*) as issue_count
        FROM example_issues 
        GROUP BY file_name 
        ORDER BY issue_count DESC 
        LIMIT 5
    """)
    print("\nFiles with most issues:")
    for row in cursor.fetchall():
        print(f"  {row[0]}: {row[1]} issues")
    
    # Sample issues
    cursor.execute("""
        SELECT file_name, issue_type, description 
        FROM example_issues 
        WHERE severity IN ('critical', 'high')
        LIMIT 5
    """)
    print("\nSample high-severity issues:")
    for row in cursor.fetchall():
        print(f"  [{row[0]}] {row[1]}: {row[2]}")
    
    conn.close()

if __name__ == "__main__":
    print("Analyzing YAPP example files...")
    analyze_all_examples()
    print_summary()

