package ecs

import (
	ecspb "github.com/delaneyj/geck/cmd/example/ecs/pb/gen/ecs/v1"
)

type Faction Entity

func FactionFromEntity(c Entity) Faction {
	return Faction(c)
}

func (c Faction) ToEntity() Entity {
	return Entity(c)
}

func (w *World) ResetFaction() Entity {
	return w.EntityFromU32(0)
}

func (c Faction) FromEntity(e Entity) Faction {
	return Faction(e)
}

//#region Events
//#endregion

func (e Entity) HasFaction() bool {
	return e.w.factionsStore.Has(e)
}

func (e Entity) ReadFaction() (Entity, bool) {
	val, ok := e.w.factionsStore.Read(e)
	if !ok {
		return Entity{}, false
	}
	return Entity(val), true
}

func (e Entity) RemoveFaction() Entity {
	e.w.factionsStore.Remove(e)

	return e
}

func (e Entity) WritableFaction() (c *Faction, done func()) {
	var ok bool
	c, ok = e.w.factionsStore.Writeable(e)
	if !ok {
		return nil, nil
	}
	return c, func() {
		e.w.patch.FactionComponents[e.val] = c.ToPB()
	}
}

func (e Entity) SetFaction(other Entity) Entity {
	e.w.factionsStore.Set(Faction(other), e)

	e.w.patch.FactionComponents[e.val] = Faction(other).ToPB()
	return e
}

func (w *World) SetFactions(c Faction, entities ...Entity) {
	if len(entities) == 0 {
		panic("no entities provided, are you sure you didn't mean to call SetFactionResource?")
	}
	w.factionsStore.Set(c, entities...)
	w.patch.FactionComponents[w.resourceEntity.val] = c.ToPB()
}

func (w *World) RemoveFactions(entities ...Entity) {
	w.factionsStore.Remove(entities...)
	for _, entity := range entities {
		w.patch.FactionComponents[entity.val] = nil
	}
}

//#region Resources

// HasFactionResource checks if the world has a Faction}}
func (w *World) HasFactionResource() bool {
	return w.resourceEntity.HasFaction()
}

// FactionResource Retrieve the  resource from the world
func (w *World) FactionResource() (Entity, bool) {
	return w.resourceEntity.ReadFaction()
}

// SetFactionResource set the resource in the world
func (w *World) SetFactionResource(e Entity) Entity {
	w.resourceEntity.SetFaction(e)
	return w.resourceEntity
}

// RemoveFactionResource removes the resource from the world
func (w *World) RemoveFactionResource() Entity {
	w.resourceEntity.RemoveFaction()

	return w.resourceEntity
}

// WriteableFactionResource returns a writable reference to the resource
func (w *World) WriteableFactionResource() (c *Faction, done func()) {
	return w.resourceEntity.WritableFaction()
}

//#endregion

//#region Iterators

type FactionReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Faction]
}

func (iter *FactionReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *FactionReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *FactionReadIterator) NextFaction() (Entity, Faction) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *FactionReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) FactionReadIter() *FactionReadIterator {
	iter := &FactionReadIterator{
		w:     w,
		store: w.factionsStore,
	}
	iter.Reset()
	return iter
}

type FactionWriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Faction]
}

func (iter *FactionWriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *FactionWriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *FactionWriteIterator) NextFaction() (Entity, *Faction, func()) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--
	done := func() {
		iter.w.patch.FactionComponents[e.val] = c.ToPB()
	}

	return e, c, done
}

func (iter *FactionWriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) FactionWriteIter() *FactionWriteIterator {
	iter := &FactionWriteIterator{
		w:     w,
		store: w.factionsStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) FactionEntities() []Entity {
	return w.factionsStore.entities()
}

func (w *World) SetFactionSortFn(lessThan func(a, b Entity) bool) {
	w.factionsStore.LessThan = lessThan
}

func (w *World) SortFactions() {
	w.factionsStore.Sort()
}

func (w *World) ApplyFactionPatch(e Entity, patch *ecspb.FactionComponent) Entity {
	c := Faction(w.EntityFromU32(patch.Entity))
	e.w.factionsStore.Set(c, e)
	return e
}

func (c Faction) ToPB() *ecspb.FactionComponent {
	pb := &ecspb.FactionComponent{
		Entity: c.val,
	}
	return pb
}
