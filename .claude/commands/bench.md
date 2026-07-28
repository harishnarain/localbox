---
description: Run the sandbox boot-time benchmark suite and compare against CLAUDE.md's per-platform budgets
---

Run LocalBox's boot-time benchmarks and report against the budgets in
CLAUDE.md Principle 5 (macOS < 500ms, Linux < 100ms, Windows < 1.5s).

1. Look for benchmark tests (`func Benchmark...` in `_test.go` files, likely
   under `internal/drivers/` or a top-level `bench/` once it exists). If none
   exist yet, say so explicitly rather than fabricating numbers — don't
   report a budget as "met" without having actually run something.
2. If benchmarks exist, run them for the current platform:
   `go test ./... -bench=Boot -benchtime=10x -run=^$`
3. Compare the result against this platform's budget and report pass/fail
   plainly, including the actual measured number.
4. If this is being run to validate a specific PR/change, note whether the
   change is likely responsible for any regression, and point at the
   specific code path if so.
