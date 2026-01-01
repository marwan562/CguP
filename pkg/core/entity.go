package core

import (
	"math"
)

type Entity struct {
	Position [3]float64
	Velocity [3]float64
	Data     [100]float64
}

// ComplexUpdate performs a CPU-intensive update on the entity.
func (e *Entity) ComplexUpdate(dt float64) {
	for i := 0; i < 3; i++ {
		e.Position[i] += e.Velocity[i] * dt
	}

	for i := 0; i < 100; i++ {
		e.Data[i] = math.Sin(e.Data[i]) * math.Cos(e.Position[0])
		e.Data[i] = math.Sqrt(math.Abs(e.Data[i])) + 0.001
	}
}

// SimpleUpdate performs a simple physics update on the entity.
func (e *Entity) SimpleUpdate(dt float64) {
	for i := 0; i < 3; i++ {
		e.Position[i] += e.Velocity[i] * dt
	}
}
