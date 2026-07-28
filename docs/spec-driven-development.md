# Spec-driven development

LocalBox uses a pattern adapted from GitHub's [Spec Kit](https://github.com/github/spec-kit):
every non-trivial change goes through four stages — **Specify → Plan → Tasks
→ Implement** — before code lands. Spec Kit's own convention stores each
stage as a file (`spec.md`, `plan.md`, `tasks.md`) in a `specs/<feature>/`
directory. LocalBox adapts that to live natively on GitHub instead: each
spec is one **GitHub Issue**, and the stages are sections within that
issue's body. That keeps the workflow in the same place as PRs, review, and
CI — and reachable by Claude Code via `gh` — without a second file-based
store that can drift out of sync with issue state.

## The pipeline

| Stage | Command | What happens | Label after |
|---|---|---|---|
| 1. Specify | `/specify <description>` | Opens a GitHub Issue from the `spec` template: problem, scope, non-goals, acceptance criteria. Implementation-free — what/why, not how. | `spec:draft` |
| 2. Plan | `/plan <issue#>` | Fills in the issue's `### Plan` section: approach, files touched, and an explicit position on CLAUDE.md's cross-platform (Principle 4), sandbox-boundary (Principle 1), and performance-budget (Principle 5) requirements where relevant. | `spec:planned` |
| 3. Tasks | `/tasks <issue#>` | Breaks the plan into an ordered GitHub task-list checklist in the issue's `### Tasks` section. Each task should leave the repo in a working state. | `spec:tasked` |
| 4. Implement | `/implement <issue#>` | Works through the checklist, running `make check` and checking off boxes as tasks land; commits reference the issue (`Refs #n`). | `spec:in-progress` → `spec:done` (issue closed) |
| — | `/spec-status` | Rolls up every spec by stage and task-checklist completion — the end-to-end progress view. | — |

All of this is driven by the `spec-manager` agent
(`.claude/agents/spec-manager.md`); the five commands are thin wrappers
around its five stages. See that file for the exact editing protocol.

## Why the issue body carries state, not a `specs/` folder

- GitHub renders a **completion count** (e.g. "3 of 7") directly from
  `- [ ]` / `- [x]` checkboxes in an issue body — that count shows up in the
  issue list, in cross-references from PRs/commits, and on Project board
  cards, for free.
- One issue = one spec means existing GitHub mechanisms — labels,
  `Fixes #n` / `Refs #n` linking, auto-close on merge — apply without extra
  tooling.
- It avoids a second source of truth: a `specs/*.md` file and a GitHub Issue
  can drift apart; a single issue can't drift from itself.

## Label taxonomy

| Label | Meaning |
|---|---|
| `spec` | Marks an issue as a tracked spec (vs. a plain bug/feature issue) |
| `spec:draft` | Specified, not yet planned |
| `spec:planned` | Plan attached |
| `spec:tasked` | Task checklist attached, ready to implement |
| `spec:in-progress` | Implementation underway |
| `spec:done` | All tasks complete, issue closed |
| `spec:blocked` | Stalled — see the latest comment for why |

These must exist on the repo before the workflow can apply them — see
"One-time repo setup" below.

## Issue body structure

Produced by [`.github/ISSUE_TEMPLATE/spec.yml`](../.github/ISSUE_TEMPLATE/spec.yml):

```
### Problem / why
### Scope
### Non-goals
### Acceptance criteria
### Cross-platform scope
### Touches the sandbox boundary?
### Plan              <- filled in by /plan
### Tasks             <- filled in by /tasks, then checked off by /implement
```

`spec-manager` edits one `### Section` at a time (fetch body → replace
between headings → write back), so earlier sections are never disturbed by
later stages.

## GitHub Project board

Labels and checklists work standalone, but the
[**LocalBox Specs**](https://github.com/users/harishnarain/projects/2)
[Projects v2](https://docs.github.com/en/issues/planning-and-tracking-with-projects)
board gives a kanban view of every spec's stage. Set up so far, scripted via
`gh project`:

- Project created and linked to this repo (`gh project link`), so spec
  issues can be added to it.
- A single-select **Stage** field with options matching the label taxonomy
  above (Draft / Planned / Tasked / In Progress / Done / Blocked).

GitHub doesn't expose Project workflow-automation rules through the API or
`gh` — these two steps are UI-only and still need to be done by a
maintainer, once, in the project's **Settings → Workflows**:

1. **Auto-add**: filter `label:spec`, so every spec issue lands on the board
   without a manual step.
2. **Item added / label sync**: a workflow that sets **Stage** from the
   issue's `spec:*` label (six rules, one per label), and one that sets
   **Stage** to Done when the issue closes.

Until those are configured, `spec-manager` labeling an issue won't move its
card — the label (source of truth) and the board (visualization) work
independently either way.

This is one-time repo configuration, not something the harness scripts —
it's a project setting, not per-spec work, and isn't required for the
label/checklist workflow to function.

## When to use this

Use `/specify` for anything that touches more than one file/package,
changes a public interface (CLI flags, the driver interface), or touches a
sandbox-boundary or cross-platform concern. Skip it for typos, doc fixes,
and single-line bug fixes — the ceremony should be proportional to the
change; see [CONTRIBUTING.md](../CONTRIBUTING.md).

## One-time repo setup

Run once, by a maintainer — not part of the per-spec workflow:

```sh
gh label create spec               --color 5319e7 --description "Tracks a feature through the spec-driven workflow"
gh label create spec:draft         --color ededed --description "Spec written, not yet planned"
gh label create spec:planned       --color fbca04 --description "Technical plan attached"
gh label create spec:tasked        --color 1d76db --description "Task checklist attached, ready to implement"
gh label create spec:in-progress   --color 0e8a16 --description "Implementation underway"
gh label create spec:done          --color 6f42c1 --description "All tasks complete"
gh label create spec:blocked       --color b60205 --description "Stalled — see latest comment for why"
```
