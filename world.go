package geck

type archetypeToRowMap map[ID]int
type componentToArchetypeMap map[uint64]archetypeToRowMap
type entityRecord struct {
	archetype *Archetype
	row       int
}

type World struct {
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
	w.identfierID = RegisterComponent(w, "", "Identifier")
	w.nameID = w.CreateEntity()
	w.identifierNameID = w.CreatePair(w.identfierID, w.nameID)
	w.SetEntityName(w.identfierID, "Identifier")
	w.SetEntityName(w.nameID, "Name")

	// create internal tags
	w.internalID = w.CreateEntity("Internal")
	w.wildcardID = w.CreateEntity("Wildcard")
	w.childOfID = w.CreateEntity("ChildOf")
	w.instanceOfID = w.CreateEntity("InstanceOf")
	w.componentID = w.CreateEntity("Component")

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
	w.nextID = UserDefined

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

func (w *World) CreateEntity(name ...string) ID {
	id := w.nextID
	w.nextID++
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
