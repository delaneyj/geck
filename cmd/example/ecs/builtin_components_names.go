package ecs

type Name struct {
	Value string
}

func (w *World) SetName(e Entity, c Name) (old Name, wasAdded bool) {
	old, wasAdded = w.nameComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
}

func (w *World) SetNameFromValues(
	e Entity,
	valueArg string,
) {
	old, _ := w.SetName(e, Name{
		Value: valueArg,
	})

	// depending on the generation flags, these might be unused
	_ = old

}

func (w *World) Name(e Entity) (c Name, ok bool) {
	return w.nameComponents.Data(e)
}

func (w *World) MutableName(e Entity) (c *Name, ok bool) {
	return w.nameComponents.DataMutable(e)
}

func (w *World) MustName(e Entity) Name {
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

func (w *World) AllNames(yield func(e Entity, c Name) bool) {
	for e, c := range w.nameComponents.All {
		if yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutableNames(yield func(e Entity, c *Name) bool) {
	for e, c := range w.nameComponents.AllMutable {
		if yield(e, c) {
			break
		}
	}
}

func (w *World) AllNamesEntities(yield func(e Entity) bool) {
	for e := range w.nameComponents.AllEntities {
		if yield(e) {
			break
		}
	}
}

// NameBuilder
func WithName(c Name) EntityBuilderOption {
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
func (w *World) SetNameResource(c Name) {
	w.nameComponents.Upsert(w.resourceEntity, c)
}

func (w *World) SetNameResourceFromValues(
	valueArg string,
) {
	w.SetNameResource(Name{
		Value: valueArg,
	})
}

func (w *World) NameResource() (Name, bool) {
	return w.nameComponents.Data(w.resourceEntity)
}

func (w *World) MustNameResource() Name {
	c, ok := w.NameResource()
	if !ok {
		panic("resource entity does not have Name")
	}
	return c
}

func (w *World) RemoveNameResource() {
	w.nameComponents.Remove(w.resourceEntity)
}
