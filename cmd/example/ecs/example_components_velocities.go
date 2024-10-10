package ecs

type Velocity struct {
	X float32
	Y float32
	Z float32
}

func (w *World) SetVelocity(e Entity, c Velocity) (old Velocity, wasAdded bool) {
	old, wasAdded = w.velocityComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
}

func (w *World) SetVelocityFromValues(
	e Entity,
	xArg float32,
	yArg float32,
	zArg float32,
) {
	old, _ := w.SetVelocity(e, Velocity{
		X: xArg,
		Y: yArg,
		Z: zArg,
	})

	// depending on the generation flags, these might be unused
	_ = old

}

func (w *World) Velocity(e Entity) (c Velocity, ok bool) {
	return w.velocityComponents.Data(e)
}

func (w *World) MutableVelocity(e Entity) (c *Velocity, ok bool) {
	return w.velocityComponents.DataMutable(e)
}

func (w *World) MustVelocity(e Entity) Velocity {
	c, ok := w.velocityComponents.Data(e)
	if !ok {
		panic("entity does not have Velocity")
	}
	return c
}

func (w *World) RemoveVelocity(e Entity) {
	wasRemoved := w.velocityComponents.Remove(e)

	// depending on the generation flags, these might be unused
	_ = wasRemoved

}

func (w *World) HasVelocity(e Entity) bool {
	return w.velocityComponents.Contains(e)
}

func (w *World) AllVelocities(yield func(e Entity, c Velocity) bool) {
	for e, c := range w.velocityComponents.All {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutableVelocities(yield func(e Entity, c *Velocity) bool) {
	for e, c := range w.velocityComponents.AllMutable {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllVelocitiesEntities(yield func(e Entity) bool) {
	for e := range w.velocityComponents.AllEntities {
		if !yield(e) {
			break
		}
	}
}

// VelocityBuilder
func WithVelocity(c Velocity) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.velocityComponents.Upsert(e, c)
	}
}

func WithVelocityFromValues(
	xArg float32,
	yArg float32,
	zArg float32,
) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.SetVelocityFromValues(e,
			xArg,
			yArg,
			zArg,
		)
	}
}

// Events

// Resource methods
func (w *World) SetVelocityResource(c Velocity) {
	w.velocityComponents.Upsert(w.resourceEntity, c)
}

func (w *World) SetVelocityResourceFromValues(
	xArg float32,
	yArg float32,
	zArg float32,
) {
	w.SetVelocityResource(Velocity{
		X: xArg,
		Y: yArg,
		Z: zArg,
	})
}

func (w *World) VelocityResource() (Velocity, bool) {
	return w.velocityComponents.Data(w.resourceEntity)
}

func (w *World) MustVelocityResource() Velocity {
	c, ok := w.VelocityResource()
	if !ok {
		panic("resource entity does not have Velocity")
	}
	return c
}

func (w *World) RemoveVelocityResource() {
	w.velocityComponents.Remove(w.resourceEntity)
}
