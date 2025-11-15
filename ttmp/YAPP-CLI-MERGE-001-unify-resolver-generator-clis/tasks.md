# Tasks

## TODO

- [x] Extract generator helpers from `cmd/yapp-gen/main.go` into a reusable package (SCAD + STL rendering)
- [x] Scaffold unified Glazed CLI (`cmd/yappctl`) with `resolve` and `generate` subcommands
- [x] Load `pkg/docs` help sections in the new CLI root and expose migration notes
- [x] Smoke-test yappctl resolve/generate using examples/yapp-demo-buttons and record STL size check (temporary stand-in for CLI integration tests)
- [x] Update README/playbooks to document the new CLI and removal timeline for legacy binaries
- [ ] Add automated CLI integration tests/goldens for resolve & generate outputs
