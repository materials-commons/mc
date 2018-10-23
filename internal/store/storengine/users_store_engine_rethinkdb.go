package storengine

import (
	"fmt"
	"time"

	"github.com/materials-commons/mc/internal/store/model"

	r "gopkg.in/gorethink/gorethink.v4"
	"gopkg.in/gorethink/gorethink.v4/encoding"
)

type UsersStoreEngineRethinkdb struct {
	Session *r.Session
}

func NewUsersStoreEngineRethinkdb(session *r.Session) *UsersStoreEngineRethinkdb {
	return &UsersStoreEngineRethinkdb{Session: session}
}

func (e *UsersStoreEngineRethinkdb) AddUser(user model.UserSchema) (model.UserSchema, error) {
	errMsg := fmt.Sprintf("Unable to add user %+v", user)
	resp, err := r.Table("users").Insert(user, r.InsertOpts{ReturnChanges: true}).RunWrite(e.Session)
	return user, checkRethinkdbInsertError(resp, err, errMsg)
}

func (e *UsersStoreEngineRethinkdb) GetUserByID(id string) (model.UserSchema, error) {
	var user model.UserSchema
	errMsg := fmt.Sprintf("No such user %s", id)
	res, err := r.Table("users").Get(id).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return user, err
	}
	defer res.Close()

	err = res.One(&user)
	return user, err
}

func (e *UsersStoreEngineRethinkdb) GetUserByAPIKey(apikey string) (model.UserSchema, error) {
	var user model.UserSchema
	errMsg := fmt.Sprintf("No such apikey %s", apikey)
	res, err := r.Table("users").GetAllByIndex("apikey", apikey).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return user, err
	}
	defer res.Close()

	err = res.One(&user)
	return user, err
}

func (e *UsersStoreEngineRethinkdb) ModifyUserFullname(id, fullname string, updatedAt time.Time) (model.UserSchema, error) {
	return e.modifyUser(id, map[string]interface{}{"fullname": fullname, "updated_at": updatedAt})
}

func (e *UsersStoreEngineRethinkdb) ModifyUserPassword(id, password string, updatedAt time.Time) (model.UserSchema, error) {
	return e.modifyUser(id, map[string]interface{}{"password": password, "updated_at": updatedAt})
}

func (e *UsersStoreEngineRethinkdb) ModifyUserAPIKey(id, apikey string, updatedAt time.Time) (model.UserSchema, error) {
	return e.modifyUser(id, map[string]interface{}{"apikey": apikey, "updated_at": updatedAt})
}

func (e *UsersStoreEngineRethinkdb) modifyUser(id string, what map[string]interface{}) (model.UserSchema, error) {
	resp, err := r.Table("users").Get(id).Update(what, r.UpdateOpts{ReturnChanges: true}).RunWrite(e.Session)
	switch {
	case err != nil:
		return model.UserSchema{}, err
	case resp.Errors != 0:
		return model.UserSchema{}, fmt.Errorf("unable to update user %s", id)
	default:
		var u model.UserSchema
		if len(resp.Changes) == 0 {
			return u, fmt.Errorf("unable to modify %s with %#v", id, what)
		}
		err := encoding.Decode(&u, resp.Changes[0].NewValue)
		return u, err
	}
}
