package geck

import (
	"math"

	geckpb "github.com/delaneyj/geck/pb/gen/geck/v1"
)

type archetypeToRowMap map[ID]int
type componentToArchetypeMap map[uint64]archetypeToRowMap
type entityRecord struct {
	archetype *Archetype
	row       int
}

const (
	IdentifierName = "Identifier"
	NameName       = "Name"
	InternalName   = "Internal"
	WildcardName   = "Wildcard"
	ChildOfName    = "ChildOf"
	InstanceOfName = "InstanceOf"
	ComponentName  = "Component"
)

type World struct {
	availableIDs                     []ID
	nextID                           ID
	emptyArchetype                   *Archetype
	finishedIter                     *EntityIterator
	componentMetadatas               map[ID]*componentMetadata
	entityRecords                    map[ID]*entityRecord
	archetypes                       map[uint64]*Archetype   // hash to archetype
	archetypeComponentColumnIndicies componentToArchetypeMap // component to archetype to column
	archetypeSet                     *IDSet

	// internal tags
	internalID   ID
	wildcardID   ID
	childOfID    ID
	instanceOfID ID
	componentID  ID
	nameID       ID

	// internal components
	identfierID ID

	// internal pairs
	wildcardAllID    ID
	identifierNameID ID
}

func NewWorld() *World {
	w := &World{
		componentMetadatas:               map[ID]*componentMetadata{},
		entityRecords:                    map[ID]*entityRecord{},
		archetypes:                       map[uint64]*Archetype{},   // hash to archetype
		archetypeComponentColumnIndicies: componentToArchetypeMap{}, // find an archetype for a list of components
		archetypeSet:                     NewIDSet(),                // component ids used to create archetypes
		nextID:                           1,
	}
	w.emptyArchetype = NewArchetype(w, nil, nil)
	w.finishedIter = &EntityIterator{
		w:      w,
		isDone: true,
	}

	// bootstrar naming
	w.identfierID = RegisterComponent(w, "", IdentifierName)
	w.nameID = w.CreateEntity()
	w.identifierNameID = w.CreatePair(w.identfierID, w.nameID)
	w.SetEntityName(w.identfierID, IdentifierName)
	w.SetEntityName(w.nameID, NameName)

	// create internal tags
	w.internalID = w.CreateEntity(InternalName)
	w.wildcardID = w.CreateEntity(WildcardName)
	w.childOfID = w.CreateEntity(ChildOfName)
	w.instanceOfID = w.CreateEntity(InstanceOfName)
	w.componentID = w.CreateEntity(ComponentName)

	// create internal pairs
	w.wildcardAllID = NewPair(w.wildcardID, w.wildcardID)

	// add internal tags to internal tags
	internalTags := NewIDSet(
		w.internalID,
		w.wildcardID,
		w.childOfID,
		w.instanceOfID,
		w.componentID,
		w.nameID,
	)
	internalTagTags := NewIDSet(w.internalID, w.wildcardAllID)
	// log.Printf("Add %s tags to entities %s", internalTagTags, internalTags)
	AddComponentsTo(w, internalTagTags, internalTags)
	w.nextID = 1000

	// add internal components to internal tags
	internalComponents := NewIDSet(w.identfierID)
	internalComponentTags := internalTagTags.Clone().Add(w.componentID)
	AddComponentsTo(w, internalComponentTags, internalComponents)

	return w
}

func (w *World) Reset() {
	w.nextID = 0
	w.componentMetadatas = make(map[ID]*componentMetadata, len(w.componentMetadatas))
	w.entityRecords = make(map[ID]*entityRecord, len(w.entityRecords))
	w.archetypes = make(map[uint64]*Archetype, len(w.archetypes))
	w.archetypeComponentColumnIndicies = make(componentToArchetypeMap, len(w.archetypeComponentColumnIndicies))
	w.archetypeSet.Clear()
}

func (w *World) EntityCount() int {
	return len(w.entityRecords)
}

func (w *World) ArchetypeCount() int {
	return len(w.archetypes)
}

func (w *World) HasComponents(entities, components *IDSet) (allExist bool) {
	componentsCount := components.Cardinality()
	entities.ConditionalRange(func(entity ID, i int) bool {
		record, ok := w.entityRecords[entity]
		if !ok {
			panic("entity not found in any archetype")
		}
		allExist = record.archetype.componentIDs.AndCardinality(components) == componentsCount
		return allExist
	})
	return allExist
}

func (w *World) createEntity(id ID, name ...string) {
	w.entityRecords[id] = &entityRecord{
		archetype: w.emptyArchetype,
		row:       -1,
	}
	w.emptyArchetype.entities = append(w.emptyArchetype.entities, id)

	if len(name) > 0 {
		idSet := NewIDSet(id)
		SetComponentData(w, w.identifierNameID, name[0], idSet)

		// target, source, _ := id.SplitPair()
		// log.Printf("created source %d->%d", source, target)
	}
	//  else {
	// log.Printf("created entity %d", id)
	// }
}

