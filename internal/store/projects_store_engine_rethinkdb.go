package store

import (
	"time"

	"github.com/pkg/errors"

	r "gopkg.in/gorethink/gorethink.v4"
)

type ProjectsStoreEngineRethinkdb struct {
	session *r.Session
}

func (e *ProjectsStoreEngineRethinkdb) AddProject(project ProjectSchema) (ProjectSchema, error) {
	return ProjectSchema{}, nil
}

func (e *ProjectsStoreEngineRethinkdb) GetProject(id string) (ProjectExtendedModel, error) {
	var project ProjectExtendedModel
	res, err := r.Table("projects").Get(id).Merge(func(p r.Term) interface{} {
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
	}).Run(e.session)

	switch {
	case err != nil:
		return project, err
	case res.IsNil():
		return project, errors.Wrapf(ErrNotFound, "No such project %s", id)
	default:
		err = res.One(&project)
		return project, err
	}
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
