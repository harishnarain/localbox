package orchestrator

import (
	"context"
	"errors"
	"testing"

	"github.com/harishnarain/localbox/internal/drivers"
	"github.com/harishnarain/localbox/internal/drivers/linux"
	"github.com/harishnarain/localbox/internal/drivers/macos"
)

func TestSelectFor(t *testing.T) {
	tests := []struct {
		goos     string
		wantType string
	}{
		{"darwin", "*macos.Driver"},
		{"linux", "*linux.Driver"},
		{"windows", "*linux.Driver"}, // WSL2: Windows shares the Linux driver.
	}

	for _, tt := range tests {
		t.Run(tt.goos, func(t *testing.T) {
			d, err := selectFor(tt.goos)
			if err != nil {
				t.Fatalf("selectFor(%q) returned error: %v", tt.goos, err)
			}
			switch tt.goos {
			case "darwin":
				if _, ok := d.(*macos.Driver); !ok {
					t.Fatalf("selectFor(%q) = %T, want *macos.Driver", tt.goos, d)
				}
			case "linux", "windows":
				if _, ok := d.(*linux.Driver); !ok {
					t.Fatalf("selectFor(%q) = %T, want *linux.Driver", tt.goos, d)
				}
			}
		})
	}
}

func TestSelectForUnsupported(t *testing.T) {
	d, err := selectFor("plan9")
	if err == nil {
		t.Fatalf("selectFor(\"plan9\") = %v, nil; want an error", d)
	}
	if d != nil {
		t.Fatalf("selectFor(\"plan9\") driver = %v, want nil", d)
	}
}

func TestNewForCurrentPlatform(t *testing.T) {
	o, err := NewForCurrentPlatform()
	if err != nil {
		t.Fatalf("NewForCurrentPlatform() returned error: %v", err)
	}
	if o == nil {
		t.Fatal("NewForCurrentPlatform() returned nil Orchestrator")
	}
	if o.driver == nil {
		t.Fatal("NewForCurrentPlatform() orchestrator has nil driver")
	}
}

func TestOrchestratorDriver(t *testing.T) {
	fd := &fakeDriver{name: "fake"}
	o := New(fd)

	if o.Driver() != fd {
		t.Fatalf("Driver() = %v, want %v", o.Driver(), fd)
	}
}

func TestOrchestratorBootPassThrough(t *testing.T) {
	fd := &fakeDriver{name: "fake"}
	o := New(fd)

	spec := drivers.SandboxSpec{WorkspaceDir: "/tmp/workspace"}
	sb, err := o.Boot(context.Background(), spec)
	if err != nil {
		t.Fatalf("Boot() returned error: %v", err)
	}
	if sb == nil {
		t.Fatal("Boot() returned nil Sandbox")
	}
	if sb.ID() != "fake-sandbox" {
		t.Fatalf("Boot() sandbox ID = %q, want %q", sb.ID(), "fake-sandbox")
	}

	if len(fd.bootCalls) != 1 {
		t.Fatalf("fakeDriver.Boot called %d times, want 1", len(fd.bootCalls))
	}
	if fd.bootCalls[0].WorkspaceDir != spec.WorkspaceDir {
		t.Fatalf("Boot() passed spec %+v, want %+v", fd.bootCalls[0], spec)
	}
}

func TestOrchestratorBootError(t *testing.T) {
	wantErr := errors.New("boom")
	fd := &fakeDriver{name: "fake", bootErr: wantErr}
	o := New(fd)

	sb, err := o.Boot(context.Background(), drivers.SandboxSpec{})
	if !errors.Is(err, wantErr) {
		t.Fatalf("Boot() error = %v, want %v", err, wantErr)
	}
	if sb != nil {
		t.Fatalf("Boot() sandbox = %v, want nil on error", sb)
	}
}
