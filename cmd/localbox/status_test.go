package main

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/harishnarain/localbox/internal/drivers"
)

// fakeDriver is a minimal test double for drivers.Driver, local to this
// package (mirrors internal/orchestrator's unexported fakeDriver, which
// isn't reachable from here since it lives in a _test.go file in a
// different package).
type fakeDriver struct {
	name       string
	bootBudget time.Duration
}

var _ drivers.Driver = (*fakeDriver)(nil)

func (d *fakeDriver) Name() string              { return d.name }
func (d *fakeDriver) BootBudget() time.Duration { return d.bootBudget }
func (d *fakeDriver) Boot(ctx context.Context, spec drivers.SandboxSpec) (drivers.Sandbox, error) {
	return nil, drivers.ErrNotImplemented
}

func TestRenderStatus(t *testing.T) {
	got := renderStatus(&fakeDriver{name: "linux", bootBudget: 100 * time.Millisecond})
	want := "driver: linux (boot budget: 100ms)"
	if got != want {
		t.Fatalf("renderStatus() = %q, want %q", got, want)
	}
}

func TestRunStatus(t *testing.T) {
	line, err := runStatus()
	if err != nil {
		t.Fatalf("runStatus() returned error: %v", err)
	}
	if !strings.HasPrefix(line, "driver: ") {
		t.Fatalf("runStatus() = %q, want prefix %q", line, "driver: ")
	}
	if !strings.Contains(line, "boot budget:") {
		t.Fatalf("runStatus() = %q, want it to contain %q", line, "boot budget:")
	}
}
