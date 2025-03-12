package ecs

type FactionComponent struct {
	Entity Entity
}

func FactionComponentFromValues(
	entityArg Entity,
) FactionComponent {
	return FactionComponent{
		Entity: entityArg,
	}
}

func DefaultFactionComponent() FactionComponent {
	return FactionComponent{
		Entity: EntityFromU32(0),
	}
}

func (c FactionComponent) Clone() FactionComponent {
	return FactionComponent{
		Entity: c.Entity,
	}
}

func (w *World) SetFaction(e Entity, arg Entity) (old FactionComponent, wasAdded bool) {
	c := FactionComponent{
		Entity: arg,
	}
	old, wasAdded = w.factionComponents.Upsert(e, c)

	// depending on the generation flags, these might be unused
	_, _ = old, wasAdded

	return old, wasAdded
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

func (w *World) FactionsCount() int {
	return w.factionComponents.Len()
}

func (w *World) FactionsCapacity() int {
	return w.factionComponents.Cap()
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

func WithFaction(arg Entity) EntityBuilderOption {
	c := FactionComponent{
		Entity: arg,
	}

	return func(w *World, e Entity) {
		w.factionComponents.Upsert(e, c)
	}
}

// Events

// Resource methods

func (w *World) SetFactionResource(arg Entity) {
	w.SetFaction(w.resourceEntity, arg)
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

func (w *World) HasFactionResource() bool {
	return w.factionComponents.Contains(w.resourceEntity)
}
