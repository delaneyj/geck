package ecs

import (
	ecspb "github.com/delaneyj/geck/cmd/example/ecs/pb/gen/ecs/v1"
)

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

	e.w.patch.RotationComponents[e.val] = nil
	return e
}

func (e Entity) WritableRotation() (*Rotation, bool) {
	return e.w.rotationsStore.Writeable(e)
}

func (e Entity) SetRotation(other Rotation) Entity {
	e.w.rotationsStore.Set(other, e)

	e.w.patch.RotationComponents[e.w.resourceEntity.val] = Rotation(other).ToPB()
	return e
}

func (e Entity) SetRotationValues(
	x0 float32,
	y1 float32,
	z2 float32,
	w3 float32,
) Entity {
	err := e.w.rotationsStore.Set(Rotation{
		X: x0,
		Y: y1,
		Z: z2,
		W: w3,
	}, e)
	if err != nil {
		panic(err)
	}
	pb := &ecspb.RotationComponent{
		X: x0,
		Y: y1,
		Z: z2,
		W: w3,
	}
	e.w.patch.RotationComponents[e.w.resourceEntity.val] = pb
	return e
}

func (w *World) SetRotations(c Rotation, entities ...Entity) {
	if len(entities) == 0 {
		panic("no entities provided, are you sure you didn't mean to call SetRotationResource?")
	}
	w.rotationsStore.Set(c, entities...)
	w.patch.RotationComponents[w.resourceEntity.val] = c.ToPB()
}

func (w *World) RemoveRotations(entities ...Entity) {
	w.rotationsStore.Remove(entities...)
	for _, entity := range entities {
		w.patch.RotationComponents[entity.val] = nil
	}
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
func (w *World) SetRotationResourceValues(
	x0 float32,
	y1 float32,
	z2 float32,
	w3 float32,
) Entity {
	w.resourceEntity.SetRotationValues(
		x0,
		y1,
		z2,
		w3,
	)
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

func (w *World) ApplyRotationPatch(e Entity, patch *ecspb.RotationComponent) Entity {
	c := Rotation{
		X: patch.X,
		Y: patch.Y,
		Z: patch.Z,
		W: patch.W,
	}
	e.w.rotationsStore.Set(c, e)
	return e
}

func (c Rotation) ToPB() *ecspb.RotationComponent {
	pb := &ecspb.RotationComponent{
		X: c.X,
		Y: c.Y,
		Z: c.Z,
		W: c.W,
	}
	return pb
}
