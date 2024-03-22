package generator

import (
	"testing"

	"github.com/delaneyj/geck/generator/out"
)

func TestOut(t *testing.T) {
	w := out.NewWorld()

	foo := w.CreateEntity("foo")
	w.SetPositions(out.Position{
		X: 1,
		Y: 2,
		Z: 3,
	}, foo)

	w.MarshalAll()
}
