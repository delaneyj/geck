package ecs

import (
	"github.com/tidwall/btree"
)

type GrowsRelationshipPair struct {
	From, To Entity
}

type GrowsRelationship struct {
	btree *btree.BTreeG[GrowsRelationshipPair]
}

func NewGrowsRelationship() *GrowsRelationship {
	return &GrowsRelationship{
		btree: btree.NewBTreeG(func(a, b GrowsRelationshipPair) bool {
			ati, bti := a.To.Index(), b.To.Index()
			if ati == bti {
				return a.From.Index() < b.From.Index()
			}
			return ati < bti
		}),
	}
}

func (r *GrowsRelationship) Clear() {
	r.btree.Clear()
}

func (w *World) LinkGrows(
	to, from Entity,
) {
	pair := GrowsRelationshipPair{
		From: from, To: to,
	}
	w.growsRelationships.btree.Set(pair)
}

func (w *World) UnlinkGrows(from, to Entity) {
	pair := GrowsRelationshipPair{From: from, To: to}
	w.growsRelationships.btree.Delete(pair)
}

func (w *World) GrowsIsLinked(from, to Entity) bool {
	pair := GrowsRelationshipPair{From: from, To: to}
	_, ok := w.growsRelationships.btree.Get(pair)
	return ok
}

func (w *World) Grows(to Entity) func(yield func(from Entity) bool) {
	return func(yield func(from Entity) bool) {
		iter := w.growsRelationships.btree.Iter()
		iter.Seek(GrowsRelationshipPair{To: to})
		end := GrowsRelationshipPair{To: to + 1}
		for iter.Next() {
			item := iter.Item()
			if item.To >= end.To {
				break
			}

			if !yield(item.From) {
				break
			}
		}
	}
}

func (w *World) RemoveGrowsRelationships(to Entity, froms ...Entity) {
	for _, from := range froms {
		pair := GrowsRelationshipPair{From: from, To: to}
		w.growsRelationships.btree.Delete(pair)
	}
}

func (w *World) RemoveAllGrowsRelationships(to Entity) {
	iter := w.growsRelationships.btree.Iter()
	end := GrowsRelationshipPair{To: to + 1}
	for iter.Next() {
		item := iter.Item()
		if item.To >= end.To {
			break
		}
		w.growsRelationships.btree.Delete(item)
	}
}
