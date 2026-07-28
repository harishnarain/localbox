## What & why

<!-- What does this change do, and why? Link any related issue. -->

## Checklist

- [ ] Commit messages follow [Conventional Commits](https://www.conventionalcommits.org/)
- [ ] Commits are signed off (`git commit -s`) — see CONTRIBUTING.md
- [ ] `make check` passes locally
- [ ] If this touches `internal/drivers/`, `internal/proxy/`, or
      `internal/snapshot/`: security review requested (CLAUDE.md Principle 1)
- [ ] If this touches a platform driver: the other two platforms are updated
      too, or the PR description explains why parity is deferred and links a
      tracked issue (CLAUDE.md Principle 4)
- [ ] If this touches the boot/exec hot path: `/bench` output included, or
      no risk of regressing the budgets in CLAUDE.md Principle 5

## How was this tested?

<!-- Commands run, platforms tested on, etc. -->
