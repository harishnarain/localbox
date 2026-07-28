package orchestrator

import (
	"context"
	"time"

	"github.com/harishnarain/localbox/internal/drivers"
)

// fakeDriver is an in-memory test double for drivers.Driver, used so
// orchestrator tests don't depend on any real platform driver.
type fakeDriver struct {
	name       string
	bootBudget time.Duration

	bootErr   error
	bootCalls []drivers.SandboxSpec
}

var _ drivers.Driver = (*fakeDriver)(nil)

func (d *fakeDriver) Name() string { return d.name }

func (d *fakeDriver) BootBudget() time.Duration { return d.bootBudget }

func (d *fakeDriver) Boot(ctx context.Context, spec drivers.SandboxSpec) (drivers.Sandbox, error) {
	d.bootCalls = append(d.bootCalls, spec)
	if d.bootErr != nil {
		return nil, d.bootErr
	}
	return &fakeSandbox{id: "fake-sandbox"}, nil
}

// fakeSandbox is an in-memory test double for drivers.Sandbox.
type fakeSandbox struct {
	id string
}

var _ drivers.Sandbox = (*fakeSandbox)(nil)

func (s *fakeSandbox) ID() string { return s.id }

func (s *fakeSandbox) Exec(ctx context.Context, cmd []string, opts drivers.ExecOptions) (*drivers.ExecResult, error) {
	return &drivers.ExecResult{ExitCode: 0}, nil
}

func (s *fakeSandbox) Diff(ctx context.Context) ([]byte, error) {
	return nil, nil
}

func (s *fakeSandbox) Teardown(ctx context.Context) error {
	return nil
}
