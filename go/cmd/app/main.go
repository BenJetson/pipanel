package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	pipanel "github.com/BenJetson/pipanel/go"
	"github.com/BenJetson/pipanel/go/frontends"
	"github.com/BenJetson/pipanel/go/server"
)

func main() {
	// Create log instances.
	logServer := log.New(os.Stdout, "server", log.LstdFlags)
	logFrontend := log.New(os.Stdout, "frontend", log.LstdFlags)

	// Create signaling channels for concurrent operations.
	interrupt := make(chan os.Signal, 1)
	shutdown := make(chan struct{}, 1)

	// Notify interrupt channel when a SIGINT is detected.
	signal.Notify(interrupt, os.Interrupt)

	// Create the frontend instance.
	var frontend *pipanel.Frontend
	// frontend = frontends.NewConsoleFrontend(logFrontend)
	frontend = frontends.NewPiPanelGTK(logFrontend)

	// Start the server.
	server := server.New(logServer, 1035, frontend)

	go server.ListenAndServe(shutdown)

	// Wait for all goroutines to shut down.
	select {
	case <-interrupt:
		fmt.Println("sigint detected")
		server.Shutdown(context.TODO()) // nolint: errcheck
	case <-shutdown:
		fmt.Println("server shutdown detected")
		server.Shutdown(context.TODO()) // nolint: errcheck
	}
}
