package globus

import "github.com/materials-commons/mc/internal/store/model"

type ProjectDownload struct {
	ProjectID string
	User      *model.UserSchema
}

func NewProjectDownload(projectID string, user *model.UserSchema) *ProjectDownload {
	return &ProjectDownload{ProjectID: projectID, User: user}
}

// CreateProjectDownloadDirectory creates a download directory for the given project. It will
// walk through the projects directories and create them in this temporary download. It will
// then create a link to in this directory that points at the file entry in the Materials Commons
// store. The reason this needs to be done is that the Materials Commons store is an object store
// (like S3), where as Globus (and users) need to see the imposed directory structure. This is
// reconstructed from the database. Links to files are used so we don't have to create copies
// of the files.
func (d *ProjectDownload) CreateProjectDownloadDirectory() {

}
