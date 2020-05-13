package ttsalerter

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"

	pipanel "github.com/BenJetson/pipanel/go"
	htgotts "github.com/hegedustibor/htgo-tts"
)

const (
	tempDirDefault  string = "/tmp/pipanel-tts/"
	languageDefault string = "en"
)

// Config specifies options that modify TTSAlerter behavior.
type Config struct {
	// NoTTSPrefix may be prepended to the text of a pipanel.AlertEvent.Message
	// to prevent TTSAlerter from reading the message out loud.
	//
	// If not set, TTSAlerter will always read the message out loud.
	// Defaults to not set, so messages are always read out loud by default.
	NoTTSPrefix string `json:"no_tts_prefix"`
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
	log         *log.Logger
	speech      *htgotts.Speech
	cfg         Config
	checkPrefix bool
}

// New creates a TTSAlerter instance.
func New() *TTSAlerter { return &TTSAlerter{} }

// ShowAlert will handle pipanel alert events by reading the alert message
// out loud to the user.
func (t *TTSAlerter) ShowAlert(e pipanel.AlertEvent) error {
	// If prefix checking is enabled and the message has the NoTTS prefix,
	// skip reading this message out loud.
	if t.checkPrefix && strings.HasPrefix(e.Message, t.cfg.NoTTSPrefix) {
		return nil
	}

	if err := t.speech.Speak(e.Message); err != nil {
		t.log.Println("Problem detected when reading alert message: ", err)
		return err
	}
	t.log.Printf("Alert message read out loud to user.")

	return nil
}

func (t *TTSAlerter) Init(log *log.Logger, rawCfg json.RawMessage) error {
	t.log = log

	// Decode config structure.
	d := json.NewDecoder(bytes.NewReader(rawCfg))
	d.DisallowUnknownFields()

	if err := d.Decode(&(t.cfg)); err != nil {
		return err
	}

	// Replace zero values with defaults.
	t.cfg.fillDefaults()

	// Enable the checkPrefix feature if a prefix is specified.
	t.checkPrefix = len(t.cfg.NoTTSPrefix) > 0

	// Create a HTGo-TTS instance to facilitate communication with Google TTS.
	t.speech = &htgotts.Speech{
		Folder:   t.cfg.TempDir,
		Language: t.cfg.Language,
	}

	return nil
}

func (t *TTSAlerter) Cleanup() error {
	return nil
}
