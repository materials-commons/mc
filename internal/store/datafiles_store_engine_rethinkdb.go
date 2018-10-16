package store

import (
	"fmt"

	r "gopkg.in/gorethink/gorethink.v4"
	"gopkg.in/gorethink/gorethink.v4/encoding"
)

type DatafilesStoreEngineRethinkdb struct {
	Session *r.Session
}

func NewDatafilesStoreEngineRethinkdb(session *r.Session) *DatafilesStoreEngineRethinkdb {
	return &DatafilesStoreEngineRethinkdb{Session: session}
}

func (e *DatafilesStoreEngineRethinkdb) AddFile(file DatafileSchema, projectID, datadirID string) (DatafileSchema, error) {
	errMsg := fmt.Sprintf("Unable to add datafile: %+v", file)
	resp, err := r.Table("datafiles").Insert(file, r.InsertOpts{ReturnChanges: true}).RunWrite(e.Session)
	if err := checkRethinkdbInsertError(resp, err, errMsg); err != nil {
		return file, err
	}

	var f DatafileSchema
	if err := encoding.Decode(&f, resp.Changes[0].NewValue); err != nil {
		return f, err
	}

	proj2datafile := map[string]interface{}{"project_id": projectID, "datafile_id": f.ID}
	resp, err = r.Table("project2datafile").Insert(proj2datafile, r.InsertOpts{ReturnChanges: true}).RunWrite(e.Session)
	if err := checkRethinkdbInsertError(resp, err, errMsg); err != nil {
		return f, err
	}

	datadir2datafile := map[string]interface{}{"datadir_id": datadirID, "datafile_id": f.ID}
	resp, err = r.Table("datadir2datafile").Insert(datadir2datafile, r.InsertOpts{ReturnChanges: true}).RunWrite(e.Session)
	return f, checkRethinkdbInsertError(resp, err, errMsg)
}

func (e *DatafilesStoreEngineRethinkdb) GetFile(id string) (DatafileSchema, error) {
	var file DatafileSchema
	errMsg := fmt.Sprintf("No such datafile %s", id)
	res, err := r.Table("datafiles").Get(id).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return file, err
	}

	err = res.One(&file)
	return file, err
}

func (e *DatafilesStoreEngineRethinkdb) GetFileWithChecksum(checksum string) (DatafileSchema, error) {
	var file DatafileSchema
	errMsg := fmt.Sprintf("No file matching checksum %s", checksum)
	res, err := r.Table("datafiles").GetAllByIndex("checksum", checksum).
		Filter(r.Row.Field("usesid").Eq("")).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return file, err
	}

	err = res.One(&file)
	return file, err
}

func (e *DatafilesStoreEngineRethinkdb) GetFileInDir(name string, dirID string) (DatafileSchema, error) {
	var file DatafileSchema
	errMsg := fmt.Sprintf("No file %s in dir %s", name, dirID)
	res, err := r.Table("datadir2datafile").GetAllByIndex("datadir_id", dirID).
		EqJoin("datafile_id", r.Table("datafiles")).Zip().
		Filter(r.Row.Field("current").Eq(true).And(r.Row.Field("name").Eq(name))).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return file, err
	}

	err = res.One(&file)
	return file, err
}

func (e *DatafilesStoreEngineRethinkdb) UpdateFileCurrentFlag(id string, current bool) error {
	errMsg := fmt.Sprintf("failed updated file %s current flag", id)

	resp, err := r.Table("datafiles").Get(id).
		Update(map[string]interface{}{"current": current}, r.UpdateOpts{ReturnChanges: true}).RunWrite(e.Session)

	return checkRethinkdbUpdateError(resp, err, errMsg)
}
