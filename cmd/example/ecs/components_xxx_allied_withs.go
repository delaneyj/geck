package ecs

type AlliedWith []Entity

func AlliedWithFromEntity(c []Entity) AlliedWith {
	return AlliedWith(c)
}

func (c AlliedWith) ToEntity() []Entity {
	return []Entity(c)
}

func (c AlliedWith) ToEntities() []Entity {
	entities := make([]Entity, len(c))
	copy(entities, c)
	return entities
}

func AlliedWithFromEntities(e ...Entity) AlliedWith {
	c := make(AlliedWith, len(e))
	copy(c, e)
	return c
}

//#region Events
//#endregion

func (e Entity) HasAlliedWith() bool {
	return e.w.alliedWithsStore.Has(e)
}

func (e Entity) ReadAlliedWith() ([]Entity, bool) {
	entities, ok := e.w.alliedWithsStore.Read(e)
	if !ok {
		return nil, false
	}
	return entities, true
}

func (e Entity) AlliedWithContains(other Entity) bool {
	entities, ok := e.w.alliedWithsStore.Read(e)
	if !ok {
		return false
	}
	for _, entity := range entities {
		if entity == other {
			return true
		}
	}
	return false
}

func (e Entity) RemoveAlliedWith(toRemove ...Entity) Entity {
	entities, ok := e.w.alliedWithsStore.Read(e)
	if !ok {
		return e
	}
	clean := make([]Entity, 0, len(entities))
	for _, tr := range toRemove {
		for _, entity := range entities {
			if entity != tr {
				clean = append(clean, entity)
			}
		}
	}
	e.w.alliedWithsStore.Set(clean, e)

	return e
}

func (e Entity) RemoveAllAlliedWith() Entity {
	e.w.alliedWithsStore.Remove(e)

	return e
}

func (e Entity) WritableAlliedWith() (*AlliedWith, bool) {
	return e.w.alliedWithsStore.Writeable(e)
}

func (e Entity) SetAlliedWith(other ...Entity) Entity {
	e.w.alliedWithsStore.Set(AlliedWith(other), e)

	return e
}

func (w *World) SetAlliedWiths(c AlliedWith, entities ...Entity) {
	w.alliedWithsStore.Set(c, entities...)
}

func (w *World) RemoveAlliedWiths(entities ...Entity) {
	w.alliedWithsStore.Remove(entities...)
}

//#region Resources

// HasAlliedWith checks if the world has a AlliedWith}}
func (w *World) HasAlliedWithResource() bool {
	return w.resourceEntity.HasAlliedWith()
}

// Retrieve the AlliedWith resource from the world
func (w *World) AlliedWithResource() ([]Entity, bool) {
	return w.resourceEntity.ReadAlliedWith()
}

// Set the AlliedWith resource in the world
func (w *World) SetAlliedWithResource(c ...Entity) Entity {
	w.resourceEntity.SetAlliedWith(c...)

	return w.resourceEntity
}

// Remove the AlliedWith resource from the world
func (w *World) RemoveAlliedWithResource() Entity {
	w.resourceEntity.RemoveAlliedWith()

	return w.resourceEntity
}

//#endregion

//#region Iterators

type AlliedWithReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[AlliedWith]
}

func (iter *AlliedWithReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *AlliedWithReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *AlliedWithReadIterator) NextAlliedWith() (Entity, AlliedWith) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *AlliedWithReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) AlliedWithReadIter() *AlliedWithReadIterator {
	iter := &AlliedWithReadIterator{
		w:     w,
		store: w.alliedWithsStore,
	}
	iter.Reset()
	return iter
}

type AlliedWithWriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[AlliedWith]
}

func (iter *AlliedWithWriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *AlliedWithWriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *AlliedWithWriteIterator) NextAlliedWith() (Entity, *AlliedWith) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--

	return e, c
}

func (iter *AlliedWithWriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) AlliedWithWriteIter() *AlliedWithWriteIterator {
	iter := &AlliedWithWriteIterator{
		w:     w,
		store: w.alliedWithsStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) AlliedWithEntities() []Entity {
	return w.alliedWithsStore.entities()
}

func (w *World) SetAlliedWithSortFn(lessThan func(a, b Entity) bool) {
	w.alliedWithsStore.LessThan = lessThan
}

func (w *World) SortAlliedWiths() {
	w.alliedWithsStore.Sort()
}
