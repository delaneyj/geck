package generator

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	geckpb "github.com/delaneyj/geck/pb/gen/geck/v1"
	"github.com/delaneyj/toolbelt"
	"github.com/go-openapi/inflect"
	"github.com/samber/lo"
)

type InflectionString struct {
	Singular toolbelt.CasedString
	Plural   toolbelt.CasedString
}

type enumEntryTmplData struct {
	Name  InflectionString
	Value int
}

type enumTmplData struct {
	PackageName string
	Folder      string
	BundleName  toolbelt.CasedString
	Name        InflectionString
	Values      []*enumEntryTmplData
	IsBitmask   bool
}
type ecsTmplData struct {
	PackageName string
	FolderPath  string
	Enums       []*enumTmplData
	Components  []*componentTmplData
	Queries     []*queryTmplData
}
type fieldTemplateData struct {
	Name                 InflectionString
	Type                 InflectionString
	Description          string
	ResetValue           string
	IsSlice, IsEntity    bool
	IsEntityRelationship bool
}

type componentTmplData struct {
	PackageName                                        string
	Folder                                             string
	BundleName                                         toolbelt.CasedString
	PBImportPath                                       string
	Name                                               InflectionString
	Description                                        string
	Fields                                             []fieldTemplateData
	IsTag, IsRelationship                              bool
	IsOnlyOneField, IsFirstFieldEntity, IsFirstSlice   bool
	ShouldGenAdded, ShouldGenRemoved, ShouldGenChanged bool
	HasAnyEvents                                       bool
	ResetValue                                         string
	OwnedBySet                                         *queryTmplData
}

type queryEntryTmplData struct {
	BundleName     toolbelt.CasedString
	Name           InflectionString
	IsMutable      bool
	ComponentOrTag *componentTmplData
}

type queryTmplData struct {
	PackageName string
	Folder      string
	Name        InflectionString
	Entries     []*queryEntryTmplData
}

func BuildECS(ctx context.Context, opts *geckpb.GeneratorOptions) error {
	log.Printf("Building ECS in '%s'", opts.FolderPath)
	start := time.Now()
	defer log.Printf("Finished building ECS in '%s'", opts.FolderPath)

	opts.Bundles = append(
		[]*geckpb.BundleDefinition{builtinBundle},
		opts.Bundles...,
	)

	// Create the folder
	if err := os.RemoveAll(opts.FolderPath); err != nil {
		return fmt.Errorf("failed to remove folder: %w", err)
	}

	if err := os.MkdirAll(opts.FolderPath, 0755); err != nil {
		return fmt.Errorf("failed to create folder: %w", err)
	}

	log.Printf("Converting options to data")
	data, err := optsToData(opts)
	if err != nil {
		return fmt.Errorf("failed to convert options to data: %w", err)
	}

	log.Printf("Generating universal files")
	if err := errors.Join(
		generateFile("world.go", data, worldTemplate),
		generateFile("sparse_set.go", data, sparseSetTemplate),
		generateFile("entities.go", data, entitiesTemplate),
		generateFile("events.go", data, eventsTemplate),
		generateFile("web.go", data, webTemplate),
		generateFile("web_templates.templ", data, templTemplate),
	); err != nil {

		return fmt.Errorf("failed to generate top level files: %w", err)
	}

	log.Printf("Generating enum files")
	for _, enum := range data.Enums {
		if err := generateEnum(enum); err != nil {
			return fmt.Errorf("failed to generate enum: %w", err)
		}
	}

	log.Printf("Generating component files")
	for _, component := range data.Components {
		if err := generateComponent(component); err != nil {
			return fmt.Errorf("failed to generate component: %w", err)
		}
	}

	log.Printf("Generating query files")
	for _, query := range data.Queries {
		if err := generateQueries(query); err != nil {
			return fmt.Errorf("failed to generate query: %w", err)
		}
	}

	log.Printf("Running post generation commands")
	type postGenCmdSet struct {
		subdir string
		cmds   []string
	}
	postGenCmds := []postGenCmdSet{
		{
			subdir: "",
			cmds: []string{
				"go install github.com/valyala/quicktemplate/qtc@latest",
				"templ generate",
				"go install golang.org/x/tools/cmd/goimports@latest",
				"goimports -w .",
				"go mod tidy",
			},
		},
	}
	for _, set := range postGenCmds {
		dir, err := filepath.Abs(filepath.Join(opts.FolderPath, set.subdir))
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}

		for _, cmd := range set.cmds {
			log.Printf("Running: '%s' inside '%s'", cmd, dir)
			parts := strings.Split(cmd, " ")
			c := exec.Command(parts[0], parts[1:]...)
			c.Dir = dir
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			if err := c.Run(); err != nil {
				log.Printf("ERROR! failed to run command: %s", err)
			}
		}
	}

	log.Printf("Generating ECS took %s", time.Since(start))
	return nil
}

