package ecs

type AlliedWithComponent struct {
	Entity []Entity
}

func (w *World) SetAlliedWith(e Entity, c AlliedWithComponent) (old AlliedWithComponent, wasAdded bool) {
	old, wasAdded = w.alliedWithComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
}

func (w *World) SetAlliedWithFromValues(
	e Entity,
	entityArg []Entity,
) {
	old, _ := w.SetAlliedWith(e, AlliedWithComponent{
		Entity: entityArg,
	})

	// depending on the generation flags, these might be unused
	_ = old

}

func (w *World) AlliedWith(e Entity) (c AlliedWithComponent, ok bool) {
	return w.alliedWithComponents.Data(e)
}

func (w *World) MutableAlliedWith(e Entity) (c *AlliedWithComponent, ok bool) {
	return w.alliedWithComponents.DataMutable(e)
}

func (w *World) MustAlliedWith(e Entity) AlliedWithComponent {
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

func (w *World) AllAlliedWiths(yield func(e Entity, c AlliedWithComponent) bool) {
	for e, c := range w.alliedWithComponents.All {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutableAlliedWiths(yield func(e Entity, c *AlliedWithComponent) bool) {
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
func WithAlliedWith(c AlliedWithComponent) EntityBuilderOption {
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
func (w *World) SetAlliedWithResource(c AlliedWithComponent) {
	w.SetAlliedWith(w.resourceEntity, c)
}

func (w *World) SetAlliedWithResourceFromValues(
	entityArg []Entity,
) {
	w.SetAlliedWithResource(AlliedWithComponent{
		Entity: entityArg,
	})
}

func (w *World) AlliedWithResource() (AlliedWithComponent, bool) {
	return w.alliedWithComponents.Data(w.resourceEntity)
}

func (w *World) MustAlliedWithResource() AlliedWithComponent {
	c, ok := w.AlliedWithResource()
	if !ok {
		panic("resource entity does not have AlliedWith")
	}
	return c
}

func (w *World) RemoveAlliedWithResource() {
	w.alliedWithComponents.Remove(w.resourceEntity)
}
