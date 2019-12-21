package audiolog

import (
	"log"

	pipanel "github.com/BenJetson/pipanel/go"
)

type AudioLog struct {
	log *log.Logger
}

func New(log *log.Logger) *AudioLog {
	return &AudioLog{log}
}

func (a *AudioLog) PlaySound(e pipanel.SoundEvent) error {
	a.log.Printf(
		"## SOUND EVENT ##\n"+
			"Sound: %s\n",
		e.Sound)

	return nil
}
