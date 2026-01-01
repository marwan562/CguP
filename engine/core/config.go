package core

type EngineConfig struct {
	AppName      string
	WindowWidth  int
	WindowHeight int
	EnableDebug  bool
}

func DefaultConfig() *EngineConfig {
	return &EngineConfig{
		AppName:      "CguP Engine",
		WindowWidth:  800,
		WindowHeight: 600,
		EnableDebug:  true,
	}
}
