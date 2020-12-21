package srclient

import "context"

// ISchemaRegistryClient provides the definition of the operations that the Schema Registry client provides.
type ISchemaRegistryClient interface {
	SetCachingEnabled(enabled bool)
	SetCodecCreationEnabled(enabled bool)
	SetBasicAuthCredentials(username, password string)

	GetSubjects() ([]string, error)
	GetSchemaTypes() ([]SchemaType, error)

	GetLatestSchema(subject string) (*Schema, error)

	GetSchemaByID(id int) (*Schema, error)
	GetSchemaBySubject(subject string, version int) (*Schema, error)

	GetSchemaVersionsByID(id int) ([]int, error)
	GetSchemaVersionsBySubject(subject string) ([]int, error)

	DeleteSubject(subject string, hardDelete bool) error
	DeleteVersion(subject string, version int, hardDelete bool) error

	CreateSchema(subject, schema string, schemaType SchemaType, references ...Reference) (*Schema, error)

	CheckIfSchemaExists(subject, schema string, schemaType SchemaType, references ...Reference) (bool, error)

	CheckIfSchemaCompatible(subject, schema string, version int, schemaType SchemaType, references ...Reference) (bool, error)

	GetConfig() (Compatibility, error)
	GetConfigOfSubject(subject string) (Compatibility, error)

	GetMode() (Mode, error)
	GetModeOfSubject(subject string) (Mode, error)

	//===================
	//WithContext methods
	//===================

	GetSubjectsWithContext(ctx context.Context) ([]string, error)
	GetSchemaTypesWithContext(ctx context.Context) ([]SchemaType, error)

	GetLatestSchemaWithContext(subject string, ctx context.Context) (*Schema, error)

	GetSchemaByIDWithContext(id int, ctx context.Context) (*Schema, error)
	GetSchemaBySubjectWithContext(subject string, ctx context.Context) (*Schema, error)

	GetSchemaVersionsByIDWithContext(id int, ctx context.Context) ([]int, error)
	GetSchemaVersionsBySubjectWithContext(subject string, ctx context.Context) ([]int, error)

	DeleteSubjectWithContext(subject string, hardDelete bool, ctx context.Context) error
	DeleteVersionWithContext(subject string, version int, hardDelete bool, ctx context.Context) error

	CreateSchemaWithContext(subject, schema string, schemaType SchemaType, ctx context.Context, references ...Reference) (*Schema, error)

	CheckIfSchemaExistsWithContext(subject, schema string, schemaType SchemaType, ctx context.Context, references ...Reference) (bool, error)

	CheckIfSchemaCompatibleWithContext(subject, schema string, version int, schemaType SchemaType, ctx context.Context, references ...Reference) (bool, error)

	GetConfigWithContext(ctx context.Context) (Compatibility, error)
	GetConfigOfSubjectWithContext(subject string, ctx context.Context) (Compatibility, error)

	GetModeWithContext(ctx context.Context) (Mode, error)
	GetModeOfSubjectWithContext(subject string, ctx context.Context) (Mode, error)
}

var _ ISchemaRegistryClient = &SchemaRegistryClient{}
var _ ISchemaRegistryClient = MockSchemaRegistryClient{}
