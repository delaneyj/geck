package ecs

import (
	ecspb "github.com/delaneyj/geck/cmd/example/ecs/pb/gen/ecs/v1"
)

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

func (e Entity) WritableDockedTo() (c *DockedTo, done func()) {
	var ok bool
	c, ok = e.w.dockedTosStore.Writeable(e)
	if !ok {
		return nil, nil
	}
	return c, func() {
		e.w.patch.DockedToComponents[e.val] = c.ToPB()
	}
}

func (e Entity) SetDockedTo(other Entity) Entity {
	e.w.dockedTosStore.Set(DockedTo(other), e)

	e.w.patch.DockedToComponents[e.val] = DockedTo(other).ToPB()
	return e
}

func (w *World) SetDockedTos(c DockedTo, entities ...Entity) {
	if len(entities) == 0 {
		panic("no entities provided, are you sure you didn't mean to call SetDockedToResource?")
	}
	w.dockedTosStore.Set(c, entities...)
	w.patch.DockedToComponents[w.resourceEntity.val] = c.ToPB()
}

func (w *World) RemoveDockedTos(entities ...Entity) {
	w.dockedTosStore.Remove(entities...)
	for _, entity := range entities {
		w.patch.DockedToComponents[entity.val] = nil
	}
}

//#region Resources

// HasDockedToResource checks if the world has a DockedTo}}
func (w *World) HasDockedToResource() bool {
	return w.resourceEntity.HasDockedTo()
}

// DockedToResource Retrieve the  resource from the world
func (w *World) DockedToResource() (Entity, bool) {
	return w.resourceEntity.ReadDockedTo()
}

// SetDockedToResource set the resource in the world
func (w *World) SetDockedToResource(e Entity) Entity {
	w.resourceEntity.SetDockedTo(e)
	return w.resourceEntity
}

// RemoveDockedToResource removes the resource from the world
func (w *World) RemoveDockedToResource() Entity {
	w.resourceEntity.RemoveDockedTo()

	return w.resourceEntity
}

// WriteableDockedToResource returns a writable reference to the resource
func (w *World) WriteableDockedToResource() (c *DockedTo, done func()) {
	return w.resourceEntity.WritableDockedTo()
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

func (iter *DockedToWriteIterator) NextDockedTo() (Entity, *DockedTo, func()) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--
	done := func() {
		iter.w.patch.DockedToComponents[e.val] = c.ToPB()
	}

	return e, c, done
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

func (w *World) ApplyDockedToPatch(e Entity, patch *ecspb.DockedToComponent) Entity {
	c := DockedTo(w.EntityFromU32(patch.Entity))
	e.w.dockedTosStore.Set(c, e)
	return e
}

func (c DockedTo) ToPB() *ecspb.DockedToComponent {
	pb := &ecspb.DockedToComponent{
		Entity: c.val,
	}
	return pb
}
