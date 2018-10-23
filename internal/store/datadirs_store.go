package store

import (
	"time"

	"github.com/materials-commons/mc/internal/store/model"
)

type DatadirsStore struct {
	ddStoreEngine DatadirsStoreEngine
}

func NewDatadirsStore(e DatadirsStoreEngine) *DatadirsStore {
	return &DatadirsStore{ddStoreEngine: e}
}

func (s *DatadirsStore) AddDatadir(ddModel model.AddDatadirModel) (model.DatadirSchema, error) {
	if err := ddModel.Validate(); err != nil {
		return model.DatadirSchema{}, err
	}

	now := time.Now()

	dd := model.DatadirSchema{
		Model: model.Model{
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

func (s *DatadirsStore) GetDatadirByPathInProject(path, projectID string) (model.DatadirSchema, error) {
	return s.ddStoreEngine.GetDatadirByPathInProject(path, projectID)
}

func (s *DatadirsStore) GetDatadirByID(id string) (model.DatadirSchema, error) {
	return s.ddStoreEngine.GetDatadir(id)
}
