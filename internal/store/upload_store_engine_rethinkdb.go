package store

import (
	"fmt"

	r "gopkg.in/gorethink/gorethink.v4"
	"gopkg.in/gorethink/gorethink.v4/encoding"
)

type UploadsStoreEngineRethinkdb struct {
	Session *r.Session
}

func NewUploadsStoreEngineRethinkdb(session *r.Session) *UploadsStoreEngineRethinkdb {
	return &UploadsStoreEngineRethinkdb{Session: session}
}

func (e *UploadsStoreEngineRethinkdb) AddUpload(upload UploadSchema) (UploadSchema, error) {
	errMsg := fmt.Sprintf("Unable to add upload request %#v", upload)
	resp, err := r.Table("uploads").Insert(upload, r.InsertOpts{ReturnChanges: true}).RunWrite(e.Session)
	if err := checkRethinkdbInsertError(resp, err, errMsg); err != nil {
		return upload, err
	}

	var uploadCreated UploadSchema
	err = encoding.Decode(&uploadCreated, resp.Changes[0].NewValue)
	return uploadCreated, err
}

func (e *UploadsStoreEngineRethinkdb) DeleteUpload(id string) error {
	errMsg := fmt.Sprintf("failed deleting upload entry %s", id)
	resp, err := r.Table("uploads").Get(id).Delete().RunWrite(e.Session)
	return checkRethinkdbWriteError(resp, err, errMsg)
}

func (e *UploadsStoreEngineRethinkdb) GetUpload(uploadID string) (UploadSchema, error) {
	var upload UploadSchema
	errMsg := fmt.Sprintf("No such upload request %s", uploadID)
	res, err := r.Table("uploads").Get(uploadID).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return upload, err
	}

	err = res.One(&upload)
	return upload, err
}

func (e *UploadsStoreEngineRethinkdb) GetAllUploads() ([]UploadSchema, error) {
	errMsg := fmt.Sprintf("Unable to retrieve uploads")
	res, err := r.Table("uploads").Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return nil, err
	}

	var uploads []UploadSchema
	err = res.All(uploads)
	return uploads, err
}
