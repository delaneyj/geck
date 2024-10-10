package ecs

func (w *World) UpsertPlanetTag(entities ...Entity) (anyUpdated bool) {
	for _, e := range entities {
		if _, updated := w.planetTags.Upsert(e, empty{}); updated {
			anyUpdated = true
		}
	}

	return anyUpdated
}

func (w *World) RemovePlanetTag(entities ...Entity) (anyRemoved bool) {
	for _, e := range entities {
		if removed := w.planetTags.Remove(e); removed {
			anyRemoved = true
		}
	}
	return anyRemoved
}

func (w *World) HasPlanetTag(entity Entity) bool {
	return w.planetTags.Contains(entity)
}

func (w *World) AllPlanetTags(yield func(e Entity) bool) {
	for e := range w.planetTags.All {
		if !yield(e) {
			break
		}
	}
}

// PlanetBuilder
func WithPlanetTag() EntityBuilderOption {
	return func(w *World, e Entity) {
		w.planetTags.Upsert(e, empty{})
	}
}

// Resource
func (w *World) ResourceUpsertPlanetTag() {
	w.planetTags.Upsert(w.resourceEntity, empty{})
}

func (w *World) ResourceRemovePlanetTag() {
	w.planetTags.Remove(w.resourceEntity)
}

func (w *World) ResourceHasPlanetTag() bool {
	return w.planetTags.Contains(w.resourceEntity)
}

// Events
