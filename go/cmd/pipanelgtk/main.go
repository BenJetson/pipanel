package main

import (
	"github.com/BenJetson/pipanel/go/cmd/launcher"
	"github.com/BenJetson/pipanel/go/frontends"
)

func main() { launcher.RunApplication(frontends.NewPiPanelGTK()) }
