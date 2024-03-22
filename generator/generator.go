package generator

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	geckpb "github.com/delaneyj/geck/pb/gen/geck/v1"
	"github.com/delaneyj/toolbelt"
	"github.com/samber/lo"
)

const (
	userDefineComponentIDStart = 1000
)

//go:embed templates/*
var templatesFS embed.FS

type OutputOptions struct {
	Path        string
	PackageName string
}

type stateTemplateData struct {
	PackageName string
	Name        toolbelt.CasedString
	Description string
	States      []toolbelt.CasedString
}

type bundleTemplateData struct {
	IsBuiltin  bool
	Name       toolbelt.CasedString
	Components []componentTemplateData
	States     []stateTemplateData
}

type fieldTemplateData struct {
	Name                   toolbelt.CasedString
	Type                   string
	Description            string
	SizeBytes, OffsetBytes int
	ResetValue             string
}

type componentTemplateData struct {
	PackageName string
	Name        toolbelt.CasedString
	Description string
	Fields      []fieldTemplateData
	SizeBytes   int
	ID          int
	IsBuiltin   bool
	IsTag       bool
	CanBePair   bool
}

func GenerateECS(ctx context.Context, opts OutputOptions, bundleDefs ...*geckpb.BundleDefinition) error {
	if opts.PackageName == "" {
		parts := strings.Split(opts.Path, string(filepath.Separator))
		last := parts[len(parts)-1]
		opts.PackageName = toolbelt.ToCasedString(last).Camel
	}

	os.RemoveAll(opts.Path)
	if err := os.MkdirAll(opts.Path, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	templatesSubFS, err := fs.Sub(templatesFS, "templates")
	if err != nil {
		return fmt.Errorf("failed to access templates: %w", err)
	}
	tmpls, err := template.ParseFS(templatesSubFS, "*.gtpl")
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	bb := bundleDefinitionToTemplateData(opts, true, builtinBundleDef)
	bb.IsBuiltin = true
	if err := generateBundle(ctx, opts, tmpls, bb); err != nil {
		return fmt.Errorf("failed to generate builtin bundle: %w", err)
	}

	userBundles := lo.Map(bundleDefs, func(bd *geckpb.BundleDefinition, i int) bundleTemplateData {
		return bundleDefinitionToTemplateData(opts, false, bd)
	})
	for _, b := range userBundles {
		if err := generateBundle(ctx, opts, tmpls, b); err != nil {
			return fmt.Errorf("failed to generate bundle: %w", err)
		}
	}

	allBundles := make([]bundleTemplateData, len(userBundles)+1)
	allBundles[0] = bb
	copy(allBundles[1:], userBundles)

	type worldTemplateData struct {
		PackageName string
		Bundles     []bundleTemplateData
	}
	worldData := worldTemplateData{
		PackageName: opts.PackageName,
		Bundles:     allBundles,
	}

	worldPath := filepath.Join(opts.Path, "world.go")
	f, err := os.Create(worldPath)
	if err != nil {
		return fmt.Errorf("failed to create world file: %w", err)
	}
	defer f.Close()

	if err := tmpls.ExecuteTemplate(f, "world.gtpl", worldData); err != nil {
		return fmt.Errorf("failed to execute world template: %w", err)
	}

	return nil
}

func bundleDefinitionToTemplateData(opts OutputOptions, isBuiltin bool, bd *geckpb.BundleDefinition) bundleTemplateData {
	return bundleTemplateData{
		IsBuiltin: isBuiltin,
		Name:      toolbelt.ToCasedString(bd.Name),
		States: lo.Map(bd.States, func(s *geckpb.StateDefinition, i int) stateTemplateData {
			data := stateTemplateData{
				PackageName: opts.PackageName,
				Name:        toolbelt.ToCasedString(s.Name),
				Description: s.Description,
				States: lo.Map(s.States, func(s string, i int) toolbelt.CasedString {
					return toolbelt.ToCasedString(s)
				}),
			}
			return data
		}),
		Components: lo.Map(bd.Components, func(c *geckpb.ComponentDefinition, i int) componentTemplateData {
			startID := 1
			if !isBuiltin {
				startID = userDefineComponentIDStart
			}
			offset := 0
			data := componentTemplateData{
				PackageName: opts.PackageName,
				Name:        toolbelt.ToCasedString(c.Name),
				ID:          startID + i,
				Description: c.Description,
				IsBuiltin:   isBuiltin,
				IsTag:       len(c.Fields) == 0,
				CanBePair:   c.CanBePair,
				Fields: lo.Map(c.Fields, func(f *geckpb.FieldDefinition, i int) fieldTemplateData {
					ftd := fieldTemplateData{
						Name:        toolbelt.ToCasedString(f.Name),
						Description: f.Description,
					}
					switch f.ResetValue.(type) {
					case *geckpb.FieldDefinition_U8:
						ftd.Type = "uint8"
						ftd.SizeBytes = 1
						ftd.ResetValue = fmt.Sprintf("%d", f.GetU8())
					case *geckpb.FieldDefinition_U16:
						ftd.Type = "uint16"
						ftd.SizeBytes = 2
						ftd.ResetValue = fmt.Sprintf("%d", f.GetU16())
					case *geckpb.FieldDefinition_U32:
						ftd.Type = "uint32"
						ftd.SizeBytes = 4
						ftd.ResetValue = fmt.Sprintf("%d", f.GetU32())
					case *geckpb.FieldDefinition_U64:
						ftd.Type = "uint64"
						ftd.SizeBytes = 8
						ftd.ResetValue = fmt.Sprintf("%d", f.GetU64())
					case *geckpb.FieldDefinition_I8:
						ftd.Type = "int8"
						ftd.SizeBytes = 1
						ftd.ResetValue = fmt.Sprintf("%d", f.GetI8())
					case *geckpb.FieldDefinition_I16:
						ftd.Type = "int16"
						ftd.SizeBytes = 2
						ftd.ResetValue = fmt.Sprintf("%d", f.GetI16())
					case *geckpb.FieldDefinition_I32:
						ftd.Type = "int32"
						ftd.SizeBytes = 4
						ftd.ResetValue = fmt.Sprintf("%d", f.GetI32())
					case *geckpb.FieldDefinition_I64:
						ftd.Type = "int64"
						ftd.SizeBytes = 8
						ftd.ResetValue = fmt.Sprintf("%d", f.GetI64())
					case *geckpb.FieldDefinition_F32:
						ftd.Type = "float32"
						ftd.SizeBytes = 4
						ftd.ResetValue = fmt.Sprintf("%f", f.GetF32())
					case *geckpb.FieldDefinition_F64:
						ftd.Type = "float64"
						ftd.SizeBytes = 8
						ftd.ResetValue = fmt.Sprintf("%f", f.GetF64())
					case *geckpb.FieldDefinition_Txt:
						ftd.Type = "string"
						ftd.SizeBytes = 16
						ftd.ResetValue = fmt.Sprintf(`"%s"`, f.GetTxt())
					case *geckpb.FieldDefinition_Bin:
						ftd.Type = "[]byte"
						ftd.SizeBytes = 16
						ftd.ResetValue = fmt.Sprintf("[]byte(%v)", f.GetBin())
					default:
						panic(fmt.Sprintf("unsupported field type: %s", f.String()))
					}

					ftd.OffsetBytes = offset
					offset += ftd.SizeBytes

					return ftd
				}),
			}

			for _, f := range data.Fields {
				data.SizeBytes += f.SizeBytes
			}

			return data
		}),
	}
}

func generateBundle(ctx context.Context, opts OutputOptions, tmpls *template.Template, bundle bundleTemplateData) error {
	if err := generateFiles(
		ctx, opts, tmpls,
		"entities",
		"components",
		"archetypes",
		"queries",
	); err != nil {
		return fmt.Errorf("failed to generate known files: %w", err)
	}

	for _, state := range bundle.States {
		if err := generateState(ctx, opts, tmpls, state); err != nil {
			return fmt.Errorf("failed to generate state: %w", err)
		}
	}

	for i, component := range bundle.Components {
		if err := generateComponent(ctx, opts, tmpls, bundle.IsBuiltin, i, component); err != nil {
			return fmt.Errorf("failed to generate component: %w", err)
		}
	}

	return nil
}

func generateFiles(ctx context.Context, opts OutputOptions, tmpls *template.Template, knownFiles ...string) error {
	type fileTemplateData struct {
		PackageName string
	}
	data := fileTemplateData{
		PackageName: opts.PackageName,
	}

	for _, filename := range knownFiles {
		n := toolbelt.ToCasedString(filename)
		path := filepath.Join(opts.Path, n.Snake+".go")
		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer f.Close()

		tmplName := n.Snake + ".gtpl"
		if err := tmpls.ExecuteTemplate(f, tmplName, data); err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}
	}
	return nil
}

