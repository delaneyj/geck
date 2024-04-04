package {{.PackageName}}

{{$nsp := .Name.Singular.Pascal -}}
{{$nspe := printf "Enum%s" $nsp -}}
type {{$nspe}} int

const (
	{{range .Values -}}
	{{$nspe}}{{.Name.Singular.Pascal}} = {{.Value}}
	{{end}}
)

func (e {{$nspe}}) String() (string,bool) {
	switch e {
	{{range .Values -}}
	case {{$nspe}}{{.Name.Singular.Pascal}}:
		return "{{.Name.Singular.Pascal}}", true
	{{end}}
	default:
		return "", false
	}
}

func (e {{$nspe}}) ToInt() int {
	return int(e)
}

func {{$nspe}}FromInt(i int) {{$nspe}} {
	return {{$nspe}}(i)
}

{{ if .IsBitmask -}}
func {{$nspe}}Set(flags ...{{$nspe}}) {{$nspe}} {
	var e {{$nspe}}
	for _, flag := range flags {
		e |= flag
	}
	return e
}

func (e {{$nspe}}) Has(flags ...{{$nspe}}) bool {
	for _, flag := range flags {
		if e & flag == 0 {
			return false
		}
	}
	return true
}

func (e {{$nspe}}) Set(flags...{{$nspe}}) {{$nspe}} {
	for _, flag := range flags {
		e |= flag
	}
	return e
}

func (e {{$nspe}}) Clear(flags...{{$nspe}}) {{$nspe}} {
	for _, flag := range flags {
		e &= ^flag
	}
	return e
}

func (e {{$nspe}}) Toggle(flags...{{$nspe}}) {{$nspe}} {
	for _, flag := range flags {
		e ^= flag
	}
	return e
}

func (e {{$nspe}}) ToggleAll() {{$nspe}} {
	return e ^ {{$nspe}}Set({{range .Values -}}
		{{$nspe}}{{.Name.Singular.Pascal}},
	{{end}})
}

func (e {{$nspe}}) AllSet() (flags []{{$nspe}}) {

	{{range .Values -}}
	if e & {{$nspe}}{{.Name.Singular.Pascal}} != 0 {
		flags = append(flags, {{$nspe}}{{.Name.Singular.Pascal}})
	}
	{{end}}
	return flags
}
{{end -}}
