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
	var ret model.BackgroundProcessSchema

	errMsg := fmt.Sprintf("Unable to add background_process: %+v", bgp)

	resp, err := r.Table("background_process").Insert(bgp, r.InsertOpts{ReturnChanges: true}).RunWrite(e.Session)
	if err := checkRethinkdbInsertError(resp, err, errMsg); err != nil {
		return bgp, err
	}

	if err := encoding.Decode(&ret, resp.Changes[0].NewValue); err != nil {
		return ret, err
	}

	return ret, nil
}

func (e *BackgroundProcessRethinkdb) GetListBackgroundProcess(glbgp model.GetListBackgroundProcessModel) ([]model.BackgroundProcessSchema, error) {
	var returnList []model.BackgroundProcessSchema

	queryIndex := []interface{}{glbgp.UserID, glbgp.ProjectID, glbgp.BackgroundProcessID}
	errMsg := fmt.Sprintf("Unable to GetAll background_process records: %+v", queryIndex)

	resp, err := r.Table("background_process").
		GetAllByIndex("user_project_process", queryIndex).Run(e.Session)
	if err := checkRethinkdbQueryError(resp, err, errMsg); err != nil {
		return returnList, err
	}

	defer resp.Close()

	err = resp.All(&returnList)

	return returnList, err
}
