package ecs

func (w *World) TagWithSpacestation(entities ...Entity) (anyUpdated bool) {
	for _, e := range entities {
		if _, updated := w.spacestationTags.Upsert(e, empty{}); updated {
			anyUpdated = true
		}
	}

	return anyUpdated
}

func (w *World) RemoveSpacestationTag(entities ...Entity) (anyRemoved bool) {
	for _, e := range entities {
		if removed := w.spacestationTags.Remove(e); removed {
			anyRemoved = true
		}
	}
	return anyRemoved
}

func (w *World) HasSpacestationTag(entity Entity) bool {
	return w.spacestationTags.Contains(entity)
}

func (w *World) AllSpacestationTags(yield func(e Entity) bool) {
	for e := range w.spacestationTags.All {
		if !yield(e) {
			break
		}
	}
}

// SpacestationBuilder
func WithSpacestationTag() EntityBuilderOption {
	return func(w *World, e Entity) {
		w.spacestationTags.Upsert(e, empty{})
	}
}

// Resource
func (w *World) ResourceUpsertSpacestationTag() {
	w.spacestationTags.Upsert(w.resourceEntity, empty{})
}

func (w *World) ResourceRemoveSpacestationTag() {
	w.spacestationTags.Remove(w.resourceEntity)
}

func (w *World) ResourceHasSpacestationTag() bool {
	return w.spacestationTags.Contains(w.resourceEntity)
}

// Events
