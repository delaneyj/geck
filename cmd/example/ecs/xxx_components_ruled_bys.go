package ecs

type RuledByComponent struct {
	Entity Entity
}

func RuledByComponentFromValues(
	entityArg Entity,
) RuledByComponent {
	return RuledByComponent{
		Entity: entityArg,
	}
}

func DefaultRuledByComponent() RuledByComponent {
	return RuledByComponent{
		Entity: EntityFromU32(0),
	}
}

func (c RuledByComponent) Clone() RuledByComponent {
	return RuledByComponent{
		Entity: c.Entity,
	}
}

func (w *World) SetRuledBy(e Entity, arg Entity) (old RuledByComponent, wasAdded bool) {
	c := RuledByComponent{
		Entity: arg,
	}
	old, wasAdded = w.ruledByComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
}

func (w *World) RuledBy(e Entity) (c RuledByComponent, ok bool) {
	return w.ruledByComponents.Data(e)
}

func (w *World) MutableRuledBy(e Entity) (c *RuledByComponent, ok bool) {
	return w.ruledByComponents.DataMutable(e)
}

func (w *World) MustMutableRuledBy(e Entity) *RuledByComponent {
	c, ok := w.MutableRuledBy(e)
	if !ok {
		panic("entity does not have RuledBy")
	}
	return c
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

func (w *World) RuledBysCount() int {
	return w.ruledByComponents.Len()
}

func (w *World) RuledBysCapacity() int {
	return w.ruledByComponents.Cap()
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

func (w *World) AllMutableRuledBysEntities(yield func(e Entity) bool) {
	w.AllRuledBysEntities(yield)
}

// RuledByBuilder
func WithRuledByDefault() EntityBuilderOption {
	return WithRuledBy(DefaultRuledByComponent().Entity)
}

func WithRuledBy(arg Entity) EntityBuilderOption {
	c := RuledByComponent{
		Entity: arg,
	}
	return func(w *World, e Entity) {
		w.ruledByComponents.Upsert(e, c)
	}
}

// Events

// Resource methods
func (w *World) SetRuledByResource(arg Entity) {
	w.SetRuledBy(w.resourceEntity, arg)
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

func (w *World) HasRuledByResource() bool {
	return w.ruledByComponents.Contains(w.resourceEntity)
}
