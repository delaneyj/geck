package ecs

type ChildOfComponent struct {
	Parent Entity
}

func ChildOfComponentFromValues(
	parentArg Entity,
) ChildOfComponent {
	return ChildOfComponent{
		Parent: parentArg,
	}
}

func DefaultChildOfComponent() ChildOfComponent {
	return ChildOfComponent{
		Parent: EntityFromU32(0),
	}
}

func (c ChildOfComponent) Clone() ChildOfComponent {
	return ChildOfComponent{
		Parent: c.Parent,
	}
}

func (w *World) SetChildOf(e Entity, arg Entity) (old ChildOfComponent, wasAdded bool) {
	c := ChildOfComponent{
		Parent: arg,
	}
	old, wasAdded = w.childOfComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
}

func (w *World) ChildOf(e Entity) (c ChildOfComponent, ok bool) {
	return w.childOfComponents.Data(e)
}

func (w *World) MutableChildOf(e Entity) (c *ChildOfComponent, ok bool) {
	return w.childOfComponents.DataMutable(e)
}

func (w *World) MustChildOf(e Entity) ChildOfComponent {
	c, ok := w.childOfComponents.Data(e)
	if !ok {
		panic("entity does not have ChildOf")
	}
	return c
}

func (w *World) RemoveChildOf(e Entity) {
	wasRemoved := w.childOfComponents.Remove(e)

	// depending on the generation flags, these might be unused
	_ = wasRemoved

}

func (w *World) HasChildOf(e Entity) bool {
	return w.childOfComponents.Contains(e)
}

func (w *World) AllChildOf(yield func(e Entity, c ChildOfComponent) bool) {
	for e, c := range w.childOfComponents.All {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutableChildOf(yield func(e Entity, c *ChildOfComponent) bool) {
	for e, c := range w.childOfComponents.AllMutable {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllChildOfEntities(yield func(e Entity) bool) {
	for e := range w.childOfComponents.AllEntities {
		if !yield(e) {
			break
		}
	}
}

// ChildOfBuilder

func WithChildOf(arg Entity) EntityBuilderOption {
	c := ChildOfComponent{
		Parent: arg,
	}

	return func(w *World, e Entity) {
		w.childOfComponents.Upsert(e, c)
	}
}

// Events

// Resource methods

func (w *World) SetChildOfResource(arg Entity) {
	w.SetChildOf(w.resourceEntity, arg)
}

func (w *World) ChildOfResource() (ChildOfComponent, bool) {
	return w.childOfComponents.Data(w.resourceEntity)
}

func (w *World) MustChildOfResource() ChildOfComponent {
	c, ok := w.ChildOfResource()
	if !ok {
		panic("resource entity does not have ChildOf")
	}
	return c
}

func (w *World) RemoveChildOfResource() {
	w.childOfComponents.Remove(w.resourceEntity)
}

func (w *World) HasChildOfResource() bool {
	return w.childOfComponents.Contains(w.resourceEntity)
}
