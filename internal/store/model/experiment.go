package model

type ExperimentSchema struct {
	Model
	Description string `db:"description" json:"description"`
	Status      string `db:"status" json:"status"`
}

type ExperimentOverviewModel struct {
	ExperimentSchema
	FilesCount   int `db:"files_count" json:"files_count" r:"files_count"`
	SamplesCount int `db:"samples_count" json:"samples_count" r:"samples_count"`
	ProcessCount int `db:"process_count" json:"process_count" r:"process_count"`
}
