package ecs

type RuledByComponent struct {
	Entity Entity
}

func (w *World) SetRuledBy(e Entity, c RuledByComponent) (old RuledByComponent, wasAdded bool) {
	old, wasAdded = w.ruledByComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
}

func (w *World) SetRuledByFromValues(
	e Entity,
	entityArg Entity,
) {
	old, _ := w.SetRuledBy(e, RuledByComponent{
		Entity: entityArg,
	})

	// depending on the generation flags, these might be unused
	_ = old

}

func (w *World) RuledBy(e Entity) (c RuledByComponent, ok bool) {
	return w.ruledByComponents.Data(e)
}

func (w *World) MutableRuledBy(e Entity) (c *RuledByComponent, ok bool) {
	return w.ruledByComponents.DataMutable(e)
}

func (w *World) MustRuledBy(e Entity) RuledByComponent {
	c, ok := w.ruledByComponents.Data(e)
	if !ok {
		panic("entity does not have RuledBy")
	}
	return c
}

func (w *World) RemoveRuledBy(e Entity) {
	wasRemoved := w.ruledByComponents.Remove(e)

	// depending on the generation flags, these might be unused
	_ = wasRemoved

}

func (w *World) HasRuledBy(e Entity) bool {
	return w.ruledByComponents.Contains(e)
}

func (w *World) AllRuledBys(yield func(e Entity, c RuledByComponent) bool) {
	for e, c := range w.ruledByComponents.All {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutableRuledBys(yield func(e Entity, c *RuledByComponent) bool) {
	for e, c := range w.ruledByComponents.AllMutable {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllRuledBysEntities(yield func(e Entity) bool) {
	for e := range w.ruledByComponents.AllEntities {
		if !yield(e) {
			break
		}
	}
}

// RuledByBuilder
func WithRuledBy(c RuledByComponent) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.ruledByComponents.Upsert(e, c)
	}
}

func WithRuledByFromValues(
	entityArg Entity,
) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.SetRuledByFromValues(e,
			entityArg,
		)
	}
}

// Events

// Resource methods
func (w *World) SetRuledByResource(c RuledByComponent) {
	w.SetRuledBy(w.resourceEntity, c)
}

func (w *World) SetRuledByResourceFromValues(
	entityArg Entity,
) {
	w.SetRuledByResource(RuledByComponent{
		Entity: entityArg,
	})
}

func (w *World) RuledByResource() (RuledByComponent, bool) {
	return w.ruledByComponents.Data(w.resourceEntity)
}

func (w *World) MustRuledByResource() RuledByComponent {
	c, ok := w.RuledByResource()
	if !ok {
		panic("resource entity does not have RuledBy")
	}
	return c
}

func (w *World) RemoveRuledByResource() {
	w.ruledByComponents.Remove(w.resourceEntity)
}
