package {{.PackageName}}

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/RoaringBitmap/roaring"
	"github.com/btvoidx/mint"
)

type System interface {
	Name() string
	ReliesOn() []string
	Tick(w *World) error
}

type systemRunner struct {
	id                       uint32
	w                        *World
	system                   System
	waitingOnTmpl, waitingOn map[uint32]*systemRunner
	hasRun, isDisabled       bool
}

type World struct {
	zeroEntity, resourceEntity, deadEntity Entity

	// maxEntity  Entity
	nextEntityID   uint32
	liveEntitieIDs *roaring.Bitmap
	freeEntitieIDs *roaring.Bitmap

	eventBus *mint.Emitter

	nextSystemID                                   	uint32
	systems, leftToRun, notRunWithDependenciesDone 	map[uint32]*systemRunner
	tickWaitGroup                                  	*sync.WaitGroup
	tickCount                   					int

	{{range .Components -}}
	{{.Name.Plural.Camel}}Store *SparseSet[{{.Name.Singular.Pascal}}]
	{{end -}}

	{{range .ComponentSets}}
	{{.Name.Singular.Pascal}} *{{.Name.Singular.Pascal}}
	{{end}}
}

func NewWorld() *World {
	w := &World{
		liveEntitieIDs: roaring.NewBitmap(),
		freeEntitieIDs: roaring.NewBitmap(),
		eventBus: &mint.Emitter{},

		nextSystemID:               1,
		systems:                    map[uint32]*systemRunner{},
		leftToRun:                  map[uint32]*systemRunner{},
		notRunWithDependenciesDone: map[uint32]*systemRunner{},
		tickWaitGroup:              &sync.WaitGroup{},
		tickCount:                  0,

		{{range .Components -}}
		{{.Name.Plural.Camel}}Store : NewSparseSet[{{.Name.Singular.Pascal}}](nil),
		{{end }}
	}

	// setup built-in entities
	w.zeroEntity = w.Entity()
	w.resourceEntity = w.Entity()
	w.deadEntity = w.EntityFromU32(DeadEntityID)

	// component sets
	{{range .ComponentSets -}}
	w.{{.Name.Singular.Pascal}} = New{{.Name.Singular.Pascal}}(w)
	{{end}}

	return w
}

//# region Systems
func (w *World) AddSystems(ss ... System) (err error) {
	for _, s := range ss {
		alreadyRegistered := false
		for _, sys := range w.systems {
			if sys.system.Name() == s.Name() {
				alreadyRegistered = true
				break
			}
		}
		if alreadyRegistered {
			return fmt.Errorf("system %s has already been added", s.Name())
		}

		sr := &systemRunner{
			id:            w.nextSystemID,
			w:             w,
			system:        s,
			waitingOnTmpl: map[uint32]*systemRunner{},
		}
		for _, r := range s.ReliesOn() {
			var dependentSystem *systemRunner
			for _, sys := range w.systems {
				if sys.system.Name() == r {
					dependentSystem = sys
					break
				}
			}
			if dependentSystem == nil {
				return fmt.Errorf(
					"system %s relies on %s, but %s has not been added",
					s.Name(), r, r,
				)
			}

			sr.waitingOnTmpl[dependentSystem.id] = dependentSystem
		}
		sr.waitingOn = map[uint32]*systemRunner{}
		for k, v := range sr.waitingOnTmpl {
			sr.waitingOn[k] = v
		}
		w.systems[sr.id] = sr
		w.nextSystemID++
	}
	return nil
}

func (w *World) RemoveSystems(ss ... System) error {
	for _, sys := range ss {
		name := sys.Name()
		var found *systemRunner
		for _, sr := range w.systems {
			if name == sr.system.Name() {
				found = sr
				break
			}
		}
		if found == nil {
			return fmt.Errorf("system %s not found", name)
		}

		reliedOnBy := []System{}
		for id, sr := range w.systems {
			if found.id == id {
				reliedOnBy = append(reliedOnBy, sr.system)
			}
		}

		if len(reliedOnBy) > 0 {
			names := []string{}
			for _, s := range reliedOnBy {
				names = append(names, s.Name())
			}

			return fmt.Errorf(
				"system %s is relied on by %s, and cannot be removed",
				name, strings.Join(names, ","),
			)
		}

		delete(w.systems, found.id)
	}

	return nil
}

