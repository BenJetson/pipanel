package frontends

import (
	"log"

	"github.com/BenJetson/pipanel/go/frontends/alerters/gtkalerter"
	"github.com/BenJetson/pipanel/go/frontends/audio_players/audiolog"
	"github.com/BenJetson/pipanel/go/frontends/display_managers/displaylog"
	"github.com/BenJetson/pipanel/go/frontends/power_managers/powerlog"
)

type PiPanelGTK struct {
	*audiolog.AudioLog
	*gtkalerter.GUI
	*displaylog.DisplayLog
	*powerlog.PowerLog
}

func NewPiPanelGTK(log *log.Logger) *PiPanelGTK {
	return &PiPanelGTK{
		GUI:        gtkalerter.New(log),
		AudioLog:   audiolog.New(log),
		DisplayLog: displaylog.New(log),
		PowerLog:   powerlog.New(log),
	}
}
