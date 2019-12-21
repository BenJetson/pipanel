package alertlog

import (
	"log"

	pipanel "github.com/BenJetson/pipanel/go"
)

type AlertLog struct {
	log *log.Logger
}

func New(log *log.Logger) *AlertLog {
	return &AlertLog{log}
}

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
