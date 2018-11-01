package store

import (
	"time"

	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/internal/store/storengine"
)

type GlobusUploadsStore struct {
	guStoreEngine storengine.GlobusUploadsStoreEngine
}

func NewGlobusUploadsStore(e storengine.GlobusUploadsStoreEngine) *GlobusUploadsStore {
	return &GlobusUploadsStore{guStoreEngine: e}
}

func (s *GlobusUploadsStore) AddGlobusUpload(upload model.AddGlobusUploadModel) (model.GlobusUploadSchema, error) {
	if err := upload.Validate(); err != nil {
		return model.GlobusUploadSchema{}, err
	}

	now := time.Now()
	gupload := model.GlobusUploadSchema{
		ModelSimple: model.ModelSimple{
			ID:        upload.ID,
			Birthtime: now,
			MTime:     now,
			OType:     "globus_upload",
		},
		Owner:            upload.Owner,
		Path:             upload.Path,
		ProjectID:        upload.ProjectID,
		GlobusAclID:      upload.GlobusAclID,
		GlobusEndpointID: upload.GlobusEndpointID,
		GlobusIdentityID: upload.GlobusIdentityID,
	}

	return s.guStoreEngine.AddGlobusUpload(gupload)
}

func (s *GlobusUploadsStore) DeleteGlobusUpload(id string) error {
	return s.guStoreEngine.DeleteGlobusUpload(id)
}

func (s *GlobusUploadsStore) GetGlobusUpload(id string) (model.GlobusUploadSchema, error) {
	return s.guStoreEngine.GetGlobusUpload(id)
}

func (s *GlobusUploadsStore) GetAllGlobusUploads() ([]model.GlobusUploadSchema, error) {
	return s.guStoreEngine.GetAllGlobusUploads()
}

func (s *GlobusUploadsStore) GetAllGlobusUploadsForUser(user string) ([]model.GlobusUploadSchema, error) {
	return s.guStoreEngine.GetAllGlobusUploadsForUser(user)
}
