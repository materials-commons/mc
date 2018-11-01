package storengine

import (
	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/pkg/mc"
)

type GlobusUploadsMemory struct {
	DB map[string]model.GlobusUploadSchema
}

func NewGlobusUploadsMemory() *GlobusUploadsMemory {
	return &GlobusUploadsMemory{
		DB: make(map[string]model.GlobusUploadSchema),
	}
}

func NewGlobusUploadsMemoryWithDB(db map[string]model.GlobusUploadSchema) *GlobusUploadsMemory {
	return &GlobusUploadsMemory{DB: db}
}

func (e *GlobusUploadsMemory) AddGlobusUpload(upload model.GlobusUploadSchema) (model.GlobusUploadSchema, error) {
	e.DB[upload.ID] = upload
	return upload, nil
}

func (e *GlobusUploadsMemory) DeleteGlobusUpload(id string) error {
	_, ok := e.DB[id]
	if !ok {
		return mc.ErrNotFound
	}
	delete(e.DB, id)
	return nil
}

func (e *GlobusUploadsMemory) GetGlobusUpload(id string) (model.GlobusUploadSchema, error) {
	upload, ok := e.DB[id]
	if !ok {
		return model.GlobusUploadSchema{}, mc.ErrNotFound
	}

	return upload, nil
}
