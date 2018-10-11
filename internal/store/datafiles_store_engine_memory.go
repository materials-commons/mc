package store

import "github.com/hashicorp/go-uuid"

type DatafileSchemaInMemory struct {
	dataFile  DatafileSchema
	datadirID string
}

type DatafilesStoreEngineMemory struct {
	DB map[string]DatafileSchemaInMemory
}

func NewDatafilesStoreEngineMemory() *DatafilesStoreEngineMemory {
	return &DatafilesStoreEngineMemory{
		DB: make(map[string]DatafileSchemaInMemory),
	}
}

func (e *DatafilesStoreEngineMemory) AddFile(file DatafileSchema, projectID, datadirID string) (DatafileSchema, error) {
	var err error
	if file.ID, err = uuid.GenerateUUID(); err != nil {
		return DatafileSchema{}, err
	}

	df := DatafileSchemaInMemory{
		dataFile:  file,
		datadirID: datadirID,
	}

	e.DB[df.dataFile.ID] = df
	return df.dataFile, nil
}

func (e *DatafilesStoreEngineMemory) GetFile(id string) (DatafileSchema, error) {
	dfEntry, ok := e.DB[id]
	if !ok {
		return DatafileSchema{}, ErrNotFound
	}

	return dfEntry.dataFile, nil
}

func (e *DatafilesStoreEngineMemory) GetFileWithChecksum(checksum string) (DatafileSchema, error) {
	for _, dfEntry := range e.DB {
		if dfEntry.dataFile.Checksum == checksum {
			return dfEntry.dataFile, nil
		}
	}

	return DatafileSchema{}, ErrNotFound
}

func (e *DatafilesStoreEngineMemory) GetFileInDir(name string, dirID string) (DatafileSchema, error) {
	for _, dfEntry := range e.DB {
		if dfEntry.datadirID == dirID && dfEntry.dataFile.Name == name && dfEntry.dataFile.Current {
			return dfEntry.dataFile, nil
		}
	}

	return DatafileSchema{}, ErrNotFound
}

func (e *DatafilesStoreEngineMemory) UpdateFileCurrentFlag(id string, current bool) error {
	dfEntry, ok := e.DB[id]
	if !ok {
		return ErrNotFound
	}

	dfEntry.dataFile.Current = current
	e.DB[id] = dfEntry
	return nil
}
