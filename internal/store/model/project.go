package model

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/materials-commons/mc/pkg/mc"
	"github.com/pkg/errors"
)

type ProjectSchema struct {
	Model
	Description string `db:"description" json:"description" r:"description"`
}

type ProjectCountModel struct {
	ProjectSchema
	OwnerDetails     OwnerDetails `json:"owner_details" r:"owner_details"`
	FilesCount       int          `json:"files_count" r:"files_count"`
	UsersCount       int          `json:"users_count" r:"users_count"`
	SamplesCount     int          `json:"samples_count" r:"samples_count"`
	ProcessesCount   int          `json:"processes_count" r:"processes_count"`
	ExperimentsCount int          `json:"experiments_count" r:"experiments_count"`
	DatasetsCount    int          `json:"datasets_count" r:"datasets_count"`
}

type OwnerDetails struct {
	Fullname string `json:"fullname" r:"fullname"`
}

type ProjectExtendedModel struct {
	ProjectSchema
	FilesCount    int                  `json:"files_count" r:"files_count"`
	Users         []ProjectUserModel   `json:"users" r:"users"`
	Samples       []SampleSchema       `json:"samples" r:"samples"`
	Processes     []ProcessSchema      `json:"processes" r:"processes"`
	Experiments   []ExperimentSchema   `json:"experiments" r:"experiments"`
	Relationships ProjectRelationships `json:"relationships" r:"relationships"`
}

type ProjectRelationships struct {
	Process2Sample    []Process2Sample    `json:"process2sample" r:"process2sample"`
	Experiment2Sample []Experiment2Sample `json:"experiment2sample" r:"experiment2sample"`
}

type ProjectUserModel struct {
	ModelSimple
	BetaUser  bool   `db:"beta_user" json:"beta_user" r:"beta_user"`
	Fullname  string `db:"fullname" json:"fullname" r:"fullname"`
	ProjectID string `db:"project_id" json:"project_id" r:"project_id"`
	UserID    string `db:"user_id" json:"user_id" r:"user_id"`
}

type ProjectNote struct {
	ModelSimple
	Note  string `json:"note" r:"note"`
	Owner string `json:"owner" r:"owner"`
	Title string `json:"title" r:"title"`
}

type ProjectTodo struct {
	Title string `json:"title" r:"title"`
	Edit  bool   `json:"edit" r:"edit"`
	Done  bool   `json:"done" r:"done"`
}

type ProjectOverviewModel struct {
	ProjectSchema
	Shortcuts      []ProjectShortcut         `json:"shortcuts" r:"shortcuts"`
	RootDir        []DatadirSimpleModel      `json:"root_dir" r:"root_dir"`
	OwnerDetails   OwnerDetails              `json:"owner_details" r:"owner_details"`
	FilesCount     int                       `json:"files_count" r:"files_count"`
	UsersCount     int                       `json:"users_count" r:"users_count"`
	SamplesCount   int                       `json:"samples_count" r:"samples_count"`
	ProcessesCount int                       `json:"processes_count" r:"processes_count"`
	DatasetsCount  int                       `json:"datasets_count" r:"datasets_count"`
	Experiments    []ExperimentOverviewModel `json:"experiments" r:"experiments"`
	Todos          []ProjectTodo             `json:"todos" r:"todos"`
}

type ProjectShortcut struct {
	Name string `json:"name" r:"name"`
	ID   string `json:"id" r:"id"`
}

type ProjectSimpleModel struct {
	ProjectSchema
	RootDir []DatadirSchema `json:"root_dir" r:"root_dir"`
}

type AddProjectModel struct {
	Name        string
	Owner       string
	Description string
}

func (p AddProjectModel) Validate() error {
	err := validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.Owner, validation.Required))

	if err != nil {
		return errors.WithMessage(mc.ErrValidation, err.Error())
	}

	return nil
}
