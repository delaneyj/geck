package ecs

import (
	"fmt"

	"github.com/RoaringBitmap/roaring"
	"github.com/btvoidx/mint"
)

type empty struct{}

type World struct {
	nextEntityID                 uint32
	livingEntities, freeEntities *roaring.Bitmap
	resourceEntity               Entity
	systems                      []SystemTicker
	eventBus                     *mint.Emitter

	// Tags
	enemyTags        *SparseSet[empty]
	spaceshipTags    *SparseSet[empty]
	spacestationTags *SparseSet[empty]
	planetTags       *SparseSet[empty]

	// Components
	nameComponents       *SparseSet[Name]
	childOfComponents    *SparseSet[ChildOf]
	isAComponents        *SparseSet[IsA]
	positionComponents   *SparseSet[Position]
	velocityComponents   *SparseSet[Velocity]
	rotationComponents   *SparseSet[Rotation]
	directionComponents  *SparseSet[Direction]
	eatsComponents       *SparseSet[Eats]
	likesComponents      *SparseSet[Likes]
	growsComponents      *SparseSet[Grows]
	gravityComponents    *SparseSet[Gravity]
	factionComponents    *SparseSet[Faction]
	dockedToComponents   *SparseSet[DockedTo]
	ruledByComponents    *SparseSet[RuledBy]
	alliedWithComponents *SparseSet[AlliedWith]
}

func NewWorld() *World {
	w := &World{
		nextEntityID:   0,
		livingEntities: roaring.New(),
		freeEntities:   roaring.New(),
		eventBus:       &mint.Emitter{},

		// Initialize tags
		enemyTags:        NewSparseSet[empty](),
		spaceshipTags:    NewSparseSet[empty](),
		spacestationTags: NewSparseSet[empty](),
		planetTags:       NewSparseSet[empty](),

		// Initialize components
		nameComponents:       NewSparseSet[Name](),
		childOfComponents:    NewSparseSet[ChildOf](),
		isAComponents:        NewSparseSet[IsA](),
		positionComponents:   NewSparseSet[Position](),
		velocityComponents:   NewSparseSet[Velocity](),
		rotationComponents:   NewSparseSet[Rotation](),
		directionComponents:  NewSparseSet[Direction](),
		eatsComponents:       NewSparseSet[Eats](),
		likesComponents:      NewSparseSet[Likes](),
		growsComponents:      NewSparseSet[Grows](),
		gravityComponents:    NewSparseSet[Gravity](),
		factionComponents:    NewSparseSet[Faction](),
		dockedToComponents:   NewSparseSet[DockedTo](),
		ruledByComponents:    NewSparseSet[RuledBy](),
		alliedWithComponents: NewSparseSet[AlliedWith](),
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
