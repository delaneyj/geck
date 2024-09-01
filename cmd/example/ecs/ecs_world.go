package ecs

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/RoaringBitmap/roaring"
	"github.com/btvoidx/mint"
	ecspb "github.com/delaneyj/geck/cmd/example/ecs/pb/gen/ecs/v1"
	"google.golang.org/protobuf/encoding/protojson"
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
	eventBus                               *mint.Emitter

	nextSystemID                                   uint32
	systems, leftToRun, notRunWithDependenciesDone map[uint32]*systemRunner
	tickWaitGroup                                  *sync.WaitGroup
	tickCount                                      int

	nextEntityID                   uint32
	liveEntitieIDs, freeEntitieIDs *roaring.Bitmap
	namesStore                     *SparseSet[Name]
	childOfStore                   *SparseSet[ChildOf]
	isAStore                       *SparseSet[IsA]
	positionsStore                 *SparseSet[Position]
	velocitiesStore                *SparseSet[Velocity]
	rotationsStore                 *SparseSet[Rotation]
	directionsStore                *SparseSet[Direction]
	eatsStore                      *SparseSet[Eats]
	likesStore                     *SparseSet[Likes]
	enemyStore                     *SparseSet[Enemy]
	growsStore                     *SparseSet[Grows]
	gravitiesStore                 *SparseSet[Gravity]
	spaceshipStore                 *SparseSet[Spaceship]
	spacestationStore              *SparseSet[Spacestation]
	factionsStore                  *SparseSet[Faction]
	dockedTosStore                 *SparseSet[DockedTo]
	planetStore                    *SparseSet[Planet]
	ruledBysStore                  *SparseSet[RuledBy]
	alliedWithsStore               *SparseSet[AlliedWith]

	PositionVelocitySet *PositionVelocitySet

	patch *ecspb.WorldPatch
}

func NewWorld() *World {
	w := &World{
		eventBus:                   &mint.Emitter{},
		nextSystemID:               1,
		systems:                    map[uint32]*systemRunner{},
		leftToRun:                  map[uint32]*systemRunner{},
		notRunWithDependenciesDone: map[uint32]*systemRunner{},
		tickWaitGroup:              &sync.WaitGroup{},
		tickCount:                  0,

		nextEntityID:      1,
		liveEntitieIDs:    roaring.NewBitmap(),
		freeEntitieIDs:    roaring.NewBitmap(),
		namesStore:        NewSparseSet[Name](nil),
		childOfStore:      NewSparseSet[ChildOf](nil),
		isAStore:          NewSparseSet[IsA](nil),
		positionsStore:    NewSparseSet[Position](nil),
		velocitiesStore:   NewSparseSet[Velocity](nil),
		rotationsStore:    NewSparseSet[Rotation](nil),
		directionsStore:   NewSparseSet[Direction](nil),
		eatsStore:         NewSparseSet[Eats](nil),
		likesStore:        NewSparseSet[Likes](nil),
		enemyStore:        NewSparseSet[Enemy](nil),
		growsStore:        NewSparseSet[Grows](nil),
		gravitiesStore:    NewSparseSet[Gravity](nil),
		spaceshipStore:    NewSparseSet[Spaceship](nil),
		spacestationStore: NewSparseSet[Spacestation](nil),
		factionsStore:     NewSparseSet[Faction](nil),
		dockedTosStore:    NewSparseSet[DockedTo](nil),
		planetStore:       NewSparseSet[Planet](nil),
		ruledBysStore:     NewSparseSet[RuledBy](nil),
		alliedWithsStore:  NewSparseSet[AlliedWith](nil),

		patch: NewWorldPatch(),
	}

	// setup built-in entities
	w.zeroEntity = w.Entity()
	w.resourceEntity = w.Entity()
	w.deadEntity = w.EntityFromU32(DeadEntityID)

	// component sets
	w.PositionVelocitySet = NewPositionVelocitySet(w)

	return w
}

//# region Systems

