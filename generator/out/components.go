package out

import (
	"fmt"
	"reflect"
	"unsafe"
)

func registerComponent[T any](w *World, resetExample T, name string) ID {
	componentID := w.CreateEntityWithoutName()

	n := fmt.Sprintf("<%s>", reflect.TypeOf(resetExample).Name())
	if len(name) > 0 {
		n = name + n
	}

	// create a column for the component
	metadata := &componentMetadata{
		id:          componentID,
		name:        n,
		elementSize: unsafe.Sizeof(resetExample),
	}
	w.componentMetadatas[componentID] = metadata

	if metadata.elementSize == 0 {
		panic("element size is zero for a component")
	}

	buf := make([]byte, metadata.elementSize)
	setColumnData(int(metadata.elementSize), 0, buf, resetExample)
	metadata.resetBuf = buf

	return componentID
}

func componentDataFromEntity[T any](w *World, cID, eID ID, data *T) {
	record, ok := w.entityRecords[eID]
	if !ok {
		// panic("entity not found in any archetype")
        return
	}

	if record.row < 0 {
		// panic("entity does not have data")
        return
	}

	// first check if the archetype has the component
	componentColumnIndicies := w.archetypeComponentColumnIndicies[record.archetype.hash]
	column := componentColumnIndicies[cID]
	componentData(w, record.archetype, column, record.row, data)
}

func componentData[T any](w *World, archeType *Archetype, colIdx, rowIdx int, data *T) {
	col := archeType.dataColumns[colIdx]
	offset := uintptr(rowIdx) * col.metadata.elementSize
	start := unsafe.Pointer(unsafe.SliceData(col.data))
	rowPosition := unsafe.Add(start, offset)
	ptr := (*T)(rowPosition)
	*data = *ptr
}

func addComponentsTo(w *World, componentIDs *IDSet, entities ...ID) {
    componentIDs.Add(WildCardAllID)
    for _, entity := range entities {
		record := w.entityRecords[entity]

		archetype, changed := w.upsertArchetype(record.archetype, componentIDs)
		if !changed {
			return
		}
		// log.Printf("Add %s components to (%d,%d)", componentIDs.String(), entity.Source(), entity.Target())

		moveEntity(w, entity, record, archetype)

		componentIDs.Range(func(cID ID) {
			componentMetadata, isComponent := w.componentMetadatas[cID]
			if isComponent {
				resetBuf := componentMetadata.resetBuf
				colIndices := w.archetypeComponentColumnIndicies[record.archetype.hash]
				colIdx := colIndices[cID]
				colData := record.archetype.dataColumns[colIdx]
				colData.data = append(colData.data, resetBuf...)
				colData.count++
			}
		})
	}
}

func removeComponentFrom(w *World, componentIDs, entities *IDSet) {
	if componentIDs.Cardinality() == 0 {
		return
	}
	entities.Range(func(entity ID) {
		record, ok := w.entityRecords[entity]
		if !ok {
			panic("entity not found in any archetype")
		}

		targetArchetype := w.backtrackArchetype(record.archetype, componentIDs)
		moveEntity(w, entity, record, targetArchetype)
	})
}

func setComponentData[T any](w *World, cID ID, data T, entities ...ID) {
	componentSet := NewIDSet(cID)
	addComponentsTo(w, componentSet, entities...)

	// source, target,wasPair := cID.SplitPair()
	for _, eID := range entities {
		record, ok := w.entityRecords[eID]
		if !ok {
			panic("entity not found in any archetype")
		}

		// first check if the archetype has the component
		componentColumnIdx, ok := w.archetypeComponentColumnIndicies[record.archetype.hash]
		if !ok {
			panic("entity does not have component")
		}
		column, ok := componentColumnIdx[cID]
		if !ok {
			panic("entity does not have component")
		}

		if len(record.archetype.dataColumns) <= column {
			panic("entity does not have enough columns")
		}

		col := record.archetype.dataColumns[column]
		if !ok {
			panic("entity does not have column")
		}
		setColumnData(int(col.metadata.elementSize), record.row, col.data, data)
	}
}