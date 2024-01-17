package geck

import (
	"fmt"
	"reflect"
	"unsafe"
)

func RegisterComponent[T any](w *World, resetExample T, name string) ID {
	componentID := w.CreateEntity()

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
	metadata.resetExample = buf

	return componentID
}

func ComponentData[T any](w *World, cID, eID ID, data *T) {
	record, ok := w.entityRecords[eID]
	if !ok {
		panic("entity not found in any archetype")
	}

	if record.row < 0 {
		panic("entity does not have data")
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

func SetComponentData[T any](w *World, cID ID, data T, entities *IDSet) {
	componentSet := NewIDSet(cID)
	AddComponentsTo(w, componentSet, entities)

	// source, target,wasPair := cID.SplitPair()
	entities.Range(func(eID ID) {
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
	})
}

func AddComponentsTo(w *World, componentIDs, entities *IDSet) {
	entities.Range(func(entity ID) {
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
				resetExample := componentMetadata.resetExample
				colIndices := w.archetypeComponentColumnIndicies[record.archetype.hash]
				colIdx := colIndices[cID]
				colData := record.archetype.dataColumns[colIdx]
				colData.data = append(colData.data, resetExample...)
				colData.count++
			}
		})

	})
}

func RemoveComponentFrom(w *World, componentIDs, entities *IDSet) {
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

func AddPairs[T any](w *World, entities *IDSet, source, target ID, data ...T) {
	pair := NewPair(source, target)
	sourceWildcard := NewPair(source, w.wildcardID)
	targetWildcard := NewPair(w.wildcardID, target)
	bothWildcard := NewPair(w.wildcardID, w.wildcardID)
	set := NewIDSet(pair, sourceWildcard, targetWildcard, bothWildcard)

	AddComponentsTo(w, set, entities)
	if len(data) > 0 {
		SetComponentData(w, pair, data[0], entities)
	}
}
