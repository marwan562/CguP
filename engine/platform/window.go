package platform

import (
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/marwan562/CguP/engine/core"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

type Window struct {
	Handle *glfw.Window
	Width  int
	Height int
}

func NewWindow(config *core.EngineConfig) (*Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, err
	}

	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI) // Prepare for Vulkan
	glfw.WindowHint(glfw.Resizable, glfw.False)

	handle, err := glfw.CreateWindow(config.WindowWidth, config.WindowHeight, config.AppName, nil, nil)
	if err != nil {
		return nil, err
	}

	return &Window{
		Handle: handle,
		Width:  config.WindowWidth,
		Height: config.WindowHeight,
	}, nil
}

func (w *Window) ShouldClose() bool {
	return w.Handle.ShouldClose()
}

func (w *Window) PollEvents() {
	glfw.PollEvents()
}

func (w *Window) Destroy() {
	w.Handle.Destroy()
	glfw.Terminate()
}
