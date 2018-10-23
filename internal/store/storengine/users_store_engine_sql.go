package storengine

import (
	"fmt"

	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/pkg/mc"

	"github.com/pkg/errors"

	"github.com/gocraft/dbr"
)

type UsersStoreEngineSQL struct {
	conn *dbr.Connection
}

func (e *UsersStoreEngineSQL) AddUser(user model.UserSchema) (model.UserSchema, error) {
	return model.UserSchema{}, nil
}

func (e *UsersStoreEngineSQL) GetUserByID(id string) (model.UserSchema, error) {
	var user model.UserSchema
	session := e.conn.NewSession(nil)
	err := session.Select("*").From("users").Where(dbr.Eq("id", id)).LoadOne(&user)
	return user, getDBError(err, fmt.Sprintf("No such user %s", id))
}

func (e *UsersStoreEngineSQL) GetUserByAPIKey(apikey string) (model.UserSchema, error) {
	var user model.UserSchema
	session := e.conn.NewSession(nil)
	err := session.Select("*").From("users").Where(dbr.Eq("apikey", apikey)).LoadOne(&user)
	return user, getDBError(err, fmt.Sprintf("No such apikey %s", apikey))
}

func getDBError(err error, msg string) error {
	switch {
	case err == nil:
		return nil
	case err == dbr.ErrNotFound:
		return errors.Wrap(mc.ErrNotFound, msg)
	default:
		return errors.Wrap(err, msg)
	}
}
