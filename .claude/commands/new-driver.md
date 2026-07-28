---
description: Scaffold a new platform driver against the shared LocalBox driver interface
argument-hint: <platform-name>
---

Scaffold a new sandbox driver for platform `$ARGUMENTS` under
`internal/drivers/$ARGUMENTS/`.

1. Read the existing drivers under `internal/drivers/{macos,linux}/`
   to find the shared interface they implement (if it doesn't exist yet as
   an extracted interface, look at how the orchestrator expects to call into
   a driver, and infer the contract from there).
2. Create `internal/drivers/$ARGUMENTS/doc.go` with a package comment
   describing: the underlying isolation primitive, the target boot budget
   (ask if not already defined in CLAUDE.md Principle 5 — don't invent a
   number), and current implementation status.
3. Implement the driver skeleton against the shared interface with clear
   `// TODO` markers for unimplemented behavior — don't fabricate working
   logic for a hypervisor/container API you haven't actually verified.
4. Add a table entry for the new platform to the boot-budget table in
   CLAUDE.md and the platform table in README.md, matching the existing
   format.
5. Use the `platform-driver` agent to sanity-check parity with the existing
   drivers before finishing.
