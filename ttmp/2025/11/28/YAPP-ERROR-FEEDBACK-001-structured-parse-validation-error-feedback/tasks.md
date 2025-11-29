# Tasks

## TODO

- [x] Add tasks here

- [x] Create pkg/resolver/errorx package with taxonomy schema (StageCode, SymptomCode, Severity enums, Taxonomy struct, TaxonomyContext interface, and context structs)
- [x] Update error producers to emit taxonomy entries (resolvercli/resolver.go, resolver/validation.go, resolver/resolver.go, resolver/strict.go)
- [x] Create pkg/resolver/rules package with registry (Renderer interface, Registry with Register/RenderAll, CLI text adapter)
- [x] Implement YamlSyntaxPointerRule (matches YAML syntax errors, shows snippet with pointer)
- [x] Implement EnumSuggestClosestRule (matches enum mismatches, suggests closest value)
- [x] Implement VarsScaffoldRule (matches missing vars, generates YAML scaffold)
- [x] Wire rule registry into yappctl CLI (catch errors, extract taxonomy, render help)
- [x] Add unit tests for taxonomy and rules
- [x] Complete position integration: Pass PositionMap through resolver pipeline and populate Line/Column fields in taxonomy contexts (currently only YAML syntax errors have line numbers)
- [x] Add unit tests for taxonomy: Test constructors, AsTaxonomy() unwrapping, and context types in pkg/resolver/errorx/
- [x] Add unit tests for rules: Test registry matching/sorting/aggregation, and test each rule (YAML syntax, enum suggest, vars scaffold) in pkg/resolver/rules/
- [ ] Enhance schema validation taxonomy: Extract enum values and min/max constraints from schema validation errors to populate SchemaConstraintContext.Allowed, Min, Max fields
- [ ] Implement ModuleDocEmbedRule: Rule that embeds module field tables from pkg/docs/schema_help.go for schema validation errors
- [ ] Implement DependencyGraphRule: Rule that shows dependency chains for missing variables to help users understand resolution order
- [ ] Implement YamlKnownFieldsRule: Rule that suggests known fields when unknown keys are detected (requires yaml.v3 KnownFields support)
