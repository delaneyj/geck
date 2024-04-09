package generator

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
	"time"

	"golang.org/x/mod/modfile"

	sprig "github.com/Masterminds/sprig/v3"
	geckpb "github.com/delaneyj/geck/pb/gen/geck/v1"
	"github.com/delaneyj/toolbelt"
	"github.com/go-openapi/inflect"
	"github.com/samber/lo"
)

//go:embed templates/*
var templatesFS embed.FS

type InflectionString struct {
	Singular toolbelt.CasedString
	Plural   toolbelt.CasedString
}

type enumEntryTmplData struct {
	Name  InflectionString
	Value int
}

type enumTmplData struct {
	PackageName  string
	PBImportPath string
	Name         InflectionString
	Values       []*enumEntryTmplData
	IsBitmask    bool
}
type ecsTmplData struct {
	PackageName   string
	PBImportPath  string
	FolderPath    string
	Enums         []*enumTmplData
	Components    []*componentTmplData
	ComponentSets []*componentSetTmplData
}
type fieldTemplateData struct {
	Name               InflectionString
	Type               InflectionString
	PBType, PBFromType string
	PBTypeSingular     string
	PBNeedsCast        bool
	Description        string
	ResetValue         string
	IsSlice, IsEntity  bool
}

type componentTmplData struct {
	PackageName                                        string
	PBImportPath                                       string
	BundleName                                         toolbelt.CasedString
	Name                                               InflectionString
	Description                                        string
	Fields                                             []fieldTemplateData
	IsTag                                              bool
	IsOnlyOneField, IsFirstFieldEntity, IsFirstSlice   bool
	ShouldGenAdded, ShouldGenRemoved, ShouldGenChanged bool
	ResetValue                                         string
	OwnedBySet                                         *componentSetTmplData
}

type componentSetEntryTmplData struct {
	Name       InflectionString
	IsWritable bool
}

