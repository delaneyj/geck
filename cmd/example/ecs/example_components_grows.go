package ecs

type Grows struct {
	Entity []Entity
}

func (w *World) SetGrows(e Entity, c Grows) (old Grows, wasAdded bool) {
	old, wasAdded = w.growsComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
}

func (w *World) SetGrowsFromValues(
	e Entity,
	entityArg []Entity,
) {
	old, _ := w.SetGrows(e, Grows{
		Entity: entityArg,
	})

	// depending on the generation flags, these might be unused
	_ = old

}

func (w *World) Grows(e Entity) (c Grows, ok bool) {
	return w.growsComponents.Data(e)
}

func (w *World) MutableGrows(e Entity) (c *Grows, ok bool) {
	return w.growsComponents.DataMutable(e)
}

func (w *World) MustGrows(e Entity) Grows {
	c, ok := w.growsComponents.Data(e)
	if !ok {
		panic("entity does not have Grows")
	}
	return c
}

func (w *World) RemoveGrows(e Entity) {
	wasRemoved := w.growsComponents.Remove(e)

	// depending on the generation flags, these might be unused
	_ = wasRemoved

}

func (w *World) HasGrows(e Entity) bool {
	return w.growsComponents.Contains(e)
}

func (w *World) AllGrows(yield func(e Entity, c Grows) bool) {
	for e, c := range w.growsComponents.All {
		if yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutableGrows(yield func(e Entity, c *Grows) bool) {
	for e, c := range w.growsComponents.AllMutable {
		if yield(e, c) {
			break
		}
	}
}

func (w *World) AllGrowsEntities(yield func(e Entity) bool) {
	for e := range w.growsComponents.AllEntities {
		if yield(e) {
			break
		}
	}
}

// GrowsBuilder
func WithGrows(c Grows) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.growsComponents.Upsert(e, c)
	}
}

func WithGrowsFromValues(
	entityArg []Entity,
) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.SetGrowsFromValues(e,
			entityArg,
		)
	}
}

// Events

// Resource methods
func (w *World) SetGrowsResource(c Grows) {
	w.growsComponents.Upsert(w.resourceEntity, c)
}

func (w *World) SetGrowsResourceFromValues(
	entityArg []Entity,
) {
	w.SetGrowsResource(Grows{
		Entity: entityArg,
	})
}

func (w *World) GrowsResource() (Grows, bool) {
	return w.growsComponents.Data(w.resourceEntity)
}

func (w *World) MustGrowsResource() Grows {
	c, ok := w.GrowsResource()
	if !ok {
		panic("resource entity does not have Grows")
	}
	return c
}

func (w *World) RemoveGrowsResource() {
	w.growsComponents.Remove(w.resourceEntity)
}
