package windows

import (
	"context"
	"time"

	"github.com/harishnarain/localbox/internal/drivers"
)

// bootBudget is Windows's target boot-time budget (CLAUDE.md Principle 5).
const bootBudget = 1500 * time.Millisecond

// Driver is the Windows sandbox driver, backed by a WSL2 utility VM /
// Hyper-V isolated runtime. See the package doc comment for
// implementation status.
type Driver struct{}

var _ drivers.Driver = (*Driver)(nil)

// Name returns "windows".
func (d *Driver) Name() string { return "windows" }

// BootBudget returns Windows's target boot-time budget (1.5s).
func (d *Driver) BootBudget() time.Duration { return bootBudget }

// Boot is not yet implemented.
//
// TODO(windows): create and start a WSL2 utility VM / Hyper-V isolated
// container, with spec.WorkspaceDir attached via a virtual disk overlay
// and spec.PersistentMounts attached as additional shares. Configure the
// guest's outbound network/proxy settings from spec.ProxyAddr.
func (d *Driver) Boot(ctx context.Context, spec drivers.SandboxSpec) (drivers.Sandbox, error) {
	return nil, drivers.ErrNotImplemented
}

// Sandbox is a running Windows sandbox instance. See the package doc
// comment for implementation status.
type Sandbox struct{}

var _ drivers.Sandbox = (*Sandbox)(nil)

// ID is not yet implemented.
//
// TODO(windows): return the WSL2 utility VM / Hyper-V container's
// identifier.
func (s *Sandbox) ID() string { return "" }

// Exec is not yet implemented.
//
// TODO(windows): run cmd inside the guest (e.g. via wsl.exe or the
// Hyper-V isolated runtime's exec API) and capture its exit
// code/stdout/stderr.
func (s *Sandbox) Exec(ctx context.Context, cmd []string, opts drivers.ExecOptions) (*drivers.ExecResult, error) {
	return nil, drivers.ErrNotImplemented
}

// Diff is not yet implemented.
//
// TODO(windows): compute a git-style diff of the workspace's virtual
// disk overlay against its base, for applying back to the host repo.
func (s *Sandbox) Diff(ctx context.Context) ([]byte, error) {
	return nil, drivers.ErrNotImplemented
}

// Teardown is not yet implemented.
//
// TODO(windows): stop and remove the WSL2 utility VM / Hyper-V container
// and release its virtual disk overlay and any other resources, leaving
// PersistentMounts untouched.
func (s *Sandbox) Teardown(ctx context.Context) error {
	return drivers.ErrNotImplemented
}
