package storengine

import (
	"github.com/hashicorp/go-uuid"
	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/pkg/mc"
)

type FileLoadsMemory struct {
	DB map[string]model.FileLoadSchema
}

func NewFileLoadsMemory() *FileLoadsMemory {
	return &FileLoadsMemory{
		DB: make(map[string]model.FileLoadSchema),
	}
}

func NewFileLoadsMemoryWithDB(db map[string]model.FileLoadSchema) *FileLoadsMemory {
	return &FileLoadsMemory{
		DB: db,
	}
}

func (e *FileLoadsMemory) AddFileLoad(fload model.FileLoadSchema) (model.FileLoadSchema, error) {
	var err error
	if fload.ID, err = uuid.GenerateUUID(); err != nil {
		return model.FileLoadSchema{}, err
	}
	e.DB[fload.ID] = fload
	return fload, nil
}

func (e *FileLoadsMemory) DeleteFileLoad(floadID string) error {
	_, ok := e.DB[floadID]
	if !ok {
		return mc.ErrNotFound
	}

	delete(e.DB, floadID)
	return nil
}

func (e *FileLoadsMemory) GetFileLoad(floadID string) (model.FileLoadSchema, error) {
	fload, ok := e.DB[floadID]
	if !ok {
		return model.FileLoadSchema{}, mc.ErrNotFound
	}

	return fload, nil
}

func (e *FileLoadsMemory) GetAllFileLoads() ([]model.FileLoadSchema, error) {
	var fileLoads []model.FileLoadSchema
	for _, entry := range e.DB {
		fileLoads = append(fileLoads, entry)
	}

	return fileLoads, nil
}

func (e *FileLoadsMemory) MarkAllNotLoading() error {
	for _, entry := range e.DB {
		entry.Loading = false
		e.DB[entry.ID] = entry
	}

	return nil
}

func (e *FileLoadsMemory) UpdateLoading(id string, loading bool) error {
	entry, ok := e.DB[id]
	if !ok {
		return mc.ErrNotFound
	}

	entry.Loading = loading
	e.DB[entry.ID] = entry
	return nil
}
