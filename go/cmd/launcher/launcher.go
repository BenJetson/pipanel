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

	// Create signaling channels for concurrent operations.
	interrupt := make(chan os.Signal, 1)
	shutdown := make(chan struct{}, 1)

	// Notify interrupt channel when a SIGINT is detected.
	signal.Notify(interrupt, os.Interrupt)

	// Initialize frontend.
	if err := frontend.Init(logFrontend); err != nil {
		panic(err)
	}

	// Determine server port.
	portStr := os.Getenv(serverPortKey)
	if len(portStr) < 1 {
		portStr = serverPortDefault
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		err = fmt.Errorf("value of %s is not a valid numeral", serverPortKey)
		panic(err)
	}

	// Start the server.
	server := server.New(logServer, port, frontend)

	go server.ListenAndServe(shutdown)

	// Create cleanup function for use upon interrupt/shutdown.
	cleanup := func() {
		if err := server.Shutdown(context.Background()); err != nil {
			panic(err)
		}

		if err := frontend.Cleanup(); err != nil {
			panic(err)
		}
	}

	// Wait for all goroutines to shut down.
	select {
	case <-interrupt:
		fmt.Println("sigint detected")
		cleanup()
	case <-shutdown:
		fmt.Println("server shutdown detected")
		cleanup()
	}
}
