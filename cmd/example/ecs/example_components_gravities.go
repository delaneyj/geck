package ecs

type Gravity struct {
	G float32
}

func (w *World) SetGravity(e Entity, c Gravity) (old Gravity, wasAdded bool) {
	old, wasAdded = w.gravityComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
}

func (w *World) SetGravityFromValues(
	e Entity,
	gArg float32,
) {
	old, _ := w.SetGravity(e, Gravity{
		G: gArg,
	})

	// depending on the generation flags, these might be unused
	_ = old

}

func (w *World) Gravity(e Entity) (c Gravity, ok bool) {
	return w.gravityComponents.Data(e)
}

func (w *World) MutableGravity(e Entity) (c *Gravity, ok bool) {
	return w.gravityComponents.DataMutable(e)
}

func (w *World) MustGravity(e Entity) Gravity {
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

func (w *World) AllGravities(yield func(e Entity, c Gravity) bool) {
	for e, c := range w.gravityComponents.All {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutableGravities(yield func(e Entity, c *Gravity) bool) {
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
func WithGravity(c Gravity) EntityBuilderOption {
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
func (w *World) SetGravityResource(c Gravity) {
	w.gravityComponents.Upsert(w.resourceEntity, c)
}

func (w *World) SetGravityResourceFromValues(
	gArg float32,
) {
	w.SetGravityResource(Gravity{
		G: gArg,
	})
}

func (w *World) GravityResource() (Gravity, bool) {
	return w.gravityComponents.Data(w.resourceEntity)
}

func (w *World) MustGravityResource() Gravity {
	c, ok := w.GravityResource()
	if !ok {
		panic("resource entity does not have Gravity")
	}
	return c
}

func (w *World) RemoveGravityResource() {
	w.gravityComponents.Remove(w.resourceEntity)
}
