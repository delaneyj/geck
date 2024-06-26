package ecs

import "google.golang.org/protobuf/types/known/emptypb"

type Spacestation struct{}

//#region Events
//#endregion

func (e Entity) HasSpacestationTag() bool {
	return e.w.spacestationStore.Has(e)
}

func (e Entity) TagWithSpacestation() Entity {
	e.w.spacestationStore.Set(e.w.spacestationStore.zero, e)
	e.w.patch.SpacestationTags[e.val] = empty
	return e
}

func (e Entity) RemoveSpacestationTag() Entity {
	e.w.spacestationStore.Remove(e)
	e.w.patch.SpacestationTags[e.val] = nil
	return e
}

func (w *World) RemoveSpacestationTags(entities ...Entity) {
	w.spacestationStore.Remove(entities...)
	for _, entity := range entities {
		w.patch.SpacestationTags[entity.val] = nil
	}
}

//#region Iterators

type SpacestationReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Spacestation]
}

func (iter *SpacestationReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *SpacestationReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *SpacestationReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) SpacestationReadIter() *SpacestationReadIterator {
	iter := &SpacestationReadIterator{
		w:     w,
		store: w.spacestationStore,
	}
	iter.Reset()
	return iter
}

type SpacestationWriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Spacestation]
}

func (iter *SpacestationWriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *SpacestationWriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *SpacestationWriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) SpacestationWriteIter() *SpacestationWriteIterator {
	iter := &SpacestationWriteIterator{
		w:     w,
		store: w.spacestationStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) SpacestationEntities() []Entity {
	return w.spacestationStore.entities()
}

func (w *World) ApplySpacestationPatch(e Entity, pb *emptypb.Empty) Entity {
	if pb == nil {
		e.RemoveSpacestationTag()
	}
	e.TagWithSpacestation()
	return e
}
