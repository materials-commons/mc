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

func (s *BackgroundProcessStore) AddBackgroundProcess(bgpModel model.AddBackgroundProcessModel) (model.BackgroundProcessSchema, error) {
	now := time.Now()
	bgp := model.BackgroundProcessSchema{
		ModelSimple: model.ModelSimple{
			Birthtime: now,
			MTime:     now,
			OType:     "background_process",
		},
		ProjectID:             bgpModel.ProjectID,
		UserID:                bgpModel.UserID,
    	BackgroundProcessID:   bgpModel.BackgroundProcessID,
    	BackgroundProcessType: bgpModel.BackgroundProcessType,
	    Status:                bgpModel.Status,
	    Message:               bgpModel.Message,
	    IsFinished:            false,
	    IsOk:                  false,
	}

	return s.bgpStoreEngine.AddBackgroundProcess(bgp)
}
