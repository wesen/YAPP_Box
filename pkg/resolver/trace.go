package resolver

// TraceKind categorizes how a value was produced during resolution.
type TraceKind string

const (
	// TraceKindLiteral indicates the value came directly from the DSL (number, string, bool).
	TraceKindLiteral TraceKind = "literal"
	// TraceKindExpression indicates the value is the result of evaluating an expression.
	TraceKindExpression TraceKind = "expression"
)

// TraceEntry describes how a single path in the resolved document was produced.
type TraceEntry struct {
	Path         string
	Kind         TraceKind
	Expression   string
	Value        any
	Dependencies []string
}

// Trace is a lookup map keyed by dotted document paths.
type Trace map[string]TraceEntry

// Get returns a trace entry for a given path.
func (t Trace) Get(path string) (TraceEntry, bool) {
	if t == nil {
		return TraceEntry{}, false
	}
	entry, ok := t[path]
	return entry, ok
}

// Result bundles the resolved document with trace metadata.
type Result struct {
	Document map[string]any
	Trace    Trace
}

type traceRecorder struct {
	entries map[string]TraceEntry
}

func newTraceRecorder() *traceRecorder {
	return &traceRecorder{
		entries: map[string]TraceEntry{},
	}
}

func (tr *traceRecorder) recordLiteral(path string, value any, literalText string) {
	if tr == nil || path == "" {
		return
	}
	if _, exists := tr.entries[path]; exists {
		// Preserve the first (typically most informative) entry.
		return
	}
	tr.entries[path] = TraceEntry{
		Path:       path,
		Kind:       TraceKindLiteral,
		Expression: literalText,
		Value:      value,
	}
}

func (tr *traceRecorder) recordExpression(path, expr string, value any, deps []string) {
	if tr == nil || path == "" {
		return
	}
	tr.entries[path] = TraceEntry{
		Path:         path,
		Kind:         TraceKindExpression,
		Expression:   expr,
		Value:        value,
		Dependencies: deps,
	}
}

func (tr *traceRecorder) Trace() Trace {
	if tr == nil {
		return nil
	}
	out := make(Trace, len(tr.entries))
	for k, v := range tr.entries {
		out[k] = v
	}
	return out
}
