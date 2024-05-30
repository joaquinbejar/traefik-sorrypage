// Package traefik_sorrypage a sorrypage page plugin.
package traefik_sorrypage

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Config the plugin configuration.
type Config struct {
	Enabled         bool   `json:"enabled"`
	RedirectService string `json:"redirectService"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Enabled:         false,
		RedirectService: "",
	}
}

// SorryPage a sorrypage page plugin.
type SorryPage struct {
	next        http.Handler
	enabled     bool
	redirectURL *url.URL
	name        string
}

// New created a new SorryPage plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.RedirectService) == 0 {
		return nil, fmt.Errorf("redirectService cannot be empty")
	}

	redirectURL, err := url.Parse(config.RedirectService)
	if err != nil {
		return nil, fmt.Errorf("invalid redirect service URL: %w", err)
	}

	return &SorryPage{
		enabled:     config.Enabled,
		redirectURL: redirectURL,
		next:        next,
		name:        name,
	}, nil
}

func (a *SorryPage) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if a.sorrypageEnabled() {
		// Redirect to the specified service
		a.redirectToService(rw, req)
		return
	}

	a.next.ServeHTTP(rw, req)
}

// Redirect to the configured service
func (a *SorryPage) redirectToService(rw http.ResponseWriter, req *http.Request) {
	proxy := httputil.NewSingleHostReverseProxy(a.redirectURL)
	req.URL.Scheme = a.redirectURL.Scheme
	req.URL.Host = a.redirectURL.Host
	req.URL.Path = a.redirectURL.Path
	req.Host = a.redirectURL.Host
	proxy.ServeHTTP(rw, req)
}

// Indicates if sorrypage mode has been enabled
func (a *SorryPage) sorrypageEnabled() bool {
	if !a.enabled {
		return false
	}

	// Additional logic to check if original servers are down can be implemented here
	// For example, checking a specific file, environment variable, etc.
	return true
}
