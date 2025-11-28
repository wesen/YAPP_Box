---
Title: Diary — SCAD provenance pass
Ticket: YAPP-SCAD-COMMENTS-001
LastUpdated: 2025-11-28T17:10:00-05:00
---

# Diary — SCAD provenance pass

## What I did
- Threaded resolver trace metadata from `resolvercli` through `yappctl` into `yappgen`.
- Added `Provenance` tracking so every SCAD scalar and array row knows its originating DSL path and expression.
- Updated the SCAD emitter to print expression + dependency comments and covered it with regression tests.

## What worked well
- The resolver already exposed variable references, so adding trace capture was straightforward.
- Centralizing comment emission inside `Provenance` kept `emit.go` changes minimal.
- Tests using `emitSCAD` snapshots caught regressions immediately.

## What didn’t work
- Initially tried to bolt provenance comments straight into `emit.go`, which duplicated logic; refactoring into a helper fixed the mess.
- Forgot to thread trace data through `generatorcli` on the first pass, causing nil-pointer panics until I added compile-time coverage.

## What I learned
- Keeping metadata in a single struct (`Provenance`) makes later enhancements (like YAML comments) much easier.
- Resolver trace paths act as a stable contract across the stack; leaning on those avoids brittle bespoke lookups.

## What to improve next time
- Start with an explicit data-flow diagram so I remember to update every hop the first time.
- Add high-level diaries as I work instead of after the fact to capture more detail (doing that now!).

