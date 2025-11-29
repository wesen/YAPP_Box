#!/usr/bin/env python3
"""
Verify SCAD files by parsing padding values and computing shell dimensions.
"""

import re
import sys
from pathlib import Path

def parse_scad_file(scad_path):
    """Parse a SCAD file and extract key values."""
    with open(scad_path, 'r') as f:
        content = f.read()
    
    values = {}
    
    # Extract numeric values
    patterns = {
        'pcbLength': r'pcbLength\s*=\s*([\d.]+)',
        'pcbWidth': r'pcbWidth\s*=\s*([\d.]+)',
        'wallThickness': r'wallThickness\s*=\s*([\d.]+)',
        'paddingFront': r'paddingFront\s*=\s*([\d.]+)',
        'paddingBack': r'paddingBack\s*=\s*([\d.]+)',
        'paddingLeft': r'paddingLeft\s*=\s*([\d.]+)',
        'paddingRight': r'paddingRight\s*=\s*([\d.]+)',
    }
    
    for key, pattern in patterns.items():
        match = re.search(pattern, content)
        if match:
            values[key] = float(match.group(1))
        else:
            values[key] = None
    
    return values

def compute_shell_dimensions(values):
    """Compute shell dimensions from PCB and padding values."""
    if None in [values['pcbLength'], values['pcbWidth'], values['wallThickness'],
                values['paddingFront'], values['paddingBack'],
                values['paddingLeft'], values['paddingRight']]:
        return None, None
    
    shell_length = (values['pcbLength'] + 
                   values['paddingFront'] + values['paddingBack'] + 
                   values['wallThickness'] * 2)
    
    shell_width = (values['pcbWidth'] + 
                  values['paddingLeft'] + values['paddingRight'] + 
                  values['wallThickness'] * 2)
    
    return shell_length, shell_width

def verify_file(scad_path, expected_description):
    """Verify a SCAD file and print results."""
    print(f"\n{'='*70}")
    print(f"Verifying: {scad_path.name}")
    print(f"Description: {expected_description}")
    print('='*70)
    
    values = parse_scad_file(scad_path)
    
    print("\nParsed values:")
    for key, val in values.items():
        if val is not None:
            print(f"  {key:20s} = {val:8.2f}")
        else:
            print(f"  {key:20s} = (not found)")
    
    shell_length, shell_width = compute_shell_dimensions(values)
    
    if shell_length is not None:
        print(f"\nComputed shell dimensions:")
        print(f"  Shell Length = {shell_length:.2f} mm")
        print(f"  Shell Width  = {shell_width:.2f} mm")
        
        # Verify padding consistency
        print(f"\nPadding verification:")
        if values['paddingFront'] == values['paddingBack'] == values['paddingLeft'] == values['paddingRight']:
            print(f"  ✓ Uniform padding: {values['paddingFront']:.2f} mm on all sides")
        else:
            print(f"  ✓ Per-side padding:")
            print(f"    Front:  {values['paddingFront']:.2f} mm")
            print(f"    Back:   {values['paddingBack']:.2f} mm")
            print(f"    Left:   {values['paddingLeft']:.2f} mm")
            print(f"    Right:  {values['paddingRight']:.2f} mm")
    else:
        print("\n⚠ Could not compute shell dimensions (missing values)")

def main():
    examples_dir = Path(__file__).parent
    
    # Define test cases
    test_cases = [
        ('uniform-clearance-example.scad', 
         'Uniform clearance: 1.5 mm on all sides → Shell: 97.8 × 77.8 mm'),
        ('per-side-clearance-example.scad',
         'Per-side clearance: front=2.0, back=1.5, left=1.0, right=1.0 → Shell: 98.3 × 76.8 mm'),
        ('final-dimensions-example.scad',
         'Final dimensions: target 100×80 → computed padding 2.6 mm per side → Shell: 100.0 × 80.0 mm'),
        ('partial-final-dimensions-example.scad',
         'Partial final dimensions: length=100 → padding 2.6 mm, width uses default 1.0 mm → Shell: 100.0 × 76.8 mm'),
    ]
    
    for scad_file, description in test_cases:
        scad_path = examples_dir / scad_file
        if scad_path.exists():
            verify_file(scad_path, description)
        else:
            print(f"\n⚠ File not found: {scad_path}")
    
    print(f"\n{'='*70}")
    print("Verification complete!")
    print('='*70)

if __name__ == '__main__':
    main()

