package file

import "github.com/materials-commons/mc/internal/store"

type DatasetDownload struct {
	downloadDir *DownloadDir
	dsStore     *store.DatasetsStore
}

func NewDatasetDownload(dsStore *store.DatasetsStore, downloadDir *DownloadDir) *DatasetDownload {
	return &DatasetDownload{downloadDir: downloadDir, dsStore: dsStore}
}

func (d *DatasetDownload) CreateDatasetDownloadDir(datasetID, basePath string) error {
	ddirs, err := d.dsStore.GetDatadirsForDataset(datasetID)
	if err != nil {
		return err
	}

	return d.downloadDir.CreateDownloadDirectory(basePath, ddirs)
}
