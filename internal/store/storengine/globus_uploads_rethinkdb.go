package storengine

import (
	"fmt"

	"github.com/materials-commons/mc/internal/store/model"
	r "gopkg.in/gorethink/gorethink.v4"
	"gopkg.in/gorethink/gorethink.v4/encoding"
)

type GlobusUploadsRethinkdb struct {
	Session *r.Session
}

func NewGlobusUploadsRethinkdb(session *r.Session) *GlobusUploadsRethinkdb {
	return &GlobusUploadsRethinkdb{Session: session}
}

func (e *GlobusUploadsRethinkdb) AddGlobusUpload(upload model.GlobusUploadSchema) (model.GlobusUploadSchema, error) {
	errMsg := fmt.Sprintf("Unable to add globus upload %#v", upload)
	resp, err := r.Table("globus_uploads").Insert(upload, r.InsertOpts{ReturnChanges: true}).RunWrite(e.Session)
	if err := checkRethinkdbInsertError(resp, err, errMsg); err != nil {
		return upload, err
	}

	var uploadCreated model.GlobusUploadSchema
	err = encoding.Decode(&uploadCreated, resp.Changes[0].NewValue)
	return uploadCreated, err
}

func (e *GlobusUploadsRethinkdb) DeleteGlobusUpload(id string) error {
	errMsg := fmt.Sprintf("failed deleting globus upload %s", id)
	resp, err := r.Table("globus_uploads").Get(id).Delete().RunWrite(e.Session)
	return checkRethinkdbDeleteError(resp, err, errMsg)
}

func (e *GlobusUploadsRethinkdb) GetGlobusUpload(id string) (model.GlobusUploadSchema, error) {
	var upload model.GlobusUploadSchema
	errMsg := fmt.Sprintf("No such globus upload %s", id)
	res, err := r.Table("globus_uploads").Get(id).Run(e.Session)
	if err = checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return upload, err
	}
	defer res.Close()

	err = res.One(&upload)
	return upload, err
}

func (e *GlobusUploadsRethinkdb) GetAllGlobusUploads() ([]model.GlobusUploadSchema, error) {
	var uploads []model.GlobusUploadSchema
	errMsg := fmt.Sprintf("Couldn't retrieve all globus uploads")
	res, err := r.Table("globus_uploads").Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return uploads, err
	}

	defer res.Close()

	err = res.All(uploads)
	return uploads, err
}

func (e *GlobusUploadsRethinkdb) GetAllGlobusUploadsForUser(user string) ([]model.GlobusUploadSchema, error) {
	var uploads []model.GlobusUploadSchema
	errMsg := fmt.Sprintf("Couldn't retrieve globus uploads for user %s", user)
	res, err := r.Table("globus_uploads").GetAllByIndex("owner", user).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return uploads, err
	}

	defer res.Close()

	err = res.All(uploads)
	return uploads, err

}
