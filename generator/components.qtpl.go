// Code generated by qtc from "components.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

// package generator
//

//line generator/components.qtpl:3
package generator

//line generator/components.qtpl:3
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line generator/components.qtpl:3
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line generator/components.qtpl:3
func streamcomponentTemplate(qw422016 *qt422016.Writer, data *componentTmplData) {
//line generator/components.qtpl:3
	qw422016.N().S(`
package `)
//line generator/components.qtpl:4
	qw422016.E().S(data.PackageName)
//line generator/components.qtpl:4
	qw422016.N().S(`
`)
//line generator/components.qtpl:6
	if data.HasAnyEvents {
//line generator/components.qtpl:6
		qw422016.N().S(`
import "github.com/btvoidx/mint"
`)
//line generator/components.qtpl:8
	}
//line generator/components.qtpl:8
	qw422016.N().S(`

`)
//line generator/components.qtpl:11
	npp := data.Name.Plural.Pascal
	nsp := data.Name.Singular.Pascal
	nsc := data.Name.Singular.Camel
	ss := nsc + "Components"

//line generator/components.qtpl:15
	qw422016.N().S(`
type `)
//line generator/components.qtpl:17
	qw422016.E().S(nsp)
//line generator/components.qtpl:17
	qw422016.N().S(` struct {
`)
//line generator/components.qtpl:18
	for _, f := range data.Fields {
//line generator/components.qtpl:18
		qw422016.N().S(`    `)
//line generator/components.qtpl:19
		qw422016.E().S(f.Name.Singular.Pascal)
//line generator/components.qtpl:19
		qw422016.N().S(` `)
//line generator/components.qtpl:19
		qw422016.E().S(f.Type.Singular.Original)
//line generator/components.qtpl:19
		qw422016.N().S(`
`)
//line generator/components.qtpl:20
	}
//line generator/components.qtpl:20
	qw422016.N().S(`}

func (w *World) Set`)
//line generator/components.qtpl:23
	qw422016.E().S(nsp)
//line generator/components.qtpl:23
	qw422016.N().S(`(e Entity, c `)
//line generator/components.qtpl:23
	qw422016.E().S(nsp)
//line generator/components.qtpl:23
	qw422016.N().S(`) (old `)
//line generator/components.qtpl:23
	qw422016.E().S(nsp)
//line generator/components.qtpl:23
	qw422016.N().S(`, wasAdded bool) {
    old, wasAdded = w.`)
//line generator/components.qtpl:24
	qw422016.E().S(ss)
//line generator/components.qtpl:24
	qw422016.N().S(`.Upsert(e, c);

    // depending on the generation flags, these might be unused
    _, _ = old, wasAdded

`)
//line generator/components.qtpl:29
	if data.ShouldGenAdded {
//line generator/components.qtpl:29
		qw422016.N().S(`    if wasAdded {
        fireEvent(w, `)
//line generator/components.qtpl:31
		qw422016.E().S(nsp)
//line generator/components.qtpl:31
		qw422016.N().S(`AddedEvent{Entity: e, Component: c})
    }
`)
//line generator/components.qtpl:33
	}
//line generator/components.qtpl:34
	if data.ShouldGenChanged {
//line generator/components.qtpl:34
		qw422016.N().S(`    fireEvent(w, `)
//line generator/components.qtpl:35
		qw422016.E().S(nsp)
//line generator/components.qtpl:35
		qw422016.N().S(`ChangedEvent{Entity: e, Old: old, New: c})
`)
//line generator/components.qtpl:36
	}
//line generator/components.qtpl:36
	qw422016.N().S(`
    return old, wasAdded
}

func (w *World) Set`)
//line generator/components.qtpl:41
	qw422016.E().S(nsp)
//line generator/components.qtpl:41
	qw422016.N().S(`FromValues(
    e Entity,
`)
//line generator/components.qtpl:43
	for _, f := range data.Fields {
//line generator/components.qtpl:43
		qw422016.N().S(`    `)
//line generator/components.qtpl:44
		qw422016.E().S(f.Name.Singular.Camel)
//line generator/components.qtpl:44
		qw422016.N().S(`Arg `)
//line generator/components.qtpl:44
		qw422016.E().S(f.Type.Singular.Original)
//line generator/components.qtpl:44
		qw422016.N().S(`,
`)
//line generator/components.qtpl:45
	}
//line generator/components.qtpl:45
	qw422016.N().S(`) {
    old, _ := w.Set`)
//line generator/components.qtpl:47
	qw422016.E().S(nsp)
//line generator/components.qtpl:47
	qw422016.N().S(`(e, `)
//line generator/components.qtpl:47
	qw422016.E().S(nsp)
//line generator/components.qtpl:47
	qw422016.N().S(`{
`)
//line generator/components.qtpl:48
	for _, f := range data.Fields {
//line generator/components.qtpl:48
		qw422016.N().S(`        `)
//line generator/components.qtpl:49
		qw422016.E().S(f.Name.Singular.Pascal)
//line generator/components.qtpl:49
		qw422016.N().S(`: `)
//line generator/components.qtpl:49
		qw422016.E().S(f.Name.Singular.Camel)
//line generator/components.qtpl:49
		qw422016.N().S(`Arg,
`)
//line generator/components.qtpl:50
	}
//line generator/components.qtpl:50
	qw422016.N().S(`    })

    // depending on the generation flags, these might be unused
    _ = old

`)
//line generator/components.qtpl:56
	if data.ShouldGenChanged {
//line generator/components.qtpl:56
		qw422016.N().S(`    fireEvent(w, `)
//line generator/components.qtpl:57
		qw422016.E().S(nsp)
//line generator/components.qtpl:57
		qw422016.N().S(`ChangedEvent{Entity: e, Old: old, New: w.Must`)
//line generator/components.qtpl:57
		qw422016.E().S(nsp)
//line generator/components.qtpl:57
		qw422016.N().S(`(e)})
`)
//line generator/components.qtpl:58
	}
//line generator/components.qtpl:58
	qw422016.N().S(`}

func (w *World) `)
//line generator/components.qtpl:61
	qw422016.E().S(nsp)
//line generator/components.qtpl:61
	qw422016.N().S(`(e Entity) (c `)
//line generator/components.qtpl:61
	qw422016.E().S(nsp)
//line generator/components.qtpl:61
	qw422016.N().S(`, ok bool) {
    return w.`)
//line generator/components.qtpl:62
	qw422016.E().S(ss)
//line generator/components.qtpl:62
	qw422016.N().S(`.Data(e)
}

func (w *World) Mutable`)
//line generator/components.qtpl:65
	qw422016.E().S(nsp)
//line generator/components.qtpl:65
	qw422016.N().S(`(e Entity) (c *`)
//line generator/components.qtpl:65
	qw422016.E().S(nsp)
//line generator/components.qtpl:65
	qw422016.N().S(`, ok bool) {
    return w.`)
//line generator/components.qtpl:66
	qw422016.E().S(ss)
//line generator/components.qtpl:66
	qw422016.N().S(`.DataMutable(e)
}

func (w *World) Must`)
//line generator/components.qtpl:69
	qw422016.E().S(nsp)
//line generator/components.qtpl:69
	qw422016.N().S(`(e Entity) `)
//line generator/components.qtpl:69
	qw422016.E().S(nsp)
//line generator/components.qtpl:69
	qw422016.N().S(` {
    c, ok := w.`)
//line generator/components.qtpl:70
	qw422016.E().S(ss)
//line generator/components.qtpl:70
	qw422016.N().S(`.Data(e)
    if !ok {
        panic("entity does not have `)
//line generator/components.qtpl:72
	qw422016.E().S(nsp)
//line generator/components.qtpl:72
	qw422016.N().S(`")
    }
    return c
}

func (w *World) Remove`)
//line generator/components.qtpl:77
	qw422016.E().S(nsp)
//line generator/components.qtpl:77
	qw422016.N().S(`(e Entity) {
    wasRemoved := w.`)
//line generator/components.qtpl:78
	qw422016.E().S(ss)
//line generator/components.qtpl:78
	qw422016.N().S(`.Remove(e)

    // depending on the generation flags, these might be unused
    _ = wasRemoved

`)
//line generator/components.qtpl:83
	if data.ShouldGenRemoved {
//line generator/components.qtpl:83
		qw422016.N().S(`    if wasRemoved {
        fireEvent(w, `)
//line generator/components.qtpl:85
		qw422016.E().S(nsp)
//line generator/components.qtpl:85
		qw422016.N().S(`RemovedEvent{Entity: e})
    }
`)
//line generator/components.qtpl:87
	}
//line generator/components.qtpl:87
	qw422016.N().S(`}

func (w *World) Has`)
//line generator/components.qtpl:90
	qw422016.E().S(nsp)
//line generator/components.qtpl:90
	qw422016.N().S(`(e Entity) bool {
    return w.`)
//line generator/components.qtpl:91
	qw422016.E().S(ss)
//line generator/components.qtpl:91
	qw422016.N().S(`.Contains(e)
}

func (w *World) All`)
//line generator/components.qtpl:94
	qw422016.E().S(npp)
//line generator/components.qtpl:94
	qw422016.N().S(`(yield func(e Entity, c `)
//line generator/components.qtpl:94
	qw422016.E().S(nsp)
//line generator/components.qtpl:94
	qw422016.N().S(`) bool) {
    for e, c := range w.`)
//line generator/components.qtpl:95
	qw422016.E().S(ss)
//line generator/components.qtpl:95
	qw422016.N().S(`.All {
        if yield(e, c) {
            break
        }
    }
}

func (w *World) AllMutable`)
//line generator/components.qtpl:102
	qw422016.E().S(npp)
//line generator/components.qtpl:102
	qw422016.N().S(`(yield func(e Entity, c *`)
//line generator/components.qtpl:102
	qw422016.E().S(nsp)
//line generator/components.qtpl:102
	qw422016.N().S(`) bool) {
    for e, c := range w.`)
//line generator/components.qtpl:103
	qw422016.E().S(ss)
//line generator/components.qtpl:103
	qw422016.N().S(`.AllMutable {
        if yield(e, c) {
            break
        }
    }
}

func (w *World) All`)
//line generator/components.qtpl:110
	qw422016.E().S(npp)
//line generator/components.qtpl:110
	qw422016.N().S(`Entities(yield func(e Entity) bool) {
    for e := range w.`)
//line generator/components.qtpl:111
	qw422016.E().S(ss)
//line generator/components.qtpl:111
	qw422016.N().S(`.AllEntities {
        if yield(e) {
            break
        }
    }
}

// `)
//line generator/components.qtpl:118
	qw422016.E().S(nsp)
//line generator/components.qtpl:118
	qw422016.N().S(`Builder
func With`)
//line generator/components.qtpl:119
	qw422016.E().S(nsp)
//line generator/components.qtpl:119
	qw422016.N().S(`(c `)
//line generator/components.qtpl:119
	qw422016.E().S(nsp)
//line generator/components.qtpl:119
	qw422016.N().S(`) EntityBuilderOption {
    return func(w *World, e Entity) {
        w.`)
//line generator/components.qtpl:121
	qw422016.E().S(ss)
//line generator/components.qtpl:121
	qw422016.N().S(`.Upsert(e, c)
    }
}

func With`)
//line generator/components.qtpl:125
	qw422016.E().S(nsp)
//line generator/components.qtpl:125
	qw422016.N().S(`FromValues(
`)
//line generator/components.qtpl:126
	for _, f := range data.Fields {
//line generator/components.qtpl:126
		qw422016.N().S(`    `)
//line generator/components.qtpl:127
		qw422016.E().S(f.Name.Singular.Camel)
//line generator/components.qtpl:127
		qw422016.N().S(`Arg `)
//line generator/components.qtpl:127
		qw422016.E().S(f.Type.Singular.Original)
//line generator/components.qtpl:127
		qw422016.N().S(`,
`)
//line generator/components.qtpl:128
	}
//line generator/components.qtpl:128
	qw422016.N().S(`) EntityBuilderOption {
    return func(w *World, e Entity) {
        w.Set`)
//line generator/components.qtpl:131
	qw422016.E().S(nsp)
//line generator/components.qtpl:131
	qw422016.N().S(`FromValues(e,
`)
//line generator/components.qtpl:132
	for _, f := range data.Fields {
//line generator/components.qtpl:132
		qw422016.N().S(`            `)
//line generator/components.qtpl:133
		qw422016.E().S(f.Name.Singular.Camel)
//line generator/components.qtpl:133
		qw422016.N().S(`Arg,
`)
//line generator/components.qtpl:134
	}
//line generator/components.qtpl:134
	qw422016.N().S(`        )
    }
}

// Events
`)
//line generator/components.qtpl:140
	if data.ShouldGenAdded {
//line generator/components.qtpl:140
		qw422016.N().S(`type `)
//line generator/components.qtpl:141
		qw422016.E().S(nsp)
//line generator/components.qtpl:141
		qw422016.N().S(`AddedEvent struct {
    Entity Entity
    Component `)
//line generator/components.qtpl:143
		qw422016.E().S(nsp)
//line generator/components.qtpl:143
		qw422016.N().S(`
}
func (w *World) On`)
//line generator/components.qtpl:145
		qw422016.E().S(nsp)
//line generator/components.qtpl:145
		qw422016.N().S(`Added(fn func(evt `)
//line generator/components.qtpl:145
		qw422016.E().S(nsp)
//line generator/components.qtpl:145
		qw422016.N().S(`AddedEvent)) UnsubscribeFunc {
    unsub := mint.On(w.eventBus, fn)
    return func() {
        unsub()
    }
}
`)
//line generator/components.qtpl:151
	}
//line generator/components.qtpl:151
	qw422016.N().S(`
`)
//line generator/components.qtpl:153
	if data.ShouldGenRemoved {
//line generator/components.qtpl:153
		qw422016.N().S(`type `)
//line generator/components.qtpl:154
		qw422016.E().S(nsp)
//line generator/components.qtpl:154
		qw422016.N().S(`RemovedEvent struct {
    Entity Entity
    Component `)
//line generator/components.qtpl:156
		qw422016.E().S(nsp)
//line generator/components.qtpl:156
		qw422016.N().S(`
}
func (w *World) On`)
//line generator/components.qtpl:158
		qw422016.E().S(nsp)
//line generator/components.qtpl:158
		qw422016.N().S(`Removed(fn func(evt `)
//line generator/components.qtpl:158
		qw422016.E().S(nsp)
//line generator/components.qtpl:158
		qw422016.N().S(`RemovedEvent)) UnsubscribeFunc {
    unsub := mint.On(w.eventBus, fn)
    return func() {
        unsub()
    }
}
`)
//line generator/components.qtpl:164
	}
//line generator/components.qtpl:164
	qw422016.N().S(`
`)
//line generator/components.qtpl:166
	if data.ShouldGenChanged {
//line generator/components.qtpl:166
		qw422016.N().S(`type `)
//line generator/components.qtpl:167
		qw422016.E().S(nsp)
//line generator/components.qtpl:167
		qw422016.N().S(`ChangedEvent struct {
    Entity Entity
    Old, New `)
//line generator/components.qtpl:169
		qw422016.E().S(nsp)
//line generator/components.qtpl:169
		qw422016.N().S(`
}
func (w *World) On`)
//line generator/components.qtpl:171
		qw422016.E().S(nsp)
//line generator/components.qtpl:171
		qw422016.N().S(`Changed(fn func(evt `)
//line generator/components.qtpl:171
		qw422016.E().S(nsp)
//line generator/components.qtpl:171
		qw422016.N().S(`ChangedEvent)) UnsubscribeFunc {
	unsub := mint.On(w.eventBus, fn)
	return func() {
		unsub()
	}
}
`)
//line generator/components.qtpl:177
	}
//line generator/components.qtpl:177
	qw422016.N().S(`
// Resource methods
func (w *World) Set`)
//line generator/components.qtpl:180
	qw422016.E().S(nsp)
//line generator/components.qtpl:180
	qw422016.N().S(`Resource(c `)
//line generator/components.qtpl:180
	qw422016.E().S(nsp)
//line generator/components.qtpl:180
	qw422016.N().S(`) {
    w.`)
//line generator/components.qtpl:181
	qw422016.E().S(ss)
//line generator/components.qtpl:181
	qw422016.N().S(`.Upsert(w.resourceEntity, c)
}

func (w *World) Set`)
//line generator/components.qtpl:184
	qw422016.E().S(nsp)
//line generator/components.qtpl:184
	qw422016.N().S(`ResourceFromValues(
`)
//line generator/components.qtpl:185
	for _, f := range data.Fields {
//line generator/components.qtpl:185
		qw422016.N().S(`    `)
//line generator/components.qtpl:186
		qw422016.E().S(f.Name.Singular.Camel)
//line generator/components.qtpl:186
		qw422016.N().S(`Arg `)
//line generator/components.qtpl:186
		qw422016.E().S(f.Type.Singular.Original)
//line generator/components.qtpl:186
		qw422016.N().S(`,
`)
//line generator/components.qtpl:187
	}
//line generator/components.qtpl:187
	qw422016.N().S(`) {
   w.Set`)
//line generator/components.qtpl:189
	qw422016.E().S(nsp)
//line generator/components.qtpl:189
	qw422016.N().S(`Resource(`)
//line generator/components.qtpl:189
	qw422016.E().S(nsp)
//line generator/components.qtpl:189
	qw422016.N().S(`{
`)
//line generator/components.qtpl:190
	for _, f := range data.Fields {
//line generator/components.qtpl:190
		qw422016.N().S(`        `)
//line generator/components.qtpl:191
		qw422016.E().S(f.Name.Singular.Pascal)
//line generator/components.qtpl:191
		qw422016.N().S(`: `)
//line generator/components.qtpl:191
		qw422016.E().S(f.Name.Singular.Camel)
//line generator/components.qtpl:191
		qw422016.N().S(`Arg,
`)
//line generator/components.qtpl:192
	}
//line generator/components.qtpl:192
	qw422016.N().S(`    })
}

func (w *World) `)
//line generator/components.qtpl:196
	qw422016.E().S(nsp)
//line generator/components.qtpl:196
	qw422016.N().S(`Resource() (`)
//line generator/components.qtpl:196
	qw422016.E().S(nsp)
//line generator/components.qtpl:196
	qw422016.N().S(`,bool) {
    return w.`)
//line generator/components.qtpl:197
	qw422016.E().S(ss)
//line generator/components.qtpl:197
	qw422016.N().S(`.Data(w.resourceEntity)
}

func (w *World) Must`)
//line generator/components.qtpl:200
	qw422016.E().S(nsp)
//line generator/components.qtpl:200
	qw422016.N().S(`Resource() `)
//line generator/components.qtpl:200
	qw422016.E().S(nsp)
//line generator/components.qtpl:200
	qw422016.N().S(` {
    c, ok := w.`)
//line generator/components.qtpl:201
	qw422016.E().S(nsp)
//line generator/components.qtpl:201
	qw422016.N().S(`Resource()
    if !ok {
        panic("resource entity does not have `)
//line generator/components.qtpl:203
	qw422016.E().S(nsp)
//line generator/components.qtpl:203
	qw422016.N().S(`")
    }
    return c
}

func (w *World) Remove`)
//line generator/components.qtpl:208
	qw422016.E().S(nsp)
//line generator/components.qtpl:208
	qw422016.N().S(`Resource() {
    w.`)
//line generator/components.qtpl:209
	qw422016.E().S(ss)
//line generator/components.qtpl:209
	qw422016.N().S(`.Remove(w.resourceEntity)
}

`)
//line generator/components.qtpl:212
}

//line generator/components.qtpl:212
func writecomponentTemplate(qq422016 qtio422016.Writer, data *componentTmplData) {
//line generator/components.qtpl:212
	qw422016 := qt422016.AcquireWriter(qq422016)
//line generator/components.qtpl:212
	streamcomponentTemplate(qw422016, data)
//line generator/components.qtpl:212
	qt422016.ReleaseWriter(qw422016)
//line generator/components.qtpl:212
}

//line generator/components.qtpl:212
func componentTemplate(data *componentTmplData) string {
//line generator/components.qtpl:212
	qb422016 := qt422016.AcquireByteBuffer()
//line generator/components.qtpl:212
	writecomponentTemplate(qb422016, data)
//line generator/components.qtpl:212
	qs422016 := string(qb422016.B)
//line generator/components.qtpl:212
	qt422016.ReleaseByteBuffer(qb422016)
//line generator/components.qtpl:212
	return qs422016
//line generator/components.qtpl:212
}
