// Package drivertest provides shared test helpers for asserting that a
// platform driver's skeleton implementation matches the shared
// internal/drivers.Driver/Sandbox contract. Used by each platform
// package's own _test.go (macos, linux) so the parity assertion
// logic isn't triplicated and can't silently drift between platforms.
package drivertest

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/harishnarain/localbox/internal/drivers"
)

// AssertDriverSkeleton asserts that d matches the shared contract: correct
// static Name/BootBudget values, and Boot returning ErrNotImplemented
// rather than silently succeeding or panicking.
//
// wantName and wantBudget are the values expected from d.Name() and
// d.BootBudget() for the platform under test.
func AssertDriverSkeleton(t *testing.T, d drivers.Driver, wantName string, wantBudget time.Duration) {
	t.Helper()

	if got := d.Name(); got != wantName {
		t.Errorf("Name() = %q, want %q", got, wantName)
	}
	if got := d.BootBudget(); got != wantBudget {
		t.Errorf("BootBudget() = %v, want %v", got, wantBudget)
	}

	ctx := context.Background()

	sb, err := d.Boot(ctx, drivers.SandboxSpec{WorkspaceDir: t.TempDir()})
	if !errors.Is(err, drivers.ErrNotImplemented) {
		t.Fatalf("Boot() error = %v, want errors.Is(err, ErrNotImplemented)", err)
	}
	if sb != nil {
		t.Fatalf("Boot() sandbox = %v, want nil alongside ErrNotImplemented", sb)
	}
}

// AssertSandboxSkeleton is the Sandbox-side counterpart of
// AssertDriverSkeleton, for platforms that want to exercise a Sandbox
// value directly (e.g. constructed without going through Boot).
func AssertSandboxSkeleton(t *testing.T, sb drivers.Sandbox) {
	t.Helper()

	ctx := context.Background()

	if _, err := sb.Exec(ctx, []string{"true"}, drivers.ExecOptions{}); !errors.Is(err, drivers.ErrNotImplemented) {
		t.Errorf("Exec() error = %v, want errors.Is(err, ErrNotImplemented)", err)
	}
	if _, err := sb.Diff(ctx); !errors.Is(err, drivers.ErrNotImplemented) {
		t.Errorf("Diff() error = %v, want errors.Is(err, ErrNotImplemented)", err)
	}
	if err := sb.Teardown(ctx); !errors.Is(err, drivers.ErrNotImplemented) {
		t.Errorf("Teardown() error = %v, want errors.Is(err, ErrNotImplemented)", err)
	}
}
