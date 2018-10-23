package store

import (
	"time"

	"github.com/materials-commons/mc/internal/store/storengine"

	"github.com/materials-commons/mc/internal/store/model"
)

type FileLoadsStore struct {
	flStoreEngine storengine.FileLoadsStoreEngine
}

func NewFileLoadsStore(e storengine.FileLoadsStoreEngine) *FileLoadsStore {
	return &FileLoadsStore{flStoreEngine: e}
}

func (s *FileLoadsStore) AddFileLoad(flModel model.AddFileLoadModel) (model.FileLoadSchema, error) {
	if err := flModel.Validate(); err != nil {
		return model.FileLoadSchema{}, err
	}

	now := time.Now()
	fl := model.FileLoadSchema{
		ModelSimple: model.ModelSimple{
			Birthtime: now,
			MTime:     now,
			OType:     "file_load",
		},
		ProjectID: flModel.ProjectID,
		Path:      flModel.Path,
		Owner:     flModel.Owner,
		Exclude:   flModel.Exclude,
	}

	return s.flStoreEngine.AddFileLoad(fl)
}

func (s *FileLoadsStore) DeleteFileLoad(id string) error {
	return s.flStoreEngine.DeleteFileLoad(id)
}

func (s *FileLoadsStore) GetFileLoad(id string) (model.FileLoadSchema, error) {
	return s.flStoreEngine.GetFileLoad(id)
}

func (s *FileLoadsStore) GetAllFileLoads() ([]model.FileLoadSchema, error) {
	return s.flStoreEngine.GetAllFileLoads()
}

func (s *FileLoadsStore) MarkAllNotLoading() error {
	return s.flStoreEngine.MarkAllNotLoading()
}

func (s *FileLoadsStore) UpdateLoading(id string, loading bool) error {
	return s.flStoreEngine.UpdateLoading(id, loading)
}
