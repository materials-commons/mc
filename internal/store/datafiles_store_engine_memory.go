package store

type DatadirSchemaInMemory struct {
	DatadirSchema
	DatadirID string
}

type DatafilesStoreEngineMemory struct {
	DB map[string]DatadirSchemaInMemory
}

func NewDatafilesStoreEngineMemory() *DatafilesStoreEngineMemory {
	return &DatafilesStoreEngineMemory{
		DB: make(map[string]DatadirSchemaInMemory),
	}
}

func (e *DatafilesStoreEngineMemory) AddFile(file DatafileSchema) (DatafileSchema, error) {
	return DatafileSchema{}, nil
}

func (e *DatafilesStoreEngineMemory) GetFile(id string) (DatafileSchema, error) {
	return DatafileSchema{}, nil
}

func (e *DatafilesStoreEngineMemory) GetFileWithChecksum(checksum string) (DatafileSchema, error) {
	return DatafileSchema{}, nil
}

func (e *DatafilesStoreEngineMemory) GetFileInDir(name string, dirID string) (DatafileSchema, error) {
	return DatafileSchema{}, nil
}
