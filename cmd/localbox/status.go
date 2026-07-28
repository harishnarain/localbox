package main

import (
	"fmt"

	"github.com/harishnarain/localbox/internal/drivers"
	"github.com/harishnarain/localbox/internal/orchestrator"
)

// renderStatus formats a driver's identity for the `status` command's
// stdout output.
func renderStatus(d drivers.Driver) string {
	return fmt.Sprintf("driver: %s (boot budget: %s)", d.Name(), d.BootBudget())
}

// runStatus selects the current platform's driver via
// orchestrator.NewForCurrentPlatform and returns the formatted status
// line, or an error if no driver is available for this platform.
func runStatus() (string, error) {
	o, err := orchestrator.NewForCurrentPlatform()
	if err != nil {
		return "", err
	}
	return renderStatus(o.Driver()), nil
}
