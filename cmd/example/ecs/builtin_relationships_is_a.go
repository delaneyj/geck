package ecs

import (
	"github.com/tidwall/btree"
)

type IsARelationshipPair struct {
	From Entity
	To   Entity
}

type IsARelationship struct {
	btree *btree.BTreeG[IsARelationshipPair]
}

func NewIsARelationship() *IsARelationship {
	return &IsARelationship{
		btree: btree.NewBTreeG[IsARelationshipPair](func(a, b IsARelationshipPair) bool {
			ati, bti := a.To.Index(), b.To.Index()
			if ati == bti {
				return a.From.Index() < b.From.Index()
			}
			return ati < bti
		}),
	}
}

func (r *IsARelationship) Clear() {
	r.btree.Clear()
}

func (w *World) LinkIsA(
	to, from Entity,
) {
	pair := IsARelationshipPair{
		From: from,
		To:   to,
	}

	w.isARelationships.btree.Set(pair)
}

func (w *World) UnlinkIsA(from, to Entity) {
	pair := IsARelationshipPair{
		From: from,
		To:   to,
	}
	w.isARelationships.btree.Delete(pair)
}

func (w *World) IsAIsLinked(from, to Entity) bool {
	pair := IsARelationshipPair{
		From: from,
		To:   to,
	}

	_, ok := w.isARelationships.btree.Get(pair)
	return ok
}

func (w *World) IsA(to Entity) func(yield func(from Entity) bool) {
	return func(yield func(from Entity) bool) {
		iter := w.isARelationships.btree.Iter()
		iter.Seek(IsARelationshipPair{To: to})
		end := IsARelationshipPair{To: to + 1}

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

func (w *World) RemoveIsARelationships(to Entity, froms ...Entity) {
	for _, from := range froms {
		pair := IsARelationshipPair{
			To:   to,
			From: from,
		}

		w.isARelationships.btree.Delete(pair)
	}
}

func (w *World) RemoveAllIsARelationships(to Entity) {
	iter := w.isARelationships.btree.Iter()
	end := IsARelationshipPair{To: to + 1}

	for iter.Next() {
		item := iter.Item()
		if item.To >= end.To {
			break
		}

		w.isARelationships.btree.Delete(item)
	}
}
