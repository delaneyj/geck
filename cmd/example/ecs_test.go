package example

// import (
// 	"context"
// 	"testing"

// 	"github.com/delaneyj/geck/cmd/example/ecs"
// 	"github.com/samber/lo"
// 	"github.com/stretchr/testify/assert"
// )

// type RelationshipSystem struct {
// 	foundCount int
// }

// func (sys *RelationshipSystem) Name() string {
// 	return "Relationship"
// }

// func (sys *RelationshipSystem) ReliesOn() []string {
// 	return nil
// }

// func (sys *RelationshipSystem) Initialize(w *ecs.World) error {
// 	return nil
// }

// // from https://ajmmertens.medium.com/building-games-in-ecs-with-entity-relationships-657275ba2c6c
// func (sys *RelationshipSystem) Tick(ctx context.Context, w *ecs.World) error {
// 	sys.foundCount = 0
// 	iter := w.SpaceshipReadIter()
// 	for iter.HasNext() {
// 		spaceship := iter.NextEntity()

// 		spacesphipFaction, ok := spaceship.ReadFaction()
// 		if !ok {
// 			continue
// 		}

// 		dockedTo, ok := spaceship.ReadDockedTo()
// 		if !ok {
// 			continue
// 		}
// 		if !dockedTo.HasPlanetTag() {
// 			continue
// 		}
// 		planet := dockedTo

// 		planetFaction, ok := planet.ReadRuledBy()
// 		if !ok {
// 			continue
// 		}

// 		planetAllies, ok := planetFaction.ReadAlliedWith()
// 		if !ok {
// 			continue
// 		}
// 		if !lo.Contains(planetAllies, spacesphipFaction) {
// 			continue
// 		}

// 		// log.Printf(
// 		// 	"Spaceship %v Faction(%v) is docked to planet %v Faction(%v)",
// 		// 	spaceship,
// 		// 	spacesphipFaction,
// 		// 	planet,
// 		// 	planetFaction,
// 		// )
// 		sys.foundCount++
// 	}

// 	return nil
// }

// func TestECSRelationships(t *testing.T) {
// 	// Ported from https://ajmmertens.medium.com/building-games-in-ecs-with-entity-relationships-657275ba2c6c
// 	w := ecs.NewWorld()
// 	defer w.Reset()

// 	relSys := &RelationshipSystem{}
// 	if err := w.AddSystems(relSys); err != nil {
// 		t.Fatal(err)
// 	}

// 	marsColonists := w.Entity()
// 	earthFederation := w.Entity()
// 	miningCorp := w.Entity()

// 	marsColonists.SetAlliedWith(earthFederation)
// 	earthFederation.SetAlliedWith(marsColonists)

// 	mars := w.Entity().TagWithPlanet().SetRuledBy(marsColonists)
// 	earth := w.Entity().TagWithPlanet().SetRuledBy(earthFederation)

// 	iss := w.Entity().TagWithSpacestation().SetRuledBy(earthFederation)

// 	w.EntityWithName("e1").TagWithSpaceship().SetFaction(marsColonists).SetDockedTo(earth)
// 	w.EntityWithName("e2").TagWithSpaceship().SetFaction(miningCorp).SetDockedTo(earth)
// 	w.EntityWithName("e3").TagWithSpaceship().SetFaction(earthFederation).SetDockedTo(mars)
// 	w.EntityWithName("e4").TagWithSpaceship().SetFaction(earthFederation).SetDockedTo(iss)

// 	if err := w.Tick(context.Background()); err != nil {
// 		t.Fatal(err)
// 	}

// 	// jsonBytes, err := w.MarshalPatchPrettyJSON()
// 	// assert.NoError(t, err)
// 	// pbBytes, err := w.MarshalPatch()
// 	// assert.NoError(t, err)
// 	// log.Printf("JSON: %d, PB: %d", len(jsonBytes), len(pbBytes))

