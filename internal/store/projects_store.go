package store

import (
	"time"

	"github.com/materials-commons/mc/internal/store/storengine"

	"github.com/materials-commons/mc/internal/store/model"
)

type ProjectsStore struct {
	pStoreEngine storengine.ProjectsStoreEngine
}

func NewProjectsStore(e storengine.ProjectsStoreEngine) *ProjectsStore {
	return &ProjectsStore{pStoreEngine: e}
}

func (s *ProjectsStore) AddProject(pModel model.AddProjectModel) (model.ProjectSchema, error) {
	if err := pModel.Validate(); err != nil {
		return model.ProjectSchema{}, err
	}

	now := time.Now()

	p := model.ProjectSchema{
		Model: model.Model{
			Name:      pModel.Name,
			Owner:     pModel.Owner,
			Birthtime: now,
			MTime:     now,
			OType:     "project",
		},
		Description: pModel.Description,
	}

	return s.pStoreEngine.AddProject(p)
}

func (s *ProjectsStore) GetProjectSimple(id string) (model.ProjectSimpleModel, error) {
	return s.pStoreEngine.GetProjectSimple(id)
}
