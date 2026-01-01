# CguP Architecture

CguP is designed as a modular, concurrent rendering engine.

## High-Level Modules

```mermaid
graph TD
    App[Game Application] --> Engine
    Engine --> Platform[Platform (GLFW)]
    Engine --> GPU[GPU (Vulkan)]
    Engine --> Core[Core (Log/Config)]
    Engine --> ECS[ECS (Entities)]
    
    GPU --> Vulkan[Vulkan SDK]
```

## Directory Structure

- **`engine/core`**: Utilities, logging, configuration.
- **`engine/platform`**: Window management, input handling (Abstracts GLFW).
- **`engine/gpu`**: Low-level Vulkan wrappers.
    - `vulkan_context.go`: Instance/Device creation.
    - Future: `swapchain.go`, `pipeline.go`.
- **`engine/ecs`**: (Planned) Data-oriented entity system.
- **`engine/renderer`**: (Planned) Render graph and pass management.

## Concurrency Model

CguP leverages Go's goroutines for:
- **Command Recording**: Parallel construction of `vk.CommandBuffer`.
- **Asset Loading**: Background streaming of textures/meshes.
- **Simulation**: Parallel entity updates (see `pkg/compute` from prototype).
