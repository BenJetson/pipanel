package pipanel

import "time"

// Frontend provides abstraction for processing PiPanel events.
type Frontend interface {
	// ShowAlert displays an alert on the screen.
	ShowAlert(e AlertEvent) error
	// PlaySound plays the specified sound.
	PlaySound(e SoundEvent) error
	// DoPowerAction performs the system power action.
	DoPowerAction(e PowerEvent) error
	// SetBrightness alters the brightness of the panel.
	SetBrightness(e BrightnessEvent) error
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
	// Tone is the name of the sound file to be played. Path is relative to
	// the sound folder configured in the preference file. Empty string will
	// result in no sound being played.
	Tone string `json:"tone"`
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
	Level int `json:"level"`
}
