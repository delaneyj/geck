package ecs

import "github.com/btvoidx/mint"

type EntitiesCreatedEvent struct {
	Entities []Entity
}

type EntitiesDestroyedEvent struct {
	Entities []Entity
}

type UnsubscribeFunc func()

func (w *World) OnEntitiesCreated(fn func(EntitiesCreatedEvent)) UnsubscribeFunc {
	stopCh := mint.On(w.eventBus, fn)
	return func() { stopCh() }
}

func (w *World) OnEntitiesDestroyed(fn func(EntitiesDestroyedEvent)) UnsubscribeFunc {
	stopCh := mint.On(w.eventBus, fn)
	return func() { stopCh() }
}

func fireEvent[T any](w *World, event T) {
	mint.Emit(w.eventBus, event)
}
