package store

import "time"

type ProjectsStore struct {
	pStoreEngine ProjectsStoreEngine
}

func NewProjectsStore(e ProjectsStoreEngine) *ProjectsStore {
	return &ProjectsStore{pStoreEngine: e}
}

func (s *ProjectsStore) AddProject(pModel AddProjectModel) (ProjectSchema, error) {
	if err := pModel.Validate(); err != nil {
		return ProjectSchema{}, err
	}

	now := time.Now()

	p := ProjectSchema{
		Model: Model{
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

func (s *ProjectsStore) GetProjectSimple(id string) (ProjectSimpleModel, error) {
	return s.pStoreEngine.GetProjectSimple(id)
}
