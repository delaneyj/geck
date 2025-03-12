package ecs

type VelocityComponent struct {
	X float32
	Y float32
	Z float32
}

func VelocityComponentFromValues(
	xArg float32,
	yArg float32,
	zArg float32,
) VelocityComponent {
	return VelocityComponent{
		X: xArg,
		Y: yArg,
		Z: zArg,
	}
}

func DefaultVelocityComponent() VelocityComponent {
	return VelocityComponent{
		X: 0.000000,
		Y: 0.000000,
		Z: 0.000000,
	}
}

func (c VelocityComponent) Clone() VelocityComponent {
	return VelocityComponent{
		X: c.X,
		Y: c.Y,
		Z: c.Z,
	}
}

func (w *World) SetVelocity(e Entity, c VelocityComponent) (old VelocityComponent, wasAdded bool) {
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
	old, _ := w.SetVelocity(e, VelocityComponent{
		X: xArg,
		Y: yArg,
		Z: zArg,
	})

	// depending on the generation flags, these might be unused
	_ = old

}

func (w *World) Velocity(e Entity) (c VelocityComponent, ok bool) {
	return w.velocityComponents.Data(e)
}

func (w *World) MutableVelocity(e Entity) (c *VelocityComponent, ok bool) {
	return w.velocityComponents.DataMutable(e)
}

func (w *World) MustVelocity(e Entity) VelocityComponent {
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

func (w *World) VelocitiesCount() int {
	return w.velocityComponents.Len()
}

func (w *World) VelocitiesCapacity() int {
	return w.velocityComponents.Cap()
}

func (w *World) AllVelocities(yield func(e Entity, c VelocityComponent) bool) {
	for e, c := range w.velocityComponents.All {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutableVelocities(yield func(e Entity, c *VelocityComponent) bool) {
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

func WithVelocity(c VelocityComponent) EntityBuilderOption {

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

func (w *World) SetVelocityResource(c VelocityComponent) {
	w.SetVelocity(w.resourceEntity, c)
}

func (w *World) SetVelocityResourceFromValues(
	xArg float32,
	yArg float32,
	zArg float32,
) {
	w.SetVelocityResource(VelocityComponent{
		X: xArg,
		Y: yArg,
		Z: zArg,
	})
}

func (w *World) VelocityResource() (VelocityComponent, bool) {
	return w.velocityComponents.Data(w.resourceEntity)
}

func (w *World) MustVelocityResource() VelocityComponent {
	c, ok := w.VelocityResource()
	if !ok {
		panic("resource entity does not have Velocity")
	}
	return c
}

func (w *World) RemoveVelocityResource() {
	w.velocityComponents.Remove(w.resourceEntity)
}

func (w *World) HasVelocityResource() bool {
	return w.velocityComponents.Contains(w.resourceEntity)
}
