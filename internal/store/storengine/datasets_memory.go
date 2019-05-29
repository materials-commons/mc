package storengine

import "github.com/materials-commons/mc/internal/store/model"

type DatasetsMemory struct {
	DB map[string]model.DatasetSchema
}

func NewDatasetsMemory() *DatasetsMemory {
	return &DatasetsMemory{
		DB: make(map[string]model.DatasetSchema),
	}
}

func NewDatasetsMemoryWithDB(db map[string]model.DatasetSchema) *DatasetsMemory {
	return &DatasetsMemory{DB: db}
}

func (e *DatasetsMemory) GetDatadirsForDataset(datasetID string) ([]model.DatadirEntryModel, error) {
	return nil, nil
}

func (e *DatasetsMemory) GetDataset(datasetID string) (model.DatasetSchema, error) {
	return model.DatasetSchema{}, nil
}
