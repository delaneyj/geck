package ecs

type Grows []Entity

func GrowsFromEntity(c []Entity) Grows {
	return Grows(c)
}

func (c Grows) ToEntity() []Entity {
	return []Entity(c)
}

func (c Grows) ToEntities() []Entity {
	entities := make([]Entity, len(c))
	copy(entities, c)
	return entities
}

func GrowsFromEntities(e ...Entity) Grows {
	c := make(Grows, len(e))
	copy(c, e)
	return c
}

//#region Events
//#endregion

func (e Entity) HasGrows() bool {
	return e.w.growsStore.Has(e)
}

func (e Entity) ReadGrows() ([]Entity, bool) {
	entities, ok := e.w.growsStore.Read(e)
	if !ok {
		return nil, false
	}
	return entities, true
}

func (e Entity) GrowsContains(other Entity) bool {
	entities, ok := e.w.growsStore.Read(e)
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

func (e Entity) RemoveGrows(toRemove ...Entity) Entity {
	entities, ok := e.w.growsStore.Read(e)
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
	e.w.growsStore.Set(clean, e)

	return e
}

func (e Entity) RemoveAllGrows() Entity {
	e.w.growsStore.Remove(e)

	return e
}

func (e Entity) WritableGrows() (*Grows, bool) {
	return e.w.growsStore.Writeable(e)
}

func (e Entity) SetGrows(other ...Entity) Entity {
	e.w.growsStore.Set(Grows(other), e)

	return e
}

func (w *World) SetGrows(c Grows, entities ...Entity) {
	w.growsStore.Set(c, entities...)
}

func (w *World) RemoveGrows(entities ...Entity) {
	w.growsStore.Remove(entities...)
}

//#region Resources

// HasGrows checks if the world has a Grows}}
func (w *World) HasGrowsResource() bool {
	return w.resourceEntity.HasGrows()
}

// Retrieve the Grows resource from the world
func (w *World) GrowsResource() ([]Entity, bool) {
	return w.resourceEntity.ReadGrows()
}

// Set the Grows resource in the world
func (w *World) SetGrowsResource(c ...Entity) Entity {
	w.resourceEntity.SetGrows(c...)

	return w.resourceEntity
}

// Remove the Grows resource from the world
func (w *World) RemoveGrowsResource() Entity {
	w.resourceEntity.RemoveGrows()

	return w.resourceEntity
}

//#endregion

//#region Iterators

type GrowsReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Grows]
}

func (iter *GrowsReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *GrowsReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *GrowsReadIterator) NextGrows() (Entity, Grows) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *GrowsReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) GrowsReadIter() *GrowsReadIterator {
	iter := &GrowsReadIterator{
		w:     w,
		store: w.growsStore,
	}
	iter.Reset()
	return iter
}

type GrowsWriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Grows]
}

func (iter *GrowsWriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *GrowsWriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *GrowsWriteIterator) NextGrows() (Entity, *Grows) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--

	return e, c
}

func (iter *GrowsWriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) GrowsWriteIter() *GrowsWriteIterator {
	iter := &GrowsWriteIterator{
		w:     w,
		store: w.growsStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) GrowsEntities() []Entity {
	return w.growsStore.entities()
}

func (w *World) SetGrowsSortFn(lessThan func(a, b Entity) bool) {
	w.growsStore.LessThan = lessThan
}

func (w *World) SortGrows() {
	w.growsStore.Sort()
}
