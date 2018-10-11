package store

import "time"

type ProjectsStoreEngineMemory struct {
	DB map[string]ProjectSchema
}

func NewProjectsStoreEngineMemory() *ProjectsStoreEngineMemory {
	return &ProjectsStoreEngineMemory{
		DB: make(map[string]ProjectSchema),
	}
}

func (e *ProjectsStoreEngineMemory) AddProject(project ProjectSchema) (ProjectSchema, error) {
	return ProjectSchema{}, nil
}

func (e *ProjectsStoreEngineMemory) GetProject(id string) (ProjectExtendedModel, error) {
	return ProjectExtendedModel{}, nil
}

func (e *ProjectsStoreEngineMemory) GetProjectSimple(id string) (ProjectSimpleModel, error) {
	return ProjectSimpleModel{}, nil
}

func (e *ProjectsStoreEngineMemory) GetAllProjectsForUser(user string) ([]ProjectExtendedModel, error) {
	return nil, nil
}

func (e *ProjectsStoreEngineMemory) DeleteProject(id string) error {
	return nil
}

func (e *ProjectsStoreEngineMemory) UpdateProjectName(id string, name string, updatedAt time.Time) error {
	return nil
}

func (e *ProjectsStoreEngineMemory) UpdateProjectDescription(id string, description string, updatedAt time.Time) error {
	return nil
}
