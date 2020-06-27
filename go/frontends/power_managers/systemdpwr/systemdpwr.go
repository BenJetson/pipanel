package systemdpwr

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	pipanel "github.com/BenJetson/pipanel/go"

	"github.com/sirupsen/logrus"
)

var _ pipanel.PowerManager = (*SystemdPowerManager)(nil)

// SystemdPowerManager handles pipanel power events for systemd-based systems
// with X display servers.
type SystemdPowerManager struct {
	log *logrus.Entry
}

// New creates a SystemdPowerManager instance.
func New() *SystemdPowerManager { return &SystemdPowerManager{} }

// DoPowerAction handles pipanel power events.
func (s *SystemdPowerManager) DoPowerAction(ctx context.Context,
	e pipanel.PowerEvent) error {

	switch e.Action {
	case pipanel.PowerActionShutdown:
		s.log.WithContext(ctx).Println("Shutting down the system NOW.")
		return exec.Command("sudo", "shutdown", "now").Run()
	case pipanel.PowerActionReboot:
		s.log.WithContext(ctx).Println("Rebooting the system NOW.")
		return exec.Command("sudo", "reboot", "now").Run()
	case pipanel.PowerActionDisplayOff:
		s.log.WithContext(ctx).Println("Turning off the display.")
		return exec.Command("xset", "dpms force off").Run()
	}

	return fmt.Errorf("command '%s' is not a known power action", e.Action)
}

// Init initializes this SystemdPowerManager by setting the logger.
func (s *SystemdPowerManager) Init(log *logrus.Entry, _ json.RawMessage) error {
	s.log = log
	return nil
}

// Cleanup tears down this SystemdPowerManager.
func (s *SystemdPowerManager) Cleanup() error { return nil }
