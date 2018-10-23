package storengine

import (
	"github.com/hashicorp/go-uuid"
	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/pkg/mc"
)

type DatadirsMemory struct {
	DB map[string]model.DatadirSchema
}

func NewDatadirsMemory() *DatadirsMemory {
	return &DatadirsMemory{
		DB: make(map[string]model.DatadirSchema),
	}
}

func NewDatadirsMemoryWithDB(db map[string]model.DatadirSchema) *DatadirsMemory {
	return &DatadirsMemory{
		DB: db,
	}
}

func (e *DatadirsMemory) AddDir(dir model.DatadirSchema) (model.DatadirSchema, error) {
	var err error
	if dir.ID, err = uuid.GenerateUUID(); err != nil {
		return model.DatadirSchema{}, err
	}
	e.DB[dir.ID] = dir

	return dir, nil
}

func (e *DatadirsMemory) GetDatadirByPathInProject(path, projectID string) (model.DatadirSchema, error) {
	for _, ddir := range e.DB {
		if ddir.Name == path && ddir.Project == projectID {
			return ddir, nil
		}
	}

	return model.DatadirSchema{}, mc.ErrNotFound
}

func (e *DatadirsMemory) GetDatadir(id string) (model.DatadirSchema, error) {
	ddir, ok := e.DB[id]
	if !ok {
		return model.DatadirSchema{}, mc.ErrNotFound
	}

	return ddir, nil
}
