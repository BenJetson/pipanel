package powerlog

import (
	"log"

	pipanel "github.com/BenJetson/pipanel/go"
)

type PowerLog struct {
	log *log.Logger
}

func New() *PowerLog { return &PowerLog{} }

func (p *PowerLog) DoPowerAction(e pipanel.PowerEvent) error {
	p.log.Printf(
		"## POWER EVENT ##\n"+
			"Action: %s\n",
		e.Action)

	return nil
}

func (p *PowerLog) Init(log *log.Logger, _ json.RawMessage) error {
	p.log = log
	return nil
}

func (p *PowerLog) Cleanup() error {
	return nil
}
