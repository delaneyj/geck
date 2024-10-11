package ecs

type GravityComponent struct {
	G float32
}

func (w *World) SetGravity(e Entity, c GravityComponent) (old GravityComponent, wasAdded bool) {
	old, wasAdded = w.gravityComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
}

func (w *World) SetGravityFromValues(
	e Entity,
	gArg float32,
) {
	old, _ := w.SetGravity(e, GravityComponent{
		G: gArg,
	})

	// depending on the generation flags, these might be unused
	_ = old

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
func WithGravity(c GravityComponent) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.gravityComponents.Upsert(e, c)
	}
}

func WithGravityFromValues(
	gArg float32,
) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.SetGravityFromValues(e,
			gArg,
		)
	}
}

// Events

// Resource methods
func (w *World) SetGravityResource(c GravityComponent) {
	w.SetGravity(w.resourceEntity, c)
}

func (w *World) SetGravityResourceFromValues(
	gArg float32,
) {
	w.SetGravityResource(GravityComponent{
		G: gArg,
	})
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
