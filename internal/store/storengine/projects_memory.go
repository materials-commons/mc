package storengine

import (
	"time"

	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/pkg/mc"

	"github.com/hashicorp/go-uuid"
)

type ProjectsStoreEngineMemory struct {
	DB map[string]model.ProjectSchema
}

func NewProjectsStoreEngineMemory() *ProjectsStoreEngineMemory {
	return &ProjectsStoreEngineMemory{
		DB: make(map[string]model.ProjectSchema),
	}
}

func NewProjectsStoreEngineMemoryWithDB(db map[string]model.ProjectSchema) *ProjectsStoreEngineMemory {
	return &ProjectsStoreEngineMemory{
		DB: db,
	}
}

func (e *ProjectsStoreEngineMemory) AddProject(project model.ProjectSchema) (model.ProjectSchema, error) {
	var err error
	if project.ID, err = uuid.GenerateUUID(); err != nil {
		return model.ProjectSchema{}, err
	}

	e.DB[project.ID] = project

	return project, nil
}

func (e *ProjectsStoreEngineMemory) GetProject(id string) (model.ProjectExtendedModel, error) {
	proj, ok := e.DB[id]
	if !ok {
		return model.ProjectExtendedModel{}, mc.ErrNotFound
	}

	p := model.ProjectExtendedModel{
		ProjectSchema: proj,
	}
	return p, nil
}

func (e *ProjectsStoreEngineMemory) GetProjectSimple(id string) (model.ProjectSimpleModel, error) {
	proj, ok := e.DB[id]
	if !ok {
		return model.ProjectSimpleModel{}, mc.ErrNotFound
	}

	p := model.ProjectSimpleModel{
		ProjectSchema: proj,
	}

	return p, nil
}

func (e *ProjectsStoreEngineMemory) GetAllProjectsForUser(user string) ([]model.ProjectExtendedModel, error) {
	var userProjects []model.ProjectExtendedModel
	for _, proj := range e.DB {
		if proj.Owner == user {
			p := model.ProjectExtendedModel{
				ProjectSchema: proj,
			}
			userProjects = append(userProjects, p)
		}
	}
	return userProjects, nil
}

func (e *ProjectsStoreEngineMemory) DeleteProject(id string) error {
	_, ok := e.DB[id]
	if !ok {
		return mc.ErrNotFound
	}

	delete(e.DB, id)

	return nil
}

func (e *ProjectsStoreEngineMemory) UpdateProjectName(id string, name string, updatedAt time.Time) error {
	proj, ok := e.DB[id]
	if !ok {
		return mc.ErrNotFound
	}

	proj.Name = name
	proj.MTime = updatedAt
	e.DB[id] = proj
	return nil
}

func (e *ProjectsStoreEngineMemory) UpdateProjectDescription(id string, description string, updatedAt time.Time) error {
	proj, ok := e.DB[id]
	if !ok {
		return mc.ErrNotFound
	}

	proj.Description = description
	proj.MTime = updatedAt
	e.DB[id] = proj
	return nil
}
