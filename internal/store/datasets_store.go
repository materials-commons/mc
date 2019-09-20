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
	return s.dsStoreEngine.GetDatadirsForDataset(datasetID)
}

func (s *DatasetsStore) GetDataset(datasetID string) (model.DatasetSchema, error) {
	return s.dsStoreEngine.GetDataset(datasetID)
}

func (s *DatasetsStore) SetDatasetZipfile(datasetID string, size int64, name string) error {
	return s.dsStoreEngine.SetDatasetZipfile(datasetID, size, name)
}
