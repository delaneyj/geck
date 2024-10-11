package ecs

import (
	"fmt"

	"github.com/btvoidx/mint"
)

type empty struct{}

type World struct {
	nextEntityID                 uint32
	livingEntities, freeEntities *SparseSet[empty]
	resourceEntity               Entity
	systems                      []SystemTicker
	eventBus                     *mint.Emitter

	// Tags
	enemyTags        *SparseSet[empty]
	spaceshipTags    *SparseSet[empty]
	spacestationTags *SparseSet[empty]
	planetTags       *SparseSet[empty]

	// Components
	nameComponents       *SparseSet[NameComponent]
	childOfComponents    *SparseSet[ChildOfComponent]
	isAComponents        *SparseSet[IsAComponent]
	positionComponents   *SparseSet[PositionComponent]
	velocityComponents   *SparseSet[VelocityComponent]
	rotationComponents   *SparseSet[RotationComponent]
	directionComponents  *SparseSet[DirectionComponent]
	eatsComponents       *SparseSet[EatsComponent]
	likesComponents      *SparseSet[LikesComponent]
	growsComponents      *SparseSet[GrowsComponent]
	gravityComponents    *SparseSet[GravityComponent]
	factionComponents    *SparseSet[FactionComponent]
	dockedToComponents   *SparseSet[DockedToComponent]
	ruledByComponents    *SparseSet[RuledByComponent]
	alliedWithComponents *SparseSet[AlliedWithComponent]
}

func NewWorld() *World {
	w := &World{
		nextEntityID:   0,
		livingEntities: NewSparseSet[empty](),
		freeEntities:   NewSparseSet[empty](),
		eventBus:       &mint.Emitter{},

		// Initialize tags
		enemyTags:        NewSparseSet[empty](),
		spaceshipTags:    NewSparseSet[empty](),
		spacestationTags: NewSparseSet[empty](),
		planetTags:       NewSparseSet[empty](),

		// Initialize components
		nameComponents:       NewSparseSet[NameComponent](),
		childOfComponents:    NewSparseSet[ChildOfComponent](),
		isAComponents:        NewSparseSet[IsAComponent](),
		positionComponents:   NewSparseSet[PositionComponent](),
		velocityComponents:   NewSparseSet[VelocityComponent](),
		rotationComponents:   NewSparseSet[RotationComponent](),
		directionComponents:  NewSparseSet[DirectionComponent](),
		eatsComponents:       NewSparseSet[EatsComponent](),
		likesComponents:      NewSparseSet[LikesComponent](),
		growsComponents:      NewSparseSet[GrowsComponent](),
		gravityComponents:    NewSparseSet[GravityComponent](),
		factionComponents:    NewSparseSet[FactionComponent](),
		dockedToComponents:   NewSparseSet[DockedToComponent](),
		ruledByComponents:    NewSparseSet[RuledByComponent](),
		alliedWithComponents: NewSparseSet[AlliedWithComponent](),
	}
	w.resourceEntity = w.CreateEntity()

	return w
}

func (w *World) AddSystems(systems ...System) error {
	for _, s := range systems {
		if err := s.Initialize(w); err != nil {
			return fmt.Errorf("failed to initialize system: %w", err)
		}

		sysTicker, ok := s.(SystemTicker)
		if !ok {
			continue
		}

		w.systems = append(w.systems, sysTicker)
	}

	return nil
}

func (w *World) Tick() error {
	for _, s := range w.systems {
		if err := s.Tick(w); err != nil {
			return err
		}
	}
	return nil
}

type System interface {
	Initialize(w *World) error
	ReliesOn(func(reliedOn System) bool)
}

type SystemTicker interface {
	System
	Tick(w *World) error
}