func (w *World) CreateEntity(name ...string) (id ID) {
	if len(w.availableIDs) > 0 {
		lastIdx := len(w.availableIDs) - 1
		id = w.availableIDs[lastIdx].UpdateGeneration()
		w.availableIDs = w.availableIDs[:lastIdx]
	} else {
		id = w.nextID
		w.nextID++
	}
	w.createEntity(id, name...)
	return id
}

func (w *World) CreatePair(source, target ID, name ...string) ID {
	pair := NewPair(source, target)
	w.createEntity(pair, name...)
	return pair
}

func (w *World) CreateEntitiesWith(count int, cIDs *IDSet) *IDSet {
	// search for an archetype that matches the components
	entities := NewIDSet()
	for i := 0; i < count; i++ {
		e := w.CreateEntity()
		entities.Add(e)
	}

	if cIDs != nil {
		AddComponentsTo(w, cIDs, entities)
	}

	return entities
}

func (w *World) SetEntityName(id ID, name string) {
	idSet := NewIDSet(id)
	SetComponentData(w, w.identifierNameID, name, idSet)
	// log.Printf("set entity %d name to %s", id, name)
}

func (w *World) EntityName(id ID) (name string) {
	ComponentData(w, w.identifierNameID, id, &name)
	return name
}

func (w *World) EntityFromName(name string) ID {
	entities := w.EntitiesFromNames(name)
	if len(entities) == 0 {
		panic("entity not found")
	}
	return entities[0]
}

func (w *World) UpsertEntityFromName(name string) ID {
	entities := w.EntitiesFromNames(name)
	if len(entities) > 0 {
		return entities[0]
	}
	return w.CreateEntity(name)
}

func (w *World) DeleteEntity(id ID) {
	record, ok := w.entityRecords[id]
	if !ok {
		panic("entity not found")
	}

	// remove from archetype
	archetype := record.archetype
	row := record.row
	archetype.entities[row] = archetype.entities[len(archetype.entities)-1]
	archetype.entities = archetype.entities[:len(archetype.entities)-1]

	// remove from world
	delete(w.entityRecords, id)

	w.availableIDs = append(w.availableIDs, id)
}

func (w *World) SetName(id ID, name string) {
	SetComponentData(w, w.identifierNameID, name, NewIDSet(id))
}

func (w *World) Name(id ID) (name string) {
	ComponentData(w, w.identifierNameID, id, &name)
	return name
}

func (w *World) EntitiesFromNames(names ...string) (ids []ID) {
	q := w.Query(&geckpb.Query_Terms{
		Terms: []*geckpb.Query_Term{
			{
				Element: &geckpb.Query_Term_Id{
					Id: uint64(w.identifierNameID),
				},
			},
		},
	})
	namesMap := map[string]uint64{}
	for q.HasNext() {
		id := q.Next()
		name := w.Name(id)
		namesMap[name] = uint64(id)
	}
	for _, name := range names {
		ids = append(ids, ID(namesMap[name]))
	}
	return ids
}

func (w *World) ToPB() (wd *geckpb.WorldDefinition, err error) {
	wd = &geckpb.WorldDefinition{
		AvailableId:                       []uint64{},
		NextId:                            uint64(w.nextID),
		ComponentMetadata:                 map[uint64]*geckpb.ComponentMetadataDefinition{},
		Archetypes:                        map[uint64]*geckpb.ArchetypeDefinition{},
		ArchetypeComponentComlumnIndicies: map[uint64]*geckpb.ArchetypeToRowMap{},
	}

	for _, id := range w.availableIDs {
		wd.AvailableId = append(wd.AvailableId, uint64(id))
	}

	for id, metadata := range w.componentMetadatas {
		wd.ComponentMetadata[id.U64()] = &geckpb.ComponentMetadataDefinition{
			Id:           uint64(metadata.id),
			Name:         metadata.name,
			ResetExample: metadata.resetExample,
			ElementSize:  uint32(metadata.elementSize),
		}
	}

	for hash, archetype := range w.archetypes {
		wdArchetype := &geckpb.ArchetypeDefinition{
			Hash:         hash,
			Depth:        uint32(archetype.depth),
			ComponentIds: archetype.componentIDs.ToUint64s(),
			DataColumns:  []*geckpb.ComponentColumnDefinition{},
			Edges:        map[uint64]*geckpb.ArchetypeDefinition_Edge{},
			Entities:     []uint64{},
		}

		for i, dc := range archetype.dataColumns {
			ccd := &geckpb.ComponentColumnDefinition{
				ComponentId:    dc.id.U64(),
				ArchetypeIndex: uint32(i),
				Count:          uint32(dc.count),
				Data:           dc.data,
			}
			wdArchetype.DataColumns = append(wdArchetype.DataColumns, ccd)
		}

		for cID, edge := range archetype.edges {
			wdEdge := &geckpb.ArchetypeDefinition_Edge{
				AddId:    math.MaxUint64,
				RemoveId: math.MaxUint64,
			}
			if edge.add != nil {
				wdEdge.AddId = edge.add.hash
			}
			if edge.remove != nil {
				wdEdge.RemoveId = edge.remove.hash
			}
			wdArchetype.Edges[cID.U64()] = wdEdge
		}

		for _, id := range archetype.entities {
			wdArchetype.Entities = append(wdArchetype.Entities, id.U64())
		}

		wd.Archetypes[hash] = wdArchetype
	}

	for archetypeHash, archetypeToRowMap := range w.archetypeComponentColumnIndicies {
		wdArchetypeToRowMap := &geckpb.ArchetypeToRowMap{
			Value: map[uint64]uint32{},
		}

		for cID, row := range archetypeToRowMap {
			wdArchetypeToRowMap.Value[cID.U64()] = uint32(row)
		}

		wd.ArchetypeComponentComlumnIndicies[archetypeHash] = wdArchetypeToRowMap
	}

	return wd, nil
}