// AddSystems adds systems to the world. Systems are run in the order they are added.
func (w *World) AddSystems(ss ...System) (err error) {
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

// RemoveSystems removes systems from the world. Systems are removed in the order they are passed.
func (w *World) RemoveSystems(ss ...System) error {
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

// Tick runs all systems in the world. Systems are run in the order they were added.
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

// DisableSystem disables systems from running. Systems are disabled in the order they are passed.
func (w *World) DisableSystem(ss ...System) error {
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

// EnableSystem enables systems to run. Systems are enabled in the order they are passed.
func (w *World) EnableSystem(ss ...System) error {
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

// TickCount returns the number of times the world has ticked.
func (w *World) TickCount() int {
	return w.tickCount
}

//# endregion

// Entity returns a new (or try to reuse dead) entity.
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
	w.namesStore.Remove(es...)
	w.childOfStore.Remove(es...)
	w.isAStore.Remove(es...)
	w.positionsStore.Remove(es...)
	w.velocitiesStore.Remove(es...)
	w.rotationsStore.Remove(es...)
	w.directionsStore.Remove(es...)
	w.eatsStore.Remove(es...)
	w.likesStore.Remove(es...)
	w.enemyStore.Remove(es...)
	w.growsStore.Remove(es...)
	w.gravitiesStore.Remove(es...)
	w.spaceshipStore.Remove(es...)
	w.spacestationStore.Remove(es...)
	w.factionsStore.Remove(es...)
	w.dockedTosStore.Remove(es...)
	w.planetStore.Remove(es...)
	w.ruledBysStore.Remove(es...)
	w.alliedWithsStore.Remove(es...)

	for _, e := range es {
		if !e.IsAlive() {
			continue
		}

		w.liveEntitieIDs.Remove(e.val)

		bumped := e.UpdateVersion()
		w.freeEntitieIDs.Add(bumped.val)

		w.patch.Entities[e.val] = nil

		fireEvent(w, EntityDestroyedEvent{e})
	}
}

func (w *World) Reset() {
	w.namesStore.Clear()
	w.childOfStore.Clear()
	w.isAStore.Clear()
	w.positionsStore.Clear()
	w.velocitiesStore.Clear()
	w.rotationsStore.Clear()
	w.directionsStore.Clear()
	w.eatsStore.Clear()
	w.likesStore.Clear()
	w.enemyStore.Clear()
	w.growsStore.Clear()
	w.gravitiesStore.Clear()
	w.spaceshipStore.Clear()
	w.spacestationStore.Clear()
	w.factionsStore.Clear()
	w.dockedTosStore.Clear()
	w.planetStore.Clear()
	w.ruledBysStore.Clear()
	w.alliedWithsStore.Clear()

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
	return &ecspb.WorldPatch{
		Entities:             map[uint32]*emptypb.Empty{},
		NameComponents:       map[uint32]*ecspb.NameComponent{},
		ChildOfComponents:    map[uint32]*ecspb.ChildOfComponent{},
		IsAComponents:        map[uint32]*ecspb.IsAComponent{},
		PositionComponents:   map[uint32]*ecspb.PositionComponent{},
		VelocityComponents:   map[uint32]*ecspb.VelocityComponent{},
		RotationComponents:   map[uint32]*ecspb.RotationComponent{},
		DirectionComponents:  map[uint32]*ecspb.DirectionComponent{},
		EatsComponents:       map[uint32]*ecspb.EatsComponent{},
		LikesComponents:      map[uint32]*ecspb.LikesComponent{},
		EnemyTags:            map[uint32]*emptypb.Empty{},
		GrowsComponents:      map[uint32]*ecspb.GrowsComponent{},
		GravityComponents:    map[uint32]*ecspb.GravityComponent{},
		SpaceshipTags:        map[uint32]*emptypb.Empty{},
		SpacestationTags:     map[uint32]*emptypb.Empty{},
		FactionComponents:    map[uint32]*ecspb.FactionComponent{},
		DockedToComponents:   map[uint32]*ecspb.DockedToComponent{},
		PlanetTags:           map[uint32]*emptypb.Empty{},
		RuledByComponents:    map[uint32]*ecspb.RuledByComponent{},
		AlliedWithComponents: map[uint32]*ecspb.AlliedWithComponent{},
	}
}

func ResetWorldPatch(patch *ecspb.WorldPatch) *ecspb.WorldPatch {
	clear(patch.Entities)
	clear(patch.NameComponents)
	clear(patch.ChildOfComponents)
	clear(patch.IsAComponents)
	clear(patch.PositionComponents)
	clear(patch.VelocityComponents)
	clear(patch.RotationComponents)
	clear(patch.DirectionComponents)
	clear(patch.EatsComponents)
	clear(patch.LikesComponents)
	clear(patch.EnemyTags)
	clear(patch.GrowsComponents)
	clear(patch.GravityComponents)
	clear(patch.SpaceshipTags)
	clear(patch.SpacestationTags)
	clear(patch.FactionComponents)
	clear(patch.DockedToComponents)
	clear(patch.PlanetTags)
	clear(patch.RuledByComponents)
	clear(patch.AlliedWithComponents)
	return patch
}

func MergeWorldWriteAheadLogs(patchs ...*ecspb.WorldPatch) *ecspb.WorldPatch {
	merged := NewWorldPatch()
	for _, patch := range patchs {
		for k, v := range patch.Entities {
			merged.Entities[k] = v
		}

		// merge Names
		for k, v := range patch.NameComponents {
			merged.NameComponents[k] = v
		}

		// merge ChildOf
		for k, v := range patch.ChildOfComponents {
			merged.ChildOfComponents[k] = v
		}

		// merge IsA
		for k, v := range patch.IsAComponents {
			merged.IsAComponents[k] = v
		}

		// merge Positions
		for k, v := range patch.PositionComponents {
			merged.PositionComponents[k] = v
		}

		// merge Velocities
		for k, v := range patch.VelocityComponents {
			merged.VelocityComponents[k] = v
		}

		// merge Rotations
		for k, v := range patch.RotationComponents {
			merged.RotationComponents[k] = v
		}

		// merge Directions
		for k, v := range patch.DirectionComponents {
			merged.DirectionComponents[k] = v
		}

		// merge Eats
		for k, v := range patch.EatsComponents {
			merged.EatsComponents[k] = v
		}

		// merge Likes
		for k, v := range patch.LikesComponents {
			merged.LikesComponents[k] = v
		}

		// merge Enemy
		for k, v := range patch.EnemyTags {
			merged.EnemyTags[k] = v
		}

		// merge Grows
		for k, v := range patch.GrowsComponents {
			merged.GrowsComponents[k] = v
		}

		// merge Gravities
		for k, v := range patch.GravityComponents {
			merged.GravityComponents[k] = v
		}

		// merge Spaceship
		for k, v := range patch.SpaceshipTags {
			merged.SpaceshipTags[k] = v
		}

		// merge Spacestation
		for k, v := range patch.SpacestationTags {
			merged.SpacestationTags[k] = v
		}

		// merge Factions
		for k, v := range patch.FactionComponents {
			merged.FactionComponents[k] = v
		}

		// merge DockedTos
		for k, v := range patch.DockedToComponents {
			merged.DockedToComponents[k] = v
		}

		// merge Planet
		for k, v := range patch.PlanetTags {
			merged.PlanetTags[k] = v
		}

		// merge RuledBys
		for k, v := range patch.RuledByComponents {
			merged.RuledByComponents[k] = v
		}

		// merge AlliedWiths
		for k, v := range patch.AlliedWithComponents {
			merged.AlliedWithComponents[k] = v
		}

	}

	return merged
}

func (w *World) ApplyPatches(patches ...*ecspb.WorldPatch) {
	for _, patch := range patches {
		for k, v := range patch.Entities {
			if v == nil {
				w.DestroyEntities(w.EntityFromU32(k))
			} else {
				w.EntityFromU32(k)
			}
		}

		// apply Names
		for val, c := range patch.NameComponents {
			e := w.EntityFromU32(val)
			w.ApplyNamePatch(e, c)
		}
		// apply ChildOf
		for val, c := range patch.ChildOfComponents {
			e := w.EntityFromU32(val)
			w.ApplyChildOfPatch(e, c)
		}
		// apply IsA
		for val, c := range patch.IsAComponents {
			e := w.EntityFromU32(val)
			w.ApplyIsAPatch(e, c)
		}
		// apply Positions
		for val, c := range patch.PositionComponents {
			e := w.EntityFromU32(val)
			w.ApplyPositionPatch(e, c)
		}
		// apply Velocities
		for val, c := range patch.VelocityComponents {
			e := w.EntityFromU32(val)
			w.ApplyVelocityPatch(e, c)
		}
		// apply Rotations
		for val, c := range patch.RotationComponents {
			e := w.EntityFromU32(val)
			w.ApplyRotationPatch(e, c)
		}
		// apply Directions
		for val, c := range patch.DirectionComponents {
			e := w.EntityFromU32(val)
			w.ApplyDirectionPatch(e, c)
		}
		// apply Eats
		for val, c := range patch.EatsComponents {
			e := w.EntityFromU32(val)
			w.ApplyEatsPatch(e, c)
		}
		// apply Likes
		for val, c := range patch.LikesComponents {
			e := w.EntityFromU32(val)
			w.ApplyLikesPatch(e, c)
		}
		// apply Enemy
		for val, c := range patch.EnemyTags {
			e := w.EntityFromU32(val)
			w.ApplyEnemyPatch(e, c)
		}
		// apply Grows
		for val, c := range patch.GrowsComponents {
			e := w.EntityFromU32(val)
			w.ApplyGrowsPatch(e, c)
		}
		// apply Gravities
		for val, c := range patch.GravityComponents {
			e := w.EntityFromU32(val)
			w.ApplyGravityPatch(e, c)
		}
		// apply Spaceship
		for val, c := range patch.SpaceshipTags {
			e := w.EntityFromU32(val)
			w.ApplySpaceshipPatch(e, c)
		}
		// apply Spacestation
		for val, c := range patch.SpacestationTags {
			e := w.EntityFromU32(val)
			w.ApplySpacestationPatch(e, c)
		}
		// apply Factions
		for val, c := range patch.FactionComponents {
			e := w.EntityFromU32(val)
			w.ApplyFactionPatch(e, c)
		}
		// apply DockedTos
		for val, c := range patch.DockedToComponents {
			e := w.EntityFromU32(val)
			w.ApplyDockedToPatch(e, c)
		}
		// apply Planet
		for val, c := range patch.PlanetTags {
			e := w.EntityFromU32(val)
			w.ApplyPlanetPatch(e, c)
		}
		// apply RuledBys
		for val, c := range patch.RuledByComponents {
			e := w.EntityFromU32(val)
			w.ApplyRuledByPatch(e, c)
		}
		// apply AlliedWiths
		for val, c := range patch.AlliedWithComponents {
			e := w.EntityFromU32(val)
			w.ApplyAlliedWithPatch(e, c)
		}

	}
}

func (w *World) MarshalPatch() ([]byte, error) {
	return w.patch.MarshalVT()
}

func (w *World) MarshalPatchJSON() ([]byte, error) {
	return w.patch.MarshalJSON()
}

func (w *World) MarshalPatchPrettyJSON() ([]byte, error) {
	return protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: false,
		UseProtoNames:   false,
		Indent:          "  ",
	}.Marshal(w.patch)
}

func (w *World) UnmarshalPatch(data []byte) error {
	return w.patch.UnmarshalVT(data)
}

func (w *World) UnmarshalPatchJSON(data []byte) error {
	return w.patch.UnmarshalJSON(data)
}
