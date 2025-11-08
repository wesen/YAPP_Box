---
Title: docmgr Workflow Playbook for This Project
Ticket: FILM-DEV-001
Status: review
Topics:
    - embedded
    - ui
    - hardware
    - timer
DocType: playbook
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2025-11-08T17:24:13.769595451-05:00
---


# docmgr Workflow Playbook for This Project

Tiny, practical guide for how we use docmgr in this repo. For full reference, see: `docmgr help how-to-use` and `docmgr COMMAND --help`. For writing guidance, run: `glaze help how-to-write-good-documentation-pages`.

## Prerequisites
- docmgr is installed and on PATH (check with: `docmgr help how-to-use`)
- If not installed, follow the setup tutorial: `docmgr help how-to-setup`

## Audience & Goal
- Audience: New contributors/interns joining the Film Development Timer project.
- Goal: Help you quickly create, update, and find documentation using docmgr the way we do it here, with consistent structure and high signal.

## Project Context (What is FILM-DEV-001?)
This ticket defines the Raspberry Pi Pico W film development timer MVP (UI on 128×64 OLED, 3 buttons + LEDs, DS18B20 temperature). Documentation under this ticket captures architecture, UI spec, timer logic, data/logging, and a future-ideas doc (temperature compensation).

Key docs to skim first:
- Architecture: `reference/01-application-architecture-and-system-overview.md`
- UI Spec: `reference/02-user-interface-specification-and-display-system.md`
- Timer System: `reference/03-timer-system-and-development-process-management.md`
- Data Management: `reference/04-data-management-and-logging-system.md`
- UI MVP Analysis: `analysis/01-ui-mvp-analysis-and-micropython-references.md`
- UI MVP Design: `design-doc/01-ui-mvp-design-and-implementation-plan.md`

## Doc types we use
- `index.md` (ticket root): overview, links, status
- `tasks.md`: actionable checklist for the ticket
- `changelog.md`: short entries after significant changes
- `analysis/` (analysis): problem framing, constraints, research
- `design-doc/` (design-doc): decisions, architecture, implementation plan
- `reference/` (reference): contracts/specs/long-lived reference
- `playbook/` (playbook): step-by-step how-tos like this
- `future-ideas/` (future-ideas): ideas explicitly out-of-scope for MVP

Doc type guidance:
- Start with `analysis` to clarify scope and constraints.
- Move to `design-doc` for concrete decisions and implementation plan.
- Use `reference` for stable specs people will depend on.
- Capture operational steps in `playbook` (like this document).
- Park non-MVP ideas in `future-ideas` to keep scope tight.

## Daily driver commands
```bash
# Discover and status
docmgr list tickets
docmgr list docs --ticket FILM-DEV-001
docmgr status --summary-only

# Create docs (inherit metadata from the ticket)
docmgr add --ticket FILM-DEV-001 --doc-type analysis --title "New Analysis"
docmgr add --ticket FILM-DEV-001 --doc-type design-doc --title "New Design"
docmgr add --ticket FILM-DEV-001 --doc-type reference --title "API Contract"
docmgr add --ticket FILM-DEV-001 --doc-type playbook --title "Runbook"
docmgr add --ticket FILM-DEV-001 --doc-type future-ideas --title "Idea"

# Update metadata (defaults to index.md when only --ticket is given)
docmgr meta update --ticket FILM-DEV-001 --field Status --value review
docmgr meta update --doc ttmp/FILM-DEV-001-*/analysis/01-ui-mvp-analysis-and-micropython-references.md \
  --field Summary --value "Scope and references for UI MVP"

# Relate files with notes (why the file matters)
docmgr relate --ticket FILM-DEV-001 \
  --file-note "ttmp/FILM-DEV-001-*/reference/02-user-interface-specification-and-display-system.md:UI spec source"

# Tasks and changelog
docmgr tasks add --ticket FILM-DEV-001 --text "Implement UI Display Layer"
docmgr tasks check --ticket FILM-DEV-001 --id 1
docmgr changelog update --ticket FILM-DEV-001 --entry "Added UI MVP analysis & design docs"

# Find docs
docmgr search --query "UI" --ticket FILM-DEV-001
```

## Minimal workflow (each change)
1) Edit or add the right doc (analysis/design/reference/playbook) under the ticket
2) Update frontmatter via `docmgr meta update` (Summary/Status if applicable)
3) Relate any important files with `docmgr relate --file-note "path:why-it-matters"`
4) Add a short `docmgr changelog update` entry
5) Keep `tasks.md` current (add/check items as you progress)

## Writing guidelines (short version)
Follow these principles (see `glaze help how-to-write-good-documentation-pages` for the long version):
- Start with context: what this is, why it exists, the scope, constraints.
- Be concise but complete: prefer bullet points, use short sections with clear headings.
- Make artifacts discoverable: link to related docs and code, and use `docmgr relate` with a reason.
- Prefer stable references: put durable contracts in `reference/` and evolving thinking in `analysis/`.
- Record decisions and rationale in `design-doc/` (what we decided and why, plus alternatives considered).
- Keep change history terse in `changelog.md` (imperative mood: “Added…”, “Moved…”, “Updated…”).
- One source of truth: avoid duplicating the same content in multiple docs—link instead.

## What “good” looks like
- `index.md` is under ~50 lines, explains the ticket and points to the right docs.
- Each subdoc has a one-line Summary and clear sections (Overview/Purpose, Scope, Decisions, Next Steps).
- Commands in playbooks are copy-pasteable and minimal; long explanations go into comments above the command blocks.
- Every significant change has a changelog entry.
- Related files have notes that explain “why this file matters”.

## Examples (copy/paste)
Update ticket summary and status:
```bash
docmgr meta update --ticket FILM-DEV-001 --field Summary \
  --value "MVP Pico W film development timer: UI, timer, logging"
docmgr meta update --ticket FILM-DEV-001 --field Status --value review
```

Add a design doc and set its Summary:
```bash
docmgr add --ticket FILM-DEV-001 --doc-type design-doc --title "Timer Loop Refinement"
DOC="ttmp/FILM-DEV-001-*/design-doc/02-timer-loop-refinement.md"
docmgr meta update --doc "$DOC" --field Summary --value "Cooperative timer loop with LED cues"
```

Relate a spec that influenced your change:
```bash
docmgr relate --ticket FILM-DEV-001 \
  --file-note "ttmp/FILM-DEV-001-*/reference/02-user-interface-specification-and-display-system.md:UI grids and button roles"
```

Record the change:
```bash
docmgr changelog update --ticket FILM-DEV-001 --entry "Added Timer Loop Refinement design doc"
```

## Common pitfalls
- Forgetting to update Summary or Status in frontmatter (use `docmgr meta update`).
- Not relating influential files (makes review harder later).
- Writing long narrative without headings—prefer scannable sections and bullets.
- Duplicating the same content across docs—link instead.

## Troubleshooting
- Not sure which command to use? `docmgr help how-to-use`
- Want exact flags for a command? `docmgr COMMAND --help`
- Folder missing for a doc type? `docmgr add` will create it on first use.
- Need to renumber files? `docmgr renumber --ticket FILM-DEV-001`
- `--version` is not supported; use `docmgr help` to verify installation

## Exit criteria (for your doc changes)
- The right doc type was used; frontmatter Summary/Status are set.
- Related files are added with notes.
- A concise changelog entry was recorded.
- The doc is skimmable, links to key resources, and leaves no major “why” unanswered.
