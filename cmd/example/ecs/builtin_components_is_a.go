package ecs

import (
	ecspb "github.com/delaneyj/geck/cmd/example/ecs/pb/gen/ecs/v1"
)

type IsA Entity

func IsAFromEntity(c Entity) IsA {
	return IsA(c)
}

func (c IsA) ToEntity() Entity {
	return Entity(c)
}

func (w *World) ResetIsA() Entity {
	return w.EntityFromU32(0)
}

func (c IsA) FromEntity(e Entity) IsA {
	return IsA(e)
}

//#region Events
//#endregion

func (e Entity) HasIsA() bool {
	return e.w.isAStore.Has(e)
}

func (e Entity) ReadIsA() (Entity, bool) {
	val, ok := e.w.isAStore.Read(e)
	if !ok {
		return Entity{}, false
	}
	return Entity(val), true
}

func (e Entity) RemoveIsA() Entity {
	e.w.isAStore.Remove(e)

	return e
}

func (e Entity) WritableIsA() (c *IsA, done func()) {
	var ok bool
	c, ok = e.w.isAStore.Writeable(e)
	if !ok {
		return nil, nil
	}
	return c, func() {
		e.w.patch.IsAComponents[e.val] = c.ToPB()
	}
}

func (e Entity) SetIsA(other Entity) Entity {
	e.w.isAStore.Set(IsA(other), e)

	e.w.patch.IsAComponents[e.val] = IsA(other).ToPB()
	return e
}

func (w *World) SetIsA(c IsA, entities ...Entity) {
	if len(entities) == 0 {
		panic("no entities provided, are you sure you didn't mean to call SetIsAResource?")
	}
	w.isAStore.Set(c, entities...)
	w.patch.IsAComponents[w.resourceEntity.val] = c.ToPB()
}

func (w *World) RemoveIsA(entities ...Entity) {
	w.isAStore.Remove(entities...)
	for _, entity := range entities {
		w.patch.IsAComponents[entity.val] = nil
	}
}

//#region Resources

// HasIsAResource checks if the world has a IsA}}
func (w *World) HasIsAResource() bool {
	return w.resourceEntity.HasIsA()
}

// IsAResource Retrieve the  resource from the world
func (w *World) IsAResource() (Entity, bool) {
	return w.resourceEntity.ReadIsA()
}

// SetIsAResource set the resource in the world
func (w *World) SetIsAResource(e Entity) Entity {
	w.resourceEntity.SetIsA(e)
	return w.resourceEntity
}

// RemoveIsAResource removes the resource from the world
func (w *World) RemoveIsAResource() Entity {
	w.resourceEntity.RemoveIsA()

	return w.resourceEntity
}

// WriteableIsAResource returns a writable reference to the resource
func (w *World) WriteableIsAResource() (c *IsA, done func()) {
	return w.resourceEntity.WritableIsA()
}

//#endregion

//#region Iterators

type IsAReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[IsA]
}

func (iter *IsAReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *IsAReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *IsAReadIterator) NextIsA() (Entity, IsA) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *IsAReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) IsAReadIter() *IsAReadIterator {
	iter := &IsAReadIterator{
		w:     w,
		store: w.isAStore,
	}
	iter.Reset()
	return iter
}

type IsAWriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[IsA]
}

func (iter *IsAWriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *IsAWriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *IsAWriteIterator) NextIsA() (Entity, *IsA, func()) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--
	done := func() {
		iter.w.patch.IsAComponents[e.val] = c.ToPB()
	}

	return e, c, done
}

func (iter *IsAWriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) IsAWriteIter() *IsAWriteIterator {
	iter := &IsAWriteIterator{
		w:     w,
		store: w.isAStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) IsAEntities() []Entity {
	return w.isAStore.entities()
}

func (w *World) SetIsASortFn(lessThan func(a, b Entity) bool) {
	w.isAStore.LessThan = lessThan
}

func (w *World) SortIsA() {
	w.isAStore.Sort()
}

func (w *World) ApplyIsAPatch(e Entity, patch *ecspb.IsAComponent) Entity {
	c := IsA(w.EntityFromU32(patch.Prototype))
	e.w.isAStore.Set(c, e)
	return e
}

func (c IsA) ToPB() *ecspb.IsAComponent {
	pb := &ecspb.IsAComponent{
		Prototype: c.val,
	}
	return pb
}
