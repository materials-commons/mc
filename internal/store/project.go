package store

type ProjectSchema struct {
	Model
	Description string `db:"description" json:"description"`
}

type ProjectExtendedModel struct {
	ProjectSchema
	FilesCount    int                  `json:"files_count"`
	Users         []ProjectUserModel   `json:"users"`
	Samples       []SampleSchema       `json:"samples"`
	Processes     []ProcessSchema      `json:"processes"`
	Experiments   []ExperimentSchema   `json:"experiments"`
	Relationships ProjectRelationships `json:"relationships"`
}

type ProjectRelationships struct {
	Process2Sample    []Process2Sample    `json:"process2sample"`
	Experiment2Sample []Experiment2Sample `json:"experiment2sample"`
}

type ProjectUserModel struct {
	ModelSimple
	BetaUser  bool   `db:"beta_user" json:"beta_user" r:"beta_user"`
	Fullname  string `db:"fullname" json:"fullname" r:"fullname"`
	ProjectID string `db:"project_id" json:"project_id" r:"project_id"`
	UserID    string `db:"user_id" json:"user_id" r:"user_id"`
}

type ProjectSimpleModel struct {
	ProjectSchema
	RootDir []DatadirSchema `json:"root_dir" r:"root_dir"`
}
