package geck

import (
	"encoding/binary"
	"slices"
	"unsafe"

	"github.com/zeebo/xxh3"
)

const PageSize = 4096

type archetypeEdge struct {
	add    *Archetype
	remove *Archetype
}

type componentMetadata struct {
	id           ID
	name         string
	resetExample []byte
	elementSize  uintptr
}

type componentColumn struct {
	metadata *componentMetadata
	id       ID
	index    int
	count    int
	data     []byte
}

type Archetype struct {
	hash         uint64
	depth        int
	componentIDs *IDSet
	canHaveData  bool
	dataColumns  []*componentColumn
	edges        map[ID]*archetypeEdge // component id to archetype
	entities     []ID
}

func NewArchetype(w *World, from *Archetype, componentID *ID) *Archetype {
	var (
		componentIDs *IDSet
		depth        int
		// name         string
	)

	if componentID == nil {
		componentIDs = NewIDSet()
		// 	// name = "Root"
	} else {
		// 	_, isComponent := w.componentMetadatas[*componentID]

		// 	if !isComponent {
		// 		// check if pair
		// 		source, target, wasPair := componentID.SplitPair()
		// 		if wasPair {
		// 			name = fmt.Sprintf("%d->%d", source, target)
		// 			_, isComponent = w.componentMetadatas[source]
		// 		}
		// 	}

		componentIDs = from.componentIDs.Clone().Add(*componentID)

		// 	componentNames := []string{}
		// 	componentIDs.Range(func(cID ID) {
		// 		name := ""
		// 		if isComponent {
		// 			name = w.componentMetadatas[cID].name
		// 		} else {

		// 			ComponentData(w, w.identifierNameID, cID, &name)
		// 			if name == "" {
		// 				name = fmt.Sprintf("tag<%d>", cID)
		// 			}
		// 		}
		// 		componentNames = append(componentNames, name)
		// 	})
		// 	name = strings.Join(componentNames, ",")
		// }
	}

	// calculate archetype id
	h := xxh3.New()
	componentIDs.Range(func(cID ID) {
		binary.Write(h, binary.LittleEndian, cID)
	})

	if from != nil {
		depth = from.depth + 1
	}

	archetype := &Archetype{
		hash:         h.Sum64(),
		depth:        depth,
		componentIDs: componentIDs,
		edges:        map[ID]*archetypeEdge{},
	}
	if from != nil {
		archetype.edges[*componentID] = &archetypeEdge{
			remove: from,
		}
	}

	componentColumnIndicies := map[ID]int{}

	i := 0
	componentIDs.Range(func(cID ID) {
		var (
			metadata    *componentMetadata
			isComponent bool
		)
		source, target, wasPair := cID.SplitPair()
		if wasPair {
			metadata, isComponent = w.componentMetadatas[source]
		} else {
			metadata, isComponent = w.componentMetadatas[target]
		}

		componentColIdx := -1

		if isComponent {
			archetype.dataColumns = append(archetype.dataColumns, &componentColumn{
				metadata: metadata,
				id:       cID,
				index:    i,
				count:    0,
				data:     make([]byte, 0, PageSize),
			})
			archetype.canHaveData = true
			componentColIdx = i
			i++
		}

		componentColumnIndicies[cID] = componentColIdx
	})
	w.archetypeComponentColumnIndicies[archetype.hash] = componentColumnIndicies
	w.archetypes[archetype.hash] = archetype
	return archetype
}

func moveEntity(w *World, entity ID, r *entityRecord, newArchetype *Archetype) {
	// log.Printf("move entity %d from archetype with %s to %s", entity, r.archetype.componentIDs, newArchetype.componentIDs)

	// TODO remove this
	if !slices.Contains(r.archetype.entities, entity) {
		panic("entity not found in archetype")
	}
	if r.archetype == newArchetype {
		panic("entity already in archetype")
	}

	// remove entity from old archetype
	prevArchetypeEntityCount := len(r.archetype.entities)
	lastEntityIdx := prevArchetypeEntityCount - 1
	if prevArchetypeEntityCount > 0 && r.row >= 0 {
		if lastEntityIdx > 0 && r.row != lastEntityIdx {
			// log.Printf("move last entity to row %d. %v", r.row, r.archetype.entities)
			lastEntityRecord := w.entityRecords[r.archetype.entities[lastEntityIdx]]
			r.archetype.entities[r.row] = r.archetype.entities[lastEntityIdx]
			lastEntityRecord.row = r.row
		}
		r.archetype.entities = r.archetype.entities[:lastEntityIdx]
		// log.Printf("move last entity to current row. %v", r.archetype.entities)
	}

	// add entity to new archetype
	newRowIdx := len(newArchetype.entities)
	newArchetype.entities = append(newArchetype.entities, entity)
	for _, col := range newArchetype.dataColumns {
		// append row to all data columns
		col.data = append(col.data, col.metadata.resetExample...)
		col.count++
	}

	// Move old data to new data
	if r.row >= 0 && r.archetype.canHaveData {
		oldRowPtr := uintptr(r.row)

		newColumns := w.archetypeComponentColumnIndicies[newArchetype.hash]

		// Has data, copy the data from the old archetype to the new archetype
		for _, oldCol := range r.archetype.dataColumns {
			if oldCol.count == 0 {
				continue
			}

			colID := oldCol.id
			size := oldCol.metadata.elementSize

			newColIdx, ok := newColumns[colID]
			if !ok {
				// new column does not exist, better be from removing a component
				continue
			}
			newCol := newArchetype.dataColumns[newColIdx]
			newRowPtr := uintptr(newRowIdx)

			if oldCol.metadata.elementSize != newCol.metadata.elementSize {
				panic("element size mismatch")
			}

			srcStart := oldRowPtr * size
			srcEnd := srcStart + size
			dstStart := newRowPtr * size
			dstEnd := dstStart + size

			if dstEnd > uintptr(len(newCol.data)) {
				panic("out of bounds")
			}

			dst := newCol.data[dstStart:dstEnd]
			src := oldCol.data[srcStart:srcEnd]
			copy(dst, src)

			// move last old row to the current row and update the index
			oldLastRowIdx := oldCol.count - 1
			if oldLastRowIdx == 0 {
				return
			}

			oldLastRowPtr := uintptr(oldLastRowIdx)
			lastSrcStart := oldLastRowPtr * oldCol.metadata.elementSize
			lastSrcEnd := lastSrcStart + oldCol.metadata.elementSize
			lastDstStart := oldRowPtr * oldCol.metadata.elementSize
			lastDstEnd := lastDstStart + oldCol.metadata.elementSize

			lastDst := oldCol.data[lastDstStart:lastDstEnd]
			lastSrc := oldCol.data[lastSrcStart:lastSrcEnd]
			copy(lastDst, lastSrc)

			oldCol.data = oldCol.data[:lastSrcStart]
			oldCol.count--
		}
	}

	r.archetype = newArchetype
	r.row = newRowIdx
}

func setColumnData[T any](elementSize int, rowIdx int, data []byte, datum T) {
	start := rowIdx * elementSize
	end := start + elementSize
	if end > len(data) {
		panic("out of bounds")
	}

	// copy the data into the pointer using unsafe
	ptr := unsafe.Pointer(&data[start])
	tPtr := (*T)(ptr)
	*tPtr = datum
}
