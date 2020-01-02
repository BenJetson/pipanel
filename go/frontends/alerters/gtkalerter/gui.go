package gtkalerter

import (
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	pipanel "github.com/BenJetson/pipanel/go"
)

type GUI struct {
	log     *log.Logger
	windows []*alertWindow
}

func New(l *log.Logger) *GUI {
	return &GUI{
		log: l,
	}
}

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

func (g *GUI) Init() error {
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

func (g *GUI) Cleanup() error {
	g.log.Println("Shutting down GUI...")

	for _, w := range g.windows {
		w.Destroy()
	}

	g.windows = nil

	gtk.MainQuit()
	return nil
}
