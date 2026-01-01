package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/marwan562/CguP/pkg/compute"
	"github.com/marwan562/CguP/pkg/core"
)

func benchmark(name string, fn func()) time.Duration {
	start := time.Now()
	fn()
	duration := time.Since(start)
	fmt.Printf("%-50s %v\n", name, duration)
	return duration
}

func createEntities(count int) []*core.Entity {
	entities := make([]*core.Entity, count)
	for i := 0; i < count; i++ {
		entities[i] = &core.Entity{
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

	fmt.Println("=== GAME ENGINE PARALLELISM BENCHMARK (REFACTORED) ===\n")

	seqUpdater := &compute.SequentialUpdater{}

	fmt.Println("TEST 1: 100 entities with SIMPLE updates")
	fmt.Println("Expected: Sequential faster (goroutine overhead > work)")
	fmt.Println(strings.Repeat("-", 70))

	entities := createEntities(100)
	seq1 := benchmark("Sequential (100 simple)", func() {
		for i := 0; i < 100; i++ {
			seqUpdater.Update(entities, dt, false)
		}
	})

	parUpdater10 := compute.NewParallelUpdater(10)
	par1 := benchmark("Parallel batch=10 (100 simple)", func() {
		for i := 0; i < 100; i++ {
			parUpdater10.Update(entities, dt, false)
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
			seqUpdater.Update(entities, dt, false)
		}
	})

	parUpdater1000 := compute.NewParallelUpdater(1000)
	par2 := benchmark("Parallel batch=1000 (10k simple)", func() {
		for i := 0; i < 10; i++ {
			parUpdater1000.Update(entities, dt, false)
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
			seqUpdater.Update(entities, dt, true)
		}
	})

	parUpdater50 := compute.NewParallelUpdater(50)
	par3_small := benchmark("Parallel batch=50 (1k complex)", func() {
		for i := 0; i < 10; i++ {
			parUpdater50.Update(entities, dt, true)
		}
	})

	parUpdater200 := compute.NewParallelUpdater(200)
	par3_large := benchmark("Parallel batch=200 (1k complex)", func() {
		for i := 0; i < 10; i++ {
			parUpdater200.Update(entities, dt, true)
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
		seqUpdater.Update(entities, dt, true)
	})

	par4 := benchmark("Parallel batch=1000 (50k complex)", func() {
		parUpdater1000.Update(entities, dt, true)
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
