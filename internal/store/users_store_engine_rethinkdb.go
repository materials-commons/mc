package store

import (
	"fmt"
	"time"

	r "gopkg.in/gorethink/gorethink.v4"
	"gopkg.in/gorethink/gorethink.v4/encoding"
)

type UsersStoreEngineRethinkdb struct {
	Session *r.Session
}

func NewUsersStoreEngineRethinkdb(session *r.Session) *UsersStoreEngineRethinkdb {
	return &UsersStoreEngineRethinkdb{Session: session}
}

func (e *UsersStoreEngineRethinkdb) AddUser(user UserSchema) (UserSchema, error) {
	errMsg := fmt.Sprintf("Unable to add user %+v", user)
	resp, err := r.Table("users").Insert(user).RunWrite(e.Session)
	return user, checkRethinkdbInsertError(resp, err, errMsg)
}

func (e *UsersStoreEngineRethinkdb) GetUserByID(id string) (UserSchema, error) {
	var user UserSchema
	errMsg := fmt.Sprintf("No such user %s", id)
	res, err := r.Table("users").Get(id).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return user, err
	}

	err = res.One(&user)
	return user, err
}

func (e *UsersStoreEngineRethinkdb) GetUserByAPIKey(apikey string) (UserSchema, error) {
	var user UserSchema
	errMsg := fmt.Sprintf("No such apikey %s", apikey)
	res, err := r.Table("users").GetAllByIndex("apikey", apikey).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return user, err
	}

	err = res.One(&user)
	return user, err
}

func (e *UsersStoreEngineRethinkdb) ModifyUserFullname(id, fullname string, updatedAt time.Time) (UserSchema, error) {
	return e.modifyUser(id, map[string]interface{}{"fullname": fullname, "updated_at": updatedAt})
}

func (e *UsersStoreEngineRethinkdb) ModifyUserPassword(id, password string, updatedAt time.Time) (UserSchema, error) {
	return e.modifyUser(id, map[string]interface{}{"password": password, "updated_at": updatedAt})
}

func (e *UsersStoreEngineRethinkdb) ModifyUserAPIKey(id, apikey string, updatedAt time.Time) (UserSchema, error) {
	return e.modifyUser(id, map[string]interface{}{"apikey": apikey, "updated_at": updatedAt})
}

func (e *UsersStoreEngineRethinkdb) modifyUser(id string, what map[string]interface{}) (UserSchema, error) {
	resp, err := r.Table("users").Get(id).Update(what, r.UpdateOpts{ReturnChanges: true}).RunWrite(e.Session)
	switch {
	case err != nil:
		return UserSchema{}, err
	case resp.Errors != 0:
		return UserSchema{}, fmt.Errorf("unable to update user %s", id)
	default:
		var u UserSchema
		if len(resp.Changes) == 0 {
			return u, fmt.Errorf("unable to modify %s with %#v", id, what)
		}
		err := encoding.Decode(&u, resp.Changes[0].NewValue)
		return u, err
	}
}
