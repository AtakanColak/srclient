package srclient

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"sync"

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
	responseContentType         = "application/vnd.schemaregistry.v1+json"
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

// MockSchemaRegistryServer is a Schema Registry implementation
// Use only for testing, it has bad performance
// Use NewMockSchemaRegistryServer() for initialization
type MockSchemaRegistryServer struct {
	Router      *mux.Router
	schemas     []mockSchemaRegistryServerSchema
	schemasLock sync.RWMutex
}

// NewMockSchemaRegistryServer constructor
func NewMockSchemaRegistryServer() *MockSchemaRegistryServer {
	server := &MockSchemaRegistryServer{
		Router:  mux.NewRouter(),
		schemas: make([]mockSchemaRegistryServerSchema, 0),
	}

	server.initializeRoutes()
	return server
}

// TestMockSchemaRegistryServer is a server with predefined simple data for testing
func TestMockSchemaRegistryServer() *MockSchemaRegistryServer {
	server := &MockSchemaRegistryServer{
		Router: mux.NewRouter(),
		schemas: []mockSchemaRegistryServerSchema{
			{
				ID:         1,
				Version:    1,
				Subject:    "test1",
				SchemaType: Avro,
				Schema:     "{\"type\":\"string\"}",
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
	server.initializeRoutes()
	return server
}

// ServeHTTP to implement http.Handler interface
func (m *MockSchemaRegistryServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Router.ServeHTTP(w, r)
}

func (m *MockSchemaRegistryServer) initializeRoutes() {
	m.Router.HandleFunc("/subjects", m.getSubjects).Methods("GET")
	m.Router.HandleFunc("/schemas/types", m.getSchemaTypes).Methods("GET")
	m.Router.HandleFunc("/schemas/ids/{id}", m.getSchemaWithID).Methods("GET")
	m.Router.HandleFunc("/schemas/ids/{id}/versions", m.getVersionByID).Methods("GET")
	m.Router.HandleFunc("/subjects/{subject}", m.checkIfSchemaExists).Methods("POST")
	m.Router.HandleFunc("/subjects/{subject}", m.deleteSubject).Methods("DELETE")
	m.Router.HandleFunc("/subjects/{subject}/versions", m.getVersions).Methods("GET")
	m.Router.HandleFunc("/subjects/{subject}/versions", m.createSchema).Methods("POST")
	m.Router.HandleFunc("/subjects/{subject}/versions/{version}", m.deleteVersion).Methods("DELETE")
	m.Router.HandleFunc("/subjects/{subject}/versions/{version}", m.getSchemaWithVersion).Methods("GET")
	m.Router.HandleFunc("/subjects/{subject}/versions/{version}/schema", m.getSchemaWithVersionUnescaped).Methods("GET")
	m.Router.HandleFunc("/compatibility/subjects/{subject}/versions/{version}", m.checkIfSchemaCompatible).Methods("POST")
	m.Router.HandleFunc("/mode", m.handleUnimplementedModeRequest)
	m.Router.HandleFunc("/config", m.getConfig).Methods("GET")
}

func (m *MockSchemaRegistryServer) getSchemas() []mockSchemaRegistryServerSchema {
	m.schemasLock.RLock()
	defer m.schemasLock.RUnlock()
	return m.schemas
}

func (m *MockSchemaRegistryServer) setSchemas(newSchemas []mockSchemaRegistryServerSchema) {
	m.schemasLock.Lock()
	defer m.schemasLock.Unlock()
	m.schemas = newSchemas
}

func (m *MockSchemaRegistryServer) getLatestVersionBySubject(subject string) (mockSchemaRegistryServerSchema, error) {
	latest := mockSchemaRegistryServerSchema{Version: -1}
	for _, schema := range m.getSchemas() {
		if schema.Subject == subject && schema.Version > latest.Version {
			latest = schema
		}
	}
	if latest.Version == -1 {
		return latest, errors.New("Schema not found")
	}
	return latest, nil
}

func (m *MockSchemaRegistryServer) getLatestVersionByID(id int) (mockSchemaRegistryServerSchema, error) {
	latest := mockSchemaRegistryServerSchema{Version: -1}
	for _, schema := range m.getSchemas() {
		if schema.ID == id && schema.Version > latest.Version {
			latest = schema
		}
	}
	if latest.Version == -1 {
		return latest, errors.New("Schema not found")
	}
	return latest, nil
}

func (m *MockSchemaRegistryServer) getSubjects(w http.ResponseWriter, r *http.Request) {
	subjectsMap := make(map[string]bool)
	for _, schema := range m.getSchemas() {
		subjectsMap[schema.Subject] = true
	}

	subjects := make([]string, 0)
	for subject := range subjectsMap {
		subjects = append(subjects, subject)
	}

	sort.Strings(subjects)

	respond(w, http.StatusOK, subjects)
}

func (m *MockSchemaRegistryServer) getVersions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subject, exists := vars["subject"]
	if !exists {
		respondWithError(w, http.StatusNotFound, 404, "HTTP 404 Not Found")
		return
	}

	versions := make([]int, 0)
	for _, schema := range m.getSchemas() {
		if schema.Subject == subject {
			versions = append(versions, schema.Version)
		}
	}

	respond(w, http.StatusOK, versions)
}

func (m *MockSchemaRegistryServer) getVersionByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, 400, "ID needs to be n integer")
		return
	}

	schema, err := m.getLatestVersionByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, 40403, "Schema not found")
		return
	}

	respond(w, http.StatusOK, []map[string]interface{}{map[string]interface{}{"subject": schema.Subject, "version": schema.Version}})
}

