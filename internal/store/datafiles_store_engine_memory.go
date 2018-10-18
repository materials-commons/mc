package store

import "github.com/hashicorp/go-uuid"

type DatafileSchemaInMemory struct {
	DataFile  DatafileSchema
	DatadirID string
}

type DatafilesStoreEngineMemory struct {
	DB map[string]DatafileSchemaInMemory
}

func NewDatafilesStoreEngineMemory() *DatafilesStoreEngineMemory {
	return &DatafilesStoreEngineMemory{
		DB: make(map[string]DatafileSchemaInMemory),
	}
}

func NewDatafilesStoreEngineMemoryWithDB(db map[string]DatafileSchemaInMemory) *DatafilesStoreEngineMemory {
	return &DatafilesStoreEngineMemory{
		DB: db,
	}
}

func (e *DatafilesStoreEngineMemory) AddFile(file DatafileSchema, projectID, datadirID string) (DatafileSchema, error) {
	var err error
	if file.ID, err = uuid.GenerateUUID(); err != nil {
		return DatafileSchema{}, err
	}

	df := DatafileSchemaInMemory{
		DataFile:  file,
		DatadirID: datadirID,
	}

	e.DB[df.DataFile.ID] = df
	return df.DataFile, nil
}

func (e *DatafilesStoreEngineMemory) GetFile(id string) (DatafileSchema, error) {
	dfEntry, ok := e.DB[id]
	if !ok {
		return DatafileSchema{}, ErrNotFound
	}

	return dfEntry.DataFile, nil
}

func (e *DatafilesStoreEngineMemory) GetFileWithChecksum(checksum string) (DatafileSchema, error) {
	for _, dfEntry := range e.DB {
		if dfEntry.DataFile.Checksum == checksum {
			return dfEntry.DataFile, nil
		}
	}

	return DatafileSchema{}, ErrNotFound
}

func (e *DatafilesStoreEngineMemory) GetFileInDir(name string, dirID string) (DatafileSchema, error) {
	for _, dfEntry := range e.DB {
		if dfEntry.DatadirID == dirID && dfEntry.DataFile.Name == name && dfEntry.DataFile.Current {
			return dfEntry.DataFile, nil
		}
	}

	return DatafileSchema{}, ErrNotFound
}

func (e *DatafilesStoreEngineMemory) UpdateFileCurrentFlag(id string, current bool) error {
	dfEntry, ok := e.DB[id]
	if !ok {
		return ErrNotFound
	}

	dfEntry.DataFile.Current = current
	e.DB[id] = dfEntry
	return nil
}
