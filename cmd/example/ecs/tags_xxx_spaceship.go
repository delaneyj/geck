package ecs

type Spaceship struct{}

//#region Events
//#endregion

func (e Entity) HasSpaceshipTag() bool {
	return e.w.spaceshipStore.Has(e)
}

func (e Entity) TagWithSpaceship() Entity {
	e.w.spaceshipStore.Set(e.w.spaceshipStore.zero, e)
	return e
}

func (e Entity) RemoveSpaceshipTag() Entity {
	e.w.spaceshipStore.Remove(e)
	return e
}

func (w *World) RemoveSpaceshipTags(entities ...Entity) {
	w.spaceshipStore.Remove(entities...)
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
