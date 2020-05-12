package gtkttsalerter

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	pipanel "github.com/BenJetson/pipanel/go"
	"github.com/BenJetson/pipanel/go/frontends/alerters/gtkalerter"
	"github.com/BenJetson/pipanel/go/frontends/alerters/ttsalerter"
)

const noTTSPrefixKey string = "PIPANEL_NOTTS_PREFIX"
const noTTSPrefixDefault string = "@NOTTS@"

// GTKTTSAlerter handles PiPanel alert events by displaying them on-screen and
// reading them out loud.
type GTKTTSAlerter struct {
	*gtkalerter.GUI
	*ttsalerter.TTSAlerter
	log         *log.Logger
	noTTSPrefix string
}

func New() *GTKTTSAlerter {
	return &GTKTTSAlerter{
		GUI:        gtkalerter.New(),
		TTSAlerter: ttsalerter.New(),
	}
}

func (g *GTKTTSAlerter) Init(log *log.Logger, config json.RawMessage) error {
	g.log = log

	// Fetch NoTTS prefix from environment.
	g.noTTSPrefix = os.Getenv(noTTSPrefixKey)
	if len(g.noTTSPrefix) < 1 {
		g.noTTSPrefix = noTTSPrefixDefault
	}

	if err := g.GUI.Init(log, config); err != nil {
		return err
	}

	if err := g.TTSAlerter.Init(log, config); err != nil {
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
	if strings.HasPrefix(e.Message, g.noTTSPrefix) {
		shouldReadMsg = false
		e.Message = e.Message[len(g.noTTSPrefix):]
		g.log.Println("Detected No TTS flag; skipping alert read-out.")
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
