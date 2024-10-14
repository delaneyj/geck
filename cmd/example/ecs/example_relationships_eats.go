package ecs

import (
	"github.com/tidwall/btree"
)

type EatsRelationshipPair struct {
	From   Entity
	To     Entity
	Amount uint8
}

type EatsRelationship struct {
	btree *btree.BTreeG[EatsRelationshipPair]
}

func NewEatsRelationship() *EatsRelationship {
	return &EatsRelationship{
		btree: btree.NewBTreeG[EatsRelationshipPair](func(a, b EatsRelationshipPair) bool {
			ati, bti := a.To.Index(), b.To.Index()
			if ati == bti {
				return a.From.Index() < b.From.Index()
			}
			return ati < bti
		}),
	}
}

func (r *EatsRelationship) Clear() {
	r.btree.Clear()
}

func (w *World) LinkEats(
	to, from Entity,
	amountArg uint8,
) {
	pair := EatsRelationshipPair{
		From:   from,
		To:     to,
		Amount: amountArg,
	}

	w.eatsRelationships.btree.Set(pair)
}

func (w *World) UnlinkEats(from, to Entity) {
	pair := EatsRelationshipPair{
		From: from,
		To:   to,
	}
	w.eatsRelationships.btree.Delete(pair)
}

func (w *World) EatsIsLinked(from, to Entity) bool {
	pair := EatsRelationshipPair{
		From: from,
		To:   to,
	}

	_, ok := w.eatsRelationships.btree.Get(pair)
	return ok
}

func (w *World) Eats(to Entity) func(yield func(from Entity) bool) {
	return func(yield func(from Entity) bool) {
		iter := w.eatsRelationships.btree.Iter()
		iter.Seek(EatsRelationshipPair{To: to})
		end := EatsRelationshipPair{To: to + 1}

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

func (w *World) RemoveEatsRelationships(to Entity, froms ...Entity) {
	for _, from := range froms {
		pair := EatsRelationshipPair{
			To:   to,
			From: from,
		}

		w.eatsRelationships.btree.Delete(pair)
	}
}

func (w *World) RemoveAllEatsRelationships(to Entity) {
	iter := w.eatsRelationships.btree.Iter()
	end := EatsRelationshipPair{To: to + 1}

	for iter.Next() {
		item := iter.Item()
		if item.To >= end.To {
			break
		}

		w.eatsRelationships.btree.Delete(item)
	}
}
