package powerlog

import (
	"encoding/json"
	"log"

	pipanel "github.com/BenJetson/pipanel/go"
)

// PowerLog implements pipanel.PowerManager and handles power events by writing
// the details to the console. Useful for testing purposes.
type PowerLog struct {
	log *log.Logger
}

// New creates a fresh PowerLog instance.
func New() *PowerLog { return &PowerLog{} }

// DoPowerAction handles power events by writing the details to the console.
func (p *PowerLog) DoPowerAction(e pipanel.PowerEvent) error {
	p.log.Printf(
		"## POWER EVENT ##\n"+
			"Action: %s\n",
		e.Action)

	return nil
}

// Init initializes this PowerLog by setting the logger.
func (p *PowerLog) Init(log *log.Logger, _ json.RawMessage) error {
	p.log = log
	return nil
}

// Cleanup tears down this DisplayLog.
func (p *PowerLog) Cleanup() error { return nil }
