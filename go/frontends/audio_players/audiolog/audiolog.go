package audiolog

import (
	"context"
	"encoding/json"

	"github.com/sirupsen/logrus"

	pipanel "github.com/BenJetson/pipanel/go"
)

var _ pipanel.AudioPlayer = (*AudioLog)(nil)

// AudioLog implements pipanel.AudioPlayer and handles sound events by writing
// the details to the console. Useful for testing purposes.
type AudioLog struct {
	log *logrus.Entry
}

// New creates a fresh AudioLog instance.
func New() *AudioLog { return &AudioLog{} }

// PlaySound handles sound events by writing the details to the console.
func (a *AudioLog) PlaySound(ctx context.Context, e pipanel.SoundEvent) error {
	a.log.WithContext(ctx).WithFields(logrus.Fields{
		"sound": e.Sound,
	}).Println("Received sound event.")

	return nil
}

// Init initailizes this AlertLog by setting the logger.
func (a *AudioLog) Init(log *logrus.Entry, _ json.RawMessage) error {
	a.log = log
	return nil
}

// Cleanup tears down this AudioLog.
func (a *AudioLog) Cleanup() error { return nil }
