package ecs

type ChildOf struct {
	Parent Entity
}

func (w *World) SetChildOf(e Entity, c ChildOf) (old ChildOf, wasAdded bool) {
	old, wasAdded = w.childOfComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
}

func (w *World) SetChildOfFromValues(
	e Entity,
	parentArg Entity,
) {
	old, _ := w.SetChildOf(e, ChildOf{
		Parent: parentArg,
	})

	// depending on the generation flags, these might be unused
	_ = old

}

func (w *World) ChildOf(e Entity) (c ChildOf, ok bool) {
	return w.childOfComponents.Data(e)
}

func (w *World) MutableChildOf(e Entity) (c *ChildOf, ok bool) {
	return w.childOfComponents.DataMutable(e)
}

func (w *World) MustChildOf(e Entity) ChildOf {
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

func (w *World) AllChildOf(yield func(e Entity, c ChildOf) bool) {
	for e, c := range w.childOfComponents.All {
		if yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutableChildOf(yield func(e Entity, c *ChildOf) bool) {
	for e, c := range w.childOfComponents.AllMutable {
		if yield(e, c) {
			break
		}
	}
}

func (w *World) AllChildOfEntities(yield func(e Entity) bool) {
	for e := range w.childOfComponents.AllEntities {
		if yield(e) {
			break
		}
	}
}

// ChildOfBuilder
func WithChildOf(c ChildOf) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.childOfComponents.Upsert(e, c)
	}
}

func WithChildOfFromValues(
	parentArg Entity,
) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.SetChildOfFromValues(e,
			parentArg,
		)
	}
}

// Events

// Resource methods
func (w *World) SetChildOfResource(c ChildOf) {
	w.childOfComponents.Upsert(w.resourceEntity, c)
}

func (w *World) SetChildOfResourceFromValues(
	parentArg Entity,
) {
	w.SetChildOfResource(ChildOf{
		Parent: parentArg,
	})
}

func (w *World) ChildOfResource() (ChildOf, bool) {
	return w.childOfComponents.Data(w.resourceEntity)
}

func (w *World) MustChildOfResource() ChildOf {
	c, ok := w.ChildOfResource()
	if !ok {
		panic("resource entity does not have ChildOf")
	}
	return c
}

func (w *World) RemoveChildOfResource() {
	w.childOfComponents.Remove(w.resourceEntity)
}
