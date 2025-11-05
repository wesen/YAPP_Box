#!/usr/bin/env python3
"""
YAPP GitBook Scraper
Downloads and analyzes GitBook documentation
"""

import sqlite3
import requests
from bs4 import BeautifulSoup
from pathlib import Path
import time
import re
from urllib.parse import urljoin, urlparse
import json

DB_PATH = Path(__file__).parent / "yapp_analysis.db"
GITBOOK_BASE = "https://mrwheel-docs.gitbook.io/yappgenerator_en/"

def get_db():
    """Get database connection"""
    return sqlite3.connect(DB_PATH)

def init_gitbook_tables():
    """Create tables for GitBook content"""
    conn = get_db()
    cursor = conn.cursor()
    
    # Table for scraped pages
    cursor.execute("""
        CREATE TABLE IF NOT EXISTS gitbook_pages (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            url TEXT NOT NULL UNIQUE,
            title TEXT,
            content TEXT,
            html_content TEXT,
            parent_url TEXT,
            depth INTEGER DEFAULT 0,
            scraped_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            word_count INTEGER,
            code_blocks INTEGER,
            has_examples BOOLEAN DEFAULT 0
        )
    """)
    
    # Table for code examples found in GitBook
    cursor.execute("""
        CREATE TABLE IF NOT EXISTS gitbook_code_examples (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            page_id INTEGER NOT NULL,
            code_block TEXT NOT NULL,
            language TEXT,
            line_number INTEGER,
            context TEXT,
            has_syntax_issue BOOLEAN DEFAULT 0,
            issue_description TEXT,
            FOREIGN KEY (page_id) REFERENCES gitbook_pages(id)
        )
    """)
    
    # Table for links between pages
    cursor.execute("""
        CREATE TABLE IF NOT EXISTS gitbook_links (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            from_page_id INTEGER NOT NULL,
            to_url TEXT NOT NULL,
            link_text TEXT,
            is_internal BOOLEAN DEFAULT 1,
            FOREIGN KEY (from_page_id) REFERENCES gitbook_pages(id)
        )
    """)
    
    # Index for faster queries
    cursor.execute("CREATE INDEX IF NOT EXISTS idx_gitbook_url ON gitbook_pages(url)")
    cursor.execute("CREATE INDEX IF NOT EXISTS idx_gitbook_title ON gitbook_pages(title)")
    cursor.execute("CREATE INDEX IF NOT EXISTS idx_code_page ON gitbook_code_examples(page_id)")
    
    conn.commit()
    conn.close()
    print("✓ GitBook tables initialized")

def scrape_page(url, depth=0, max_depth=3):
    """Scrape a single GitBook page"""
    try:
        print(f"  {'  ' * depth}Scraping: {url}")
        
        # Add delay to be respectful
        time.sleep(0.5)
        
        headers = {
            'User-Agent': 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36'
        }
        response = requests.get(url, headers=headers, timeout=10)
        response.raise_for_status()
        
        soup = BeautifulSoup(response.text, 'html.parser')
        
        # Extract title
        title = None
        title_elem = soup.find('h1')
        if title_elem:
            title = title_elem.get_text().strip()
        else:
            # Try meta title
            meta_title = soup.find('meta', property='og:title')
            if meta_title:
                title = meta_title.get('content', '').strip()
        
        # Extract main content
        content_elem = soup.find('article') or soup.find('main') or soup.find('div', class_='content')
        
        if not content_elem:
            # Try to find the main content area
            content_elem = soup.find('div', {'role': 'main'})
        
        text_content = ""
        html_content = ""
        code_blocks = []
        
        if content_elem:
            # Get text content
            text_content = content_elem.get_text(separator='\n').strip()
            html_content = str(content_elem)
            
            # Extract code blocks
            for code_elem in content_elem.find_all(['code', 'pre']):
                code_text = code_elem.get_text().strip()
                if len(code_text) > 20:  # Ignore very short snippets
                    language = None
                    # Try to detect language from class
                    classes = code_elem.get('class', [])
                    for cls in classes:
                        if 'language-' in cls:
                            language = cls.replace('language-', '')
                        elif cls in ['openscad', 'scad', 'javascript', 'bash', 'python']:
                            language = cls
                    
                    # Get context (surrounding text)
                    context = ""
                    parent = code_elem.find_parent(['p', 'div', 'section'])
                    if parent:
                        context = parent.get_text()[:200]
                    
                    code_blocks.append({
                        'code': code_text,
                        'language': language,
                        'context': context
                    })
        
        # Count words
        word_count = len(text_content.split())
        
        # Extract internal links
        links = []
        for link in soup.find_all('a', href=True):
            href = link['href']
            link_text = link.get_text().strip()
            
            # Resolve relative URLs
            full_url = urljoin(url, href)
            
            # Check if internal (same domain)
            is_internal = urlparse(full_url).netloc == urlparse(GITBOOK_BASE).netloc
            
            if is_internal and full_url.startswith(GITBOOK_BASE):
                links.append({
                    'url': full_url,
                    'text': link_text,
                    'is_internal': True
                })
        
        return {
            'url': url,
            'title': title,
            'content': text_content,
            'html_content': html_content,
            'word_count': word_count,
            'code_blocks': code_blocks,
            'links': links,
            'depth': depth
        }
        
    except Exception as e:
        print(f"    Error scraping {url}: {e}")
        return None

