package misc

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-slim.dev/nego"
)

func TestIsTLS(t *testing.T) {
	tests := []struct {
		name string
		tls  *tls.ConnectionState
		want bool
	}{
		{name: "plain HTTP", want: false},
		{name: "TLS connection", tls: new(tls.ConnectionState), want: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			req.TLS = test.tls
			if got := IsTLS(req); got != test.want {
				t.Fatalf("IsTLS() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestIsWebSocket(t *testing.T) {
	tests := []struct {
		name       string
		upgrade    string
		connection string
		want       bool
	}{
		{name: "upgrade request", upgrade: "websocket", connection: "Upgrade", want: true},
		{name: "case insensitive values", upgrade: "WebSocket", connection: "keep-alive, UPGRADE", want: true},
		{name: "missing upgrade header", connection: "upgrade", want: false},
		{name: "missing connection upgrade", upgrade: "websocket", connection: "keep-alive", want: false},
		{name: "different protocol", upgrade: "h2c", connection: "upgrade", want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			req.Header.Set(nego.HeaderUpgrade, test.upgrade)
			req.Header.Set(nego.HeaderConnection, test.connection)
			if got := IsWebSocket(req); got != test.want {
				t.Fatalf("IsWebSocket() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestScheme(t *testing.T) {
	tests := []struct {
		name    string
		tls     bool
		headers http.Header
		want    string
	}{
		{name: "plain HTTP", want: "http"},
		{
			name:    "direct TLS takes precedence",
			tls:     true,
			headers: http.Header{nego.HeaderXForwardedProto: {"http"}},
			want:    "https",
		},
		{
			name:    "forwarded proto",
			headers: http.Header{nego.HeaderXForwardedProto: {"HTTPS"}},
			want:    "HTTPS",
		},
		{
			name: "forwarded protocol after invalid forwarded proto",
			headers: http.Header{
				nego.HeaderXForwardedProto:    {"ftp"},
				nego.HeaderXForwardedProtocol: {"wss"},
			},
			want: "wss",
		},
		{
			name:    "forwarded SSL",
			headers: http.Header{nego.HeaderXForwardedSsl: {"on"}},
			want:    "https",
		},
		{
			name:    "URL scheme",
			headers: http.Header{nego.HeaderXUrlScheme: {"ws"}},
			want:    "ws",
		},
		{
			name: "header precedence",
			headers: http.Header{
				nego.HeaderXForwardedProto:    {"http"},
				nego.HeaderXForwardedProtocol: {"https"},
				nego.HeaderXForwardedSsl:      {"on"},
				nego.HeaderXUrlScheme:         {"wss"},
			},
			want: "http",
		},
		{
			name: "invalid forwarded values",
			headers: http.Header{
				nego.HeaderXForwardedProto:    {"ftp"},
				nego.HeaderXForwardedProtocol: {"smtp"},
				nego.HeaderXForwardedSsl:      {"off"},
				nego.HeaderXUrlScheme:         {"file"},
			},
			want: "http",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			req.Header = test.headers.Clone()
			if test.tls {
				req.TLS = new(tls.ConnectionState)
			}
			if got := Scheme(req); got != test.want {
				t.Fatalf("Scheme() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestRedirect(t *testing.T) {
	for _, code := range []int{http.StatusMultipleChoices, http.StatusPermanentRedirect} {
		t.Run(http.StatusText(code), func(t *testing.T) {
			recorder := httptest.NewRecorder()
			if err := Redirect(recorder, code, "https://example.com/next"); err != nil {
				t.Fatalf("Redirect() error = %v", err)
			}
			if recorder.Code != code {
				t.Fatalf("status = %d, want %d", recorder.Code, code)
			}
			if got := recorder.Header().Get(nego.HeaderLocation); got != "https://example.com/next" {
				t.Fatalf("Location = %q, want %q", got, "https://example.com/next")
			}
			if recorder.Body.Len() != 0 {
				t.Fatalf("body = %q, want empty", recorder.Body.String())
			}
		})
	}

	for _, code := range []int{http.StatusMultipleChoices - 1, http.StatusPermanentRedirect + 1} {
		t.Run("invalid", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			err := Redirect(recorder, code, "https://example.com/next")
			if !errors.Is(err, ErrInvalidRedirectCode) {
				t.Fatalf("Redirect() error = %v, want ErrInvalidRedirectCode", err)
			}
			if recorder.Code != http.StatusOK {
				t.Fatalf("status = %d, want unmodified recorder status %d", recorder.Code, http.StatusOK)
			}
			if got := recorder.Header().Get(nego.HeaderLocation); got != "" {
				t.Fatalf("Location = %q, want empty", got)
			}
		})
	}
}

func TestRealIP(t *testing.T) {
	t.Run("uses first non-empty extractor", func(t *testing.T) {
		req := &http.Request{RemoteAddr: "203.0.113.10:8080"}
		called := false
		got := RealIP(req,
			func(*http.Request) string { return "" },
			func(*http.Request) string { return "198.51.100.20" },
			func(*http.Request) string {
				called = true
				return "192.0.2.30"
			},
		)
		if got != "198.51.100.20" {
			t.Fatalf("RealIP() = %q, want %q", got, "198.51.100.20")
		}
		if called {
			t.Fatal("RealIP() called an extractor after finding a non-empty address")
		}
	})

	t.Run("falls back to RemoteAddr", func(t *testing.T) {
		req := &http.Request{RemoteAddr: "[2001:db8::1]:8080"}
		if got := RealIP(req, func(*http.Request) string { return "" }); got != "2001:db8::1" {
			t.Fatalf("RealIP() = %q, want %q", got, "2001:db8::1")
		}
	})
}

func TestIPCheckerTrustOptions(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		options []TrustOption
		want    bool
	}{
		{name: "loopback trusted by default", ip: "127.0.0.1", want: true},
		{name: "link-local trusted by default", ip: "169.254.1.1", want: true},
		{name: "private network trusted by default", ip: "10.0.0.1", want: true},
		{name: "public network untrusted by default", ip: "203.0.113.1", want: false},
		{name: "loopback can be disabled", ip: "127.0.0.1", options: []TrustOption{TrustLoopback(false)}, want: false},
		{name: "link-local can be disabled", ip: "fe80::1", options: []TrustOption{TrustLinkLocal(false)}, want: false},
		{name: "private network can be disabled", ip: "fd00::1", options: []TrustOption{TrustPrivateNet(false)}, want: false},
		{
			name:    "additional range",
			ip:      "203.0.113.42",
			options: []TrustOption{TrustIPRange(mustParseIPNet(t, "203.0.113.0/24"))},
			want:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := newIPChecker(test.options).trust(net.ParseIP(test.ip))
			if got != test.want {
				t.Fatalf("trust(%q) = %t, want %t", test.ip, got, test.want)
			}
		})
	}
}

func TestExtractIPDirect(t *testing.T) {
	tests := []struct {
		name       string
		remoteAddr string
		want       string
	}{
		{name: "IPv4 with port", remoteAddr: "203.0.113.1:8080", want: "203.0.113.1"},
		{name: "bare IPv4", remoteAddr: "203.0.113.1", want: "203.0.113.1"},
		{name: "IPv6 with port", remoteAddr: "[2001:db8::1]:8080", want: "2001:db8::1"},
		{name: "bare IPv6", remoteAddr: "2001:db8::1", want: "2001:db8::1"},
		{name: "invalid address", remoteAddr: "not-an-address", want: ""},
	}

	extract := ExtractIPDirect()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := &http.Request{
				Header:     make(http.Header),
				RemoteAddr: test.remoteAddr,
			}
			req.Header.Set(nego.HeaderXRealIP, "10.0.0.1")
			if got := extract(req); got != test.want {
				t.Fatalf("ExtractIPDirect() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestExtractIPFromRealIPHeader(t *testing.T) {
	publicRange := mustParseIPNet(t, "203.0.113.0/24")
	ipv6Range := mustParseIPNet(t, "2001:db8::/64")
	tests := []struct {
		name       string
		realIP     string
		remoteAddr string
		options    []TrustOption
		want       string
	}{
		{name: "missing header", remoteAddr: "203.0.113.1:8080", want: "203.0.113.1"},
		{name: "malformed header", realIP: "invalid", remoteAddr: "203.0.113.1:8080", want: "203.0.113.1"},
		{name: "untrusted public header", realIP: "198.51.100.10", remoteAddr: "203.0.113.1:8080", want: "203.0.113.1"},
		{name: "private header trusted by default", realIP: "10.0.0.2", remoteAddr: "203.0.113.1:8080", want: "10.0.0.2"},
		{
			name:       "private trust disabled",
			realIP:     "10.0.0.2",
			remoteAddr: "203.0.113.1:8080",
			options:    []TrustOption{TrustPrivateNet(false)},
			want:       "203.0.113.1",
		},
		{
			name:       "additional public range",
			realIP:     "203.0.113.42",
			remoteAddr: "198.51.100.1:8080",
			options:    []TrustOption{TrustIPRange(publicRange)},
			want:       "203.0.113.42",
		},
		{
			name:       "bracketed IPv6",
			realIP:     "[2001:db8::42]",
			remoteAddr: "[2001:db9::1]:8080",
			options:    []TrustOption{TrustIPRange(ipv6Range)},
			want:       "2001:db8::42",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := &http.Request{
				Header:     make(http.Header),
				RemoteAddr: test.remoteAddr,
			}
			req.Header.Set(nego.HeaderXRealIP, test.realIP)
			if got := ExtractIPFromRealIPHeader(test.options...)(req); got != test.want {
				t.Fatalf("ExtractIPFromRealIPHeader() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestExtractIPFromXFFHeader(t *testing.T) {
	publicRange := mustParseIPNet(t, "203.0.113.0/24")
	tests := []struct {
		name       string
		xff        []string
		remoteAddr string
		options    []TrustOption
		want       string
	}{
		{name: "missing header", remoteAddr: "203.0.113.1:8080", want: "203.0.113.1"},
		{
			name:       "all hops trusted",
			xff:        []string{"10.0.0.3, 192.168.0.2"},
			remoteAddr: "127.0.0.1:8080",
			want:       "10.0.0.3",
		},
		{
			name:       "direct peer is untrusted",
			xff:        []string{"10.0.0.3"},
			remoteAddr: "203.0.113.1:8080",
			want:       "203.0.113.1",
		},
		{
			name:       "nearest untrusted hop",
			xff:        []string{"192.0.2.10, 198.51.100.20, 203.0.113.30, 10.0.0.2"},
			remoteAddr: "127.0.0.1:8080",
			options:    []TrustOption{TrustIPRange(publicRange)},
			want:       "198.51.100.20",
		},
		{
			name:       "multiple field lines",
			xff:        []string{"10.0.0.3", "192.168.0.2"},
			remoteAddr: "127.0.0.1:8080",
			want:       "10.0.0.3",
		},
		{
			name:       "bracketed IPv6",
			xff:        []string{"[fe80::3], [fe80::2]"},
			remoteAddr: "[fe80::1]:8080",
			want:       "fe80::3",
		},
		{
			name:       "malformed chain",
			xff:        []string{"invalid, 10.0.0.2"},
			remoteAddr: "127.0.0.1:8080",
			want:       "127.0.0.1",
		},
		{
			name:       "default trust can be disabled",
			xff:        []string{"10.0.0.3"},
			remoteAddr: "127.0.0.1:8080",
			options: []TrustOption{
				TrustLoopback(false),
				TrustLinkLocal(false),
				TrustPrivateNet(false),
			},
			want: "127.0.0.1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := &http.Request{
				Header:     http.Header{nego.HeaderXForwardedFor: test.xff},
				RemoteAddr: test.remoteAddr,
			}
			if got := ExtractIPFromXFFHeader(test.options...)(req); got != test.want {
				t.Fatalf("ExtractIPFromXFFHeader() = %q, want %q", got, test.want)
			}
		})
	}
}

func mustParseIPNet(t *testing.T, value string) *net.IPNet {
	t.Helper()
	_, network, err := net.ParseCIDR(value)
	if err != nil {
		t.Fatalf("net.ParseCIDR(%q): %v", value, err)
	}
	return network
}
