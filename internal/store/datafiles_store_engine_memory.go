package store

import (
	"github.com/hashicorp/go-uuid"
	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/pkg/mc"
)

type DatafileSchemaInMemory struct {
	DataFile  model.DatafileSchema
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

func (e *DatafilesStoreEngineMemory) AddFile(file model.DatafileSchema, projectID, datadirID string) (model.DatafileSchema, error) {
	var err error
	if file.ID, err = uuid.GenerateUUID(); err != nil {
		return model.DatafileSchema{}, err
	}

	df := DatafileSchemaInMemory{
		DataFile:  file,
		DatadirID: datadirID,
	}

	e.DB[df.DataFile.ID] = df
	return df.DataFile, nil
}

func (e *DatafilesStoreEngineMemory) GetFile(id string) (model.DatafileSchema, error) {
	dfEntry, ok := e.DB[id]
	if !ok {
		return model.DatafileSchema{}, mc.ErrNotFound
	}

	return dfEntry.DataFile, nil
}

func (e *DatafilesStoreEngineMemory) GetFileWithChecksum(checksum string) (model.DatafileSchema, error) {
	for _, dfEntry := range e.DB {
		if dfEntry.DataFile.Checksum == checksum {
			return dfEntry.DataFile, nil
		}
	}

	return model.DatafileSchema{}, mc.ErrNotFound
}

func (e *DatafilesStoreEngineMemory) GetFileInDir(name string, dirID string) (model.DatafileSchema, error) {
	for _, dfEntry := range e.DB {
		if dfEntry.DatadirID == dirID && dfEntry.DataFile.Name == name && dfEntry.DataFile.Current {
			return dfEntry.DataFile, nil
		}
	}

	return model.DatafileSchema{}, mc.ErrNotFound
}

func (e *DatafilesStoreEngineMemory) UpdateFileCurrentFlag(id string, current bool) error {
	dfEntry, ok := e.DB[id]
	if !ok {
		return mc.ErrNotFound
	}

	dfEntry.DataFile.Current = current
	e.DB[id] = dfEntry
	return nil
}
