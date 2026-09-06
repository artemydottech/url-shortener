package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func post(t *testing.T, body string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, "http://short.test/api/shorten", strings.NewReader(body))
	recorder := httptest.NewRecorder()
	shortenHandler(recorder, request)
	return recorder
}

func TestGetPort(t *testing.T) {
	cases := map[string]string{
		"":      ":8080",
		"8080":  ":8080",
		":2020": ":2020",
	}

	for value, want := range cases {
		if value == "" {
			t.Setenv("PORT", "")
		} else {
			t.Setenv("PORT", value)
		}
		if got := getPort(); got != want {
			t.Errorf("PORT=%q: got %q, want %q", value, got, want)
		}
	}
}

func TestShortenReturnsCodeAndLinkOnTheRequestHost(t *testing.T) {
	storage = NewStorage()

	recorder := post(t, `{"url":"https://example.com/page"}`)
	if recorder.Code != http.StatusOK {
		t.Fatalf("got %d, want 200", recorder.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("bad json: %v", err)
	}

	code := response["short_code"]
	if len(code) != codeLength {
		t.Errorf("short_code %q is %d chars, want %d", code, len(code), codeLength)
	}
	if want := "http://short.test/" + code; response["short_url"] != want {
		t.Errorf("short_url = %q, want %q", response["short_url"], want)
	}
	if stored, ok := storage.Get(code); !ok || stored != "https://example.com/page" {
		t.Errorf("storage holds %q (found=%v)", stored, ok)
	}
}

func TestShortenRejectsBadInput(t *testing.T) {
	storage = NewStorage()

	if got := post(t, `not json`).Code; got != http.StatusBadRequest {
		t.Errorf("broken json: got %d, want 400", got)
	}
	if got := post(t, `{"url":""}`).Code; got != http.StatusBadRequest {
		t.Errorf("empty url: got %d, want 400", got)
	}
}

func TestShortenAnswersPreflightAndRejectsOtherMethods(t *testing.T) {
	preflight := httptest.NewRequest(http.MethodOptions, "http://short.test/api/shorten", nil)
	recorder := httptest.NewRecorder()
	shortenHandler(recorder, preflight)

	if recorder.Code != http.StatusNoContent {
		t.Errorf("preflight: got %d, want 204", recorder.Code)
	}
	if origin := recorder.Header().Get("Access-Control-Allow-Origin"); origin != "*" {
		t.Errorf("preflight origin header = %q", origin)
	}

	get := httptest.NewRequest(http.MethodGet, "http://short.test/api/shorten", nil)
	recorder = httptest.NewRecorder()
	shortenHandler(recorder, get)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("GET: got %d, want 405", recorder.Code)
	}
	if origin := recorder.Header().Get("Access-Control-Allow-Origin"); origin != "*" {
		t.Errorf("405 origin header = %q, want it set so the browser can read the error", origin)
	}
}

func TestRedirect(t *testing.T) {
	storage = NewStorage()
	storage.Save("abc123", "https://example.com/target")

	request := httptest.NewRequest(http.MethodGet, "http://short.test/abc123", nil)
	recorder := httptest.NewRecorder()
	redirectHandler(recorder, request)

	if recorder.Code != http.StatusFound {
		t.Errorf("got %d, want 302", recorder.Code)
	}
	if location := recorder.Header().Get("Location"); location != "https://example.com/target" {
		t.Errorf("Location = %q", location)
	}
}

func TestRedirectMisses(t *testing.T) {
	storage = NewStorage()

	for path, want := range map[string]int{
		"http://short.test/":       http.StatusBadRequest,
		"http://short.test/nope42": http.StatusNotFound,
	} {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		recorder := httptest.NewRecorder()
		redirectHandler(recorder, request)

		if recorder.Code != want {
			t.Errorf("%s: got %d, want %d", path, recorder.Code, want)
		}
	}
}

func TestBaseURLFollowsTheForwardedScheme(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "http://short.test/", nil)
	request.Header.Set("X-Forwarded-Proto", "https")

	if got := baseURL(request); got != "https://short.test" {
		t.Errorf("got %q, want https://short.test", got)
	}
}
