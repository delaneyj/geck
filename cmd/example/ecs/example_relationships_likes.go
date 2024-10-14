package ecs

import (
	"github.com/tidwall/btree"
)

type LikesRelationshipPair struct {
	From Entity
	To   Entity
}

type LikesRelationship struct {
	btree *btree.BTreeG[LikesRelationshipPair]
}

func NewLikesRelationship() *LikesRelationship {
	return &LikesRelationship{
		btree: btree.NewBTreeG[LikesRelationshipPair](func(a, b LikesRelationshipPair) bool {
			ati, bti := a.To.Index(), b.To.Index()
			if ati == bti {
				return a.From.Index() < b.From.Index()
			}
			return ati < bti
		}),
	}
}

func (r *LikesRelationship) Clear() {
	r.btree.Clear()
}

func (w *World) LinkLikes(
	to, from Entity,
) {
	pair := LikesRelationshipPair{
		From: from,
		To:   to,
	}

	w.likesRelationships.btree.Set(pair)
}

func (w *World) UnlinkLikes(from, to Entity) {
	pair := LikesRelationshipPair{
		From: from,
		To:   to,
	}
	w.likesRelationships.btree.Delete(pair)
}

func (w *World) LikesIsLinked(from, to Entity) bool {
	pair := LikesRelationshipPair{
		From: from,
		To:   to,
	}

	_, ok := w.likesRelationships.btree.Get(pair)
	return ok
}

func (w *World) Likes(to Entity) func(yield func(from Entity) bool) {
	return func(yield func(from Entity) bool) {
		iter := w.likesRelationships.btree.Iter()
		iter.Seek(LikesRelationshipPair{To: to})
		end := LikesRelationshipPair{To: to + 1}

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

func (w *World) RemoveLikesRelationships(to Entity, froms ...Entity) {
	for _, from := range froms {
		pair := LikesRelationshipPair{
			To:   to,
			From: from,
		}

		w.likesRelationships.btree.Delete(pair)
	}
}

func (w *World) RemoveAllLikesRelationships(to Entity) {
	iter := w.likesRelationships.btree.Iter()
	end := LikesRelationshipPair{To: to + 1}

	for iter.Next() {
		item := iter.Item()
		if item.To >= end.To {
			break
		}

		w.likesRelationships.btree.Delete(item)
	}
}
