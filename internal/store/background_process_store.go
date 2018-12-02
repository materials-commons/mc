package store

import (
	"time"

	"github.com/materials-commons/mc/internal/store/storengine"

	"github.com/materials-commons/mc/internal/store/model"
)

type BackgroundProcessStore struct {
	bgpStoreEngine storengine.BackgroundProcessStoreEngine
}

func NewBackgroundProcessStore(e storengine.BackgroundProcessStoreEngine) *BackgroundProcessStore {
	return &BackgroundProcessStore{bgpStoreEngine: e}
}

func (s *BackgroundProcessStore) AddBackgroundProcess(bgpAddModel model.AddBackgroundProcessModel) (model.BackgroundProcessSchema, error) {
	now := time.Now()
	schema := model.BackgroundProcessSchema{
		ModelSimple: model.ModelSimple{
			Birthtime: now,
			MTime:     now,
			OType:     "background_process",
		},
		UserID:                bgpAddModel.UserID,
		ProjectID:             bgpAddModel.ProjectID,
		BackgroundProcessID:   bgpAddModel.BackgroundProcessID,
		BackgroundProcessType: bgpAddModel.BackgroundProcessType,
		Status:                bgpAddModel.Status,
		Message:               bgpAddModel.Message,
		IsFinished:            false,
		IsOk:                  false,
	}

	return s.bgpStoreEngine.AddBackgroundProcess(schema)
}
