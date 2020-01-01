package ttsalerter

import (
	"log"

	pipanel "github.com/BenJetson/pipanel/go"
	htgotts "github.com/hegedustibor/htgo-tts"
)

// TTSAlerter is an implementation of pipanel.Alerter that reads alerts
// out loud via text-to-speech.
type TTSAlerter struct {
	log    *log.Logger
	speech *htgotts.Speech
}

// New creates a TTSAlerter instance given the logger to use, path to a suitable
// temporary directory for audio clips, and the TTS language.
func New(log *log.Logger, tmpDir, language string) *TTSAlerter {
	return &TTSAlerter{
		log: log,
		speech: &htgotts.Speech{
			Folder:   tmpDir,
			Language: language,
		},
	}
}

// ShowAlert will handle pipanel alert events by reading the alert message
// out loud to the user.
func (t *TTSAlerter) ShowAlert(e pipanel.AlertEvent) error {
	t.speech.Speak(e.Message)
	t.log.Printf("Alert message read out loud to user.")

	return nil
}

func (t *TTSAlerter) Init() error {
	return nil
}

func (t *TTSAlerter) Cleanup() error {
	return nil
}
