package srclient_test

import (
	"io"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/riferrei/srclient"
)

var dummy = srclient.TestMockSchemaRegistryServer()

func doRequest(t testing.TB, method, url string, body io.Reader) string {
	req := httptest.NewRequest(method, url, body)
	recorder := httptest.NewRecorder()
	dummy.Router.ServeHTTP(recorder, req)
	resp := recorder.Result()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err.Error())
	} else {
		t.Log(string(result))
	}
	return string(result)
}

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

func TestMockSchemaRegistryServer(t *testing.T) {

	t.Run("GetSubjects", func(t *testing.T) {
		doRequest(t, "GET", "http://example.com/subjects", nil)
	})

	t.Run("GetVersions", func(t *testing.T) {
		doRequest(t, "GET", "http://example.com/subjects/test1/versions", nil)
	})

	t.Run("GetSchemaTypes", func(t *testing.T) {
		doRequest(t, "GET", "http://example.com/schemas/types", nil)
	})

	t.Run("GetSchemaWithID", func(t *testing.T) {
		doRequest(t, "GET", "http://example.com/schemas/ids/1", nil)
	})

	t.Run("GetVersionByID", func(t *testing.T) {
		doRequest(t, "GET", "http://example.com/schemas/ids/1/versions", nil)
	})
}
