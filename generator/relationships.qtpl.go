// Code generated by qtc from "relationships.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

// package generator
//

//line generator/relationships.qtpl:3
package generator

//line generator/relationships.qtpl:3
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line generator/relationships.qtpl:3
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line generator/relationships.qtpl:3
func streamrelationshipTemplate(qw422016 *qt422016.Writer, data *componentTmplData) {
//line generator/relationships.qtpl:3
	qw422016.N().S(`
package `)
//line generator/relationships.qtpl:4
	qw422016.E().S(data.PackageName)
//line generator/relationships.qtpl:4
	qw422016.N().S(`
import (
    "github.com/tidwall/btree"
`)
//line generator/relationships.qtpl:8
	if data.HasAnyEvents {
//line generator/relationships.qtpl:8
		qw422016.N().S(`
    import "github.com/btvoidx/mint"
`)
//line generator/relationships.qtpl:10
	}
//line generator/relationships.qtpl:10
	qw422016.N().S(`
)

`)
//line generator/relationships.qtpl:14
	nsp := data.Name.Singular.Pascal
	nsc := data.Name.Singular.Camel
	pairName := data.Name.Singular.Pascal + "RelationshipPair"

//line generator/relationships.qtpl:17
	qw422016.N().S(`
type `)
//line generator/relationships.qtpl:19
	qw422016.E().S(pairName)
//line generator/relationships.qtpl:19
	qw422016.N().S(` struct {
    Source Entity
    Target Entity
}

type `)
//line generator/relationships.qtpl:24
	qw422016.E().S(nsp)
//line generator/relationships.qtpl:24
	qw422016.N().S(`Relationship struct {
    btree *btree.BTreeG[`)
//line generator/relationships.qtpl:25
	qw422016.E().S(pairName)
//line generator/relationships.qtpl:25
	qw422016.N().S(`]
}

func New`)
//line generator/relationships.qtpl:28
	qw422016.E().S(nsp)
//line generator/relationships.qtpl:28
	qw422016.N().S(`Relationship() *`)
//line generator/relationships.qtpl:28
	qw422016.E().S(nsp)
//line generator/relationships.qtpl:28
	qw422016.N().S(`Relationship {
    return &`)
//line generator/relationships.qtpl:29
	qw422016.E().S(nsp)
//line generator/relationships.qtpl:29
	qw422016.N().S(`Relationship{
        btree: btree.NewBTreeG[`)
//line generator/relationships.qtpl:30
	qw422016.E().S(pairName)
//line generator/relationships.qtpl:30
	qw422016.N().S(`](func(a, b `)
//line generator/relationships.qtpl:30
	qw422016.E().S(pairName)
//line generator/relationships.qtpl:30
	qw422016.N().S(`) bool {
            ati, bti := a.Target.Index(), b.Target.Index()
            if ati == bti {
                return a.Source.Index() < b.Source.Index()
            }
            return ati < bti
        }),
    }
}

func (r *`)
//line generator/relationships.qtpl:40
	qw422016.E().S(nsp)
//line generator/relationships.qtpl:40
	qw422016.N().S(`Relationship) Clear() {
    r.btree.Clear()
}

func(w *World) Link`)
//line generator/relationships.qtpl:44
	qw422016.E().S(nsp)
//line generator/relationships.qtpl:44
	qw422016.N().S(`(target Entity, sources ... Entity) {
    for _, source := range sources {
        pair := `)
//line generator/relationships.qtpl:46
	qw422016.E().S(pairName)
//line generator/relationships.qtpl:46
	qw422016.N().S(`{
            Target: target,
            Source: source,
        }

        w.`)
//line generator/relationships.qtpl:51
	qw422016.E().S(nsc)
//line generator/relationships.qtpl:51
	qw422016.N().S(`Relationships.btree.Set(pair)
    }
}

func(w *World) Unlink`)
//line generator/relationships.qtpl:55
	qw422016.E().S(nsp)
//line generator/relationships.qtpl:55
	qw422016.N().S(`(target Entity, sources ... Entity) {
    for _, source := range sources {
        pair := `)
//line generator/relationships.qtpl:57
	qw422016.E().S(pairName)
//line generator/relationships.qtpl:57
	qw422016.N().S(`{
            Target: target,
            Source: source,
        }

        w.`)
//line generator/relationships.qtpl:62
	qw422016.E().S(nsc)
//line generator/relationships.qtpl:62
	qw422016.N().S(`Relationships.btree.Delete(pair)
    }
}

func (w *World) `)
//line generator/relationships.qtpl:66
	qw422016.E().S(nsp)
//line generator/relationships.qtpl:66
	qw422016.N().S(`IsLinked(source, target Entity) bool {
    pair := `)
//line generator/relationships.qtpl:67
	qw422016.E().S(pairName)
//line generator/relationships.qtpl:67
	qw422016.N().S(`{
        Source: source,
        Target: target,
    }

    _, ok := w.`)
//line generator/relationships.qtpl:72
	qw422016.E().S(nsc)
//line generator/relationships.qtpl:72
	qw422016.N().S(`Relationships.btree.Get(pair)
    return ok
}

func (w *World) `)
//line generator/relationships.qtpl:76
	qw422016.E().S(nsp)
//line generator/relationships.qtpl:76
	qw422016.N().S(`Sources(target Entity) func(yield func(source Entity) bool) {
    return func(yield func(source Entity) bool) {
        iter := w.`)
//line generator/relationships.qtpl:78
	qw422016.E().S(nsc)
//line generator/relationships.qtpl:78
	qw422016.N().S(`Relationships.btree.Iter()
        iter.Seek(`)
//line generator/relationships.qtpl:79
	qw422016.E().S(pairName)
//line generator/relationships.qtpl:79
	qw422016.N().S(`{ Target: target })
        end := `)
//line generator/relationships.qtpl:80
	qw422016.E().S(pairName)
//line generator/relationships.qtpl:80
	qw422016.N().S(`{ Target: target + 1 }

        for iter.Next() {
            item := iter.Item()
            if item.Target >= end.Target {
                break
            }

            if !yield(item.Source) {
                break
            }
        }
    }
}

func (w *World) Remove`)
//line generator/relationships.qtpl:95
	qw422016.E().S(nsp)
//line generator/relationships.qtpl:95
	qw422016.N().S(`Relationships(target Entity, sources ... Entity) {
    for _, source := range sources {
        pair := `)
//line generator/relationships.qtpl:97
	qw422016.E().S(pairName)
//line generator/relationships.qtpl:97
	qw422016.N().S(`{
            Target: target,
            Source: source,
        }

        w.`)
//line generator/relationships.qtpl:102
	qw422016.E().S(nsc)
//line generator/relationships.qtpl:102
	qw422016.N().S(`Relationships.btree.Delete(pair)
    }
}

func (w *World) RemoveAll`)
//line generator/relationships.qtpl:106
	qw422016.E().S(nsp)
//line generator/relationships.qtpl:106
	qw422016.N().S(`Relationships(target Entity) {
    iter := w.`)
//line generator/relationships.qtpl:107
	qw422016.E().S(nsc)
//line generator/relationships.qtpl:107
	qw422016.N().S(`Relationships.btree.Iter()
    iter.Seek(`)
//line generator/relationships.qtpl:108
	qw422016.E().S(pairName)
//line generator/relationships.qtpl:108
	qw422016.N().S(`{ Target: target })
    end := `)
//line generator/relationships.qtpl:109
	qw422016.E().S(pairName)
//line generator/relationships.qtpl:109
	qw422016.N().S(`{ Target: target + 1 }

    for iter.Next() {
        item := iter.Item()
        if item.Target >= end.Target {
            break
        }

        w.`)
//line generator/relationships.qtpl:117
	qw422016.E().S(nsc)
//line generator/relationships.qtpl:117
	qw422016.N().S(`Relationships.btree.Delete(item)
    }
}


`)
//line generator/relationships.qtpl:122
}

//line generator/relationships.qtpl:122
func writerelationshipTemplate(qq422016 qtio422016.Writer, data *componentTmplData) {
//line generator/relationships.qtpl:122
	qw422016 := qt422016.AcquireWriter(qq422016)
//line generator/relationships.qtpl:122
	streamrelationshipTemplate(qw422016, data)
//line generator/relationships.qtpl:122
	qt422016.ReleaseWriter(qw422016)
//line generator/relationships.qtpl:122
}

//line generator/relationships.qtpl:122
func relationshipTemplate(data *componentTmplData) string {
//line generator/relationships.qtpl:122
	qb422016 := qt422016.AcquireByteBuffer()
//line generator/relationships.qtpl:122
	writerelationshipTemplate(qb422016, data)
//line generator/relationships.qtpl:122
	qs422016 := string(qb422016.B)
//line generator/relationships.qtpl:122
	qt422016.ReleaseByteBuffer(qb422016)
//line generator/relationships.qtpl:122
	return qs422016
//line generator/relationships.qtpl:122
}
