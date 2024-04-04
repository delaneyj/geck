package {{.PackageName}}

import (
    "github.com/RoaringBitmap/roaring"
)

// {{.Description}}
func(w *World) {{.Name}}(
    fn func(
        {{range .Params -}}
        {{.VariableName}}Entity Entity,
        {{end -}}
    ) error,
) error {

    // variables := map[string]Entity{}

    {{ range .Steps }}
    // Step {{.Index}} [{{.FromStore.Singular.Pascal}}] - "{{.Expression}}" {{if .ToEntityVariable}}--> {{.ToEntityVariable}}{{end}}
    {{ end -}}

    return nil
}
