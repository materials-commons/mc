package storengine

import (
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

func (e *ProjectsMemory) DeleteProject(id string) error {
	_, ok := e.DB[id]
	if !ok {
		return mc.ErrNotFound
	}

	delete(e.DB, id)

	return nil
}

func (e *ProjectsMemory) GetProjectUsers(id string) ([]model.UserSchema, error) {
	return []model.UserSchema{}, nil
}
