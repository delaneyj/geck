package ecs

type IsAComponent struct {
	Prototype Entity
}

func IsAComponentFromValues(
	prototypeArg Entity,
) IsAComponent {
	return IsAComponent{
		Prototype: prototypeArg,
	}
}

func DefaultIsAComponent() IsAComponent {
	return IsAComponent{
		Prototype: EntityFromU32(0),
	}
}

func (c IsAComponent) Clone() IsAComponent {
	return IsAComponent{
		Prototype: c.Prototype,
	}
}

func (w *World) SetIsA(e Entity, arg Entity) (old IsAComponent, wasAdded bool) {
	c := IsAComponent{
		Prototype: arg,
	}
	old, wasAdded = w.isAComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
}

func (w *World) IsA(e Entity) (c IsAComponent, ok bool) {
	return w.isAComponents.Data(e)
}

func (w *World) MutableIsA(e Entity) (c *IsAComponent, ok bool) {
	return w.isAComponents.DataMutable(e)
}

func (w *World) MustIsA(e Entity) IsAComponent {
	c, ok := w.isAComponents.Data(e)
	if !ok {
		panic("entity does not have IsA")
	}
	return c
}

func (w *World) RemoveIsA(e Entity) {
	wasRemoved := w.isAComponents.Remove(e)

	// depending on the generation flags, these might be unused
	_ = wasRemoved

}

func (w *World) HasIsA(e Entity) bool {
	return w.isAComponents.Contains(e)
}

func (w *World) AllIsA(yield func(e Entity, c IsAComponent) bool) {
	for e, c := range w.isAComponents.All {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutableIsA(yield func(e Entity, c *IsAComponent) bool) {
	for e, c := range w.isAComponents.AllMutable {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllIsAEntities(yield func(e Entity) bool) {
	for e := range w.isAComponents.AllEntities {
		if !yield(e) {
			break
		}
	}
}

// IsABuilder

func WithIsA(arg Entity) EntityBuilderOption {
	c := IsAComponent{
		Prototype: arg,
	}

	return func(w *World, e Entity) {
		w.isAComponents.Upsert(e, c)
	}
}

// Events

// Resource methods

func (w *World) SetIsAResource(arg Entity) {
	w.SetIsA(w.resourceEntity, arg)
}

func (w *World) IsAResource() (IsAComponent, bool) {
	return w.isAComponents.Data(w.resourceEntity)
}

func (w *World) MustIsAResource() IsAComponent {
	c, ok := w.IsAResource()
	if !ok {
		panic("resource entity does not have IsA")
	}
	return c
}

func (w *World) RemoveIsAResource() {
	w.isAComponents.Remove(w.resourceEntity)
}

func (w *World) HasIsAResource() bool {
	return w.isAComponents.Contains(w.resourceEntity)
}
