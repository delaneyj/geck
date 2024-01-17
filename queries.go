package geck

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

func Data[T any](eIter *EntityIterator, data *T, colIdx int) {
	// colIdx := eIter.currentComponentToColumn[componentID]
	componentData(eIter.w, eIter.currentArchetype, colIdx, eIter.rowIdx, data)
}

func Query(w *World, componentIDs ...ID) (iter *EntityIterator) {

	if len(componentIDs) == 0 {
		return w.finishedIter
	}

	queryCIDs := NewIDSet(componentIDs...)
	componentCount := len(componentIDs)

	archetypes := []*Archetype{}
	archetypeComponentColumns := [][]ID{}
	for archetypeID, components := range w.archetypeComponentColumnIndicies {
		if len(components) == 0 {
			continue
		}

		archetype := w.archetypes[archetypeID]
		sharedCount := queryCIDs.AndCardinality(archetype.componentIDs)
		if sharedCount != componentCount {
			continue
		}
		if len(archetype.entities) == 0 {
			continue
		}
		archetypes = append(archetypes, archetype)

		dataCols := make([]ID, len(archetype.dataColumns))
		for i, col := range archetype.dataColumns {
			for _, cID := range componentIDs {
				if col.metadata.id == cID {
					dataCols[i] = cID
					break
				}
			}

		}
		archetypeComponentColumns = append(archetypeComponentColumns, dataCols)
	}

	if len(archetypes) == 0 {
		return w.finishedIter
	}

	iter = &EntityIterator{
		w:             w,
		archetypes:    archetypes,
		componentCols: archetypeComponentColumns,
	}
	iter.Reset()
	return iter
}
