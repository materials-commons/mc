package storengine

import (
	"fmt"
	"time"

	"github.com/materials-commons/mc/internal/store/model"

	"gopkg.in/gorethink/gorethink.v4/encoding"

	r "gopkg.in/gorethink/gorethink.v4"
)

type DatadirsRethinkdb struct {
	Session *r.Session
}

func NewDatadirsRethinkdb(session *r.Session) *DatadirsRethinkdb {
	return &DatadirsRethinkdb{Session: session}
}

func (e *DatadirsRethinkdb) AddDir(dir model.DatadirSchema) (model.DatadirSchema, error) {
	return AddDatadir(dir, e.Session)
}

func AddDatadir(dir model.DatadirSchema, session *r.Session) (model.DatadirSchema, error) {
	errMsg := fmt.Sprintf("Unable to add datadir: %+v", dir)
	resp, err := r.Table("datadirs").Insert(dir, r.InsertOpts{ReturnChanges: true}).RunWrite(session)
	if err := checkRethinkdbInsertError(resp, err, errMsg); err != nil {
		return dir, err
	}

	var d model.DatadirSchema
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

func ToDatadirSchema(ddModel model.AddDatadirModel) model.DatadirSchema {
	now := time.Now()

	dd := model.DatadirSchema{
		Model: model.Model{
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

func (e *DatadirsRethinkdb) GetDatadirByPathInProject(path, projectID string) (model.DatadirSchema, error) {
	var dir model.DatadirSchema
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

func (e *DatadirsRethinkdb) GetDatadir(id string) (model.DatadirSchema, error) {
	var dir model.DatadirSchema
	errMsg := fmt.Sprintf("Unable to find datadir %s", id)
	res, err := r.Table("datadirs").Get(id).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return dir, err
	}
	defer res.Close()

	err = res.One(&dir)
	return dir, err
}
