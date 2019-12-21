package powerlog

import (
	"log"

	pipanel "github.com/BenJetson/pipanel/go"
)

type PowerLog struct {
	log *log.Logger
}

func New(log *log.Logger) *PowerLog {
	return &PowerLog{log}
}

func (p *PowerLog) DoPowerAction(e pipanel.PowerEvent) error {
	p.log.Printf(
		"## POWER EVENT ##\n"+
			"Action: %s\n",
		e.Action)

	return nil
}
