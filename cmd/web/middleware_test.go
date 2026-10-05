package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kyzercmd/snipio/internal/assert"
)

func TestCommonHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	commonHeaders(next).ServeHTTP(rr, r)

	result := rr.Result()

	expected := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, result.Header.Get("Content-Security-Policy"), expected)

	expected = "origin-when-cross-origin"
	assert.Equal(t, result.Header.Get("Referrer-Policy"), expected)

	expected = "nosniff"
	assert.Equal(t, result.Header.Get("X-Content-Type-Options"), expected)

	expected = "deny"
	assert.Equal(t, result.Header.Get("X-Frame-Options"), expected)

	expected = "0"
	assert.Equal(t, result.Header.Get("X-XSS-Protection"), expected)

	expected = "Go"
	assert.Equal(t, result.Header.Get("Server"), expected)

	assert.Equal(t, result.StatusCode, http.StatusOK)

	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
