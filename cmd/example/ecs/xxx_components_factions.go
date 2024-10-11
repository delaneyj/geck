package ecs

type FactionComponent struct {
	Entity Entity
}

func (w *World) SetFaction(e Entity, c FactionComponent) (old FactionComponent, wasAdded bool) {
	old, wasAdded = w.factionComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
}

func (w *World) SetFactionFromValues(
	e Entity,
	entityArg Entity,
) {
	old, _ := w.SetFaction(e, FactionComponent{
		Entity: entityArg,
	})

	// depending on the generation flags, these might be unused
	_ = old

}

func (w *World) Faction(e Entity) (c FactionComponent, ok bool) {
	return w.factionComponents.Data(e)
}

func (w *World) MutableFaction(e Entity) (c *FactionComponent, ok bool) {
	return w.factionComponents.DataMutable(e)
}

func (w *World) MustFaction(e Entity) FactionComponent {
	c, ok := w.factionComponents.Data(e)
	if !ok {
		panic("entity does not have Faction")
	}
	return c
}

func (w *World) RemoveFaction(e Entity) {
	wasRemoved := w.factionComponents.Remove(e)

	// depending on the generation flags, these might be unused
	_ = wasRemoved

}

func (w *World) HasFaction(e Entity) bool {
	return w.factionComponents.Contains(e)
}

func (w *World) AllFactions(yield func(e Entity, c FactionComponent) bool) {
	for e, c := range w.factionComponents.All {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutableFactions(yield func(e Entity, c *FactionComponent) bool) {
	for e, c := range w.factionComponents.AllMutable {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllFactionsEntities(yield func(e Entity) bool) {
	for e := range w.factionComponents.AllEntities {
		if !yield(e) {
			break
		}
	}
}

// FactionBuilder
func WithFaction(c FactionComponent) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.factionComponents.Upsert(e, c)
	}
}

func WithFactionFromValues(
	entityArg Entity,
) EntityBuilderOption {
	return func(w *World, e Entity) {
		w.SetFactionFromValues(e,
			entityArg,
		)
	}
}

// Events

// Resource methods
func (w *World) SetFactionResource(c FactionComponent) {
	w.SetFaction(w.resourceEntity, c)
}

func (w *World) SetFactionResourceFromValues(
	entityArg Entity,
) {
	w.SetFactionResource(FactionComponent{
		Entity: entityArg,
	})
}

func (w *World) FactionResource() (FactionComponent, bool) {
	return w.factionComponents.Data(w.resourceEntity)
}

func (w *World) MustFactionResource() FactionComponent {
	c, ok := w.FactionResource()
	if !ok {
		panic("resource entity does not have Faction")
	}
	return c
}

func (w *World) RemoveFactionResource() {
	w.factionComponents.Remove(w.resourceEntity)
}
