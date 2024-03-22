package {{.PackageName}}

type archetypeToRowMap map[ID]int
type componentToArchetypeMap map[uint64]archetypeToRowMap
type entityRecord struct {
	archetype *Archetype
	row       int
}

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
}

var(
	WildCardAllID = NewPair(WildcardID, WildcardID)
	IdentifierNameID = NewPair(IdentifierID, NameID)
)

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

	// bootstrap naming
	registerComponent(w, "", IdentifierName)
	nameID := w.CreateEntityWithoutName()
	if nameID != NameID {
		panic("name id mismatch")
	}
	w.SetEntityName(IdentifierID, IdentifierName)
	w.SetEntityName(NameID, NameName)

	// create internal tags
	{{ range .Bundles -}}
	{{ if .IsBuiltin -}}
		{{ range slice .Components 2 -}}
			{{.Name.Camel}}ID := w.CreateEntity({{.Name.Pascal}}Name)
			if {{.Name.Camel}}ID != {{.Name.Pascal}}ID {
				panic("id mismatch")
			}
		{{ end -}}
	{{ end -}}
	{{ end -}}

	// add "internal" tag to internal tags
	internalTags := []ID{
		{{ range .Bundles -}}
			{{if .IsBuiltin -}}
				{{ range .Components -}}
		{{.Name.Pascal}}ID,
				{{ end -}}
			{{ end -}}
		{{ end -}}
	}
	internalTagTags := NewIDSet(InternalID, WildCardAllID)
	addComponentsTo(w, internalTagTags, internalTags...)
	w.nextID = 1000

	{{ range .Bundles -}}
		{{if not .IsBuiltin -}}
			{{ range .Components -}}
	registerComponent(w, {{.Name.Pascal}}ResetValue, {{.Name.Pascal}}Name)
			{{ end -}}
		{{ end -}}
	{{ end -}}


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
		setComponentData(w, IdentifierNameID, name[0], id)

		// target, source, _ := id.SplitPair()
		// log.Printf("created source %d->%d", source, target)
	}
	//  else {
	// log.Printf("created entity %d", id)
	// }
}

func (w *World) nextAvailableID() (id ID) {
	if len(w.availableIDs) > 0 {
		lastIdx := len(w.availableIDs) - 1
		id = w.availableIDs[lastIdx].UpdateGeneration()
		w.availableIDs = w.availableIDs[:lastIdx]
	} else {
		id = w.nextID
		w.nextID++
	}
	return id
}

func (w *World) CreateEntity(name string) (id ID) {
	id = w.nextAvailableID()
	w.createEntity(id, name)
	return id
}

func (w *World) CreateEntityWithoutName() (id ID) {
	id = w.nextAvailableID()
	w.createEntity(id)
	return id
}

func (w *World) CreatePair(source, target ID, name string) ID {
	pair := NewPair(source, target)
	w.createEntity(pair, name)
	return pair
}

func (w *World) CreatePairWithoutName(source, target ID) ID {
	pair := NewPair(source, target)
	w.createEntity(pair)
	return pair
}

func (w *World) CreateEntitiesWith(count int, cIDs *IDSet) []ID {
	// search for an archetype that matches the components
	entities := make([]ID, count)
	for i := 0; i < count; i++ {
		e := w.CreateEntityWithoutName()
		entities[i] = e
	}

	if cIDs != nil {
		addComponentsTo(w, cIDs, entities...)
	}

	return entities
}

func (w *World) SetEntityName(id ID, name string) {
	setComponentData(w, IdentifierNameID, name, id)
	// log.Printf("set entity %d name to %s", id, name)
}

func (w *World) EntityName(id ID) (name string) {
	componentDataFromEntity(w, IdentifierNameID, id, &name)
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
	setComponentData(w, IdentifierNameID, name, id)
}

func (w *World) Name(id ID) (name string) {
	componentDataFromEntity(w, IdentifierNameID, id, &name)
	return name
}

func (w *World) EntitiesFromNames(names ...string) (ids []ID) {
	panic("not implemented")
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

func (w *World) MarshalAll() () {
	{{ range .Bundles -}}
		{{ range .Components -}}
			w.MarshalAll{{.Name.Pascal}}()
		{{ end -}}
	{{ end -}}
}