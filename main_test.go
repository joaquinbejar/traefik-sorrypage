package traefik_sorrypage

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

const sorryPageBody = "<html><head></head><body>SorryPage</body></html>"

// newSorryPageServer starts a fake "sorry page" backend that records the
// requests it receives and answers with a fixed 503 HTML page.
func newSorryPageServer(t *testing.T, received *[]*http.Request) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		*received = append(*received, req.Clone(req.Context()))
		rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		rw.WriteHeader(http.StatusServiceUnavailable)
		_, _ = rw.Write([]byte(sorryPageBody))
	}))
	t.Cleanup(server.Close)

	return server
}

func TestSorryPageEnabledProxiesToService(t *testing.T) {
	var received []*http.Request
	sorryServer := newSorryPageServer(t, &received)

	cfg := CreateConfig()
	cfg.Enabled = true
	cfg.RedirectService = sorryServer.URL

	ctx := context.Background()
	nextCalled := false
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) { nextCalled = true })

	handler, err := New(ctx, next, cfg, "traefik-sorrypage")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/some/path", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	if nextCalled {
		t.Error("next handler must not be called when sorrypage mode is enabled")
	}

	if len(received) != 1 {
		t.Fatalf("expected 1 proxied request, got %d", len(received))
	}

	assertResponseStatus(t, recorder, http.StatusServiceUnavailable)
	assertResponseHeader(t, recorder, "Content-Type", "text/html; charset=utf-8")
	assertResponseBody(t, recorder, sorryPageBody)
}

func TestSorryPageDisabledCallsNext(t *testing.T) {
	var received []*http.Request
	sorryServer := newSorryPageServer(t, &received)

	cfg := CreateConfig()
	cfg.Enabled = false
	cfg.RedirectService = sorryServer.URL

	ctx := context.Background()
	nextCalled := false
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		nextCalled = true
		rw.WriteHeader(http.StatusOK)
	})

	handler, err := New(ctx, next, cfg, "traefik-sorrypage")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	if !nextCalled {
		t.Error("next handler must be called when sorrypage mode is disabled")
	}

	if len(received) != 0 {
		t.Errorf("expected no proxied requests, got %d", len(received))
	}

	assertResponseStatus(t, recorder, http.StatusOK)
	assertEmptyContentTypeHeader(t, recorder)
	assertEmptyResponseBody(t, recorder)
}

func TestNewRejectsEmptyRedirectService(t *testing.T) {
	cfg := CreateConfig()
	cfg.Enabled = true

	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	if _, err := New(context.Background(), next, cfg, "traefik-sorrypage"); err == nil {
		t.Fatal("expected an error for an empty redirectService")
	}
}

func TestNewRejectsInvalidRedirectService(t *testing.T) {
	cfg := CreateConfig()
	cfg.Enabled = true
	cfg.RedirectService = "http://[::1"

	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	if _, err := New(context.Background(), next, cfg, "traefik-sorrypage"); err == nil {
		t.Fatal("expected an error for an invalid redirectService URL")
	}
}

func assertEmptyResponseBody(t *testing.T, recorder *httptest.ResponseRecorder) {
	t.Helper()

	responseBodyValue := recorder.Body.String()
	if responseBodyValue != "" {
		t.Errorf("unexpected response body value: %s", responseBodyValue)
	}
}

func assertEmptyContentTypeHeader(t *testing.T, recorder *httptest.ResponseRecorder) {
	t.Helper()

	contentTypeHeaderValue := recorder.Header().Get("Content-Type")
	if contentTypeHeaderValue != "" {
		t.Errorf("unexpected header value: %s", contentTypeHeaderValue)
	}
}

func assertResponseStatus(t *testing.T, resp *httptest.ResponseRecorder, expected int) {
	t.Helper()

	if resp.Code != expected {
		t.Errorf("invalid response status [%d] was expecting [%d]", resp.Code, expected)
	}
}

func assertResponseHeader(t *testing.T, resp *httptest.ResponseRecorder, key, expected string) {
	t.Helper()

	if resp.Header().Get(key) != expected {
		t.Errorf("invalid header value [%s] was expecting [%s]", resp.Header().Get(key), expected)
	}
}

func assertResponseBody(t *testing.T, resp *httptest.ResponseRecorder, expected string) {
	t.Helper()

	if resp.Body.String() != expected {
		t.Errorf("invalid response value [%s] was expecting [%s]", resp.Body.String(), expected)
	}
}
