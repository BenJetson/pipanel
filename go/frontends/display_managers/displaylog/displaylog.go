package displaylog

import (
	"log"

	pipanel "github.com/BenJetson/pipanel/go"
)

type DisplayLog struct {
	log *log.Logger
}

func New() *DisplayLog { return &DisplayLog{} }

func (d *DisplayLog) SetBrightness(e pipanel.BrightnessEvent) error {
	d.log.Printf(
		"## BRIGHTNESS EVENT ##\n"+
			"Level: %d\n",
		e.Level)

	return nil
}

func (d *DisplayLog) Init(log *log.Logger, _ json.RawMessage) error {
	d.log = log
	return nil
}

func (d *DisplayLog) Cleanup() error {
	return nil
}
