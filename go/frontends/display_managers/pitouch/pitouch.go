package pitouch

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	pipanel "github.com/BenJetson/pipanel/go"
)

var _ pipanel.DisplayManager = (*TouchDisplayManager)(nil)

const brightFile string = "/sys/class/backlight/rpi_backlight/brightness"

// TouchDisplayManager implements pipanel.DisplayManager for the Raspberry Pi
// official 7" touchscreen device.
type TouchDisplayManager struct {
	log *logrus.Entry
}

// New creates a TouchDisplayManager instance.
func New() *TouchDisplayManager { return &TouchDisplayManager{} }

// SetBrightness handles pipanel brightness events.
func (t *TouchDisplayManager) SetBrightness(e pipanel.BrightnessEvent) error {
	// Setting the brightness less than ten will cause the screen to blank.
	if e.Level < 10 {
		return errors.New("device does not support brigtness values < 10")
	}

	// Open brightness register file in write-only mode.
	f, err := os.OpenFile(brightFile, os.O_WRONLY, 0666)

	if err != nil {
		return errors.Wrap(err, "could not open brightness register file")
	}

	t.log.Printf("Setting RPi touchscreen brightness to %d.", e.Level)

	if _, err = f.WriteString(strconv.Itoa(int(e.Level))); err != nil {
		return errors.Wrap(err, "could not write to brightness register file")
	}

	err = f.Close()
	return errors.Wrap(err, "failed to close brightness register file")
}

// Init initializes this TouchDisplayManager.
func (t *TouchDisplayManager) Init(log *logrus.Entry, _ json.RawMessage) error {
	// TODO might be a good idea to set a default brightness in here.
	t.log = log
	return nil
}

// Cleanup tears down this TouchDisplayManager.
func (t *TouchDisplayManager) Cleanup() error { return nil }
