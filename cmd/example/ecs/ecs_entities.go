package ecs

import "slices"

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
	return Entity((generation << generationBits) | index)
}

func (e Entity) Index() int {
	return int(e) & indexMask
}

func (e Entity) Generation() int {
	return int(e) >> indexBits
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

func (w *World) CreateEntities(count int, opts ...EntityBuilderOption) []Entity {
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
