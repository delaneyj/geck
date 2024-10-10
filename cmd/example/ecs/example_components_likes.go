package ecs

type Likes struct {
	Entity []Entity
}

func (w *World) SetLikes(e Entity, c Likes) (old Likes, wasAdded bool) {
	old, wasAdded = w.likesComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
}

func (w *World) SetLikesFromValues(
	e Entity,
	entityArg []Entity,
) {
	old, _ := w.SetLikes(e, Likes{
		Entity: entityArg,
	})

	// depending on the generation flags, these might be unused
	_ = old

}

func (w *World) Likes(e Entity) (c Likes, ok bool) {
	return w.likesComponents.Data(e)
}

func (w *World) MutableLikes(e Entity) (c *Likes, ok bool) {
	return w.likesComponents.DataMutable(e)
}

func (w *World) MustLikes(e Entity) Likes {
	c, ok := w.likesComponents.Data(e)
	if !ok {
		panic("entity does not have Likes")
	}
	return c
}

func (w *World) RemoveLikes(e Entity) {
	wasRemoved := w.likesComponents.Remove(e)

	// depending on the generation flags, these might be unused
	_ = wasRemoved

}

func (w *World) HasLikes(e Entity) bool {
	return w.likesComponents.Contains(e)
}

func (w *World) AllLikes(yield func(e Entity, c Likes) bool) {
	for e, c := range w.likesComponents.All {
		if yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutableLikes(yield func(e Entity, c *Likes) bool) {
	for e, c := range w.likesComponents.AllMutable {
		if yield(e, c) {
			break
		}
	}
}

func (w *World) AllLikesEntities(yield func(e Entity) bool) {
	for e := range w.likesComponents.AllEntities {
		if yield(e) {
			break
		}
	}
}

// LikesBuilder
func WithLikes(c Likes) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.likesComponents.Upsert(e, c)
	}
}

func WithLikesFromValues(
	entityArg []Entity,
) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.SetLikesFromValues(e,
			entityArg,
		)
	}
}

// Events

// Resource methods
func (w *World) SetLikesResource(c Likes) {
	w.likesComponents.Upsert(w.resourceEntity, c)
}

func (w *World) SetLikesResourceFromValues(
	entityArg []Entity,
) {
	w.SetLikesResource(Likes{
		Entity: entityArg,
	})
}

func (w *World) LikesResource() (Likes, bool) {
	return w.likesComponents.Data(w.resourceEntity)
}

func (w *World) MustLikesResource() Likes {
	c, ok := w.LikesResource()
	if !ok {
		panic("resource entity does not have Likes")
	}
	return c
}

func (w *World) RemoveLikesResource() {
	w.likesComponents.Remove(w.resourceEntity)
}