func optsToData(opts *geckpb.GeneratorOptions) (data *ecsTmplData, err error) {
	if opts.PackageName == "" {
		opts.PackageName = filepath.Base(opts.FolderPath)
	}

	data = &ecsTmplData{
		PackageName: opts.PackageName,
		FolderPath:  opts.FolderPath,
	}

	inflectionStrings := func(s string, shouldInflect bool) InflectionString {
		singular, plural := s, s
		if shouldInflect {
			singular = inflect.Singularize(s)
			plural = inflect.Pluralize(s)
		}
		return InflectionString{
			Singular: toolbelt.ToCasedString(singular),
			Plural:   toolbelt.ToCasedString(plural),
		}
	}

	componentByNames := map[string]map[string]*componentTmplData{}
	for _, bundleDef := range opts.Bundles {
		bundleName := toolbelt.ToCasedString(bundleDef.Name)

		for _, ed := range bundleDef.Enums {
			if len(ed.Values) == 0 {
				return nil, fmt.Errorf("enum '%s' must have at least one value", ed.Name)
			}

			enum := &enumTmplData{
				PackageName: data.PackageName,
				Folder:      opts.FolderPath,
				BundleName:  bundleName,
				Name:        inflectionStrings(ed.Name, true),
				IsBitmask:   ed.IsBitmask,
				Values: lo.Map(ed.Values, func(v *geckpb.Enum_Value, i int) *enumEntryTmplData {
					return &enumEntryTmplData{
						Name:  inflectionStrings(v.Name, true),
						Value: int(v.Value),
					}
				}),
			}
			slices.SortFunc(enum.Values, func(a, b *enumEntryTmplData) int {
				return a.Value - b.Value
			})

			if len(enum.Values) != len(lo.UniqBy(enum.Values, func(e *enumEntryTmplData) int {
				return e.Value
			})) {
				return nil, fmt.Errorf("enum values must be unique")
			}

			if enum.Values[0].Value != 0 {
				enum.Values = append([]*enumEntryTmplData{
					{
						Name:  inflectionStrings("Unknown", false),
						Value: 0,
					},
				}, enum.Values...)
			}

			data.Enums = append(data.Enums, enum)
		}

		bundleComponentNames := map[string]*componentTmplData{}
		for _, cd := range bundleDef.Components {
			isTag := len(cd.Fields) == 0 && !cd.IsRelationship

			component := &componentTmplData{
				PackageName:      data.PackageName,
				Folder:           opts.FolderPath,
				BundleName:       bundleName,
				Name:             inflectionStrings(cd.Name, !isTag && !cd.ShouldNotInflect),
				Description:      cd.Description,
				IsTag:            isTag,
				IsRelationship:   cd.IsRelationship,
				ShouldGenAdded:   cd.ShouldGenerateAddedEvent,
				ShouldGenRemoved: cd.ShouldGenerateRemovedEvent,
				ShouldGenChanged: cd.ShouldGenerateChangedEvent,
			}

			if component.ShouldGenAdded || component.ShouldGenRemoved || component.ShouldGenChanged {
				component.HasAnyEvents = true
			}

			bundleComponentNames[cd.Name] = component

			if len(cd.Fields) == 1 {
				component.IsOnlyOneField = true
				if cd.Fields[0].HasMultiple {
					component.IsFirstSlice = true
				}
				switch cd.Fields[0].ResetValue.(type) {
				case *geckpb.FieldDefinition_Entity:
					component.IsFirstFieldEntity = true
				}
			}

			for _, f := range cd.Fields {
				ftd := fieldTemplateData{
					Name:        inflectionStrings(f.Name, false),
					Description: f.Description,
					IsSlice:     f.HasMultiple,
				}

				var typ string
				switch f.ResetValue.(type) {
				case *geckpb.FieldDefinition_U8:
					typ = "uint8"
					ftd.ResetValue = fmt.Sprintf("%d", f.GetU8())
				case *geckpb.FieldDefinition_U16:
					typ = "uint16"
					ftd.ResetValue = fmt.Sprintf("%d", f.GetU16())
				case *geckpb.FieldDefinition_U32:
					typ = "uint32"
					ftd.ResetValue = fmt.Sprintf("%d", f.GetU32())
				case *geckpb.FieldDefinition_U64:
					typ = "uint64"
					ftd.ResetValue = fmt.Sprintf("%d", f.GetU64())
				case *geckpb.FieldDefinition_I8:
					typ = "int8"
					ftd.ResetValue = fmt.Sprintf("%d", f.GetI8())
				case *geckpb.FieldDefinition_I16:
					typ = "int16"
					ftd.ResetValue = fmt.Sprintf("%d", f.GetI16())
				case *geckpb.FieldDefinition_I32:
					typ = "int32"
					ftd.ResetValue = fmt.Sprintf("%d", f.GetI32())
				case *geckpb.FieldDefinition_I64:
					typ = "int64"
					ftd.ResetValue = fmt.Sprintf("%d", f.GetI64())
				case *geckpb.FieldDefinition_F32:
					typ = "float32"
					ftd.ResetValue = fmt.Sprintf("%f", f.GetF32())
				case *geckpb.FieldDefinition_F64:
					typ = "float64"
					ftd.ResetValue = fmt.Sprintf("%f", f.GetF64())
				case *geckpb.FieldDefinition_Txt:
					typ = "string"
					ftd.ResetValue = fmt.Sprintf(`"%s"`, f.GetTxt())
				case *geckpb.FieldDefinition_Bin:
					typ = "[]byte"
					ftd.ResetValue = fmt.Sprintf("[]byte(%v)", f.GetBin())
				case *geckpb.FieldDefinition_Entity:
					typ = "Entity"
					ftd.ResetValue = "EntityFromU32(0)"
					ftd.IsEntity = true
				case *geckpb.FieldDefinition_Enum:
					e := f.GetEnum()
					typ = e.Name

					var enum *enumTmplData
					for _, e := range data.Enums {
						if e.Name.Singular.Original == typ {
							enum = e
							break
						}
					}
					if enum == nil {
						return nil, fmt.Errorf("enum not found: %s", f.Name)
					}
					typ = "Enum" + typ
					ftd.ResetValue = fmt.Sprintf("%s(%d)", typ, e.Value)
				default:
					return nil, fmt.Errorf("unknown field type: %s %T", f.Name, f.ResetValue)
				}

				if f.HasMultiple {
					typ = "[]" + typ
					ftd.ResetValue = "nil"
				}

				ftd.Type = inflectionStrings(typ, cd.ShouldNotInflect)
				component.Fields = append(component.Fields, ftd)
			}

			fieldCount := len(component.Fields)
			if fieldCount > 0 {
				if fieldCount == 1 {
					// component.ResetValue = fmt.Sprintf(
					// 	"%s(%s)",
					// 	component.Name.Singular.Pascal,
					// 	component.Fields[0].ResetValue,
					// )
					component.ResetValue = component.Fields[0].ResetValue
				} else {
					sb := strings.Builder{}
					sb.WriteString(component.Name.Singular.Pascal + " {")
					for i, f := range component.Fields {
						sb.WriteString(f.Name.Singular.Pascal + ": " + f.ResetValue)
						if i < fieldCount-1 {
							sb.WriteString(", ")
						}
					}
					sb.WriteString("}")
					component.ResetValue = sb.String()
				}

			}
			data.Components = append(data.Components, component)
		}

		componentByNames[bundleName.Pascal] = bundleComponentNames
	}

	log.Printf("%+v", componentByNames)

	for _, queryDef := range opts.Queries {
		if len(queryDef.Entries) == 0 {
			return nil, fmt.Errorf("query must have at least one component or tag")
		}

		query := &queryTmplData{
			PackageName: opts.PackageName,
			Folder:      opts.FolderPath,
		}

		type Name struct {
			Bundle string
			Name   string
		}
		names := []Name{}
		componentOrTagEntry := func(cd *geckpb.QueryDefinition_ComponentOrTag) (*queryEntryTmplData, error) {
			bundleName := toolbelt.ToCasedString(cd.BundleName)
			bundleNames := componentByNames[bundleName.Pascal]
			if bundleNames == nil {
				return nil, fmt.Errorf("bundle not found: %s", cd.BundleName)
			}

			c, ok := bundleNames[cd.Name]
			if !ok {
				return nil, fmt.Errorf("component not found: %s", cd.Name)
			}

			if cd.IsMutable && c.IsTag {
				return nil, fmt.Errorf("tags cannot be mutable")
			}

			names = append(names, Name{
				Bundle: bundleName.Pascal,
				Name:   c.Name.Singular.Original,
			})
			ce := &queryEntryTmplData{
				BundleName:     bundleName,
				Name:           c.Name,
				IsMutable:      cd.IsMutable,
				ComponentOrTag: c,
			}

			return ce, nil
		}

		for _, def := range queryDef.Entries {
			componentEntry, err := componentOrTagEntry(def)
			if err != nil {
				return nil, err
			}
			query.Entries = append(query.Entries, componentEntry)
		}

		if queryDef.Alias != "" {
			query.Name = inflectionStrings(queryDef.Alias, true)
		} else {
			names = lo.Uniq(names)
			slices.SortFunc(names, func(a, b Name) int {
				bundle := strings.Compare(a.Bundle, b.Bundle)
				if bundle != 0 {
					return bundle
				}

				return strings.Compare(a.Name, b.Name)
			})

			currentBundle := ""
			nameBuilder := strings.Builder{}
			for i, n := range names {
				if n.Bundle != currentBundle {
					if i > 0 {
						nameBuilder.WriteString("_")
					}
					nameBuilder.WriteString(n.Bundle)
					currentBundle = n.Bundle
				}

				if i > 0 {
					nameBuilder.WriteString("_")
				}
				nameBuilder.WriteString(n.Name)
			}
			name := nameBuilder.String()
			query.Name = inflectionStrings(name, true)
		}

		data.Queries = append(data.Queries, query)
	}

	return data, nil
}

