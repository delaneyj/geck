package ecs

type PositionComponent struct {
	X float32
	Y float32
	Z float32
}

func PositionComponentFromValues(
	xArg float32,
	yArg float32,
	zArg float32,
) PositionComponent {
	return PositionComponent{
		X: xArg,
		Y: yArg,
		Z: zArg,
	}
}

func DefaultPositionComponent() PositionComponent {
	return PositionComponent{
		X: 0.000000,
		Y: 0.000000,
		Z: 0.000000,
	}
}

func (c PositionComponent) Clone() PositionComponent {
	return PositionComponent{
		X: c.X,
		Y: c.Y,
		Z: c.Z,
	}
}

func (w *World) SetPosition(e Entity, c PositionComponent) (old PositionComponent, wasAdded bool) {
	old, wasAdded = w.positionComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
}

func (w *World) SetPositionFromValues(
	e Entity,
	xArg float32,
	yArg float32,
	zArg float32,
) {
	old, _ := w.SetPosition(e, PositionComponent{
		X: xArg,
		Y: yArg,
		Z: zArg,
	})

	// depending on the generation flags, these might be unused
	_ = old

}

func (w *World) Position(e Entity) (c PositionComponent, ok bool) {
	return w.positionComponents.Data(e)
}

func (w *World) MutablePosition(e Entity) (c *PositionComponent, ok bool) {
	return w.positionComponents.DataMutable(e)
}

func (w *World) MustMutablePosition(e Entity) *PositionComponent {
	c, ok := w.MutablePosition(e)
	if !ok {
		panic("entity does not have Position")
	}
	return c
}

func (w *World) MustPosition(e Entity) PositionComponent {
	c, ok := w.positionComponents.Data(e)
	if !ok {
		panic("entity does not have Position")
	}
	return c
}

func (w *World) RemovePosition(e Entity) {
	wasRemoved := w.positionComponents.Remove(e)

	// depending on the generation flags, these might be unused
	_ = wasRemoved

}

func (w *World) HasPosition(e Entity) bool {
	return w.positionComponents.Contains(e)
}

func (w *World) PositionsCount() int {
	return w.positionComponents.Len()
}

func (w *World) PositionsCapacity() int {
	return w.positionComponents.Cap()
}

func (w *World) AllPositions(yield func(e Entity, c PositionComponent) bool) {
	for e, c := range w.positionComponents.All {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutablePositions(yield func(e Entity, c *PositionComponent) bool) {
	for e, c := range w.positionComponents.AllMutable {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllPositionsEntities(yield func(e Entity) bool) {
	for e := range w.positionComponents.AllEntities {
		if !yield(e) {
			break
		}
	}
}

func (w *World) AllMutablePositionsEntities(yield func(e Entity) bool) {
	w.AllPositionsEntities(yield)
}

// PositionBuilder
func WithPositionDefault() EntityBuilderOption {
	return WithPosition(DefaultPositionComponent())
}

func WithPosition(c PositionComponent) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.positionComponents.Upsert(e, c)
	}
}

func WithPositionFromValues(
	xArg float32,
	yArg float32,
	zArg float32,
) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.SetPositionFromValues(e,
			xArg,
			yArg,
			zArg,
		)
	}
}

// Events

// Resource methods
func (w *World) SetPositionResource(c PositionComponent) {
	w.SetPosition(w.resourceEntity, c)
}

func (w *World) SetPositionResourceFromValues(
	xArg float32,
	yArg float32,
	zArg float32,
) {
	w.SetPositionResource(PositionComponent{
		X: xArg,
		Y: yArg,
		Z: zArg,
	})
}

func (w *World) PositionResource() (PositionComponent, bool) {
	return w.positionComponents.Data(w.resourceEntity)
}

func (w *World) MustPositionResource() PositionComponent {
	c, ok := w.PositionResource()
	if !ok {
		panic("resource entity does not have Position")
	}
	return c
}

func (w *World) RemovePositionResource() {
	w.positionComponents.Remove(w.resourceEntity)
}

func (w *World) HasPositionResource() bool {
	return w.positionComponents.Contains(w.resourceEntity)
}
