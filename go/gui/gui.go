package gui

import (
	"log"

	"github.com/gotk3/gotk3/glib"

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

func (g *GUI) ShowAlert(a pipanel.AlertEvent) error {
	_, err := glib.IdleAdd(func() {
		w, err := newAlertWindow(a)

		if err != nil {
			panic(err)
		}

		w.ShowAll()

		g.windows = append(g.windows, w)
	})

	return err
}

func (g *GUI) Shutdown() {
	g.log.Println("Shutting down GUI...")

	for _, w := range g.windows {
		w.Destroy()
	}

	g.windows = nil
}
