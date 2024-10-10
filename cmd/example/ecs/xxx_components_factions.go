package ecs

type Faction struct {
	Entity Entity
}

func (w *World) SetFaction(e Entity, c Faction) (old Faction, wasAdded bool) {
	old, wasAdded = w.factionComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
}

func (w *World) SetFactionFromValues(
	e Entity,
	entityArg Entity,
) {
	old, _ := w.SetFaction(e, Faction{
		Entity: entityArg,
	})

	// depending on the generation flags, these might be unused
	_ = old

}

func (w *World) Faction(e Entity) (c Faction, ok bool) {
	return w.factionComponents.Data(e)
}

func (w *World) MutableFaction(e Entity) (c *Faction, ok bool) {
	return w.factionComponents.DataMutable(e)
}

func (w *World) MustFaction(e Entity) Faction {
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

func (w *World) AllFactions(yield func(e Entity, c Faction) bool) {
	for e, c := range w.factionComponents.All {
		if !yield(e, c) {
			break
		}
	}
}

func (w *World) AllMutableFactions(yield func(e Entity, c *Faction) bool) {
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
func WithFaction(c Faction) EntityBuilderOption {
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
func (w *World) SetFactionResource(c Faction) {
	w.factionComponents.Upsert(w.resourceEntity, c)
}

func (w *World) SetFactionResourceFromValues(
	entityArg Entity,
) {
	w.SetFactionResource(Faction{
		Entity: entityArg,
	})
}

func (w *World) FactionResource() (Faction, bool) {
	return w.factionComponents.Data(w.resourceEntity)
}

func (w *World) MustFactionResource() Faction {
	c, ok := w.FactionResource()
	if !ok {
		panic("resource entity does not have Faction")
	}
	return c
}

func (w *World) RemoveFactionResource() {
	w.factionComponents.Remove(w.resourceEntity)
}
