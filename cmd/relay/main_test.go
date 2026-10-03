// Readiness checks prove DB failure is safe and business paths remain unavailable.
package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReadinessAndUnavailableBusinessPaths(t *testing.T) {
	for _, test := range []struct {
		method, path string
		pingError    error
		status       int
	}{
		{http.MethodGet, "/healthz", nil, http.StatusOK},
		{http.MethodGet, "/healthz", errors.New("example-secret"), http.StatusServiceUnavailable},
		{http.MethodPost, "/healthz", nil, http.StatusMethodNotAllowed},
		{http.MethodPost, "/ingest", nil, http.StatusNotFound},
		{http.MethodPost, "/approve", nil, http.StatusNotFound},
		{http.MethodPost, "/exec", nil, http.StatusNotFound},
	} {
		response := httptest.NewRecorder()
		healthHandler(func(context.Context) error { return test.pingError }).ServeHTTP(response, httptest.NewRequest(test.method, test.path, nil))
		if response.Code != test.status || strings.Contains(response.Body.String(), "example-secret") {
			t.Fatalf("%s %s: unexpected response %d %s", test.method, test.path, response.Code, response.Body.String())
		}
	}
}
