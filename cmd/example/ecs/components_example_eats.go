package ecs

import (
	ecspb "github.com/delaneyj/geck/cmd/example/ecs/pb/gen/ecs/v1"
)

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

	e.w.patch.EatsComponents[e.val] = nil
	return e
}

func (e Entity) WritableEats() (c *Eats, done func()) {
	var ok bool
	c, ok = e.w.eatsStore.Writeable(e)
	if !ok {
		return nil, nil
	}
	return c, func() {
		e.w.patch.EatsComponents[e.val] = c.ToPB()
	}
}

func (e Entity) SetEats(other Eats) Entity {
	e.w.eatsStore.Set(other, e)

	e.w.patch.EatsComponents[e.val] = Eats(other).ToPB()
	return e
}

func (e Entity) SetEatsValues(
	entities0 []Entity,
	amounts1 []uint8,
) Entity {
	err := e.w.eatsStore.Set(Eats{
		Entities: entities0,
		Amounts:  amounts1,
	}, e)
	if err != nil {
		panic(err)
	}
	pb := &ecspb.EatsComponent{

		Entities: EntitiesToU32s(entities0...),
		Amounts:  make([]uint32, len(amounts1)),
	}
	for i, v := range amounts1 {
		pb.Amounts[i] = uint32(v)
	}
	e.w.patch.EatsComponents[e.w.resourceEntity.val] = pb
	return e
}

func (w *World) SetEats(c Eats, entities ...Entity) {
	if len(entities) == 0 {
		panic("no entities provided, are you sure you didn't mean to call SetEatsResource?")
	}
	w.eatsStore.Set(c, entities...)
	w.patch.EatsComponents[w.resourceEntity.val] = c.ToPB()
}

func (w *World) RemoveEats(entities ...Entity) {
	w.eatsStore.Remove(entities...)
	for _, entity := range entities {
		w.patch.EatsComponents[entity.val] = nil
	}
}

//#region Resources

// HasEatsResource checks if the world has a Eats}}
func (w *World) HasEatsResource() bool {
	return w.resourceEntity.HasEats()
}

// EatsResource Retrieve the  resource from the world
func (w *World) EatsResource() (Eats, bool) {
	return w.resourceEntity.ReadEats()
}

// SetEatsResource set the resource in the world
func (w *World) SetEatsResource(c Eats) Entity {
	w.resourceEntity.SetEats(c)
	return w.resourceEntity
}
func (w *World) SetEatsResourceValues(
	entities0 []Entity,
	amounts1 []uint8,
) Entity {
	w.resourceEntity.SetEatsValues(
		entities0,
		amounts1,
	)
	return w.resourceEntity
}

// RemoveEatsResource removes the resource from the world
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

func (iter *EatsWriteIterator) NextEats() (Entity, *Eats, func()) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--
	done := func() {
		iter.w.patch.EatsComponents[e.val] = c.ToPB()
	}

	return e, c, done
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

func (w *World) ApplyEatsPatch(e Entity, patch *ecspb.EatsComponent) Entity {
	c := Eats{
		Entities: w.EntitiesFromU32s(patch.Entities...),
		Amounts:  make([]uint8, len(patch.Amounts)),
	}
	for i, v := range patch.Amounts {
		c.Amounts[i] = uint8(v)
	}
	e.w.eatsStore.Set(c, e)
	return e
}

func (c Eats) ToPB() *ecspb.EatsComponent {
	pb := &ecspb.EatsComponent{
		Entities: EntitiesToU32s(c.Entities...),
		Amounts:  make([]uint32, len(c.Amounts)),
	}
	for i, v := range c.Entities {
		pb.Entities[i] = v.val
	}
	for i, v := range c.Amounts {
		pb.Amounts[i] = uint32(v)
	}
	return pb
}
