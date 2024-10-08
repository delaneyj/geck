package ecs

type PositionVelocitySet struct {
	lastIdx int

	// owned components
	ownedVelocitiesStore *SparseSet[Velocity]
	ownedPositionsStore  *SparseSet[Position]
}

func NewPositionVelocitySet(w *World) *PositionVelocitySet {
	set := &PositionVelocitySet{
		lastIdx: -1,

		ownedVelocitiesStore: w.velocitiesStore,
		ownedPositionsStore:  w.positionsStore,
	}
	return set
}

func (set *PositionVelocitySet) PossibleUpdate(entities ...Entity) {
	for _, e := range entities {
		hasAllOwned := true

		if !set.ownedVelocitiesStore.Has(e) {
			hasAllOwned = false
			break
		}

		if !set.ownedPositionsStore.Has(e) {
			hasAllOwned = false
			break
		}

		sparseIdx := e.Index()

		if hasAllOwned {
			// swap with next after last
			set.lastIdx++

			wasSwapped := false

			if set.ownedVelocitiesStore.ownedSetSwap(set.lastIdx, sparseIdx, false) {
				wasSwapped = true
			}

			if set.ownedPositionsStore.ownedSetSwap(set.lastIdx, sparseIdx, false) {
				wasSwapped = true
			}

			if !wasSwapped {
				set.lastIdx--
			}
		} else {
			// swap with last
			wasSwapped := false

			if set.ownedVelocitiesStore.ownedSetSwap(set.lastIdx, sparseIdx, true) {
				wasSwapped = true
			}

			if set.ownedPositionsStore.ownedSetSwap(set.lastIdx, sparseIdx, true) {
				wasSwapped = true
			}

			if wasSwapped {
				set.lastIdx--
			}
		}

		// do something with
		// hasAllBorrowed := true

	}
}

func (set *PositionVelocitySet) Len() int {
	return set.lastIdx + 1
}

func (set *PositionVelocitySet) All() []Entity {
	entities := make([]Entity, 0, set.Len())
	iter := set.NewIterator()
	for iter.HasNext() {
		e, _, _ := iter.Next()
		entities = append(entities, e)
	}
	return entities
}

func (set *PositionVelocitySet) NewIterator() *PositionVelocitySetIter {
	iter := &PositionVelocitySetIter{set: set}
	iter.Reset()
	return iter
}

type PositionVelocitySetIter struct {
	set     *PositionVelocitySet
	current int
}

func (iter *PositionVelocitySetIter) Reset() {
	iter.current = iter.set.lastIdx
}

func (iter *PositionVelocitySetIter) HasNext() bool {
	return iter.current >= 0
}

func (iter *PositionVelocitySetIter) Next() (
	Entity,
	Velocity,
	*Position,
) {
	e := iter.set.ownedVelocitiesStore.dense[iter.current]
	comp0 := iter.set.ownedVelocitiesStore.components[iter.current]
	comp1 := &iter.set.ownedPositionsStore.components[iter.current]
	iter.current--
	return e, comp0, comp1
}
