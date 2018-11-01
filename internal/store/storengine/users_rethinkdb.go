package storengine

import (
	"fmt"
	"time"

	"github.com/materials-commons/mc/internal/store/model"

	r "gopkg.in/gorethink/gorethink.v4"
	"gopkg.in/gorethink/gorethink.v4/encoding"
)

type UsersRethinkdb struct {
	Session *r.Session
}

func NewUsersRethinkdb(session *r.Session) *UsersRethinkdb {
	return &UsersRethinkdb{Session: session}
}

func (e *UsersRethinkdb) AddUser(user model.UserSchema) (model.UserSchema, error) {
	errMsg := fmt.Sprintf("Unable to add user %+v", user)
	resp, err := r.Table("users").Insert(user, r.InsertOpts{ReturnChanges: true}).RunWrite(e.Session)
	return user, checkRethinkdbInsertError(resp, err, errMsg)
}

func (e *UsersRethinkdb) GetUserByID(id string) (model.UserSchema, error) {
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

func (e *UsersRethinkdb) GetUserByAPIKey(apikey string) (model.UserSchema, error) {
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

func (e *UsersRethinkdb) UpdateUserFullname(id, fullname string, updatedAt time.Time) (model.UserSchema, error) {
	return e.updateUser(id, map[string]interface{}{"fullname": fullname, "updated_at": updatedAt})
}

func (e *UsersRethinkdb) UpdateUserPassword(id, password string, updatedAt time.Time) (model.UserSchema, error) {
	return e.updateUser(id, map[string]interface{}{"password": password, "updated_at": updatedAt})
}

func (e *UsersRethinkdb) UpdateUserAPIKey(id, apikey string, updatedAt time.Time) (model.UserSchema, error) {
	return e.updateUser(id, map[string]interface{}{"apikey": apikey, "updated_at": updatedAt})
}

func (e *UsersRethinkdb) UpdateUserGlobusUser(id, globusUser string) error {
	errMsg := fmt.Sprintf("Unable to modify user %s to set globus user to %s", id, globusUser)
	resp, err := r.Table("users").Update(map[string]interface{}{"globus_user": globusUser}).RunWrite(e.Session)
	return checkRethinkdbWriteError(resp, err, errMsg)
}

func (e *UsersRethinkdb) updateUser(id string, what map[string]interface{}) (model.UserSchema, error) {
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
