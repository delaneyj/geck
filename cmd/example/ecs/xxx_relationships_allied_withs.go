package ecs

import (
	"github.com/tidwall/btree"
)

type AlliedWithRelationshipPair struct {
	From, To Entity
}

type AlliedWithRelationship struct {
	btree *btree.BTreeG[AlliedWithRelationshipPair]
}

func NewAlliedWithRelationship() *AlliedWithRelationship {
	return &AlliedWithRelationship{
		btree: btree.NewBTreeG(func(a, b AlliedWithRelationshipPair) bool {
			ati, bti := a.To.Index(), b.To.Index()
			if ati == bti {
				return a.From.Index() < b.From.Index()
			}
			return ati < bti
		}),
	}
}

func (r *AlliedWithRelationship) Clear() {
	r.btree.Clear()
}

func (w *World) LinkAlliedWith(
	to, from Entity,
) {
	pair := AlliedWithRelationshipPair{
		From: from, To: to,
	}
	w.alliedWithRelationships.btree.Set(pair)
}

func (w *World) UnlinkAlliedWith(from, to Entity) {
	pair := AlliedWithRelationshipPair{From: from, To: to}
	w.alliedWithRelationships.btree.Delete(pair)
}

func (w *World) AlliedWithIsLinked(from, to Entity) bool {
	pair := AlliedWithRelationshipPair{From: from, To: to}
	_, ok := w.alliedWithRelationships.btree.Get(pair)
	return ok
}

func (w *World) AlliedWith(to Entity) func(yield func(from Entity) bool) {
	return func(yield func(from Entity) bool) {
		iter := w.alliedWithRelationships.btree.Iter()
		iter.Seek(AlliedWithRelationshipPair{To: to})
		end := AlliedWithRelationshipPair{To: to + 1}
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

func (w *World) RemoveAlliedWithRelationships(to Entity, froms ...Entity) {
	for _, from := range froms {
		pair := AlliedWithRelationshipPair{From: from, To: to}
		w.alliedWithRelationships.btree.Delete(pair)
	}
}

func (w *World) RemoveAllAlliedWithRelationships(to Entity) {
	iter := w.alliedWithRelationships.btree.Iter()
	end := AlliedWithRelationshipPair{To: to + 1}
	for iter.Next() {
		item := iter.Item()
		if item.To >= end.To {
			break
		}
		w.alliedWithRelationships.btree.Delete(item)
	}
}
