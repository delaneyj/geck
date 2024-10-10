package ecs

type AlliedWith struct {
	Entity []Entity
}

func (w *World) SetAlliedWith(e Entity, c AlliedWith) (old AlliedWith, wasAdded bool) {
	old, wasAdded = w.alliedWithComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
}

func (w *World) SetAlliedWithFromValues(
	e Entity,
	entityArg []Entity,
) {
	old, _ := w.SetAlliedWith(e, AlliedWith{
		Entity: entityArg,
	})

	// depending on the generation flags, these might be unused
	_ = old

}

func (w *World) AlliedWith(e Entity) (c AlliedWith, ok bool) {
	return w.alliedWithComponents.Data(e)
}

func (w *World) MutableAlliedWith(e Entity) (c *AlliedWith, ok bool) {
	return w.alliedWithComponents.DataMutable(e)
}

func (w *World) MustAlliedWith(e Entity) AlliedWith {
	c, ok := w.alliedWithComponents.Data(e)
	if !ok {
		panic("entity does not have AlliedWith")
	}
	return c
}

func (w *World) RemoveAlliedWith(e Entity) {
	wasRemoved := w.alliedWithComponents.Remove(e)

	// depending on the generation flags, these might be unused
	_ = wasRemoved

}

func (w *World) HasAlliedWith(e Entity) bool {
	return w.alliedWithComponents.Contains(e)
}

func (w *World) AllAlliedWiths(yield func(e Entity, c AlliedWith) bool) {
	for e, c := range w.alliedWithComponents.All {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutableAlliedWiths(yield func(e Entity, c *AlliedWith) bool) {
	for e, c := range w.alliedWithComponents.AllMutable {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllAlliedWithsEntities(yield func(e Entity) bool) {
	for e := range w.alliedWithComponents.AllEntities {
		if !yield(e) {
			break
		}
	}
}

// AlliedWithBuilder
func WithAlliedWith(c AlliedWith) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.alliedWithComponents.Upsert(e, c)
	}
}

func WithAlliedWithFromValues(
	entityArg []Entity,
) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.SetAlliedWithFromValues(e,
			entityArg,
		)
	}
}

// Events

// Resource methods
func (w *World) SetAlliedWithResource(c AlliedWith) {
	w.alliedWithComponents.Upsert(w.resourceEntity, c)
}

func (w *World) SetAlliedWithResourceFromValues(
	entityArg []Entity,
) {
	w.SetAlliedWithResource(AlliedWith{
		Entity: entityArg,
	})
}

func (w *World) AlliedWithResource() (AlliedWith, bool) {
	return w.alliedWithComponents.Data(w.resourceEntity)
}

func (w *World) MustAlliedWithResource() AlliedWith {
	c, ok := w.AlliedWithResource()
	if !ok {
		panic("resource entity does not have AlliedWith")
	}
	return c
}

func (w *World) RemoveAlliedWithResource() {
	w.alliedWithComponents.Remove(w.resourceEntity)
}
