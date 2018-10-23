package store

import (
	"fmt"

	r "gopkg.in/gorethink/gorethink.v4"
	"gopkg.in/gorethink/gorethink.v4/encoding"
)

type AccessStoreEngineRethinkdb struct {
	Session *r.Session
}

func NewAccessStoreEngineRethinkdb(session *r.Session) *AccessStoreEngineRethinkdb {
	return &AccessStoreEngineRethinkdb{Session: session}
}

func (e *AccessStoreEngineRethinkdb) AddAccessEntry(entry AccessSchema) (AccessSchema, error) {
	errMsg := fmt.Sprintf("Unable to insert access entry %+v", entry)

	resp, err := r.Table("access").Insert(entry, r.InsertOpts{ReturnChanges: true}).RunWrite(e.Session)
	if err := checkRethinkdbInsertError(resp, err, errMsg); err != nil {
		return entry, err
	}

	var accessEntry AccessSchema
	err = encoding.Decode(&accessEntry, resp.Changes[0].NewValue)
	return accessEntry, err
}

func (e *AccessStoreEngineRethinkdb) DeleteAccess(projectID, userID string) error {
	errMsg := fmt.Sprintf("Unable to delete access entry for user %s from project %s", userID, projectID)
	resp, err := r.Table("access").GetAllByIndex("user_project", []interface{}{userID, projectID}).
		Delete().RunWrite(e.Session)
	return checkRethinkdbDeleteError(resp, err, errMsg)
}

func (e *AccessStoreEngineRethinkdb) DeleteAllAccessForProject(projectID string) error {
	errMsg := fmt.Sprintf("Unable to delete access entries for project %s", projectID)
	resp, err := r.Table("access").GetAllByIndex("project_id", projectID).Delete().RunWrite(e.Session)
	return checkRethinkdbDeleteError(resp, err, errMsg)
}

func (e *AccessStoreEngineRethinkdb) GetProjectAccessEntries(projectID string) ([]AccessSchema, error) {
	var entries []AccessSchema
	errMsg := fmt.Sprintf("Can't retrieve entries for project %s", projectID)

	res, err := r.Table("access").GetAllByIndex("project_id", projectID).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return entries, err
	}
	defer res.Close()

	err = res.All(entries)
	return entries, err
}

func (e *AccessStoreEngineRethinkdb) GetUserAccessEntries(userID string) ([]AccessSchema, error) {
	var entries []AccessSchema
	errMsg := fmt.Sprintf("Can't retrieve access entries for user %s", userID)
	res, err := r.Table("access").GetAllByIndex("user_id", userID).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return entries, err
	}
	defer res.Close()

	err = res.All(entries)
	return entries, err
}
