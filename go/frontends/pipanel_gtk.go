package frontends

import (
	"log"

	pipanel "github.com/BenJetson/pipanel/go"

	"github.com/BenJetson/pipanel/go/frontends/alerters/gtkttsalerter"
	"github.com/BenJetson/pipanel/go/frontends/audio_players/beeper"
	"github.com/BenJetson/pipanel/go/frontends/display_managers/displaylog"
	"github.com/BenJetson/pipanel/go/frontends/power_managers/powerlog"
)

// NewPiPanelGTK creates a pipanel.Frontend that supports the RPi official
// touch display, includes GTK/TTS alerts, and systemd power management.
func NewPiPanelGTK(log *log.Logger) *pipanel.Frontend {
	return &pipanel.Frontend{
		Alerter:        gtkttsalerter.New(log, "/tmp/alert-tts/", "en"),
		AudioPlayer:    beeper.New(log, "/tmp/pipanel-sounds/"),
		DisplayManager: displaylog.New(log),
		PowerManager:   powerlog.New(log),
	}
}
