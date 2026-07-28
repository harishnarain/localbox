// Package orchestrator selects the right platform Driver (see
// internal/drivers) for the host it's running on and provides a thin
// pass-through lifecycle wrapper around it. It contains no
// platform-specific code itself: Windows is handled by detecting it
// runs under WSL2 and selecting the Linux driver, per CLAUDE.md
// Principle 4 (Windows shares the Linux driver rather than having its
// own).
//
// This is an MVP: it does not yet wire in internal/proxy or
// internal/snapshot, hold any sandbox registry/state, or support
// multiple concurrent sandboxes. Those are tracked as future work.
package orchestrator

import (
	"context"
	"fmt"
	"runtime"

	"github.com/harishnarain/localbox/internal/drivers"
	"github.com/harishnarain/localbox/internal/drivers/linux"
	"github.com/harishnarain/localbox/internal/drivers/macos"
)

// Select returns the drivers.Driver appropriate for the current host, based
// on runtime.GOOS.
func Select() (drivers.Driver, error) {
	return selectFor(runtime.GOOS)
}

// selectFor is the testable seam behind Select: runtime.GOOS can't be
// mocked directly, so tests exercise this with explicit GOOS-like values
// instead.
func selectFor(goos string) (drivers.Driver, error) {
	switch goos {
	case "darwin":
		return &macos.Driver{}, nil
	case "linux", "windows":
		// Windows runs the Linux driver inside WSL2 rather than having
		// its own native driver (CLAUDE.md Principle 4).
		return &linux.Driver{}, nil
	default:
		return nil, fmt.Errorf("orchestrator: unsupported platform %q", goos)
	}
}

// Orchestrator is a thin lifecycle wrapper around a single selected
// drivers.Driver.
type Orchestrator struct {
	driver drivers.Driver
}

// New returns an Orchestrator wrapping the given driver.
func New(d drivers.Driver) *Orchestrator {
	return &Orchestrator{driver: d}
}

// Driver returns the drivers.Driver this Orchestrator wraps, e.g. so a
// caller can report its Name/BootBudget without booting a sandbox.
func (o *Orchestrator) Driver() drivers.Driver {
	return o.driver
}

// NewForCurrentPlatform selects the appropriate driver for the current
// host (via Select) and returns an Orchestrator wrapping it.
func NewForCurrentPlatform() (*Orchestrator, error) {
	d, err := Select()
	if err != nil {
		return nil, err
	}
	return New(d), nil
}

// Boot creates and starts a new sandbox from spec, delegating to the
// wrapped driver.
func (o *Orchestrator) Boot(ctx context.Context, spec drivers.SandboxSpec) (drivers.Sandbox, error) {
	return o.driver.Boot(ctx, spec)
}
