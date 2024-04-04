package {{.PackageName}}

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
	e.w.DestroyEntities(e)
}

func EntitiesToU32s(entities ...Entity) []uint32 {
	u32s := make([]uint32, len(entities))
	for i, e := range entities {
		u32s[i] = e.val
	}
	return u32s
}