---
Title: 2025-11-16 Step 21 - module authoring guide
Ticket: YAPP-MODULE-SYSTEM-001
Status: active
Topics:
    - yapp
    - dsl
    - codegen
    - architecture
DocType: log
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Created comprehensive guide for adding new modules
LastUpdated: 2025-11-15T21:53:30.510204238-05:00
---


# 2025-11-16 Step 21 - module authoring guide

<!-- Log entries in reverse chronological order (newest first) -->

## 2025-11-16 - Created module authoring guide

Wrote comprehensive developer guide for adding new YAPP DSL modules:

**Location:** `pkg/docs/tutorials/yapp-module-authoring-guide.md`

**Content covers:**
1. **Overview** - What you'll create, time estimates
2. **Prerequisites** - Go version, tools, knowledge needed
3. **Step-by-step workflow:**
   - Create module directory
   - Write schema.yaml
   - Validate schema
   - Generate code
   - Implement builder
   - Implement registry boilerplate
   - Test and verify
4. **Common patterns:**
   - Optional fields with defaults
   - Enum validation
   - Nested objects
   - Shape-specific logic
5. **Troubleshooting** - Common errors and fixes
6. **Best practices** - Naming conventions, documentation, testing
7. **Complete example** - LED indicators module showing all concepts

**Features:**
- Real code examples from existing modules
- Command-line snippets for each step
- Troubleshooting section with error messages and fixes
- Best practices for naming, testing, documentation
- Links to existing modules as references

**Testing:**
```bash
go run ./cmd/yappctl help yapp-module-authoring-guide
```

Guide is now available in the CLI help system!
