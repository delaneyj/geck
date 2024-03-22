package ecs

type Rotation struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
	W float32 `json:"w"`
}

func (w *World) ResetRotation() Rotation {
	return Rotation{
		X: 0.000000,
		Y: 0.000000,
		Z: 0.000000,
		W: 1.000000,
	}
}

//#region Events
//#endregion

func (e Entity) HasRotation() bool {
	return e.w.rotationsStore.Has(e)
}

func (e Entity) ReadRotation() (Rotation, bool) {
	return e.w.rotationsStore.Read(e)
}

func (e Entity) RemoveRotation() Entity {
	e.w.rotationsStore.Remove(e)

	return e
}

func (e Entity) WritableRotation() (*Rotation, bool) {
	return e.w.rotationsStore.Writeable(e)
}

func (e Entity) SetRotation(other Rotation) Entity {
	e.w.rotationsStore.Set(other, e)

	return e
}

func (w *World) SetRotations(c Rotation, entities ...Entity) {
	w.rotationsStore.Set(c, entities...)
}

func (w *World) RemoveRotations(entities ...Entity) {
	w.rotationsStore.Remove(entities...)
}

//#region Resources

// HasRotation checks if the world has a Rotation}}
func (w *World) HasRotationResource() bool {
	return w.resourceEntity.HasRotation()
}

// Retrieve the Rotation resource from the world
func (w *World) RotationResource() (Rotation, bool) {
	return w.resourceEntity.ReadRotation()
}

// Set the Rotation resource in the world
func (w *World) SetRotationResource(c Rotation) Entity {
	w.resourceEntity.SetRotation(c)
	return w.resourceEntity
}

// Remove the Rotation resource from the world
func (w *World) RemoveRotationResource() Entity {
	w.resourceEntity.RemoveRotation()

	return w.resourceEntity
}

//#endregion

//#region Iterators

type RotationReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Rotation]
}

func (iter *RotationReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *RotationReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *RotationReadIterator) NextRotation() (Entity, Rotation) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *RotationReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) RotationReadIter() *RotationReadIterator {
	iter := &RotationReadIterator{
		w:     w,
		store: w.rotationsStore,
	}
	iter.Reset()
	return iter
}

type RotationWriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Rotation]
}

func (iter *RotationWriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *RotationWriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *RotationWriteIterator) NextRotation() (Entity, *Rotation) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--

	return e, c
}

func (iter *RotationWriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) RotationWriteIter() *RotationWriteIterator {
	iter := &RotationWriteIterator{
		w:     w,
		store: w.rotationsStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) RotationEntities() []Entity {
	return w.rotationsStore.entities()
}

func (w *World) SetRotationSortFn(lessThan func(a, b Entity) bool) {
	w.rotationsStore.LessThan = lessThan
}

func (w *World) SortRotations() {
	w.rotationsStore.Sort()
}
