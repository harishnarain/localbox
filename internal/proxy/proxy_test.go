package proxy

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const testSecretEnvVar = "LOCALBOX_TEST_PROXY_SECRET"

func TestForwardsWithSecretSubstitution(t *testing.T) {
	const realSecret = "sk-real-secret-value"
	t.Setenv(testSecretEnvVar, realSecret)

	var gotAuth string
	var gotPath string
	upstream := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotPath = r.URL.Path
		w.Header().Set("X-Upstream", "yes")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("upstream response body"))
	}))
	defer upstream.Close()

	svc := Service{
		Name:             "anthropic",
		UpstreamBaseURL:  upstream.URL, // httptest TLS server URL is https://127.0.0.1:port
		SecretEnvVar:     testSecretEnvVar,
		PlaceholderToken: "placeholder-token",
	}

	// buildUpstreamRequest requires an https:// URL; httptest.NewTLSServer
	// produces exactly that (https://127.0.0.1:<port>).
	if !strings.HasPrefix(svc.UpstreamBaseURL, "https://") {
		t.Fatalf("expected httptest TLS server URL to be https, got %q", svc.UpstreamBaseURL)
	}

	p, err := New([]Service{svc})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	// Use the upstream's own client (which trusts its self-signed cert)
	// as the proxy's outbound client so the TLS handshake succeeds.
	p.client = upstream.Client()
	if err := p.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = p.Stop(ctx)
	}()

	req, _ := http.NewRequest(http.MethodGet, "http://"+p.Addr().String()+"/anthropic/v1/messages", nil)
	req.Header.Set("Authorization", "placeholder-token")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request to proxy failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if gotAuth != realSecret {
		t.Fatalf("upstream saw Authorization=%q, want real secret", gotAuth)
	}
	if gotPath != "/v1/messages" {
		t.Fatalf("upstream saw path=%q, want /v1/messages", gotPath)
	}
	if got := resp.Header.Get("X-Upstream"); got != "yes" {
		t.Fatalf("expected relayed X-Upstream header, got %q", got)
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "upstream response body" {
		t.Fatalf("expected relayed body, got %q", string(body))
	}

	// The real secret must never be observable to the caller: it's an
	// inbound-only substitution, so the response the caller sees (body,
	// headers, status) must not contain it.
	if strings.Contains(string(body), realSecret) || strings.Contains(resp.Header.Get("X-Upstream"), realSecret) {
		t.Fatalf("real secret leaked into response observable by caller")
	}
}

func TestRejectsMissingOrWrongToken(t *testing.T) {
	t.Setenv(testSecretEnvVar, "sk-real-secret-value")

	upstreamCalled := false
	upstream := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upstreamCalled = true
		w.WriteHeader(http.StatusOK)
	}))
	defer upstream.Close()

	svc := Service{
		Name:             "github",
		UpstreamBaseURL:  upstream.URL,
		SecretEnvVar:     testSecretEnvVar,
		PlaceholderToken: "correct-token",
	}
	p, err := New([]Service{svc})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	p.client = upstream.Client()
	if err := p.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = p.Stop(ctx)
	}()

	base := "http://" + p.Addr().String()

	cases := []struct {
		name  string
		token string
	}{
		{"missing token", ""},
		{"wrong token", "wrong-token"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, base+"/github/repos", nil)
			if tc.token != "" {
				req.Header.Set("Authorization", tc.token)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusForbidden {
				t.Fatalf("expected 403, got %d", resp.StatusCode)
			}
		})
	}
	if upstreamCalled {
		t.Fatalf("upstream must not be called for an unauthenticated request")
	}
}

func TestRejectsUnknownService(t *testing.T) {
	upstream := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("upstream must not be called for an unknown service")
	}))
	defer upstream.Close()

	svc := Service{
		Name:             "anthropic",
		UpstreamBaseURL:  upstream.URL,
		SecretEnvVar:     testSecretEnvVar,
		PlaceholderToken: "placeholder-token",
	}
	p, err := New([]Service{svc})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	p.client = upstream.Client()
	if err := p.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = p.Stop(ctx)
	}()

	base := "http://" + p.Addr().String()

	for _, path := range []string{"/openai/v1/models", "/"} {
		req, _ := http.NewRequest(http.MethodGet, base+path, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("request to %q failed: %v", path, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("path %q: expected 404, got %d", path, resp.StatusCode)
		}
	}
}

func TestNewRejectsInvalidServices(t *testing.T) {
	cases := []struct {
		name string
		svc  Service
	}{
		{"empty name", Service{Name: "", UpstreamBaseURL: "https://example.com", SecretEnvVar: "X", PlaceholderToken: "t"}},
		{"name with slash", Service{Name: "a/b", UpstreamBaseURL: "https://example.com", SecretEnvVar: "X", PlaceholderToken: "t"}},
		{"empty placeholder token", Service{Name: "a", UpstreamBaseURL: "https://example.com", SecretEnvVar: "X", PlaceholderToken: ""}},
		{"empty secret env var", Service{Name: "a", UpstreamBaseURL: "https://example.com", SecretEnvVar: "", PlaceholderToken: "t"}},
		{"non-https upstream", Service{Name: "a", UpstreamBaseURL: "http://example.com", SecretEnvVar: "X", PlaceholderToken: "t"}},
		{"relative upstream", Service{Name: "a", UpstreamBaseURL: "/not-absolute", SecretEnvVar: "X", PlaceholderToken: "t"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := New([]Service{tc.svc}); err == nil {
				t.Fatalf("expected New to reject invalid service, got nil error")
			}
		})
	}

	// Duplicate names.
	dup := Service{Name: "a", UpstreamBaseURL: "https://example.com", SecretEnvVar: "X", PlaceholderToken: "t"}
	if _, err := New([]Service{dup, dup}); err == nil {
		t.Fatalf("expected New to reject duplicate service names, got nil error")
	}
}

func TestFailsClosedWhenSecretEnvVarUnset(t *testing.T) {
	// Deliberately do not set the env var.
	svc := Service{
		Name:             "anthropic",
		UpstreamBaseURL:  "https://example.invalid",
		SecretEnvVar:     "LOCALBOX_TEST_PROXY_SECRET_UNSET",
		PlaceholderToken: "placeholder-token",
	}
	p, err := New([]Service{svc})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := p.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = p.Stop(ctx)
	}()

	req, _ := http.NewRequest(http.MethodGet, "http://"+p.Addr().String()+"/anthropic/v1/messages", nil)
	req.Header.Set("Authorization", "placeholder-token")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 when secret env var is unset, got %d", resp.StatusCode)
	}
}
