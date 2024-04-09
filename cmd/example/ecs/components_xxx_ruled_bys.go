package ecs

import (
	ecspb "github.com/delaneyj/geck/cmd/example/ecs/pb/gen/ecs/v1"
)

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

func (e Entity) WritableRuledBy() (c *RuledBy, done func()) {
	var ok bool
	c, ok = e.w.ruledBysStore.Writeable(e)
	if !ok {
		return nil, nil
	}
	return c, func() {
		e.w.patch.RuledByComponents[e.val] = c.ToPB()
	}
}

func (e Entity) SetRuledBy(other Entity) Entity {
	e.w.ruledBysStore.Set(RuledBy(other), e)

	e.w.patch.RuledByComponents[e.val] = RuledBy(other).ToPB()
	return e
}

func (w *World) SetRuledBys(c RuledBy, entities ...Entity) {
	if len(entities) == 0 {
		panic("no entities provided, are you sure you didn't mean to call SetRuledByResource?")
	}
	w.ruledBysStore.Set(c, entities...)
	w.patch.RuledByComponents[w.resourceEntity.val] = c.ToPB()
}

func (w *World) RemoveRuledBys(entities ...Entity) {
	w.ruledBysStore.Remove(entities...)
	for _, entity := range entities {
		w.patch.RuledByComponents[entity.val] = nil
	}
}

//#region Resources

// HasRuledByResource checks if the world has a RuledBy}}
func (w *World) HasRuledByResource() bool {
	return w.resourceEntity.HasRuledBy()
}

// RuledByResource Retrieve the  resource from the world
func (w *World) RuledByResource() (Entity, bool) {
	return w.resourceEntity.ReadRuledBy()
}

// SetRuledByResource set the resource in the world
func (w *World) SetRuledByResource(e Entity) Entity {
	w.resourceEntity.SetRuledBy(e)
	return w.resourceEntity
}

// RemoveRuledByResource removes the resource from the world
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

func (iter *RuledByWriteIterator) NextRuledBy() (Entity, *RuledBy, func()) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--
	done := func() {
		iter.w.patch.RuledByComponents[e.val] = c.ToPB()
	}

	return e, c, done
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

func (w *World) ApplyRuledByPatch(e Entity, patch *ecspb.RuledByComponent) Entity {
	c := RuledBy(w.EntityFromU32(patch.Entity))
	e.w.ruledBysStore.Set(c, e)
	return e
}

func (c RuledBy) ToPB() *ecspb.RuledByComponent {
	pb := &ecspb.RuledByComponent{
		Entity: c.val,
	}
	return pb
}
