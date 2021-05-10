package pipanel

import "time"

type Alerter interface {
	// ShowAlert displays an alert on the screen.
	ShowAlert(e AlertEvent) error
	// Shutdown closes all active alerts and prepares for application shutdown.
	Shutdown() error
}

type Sounder interface {
	// PlaySound plays the specified sound.
	PlaySound(e SoundEvent) error
	// Stop causes all currently playing sounds to stop immediately.
	Stop() error
	// Shutdown stops all currently playing sounds and prepares for shutdown.
	Shutdown() error
}

type Commander interface {
	// DoPowerAction performs the system power action.
	DoSystemAction(e PowerEvent) error
}

// BrigtnessController handles changing the brightness of the display panel.
//   Further, it should also handle brightening the screen when an alert is
//   triggered. A counter should be kept of the number of active alerts, and
//   when that counter reaches zero, dim the display back to the previous
//   brightness level set.
type BrightnessController interface {
	// SetBrightness alters the brightness of the panel.
	SetBrightness(e BrightnessEvent) error
	// AlertAdd should be called upon displaying an alert to the user.
	AlertAdd()
	// AlertDismissed should be called when an alert is dimsmissed by the user
	//   or when it times out.
	AlertDismiss()
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
	Level int `json:"level"`
}
