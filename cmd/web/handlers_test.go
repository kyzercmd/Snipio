package main

import (
	"net/http"
	"testing"

	"github.com/kyzercmd/snipio/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	status, _, body := ts.get(t, "/ping")

	assert.Equal(t, status, http.StatusOK)
	assert.Equal(t, body, "OK")
}

func TestSnipperView(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "The sky is beautiful",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/2",
			wantCode: http.StatusNotFound,
			wantBody: "",
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-1",
			wantCode: http.StatusNotFound,
			wantBody: "",
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.2",
			wantCode: http.StatusNotFound,
			wantBody: "",
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/test",
			wantCode: http.StatusNotFound,
			wantBody: "",
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
			wantBody: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, _, body := ts.get(t, tt.urlPath)

			assert.Equal(t, statusCode, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}
