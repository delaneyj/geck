package {{.PackageName}}

import (
    "github.com/btvoidx/mint"
)

{{$np := .Name.Singular.Pascal -}}
{{$nc := .Name.Singular.Camel -}}

{{ if not .IsTag}}
{{ if and .IsOnlyOneField -}}
{{ $f := .Fields | first -}}
{{ $fp := $f.Type.Singular.Pascal -}}
{{ $fo := $f.Type.Singular.Original -}}

type {{$np}} {{$fo}}

{{ if .IsFirstSlice -}}
func {{$np}}From{{$fp}}(c {{ $fo }}) {{$np}} {
    return {{$np}}(c)
}

func (c {{$np}}) To{{$fp}}() {{ $fo }} {
    return {{ $fo }}(c)
}
{{ else -}}
func {{$np}}From{{$fp}}(c {{ $fo }}) {{$np}} {
    return {{$np}}(c)
}

func (c {{$np}}) To{{$fp}}() {{ $fo }} {
    return {{ $fo }}(c)
}

func (w *World) Reset{{.Name.Singular.Pascal}}() {{$fo}}{
    return {{.ResetValue}}
}

{{ end -}}
{{else -}}
type {{$np}} struct {
    {{- range .Fields }}
    {{.Name.Singular.Pascal}} {{.Type.Singular.Original}} `json:"{{.Name.Singular.Camel}}"`
    {{- end }}
}

func (w *World) Reset{{.Name.Singular.Pascal}}() {{.Name.Singular.Pascal}}{
    return {{.Name.Singular.Pascal}}{
        {{- range .Fields }}
        {{.Name.Singular.Pascal}}: {{.ResetValue}},
        {{- end }}
    }
}

{{end -}}

{{ if and .IsFirstFieldEntity .IsOnlyOneField  }}
{{range .Fields}}
{{ if (hasSuffix "Entity" .Type.Singular.Pascal)  }}
{{ if .IsSlice  -}}
func (c {{$np}}) ToEntities() []Entity {
    entities := make([]Entity, len(c))
    copy(entities, c)
    return entities
}

func {{$np}}FromEntities(e ...Entity) {{$np}} {
    c := make({{$np}}, len(e))
    copy(c, e)
    return c
}
{{ else -}}


func (c {{$np}}) FromEntity(e Entity) {{$np}} {
    return {{$np}}(e)
}
{{ end -}}
{{ end -}}

{{end}}
{{end -}}

{{else -}}
type {{$np}} struct {}
{{ end }}

//#region Events
{{if .ShouldGenAdded -}}
type {{$np}}AddedEvent struct {
    Entity Entity
    {{$np}} {{$np}}
}
func (w *World) On{{$np}}Added(fn func({{$np}}AddedEvent)) UnsubscribeFunc {
    stopCh := mint.On(w.eventBus, fn)
	return func() { stopCh() }
}
{{end -}}
{{if .ShouldGenRemoved -}}
type {{$np}}RemovedEvent struct {
    Entity Entity
}
func (w *World) On{{$np}}Removed(fn func({{$np}}RemovedEvent)) UnsubscribeFunc {
	stopCh := mint.On(w.eventBus, fn)
	return func() { stopCh() }
}
{{end -}}
{{if .ShouldGenChanged -}}
type {{$np}}ChangedEvent struct {
    Entity Entity
    {{$np}} {{$np}}
}
func (w *World) On{{$np}}Changed(fn func({{$np}}ChangedEvent)) UnsubscribeFunc {
	stopCh := mint.On(w.eventBus, fn)
	return func() { stopCh() }
}
{{end -}}
//#endregion

