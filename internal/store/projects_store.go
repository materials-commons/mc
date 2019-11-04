package store

import (
	"github.com/materials-commons/mc/internal/store/storengine"

	"github.com/materials-commons/mc/internal/store/model"
)

type ProjectsStore struct {
	pStoreEngine storengine.ProjectsStoreEngine
}

func NewProjectsStore(e storengine.ProjectsStoreEngine) *ProjectsStore {
	return &ProjectsStore{pStoreEngine: e}
}

func (s *ProjectsStore) GetProjectSimple(id string) (model.ProjectSimpleModel, error) {
	return s.pStoreEngine.GetProjectSimple(id)
}

func (s *ProjectsStore) GetProjectUsers(id string) ([]model.UserSchema, error) {
	return s.pStoreEngine.GetProjectUsers(id)
}
