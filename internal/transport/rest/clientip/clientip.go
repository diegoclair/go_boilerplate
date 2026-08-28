// Package clientip resolves the address a request actually came from when the
// application sits behind one or more proxies.
//
// It exists because the framework cannot answer this: the socket is the last
// proxy, and every header that carries the original address is written by
// whoever is in front. Trust has to be declared, never discovered — a value
// picked out of an arbitrary header is a value the caller chose.
package clientip

import (
	"net"
	"net/http"
	"strings"

	echo "github.com/labstack/echo/v5"
)

// Headers a CDN or ingress overwrites on every request it forwards, ordered
// from the most specific to the most general. Each is a single address, so
// none needs a chain parsed. Overriding the order is how a deployment says
// which edge is actually in front.
//
// X-Forwarded-For is deliberately absent: it is a chain a caller can prepend
// to, so reading it means declaring how many hops to discard — a separate
// decision, and one no default can make.
var DefaultHeaders = []string{
	"CF-Connecting-IP", // Cloudflare
	"True-Client-IP",   // Cloudflare enterprise, Akamai
	"Fly-Client-IP",    // Fly.io
	"X-Azure-ClientIP", // Azure Front Door
	"X-Real-IP",        // nginx ingress, and common convention
}

// SourceObserver is told which header answered, or an empty name when nothing
// did and the socket was used. Optional, and nil by default: an application
// that wants to see a change of edge coming decides for itself whether that is
// a metric, a log line or nothing.
type SourceObserver func(header string)

// Extractor resolves the client address from the first header that carries
// one, falling back to the socket. The fallback is the honest answer rather
// than the right one: it names the proxy, but unlike a header it cannot be
// forged by the caller.
//
// Trailing hops (AWS ALB, GCP load balancers, Heroku) publish nothing but
// X-Forwarded-For; those deployments pass echo's own XFF extractor instead,
// which takes the trusted-range options this cannot infer.
func Extractor(headers []string, observed SourceObserver) echo.IPExtractor {
	if len(headers) == 0 {
		headers = DefaultHeaders
	}

	return func(r *http.Request) string {
		for _, header := range headers {
			if ip := strings.TrimSpace(r.Header.Get(header)); ip != "" {
				observe(observed, header)
				return ip
			}
		}

		observe(observed, "")

		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return r.RemoteAddr
		}
		return host
	}
}

func observe(observed SourceObserver, header string) {
	if observed != nil {
		observed(strings.ToLower(header))
	}
}
