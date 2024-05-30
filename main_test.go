package traefik_sorrypage

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSorryPagePage(t *testing.T) {
	sorrypagePage := "sorrypage_service"

	cfg := CreateConfig()
	cfg.Enabled = true
	cfg.RedirectService = sorrypagePage

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

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

	assertResponseStatus(t, recorder, http.StatusServiceUnavailable)
	assertResponseHeader(t, recorder, "Content-Type", "text/html; charset=utf-8")
	assertResponseBody(t, recorder, "<html><head></head><body>SorryPage</body></html>")
}

func TestSorryPagePageWithOtherStatusCodeAndContentType(t *testing.T) {
	sorrypagePage := "sorrypage_service"

	cfg := CreateConfig()
	cfg.Enabled = true
	cfg.RedirectService = sorrypagePage

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

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

	assertResponseStatus(t, recorder, http.StatusTeapot)
	assertResponseHeader(t, recorder, "Content-Type", "application/json; charset=utf-8")
	assertResponseBody(t, recorder, "{\"detail\": \"This endpoint is currently in sorrypage mode\"}")
}

func TestSorryPagePageWithoutTrigger(t *testing.T) {
	cfg := CreateConfig()
	cfg.Enabled = true
	sorrypagePage := "sorrypage_service"
	cfg.RedirectService = sorrypagePage

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

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

	assertEmptyContentTypeHeader(t, recorder)
	assertEmptyResponseBody(t, recorder)
}

func TestSorryPagePageWithMissingTrigger(t *testing.T) {
	cfg := CreateConfig()
	cfg.Enabled = true
	sorrypagePage := "sorrypage_service"
	cfg.RedirectService = sorrypagePage

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

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

	assertResponseHeader(t, recorder, "Content-Type", "text/html; charset=utf-8")
	assertResponseBody(t, recorder, "<html><head></head><body>SorryPage</body></html>")
}

func TestDisabledSorryPagePage(t *testing.T) {

	cfg := CreateConfig()
	cfg.Enabled = false
	sorrypagePage := "sorrypage_service"
	cfg.RedirectService = sorrypagePage

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

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

	assertEmptyContentTypeHeader(t, recorder)
	assertEmptyResponseBody(t, recorder)
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
		t.Errorf("invalid resonse status [%d] was expecting [%d]", resp.Code, expected)
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
