package pipanel

type Config struct {
	Application ApplicationConfig `json:"application"`
	Server      ServerConfig      `json:"server"`
}

type ApplicationConfig struct {
	Alert struct {
		Enabled  bool `json:"enabled"`
		FontSize int  `json:"font_size"`
		IconSize int  `json:"icon_size"`
	} `json:"alert"`
	Sound struct {
		Enabled bool              `json:"enabled"`
		Library map[string]string `json:"library"`
	} `json:"sound"`
	SystemAction struct {
		Enabled  bool              `json:"enabled"`
		Commands map[string]string `json:"commands"`
	} `json:"system_action"`
	Brightness struct {
		Enabled bool `json:"enabled"`
		Min     int  `json:"min"`
		Max     int  `json:"max"`
	} `json:"brightness"`
}

type ServerConfig struct {
	Port int `json:"port"`
}
