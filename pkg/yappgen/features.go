package yappgen

import (
	"context"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/pushbuttons"
)

// FeatureModule describes how a DSL feature is collected and emitted.
type FeatureModule interface {
	Name() string
	Collect(resolved map[string]any, features map[string]any, model *Model) error
	Emit(ctx context.Context, model *Model, b *strings.Builder) error
}

var featureModules = []FeatureModule{
	newArrayFeatureModule("pcb_stands", "pcbStands",
		func(m *Model) *[]map[string]any { return &m.PcbStands },
		buildPcbStands, nil),
	newArrayFeatureModule("connectors", "connectors",
		func(m *Model) *[]map[string]any { return &m.Connectors },
		buildConnectors, nil),
	newArrayFeatureModule("push_buttons", "pushButtons",
		func(m *Model) *[]map[string]any { return &m.PushButtons },
		pushbuttons.Build,
		func(m *Model, items []map[string]any) {
			m.PrintSwitchExtenders = len(items) > 0
		}),
	newArrayFeatureModule("snap_joins", "snapJoins",
		func(m *Model) *[]map[string]any { return &m.SnapJoins },
		buildSnapJoins, nil),
	newCutoutFeatureModule(),
}

func collectFeatureModules(resolved map[string]any, features map[string]any, model *Model) error {
	for _, module := range featureModules {
		if err := module.Collect(resolved, features, model); err != nil {
			return errors.Wrapf(err, "feature %s", module.Name())
		}
	}
	return nil
}

func emitFeatureModules(ctx context.Context, model *Model, b *strings.Builder) error {
	for _, module := range featureModules {
		if err := module.Emit(ctx, model, b); err != nil {
			return errors.Wrapf(err, "feature %s", module.Name())
		}
	}
	return nil
}

type arrayFeatureModule struct {
	key          string
	scadName     string
	field        func(*Model) *[]map[string]any
	builder      func([]map[string]any) ([][]any, error)
	afterCollect func(*Model, []map[string]any)
}

func newArrayFeatureModule(
	key string,
	scadName string,
	field func(*Model) *[]map[string]any,
	builder func([]map[string]any) ([][]any, error),
	afterCollect func(*Model, []map[string]any),
) FeatureModule {
	return &arrayFeatureModule{
		key:          key,
		scadName:     scadName,
		field:        field,
		builder:      builder,
		afterCollect: afterCollect,
	}
}

func (m *arrayFeatureModule) Name() string {
	return m.key
}

func (m *arrayFeatureModule) Collect(resolved map[string]any, features map[string]any, model *Model) error {
	ptr := m.field(model)
	if features == nil {
		*ptr = nil
		if m.afterCollect != nil {
			m.afterCollect(model, nil)
		}
		return nil
	}
	arr, ok := getArray(features, m.key)
	if !ok {
		*ptr = nil
		if m.afterCollect != nil {
			m.afterCollect(model, nil)
		}
		return nil
	}
	items := normalizeArrayOfMaps(arr)
	*ptr = items
	if m.afterCollect != nil {
		m.afterCollect(model, items)
	}
	return nil
}

func (m *arrayFeatureModule) Emit(ctx context.Context, model *Model, b *strings.Builder) error {
	ptr := m.field(model)
	if len(*ptr) == 0 {
		return nil
	}
	rows, err := m.builder(*ptr)
	if err != nil {
		return err
	}
	writeArrayDecl(b, m.scadName, rows)
	b.WriteString("\n")
	return nil
}

type cutoutFeatureModule struct{}

func newCutoutFeatureModule() FeatureModule {
	return &cutoutFeatureModule{}
}

func (m *cutoutFeatureModule) Name() string {
	return "cutouts"
}

func (m *cutoutFeatureModule) Collect(resolved map[string]any, features map[string]any, model *Model) error {
	model.Cutouts = nil
	if features == nil {
		return nil
	}
	arr, ok := getArray(features, "cutouts")
	if !ok {
		return nil
	}
	for _, it := range arr {
		item, ok := it.(map[string]any)
		if !ok {
			return errors.Errorf("cutouts items must be objects, got %T", it)
		}
		face, _ := item["face"].(string)
		if face == "" {
			return errors.Errorf("cutout missing required 'face'")
		}
		model.Cutouts = append(model.Cutouts, Cutout{Face: face, Item: item})
	}
	return nil
}

func (m *cutoutFeatureModule) Emit(ctx context.Context, model *Model, b *strings.Builder) error {
	if len(model.Cutouts) == 0 {
		return nil
	}
	byFace, err := distributeCutouts(model.Cutouts)
	if err != nil {
		return err
	}
	keys := make([]string, 0, len(byFace))
	for k := range byFace {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if len(byFace[k]) == 0 {
			continue
		}
		writeArrayDecl(b, k, byFace[k])
		b.WriteString("\n")
	}
	return nil
}
