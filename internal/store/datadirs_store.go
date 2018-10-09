package store

import "time"

type DatadirsStore struct {
	ddStoreEngine DatadirsStoreEngine
}

func NewDatadirsStore(e DatadirsStoreEngine) *DatadirsStore {
	return &DatadirsStore{ddStoreEngine: e}
}

func (s *DatadirsStore) AddDatadir(ddModel AddDatadirModel) (DatadirSchema, error) {
	if err := ddModel.Validate(); err != nil {
		return DatadirSchema{}, err
	}

	now := time.Now()

	dd := DatadirSchema{
		Model: Model{
			Name:      ddModel.Name,
			Owner:     ddModel.Owner,
			Birthtime: now,
			MTime:     now,
			OType:     "datadir",
		},
		Parent:  ddModel.Parent,
		Project: ddModel.ProjectID,
	}

	return s.ddStoreEngine.AddDir(dd)
}

func (s *DatadirsStore) GetDatadirByPathInProject(path, projectID string) (DatadirSchema, error) {
	return s.ddStoreEngine.GetDatadirByPathInProject(path, projectID)
}

func (s *DatadirsStore) GetDatadirByID(id string) (DatadirSchema, error) {
	return s.ddStoreEngine.GetDatadir(id)
}
