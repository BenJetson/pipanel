package frontends

import (
	"log"

	pipanel "github.com/BenJetson/pipanel/go"
	"github.com/BenJetson/pipanel/go/frontends/alerters/ttsalerter"

	"github.com/BenJetson/pipanel/go/frontends/alerters/gtkalerter"
	"github.com/BenJetson/pipanel/go/frontends/audio_players/audiolog"
	"github.com/BenJetson/pipanel/go/frontends/display_managers/displaylog"
	"github.com/BenJetson/pipanel/go/frontends/power_managers/powerlog"
)

type PiPanelGTK struct {
	*audiolog.AudioLog
	*gtkalerter.GUI
	*ttsalerter.TTSAlerter
	*displaylog.DisplayLog
	*powerlog.PowerLog
}

func NewPiPanelGTK(log *log.Logger) *PiPanelGTK {
	return &PiPanelGTK{
		GUI:        gtkalerter.New(log),
		TTSAlerter: ttsalerter.New(log, "/tmp/alert-tts/", "en"),
		AudioLog:   audiolog.New(log),
		DisplayLog: displaylog.New(log),
		PowerLog:   powerlog.New(log),
	}
}

// ShowAlert will present both the GTK alert and read the alert out loud.
func (p *PiPanelGTK) ShowAlert(e pipanel.AlertEvent) error {
	err := p.GUI.ShowAlert(e)

	if err != nil {
		return err
	}

	return p.TTSAlerter.ShowAlert(e)
}
