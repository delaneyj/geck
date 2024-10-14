package ecs

type EatsComponent struct {
	Entities []Entity
	Amounts  []uint8
}

func EatsComponentFromValues(
	entitiesArg []Entity,
	amountsArg []uint8,
) EatsComponent {
	return EatsComponent{
		Entities: entitiesArg,
		Amounts:  amountsArg,
	}
}

func DefaultEatsComponent() EatsComponent {
	return EatsComponent{
		Entities: nil,
		Amounts:  nil,
	}
}

func (c EatsComponent) Clone() EatsComponent {
	return EatsComponent{
		Entities: c.Entities,
		Amounts:  c.Amounts,
	}
}

func (w *World) SetEats(e Entity, c EatsComponent) (old EatsComponent, wasAdded bool) {
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
	old, _ := w.SetEats(e, EatsComponent{
		Entities: entitiesArg,
		Amounts:  amountsArg,
	})

	// depending on the generation flags, these might be unused
	_ = old

}

func (w *World) Eats(e Entity) (c EatsComponent, ok bool) {
	return w.eatsComponents.Data(e)
}

func (w *World) MutableEats(e Entity) (c *EatsComponent, ok bool) {
	return w.eatsComponents.DataMutable(e)
}

func (w *World) MustEats(e Entity) EatsComponent {
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

func (w *World) AllEats(yield func(e Entity, c EatsComponent) bool) {
	for e, c := range w.eatsComponents.All {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutableEats(yield func(e Entity, c *EatsComponent) bool) {
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

func WithEats(c EatsComponent) EntityBuilderOption {

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

func (w *World) SetEatsResource(c EatsComponent) {
	w.SetEats(w.resourceEntity, c)
}

func (w *World) SetEatsResourceFromValues(
	entitiesArg []Entity,
	amountsArg []uint8,
) {
	w.SetEatsResource(EatsComponent{
		Entities: entitiesArg,
		Amounts:  amountsArg,
	})
}

func (w *World) EatsResource() (EatsComponent, bool) {
	return w.eatsComponents.Data(w.resourceEntity)
}

func (w *World) MustEatsResource() EatsComponent {
	c, ok := w.EatsResource()
	if !ok {
		panic("resource entity does not have Eats")
	}
	return c
}

func (w *World) RemoveEatsResource() {
	w.eatsComponents.Remove(w.resourceEntity)
}

func (w *World) HasEatsResource() bool {
	return w.eatsComponents.Contains(w.resourceEntity)
}
