package ecs

type QueryExamplePositionVelocitiesArgs struct {
	Velocity VelocityComponent

	Position *PositionComponent
}

type queryExamplePositionVelocitiesIter func(e Entity, args QueryExamplePositionVelocitiesArgs) bool

func (w *World) QueryExamplePositionVelocity(yield queryExamplePositionVelocitiesIter) {
	args := QueryExamplePositionVelocitiesArgs{}

	var ok bool

	for e, first := range w.AllVelocities {
		args.Velocity = first
		ok = true

		args.Position, ok = w.MutablePosition(e)

		if !ok {
			continue
		}

		if !yield(e, args) {
			break
		}
	}
}

func (w *World) QueryExamplePositionVelocityEntities(yield func(e Entity) bool) {
	for e := range w.AllVelocitiesEntities {
		if !w.HasPosition(e) {
			continue
		}

		if !yield(e) {
			break
		}
	}
}
