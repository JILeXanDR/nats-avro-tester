package avro

type CodecStorage interface {
	GetAll() ([]*CodecWrapper, error)
	FindByNamespace(subject string) (*CodecWrapper, error)
	SyncSchemas(schemas ...string) error
}
