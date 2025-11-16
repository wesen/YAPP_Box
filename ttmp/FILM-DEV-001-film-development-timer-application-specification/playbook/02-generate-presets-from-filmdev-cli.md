---
Title: Generate presets from filmdev CLI
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
Summary: Use the filmdev CLI (SQLite) to export developer/film/ISO times and convert them into film-developer/data/presets.json entries
LastUpdated: 2025-11-15T21:47:41.722597366-05:00
---


# Generate presets from filmdev CLI
 
## Purpose
Convert authoritative film development times from the `filmdev` SQLite/CLI project into presets consumable by the Pico W app (`film-developer/data/presets.json`).
 
## Environment
- filmdev project available locally:
  - DB: `/home/manuel/code/wesen/corporate-headquarters/vibes/2025-11-15/filmdev-project/filmdev.db`
  - CLI: `/home/manuel/code/wesen/corporate-headquarters/vibes/2025-11-15/filmdev-project/filmdev-cli/filmdev`
 
## Quick reference (Tri‑X + XTOL/D‑76 @ ISO 400/800/1600)
 
```bash
DB=/home/manuel/code/wesen/corporate-headquarters/vibes/2025-11-15/filmdev-project/filmdev.db
CLI=/home/manuel/code/wesen/corporate-headquarters/vibes/2025-11-15/filmdev-project/filmdev-cli/filmdev
 
# Tri‑X 400 in XTOL at 20C (selected examples)
$CLI query --db-path "$DB" --film "Tri-X" --developer "Xtol" --iso 400 --output json
$CLI query --db-path "$DB" --film "Tri-X" --developer "Xtol" --iso 800 --output json
$CLI query --db-path "$DB" --film "Tri-X" --developer "Xtol" --iso 1600 --output json
 
# Tri‑X 400 in D‑76 at 20C (selected examples)
$CLI query --db-path "$DB" --film "Tri-X" --developer "D-76" --iso 400 --output json
$CLI query --db-path "$DB" --film "Tri-X" --developer "D-76" --iso 800 --output json
$CLI query --db-path "$DB" --film "Tri-X" --developer "D-76" --iso 1600 --output json
 
# HP5+ in XTOL/D‑76 (push +2 to ISO1600 + others)
$CLI query --db-path "$DB" --film "HP5"  --developer "Xtol" --iso 400  --output json
$CLI query --db-path "$DB" --film "HP5"  --developer "Xtol" --iso 1600 --output json
$CLI query --db-path "$DB" --film "HP5"  --developer "D-76" --iso 400  --output json
$CLI query --db-path "$DB" --film "HP5"  --developer "D-76" --iso 800  --output json
$CLI query --db-path "$DB" --film "HP5"  --developer "D-76" --iso 1600 --output json
```
 
## Converting decimal minutes → MM:SS
`filmdev` returns minutes as decimals (e.g., 11.75). Convert with:
 
```bash
python3 - <<'PY'
def mmss(x):
    m = int(x)
    s = int(round((x - m) * 60))
    return f"{m:02d}:{s:02d}"
for v in [6.75, 8.75, 9.75, 10, 11.75, 13.25, 16.5, 18]:
    print(v, "->", mmss(float(v)))
PY
```
 
Examples:
- 8.75 → 08:45
- 9.75 → 09:45
- 10 → 10:00
- 11.75 → 11:45
- 13.25 → 13:15
- 16.5 → 16:30
- 18 → 18:00
 
## Preset JSON shape
 
```json
{
  "id": "kodak_tri_x_400_xtol_1_1_iso_1600_20c",
  "film": "Kodak Tri-X 400",
  "developer": "Xtol",
  "dilution": "1+1",
  "iso": 1600,
  "temp_c": 20.0,
  "developer_time": "11:45",
  "stages": { "stop_bath": "00:30", "fixer": "10:00", "wash": "15:00" },
  "label": "TRI-X 400 XTOL 1+1 ISO1600"
}
```
 
## Workflow
1) Query filmdev for the film/developer/ISO tuples you care about (prefer 20C entries and dilution 1+1 where available).
2) Convert `time_35mm` decimal to `developer_time` MM:SS.
3) Add/update objects in `film-developer/data/presets.json` using the shape above.
4) In the device UI:
   - Developer → Film → Target ISO
   - The app picks the “best” preset (prefers 1+1 @ ~20C) and applies the developer stage time
 
## Notes
- If multiple entries exist (stock/1+2/1+3), keep at least one “canonical” (1+1) and optionally add others.
- Keep labels ≤ 16 chars when possible; long labels will scroll (marquee).
- Stages (stop/fixer/wash) are constant here; adjust to your process as needed.