---
description: Cut a LocalBox release — changelog, version bump, and pre-release checklist (draft only, no tagging/pushing)
argument-hint: [major|minor|patch]
---

Use the `release-manager` agent to prepare a release for LocalBox
(bump hint, if given: `$ARGUMENTS`).

The agent should draft the changelog and version bump and walk the
pre-release checklist (CI status, open-core boundary, cross-platform
parity), then present the result for the maintainer to confirm. Do not tag
or push anything without explicit confirmation — releasing is a one-way,
externally-visible action.