func generateFile(templateName string, data *ecsTmplData, templates func(data *ecsTmplData) string) error {
	fn := fmt.Sprintf("ecs_%s", templateName)
	fp := filepath.Join(data.FolderPath, fn)

	contents := templates(data)

	if err := os.WriteFile(fp, []byte(contents), 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func generateEnum(enum *enumTmplData) error {
	fp := filepath.Join(
		enum.Folder,
		fmt.Sprintf(
			"%s_enums_%s.go",
			enum.BundleName.Snake,
			enum.Name.Plural.Snake,
		),
	)
	contents := enumTemplate(enum)
	if err := os.WriteFile(fp, []byte(contents), 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func generateComponent(component *componentTmplData) error {

	var prefix, contents string
	switch {
	case component.IsRelationship:
		prefix = "relationships"
		contents = relationshipTemplate(component)
	case component.IsTag:
		prefix = "tags"
		contents = tagTemplate(component)
	default:
		prefix = "components"
		contents = componentTemplate(component)
	}

	fp := filepath.Join(
		component.Folder,
		fmt.Sprintf(
			"%s_%s_%s.go",
			component.BundleName.Snake,
			prefix,
			component.Name.Plural.Snake,
		),
	)

	if err := os.WriteFile(fp, []byte(contents), 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func generateQueries(query *queryTmplData) error {
	fp := filepath.Join(
		query.Folder,
		fmt.Sprintf(
			"queries_%s.go",
			query.Name.Plural.Snake,
		),
	)

	contents := queryTemplate(query)
	if err := os.WriteFile(fp, []byte(contents), 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

var builtinBundle = &geckpb.BundleDefinition{
	Name:        "Builtin",
	Description: "The built-in bundle",
	Components: []*geckpb.ComponentDefinition{
		{
			Name: "Name",
			Fields: []*geckpb.FieldDefinition{
				{
					Name:        "Value",
					Description: "The name of the entity",
					ResetValue:  &geckpb.FieldDefinition_Txt{},
				},
			},
		},
		{
			Name:             "ChildOf",
			ShouldNotInflect: true,
			IsRelationship:   true,
		},
		{
			Name:             "IsA",
			ShouldNotInflect: true,
			IsRelationship:   true,
		},
	},
}
