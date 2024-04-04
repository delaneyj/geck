package ecs

import "google.golang.org/protobuf/types/known/emptypb"

type Planet struct{}

//#region Events
//#endregion

func (e Entity) HasPlanetTag() bool {
	return e.w.planetStore.Has(e)
}

func (e Entity) TagWithPlanet() Entity {
	e.w.planetStore.Set(e.w.planetStore.zero, e)
	e.w.patch.PlanetTags[e.val] = empty
	return e
}

func (e Entity) RemovePlanetTag() Entity {
	e.w.planetStore.Remove(e)
	e.w.patch.PlanetTags[e.val] = nil
	return e
}

func (w *World) RemovePlanetTags(entities ...Entity) {
	w.planetStore.Remove(entities...)
	for _, entity := range entities {
		w.patch.PlanetTags[entity.val] = nil
	}
}

//#region Iterators

type PlanetReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Planet]
}

func (iter *PlanetReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *PlanetReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *PlanetReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) PlanetReadIter() *PlanetReadIterator {
	iter := &PlanetReadIterator{
		w:     w,
		store: w.planetStore,
	}
	iter.Reset()
	return iter
}

type PlanetWriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Planet]
}

func (iter *PlanetWriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *PlanetWriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *PlanetWriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) PlanetWriteIter() *PlanetWriteIterator {
	iter := &PlanetWriteIterator{
		w:     w,
		store: w.planetStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) PlanetEntities() []Entity {
	return w.planetStore.entities()
}

func (w *World) ApplyPlanetPatch(e Entity, pb *emptypb.Empty) Entity {
	if pb == nil {
		e.RemovePlanetTag()
	}
	e.TagWithPlanet()
	return e
}
