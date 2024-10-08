package ecs

import (
	ecspb "github.com/delaneyj/geck/cmd/example/ecs/pb/gen/ecs/v1"
)

type ChildOf Entity

func ChildOfFromEntity(c Entity) ChildOf {
	return ChildOf(c)
}

func (c ChildOf) ToEntity() Entity {
	return Entity(c)
}

func (w *World) ResetChildOf() Entity {
	return w.EntityFromU32(0)
}

func (c ChildOf) FromEntity(e Entity) ChildOf {
	return ChildOf(e)
}

//#region Events
//#endregion

func (e Entity) HasChildOf() bool {
	return e.w.childOfStore.Has(e)
}

func (e Entity) ReadChildOf() (Entity, bool) {
	val, ok := e.w.childOfStore.Read(e)
	if !ok {
		return Entity{}, false
	}
	return Entity(val), true
}

func (e Entity) RemoveChildOf() Entity {
	e.w.childOfStore.Remove(e)

	return e
}

func (e Entity) WritableChildOf() (c *ChildOf, done func()) {
	var ok bool
	c, ok = e.w.childOfStore.Writeable(e)
	if !ok {
		return nil, nil
	}
	return c, func() {
		e.w.patch.ChildOfComponents[e.val] = c.ToPB()
	}
}

func (e Entity) SetChildOf(other Entity) Entity {
	e.w.childOfStore.Set(ChildOf(other), e)

	e.w.patch.ChildOfComponents[e.val] = ChildOf(other).ToPB()
	return e
}

func (w *World) SetChildOf(c ChildOf, entities ...Entity) {
	if len(entities) == 0 {
		panic("no entities provided, are you sure you didn't mean to call SetChildOfResource?")
	}
	w.childOfStore.Set(c, entities...)
	w.patch.ChildOfComponents[w.resourceEntity.val] = c.ToPB()
}

func (w *World) RemoveChildOf(entities ...Entity) {
	w.childOfStore.Remove(entities...)
	for _, entity := range entities {
		w.patch.ChildOfComponents[entity.val] = nil
	}
}

//#region Resources

// HasChildOfResource checks if the world has a ChildOf}}
func (w *World) HasChildOfResource() bool {
	return w.resourceEntity.HasChildOf()
}

// ChildOfResource Retrieve the  resource from the world
func (w *World) ChildOfResource() (Entity, bool) {
	return w.resourceEntity.ReadChildOf()
}

// SetChildOfResource set the resource in the world
func (w *World) SetChildOfResource(e Entity) Entity {
	w.resourceEntity.SetChildOf(e)
	return w.resourceEntity
}

// RemoveChildOfResource removes the resource from the world
func (w *World) RemoveChildOfResource() Entity {
	w.resourceEntity.RemoveChildOf()

	return w.resourceEntity
}

// WriteableChildOfResource returns a writable reference to the resource
func (w *World) WriteableChildOfResource() (c *ChildOf, done func()) {
	return w.resourceEntity.WritableChildOf()
}

//#endregion

//#region Iterators

type ChildOfReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[ChildOf]
}

func (iter *ChildOfReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *ChildOfReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *ChildOfReadIterator) NextChildOf() (Entity, ChildOf) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *ChildOfReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) ChildOfReadIter() *ChildOfReadIterator {
	iter := &ChildOfReadIterator{
		w:     w,
		store: w.childOfStore,
	}
	iter.Reset()
	return iter
}

type ChildOfWriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[ChildOf]
}

func (iter *ChildOfWriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *ChildOfWriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *ChildOfWriteIterator) NextChildOf() (Entity, *ChildOf, func()) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--
	done := func() {
		iter.w.patch.ChildOfComponents[e.val] = c.ToPB()
	}

	return e, c, done
}

func (iter *ChildOfWriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) ChildOfWriteIter() *ChildOfWriteIterator {
	iter := &ChildOfWriteIterator{
		w:     w,
		store: w.childOfStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) ChildOfEntities() []Entity {
	return w.childOfStore.entities()
}

func (w *World) SetChildOfSortFn(lessThan func(a, b Entity) bool) {
	w.childOfStore.LessThan = lessThan
}

func (w *World) SortChildOf() {
	w.childOfStore.Sort()
}

func (w *World) ApplyChildOfPatch(e Entity, patch *ecspb.ChildOfComponent) Entity {
	c := ChildOf(w.EntityFromU32(patch.Parent))
	e.w.childOfStore.Set(c, e)
	return e
}

func (c ChildOf) ToPB() *ecspb.ChildOfComponent {
	pb := &ecspb.ChildOfComponent{
		Parent: c.val,
	}
	return pb
}
