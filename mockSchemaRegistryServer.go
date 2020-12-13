package srclient

type mockSchemaRegistryServerSchema struct {
	ID         int
	Version    int
	Subject    string
	SchemaType SchemaType
	Schema     string
}

type mockSchemaRegistryServer struct {
	schemas []mockSchemaRegistryServerSchema
}

type errorResponse struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

// MockSchemaRegistryServer returns an unstarted httptest.Server
// that satisfies Schema Registry 5.5.0 API for testing
// refer to https://docs.confluent.io/5.5.0/schema-registry/schema-validation.html
// Compatibility is
// func MockSchemaRegistryServer() *httptest.Server {

// 	s := mockSchemaRegistryServer{
// 		schemas: make([]mockSchemaRegistryServerSchema, 0),
// 	}
// 	r := mux.NewRouter()
// 	r.Headers("Content-Type", "application/(vnd\\.schemaregistry\\.v1\\+json|vnd\\.schemaregistry\\+json|octet\\-stream|json)")
// 	// handler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request {

// 	// })
// 	_ =
// 	return nil
// }
