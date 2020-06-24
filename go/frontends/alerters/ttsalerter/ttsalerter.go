package ttsalerter

import (
	"bytes"
	"encoding/json"

	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	pipanel "github.com/BenJetson/pipanel/go"
)

var _ pipanel.Alerter = (*TTSAlerter)(nil)

const (
	tempDirDefault  string = "/tmp/pipanel-tts/"
	languageDefault string = "en"
)

// Config specifies options that modify TTSAlerter behavior.
type Config struct {
	// TempDir will be used by TTSAlerter as a cache for TTS clips.
	//
	// Defaults to "/tmp/pipanel-tts" if not net.
	TempDir string `json:"temp_dir"`
	// Language determines what language will be used by TTSAlerter when
	// converting the message to speech.
	//
	// Defaults to "en" if not set.
	Language string `json:"language"`
}

// fillDefaults will overwrite zero values with the default configuration.
func (cfg *Config) fillDefaults() {
	if len(cfg.TempDir) < 1 {
		cfg.TempDir = tempDirDefault
	}

	if len(cfg.Language) < 1 {
		cfg.Language = languageDefault
	}
}

// TTSAlerter is an implementation of pipanel.Alerter that reads alerts
// out loud via text-to-speech.
type TTSAlerter struct {
	log    *logrus.Entry
	speech *htgotts.Speech
	cfg    Config
}

// New creates a TTSAlerter instance.
func New() *TTSAlerter { return &TTSAlerter{} }

// ShowAlert will handle pipanel alert events by reading the alert message
// out loud to the user.
func (t *TTSAlerter) ShowAlert(e pipanel.AlertEvent) error {
	t.log.Println("Starting to read alert message out loud to user.")

	// Since the Speak method blocks while reading to the user, it will be run
	// asynchronously. Consequentially, all ShowAlert invocations upon a
	// TTSAlerter will always return with success. Errors are logged only.
	go func() {
		if err := t.speech.Speak(e.Message); err != nil {
			err = errors.Wrap(err, "failed to read alert message out loud")
			t.log.WithFields(logrus.Fields{
				"error": err,
			}).Errorln("Problem when reading alert message out loud.")
		}
		t.log.Infoln("Reading alert message out loud has finished.")
	}()

	return nil
}

// Init initializes this TTSAlerter, loading the configuration from the
// provided JSON.
func (t *TTSAlerter) Init(log *logrus.Entry, rawCfg json.RawMessage) error {
	t.log = log

	// Decode config structure.
	d := json.NewDecoder(bytes.NewReader(rawCfg))
	d.DisallowUnknownFields()

	if err := d.Decode(&t.cfg); err != nil {
		return errors.Wrap(err, "malformed JSON for TTSAlerter configuration")
	}

	// Replace zero values with defaults.
	t.cfg.fillDefaults()

	// Create a HTGo-TTS instance to facilitate communication with Google TTS.
	t.speech = &htgotts.Speech{
		Folder:   t.cfg.TempDir,
		Language: t.cfg.Language,
	}

	return nil
}

// Cleanup tears down this TTSAlerter.
func (t *TTSAlerter) Cleanup() error { return nil }
