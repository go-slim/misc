package misc

import (
	"errors"
	"net"
	"net/http"
	"strings"

	"go-slim.dev/nego"
)

// ErrInvalidRedirectCode is returned when [Redirect] receives a status code
// outside the supported 300 through 308 range.
var ErrInvalidRedirectCode = errors.New("invalid redirect status code")

// IsTLS reports whether req was received over a TLS connection, as indicated
// by a non-nil [http.Request.TLS] field. It does not inspect forwarded headers.
func IsTLS(req *http.Request) bool {
	return req.TLS != nil
}

// IsWebSocket reports whether req contains a WebSocket upgrade request. Header
// names and values are matched case-insensitively.
func IsWebSocket(req *http.Request) bool {
	upgrade := req.Header.Get(nego.HeaderUpgrade)
	connection := req.Header.Get(nego.HeaderConnection)
	return strings.EqualFold(upgrade, "websocket") && strings.Contains(strings.ToLower(connection), "upgrade")
}

// Scheme returns the request scheme.
//
// A direct TLS connection always returns "https". Otherwise Scheme checks, in
// order, X-Forwarded-Proto, X-Forwarded-Protocol, X-Forwarded-Ssl, and
// X-Url-Scheme before falling back to "http". Forwarded scheme values may be
// "http", "https", "ws", or "wss" and preserve their original case.
//
// Scheme does not verify that forwarded headers came from a trusted proxy.
// Applications must ignore or remove client-supplied forwarding headers at
// their trust boundary.
func Scheme(req *http.Request) string {
	// Request.URL.Scheme is normally empty for incoming server requests.
	if IsTLS(req) {
		return "https"
	}
	if scheme := req.Header.Get(nego.HeaderXForwardedProto); isValidProto(scheme) {
		return scheme
	}
	if scheme := req.Header.Get(nego.HeaderXForwardedProtocol); isValidProto(scheme) {
		return scheme
	}
	if ssl := req.Header.Get(nego.HeaderXForwardedSsl); ssl == "on" {
		return "https"
	}
	if scheme := req.Header.Get(nego.HeaderXUrlScheme); isValidProto(scheme) {
		return scheme
	}
	return "http"
}

func isValidProto(proto string) bool {
	if proto == "" {
		return false
	}
	for _, p := range []string{"http", "https", "ws", "wss"} {
		if strings.EqualFold(proto, p) {
			return true
		}
	}
	return false
}

// Redirect sets the Location response header and writes code.
//
// Codes from 300 through 308 are accepted. An unsupported code returns
// [ErrInvalidRedirectCode] without modifying the response. Redirect writes no
// response body; use [http.Redirect] when its escaping and response-body
// behavior is desired.
func Redirect(w http.ResponseWriter, code int, url string) error {
	if code < 300 || code > 308 {
		return ErrInvalidRedirectCode
	}
	w.Header().Set(nego.HeaderLocation, url)
	w.WriteHeader(code)
	return nil
}

// RealIP returns the first non-empty client address produced by extractors, in
// argument order. If no extractor is supplied, or every extractor returns an
// empty string, it falls back to the host portion of req.RemoteAddr.
//
// RealIP does not read proxy headers by default. Applications behind a proxy
// must explicitly supply [ExtractIPFromXFFHeader] or
// [ExtractIPFromRealIPHeader] and configure their trust options to match the
// deployment topology.
func RealIP(req *http.Request, extractors ...IPExtractor) string {
	for _, extractor := range extractors {
		if ip := extractor(req); ip != "" {
			return ip
		}
	}

	// RemoteAddr identifies the directly connected peer, which may be a proxy.
	return extractIP(req)
}

type ipChecker struct {
	trustExtraRanges []*net.IPNet
	trustLoopback    bool
	trustLinkLocal   bool
	trustPrivateNet  bool
}

// TrustOption configures which addresses proxy-aware [IPExtractor] functions
// treat as trusted. Loopback, link-local unicast, and private addresses are
// trusted by default; the corresponding options can disable each category.
type TrustOption func(*ipChecker)

// TrustLoopback controls whether loopback addresses are trusted. The default is
// true.
func TrustLoopback(v bool) TrustOption {
	return func(c *ipChecker) {
		c.trustLoopback = v
	}
}

// TrustLinkLocal controls whether link-local unicast addresses are trusted. The
// default is true.
func TrustLinkLocal(v bool) TrustOption {
	return func(c *ipChecker) {
		c.trustLinkLocal = v
	}
}

