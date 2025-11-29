# Tasks

## TODO

- [x] Add tasks here

- [x] Create pkg/resolver/errorx package with taxonomy schema (StageCode, SymptomCode, Severity enums, Taxonomy struct, TaxonomyContext interface, and context structs)
- [x] Update error producers to emit taxonomy entries (resolvercli/resolver.go, resolver/validation.go, resolver/resolver.go, resolver/strict.go)
- [x] Create pkg/resolver/rules package with registry (Renderer interface, Registry with Register/RenderAll, CLI text adapter)
- [x] Implement YamlSyntaxPointerRule (matches YAML syntax errors, shows snippet with pointer)
- [x] Implement EnumSuggestClosestRule (matches enum mismatches, suggests closest value)
- [x] Implement VarsScaffoldRule (matches missing vars, generates YAML scaffold)
- [ ] Wire rule registry into yappctl CLI (catch errors, extract taxonomy, render help)
- [ ] Add unit tests for taxonomy and rules
