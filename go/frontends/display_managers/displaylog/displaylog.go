package displaylog

import (
	"context"
	"encoding/json"

	pipanel "github.com/BenJetson/pipanel/go"

	"github.com/sirupsen/logrus"
)

var _ pipanel.DisplayManager = (*DisplayLog)(nil)

// DisplayLog implements pipanel.DisplayManager and handles brightness events
// by writing the details to the console. Useful for testing purposes.
type DisplayLog struct {
	log *logrus.Entry
}

// New creates a fresh DisplayLog instance.
func New() *DisplayLog { return &DisplayLog{} }

// SetBrightness handles brightness events by writing the details to the
// console.
func (d *DisplayLog) SetBrightness(ctx context.Context,
	e pipanel.BrightnessEvent) error {

	d.log.WithContext(ctx).WithFields(logrus.Fields{
		"level": e.Level,
	}).Println("Received brightness event.")

	return nil
}

// Init initializes this DisplayLog by setting the logger.
func (d *DisplayLog) Init(log *logrus.Entry, _ json.RawMessage) error {
	d.log = log
	return nil
}

// Cleanup tears down this DisplayLog.
func (d *DisplayLog) Cleanup() error { return nil }
