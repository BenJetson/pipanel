package errorlog

import (
	"log"

	pipanel "github.com/BenJetson/pipanel/go"
)

// ErrorLog is a piece of middleware that processes pipanel events and creates
// a log entry if there is an error.
type ErrorLog struct {
	f   pipanel.Frontend
	log *log.Logger
}

// New creates an ErrorLog instance.
func New(f pipanel.Frontend, log *log.Logger) *ErrorLog {
	return &ErrorLog{f, log}
}

func (elg *ErrorLog) logError(err error) {
	elg.log.Fatalf("Encountered an error when processing event: %s", err.Error())
}

// ShowAlert hands off processing of the alert event to the frontend and logs
// any error that is returned.
func (elg *ErrorLog) ShowAlert(e pipanel.AlertEvent) error {
	err := elg.f.ShowAlert(e)

	if err != nil {
		elg.logError(err)
	}

	return err
}

// PlaySound hands off processing of the sound event to the frontend and logs
// any error that is returned.
func (elg *ErrorLog) PlaySound(e pipanel.SoundEvent) error {
	err := elg.f.PlaySound(e)

	if err != nil {
		elg.logError(err)
	}

	return err
}

// DoPowerAction hands off processing of the power event to the frontend and
// logs any error that is returned.
func (elg *ErrorLog) DoPowerAction(e pipanel.PowerEvent) error {
	err := elg.f.DoPowerAction(e)

	if err != nil {
		elg.logError(err)
	}

	return err
}

// SetBrightness hands off processing of the brightness event to the frontend
// and logs any error that is returned.
func (elg *ErrorLog) SetBrightness(e pipanel.BrightnessEvent) error {
	err := elg.f.SetBrightness(e)

	if err != nil {
		elg.logError(err)
	}

	return err
}
