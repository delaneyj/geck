package ecs

import (
	ecspb "github.com/delaneyj/geck/cmd/example/ecs/pb/gen/ecs/v1"
)

type Position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

func (w *World) ResetPosition() Position {
	return Position{
		X: 0.000000,
		Y: 0.000000,
		Z: 0.000000,
	}
}

//#region Events
//#endregion

func (e Entity) HasPosition() bool {
	return e.w.positionsStore.Has(e)
}

func (e Entity) ReadPosition() (Position, bool) {
	return e.w.positionsStore.Read(e)
}

func (e Entity) RemovePosition() Entity {
	e.w.positionsStore.Remove(e)

	e.w.patch.PositionComponents[e.val] = nil
	return e
}

func (e Entity) WritablePosition() (c *Position, done func()) {
	var ok bool
	c, ok = e.w.positionsStore.Writeable(e)
	if !ok {
		return nil, nil
	}
	return c, func() {
		e.w.patch.PositionComponents[e.val] = c.ToPB()
	}
}

func (e Entity) SetPosition(other Position) Entity {
	e.w.positionsStore.Set(other, e)

	e.w.PositionVelocitySet.PossibleUpdate(e)
	e.w.patch.PositionComponents[e.val] = Position(other).ToPB()
	return e
}

func (e Entity) SetPositionValues(
	x0 float32,
	y1 float32,
	z2 float32,
) Entity {
	err := e.w.positionsStore.Set(Position{
		X: x0,
		Y: y1,
		Z: z2,
	}, e)
	if err != nil {
		panic(err)
	}
	pb := &ecspb.PositionComponent{
		X: x0,
		Y: y1,
		Z: z2,
	}
	e.w.patch.PositionComponents[e.w.resourceEntity.val] = pb
	return e
}

func (w *World) SetPositions(c Position, entities ...Entity) {
	if len(entities) == 0 {
		panic("no entities provided, are you sure you didn't mean to call SetPositionResource?")
	}
	w.positionsStore.Set(c, entities...)
	w.PositionVelocitySet.PossibleUpdate(entities...)
	w.patch.PositionComponents[w.resourceEntity.val] = c.ToPB()
}

func (w *World) RemovePositions(entities ...Entity) {
	w.positionsStore.Remove(entities...)
	w.PositionVelocitySet.PossibleUpdate(entities...)
	for _, entity := range entities {
		w.patch.PositionComponents[entity.val] = nil
	}
}

//#region Resources

// HasPositionResource checks if the world has a Position}}
func (w *World) HasPositionResource() bool {
	return w.resourceEntity.HasPosition()
}

// PositionResource Retrieve the  resource from the world
func (w *World) PositionResource() (Position, bool) {
	return w.resourceEntity.ReadPosition()
}

// SetPositionResource set the resource in the world
func (w *World) SetPositionResource(c Position) Entity {
	w.resourceEntity.SetPosition(c)
	return w.resourceEntity
}
func (w *World) SetPositionResourceValues(
	x0 float32,
	y1 float32,
	z2 float32,
) Entity {
	w.resourceEntity.SetPositionValues(
		x0,
		y1,
		z2,
	)
	return w.resourceEntity
}

// RemovePositionResource removes the resource from the world
func (w *World) RemovePositionResource() Entity {
	w.resourceEntity.RemovePosition()

	w.PositionVelocitySet.PossibleUpdate(w.resourceEntity)
	return w.resourceEntity
}

//#endregion

//#region Iterators

type PositionReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Position]
}

func (iter *PositionReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *PositionReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *PositionReadIterator) NextPosition() (Entity, Position) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *PositionReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) PositionReadIter() *PositionReadIterator {
	iter := &PositionReadIterator{
		w:     w,
		store: w.positionsStore,
	}
	iter.Reset()
	return iter
}

type PositionWriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Position]
}

func (iter *PositionWriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *PositionWriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *PositionWriteIterator) NextPosition() (Entity, *Position, func()) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--
	done := func() {
		iter.w.patch.PositionComponents[e.val] = c.ToPB()
	}

	return e, c, done
}

func (iter *PositionWriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) PositionWriteIter() *PositionWriteIterator {
	iter := &PositionWriteIterator{
		w:     w,
		store: w.positionsStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) PositionEntities() []Entity {
	return w.positionsStore.entities()
}

func (w *World) SetPositionSortFn(lessThan func(a, b Entity) bool) {
	w.positionsStore.LessThan = lessThan
}

func (w *World) SortPositions() {
	w.positionsStore.Sort()
}

func (w *World) ApplyPositionPatch(e Entity, patch *ecspb.PositionComponent) Entity {
	c := Position{
		X: patch.X,
		Y: patch.Y,
		Z: patch.Z,
	}
	e.w.positionsStore.Set(c, e)
	return e
}

func (c Position) ToPB() *ecspb.PositionComponent {
	pb := &ecspb.PositionComponent{
		X: c.X,
		Y: c.Y,
		Z: c.Z,
	}
	return pb
}
