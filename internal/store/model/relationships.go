package model

type Process2Sample struct {
	SampleID      string `db:"sample_id" json:"sample_id"`
	PropertySetID string `db:"property_set_id" json:"property_set_id"`
	ProcessID     string `db:"process_id" json:"process_id"`
	Direction     string `db:"direction" json:"direction"`
}

type Experiment2Sample struct {
	ExperimentID string `db:"experiment_id" json:"experiment_id"`
	SampleID     string `db:"sample_id" json:"sample_id"`
}
