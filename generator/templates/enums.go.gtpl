package {{.PackageName}}

import (
    {{.PackageName}}pb "{{.PBImportPath}}"
)

{{$nsp := .Name.Singular.Pascal -}}
{{$ensp := printf "Enum%s" $nsp -}}
{{$nspe := printf "%sEnum" $nsp -}}
type {{$ensp}} int

const (
	{{range .Values -}}
	{{$ensp}}{{.Name.Singular.Pascal}} = {{.Value}}
	{{end}}
)

func (e {{$ensp}}) String() (string,bool) {
	switch e {
	{{range .Values -}}
	case {{$ensp}}{{.Name.Singular.Pascal}}:
		return "{{.Name.Singular.Pascal}}", true
	{{end}}
	default:
		return "", false
	}
}

func (e {{$ensp}}) ToInt() int {
	return int(e)
}

func {{$ensp}}FromInt(i int) {{$ensp}} {
	return {{$ensp}}(i)
}

func (e {{$ensp}}) ToPB() ({{.PackageName}}pb.{{$nspe}}) {
	return {{.PackageName}}pb.{{$nspe}}(e.ToInt())
}

{{ if .IsBitmask -}}
func {{$ensp}}Set(flags ...{{$ensp}}) {{$ensp}} {
	var e {{$ensp}}
	for _, flag := range flags {
		e |= flag
	}
	return e
}

func (e {{$ensp}}) Has(flags ...{{$ensp}}) bool {
	for _, flag := range flags {
		if e & flag == 0 {
			return false
		}
	}
	return true
}

func (e {{$ensp}}) Set(flags...{{$ensp}}) {{$ensp}} {
	for _, flag := range flags {
		e |= flag
	}
	return e
}

func (e {{$ensp}}) Clear(flags...{{$ensp}}) {{$ensp}} {
	for _, flag := range flags {
		e &= ^flag
	}
	return e
}

func (e {{$ensp}}) Toggle(flags...{{$ensp}}) {{$ensp}} {
	for _, flag := range flags {
		e ^= flag
	}
	return e
}

func (e {{$ensp}}) ToggleAll() {{$ensp}} {
	return e ^ {{$ensp}}Set({{range .Values -}}
		{{$ensp}}{{.Name.Singular.Pascal}},
	{{end}})
}

func (e {{$ensp}}) AllSet() (flags []{{$ensp}}) {

	{{range .Values -}}
	if e & {{$ensp}}{{.Name.Singular.Pascal}} != 0 {
		flags = append(flags, {{$ensp}}{{.Name.Singular.Pascal}})
	}
	{{end}}
	return flags
}
{{end -}}
