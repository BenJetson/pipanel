package audiolog

import (
	"encoding/json"
	"log"

	pipanel "github.com/BenJetson/pipanel/go"
)

// AudioLog implements pipanel.AudioPlayer and handles sound events by writing
// the details to the console. Useful for testing purposes.
type AudioLog struct {
	log *log.Logger
}

// New creates a fresh AudioLog instance.
func New() *AudioLog { return &AudioLog{} }

// PlaySound handles sound events by writing the details to the console.
func (a *AudioLog) PlaySound(e pipanel.SoundEvent) error {
	a.log.Printf(
		"## SOUND EVENT ##\n"+
			"Sound: %s\n",
		e.Sound)

	return nil
}

// Init initailizes this AlertLog by setting the logger.
func (a *AudioLog) Init(log *log.Logger, _ json.RawMessage) error {
	a.log = log
	return nil
}

// Cleanup tears down this AudioLog.
func (a *AudioLog) Cleanup() error { return nil }