func (m *MockSchemaRegistryServer) getSchemaTypes(w http.ResponseWriter, r *http.Request) {
	types := make(map[SchemaType]bool)
	for _, schema := range m.getSchemas() {
		types[schema.SchemaType] = true
	}

	typesArr := make([]string, 0)
	for t := range types {
		typesArr = append(typesArr, string(t))
	}

	respond(w, http.StatusOK, typesArr)
}

func (m *MockSchemaRegistryServer) getSchemaWithID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, 400, "ID needs to be n integer")
		return
	}

	schema, err := m.getLatestVersionByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, 40403, "Schema not found")
		return
	}

	respond(w, http.StatusOK, map[string]string{"schema": schema.Schema})
}

func (m *MockSchemaRegistryServer) getSchemaWithVersion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subject, subjectExists := vars["subject"]
	version, versionExists := vars["version"]
	if !subjectExists || !versionExists {
		respondWithError(w, http.StatusNotFound, 404, "HTTP 404 Not Found")
		return
	}

	v := -1
	if version == "latest" {
		for _, schema := range m.getSchemas() {
			if schema.Subject == subject && schema.Version > v {
				v = schema.Version
			}
		}
		if v == -1 {
			respondWithError(w, http.StatusNotFound, 40402, "Version Not Found")
			return
		}
	} else {
		vNum, err := strconv.Atoi(version)
		if err != nil {
			respondWithError(w, http.StatusNotFound, 40402, "Version Not Found")
			return
		}
		v = vNum
	}

	var schemaToReturn mockSchemaRegistryServerSchema

	purified := make([]mockSchemaRegistryServerSchema, 0)
	for _, schema := range m.getSchemas() {
		if schema.Subject == subject && schema.Version == v {
			schemaToReturn = schema
			continue
		}

		purified = append(purified, schema)
	}
	m.setSchemas(purified)
	respond(w, http.StatusOK, schemaResponse{
		Subject: schemaToReturn.Subject,
		ID:      schemaToReturn.ID,
		Version: schemaToReturn.Version,
		Schema:  schemaToReturn.Schema,
	})
}

