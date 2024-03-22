package ecs

type DockedTo Entity

func DockedToFromEntity(c Entity) DockedTo {
	return DockedTo(c)
}

func (c DockedTo) ToEntity() Entity {
	return Entity(c)
}

func (w *World) ResetDockedTo() Entity {
	return w.EntityFromU32(0)
}

func (c DockedTo) FromEntity(e Entity) DockedTo {
	return DockedTo(e)
}

//#region Events
//#endregion

func (e Entity) HasDockedTo() bool {
	return e.w.dockedTosStore.Has(e)
}

func (e Entity) ReadDockedTo() (Entity, bool) {
	val, ok := e.w.dockedTosStore.Read(e)
	if !ok {
		return Entity{}, false
	}
	return Entity(val), true
}

func (e Entity) RemoveDockedTo() Entity {
	e.w.dockedTosStore.Remove(e)

	return e
}

func (e Entity) WritableDockedTo() (*DockedTo, bool) {
	return e.w.dockedTosStore.Writeable(e)
}

func (e Entity) SetDockedTo(other Entity) Entity {
	e.w.dockedTosStore.Set(DockedTo(other), e)

	return e
}

func (w *World) SetDockedTos(c DockedTo, entities ...Entity) {
	w.dockedTosStore.Set(c, entities...)
}

func (w *World) RemoveDockedTos(entities ...Entity) {
	w.dockedTosStore.Remove(entities...)
}

//#region Resources

// HasDockedTo checks if the world has a DockedTo}}
func (w *World) HasDockedToResource() bool {
	return w.resourceEntity.HasDockedTo()
}

// Retrieve the DockedTo resource from the world
func (w *World) DockedToResource() (Entity, bool) {
	return w.resourceEntity.ReadDockedTo()
}

// Set the DockedTo resource in the world
func (w *World) SetDockedToResource(c Entity) Entity {
	w.resourceEntity.SetDockedTo(c)

	return w.resourceEntity
}

// Remove the DockedTo resource from the world
func (w *World) RemoveDockedToResource() Entity {
	w.resourceEntity.RemoveDockedTo()

	return w.resourceEntity
}

//#endregion

//#region Iterators

type DockedToReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[DockedTo]
}

func (iter *DockedToReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *DockedToReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *DockedToReadIterator) NextDockedTo() (Entity, DockedTo) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *DockedToReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) DockedToReadIter() *DockedToReadIterator {
	iter := &DockedToReadIterator{
		w:     w,
		store: w.dockedTosStore,
	}
	iter.Reset()
	return iter
}

type DockedToWriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[DockedTo]
}

func (iter *DockedToWriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *DockedToWriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *DockedToWriteIterator) NextDockedTo() (Entity, *DockedTo) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--

	return e, c
}

func (iter *DockedToWriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) DockedToWriteIter() *DockedToWriteIterator {
	iter := &DockedToWriteIterator{
		w:     w,
		store: w.dockedTosStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) DockedToEntities() []Entity {
	return w.dockedTosStore.entities()
}

func (w *World) SetDockedToSortFn(lessThan func(a, b Entity) bool) {
	w.dockedTosStore.LessThan = lessThan
}

func (w *World) SortDockedTos() {
	w.dockedTosStore.Sort()
}
