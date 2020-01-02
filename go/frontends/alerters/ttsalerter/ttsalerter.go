package ttsalerter

import (
	"log"
	"os"

	pipanel "github.com/BenJetson/pipanel/go"
	htgotts "github.com/hegedustibor/htgo-tts"
)

const (
	tmpDirKey     string = "PIPANEL_TTS_TMP_DIR"
	tmpDirDefault string = "/tmp/pipanel-tts/"

	languageKey     string = "PIPANEL_TTS_LANGUAGE"
	languageDefault string = "en"
)

// TTSAlerter is an implementation of pipanel.Alerter that reads alerts
// out loud via text-to-speech.
type TTSAlerter struct {
	log    *log.Logger
	speech *htgotts.Speech
}

// New creates a TTSAlerter instance given the logger to use.
func New(log *log.Logger) *TTSAlerter {
	return &TTSAlerter{log: log}
}

// ShowAlert will handle pipanel alert events by reading the alert message
// out loud to the user.
func (t *TTSAlerter) ShowAlert(e pipanel.AlertEvent) error {
	t.speech.Speak(e.Message)
	t.log.Printf("Alert message read out loud to user.")

	return nil
}

func (t *TTSAlerter) Init() error {
	tmpDir := os.Getenv(tmpDirKey)
	if len(tmpDir) < 1 {
		tmpDir = tmpDirDefault
	}

	language := os.Getenv(languageKey)
	if len(language) < 1 {
		language = languageDefault
	}

	t.speech = &htgotts.Speech{
		Folder:   tmpDir,
		Language: language,
	}

	return nil
}

func (t *TTSAlerter) Cleanup() error {
	return nil
}