type componentSetTmplData struct {
	PackageName            string
	Name                   InflectionString
	HasWriteableComponents bool
	OwnedComponents        []*componentSetEntryTmplData
	BorrowedComponents     []*componentSetEntryTmplData
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

	templatesSubFS, err := fs.Sub(templatesFS, "templates")
	if err != nil {
		return fmt.Errorf("failed to access templates: %w", err)
	}
	tmpls, err := template.New("root").
		Funcs(sprig.FuncMap()).
		ParseFS(templatesSubFS, "*.gtpl")
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}
	tmpls.Funcs(sprig.FuncMap())

	log.Printf("Converting options to data")
	data, err := optsToData(opts)
	if err != nil {
		return fmt.Errorf("failed to convert options to data: %w", err)
	}

	log.Printf("Generating universal files")
	if err := generateFiles(tmpls, data,
		"entities.go",
		"events.go",
		"sparse_sets.go",
		"sparse_sets_test.go",
		"sparse_sets_timsort.go",
		"world.go",
	); err != nil {
		return fmt.Errorf("failed to generate top level files: %w", err)
	}

	log.Printf("Generating enum files")
	for _, enum := range data.Enums {
		log.Printf("\t%s", enum.Name.Plural.Pascal)
		if err := generateEnum(tmpls, data, enum); err != nil {
			return fmt.Errorf("failed to generate bundle: %w", err)
		}
	}

	log.Printf("Generating component files")
	for _, component := range data.Components {
		log.Printf("\t%s", component.Name.Plural.Pascal)
		if err := generateComponent(tmpls, data, component); err != nil {
			return fmt.Errorf("failed to generate bundle: %w", err)
		}
	}

	log.Printf("Generating component set files")
	for _, set := range data.ComponentSets {
		log.Printf("\t%s", set.Name.Plural.Pascal)
		if err := generateComponentSet(tmpls, data, set); err != nil {
			return fmt.Errorf("failed to generate bundle: %w", err)
		}
	}

	log.Printf("Generate protobufs")
	if err := generateProtobufs(tmpls, data); err != nil {
		return fmt.Errorf("failed to generate protobufs: %w", err)
	}

	log.Printf("Running post generation commands")
	type postGenCmdSet struct {
		subdir string
		cmds   []string
	}
	postGenCmds := []postGenCmdSet{
		{
			subdir: "pb",
			cmds: []string{
				"go install github.com/bufbuild/buf/cmd/buf@latest",
				"clang-format -i ecs/v1/ecs.proto",
				"buf mod update",
				"buf generate",
			},
		},
		{
			subdir: "",
			cmds: []string{
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

	startPath, err := filepath.Abs(".")
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}
	currPath, err := filepath.Abs(opts.FolderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}
	var goModPath, goModFolder string
	for goModPath == "" || currPath != "/" {
		goModPossiblePath := filepath.Join(currPath, "go.mod")
		if _, err := os.Stat(goModPossiblePath); err == nil {
			goModPath = goModPossiblePath
			goModFolder = currPath
			break
		}
		currPath = filepath.Dir(currPath)
	}

	goModBytes, err := os.ReadFile(goModPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read go.mod: %w", err)
	}
	modFile, err := modfile.Parse("go.mod", goModBytes, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse go.mod: %w", err)
	}
	genFolder := filepath.Join(startPath, opts.FolderPath)

	pbFolder := filepath.Join(genFolder, "pb", "gen", opts.PackageName, fmt.Sprintf("v%d", opts.Version))
	pbBasePath, err := filepath.Rel(goModFolder, pbFolder)
	if err != nil {
		return nil, fmt.Errorf("failed to get relative path: %w", err)
	}

	data.PBImportPath = filepath.Join(modFile.Module.Mod.Path, pbBasePath)

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

	for _, ed := range opts.Enums {
		if len(ed.Values) == 0 {
			return nil, fmt.Errorf("enum must have at least one value")
		}

		enum := &enumTmplData{
			PackageName:  opts.PackageName,
			PBImportPath: data.PBImportPath,
			Name:         inflectionStrings(ed.Name, true),
			IsBitmask:    ed.IsBitmask,
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
					Name:  inflectionStrings(enum.Name.Singular.Pascal+"_Unknown", false),
					Value: 0,
				},
			}, enum.Values...)
		}

		data.Enums = append(data.Enums, enum)
	}

	componentByNames := map[string]*componentTmplData{}
	for _, bundleDef := range opts.Bundles {
		bundleName := toolbelt.ToCasedString(bundleDef.Name)
		for _, cd := range bundleDef.Components {
			isTag := len(cd.Fields) == 0

			component := &componentTmplData{
				BundleName:       bundleName,
				PackageName:      opts.PackageName,
				PBImportPath:     data.PBImportPath,
				Name:             inflectionStrings(cd.Name, !isTag && !cd.ShouldNotInflect),
				Description:      cd.Description,
				IsTag:            isTag,
				ShouldGenAdded:   cd.ShouldGenerateAddedEvent,
				ShouldGenRemoved: cd.ShouldGenerateRemovedEvent,
				ShouldGenChanged: cd.ShouldGenerateChangedEvent,
			}

			componentByNames[cd.Name] = component

			if len(cd.Fields) == 1 {
				component.IsOnlyOneField = true
				switch cd.Fields[0].ResetValue.(type) {
				case *geckpb.FieldDefinition_Entity:
					component.IsFirstFieldEntity = true

					if cd.Fields[0].HasMultiple {
						component.IsFirstSlice = true
					}
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
					ftd.PBType = "uint32"
					ftd.ResetValue = fmt.Sprintf("%d", f.GetU8())
					ftd.PBNeedsCast = true
				case *geckpb.FieldDefinition_U16:
					typ = "uint16"
					ftd.PBType = "uint32"
					ftd.ResetValue = fmt.Sprintf("%d", f.GetU16())
					ftd.PBNeedsCast = true
				case *geckpb.FieldDefinition_U32:
					typ = "uint32"
					ftd.ResetValue = fmt.Sprintf("%d", f.GetU32())
				case *geckpb.FieldDefinition_U64:
					typ = "uint64"
					ftd.ResetValue = fmt.Sprintf("%d", f.GetU64())
				case *geckpb.FieldDefinition_I8:
					typ = "int8"
					ftd.PBType = "sint32"
					ftd.PBFromType = "int32"
					ftd.ResetValue = fmt.Sprintf("%d", f.GetI8())
					ftd.PBNeedsCast = true
				case *geckpb.FieldDefinition_I16:
					typ = "int16"
					ftd.PBType = "sint32"
					ftd.PBFromType = "int32"
					ftd.ResetValue = fmt.Sprintf("%d", f.GetI16())
					ftd.PBNeedsCast = true
				case *geckpb.FieldDefinition_I32:
					typ = "int32"
					ftd.PBType = "sint32"
					ftd.PBFromType = "int32"
					ftd.ResetValue = fmt.Sprintf("%d", f.GetI32())
				case *geckpb.FieldDefinition_I64:
					typ = "int64"
					ftd.PBType = "sint64"
					ftd.PBFromType = "int64"
					ftd.ResetValue = fmt.Sprintf("%d", f.GetI64())
				case *geckpb.FieldDefinition_F32:
					typ = "float32"
					ftd.PBType = "float"
					ftd.ResetValue = fmt.Sprintf("%f", f.GetF32())
				case *geckpb.FieldDefinition_F64:
					typ = "float64"
					ftd.PBType = "double"
					ftd.ResetValue = fmt.Sprintf("%f", f.GetF64())
				case *geckpb.FieldDefinition_Txt:
					typ = "string"
					ftd.ResetValue = fmt.Sprintf(`"%s"`, f.GetTxt())
				case *geckpb.FieldDefinition_Bin:
					typ = "[]byte"
					ftd.PBType = "bytes"
					ftd.ResetValue = fmt.Sprintf("[]byte(%v)", f.GetBin())
				case *geckpb.FieldDefinition_Entity:
					typ = "Entity"
					ftd.PBType = "uint32"
					ftd.ResetValue = "w.EntityFromU32(0)"
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
					ftd.PBType = typ + "Enum"
					typ = "Enum" + typ
					ftd.ResetValue = fmt.Sprintf("%s(%d)", typ, e.Value)

				default:
					return nil, fmt.Errorf("unknown field type: %T", f.ResetValue)
				}

				if ftd.PBType == "" {
					ftd.PBType = typ
				}
				ftd.PBTypeSingular = ftd.PBType

				if f.HasMultiple {
					typ = "[]" + typ
					ftd.ResetValue = "nil"
					ftd.PBType = "repeated " + ftd.PBType
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
	}

	for _, sd := range opts.ComponentSets {
		set := &componentSetTmplData{
			PackageName: opts.PackageName,
		}

		componentNames := []string{}
		componentEntry := func(cd *geckpb.ComponentSetDefinition_Component) (*componentSetEntryTmplData, error) {
			componentNames = append(componentNames, cd.Name)
			c, ok := componentByNames[cd.Name]
			if !ok {
				return nil, fmt.Errorf("component not found: %s", cd.Name)
			}
			ce := &componentSetEntryTmplData{
				Name:       c.Name,
				IsWritable: cd.IsWriteable,
			}

			if cd.IsWriteable {
				set.HasWriteableComponents = true
			}

			return ce, nil
		}

		for _, ed := range sd.Owned {
			componentEntry, err := componentEntry(ed)
			if err != nil {
				return nil, err
			}
			set.OwnedComponents = append(set.OwnedComponents, componentEntry)
		}
		for _, ed := range sd.Borrowed {
			componentEntry, err := componentEntry(ed)
			if err != nil {
				return nil, err
			}
			set.OwnedComponents = append(set.OwnedComponents, componentEntry)
		}
		componentNames = lo.Uniq(componentNames)
		slices.Sort(componentNames)
		name := strings.Join(componentNames, "_") + "_set"
		set.Name = inflectionStrings(name, true)

		data.ComponentSets = append(data.ComponentSets, set)
		for _, ce := range set.OwnedComponents {
			c, ok := componentByNames[ce.Name.Singular.Original]
			if !ok {
				return nil, fmt.Errorf("component not found: %s", ce.Name.Singular.Original)
			}

			c.OwnedBySet = set
		}
	}

	return data, nil
}

func generateFiles(tmpls *template.Template, data *ecsTmplData, templateNames ...string) error {
	for _, templateName := range templateNames {
		fn := fmt.Sprintf("ecs_%s", templateName)
		fp := filepath.Join(data.FolderPath, fn)
		f, err := os.Create(fp)
		if err != nil {
			return fmt.Errorf("failed to create world.go: %w", err)
		}
		defer f.Close()

		if err := tmpls.ExecuteTemplate(f, templateName+".gtpl", data); err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}
	}
	return nil
}

func generateEnum(tmpls *template.Template, data *ecsTmplData, enum *enumTmplData) error {
	fp := filepath.Join(
		data.FolderPath,
		fmt.Sprintf(
			"enums_%s.go",
			enum.Name.Plural.Snake,
		),
	)
	enumFile, err := os.Create(fp)
	if err != nil {
		return fmt.Errorf("failed to create enum file: %w", err)
	}
	defer enumFile.Close()

	return tmpls.ExecuteTemplate(enumFile, "enums.go.gtpl", enum)
}

func generateComponent(tmpls *template.Template, data *ecsTmplData, component *componentTmplData) error {
	prefix := "components"
	if component.IsTag {
		prefix = "tags"
	}
	fp := filepath.Join(
		data.FolderPath,
		fmt.Sprintf(
			"%s_%s_%s.go",
			prefix,
			component.BundleName.Camel,
			component.Name.Plural.Snake,
		),
	)
	componentFile, err := os.Create(fp)
	if err != nil {
		return fmt.Errorf("failed to create component file: %w", err)
	}
	defer componentFile.Close()

	return tmpls.ExecuteTemplate(componentFile, "components.go.gtpl", component)
}

func generateComponentSet(tmpls *template.Template, data *ecsTmplData, componentSet *componentSetTmplData) error {
	fp := filepath.Join(
		data.FolderPath,
		fmt.Sprintf(
			"components_%s.go",
			toolbelt.Snake(componentSet.Name.Plural.Snake),
		),
	)
	setFile, err := os.Create(fp)
	if err != nil {
		return fmt.Errorf("failed to create component set file: %w", err)
	}
	defer setFile.Close()

	return tmpls.ExecuteTemplate(setFile, "component_sets.go.gtpl", componentSet)
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
			Fields: []*geckpb.FieldDefinition{
				{
					Name:        "Parent",
					Description: "The parent entity",
					ResetValue:  &geckpb.FieldDefinition_Entity{},
				},
			},
		},
		{
			Name:             "IsA",
			ShouldNotInflect: true,
			Fields: []*geckpb.FieldDefinition{
				{
					Name:        "Prototype",
					Description: "The prototype entity",
					ResetValue:  &geckpb.FieldDefinition_Entity{},
				},
			},
		},
	},
}

