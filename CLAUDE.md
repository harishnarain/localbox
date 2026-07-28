# CLAUDE.md

Operating instructions for Claude Code (and any other agent framework) working
in this repository. This file is the source of truth: the principles in Part 1
are non-negotiable — every design decision, PR, and agent action is expected
to be checked against them. When this document and convenience disagree, this
document wins.

Amending Part 1 requires a PR that a human maintainer explicitly approves —
agents should treat it as read-only context and propose changes via PR rather
than editing it as a side effect of unrelated work.

## What this repo is

LocalBox is a local-first, cross-platform orchestrator that gives AI coding
agents (Claude Code, Codex CLI, Copilot CLI, etc.) fast, ephemeral, isolated
sandboxes instead of raw host access — see [README.md](README.md) for the
full pitch.

---

## Part 1 — Principles

### 1. The sandbox boundary is the product

LocalBox exists to give AI agents a blast-radius limit. Any regression in
isolation — filesystem escape, privilege escalation, network exfiltration,
credential leakage — is a Sev-1, not a routine bug. Any change touching
`internal/drivers/*`, `internal/proxy/*`, or `internal/snapshot/*` requires a
security review before merge (the `security-reviewer` agent, or a human, per
[CONTRIBUTING.md](CONTRIBUTING.md)).

### 2. Local-first, zero required cloud dependency

Community-tier functionality must work fully offline, on-host, with no
telemetry and no calls to LocalBox-operated infrastructure. Anything that
phones home is opt-in, off by default, and disclosed in the README.

### 3. Zero-trust credentials by default

Real secrets (Anthropic/OpenAI/GitHub API keys, SSH keys) never enter a
sandbox. The credential proxy substitutes scoped dummy tokens at the host
boundary. A driver or feature that hands a sandbox a real secret to save a
round trip is not acceptable, no matter how convenient.

### 4. Cross-platform parity is a requirement, not a nice-to-have

LocalBox ships drivers for macOS (Virtualization.framework) and Linux
(LXC/Incus) behind one orchestration interface. Windows support runs the
Linux driver inside WSL2 rather than a third native driver — WSL2 is a
real Linux kernel, so a sandbox "on Windows" is a Linux container one
layer deeper, not a distinct isolation primitive to maintain at parity. A
feature isn't "done" when it works on one platform — it's done when it
works on both drivers (and, transitively, on Windows via WSL2), or parity
is explicitly deferred with a tracked issue and a documented reason. See
the `platform-driver` agent and `driver-parity-check` skill.

### 5. Performance budgets are product requirements

| Platform               | Boot budget |
|-------------------------|------------|
| macOS (Apple Silicon)   | < 500 ms   |
| Linux                   | < 100 ms   |

Windows (via WSL2) uses the Linux driver directly and shares its < 100 ms
budget — WSL2 runs a real Linux kernel, so there is no separate native
primitive to budget for.

A change that regresses these numbers needs an explicit justification in the
PR description, not a silent trade-off. Use `/bench` before merging anything
in the boot/exec hot path.

### 6. Open-core boundary stays honest

The Community tier (Apache 2.0) must remain genuinely useful on its own —
never crippled to force an upgrade. Pro/Enterprise-only code (multi-node
offload, SSO, policy engine, audit reporting) lives in clearly separated
packages/build tags, never entangled with core sandboxing logic.

### 7. Single static binary, minimal runtime surface

LocalBox distributes as one compiled binary per platform. Avoid introducing
runtime dependencies (interpreters, dynamic libs, background daemons) the
user didn't ask to install.

### 8. Say what you actually verified

Boot-time numbers, "parity achieved," and "isolation holds" are claims about
reality, not aspirations. Don't report a performance budget met or a sandbox
escape closed without having actually run the check that would catch it.

---

## Part 2 — Operating instructions

### Repo layout

Current:

```
CLAUDE.md            — this file
README.md            — public project overview
CONTRIBUTING.md       — dev setup, git workflow, review policy
SECURITY.md           — vulnerability disclosure
CODE_OF_CONDUCT.md
docs/spec-driven-development.md — the specify → plan → tasks → implement workflow
.claude/agents/       — specialized subagents (security-reviewer, platform-driver, release-manager, spec-manager)
.claude/commands/     — slash commands (/new-driver, /bench, /release, /specify, /plan, /tasks, /implement, /spec-status)
.claude/skills/       — packaged workflows (threat-model-review, driver-parity-check)
.claude/settings.json — shared permission allowlist
.github/ISSUE_TEMPLATE/spec.yml — issue form backing the spec workflow
cmd/localbox/         — CLI entrypoint (stub)
internal/drivers/{macos,linux}/ — platform drivers (macOS, Linux; Windows
                          uses the Linux driver directly, via WSL2)
internal/proxy/       — credential proxy (stub)
internal/snapshot/    — CoW workspace snapshotting (stub)
Makefile
```

Planned as the product grows (not yet created): `internal/orchestrator/` for
the driver-selection/lifecycle engine, `internal/config/` for
`~/.localbox` config handling, `pkg/` only if a public Go API is ever needed.
Don't pre-create packages "for later" — add them when a change actually needs
them.

### Build & test

```
make build   # go build ./...
make test    # go test ./...
make vet     # go vet ./...
make fmt     # gofmt -l -w .
make check   # fmt + vet + test, what CI runs
```

### Git workflow

See [CONTRIBUTING.md](CONTRIBUTING.md) for the full policy. Summary: trunk-based
on `main` (protected), branches named `type/short-desc` (e.g.
`fix/proxy-token-leak`), Conventional Commits, PRs require CI green + one
approval, security-sensitive paths require the security review in Principle 1.

### Using the agent harness

- **`security-reviewer` agent** — invoke (or let it trigger proactively) for
  any change under `internal/drivers/`, `internal/proxy/`, or
  `internal/snapshot/`.
- **`platform-driver` agent** — invoke when implementing or modifying a
  platform driver, to check cross-platform parity per Principle 4.
- **`release-manager` agent** — invoke via `/release` when cutting a release.
- **`spec-manager` agent** — drives the spec-driven workflow below; invoke
  via its five commands rather than directly.
- **`/new-driver`** — scaffold a new platform driver against the shared
  interface.
- **`/bench`** — run the boot-time benchmark suite against the budgets in
  Principle 5.
- **`threat-model-review` skill** — structured STRIDE-style walkthrough for
  sandbox-boundary changes.
- **`driver-parity-check` skill** — verifies a change landed on all three
  drivers or has a documented, tracked deferral.

Full definitions live under `.claude/`.

### Spec-driven development

Non-trivial changes go through **Specify → Plan → Tasks → Implement**,
tracked as a GitHub Issue per spec (adapted from GitHub's Spec Kit — see
[docs/spec-driven-development.md](docs/spec-driven-development.md) for the
full pattern, label taxonomy, and issue structure):

- **`/specify <description>`** — open a spec issue (problem, scope,
  acceptance criteria).
- **`/plan <issue#>`** — draft the technical plan, addressing Principles 1/4/5
  where relevant.
- **`/tasks <issue#>`** — break the plan into an ordered, checkable task list.
- **`/implement <issue#>`** — work the checklist, checking items off as they
  land; closes the issue when done.
- **`/spec-status`** — end-to-end progress across all specs.

Skip this ceremony for typos, doc fixes, and single-line bug fixes — see
"When to use this" in the doc above.
