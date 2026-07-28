---
name: threat-model-review
description: Structured threat-model walkthrough for any change touching the sandbox boundary — internal/drivers/, internal/proxy/, or internal/snapshot/. Use before merging, or when asked to threat-model a change.
---

# Threat-model review

A structured checklist for reviewing changes that touch LocalBox's sandbox
boundary (CLAUDE.md Principle 1). This is a thinking tool, not a substitute
for the `security-reviewer` agent — use it to structure that review, or run
it standalone for a quick self-check before opening a PR.

Walk through each category for the change under review. For each, either
state "not applicable, because ___" or name the concrete risk.

## 1. Spoofing
Can something inside the sandbox impersonate something it shouldn't — the
host user, another sandbox, a whitelisted network endpoint?

## 2. Tampering
Can data crossing the sandbox boundary (the workspace mount, the extracted
git diff, config passed in) be modified in a way that isn't expected —
path traversal in a diff, a symlink escaping the mounted workspace, a
config value injected via an unexpected channel?

## 3. Repudiation
If something goes wrong, can you tell what the sandbox actually did? Does
this change preserve enough of an audit trail (logs, the diff itself) to
reconstruct what happened?

## 4. Information disclosure
Does anything reach the sandbox that shouldn't — a real credential instead
of a proxy token, host environment variables beyond what's needed, host
filesystem paths outside the workspace mount?

## 5. Denial of service
Can a sandboxed agent (accidentally or via prompt injection) exhaust host
resources — disk via snapshot growth, memory/CPU without limits, ports?

## 6. Elevation of privilege
Does this change give the sandbox — or a component acting on its behalf —
more capability on the host than the minimum needed for the feature?

## Output

For each category with a real finding: state the concrete scenario (inputs
→ outcome), not just the category name. If nothing applies, say so briefly
and move on — don't pad the review to look thorough.
