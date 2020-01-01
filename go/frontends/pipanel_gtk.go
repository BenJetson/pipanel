package frontends

import (
	"log"
	"time"

	pipanel "github.com/BenJetson/pipanel/go"
	"github.com/BenJetson/pipanel/go/frontends/alerters/ttsalerter"
	"github.com/faiface/beep/speaker"

	"github.com/BenJetson/pipanel/go/frontends/alerters/gtkalerter"
	"github.com/BenJetson/pipanel/go/frontends/audio_players/beeper"
	"github.com/BenJetson/pipanel/go/frontends/display_managers/displaylog"
	"github.com/BenJetson/pipanel/go/frontends/power_managers/powerlog"
)

type PiPanelGTK struct {
	*beeper.Beeper
	*gtkalerter.GUI
	*ttsalerter.TTSAlerter
	*displaylog.DisplayLog
	*powerlog.PowerLog
}

func NewPiPanelGTK(log *log.Logger) *PiPanelGTK {
	speaker.Init(beeper.SampleRate, beeper.SampleRate.N(time.Second/10))

	return &PiPanelGTK{
		GUI:        gtkalerter.New(log),
		TTSAlerter: ttsalerter.New(log, "/tmp/alert-tts/", "en"),
		Beeper:     beeper.New(log, "/tmp/pipanel-sounds/"),
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
