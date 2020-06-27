package powerlog

import (
	"context"
	"encoding/json"

	pipanel "github.com/BenJetson/pipanel/go"

	"github.com/sirupsen/logrus"
)

var _ pipanel.PowerManager = (*PowerLog)(nil)

// PowerLog implements pipanel.PowerManager and handles power events by writing
// the details to the console. Useful for testing purposes.
type PowerLog struct {
	log *logrus.Entry
}

// New creates a fresh PowerLog instance.
func New() *PowerLog { return &PowerLog{} }

// DoPowerAction handles power events by writing the details to the console.
func (p *PowerLog) DoPowerAction(ctx context.Context, e pipanel.PowerEvent) error {
	p.log.WithContext(ctx).WithFields(logrus.Fields{
		"action": e.Action,
	}).Println("Received power action event.")

	return nil
}

// Init initializes this PowerLog by setting the logger.
func (p *PowerLog) Init(log *logrus.Entry, _ json.RawMessage) error {
	p.log = log
	return nil
}

// Cleanup tears down this DisplayLog.
func (p *PowerLog) Cleanup() error { return nil }
