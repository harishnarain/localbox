// Package linux implements the LocalBox sandbox driver for Linux, backed
// by LXC/Incus or rootless Podman (cgroups v2, user namespaces, seccomp).
// Target boot budget: < 100ms (CLAUDE.md Principle 5). Not yet implemented.
package linux
