package store

import (
	"fmt"

	"gopkg.in/gorethink/gorethink.v4/encoding"

	r "gopkg.in/gorethink/gorethink.v4"
)

type DatadirsStoreEngineRethinkdb struct {
	Session *r.Session
}

func NewDatadirsStoreEngineRethinkdb(session *r.Session) *DatadirsStoreEngineRethinkdb {
	return &DatadirsStoreEngineRethinkdb{Session: session}
}

func (e *DatadirsStoreEngineRethinkdb) AddDir(dir DatadirSchema) (DatadirSchema, error) {
	errMsg := fmt.Sprintf("Unable to add datadir: %+v", dir)
	resp, err := r.Table("datadirs").Insert(dir, r.InsertOpts{ReturnChanges: true}).RunWrite(e.Session)
	if err := checkRethinkdbInsertError(resp, err, errMsg); err != nil {
		return dir, err
	}

	var d DatadirSchema
	err = encoding.Decode(&d, resp.Changes[0].NewValue)
	return d, err
}

func (e *DatadirsStoreEngineRethinkdb) GetDatadirByPathInProject(path, projectID string) (DatadirSchema, error) {
	var dir DatadirSchema
	errMsg := fmt.Sprintf("Unable to find datadir path %s in project %s", path, projectID)
	res, err := r.Table("datadirs").
		GetAllByIndex("datadir_project_name", []interface{}{projectID, path}).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return dir, err
	}

	err = res.One(&dir)
	return dir, err
}

func (e *DatadirsStoreEngineRethinkdb) GetDatadir(id string) (DatadirSchema, error) {
	var dir DatadirSchema
	errMsg := fmt.Sprintf("Unable to find datadir %s", id)
	res, err := r.Table("datadirs").Get(id).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return dir, err
	}

	err = res.One(&dir)
	return dir, err
}