func (m *MockSchemaRegistryServer) getSchemaWithVersionUnescaped(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (m *MockSchemaRegistryServer) createSchema(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subject, exists := vars["subject"]
	if !exists {
		respondWithError(w, http.StatusNotFound, 404, "HTTP 404 Not Found")
		return
	}

	req := new(schemaRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		respondWithError(w, http.StatusUnprocessableEntity, 42201, "Invalid schema request")
		return
	}
	newSchema := mockSchemaRegistryServerSchema{
		Subject:    subject,
		Schema:     req.Schema,
		SchemaType: SchemaType(req.SchemaType),
	}

	for _, schema := range m.getSchemas() {
		if schema.ID >= newSchema.ID {
			newSchema.ID = schema.ID + 1
		}

		if schema.Subject == newSchema.Subject && schema.Version >= newSchema.Version {
			newSchema.Version = schema.Version + 1
		}
	}

	newSchemas := append(m.getSchemas(), newSchema)

	m.setSchemas(newSchemas)
	respond(w, http.StatusOK, schemaResponse{
		Subject: newSchema.Subject,
		ID:      newSchema.ID,
		Version: newSchema.Version,
		Schema:  newSchema.Schema,
	})
}

func (m *MockSchemaRegistryServer) checkIfSchemaExists(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subject, exists := vars["subject"]
	if !exists {
		respondWithError(w, http.StatusNotFound, 404, "HTTP 404 Not Found")
		return
	}

	sr := new(schemaRequest)
	if err := json.NewDecoder(r.Body).Decode(sr); err != nil {
		respondWithError(w, http.StatusInternalServerError, 500, "Internal server error")
		return
	}

	requested := make(map[string]interface{})
	if err := json.Unmarshal([]byte(sr.Schema), &requested); err != nil {
		respondWithError(w, http.StatusInternalServerError, 500, "Internal server error")
		return
	}

	for _, schema := range m.getSchemas() {
		if schema.Subject == subject && string(schema.SchemaType) == sr.SchemaType {
			existing := make(map[string]interface{})
			if err := json.Unmarshal([]byte(schema.Schema), &existing); err != nil {
				continue
			}

			if reflect.DeepEqual(requested, existing) {
				respond(w, http.StatusOK, schemaResponse{
					Subject: subject,
					ID:      schema.ID,
					Version: schema.Version,
					Schema:  schema.Schema,
				})
				return
			}
		}
	}

	respondWithError(w, http.StatusNotFound, 40403, "Schema not found")
}

func (m *MockSchemaRegistryServer) checkIfSchemaCompatible(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// deleteSubject doesn't differentiate between permanent or not
func (m *MockSchemaRegistryServer) deleteSubject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subject, exists := vars["subject"]
	if !exists {
		respondWithError(w, http.StatusNotFound, 404, "HTTP 404 Not Found")
		return
	}

	purified := make([]mockSchemaRegistryServerSchema, 0)
	deletedVersions := make([]int, 0)
	for _, schema := range m.getSchemas() {
		if schema.Subject != subject {
			purified = append(purified, schema)
		} else {
			deletedVersions = append(deletedVersions, schema.Version)
		}
	}
	m.setSchemas(purified)
	respond(w, http.StatusOK, deletedVersions)
}

func (m *MockSchemaRegistryServer) deleteVersion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subject, subjectExists := vars["subject"]
	version, versionExists := vars["version"]
	if !subjectExists || !versionExists {
		respondWithError(w, http.StatusNotFound, 404, "HTTP 404 Not Found")
		return
	}

	versions := make([]mockSchemaRegistryServerSchema, 0)
	for _, schema := range m.getSchemas() {
		if schema.Subject == subject {
			versions = append(versions, schema)
		}
	}

	if len(versions) == 0 {
		respondWithError(w, http.StatusNotFound, 40401, "Subject Not Found")
		return
	}

	purified := make([]mockSchemaRegistryServerSchema, 0)

	if version == "latest" {
		latest := versions[0]
		for _, schema := range versions {
			if schema.Version > latest.Version {
				latest = schema
			}
		}

		for _, schema := range m.getSchemas() {
			if schema.Subject == subject && schema.Version == latest.Version {
				continue
			}
			purified = append(purified, schema)
		}
		m.setSchemas(purified)
		respond(w, http.StatusOK, latest.Version)
		return
	}

	v, err := strconv.Atoi(version)
	if err != nil {
		respondWithError(w, http.StatusNotFound, 40402, "Version Not Found")
		return
	}

	for _, schema := range m.getSchemas() {
		if schema.Version != v {
			purified = append(purified, schema)
		}
	}
	m.setSchemas(purified)
	respond(w, http.StatusOK, v)
}

func (m *MockSchemaRegistryServer) handleUnimplementedModeRequest(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (m *MockSchemaRegistryServer) getConfig(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (m *MockSchemaRegistryServer) handleUnimplementedConfigRequest(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func respondWithError(w http.ResponseWriter, statusCode, errorCode int, message string) {
	respond(w, statusCode, errorResponse{ErrorCode: errorCode, Message: message})
}

func respond(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", responseContentType)
	w.WriteHeader(statusCode)
	w.Write(response)
}
