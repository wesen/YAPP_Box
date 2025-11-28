package yappgen

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/wesen/yapp-encl-resolver/pkg/registry"
	boxmounts "github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/boxmounts"
	connectors "github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/connectors"
	cutouts "github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/cutouts"
	lighttubes "github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/lighttubes"
	pcbstands "github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/pcbstands"
	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/pushbuttons"
	snapjoins "github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/snapjoins"
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
		pcbstands.NewModule().Build, nil),
	newArrayFeatureModule("connectors", "connectors",
		func(m *Model) *[]map[string]any { return &m.Connectors },
		connectors.NewModule().Build, nil),
	newArrayFeatureModule("box_mounts", "boxMounts",
		func(m *Model) *[]map[string]any { return &m.BoxMounts },
		boxmounts.NewModule().Build, nil),
	newArrayFeatureModule("push_buttons", "pushButtons",
		func(m *Model) *[]map[string]any { return &m.PushButtons },
		pushbuttons.NewModule().Build,
		func(m *Model, items []map[string]any) {
			m.PrintSwitchExtenders = len(items) > 0
		}),
	newArrayFeatureModule("snap_joins", "snapJoins",
		func(m *Model) *[]map[string]any { return &m.SnapJoins },
		snapjoins.NewModule().Build, nil),
	newArrayFeatureModule("light_tubes", "lightTubes",
		func(m *Model) *[]map[string]any { return &m.LightTubes },
		lighttubes.NewModule().Build, nil),
	newMultiArrayFeatureModule("cutouts",
		func(m *Model) *[]map[string]any { return &m.Cutouts },
		cutouts.NewModule().Build),
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
	builder      func([]map[string]any) ([]registry.ArrayDecl, error)
	afterCollect func(*Model, []map[string]any)
}

func newArrayFeatureModule(
	key string,
	scadName string,
	field func(*Model) *[]map[string]any,
	builder func([]map[string]any) ([]registry.ArrayDecl, error),
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
	if model != nil && model.Provenance != nil {
		model.Provenance.RegisterFeature(m.key, len(items))
	}
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

	decls, err := m.builder(*ptr)
	if err != nil {
		return err
	}

	if len(decls) != 1 {
		return fmt.Errorf("arrayFeatureModule %s returned %d arrays, expected 1", m.key, len(decls))
	}

	if decls[0].Name != m.scadName {
		return fmt.Errorf("arrayFeatureModule %s returned array %s, expected %s", m.key, decls[0].Name, m.scadName)
	}

	if len(decls[0].Rows) == 0 {
		return nil
	}

	var rowComments [][]string
	if model != nil && model.Provenance != nil {
		items := *m.field(model)
		if len(items) == len(decls[0].Rows) {
			rowComments = make([][]string, len(items))
			for idx := range items {
				rowComments[idx] = model.Provenance.DescribeFeatureRow(m.key, decls[0].Name, idx, items[idx])
			}
		}
	}
	writeArrayDecl(b, decls[0].Name, decls[0].Rows, rowComments)
	b.WriteString("\n")
	return nil
}

type multiArrayFeatureModule struct {
	key     string
	field   func(*Model) *[]map[string]any
	builder func([]map[string]any) ([]registry.ArrayDecl, error)
}

func newMultiArrayFeatureModule(
	key string,
	field func(*Model) *[]map[string]any,
	builder func([]map[string]any) ([]registry.ArrayDecl, error),
) FeatureModule {
	return &multiArrayFeatureModule{
		key:     key,
		field:   field,
		builder: builder,
	}
}

func (m *multiArrayFeatureModule) Name() string {
	return m.key
}

func (m *multiArrayFeatureModule) Collect(resolved map[string]any, features map[string]any, model *Model) error {
	ptr := m.field(model)
	if features == nil {
		*ptr = nil
		return nil
	}
	arr, ok := getArray(features, m.key)
	if !ok {
		*ptr = nil
		return nil
	}
	items := normalizeArrayOfMaps(arr)
	*ptr = items
	if model != nil && model.Provenance != nil {
		model.Provenance.RegisterFeature(m.key, len(items))
	}
	return nil
}

func (m *multiArrayFeatureModule) Emit(ctx context.Context, model *Model, b *strings.Builder) error {
	ptr := m.field(model)
	if len(*ptr) == 0 {
		return nil
	}

	decls, err := m.builder(*ptr)
	if err != nil {
		return err
	}

	for _, decl := range decls {
		if len(decl.Rows) == 0 {
			continue
		}
		writeArrayDecl(b, decl.Name, decl.Rows, nil)
		b.WriteString("\n")
	}
	return nil
}
