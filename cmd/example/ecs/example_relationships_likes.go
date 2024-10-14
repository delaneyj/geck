package ecs

import (
	"github.com/tidwall/btree"
)

type LikesRelationshipPair struct {
	Source Entity
	Target Entity
}

type LikesRelationship struct {
	btree *btree.BTreeG[LikesRelationshipPair]
}

func NewLikesRelationship() *LikesRelationship {
	return &LikesRelationship{
		btree: btree.NewBTreeG[LikesRelationshipPair](func(a, b LikesRelationshipPair) bool {
			ati, bti := a.Target.Index(), b.Target.Index()
			if ati == bti {
				return a.Source.Index() < b.Source.Index()
			}
			return ati < bti
		}),
	}
}

func (r *LikesRelationship) Clear() {
	r.btree.Clear()
}

func (w *World) LinkLikes(target Entity, sources ...Entity) {
	for _, source := range sources {
		pair := LikesRelationshipPair{
			Target: target,
			Source: source,
		}

		w.likesRelationships.btree.Set(pair)
	}
}

func (w *World) UnlinkLikes(target Entity, sources ...Entity) {
	for _, source := range sources {
		pair := LikesRelationshipPair{
			Target: target,
			Source: source,
		}

		w.likesRelationships.btree.Delete(pair)
	}
}

func (w *World) LikesIsLinked(source, target Entity) bool {
	pair := LikesRelationshipPair{
		Source: source,
		Target: target,
	}

	_, ok := w.likesRelationships.btree.Get(pair)
	return ok
}

func (w *World) LikesSources(target Entity) func(yield func(source Entity) bool) {
	return func(yield func(source Entity) bool) {
		iter := w.likesRelationships.btree.Iter()
		iter.Seek(LikesRelationshipPair{Target: target})
		end := LikesRelationshipPair{Target: target + 1}

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

func (w *World) RemoveLikesRelationships(target Entity, sources ...Entity) {
	for _, source := range sources {
		pair := LikesRelationshipPair{
			Target: target,
			Source: source,
		}

		w.likesRelationships.btree.Delete(pair)
	}
}

func (w *World) RemoveAllLikesRelationships(target Entity) {
	iter := w.likesRelationships.btree.Iter()
	iter.Seek(LikesRelationshipPair{Target: target})
	end := LikesRelationshipPair{Target: target + 1}

	for iter.Next() {
		item := iter.Item()
		if item.Target >= end.Target {
			break
		}

		w.likesRelationships.btree.Delete(item)
	}
}
