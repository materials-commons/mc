package store

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	r "gopkg.in/gorethink/gorethink.v4"
)

type ProjectsStoreEngineRethinkdb struct {
	Session *r.Session
}

func (e *ProjectsStoreEngineRethinkdb) AddProject(project ProjectSchema) (ProjectSchema, error) {
	return ProjectSchema{}, nil
}

func (e *ProjectsStoreEngineRethinkdb) GetProject(id string) (ProjectExtendedModel, error) {
	var project ProjectExtendedModel

	res, err := r.Table("projects").Get(id).Merge(projectDetails).Run(e.Session)

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

func (e *ProjectsStoreEngineRethinkdb) GetAllProjectsForUser(user string) ([]ProjectExtendedModel, error) {
	var (
		userProjects     []ProjectExtendedModel
		projectsMemberOf []ProjectExtendedModel
	)

	res, err := r.Table("projects").GetAllByIndex("owner", user).Merge(projectDetails).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, fmt.Sprintf("Can't retrieve projects for user %s", user)); err != nil {
		return userProjects, err
	} else if err := res.All(userProjects); err != nil {
		return userProjects, err
	}

	res, err = r.Table("access").GetAllByIndex("user_id", user).
		EqJoin("project_id", r.Table("projects")).Zip().Filter(r.Row.Field("owner").Ne(user)).
		Merge(projectDetails).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, fmt.Sprintf("Can't retrieve projects user (%s) is member of", user)); err != nil {
		return userProjects, err
	} else if err := res.All(projectsMemberOf); err != nil {
		return userProjects, err
	}

	userProjects = append(userProjects, projectsMemberOf...)
	return userProjects, nil
}

func projectDetails(p r.Term) interface{} {
	return map[string]interface{}{
		"owner_details": r.Table("users").Get(p.Field("owner")).Pluck("fullname"),
		"users": r.Table("access").GetAllByIndex("project_id", p.Field("id")).
			EqJoin("user_id", r.Table("users")).Zip().CoerceTo("array"),
		"samples": r.Table("project2sample").GetAllByIndex("project_id", p.Field("id")).
			EqJoin("sample_id", r.Table("samples")).Zip().CoerceTo("array"),
		"processes": r.Table("project2process").GetAllByIndex("project_id", p.Field("id")).
			EqJoin("process_id", r.Table("processes")).Zip().CoerceTo("array"),
		"experiments": r.Table("project2experiment").GetAllByIndex("project_id", p.Field("id")).
			EqJoin("experiment_id", r.Table("experiments")).Zip().CoerceTo("array"),
		"files_count": r.Table("project2datafile").GetAllByIndex("project_id", p.Field("id")).Count(),
		"relationships": map[string]interface{}{
			"process2sample": r.Table("project2process").GetAllByIndex("project_id", p.Field("id")).
				EqJoin("process_id", r.Table("process2sample"), r.EqJoinOpts{Index: "process_id"}).Zip().
				Pluck("sample_id", "property_set_id", "process_id", "direction").CoerceTo("array"),
			"experiment2sample": r.Table("project2experiment").GetAllByIndex("project_id", p.Field("id")).
				EqJoin("experiment_id", r.Table("experiment2sample"), r.EqJoinOpts{Index: "experiment_id"}).Zip().
				Pluck("experiment_id", "sample_id").CoerceTo("array"),
		},
	}
}

func (e *ProjectsStoreEngineRethinkdb) DeleteProject(id string) error {
	errMsg := fmt.Sprintf("failed deleting project %s", id)

	resp, err := r.Table("projects").Get(id).
		Update(map[string]interface{}{"owner": "delete@materialscommons.org"}).RunWrite(e.Session)

	if err := checkRethinkdbWriteError(resp, err, errMsg); err != nil {
		return err
	}

	resp, err = r.Table("access").GetAllByIndex("project_id", id).Delete().RunWrite(e.Session)
	return checkRethinkdbWriteError(resp, err, errMsg)
}

func (e *ProjectsStoreEngineRethinkdb) UpdateProjectName(id string, name string, updatedAt time.Time) error {
	resp, err := r.Table("projects").Get(id).
		Update(map[string]interface{}{"name": name, "updated_at": updatedAt}).RunWrite(e.Session)
	return checkRethinkdbWriteError(resp, err, fmt.Sprintf("Unable to update name to '%s' for project %s", name, id))
}

func (e *ProjectsStoreEngineRethinkdb) UpdateProjectDescription(id string, description string, updatedAt time.Time) error {
	resp, err := r.Table("projects").Get(id).
		Update(map[string]interface{}{"description": description, "updated_at": updatedAt}).RunWrite(e.Session)
	return checkRethinkdbWriteError(resp, err, fmt.Sprintf("Unable to update desciption to '%s' for project %s", description, id))
}

func (e *ProjectsStoreEngineRethinkdb) Name() string {
	return "ProjectsStoreEngineRethinkdb"
}

func checkRethinkdbQueryError(res *r.Cursor, err error, msg string) error {
	switch {
	case err != nil:
		return errors.Wrapf(err, msg)
	case res.IsNil():
		return errors.Wrapf(ErrNotFound, msg)
	default:
		return nil
	}
}

func checkRethinkdbWriteError(resp r.WriteResponse, err error, msg string) error {
	switch {
	case err != nil:
		return errors.Wrapf(err, msg)
	case resp.Errors != 0:
		return fmt.Errorf("%s: %s", msg, resp.FirstError)
	default:
		return nil
	}
}
