package out



type EntityIterator struct {
	w                    *World
	archetypes           []*Archetype
	componentCols        [][]ID
	currentArchetype     *Archetype
	currentColumns       []ID
	archetypeIdx, rowIdx int
	isDone               bool
}


func (eIter *EntityIterator) Next() ID {
	if eIter.isDone {
		panic("iterator is done")
	}

	entity := eIter.currentArchetype.entities[eIter.rowIdx]
	eIter.rowIdx++

	if eIter.rowIdx >= len(eIter.currentArchetype.entities) {
		eIter.archetypeIdx++
		if eIter.archetypeIdx >= len(eIter.archetypes) {
			eIter.isDone = true
		} else {
			eIter.updateToCurrentArchetype()
		}
	}

	return entity
}

func (eIter *EntityIterator) updateToCurrentArchetype() {
	eIter.currentArchetype = eIter.archetypes[eIter.archetypeIdx]
	eIter.currentColumns = eIter.componentCols[eIter.archetypeIdx]
	eIter.rowIdx = 0
}

func (eIter *EntityIterator) HasNext() bool {
	return !eIter.isDone
}

func (eIter *EntityIterator) Reset() {
	eIter.archetypeIdx = 0
	eIter.updateToCurrentArchetype()
	eIter.isDone = false
}