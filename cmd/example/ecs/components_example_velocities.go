package ecs

import (
	ecspb "github.com/delaneyj/geck/cmd/example/ecs/pb/gen/ecs/v1"
)

type Velocity struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

func (w *World) ResetVelocity() Velocity {
	return Velocity{
		X: 0.000000,
		Y: 0.000000,
		Z: 0.000000,
	}
}

//#region Events
//#endregion

func (e Entity) HasVelocity() bool {
	return e.w.velocitiesStore.Has(e)
}

func (e Entity) ReadVelocity() (Velocity, bool) {
	return e.w.velocitiesStore.Read(e)
}

func (e Entity) RemoveVelocity() Entity {
	e.w.velocitiesStore.Remove(e)

	e.w.patch.VelocityComponents[e.val] = nil
	return e
}

func (e Entity) WritableVelocity() (c *Velocity, done func()) {
	var ok bool
	c, ok = e.w.velocitiesStore.Writeable(e)
	if !ok {
		return nil, nil
	}
	return c, func() {
		e.w.patch.VelocityComponents[e.val] = c.ToPB()
	}
}

func (e Entity) SetVelocity(other Velocity) Entity {
	e.w.velocitiesStore.Set(other, e)

	e.w.PositionVelocitySet.PossibleUpdate(e)
	e.w.patch.VelocityComponents[e.val] = Velocity(other).ToPB()
	return e
}

func (e Entity) SetVelocityValues(
	x0 float32,
	y1 float32,
	z2 float32,
) Entity {
	err := e.w.velocitiesStore.Set(Velocity{
		X: x0,
		Y: y1,
		Z: z2,
	}, e)
	if err != nil {
		panic(err)
	}
	pb := &ecspb.VelocityComponent{
		X: x0,
		Y: y1,
		Z: z2,
	}
	e.w.patch.VelocityComponents[e.w.resourceEntity.val] = pb
	return e
}

func (w *World) SetVelocities(c Velocity, entities ...Entity) {
	if len(entities) == 0 {
		panic("no entities provided, are you sure you didn't mean to call SetVelocityResource?")
	}
	w.velocitiesStore.Set(c, entities...)
	w.PositionVelocitySet.PossibleUpdate(entities...)
	w.patch.VelocityComponents[w.resourceEntity.val] = c.ToPB()
}

func (w *World) RemoveVelocities(entities ...Entity) {
	w.velocitiesStore.Remove(entities...)
	w.PositionVelocitySet.PossibleUpdate(entities...)
	for _, entity := range entities {
		w.patch.VelocityComponents[entity.val] = nil
	}
}

//#region Resources

// HasVelocityResource checks if the world has a Velocity}}
func (w *World) HasVelocityResource() bool {
	return w.resourceEntity.HasVelocity()
}

// VelocityResource Retrieve the  resource from the world
func (w *World) VelocityResource() (Velocity, bool) {
	return w.resourceEntity.ReadVelocity()
}

// SetVelocityResource set the resource in the world
func (w *World) SetVelocityResource(c Velocity) Entity {
	w.resourceEntity.SetVelocity(c)
	return w.resourceEntity
}
func (w *World) SetVelocityResourceValues(
	x0 float32,
	y1 float32,
	z2 float32,
) Entity {
	w.resourceEntity.SetVelocityValues(
		x0,
		y1,
		z2,
	)
	return w.resourceEntity
}

// RemoveVelocityResource removes the resource from the world
func (w *World) RemoveVelocityResource() Entity {
	w.resourceEntity.RemoveVelocity()

	w.PositionVelocitySet.PossibleUpdate(w.resourceEntity)
	return w.resourceEntity
}

//#endregion

//#region Iterators

type VelocityReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Velocity]
}

func (iter *VelocityReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *VelocityReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *VelocityReadIterator) NextVelocity() (Entity, Velocity) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *VelocityReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) VelocityReadIter() *VelocityReadIterator {
	iter := &VelocityReadIterator{
		w:     w,
		store: w.velocitiesStore,
	}
	iter.Reset()
	return iter
}

type VelocityWriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Velocity]
}

func (iter *VelocityWriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *VelocityWriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *VelocityWriteIterator) NextVelocity() (Entity, *Velocity, func()) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--
	done := func() {
		iter.w.patch.VelocityComponents[e.val] = c.ToPB()
	}

	return e, c, done
}

func (iter *VelocityWriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) VelocityWriteIter() *VelocityWriteIterator {
	iter := &VelocityWriteIterator{
		w:     w,
		store: w.velocitiesStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) VelocityEntities() []Entity {
	return w.velocitiesStore.entities()
}

func (w *World) SetVelocitySortFn(lessThan func(a, b Entity) bool) {
	w.velocitiesStore.LessThan = lessThan
}

func (w *World) SortVelocities() {
	w.velocitiesStore.Sort()
}

func (w *World) ApplyVelocityPatch(e Entity, patch *ecspb.VelocityComponent) Entity {
	c := Velocity{
		X: patch.X,
		Y: patch.Y,
		Z: patch.Z,
	}
	e.w.velocitiesStore.Set(c, e)
	return e
}

func (c Velocity) ToPB() *ecspb.VelocityComponent {
	pb := &ecspb.VelocityComponent{
		X: c.X,
		Y: c.Y,
		Z: c.Z,
	}
	return pb
}
