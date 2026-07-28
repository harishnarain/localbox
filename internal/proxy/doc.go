// Package proxy implements LocalBox's zero-trust credential proxy: a
// host-side loopback proxy that substitutes scoped dummy tokens for real
// secrets (Anthropic, OpenAI, GitHub, SSH, ...) in transit, and enforces
// outbound network domain whitelisting. See CLAUDE.md Principle 3.
// Not yet implemented.
package proxy
