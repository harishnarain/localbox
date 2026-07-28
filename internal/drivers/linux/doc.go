// Package linux implements the LocalBox sandbox driver for Linux, backed
// by LXC/Incus or rootless Podman (cgroups v2, user namespaces, seccomp).
// Target boot budget: < 100ms (CLAUDE.md Principle 5).
//
// Driver implements the shared internal/drivers.Driver interface; boot,
// exec, diff, and teardown are not yet implemented (every method returns
// drivers.ErrNotImplemented) — no LXC/Incus/Podman code has been written
// or verified yet.
package linux
