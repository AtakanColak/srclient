package srclient

import "github.com/linkedin/goavro/v2"

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

func (c Compatibility) String() string {
	return string(c)
}

// SchemaType of Schemas
type SchemaType string

const (
	SchemaTypeProtobuf SchemaType = "PROTOBUF"
	SchemaTypeAvro     SchemaType = "AVRO"
	SchemaTypeJson     SchemaType = "JSON"
)

func (s SchemaType) String() string {
	return string(s)
}

// Mode resource
type Mode string

const (
	ModeImport    Mode = "IMPORT"
	ModeReadOnly  Mode = "READONLY"
	ModeReadWrite Mode = "READWRITE"
)

func (m Mode) String() string {
	return string(m)
}

// Reference of Schemas use the import statement of Protobuf and
// the $ref field of JSON Schema. They are defined by the name
// of the import or $ref and the associated subject in the registry.
type Reference struct {
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Version int    `json:"version"`
}

// Schema is a data structure that holds all
// the relevant information about schemas.
type Schema struct {
	id      int
	schema  string
	version int
	codec   *goavro.Codec
}

// ID ensures access to ID
func (schema *Schema) ID() int {
	return schema.id
}

// Schema ensures access to Schema
func (schema *Schema) Schema() string {
	return schema.schema
}

// Version ensures access to Version
func (schema *Schema) Version() int {
	return schema.version
}

// Codec ensures access to Codec
// Will try to initialize a new one if it hasn't been initialized before
// Will return nil if it can't initialize a codec from the schema
func (schema *Schema) Codec() *goavro.Codec {
	if schema.codec == nil {
		codec, err := goavro.NewCodec(schema.Schema())
		if err == nil {
			schema.codec = codec
		}
	}
	return schema.codec
}

type errorResponse struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

type schemaRequest struct {
	Schema     string      `json:"schema"`
	SchemaType string      `json:"schemaType"`
	References []Reference `json:"references,omitempty"`
}

type schemaResponse struct {
	Subject string `json:"subject"`
	Version int    `json:"version"`
	Schema  string `json:"schema"`
	ID      int    `json:"id"`
}

type isCompatibleResponse struct {
	IsCompatible bool `json:"is_compatible"`
}
