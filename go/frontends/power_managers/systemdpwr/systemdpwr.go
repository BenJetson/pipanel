package systemdpwr

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	pipanel "github.com/BenJetson/pipanel/go"
)

// SystemdPowerManager handles pipanel power events for systemd-based systems
// with X display servers.
type SystemdPowerManager struct {
	log *log.Logger
}

// New creates a SystemdPowerManager instance.
func New() *SystemdPowerManager { return &SystemdPowerManager{} }

// DoPowerAction handles pipanel power events.
func (s *SystemdPowerManager) DoPowerAction(e pipanel.PowerEvent) error {
	switch e.Action {
	case pipanel.PowerActionShutdown:
		s.log.Println("Shutting down the system NOW.")
		return exec.Command("sudo", "shutdown", "now").Run()
	case pipanel.PowerActionReboot:
		s.log.Println("Rebooting the system NOW.")
		return exec.Command("sudo", "reboot", "now").Run()
	case pipanel.PowerActionDisplayOff:
		s.log.Println("Turning off the display.")
		return exec.Command("xset", "dpms force off").Run()
	}

	return fmt.Errorf("command '%s' is not a known power action", e.Action)
}

func (s *SystemdPowerManager) Init(log *log.Logger, _ json.RawMessage) error {
	s.log = log
	return nil
}

func (s *SystemdPowerManager) Cleanup() error { return nil }
