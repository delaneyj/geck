package geck

import (
	geckpb "github.com/delaneyj/geck/pb/gen/geck/v1"
)

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

func (w *World) Query(terms *geckpb.Query_Terms) (iter *EntityIterator) {
	// queryCIDs := NewIDSet(componentIDs...)
	// componentCount := len(componentIDs)

	archetypes := []*Archetype{}
	archetypeComponentsColumns := [][]ID{}
	for archetypeID, components := range w.archetypeComponentColumnIndicies {
		if len(components) == 0 {
			continue
		}

		archetype, ok := w.archetypes[archetypeID]

		if !ok || len(archetype.entities) == 0 {
			continue
		}

		dataCols := make([]ID, len(archetype.dataColumns))

		// check that query terms are satisfied
		var checkValidity func(depth int, terms *geckpb.Query_Terms) bool
		checkValidity = func(depth int, terms *geckpb.Query_Terms) bool {
			isValid := false
		valid:
			for _, term := range terms.Terms {
				// check that the archetype has the component
				switch el := term.Element.(type) {
				case *geckpb.Query_Term_Id:
					elID := ID(el.Id)
					hasComponent := archetype.componentIDs.Contains(elID)

					shouldBreak := false
					if hasComponent {
						switch terms.Op {
						case geckpb.Query_AND:
							isValid = true
						case geckpb.Query_OR:
							isValid = true
							shouldBreak = true
						case geckpb.Query_NOT:
							isValid = false
							shouldBreak = true
						}
					} else {
						switch terms.Op {
						case geckpb.Query_AND:
							isValid = false
							shouldBreak = true
						case geckpb.Query_NOT:
							isValid = true
							shouldBreak = true
						}
					}

					if depth == 0 && isValid {
						for j, dc := range archetype.dataColumns {
							if dc.metadata.id == elID {
								dataCols[j] = elID
								break
							}
						}
					}

					if shouldBreak {
						break valid
					}

				case *geckpb.Query_Term_Terms:
					isValid = checkValidity(depth+1, el.Terms)
				}
			}

			return isValid
		}

		if !checkValidity(0, terms) {
			continue
		}

		archetypes = append(archetypes, archetype)
		archetypeComponentsColumns = append(archetypeComponentsColumns, dataCols)
	}

	if len(archetypes) == 0 {
		return w.finishedIter
	}

	iter = &EntityIterator{
		w:             w,
		archetypes:    archetypes,
		componentCols: archetypeComponentsColumns,
	}
	iter.Reset()
	return iter
}

func idToTerms(op geckpb.Query_Op, ids ...ID) *geckpb.Query_Terms {
	terms := make([]*geckpb.Query_Term, len(ids))
	for i, id := range ids {
		terms[i] = &geckpb.Query_Term{
			Element: &geckpb.Query_Term_Id{
				Id: uint64(id),
			},
		}
	}
	return &geckpb.Query_Terms{
		Op:    op,
		Terms: terms,
	}
}

func (w *World) QueryAnd(componentIDs ...ID) *EntityIterator {
	qt := idToTerms(geckpb.Query_AND, componentIDs...)
	return w.Query(qt)
}

func (w *World) QueryOr(componentIDs ...ID) *EntityIterator {
	qt := idToTerms(geckpb.Query_OR, componentIDs...)
	return w.Query(qt)
}

func (w *World) QueryNot(componentIDs ...ID) *EntityIterator {
	qt := idToTerms(geckpb.Query_NOT, componentIDs...)
	return w.Query(qt)
}
