# CguP - Concurrent Game Update Package

A high-performance, concurrent entity update system for Go game engines.

## Overview

This package provides a modular way to handle entity updates, specifically designed to leverage multi-core CPUs for complex simulation logic.

## Structure

- `pkg/core`: Core entity definitions.
- `pkg/compute`: Update strategies (Sequential and Parallel).
- `cmd/benchmark`: Performance benchmarking tool.

## Usage

```go
import (
    "github.com/marwan562/CguP/pkg/core"
    "github.com/marwan562/CguP/pkg/compute"
)

// Create entities
entities := make([]*core.Entity, 1000)
// ... init entities

// Create a parallel updater with batch size 100
updater := compute.NewParallelUpdater(100)

// Update loop
updater.Update(entities, 0.016, true)
```

## Benchmarks

Run the included benchmark to see the speedup on your machine:

```bash
go run cmd/benchmark/main.go
```
