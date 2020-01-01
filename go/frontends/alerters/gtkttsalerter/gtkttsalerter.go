package gtkttsalerter

import (
	"log"
	"strings"

	pipanel "github.com/BenJetson/pipanel/go"
	"github.com/BenJetson/pipanel/go/frontends/alerters/gtkalerter"
	"github.com/BenJetson/pipanel/go/frontends/alerters/ttsalerter"
)

const noTTSPrefix string = "@NOTTS@"

// GTKTTSAlerter handles PiPanel alert events by displaying them on-screen and
// reading them out loud.
type GTKTTSAlerter struct {
	*gtkalerter.GUI
	*ttsalerter.TTSAlerter
}

func New(log *log.Logger, tmpDir, language string) *GTKTTSAlerter {
	return &GTKTTSAlerter{
		GUI:        gtkalerter.New(log),
		TTSAlerter: ttsalerter.New(log, tmpDir, language),
	}
}

func (g *GTKTTSAlerter) Init() error {
	if err := g.GUI.Init(); err != nil {
		return err
	}

	if err := g.TTSAlerter.Init(); err != nil {
		return err
	}

	return nil
}

func (g *GTKTTSAlerter) Cleanup() error {
	if err := g.GUI.Cleanup(); err != nil {
		return err
	}

	if err := g.TTSAlerter.Cleanup(); err != nil {
		return err
	}

	return nil
}

// ShowAlert displays the alert on the screen using gtkalerter.GUI and
// (provided No TTS prefix is not present) reads the message out loud.
func (g *GTKTTSAlerter) ShowAlert(e pipanel.AlertEvent) error {
	// If a message has the no TTS prefix, it should not be read out loud.
	shouldReadMsg := true
	if strings.HasPrefix(e.Message, noTTSPrefix) {
		shouldReadMsg = false
		e.Message = e.Message[len(noTTSPrefix):]
	}

	err := g.GUI.ShowAlert(e)

	if err != nil {
		return err
	}

	if shouldReadMsg {
		return g.TTSAlerter.ShowAlert(e)
	}
	return nil
}
