// Package drivers defines the shared contract every LocalBox platform
// driver (macOS, Linux, Windows — see internal/drivers/{macos,linux,windows})
// implements. It contains no platform-specific code: it exists so the
// planned orchestrator (internal/orchestrator, not yet created) can select
// and drive any platform's sandbox without knowing which one it's talking
// to, per CLAUDE.md Principle 4 (cross-platform parity).
package drivers

import (
	"context"
	"errors"
	"time"
)

// ErrNotImplemented is returned by driver/sandbox methods that don't yet
// have a real implementation for their platform. Callers should use
// errors.Is(err, ErrNotImplemented) rather than comparing errors directly,
// since a driver may wrap it with additional context.
var ErrNotImplemented = errors.New("drivers: not yet implemented")

// Driver boots sandboxes for one platform's native isolation primitive
// (Virtualization.framework, LXC/Incus, WSL2, ...).
type Driver interface {
	// Name returns the platform name this driver implements, e.g.
	// "macos", "linux", "windows".
	Name() string

	// BootBudget returns this platform's target boot-time budget, per
	// CLAUDE.md Principle 5.
	BootBudget() time.Duration

	// Boot creates and starts a new sandbox from spec. The returned
	// Sandbox is ready to accept Exec calls.
	Boot(ctx context.Context, spec SandboxSpec) (Sandbox, error)
}

// Sandbox is a single running, isolated environment created by a Driver.
type Sandbox interface {
	// ID returns an identifier for this sandbox instance, unique among
	// currently-running sandboxes for its driver.
	ID() string

	// Exec runs cmd inside the sandbox and returns its result.
	Exec(ctx context.Context, cmd []string, opts ExecOptions) (*ExecResult, error)

	// Diff returns a git-style diff of changes made inside the sandbox's
	// workspace, suitable for applying back to the real, host-side repo.
	Diff(ctx context.Context) ([]byte, error)

	// Teardown destroys the sandbox and releases any resources it holds.
	// Persistent mounts (see SandboxSpec.PersistentMounts) are not
	// affected by teardown.
	Teardown(ctx context.Context) error
}

// SandboxSpec configures a sandbox at boot time.
type SandboxSpec struct {
	// WorkspaceDir is the host directory to copy-on-write-snapshot into
	// the sandbox as its working directory.
	WorkspaceDir string

	// PersistentMounts lists host directories (e.g. ~/.claude,
	// ~/.config/github-copilot) that are attached persistently and
	// survive sandbox teardown, so credentials/config don't need
	// re-authentication on every run.
	PersistentMounts []string

	// ProxyAddr is the loopback address of the host-side credential
	// proxy (see internal/proxy) the sandbox should route outbound
	// requests through, so real secrets never enter the sandbox.
	ProxyAddr string
}

// ExecOptions configures a single Sandbox.Exec call.
type ExecOptions struct {
	// Dir is the working directory inside the sandbox, relative to the
	// workspace root. Empty means the workspace root itself.
	Dir string

	// Env holds additional "KEY=VALUE" environment variables, appended
	// to the sandbox's default environment.
	Env []string
}

// ExecResult is the outcome of a Sandbox.Exec call.
type ExecResult struct {
	ExitCode int
	Stdout   []byte
	Stderr   []byte
}
