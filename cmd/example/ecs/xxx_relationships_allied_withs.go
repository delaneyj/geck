package ecs

import (
	"github.com/tidwall/btree"
)

type AlliedWithRelationshipPair struct {
	Source Entity
	Target Entity
}

type AlliedWithRelationship struct {
	btree *btree.BTreeG[AlliedWithRelationshipPair]
}

func NewAlliedWithRelationship() *AlliedWithRelationship {
	return &AlliedWithRelationship{
		btree: btree.NewBTreeG[AlliedWithRelationshipPair](func(a, b AlliedWithRelationshipPair) bool {
			ati, bti := a.Target.Index(), b.Target.Index()
			if ati == bti {
				return a.Source.Index() < b.Source.Index()
			}
			return ati < bti
		}),
	}
}

func (r *AlliedWithRelationship) Clear() {
	r.btree.Clear()
}

func (w *World) LinkAlliedWith(target Entity, sources ...Entity) {
	for _, source := range sources {
		pair := AlliedWithRelationshipPair{
			Target: target,
			Source: source,
		}

		w.alliedWithRelationships.btree.Set(pair)
	}
}

func (w *World) UnlinkAlliedWith(target Entity, sources ...Entity) {
	for _, source := range sources {
		pair := AlliedWithRelationshipPair{
			Target: target,
			Source: source,
		}

		w.alliedWithRelationships.btree.Delete(pair)
	}
}

func (w *World) AlliedWithIsLinked(source, target Entity) bool {
	pair := AlliedWithRelationshipPair{
		Source: source,
		Target: target,
	}

	_, ok := w.alliedWithRelationships.btree.Get(pair)
	return ok
}

func (w *World) AlliedWithSources(target Entity) func(yield func(source Entity) bool) {
	return func(yield func(source Entity) bool) {
		iter := w.alliedWithRelationships.btree.Iter()
		iter.Seek(AlliedWithRelationshipPair{Target: target})
		end := AlliedWithRelationshipPair{Target: target + 1}

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

func (w *World) RemoveAlliedWithRelationships(target Entity, sources ...Entity) {
	for _, source := range sources {
		pair := AlliedWithRelationshipPair{
			Target: target,
			Source: source,
		}

		w.alliedWithRelationships.btree.Delete(pair)
	}
}

func (w *World) RemoveAllAlliedWithRelationships(target Entity) {
	iter := w.alliedWithRelationships.btree.Iter()
	iter.Seek(AlliedWithRelationshipPair{Target: target})
	end := AlliedWithRelationshipPair{Target: target + 1}

	for iter.Next() {
		item := iter.Item()
		if item.Target >= end.Target {
			break
		}

		w.alliedWithRelationships.btree.Delete(item)
	}
}
