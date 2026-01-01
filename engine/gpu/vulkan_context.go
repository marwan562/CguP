package gpu

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/go-gl/glfw/v3.3/glfw"
	vk "github.com/goki/vulkan"
	"github.com/marwan562/CguP/engine/core"
)

type VulkanContext struct {
	Instance       vk.Instance
	PhysicalDevice vk.PhysicalDevice
	Device         vk.Device

	// Debug
	DebugCallback vk.DebugReportCallback
}

func NewVulkanContext(window *glfw.Window, appName string) (*VulkanContext, error) {
	// 1. Enable MoltenVK on macOS
	if runtime.GOOS == "darwin" {
		enableMoltenVK()
	}

	// 2. Set Instance Extensions
	requiredExtensions := getRequiredExtensions(window)

	// 3. Set Layers
	var validationLayers []string
	// validationLayers := []string{"VK_LAYER_KHRONOS_validation"}

	// Check layer support (simplified for now, assume they exist or we fail gracefully if not found)
	// In production we would check vk.EnumerateInstanceLayerProperties

	// 4. Create Instance
	appInfo := vk.ApplicationInfo{
		SType:              vk.StructureTypeApplicationInfo,
		PApplicationName:   appName + "\x00",
		ApplicationVersion: vk.MakeVersion(1, 0, 0),
		PEngineName:        "CguP Engine\x00",
		EngineVersion:      vk.MakeVersion(1, 0, 0),
		ApiVersion:         vk.MakeVersion(1, 2, 0),
	}

	var flags vk.InstanceCreateFlags
	if runtime.GOOS == "darwin" {
		flags |= vk.InstanceCreateFlags(vk.InstanceCreateEnumeratePortabilityBit)
	}

	core.LogInfo("Active Extensions: %v", requiredExtensions)

	instanceCreateInfo := vk.InstanceCreateInfo{
		SType:                   vk.StructureTypeInstanceCreateInfo,
		PApplicationInfo:        &appInfo,
		EnabledExtensionCount:   uint32(len(requiredExtensions)),
		PpEnabledExtensionNames: requiredExtensions,
		EnabledLayerCount:       uint32(len(validationLayers)),
		PpEnabledLayerNames:     validationLayers,
		Flags:                   flags,
	}

	var instance vk.Instance
	if err := vk.Error(vk.CreateInstance(&instanceCreateInfo, nil, &instance)); err != nil {
		core.LogError("Failed to create Vulkan instance: %v", err)
		return nil, err
	}
	vk.InitInstance(instance)
	core.LogInfo("Vulkan Instance created successfully.")

	// 5. Select Physical Device
	physicalDevice, err := pickPhysicalDevice(instance)
	if err != nil {
		vk.DestroyInstance(instance, nil)
		return nil, err
	}

	return &VulkanContext{
		Instance:       instance,
		PhysicalDevice: physicalDevice,
	}, nil
}

func (ctx *VulkanContext) Destroy() {
	var emptyInstance vk.Instance
	if ctx.Instance != emptyInstance {
		vk.DestroyInstance(ctx.Instance, nil)
		ctx.Instance = emptyInstance
	}
}

func pickPhysicalDevice(instance vk.Instance) (vk.PhysicalDevice, error) {
	var emptyPhysicalDevice vk.PhysicalDevice
	var count uint32
	if err := vk.Error(vk.EnumeratePhysicalDevices(instance, &count, nil)); err != nil {
		return emptyPhysicalDevice, err
	}

	if count == 0 {
		return emptyPhysicalDevice, fmt.Errorf("failed to find GPUs with Vulkan support")
	}

	devices := make([]vk.PhysicalDevice, count)
	if err := vk.Error(vk.EnumeratePhysicalDevices(instance, &count, devices)); err != nil {
		return emptyPhysicalDevice, err
	}

	// Just pick the first one for now
	device := devices[0]

	var props vk.PhysicalDeviceProperties
	vk.GetPhysicalDeviceProperties(device, &props)
	props.Deref()

	// Convert null-terminated byte array to string
	deviceName := strings.TrimRight(string(props.DeviceName[:]), "\x00")
	core.LogInfo("Selected GPU: %s", deviceName)

	return device, nil
}

func getRequiredExtensions(window *glfw.Window) []string {
	glfwExtensions := window.GetRequiredInstanceExtensions()

	// Standardize extensions for Vulkan-Go
	var extensions []string
	for _, ext := range glfwExtensions {
		extensions = append(extensions, ext)
	}

	// Add MacOS specific extensions
	if runtime.GOOS == "darwin" {
		extensions = append(extensions, "VK_KHR_portability_enumeration")
		extensions = append(extensions, "VK_KHR_get_physical_device_properties2")
		// Note: We might need to set the flag VK_INSTANCE_CREATE_ENUMERATE_PORTABILITY_BIT_KHR
		// but vulkan-go struct might handle it or we need to look into how to pass flags if needed.
		// For now let's hope the MoltenVK enable function covers enough.
	}

	// Add Debug utils if we want them (VK_EXT_debug_utils or VK_EXT_debug_report)
	// extensions = append(extensions, "VK_EXT_debug_report")

	return extensions
}

func enableMoltenVK() {
	core.LogInfo("Enabling MoltenVK...")
	vk.SetGetInstanceProcAddr(glfw.GetVulkanGetInstanceProcAddress())
	if err := vk.Init(); err != nil {
		core.LogError("Failed to initialize Vulkan bindings: %v", err)
	}
	core.LogInfo("Vulkan bindings initialized.")
}
