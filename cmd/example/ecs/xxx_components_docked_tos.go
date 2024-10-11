package ecs

type DockedToComponent struct {
	Entity Entity
}

func (w *World) SetDockedTo(e Entity, c DockedToComponent) (old DockedToComponent, wasAdded bool) {
	old, wasAdded = w.dockedToComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
}

func (w *World) SetDockedToFromValues(
	e Entity,
	entityArg Entity,
) {
	old, _ := w.SetDockedTo(e, DockedToComponent{
		Entity: entityArg,
	})

	// depending on the generation flags, these might be unused
	_ = old

}

func (w *World) DockedTo(e Entity) (c DockedToComponent, ok bool) {
	return w.dockedToComponents.Data(e)
}

func (w *World) MutableDockedTo(e Entity) (c *DockedToComponent, ok bool) {
	return w.dockedToComponents.DataMutable(e)
}

func (w *World) MustDockedTo(e Entity) DockedToComponent {
	c, ok := w.dockedToComponents.Data(e)
	if !ok {
		panic("entity does not have DockedTo")
	}
	return c
}

func (w *World) RemoveDockedTo(e Entity) {
	wasRemoved := w.dockedToComponents.Remove(e)

	// depending on the generation flags, these might be unused
	_ = wasRemoved

}

func (w *World) HasDockedTo(e Entity) bool {
	return w.dockedToComponents.Contains(e)
}

func (w *World) AllDockedTos(yield func(e Entity, c DockedToComponent) bool) {
	for e, c := range w.dockedToComponents.All {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutableDockedTos(yield func(e Entity, c *DockedToComponent) bool) {
	for e, c := range w.dockedToComponents.AllMutable {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllDockedTosEntities(yield func(e Entity) bool) {
	for e := range w.dockedToComponents.AllEntities {
		if !yield(e) {
			break
		}
	}
}

// DockedToBuilder
func WithDockedTo(c DockedToComponent) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.dockedToComponents.Upsert(e, c)
	}
}

func WithDockedToFromValues(
	entityArg Entity,
) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.SetDockedToFromValues(e,
			entityArg,
		)
	}
}

// Events

// Resource methods
func (w *World) SetDockedToResource(c DockedToComponent) {
	w.SetDockedTo(w.resourceEntity, c)
}

func (w *World) SetDockedToResourceFromValues(
	entityArg Entity,
) {
	w.SetDockedToResource(DockedToComponent{
		Entity: entityArg,
	})
}

func (w *World) DockedToResource() (DockedToComponent, bool) {
	return w.dockedToComponents.Data(w.resourceEntity)
}

func (w *World) MustDockedToResource() DockedToComponent {
	c, ok := w.DockedToResource()
	if !ok {
		panic("resource entity does not have DockedTo")
	}
	return c
}

func (w *World) RemoveDockedToResource() {
	w.dockedToComponents.Remove(w.resourceEntity)
}
