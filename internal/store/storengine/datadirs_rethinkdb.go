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

func (e *DatadirsRethinkdb) GetFilesForDatadir(projectID, userID, dirID string) ([]model.DatafileSimpleModel, error) {
	var files []model.DatafileSimpleModel
	errMsg := fmt.Sprintf("No files for directory %s in project %s for user %s", dirID, projectID, userID)
	res, err := r.Table("access").GetAllByIndex("user_project", []interface{}{userID, projectID}).
		EqJoin([]interface{}{r.Row.Field("project_id"), dirID}, r.Table("project2datadir"), r.EqJoinOpts{Index: "project_datadir"}).Zip().
		EqJoin("datadir_id", r.Table("datadir2datafile"), r.EqJoinOpts{Index: "datadir_id"}).Zip().
		EqJoin("datafile_id", r.Table("datafiles")).Zip().Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return files, err
	}
	defer res.Close()

	err = res.All(&files)
	return files, err
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

func (e *DatadirsRethinkdb) GetDatadirForProject(projectID, userID, dirID string) (model.DatadirEntryModel, error) {
	var dir model.DatadirEntryModel
	errMsg := fmt.Sprintf("No such directory %s in project %s for user %s", dirID, projectID, userID)
	res, err := r.Table("access").GetAllByIndex("user_project", []interface{}{userID, projectID}).
		EqJoin([]interface{}{r.Row.Field("project_id"), dirID}, r.Table("project2datadir"), r.EqJoinOpts{Index: "project_datadir"}).Zip().
		EqJoin("datadir_id", r.Table("datadirs")).Zip().Merge(directoryEntries).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return dir, err
	}
	defer res.Close()

	err = res.One(&dir)
	return dir, err
}

func directoryEntries(p r.Term) interface{} {
	return map[string]interface{}{
		"directories": r.Table("datadirs").GetAllByIndex("parent", p.Field("datadir_id")).CoerceTo("array"),
		"files": r.Table("datadir2datafile").GetAllByIndex("datadir_id", p.Field("datadir_id")).
			EqJoin("datafile_id", r.Table("datafiles")).Zip().CoerceTo("array"),
	}
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
