// Package proxy implements LocalBox's zero-trust credential proxy: a
// host-side, loopback-only HTTP reverse proxy that lets sandbox processes
// call a placeholder-token-protected local endpoint per configured
// service, while the proxy substitutes the real secret (read from a
// host-side environment variable, never mounted into the sandbox) and
// forwards the request to the real HTTPS upstream. See CLAUDE.md
// Principle 3.
//
// This is an MVP: the service list is supplied by the caller as a static
// []Service at construction time. Config-file/env-var-driven service
// loading, driver/CLI wiring, credential refresh/rotation, and network
// domain whitelisting beyond the configured service allowlist are not yet
// implemented.
package proxy
