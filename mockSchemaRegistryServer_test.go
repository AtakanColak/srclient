package srclient

import (
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestContentTypeHeader(t *testing.T) {
	r := mux.NewRouter()
	route := r.NewRoute().HeadersRegexp("Content-Type", "application/(vnd\\.schemaregistry\\.v1\\+json|vnd\\.schemaregistry\\+json|octet\\-stream|json)")

	types := []string{
		"application/vnd.schemaregistry.v1+json",
		"application/vnd.schemaregistry+json",
		"application/octet-stream",
		"application/json",
	}

	for _, ctype := range types {
		req := httptest.NewRequest("GET", "http://example.com", nil)
		req.Header.Set("Content-Type", ctype)

		if !route.Match(req, &mux.RouteMatch{}) {
			t.Fatalf("mismatch for %s", ctype)
		}
	}
}
