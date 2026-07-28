package proxy

import (
	"context"
	"crypto/hmac"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// defaultHeader is the header substituted with the real secret when a
// Service does not specify one.
const defaultHeader = "Authorization"

// Service configures one upstream that sandbox processes may reach through
// the proxy under the path prefix "/<Name>/...".
type Service struct {
	// Name selects this service in the request path: "/<Name>/...".
	// Must be non-empty and unique across the services passed to New.
	Name string

	// UpstreamBaseURL is the real, real HTTPS base URL requests are
	// forwarded to. The remainder of the incoming request path (after
	// the "/<Name>" prefix) is appended to it. Must be an absolute
	// "https://" URL.
	UpstreamBaseURL string

	// SecretEnvVar is the name of a host-side environment variable
	// holding the real secret. It is read at request time (never cached
	// in a log-visible field) and is never mounted into a sandbox.
	SecretEnvVar string

	// PlaceholderToken is the dummy value a sandbox process must present
	// in Header for the request to be forwarded. Must be non-empty.
	PlaceholderToken string

	// Header is the HTTP header carrying the placeholder token inbound,
	// and the real secret outbound. Defaults to "Authorization".
	Header string
}

func (s Service) header() string {
	if s.Header == "" {
		return defaultHeader
	}
	return s.Header
}

// Proxy is a loopback-only HTTP reverse proxy that substitutes real
// secrets for placeholder tokens at the host boundary. See the package
// doc comment for the security model.
type Proxy struct {
	services map[string]Service
	client   *http.Client

	server   *http.Server
	listener net.Listener
}

// New builds a Proxy for the given services. It fails fast (returns an
// error) if any service is misconfigured: empty/duplicate name, empty
// placeholder token, or an UpstreamBaseURL that isn't an absolute
// "https://" URL.
func New(services []Service) (*Proxy, error) {
	m := make(map[string]Service, len(services))
	for _, svc := range services {
		if svc.Name == "" {
			return nil, fmt.Errorf("proxy: service has empty Name")
		}
		if strings.Contains(svc.Name, "/") {
			return nil, fmt.Errorf("proxy: service %q: Name must not contain '/'", svc.Name)
		}
		if _, exists := m[svc.Name]; exists {
			return nil, fmt.Errorf("proxy: duplicate service name %q", svc.Name)
		}
		if svc.PlaceholderToken == "" {
			return nil, fmt.Errorf("proxy: service %q: PlaceholderToken must not be empty", svc.Name)
		}
		if svc.SecretEnvVar == "" {
			return nil, fmt.Errorf("proxy: service %q: SecretEnvVar must not be empty", svc.Name)
		}
		u, err := url.Parse(svc.UpstreamBaseURL)
		if err != nil {
			return nil, fmt.Errorf("proxy: service %q: invalid UpstreamBaseURL: %w", svc.Name, err)
		}
		if !u.IsAbs() || u.Scheme != "https" {
			return nil, fmt.Errorf("proxy: service %q: UpstreamBaseURL must be an absolute https:// URL", svc.Name)
		}
		m[svc.Name] = svc
	}

	p := &Proxy{
		services: m,
		client:   &http.Client{Timeout: 30 * time.Second},
	}
	return p, nil
}

// Start binds a loopback-only listener on an OS-assigned ephemeral port
// and begins serving in the background. Callers must call Addr to
// discover the bound port, and Stop to shut the proxy down.
func (p *Proxy) Start() error {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("proxy: listen: %w", err)
	}
	p.listener = ln
	p.server = &http.Server{Handler: p}

	go func() {
		// ErrServerClosed is the expected return from a graceful Stop.
		_ = p.server.Serve(ln)
	}()
	return nil
}

// Addr returns the address the proxy is listening on. Only valid after a
// successful Start.
func (p *Proxy) Addr() net.Addr {
	return p.listener.Addr()
}

// Stop gracefully shuts the proxy down, waiting for in-flight requests to
// complete or ctx to be done, whichever comes first.
func (p *Proxy) Stop(ctx context.Context) error {
	if p.server == nil {
		return nil
	}
	return p.server.Shutdown(ctx)
}

// hopByHopHeaders lists headers that are connection-specific and must not
// be forwarded, per RFC 7230 section 6.1.
var hopByHopHeaders = []string{
	"Connection",
	"Proxy-Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",
	"Trailer",
	"Transfer-Encoding",
	"Upgrade",
}

// ServeHTTP implements http.Handler. It selects a Service from the first
// path segment, validates the caller's placeholder token in constant
// time, and — only on success — substitutes the real secret and forwards
// the request to the configured upstream, relaying the response back
// unmodified.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name, rest, ok := splitServicePath(r.URL.Path)
	if !ok {
		http.NotFound(w, r)
		return
	}
	svc, ok := p.services[name]
	if !ok {
		http.NotFound(w, r)
		return
	}

	header := svc.header()
	presented := r.Header.Get(header)
	if presented == "" || !constantTimeEqual(presented, svc.PlaceholderToken) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	secret := os.Getenv(svc.SecretEnvVar)
	if secret == "" {
		// Fail closed: an unconfigured host-side secret must not result
		// in a request going out with no auth header at all.
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}

	upstreamReq, err := buildUpstreamRequest(r, svc.UpstreamBaseURL, rest, header, secret)
	if err != nil {
		http.Error(w, "bad gateway", http.StatusBadGateway)
		return
	}

	resp, err := p.client.Do(upstreamReq)
	if err != nil {
		http.Error(w, "bad gateway", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	relayResponse(w, resp)
}

// splitServicePath splits "/<name>/<rest>" into name and "/<rest>". ok is
// false if the path has no service segment (e.g. "/" or "").
func splitServicePath(path string) (name, rest string, ok bool) {
	trimmed := strings.TrimPrefix(path, "/")
	if trimmed == "" {
		return "", "", false
	}
	idx := strings.IndexByte(trimmed, '/')
	if idx == -1 {
		return trimmed, "", true
	}
	return trimmed[:idx], trimmed[idx:], true
}

// constantTimeEqual reports whether a and b are equal, without leaking
// timing information about a shared prefix.
func constantTimeEqual(a, b string) bool {
	return hmac.Equal([]byte(a), []byte(b))
}

// buildUpstreamRequest constructs the outbound request to the real
// upstream, copying the inbound method/body/headers (minus hop-by-hop
// headers), and substituting the real secret into header.
func buildUpstreamRequest(r *http.Request, upstreamBase, rest, header, secret string) (*http.Request, error) {
	base, err := url.Parse(upstreamBase)
	if err != nil {
		return nil, err
	}
	restURL, err := url.Parse(rest)
	if err != nil {
		return nil, err
	}
	target := base.ResolveReference(&url.URL{Path: base.Path + restURL.Path, RawQuery: r.URL.RawQuery})

	req, err := http.NewRequestWithContext(r.Context(), r.Method, target.String(), r.Body)
	if err != nil {
		return nil, err
	}

	req.Header = r.Header.Clone()
	for _, h := range hopByHopHeaders {
		req.Header.Del(h)
	}
	req.Header.Set(header, secret)
	req.ContentLength = r.ContentLength
	req.Host = base.Host

	return req, nil
}

// relayResponse copies status, headers (minus hop-by-hop headers), and
// body from resp to w unmodified.
func relayResponse(w http.ResponseWriter, resp *http.Response) {
	dst := w.Header()
	for k, vv := range resp.Header {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
	for _, h := range hopByHopHeaders {
		dst.Del(h)
	}
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}
