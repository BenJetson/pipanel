package launcher

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"

	pipanel "github.com/BenJetson/pipanel/go"
	"github.com/BenJetson/pipanel/go/server"
)

const serverPortKey string = "PIPANEL_SERVER_PORT"
const serverPortDefault string = "1035"

func RunApplication(frontend *pipanel.Frontend) {
	// Create log instances.
	logServer := log.New(os.Stdout, "server ", log.LstdFlags)
	logFrontend := log.New(os.Stdout, "frontend ", log.LstdFlags)
	logMain := log.New(os.Stdout, "main ", log.LstdFlags)

	// Create signaling channels for concurrent operations.
	interrupt := make(chan os.Signal, 1)
	shutdown := make(chan struct{}, 1)

	// Notify interrupt channel when a SIGINT is detected.
	signal.Notify(interrupt, os.Interrupt)

	// Initialize frontend.
	logMain.Println("Initializing frontend...")
	if err := frontend.Init(logFrontend); err != nil {
		panic(err)
	}

	// Determine server port.
	logMain.Print("Fetching server port from environment... ")
	portStr := os.Getenv(serverPortKey)
	if len(portStr) < 1 {
		portStr = serverPortDefault
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		err = fmt.Errorf("value of %s is not a valid numeral", serverPortKey)
		panic(err)
	}

	logMain.Printf("Will use port %d.\n", port)

	// Start the server.
	logMain.Println("Starting the server...")
	server := server.New(logServer, port, frontend)

	go server.ListenAndServe(shutdown)

	// Create cleanup function for use upon interrupt/shutdown.
	cleanup := func(reason string) {
		logMain.Printf("Terminating: %s\n", reason)

		logMain.Println("Shutting down the server...")
		if err := server.Shutdown(context.Background()); err != nil {
			panic(err)
		}

		logMain.Println("Clearing frontend resources...")
		if err := frontend.Cleanup(); err != nil {
			panic(err)
		}
	}

	logMain.Println("Ready to receive events.")

	// Wait for all goroutines to shut down.
	select {
	case <-interrupt:
		cleanup("sigint detected")
	case <-shutdown:
		cleanup("server shutdown detected")
	}
}
