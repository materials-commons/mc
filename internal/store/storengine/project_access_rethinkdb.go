package storengine

import (
	"fmt"

	"github.com/materials-commons/mc/internal/store/model"

	r "gopkg.in/gorethink/gorethink.v4"
	"gopkg.in/gorethink/gorethink.v4/encoding"
)

type ProjectAccessRethinkdb struct {
	Session *r.Session
}

func NewProjectAccessRethinkdb(session *r.Session) *ProjectAccessRethinkdb {
	return &ProjectAccessRethinkdb{Session: session}
}

func (e *ProjectAccessRethinkdb) AddAccessEntry(entry model.ProjectAccessSchema) (model.ProjectAccessSchema, error) {
	errMsg := fmt.Sprintf("Unable to insert access entry %+v", entry)

	resp, err := r.Table("access").Insert(entry, r.InsertOpts{ReturnChanges: true}).RunWrite(e.Session)
	if err := checkRethinkdbInsertError(resp, err, errMsg); err != nil {
		return entry, err
	}

	var accessEntry model.ProjectAccessSchema
	err = encoding.Decode(&accessEntry, resp.Changes[0].NewValue)
	return accessEntry, err
}

func (e *ProjectAccessRethinkdb) DeleteAccess(projectID, userID string) error {
	errMsg := fmt.Sprintf("Unable to delete access entry for user %s from project %s", userID, projectID)
	resp, err := r.Table("access").GetAllByIndex("user_project", []interface{}{userID, projectID}).
		Delete().RunWrite(e.Session)
	return checkRethinkdbDeleteError(resp, err, errMsg)
}

func (e *ProjectAccessRethinkdb) DeleteAllAccessForProject(projectID string) error {
	errMsg := fmt.Sprintf("Unable to delete access entries for project %s", projectID)
	resp, err := r.Table("access").GetAllByIndex("project_id", projectID).Delete().RunWrite(e.Session)
	return checkRethinkdbDeleteError(resp, err, errMsg)
}

func (e *ProjectAccessRethinkdb) GetProjectAccess(projectID string) (model.ProjectAccessModel, error) {
	var projectAccessModel model.ProjectAccessModel
	errMsg := fmt.Sprintf("Couldn't retrieve project access for project %s", projectID)

	res, err := r.Table("projects").Get(projectID).Pluck("id", "owner").
		Merge(accessEntries).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return projectAccessModel, err
	}

	defer res.Close()

	err = res.One(&projectAccessModel)
	return projectAccessModel, err
}

func accessEntries(p r.Term) interface{} {
	return map[string]interface{}{
		"access_entries": r.Table("access").GetAllByIndex("project_id", p.Field("id")).
			EqJoin("user_id", r.Table("users")).
			Without(map[string]interface{}{"right": map[string]interface{}{"id": true}}).Zip().CoerceTo("array"),
	}
}

func (e *ProjectAccessRethinkdb) GetProjectAccessEntries(projectID string) ([]model.ProjectAccessSchema, error) {
	var entries []model.ProjectAccessSchema
	errMsg := fmt.Sprintf("Can't retrieve entries for project %s", projectID)

	res, err := r.Table("access").GetAllByIndex("project_id", projectID).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return entries, err
	}
	defer res.Close()

	err = res.All(&entries)
	return entries, err
}

func (e *ProjectAccessRethinkdb) GetUserAccessEntries(userID string) ([]model.ProjectAccessSchema, error) {
	var entries []model.ProjectAccessSchema
	errMsg := fmt.Sprintf("Can't retrieve access entries for user %s", userID)
	res, err := r.Table("access").GetAllByIndex("user_id", userID).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return entries, err
	}
	defer res.Close()

	err = res.All(&entries)
	return entries, err
}
