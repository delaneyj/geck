package ecs

import (
	ecspb "github.com/delaneyj/geck/cmd/example/ecs/pb/gen/ecs/v1"
)

type Direction EnumDirection

func DirectionFromEnumDirection(c EnumDirection) Direction {
	return Direction(c)
}

func (c Direction) ToEnumDirection() EnumDirection {
	return EnumDirection(c)
}

func (w *World) ResetDirection() EnumDirection {
	return EnumDirection(0)
}

//#region Events
//#endregion

func (e Entity) HasDirection() bool {
	return e.w.directionsStore.Has(e)
}

func (e Entity) ReadDirection() (Direction, bool) {
	return e.w.directionsStore.Read(e)
}

func (e Entity) MustReadDirection() Direction {
	c, ok := e.w.directionsStore.Read(e)
	if !ok {
		panic("Direction not found")
	}
	return c
}

func (e Entity) RemoveDirection() Entity {
	e.w.directionsStore.Remove(e)

	e.w.patch.DirectionComponents[e.val] = nil
	return e
}

func (e Entity) WritableDirection() (c *Direction, done func()) {
	var ok bool
	c, ok = e.w.directionsStore.Writeable(e)
	if !ok {
		return nil, nil
	}
	return c, func() {
		e.w.patch.DirectionComponents[e.val] = c.ToPB()
	}
}

func (e Entity) SetDirection(other Direction) Entity {
	e.w.directionsStore.Set(other, e)

	e.w.patch.DirectionComponents[e.val] = Direction(other).ToPB()
	return e
}

func (w *World) SetDirections(c Direction, entities ...Entity) {
	if len(entities) == 0 {
		panic("no entities provided, are you sure you didn't mean to call SetDirectionResource?")
	}
	w.directionsStore.Set(c, entities...)
	w.patch.DirectionComponents[w.resourceEntity.val] = c.ToPB()
}

func (w *World) RemoveDirections(entities ...Entity) {
	w.directionsStore.Remove(entities...)
	for _, entity := range entities {
		w.patch.DirectionComponents[entity.val] = nil
	}
}

//#region Resources

// HasDirectionResource checks if the world has a Direction}}
func (w *World) HasDirectionResource() bool {
	return w.resourceEntity.HasDirection()
}

// DirectionResource Retrieve the  resource from the world
func (w *World) DirectionResource() (Direction, bool) {
	return w.resourceEntity.ReadDirection()
}

// SetDirectionResource set the resource in the world
func (w *World) SetDirectionResource(c Direction) Entity {
	w.resourceEntity.SetDirection(c)
	return w.resourceEntity
}

// RemoveDirectionResource removes the resource from the world
func (w *World) RemoveDirectionResource() Entity {
	w.resourceEntity.RemoveDirection()

	return w.resourceEntity
}

// WriteableDirectionResource returns a writable reference to the resource
func (w *World) WriteableDirectionResource() (c *Direction, done func()) {
	return w.resourceEntity.WritableDirection()
}

//#endregion

//#region Iterators

type DirectionReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Direction]
}

func (iter *DirectionReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *DirectionReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *DirectionReadIterator) NextDirection() (Entity, Direction) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *DirectionReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) DirectionReadIter() *DirectionReadIterator {
	iter := &DirectionReadIterator{
		w:     w,
		store: w.directionsStore,
	}
	iter.Reset()
	return iter
}

type DirectionWriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Direction]
}

func (iter *DirectionWriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *DirectionWriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *DirectionWriteIterator) NextDirection() (Entity, *Direction, func()) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--
	done := func() {
		iter.w.patch.DirectionComponents[e.val] = c.ToPB()
	}

	return e, c, done
}

func (iter *DirectionWriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) DirectionWriteIter() *DirectionWriteIterator {
	iter := &DirectionWriteIterator{
		w:     w,
		store: w.directionsStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) DirectionEntities() []Entity {
	return w.directionsStore.entities()
}

func (w *World) SetDirectionSortFn(lessThan func(a, b Entity) bool) {
	w.directionsStore.LessThan = lessThan
}

func (w *World) SortDirections() {
	w.directionsStore.Sort()
}

func (w *World) ApplyDirectionPatch(e Entity, patch *ecspb.DirectionComponent) Entity {
	c := Direction(patch.Values)
	e.w.directionsStore.Set(c, e)
	return e
}

func (c Direction) ToPB() *ecspb.DirectionComponent {
	pb := &ecspb.DirectionComponent{
		Values: ecspb.DirectionEnum(c.ToEnumDirection()),
	}
	return pb
}
