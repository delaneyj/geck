package ecs

type NameComponent struct {
	Value string
}

func NameComponentFromValues(
	valueArg string,
) NameComponent {
	return NameComponent{
		Value: valueArg,
	}
}

func DefaultNameComponent() NameComponent {
	return NameComponent{
		Value: "",
	}
}

func (c NameComponent) Clone() NameComponent {
	return NameComponent{
		Value: c.Value,
	}
}

func (w *World) SetName(e Entity, arg string) (old NameComponent, wasAdded bool) {
	c := NameComponent{
		Value: arg,
	}
	old, wasAdded = w.nameComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
}

func (w *World) Name(e Entity) (c NameComponent, ok bool) {
	return w.nameComponents.Data(e)
}

func (w *World) MutableName(e Entity) (c *NameComponent, ok bool) {
	return w.nameComponents.DataMutable(e)
}

func (w *World) MustName(e Entity) NameComponent {
	c, ok := w.nameComponents.Data(e)
	if !ok {
		panic("entity does not have Name")
	}
	return c
}

func (w *World) RemoveName(e Entity) {
	wasRemoved := w.nameComponents.Remove(e)

	// depending on the generation flags, these might be unused
	_ = wasRemoved

}

func (w *World) HasName(e Entity) bool {
	return w.nameComponents.Contains(e)
}

func (w *World) NamesCount() int {
	return w.nameComponents.Len()
}

func (w *World) NamesCapacity() int {
	return w.nameComponents.Cap()
}

func (w *World) AllNames(yield func(e Entity, c NameComponent) bool) {
	for e, c := range w.nameComponents.All {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutableNames(yield func(e Entity, c *NameComponent) bool) {
	for e, c := range w.nameComponents.AllMutable {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllNamesEntities(yield func(e Entity) bool) {
	for e := range w.nameComponents.AllEntities {
		if !yield(e) {
			break
		}
	}
}

// NameBuilder

func WithName(arg string) EntityBuilderOption {
	c := NameComponent{
		Value: arg,
	}

	return func(w *World, e Entity) {
		w.nameComponents.Upsert(e, c)
	}
}

// Events

// Resource methods

func (w *World) SetNameResource(arg string) {
	w.SetName(w.resourceEntity, arg)
}

func (w *World) NameResource() (NameComponent, bool) {
	return w.nameComponents.Data(w.resourceEntity)
}

func (w *World) MustNameResource() NameComponent {
	c, ok := w.NameResource()
	if !ok {
		panic("resource entity does not have Name")
	}
	return c
}

func (w *World) RemoveNameResource() {
	w.nameComponents.Remove(w.resourceEntity)
}

func (w *World) HasNameResource() bool {
	return w.nameComponents.Contains(w.resourceEntity)
}
