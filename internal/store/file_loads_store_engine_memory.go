package store

import "github.com/hashicorp/go-uuid"

type FileLoadsStoreEngineMemory struct {
	DB map[string]FileLoadSchema
}

func NewFileLoadsStoreEngineMemory() *FileLoadsStoreEngineMemory {
	return &FileLoadsStoreEngineMemory{
		DB: make(map[string]FileLoadSchema),
	}
}

func (e *FileLoadsStoreEngineMemory) AddFileLoad(fload FileLoadSchema) (FileLoadSchema, error) {
	var err error
	if fload.ID, err = uuid.GenerateUUID(); err != nil {
		return FileLoadSchema{}, err
	}
	e.DB[fload.ID] = fload
	return fload, nil
}

func (e *FileLoadsStoreEngineMemory) DeleteFileLoad(floadID string) error {
	_, ok := e.DB[floadID]
	if !ok {
		return ErrNotFound
	}

	delete(e.DB, floadID)
	return nil
}

func (e *FileLoadsStoreEngineMemory) GetFileLoad(floadID string) (FileLoadSchema, error) {
	fload, ok := e.DB[floadID]
	if !ok {
		return FileLoadSchema{}, ErrNotFound
	}

	return fload, nil
}

func (e *FileLoadsStoreEngineMemory) GetAllFileLoads() ([]FileLoadSchema, error) {
	var fileLoads []FileLoadSchema
	for _, entry := range e.DB {
		fileLoads = append(fileLoads, entry)
	}

	return fileLoads, nil
}
