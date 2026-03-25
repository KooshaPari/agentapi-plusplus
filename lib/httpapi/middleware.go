package httpapi

import (
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

// hostAuthorizationMiddleware enforces that the request Host header matches one of the allowed
// hosts, ignoring any port in the comparison. If allowedHosts is empty, all hosts are allowed.
// Always uses url.Parse("http://" + r.Host) to robustly extract the hostname (handles IPv6).
func hostAuthorizationMiddleware(allowedHosts []string, badHostHandler http.Handler) func(next http.Handler) http.Handler {
	// Copy for safety; also build a map for O(1) lookups with case-insensitive keys.
	allowed := make(map[string]struct{}, len(allowedHosts))
	for _, h := range allowedHosts {
		allowed[strings.ToLower(h)] = struct{}{}
	}
	wildcard := slices.Contains(allowedHosts, "*")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if wildcard { // wildcard semantics: allow all
				next.ServeHTTP(w, r)
				return
			}
			// Extract hostname from the Host header using url.Parse; ignore any port.
			hostHeader := r.Host
			if hostHeader == "" {
				badHostHandler.ServeHTTP(w, r)
				return
			}
			if u, err := url.Parse("http://" + hostHeader); err == nil {
				hostname := u.Hostname()
				if _, ok := allowed[strings.ToLower(hostname)]; ok {
					next.ServeHTTP(w, r)
					return
				}
			}
			badHostHandler.ServeHTTP(w, r)
		})
	}
}

// sseMiddleware creates middleware that prevents proxy buffering for SSE endpoints
func sseMiddleware(ctx huma.Context, next func(huma.Context)) {
	// Disable proxy buffering for SSE endpoints
	ctx.SetHeader("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.SetHeader("Pragma", "no-cache")
	ctx.SetHeader("Expires", "0")
	ctx.SetHeader("X-Accel-Buffering", "no") // nginx
	ctx.SetHeader("X-Proxy-Buffering", "no") // generic proxy
	ctx.SetHeader("Connection", "keep-alive")

	next(ctx)
}
