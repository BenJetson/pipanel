package frontends

import (
	"log"

	"github.com/BenJetson/pipanel/go/frontends/alerters/alertlog"
	"github.com/BenJetson/pipanel/go/frontends/audio_players/audiolog"
	"github.com/BenJetson/pipanel/go/frontends/display_managers/displaylog"
	"github.com/BenJetson/pipanel/go/frontends/power_managers/powerlog"
)

type ConsoleFrontend struct {
	*alertlog.AlertLog
	*audiolog.AudioLog
	*displaylog.DisplayLog
	*powerlog.PowerLog
}

func NewConsoleFrontend(log *log.Logger) *ConsoleFrontend {
	return &ConsoleFrontend{
		AlertLog:   alertlog.New(log),
		AudioLog:   audiolog.New(log),
		DisplayLog: displaylog.New(log),
		PowerLog:   powerlog.New(log),
	}
}
