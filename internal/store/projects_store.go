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

func (s *ProjectsStore) GetProjectsForUser(userID string) ([]model.ProjectCountModel, error) {
	return s.pStoreEngine.GetAllProjectsForUser(userID)
}

func (s *ProjectsStore) GetProjectOverview(projectID, userID string) (model.ProjectOverviewModel, error) {
	return s.pStoreEngine.GetProjectOverview(projectID, userID)
}

func (s *ProjectsStore) GetProjectAccessEntries(projectID string) ([]model.ProjectUserAccessModel, error) {
	return s.pStoreEngine.GetProjectAccessEntries(projectID)
}

func (s *ProjectsStore) AddUserToProject(projectID, userID string) (model.ProjectAccessEntry, error) {
	return s.pStoreEngine.AddAccessToProject(projectID, userID)
}

func (s *ProjectsStore) DeleteAccessEntry(id string) error {
	return s.pStoreEngine.DeleteAccessEntry(id)
}
