package {{.PackageName}}

import (
    "github.com/goccy/go-json"
    "unsafe"
    "log"
    "fmt"
)

{{if .Description}}
// {{.Description}}
{{ end -}}
type {{.Name.Pascal}} struct {
{{- range .Fields -}}
    {{- if len .Description}}
    // {{.Description}}{{ end}}
    {{.Name.Pascal}} {{.Type}} `json:"{{.Name.Camel}},omitempty"`
{{- end }}
}

const (
    {{.Name.Pascal}}ID ID = {{.ID}}
    {{.Name.Pascal}}Name = "{{.Name.Pascal}}"
    {{.Name.Pascal}}SizeBytes = {{.SizeBytes}}
    {{.Name.Pascal}}IsTag = {{eq (len .Fields) 0}}
    {{.Name.Pascal}}IsBuiltin = {{.IsBuiltin}}
)

var (
    {{.Name.Pascal}}ResetValue = {{.Name.Pascal}}{
{{- range .Fields}}
        {{.Name.Pascal}}: {{.ResetValue}},
{{- end}}
    }
    {{.Name.Pascal}}ByteOffsets = []int{
{{- range .Fields}}
        {{.OffsetBytes}}, // {{.Name.Pascal}}
{{- end}}
    }
    {{.Name.Pascal}}IDSet = NewIDSet( {{.Name.Pascal}}ID )
)

func (c *{{.Name.Pascal}}) Copy(other {{.Name.Pascal}}) {
{{- range .Fields}}
    c.{{.Name.Pascal}} = other.{{.Name.Pascal}}
{{- end}}
}

func (c *{{.Name.Pascal}}) Reset() {
    c.Copy({{.Name.Pascal}}ResetValue)
}

func (c *{{.Name.Pascal}}) Bytes() []byte {
    // copy the data into the pointer using unsafe
    buf := make([]byte, {{.Name.Pascal}}SizeBytes)
	ptr := unsafe.Pointer(&buf[0])
    tPtr := (*{{.Name.Pascal}})(ptr)
    *tPtr = *c
    return buf
}

func (c *{{.Name.Pascal}}) FromBytes(buf []byte) {
    if len(buf) < {{.Name.Pascal}}SizeBytes {
        panic("Invalid buffer size")
    }
    p := unsafe.Pointer(unsafe.SliceData(buf))
    ptr := (*{{.Name.Pascal}})(p)
	*c = *ptr
}

func (c *{{.Name.Pascal}}) ToJSON() ([]byte, error) {
    return json.Marshal(c)
}

func (c *{{.Name.Pascal}}) FromJSON(data []byte) error {
    return json.Unmarshal(data, c)
}

func {{.Name.Pascal}}FromJSON(data []byte) (*{{.Name.Pascal}}, error) {
    c := &{{.Name.Pascal}}{}
    if err := c.FromJSON(data); err != nil {
        return nil, err
    }
    return c, nil
}

{{ if not .IsBuiltin -}}
func (w *World) Set{{.Name.Pascal}}s(component {{.Name.Pascal}}, entities ...ID) {
    setComponentData(w, {{.Name.Pascal}}ID, component, entities...)
}

func (w *World) {{.Name.Pascal}}(entity ID) *{{.Name.Pascal}} {
    data := &{{.Name.Pascal}}{}
    w.{{.Name.Pascal}}CopyTo(entity, data)
    return data
}

func (w *World) {{.Name.Pascal}}CopyTo(entity ID, copyTo *{{.Name.Pascal}}) {
    componentDataFromEntity(w, {{.Name.Pascal}}ID, entity, copyTo)
}

func (w *World) Remove{{.Name.Pascal}}(entities ...ID) {
    removeComponentFrom(w,  {{.Name.Pascal}}IDSet, NewIDSet(entities...))
}
{{- end }}

func (w *World) MarshalAll{{.Name.Pascal}}() {

    log.Print("Marshaling {{.Name.Pascal}}")
    hasValues := false
    for _, a := range w.archetypes {
        if len(a.entities) == 0 {
			continue
		}

		count := a.componentIDs.Cardinality()
		if count == 0 {
			continue
		}

		validCIDs := []ID{}
        validNames := []string{}
		a.componentIDs.Range(func(cid ID) {
			source, target, _ := cid.SplitPair()
            if source == {{.Name.Pascal}}ID || target == {{.Name.Pascal}}ID {
				validCIDs = append(validCIDs, cid)

                sn := w.EntityName(source)
                if sn == "" {
                    sn = "_"
                }
                tn := w.EntityName(target)
                if tn == "" {
                    tn = "_"
                }
                n := fmt.Sprintf("%s,%s", sn, tn)
                validNames = append(validNames, n)
			}

		})
		if len(validCIDs) == 0 {
			continue
		}

        {{- if not .IsTag}}
		colIndicies := w.archetypeComponentColumnIndicies[a.hash]
		for i, cid := range validCIDs {
            log.Printf("Marshaling {{.Name.Pascal}} for %s", validNames[i])

			colIdx, ok := colIndicies[cid]
			if !ok || colIdx < 0 {
				continue
			}

			col := a.dataColumns[colIdx]
            c := &{{.Name.Pascal}}{}

			for i, e := range a.entities {
				start := uintptr(i) * col.metadata.elementSize
				end := start + col.metadata.elementSize
				buf := col.data[start:end]
				c.FromBytes(buf)

                log.Printf("%d : %+v", e,c)
                hasValues = true
			}
		}
        {{- else}}
        for i := range validCIDs {
            log.Printf("Marshaling {{.Name.Pascal}} for %s", validNames[i])
            for _, e := range a.entities {
                log.Printf("%d", e)
                hasValues = true
            }
        }
        {{- end}}

    }

    if !hasValues {
        log.Print("No values for {{.Name.Pascal}}")
    }
}

