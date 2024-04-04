syntax = "proto3";

package {{.PackageName}}.v1;

option go_package = "{{.PackageName}}/pb/gen/{{.PackageName}}/v1;{{.PackageName}}pb";

import "google/protobuf/empty.proto";

{{if .Enums -}}
// Enums
{{range .Enums -}}
enum {{.Name.Singular.Pascal}}Enum {
  {{- range .Values}}
  {{.Name.Singular.Snake | upper }} = {{.Value}};
  {{- end}}
}
{{end -}}
{{end }}

// Components
{{range .Components -}}
{{if not .IsTag -}}
message {{.Name.Singular.Pascal}}Component {
{{- range $i, $f := .Fields }}
  {{$f.PBType}} {{$f.Name.Singular.Snake}} = {{add $i 1}};
{{- end }}
}
{{end }}
{{end -}}


message WorldPatch {
  map<uint32, google.protobuf.Empty> entities = 1;
  {{range $i,$c := .Components -}}
  {{ $v := add $i 2 -}}
  {{if .IsTag -}}
  map<uint32, google.protobuf.Empty> {{.Name.Singular.Snake}}_tags = {{$v}};
  {{else -}}
  map<uint32, {{.Name.Singular.Pascal}}Component> {{.Name.Singular.Snake}}_components = {{$v}};
  {{end -}}
  {{end -}}
}