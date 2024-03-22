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
	if !e.IsAlive() || e.IsResourceEntity() {
		return
	}

	{{range .Components -}}
	e.w.{{.Name.Plural.Camel}}Store.Remove(e)
	{{end }}

	e.w.liveEntitieIDs.Remove(e.val)
	bumped := e.UpdateVersion().val
	e.w.freeEntitieIDs.Add(bumped)

	fireEvent(e.w, EntityDestroyedEvent{e})
}
