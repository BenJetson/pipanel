package console

import (
	"log"

	pipanel "github.com/BenJetson/pipanel/go"
)

type ConsoleFrontend struct {
	log *log.Logger
}

func New(l *log.Logger) *ConsoleFrontend {
	return &ConsoleFrontend{l}
}

func (c *ConsoleFrontend) ShowAlert(e pipanel.AlertEvent) error {
	c.log.Printf("## ALERT EVENT ##\nMessage: %s\nTimeout: %d\n"+
		"AutoDismiss: %t\nIcon:%s\n",
		e.Message, e.Timeout, e.Perpetual, e.Icon)
	return nil
}

func (c *ConsoleFrontend) PlaySound(e pipanel.SoundEvent) error {
	c.log.Printf("## SOUND EVENT ##\nTone: %s\n", e.Tone)
	return nil
}

func (c *ConsoleFrontend) DoPowerAction(e pipanel.PowerEvent) error {
	c.log.Printf("## POWER EVENT ##\nAction: %s\n", e.Action)
	return nil
}

func (c *ConsoleFrontend) SetBrightness(e pipanel.BrightnessEvent) error {
	c.log.Printf("## BRIGHTNESS EVENT ##\nLevel: %d\n", e.Level)
	return nil
}
