package ecs

import (
	"context"
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
	nameComponents      *SparseSet[NameComponent]
	childOfComponents   *SparseSet[ChildOfComponent]
	isAComponents       *SparseSet[IsAComponent]
	positionComponents  *SparseSet[PositionComponent]
	velocityComponents  *SparseSet[VelocityComponent]
	rotationComponents  *SparseSet[RotationComponent]
	directionComponents *SparseSet[DirectionComponent]
	eatsComponents      *SparseSet[EatsComponent]
	gravityComponents   *SparseSet[GravityComponent]
	factionComponents   *SparseSet[FactionComponent]
	dockedToComponents  *SparseSet[DockedToComponent]
	ruledByComponents   *SparseSet[RuledByComponent]

	// Relationships
	likesRelationships      *LikesRelationship
	growsRelationships      *GrowsRelationship
	alliedWithRelationships *AlliedWithRelationship
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
		nameComponents:      NewSparseSet[NameComponent](),
		childOfComponents:   NewSparseSet[ChildOfComponent](),
		isAComponents:       NewSparseSet[IsAComponent](),
		positionComponents:  NewSparseSet[PositionComponent](),
		velocityComponents:  NewSparseSet[VelocityComponent](),
		rotationComponents:  NewSparseSet[RotationComponent](),
		directionComponents: NewSparseSet[DirectionComponent](),
		eatsComponents:      NewSparseSet[EatsComponent](),
		gravityComponents:   NewSparseSet[GravityComponent](),
		factionComponents:   NewSparseSet[FactionComponent](),
		dockedToComponents:  NewSparseSet[DockedToComponent](),
		ruledByComponents:   NewSparseSet[RuledByComponent](),

		// Initialize relationships
		likesRelationships:      NewLikesRelationship(),
		growsRelationships:      NewGrowsRelationship(),
		alliedWithRelationships: NewAlliedWithRelationship(),
	}

	w.Reset()

	return w
}

func (w *World) Reset() {
	w.nextEntityID = 0
	w.livingEntities.Clear()
	w.freeEntities.Clear()
	w.resourceEntity = w.NextEntity()

	// Reset tags
	w.enemyTags.Clear()
	w.spaceshipTags.Clear()
	w.spacestationTags.Clear()
	w.planetTags.Clear()

	// Reset components
	w.nameComponents.Clear()
	w.childOfComponents.Clear()
	w.isAComponents.Clear()
	w.positionComponents.Clear()
	w.velocityComponents.Clear()
	w.rotationComponents.Clear()
	w.directionComponents.Clear()
	w.eatsComponents.Clear()
	w.gravityComponents.Clear()
	w.factionComponents.Clear()
	w.dockedToComponents.Clear()
	w.ruledByComponents.Clear()

	// Reset relationships
	w.likesRelationships.Clear()
	w.growsRelationships.Clear()
	w.alliedWithRelationships.Clear()
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

func (w *World) Tick(ctx context.Context) error {
	for _, s := range w.systems {
		if err := s.Tick(ctx, w); err != nil {
			return err
		}
	}
	return nil
}

type ReliedOnIter func(reliedOn System) bool

type System interface {
	Initialize(w *World) error
	ReliesOn() ReliedOnIter
}

type SystemTicker interface {
	System
	Tick(ctx context.Context, w *World) error
}
