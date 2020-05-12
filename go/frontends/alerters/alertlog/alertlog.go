package alertlog

import (
	"encoding/json"
	"log"

	pipanel "github.com/BenJetson/pipanel/go"
)

type AlertLog struct {
	log *log.Logger
}

func New() *AlertLog { return &AlertLog{} }

func (a *AlertLog) ShowAlert(e pipanel.AlertEvent) error {
	a.log.Printf(
		"## ALERT EVENT ##\n"+
			"Message: %s\n"+
			"Timeout: %d\n"+
			"AutoDismiss: %t\n"+
			"Icon:%s\n",
		e.Message, e.Timeout, e.Perpetual, e.Icon)

	return nil
}

func (a *AlertLog) Init(log *log.Logger, _ json.RawMessage) error {
	a.log = log
	return nil
}

func (a *AlertLog) Cleanup() error {
	return nil
}