func generateProtobufs(tmpls *template.Template, data *ecsTmplData) error {
	pbFolder := filepath.Join(data.FolderPath, "pb")
	if err := os.MkdirAll(pbFolder, 0755); err != nil {
		return fmt.Errorf("failed to create pb folder: %w", err)
	}

	topLevelFiles := []string{
		"buf.gen.yaml",
		"buf.yaml",
	}

	for _, filename := range topLevelFiles {
		fp := filepath.Join(
			pbFolder,
			filename,
		)
		f, err := os.Create(fp)
		if err != nil {
			return fmt.Errorf("failed to create protobuf file: %w", err)
		}
		defer f.Close()

		if err := tmpls.ExecuteTemplate(f, filename+".gtpl", data); err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}
	}

	// Generate the protobuf files
	protoPath := filepath.Join(pbFolder, data.PackageName, "v1")
	if err := os.MkdirAll(protoPath, 0755); err != nil {
		return fmt.Errorf("failed to create proto folder: %w", err)
	}

	ecsProto := filepath.Join(protoPath, "ecs.proto")
	ecsProtoFile, err := os.Create(ecsProto)
	if err != nil {
		return fmt.Errorf("failed to create ecs.proto: %w", err)
	}

	if err := tmpls.ExecuteTemplate(ecsProtoFile, "ecs.proto.gtpl", data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
