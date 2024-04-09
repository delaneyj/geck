package ecs

import (
	ecspb "github.com/delaneyj/geck/cmd/example/ecs/pb/gen/ecs/v1"
)

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

func (c Grows) ToU32s() []uint32 {
	u32s := make([]uint32, len(c))
	for i, e := range c {
		u32s[i] = e.val
	}
	return u32s
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

	e.w.patch.GrowsComponents[e.val] = nil

	return e
}

func (e Entity) RemoveAllGrows() Entity {
	e.w.growsStore.Remove(e)

	return e
}

func (e Entity) WritableGrows() (c *Grows, done func()) {
	var ok bool
	c, ok = e.w.growsStore.Writeable(e)
	if !ok {
		return nil, nil
	}
	return c, func() {
		e.w.patch.GrowsComponents[e.val] = c.ToPB()
	}
}

func (e Entity) SetGrows(other ...Entity) Entity {
	e.w.growsStore.Set(Grows(other), e)

	e.w.patch.GrowsComponents[e.val] = Grows(other).ToPB()
	return e
}

func (w *World) SetGrows(c Grows, entities ...Entity) {
	if len(entities) == 0 {
		panic("no entities provided, are you sure you didn't mean to call SetGrowsResource?")
	}
	w.growsStore.Set(c, entities...)
	w.patch.GrowsComponents[w.resourceEntity.val] = c.ToPB()
}

func (w *World) RemoveGrows(entities ...Entity) {
	w.growsStore.Remove(entities...)
	for _, entity := range entities {
		w.patch.GrowsComponents[entity.val] = nil
	}
}

//#region Resources

// HasGrowsResource checks if the world has a Grows}}
func (w *World) HasGrowsResource() bool {
	return w.resourceEntity.HasGrows()
}

// GrowsResource Retrieve the  resource from the world
func (w *World) GrowsResource() ([]Entity, bool) {
	return w.resourceEntity.ReadGrows()
}

// SetGrowsResource set the resource in the world
func (w *World) SetGrowsResource(e ...Entity) Entity {
	w.resourceEntity.SetGrows(e...)
	return w.resourceEntity
}

// RemoveGrowsResource removes the resource from the world
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

func (iter *GrowsWriteIterator) NextGrows() (Entity, *Grows, func()) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--
	done := func() {
		iter.w.patch.GrowsComponents[e.val] = c.ToPB()
	}

	return e, c, done
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

func (w *World) ApplyGrowsPatch(e Entity, patch *ecspb.GrowsComponent) Entity {

	c := Grows(w.EntitiesFromU32s(patch.Entity...))
	e.w.growsStore.Set(c, e)
	return e
}

func (c Grows) ToPB() *ecspb.GrowsComponent {
	pb := &ecspb.GrowsComponent{
		Entity: c.ToU32s(),
	}
	return pb
}
