package ecs

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

func (e Entity) WritableFaction() (*Faction, bool) {
	return e.w.factionsStore.Writeable(e)
}

func (e Entity) SetFaction(other Entity) Entity {
	e.w.factionsStore.Set(Faction(other), e)

	return e
}

func (w *World) SetFactions(c Faction, entities ...Entity) {
	w.factionsStore.Set(c, entities...)
}

func (w *World) RemoveFactions(entities ...Entity) {
	w.factionsStore.Remove(entities...)
}

//#region Resources

// HasFaction checks if the world has a Faction}}
func (w *World) HasFactionResource() bool {
	return w.resourceEntity.HasFaction()
}

// Retrieve the Faction resource from the world
func (w *World) FactionResource() (Entity, bool) {
	return w.resourceEntity.ReadFaction()
}

// Set the Faction resource in the world
func (w *World) SetFactionResource(c Entity) Entity {
	w.resourceEntity.SetFaction(c)

	return w.resourceEntity
}

// Remove the Faction resource from the world
func (w *World) RemoveFactionResource() Entity {
	w.resourceEntity.RemoveFaction()

	return w.resourceEntity
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

func (iter *FactionWriteIterator) NextFaction() (Entity, *Faction) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--

	return e, c
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
