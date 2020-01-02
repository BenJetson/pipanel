package systemdpwr

import (
	"fmt"
	"log"
	"os/exec"

	pipanel "github.com/BenJetson/pipanel/go"
)

var shutdownCmd *exec.Cmd = exec.Command("sudo /sbin/shutdown", "now")
var rebootCmd *exec.Cmd = exec.Command("sudo /sbin/reboot", "now")
var displayOffCmd *exec.Cmd = exec.Command("xset dpms force off")

// SystemdPowerManager handles pipanel power events for systemd-based systems
// with X display servers.
type SystemdPowerManager struct {
	log *log.Logger
}

// New creates a SystemdPowerManager instance.
func New(log *log.Logger) *SystemdPowerManager {
	return &SystemdPowerManager{log}
}

// DoPowerAction handles pipanel power events.
func (s *SystemdPowerManager) DoPowerAction(e pipanel.PowerEvent) error {
	switch e.Action {
	case pipanel.PowerActionShutdown:
		s.log.Println("Shutting down the system NOW.")
		return shutdownCmd.Run()
	case pipanel.PowerActionReboot:
		s.log.Println("Rebooting the system NOW.")
		return rebootCmd.Run()
	case pipanel.PowerActionDisplayOff:
		s.log.Println("Turning off the display.")
		return displayOffCmd.Run()
	}

	return fmt.Errorf("command '%s' is not a known power action", e.Action)
}

func (s *SystemdPowerManager) Init() error { return nil }

func (s *SystemdPowerManager) Cleanup() error { return nil }
