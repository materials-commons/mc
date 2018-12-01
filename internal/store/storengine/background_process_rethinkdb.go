package storengine

import (
	"fmt"

	"github.com/materials-commons/mc/internal/store/model"

	"gopkg.in/gorethink/gorethink.v4/encoding"

	r "gopkg.in/gorethink/gorethink.v4"
)

type BackgroundProcessRethinkdb struct {
	Session *r.Session
}

func NewBackgroundProcessRethinkdb(session *r.Session) *BackgroundProcessRethinkdb {
	return &BackgroundProcessRethinkdb{Session: session}
}

func (e *BackgroundProcessRethinkdb) AddBackgroundProcess(bgp model.BackgroundProcessSchema) (model.BackgroundProcessSchema, error) {
	errMsg := fmt.Sprintf("Unable to add background_process: %+v", bgp)

	resp, err := r.Table("background_process").Insert(bgp, r.InsertOpts{ReturnChanges: true}).RunWrite(e.Session)
	if err := checkRethinkdbInsertError(resp, err, errMsg); err != nil {
		return bgp, err
	}

	var ret model.BackgroundProcessSchema
	if err := encoding.Decode(&ret, resp.Changes[0].NewValue); err != nil {
		return ret, err
	}

	return ret, nil
}
