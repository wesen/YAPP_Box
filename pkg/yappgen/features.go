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
		pcbstands.NewModule().Build, nil, nil),
	newArrayFeatureModule("connectors", "connectors",
		func(m *Model) *[]map[string]any { return &m.Connectors },
		connectors.NewModule().Build, nil, nil),
	newArrayFeatureModule("box_mounts", "boxMounts",
		func(m *Model) *[]map[string]any { return &m.BoxMounts },
		boxmounts.NewModule().Build, nil, nil),
	newArrayFeatureModule("push_buttons", "pushButtons",
		func(m *Model) *[]map[string]any { return &m.PushButtons },
		pushbuttons.NewModule().Build,
		func(m *Model, items []map[string]any) {
			m.PrintSwitchExtenders = len(items) > 0
		},
		pushButtonsHeader()),
	newArrayFeatureModule("snap_joins", "snapJoins",
		func(m *Model) *[]map[string]any { return &m.SnapJoins },
		snapjoins.NewModule().Build, nil, nil),
	newArrayFeatureModule("light_tubes", "lightTubes",
		func(m *Model) *[]map[string]any { return &m.LightTubes },
		lighttubes.NewModule().Build, nil, nil),
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

func yappArrayHeader(arrayName string, override []string) []string {
	if len(override) > 0 {
		return override
	}
	switch arrayName {
	case "pcbStands":
		return paramSpecHeader("pcbStands", pcbStandsSchema)
	case "connectors":
		return paramSpecHeader("connectors", connectorsSchema)
	case "snapJoins":
		header := paramSpecHeader("snapJoins", snapJoinsSchema)
		header = append(header, "// Followed by required yapp<Side> flag (left/right/front/back) and optional alignment/origin flags.")
		return header
	case "boxMounts":
		return []string{
			"// YAPP boxMounts row: [pos | [pos, offset], screw_d, slot_width, height, fillet_radius?, face_flags..., optional flags: yappNoFillet, yappLid, yappCenter, yappAltOrigin]",
		}
	case "lightTubes":
		return []string{
			"// YAPP lightTubes row: [x, y, tube_length, tube_width, tube_wall, gap_above_pcb, shape_flag, lens_thickness?, height?, fillet_radius?, flags...]",
		}
	case "cutoutsFront", "cutoutsBack", "cutoutsLeft", "cutoutsRight", "cutoutsLid", "cutoutsBase":
		return []string{
			"// YAPP cutouts row: [pos0, pos1, width, length, radius, shape_flag, depth?, angle?, extras...]",
			"// Extras may include polygon presets, mask definitions, coordinate/origin flags, yappNoFillet, [yappPCBName, pcb_name].",
		}
	case "pushButtons":
		return pushButtonsHeader()
	default:
		return nil
	}
}

func paramSpecHeader(scadName string, specs []ParamSpec) []string {
	if len(specs) == 0 {
		return nil
	}
	names := make([]string, len(specs))
	for i, spec := range specs {
		name := spec.Name
		if !spec.Required {
			name += "?"
		}
		names[i] = name
	}
	return []string{
		fmt.Sprintf("// YAPP %s row: [%s]", scadName, strings.Join(names, ", ")),
	}
}

func pushButtonsHeader() []string {
	return []string{
		"// YAPP pushButtons row: [x, y, cap_length, cap_width, cap_radius, lid_protrusion, switch_height, switch_travel, switch_pole_diameter, top_height?, shape_flag, angle?, cap_fillet?, lid_wall?, lid_plate_thickness?, lid_slack?, lid_snap_slack?, extras...]",
		"// Extras include coordinate/origin flags, [yappPCBName, pcb_name], yappNoFillet, shape presets, mask definitions, and other generator tokens.",
	}
}

type arrayFeatureModule struct {
	key          string
	scadName     string
	field        func(*Model) *[]map[string]any
	builder      func([]map[string]any) ([]registry.ArrayDecl, error)
	afterCollect func(*Model, []map[string]any)
	header       []string
}

func newArrayFeatureModule(
	key string,
	scadName string,
	field func(*Model) *[]map[string]any,
	builder func([]map[string]any) ([]registry.ArrayDecl, error),
	afterCollect func(*Model, []map[string]any),
	header []string,
) FeatureModule {
	return &arrayFeatureModule{
		key:          key,
		scadName:     scadName,
		field:        field,
		builder:      builder,
		afterCollect: afterCollect,
		header:       header,
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
	header := yappArrayHeader(decls[0].Name, m.header)
	writeArrayDecl(b, decls[0].Name, decls[0].Rows, header, rowComments)
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

	globalIdx := 0
	for _, decl := range decls {
		if len(decl.Rows) == 0 {
			continue
		}
		header := yappArrayHeader(decl.Name, nil)
		var rowComments [][]string
		if model != nil && model.Provenance != nil {
			rowComments = make([][]string, len(decl.Rows))
			for i := range decl.Rows {
				if globalIdx < len(*ptr) {
					rowComments[i] = model.Provenance.DescribeFeatureRow(m.key, decl.Name, globalIdx, (*ptr)[globalIdx])
				}
				globalIdx++
			}
		} else {
			globalIdx += len(decl.Rows)
		}
		writeArrayDecl(b, decl.Name, decl.Rows, header, rowComments)
		b.WriteString("\n")
	}
	return nil
}
