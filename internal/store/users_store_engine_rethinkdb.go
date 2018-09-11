package store

import (
	"time"

	"github.com/pkg/errors"
	r "gopkg.in/gorethink/gorethink.v4"
)

type UsersStoreEngineRethinkdb struct {
	session *r.Session
}

func (e *UsersStoreEngineRethinkdb) AddUser(user UserSchema) (UserSchema, error) {
	return UserSchema{}, nil
}

func (e *UsersStoreEngineRethinkdb) GetUserByID(id string) (UserSchema, error) {
	var user UserSchema
	res, err := r.Table("users").Get(id).Run(e.session)
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
	res, err := r.Table("users").GetAllByIndex("apikey", apikey).Run(e.session)
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
	return UserSchema{}, nil
}

func (e *UsersStoreEngineRethinkdb) ModifyUserPassword(id, password string, updatedAt time.Time) (UserSchema, error) {
	return UserSchema{}, nil
}

func (e *UsersStoreEngineRethinkdb) ModifyUserAPIKey(id, apikey string, updatedAt time.Time) (UserSchema, error) {
	return UserSchema{}, nil
}

func (e *UsersStoreEngineRethinkdb) Name() string {
	return "UsersStoreEngineRethinkdb"
}
