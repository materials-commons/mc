package store

import (
	"time"

	"github.com/hashicorp/go-uuid"
)

type ProjectsStoreEngineMemory struct {
	DB map[string]ProjectSchema
}

func NewProjectsStoreEngineMemory() *ProjectsStoreEngineMemory {
	return &ProjectsStoreEngineMemory{
		DB: make(map[string]ProjectSchema),
	}
}

func NewProjectsStoreEngineMemoryWithDB(db map[string]ProjectSchema) *ProjectsStoreEngineMemory {
	return &ProjectsStoreEngineMemory{
		DB: db,
	}
}

func (e *ProjectsStoreEngineMemory) AddProject(project ProjectSchema) (ProjectSchema, error) {
	var err error
	if project.ID, err = uuid.GenerateUUID(); err != nil {
		return ProjectSchema{}, err
	}

	e.DB[project.ID] = project

	return project, nil
}

func (e *ProjectsStoreEngineMemory) GetProject(id string) (ProjectExtendedModel, error) {
	proj, ok := e.DB[id]
	if !ok {
		return ProjectExtendedModel{}, ErrNotFound
	}

	p := ProjectExtendedModel{
		ProjectSchema: proj,
	}
	return p, nil
}

func (e *ProjectsStoreEngineMemory) GetProjectSimple(id string) (ProjectSimpleModel, error) {
	proj, ok := e.DB[id]
	if !ok {
		return ProjectSimpleModel{}, ErrNotFound
	}

	p := ProjectSimpleModel{
		ProjectSchema: proj,
	}

	return p, nil
}

func (e *ProjectsStoreEngineMemory) GetAllProjectsForUser(user string) ([]ProjectExtendedModel, error) {
	var userProjects []ProjectExtendedModel
	for _, proj := range e.DB {
		if proj.Owner == user {
			p := ProjectExtendedModel{
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
		return ErrNotFound
	}

	delete(e.DB, id)

	return nil
}

func (e *ProjectsStoreEngineMemory) UpdateProjectName(id string, name string, updatedAt time.Time) error {
	proj, ok := e.DB[id]
	if !ok {
		return ErrNotFound
	}

	proj.Name = name
	proj.MTime = updatedAt
	e.DB[id] = proj
	return nil
}

func (e *ProjectsStoreEngineMemory) UpdateProjectDescription(id string, description string, updatedAt time.Time) error {
	proj, ok := e.DB[id]
	if !ok {
		return ErrNotFound
	}

	proj.Description = description
	proj.MTime = updatedAt
	e.DB[id] = proj
	return nil
}
