# Changelog

## 2025-11-16

- Initial workspace created


## 2025-11-16

Implemented embedded generator: YAPPgenerator_v3.scad now embedded in yappctl binary, auto-extracted when rendering STLs or with --copy-generator flag. Creates standalone SCAD files without manual sed/path manipulation

### Related Files

- /home/manuel/code/others/YAPP_Box/cmd/yappctl/generate_command.go — --copy-generator and --generator-path flags
- /home/manuel/code/others/YAPP_Box/pkg/cli/generatorcli/generator.go — copyGeneratorToOutput function
- /home/manuel/code/others/YAPP_Box/pkg/yappgen/assets/assets.go — Go embed directive


## 2025-11-16

Created scripts/compare_all.sh for automated DSL vs legacy SCAD comparison with STL generation

### Related Files

- /home/manuel/code/others/YAPP_Box/scripts/compare_all.sh — automated comparison script


## 2025-11-16

Updated yappctl-cli-overview docs with embedded generator documentation and new flags

### Related Files

- /home/manuel/code/others/YAPP_Box/pkg/docs/tutorials/yappctl-cli-overview.md — document embedded generator

