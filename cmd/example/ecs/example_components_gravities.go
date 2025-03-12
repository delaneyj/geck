package ecs

type GravityComponent struct {
	G float32
}

func GravityComponentFromValues(
	gArg float32,
) GravityComponent {
	return GravityComponent{
		G: gArg,
	}
}

func DefaultGravityComponent() GravityComponent {
	return GravityComponent{
		G: -9.800000,
	}
}

func (c GravityComponent) Clone() GravityComponent {
	return GravityComponent{
		G: c.G,
	}
}

func (w *World) SetGravity(e Entity, arg float32) (old GravityComponent, wasAdded bool) {
	c := GravityComponent{
		G: arg,
	}
	old, wasAdded = w.gravityComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
}

func (w *World) Gravity(e Entity) (c GravityComponent, ok bool) {
	return w.gravityComponents.Data(e)
}

func (w *World) MutableGravity(e Entity) (c *GravityComponent, ok bool) {
	return w.gravityComponents.DataMutable(e)
}

func (w *World) MustGravity(e Entity) GravityComponent {
	c, ok := w.gravityComponents.Data(e)
	if !ok {
		panic("entity does not have Gravity")
	}
	return c
}

func (w *World) RemoveGravity(e Entity) {
	wasRemoved := w.gravityComponents.Remove(e)

	// depending on the generation flags, these might be unused
	_ = wasRemoved

}

func (w *World) HasGravity(e Entity) bool {
	return w.gravityComponents.Contains(e)
}

func (w *World) GravitiesCount() int {
	return w.gravityComponents.Len()
}

func (w *World) GravitiesCapacity() int {
	return w.gravityComponents.Cap()
}

func (w *World) AllGravities(yield func(e Entity, c GravityComponent) bool) {
	for e, c := range w.gravityComponents.All {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutableGravities(yield func(e Entity, c *GravityComponent) bool) {
	for e, c := range w.gravityComponents.AllMutable {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllGravitiesEntities(yield func(e Entity) bool) {
	for e := range w.gravityComponents.AllEntities {
		if !yield(e) {
			break
		}
	}
}

// GravityBuilder

func WithGravity(arg float32) EntityBuilderOption {
	c := GravityComponent{
		G: arg,
	}

	return func(w *World, e Entity) {
		w.gravityComponents.Upsert(e, c)
	}
}

// Events

// Resource methods

func (w *World) SetGravityResource(arg float32) {
	w.SetGravity(w.resourceEntity, arg)
}

func (w *World) GravityResource() (GravityComponent, bool) {
	return w.gravityComponents.Data(w.resourceEntity)
}

func (w *World) MustGravityResource() GravityComponent {
	c, ok := w.GravityResource()
	if !ok {
		panic("resource entity does not have Gravity")
	}
	return c
}

func (w *World) RemoveGravityResource() {
	w.gravityComponents.Remove(w.resourceEntity)
}

func (w *World) HasGravityResource() bool {
	return w.gravityComponents.Contains(w.resourceEntity)
}
