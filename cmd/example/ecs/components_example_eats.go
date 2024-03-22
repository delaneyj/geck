package ecs

type Eats struct {
	Entities []Entity `json:"entities"`
	Amounts  []uint8  `json:"amounts"`
}

func (w *World) ResetEats() Eats {
	return Eats{
		Entities: nil,
		Amounts:  nil,
	}
}

//#region Events
//#endregion

func (e Entity) HasEats() bool {
	return e.w.eatsStore.Has(e)
}

func (e Entity) ReadEats() (Eats, bool) {
	return e.w.eatsStore.Read(e)
}

func (e Entity) RemoveEats() Entity {
	e.w.eatsStore.Remove(e)

	return e
}

func (e Entity) WritableEats() (*Eats, bool) {
	return e.w.eatsStore.Writeable(e)
}

func (e Entity) SetEats(other Eats) Entity {
	e.w.eatsStore.Set(other, e)

	return e
}

func (w *World) SetEats(c Eats, entities ...Entity) {
	w.eatsStore.Set(c, entities...)
}

func (w *World) RemoveEats(entities ...Entity) {
	w.eatsStore.Remove(entities...)
}

//#region Resources

// HasEats checks if the world has a Eats}}
func (w *World) HasEatsResource() bool {
	return w.resourceEntity.HasEats()
}

// Retrieve the Eats resource from the world
func (w *World) EatsResource() (Eats, bool) {
	return w.resourceEntity.ReadEats()
}

// Set the Eats resource in the world
func (w *World) SetEatsResource(c Eats) Entity {
	w.resourceEntity.SetEats(c)
	return w.resourceEntity
}

// Remove the Eats resource from the world
func (w *World) RemoveEatsResource() Entity {
	w.resourceEntity.RemoveEats()

	return w.resourceEntity
}

//#endregion

//#region Iterators

type EatsReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Eats]
}

func (iter *EatsReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *EatsReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *EatsReadIterator) NextEats() (Entity, Eats) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *EatsReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) EatsReadIter() *EatsReadIterator {
	iter := &EatsReadIterator{
		w:     w,
		store: w.eatsStore,
	}
	iter.Reset()
	return iter
}

type EatsWriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Eats]
}

func (iter *EatsWriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *EatsWriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *EatsWriteIterator) NextEats() (Entity, *Eats) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--

	return e, c
}

func (iter *EatsWriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) EatsWriteIter() *EatsWriteIterator {
	iter := &EatsWriteIterator{
		w:     w,
		store: w.eatsStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) EatsEntities() []Entity {
	return w.eatsStore.entities()
}

func (w *World) SetEatsSortFn(lessThan func(a, b Entity) bool) {
	w.eatsStore.LessThan = lessThan
}

func (w *World) SortEats() {
	w.eatsStore.Sort()
}
