---
Title: Diary — YAML comment carryover
Ticket: YAPP-SCAD-COMMENTS-001
LastUpdated: 2025-11-28T17:10:00-05:00
---

# Diary — YAML comment carryover

## What I did
- Replaced the old `yaml.Unmarshal` loader with a `yaml.Node` walker that preserves head/line/foot comments for every path.
- Extended the `resolvercli → yappctl → generatorcli → yappgen` pipeline to pass a `map[path][]string` of comments.
- Updated `Provenance` and the SCAD emitter so YAML comments appear before the computed provenance lines, and expanded tests accordingly.

## What worked well
- The node walker doubles as the numeric decoder, so we only parse the document once.
- Augmenting `Provenance` kept the SCAD emission changes localized; no need to touch every feature module separately.
- Tests were easy to extend by injecting synthetic comment maps.

## What didn’t work
- First attempt at comment extraction reattached root-level comments to an empty path, so nothing printed; fixed by skipping blank paths.
- Forgot that YAML comments still contain `#` prefixes; the SCAD output looked noisy until I stripped them during normalization.

## What I learned
- `yaml.Node` carries enough metadata (comments + line numbers) to unlock future features (line anchors, error reporting) with little extra work.
- Prepending YAML comments keeps user intent front-and-center, making downstream debugging far easier.

## What to improve next time
- Add integration tests that run `yappctl generate` on a real fixture to catch wiring mistakes earlier.
- Explore deduplicating repeated comments when the same path feeds several SCAD artifacts.

