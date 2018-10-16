package store

import (
	"fmt"

	r "gopkg.in/gorethink/gorethink.v4"
	"gopkg.in/gorethink/gorethink.v4/encoding"
)

type FileLoadsStoreEngineRethinkdb struct {
	Session *r.Session
}

func NewFileLoadsStoreEngineRethinkdb(session *r.Session) *FileLoadsStoreEngineRethinkdb {
	return &FileLoadsStoreEngineRethinkdb{Session: session}
}

func (e *FileLoadsStoreEngineRethinkdb) AddFileLoad(upload FileLoadSchema) (FileLoadSchema, error) {
	errMsg := fmt.Sprintf("Unable to add file load request %#v", upload)
	resp, err := r.Table("file_loads").Insert(upload, r.InsertOpts{ReturnChanges: true}).RunWrite(e.Session)
	if err := checkRethinkdbInsertError(resp, err, errMsg); err != nil {
		return upload, err
	}

	var uploadCreated FileLoadSchema
	err = encoding.Decode(&uploadCreated, resp.Changes[0].NewValue)
	return uploadCreated, err
}

func (e *FileLoadsStoreEngineRethinkdb) DeleteFileLoad(id string) error {
	errMsg := fmt.Sprintf("failed deleting file load entry %s", id)
	resp, err := r.Table("file_loads").Get(id).Delete().RunWrite(e.Session)
	return checkRethinkdbWriteError(resp, err, errMsg)
}

func (e *FileLoadsStoreEngineRethinkdb) GetFileLoad(uploadID string) (FileLoadSchema, error) {
	var upload FileLoadSchema
	errMsg := fmt.Sprintf("No such file load request %s", uploadID)
	res, err := r.Table("file_loads").Get(uploadID).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return upload, err
	}

	err = res.One(&upload)
	return upload, err
}

func (e *FileLoadsStoreEngineRethinkdb) GetAllFileLoads() ([]FileLoadSchema, error) {
	errMsg := fmt.Sprintf("Unable to retrieve file loads")
	res, err := r.Table("file_loads").Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return nil, err
	}

	var uploads []FileLoadSchema
	err = res.All(uploads)
	return uploads, err
}

func (e *FileLoadsStoreEngineRethinkdb) MarkAllNotLoading() error {
	errMsg := fmt.Sprintf("Unable to upload file loads")
	resp, err := r.Table("file_loads").Update(map[string]interface{}{"loading": false}).RunWrite(e.Session)
	return checkRethinkdbWriteError(resp, err, errMsg)
}

func (e *FileLoadsStoreEngineRethinkdb) UpdateLoading(id string, loading bool) error {
	errMsg := fmt.Sprintf("Unable to update file load %s", id)
	resp, err := r.Table("file_loads").
		Update(map[string]interface{}{"loading": loading}, r.UpdateOpts{ReturnChanges: true}).RunWrite(e.Session)
	return checkRethinkdbWriteError(resp, err, errMsg)
}
