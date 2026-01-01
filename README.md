# CguP - Concurrent Go Rendering Engine

<div align="center">

[![Go](https://github.com/marwan562/CguP/actions/workflows/go.yml/badge.svg)](https://github.com/marwan562/CguP/actions/workflows/go.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**A modern, data-oriented rendering engine written in Go, powered by Vulkan.**
Build for learning, research, and high-performance indie games.

[Contributing](CONTRIBUTING.md) • [Architecture](ARCHITECTURE.md) • [Roadmap](ROADMAP.md)

</div>

## Vision

CguP aims to demonstrate how Go's currency primitives (Goroutines, Channels) can be effectively used in a high-performance graphics engine. We prioritize **architecture** and **parallelism** over raw single-threaded speed.

## Features (Planned)

- 🚀 **Vulkan Backend**: Modern, explicit GPU control.
- 🧵 **Concurrent Architecture**: Parallel command recording and asset loading.
- 🧩 **Modular Design**: Clear separation of Core, Platform, and GPU layers.
- 📄 **Render Graph**: Automatic resource barrier management and pass reordering.

## Getting Started

### Prerequisites
- **Go 1.25+**
- **Vulkan SDK** (Required for MoltenVK on macOS)

### Running the Demo
```bash
# On macOS
go mod tidy
CGO_CFLAGS="-I/opt/homebrew/include" CGO_LDFLAGS="-L/opt/homebrew/lib" go run cmd/demo/main.go
```

## Structure

- `engine/core`: Basic utilities (Logging, Config).
- `engine/platform`: Windowing and Input (GLFW).
- `engine/gpu`: Vulkan wrappers and context management.
- `pkg/compute`: (Legacy) CPU-bound entity update prototypes.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