func (w *World) Tick() error {
	// fill leftToRun
	for _, sr := range w.systems {
		if !sr.isDisabled {
			w.leftToRun[sr.id] = sr
		}
	}

	for len(w.leftToRun) > 0 {
		for _, sr := range w.leftToRun {
			if !sr.hasRun && len(sr.waitingOn) == 0 {
				w.notRunWithDependenciesDone[sr.id] = sr
			}
		}

		toRunConcurrentlyCount := len(w.notRunWithDependenciesDone)
		w.tickWaitGroup.Add(toRunConcurrentlyCount)
		for _, sr := range w.notRunWithDependenciesDone {
			go func(sr *systemRunner) {
				defer w.tickWaitGroup.Done()
				if err := sr.system.Tick(w); err != nil {
					log.Printf("system %s failed: %s", sr.system.Name(), err)
				}
				sr.hasRun = true
			}(sr)
		}
		w.tickWaitGroup.Wait()

		for _, ranSR := range w.notRunWithDependenciesDone {
			for _, sr := range w.leftToRun {
				delete(sr.waitingOn, ranSR.id)
			}
			delete(w.leftToRun, ranSR.id)
		}
	}

	// reset for next tick
	clear(w.leftToRun)
	clear(w.notRunWithDependenciesDone)
	for _, sr := range w.systems {
		for k, v := range sr.waitingOnTmpl {
			sr.waitingOn[k] = v
		}
		sr.hasRun = false
	}
	w.tickCount++

	return nil
}

func (w *World) DisableSystem(ss ... System) error {
	for _, sys := range ss {
		name := sys.Name()
		var found *systemRunner
		for _, sr := range w.systems {
			if name == sr.system.Name() {
				found = sr
				break
			}
		}
		if found == nil {
			return fmt.Errorf("system %s not found", name)
		}

		found.isDisabled = true
	}

	return nil
}

func (w *World) EnableSystem(ss ... System) error {
	for _, sys := range ss {
		name := sys.Name()
		var found *systemRunner
		for _, sr := range w.systems {
			if name == sr.system.Name() {
				found = sr
				break
			}
		}
		if found == nil {
			return fmt.Errorf("system %s not found", name)
		}

		found.isDisabled = false
	}

	return nil
}

func (w *World) TickCount() int {
	return w.tickCount
}

//# endregion

func (w *World) Entity() (e Entity) {
	e.w = w

	if w.freeEntitieIDs.IsEmpty() {
		e.val = w.nextEntityID
		w.nextEntityID++
	} else {
		last := w.freeEntitieIDs.Maximum()
		e.val = last
		w.freeEntitieIDs.Remove(last)
	}

	w.liveEntitieIDs.Add(e.val)
	fireEvent(w, EntityCreatedEvent{e})

	return e
}

func (w *World) EntityWithName(name string) Entity {
	return w.Entity().SetName(Name(name))
}

func (w *World) EntityFromU32(val uint32) Entity {
	e := Entity{w: w, val: val}
	if e.IsAlive() {
		return e
	}

	w.freeEntitieIDs.Remove(val)
	w.liveEntitieIDs.Add(val)
	fireEvent(w, EntityCreatedEvent{e})

	return e
}

func (w *World) Entities(count int) []Entity {
	entities := make([]Entity, count)
	for i := 0; i < count; i++ {
		entities[i] = w.Entity()
	}
	return entities
}

func (w *World) Reset() {
	{{range .Components -}}
	w.{{.Name.Plural.Camel}}Store.Clear()
	{{end }}

	iter := w.liveEntitieIDs.Iterator()
	for iter.HasNext() {
		id := iter.Next()
		e := w.EntityFromU32(id)
		fireEvent(w, EntityDestroyedEvent{e})
	}

	w.liveEntitieIDs.Clear()
	w.freeEntitieIDs.Clear()
}




