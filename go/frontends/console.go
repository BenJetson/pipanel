package frontends

import (
	pipanel "github.com/BenJetson/pipanel/go"
	"github.com/BenJetson/pipanel/go/frontends/alerters/alertlog"
	"github.com/BenJetson/pipanel/go/frontends/audio_players/audiolog"
	"github.com/BenJetson/pipanel/go/frontends/display_managers/displaylog"
	"github.com/BenJetson/pipanel/go/frontends/power_managers/powerlog"
)

// NewConsoleFrontend creates a pipanel.Frontend that writes all events to the
// provided logger instance.
func NewConsoleFrontend() *pipanel.Frontend {
	return &pipanel.Frontend{
		Alerter:        alertlog.New(),
		AudioPlayer:    audiolog.New(),
		DisplayManager: displaylog.New(),
		PowerManager:   powerlog.New(),
	}
}
