package ecs

type RuledBy Entity

func RuledByFromEntity(c Entity) RuledBy {
	return RuledBy(c)
}

func (c RuledBy) ToEntity() Entity {
	return Entity(c)
}

func (w *World) ResetRuledBy() Entity {
	return w.EntityFromU32(0)
}

func (c RuledBy) FromEntity(e Entity) RuledBy {
	return RuledBy(e)
}

//#region Events
//#endregion

func (e Entity) HasRuledBy() bool {
	return e.w.ruledBysStore.Has(e)
}

func (e Entity) ReadRuledBy() (Entity, bool) {
	val, ok := e.w.ruledBysStore.Read(e)
	if !ok {
		return Entity{}, false
	}
	return Entity(val), true
}

func (e Entity) RemoveRuledBy() Entity {
	e.w.ruledBysStore.Remove(e)

	return e
}

func (e Entity) WritableRuledBy() (*RuledBy, bool) {
	return e.w.ruledBysStore.Writeable(e)
}

func (e Entity) SetRuledBy(other Entity) Entity {
	e.w.ruledBysStore.Set(RuledBy(other), e)

	return e
}

func (w *World) SetRuledBys(c RuledBy, entities ...Entity) {
	w.ruledBysStore.Set(c, entities...)
}

func (w *World) RemoveRuledBys(entities ...Entity) {
	w.ruledBysStore.Remove(entities...)
}

//#region Resources

// HasRuledBy checks if the world has a RuledBy}}
func (w *World) HasRuledByResource() bool {
	return w.resourceEntity.HasRuledBy()
}

// Retrieve the RuledBy resource from the world
func (w *World) RuledByResource() (Entity, bool) {
	return w.resourceEntity.ReadRuledBy()
}

// Set the RuledBy resource in the world
func (w *World) SetRuledByResource(c Entity) Entity {
	w.resourceEntity.SetRuledBy(c)

	return w.resourceEntity
}

// Remove the RuledBy resource from the world
func (w *World) RemoveRuledByResource() Entity {
	w.resourceEntity.RemoveRuledBy()

	return w.resourceEntity
}

//#endregion

//#region Iterators

type RuledByReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[RuledBy]
}

func (iter *RuledByReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *RuledByReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *RuledByReadIterator) NextRuledBy() (Entity, RuledBy) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *RuledByReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) RuledByReadIter() *RuledByReadIterator {
	iter := &RuledByReadIterator{
		w:     w,
		store: w.ruledBysStore,
	}
	iter.Reset()
	return iter
}

type RuledByWriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[RuledBy]
}

func (iter *RuledByWriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *RuledByWriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *RuledByWriteIterator) NextRuledBy() (Entity, *RuledBy) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--

	return e, c
}

func (iter *RuledByWriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) RuledByWriteIter() *RuledByWriteIterator {
	iter := &RuledByWriteIterator{
		w:     w,
		store: w.ruledBysStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) RuledByEntities() []Entity {
	return w.ruledBysStore.entities()
}

func (w *World) SetRuledBySortFn(lessThan func(a, b Entity) bool) {
	w.ruledBysStore.LessThan = lessThan
}

func (w *World) SortRuledBys() {
	w.ruledBysStore.Sort()
}
