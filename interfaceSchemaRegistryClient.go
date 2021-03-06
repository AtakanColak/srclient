package srclient

// ISchemaRegistryClient provides the
// definition of the operations that
// this Schema Registry client provides.
type ISchemaRegistryClient interface {
	GetSubjects() ([]string, error)
	GetLatestSchema(subject string, isKey bool) (*Schema, error)
	GetSchemaVersions(subject string, isKey bool) ([]int, error)

	GetSchemaByID(schemaID int) (*Schema, error)
	GetSchemaBySubject(subject string, isKey bool) (*Schema, error)
	GetSchemaByVersion(subject string, version string, isKey bool) (*Schema, error)

	CreateSchema(subject string, schema string, schemaType SchemaType, isKey bool, references ...Reference) (*Schema, error)
	DeleteSubject(subject string, permanent bool) error

	SetCachingEnabled(value bool)
	SetCodecCreationEnabled(value bool)

	IsSchemaCompatible(subject, schema, version string, schemaType SchemaType, isKey bool) (bool, error)
}

// ensure interface is implemented
var _ ISchemaRegistryClient = &SchemaRegistryClient{}
var _ ISchemaRegistryClient = MockSchemaRegistryClient{}
