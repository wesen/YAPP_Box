# Changelog

## 2025-11-15

- Initial workspace created
- Drafted CLI merge assessment (differences, architecture, tasks) and populated tasks backlog
- Implemented unified `yappctl` CLI (resolve + generate verbs), extracted resolver/generator helpers, loaded help docs, and removed legacy binaries per "no shims" decision

## 2025-11-15 - Documented yappctl usage + manual smoke test

Recorded the yappctl generate smoke-playbook after verifying STL outputs, refreshed README with the new CLI workflow, and added a help-system tutorial so users find the unified verbs directly from the CLI.

### Related Files

- README.md — New yappctl section with resolve/generate examples
- pkg/docs/tutorials/yappctl-cli-overview.md — Help tutorial surfaced via yappctl help
- ttmp/YAPP-CLI-MERGE-001-unify-resolver-generator-clis/playbooks/yappctl-manual-verification.md — Manual smoke test instructions

