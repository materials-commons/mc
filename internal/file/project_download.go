package file

import (
	"github.com/apex/log"
)

type ProjectDownload struct {
	downloadDir *DownloadDir
}

func NewProjectDownload(downloadDir *DownloadDir) *ProjectDownload {
	return &ProjectDownload{downloadDir: downloadDir}
}

// CreateProjectDownloadDirectory creates a download directory for the given project.
func (d *ProjectDownload) CreateProjectDownloadDirectory(basePath string) error {
	log.Infof("CreateProjectDownloadDirectory %s", basePath)
	ddirs, err := d.downloadDir.ddirsStore.GetDatadirsForProject(d.downloadDir.ProjectID, d.downloadDir.User.ID)
	if err != nil {
		log.Infof("GetDatadirsForProject returned %s", err)
		return err
	}
	return d.downloadDir.CreateDownloadDirectory(basePath, ddirs)
}
