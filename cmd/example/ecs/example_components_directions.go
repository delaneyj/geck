package ecs

type Direction struct {
	Values EnumDirection
}

func (w *World) SetDirection(e Entity, c Direction) (old Direction, wasAdded bool) {
	old, wasAdded = w.directionComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
}

func (w *World) SetDirectionFromValues(
	e Entity,
	valuesArg EnumDirection,
) {
	old, _ := w.SetDirection(e, Direction{
		Values: valuesArg,
	})

	// depending on the generation flags, these might be unused
	_ = old

}

func (w *World) Direction(e Entity) (c Direction, ok bool) {
	return w.directionComponents.Data(e)
}

func (w *World) MutableDirection(e Entity) (c *Direction, ok bool) {
	return w.directionComponents.DataMutable(e)
}

func (w *World) MustDirection(e Entity) Direction {
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

func (w *World) AllDirections(yield func(e Entity, c Direction) bool) {
	for e, c := range w.directionComponents.All {
		if yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutableDirections(yield func(e Entity, c *Direction) bool) {
	for e, c := range w.directionComponents.AllMutable {
		if yield(e, c) {
			break
		}
	}
}

func (w *World) AllDirectionsEntities(yield func(e Entity) bool) {
	for e := range w.directionComponents.AllEntities {
		if yield(e) {
			break
		}
	}
}

// DirectionBuilder
func WithDirection(c Direction) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.directionComponents.Upsert(e, c)
	}
}

func WithDirectionFromValues(
	valuesArg EnumDirection,
) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.SetDirectionFromValues(e,
			valuesArg,
		)
	}
}

// Events

// Resource methods
func (w *World) SetDirectionResource(c Direction) {
	w.directionComponents.Upsert(w.resourceEntity, c)
}

func (w *World) SetDirectionResourceFromValues(
	valuesArg EnumDirection,
) {
	w.SetDirectionResource(Direction{
		Values: valuesArg,
	})
}

func (w *World) DirectionResource() (Direction, bool) {
	return w.directionComponents.Data(w.resourceEntity)
}

func (w *World) MustDirectionResource() Direction {
	c, ok := w.DirectionResource()
	if !ok {
		panic("resource entity does not have Direction")
	}
	return c
}

func (w *World) RemoveDirectionResource() {
	w.directionComponents.Remove(w.resourceEntity)
}
