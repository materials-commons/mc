package model

type DataclipSchema struct {
	Model
	Description  string
	ProjectID    string
	ExperimentID string
	ProcessIDs   []string
	SampleIDs    []string
	Access       string
}

type AddDataclipModel struct {
	Name         string
	Owner        string
	ProjectID    string
	ExperimentID string
	Access       string
}
