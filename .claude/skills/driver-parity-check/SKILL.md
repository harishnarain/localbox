---
name: driver-parity-check
description: Verify a change to internal/drivers/{macos,linux}/ has been mirrored across both platform drivers (macOS, Linux — Windows uses the Linux driver via WSL2), or has an explicit, tracked deferral. Use before merging any driver change, per CLAUDE.md Principle 4.
---

# Driver parity check

CLAUDE.md Principle 4: a feature isn't "done" on one platform driver — it's
done on both, or parity is explicitly deferred with a tracked issue.
This skill is the check that enforces that before merge.

## Steps

1. Identify what changed in `internal/drivers/<platform>/` for this PR/diff.
2. For the other platform, check whether the equivalent
   behavior already exists:
   - If yes: confirm it's actually equivalent (not just superficially
     similar) — behavior, error handling, and any config surface should
     match in intent even if the implementation differs per platform.
   - If no: this is a parity gap. Check whether the PR description names a
     tracked issue and a reason for deferral (CONTRIBUTING.md requires
     this). If not, flag it — don't silently let a gap through.
3. If the change is to the *shared driver interface* itself (not a single
   platform's implementation), both drivers must compile against it;
   confirm neither was left implementing a stale version of the interface.
4. If a platform genuinely cannot support the feature (a primitive that
   doesn't exist on that OS), that's not a "gap to fill later" — it should
   be modeled explicitly (a capability flag, a documented limitation), not
   left as an implicit TODO.

## Output

Report per-platform status (mirrored / not applicable / gap-with-tracked-
issue / unflagged gap). Unflagged gaps are the finding that matters most —
say so plainly.
