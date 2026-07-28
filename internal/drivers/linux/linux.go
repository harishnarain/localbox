package linux

import (
	"context"
	"time"

	"github.com/harishnarain/localbox/internal/drivers"
)

// bootBudget is Linux's target boot-time budget (CLAUDE.md Principle 5).
const bootBudget = 100 * time.Millisecond

// Driver is the Linux sandbox driver, backed by LXC/Incus or rootless
// Podman. See the package doc comment for implementation status.
type Driver struct{}

var _ drivers.Driver = (*Driver)(nil)

// Name returns "linux".
func (d *Driver) Name() string { return "linux" }

// BootBudget returns Linux's target boot-time budget (100ms).
func (d *Driver) BootBudget() time.Duration { return bootBudget }

// Boot is not yet implemented.
//
// TODO(linux): create and start an LXC/Incus container (or rootless
// Podman container) with cgroups v2, user namespaces, and seccomp
// confinement, with spec.WorkspaceDir attached as a CoW Btrfs/ZFS
// snapshot and spec.PersistentMounts bind-mounted in. Configure the
// container's outbound network/proxy settings from spec.ProxyAddr.
func (d *Driver) Boot(ctx context.Context, spec drivers.SandboxSpec) (drivers.Sandbox, error) {
	return nil, drivers.ErrNotImplemented
}

// Sandbox is a running Linux sandbox instance. See the package doc
// comment for implementation status.
type Sandbox struct{}

var _ drivers.Sandbox = (*Sandbox)(nil)

// ID is not yet implemented.
//
// TODO(linux): return the container's identifier.
func (s *Sandbox) ID() string { return "" }

// Exec is not yet implemented.
//
// TODO(linux): run cmd inside the container (e.g. via lxc exec / podman
// exec) and capture its exit code/stdout/stderr.
func (s *Sandbox) Exec(ctx context.Context, cmd []string, opts drivers.ExecOptions) (*drivers.ExecResult, error) {
	return nil, drivers.ErrNotImplemented
}

// Diff is not yet implemented.
//
// TODO(linux): compute a git-style diff of the workspace's Btrfs/ZFS
// snapshot against its base, for applying back to the host repo.
func (s *Sandbox) Diff(ctx context.Context) ([]byte, error) {
	return nil, drivers.ErrNotImplemented
}

// Teardown is not yet implemented.
//
// TODO(linux): stop and remove the container and release its snapshot
// and any other resources, leaving PersistentMounts untouched.
func (s *Sandbox) Teardown(ctx context.Context) error {
	return drivers.ErrNotImplemented
}
