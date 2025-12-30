package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

type Entity struct {
	Position [3]float64
	Velocity [3]float64
	Data     [100]float64
}

func (e *Entity) complexUpdate(dt float64) {
	for i := 0; i < 3; i++ {
		e.Position[i] += e.Velocity[i] * dt
	}

	for i := 0; i < 100; i++ {
		e.Data[i] = math.Sin(e.Data[i]) * math.Cos(e.Position[0])
		e.Data[i] = math.Sqrt(math.Abs(e.Data[i])) + 0.001
	}
}

func (e *Entity) simpleUpdate(dt float64) {
	for i := 0; i < 3; i++ {
		e.Position[i] += e.Velocity[i] * dt
	}
}

func updateSequential(entities []*Entity, dt float64, complex bool) {
	for _, e := range entities {
		if complex {
			e.complexUpdate(dt)
		} else {
			e.simpleUpdate(dt)
		}
	}
}

func updateParallel(entities []*Entity, dt float64, batchSize int, complex bool) {
	var wg sync.WaitGroup

	for i := 0; i < len(entities); i += batchSize {
		end := i + batchSize
		if end > len(entities) {
			end = len(entities)
		}

		wg.Add(1)
		go func(batch []*Entity) {
			defer wg.Done()
			for _, e := range batch {
				if complex {
					e.complexUpdate(dt)
				} else {
					e.simpleUpdate(dt)
				}
			}
		}(entities[i:end])
	}

	wg.Wait()
}

func benchmark(name string, fn func()) time.Duration {
	start := time.Now()
	fn()
	duration := time.Since(start)
	fmt.Printf("%-50s %v\n", name, duration)
	return duration
}

func createEntities(count int) []*Entity {
	entities := make([]*Entity, count)
	for i := 0; i < count; i++ {
		entities[i] = &Entity{
			Position: [3]float64{rand.Float64(), rand.Float64(), rand.Float64()},
			Velocity: [3]float64{rand.Float64(), rand.Float64(), rand.Float64()},
		}
		for j := 0; j < 100; j++ {
			entities[i].Data[j] = rand.Float64()
		}
	}
	return entities
}

func main() {
	rand.Seed(time.Now().UnixNano())
	dt := 0.016

	fmt.Println("=== GAME ENGINE PARALLELISM BENCHMARK ===\n")

	fmt.Println("TEST 1: 100 entities with SIMPLE updates")
	fmt.Println("Expected: Sequential faster (goroutine overhead > work)")
	fmt.Println(strings.Repeat("-", 70))

	entities := createEntities(100)
	seq1 := benchmark("Sequential (100 simple)", func() {
		for i := 0; i < 100; i++ {
			updateSequential(entities, dt, false)
		}
	})

	par1 := benchmark("Parallel batch=10 (100 simple)", func() {
		for i := 0; i < 100; i++ {
			updateParallel(entities, dt, 10, false)
		}
	})

	speedup := float64(seq1) / float64(par1)
	fmt.Printf("Speedup: %.2fx %s\n\n", speedup, getVerdict(speedup))

	fmt.Println("TEST 2: 10,000 entities with SIMPLE updates")
	fmt.Println("Expected: Parallel slightly faster")
	fmt.Println(strings.Repeat("-", 70))

	entities = createEntities(10000)
	seq2 := benchmark("Sequential (10k simple)", func() {
		for i := 0; i < 10; i++ {
			updateSequential(entities, dt, false)
		}
	})

	par2 := benchmark("Parallel batch=1000 (10k simple)", func() {
		for i := 0; i < 10; i++ {
			updateParallel(entities, dt, 1000, false)
		}
	})

	speedup = float64(seq2) / float64(par2)
	fmt.Printf("Speedup: %.2fx %s\n\n", speedup, getVerdict(speedup))

	fmt.Println("TEST 3: 1,000 entities with COMPLEX updates (CPU-heavy)")
	fmt.Println("Expected: Parallel MUCH faster (real CPU work)")
	fmt.Println(strings.Repeat("-", 70))

	entities = createEntities(1000)
	seq3 := benchmark("Sequential (1k complex)", func() {
		for i := 0; i < 10; i++ {
			updateSequential(entities, dt, true)
		}
	})

	par3_small := benchmark("Parallel batch=50 (1k complex)", func() {
		for i := 0; i < 10; i++ {
			updateParallel(entities, dt, 50, true)
		}
	})

	par3_large := benchmark("Parallel batch=200 (1k complex)", func() {
		for i := 0; i < 10; i++ {
			updateParallel(entities, dt, 200, true)
		}
	})

	speedup1 := float64(seq3) / float64(par3_small)
	speedup2 := float64(seq3) / float64(par3_large)
	fmt.Printf("Speedup (batch=50):  %.2fx %s\n", speedup1, getVerdict(speedup1))
	fmt.Printf("Speedup (batch=200): %.2fx %s\n\n", speedup2, getVerdict(speedup2))

	fmt.Println("TEST 4: 50,000 entities with COMPLEX updates")
	fmt.Println("Expected: Parallel SIGNIFICANTLY faster")
	fmt.Println(strings.Repeat("-", 70))

	entities = createEntities(50000)
	seq4 := benchmark("Sequential (50k complex)", func() {
		updateSequential(entities, dt, true)
	})

	par4 := benchmark("Parallel batch=1000 (50k complex)", func() {
		updateParallel(entities, dt, 1000, true)
	})

	speedup = float64(seq4) / float64(par4)
	fmt.Printf("Speedup: %.2fx %s\n\n", speedup, getVerdict(speedup))

	// Summary
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("SUMMARY:")
	fmt.Println("• Small workloads: Parallel SLOWER (overhead > benefit)")
	fmt.Println("• Simple operations: Parallel slightly better with many entities")
	fmt.Println("• Complex CPU-bound work: Parallel 2-4x faster!")
	fmt.Println("• Batch size matters: Too small = overhead, too large = less parallelism")
	fmt.Println(strings.Repeat("=", 70))
}

func getVerdict(speedup float64) string {
	if speedup < 0.8 {
		return "SLOWER (overhead too high)"
	} else if speedup < 1.2 {
		return "SIMILAR (marginal difference)"
	} else if speedup < 2.0 {
		return "FASTER (good speedup)"
	} else {
		return "MUCH FASTER (excellent speedup)"
	}
}

// Helper to avoid import
var strings = struct {
	Repeat func(string, int) string
}{
	Repeat: func(s string, count int) string {
		result := ""
		for i := 0; i < count; i++ {
			result += s
		}
		return result
	},
}
