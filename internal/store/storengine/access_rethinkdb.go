package storengine

import (
	"fmt"

	"github.com/materials-commons/mc/internal/store/model"

	r "gopkg.in/gorethink/gorethink.v4"
	"gopkg.in/gorethink/gorethink.v4/encoding"
)

type AccessRethinkdb struct {
	Session *r.Session
}

func NewAccessRethinkdb(session *r.Session) *AccessRethinkdb {
	return &AccessRethinkdb{Session: session}
}

func (e *AccessRethinkdb) AddAccessEntry(entry model.AccessSchema) (model.AccessSchema, error) {
	errMsg := fmt.Sprintf("Unable to insert access entry %+v", entry)

	resp, err := r.Table("access").Insert(entry, r.InsertOpts{ReturnChanges: true}).RunWrite(e.Session)
	if err := checkRethinkdbInsertError(resp, err, errMsg); err != nil {
		return entry, err
	}

	var accessEntry model.AccessSchema
	err = encoding.Decode(&accessEntry, resp.Changes[0].NewValue)
	return accessEntry, err
}

func (e *AccessRethinkdb) DeleteAccess(projectID, userID string) error {
	errMsg := fmt.Sprintf("Unable to delete access entry for user %s from project %s", userID, projectID)
	resp, err := r.Table("access").GetAllByIndex("user_project", []interface{}{userID, projectID}).
		Delete().RunWrite(e.Session)
	return checkRethinkdbDeleteError(resp, err, errMsg)
}

func (e *AccessRethinkdb) DeleteAllAccessForProject(projectID string) error {
	errMsg := fmt.Sprintf("Unable to delete access entries for project %s", projectID)
	resp, err := r.Table("access").GetAllByIndex("project_id", projectID).Delete().RunWrite(e.Session)
	return checkRethinkdbDeleteError(resp, err, errMsg)
}

func (e *AccessRethinkdb) GetProjectAccessEntries(projectID string) ([]model.AccessSchema, error) {
	var entries []model.AccessSchema
	errMsg := fmt.Sprintf("Can't retrieve entries for project %s", projectID)

	res, err := r.Table("access").GetAllByIndex("project_id", projectID).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return entries, err
	}
	defer res.Close()

	err = res.All(entries)
	return entries, err
}

func (e *AccessRethinkdb) GetUserAccessEntries(userID string) ([]model.AccessSchema, error) {
	var entries []model.AccessSchema
	errMsg := fmt.Sprintf("Can't retrieve access entries for user %s", userID)
	res, err := r.Table("access").GetAllByIndex("user_id", userID).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return entries, err
	}
	defer res.Close()

	err = res.All(entries)
	return entries, err
}