// 	assert.Equal(t, relSys.foundCount, 2)
// }

// func TestECS(t *testing.T) {
// 	// The world is the container for all ECS data.
// 	// It stores the entities and their components, does queries and runs systems.
// 	// Typically there is only a single world, but there is no limit on the number of worlds an application can create.
// 	w := ecs.NewWorld()
// 	defer w.Reset()

// 	// // An entity is a unique thing in the world, and is represented by a 64 bit id.
// 	// // Entities can be created and deleted.
// 	// // If an entity is deleted it is no longer considered "alive".
// 	// // A world can contain up to 2 billion(!) alive entities.
// 	// // Entity identifiers contain a few bits that make it possible to check whether an entity is alive or not.
// 	e := w.Entity()
// 	assert.True(t, e.IsAlive())
// 	assert.Equal(t, e.Version(), uint32(0))
// 	e.Destroy()
// 	assert.False(t, e.IsAlive())

// 	e = w.Entity()
// 	assert.True(t, e.IsAlive())
// 	assert.Equal(t, e.Version(), uint32(1))

// 	// A component is a type of which instances can be added and removed to entities.
// 	// Each component can be added only once to an entity (though not really, see Relation).
// 	// In C applications components must be registered before use.
// 	// By default in C++ this happens automatically.
// 	e.
// 		SetPosition(ecs.Position{X: 1, Y: 2, Z: 3}).
// 		SetVelocity(ecs.Velocity{X: 4, Y: 5, Z: 6}).
// 		TagWithEnemy()

// 	pos, hasPos := e.ReadPosition()
// 	assert.True(t, hasPos)
// 	assert.Equal(t, pos.X, float32(1))
// 	assert.Equal(t, pos.Y, float32(2))
// 	assert.Equal(t, pos.Z, float32(3))

// 	assert.True(t, e.HasPosition())
// 	e.RemovePosition()
// 	assert.False(t, e.HasPosition())

// 	// A tag is a component that does not have any data.
// 	// In Flecs tags can be either empty types (in C++) or
// 	// regular entities (C & C++) that do not have the EcsComponent component (or have an EcsComponent component with size 0).
// 	// Tags can be added & removed using the same APIs as adding & removing components, but because tags have no data, they cannot be assigned a value.
// 	// Because tags (like components) are regular entities, they can be created & deleted at runtime.
// 	assert.True(t, e.HasEnemyTag())
// 	e.RemoveEnemyTag()
// 	assert.False(t, e.HasEnemyTag())

// 	// A pair is a combination of two entity ids.
// 	// Relations can be used to store entity relationships,
// 	// where the first id represents the relationship kind and the second id represents the relationship target (called "object").
// 	// This is best explained by an example:

// 	bob := w.Entity()
// 	alice := w.Entity()
// 	bob.SetLikes(alice) // Bob likes Alice
// 	alice.SetLikes(bob) // Alice likes Bob
// 	assert.True(t, bob.LikesContains(alice))

// 	bob.RemoveLikes(alice)
// 	assert.False(t, bob.LikesContains(alice))

// 	apples := w.Entity()
// 	pears := w.Entity()

// 	bob.SetEats(ecs.Eats{
// 		Entities: []ecs.Entity{apples, pears},
// 		Amounts:  []uint8{3, 2},
// 	})
// 	bob.SetGrows(pears)

// 	bobEats, _ := bob.ReadEats()
// 	assert.True(t, lo.Contains(bobEats.Entities, apples))
// 	assert.True(t, lo.Contains(bobEats.Entities, pears))

// 	w.SetGravityResource(ecs.Gravity(-9.8))
// 	assert.True(t, w.HasGravityResource())
// 	g, ok := w.GravityResource()
// 	assert.True(t, ok)
// 	assert.Equal(t, g, ecs.Gravity(-9.8))

// 	w.Entities(1000)
// 	// log.Print(ee)
// }
