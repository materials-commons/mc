package store

import "time"

type FileLoadsStore struct {
	flStoreEngine FileLoadsStoreEngine
}

func NewFileLoadsStore(e FileLoadsStoreEngine) *FileLoadsStore {
	return &FileLoadsStore{flStoreEngine: e}
}

func (s *FileLoadsStore) AddFileLoad(flModel AddFileLoadModel) (FileLoadSchema, error) {
	if err := flModel.Validate(); err != nil {
		return FileLoadSchema{}, err
	}

	now := time.Now()
	fl := FileLoadSchema{
		ModelSimple: ModelSimple{
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

func (s *FileLoadsStore) GetFileLoad(id string) (FileLoadSchema, error) {
	return s.flStoreEngine.GetFileLoad(id)
}

func (s *FileLoadsStore) GetAllFileLoads() ([]FileLoadSchema, error) {
	return s.flStoreEngine.GetAllFileLoads()
}

func (s *FileLoadsStore) MarkAllNotLoading() error {
	return s.flStoreEngine.MarkAllNotLoading()
}

func (s *FileLoadsStore) UpdateLoading(id string, loading bool) error {
	return s.flStoreEngine.UpdateLoading(id, loading)
}
