package store

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
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
	resp, err := r.Table("users").Insert(user).RunWrite(e.Session)
	switch {
	case err != nil:
		return user, err
	case resp.Errors != 0:
		return user, fmt.Errorf("failed writing user %s", resp.FirstError)
	default:
		return user, nil
	}
}

func (e *UsersStoreEngineRethinkdb) GetUserByID(id string) (UserSchema, error) {
	var user UserSchema
	res, err := r.Table("users").Get(id).Run(e.Session)
	switch {
	case err != nil:
		return user, err
	case res.IsNil():
		return user, errors.Wrapf(ErrNotFound, "No such user %s", id)
	default:
		err = res.One(&user)
		return user, err
	}
}

func (e *UsersStoreEngineRethinkdb) GetUserByAPIKey(apikey string) (UserSchema, error) {
	var user UserSchema
	res, err := r.Table("users").GetAllByIndex("apikey", apikey).Run(e.Session)
	switch {
	case err != nil:
		return user, err
	case res.IsNil():
		return user, errors.Wrapf(ErrNotFound, "No such apikey %s", apikey)
	default:
		err = res.One(&user)
		return user, err
	}
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
		fmt.Println("len(resp.Changes) = ", len(resp.Changes))
		fmt.Printf("resp %#v\n\n", resp.Changes)
		err := encoding.Decode(&u, resp.Changes[0].NewValue)
		return u, err
	}
}

func (e *UsersStoreEngineRethinkdb) Name() string {
	return "UsersStoreEngineRethinkdb"
}
