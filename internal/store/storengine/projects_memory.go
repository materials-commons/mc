package storengine

import (
	"time"

	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/pkg/mc"

	"github.com/hashicorp/go-uuid"
)

type ProjectsMemory struct {
	DB map[string]model.ProjectSchema
}

func NewProjectsMemory() *ProjectsMemory {
	return &ProjectsMemory{
		DB: make(map[string]model.ProjectSchema),
	}
}

func NewProjectsMemoryWithDB(db map[string]model.ProjectSchema) *ProjectsMemory {
	return &ProjectsMemory{
		DB: db,
	}
}

func (e *ProjectsMemory) AddProject(project model.ProjectSchema) (model.ProjectSchema, error) {
	var err error
	if project.ID, err = uuid.GenerateUUID(); err != nil {
		return model.ProjectSchema{}, err
	}

	e.DB[project.ID] = project

	return project, nil
}

func (e *ProjectsMemory) GetProjectAccessEntries(id string) ([]model.ProjectUserAccessModel, error) {
	return nil, mc.ErrNotFound
}

func (e *ProjectsMemory) GetProjectOverview(projectID, userID string) (model.ProjectOverviewModel, error) {
	return model.ProjectOverviewModel{}, mc.ErrNotFound
}

func (e *ProjectsMemory) GetProjectSimple(id string) (model.ProjectSimpleModel, error) {
	proj, ok := e.DB[id]
	if !ok {
		return model.ProjectSimpleModel{}, mc.ErrNotFound
	}

	p := model.ProjectSimpleModel{
		ProjectSchema: proj,
	}

	return p, nil
}

func (e *ProjectsMemory) GetProjectNotes(projectID, userID string) ([]model.ProjectNote, error) {
	return nil, mc.ErrNotFound
}

func (e *ProjectsMemory) GetAllProjectsForUser(user string) ([]model.ProjectCountModel, error) {
	var userProjects []model.ProjectCountModel
	for _, proj := range e.DB {
		if proj.Owner == user {
			p := model.ProjectCountModel{
				ProjectSchema: proj,
			}
			userProjects = append(userProjects, p)
		}
	}
	return userProjects, nil
}

func (e *ProjectsMemory) DeleteProject(id string) error {
	_, ok := e.DB[id]
	if !ok {
		return mc.ErrNotFound
	}

	delete(e.DB, id)

	return nil
}

func (e *ProjectsMemory) UpdateProjectName(id string, name string, updatedAt time.Time) error {
	proj, ok := e.DB[id]
	if !ok {
		return mc.ErrNotFound
	}

	proj.Name = name
	proj.MTime = updatedAt
	e.DB[id] = proj
	return nil
}

func (e *ProjectsMemory) UpdateProjectDescription(id string, description string, updatedAt time.Time) error {
	proj, ok := e.DB[id]
	if !ok {
		return mc.ErrNotFound
	}

	proj.Description = description
	proj.MTime = updatedAt
	e.DB[id] = proj
	return nil
}
