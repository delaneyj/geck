package ecs

type Eats struct {
	Entities []Entity
	Amounts  []uint8
}

func (w *World) SetEats(e Entity, c Eats) (old Eats, wasAdded bool) {
	old, wasAdded = w.eatsComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
}

func (w *World) SetEatsFromValues(
	e Entity,
	entitiesArg []Entity,
	amountsArg []uint8,
) {
	old, _ := w.SetEats(e, Eats{
		Entities: entitiesArg,
		Amounts:  amountsArg,
	})

	// depending on the generation flags, these might be unused
	_ = old

}

func (w *World) Eats(e Entity) (c Eats, ok bool) {
	return w.eatsComponents.Data(e)
}

func (w *World) MutableEats(e Entity) (c *Eats, ok bool) {
	return w.eatsComponents.DataMutable(e)
}

func (w *World) MustEats(e Entity) Eats {
	c, ok := w.eatsComponents.Data(e)
	if !ok {
		panic("entity does not have Eats")
	}
	return c
}

func (w *World) RemoveEats(e Entity) {
	wasRemoved := w.eatsComponents.Remove(e)

	// depending on the generation flags, these might be unused
	_ = wasRemoved

}

func (w *World) HasEats(e Entity) bool {
	return w.eatsComponents.Contains(e)
}

func (w *World) AllEats(yield func(e Entity, c Eats) bool) {
	for e, c := range w.eatsComponents.All {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutableEats(yield func(e Entity, c *Eats) bool) {
	for e, c := range w.eatsComponents.AllMutable {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllEatsEntities(yield func(e Entity) bool) {
	for e := range w.eatsComponents.AllEntities {
		if !yield(e) {
			break
		}
	}
}

// EatsBuilder
func WithEats(c Eats) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.eatsComponents.Upsert(e, c)
	}
}

func WithEatsFromValues(
	entitiesArg []Entity,
	amountsArg []uint8,
) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.SetEatsFromValues(e,
			entitiesArg,
			amountsArg,
		)
	}
}

// Events

// Resource methods
func (w *World) SetEatsResource(c Eats) {
	w.eatsComponents.Upsert(w.resourceEntity, c)
}

func (w *World) SetEatsResourceFromValues(
	entitiesArg []Entity,
	amountsArg []uint8,
) {
	w.SetEatsResource(Eats{
		Entities: entitiesArg,
		Amounts:  amountsArg,
	})
}

func (w *World) EatsResource() (Eats, bool) {
	return w.eatsComponents.Data(w.resourceEntity)
}

func (w *World) MustEatsResource() Eats {
	c, ok := w.EatsResource()
	if !ok {
		panic("resource entity does not have Eats")
	}
	return c
}

func (w *World) RemoveEatsResource() {
	w.eatsComponents.Remove(w.resourceEntity)
}
