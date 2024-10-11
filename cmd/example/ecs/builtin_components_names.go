package ecs

type NameComponent struct {
	Value string
}

func (w *World) SetName(e Entity, c NameComponent) (old NameComponent, wasAdded bool) {
	old, wasAdded = w.nameComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
}

func (w *World) SetNameFromValues(
	e Entity,
	valueArg string,
) {
	old, _ := w.SetName(e, NameComponent{
		Value: valueArg,
	})

	// depending on the generation flags, these might be unused
	_ = old

}

func (w *World) Name(e Entity) (c NameComponent, ok bool) {
	return w.nameComponents.Data(e)
}

func (w *World) MutableName(e Entity) (c *NameComponent, ok bool) {
	return w.nameComponents.DataMutable(e)
}

func (w *World) MustName(e Entity) NameComponent {
	c, ok := w.nameComponents.Data(e)
	if !ok {
		panic("entity does not have Name")
	}
	return c
}

func (w *World) RemoveName(e Entity) {
	wasRemoved := w.nameComponents.Remove(e)

	// depending on the generation flags, these might be unused
	_ = wasRemoved

}

func (w *World) HasName(e Entity) bool {
	return w.nameComponents.Contains(e)
}

func (w *World) AllNames(yield func(e Entity, c NameComponent) bool) {
	for e, c := range w.nameComponents.All {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutableNames(yield func(e Entity, c *NameComponent) bool) {
	for e, c := range w.nameComponents.AllMutable {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllNamesEntities(yield func(e Entity) bool) {
	for e := range w.nameComponents.AllEntities {
		if !yield(e) {
			break
		}
	}
}

// NameBuilder
func WithName(c NameComponent) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.nameComponents.Upsert(e, c)
	}
}

func WithNameFromValues(
	valueArg string,
) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.SetNameFromValues(e,
			valueArg,
		)
	}
}

// Events

// Resource methods
func (w *World) SetNameResource(c NameComponent) {
	w.SetName(w.resourceEntity, c)
}

func (w *World) SetNameResourceFromValues(
	valueArg string,
) {
	w.SetNameResource(NameComponent{
		Value: valueArg,
	})
}

func (w *World) NameResource() (NameComponent, bool) {
	return w.nameComponents.Data(w.resourceEntity)
}

func (w *World) MustNameResource() NameComponent {
	c, ok := w.NameResource()
	if !ok {
		panic("resource entity does not have Name")
	}
	return c
}

func (w *World) RemoveNameResource() {
	w.nameComponents.Remove(w.resourceEntity)
}
