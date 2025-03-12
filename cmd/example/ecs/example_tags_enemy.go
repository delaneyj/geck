package ecs

func (w *World) TagWithEnemy(entities ...Entity) (anyUpdated bool) {
	for _, e := range entities {
		if _, updated := w.enemyTags.Upsert(e, empty{}); updated {
			anyUpdated = true
		}
	}

	return anyUpdated
}

func (w *World) RemoveEnemyTag(entities ...Entity) (anyRemoved bool) {
	for _, e := range entities {
		if removed := w.enemyTags.Remove(e); removed {
			anyRemoved = true
		}
	}
	return anyRemoved
}

func (w *World) HasEnemyTag(entity Entity) bool {
	return w.enemyTags.Contains(entity)
}

func (w *World) EnemyTagCount() int {
	return w.enemyTags.Len()
}

func (w *World) EnemyTagCapacity() int {
	return w.enemyTags.Cap()
}

func (w *World) AllEnemyEntities(yield func(e Entity) bool) {
	for e := range w.enemyTags.All {
		if !yield(e) {
			break
		}
	}
}

// EnemyBuilder
func WithEnemyTag() EntityBuilderOption {
	return func(w *World, e Entity) {
		w.enemyTags.Upsert(e, empty{})
	}
}

// Resource
func (w *World) ResourceUpsertEnemyTag() {
	w.enemyTags.Upsert(w.resourceEntity, empty{})
}

func (w *World) ResourceRemoveEnemyTag() {
	w.enemyTags.Remove(w.resourceEntity)
}

func (w *World) ResourceHasEnemyTag() bool {
	return w.enemyTags.Contains(w.resourceEntity)
}

// Events
