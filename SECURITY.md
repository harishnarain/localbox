# Security Policy

LocalBox's entire value proposition is the sandbox boundary — see
[CLAUDE.md](CLAUDE.md) Principle 1. We treat isolation failures, credential
leaks, and network-whitelist bypasses as critical.

## Reporting a vulnerability

**Do not open a public GitHub issue for a suspected vulnerability.**

Instead, report it privately using one of:

1. [GitHub Security Advisories](https://github.com/harishnarain/localbox/security/advisories/new)
   for this repo (preferred — keeps discussion and any fix coordinated in
   one place).
2. Email **harishnarain@gmail.com** with a description, reproduction steps,
   and impact assessment if you have one.

Please include:

- The affected component (driver, credential proxy, snapshot engine, CLI)
  and platform.
- Reproduction steps or a proof of concept.
- What you were able to access or do as a result (sandbox escape, host
  filesystem access, credential exposure, etc.).

## Scope

In scope:

- Sandbox escape (breaking out of the macOS/Linux/Windows isolation
  boundary).
- Credential proxy bypass or leakage of real secrets into a sandbox.
- Network whitelist bypass (sandbox reaching a non-whitelisted endpoint).
- Privilege escalation from inside a sandbox to the host.

Out of scope (for now, pre-1.0): denial-of-service against your own local
machine, issues requiring physical access, and social engineering.

## Response

This is a young, pre-1.0 project maintained part-time. Expect an
acknowledgment within a few days. There's no bug bounty program at this
stage — credit in the fix's release notes, if you want it, is what's on
offer.

## Supported versions

Pre-1.0: only the latest release/`main` is supported. A versioned support
table will replace this once there's a 1.0.
