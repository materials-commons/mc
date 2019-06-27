package cmd

import (
	"fmt"
	"path/filepath"
)

const publicDatasetsDir = "__published_datasets"
const privateDatasetsDir = "__datasets"

func publicDatasetPath(mcdir string, datasetID string) string {
	return filepath.Join(mcdir, publicDatasetsDir, datasetID)
}

func privateDatasetPath(mcdir string, datasetID string) string {
	return filepath.Join(mcdir, privateDatasetsDir, datasetID)
}

func publicGlobusDatasetPath(datasetID string) string {
	return fmt.Sprintf("/%s/%s/", publicDatasetsDir, datasetID)
}

func privateGlobusDatasetPath(datasetID string) string {
	return fmt.Sprintf("/%s/%s/", privateDatasetsDir, datasetID)
}
