package windows

import (
	"testing"

	"github.com/harishnarain/localbox/internal/drivers/drivertest"
)

func TestDriverSkeleton(t *testing.T) {
	drivertest.AssertDriverSkeleton(t, &Driver{}, "windows", bootBudget)
}

func TestSandboxSkeleton(t *testing.T) {
	drivertest.AssertSandboxSkeleton(t, &Sandbox{})
}
