package store

import (
	"time"
)

type DatafilesStore struct {
	dfStoreEngine DatafilesStoreEngine
}

func NewDatafilesStore(e DatafilesStoreEngine) *DatafilesStore {
	return &DatafilesStore{dfStoreEngine: e}
}

func (s *DatafilesStore) AddDatafile(dfModel AddDatafileModel) (DatafileSchema, error) {
	if err := dfModel.Validate(); err != nil {
		return DatafileSchema{}, err
	}

	now := time.Now()
	df := DatafileSchema{
		Model: Model{
			Owner:     dfModel.Owner,
			Name:      dfModel.Name,
			Birthtime: now,
			MTime:     now,
			OType:     "datafile",
		},
		Description: dfModel.Description,
		Size:        dfModel.Size,
		Checksum:    dfModel.Checksum,
		UsesID:      dfModel.UsesID,
		Parent:      dfModel.Parent,
		Current:     true,
		MediaType:   dfModel.MediaType,
	}

	return s.dfStoreEngine.AddFile(df, dfModel.ProjectID, dfModel.DatadirID)
}

func (s *DatafilesStore) GetDatafileByID(id string) (DatafileSchema, error) {
	return s.dfStoreEngine.GetFile(id)
}

func (s *DatafilesStore) GetDatafileWithChecksum(checksum string) (DatafileSchema, error) {
	return s.dfStoreEngine.GetFileWithChecksum(checksum)
}

func (s *DatafilesStore) GetDatafileInDir(name, datadirID string) (DatafileSchema, error) {
	return s.dfStoreEngine.GetFileInDir(name, datadirID)
}

func (s *DatafilesStore) UpdateDatafileCurrentFlag(id string, current bool) error {
	return s.dfStoreEngine.UpdateFileCurrentFlag(id, current)
}
