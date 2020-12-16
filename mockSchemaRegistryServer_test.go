package srclient_test

import (
	"io"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/riferrei/srclient"
	"github.com/stretchr/testify/assert"
)

var dummy = srclient.TestMockSchemaRegistryServer()

func doRequest(t testing.TB, method, url string, body io.Reader, expected []byte) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatal(r)
		}
	}()

	req := httptest.NewRequest(method, url, body)
	recorder := httptest.NewRecorder()
	dummy.Router.ServeHTTP(recorder, req)
	resp := recorder.Result()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err.Error())
	}
	assert.Equal(t, expected, result)
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
	doRequest(t, "GET", "http://example.com/subjects", nil, []byte(`["test1","test2"]`))
}

func TestGetVersions(t *testing.T) {
	doRequest(t, "GET", "http://example.com/subjects/test1/versions", nil, []byte(`[1]`))
}

func TestGetSchemaTypes(t *testing.T) {
	doRequest(t, "GET", "http://example.com/schemas/types", nil, []byte(`["AVRO"]`))
}

func TestGetSchemaWithID(t *testing.T) {
	doRequest(t, "GET", "http://example.com/schemas/ids/1", nil, []byte(`{"schema":"{\"type\":\"string\"}"}`))
}

func TestGetVersionByID(t *testing.T) {
	doRequest(t, "GET", "http://example.com/schemas/ids/1/versions", nil, []byte(`[{"subject":"test1","version":1}]`))
}

func TestCheckIfSchemaExists(t *testing.T) {
	doRequest(t, "POST", "http://example.com/subjects/test1", nil, []byte(``))
}
