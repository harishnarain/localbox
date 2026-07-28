package macos

import (
	"context"
	"time"

	"github.com/harishnarain/localbox/internal/drivers"
)

// bootBudget is macOS's target boot-time budget (CLAUDE.md Principle 5).
const bootBudget = 500 * time.Millisecond

// Driver is the macOS sandbox driver, backed by Apple's
// Virtualization.framework. See the package doc comment for
// implementation status.
type Driver struct{}

var _ drivers.Driver = (*Driver)(nil)

// Name returns "macos".
func (d *Driver) Name() string { return "macos" }

// BootBudget returns macOS's target boot-time budget (500ms).
func (d *Driver) BootBudget() time.Duration { return bootBudget }

// Boot is not yet implemented.
//
// TODO(macos): create and start a VZVirtualMachine backed by an
// ARM64 Linux microVM image, with spec.WorkspaceDir attached as a CoW
// APFS snapshot and spec.PersistentMounts attached as additional
// virtiofs shares. Configure the guest's outbound network/proxy settings
// from spec.ProxyAddr.
func (d *Driver) Boot(ctx context.Context, spec drivers.SandboxSpec) (drivers.Sandbox, error) {
	return nil, drivers.ErrNotImplemented
}

// Sandbox is a running macOS sandbox instance. See the package doc
// comment for implementation status.
type Sandbox struct{}

var _ drivers.Sandbox = (*Sandbox)(nil)

// ID is not yet implemented.
//
// TODO(macos): return the VZVirtualMachine instance's identifier.
func (s *Sandbox) ID() string { return "" }

// Exec is not yet implemented.
//
// TODO(macos): run cmd inside the guest microVM (e.g. via a guest agent
// over virtio-vsock) and capture its exit code/stdout/stderr.
func (s *Sandbox) Exec(ctx context.Context, cmd []string, opts drivers.ExecOptions) (*drivers.ExecResult, error) {
	return nil, drivers.ErrNotImplemented
}

// Diff is not yet implemented.
//
// TODO(macos): compute a git-style diff of the workspace's APFS snapshot
// against its base, for applying back to the host repo.
func (s *Sandbox) Diff(ctx context.Context) ([]byte, error) {
	return nil, drivers.ErrNotImplemented
}

// Teardown is not yet implemented.
//
// TODO(macos): stop the VZVirtualMachine and release its APFS snapshot
// and any other resources, leaving PersistentMounts untouched.
func (s *Sandbox) Teardown(ctx context.Context) error {
	return drivers.ErrNotImplemented
}
