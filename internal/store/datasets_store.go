package store

import (
	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/internal/store/storengine"
)

type DatasetsStore struct {
	dsStoreEngine storengine.DatasetsStoreEngine
}

func NewDatasetsStore(e storengine.DatasetsStoreEngine) *DatasetsStore {
	return &DatasetsStore{dsStoreEngine: e}
}

func (s *DatasetsStore) GetDatadirsForDataset(datasetID string) ([]model.DatadirEntryModel, error) {
	return s.GetDatadirsForDataset(datasetID)
}

func (s *DatasetsStore) GetDataset(datasetID string) (model.DatasetSchema, error) {
	return s.GetDataset(datasetID)
}
