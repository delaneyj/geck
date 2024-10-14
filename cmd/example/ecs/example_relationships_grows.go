package ecs

import (
	"github.com/tidwall/btree"
)

type GrowsRelationshipPair struct {
	Source Entity
	Target Entity
}

type GrowsRelationship struct {
	btree *btree.BTreeG[GrowsRelationshipPair]
}

func NewGrowsRelationship() *GrowsRelationship {
	return &GrowsRelationship{
		btree: btree.NewBTreeG[GrowsRelationshipPair](func(a, b GrowsRelationshipPair) bool {
			ati, bti := a.Target.Index(), b.Target.Index()
			if ati == bti {
				return a.Source.Index() < b.Source.Index()
			}
			return ati < bti
		}),
	}
}

func (r *GrowsRelationship) Clear() {
	r.btree.Clear()
}

func (w *World) LinkGrows(target Entity, sources ...Entity) {
	for _, source := range sources {
		pair := GrowsRelationshipPair{
			Target: target,
			Source: source,
		}

		w.growsRelationships.btree.Set(pair)
	}
}

func (w *World) UnlinkGrows(target Entity, sources ...Entity) {
	for _, source := range sources {
		pair := GrowsRelationshipPair{
			Target: target,
			Source: source,
		}

		w.growsRelationships.btree.Delete(pair)
	}
}

func (w *World) GrowsIsLinked(source, target Entity) bool {
	pair := GrowsRelationshipPair{
		Source: source,
		Target: target,
	}

	_, ok := w.growsRelationships.btree.Get(pair)
	return ok
}

func (w *World) GrowsSources(target Entity) func(yield func(source Entity) bool) {
	return func(yield func(source Entity) bool) {
		iter := w.growsRelationships.btree.Iter()
		iter.Seek(GrowsRelationshipPair{Target: target})
		end := GrowsRelationshipPair{Target: target + 1}

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

func (w *World) RemoveGrowsRelationships(target Entity, sources ...Entity) {
	for _, source := range sources {
		pair := GrowsRelationshipPair{
			Target: target,
			Source: source,
		}

		w.growsRelationships.btree.Delete(pair)
	}
}

func (w *World) RemoveAllGrowsRelationships(target Entity) {
	iter := w.growsRelationships.btree.Iter()
	iter.Seek(GrowsRelationshipPair{Target: target})
	end := GrowsRelationshipPair{Target: target + 1}

	for iter.Next() {
		item := iter.Item()
		if item.Target >= end.Target {
			break
		}

		w.growsRelationships.btree.Delete(item)
	}
}