func (w *World) FromPB(wd *geckpb.WorldDefinition) (err error) {
	w.Reset()
	w.availableIDs = make([]ID, len(wd.AvailableId))
	for i, id := range wd.AvailableId {
		w.availableIDs[i] = ID(id)
	}

	w.nextID = ID(wd.NextId)

	w.emptyArchetype = NewArchetype(w, nil, nil)
	w.finishedIter = &EntityIterator{
		w:      w,
		isDone: true,
	}

	for cidRaw, cmd := range wd.ComponentMetadata {
		cID := ID(cidRaw)
		w.componentMetadatas[cID] = &componentMetadata{
			id:           ID(cmd.Id),
			name:         cmd.Name,
			resetExample: cmd.ResetExample,
			elementSize:  uintptr(cmd.ElementSize),
		}
	}

	for hash, wdArchetype := range wd.Archetypes {
		archetype := &Archetype{
			hash:         hash,
			depth:        int(wdArchetype.Depth),
			componentIDs: NewIDSetFromUint64s(wdArchetype.ComponentIds...),
			dataColumns:  []*componentColumn{},
			edges:        map[ID]*archetypeEdge{},
			entities:     []ID{},
		}

		for _, wdColumn := range wdArchetype.DataColumns {
			archetype.dataColumns = append(archetype.dataColumns, &componentColumn{
				metadata: w.componentMetadatas[ID(wdColumn.ComponentId)],
				id:       ID(wdColumn.ComponentId),
				index:    int(wdColumn.ArchetypeIndex),
				count:    int(wdColumn.Count),
				data:     wdColumn.Data,
			})
		}

		for cIDRaw, wdEdge := range wdArchetype.Edges {
			cID := ID(cIDRaw)
			archetype.edges[cID] = &archetypeEdge{
				add:    w.archetypes[wdEdge.AddId],
				remove: w.archetypes[wdEdge.RemoveId],
			}
		}

		for _, idRaw := range wdArchetype.Entities {
			archetype.entities = append(archetype.entities, ID(idRaw))
		}

		w.archetypes[hash] = archetype
	}

	for archetypeHash, wdArchetypeToRowMap := range wd.ArchetypeComponentComlumnIndicies {
		archetypeToRowMap := archetypeToRowMap{}
		for cIDRaw, row := range wdArchetypeToRowMap.Value {
			archetypeToRowMap[ID(cIDRaw)] = int(row)
		}
		w.archetypeComponentColumnIndicies[archetypeHash] = archetypeToRowMap
	}

	ids := w.EntitiesFromNames(
		"Internal",
		"Wildcard",
		"ChildOf",
		"InstanceOf",
		"Component",
		"Name",
		"Identifier",
	)
	w.internalID = ids[0]
	w.wildcardID = ids[1]
	w.childOfID = ids[2]
	w.instanceOfID = ids[3]
	w.componentID = ids[4]
	w.nameID = ids[5]
	w.identfierID = ids[6]

	w.wildcardAllID = NewPair(w.wildcardID, w.wildcardID)
	w.identifierNameID = NewPair(w.identfierID, w.nameID)

	return nil
}

func (w *World) upsertArchetype(from *Archetype, cIDs *IDSet) (current *Archetype, changed bool) {
	current = from
	cIDs.Range(func(cID ID) {
		if current.componentIDs.Contains(cID) {
			return
		}

		edge, ok := current.edges[cID]
		if !ok {
			edge = &archetypeEdge{}
			current.edges[cID] = edge
		}

		if edge.add == nil {
			edge.add = NewArchetype(w, current, &cID)
		}

		current = edge.add
		changed = true
	})

	return current, changed
}

func (w *World) backtrackArchetype(from *Archetype, removedComponents *IDSet) (current *Archetype) {
	current = from
	componentsLeft := removedComponents.Clone()

	for componentsLeft.Cardinality() > 0 {
		cID := componentsLeft.Pop()
		edge, ok := current.edges[cID]
		if !ok || edge.remove == nil {
			panic("edge not found")
		}
		current = edge.remove
	}

	return current
}
