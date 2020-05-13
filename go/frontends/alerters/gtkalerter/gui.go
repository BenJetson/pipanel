package gtkalerter

import (
	"encoding/json"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	pipanel "github.com/BenJetson/pipanel/go"
)

// GUI is a GTK application that is capable of responding to PiPanel alert
// events by showing them on the screen.
type GUI struct {
	log     *log.Logger
	windows []*alertWindow
}

// New creates a fresh GUI instance.
func New() *GUI { return &GUI{} }

// ShowAlert handles alert events by displaying a window to alert the user.
func (g *GUI) ShowAlert(e pipanel.AlertEvent) error {
	_, err := glib.IdleAdd(func() {
		w, err := newAlertWindow(e, g.removeInactiveWindows)

		if err != nil {
			panic(err)
		}

		g.log.Println("Displaying alert window to user.")
		w.ShowAll()

		g.windows = append(g.windows, w)
	})

	return err
}

// Init initializes this GUI instance, setting the logger and starting the GTK
// main event loop in a separate goroutine.
func (g *GUI) Init(log *log.Logger, _ json.RawMessage) error {
	g.log = log

	gtk.Init(nil)
	go gtk.Main()

	return nil
}

func (g *GUI) removeInactiveWindows() {
	g.log.Println("Clearing inactive windows...")

	count := 0

	for i := len(g.windows) - 1; i > -1; i-- {
		if g.windows[i].inactive {
			g.windows = append(g.windows[:i], g.windows[i+1:]...)

			count++
		}
	}

	g.log.Printf("Cleared all inactive windows: %d total.\n", count)
}

// Cleanup tears down this GUI instance, destroying all windows and halting the
// GTK main event loop.
func (g *GUI) Cleanup() error {
	g.log.Println("Shutting down GUI...")

	for _, w := range g.windows {
		w.Destroy()
	}

	g.windows = nil

	gtk.MainQuit()
	return nil
}