// TrustPrivateNet controls whether addresses reported by [net.IP.IsPrivate] are
// trusted. The default is true.
func TrustPrivateNet(v bool) TrustOption {
	return func(c *ipChecker) {
		c.trustPrivateNet = v
	}
}

// TrustIPRange adds an IP network to the trusted ranges. ipRange must not be
// nil.
func TrustIPRange(ipRange *net.IPNet) TrustOption {
	return func(c *ipChecker) {
		c.trustExtraRanges = append(c.trustExtraRanges, ipRange)
	}
}

func newIPChecker(configs []TrustOption) *ipChecker {
	checker := &ipChecker{trustLoopback: true, trustLinkLocal: true, trustPrivateNet: true}
	for _, configure := range configs {
		configure(checker)
	}
	return checker
}

func (c *ipChecker) trust(ip net.IP) bool {
	if c.trustLoopback && ip.IsLoopback() {
		return true
	}
	if c.trustLinkLocal && ip.IsLinkLocalUnicast() {
		return true
	}
	if c.trustPrivateNet && ip.IsPrivate() {
		return true
	}
	for _, trustedRange := range c.trustExtraRanges {
		if trustedRange.Contains(ip) {
			return true
		}
	}
	return false
}

// IPExtractor resolves a client address from an HTTP request. Returning an
// empty string allows [RealIP] to try the next extractor before falling back to
// the directly connected peer.
//
// Extractors that consume proxy headers must reflect the application's actual
// proxy topology. The edge proxy must replace untrusted client-supplied
// forwarding headers, otherwise those headers can be spoofed.
//
// The trust-chain model used by the provided extractors is adapted from Echo's
// IP address extraction design: https://echo.labstack.com/docs/ip-address
type IPExtractor func(*http.Request) string

// ExtractIPDirect returns an extractor that reads only req.RemoteAddr. It
// returns the host from a host:port address, accepts a bare IP address, and
// returns an empty string for malformed input. It never reads forwarding
// headers.
//
// Use this extractor when the server is directly exposed to clients or when
// the address of the directly connected proxy is desired.
func ExtractIPDirect() IPExtractor {
	return extractIP
}

func extractIP(req *http.Request) string {
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		if net.ParseIP(req.RemoteAddr) != nil {
			return req.RemoteAddr
		}
		return ""
	}
	return host
}

// ExtractIPFromRealIPHeader returns an extractor for X-Real-IP.
//
// A syntactically valid header value is returned only when it belongs to a
// trusted address category or a range added with [TrustIPRange]. Otherwise the
// extractor falls back to req.RemoteAddr. Surrounding brackets on an IPv6
// address are accepted.
func ExtractIPFromRealIPHeader(options ...TrustOption) IPExtractor {
	checker := newIPChecker(options)
	return func(req *http.Request) string {
		realIP := req.Header.Get(nego.HeaderXRealIP)
		if realIP != "" {
			realIP = strings.TrimPrefix(realIP, "[")
			realIP = strings.TrimSuffix(realIP, "]")
			if ip := net.ParseIP(realIP); ip != nil && checker.trust(ip) {
				return realIP
			}
		}
		return extractIP(req)
	}
}

// ExtractIPFromXFFHeader returns an extractor for X-Forwarded-For.
//
// The extractor combines every X-Forwarded-For field value, appends the direct
// peer from req.RemoteAddr, and scans the resulting chain from right to left.
// It returns the first untrusted address. If every address is trusted, it
// returns the leftmost address as a best-effort client address. A malformed
// address anywhere in the chain causes the extractor to return the direct peer.
func ExtractIPFromXFFHeader(options ...TrustOption) IPExtractor {
	checker := newIPChecker(options)
	return func(req *http.Request) string {
		directIP := extractIP(req)
		xffs := req.Header[nego.HeaderXForwardedFor]
		if len(xffs) == 0 {
			return directIP
		}
		ips := append(strings.Split(strings.Join(xffs, ","), ","), directIP)
		for i := len(ips) - 1; i >= 0; i-- {
			ips[i] = strings.TrimSpace(ips[i])
			ips[i] = strings.TrimPrefix(ips[i], "[")
			ips[i] = strings.TrimSuffix(ips[i], "]")
			ip := net.ParseIP(ips[i])
			if ip == nil {
				// One malformed hop makes the complete forwarded chain unreliable.
				return directIP
			}
			if !checker.trust(ip) {
				return ip.String()
			}
		}
		// The leftmost address is furthest from the server.
		return strings.TrimSpace(ips[0])
	}
}
