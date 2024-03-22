package ecs

import (
	"errors"
	"math"
)

const (
	EntityIndexBits   = 20
	EntityIndexMask   = 1<<EntityIndexBits - 1
	EntityVersionBits = 32 - EntityIndexBits
	EntityVersionMask = math.MaxUint32 ^ EntityIndexMask
	DeadEntityID      = EntityIndexMask
)

var (
	ErrEntityVersionMismatch = errors.New("entity version mismatch")
)

type Entity struct {
	w   *World
	val uint32
}

func (e Entity) World() *World {
	return e.w
}

func (e Entity) Index() int {
	return int(e.val & EntityIndexMask)
}

func (e Entity) IndexU32() uint32 {
	return e.val & EntityIndexMask
}

func (e Entity) Version() uint32 {
	return (e.val & EntityVersionMask) >> EntityIndexBits
}

func (e Entity) UpdateVersion() Entity {
	id := e.Index()
	version := e.Version()

	updated := uint32(id) + ((version + 1) << EntityIndexBits)
	return Entity{w: e.w, val: updated}
}

func (e Entity) Raw() uint32 {
	return e.val
}

func (e Entity) IsAlive() bool {
	return e.w.liveEntitieIDs.Contains(e.val)
}

func (e Entity) IsResourceEntity() bool {
	return e.val == e.w.resourceEntity.val
}

func (e Entity) Destroy() {
	if !e.IsAlive() || e.IsResourceEntity() {
		return
	}

	e.w.namesStore.Remove(e)
	e.w.childOfStore.Remove(e)
	e.w.isAStore.Remove(e)
	e.w.positionsStore.Remove(e)
	e.w.velocitiesStore.Remove(e)
	e.w.rotationsStore.Remove(e)
	e.w.eatsStore.Remove(e)
	e.w.likesStore.Remove(e)
	e.w.enemyStore.Remove(e)
	e.w.growsStore.Remove(e)
	e.w.gravitiesStore.Remove(e)
	e.w.spaceshipStore.Remove(e)
	e.w.spacestationStore.Remove(e)
	e.w.factionsStore.Remove(e)
	e.w.dockedTosStore.Remove(e)
	e.w.planetStore.Remove(e)
	e.w.ruledBysStore.Remove(e)
	e.w.alliedWithsStore.Remove(e)

	e.w.liveEntitieIDs.Remove(e.val)
	bumped := e.UpdateVersion().val
	e.w.freeEntitieIDs.Add(bumped)

	fireEvent(e.w, EntityDestroyedEvent{e})
}
