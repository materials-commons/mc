package store

import (
	"time"

	r "gopkg.in/gorethink/gorethink.v4"
)

type ProjectsStoreEngineRethinkdb struct {
	s *r.Session
}

func (e *ProjectsStoreEngineRethinkdb) AddProject(project ProjectSchema) (ProjectSchema, error) {
	return ProjectSchema{}, nil
}

func (e *ProjectsStoreEngineRethinkdb) GetProject(id string) (ProjectExtendedModel, error) {
	rql := r.Table("projects").Get(id).Merge(func(p r.Term) interface{} {
		return map[string]interface{}{
			"samples": r.Table("project2sample").GetAllByIndex("project_id", p.Field("id")).
				EqJoin("sample_id", r.Table("samples")).Zip().CoerceTo("array"),
			"processes": r.Table("project2process").GetAllByIndex("project_id", p.Field("id")).
				EqJoin("process_id", r.Table("processes")).Zip().CoerceTo("array"),
			"experiments": r.Table("project2experiment").GetAllByIndex("project_id", p.Field("id")).
				EqJoin("experiment_id", r.Table("experiments")).Zip().CoerceTo("array"),
			"relationships": map[string]interface{}{
				"process2sample": r.Table("project2process").GetAllByIndex("project_id", p.Field("id")).
					EqJoin("process_id", r.Table("process2sample"), r.EqJoinOpts{Index: "process_id"}).Zip().
					Pluck("sample_id", "property_set_id", "process_id", "direction").CoerceTo("array"),
				"experiment2sample": r.Table("project2experiment").GetAllByIndex("project_id", p.Field("id")).
					EqJoin("experiment_id", r.Table("experiment2sample"), r.EqJoinOpts{Index: "experiment_id"}).Zip().
					Pluck("experiment_id", "sample_id").CoerceTo("array"),
			},
		}
	})
	res, err := rql.Run(e.s)
	if err != nil {
		return ProjectExtendedModel{}, err
	}

	var project ProjectExtendedModel
	err = res.One(&project)
	return project, err
}

func (e *ProjectsStoreEngineRethinkdb) GetAllProjectsForUser(user string) ([]ProjectSchema, error) {
	return nil, nil
}

func (e *ProjectsStoreEngineRethinkdb) DeleteProject(id string) error {
	return nil
}

func (e *ProjectsStoreEngineRethinkdb) UpdateProjectName(id string, name string, updatedAt time.Time) error {
	return nil
}

func (e *ProjectsStoreEngineRethinkdb) UpdateProjectDescription(id string, description string, updatedAt time.Time) error {
	return nil
}

func (e *ProjectsStoreEngineRethinkdb) Name() string {
	return "ProjectsStoreEngineRethinkdb"
}
