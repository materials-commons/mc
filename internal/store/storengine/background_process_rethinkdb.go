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

	defer resp.Close()

	if err := checkRethinkdbQueryError(resp, err, errMsg); err != nil {
		return returnList, err
	}

	err = resp.All(&returnList)

	return returnList, err
}

func (e *BackgroundProcessRethinkdb) DeleteBackgroundProcess(id string) error {
	errMsg := fmt.Sprintf("failed deleting background_process record %s", id)
	resp, err := r.Table("background_process").Get(id).Delete().RunWrite(e.Session)

	return checkRethinkdbDeleteError(resp, err, errMsg)
}

func (e *BackgroundProcessRethinkdb) GetBackgroundProcess(id string) (model.BackgroundProcessSchema, error) {
	var bgp model.BackgroundProcessSchema
	errMsg := fmt.Sprintf("failed to get background_process record %s", id)

	resp, err := r.Table("background_process").Get(id).Run(e.Session)
	defer resp.Close()

	if err := checkRethinkdbQueryError(resp, err, errMsg); err != nil {
		return bgp, err
	}

	err = resp.One(&bgp)

	return bgp, err
}

func (e *BackgroundProcessRethinkdb) SetFinishedBackgroundProcess(id string, done bool) error {
	errMsg := fmt.Sprintf("failed update on %s: IsFinished = %t", id, done)

	resp, err := r.Table("background_process").Get(id).
		Update(map[string]interface{}{"is_finished": done}, r.UpdateOpts{ReturnChanges: true}).RunWrite(e.Session)

	return checkRethinkdbUpdateError(resp, err, errMsg)
}

func (e *BackgroundProcessRethinkdb) SetOkBackgroundProcess(id string, success bool) error {
	errMsg := fmt.Sprintf("failed update on %s: IsOk = %t", id, success)

	resp, err := r.Table("background_process").Get(id).
		Update(map[string]interface{}{"is_ok": success}, r.UpdateOpts{ReturnChanges: true}).RunWrite(e.Session)

	return checkRethinkdbUpdateError(resp, err, errMsg)
}

func (e *BackgroundProcessRethinkdb) UpdateStatusBackgroundProcess(id string, status string, message string) error {
	errMsg := fmt.Sprintf("failed update on %s: status = %s, message = %s", id, status, message)

	resp, err := r.Table("background_process").Get(id).
		Update(map[string]interface{}{"status": status, "message": message}, r.UpdateOpts{ReturnChanges: true}).RunWrite(e.Session)

	return checkRethinkdbUpdateError(resp, err, errMsg)
}
