---
Title: yappctl CLI Overview
Slug: yappctl-cli-overview
Short: Quick reference for the unified resolver/generator CLI.
Topics:
  - yapp
  - cli
  - generator
  - resolver
Commands:
  - yappctl
  - go
Flags: []
IsTemplate: false
IsTopLevel: true
ShowPerDefault: true
SectionType: Tutorial
Order: 10
---

## Why yappctl?

`yappctl` unifies the old `encl-resolve` and `yapp-gen` binaries. One Glazed root command now loads tutorials (this page, documentation study guide, etc.) and exposes two verbs:

- `resolve`: reads the YAML DSL, evaluates expressions/vars, emits YAML or JSON.
- `generate`: resolves the YAML, writes SCAD with proper includes, and can render STL meshes via OpenSCAD.

## Running from Source

```bash
go run ./cmd/yappctl --help
```

The help system lists all embedded tutorials; navigate with `yappctl help tutorial <slug>`.

## Resolve Workflow

```bash
go run ./cmd/yappctl resolve \
  --input examples/yapp-demo-buttons.yaml \
  --out-file /tmp/buttons-resolved.yaml \
  --format json \
  --strict
```

- `--max-iterations` controls resolver passes (default 16).
- Accepts YAML or JSON output (`--format`).

## Generate Workflow

```bash
go run ./cmd/yappctl generate \
  --input examples/yapp-demo-buttons.yaml \
  --scad-out examples/yapp-demo-buttons.scad \
  --stl-base examples/yapp-demo-buttons-base.stl \
  --stl-lid examples/yapp-demo-buttons-lid.stl \
  --render-timeout 45s
```

**Key features:**

- **Embedded generator:** `yappctl` embeds `YAPPgenerator_v3.scad` in the binary. When rendering STLs or using `--copy-generator`, it automatically extracts the generator to the output directory, creating a standalone SCAD file.
- **Automatic generator copying:** When `--stl-base` or `--stl-lid` is specified, the generator is automatically copied (no manual `sed` commands needed). The include path is rewritten to `include <YAPPgenerator_v3.scad>` (same directory).
- **Manual copy mode:** Use `--copy-generator` to extract the generator even when not rendering STLs (useful for sharing standalone SCAD files).
- **Custom generator:** Use `--generator-path /path/to/custom/YAPPgenerator_v3.scad` to use a modified generator version.
- **Quality overrides:** Use `--quality-value 4` to force the YAPP generator’s `renderQuality` to a lower value during STL generation for faster iteration.

**Flags:**

- `--input/-i`: Path to input YAML file (required)
- `--scad-out/-o`: Path to write generated SCAD file (required)
- `--stl-base`: Path for base STL output (auto-copies generator)
- `--stl-lid`: Path for lid STL output (auto-copies generator)
- `--copy-generator`: Extract generator to output directory
- `--generator-path`: Use custom generator instead of embedded version
- `--openscad-bin`: OpenSCAD executable path (default: `openscad`)
- `--render-timeout`: Rendering timeout duration (default: `30s`)
- `--quality-value`: Override the generator `renderQuality` (set to `0` to keep defaults)
- `--max-iterations`: Resolver passes (default: 16)
- `--strict`: Enable strict validation

## Manual Smoke Test (Temporary)

Until automated goldens land, run the smoke-playbook recorded at `ttmp/YAPP-CLI-MERGE-001-.../playbooks/yappctl-manual-verification.md`:

1. Execute the `generate` command above.
2. Confirm the STL files report non-zero sizes:

   ```bash
   stat -c '%n %s' examples/yapp-demo-buttons-*.stl
   ```

Record outcomes in the ticket when debugging CLI regressions.
