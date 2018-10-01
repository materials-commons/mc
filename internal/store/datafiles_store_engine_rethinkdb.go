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
	err = encoding.Decode(&f, resp.Changes[0].NewValue)
	return f, err
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
