package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"sort"

	"github.com/sirupsen/logrus"

	pipanel "github.com/BenJetson/pipanel/go"
	"github.com/BenJetson/pipanel/go/frontends"
	"github.com/BenJetson/pipanel/go/server"
)

// Command line flag constants.
const (
	serverPortFlag    = "port"
	serverPortDefault = -1
	serverPortDesc    = "the port for the server to listen on; " +
		"if this flag is set, will override port from config file"

	cfgPathFlag    = "config"
	cfgPathDefault = "~/pipanel_config.json"
	cfgPathDesc    = "absolute path to the configuration JSON file"
)

// frontendRegister is map from frontend name to a function that creates a new
// instance of that particular frontend type.
var frontendRegister = map[string]func() *pipanel.Frontend{
	"console":     frontends.NewConsoleFrontend,
	"pipanel-gtk": frontends.NewPiPanelGTK,
}

func checkConfig(log *logrus.Entry, cfg *pipanel.Config) {
	if cfg.Server.Port < 0 {
		log.Fatalln("Port number cannot be negative.")
	} else if cfg.Server.Port < 1024 {
		log.Fatalln("Port numbers 0-1023 are reserved by the system.")
	}

	if _, ok := frontendRegister[cfg.Frontend.Name]; !ok {
		log.Fatalf("No such frontend '%s' registered.\n", cfg.Frontend.Name)
	}

	log.Println("Configuration accepted.")
}

func loadConfig(log *logrus.Entry) *pipanel.Config {
	// Set up command line flags.
	var port int
	flag.IntVar(&port, serverPortFlag, serverPortDefault, serverPortDesc)

	var cfgPath string
	flag.StringVar(&cfgPath, cfgPathFlag, cfgPathDefault, cfgPathDesc)

	// Read command line flags.
	flag.Parse()

	// Load config from disk.
	log.Println("Loading configuration from disk.")
	file, err := os.Open(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load configuration file at path '%s'.\n", cfgPath)
	}

	// Decode JSON into configuration structure.
	var cfg pipanel.Config

	d := json.NewDecoder(file)
	d.DisallowUnknownFields()

	if err = d.Decode(&cfg); err != nil {
		log.Fatalln("Failed to read configuration file: bad JSON formatting.")
	}

	// If a port is specified at the shell prompt, overwrite the config.
	if port != -1 {
		log.Println("Port flag set: overriding configuration file preference.")
		cfg.Server.Port = port
	}

	// Validate the configuration.
	checkConfig(log, &cfg)

	return &cfg
}

const msgSrcLogKey = "msg-src"

func main() {
	// Create log instances.
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		SortingFunc: func(keys []string) {
			sort.Slice(keys, func(i, j int) bool {
				// Ensure that the component key is always at the top.
				if keys[i] == msgSrcLogKey {
					return true
				} else if keys[j] == msgSrcLogKey {
					return false
				}

				// All keys after component are sorted alphabetically.
				return keys[i] < keys[j]
			})
		},
	})

	logServer := logger.WithFields(logrus.Fields{msgSrcLogKey: "server"})
	logFrontend := logger.WithFields(logrus.Fields{msgSrcLogKey: "frontend"})
	logMain := logger.WithFields(logrus.Fields{msgSrcLogKey: "main"})

	// Load configuration.
	cfg := loadConfig(logMain)

	// Create signaling channels for concurrent operations.
	interrupt := make(chan os.Signal, 1)
	shutdown := make(chan struct{}, 1)

	// Notify interrupt channel when a SIGINT is detected.
	signal.Notify(interrupt, os.Interrupt)

	// Create new frontend instance and initialize it.
	logMain.Println("Initializing frontend...")
	frontend := frontendRegister[cfg.Frontend.Name]()

	err := frontend.Init(logFrontend, &cfg.Frontend)
	if err != nil {
		logMain.WithFields(logrus.Fields{
			"error": err,
		}).Fatalln("Problem when initializing frontend.")
	}

	// Start the server.
	logMain.Println("Starting the server...")
	server := server.New(logServer, cfg.Server.Port, frontend)

	go server.ListenAndServe(shutdown)

	// Create cleanup function for use upon interrupt/shutdown.
	cleanup := func(reason string) {
		logMain.Printf("Terminating: %s\n", reason)

		logMain.Println("Shutting down the server...")
		if err = server.Shutdown(context.Background()); err != nil {
			logMain.WithFields(logrus.Fields{
				"error": err,
			}).Errorln("Shutting down server failed.")
		}

		logMain.Println("Clearing frontend resources...")
		if err = frontend.Cleanup(); err != nil {
			logMain.WithFields(logrus.Fields{
				"error": err,
			}).Errorln("Clearing frontend resources failed.")
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
