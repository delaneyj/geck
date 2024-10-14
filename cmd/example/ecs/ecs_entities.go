package ecs

import "slices"

const (
	indexBits      = 20
	generationBits = 12
	indexMask      = (1 << indexBits) - 1
	generationMask = (1 << generationBits) - 1
	maxEntities    = 1 << indexBits
)

var Tombstone = Entity(maxEntities)

type Entity uint32

func NewEntity(index, generation int) Entity {
	return Entity((index << generationBits) | generation)
}

func (e Entity) Index() int {
	return int(e) >> generationBits
}

func (e Entity) Generation() int {
	return int(e) & generationMask
}

func EntityFromU32(u uint32) Entity {
	return Entity(u)
}

func (e Entity) InSlice(entities ...Entity) bool {
	for _, entity := range entities {
		if e == entity {
			return true
		}
	}
	return false
}

func (e Entity) InIter(iter func(yield func(entity Entity) bool)) func(yield func(entity Entity) bool) {
	return func(yield func(entity Entity) bool) {
		iter(func(entity Entity) bool {
			if e == entity {
				return yield(entity)
			}
			return true
		})
	}
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

func (w *World) NextEntities(count int, opts ...EntityBuilderOption) []Entity {
	entities := make([]Entity, count)
	for i := range entities {
		var entity Entity

		if w.freeEntities.Len() == 0 {
			entity = Entity(w.nextEntityID)
			w.nextEntityID++
		} else {
			entity = w.freeEntities.dense[0]
			w.freeEntities.Remove(entity)
		}
		w.livingEntities.Upsert(entity, empty{})
		entities[i] = entity

		for _, opt := range opts {
			opt(w, entity)
		}
	}
	return entities
}

func (w *World) NextEntity(opts ...EntityBuilderOption) Entity {
	return w.NextEntities(1, opts...)[0]
}

func (w *World) DestroyEntities(entities ...Entity) {
	for _, entity := range entities {
		w.livingEntities.Remove(entity)
		w.freeEntities.Upsert(entity, empty{})

		w.nameComponents.Remove(entity)
		w.childOfComponents.Remove(entity)
		w.isAComponents.Remove(entity)
		w.positionComponents.Remove(entity)
		w.velocityComponents.Remove(entity)
		w.rotationComponents.Remove(entity)
		w.directionComponents.Remove(entity)
		w.eatsComponents.Remove(entity)
		w.likesRelationships.Clear()
		w.enemyTags.Remove(entity)
		w.growsRelationships.Clear()
		w.gravityComponents.Remove(entity)
		w.spaceshipTags.Remove(entity)
		w.spacestationTags.Remove(entity)
		w.factionComponents.Remove(entity)
		w.dockedToComponents.Remove(entity)
		w.planetTags.Remove(entity)
		w.ruledByComponents.Remove(entity)
		w.alliedWithRelationships.Clear()
	}
}

func (w *World) IsAlive(entity Entity) bool {
	return w.livingEntities.Contains(entity)
}

func (w *World) All(yield func(entity Entity) bool) {
	for e := range w.livingEntities.AllEntities {
		if !yield(e) {
			break
		}
	}
}
