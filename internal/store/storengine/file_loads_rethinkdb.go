package storengine

import (
	"fmt"

	"github.com/materials-commons/mc/internal/store/model"

	r "gopkg.in/gorethink/gorethink.v4"
	"gopkg.in/gorethink/gorethink.v4/encoding"
)

type FileLoadsRethinkdb struct {
	Session *r.Session
}

func NewFileLoadsRethinkdb(session *r.Session) *FileLoadsRethinkdb {
	return &FileLoadsRethinkdb{Session: session}
}

func (e *FileLoadsRethinkdb) AddFileLoad(upload model.FileLoadSchema) (model.FileLoadSchema, error) {
	errMsg := fmt.Sprintf("Unable to add file load request %#v", upload)
	resp, err := r.Table("file_loads").Insert(upload, r.InsertOpts{ReturnChanges: true}).RunWrite(e.Session)
	if err := checkRethinkdbInsertError(resp, err, errMsg); err != nil {
		return upload, err
	}

	var uploadCreated model.FileLoadSchema
	err = encoding.Decode(&uploadCreated, resp.Changes[0].NewValue)
	return uploadCreated, err
}

func (e *FileLoadsRethinkdb) DeleteFileLoad(id string) error {
	errMsg := fmt.Sprintf("failed deleting file load entry %s", id)
	resp, err := r.Table("file_loads").Get(id).Delete().RunWrite(e.Session)
	return checkRethinkdbWriteError(resp, err, errMsg)
}

func (e *FileLoadsRethinkdb) GetFileLoad(uploadID string) (model.FileLoadSchema, error) {
	var upload model.FileLoadSchema
	errMsg := fmt.Sprintf("No such file load request %s", uploadID)
	res, err := r.Table("file_loads").Get(uploadID).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return upload, err
	}
	defer res.Close()

	err = res.One(&upload)
	return upload, err
}

func (e *FileLoadsRethinkdb) GetAllFileLoads() ([]model.FileLoadSchema, error) {
	errMsg := fmt.Sprintf("Unable to retrieve file loads")
	res, err := r.Table("file_loads").Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return nil, err
	}
	defer res.Close()

	var uploads []model.FileLoadSchema
	err = res.All(uploads)
	return uploads, err
}

func (e *FileLoadsRethinkdb) MarkAllNotLoading() error {
	errMsg := fmt.Sprintf("Unable to upload file loads")
	resp, err := r.Table("file_loads").Update(map[string]interface{}{"loading": false}).RunWrite(e.Session)
	return checkRethinkdbWriteError(resp, err, errMsg)
}

func (e *FileLoadsRethinkdb) UpdateLoading(id string, loading bool) error {
	errMsg := fmt.Sprintf("Unable to update file load %s", id)
	resp, err := r.Table("file_loads").
		Update(map[string]interface{}{"loading": loading}, r.UpdateOpts{ReturnChanges: true}).RunWrite(e.Session)
	return checkRethinkdbWriteError(resp, err, errMsg)
}
