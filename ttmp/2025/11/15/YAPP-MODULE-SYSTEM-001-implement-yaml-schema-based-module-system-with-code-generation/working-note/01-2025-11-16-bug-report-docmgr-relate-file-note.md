---
Title: 2025-11-16 Bug report - docmgr relate file-note
Ticket: YAPP-MODULE-SYSTEM-001
Status: active
Topics:
    - yapp
    - dsl
    - codegen
    - architecture
DocType: working-note
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: relate command rejects --file-note with 'no changes specified'
LastUpdated: 2025-11-15T20:08:10.65942996-05:00
---


# 2025-11-16 Bug report - docmgr relate file-note

## Summary

- `docmgr relate` (latest CLI from repo) refuses any `--file-note` arguments after a recent change that removed `--files` but made `--file-note` mandatory. The command always returns `Error: no changes specified` even though file-note pairs are provided.

## Notes

### Reproduction Steps

1. From repo root: `docmgr relate --ticket YAPP-MODULE-SYSTEM-001 --file-note "pkg/schemagen/validate.go:Schema validator with bool support"`
2. Tried alternate syntaxes:
   - `--file-note pkg/schemagen/validate.go:Schema validator with bool support`
   - `--file-note pkg/schemagen/validate.go=Schema validator with bool support`
   - `--doc <path-to-index.md> --file-note ...`
3. CLI output is consistently:
   ```
   Error: no changes specified. Use --file-note 'path:note' to add/update, --remove-files to remove, or --suggest --apply-suggestions to apply suggestions
   ```

### Expected Behavior

- Per `docmgr relate --help`, providing at least one `--file-note 'path:note'` should update the ticket index’s `RelatedFiles` list (or the specified doc) with that entry.

### Actual Behavior

- Command exits before writing anything, complains that no changes were specified, so we cannot add new files/notes via CLI anymore.

### Environment

- Repo: `github.com/wesen/yapp-encl-resolver`
- Root: `/home/manuel/code/others/YAPP_Box`
- Tickets stored under `ttmp/`
- docmgr version: bundled with repo (not separately versioned); issue started 2025-11-16.

### Impact

- Unable to keep `RelatedFiles` in sync via CLI, which blocks required documentation workflow (ticket instructions mandate using docmgr for file relationships).

## Decisions

- Log bug and escalate so docmgr maintainers can fix relate’s CLI argument parsing.

## Next Steps

- [ ] Share this bug report with docmgr maintainers / tooling channel.
- [ ] After fix lands, re-run `docmgr relate` to attach `pkg/schemagen/validate.go`, `pkg/yappgen/modules/pushbuttons/schema.yaml`, etc.
