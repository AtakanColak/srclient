package srclient

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Compatibility enforces compatibility rules when newer schemas are registered in the same subject
// See https://docs.confluent.io/5.5.0/schema-registry/develop/api.html#compatibility
type Compatibility string

const (
	CompatibilityBackward           Compatibility = "BACKWARD"
	CompatibilityBackwardTransitive Compatibility = "BACKWARD_TRANSITIVE"
	CompatibilityForward            Compatibility = "FORWARD"
	CompatibilityForwardTransitive  Compatibility = "FORWARD_TRANSITIVE"
	CompatibilityFull               Compatibility = "FULL"
	CompatibilityFullTransitive     Compatibility = "FULL_TRANSITIVE"
	CompatibilityNone               Compatibility = "NONE"
)

const (
	acceptableContentTypesRegex = "application/(vnd\\.schemaregistry\\.v1\\+json|vnd\\.schemaregistry\\+json|octet\\-stream|json)"
)

type errorResponse struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

type mockSchemaRegistryServerSchema struct {
	ID         int
	Version    int
	Subject    string
	SchemaType SchemaType
	Schema     string
}

// MockSchemaRegistryServer is a Schema Registry implementation for testing
// Use NewMockSchemaRegistryServer() for initialization
type MockSchemaRegistryServer struct {
	router  *mux.Router
	schemas []mockSchemaRegistryServerSchema
}

// NewMockSchemaRegistryServer constructor
func NewMockSchemaRegistryServer() *MockSchemaRegistryServer {
	return &MockSchemaRegistryServer{
		router:  mux.NewRouter(),
		schemas: make([]mockSchemaRegistryServerSchema, 0),
	}
}

func (m *MockSchemaRegistryServer) initializeRoutes() {
	m.router.HandleFunc("/subjects", m.getSubjects).Methods("GET")
	m.router.HandleFunc("/subjects/{subject}/versions", m.getVersions).Methods("GET")
}

func (m *MockSchemaRegistryServer) getSubjects(w http.ResponseWriter, r *http.Request) {
	subjectsMap := make(map[string]bool)
	for _, schema := range m.schemas {
		subjectsMap[schema.Subject] = true
	}

	subjects := make([]string, 0)
	for subject := range subjectsMap {
		subjects = append(subjects, subject)
	}

	respond(w, http.StatusOK, subjects)
}

func (m *MockSchemaRegistryServer) getVersions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subject, exists := vars["subject"]
	if !exists {
		log.Println(vars)
		respondWithError(w, http.StatusNotFound, 404, "HTTP 404 Not Found")
	}

	versions := make([]int, 0)
	for _, schema := range m.schemas {
		if schema.Subject == subject {
			versions = append(versions, schema.Version)
		}
	}

	respond(w, http.StatusOK, versions)
}

func respondWithError(w http.ResponseWriter, statusCode, errorCode int, message string) {
	respond(w, statusCode, errorResponse{ErrorCode: errorCode, Message: message})
}

func respond(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/vnd.schemaregistry.v1+json")
	w.WriteHeader(statusCode)
	w.Write(response)
}
