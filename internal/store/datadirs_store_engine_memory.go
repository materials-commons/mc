package store

import "github.com/hashicorp/go-uuid"

type DatadirsStoreEngineMemory struct {
	DB map[string]DatadirSchema
}

func NewDatadirsStoreEngineMemory() *DatadirsStoreEngineMemory {
	return &DatadirsStoreEngineMemory{
		DB: make(map[string]DatadirSchema),
	}
}

func NewDatadirsStoreEngineMemoryWithDB(db map[string]DatadirSchema) *DatadirsStoreEngineMemory {
	return &DatadirsStoreEngineMemory{
		DB: db,
	}
}

func (e *DatadirsStoreEngineMemory) AddDir(dir DatadirSchema) (DatadirSchema, error) {
	var err error
	if dir.ID, err = uuid.GenerateUUID(); err != nil {
		return DatadirSchema{}, err
	}
	e.DB[dir.ID] = dir

	return dir, nil
}

func (e *DatadirsStoreEngineMemory) GetDatadirByPathInProject(path, projectID string) (DatadirSchema, error) {
	for _, ddir := range e.DB {
		if ddir.Name == path && ddir.Project == projectID {
			return ddir, nil
		}
	}

	return DatadirSchema{}, ErrNotFound
}

func (e *DatadirsStoreEngineMemory) GetDatadir(id string) (DatadirSchema, error) {
	ddir, ok := e.DB[id]
	if !ok {
		return DatadirSchema{}, ErrNotFound
	}

	return ddir, nil
}
