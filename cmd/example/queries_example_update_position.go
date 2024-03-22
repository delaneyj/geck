package example

import "github.com/delaneyj/geck/cmd/example/ecs"

type UpdatePositionsSystem struct {
}

func (sys *UpdatePositionsSystem) Name() string {
	return "UpdatePositions"
}

func (sys *UpdatePositionsSystem) ReliesOn() []string {
	return nil
}

func (sys *UpdatePositionsSystem) Tick(w *ecs.World) error {
	// variables := map[string]WorldEntity{}
	iter := w.VelocityWriteIter()
	for iter.HasNext() {
		e, v := iter.NextVelocity()
		p, pOk := e.WritablePosition()
		if !pOk {
			continue
		}
		p.X += v.X
		p.Y += v.Y
	}
	return nil
}
