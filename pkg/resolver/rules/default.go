package rules

// DefaultRegistry returns a registry with all built-in rules registered.
func DefaultRegistry() *Registry {
	reg := NewRegistry()
	reg.Register(&YamlSyntaxPointerRule{})
	reg.Register(&EnumSuggestClosestRule{})
	reg.Register(&VarsScaffoldRule{})
	reg.Register(&ModuleDocEmbedRule{})
	reg.Register(&ModuleDocEmbedRule{})
	return reg
}

