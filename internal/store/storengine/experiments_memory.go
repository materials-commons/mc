package storengine

import "github.com/materials-commons/mc/internal/store/model"

type ExperimentsMemory struct {
	DB map[string]model.ExperimentSchema
}

func NewExperimentsMemory() *ExperimentsMemory {
	return &ExperimentsMemory{
		DB: make(map[string]model.ExperimentSchema),
	}
}

func NewExperimentsMemoryWithDB(db map[string]model.ExperimentSchema) *ExperimentsMemory {
	return &ExperimentsMemory{
		DB: db,
	}
}

func (e *ExperimentsMemory) GetExperimentOverviewsForProject(projectID string) ([]model.ExperimentOverviewModel, error) {
	return nil, nil
}
