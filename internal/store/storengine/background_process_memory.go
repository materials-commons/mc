package storengine

import (
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/pkg/mc"
)

type BackgroundProcessMemory struct {
	DB map[string]model.BackgroundProcessSchema
}

func NewBackgroundProcessMemory() *BackgroundProcessMemory {
	return &BackgroundProcessMemory{
		DB: make(map[string]model.BackgroundProcessSchema),
	}
}

func NewBackgroundProcessMemoryWithDB(db map[string]model.BackgroundProcessSchema) *BackgroundProcessMemory {
	return &BackgroundProcessMemory{
		DB: db,
	}
}

func (e *BackgroundProcessMemory) AddBackgroundProcess(bgp model.BackgroundProcessSchema) (model.BackgroundProcessSchema, error) {
	var err error
	if bgp.ID == "" {
		if bgp.ID, err = uuid.GenerateUUID(); err != nil {
			return bgp, err
		}
	}
	_, ok := e.DB[bgp.ID]
	if !ok {
		e.DB[bgp.ID] = bgp
	} else {
		err := fmt.Errorf("BackgroundProcessMemory - ID already exists: %s", bgp.ID)
		return bgp, err
	}
	return bgp, nil
}

func (e *BackgroundProcessMemory) GetBackgroundProcess(id string) (model.BackgroundProcessSchema, error) {
	bgp, ok := e.DB[id]
	if ok {
		return bgp, nil
	}
	return model.BackgroundProcessSchema{}, mc.ErrNotFound
}

func (e *BackgroundProcessMemory) DeleteBackgroundProcess(id string) error {
	_, ok := e.DB[id]
	if !ok {
		return mc.ErrNotFound
	}

	delete(e.DB, id)
	return nil
}

func (e *BackgroundProcessMemory) getListMatches(candidate model.BackgroundProcessSchema, template model.GetListBackgroundProcessModel) bool {
	return (candidate.UserID == template.UserID &&
		candidate.ProjectID == template.ProjectID &&
		candidate.BackgroundTaskID == template.BackgroundTaskID)
}

func (e *BackgroundProcessMemory) GetListBackgroundProcess(glbpg model.GetListBackgroundProcessModel) ([]model.BackgroundProcessSchema, error) {
	var ret []model.BackgroundProcessSchema

	for id := range e.DB {
		probe := e.DB[id]
		if e.getListMatches(probe, glbpg) {
			ret = append(ret, probe)
		}
	}
	return ret, nil
}

func (e *BackgroundProcessMemory) UpdateStatusBackgroundProcess(id string, status string, message string) error {
	entry, err := e.GetBackgroundProcess(id)
	if err != nil {
		return err
	}
	entry.Status = status
	entry.Message = message
	e.DB[id] = entry
	return nil
}

func (e *BackgroundProcessMemory) SetFinishedBackgroundProcess(id string, done bool) error {
	entry, err := e.GetBackgroundProcess(id)
	if err != nil {
		return err
	}
	entry.IsFinished = done
	e.DB[id] = entry
	return nil
}

func (e *BackgroundProcessMemory) SetOkBackgroundProcess(id string, success bool) error {
	entry, err := e.GetBackgroundProcess(id)
	if err != nil {
		return err
	}
	entry.IsOk = success
	e.DB[id] = entry
	return nil
}
