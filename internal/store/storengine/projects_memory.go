package storengine

import (
	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/pkg/mc"
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

func (e *ProjectsMemory) GetProjectUsers(id string) ([]model.UserSchema, error) {
	return []model.UserSchema{}, nil
}
