package {{.PackageName}}

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/RoaringBitmap/roaring"
	"github.com/btvoidx/mint"
	ecspb "github.com/delaneyj/geck/cmd/example/ecs/pb/gen/ecs/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

var empty = &emptypb.Empty{}

type System interface {
	Name() string
	ReliesOn() []string
	Initialize(w *World) error
	Tick(ctx context.Context, w *World) error
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
	eventBus *mint.Emitter

	nextSystemID                                   	uint32
	systems, leftToRun, notRunWithDependenciesDone 	map[uint32]*systemRunner
	tickWaitGroup                                  	*sync.WaitGroup
	tickCount                   					int

	nextEntityID                   uint32
	liveEntitieIDs, freeEntitieIDs *roaring.Bitmap
	{{range .Components -}}
	{{.Name.Plural.Camel}}Store *SparseSet[{{.Name.Singular.Pascal}}]
	{{end -}}

	{{range .ComponentSets}}
	{{.Name.Singular.Pascal}} *{{.Name.Singular.Pascal}}
	{{end}}

	patch *ecspb.WorldPatch
}

func NewWorld() *World {
	w := &World{
		eventBus: &mint.Emitter{},
		nextSystemID:               1,
		systems:                    map[uint32]*systemRunner{},
		leftToRun:                  map[uint32]*systemRunner{},
		notRunWithDependenciesDone: map[uint32]*systemRunner{},
		tickWaitGroup:              &sync.WaitGroup{},
		tickCount:                  0,

		nextEntityID: 1,
		liveEntitieIDs: roaring.NewBitmap(),
		freeEntitieIDs: roaring.NewBitmap(),
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
		if err := sr.system.Initialize(w); err != nil {
			return fmt.Errorf("system %s failed to initialize: %s", s.Name(), err)
		}
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

func (w *World) Tick(ctx context.Context) error {
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
				if err := sr.system.Tick(ctx, w); err != nil {
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

	w.patch.Entities[e.val] = empty

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

func (w *World) EntitiesFromU32s(vals ...uint32) (entities []Entity) {
	entities = make([]Entity, len(vals))
	for i, val := range vals {
		e := Entity{w: w, val: val}
		if !e.IsAlive() {
			fireEvent(w, EntityCreatedEvent{e})
		}
		entities[i] = e
		w.freeEntitieIDs.Remove(val)
		w.liveEntitieIDs.Add(val)

		w.patch.Entities[val] = empty
	}
	return entities
}

func (w *World) Entities(count int) []Entity {
	entities := make([]Entity, count)
	for i := 0; i < count; i++ {
		entities[i] = w.Entity()
	}
	return entities
}

func (w *World) DestroyEntities(es ...Entity) {
	{{- range .Components}}
	w.{{.Name.Plural.Camel}}Store.Remove(es...)
	{{- end}}

	for _, e := range es {
		if !e.IsAlive() {
			continue
		}

		fireEvent(w, EntityDestroyedEvent{e})
		w.liveEntitieIDs.Remove(e.val)
		w.freeEntitieIDs.Add(e.val)

		w.patch.Entities[e.val] = nil
	}
}


func(w *World) Reset(){
	{{- range .Components}}
	w.{{.Name.Plural.Camel}}Store.Clear()
	{{- end}}

	liveEntitieIDs := w.liveEntitieIDs.ToArray()
	w.liveEntitieIDs.Clear()
	w.freeEntitieIDs.Clear()

	for _, id := range liveEntitieIDs {
		e := w.EntityFromU32(id)
		fireEvent(w, EntityDestroyedEvent{e})
	}
	ResetWorldPatch(w.patch)
}

func NewWorldPatch() *ecspb.WorldPatch {
	return &{{.PackageName}}pb.WorldPatch{
		Entities: map[uint32]*emptypb.Empty{},
		{{range .Components -}}
		{{if .IsTag -}}
			{{.Name.Singular.Pascal}}Tags: map[uint32]*emptypb.Empty{},
		{{else -}}
			{{.Name.Singular.Pascal}}Components: map[uint32]*{{.PackageName}}pb.{{.Name.Singular.Pascal}}Component{},
		{{end -}}
		{{end -}}
	}
}

func ResetWorldPatch(patch *ecspb.WorldPatch) *ecspb.WorldPatch {
	clear(patch.Entities)
	{{range .Components -}}
	{{if .IsTag -}}
	clear(patch.{{.Name.Singular.Pascal}}Tags)
	{{else -}}
	clear(patch.{{.Name.Singular.Pascal}}Components)
	{{end -}}
	{{end -}}
	return patch
}

func MergeWorldWriteAheadLogs(patchs ...*ecspb.WorldPatch) *ecspb.WorldPatch {
	merged := NewWorldPatch()
	for _, patch := range patchs {
		for k,v := range patch.Entities {
			merged.Entities[k] = v
		}

		{{range .Components -}}
		// merge {{.Name.Plural.Pascal}}
		{{if .IsTag -}}
		for k,v := range patch.{{.Name.Singular.Pascal}}Tags {
			merged.{{.Name.Singular.Pascal}}Tags[k] = v
		{{else -}}
		for k,v := range patch.{{.Name.Singular.Pascal}}Components {
			merged.{{.Name.Singular.Pascal}}Components[k] = v
		{{end -}}
		}

		{{end }}
	}

	return merged
}

func(w *World) ApplyPatches(patches ...*ecspb.WorldPatch){
	for _, patch := range patches {
		for k,v := range patch.Entities {
			if v == nil {
				w.DestroyEntities(w.EntityFromU32(k))
			} else {
				w.EntityFromU32(k)
			}
		}

		{{range .Components -}}
		// apply {{.Name.Plural.Pascal}}
		{{if .IsTag -}}
		for val,c := range patch.{{.Name.Singular.Pascal}}Tags {
		{{else -}}
		for val,c := range patch.{{.Name.Singular.Pascal}}Components {
		{{end -}}
			e := w.EntityFromU32(val)
			w.Apply{{.Name.Singular.Pascal}}Patch(e,c)
		}
		{{end }}
	}
}