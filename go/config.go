package pipanel

import "encoding/json"

// ServerConfig contains configuration for
type ServerConfig struct {
	Port int
}

// FrontendConfig contains configuration for each of the frontend compoments.
type FrontendConfig struct {
	// Name is the name of the desired frontend implementation.
	Name string `json:"name"`
	// AlerterConfig is the raw JSON object that will be passed to the chosen
	// Alerter implementation upon instantiation.
	AlerterConfig json.RawMessage `json:"alerter,omitempty"`
	// AudioPlayerConfig is the raw JSON object that will be passed to the
	// chosen AudioPlayer implementation upon instantiation.
	AudioPlayerConfig json.RawMessage `json:"audio_player,omitempty"`
	// PowerManagerConfig is the raw JSON object that will be passed to the
	// chosen PowerManager implementation upon instantiation.
	PowerManagerConfig json.RawMessage `json:"power_manager,omitempty"`
	// DisplayManagerConfig is the raw JSON object that will be passed to the
	// chosen DisplayManager implementation upon instantiation.
	DisplayManagerConfig json.RawMessage `json:"display_manager,omitempty"`
}

// Config is the format for the program's configuration file.
type Config struct {
	// Server contains the configuration that will be passed to the PiPanel
	// server upon instantiation.
	Server ServerConfig `json:"server"`
	// Frontend contains the configuratin that will be passed to the PiPanel
	// frontend upon instantiation.
	Frontend FrontendConfig `json:"frontend"`
}