def save_page_to_db(page_data):
    """Save scraped page to database"""
    conn = get_db()
    cursor = conn.cursor()
    
    try:
        # Insert page
        cursor.execute("""
            INSERT OR REPLACE INTO gitbook_pages 
            (url, title, content, html_content, depth, word_count, code_blocks, has_examples)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?)
        """, (
            page_data['url'],
            page_data['title'],
            page_data['content'],
            page_data['html_content'],
            page_data['depth'],
            page_data['word_count'],
            len(page_data['code_blocks']),
            len(page_data['code_blocks']) > 0
        ))
        
        page_id = cursor.lastrowid
        
        # Insert code blocks
        for code_block in page_data['code_blocks']:
            cursor.execute("""
                INSERT INTO gitbook_code_examples 
                (page_id, code_block, language, context)
                VALUES (?, ?, ?, ?)
            """, (
                page_id,
                code_block['code'],
                code_block['language'],
                code_block['context']
            ))
        
        # Insert links
        for link in page_data['links']:
            cursor.execute("""
                INSERT INTO gitbook_links 
                (from_page_id, to_url, link_text, is_internal)
                VALUES (?, ?, ?, ?)
            """, (
                page_id,
                link['url'],
                link['text'],
                link['is_internal']
            ))
        
        conn.commit()
        return page_id
        
    except Exception as e:
        print(f"    Error saving to DB: {e}")
        conn.rollback()
        return None
    finally:
        conn.close()

def scrape_gitbook_recursive(start_url, max_depth=2, visited=None):
    """Recursively scrape GitBook starting from a URL"""
    if visited is None:
        visited = set()
    
    if start_url in visited or len(visited) > 50:  # Limit total pages
        return
    
    visited.add(start_url)
    
    # Scrape current page
    page_data = scrape_page(start_url, depth=len(visited))
    
    if not page_data:
        return
    
    # Save to database
    save_page_to_db(page_data)
    
    # If not at max depth, follow internal links
    if page_data['depth'] < max_depth:
        for link in page_data['links']:
            if link['is_internal'] and link['url'] not in visited:
                # Only follow links that look like documentation pages
                if '/yappgenerator_en/' in link['url']:
                    scrape_gitbook_recursive(link['url'], max_depth, visited)

def analyze_code_examples():
    """Analyze code examples for potential issues"""
    conn = get_db()
    cursor = conn.cursor()
    
    print("\nAnalyzing code examples...")
    
    # Get all code examples
    cursor.execute("""
        SELECT id, code_block, language, page_id
        FROM gitbook_code_examples
    """)
    
    issues_found = 0
    
    for row in cursor.fetchall():
        code_id, code, language, page_id = row
        
        issues = []
        
        # Check for common issues
        if language in ['openscad', 'scad', None]:
            # Check for variables defined before include
            if 'pcbLength' in code or 'pcbWidth' in code:
                lines = code.split('\n')
                include_line = None
                for i, line in enumerate(lines):
                    if 'include' in line and 'YAPP' in line:
                        include_line = i
                        break
                
                if include_line is not None:
                    for i, line in enumerate(lines[:include_line]):
                        if re.match(r'^\s*(pcbLength|pcbWidth|wallThickness)\s*=', line):
                            issues.append('Parameter defined before include statement')
                            break
            
            # Check for deprecated patterns
            if 'printSideBySide' in code:
                issues.append('Uses deprecated printSideBySide (should be showSideBySide)')
            
            # Check for missing yappCenter with circles
            if 'yappCircle' in code and 'yappCenter' not in code:
                issues.append('yappCircle without yappCenter - may need center-based positioning')
        
        # Update database if issues found
        if issues:
            cursor.execute("""
                UPDATE gitbook_code_examples 
                SET has_syntax_issue = 1, issue_description = ?
                WHERE id = ?
            """, ('; '.join(issues), code_id))
            issues_found += 1
    
    conn.commit()
    conn.close()
    
    print(f"✓ Analyzed code examples, found {issues_found} with potential issues")

