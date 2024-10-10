package ecs

type Rotation struct {
	X float32
	Y float32
	Z float32
	W float32
}

func (w *World) SetRotation(e Entity, c Rotation) (old Rotation, wasAdded bool) {
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
	old, _ := w.SetRotation(e, Rotation{
		X: xArg,
		Y: yArg,
		Z: zArg,
		W: wArg,
	})

	// depending on the generation flags, these might be unused
	_ = old

}

func (w *World) Rotation(e Entity) (c Rotation, ok bool) {
	return w.rotationComponents.Data(e)
}

func (w *World) MutableRotation(e Entity) (c *Rotation, ok bool) {
	return w.rotationComponents.DataMutable(e)
}

func (w *World) MustRotation(e Entity) Rotation {
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

func (w *World) AllRotations(yield func(e Entity, c Rotation) bool) {
	for e, c := range w.rotationComponents.All {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutableRotations(yield func(e Entity, c *Rotation) bool) {
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

// RotationBuilder
func WithRotation(c Rotation) EntityBuilderOption {
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
func (w *World) SetRotationResource(c Rotation) {
	w.rotationComponents.Upsert(w.resourceEntity, c)
}

func (w *World) SetRotationResourceFromValues(
	xArg float32,
	yArg float32,
	zArg float32,
	wArg float32,
) {
	w.SetRotationResource(Rotation{
		X: xArg,
		Y: yArg,
		Z: zArg,
		W: wArg,
	})
}

func (w *World) RotationResource() (Rotation, bool) {
	return w.rotationComponents.Data(w.resourceEntity)
}

func (w *World) MustRotationResource() Rotation {
	c, ok := w.RotationResource()
	if !ok {
		panic("resource entity does not have Rotation")
	}
	return c
}

func (w *World) RemoveRotationResource() {
	w.rotationComponents.Remove(w.resourceEntity)
}
