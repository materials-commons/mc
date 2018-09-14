package store

import (
	"fmt"
	"time"

	r "gopkg.in/gorethink/gorethink.v4"
	"gopkg.in/gorethink/gorethink.v4/encoding"
)

type ProjectsStoreEngineRethinkdb struct {
	Session *r.Session
}

func NewProjectsStoreEngineRethinkdb(session *r.Session) *ProjectsStoreEngineRethinkdb {
	return &ProjectsStoreEngineRethinkdb{Session: session}
}

func (e *ProjectsStoreEngineRethinkdb) AddProject(project ProjectSchema) (ProjectSchema, error) {
	errMsg := fmt.Sprintf("Unable to add project %+v", project)
	resp, err := r.Table("projects").Insert(project, r.InsertOpts{ReturnChanges: true}).RunWrite(e.Session, r.RunOpts{})
	if err := checkRethinkdbInsertError(resp, err, errMsg); err != nil {
		return project, err
	}

	var proj ProjectSchema
	err = encoding.Decode(&proj, resp.Changes[0].NewValue)
	return proj, err
}

func (e *ProjectsStoreEngineRethinkdb) GetProject(id string) (ProjectExtendedModel, error) {
	var project ProjectExtendedModel
	errMsg := fmt.Sprintf("No such project %s", id)
	res, err := r.Table("projects").Get(id).Merge(projectDetails).Run(e.Session)

	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return project, err
	}

	err = res.One(&project)
	return project, err
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
	fmt.Println("Delete Project", id)
	errMsg := fmt.Sprintf("failed deleting project %s", id)

	resp, err := r.Table("projects").Get(id).
		Update(map[string]interface{}{"owner": "delete@materialscommons.org"}).RunWrite(e.Session)

	if err := checkRethinkdbUpdateError(resp, err, errMsg); err != nil {
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