{{if not .IsTag -}}
func (e Entity) Has{{$np}}() bool {
{{else -}}
func (e Entity) Has{{$np}}Tag() bool {
{{end -}}
    return e.w.{{.Name.Plural.Camel}}Store.Has(e)
}

{{if not .IsTag -}}
{{ if and .IsOnlyOneField .IsFirstFieldEntity -}}
{{ if .IsFirstSlice -}}
func (e Entity) Read{{$np}}() ([]Entity, bool) {
    entities, ok := e.w.{{.Name.Plural.Camel}}Store.Read(e)
    if !ok {
        return nil, false
    }
    return entities, true
}

func (e Entity) {{$np}}Contains(other Entity) bool {
    entities, ok := e.w.{{.Name.Plural.Camel}}Store.Read(e)
    if !ok {
        return false
    }
    for _, entity := range entities {
        if entity == other {
            return true
        }
    }
    return false
}

func (e Entity) Remove{{$np}}(toRemove ...Entity) Entity {
    entities, ok := e.w.{{.Name.Plural.Camel}}Store.Read(e)
    if !ok {
        return e
    }
    clean := make([]Entity, 0, len(entities))
    for _, tr := range toRemove {
        for _, entity := range entities {
            if entity != tr {
                clean = append(clean, entity)
            }
        }
    }
    e.w.{{.Name.Plural.Camel}}Store.Set(clean, e)
    {{if .ShouldGenRemoved}}
    for _, entity := range toRemove {
        fireEvent(e.w, {{$np}}RemovedEvent{ entity })
    }
    {{end}}
    {{if .OwnedBySet -}}
    e.w.{{.OwnedBySet.Name.Singular.Pascal}}.PossibleUpdate(entities...)
    {{end -}}
    return e
}

func (e Entity) RemoveAll{{$np}}() Entity {
    e.w.{{.Name.Plural.Camel}}Store.Remove(e)
    {{if .ShouldGenRemoved}}
    fireEvent(e.w, {{$np}}RemovedEvent{e })
    {{end}}
    return e
}

{{ else -}}
func (e Entity) Read{{$np}}() (Entity, bool) {
    val,ok := e.w.{{.Name.Plural.Camel}}Store.Read(e)
    if !ok {
        return Entity{}, false
    }
    return Entity(val), true
}

func (e Entity) Remove{{$np}}() Entity {
    e.w.{{.Name.Plural.Camel}}Store.Remove(e)
    {{if .ShouldGenRemoved}}
    fireEvent(e.w, {{$np}}RemovedEvent{e })
    {{end}}
    return e
}
{{ end -}}
{{else}}
func (e Entity) Read{{$np}}() ({{$np}}, bool) {
    return e.w.{{.Name.Plural.Camel}}Store.Read(e)
}

func (e Entity) Remove{{$np}}() Entity {
    e.w.{{.Name.Plural.Camel}}Store.Remove(e)
    {{if .ShouldGenRemoved}}
    fireEvent(e.w, {{$np}}RemovedEvent{e })
    {{end}}
    return e
}
{{end }}


func (e Entity) Writable{{$np}}() (*{{$np}}, bool) {
    return e.w.{{.Name.Plural.Camel}}Store.Writeable(e)
}

{{ if and .IsOnlyOneField .IsFirstFieldEntity -}}
func (e Entity) Set{{$np}}(other {{if .IsFirstSlice}}...{{end}}Entity) Entity {
    e.w.{{.Name.Plural.Camel}}Store.Set({{$np}}(other), e)
{{ else -}}
func (e Entity) Set{{$np}}(other {{$np}}) Entity {
    e.w.{{.Name.Plural.Camel}}Store.Set(other,e)
{{ end -}}
    {{if .ShouldGenChanged}}
    fireEvent(e.w, {{$np}}ChangedEvent{e, {{$np}}(other)})
    {{end}}
    {{if .OwnedBySet -}}
    e.w.{{.OwnedBySet.Name.Singular.Pascal}}.PossibleUpdate(e)
    {{end -}}
    return e
}

func (w *World) Set{{.Name.Plural.Pascal}}(c {{$np}}, entities ...Entity) {
    w.{{.Name.Plural.Camel}}Store.Set(c, entities...)
    {{if .ShouldGenChanged -}}
    for _, entity := range entities {
        fireEvent(w, {{$np}}ChangedEvent{entity, c})
    }
    {{end -}}
    {{if .OwnedBySet -}}
    w.{{.OwnedBySet.Name.Singular.Pascal}}.PossibleUpdate(entities...)
    {{end -}}
}

{{else -}}
func (e Entity) TagWith{{$np}}() Entity {
    e.w.{{.Name.Plural.Camel}}Store.Set(e.w.{{.Name.Plural.Camel}}Store.zero, e)
    {{if .ShouldGenAdded -}}
    fireEvent(e.w, {{$np}}AddedEvent{e, {{$np}}{}})
    {{end -}}
    {{if .OwnedBySet -}}
    e.w.{{.OwnedBySet.Name.Singular.Pascal}}.PossibleUpdate(e)
    {{end -}}
    return e
}

func (e Entity) Remove{{$np}}Tag() Entity {
    e.w.{{.Name.Plural.Camel}}Store.Remove(e)
    {{if .ShouldGenRemoved -}}
    fireEvent(e.w, {{$np}}RemovedEvent{e })
    {{end -}}
    {{if .OwnedBySet -}}
    e.w.{{.OwnedBySet.Name.Singular.Pascal}}.PossibleUpdate(e)
    {{end -}}
    return e
}
{{ end }}

{{if not .IsTag -}}
func (w *World) Remove{{.Name.Plural.Pascal}}(entities ...Entity) {
{{else -}}
func (w *World) Remove{{.Name.Singular.Pascal}}Tags(entities ...Entity) {
{{end -}}
    w.{{.Name.Plural.Camel}}Store.Remove(entities...)
    {{if .ShouldGenRemoved -}}
    for _, entity := range entities {
        fireEvent(w, {{$np}}RemovedEvent{entity })
    }
    {{end -}}
    {{if .OwnedBySet -}}
    w.{{.OwnedBySet.Name.Singular.Pascal}}.PossibleUpdate(entities...)
    {{end -}}
}

{{if not .IsTag -}}
//#region Resources

// Has{{$np}} checks if the world has a {{$np}}}}
func (w *World) Has{{$np}}Resource() bool {
    return w.resourceEntity.Has{{$np}}()
}

// Retrieve the {{$np}} resource from the world
{{if and .IsOnlyOneField .IsFirstFieldEntity -}}
{{ if .IsFirstSlice -}}
func (w *World) {{$np}}Resource() ([]Entity, bool) {
{{else -}}
func (w *World) {{$np}}Resource() (Entity, bool) {
{{end -}}
{{else -}}
func (w *World) {{$np}}Resource() ({{$np}}, bool) {
{{end -}}
    return w.resourceEntity.Read{{$np}}()
}

// Set the {{$np}} resource in the world
{{if and .IsOnlyOneField .IsFirstFieldEntity -}}
func (w *World) Set{{$np}}Resource(c {{if .IsFirstSlice}}...{{end}}Entity) Entity {
    w.resourceEntity.Set{{$np}}(c{{if .IsFirstSlice}}...{{end}})
    {{if .ShouldGenChanged}}
    fireEvent(w, {{$np}}ChangedEvent{w.resourceEntity, {{$np}}(c)})
    {{end}}
    {{if .OwnedBySet -}}
    e.w.{{.OwnedBySet.Name.Singular.Pascal}}.PossibleUpdate(w.resourceEntity)
    {{end -}}
    return w.resourceEntity
}
{{else -}}
func (w *World) Set{{$np}}Resource(c {{$np}}) Entity {
    w.resourceEntity.Set{{$np}}(c)
    {{if .ShouldGenChanged -}}
    fireEvent(w, {{$np}}ChangedEvent{w.resourceEntity, c})
    {{end -}}
    {{if .OwnedBySet -}}
    w.{{.OwnedBySet.Name.Singular.Pascal}}.PossibleUpdate(w.resourceEntity)
    {{end -}}
    return w.resourceEntity
}
{{end -}}

// Remove the {{$np}} resource from the world
func (w *World) Remove{{$np}}Resource() Entity {
    w.resourceEntity.Remove{{$np}}()
    {{if .ShouldGenRemoved}}
    fireEvent(w, {{$np}}RemovedEvent{w.resourceEntity })
    {{end}}
    {{if .OwnedBySet -}}
    w.{{.OwnedBySet.Name.Singular.Pascal}}.PossibleUpdate(w.resourceEntity)
    {{end -}}
    return w.resourceEntity
}

//#endregion
{{end }}

//#region Iterators

type {{.Name.Singular.Pascal}}ReadIterator struct {
    w *World
    currIdx int
    store *SparseSet[{{.Name.Singular.Pascal}}]
}

func (iter *{{.Name.Singular.Pascal}}ReadIterator) HasNext() bool {
    return iter.currIdx < iter.store.Len()
}

func (iter *{{.Name.Singular.Pascal}}ReadIterator) NextEntity() Entity {
    e := iter.store.dense[iter.currIdx]
    iter.currIdx++
    return e
}

{{if not .IsTag -}}
func (iter *{{.Name.Singular.Pascal}}ReadIterator) Next{{.Name.Singular.Pascal}}() (Entity, {{.Name.Singular.Pascal}}) {
    e := iter.store.dense[iter.currIdx]
    c := iter.store.components[iter.currIdx]
    iter.currIdx++
    return e, c
}
{{end }}

func (iter *{{.Name.Singular.Pascal}}ReadIterator) Reset() {
    iter.currIdx = 0
}

func (w *World) {{.Name.Singular.Pascal}}ReadIter() *{{.Name.Singular.Pascal}}ReadIterator {
    iter := &{{.Name.Singular.Pascal}}ReadIterator{
        w: w,
        store: w.{{.Name.Plural.Camel}}Store,
    }
    iter.Reset()
    return iter
}

type {{.Name.Singular.Pascal}}WriteIterator struct {
    w *World
    currIdx int
    store *SparseSet[{{.Name.Singular.Pascal}}]
}

func (iter *{{.Name.Singular.Pascal}}WriteIterator) HasNext() bool {
    return iter.currIdx >= 0
}

func (iter *{{.Name.Singular.Pascal}}WriteIterator) NextEntity() Entity {
    e := iter.store.dense[iter.currIdx]
    iter.currIdx--

    return e
}

{{if not .IsTag -}}
func (iter *{{.Name.Singular.Pascal}}WriteIterator) Next{{.Name.Singular.Pascal}}() (Entity, *{{.Name.Singular.Pascal}}) {
    e := iter.store.dense[iter.currIdx]
    c := &iter.store.components[iter.currIdx]
    iter.currIdx--

    return e, c
}
{{end }}

func (iter *{{.Name.Singular.Pascal}}WriteIterator) Reset() {
    iter.currIdx = iter.store.Len() - 1
}

func (w *World) {{.Name.Singular.Pascal}}WriteIter() *{{.Name.Singular.Pascal}}WriteIterator {
    iter := &{{.Name.Singular.Pascal}}WriteIterator{
        w: w,
        store: w.{{.Name.Plural.Camel}}Store,
    }
    iter.Reset()
    return iter
}

//#endregion

func (w *World) {{.Name.Singular.Pascal}}Entities() []Entity {
    return w.{{.Name.Plural.Camel}}Store.entities()
}

{{if not .IsTag -}}
func (w *World) Set{{.Name.Singular.Pascal}}SortFn(lessThan func(a, b Entity) bool) {
    w.{{.Name.Plural.Camel}}Store.LessThan = lessThan
}

func (w *World) Sort{{.Name.Plural.Pascal}}() {
    w.{{.Name.Plural.Camel}}Store.Sort()
}
{{end -}}

