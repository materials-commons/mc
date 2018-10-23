package store

import (
	"fmt"
	"time"

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
	return addDatadir(dir, e.Session)
}

func addDatadir(dir DatadirSchema, session *r.Session) (DatadirSchema, error) {
	errMsg := fmt.Sprintf("Unable to add datadir: %+v", dir)
	resp, err := r.Table("datadirs").Insert(dir, r.InsertOpts{ReturnChanges: true}).RunWrite(session)
	if err := checkRethinkdbInsertError(resp, err, errMsg); err != nil {
		return dir, err
	}

	var d DatadirSchema
	if err := encoding.Decode(&d, resp.Changes[0].NewValue); err != nil {
		return d, err
	}

	proj2datadir := map[string]interface{}{"project_id": dir.Project, "datadir_id": d.ID}

	resp, err = r.Table("project2datadir").Insert(proj2datadir, r.InsertOpts{ReturnChanges: true}).RunWrite(session)

	if err := checkRethinkdbInsertError(resp, err, errMsg); err != nil {
		return d, err
	}

	return d, nil
}

func toDatadirSchema(ddModel AddDatadirModel) DatadirSchema {
	now := time.Now()

	dd := DatadirSchema{
		Model: Model{
			Name:      ddModel.Name,
			Owner:     ddModel.Owner,
			Birthtime: now,
			MTime:     now,
			OType:     "datadir",
		},
		Parent:  ddModel.Parent,
		Project: ddModel.ProjectID,
	}

	return dd
}

func (e *DatadirsStoreEngineRethinkdb) GetDatadirByPathInProject(path, projectID string) (DatadirSchema, error) {
	var dir DatadirSchema
	errMsg := fmt.Sprintf("Unable to find datadir path %s in project %s", path, projectID)
	res, err := r.Table("datadirs").
		GetAllByIndex("datadir_project_name", []interface{}{projectID, path}).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return dir, err
	}
	defer res.Close()

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
	defer res.Close()

	err = res.One(&dir)
	return dir, err
}
