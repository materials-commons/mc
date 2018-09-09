package store

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/gocraft/dbr"
)

type UserStoreEngineSQL struct {
	conn *dbr.Connection
}

func (e *UserStoreEngineSQL) AddUser(user UserSchema) (UserSchema, error) {
	return UserSchema{}, nil
}

func (e *UserStoreEngineSQL) GetUserByID(id string) (UserSchema, error) {
	var user UserSchema
	session := e.conn.NewSession(nil)
	err := session.Select("*").From("users").Where(dbr.Eq("id", id)).LoadOne(&user)
	return user, getDBError(err, fmt.Sprintf("No such user %s", id))
}

func (e *UserStoreEngineSQL) GetUserByAPIKey(apikey string) (UserSchema, error) {
	var user UserSchema
	session := e.conn.NewSession(nil)
	err := session.Select("*").From("users").Where(dbr.Eq("apikey", apikey)).LoadOne(&user)
	return user, getDBError(err, fmt.Sprintf("No such apikey %s", apikey))
}

func getDBError(err error, msg string) error {
	switch {
	case err == nil:
		return nil
	case err == dbr.ErrNotFound:
		return errors.Wrap(ErrNotFound, msg)
	default:
		return errors.Wrap(err, msg)
	}
}
