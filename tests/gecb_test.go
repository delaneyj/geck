package tests

import (
	"log"
	"testing"

	"github.com/delaneyj/geck"
	"github.com/stretchr/testify/assert"
)

func TestWorld(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	w := geck.NewWorld()
	assert.NotNil(t, w)

	type Vector2 struct {
		X, Y int
	}
	badID := geck.ID(1234)
	badIDSet := geck.NewIDSet(badID)
	resetPosition := Vector2{X: 1, Y: 2}

	assert.Panics(t, func() { geck.ComponentData(w, badID, badID, &resetPosition) })
	assert.Panics(t, func() { geck.SetComponentData(w, badID, resetPosition, badIDSet) })
	assert.Panics(t, func() { geck.RemoveComponentFrom(w, badIDSet, badIDSet) })

	badID = w.CreateEntity()

	assert.Panics(t, func() { geck.ComponentData(w, badID, badID, &resetPosition) })
	assert.Panics(t, func() { geck.SetComponentData(w, badID, resetPosition, badIDSet) })

	position := geck.RegisterComponent(w, resetPosition, "position")
	// velocity := RegisterComponent(w, Vector2{}, 9000)

	p := &Vector2{}
	assert.Panics(t, func() { geck.ComponentData(w, badID, badID, p) })

	e1 := w.CreateEntity()
	e1Set := geck.NewIDSet(e1)
	assert.Panics(t, func() { geck.ComponentData(w, badID, badID, p) })

	positionSet := geck.NewIDSet(position)
	geck.AddComponentsTo(w, positionSet, e1Set)
	assert.True(t, w.HasComponents(e1Set, positionSet))

	geck.ComponentData(w, position, e1, p)

	assert.Equal(t, *p, resetPosition)

	assert.Panics(t, func() { geck.ComponentData(w, badID, 123, p) })

	p2 := Vector2{X: 3, Y: 4}
	geck.SetComponentData(w, position, p2, e1Set)

	assert.Panics(t, func() { geck.SetComponentData(w, position, p2, badIDSet) })

	geck.RemoveComponentFrom(w, positionSet, e1Set)
	assert.False(t, w.HasComponents(e1Set, positionSet))

	geck.AddComponentsTo(w, positionSet, e1Set)

	velocity := geck.RegisterComponent(w, Vector2{}, "velocity")
	velocitySet := geck.NewIDSet(velocity)
	geck.AddComponentsTo(w, velocitySet, e1Set)

	tag := w.CreateEntity("tag")
	tagSet := geck.NewIDSet(tag)
	geck.AddComponentsTo(w, tagSet, e1Set)

}

func BenchmarkArchetype(b *testing.B) {
	b.StopTimer()
	const nPos = 9000
	const nPosVel = 1000

	w := geck.NewWorld()

	type Vector2 struct {
		X, Y int
	}
	position := geck.RegisterComponent(w, Vector2{}, "position")
	velocity := geck.RegisterComponent(w, Vector2{1, 1}, "velocity")

	w.CreateEntitiesWith(nPos, geck.NewIDSet(position))
	w.CreateEntitiesWith(nPosVel, geck.NewIDSet(position, velocity))

	p, v := &Vector2{}, &Vector2{}
	iter := w.QueryAnd(position, velocity)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		iter.Reset()
		for iter.HasNext() {
			iter.Next()
			geck.Data(iter, p, 0)
			geck.Data(iter, v, 1)
			p.X += v.X
			p.Y += v.Y
		}
	}
}
