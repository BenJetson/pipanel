package alertlog

import (
	"encoding/json"

	"github.com/sirupsen/logrus"

	pipanel "github.com/BenJetson/pipanel/go"
)

var _ pipanel.Alerter = (*AlertLog)(nil)

// AlertLog implements pipanel.Alerter and handles alert events by writing the
// details to the console. Useful for testing purposes.
type AlertLog struct {
	log *logrus.Entry
}

// New creats a fresh AlertLog instance.
func New() *AlertLog { return &AlertLog{} }

// ShowAlert handles alert events by writing the details to the console.
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

// Init initializes this AlertLog by setting the logger.
func (a *AlertLog) Init(log *logrus.Entry, _ json.RawMessage) error {
	a.log = log
	return nil
}

// Cleanup tears down this AlertLog.
func (a *AlertLog) Cleanup() error { return nil }
