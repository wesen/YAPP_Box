# Tasks

## TODO

- [x] Finalize push_buttons DSL schema (field names, enums, defaults) covering circle, rectangle, rounded, and polygon cases.
- [x] Implement `pushButtons` builder in `pkg/yappgen` and wire to `Model` so yappctl emit includes SCAD arrays and toggles `printSwitchExtenders`.
- [x] Extend examples (`examples/yapp-demo-buttons.yaml`, new buttons2 YAML) plus regression tests to prove parity with SCAD outputs.
- [x] Document the DSL usage in `pkg/docs/tutorials/yapp-dsl-getting-started.md` and `pkg/docs/tutorials/yapp-dsl-reference.md`, plus ticket playbooks.
- [x] Capture the feature-module registry design + developer tutorial so future DSL extensions follow the new workflow.
- [x] Move the `push_buttons` implementation into its own module package to validate the registry/plugin structure.
- [ ] Run the DSL module schema debate (candidates, questions, rounds, and synthesis) to lock in the typed struct + schema registration plan.
- [ ] Update module developer documentation and registry implementation once the typed schema workflow is chosen.
