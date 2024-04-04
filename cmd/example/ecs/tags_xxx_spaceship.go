package ecs

import "google.golang.org/protobuf/types/known/emptypb"

type Spaceship struct{}

//#region Events
//#endregion

func (e Entity) HasSpaceshipTag() bool {
	return e.w.spaceshipStore.Has(e)
}

func (e Entity) TagWithSpaceship() Entity {
	e.w.spaceshipStore.Set(e.w.spaceshipStore.zero, e)
	e.w.patch.SpaceshipTags[e.val] = empty
	return e
}

func (e Entity) RemoveSpaceshipTag() Entity {
	e.w.spaceshipStore.Remove(e)
	e.w.patch.SpaceshipTags[e.val] = nil
	return e
}

func (w *World) RemoveSpaceshipTags(entities ...Entity) {
	w.spaceshipStore.Remove(entities...)
	for _, entity := range entities {
		w.patch.SpaceshipTags[entity.val] = nil
	}
}

//#region Iterators

type SpaceshipReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Spaceship]
}

func (iter *SpaceshipReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *SpaceshipReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *SpaceshipReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) SpaceshipReadIter() *SpaceshipReadIterator {
	iter := &SpaceshipReadIterator{
		w:     w,
		store: w.spaceshipStore,
	}
	iter.Reset()
	return iter
}

type SpaceshipWriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Spaceship]
}

func (iter *SpaceshipWriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *SpaceshipWriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *SpaceshipWriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) SpaceshipWriteIter() *SpaceshipWriteIterator {
	iter := &SpaceshipWriteIterator{
		w:     w,
		store: w.spaceshipStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) SpaceshipEntities() []Entity {
	return w.spaceshipStore.entities()
}

func (w *World) ApplySpaceshipPatch(e Entity, pb *emptypb.Empty) Entity {
	if pb == nil {
		e.RemoveSpaceshipTag()
	}
	e.TagWithSpaceship()
	return e
}
