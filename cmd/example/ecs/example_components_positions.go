package ecs

type Position struct {
	X float32
	Y float32
	Z float32
}

func (w *World) SetPosition(e Entity, c Position) (old Position, wasAdded bool) {
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
	old, _ := w.SetPosition(e, Position{
		X: xArg,
		Y: yArg,
		Z: zArg,
	})

	// depending on the generation flags, these might be unused
	_ = old

}

func (w *World) Position(e Entity) (c Position, ok bool) {
	return w.positionComponents.Data(e)
}

func (w *World) MutablePosition(e Entity) (c *Position, ok bool) {
	return w.positionComponents.DataMutable(e)
}

func (w *World) MustPosition(e Entity) Position {
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

func (w *World) AllPositions(yield func(e Entity, c Position) bool) {
	for e, c := range w.positionComponents.All {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutablePositions(yield func(e Entity, c *Position) bool) {
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

// PositionBuilder
func WithPosition(c Position) EntityBuilderOption {
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
func (w *World) SetPositionResource(c Position) {
	w.positionComponents.Upsert(w.resourceEntity, c)
}

func (w *World) SetPositionResourceFromValues(
	xArg float32,
	yArg float32,
	zArg float32,
) {
	w.SetPositionResource(Position{
		X: xArg,
		Y: yArg,
		Z: zArg,
	})
}

func (w *World) PositionResource() (Position, bool) {
	return w.positionComponents.Data(w.resourceEntity)
}

func (w *World) MustPositionResource() Position {
	c, ok := w.PositionResource()
	if !ok {
		panic("resource entity does not have Position")
	}
	return c
}

func (w *World) RemovePositionResource() {
	w.positionComponents.Remove(w.resourceEntity)
}
