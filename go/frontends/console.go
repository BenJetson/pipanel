package frontends

import (
	"log"

	pipanel "github.com/BenJetson/pipanel/go"
	"github.com/BenJetson/pipanel/go/frontends/alerters/alertlog"
	"github.com/BenJetson/pipanel/go/frontends/audio_players/audiolog"
	"github.com/BenJetson/pipanel/go/frontends/display_managers/displaylog"
	"github.com/BenJetson/pipanel/go/frontends/power_managers/powerlog"
)

// NewConsoleFrontend creates a pipanel.Frontend that writes all events to the
// provided logger instance.
func NewConsoleFrontend(log *log.Logger) *pipanel.Frontend {
	return &pipanel.Frontend{
		Alerter:        alertlog.New(log),
		AudioPlayer:    audiolog.New(log),
		DisplayManager: displaylog.New(log),
		PowerManager:   powerlog.New(log),
	}
}
