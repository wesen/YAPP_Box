package resolver

import (
	"github.com/pkg/errors"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver/errorx"
)

// validateTopLevelKeys checks for unknown top-level keys in strict mode.
func validateTopLevelKeys(state map[string]any) error {
	allowed := map[string]struct{}{
		"project": {}, "version": {}, "units": {}, "yapp_version": {},
		"vars": {}, "pcb": {}, "enclosure": {}, "features": {}, "coordinates": {}, "tolerances": {}, "show": {},
	}
	var unknown []string
	for k := range state {
		if _, ok := allowed[k]; !ok {
			unknown = append(unknown, k)
		}
	}
	if len(unknown) > 0 {
		taxonomy := errorx.NewStrictModeTaxonomy("unknown-key", unknown)
		return errors.Wrap(taxonomy, "unknown top-level keys (strict mode)")
	}
	return nil
}

// validateUnusedVars checks for vars.* that were never referenced.
func validateUnusedVars(state map[string]any, used map[string]struct{}) error {
	varsMap, ok := state["vars"].(map[string]any)
	if !ok {
		return nil
	}
	var unused []string
	for k := range varsMap {
		full := "vars." + k
		if _, ok := used[full]; !ok {
			unused = append(unused, full)
		}
	}
	if len(unused) > 0 {
		taxonomy := errorx.NewStrictModeTaxonomy("unused-var", unused)
		return errors.Wrap(taxonomy, "unused variables (strict mode)")
	}
	return nil
}
