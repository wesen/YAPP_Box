# Tasks

## TODO

- [ ] Finalize push_buttons DSL schema (field names, enums, defaults) covering circle, rectangle, rounded, and polygon cases.
- [ ] Implement `pushButtons` builder in `pkg/yappgen` and wire to `Model` so yappctl emit includes SCAD arrays and toggles `printSwitchExtenders`.
- [ ] Extend examples (`examples/yapp-demo-buttons.yaml`, new buttons2 YAML) plus regression tests to prove parity with SCAD outputs.
- [x] Document the DSL usage via `pkg/docs/tutorials/push-buttons-dsl.md` and ticket playbooks.
