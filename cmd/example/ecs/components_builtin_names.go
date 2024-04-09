package ecs

import (
	ecspb "github.com/delaneyj/geck/cmd/example/ecs/pb/gen/ecs/v1"
)

type Name string

func NameFromString(c string) Name {
	return Name(c)
}

func (c Name) ToString() string {
	return string(c)
}

func (w *World) ResetName() string {
	return ""
}

//#region Events
//#endregion

func (e Entity) HasName() bool {
	return e.w.namesStore.Has(e)
}

func (e Entity) ReadName() (Name, bool) {
	return e.w.namesStore.Read(e)
}

func (e Entity) RemoveName() Entity {
	e.w.namesStore.Remove(e)

	e.w.patch.NameComponents[e.val] = nil
	return e
}

func (e Entity) WritableName() (c *Name, done func()) {
	var ok bool
	c, ok = e.w.namesStore.Writeable(e)
	if !ok {
		return nil, nil
	}
	return c, func() {
		e.w.patch.NameComponents[e.val] = c.ToPB()
	}
}

func (e Entity) SetName(other Name) Entity {
	e.w.namesStore.Set(other, e)

	e.w.patch.NameComponents[e.val] = Name(other).ToPB()
	return e
}

func (w *World) SetNames(c Name, entities ...Entity) {
	if len(entities) == 0 {
		panic("no entities provided, are you sure you didn't mean to call SetNameResource?")
	}
	w.namesStore.Set(c, entities...)
	w.patch.NameComponents[w.resourceEntity.val] = c.ToPB()
}

func (w *World) RemoveNames(entities ...Entity) {
	w.namesStore.Remove(entities...)
	for _, entity := range entities {
		w.patch.NameComponents[entity.val] = nil
	}
}

//#region Resources

// HasNameResource checks if the world has a Name}}
func (w *World) HasNameResource() bool {
	return w.resourceEntity.HasName()
}

// NameResource Retrieve the  resource from the world
func (w *World) NameResource() (Name, bool) {
	return w.resourceEntity.ReadName()
}

// SetNameResource set the resource in the world
func (w *World) SetNameResource(c Name) Entity {
	w.resourceEntity.SetName(c)
	return w.resourceEntity
}

// RemoveNameResource removes the resource from the world
func (w *World) RemoveNameResource() Entity {
	w.resourceEntity.RemoveName()

	return w.resourceEntity
}

//#endregion

//#region Iterators

type NameReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Name]
}

func (iter *NameReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *NameReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *NameReadIterator) NextName() (Entity, Name) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *NameReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) NameReadIter() *NameReadIterator {
	iter := &NameReadIterator{
		w:     w,
		store: w.namesStore,
	}
	iter.Reset()
	return iter
}

type NameWriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Name]
}

func (iter *NameWriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *NameWriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *NameWriteIterator) NextName() (Entity, *Name, func()) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--
	done := func() {
		iter.w.patch.NameComponents[e.val] = c.ToPB()
	}

	return e, c, done
}

func (iter *NameWriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) NameWriteIter() *NameWriteIterator {
	iter := &NameWriteIterator{
		w:     w,
		store: w.namesStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) NameEntities() []Entity {
	return w.namesStore.entities()
}

func (w *World) SetNameSortFn(lessThan func(a, b Entity) bool) {
	w.namesStore.LessThan = lessThan
}

func (w *World) SortNames() {
	w.namesStore.Sort()
}

func (w *World) ApplyNamePatch(e Entity, patch *ecspb.NameComponent) Entity {
	c := Name(patch.Value)
	e.w.namesStore.Set(c, e)
	return e
}

func (c Name) ToPB() *ecspb.NameComponent {
	pb := &ecspb.NameComponent{
		Value: c.ToString(),
	}
	return pb
}
