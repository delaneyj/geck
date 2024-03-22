package ecs

type Gravity float32

func GravityFromFloat32(c float32) Gravity {
	return Gravity(c)
}

func (c Gravity) ToFloat32() float32 {
	return float32(c)
}

func (w *World) ResetGravity() float32 {
	return -9.800000
}

//#region Events
//#endregion

func (e Entity) HasGravity() bool {
	return e.w.gravitiesStore.Has(e)
}

func (e Entity) ReadGravity() (Gravity, bool) {
	return e.w.gravitiesStore.Read(e)
}

func (e Entity) RemoveGravity() Entity {
	e.w.gravitiesStore.Remove(e)

	return e
}

func (e Entity) WritableGravity() (*Gravity, bool) {
	return e.w.gravitiesStore.Writeable(e)
}

func (e Entity) SetGravity(other Gravity) Entity {
	e.w.gravitiesStore.Set(other, e)

	return e
}

func (w *World) SetGravities(c Gravity, entities ...Entity) {
	w.gravitiesStore.Set(c, entities...)
}

func (w *World) RemoveGravities(entities ...Entity) {
	w.gravitiesStore.Remove(entities...)
}

//#region Resources

// HasGravity checks if the world has a Gravity}}
func (w *World) HasGravityResource() bool {
	return w.resourceEntity.HasGravity()
}

// Retrieve the Gravity resource from the world
func (w *World) GravityResource() (Gravity, bool) {
	return w.resourceEntity.ReadGravity()
}

// Set the Gravity resource in the world
func (w *World) SetGravityResource(c Gravity) Entity {
	w.resourceEntity.SetGravity(c)
	return w.resourceEntity
}

// Remove the Gravity resource from the world
func (w *World) RemoveGravityResource() Entity {
	w.resourceEntity.RemoveGravity()

	return w.resourceEntity
}

//#endregion

//#region Iterators

type GravityReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Gravity]
}

func (iter *GravityReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *GravityReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *GravityReadIterator) NextGravity() (Entity, Gravity) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *GravityReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) GravityReadIter() *GravityReadIterator {
	iter := &GravityReadIterator{
		w:     w,
		store: w.gravitiesStore,
	}
	iter.Reset()
	return iter
}

type GravityWriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Gravity]
}

func (iter *GravityWriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *GravityWriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *GravityWriteIterator) NextGravity() (Entity, *Gravity) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--

	return e, c
}

func (iter *GravityWriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) GravityWriteIter() *GravityWriteIterator {
	iter := &GravityWriteIterator{
		w:     w,
		store: w.gravitiesStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) GravityEntities() []Entity {
	return w.gravitiesStore.entities()
}

func (w *World) SetGravitySortFn(lessThan func(a, b Entity) bool) {
	w.gravitiesStore.LessThan = lessThan
}

func (w *World) SortGravities() {
	w.gravitiesStore.Sort()
}
