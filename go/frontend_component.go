package pipanel

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

// InitCleaner allows PiPanel components to perform actions prior to their first
// usage or the final teardown.
type InitCleaner interface {
	// Init performs any necessary setup that needs to be done prior to starting
	// a PiPanel component and sets the logger instance to be used.
	//
	// The raw JSON from the configuration file for this component is provided.
	// Implementors may decode this config object into the appropriate format
	// for their component.
	Init(log *logrus.Entry, cfg json.RawMessage) error
	// Cleanup performs any necessary cleanup and teardown that needs to be done
	// prior to halting a PiPanel component.
	Cleanup() error
}

// An Alerter presents alert messages to the user.
type Alerter interface {
	InitCleaner
	// ShowAlert displays an alert on the screen.
	ShowAlert(e AlertEvent) error
}

// An AudioPlayer plays audio clips stored on the system.
type AudioPlayer interface {
	InitCleaner
	// PlaySound plays the specified sound.
	PlaySound(e SoundEvent) error
}

// A PowerManager controls system power functions.
type PowerManager interface {
	InitCleaner
	// DoPowerAction performs the system power action.
	DoPowerAction(e PowerEvent) error
}

// A DisplayManager controls properties of the display.
type DisplayManager interface {
	InitCleaner
	// SetBrightness alters the brightness of the panel.
	SetBrightness(e BrightnessEvent) error
}
