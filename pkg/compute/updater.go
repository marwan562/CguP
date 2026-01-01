package compute

import (
	"sync"

	"github.com/marwan562/CguP/pkg/core"
)

// Updater defines the interface for updating entities.
type Updater interface {
	Update(entities []*core.Entity, dt float64, complex bool)
}

// SequentialUpdater updates entities one by one.
type SequentialUpdater struct{}

func (u *SequentialUpdater) Update(entities []*core.Entity, dt float64, complex bool) {
	for _, e := range entities {
		if complex {
			e.ComplexUpdate(dt)
		} else {
			e.SimpleUpdate(dt)
		}
	}
}

// ParallelUpdater updates entities concurrently using batches.
type ParallelUpdater struct {
	BatchSize int
}

// NewParallelUpdater creates a new ParallelUpdater with the specified batch size.
func NewParallelUpdater(batchSize int) *ParallelUpdater {
	if batchSize <= 0 {
		batchSize = 100 // Default batch size
	}
	return &ParallelUpdater{
		BatchSize: batchSize,
	}
}

func (u *ParallelUpdater) Update(entities []*core.Entity, dt float64, complex bool) {
	var wg sync.WaitGroup

	for i := 0; i < len(entities); i += u.BatchSize {
		end := i + u.BatchSize
		if end > len(entities) {
			end = len(entities)
		}

		wg.Add(1)
		go func(batch []*core.Entity) {
			defer wg.Done()
			for _, e := range batch {
				if complex {
					e.ComplexUpdate(dt)
				} else {
					e.SimpleUpdate(dt)
				}
			}
		}(entities[i:end])
	}

	wg.Wait()
}
