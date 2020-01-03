package audiolog

import (
	"log"

	pipanel "github.com/BenJetson/pipanel/go"
)

type AudioLog struct {
	log *log.Logger
}

func New() *AudioLog { return &AudioLog{} }

func (a *AudioLog) PlaySound(e pipanel.SoundEvent) error {
	a.log.Printf(
		"## SOUND EVENT ##\n"+
			"Sound: %s\n",
		e.Sound)

	return nil
}

func (a *AudioLog) Init(log *log.Logger) error {
	a.log = log
	return nil
}

func (a *AudioLog) Cleanup() error {
	return nil
}
