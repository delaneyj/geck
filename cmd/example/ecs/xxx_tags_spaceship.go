package ecs

func (w *World) TagWithSpaceship(entities ...Entity) (anyUpdated bool) {
	for _, e := range entities {
		if _, updated := w.spaceshipTags.Upsert(e, empty{}); updated {
			anyUpdated = true
		}
	}

	return anyUpdated
}

func (w *World) RemoveSpaceshipTag(entities ...Entity) (anyRemoved bool) {
	for _, e := range entities {
		if removed := w.spaceshipTags.Remove(e); removed {
			anyRemoved = true
		}
	}
	return anyRemoved
}

func (w *World) HasSpaceshipTag(entity Entity) bool {
	return w.spaceshipTags.Contains(entity)
}

func (w *World) SpaceshipTagCount() int {
	return w.spaceshipTags.Len()
}

func (w *World) SpaceshipTagCapacity() int {
	return w.spaceshipTags.Cap()
}

func (w *World) AllSpaceshipEntities(yield func(e Entity) bool) {
	for e := range w.spaceshipTags.All {
		if !yield(e) {
			break
		}
	}
}

// SpaceshipBuilder
func WithSpaceshipTag() EntityBuilderOption {
	return func(w *World, e Entity) {
		w.spaceshipTags.Upsert(e, empty{})
	}
}

// Resource
func (w *World) ResourceUpsertSpaceshipTag() {
	w.spaceshipTags.Upsert(w.resourceEntity, empty{})
}

func (w *World) ResourceRemoveSpaceshipTag() {
	w.spaceshipTags.Remove(w.resourceEntity)
}

func (w *World) ResourceHasSpaceshipTag() bool {
	return w.spaceshipTags.Contains(w.resourceEntity)
}

// Events