def generate_doc_issues_from_scrape():
    """Generate documentation issues based on scraped content"""
    conn = get_db()
    cursor = conn.cursor()
    
    print("\nGenerating documentation issues from scraped content...")
    
    # Find pages with problematic code examples
    cursor.execute("""
        SELECT DISTINCT p.url, p.title, c.issue_description
        FROM gitbook_pages p
        JOIN gitbook_code_examples c ON c.page_id = p.id
        WHERE c.has_syntax_issue = 1
    """)
    
    for row in cursor.fetchall():
        url, title, issue_desc = row
        
        # Check if issue already exists
        cursor.execute("""
            SELECT id FROM doc_issues 
            WHERE page_url = ? AND description LIKE ?
        """, (url, f'%{issue_desc[:30]}%'))
        
        if not cursor.fetchone():
            # Add new issue
            cursor.execute("""
                INSERT INTO doc_issues 
                (page_name, page_url, issue_type, description, severity, status)
                VALUES (?, ?, 'incorrect_example', ?, 'high', 'open')
            """, (title or 'Unknown Page', url, f'Code example issue: {issue_desc}'))
    
    conn.commit()
    conn.close()
    
    print("✓ Generated documentation issues from scraped content")

def print_scrape_summary():
    """Print summary of scraped content"""
    conn = get_db()
    cursor = conn.cursor()
    
    print("\n" + "="*70)
    print("GitBook Scrape Summary")
    print("="*70)
    
    # Total pages
    cursor.execute("SELECT COUNT(*) FROM gitbook_pages")
    print(f"\nPages scraped: {cursor.fetchone()[0]}")
    
    # Pages with code
    cursor.execute("SELECT COUNT(*) FROM gitbook_pages WHERE has_examples = 1")
    print(f"Pages with code examples: {cursor.fetchone()[0]}")
    
    # Total code blocks
    cursor.execute("SELECT COUNT(*) FROM gitbook_code_examples")
    print(f"Total code blocks: {cursor.fetchone()[0]}")
    
    # Code blocks with issues
    cursor.execute("SELECT COUNT(*) FROM gitbook_code_examples WHERE has_syntax_issue = 1")
    print(f"Code blocks with issues: {cursor.fetchone()[0]}")
    
    # Top pages by word count
    cursor.execute("""
        SELECT title, word_count, code_blocks 
        FROM gitbook_pages 
        ORDER BY word_count DESC 
        LIMIT 5
    """)
    print("\nTop pages by content:")
    for title, words, code in cursor.fetchall():
        print(f"  {title}: {words} words, {code} code blocks")
    
    # Pages with most code
    cursor.execute("""
        SELECT title, code_blocks 
        FROM gitbook_pages 
        WHERE code_blocks > 0
        ORDER BY code_blocks DESC 
        LIMIT 5
    """)
    print("\nPages with most code examples:")
    for title, code in cursor.fetchall():
        print(f"  {title}: {code} code blocks")
    
    conn.close()

def main():
    """Main scraper function"""
    print("YAPP GitBook Scraper")
    print("="*70)
    
    # Initialize tables
    init_gitbook_tables()
    
    # Start scraping from main page
    print("\nStarting scrape from:", GITBOOK_BASE)
    print("This may take a few minutes...")
    
    # Start analysis run
    conn = get_db()
    cursor = conn.cursor()
    cursor.execute("""
        INSERT INTO analysis_runs (run_type, status)
        VALUES ('gitbook_scrape', 'running')
    """)
    run_id = cursor.lastrowid
    conn.commit()
    conn.close()
    
    try:
        # Scrape the GitBook
        scrape_gitbook_recursive(GITBOOK_BASE, max_depth=2)
        
        # Analyze code examples
        analyze_code_examples()
        
        # Generate issues
        generate_doc_issues_from_scrape()
        
        # Update analysis run
        conn = get_db()
        cursor = conn.cursor()
        cursor.execute("""
            UPDATE analysis_runs 
            SET status = 'completed', end_time = CURRENT_TIMESTAMP
            WHERE id = ?
        """, (run_id,))
        conn.commit()
        conn.close()
        
        # Print summary
        print_scrape_summary()
        
        print("\n✓ GitBook scrape completed successfully!")
        print(f"\nQuery the data with:")
        print(f"  ./query_db.py query \"SELECT title, url FROM gitbook_pages\"")
        print(f"  sqlite3 yapp_analysis.db \"SELECT * FROM gitbook_code_examples WHERE has_syntax_issue = 1\"")
        
    except Exception as e:
        print(f"\n✗ Error during scrape: {e}")
        conn = get_db()
        cursor = conn.cursor()
        cursor.execute("""
            UPDATE analysis_runs 
            SET status = 'failed', end_time = CURRENT_TIMESTAMP, errors = ?
            WHERE id = ?
        """, (str(e), run_id))
        conn.commit()
        conn.close()
        raise

if __name__ == "__main__":
    main()

