// Package windows implements the LocalBox sandbox driver for Windows,
// backed by a WSL2 utility VM / Hyper-V isolated runtime. Target boot
// budget: < 1.5s (CLAUDE.md Principle 5).
//
// Driver implements the shared internal/drivers.Driver interface; boot,
// exec, diff, and teardown are not yet implemented (every method returns
// drivers.ErrNotImplemented) — no WSL2/Hyper-V code has been written or
// verified yet.
package windows
