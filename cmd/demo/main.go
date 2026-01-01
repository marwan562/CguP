package main

import (
	"github.com/marwan562/CguP/engine/core"
	"github.com/marwan562/CguP/engine/gpu"
	"github.com/marwan562/CguP/engine/platform"
)

func main() {
	// 1. Init Config & Logs
	core.InitLogger()
	config := core.DefaultConfig()

	// 2. Create Window
	window, err := platform.NewWindow(config)
	if err != nil {
		core.LogFatal("Failed to create window: %v", err)
	}
	defer window.Destroy()
	core.LogInfo("Window created: %dx%d", config.WindowWidth, config.WindowHeight)

	// 3. Init Vulkan
	vkCtx, err := gpu.NewVulkanContext(window.Handle, config.AppName)
	if err != nil {
		core.LogFatal("Failed to initialize Vulkan: %v", err)
	}
	defer vkCtx.Destroy()

	// 4. Main Loop
	core.LogInfo("Starting Main Loop...")
	for !window.ShouldClose() {
		window.PollEvents()
		// Rendering happens here later
	}
	core.LogInfo("Exiting...")
}
