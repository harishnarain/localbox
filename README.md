# LocalBox

**Local-first, sub-second sandboxes for AI coding agents.**

> Status: pre-alpha. No releases yet — see the [roadmap](#roadmap) below.
> This repo currently contains the project's engineering constitution and
> agent-collaboration harness; the sandboxing engine itself is not yet
> implemented.

Autonomous coding agents (Claude Code, Codex CLI, Copilot CLI, and similar)
need to run shell commands, tests, and builds without supervision. Giving
them unconstrained access to your filesystem, environment variables, and SSH
keys is a real risk — a bad command or a prompt-injected one can do lasting
damage. LocalBox spins up an isolated, ephemeral sandbox in under a second,
using your OS's own virtualization primitives, so an agent can work freely
without touching the host.

## Why not just use a cloud sandbox?

Cloud microVM sandboxes (E2B, Daytona, etc.) work, but they ship your code
off-host, add network latency to every test run, and bill you monthly for
compute you already own. LocalBox runs entirely on your machine (or one you
control on your LAN), using the hypervisor/container primitives each OS
already ships.

## How it works

LocalBox is a single compiled binary that talks directly to the native
isolation primitive on your OS:

| Host platform | Runtime | Boot budget | Isolation model |
|---|---|---|---|
| macOS (Apple Silicon) | `Virtualization.framework` | < 500 ms | Hardware-accelerated ARM64 Linux microVM |
| Linux | LXC/Incus or rootless Podman | < 100 ms | cgroups v2, user namespaces, seccomp |

Windows 11 support runs the same Linux driver inside WSL2 — WSL2 is a real
Linux VM, so a sandbox on Windows is a Linux container one layer deeper
(cgroups v2, user namespaces, seccomp), not a separate Hyper-V/utility-VM
driver, and shares the Linux driver's < 100 ms boot budget.

Key mechanisms:

- **Copy-on-write workspace snapshots** — APFS snapshots (macOS) or
  Btrfs/ZFS snapshots (Linux, including inside WSL2 on Windows) fork your
  working directory instantly. Git diffs produced inside the sandbox are
  extracted and applied back to your real repo.
- **Zero-trust credential proxy** — a host-side loopback proxy intercepts
  outbound requests from the sandbox. Real API keys (Anthropic, OpenAI,
  GitHub, …) never leave the host; the sandbox gets scoped dummy tokens
  swapped in-transit.
- **Network domain whitelisting** — agents can reach only the endpoints you
  allow (e.g. `registry.npmjs.org`, `api.anthropic.com`); everything else is
  blocked outbound.
- **Persistent config, ephemeral everything else** — directories like
  `~/.claude`, `~/.gemini`, or `~/.config/github-copilot` are attached
  persistently across sandbox teardowns so you don't re-authenticate every
  run; the rest of the sandbox disappears when the task ends.

See [CLAUDE.md](CLAUDE.md) for the full set of engineering principles this
project holds itself to.

## Quick start

Not yet available — there's no buildable CLI yet. Once `cmd/localbox` has a
real implementation, this section will cover `go install` / release binaries
and `localbox run --agent claude --repo .`.

## Pricing / editions

LocalBox is open-core:

| Tier | Price | Includes |
|---|---|---|
| Community | Free, Apache 2.0 | Single-node CLI, native hypervisor drivers, local credential proxy, basic snapshotting |
| Pro | $15/mo | Multi-node offload (e.g. Mac → home Proxmox box over Tailscale), advanced secret masking, persistent template library |
| Enterprise | $40/user/mo | Central policy engine, audit logs, SSO, shared team environment templates, compliance reporting |

The Community tier is not a crippled trial — see [CLAUDE.md](CLAUDE.md)
Principle 6.

## Roadmap

- **Phase 1 (MVP)** — macOS driver (Virtualization.framework), Linux driver
  (LXC/Podman), basic CLI, initial credential proxy.
- **Phase 2** — Verify/adapt the Linux driver running inside WSL2
  (network/proxy interop, snapshot filesystem support) for Windows 11, VS
  Code Remote-SSH / JetBrains integration, Pro remote-offload over
  WireGuard/Tailscale.
- **Phase 3** — Enterprise policy dashboard + SSO, partnerships with agent
  frameworks (Claude Code extensions, LangChain, AutoGen).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for dev setup, git workflow, and
review policy, and [CLAUDE.md](CLAUDE.md) for the principles all changes are
expected to hold to. Security issues go through [SECURITY.md](SECURITY.md),
not public issues.

## License

Apache 2.0 — see [LICENSE](LICENSE).
