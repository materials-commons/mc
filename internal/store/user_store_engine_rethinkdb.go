package store

import (
	"github.com/pkg/errors"
	r "gopkg.in/gorethink/gorethink.v4"
)

type UserStoreEngineRethinkdb struct {
	s *r.Session
}

func (e *UserStoreEngineRethinkdb) GetUserByID(id string) (UserSchema, error) {
	var user UserSchema
	res, err := r.Table("users").Get(id).Run(e.s)
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

func (e *UserStoreEngineRethinkdb) GetUserByAPIKey(apikey string) (UserSchema, error) {
	var user UserSchema
	res, err := r.Table("users").GetAllByIndex("apikey", apikey).Run(e.s)
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
