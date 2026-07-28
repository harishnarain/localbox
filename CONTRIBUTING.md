# Contributing to LocalBox

Thanks for considering a contribution. Read [CLAUDE.md](CLAUDE.md) first —
it's the constitution every change is expected to hold to, not just guidance
for AI agents.

## Dev setup

- Go 1.23+ (latest stable recommended)
- `make check` runs everything CI runs: `gofmt`, `go vet`, `go test ./...`

```
git clone https://github.com/harishnarain/localbox.git
cd localbox
make check
```

## Git workflow

- **`main` is protected.** No direct pushes; all changes land via PR.
- **Branch naming:** `type/short-description`, e.g. `feat/linux-lxc-driver`,
  `fix/proxy-token-leak`, `docs/readme-quickstart`. Types match the commit
  convention below.
- **Commits:** [Conventional Commits](https://www.conventionalcommits.org/)
  (`feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`, `ci:`). Release
  notes are generated from these, so keep the subject line accurate.
- **Sign off your commits** (`git commit -s`) — adds a `Signed-off-by`
  trailer certifying you have the right to submit the contribution under
  Apache 2.0 (Developer Certificate of Origin).
- **PRs:** use the PR template, keep them scoped to one logical change, and
  make sure `make check` passes locally before requesting review.

## Review policy

- One maintainer approval required to merge.
- CI (`gofmt`, `go vet`, `go test`, CodeQL) must be green.
- **Security-sensitive paths** — anything under `internal/drivers/`,
  `internal/proxy/`, or `internal/snapshot/` — require an explicit security
  review per [CLAUDE.md](CLAUDE.md) Principle 1. Use the `security-reviewer`
  agent locally before opening the PR; it's not a substitute for the human
  review, but it catches the obvious stuff first.
- **Driver changes** — anything under `internal/drivers/<platform>/` should
  either land alongside the equivalent change for the other two platforms, or
  the PR description must explain why parity is deferred and link a tracked
  issue. The `platform-driver` agent and `driver-parity-check` skill help
  with this.
- Changes that touch the boot-time hot path should include `/bench` output
  in the PR description if there's any risk of regressing the budgets in
  [CLAUDE.md](CLAUDE.md) Principle 5.

## Spec-driven development

Non-trivial changes — anything touching more than one file/package, a public
interface, or a sandbox-boundary/cross-platform concern — should go through
`/specify → /plan → /tasks → /implement` before code lands, tracked as a
GitHub Issue per spec. See
[docs/spec-driven-development.md](docs/spec-driven-development.md) for the
full pattern. Skip it for typos, doc fixes, and single-line bug fixes — the
ceremony should match the size of the change.

## Using the agent harness

This repo is set up to be worked on collaboratively with Claude Code (or
other agent CLIs). `.claude/agents`, `.claude/commands`, and `.claude/skills`
encode project-specific review and workflow logic — see the "Using the agent
harness" section of [CLAUDE.md](CLAUDE.md) for what's available and when to
reach for each one. You don't have to use them, but they exist to catch the
same things a human reviewer would flag.

## Reporting bugs / requesting features

Use the GitHub issue templates. For anything that looks like a security
vulnerability (a way to escape the sandbox, leak a credential, or defeat
network whitelisting), do **not** open a public issue — follow
[SECURITY.md](SECURITY.md) instead.

## Code of conduct

This project follows the [Code of Conduct](CODE_OF_CONDUCT.md).
