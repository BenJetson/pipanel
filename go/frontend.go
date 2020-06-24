package pipanel

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Frontend provides abstraction for processing PiPanel events.
type Frontend struct {
	Alerter
	AudioPlayer
	PowerManager
	DisplayManager
}

// Init initializes all components of the Frontend.
func (f *Frontend) Init(log *logrus.Entry, cfg *FrontendConfig) error {
	if f.Alerter != nil {
		if err := f.Alerter.Init(log, cfg.AlerterConfig); err != nil {
			return errors.Wrap(err, "failed to initialize Alerter")
		}
	}

	if f.AudioPlayer != nil {
		if err := f.AudioPlayer.Init(log, cfg.AudioPlayerConfig); err != nil {
			return errors.Wrap(err, "failed to initialize AudioPlayer")
		}
	}

	if f.PowerManager != nil {
		if err := f.PowerManager.Init(log, cfg.PowerManagerConfig); err != nil {
			return errors.Wrap(err, "failed to initialize PowerManager")
		}
	}

	if f.DisplayManager != nil {
		if err := f.DisplayManager.Init(log, cfg.DisplayManagerConfig); err != nil {
			return errors.Wrap(err, "failed initialize DisplayManager")
		}
	}

	return nil
}

// Cleanup tears down all components of the Frontend.
func (f *Frontend) Cleanup() error {
	// FIXME should probably continue trying to clean up other components  even
	// if one above it returns with error

	if f.Alerter != nil {
		if err := f.Alerter.Cleanup(); err != nil {
			return errors.Wrap(err, "failed to cleanup Alerter")
		}
	}

	if f.AudioPlayer != nil {
		if err := f.AudioPlayer.Cleanup(); err != nil {
			return errors.Wrap(err, "failed to cleanup AudioPlayer")
		}
	}

	if f.PowerManager != nil {
		if err := f.PowerManager.Cleanup(); err != nil {
			return errors.Wrap(err, "failed to cleanup PowerManager")
		}
	}

	if f.DisplayManager != nil {
		if err := f.DisplayManager.Cleanup(); err != nil {
			return errors.Wrap(err, "failed to cleanup DisplayManager")
		}
	}

	return nil
}
