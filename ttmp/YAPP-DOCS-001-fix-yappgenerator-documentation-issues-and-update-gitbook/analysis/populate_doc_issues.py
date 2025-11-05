#!/usr/bin/env python3
"""
Populate documentation issues from known problems
Based on yapp-llm-guidelines.md analysis
"""

import sqlite3
from pathlib import Path

DB_PATH = Path(__file__).parent / "yapp_analysis.db"

def get_db():
    return sqlite3.connect(DB_PATH)

def populate_known_issues():
    """Add known documentation issues from the LLM guidelines"""
    conn = get_db()
    cursor = conn.cursor()
    
    # Known issues from yapp-llm-guidelines.md
    issues = [
        {
            'page_name': 'Version Information',
            'page_url': 'https://mrwheel-docs.gitbook.io/yappgenerator_en/',
            'section': 'Header',
            'issue_type': 'wrong_version',
            'description': 'Documentation claims to be for v3.0 (from Februari 2024) but repository is v3.3.8 (2025-10-24)',
            'current_content': 'This documentation applies to v3.0 (from Februari 2024)',
            'suggested_fix': 'Update to: This documentation applies to v3.3.8 (from October 2024). Note: Some examples may reference older versions.',
            'severity': 'critical',
        },
        {
            'page_name': 'Cutouts',
            'page_url': 'https://mrwheel-docs.gitbook.io/yappgenerator_en/cutouts',
            'section': 'yappCenter usage',
            'issue_type': 'incomplete_documentation',
            'description': 'Does not clearly explain when to use yappCenter vs yappOrigin for circular cutouts',
            'current_content': 'Limited explanation of origin modifiers',
            'suggested_fix': 'Add section: "For yappCircle cutouts, use yappCenter to specify the center point directly. Without yappCenter, the position refers to the bounding box corner, requiring manual offset calculations."',
            'severity': 'high',
        },
        {
            'page_name': 'Parameters',
            'page_url': 'https://mrwheel-docs.gitbook.io/yappgenerator_en/parameters',
            'section': 'p(n) vs n(a) notation',
            'issue_type': 'incomplete_documentation',
            'description': 'Uses p(n) and n(a) notation but does not clearly explain the difference between positional and named parameters',
            'current_content': 'Shows p(n) and n(a) without clear explanation',
            'suggested_fix': 'Add explanation: "p(n) = positional parameters (must be in order, use undef to skip). n(a) = named parameters (flags/enums, can appear anywhere after positional parameters)."',
            'severity': 'high',
        },
        {
            'page_name': 'Default Values',
            'page_url': 'https://mrwheel-docs.gitbook.io/yappgenerator_en/',
            'section': 'Various',
            'issue_type': 'missing_information',
            'description': 'Many parameter descriptions do not state default values',
            'current_content': 'Parameter descriptions without defaults',
            'suggested_fix': 'Add default values for all optional parameters. Example: standoffHeight (default: 1.0mm)',
            'severity': 'medium',
        },
        {
            'page_name': 'Coordinate Systems',
            'page_url': 'https://mrwheel-docs.gitbook.io/yappgenerator_en/coordinate-systems',
            'section': 'Diagrams',
            'issue_type': 'incomplete_documentation',
            'description': 'Diagrams do not show all three coordinate systems side-by-side for easy comparison',
            'current_content': 'Individual diagrams for each system',
            'suggested_fix': 'Add comparison diagram showing yappCoordPCB, yappCoordBox, and yappCoordBoxInside with conversion formulas',
            'severity': 'medium',
        },
        {
            'page_name': 'Examples',
            'page_url': 'https://mrwheel-docs.gitbook.io/yappgenerator_en/examples',
            'section': 'Various',
            'issue_type': 'outdated_syntax',
            'description': 'Some examples may use outdated syntax from v3.0 era',
            'current_content': 'Examples potentially from v3.0',
            'suggested_fix': 'Review all examples and update to v3.3.8 syntax. Add version tags to examples.',
            'severity': 'medium',
        },
        {
            'page_name': 'PCB Stands',
            'page_url': 'https://mrwheel-docs.gitbook.io/yappgenerator_en/pcb-stands',
            'section': 'Parameter order',
            'issue_type': 'outdated_syntax',
            'description': 'Parameter order may differ between documentation and actual v3.3.8 implementation',
            'current_content': 'Parameter order from v3.0',
            'suggested_fix': 'Verify parameter order against YAPP_Template_v3.scad and update if needed',
            'severity': 'high',
        },
        {
            'page_name': 'Include Order',
            'page_url': 'https://mrwheel-docs.gitbook.io/yappgenerator_en/',
            'section': 'Getting Started',
            'issue_type': 'missing_information',
            'description': 'Does not explain that include must come BEFORE parameter definitions, and that YAPP evaluates code at include time',
            'current_content': 'No mention of include order requirements',
            'suggested_fix': 'Add critical note: "The include statement must come FIRST in your file, before any parameter definitions. YAPPgenerator executes code at include time, so parameters must be defined after the include."',
            'severity': 'critical',
        },
    ]
    
    for issue in issues:
        cursor.execute("""
            INSERT INTO doc_issues 
            (page_name, page_url, section, issue_type, description, 
             current_content, suggested_fix, severity, status)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'open')
        """, (
            issue['page_name'],
            issue['page_url'],
            issue['section'],
            issue['issue_type'],
            issue['description'],
            issue['current_content'],
            issue['suggested_fix'],
            issue['severity']
        ))
    
    conn.commit()
    
    # Add fix recommendations for the critical issues
    cursor.execute("""
        SELECT id FROM doc_issues WHERE severity = 'critical'
    """)
    critical_issues = cursor.fetchall()
    
    for (issue_id,) in critical_issues:
        cursor.execute("""
            INSERT INTO fix_recommendations 
            (issue_type, issue_id, priority, recommendation, estimated_effort)
            VALUES ('doc_issue', ?, 10, 
                    'Update GitBook documentation to reflect current v3.3.8 API and include critical usage notes',
                    'medium')
        """, (issue_id,))
    
    conn.commit()
    conn.close()
    
    print(f"✓ Added {len(issues)} known documentation issues")
    print(f"✓ Added {len(critical_issues)} fix recommendations for critical issues")

if __name__ == "__main__":
    populate_known_issues()

