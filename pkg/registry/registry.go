package registry

import (
	"fmt"
	"sort"
	"sync"
)

var (
	modulesMu   sync.RWMutex
	modules     = make(map[string]FeatureModule)
	moduleOrder []string
)

// Register adds a module to the global registry, enforcing deterministic order
// and detecting duplicate paths.
func Register(module FeatureModule) {
	if module == nil {
		panic("registry: cannot register a nil module")
	}
	schema := module.Schema()
	if schema == nil {
		panic("registry: module provided a nil schema")
	}
	path := schema.Path()
	if path == "" {
		panic("registry: module schema returned an empty path")
	}

	modulesMu.Lock()
	defer modulesMu.Unlock()

	if _, exists := modules[path]; exists {
		panic(fmt.Sprintf("registry: duplicate registration for path %q", path))
	}

	modules[path] = module
	moduleOrder = append(moduleOrder, path)
}

// Get retrieves a module by registry path.
func Get(path string) (FeatureModule, bool) {
	modulesMu.RLock()
	defer modulesMu.RUnlock()
	mod, ok := modules[path]
	return mod, ok
}

// All returns modules in deterministic registration order.
func All() []FeatureModule {
	modulesMu.RLock()
	defer modulesMu.RUnlock()

	out := make([]FeatureModule, 0, len(moduleOrder))
	for _, path := range moduleOrder {
		if mod, ok := modules[path]; ok {
			out = append(out, mod)
		}
	}
	return out
}

// SetTestRegistry replaces the registry contents for tests.
func SetTestRegistry(testModules map[string]FeatureModule) {
	modulesMu.Lock()
	defer modulesMu.Unlock()

	modules = make(map[string]FeatureModule, len(testModules))
	moduleOrder = moduleOrder[:0]

	if len(testModules) == 0 {
		return
	}

	keys := make([]string, 0, len(testModules))
	for path, mod := range testModules {
		if mod == nil {
			continue
		}
		modules[path] = mod
		keys = append(keys, path)
	}

	sort.Strings(keys)
	moduleOrder = append(moduleOrder, keys...)
}
