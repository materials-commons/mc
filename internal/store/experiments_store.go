package store

import (
	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/internal/store/storengine"
)

type ExperimentsStore struct {
	eStoreEngine storengine.ExperimentsStoreEngine
}

func NewExperimentsStore(e storengine.ExperimentsStoreEngine) *ExperimentsStore {
	return &ExperimentsStore{eStoreEngine: e}
}

func (s *ExperimentsStore) GetExperimentOverviewsForProject(projectID string) ([]model.ExperimentOverviewModel, error) {
	return s.eStoreEngine.GetExperimentOverviewsForProject(projectID)
}
