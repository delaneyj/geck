package ecs

type DockedToComponent struct {
	Entity Entity
}

func DockedToComponentFromValues(
	entityArg Entity,
) DockedToComponent {
	return DockedToComponent{
		Entity: entityArg,
	}
}

func DefaultDockedToComponent() DockedToComponent {
	return DockedToComponent{
		Entity: EntityFromU32(0),
	}
}

func (c DockedToComponent) Clone() DockedToComponent {
	return DockedToComponent{
		Entity: c.Entity,
	}
}

func (w *World) SetDockedTo(e Entity, arg Entity) (old DockedToComponent, wasAdded bool) {
	c := DockedToComponent{
		Entity: arg,
	}
	old, wasAdded = w.dockedToComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
}

func (w *World) DockedTo(e Entity) (c DockedToComponent, ok bool) {
	return w.dockedToComponents.Data(e)
}

func (w *World) MutableDockedTo(e Entity) (c *DockedToComponent, ok bool) {
	return w.dockedToComponents.DataMutable(e)
}

func (w *World) MustMutableDockedTo(e Entity) *DockedToComponent {
	c, ok := w.MutableDockedTo(e)
	if !ok {
		panic("entity does not have DockedTo")
	}
	return c
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

func (w *World) DockedTosCount() int {
	return w.dockedToComponents.Len()
}

func (w *World) DockedTosCapacity() int {
	return w.dockedToComponents.Cap()
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

func (w *World) AllMutableDockedTosEntities(yield func(e Entity) bool) {
	w.AllDockedTosEntities(yield)
}

// DockedToBuilder
func WithDockedToDefault() EntityBuilderOption {
	return WithDockedTo(DefaultDockedToComponent().Entity)
}

func WithDockedTo(arg Entity) EntityBuilderOption {
	c := DockedToComponent{
		Entity: arg,
	}
	return func(w *World, e Entity) {
		w.dockedToComponents.Upsert(e, c)
	}
}

// Events

// Resource methods
func (w *World) SetDockedToResource(arg Entity) {
	w.SetDockedTo(w.resourceEntity, arg)
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

func (w *World) HasDockedToResource() bool {
	return w.dockedToComponents.Contains(w.resourceEntity)
}
