# Tasks

## TODO

- [ ] Draft playbook structure and scope
- [ ] Design YAML DSL schema v0
- [ ] Collect relevant docs from ticket YAPP-DOCS-001
- [ ] Update ticket index.md with links and next steps
- [ ] Draft analysis document summarizing findings and requirements
- [ ] Add example YAML and OpenSCAD mapping strategy
- [ ] Update vocabulary if new docTypes/topics needed
- [x] Scaffold Go CLI resolver (cobra) to parse DSL and output resolved YAML
- [x] Implement expression parsing, dotted-path resolution, and fixed-point evaluation
- [ ] Add schema validation, error reporting, and CLI flags (input/output)
- [x] Write unit tests incl. cycles, unresolved refs, and numeric ops
- [ ] Document CLI usage and integrate into playbook Commands section
- [x] Update DSL spec: remove coordinates section, document YAPP defaults (1-2h)
- [x] Build parameter schemas for pcbStands, connectors, snapJoins, cutouts (4-6h)
- [x] Implement generator with semantic fixes (6-10h)
- [x] Test generated SCAD against real examples (2-4h)
- [x] Emit SCAD header/globals/features/footer; prepare template bindings
- [x] Add unit tests and golden outputs for generator mapping
- [x] Integrate Cobra CLI yapp-gen end-to-end (resolver+emitter)
- [x] Add YAML examples under examples/ and generate SCAD outputs
- [x] Write unit tests for mapping, cutout shapes, snap sides, emitter
- [ ] Extend DSL/generator to support push buttons (from YAPP_Demo_buttons_v31.scad parity)
