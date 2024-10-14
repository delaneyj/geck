package ecs

type DirectionComponent struct {
	Values EnumDirection
}

func DirectionComponentFromValues(
	valuesArg EnumDirection,
) DirectionComponent {
	return DirectionComponent{
		Values: valuesArg,
	}
}

func DefaultDirectionComponent() DirectionComponent {
	return DirectionComponent{
		Values: EnumDirection(0),
	}
}

func (c DirectionComponent) Clone() DirectionComponent {
	return DirectionComponent{
		Values: c.Values,
	}
}

func (w *World) SetDirection(e Entity, arg EnumDirection) (old DirectionComponent, wasAdded bool) {
	c := DirectionComponent{
		Values: arg,
	}
	old, wasAdded = w.directionComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
}

func (w *World) Direction(e Entity) (c DirectionComponent, ok bool) {
	return w.directionComponents.Data(e)
}

func (w *World) MutableDirection(e Entity) (c *DirectionComponent, ok bool) {
	return w.directionComponents.DataMutable(e)
}

func (w *World) MustDirection(e Entity) DirectionComponent {
	c, ok := w.directionComponents.Data(e)
	if !ok {
		panic("entity does not have Direction")
	}
	return c
}

func (w *World) RemoveDirection(e Entity) {
	wasRemoved := w.directionComponents.Remove(e)

	// depending on the generation flags, these might be unused
	_ = wasRemoved

}

func (w *World) HasDirection(e Entity) bool {
	return w.directionComponents.Contains(e)
}

func (w *World) AllDirections(yield func(e Entity, c DirectionComponent) bool) {
	for e, c := range w.directionComponents.All {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutableDirections(yield func(e Entity, c *DirectionComponent) bool) {
	for e, c := range w.directionComponents.AllMutable {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllDirectionsEntities(yield func(e Entity) bool) {
	for e := range w.directionComponents.AllEntities {
		if !yield(e) {
			break
		}
	}
}

// DirectionBuilder

func WithDirection(arg EnumDirection) EntityBuilderOption {
	c := DirectionComponent{
		Values: arg,
	}

	return func(w *World, e Entity) {
		w.directionComponents.Upsert(e, c)
	}
}

// Events

// Resource methods

func (w *World) SetDirectionResource(arg EnumDirection) {
	w.SetDirection(w.resourceEntity, arg)
}

func (w *World) DirectionResource() (DirectionComponent, bool) {
	return w.directionComponents.Data(w.resourceEntity)
}

func (w *World) MustDirectionResource() DirectionComponent {
	c, ok := w.DirectionResource()
	if !ok {
		panic("resource entity does not have Direction")
	}
	return c
}

func (w *World) RemoveDirectionResource() {
	w.directionComponents.Remove(w.resourceEntity)
}

func (w *World) HasDirectionResource() bool {
	return w.directionComponents.Contains(w.resourceEntity)
}
