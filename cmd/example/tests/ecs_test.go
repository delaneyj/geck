package example

import (
	"context"
	"slices"
	"testing"

	"github.com/delaneyj/geck/cmd/example/ecs"
	"github.com/stretchr/testify/assert"
)

type RelationshipSystem struct {
	foundCount int
}

func (sys *RelationshipSystem) Name() string {
	return "Relationship"
}

func (sys *RelationshipSystem) ReliesOn() ecs.ReliedOnIter {
	return func(yield ecs.System) bool {
		return false
	}
}

func (sys *RelationshipSystem) Initialize(ctx context.Context, w *ecs.World) error {
	return nil
}

// from https://ajmmertens.medium.com/building-games-in-ecs-with-entity-relationships-657275ba2c6c
func (sys *RelationshipSystem) Tick(ctx context.Context, w *ecs.World) error {
	sys.foundCount = 0
	for spaceship := range w.AllSpaceshipEntities {

		spacesphipFaction, ok := w.Faction(spaceship)
		if !ok {
			continue
		}

		dockedTo, ok := w.DockedTo(spaceship)
		if !ok {
			continue
		}
		if !w.HasPlanetTag(dockedTo.Entity) {
			continue
		}
		planet := dockedTo.Entity

		planetFaction, ok := w.RuledBy(planet)
		if !ok {
			continue
		}

		hasAllies := false
		for ally := range w.AlliedWith(spacesphipFaction.Entity) {
			if ally == planetFaction.Entity {
				sys.foundCount++
				hasAllies = true
				break
			}
		}

		if !hasAllies {
			continue
		}

		// log.Printf(
		// 	"Spaceship %v Faction(%v) is docked to planet %v Faction(%v)",
		// 	spaceship,
		// 	spacesphipFaction,
		// 	planet,
		// 	planetFaction,
		// )
		sys.foundCount++
	}

	return nil
}

func TestECSRelationships(t *testing.T) {
	// Ported from https://ajmmertens.medium.com/building-games-in-ecs-with-entity-relationships-657275ba2c6c
	w := ecs.NewWorld()
	defer w.Reset()

	relSys := &RelationshipSystem{}
	if err := w.AddSystems(t.Context(), relSys); err != nil {
		t.Fatal(err)
	}

	marsColonists := w.NextEntity()
	earthFederation := w.NextEntity()
	miningCorp := w.NextEntity()

	w.LinkAlliedWith(earthFederation, marsColonists)
	w.LinkAlliedWith(marsColonists, earthFederation)

	mars := w.NextEntity(
		ecs.WithPlanetTag(),
		ecs.WithRuledBy(marsColonists),
	)
	earth := w.NextEntity(
		ecs.WithPlanetTag(),
		ecs.WithRuledBy(earthFederation),
	)

	iss := w.NextEntity(
		ecs.WithSpacestationTag(),
		ecs.WithRuledBy(earthFederation),
	)

	w.NextEntity(
		ecs.WithName("e1"),
		ecs.WithSpaceshipTag(),
		ecs.WithFaction(marsColonists),
		ecs.WithDockedTo(earth),
	)
	w.NextEntity(
		ecs.WithName("e2"),
		ecs.WithSpaceshipTag(),
		ecs.WithFaction(miningCorp),
		ecs.WithDockedTo(earth),
	)
	w.NextEntity(
		ecs.WithName("e3"),
		ecs.WithSpaceshipTag(),
		ecs.WithFaction(earthFederation),
		ecs.WithDockedTo(mars),
	)
	w.NextEntity(
		ecs.WithName("e4"),
		ecs.WithSpaceshipTag(),
		ecs.WithFaction(earthFederation),
		ecs.WithDockedTo(iss),
	)

	ctx := context.Background()
	if err := w.Tick(ctx); err != nil {
		t.Fatal(err)
	}

	// jsonBytes, err := w.MarshalPatchPrettyJSON()
	// assert.NoError(t, err)
	// pbBytes, err := w.MarshalPatch()
	// assert.NoError(t, err)
	// log.Printf("JSON: %d, PB: %d", len(jsonBytes), len(pbBytes))

	assert.Equal(t, relSys.foundCount, 2)
}

func TestECS(t *testing.T) {
	// The world is the container for all ECS data.
	// It stores the entities and their components, does queries and runs systems.
	// Typically there is only a single world, but there is no limit on the number of worlds an application can create.
	w := ecs.NewWorld()
	defer w.Reset()

	// // An entity is a unique thing in the world, and is represented by a 64 bit id.
	// // Entities can be created and deleted.
	// // If an entity is deleted it is no longer considered "alive".
	// // A world can contain up to 2 billion(!) alive entities.
	// // Entity identifiers contain a few bits that make it possible to check whether an entity is alive or not.
	e := w.NextEntity()
	assert.True(t, w.IsAlive(e))
	assert.Equal(t, e.Generation(), uint32(0))
	w.DestroyEntities(e)
	assert.False(t, w.IsAlive(e))

	e = w.NextEntity()
	assert.True(t, w.IsAlive(e))
	assert.Equal(t, e.Generation(), uint32(1))

	// A component is a type of which instances can be added and removed to entities.
	// Each component can be added only once to an entity (though not really, see Relation).
	// In C applications components must be registered before use.
	// By default in C++ this happens automatically.
	w.SetPositionFromValues(e, 1, 2, 3)
	w.SetVelocityFromValues(e, 4, 5, 6)
	w.TagWithEnemy(e)

	pos, hasPos := w.Position(e)
	assert.True(t, hasPos)
	assert.Equal(t, pos.X, float32(1))
	assert.Equal(t, pos.Y, float32(2))
	assert.Equal(t, pos.Z, float32(3))

	assert.True(t, w.HasPosition(e))
	w.RemovePosition(e)
	assert.False(t, w.HasPosition(e))

	// A tag is a component that does not have any data.
	// In Flecs tags can be either empty types (in C++) or
	// regular entities (C & C++) that do not have the EcsComponent component (or have an EcsComponent component with size 0).
	// Tags can be added & removed using the same APIs as adding & removing components, but because tags have no data, they cannot be assigned a value.
	// Because tags (like components) are regular entities, they can be created & deleted at runtime.
	assert.True(t, w.HasEnemyTag(e))
	w.RemoveEnemyTag(e)
	assert.False(t, w.HasEnemyTag(e))

	// A pair is a combination of two entity ids.
	// Relations can be used to store entity relationships,
	// where the first id represents the relationship kind and the second id represents the relationship target (called "object").
	// This is best explained by an example:

	bob := w.NextEntity()
	alice := w.NextEntity()
	w.LinkLikes(alice, bob) // Bob likes Alice
	assert.True(t, w.LikesIsLinked(bob, alice))

	w.RemoveLikesRelationships(bob, alice)
	assert.False(t, w.LikesIsLinked(bob, alice))

	apples := w.NextEntity()
	pears := w.NextEntity()

	w.LinkEats(bob, apples, 3)
	w.LinkEats(bob, pears, 2)
	w.LinkGrows(bob, pears)

	bobEats := slices.Collect(w.Eats(bob))
	assert.True(t, apples.InSlice(bobEats...))
	assert.True(t, pears.InSlice(bobEats...))

	w.SetGravityResource(-9.8)
	assert.True(t, w.HasGravityResource())
	g, ok := w.GravityResource()
	assert.True(t, ok)
	assert.Equal(t, g, ecs.GravityComponent{G: -9.8})

	w.NextEntities(1000)
	// log.Print(ee)
}
