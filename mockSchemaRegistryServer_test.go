package srclient

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

var dummy *MockSchemaRegistryServer

func init() {
	dummy = &MockSchemaRegistryServer{
		router: mux.NewRouter(),
		schemas: []mockSchemaRegistryServerSchema{
			{
				ID:         1,
				Version:    1,
				Subject:    "test1",
				SchemaType: Avro,
				Schema:     "\"string\"",
			},
			{
				ID:         2,
				Version:    1,
				Subject:    "test2",
				SchemaType: Avro,
				Schema:     "\"int\"",
			},
		},
	}
	dummy.initializeRoutes()
}

func doRequest(t testing.TB, req *http.Request, handlerFunc http.HandlerFunc) string {
	recorder := httptest.NewRecorder()
	handlerFunc(recorder, req)
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

func TestGetSubjects(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/subjects", nil)
	doRequest(t, req, dummy.getSubjects)
}

func TestGetVersions(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/subjects/test1/versions", nil)
	doRequest(t, req, dummy.getVersions)
}
