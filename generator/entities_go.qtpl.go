// Code generated by qtc from "entities_go.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

// package generator
//

//line generator/entities_go.qtpl:3
package generator

//line generator/entities_go.qtpl:3
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line generator/entities_go.qtpl:3
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line generator/entities_go.qtpl:3
func streamentitiesTemplate(qw422016 *qt422016.Writer, data *ecsTmplData) {
//line generator/entities_go.qtpl:3
	qw422016.N().S(`
package `)
//line generator/entities_go.qtpl:4
	qw422016.E().S(data.PackageName)
//line generator/entities_go.qtpl:4
	qw422016.N().S(`

const (
	indexBits      = 12
	generationBits = 20
	indexMask      = (1 << indexBits) - 1
	generationMask = (1 << generationBits) - 1
	maxEntities    = 1 << indexBits
)

var Tombstone = Entity(maxEntities)

type Entity uint32

func NewEntity(index, generation int) Entity {
	return Entity((generation & generationMask) | ((index & indexMask) << generationBits))
}

func (e Entity) Index() int {
	return int(e>>generationBits) & indexMask
}


func (e Entity) Generation() int {
	return int(e) & generationMask
}

func (e Entity) In(entities ...Entity) bool {
	for _, entity := range entities {
		if e == entity {
			return true
		}
	}
	return false
}

func SortEntities(fn func(yield func(entity Entity) bool)) []Entity {
	entities := make([]Entity, 0, 4096)
	for e := range fn {
		entities = append(entities, e)
	}

	slices.Sort(entities)

	return entities
}

type EntityBuilderOption func(w *World, entity Entity)

func(w *World) CreateEntities(count int, opts ...EntityBuilderOption) []Entity{
    entities := make([]Entity, count)
    for i := range entities {
        var entity Entity
        
        if w.freeEntities.IsEmpty() {
            entity = Entity(w.nextEntityID)
            w.nextEntityID++
        } else {
            entity = Entity(w.freeEntities.Minimum())
            w.freeEntities.Remove(uint32(entity))
        }
        w.livingEntities.Add(uint32(entity))
        entities[i] = entity

		for _, opt := range opts {
			opt(w, entity)
		}
    }
    return entities
}

func (w *World) CreateEntity(opts ...EntityBuilderOption) Entity {
	return w.CreateEntities(1, opts...)[0]
}

func (w *World) DestroyEntities(entities ...Entity) {
	for _, entity := range entities {
		w.livingEntities.Remove(uint32(entity))
		w.freeEntities.Add(uint32(entity))
	}
}

func (w *World) IsAlive(entity Entity) bool {
	return w.livingEntities.Contains(uint32(entity))
}

func (w *World) All(yield func(entity Entity) bool) {
	for _, entity := range w.livingEntities.ToArray() {
		if !yield(Entity(entity)) {
			break
		}
	}
}

`)
//line generator/entities_go.qtpl:98
}

//line generator/entities_go.qtpl:98
func writeentitiesTemplate(qq422016 qtio422016.Writer, data *ecsTmplData) {
//line generator/entities_go.qtpl:98
	qw422016 := qt422016.AcquireWriter(qq422016)
//line generator/entities_go.qtpl:98
	streamentitiesTemplate(qw422016, data)
//line generator/entities_go.qtpl:98
	qt422016.ReleaseWriter(qw422016)
//line generator/entities_go.qtpl:98
}

//line generator/entities_go.qtpl:98
func entitiesTemplate(data *ecsTmplData) string {
//line generator/entities_go.qtpl:98
	qb422016 := qt422016.AcquireByteBuffer()
//line generator/entities_go.qtpl:98
	writeentitiesTemplate(qb422016, data)
//line generator/entities_go.qtpl:98
	qs422016 := string(qb422016.B)
//line generator/entities_go.qtpl:98
	qt422016.ReleaseByteBuffer(qb422016)
//line generator/entities_go.qtpl:98
	return qs422016
//line generator/entities_go.qtpl:98
}
