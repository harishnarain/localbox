// Package macos implements the LocalBox sandbox driver for macOS,
// backed by Apple's Virtualization.framework. Target boot budget: < 500ms
// (CLAUDE.md Principle 5).
//
// Driver implements the shared internal/drivers.Driver interface; boot,
// exec, diff, and teardown are not yet implemented (every method returns
// drivers.ErrNotImplemented) — no Virtualization.framework code has been
// written or verified yet.
package macos
