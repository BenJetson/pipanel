package displaylog

import (
	"log"

	pipanel "github.com/BenJetson/pipanel/go"
)

type DisplayLog struct {
	log *log.Logger
}

func New(log *log.Logger) *DisplayLog {
	return &DisplayLog{log}
}

func (d *DisplayLog) SetBrightness(e pipanel.BrightnessEvent) error {
	d.log.Printf(
		"## BRIGHTNESS EVENT ##\n"+
			"Level: %d\n",
		e.Level)

	return nil
}

func (d *DisplayLog) Init() error {
	return nil
}

func (d *DisplayLog) Cleanup() error {
	return nil
}
