package ecs

import (
	"github.com/tidwall/btree"
)

type ChildOfRelationshipPair struct {
	From, To Entity
}

type ChildOfRelationship struct {
	btree *btree.BTreeG[ChildOfRelationshipPair]
}

func NewChildOfRelationship() *ChildOfRelationship {
	return &ChildOfRelationship{
		btree: btree.NewBTreeG(func(a, b ChildOfRelationshipPair) bool {
			ati, bti := a.To.Index(), b.To.Index()
			if ati == bti {
				return a.From.Index() < b.From.Index()
			}
			return ati < bti
		}),
	}
}

func (r *ChildOfRelationship) Clear() {
	r.btree.Clear()
}

func (w *World) LinkChildOf(
	to, from Entity,
) {
	pair := ChildOfRelationshipPair{
		From: from, To: to,
	}
	w.childOfRelationships.btree.Set(pair)
}

func (w *World) UnlinkChildOf(from, to Entity) {
	pair := ChildOfRelationshipPair{From: from, To: to}
	w.childOfRelationships.btree.Delete(pair)
}

func (w *World) ChildOfIsLinked(from, to Entity) bool {
	pair := ChildOfRelationshipPair{From: from, To: to}
	_, ok := w.childOfRelationships.btree.Get(pair)
	return ok
}

func (w *World) ChildOf(to Entity) func(yield func(from Entity) bool) {
	return func(yield func(from Entity) bool) {
		iter := w.childOfRelationships.btree.Iter()
		iter.Seek(ChildOfRelationshipPair{To: to})
		end := ChildOfRelationshipPair{To: to + 1}
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

func (w *World) RemoveChildOfRelationships(to Entity, froms ...Entity) {
	for _, from := range froms {
		pair := ChildOfRelationshipPair{From: from, To: to}
		w.childOfRelationships.btree.Delete(pair)
	}
}

func (w *World) RemoveAllChildOfRelationships(to Entity) {
	iter := w.childOfRelationships.btree.Iter()
	end := ChildOfRelationshipPair{To: to + 1}
	for iter.Next() {
		item := iter.Item()
		if item.To >= end.To {
			break
		}
		w.childOfRelationships.btree.Delete(item)
	}
}