func generateState(ctx context.Context, opts OutputOptions, tmpls *template.Template, state stateTemplateData) error {

	statePath := filepath.Join(opts.Path, fmt.Sprintf("state_%s.go", state.Name.Camel))
	f, err := os.Create(statePath)
	if err != nil {
		return fmt.Errorf("failed to create state file: %w", err)
	}
	defer f.Close()

	if err := tmpls.ExecuteTemplate(f, "state.gtpl", state); err != nil {
		return fmt.Errorf("failed to execute state template: %w", err)
	}

	return nil
}

func generateComponent(ctx context.Context, opts OutputOptions, tmpls *template.Template, isBuiltin bool, i int, c componentTemplateData) error {

	prefix := "component"
	if isBuiltin {
		prefix += "_builtin"
	}
	componentPath := filepath.Join(opts.Path, fmt.Sprintf("%s_%s.go", prefix, c.Name.Snake))
	f, err := os.Create(componentPath)
	if err != nil {
		return fmt.Errorf("failed to create component file: %w", err)
	}
	defer f.Close()

	if err := tmpls.ExecuteTemplate(f, "component.gtpl", c); err != nil {
		return fmt.Errorf("failed to execute component template: %w", err)
	}

	return nil
}

var (
	builtinBundleDef = &geckpb.BundleDefinition{
		Name:        "Builtin",
		Description: "The builtin bundle",
		Components: []*geckpb.ComponentDefinition{
			{
				Name:        "Identifier",
				Description: "A string identifier for an entity. This is used to identify an entity across the network and UI.",
				CanBePair:   true,
				Fields: []*geckpb.FieldDefinition{
					{
						Name: "Value",
						ResetValue: &geckpb.FieldDefinition_Txt{
							Txt: "!!!UNKNOWN IDENTIFIER!!!",
						},
					},
				},
			},
			{
				Name:        "Name",
				Description: "A name tag",
			},
			{
				Name:        "Internal",
				Description: "Internal components are used by the ECS system and are not user-defined.",
			},
			{
				Name:        "Wildcard",
				Description: "A wildcard tag. This is used to match any entity, used with pairs",
			},
			{
				Name:        "ChildOf",
				Description: "Allows parent-child relationships between entities",
				CanBePair:   true,
				Fields: []*geckpb.FieldDefinition{
					{
						Name: "Parent",
						ResetValue: &geckpb.FieldDefinition_U32{
							U32: 0,
						},
					},
				},
			},
			{
				Name:        "InstanceOf",
				Description: "Allows instance-of relationships between entities",
				CanBePair:   true,
				Fields: []*geckpb.FieldDefinition{
					{
						Name: "Value",
						ResetValue: &geckpb.FieldDefinition_U32{
							U32: 0,
						},
					},
				},
			},
			{
				Name:        "Component",
				Description: "A component tag",
			},
			{
				Name:        "Tag",
				Description: "A tag",
			},
		},
	}
)
