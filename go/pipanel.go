package pipanel

import "time"

// InitCleaner allows PiPanel components to perform actions prior to their first
// usage or the final teardown.
type InitCleaner interface {
	// Init performs any necessary setup that needs to be done prior to starting
	// a PiPanel component.
	Init() error
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

// Frontend provides abstraction for processing PiPanel events.
type Frontend struct {
	InitCleaner
	Alerter
	AudioPlayer
	PowerManager
	DisplayManager
}

// Init initializes all components of the Frontend.
func (f *Frontend) Init() {
	if f.InitCleaner != nil {
		f.InitCleaner.Init()
	}

	if f.Alerter != nil {
		f.Alerter.Init()
	}

	if f.AudioPlayer != nil {
		f.AudioPlayer.Init()
	}

	if f.PowerManager != nil {
		f.PowerManager.Init()
	}

	if f.DisplayManager != nil {
		f.DisplayManager.Init()
	}
}

// Cleanup tears down all components of the Frontend.
func (f *Frontend) Cleanup() {
	if f.InitCleaner != nil {
		f.InitCleaner.Cleanup()
	}

	if f.Alerter != nil {
		f.Alerter.Cleanup()
	}

	if f.AudioPlayer != nil {
		f.AudioPlayer.Cleanup()
	}

	if f.PowerManager != nil {
		f.PowerManager.Cleanup()
	}

	if f.DisplayManager != nil {
		f.DisplayManager.Cleanup()
	}
}

// An AlertEvent contains information about an alert display request.
type AlertEvent struct {
	SoundEvent
	// Message is the text content of an alert.
	Message string `json:"message"`
	// Timeout is the number of seconds until an alert will auto-dismiss.
	// This field is ignored when Perpetual == true.
	Timeout time.Duration `json:"timeout"`
	// Perpetual is true when the alert should not automatically dismiss.
	Perpetual bool `json:"perpetual"`
	// Icon is the name of a gtk icon that should be displayed on-screen.
	// If this is blank, a default icon is used.
	Icon string `json:"icon"`
	// // FontSize is the font size that should be used for the alert text content,
	// // measured in in points.
	// FontSize int `json:"font_size"`
}

// A SoundEvent contains information about a tone that will be played on the panel.
type SoundEvent struct {
	// Sound is the name of the sound file to be played. Path is relative to
	// the sound folder configured in the preference file. Empty string will
	// result in no sound being played.
	Sound string `json:"sound"`
}

// PowerAction describes a system power action to be taken by the panel.
type PowerAction string

const (
	// PowerActionShutdown instructs the panel to shut down.
	PowerActionShutdown PowerAction = "shutdown"
	// PowerActionReboot instructs the panel to reboot.
	PowerActionReboot PowerAction = "reboot"
	// PowerActionDisplayOff instructs the panel to blank the display.
	PowerActionDisplayOff PowerAction = "displayOff"
)

// A PowerEvent contains information about a system power request.
type PowerEvent struct {
	// Action is the power action that should be performed by the panel.
	Action PowerAction `json:"action"`
}

// A BrightnessEvent contains information about a brightness change request.
type BrightnessEvent struct {
	// Level is the level that the brightness of the panel should be set to.
	// This must be on the range [0,255].
	Level uint8 `json:"level"`
}
