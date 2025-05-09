package ecs

type RotationComponent struct {
	X float32
	Y float32
	Z float32
	W float32
}

func RotationComponentFromValues(
	xArg float32,
	yArg float32,
	zArg float32,
	wArg float32,
) RotationComponent {
	return RotationComponent{
		X: xArg,
		Y: yArg,
		Z: zArg,
		W: wArg,
	}
}

func DefaultRotationComponent() RotationComponent {
	return RotationComponent{
		X: 0.000000,
		Y: 0.000000,
		Z: 0.000000,
		W: 1.000000,
	}
}

func (c RotationComponent) Clone() RotationComponent {
	return RotationComponent{
		X: c.X,
		Y: c.Y,
		Z: c.Z,
		W: c.W,
	}
}

func (w *World) SetRotation(e Entity, c RotationComponent) (old RotationComponent, wasAdded bool) {
	old, wasAdded = w.rotationComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
}

func (w *World) SetRotationFromValues(
	e Entity,
	xArg float32,
	yArg float32,
	zArg float32,
	wArg float32,
) {
	old, _ := w.SetRotation(e, RotationComponent{
		X: xArg,
		Y: yArg,
		Z: zArg,
		W: wArg,
	})

	// depending on the generation flags, these might be unused
	_ = old

}

func (w *World) Rotation(e Entity) (c RotationComponent, ok bool) {
	return w.rotationComponents.Data(e)
}

func (w *World) MutableRotation(e Entity) (c *RotationComponent, ok bool) {
	return w.rotationComponents.DataMutable(e)
}

func (w *World) MustMutableRotation(e Entity) *RotationComponent {
	c, ok := w.MutableRotation(e)
	if !ok {
		panic("entity does not have Rotation")
	}
	return c
}

func (w *World) MustRotation(e Entity) RotationComponent {
	c, ok := w.rotationComponents.Data(e)
	if !ok {
		panic("entity does not have Rotation")
	}
	return c
}

func (w *World) RemoveRotation(e Entity) {
	wasRemoved := w.rotationComponents.Remove(e)

	// depending on the generation flags, these might be unused
	_ = wasRemoved

}

func (w *World) HasRotation(e Entity) bool {
	return w.rotationComponents.Contains(e)
}

func (w *World) RotationsCount() int {
	return w.rotationComponents.Len()
}

func (w *World) RotationsCapacity() int {
	return w.rotationComponents.Cap()
}

func (w *World) AllRotations(yield func(e Entity, c RotationComponent) bool) {
	for e, c := range w.rotationComponents.All {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutableRotations(yield func(e Entity, c *RotationComponent) bool) {
	for e, c := range w.rotationComponents.AllMutable {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllRotationsEntities(yield func(e Entity) bool) {
	for e := range w.rotationComponents.AllEntities {
		if !yield(e) {
			break
		}
	}
}

func (w *World) AllMutableRotationsEntities(yield func(e Entity) bool) {
	w.AllRotationsEntities(yield)
}

// RotationBuilder
func WithRotationDefault() EntityBuilderOption {
	return WithRotation(DefaultRotationComponent())
}

func WithRotation(c RotationComponent) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.rotationComponents.Upsert(e, c)
	}
}

func WithRotationFromValues(
	xArg float32,
	yArg float32,
	zArg float32,
	wArg float32,
) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.SetRotationFromValues(e,
			xArg,
			yArg,
			zArg,
			wArg,
		)
	}
}

// Events

// Resource methods
func (w *World) SetRotationResource(c RotationComponent) {
	w.SetRotation(w.resourceEntity, c)
}

func (w *World) SetRotationResourceFromValues(
	xArg float32,
	yArg float32,
	zArg float32,
	wArg float32,
) {
	w.SetRotationResource(RotationComponent{
		X: xArg,
		Y: yArg,
		Z: zArg,
		W: wArg,
	})
}

func (w *World) RotationResource() (RotationComponent, bool) {
	return w.rotationComponents.Data(w.resourceEntity)
}

func (w *World) MustRotationResource() RotationComponent {
	c, ok := w.RotationResource()
	if !ok {
		panic("resource entity does not have Rotation")
	}
	return c
}

func (w *World) RemoveRotationResource() {
	w.rotationComponents.Remove(w.resourceEntity)
}

func (w *World) HasRotationResource() bool {
	return w.rotationComponents.Contains(w.resourceEntity)
}
